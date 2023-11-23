package isunippets

import (
	"github.com/stretchr/testify/assert"
	"sort"
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

	start := time.Now()
	err := RunConcurrent(requests, 2, func(r *MyConcurrencyRequest, i int) *MyConcurrencyRequest {
		time.Sleep(1 * time.Second)

		r.index = i
		r.result = r.params + i
		return r
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
