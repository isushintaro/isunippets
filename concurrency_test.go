package isunippets

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestRunConcurrent(t *testing.T) {
	assert := assert.New(t)

	type MyConcurrencyRequest struct {
		index  int
		params int
		result int
	}

	requests := []*MyConcurrencyRequest{
		{params: 0},
		{params: 2},
		{params: 4},
		{params: 6},
	}

	opts := &RunConcurrentOptions{
		Concurrency: 2,
	}

	start := time.Now()
	err := RunConcurrent(requests, opts, func(c context.Context, r *MyConcurrencyRequest, i int) (*MyConcurrencyRequest, error) {
		r.index = i
		r.result = r.params + i

		time.Sleep(1 * time.Second)

		return r, nil
	})
	assert.NoError(err)

	elapsed := time.Since(start)
	assert.LessOrEqual(elapsed.Seconds(), 2.5)

	sort.Slice(requests, func(i, j int) bool {
		return requests[i].index < requests[j].index
	})

	expected := []*MyConcurrencyRequest{
		{index: 0, params: 0, result: 0},
		{index: 1, params: 2, result: 3},
		{index: 2, params: 4, result: 6},
		{index: 3, params: 6, result: 9},
	}

	assert.Equal(expected, requests)
}

func TestRunConcurrent_WithMutex(t *testing.T) {
	assert := assert.New(t)

	var requests []*int
	for i := 0; i < 100000; i++ {
		requests = append(requests, &i)
	}

	var mu sync.Mutex
	cnt := 0

	start := time.Now()
	err := RunConcurrent(requests, nil, func(c context.Context, r *int, i int) (*int, error) {
		time.Sleep(50 * time.Millisecond)
		mu.Lock()
		cnt++
		mu.Unlock()
		time.Sleep(50 * time.Millisecond)
		return nil, nil
	})
	assert.NoError(err)

	elapsed := time.Since(start)
	assert.LessOrEqual(elapsed.Seconds(), 0.5)
	assert.Equal(100000, cnt)
}

func TestRunConcurrent_With(t *testing.T) {
	// FIXME: 期待通りの動作（エラー時に即時実行終了する）をしていない
	assert := assert.New(t)

	type MyConcurrencyRequest struct {
		index  int
		params int
		result int
	}

	requests := []*MyConcurrencyRequest{
		{params: 0},
		{params: 2},
		{params: 4},
		{params: 6},
	}

	opts := &RunConcurrentOptions{
		CancelOnError: true,
	}

	start := time.Now()
	err := RunConcurrent(requests, opts, func(c context.Context, r *MyConcurrencyRequest, i int) (*MyConcurrencyRequest, error) {
		r.index = i
		r.result = r.params + i

		time.Sleep(1 * time.Second)

		if i%2 == 0 {
			return r, errors.New(fmt.Sprintf("error: %d", i))
		}

		return r, nil
	})
	assert.Error(err)
	assert.Contains(err.Error(), "error: 0")
	assert.Contains(err.Error(), "error: 2")

	elapsed := time.Since(start)
	assert.LessOrEqual(elapsed.Seconds(), 1.5)

	sort.Slice(requests, func(i, j int) bool {
		return requests[i].index < requests[j].index
	})

	expected := []*MyConcurrencyRequest{
		{index: 0, params: 0, result: 0},
		{index: 1, params: 2, result: 3},
		{index: 2, params: 4, result: 6},
		{index: 3, params: 6, result: 9},
	}

	assert.Equal(expected, requests)
}
