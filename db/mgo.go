package db

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"baselib/config"
	"strings"
	"baselib/logger"
	"time"
	"sync"
)

type mongoWrapper struct {
	session	*mgo.Session
	db   	*mgo.Database
	mux   	*sync.RWMutex
	status 	int8
}

//status mongo状态
const (
	NOT_INIT = iota	//0 未初始化
	CONNECTING		//1 正在连接
	CONNECTED		//2连接成功
)

var (
	mongo = &mongoWrapper{
		mux: new(sync.RWMutex),
		status: NOT_INIT,
	}
)

func (m *mongoWrapper) getDb() *mgo.Database {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.db
}

func (m *mongoWrapper) getMs() *mgo.Session {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.session
}

func (m *mongoWrapper) setMs(ms *mgo.Session) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.db = ms.DB(config.GetConf().MongoDb.Dbname)
	m.session = ms
	m.status = CONNECTED
}

func (m *mongoWrapper) getStatus() int8 {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.status
}

func (m *mongoWrapper) setStatus(status int8) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.status = status
}

func GetMs() *mgo.Session {
	err := mongo.getMs().Ping()
	if err != nil {
		fmt.Errorf("%s", err.Error())
		if mongo.getStatus() != CONNECTING {
			InitMgo()
		} else {
			logger.Info("other thread initing mongo, wait")
			ticker := time.NewTicker(time.Second * 2)
			for _ = range ticker.C {
				logger.Info("waiting")
				if mongo.getMs().Ping() == nil {
					ticker.Stop()
					break
				} else {
					logger.Info("waiting for other thread initing mongo...")
				}
			}
		}
	}

	return mongo.getMs()
}

func GetDb() *mgo.Database {
	if mongo.getDb() == nil || mongo.getMs() == nil {
		GetMs()
	}
	if mongo.getMs().Ping() != nil {
		logger.Error("conn error")
		GetMs()
	}
	return mongo.getDb()
}

func InitMgo() {
	mongo.setStatus(CONNECTING)

	//格式 mongodb://user:passwd@localhost:40001,otherhost:40001/mydb
	var(
		connstr = "mongodb://"
		role = config.GetConf().MongoDb.Role
		userpass = config.GetConf().MongoDb.Userpass
		hostsports = config.GetConf().MongoDb.Hostsports
		replicaset = config.GetConf().MongoDb.ReplicaSet
	)

	if !strings.EqualFold("", userpass) {
		connstr += userpass + "@"
	}

	connstr += fmt.Sprintf("%s/%s", hostsports, role)
	if !strings.EqualFold("", replicaset) {
		connstr += "?replicaSet=" + replicaset
	}

	session, err := mgo.DialWithTimeout(connstr, time.Second*5)
	for {
		if err != nil {
			logger.Info(err.Error(), "retry to conn... ")
			session, err = mgo.DialWithTimeout(connstr, time.Second*1)
		} else {
			logger.Info("conn success. prepare to ping")
			err = session.Ping()
			if err != nil {
				logger.Error(err.Error(), "ping Mongo error")
			} else {
				logger.Info("ping Mongo success.")
				break
			}
		}
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	logger.Info("succeed connect to mongo")

	mongo.setMs(session)
}
