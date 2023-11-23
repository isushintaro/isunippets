package isunippets

import (
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
)

func Compare(lhs, rhs interface{}) error {
	result := cmp.Diff(lhs, rhs)
	if result != "" {
		return errors.New(fmt.Sprintf("has diff: lhs=%p, rhs=%p, diff=%s", lhs, rhs, result))
	} else {
		return nil
	}
}
