package isunippets

import (
	"encoding/json"
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

func TestCompareJson(t *testing.T) {
	assert := assert.New(t)

	type MyStruct struct {
		String            string            `json:"string"`
		Int               int               `json:"int"`
		Bool              bool              `json:"bool"`
		Slice             []string          `json:"slice"`
		Map               map[string]string `json:"map"`
		MyStruct          *MyStruct         `json:"my_struct"`
		IgnoreString      string            `json:"-"`
		UnSupportedString string
	}

	myStruct := MyStruct{
		String: "string",
		Int:    1,
		Bool:   true,
		Slice:  []string{"slice1", "slice2"},
		Map: map[string]string{
			"key1": "value1",
		},
		MyStruct: &MyStruct{
			String: "stringInner",
			Int:    2,
			Bool:   false,
			Slice:  []string{"sliceInner1", "sliceInner2"},
			Map: map[string]string{
				"keyInner1": "valueInner1",
			},
			MyStruct: nil,
		},
		IgnoreString:      "ignore",
		UnSupportedString: "unSupported",
	}
	marshalledMyStruct, err := json.Marshal(myStruct)
	assert.NoError(err)

	differentMyStruct := MyStruct{
		String: "stringDiff",
		Int:    1,
		Bool:   true,
		Slice:  []string{"slice1", "slice2"},
		Map: map[string]string{
			"key1": "value1Diff",
		},
		MyStruct: &MyStruct{
			String: "stringInnerDiff",
			Int:    2,
			Bool:   false,
			Slice:  []string{"sliceInner1", "sliceInner2"},
			Map: map[string]string{
				"keyInner1": "valueInner1Diff",
			},
			MyStruct: nil,
		},
		IgnoreString:      "ignoreDiff",
		UnSupportedString: "unSupportedDiff",
	}
	marshalledDifferentMyStruct, err := json.Marshal(differentMyStruct)
	assert.NoError(err)

	cases := []struct {
		lhs      interface{}
		rhs      interface{}
		expected bool
	}{
		// 基本型
		{lhs: nil, rhs: nil, expected: true},
		{lhs: 1, rhs: 1, expected: true},
		{lhs: "1", rhs: "1", expected: true},
		{lhs: "\"\"", rhs: "\"\"", expected: true},
		{lhs: "", rhs: "", expected: false},
		{lhs: nil, rhs: "", expected: false},
		{lhs: "", rhs: nil, expected: false},

		{lhs: myStruct, rhs: myStruct, expected: true},
		{lhs: myStruct, rhs: string(marshalledMyStruct), expected: true},
		{lhs: myStruct, rhs: []byte(string(marshalledMyStruct)), expected: true},

		{lhs: myStruct, rhs: nil, expected: false},
		{lhs: myStruct, rhs: differentMyStruct, expected: false},
		{lhs: myStruct, rhs: string(marshalledDifferentMyStruct), expected: false},
		{lhs: myStruct, rhs: []byte(string(marshalledDifferentMyStruct)), expected: false},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := CompareJson(tt.lhs, tt.rhs)
			if tt.expected {
				assert.NoError(err)
			} else {
				assert.Error(err)
			}
		})
	}
}
