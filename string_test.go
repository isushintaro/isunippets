package isunippets

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegexpIsMatch(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		condition string
		expected  bool
	}{
		{
			condition: "is_dirty=true",
			expected:  true,
		},
		{
			condition: "is_dirty=false",
			expected:  true,
		},
		{
			condition: "is_dirty=unknown",
			expected:  false,
		},
		{
			condition: "",
			expected:  false,
		},
	}

	for i, tt := range cases {
		tt := tt
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()
			actual := RegexpIsMatch(tt.condition)
			assert.Equal(tt.expected, actual)
		})
	}
}

func TestGenerateUUID(t *testing.T) {
	assert := assert.New(t)

	actual := GenerateUUID()
	assert.NotEmpty(actual)
}
