package isunippets

import (
	"fmt"
	"github.com/sourcegraph/conc/pool"
)

func RunConcurrent[T any](requests []*T, concurrency int, execute func(*T, int) *T) error {
	fmt.Printf("requests: %v\n", requests)
	p := pool.NewWithResults[*T]().WithMaxGoroutines(concurrency)
	for i, r := range requests {
		i := i
		r := r
		fmt.Printf("i: %d, %p\n", i, r)
		p.Go(func() *T {
			return execute(r, i)
		})
	}
	res := p.Wait()
	fmt.Printf("res: %v\n", res)

	return nil
}
