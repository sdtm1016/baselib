package cache

import (
	"github.com/muesli/cache2go"
	"time"
)

const (
	//a new cache table
	GB_CACHE_NAME = "gb_cache"
)

var gbCache *cache2go.CacheTable = cache2go.Cache(GB_CACHE_NAME)

func GetCache() *cache2go.CacheTable{
	return gbCache
}

func Get(name interface{}) interface{} {
	item, _ := gbCache.Value(name)
	if item != nil {
		return item.Data()
	}

	return nil
}

func Put(key interface{}, val interface{}) {
	gbCache.Add(key, 0, val)
}

func PutWithExpire(key interface{}, val interface{}, expire time.Duration) *cache2go.CacheItem{
	return gbCache.Add(key, expire, val)
}