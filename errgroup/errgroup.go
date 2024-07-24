package errgroup

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

type token struct{}

type PoolService interface {
	Submit(task func())
}

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid, has no limit on the number of active goroutines,
// and does not cancel on error.
type Group struct {
	pool PoolService

	ctx    context.Context
	cancel func(error)

	wg sync.WaitGroup

	sem chan token

	errOnce sync.Once
	err     error
}

func (g *Group) done() {
	if g.sem != nil {
		<-g.sem
	}
	g.wg.Done()
}

func (g *Group) setErrAndCancel(err error) {
	g.errOnce.Do(func() {
		g.err = err
		if g.cancel != nil {
			g.cancel(g.err)
		}
	})
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context, pool PoolService) *Group {
	ctx, cancel := withCancelCause(ctx)
	return &Group{pool: pool, ctx: ctx, cancel: cancel}
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel(g.err)
	}
	return g.err
}

// Go calls the given function in a new goroutine.
// It blocks until the new goroutine can be added without the number of
// active goroutines in the group exceeding the configured limit.
//
// The first call to return a non-nil error cancels the group's context, if the
// group was created by calling WithContext. The error will be returned by Wait.
func (g *Group) Go(f func(context.Context) error) {
	if g.sem != nil {
		g.sem <- token{}
	}

	g.wg.Add(1)
	f2 := func() {
		defer g.done()
		if err := g.ctx.Err(); err != nil { // cancellation from the parent context or the current context
			g.setErrAndCancel(errors.WithStack(err))
			return
		}
		if err := f(g.ctx); err != nil {
			g.setErrAndCancel(err)
		}
	}
	if g.pool != nil {
		g.pool.Submit(f2)
	} else {
		go f2()
	}
}

// SetLimit limits the number of active goroutines in this group to at most n.
// A negative value indicates no limit.
//
// Any subsequent call to the Go method will block until it can add an active
// goroutine without exceeding the configured limit.
//
// The limit must not be modified while any goroutines in the group are active.
func (g *Group) SetLimit(n int) {
	if n < 0 {
		g.sem = nil
		return
	}
	if len(g.sem) != 0 {
		panic(fmt.Errorf("errgroup: modify limit while %v goroutines in the group are still active", len(g.sem)))
	}
	g.sem = make(chan token, n)
}

func withCancelCause(parent context.Context) (context.Context, func(error)) {
	return context.WithCancelCause(parent)
}
