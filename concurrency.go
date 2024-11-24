package isunippets

import (
	"context"
	"github.com/sourcegraph/conc/pool"
)

type RunConcurrentOptions struct {
	Concurrency int
}

func RunConcurrent[T any](requests []*T, options *RunConcurrentOptions, execute func(context.Context, *T, int) (*T, error)) error {
	ctx := context.Background()
	p := pool.NewWithResults[*T]()
	if options != nil {
		if options.Concurrency > 0 {
			p = p.WithMaxGoroutines(options.Concurrency)
		}
	}

	wp := p.WithContext(ctx).WithCancelOnError()

	for i, r := range requests {
		i := i
		r := r
		wp.Go(func(c context.Context) (*T, error) {
			return execute(c, r, i)
		})
	}
	_, err := wp.Wait()
	return err
}
