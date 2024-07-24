package wpool

import (
	stderrors "errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/zzzzer91/gopkg/logx"
	"github.com/zzzzer91/gopkg/stackx"
)

var (
	ErrWpoolIsFull = stderrors.New("wpool is full")
)

// CachedGoroutinePool is a worker pool bind with some idle goroutines.
// Its behavior is like CachedThreadPool in Java.
type CachedGoroutinePool struct {
	wg sync.WaitGroup

	size  int32
	tasks chan func()

	// maxIdle is the number of the max idle workers in the pool.
	// if maxIdle too small, the pool works like a native 'go func()'.
	maxIdle int32
	// maxIdleTime is the max idle time that the worker will wait for the new task.
	maxIdleTime time.Duration
}

func NewCachedGoroutinePool(maxIdle int, maxIdleTime time.Duration) *CachedGoroutinePool {
	return &CachedGoroutinePool{
		tasks:       make(chan func()),
		maxIdle:     int32(maxIdle),
		maxIdleTime: maxIdleTime,
	}
}

// Size returns the number of the running workers.
func (p *CachedGoroutinePool) Size() int32 {
	return atomic.LoadInt32(&p.size)
}

// Submit creates/reuses a worker to run task.
func (p *CachedGoroutinePool) Submit(task func()) {
	select {
	case p.tasks <- task:
		// reuse exist worker
		return
	default:
	}

	// create new worker
	p.wg.Add(1)
	atomic.AddInt32(&p.size, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logx.Errorf("panic in CachedGoroutinePool: error=%v: stack=%s", r, stackx.StackToString(stackx.Callers(2)))
			}
			atomic.AddInt32(&p.size, -1)
			p.wg.Done()
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

func (p *CachedGoroutinePool) Wait() {
	p.wg.Wait()
}

// NewFixedGoroutinePool is a worker pool bind with some idle goroutines.
// Its behavior is like FixedThreadPool in Java, but it supports setting idle time.
// And the task queue has a capacity limit.
type FixedGoroutinePool struct {
	sync.Mutex
	wg sync.WaitGroup

	size  int32
	tasks chan func()

	maxSize     int32
	maxIdleTime time.Duration
	allowBlock  bool
}

func NewFixedGoroutinePool(maxSize int, maxIdleTime time.Duration, maxTaskQueueSize int, allowBlock bool) *FixedGoroutinePool {
	return &FixedGoroutinePool{
		tasks:       make(chan func(), maxTaskQueueSize),
		maxSize:     int32(maxSize),
		maxIdleTime: maxIdleTime,
		allowBlock:  allowBlock,
	}
}

// Size returns the number of the running workers.
func (p *FixedGoroutinePool) Size() int32 {
	return atomic.LoadInt32(&p.size)
}

// Submit creates/reuses a worker to run task.
// If the task queue is full, returns error.
func (p *FixedGoroutinePool) Submit(task func()) error {
	if atomic.LoadInt32(&p.size) < p.maxSize {
		p.createWorker()
	}
	if p.allowBlock {
		p.tasks <- task
	} else {
		select {
		case p.tasks <- task:
		default:
			return errors.WithStack(ErrWpoolIsFull)
		}
	}
	return nil
}

func (p *FixedGoroutinePool) createWorker() {
	p.Lock()
	defer p.Unlock()
	if p.size < p.maxSize {
		// create new worker
		p.wg.Add(1)
		atomic.AddInt32(&p.size, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logx.Errorf("panic in FixedGoroutinePool: error=%v: stack=%s", r, stackx.StackToString(stackx.Callers(2)))
				}
				atomic.AddInt32(&p.size, -1)
				p.wg.Done()
			}()
			// waiting for new task
			idleTimer := time.NewTimer(p.maxIdleTime)
			for {
				select {
				case task := <-p.tasks:
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
}

func (p *FixedGoroutinePool) Wait() {
	p.wg.Wait()
}
