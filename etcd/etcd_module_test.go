package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/stretchr/testify/assert"
)

func TestWriteAndRead(t *testing.T) {
	etcd, err := NewIEtcd([]string{"192.168.96.6:2379"}, true)
	assert.NotNil(t, etcd)
	assert.NoError(t, err)
	defer etcd.Close()

	err = etcd.Register("service1/192.168.96.140:2379", "192.168.96.140:2379")
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)

	err = etcd.Register("service2/192.168.96.141:2379", "192.168.96.141:2379")
	assert.NoError(t, err)
	time.Sleep(8 * time.Second)

	response, e := etcd.Cli.Get(context.Background(), "service1/", clientv3.WithPrefix())
	assert.Nil(t, e)

	assert.Equal(t, "192.168.96.140:2379", string(response.Kvs[0].Value))
	assert.Equal(t, "192.168.96.141:2379", string(response.Kvs[1].Value))
}

func TestWrite(t *testing.T) {
	etcd, err := NewIEtcd([]string{"192.168.96.6:2379"}, true)
	assert.NotNil(t, etcd)
	assert.NoError(t, err)
	defer etcd.Close()

	err = etcd.Register("service1/192.168.96.140:2379", "192.168.96.140:2379")
	assert.NoError(t, err)

	err = etcd.Register("service2/192.168.96.141:2379", "192.168.96.141:2379")
	assert.NoError(t, err)
}

func TestWatcher(t *testing.T) {
	etcd, err := NewIEtcd([]string{"192.168.96.6:2379"}, true)
	assert.NotNil(t, etcd)
	assert.NoError(t, err)
	defer etcd.Close()

	err = etcd.Register("service1/192.168.96.140:2379", "192.168.96.140:2379")
	assert.NoError(t, err)

	err = etcd.Register("service1/192.168.96.141:2379", "192.168.96.141:2379")
	assert.NoError(t, err)

	rch := etcd.Cli.Watch(context.Background(), "service1/", clientv3.WithPrefix())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for {
		select {
		case wresp := <-rch:
			for _, ev := range wresp.Events {
				fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		case <-time.After(3 * time.Second):
			err = etcd.UnRegister("service1/192.168.96.141:2379")
			assert.NoError(t, err)
		case <-ctx.Done():
			return
		}
	}

}

func TestLongtimeNoKeepalive(t *testing.T) {
	etcd, err := NewIEtcd([]string{"192.168.96.6:2379"}, true)
	assert.NotNil(t, etcd)
	assert.NoError(t, err)
	//defer etcd.Close()

	err = etcd.Register("service1/192.168.96.140:2379", "192.168.96.140:2379")
	assert.NoError(t, err)

	time.Sleep(8 * time.Second)

	err = etcd.Register("service1/192.168.96.141:2379", "192.168.96.141:2379")
	assert.NoError(t, err)

	fmt.Println("shut down the network")
	time.Sleep(10 * time.Second)
	//shut down the network
	fmt.Println("open the network")

	response, e := etcd.Cli.Get(context.Background(), "service1/", clientv3.WithPrefix())
	assert.Nil(t, e)

	assert.Equal(t, "192.168.96.140:2379", string(response.Kvs[0].Value))
	assert.Equal(t, "192.168.96.141:2379", string(response.Kvs[1].Value))

}
