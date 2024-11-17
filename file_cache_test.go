package isunippets

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
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

	f, err := os.Open("gopher.png")
	assert.NoError(err)

	err = PutFileCacheStream(f, "TestPutFileCacheStream.png")
	assert.NoError(err)

	reader, err := GetFileCacheStream("TestPutFileCacheStream.png")
	assert.NoError(err)

	actual := new(bytes.Buffer)
	_, err = io.Copy(actual, reader)
	assert.NoError(err)

	expected, err := os.ReadFile("gopher.png")
	assert.NoError(err)
	assert.Equal(expected, actual.Bytes())
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
