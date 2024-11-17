package isunippets

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	inMemoryCache          = cache.New(time.Hour, time.Hour)
	inMemoryCacheKeyPrefix = "data"
)

type InMemoryCacheData struct {
	String string
	Int    int
	Bool   bool
	Slice  []string
	Map    map[string]string
	Bytes  []byte
}

func PutInMemoryCacheData(value *InMemoryCacheData, userId int64) error {
	key := fmt.Sprintf("%s-%d", inMemoryCacheKeyPrefix, userId)
	inMemoryCache.Set(key, value, cache.DefaultExpiration)
	return nil
}

func GetInMemoryCacheData(userId int64) (*InMemoryCacheData, bool) {
	key := fmt.Sprintf("%s-%d", inMemoryCacheKeyPrefix, userId)
	rawValue, ok := inMemoryCache.Get(key)
	if !ok {
		return nil, false
	}
	return rawValue.(*InMemoryCacheData), true
}
