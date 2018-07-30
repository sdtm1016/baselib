package rmq

import (
	"encoding/json"
	"testing"
	"baselib/logger"
	"fmt"
	"time"
)

//示例的一个MQ消息结构
type ConnSucMQ struct {
	VccID     	string	`json:"vccId"`
	UserID    	string  `json:"userId"`
	AgentID   	string	`json:"agentId"`
	ConnSucTime int64  	`json:"connSuccessTime"`
}

func failOnErr(err error, msg string) {
	if err != nil {
		logger.Error(fmt.Sprintf("MQ failOnErr, %s:%s", msg, err))
	}
}

func TestPublishMsg(t *testing.T) {
	var(
		msg = ConnSucMQ{
			VccID:"vcc01",
			UserID:"user01",
			AgentID:"agent01",
			ConnSucTime:time.Now().Unix(),
		}
	)

	msgContent, err := json.Marshal(msg)
	failOnErr(err, "Publish ConnSucMQ, Failed to json format")

	sendObj := getProducerObj()
	if sendObj == nil {
		logger.Info("Publish ConnSucMQ, Fail to send for nil sendObj")
		return
	}

	logger.Info("Publish ConnSucMQ, Will send msg...")
	sendObj.TaskStream <- TaskSheet{
		RoutingKey: routeKey_connSuc,
		MsgData:    []byte(msgContent),
	}
}