package isunippets

import (
	"encoding/json"
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

func CompareJson(lhs, rhs interface{}) error {
	if lhs == nil && rhs == nil {
		return nil
	} else if lhs == nil {
		return errors.New(fmt.Sprintf("lhs is nil: rhs=%p", rhs))
	} else if rhs == nil {
		return errors.New(fmt.Sprintf("rhs is nil: lhs=%p", lhs))
	}

	var lhsBytes, rhsBytes []byte

	switch v := lhs.(type) {
	case []byte:
		lhsBytes = v
		if !json.Valid(lhsBytes) {
			return errors.New(fmt.Sprintf("invalid lhs bytes: %v", lhsBytes))
		}
	case string:
		lhsBytes = []byte(v)
		if !json.Valid(lhsBytes) {
			return errors.New(fmt.Sprintf("invalid lhs string: %v", lhsBytes))
		}
	default:
		bs, err := json.Marshal(v)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to marshal lhs: %v", err))
		}
		lhsBytes = bs
	}

	switch v := rhs.(type) {
	case []byte:
		rhsBytes = v
		if !json.Valid(rhsBytes) {
			return errors.New(fmt.Sprintf("invalid rhs bytes: %v", rhsBytes))
		}
	case string:
		rhsBytes = []byte(v)
		if !json.Valid(rhsBytes) {
			return errors.New(fmt.Sprintf("invalid rhs string: %v", rhsBytes))
		}
	default:
		bs, err := json.Marshal(v)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to marshal rhs: %v", err))
		}
		rhsBytes = bs
	}

	transformJSON := cmp.FilterValues(func(x, y []byte) bool {
		return json.Valid(x) && json.Valid(y)
	}, cmp.Transformer("ParseJSON", func(in []byte) (out interface{}) {
		if err := json.Unmarshal(in, &out); err != nil {
			panic(err)
		}
		return out
	}))

	result := cmp.Diff(lhsBytes, rhsBytes, transformJSON)
	if result != "" {
		return errors.New(fmt.Sprintf("has diff json: lhs=%p, rhs=%p, diff=%s", lhs, rhs, result))
	} else {
		return nil
	}
}
