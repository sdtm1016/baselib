package rmq

import (
	"baselib/logger"
	"baselib/config"
)

const(
	exchangePub = "com.baselibmq.pub.ex"

	//routeKey for Producer
	routeKey_connSuc = "routeKey_connSuc"

)

var producer *RmqProducer

//收到待同步的数据，按照routeKey进行区分然后做相应的同步
func InitForEmitData() {
	logger.Info("initForEmitData() enter,")

	producer = NewRmqProducer(config.GetConf().RabbitMQ.Addr, exchangePub)
}

func getProducerObj() *RmqProducer {
	if producer == nil {
		logger.Error("getProducerObj(), producer obj is nill")
		return nil
	}

	return producer
}
