/*
 * Copyright 2021 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package wpool

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCachedGoroutinePool(t *testing.T) {
	p := NewCachedGoroutinePool(1, time.Millisecond*100)
	var (
		sum  int32
		wg   sync.WaitGroup
		size = 100
	)
	assert.True(t, p.Size() == 0)
	for i := 0; i < size; i++ {
		wg.Add(1)
		p.Submit(func() {
			defer wg.Done()
			panic("test")
		})
	}
	assert.True(t, p.Size() != 0)

	wg.Wait()
	assert.Equal(t, atomic.LoadInt32(&sum), int32(size))
	time.Sleep(time.Millisecond * 10) // waiting for workers finished
	assert.True(t, p.Size() == 1)
	time.Sleep(time.Millisecond * 100) // waiting for idle timeout
	assert.True(t, p.Size() == 0)
}

func TestFixedGoroutinePool(t *testing.T) {
	p := NewFixedGoroutinePool(4, time.Millisecond*100, 1000)
	var (
		sum  int32
		wg   sync.WaitGroup
		size = 100
	)
	assert.True(t, p.Size() == 0)
	for i := 0; i < size; i++ {
		wg.Add(1)
		p.Submit(func() {
			defer wg.Done()
			atomic.AddInt32(&sum, 1)
		})
	}
	assert.True(t, p.Size() == 4)

	wg.Wait()
	assert.Equal(t, atomic.LoadInt32(&sum), int32(size))
	time.Sleep(time.Millisecond * 10) // waiting for workers finished
	assert.True(t, p.Size() == 4)
	time.Sleep(time.Millisecond * 100) // waiting for idle timeout
	assert.True(t, p.Size() == 0)
}
