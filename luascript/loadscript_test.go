package luascript

import (
	"testing"
	"github.com/go-redis/redis"
	"fmt"
)

var (
	client *redis.Client
	ls     *ScriptLoader
)

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "gb",
		DB:       1,
	})
	script := []string{
		GET_BATCH_VALUES_LUA,
		CHECK_GROUP_ALL_STATUS_LUA,
	}
	ls = NewLoader(client, script)
}

func TestNewLoader(t *testing.T) {
	fmt.Printf("%s\n", ls.GetSha(GET_BATCH_VALUES_LUA))

	fmt.Printf("%s\n", ls.GetSha(CHECK_GROUP_ALL_STATUS_LUA))
}

func BenchmarkScriptLoader_GetShaParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ls.GetSha(GET_BATCH_VALUES_LUA)
			ls.GetSha(CHECK_GROUP_ALL_STATUS_LUA)
		}
	})
}

func BenchmarkScriptLoader_GetSha(b *testing.B) {
	b.SetBytes(102400)
	b.SetParallelism(1000)
	for i:=0 ; i<b.N ; i++ {
		ls.GetSha(GET_BATCH_VALUES_LUA)
		ls.GetSha(CHECK_GROUP_ALL_STATUS_LUA)
	}
}