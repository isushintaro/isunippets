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
		r.result = r.params + i

		time.Sleep(1 * time.Second)

		return r, nil
	})
	assert.NoError(err)

	elapsed := time.Since(start)
	assert.LessOrEqual(elapsed.Seconds(), 2.5)

	sort.Slice(requests, func(i, j int) bool {
		return requests[i].params < requests[j].params
	})

	expected := []*MyConcurrencyRequest{
		{params: 0, result: 0},
		{params: 2, result: 3},
		{params: 4, result: 6},
		{params: 6, result: 9},
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

func TestRunConcurrent_WithError(t *testing.T) {
	assert := assert.New(t)

	type MyConcurrencyRequest struct {
		params int
		result int
	}

	requests := []*MyConcurrencyRequest{
		{params: 0},
		{params: 2},
		{params: 4},
		{params: 6},
	}

	opts := &RunConcurrentOptions{}

	start := time.Now()
	err := RunConcurrent(requests, opts, func(c context.Context, r *MyConcurrencyRequest, i int) (*MyConcurrencyRequest, error) {
		if r.params%2 == 0 {
			return r, errors.New(fmt.Sprintf("error: %d", i))
		}

		r.result = r.params + i

		time.Sleep(1 * time.Second)

		return r, nil
	})
	assert.Error(err)
	assert.Contains(err.Error(), "error: 0")
	assert.Contains(err.Error(), "error: 2")

	elapsed := time.Since(start)
	assert.LessOrEqual(elapsed.Seconds(), 0.5)

	sort.Slice(requests, func(i, j int) bool {
		return requests[i].params < requests[j].params
	})

	expected := []*MyConcurrencyRequest{
		{params: 0, result: 0},
		{params: 2, result: 0},
		{params: 4, result: 0},
		{params: 6, result: 0},
	}

	assert.Equal(expected, requests)
}
