package isunippets

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInMemoryCacheData(t *testing.T) {
	assert := assert.New(t)

	inMemoryCache.Flush()

	_, ok := GetInMemoryCacheData()
	assert.False(ok)
}

func TestPutInMemoryCacheData(t *testing.T) {
	assert := assert.New(t)

	inMemoryCache.Flush()

	sliceValue := []string{"a", "b", "c"}
	mapValue := map[string]string{"a": "b", "c": "d"}

	expected := InMemoryCacheData{
		String: "value1",
		Int:    1,
		Bool:   true,
		Slice:  sliceValue,
		Map:    mapValue,
	}

	err := PutInMemoryCacheData(&expected)
	assert.NoError(err)

	// 値はコピーされない
	expected.String = "value2"
	expected.Int = 2
	expected.Bool = false
	sliceValue[0] = "d"
	mapValue["a"] = "e"

	raw, ok := inMemoryCache.Get(inMemoryCacheKey)
	assert.True(ok)

	actual, ok := raw.(*InMemoryCacheData)
	assert.True(ok)
	assert.Equal(&expected, actual)
}
