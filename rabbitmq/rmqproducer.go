package rmq

import (
	"net"
	"time"
	"github.com/streadway/amqp"
	"baselib/logger"
	"fmt"
)

type TaskSheet struct {
	RoutingKey string
	MsgData []byte
}

type RmqProducer struct {
	core    *amqp.Connection
	mqchan  *amqp.Channel
	connerr chan *amqp.Error

	uri       string
	exchange string
	TaskStream chan TaskSheet
}

func NewRmqProducer(uri string, exchange string) (rmq *RmqProducer) {
	rmq = &RmqProducer{
		uri:       uri,
		exchange: exchange,
		TaskStream: make(chan TaskSheet, 256),
	}

	go rmq.mainRoutine()
	return
}

func (rmq *RmqProducer) mainRoutine() {
	for {
		if rmq.core == nil {
			if err := rmq.initCore(); err != nil {
			} else if err = rmq.setupContext(); err != nil {
			}
		} else {
			select {
			case err := <-rmq.connerr:
				logger.Error(fmt.Sprintf(
					"producer(%p) lost connection lost with rabbitmq, err(%v)",
					rmq, err))
				rmq.closeConn()
			case task := <-rmq.TaskStream:
				rmq.publishMsg(task)	// will always print log
			}
		}
	}
}

func (rmq *RmqProducer) closeConn() {
	logger.Debug(fmt.Sprintf("producer(%p) closing connection", rmq))
	if rmq.core != nil {
		rmq.core.Close()
		rmq.core = nil
	}
	rmq.mqchan = nil
	rmq.connerr = nil
	<-time.After(3 * time.Second)
}

func (rmq *RmqProducer) initCore() (err error) {
	if rmq.core, err = amqp.DialConfig(rmq.uri, amqp.Config{
		Vhost: vhost,
		Heartbeat: 3 * time.Second,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 2*time.Second)
		},
	}); err != nil {
	} else {
		rmq.connerr = make(chan *amqp.Error)
		rmq.core.NotifyClose(rmq.connerr)
	}

	if err != nil {
		rmq.closeConn()
	}
	logger.Info(fmt.Sprintf(
		"producer(%p) got connection(%p) with rabbitmq(%s) err(%v)",
		rmq, rmq.core, rmq.uri, err))
	return
}

func (rmq *RmqProducer) setupContext() (err error) {
	//XXX the creation of channel and queue MUST be in same goroutine with `Consume`, otherwise the channel will disappear
	//XXX block or not, this is the only choice
	if rmq.mqchan, err = rmq.core.Channel(); err != nil {
	} else if err = rmq.mqchan.ExchangeDeclare(
		rmq.exchange,"topic",true,false,false,false,nil,
	); err != nil {
		logger.Error(fmt.Sprintf(
			"producer(%p) failed declaring exchange(%s), err(%v)",
			rmq, rmq.exchange, err))
	} else if err = rmq.mqchan.Qos(1, 0, false); err != nil {
		logger.Error(fmt.Sprintf(
			"producer(%p) failed setting Qos, err(%v)",
			rmq, err))
	}

	if err != nil {
		rmq.closeConn()
	}

	return
}

func (rmq *RmqProducer) publishMsg(task TaskSheet) (err error) {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp: time.Now(),
		ContentType: "text/plain",
		Body: task.MsgData,
	}

	err = rmq.mqchan.Publish(rmq.exchange, task.RoutingKey, false, false, msg)
	if err != nil {
		logger.Warn(fmt.Sprintf(
			"producer(%p) failed publishing message(%s) upon exchange(%s) with routingkey(%s)",
			rmq, string(task.MsgData), rmq.exchange, task.RoutingKey))
	} else {
		logger.Debug(fmt.Sprintf(
			"producer(%p) published message(%s) upon exchange(%s) with routingkey(%s)",
			rmq, string(task.MsgData), rmq.exchange, task.RoutingKey))
	}

	return
}
