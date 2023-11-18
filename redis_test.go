package isunippets

import (
	"errors"
	"github.com/alicebob/miniredis/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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

	redisLogger = echo.New().Logger
	redisLogger.SetLevel(log.DEBUG)

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

func TestFlushRedisBatchRequests(t *testing.T) {
	assert := assert.New(t)

	s, err := SetUpRedisForTesting()
	assert.NoError(err)
	defer s.Close()

	err = PutRedisBatchRequest(RedisBatchRequest{}, redisBatchRequestNormal)
	assert.NoError(err)

	err = FlushRedisBatchRequests()
	assert.NoError(err)
}
