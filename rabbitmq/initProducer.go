package rmq

import (
	"baselib/logger"
	"baselib/config"
)

var producer *RmqProducer

//收到待同步的数据，按照routeKey进行区分然后做相应的同步
func InitForEmitData() {
	logger.Info("initForEmitData() enter,")

	producer = NewRmqProducer(config.GetConf().RabbitMQ.Addr, exchange)
}

func getProducerObj() *RmqProducer {
	if producer == nil {
		logger.Error("getProducerObj(), producer obj is nill")
		return nil
	}

	return producer
}
