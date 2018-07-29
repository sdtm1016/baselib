package etcd

import (
	"fmt"
	"net/http"
	"testing"
)

func BenchmarkLong(b *testing.B) {
	var (
		etcdCli *IEtcd
		err     error
	)

	key := fmt.Sprintf("bench/%s", "192.168.96.141:123")
	if etcdCli, err = NewIEtcd([]string{"192.168.96.6:2379"}, true); err != nil {
		b.Fail()
		return
	}
	defer etcdCli.Close()

	if err = etcdCli.Register(key, "192.168.96.141:123"); err != nil {
		b.Fail()
		return
	}

	go func() {
		http.ListenAndServe("192.168.96.6:4321", nil)
	}()

	select {}
}
