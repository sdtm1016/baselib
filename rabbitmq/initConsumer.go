package rmq

import (
	"baselib/logger"
	"baselib/config"
	"bytes"
)

const(
	exchangeRecv = "com.baselibmq.recv.ex"

	//routeKey for Consumer
	routeKey_configChange = "routeKey_configChange"
	routeKey_powerChange = "routeKey_powerChange"
)

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}

//根据routeKey来区分不同的消息
func InitForRecvData() {

	//var routeKeys []string = make([]string, MAX_ROUTE_REC)
	routeKeys := []string{
		routeKey_configChange,
		routeKey_powerChange,
	}
	consumer := NewRmqConsumer(config.GetConf().RabbitMQ.Addr, exchangeRecv, routeKeys)

	//msg := <- consumer.MsgRecv
	logger.Info("InitForRecvData, consumer listen on MsgRecv ...")

	for msg := range consumer.MsgRecv {
		logger.Debug("got rmq message on exchange(%s) routingkey(%s), data(%v)",
			msg.Exchange, msg.RoutingKey, string(msg.Body))

		switch msg.RoutingKey {
		case routeKey_configChange:
			configChangeHandle(msg.Body)
		case routeKey_powerChange:
			powerChangeHandle(msg.Body)
		default:
			logger.Error("InitForRecvData, err:Unkown routeKey:", msg.RoutingKey)
		}
	}
}

func configChangeHandle(dataReceived []byte) {
	s := BytesToString(&(dataReceived))
	logger.Info("configChangeHandle, receive msg is :%s --- \n", *s)

	//do handle for received msg
	//todo:

}

func powerChangeHandle(dataReceived []byte) {
	s := BytesToString(&(dataReceived))
	logger.Info("powerChangeHandle, receive msg is :%s --- \n", *s)

	//do handle for received msg
	//todo:

}