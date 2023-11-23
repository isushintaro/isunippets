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
	marshalledMyStructString := string(marshalledMyStruct)

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
	marshalledDifferentMyStructString := string(marshalledDifferentMyStruct)

	cases := []struct {
		lhs      interface{}
		rhs      interface{}
		expected bool
	}{
		// 基本データ型
		{lhs: nil, rhs: nil, expected: true},
		{lhs: "1", rhs: "1", expected: true},
		{lhs: "\"\"", rhs: "\"\"", expected: true},
		{lhs: "\"1\"", rhs: "\"1\"", expected: true},
		{lhs: "false", rhs: "false", expected: true},
		{lhs: "null", rhs: "null", expected: true},
		{lhs: "1", rhs: "1.0", expected: true},
		{lhs: "{}", rhs: "{}", expected: true},

		// JSONリテラルになったときに妥当な文字列
		{lhs: 1, rhs: 1, expected: true},
		{lhs: -1.234, rhs: -1.234, expected: true},
		{lhs: 2, rhs: "2", expected: true},
		{lhs: nil, rhs: "null", expected: true},
		{lhs: true, rhs: true, expected: true},
		{lhs: false, rhs: "false", expected: true},

		// データ型は考慮される
		{lhs: "\"1\"", rhs: "1", expected: false},

		// 値が違えばエラー
		{lhs: "1.2345", rhs: "1.2346", expected: false},
		{lhs: "\"hello\"", rhs: "\"Hello\"", expected: false},

		// 空値はエラー
		{lhs: nil, rhs: "", expected: false},
		{lhs: "", rhs: nil, expected: false},
		{lhs: "", rhs: "", expected: false},

		// JSON構造体
		{lhs: myStruct, rhs: myStruct, expected: true},
		{lhs: myStruct, rhs: &myStruct, expected: true},
		{lhs: myStruct, rhs: marshalledMyStruct, expected: true},
		{lhs: myStruct, rhs: marshalledMyStructString, expected: true},

		// ポインタ変換された場合は不一致
		{lhs: myStruct, rhs: &marshalledMyStruct, expected: false},
		{lhs: myStruct, rhs: &marshalledMyStructString, expected: false},

		// JSON構造体の中身が違えばエラー
		{lhs: myStruct, rhs: nil, expected: false},
		{lhs: myStruct, rhs: differentMyStruct, expected: false},
		{lhs: myStruct, rhs: &differentMyStruct, expected: false},
		{lhs: myStruct, rhs: marshalledDifferentMyStruct, expected: false},
		{lhs: myStruct, rhs: &marshalledDifferentMyStruct, expected: false},
		{lhs: myStruct, rhs: marshalledDifferentMyStructString, expected: false},
		{lhs: myStruct, rhs: &marshalledDifferentMyStructString, expected: false},

		// 書式が違うだけならOK
		{
			lhs:      "{\"str\": \"慎太郎\", \"int\": 1}",
			rhs:      "{\n\"int\":1, \"str\":\t\t\"\\u614e\\u592a\\u90ce\"    }",
			expected: true,
		},

		// 値の順番が違うとエラー
		{
			lhs:      "{\"ary\": [1, 2]}",
			rhs:      "{\"ary\": [2, 1]}",
			expected: false,
		},

		// 形式不正
		{lhs: "{", rhs: "{", expected: false},
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
