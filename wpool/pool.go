package wpool

import (
	"sync/atomic"
	"time"

	"github.com/zzzzer91/gopkg/logx"
	"github.com/zzzzer91/gopkg/stackx"
)

// Task is the function that the worker will execute.
type Task func()

// CachedGoroutinePool is a worker pool bind with some idle goroutines.
// Its behavior is like CachedThreadPool in Java.
type CachedGoroutinePool struct {
	size  int32
	tasks chan Task

	// maxIdle is the number of the max idle workers in the pool.
	// if maxIdle too small, the pool works like a native 'go func()'.
	maxIdle int32
	// maxIdleTime is the max idle time that the worker will wait for the new task.
	maxIdleTime time.Duration
}

func NewCachedGoroutinePool(maxIdle int, maxIdleTime time.Duration) *CachedGoroutinePool {
	return &CachedGoroutinePool{
		tasks:       make(chan Task),
		maxIdle:     int32(maxIdle),
		maxIdleTime: maxIdleTime,
	}
}

// Size returns the number of the running workers.
func (p *CachedGoroutinePool) Size() int32 {
	return atomic.LoadInt32(&p.size)
}

// Submit creates/reuses a worker to run task.
func (p *CachedGoroutinePool) Submit(task Task) {
	select {
	case p.tasks <- task:
		// reuse exist worker
		return
	default:
	}

	// create new worker
	atomic.AddInt32(&p.size, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logx.Errorf("panic in CachedGoroutinePool: error=%v: stack=%s", r, stackx.RecordStack(2))
			}
			atomic.AddInt32(&p.size, -1)
		}()
		task()
		if atomic.LoadInt32(&p.size) > p.maxIdle {
			return
		}

		// waiting for new task
		idleTimer := time.NewTimer(p.maxIdleTime)
		for {
			select {
			case task = <-p.tasks:
				task()
			case <-idleTimer.C:
				// worker exits
				return
			}

			if !idleTimer.Stop() {
				<-idleTimer.C
			}
			idleTimer.Reset(p.maxIdleTime)
		}
	}()
}

// NewFixedGoroutinePool is a worker pool bind with some idle goroutines.
// Its behavior is like FixedThreadPool in Java, but it supports setting idle time.
// And the task queue has a capacity limit.
type FixedGoroutinePool struct {
	size  int32
	tasks chan Task

	maxSize     int32
	maxIdleTime time.Duration
}

func NewFixedGoroutinePool(maxSize int, maxIdleTime time.Duration, maxTaskQueueSize int) *FixedGoroutinePool {
	return &FixedGoroutinePool{
		tasks:       make(chan Task, maxTaskQueueSize),
		maxSize:     int32(maxSize),
		maxIdleTime: maxIdleTime,
	}
}

// Size returns the number of the running workers.
func (p *FixedGoroutinePool) Size() int32 {
	return atomic.LoadInt32(&p.size)
}

// Submit creates/reuses a worker to run task.
// If the task queue is full, it will be blocked.
func (p *FixedGoroutinePool) Submit(task Task) {
	if atomic.LoadInt32(&p.size) < p.maxSize {
		// create new worker
		atomic.AddInt32(&p.size, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logx.Errorf("panic in FixedGoroutinePool: error=%v: stack=%s", r, stackx.RecordStack(2))
				}
				atomic.AddInt32(&p.size, -1)
			}()
			// waiting for new task
			idleTimer := time.NewTimer(p.maxIdleTime)
			for {
				select {
				case task = <-p.tasks:
					task()
				case <-idleTimer.C:
					// worker exits
					return
				}

				if !idleTimer.Stop() {
					<-idleTimer.C
				}
				idleTimer.Reset(p.maxIdleTime)
			}
		}()
	}

	// WARN: If the task queue is full, it will be blocked here.
	p.tasks <- task
}
