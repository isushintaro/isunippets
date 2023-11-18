package isunippets

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	inMemoryCache    = cache.New(time.Hour, time.Hour)
	inMemoryCacheKey = "data"
)

type InMemoryCacheData struct {
	String string
	Int    int
	Bool   bool
	Slice  []string
	Map    map[string]string
}

func PutInMemoryCacheData(value *InMemoryCacheData) error {
	inMemoryCache.Set(inMemoryCacheKey, value, cache.DefaultExpiration)
	return nil
}

func GetInMemoryCacheData() (*InMemoryCacheData, bool) {
	rawValue, ok := inMemoryCache.Get(inMemoryCacheKey)
	if !ok {
		return nil, false
	}
	return rawValue.(*InMemoryCacheData), true
}
