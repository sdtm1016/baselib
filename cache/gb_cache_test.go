package cache

import (
	"fmt"
	"time"
	"github.com/muesli/cache2go"
	."baselib/util/encrypt"
	"testing"
)

type myStruct struct {
	text     string
	moreData []byte
}

var(
	val1 = myStruct{"Hello world!", []byte{}}
	val2 = myStruct{"Hello guoBin!", []byte{}}
)

func TestCache(t *testing.T) {
	cacheT := GetCache()

	// callback
	cacheT.SetAddedItemCallback(func(entry *cache2go.CacheItem) {
		fmt.Println("Added:", entry.Key(),
			entry.Data().(*myStruct).text,
			TimeFormat(entry.CreatedOn()),
			TimeFormat(time.Now()))
	})

	// callback
	cacheT.SetAboutToDeleteItemCallback(func(entry *cache2go.CacheItem) {
		fmt.Println("Deleting:", entry.Key(),
			entry.Data().(*myStruct).text,
			TimeFormat(entry.CreatedOn()),
			TimeFormat(time.Now()))
	})

	// 添加item 将触发AddedItem callback.
	Put("someKey",&val1)

	// 取回item
	fmt.Println("Get value in cache:", Get("someKey"))

	// 删除item 将触发AboutToDeleteItem callback.
	cacheT.Delete("someKey")

	// 添加item 有过期时间3s
	res := PutWithExpire("anotherKey",&val2,3*time.Second)

	// item过期回调
	res.SetAboutToExpireCallback(func(key interface{}) {
		fmt.Println("About to expire:", key.(string),TimeFormat(time.Now()))
	})

	time.Sleep(5 * time.Second)
}
