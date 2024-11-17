package isunippets

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCleanupFileCache(t *testing.T) {
	assert := assert.New(t)

	err := PutFileCacheData([]byte("test"), "TestCleanupFileCache.txt")
	assert.NoError(err)

	_, err = GetFileCacheData("TestCleanupFileCache.txt")
	assert.NoError(err)

	err = CleanupFileCache()
	assert.NoError(err)

	_, err = GetFileCacheData("TestCleanupFileCache.txt")
	assert.Error(err)
}

func TestGetFileCacheData(t *testing.T) {
	assert := assert.New(t)

	_, err := GetFileCacheData("TestGetFileCacheData.txt")
	assert.Error(err)
}

func TestPutFileCacheStream(t *testing.T) {
	assert := assert.New(t)

	expected := new(bytes.Buffer)
	expected.WriteString("test")

	err := PutFileCacheStream(expected, "TestPutFileCacheStream.txt")
	assert.NoError(err)

	stream, err := GetFileCacheStream("TestPutFileCacheStream.txt")
	assert.NoError(err)

	actual := new(bytes.Buffer)
	_, err = actual.ReadFrom(stream)
	assert.Equal([]byte("test"), actual.Bytes())
}

func TestPutFileCacheData(t *testing.T) {
	assert := assert.New(t)

	expected := []byte("test")
	err := PutFileCacheData(expected, "TestPutFileCacheData.txt")
	assert.NoError(err)

	actual, err := GetFileCacheData("TestPutFileCacheData.txt")
	assert.NoError(err)
	assert.Equal(expected, actual)
}

func TestGetFileCacheStream(t *testing.T) {
	assert := assert.New(t)

	_, err := GetFileCacheStream("TestGetFileCacheStream.txt")
	assert.Error(err)
}
