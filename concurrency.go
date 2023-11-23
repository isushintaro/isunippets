package isunippets

import (
	"context"
	"fmt"
	"github.com/sourcegraph/conc/pool"
)

type RunConcurrentOptions struct {
	Concurrency   int
	CancelOnError bool
}

func RunConcurrent[T any](requests []*T, options *RunConcurrentOptions, execute func(context.Context, *T, int) (*T, error)) error {
	fmt.Printf("requests: %v\n", requests)
	ctx := context.Background()
	p := pool.NewWithResults[*T]()
	if options != nil {
		if options.Concurrency > 0 {
			p = p.WithMaxGoroutines(options.Concurrency)
		}
	}

	wp := p.WithContext(ctx).WithCollectErrored()
	if options != nil {
		if options.CancelOnError {
			wp = wp.WithCancelOnError()
		}
	}

	for i, r := range requests {
		i := i
		r := r
		fmt.Printf("i: %d, %p\n", i, r)
		wp.Go(func(c context.Context) (*T, error) {
			return execute(c, r, i)
		})
	}
	res, err := wp.Wait()
	fmt.Printf("res: %v\n", res)
	return err
}
