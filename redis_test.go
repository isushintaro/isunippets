package isunippets

import (
	"errors"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func SetUpRedisForTesting() (*miniredis.Miniredis, error) {
	s, err := miniredis.Run()
	if err != nil {
		return nil, err
	}

	err = os.Setenv("REDIS_HOST", s.Host())
	if err != nil {
		return nil, err
	}
	err = os.Setenv("REDIS_PORT", s.Port())
	if err != nil {
		return nil, err
	}
	err = os.Setenv("REDIS_PASSWORD", "")
	if err != nil {
		return nil, err
	}

	//SetRedisLogLevel(log.DEBUG)

	return s, nil
}

func TestGetRedisClient(t *testing.T) {
	assert := assert.New(t)

	s, err := SetUpRedisForTesting()
	assert.NoError(err)
	defer s.Close()

	rdb := GetRedisClient()
	assert.NotNil(rdb)
}

func TestPutRedisBatchRequest(t *testing.T) {
	assert := assert.New(t)

	s, err := SetUpRedisForTesting()
	assert.NoError(err)
	defer s.Close()

	err = PutRedisBatchRequest(RedisBatchRequest{}, redisBatchRequestNormal)
	assert.NoError(err)
}

func TestRunRedisBatchMainLoop(t *testing.T) {
	assert := assert.New(t)

	s, err := SetUpRedisForTesting()
	assert.NoError(err)
	defer s.Close()

	for i, queueName := range []string{redisBatchRequestNormal, redisBatchRequestHigh} {
		err = PutRedisBatchRequest(RedisBatchRequest{
			String: "string",
			Int:    i,
			Bool:   true,
			Slice:  []string{"slice1", "slice2"},
			Map: map[string]string{
				"key1": "value1",
			},
		}, queueName)
		assert.NoError(err)
	}

	finished := errors.New("process finished")
	process := func(request RedisBatchRequest, requestContext RedisBatchRequestContext) error {
		assert.NotEmpty(request.QueuedAt)
		assert.Equal("slice1", request.Slice[0])
		assert.Equal("value1", request.Map["key1"])
		if request.QueueName == redisBatchRequestHigh {
			assert.Equal(1, request.Int)
			return nil
		} else if request.QueueName == redisBatchRequestNormal {
			assert.Equal(0, request.Int)
			return finished
		} else {
			return errors.New("invalid queue name")
		}
	}

	err = RunRedisBatchMainLoop(1*time.Second, false, process)
	assert.ErrorIs(err, finished)
}

func TestClearRedis(t *testing.T) {
	assert := assert.New(t)

	s, err := SetUpRedisForTesting()
	assert.NoError(err)
	defer s.Close()

	err = PutRedisBatchRequest(RedisBatchRequest{}, redisBatchRequestNormal)
	assert.NoError(err)

	err = ClearRedis()
	assert.NoError(err)
}

func TestPutRedisCache(t *testing.T) {
	assert := assert.New(t)

	s, err := SetUpRedisForTesting()
	assert.NoError(err)
	defer s.Close()

	err = PutRedisCache("key", RedisCacheValue{
		String: "string",
		Int:    1,
		Bool:   true,
		Slice:  []string{"slice1", "slice2"},
		Map: map[string]string{
			"key1": "value1",
		},
	})
	assert.NoError(err)
}

func TestGetRedisCache(t *testing.T) {
	assert := assert.New(t)

	s, err := SetUpRedisForTesting()
	assert.NoError(err)
	defer s.Close()

	for i, key := range []string{"key1", "key2"} {
		expected := RedisCacheValue{
			String: fmt.Sprintf("string %s", key),
			Int:    i,
			Bool:   true,
			Slice:  []string{"slice1", "slice2"},
			Map: map[string]string{
				"key1": "value1",
			},
		}

		err = PutRedisCache(key, expected)
		assert.NoError(err)
	}

	value, err := GetRedisCache("key1")
	assert.NoError(err)
	assert.Equal("string key1", value.String)

	value, err = GetRedisCache("key2")
	assert.NoError(err)
	assert.Equal("string key2", value.String)
}
