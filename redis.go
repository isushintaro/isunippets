package isunippets

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

var (
	redisLogger = echo.New().Logger
	redisCtx    = context.Background()
)

func SetRedisLogLevel(level log.Lvl) {
	redisLogger.SetLevel(level)
}

func GetRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	return rdb
}

type RedisCacheValue struct {
	String string            `json:"string"`
	Int    int               `json:"int"`
	Bool   bool              `json:"bool"`
	Slice  []string          `json:"slice"`
	Map    map[string]string `json:"map"`
}

type RedisBatchRequest struct {
	QueueName string `json:"queueName"`
	QueuedAt  string `json:"queuedAt"`

	String string            `json:"string"`
	Int    int               `json:"int"`
	Bool   bool              `json:"bool"`
	Slice  []string          `json:"slice"`
	Map    map[string]string `json:"map"`
}

type RedisBatchRequestContext struct {
	Logger echo.Logger
}

const (
	redisBatchRequestHigh   = "redisBatchRequest:high"
	redisBatchRequestNormal = "redisBatchRequest:normal"
)

func ClearRedis() error {
	rdb := GetRedisClient()

	statusCmd := rdb.FlushAll(redisCtx)
	if statusCmd.Err() != nil {
		redisLogger.Errorf("failed to FlushAll: %v", statusCmd.Err())
		return statusCmd.Err()
	}

	return nil
}

func RunRedisBatchMainLoop(timeout time.Duration, continueOnError bool, process func(RedisBatchRequest, RedisBatchRequestContext) error) error {
	redisLogger.Info("start runRedisBatchMainLoop")
	rdb := GetRedisClient()
	requestContext := RedisBatchRequestContext{
		Logger: redisLogger,
	}
	for {
		highLen := rdb.LLen(redisCtx, redisBatchRequestHigh).Val()
		normalLen := rdb.LLen(redisCtx, redisBatchRequestNormal).Val()
		redisLogger.Debugf("queue length high=%v normal=%v", highLen, normalLen)

		blPopCmd := rdb.BLPop(redisCtx, timeout, redisBatchRequestHigh, redisBatchRequestNormal)
		if blPopCmd.Err() != nil {
			if continueOnError {
				redisLogger.Errorf("failed to BLPop: %v", blPopCmd.Err())
				continue
			} else {
				return blPopCmd.Err()
			}
		}
		redisLogger.Debug("fetch request: %v", blPopCmd.Val())
		requestStr := blPopCmd.Val()[1]

		var request RedisBatchRequest
		if err := json.Unmarshal([]byte(requestStr), &request); err != nil {
			if continueOnError {
				redisLogger.Errorf("failed to unmarshal: %v, %v", requestStr, err)
				continue
			} else {
				return err
			}
		}
		err := process(request, requestContext)
		if err != nil {
			if continueOnError {
				redisLogger.Errorf("failed to processRedisBatchRequest: %v", err)
				continue
			} else {
				return err
			}
		}
	}
}

func PutRedisBatchRequest(request RedisBatchRequest, queueName string) error {
	client := GetRedisClient()
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			redisLogger.Errorf("failed to close: %v", err)
		}
	}(client)

	request.QueueName = queueName
	request.QueuedAt = time.Now().String()

	requestBytes, err := json.Marshal(request)
	if err != nil {
		redisLogger.Errorf("failed to marshal: %v, %v", requestBytes, err)
		return err
	}

	requestStr := string(requestBytes)
	redisLogger.Debugf("try RPush: %v, %v", queueName, requestStr)

	rPushCmd := client.RPush(redisCtx, queueName, requestStr)
	if rPushCmd.Err() != nil {
		redisLogger.Errorf("failed to RPush: %v", rPushCmd.Err())
		return rPushCmd.Err()
	}

	redisLogger.Debug("RPush complete")
	return nil
}

func PutRedisCache(key string, value RedisCacheValue) error {
	dataStr, err := json.Marshal(value)
	if err != nil {
		return err
	}

	rdb := GetRedisClient()
	err = rdb.Set(redisCtx, key, dataStr, time.Duration(0)).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetRedisCache(key string) (RedisCacheValue, error) {
	rdb := GetRedisClient()
	getCmd := rdb.Get(redisCtx, key)
	if getCmd.Err() != nil {
		return RedisCacheValue{}, getCmd.Err()
	}

	valBytes := getCmd.Val()

	var value RedisCacheValue
	err := json.Unmarshal([]byte(valBytes), &value)
	if err != nil {
		return RedisCacheValue{}, err
	}

	return value, nil
}
