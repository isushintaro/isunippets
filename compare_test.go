package isunippets

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompare(t *testing.T) {
	assert := assert.New(t)

	type MyStruct struct {
		String string
		Int    int
		Bool   bool
		Slice  []string
		Map    map[string]string
	}

	myStruct := MyStruct{
		String: "string",
		Int:    1,
		Bool:   true,
		Slice:  []string{"slice1", "slice2"},
		Map: map[string]string{
			"key1": "value1",
		},
	}
	myStructPointer := &myStruct

	anotherMyStruct := MyStruct{
		String: "string",
		Int:    1,
		Bool:   true,
		Slice:  []string{"slice1", "slice2"},
		Map: map[string]string{
			"key1": "value1",
		},
	}
	anotherMyStructPointer := &anotherMyStruct

	cases := []struct {
		lhs      interface{}
		rhs      interface{}
		expected bool
	}{
		// 基本型
		{lhs: nil, rhs: nil, expected: true},
		{lhs: "", rhs: "", expected: true},
		{lhs: nil, rhs: "", expected: false},
		{lhs: "", rhs: nil, expected: false},
		{lhs: 1, rhs: "a", expected: false},

		// 構造体
		{lhs: myStruct, rhs: myStruct, expected: true},
		{lhs: &myStruct, rhs: &myStruct, expected: true},
		{lhs: myStruct, rhs: anotherMyStruct, expected: true},
		{lhs: &myStruct, rhs: &anotherMyStruct, expected: true},
		{lhs: &myStruct, rhs: anotherMyStructPointer, expected: true},
		{lhs: myStructPointer, rhs: anotherMyStructPointer, expected: true},
		{lhs: myStruct, rhs: anotherMyStructPointer, expected: false},
		{lhs: myStruct, rhs: &myStruct, expected: false},
		{lhs: myStruct, rhs: nil, expected: false},
		{lhs: &myStruct, rhs: nil, expected: false},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Compare(tt.lhs, tt.rhs)
			if tt.expected {
				assert.NoError(err)
			} else {
				assert.Error(err)
			}
		})
	}
}
