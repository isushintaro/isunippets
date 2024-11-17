package isunippets

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInMemoryCacheData(t *testing.T) {
	assert := assert.New(t)

	inMemoryCache.Flush()

	userId := int64(1)

	_, ok := GetInMemoryCacheData(userId)
	assert.False(ok)
}

func TestPutInMemoryCacheData(t *testing.T) {
	assert := assert.New(t)

	inMemoryCache.Flush()

	userId := int64(1)
	sliceValue := []string{"a", "b", "c"}
	mapValue := map[string]string{"a": "b", "c": "d"}
	bytesValue := []byte{0x01, 0x02, 0x03}

	expected := InMemoryCacheData{
		String: "value1",
		Int:    1,
		Bool:   true,
		Slice:  sliceValue,
		Map:    mapValue,
		Bytes:  bytesValue,
	}

	err := PutInMemoryCacheData(&expected, userId)
	assert.NoError(err)

	expected.String = "value2"
	expected.Int = 2
	expected.Bool = false
	sliceValue[0] = "d"
	mapValue["a"] = "e"
	bytesValue[0] = 0x04

	raw, ok := inMemoryCache.Get(fmt.Sprintf("%s-%d", inMemoryCacheKeyPrefix, userId))
	assert.True(ok)

	actual, ok := raw.(*InMemoryCacheData)
	assert.True(ok)
	assert.Equal(&expected, actual)
}
