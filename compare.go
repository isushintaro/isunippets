package isunippets

import (
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
)

func Compare(lhs, rhs interface{}) error {
	opts := cmp.AllowUnexported(lhs, rhs)
	result := cmp.Diff(lhs, rhs, opts)
	if result != "" {
		return errors.New(fmt.Sprintf("has diff: lhs=%p, rhs=%p, diff=%s", lhs, rhs, result))
	} else {
		return nil
	}
}
