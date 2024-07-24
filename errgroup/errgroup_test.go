package errgroup

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	errSomeErr = fmt.Errorf("some error")
)

func Test_Group(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "a", 1)
	ctx, cancel := context.WithCancel(ctx)
	ctx = context.WithValue(ctx, "b", 2)
	ctx = context.WithValue(ctx, "c", 3)

	count := atomic.Int32{}

	g := WithContext(ctx, nil)
	g.Go(func(ctx context.Context) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		count.Add(1)
		return nil
	})
	g.Go(func(ctx context.Context) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		count.Add(1)
		return nil
	})
	assert.NoError(t, g.Wait())
	assert.NoError(t, ctx.Err())
	assert.Equal(t, count.Load(), int32(2))

	g = WithContext(ctx, nil)
	g.SetLimit(1)
	g.Go(func(_ context.Context) error {
		return errSomeErr
	})
	g.Go(func(_ context.Context) error { // will not be executed
		count.Add(1)
		return nil
	})
	assert.ErrorIs(t, g.Wait(), errSomeErr)
	assert.NoError(t, ctx.Err())
	assert.Equal(t, count.Load(), int32(2))

	cancel()

	g = WithContext(ctx, nil)
	g.Go(func(_ context.Context) error {
		count.Add(1)
		return nil
	})
	g.Go(func(_ context.Context) error {
		count.Add(1)
		return nil
	})
	assert.ErrorIs(t, g.Wait(), context.Canceled)
	assert.ErrorIs(t, ctx.Err(), context.Canceled)
	assert.Equal(t, count.Load(), int32(2))
}
