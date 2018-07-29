package etcd

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "baselib/logger"
	"github.com/coreos/etcd/clientv3"
)

const (
	DialTimeout    = 3
	LeaseGrantTTL  = 5
	ContextTimeout = 2
)

//IEtcd a etcd client
type IEtcd struct {
	addr      []string          //etcd server address
	grantID   clientv3.LeaseID  //lease id that communicate with etcd server
	muGrantID sync.Mutex        //protect grantID
	mu        sync.RWMutex      //prootect data
	data      map[string]string //data that write into etcd
	Cli       *clientv3.Client  //etcd client.
	exitCh    chan byte         //notice exit goroutine
}

//NewIEtcd 初始化etcd模块.
func NewIEtcd(addr []string, register bool) (*IEtcd, error) {
	var (
		err error
		i   *IEtcd
	)

	i = &IEtcd{
		addr:    addr,
		exitCh:  make(chan byte, 1),
		data:    make(map[string]string),
		grantID: clientv3.NoLease,
	}

	i.Cli, err = i.loadCli()
	if err != nil {
		log.Error(err)
	}

	if register {
		go i.registerAndKeepAlive()
	}
	return i, err
}

//注册key/value to etcd server.
//key: ie:"microservice001/136.23.3.2:6590", "microservice002/136.23.3.2:6590"
//value: ie:"136.23.3.2:6590", "www.xxx.com:9001"
func (i *IEtcd) Register(key, value string) error {
	i.writeMap(key, value)
	return i.writeEtcd(key, value)
}

//注销key/value to etcd server.
//key: ie:"microservice001/136.23.3.2:6590", "microservice002/136.23.3.2:6590"
func (i *IEtcd) UnRegister(key string) error {
	i.deleteMap(key)
	return i.deleteEtcd(key)
}

//Close release all resouce
func (i *IEtcd) Close() {
	i.exitCh <- 0
	defer close(i.exitCh)

	i.Cli.Close()
	log.Debug("exit goroutine")
}

func (i *IEtcd) writeEtcd(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout*time.Second)
	_, err := i.Cli.Put(ctx, key, value, clientv3.WithLease(i.etcdGrant()))
	cancel()
	if err != nil {
		log.Error(err)
	}
	return err
}

func (i *IEtcd) deleteEtcd(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout*time.Second)
	_, err := i.Cli.Delete(ctx, key)
	cancel()
	if err != nil {
		log.Error(err)
	}
	return err
}

func (i *IEtcd) writeMap(key, value string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.data[key] = value
}

func (i *IEtcd) deleteMap(key string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.data, key)
}

func (i *IEtcd) readMap() map[string]string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.data
}

//keepAlive 申请租约并keepAlive.
func (i *IEtcd) registerAndKeepAlive() {
	var (
		err           error
		grantID       clientv3.LeaseID
		allStart      chan byte
		putChan       chan byte
		keepAliveChan chan byte
	)

	log.Debug("start registerAndKeepAlive...")
	allStart = make(chan byte, 1)
	putChan = make(chan byte, 1)
	keepAliveChan = make(chan byte, 1)
	allStart <- 1

	defer func() {
		close(allStart)
		close(putChan)
		close(keepAliveChan)
	}()

	for {
		select {
		case <-allStart:
			log.Debug("all start")
			grantID = i.genGrantID()
			if grantID > clientv3.NoLease {
				putChan <- 1
			} else {
				log.Error("keepAlive etcdGrant id not right")
				allStart <- 1
			}
			log.Debug("leaseId %d", grantID)
		case <-putChan:
			//write etcd
			log.Debug("keepAlive put value")
			i.mu.RLock()
			for k, v := range i.data {
				ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout*time.Second)
				_, err = i.Cli.Put(ctx, k, v, clientv3.WithLease(i.etcdGrant()))
				cancel()
				if err != nil {
					log.Error(err)
					break
				}
			}
			i.mu.RUnlock()

			if err == nil {
				keepAliveChan <- 1
			} else {
				allStart <- 1
			}

		case <-keepAliveChan:
			log.Debug("etcd keep alive")
			c, err := i.Cli.KeepAlive(context.Background(), i.etcdGrant())
			if err != nil {
				log.Error("keep alive err:%+v", err)
				allStart <- 1
				break
			}
		keepAlive:
			for {
				select {
				case <-allStart:
					log.Debug("restart when keepalive")
					allStart <- 1
					break keepAlive
				case _, ok := <-c:
					if !ok {
						log.Error("keep alive no responce, should all start again")
						allStart <- 1
						break keepAlive
					}
				case <-i.exitCh:
					log.Debug("exit")
					return
				case <-time.After(5 * time.Minute):
					log.Error("should not here")
					allStart <- 1
				}
			}

		case <-i.exitCh:
			log.Debug("exit")
			return
		}
	}
}

//get grantId
func (i *IEtcd) etcdGrant() clientv3.LeaseID {
	i.muGrantID.Lock()
	defer i.muGrantID.Unlock()

	if i.grantID > clientv3.NoLease {
		return i.grantID
	}

	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout*time.Second)
	grant, grantErr := i.Cli.Grant(ctx, LeaseGrantTTL)
	cancel()
	if grantErr != nil {
		log.Error(fmt.Sprintf("etcdGrant err.%v", grantErr))
		return clientv3.NoLease
	}
	i.grantID = grant.ID

	return grant.ID
}

//generate grantId
func (i *IEtcd) genGrantID() clientv3.LeaseID {
	i.muGrantID.Lock()
	defer i.muGrantID.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout*time.Second)
	grant, grantErr := i.Cli.Grant(ctx, LeaseGrantTTL)
	cancel()
	if grantErr != nil {
		log.Error(fmt.Sprintf("etcdGrant err.%v", grantErr))
		return clientv3.NoLease
	}
	i.grantID = grant.ID

	return grant.ID
}

func (i *IEtcd) loadCli() (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   i.addr,
		DialTimeout: DialTimeout * time.Second,
	})
}
