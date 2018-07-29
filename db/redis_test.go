package db

import (
	"testing"
	"fmt"
	"baselib/logger"
)

const(
	USERINFO_HASH_KEY = "baselib.userinfo.%s.%s"
)

func TestGetRedisClient(t *testing.T) {
	vccId := "vcc001"
	userId := "421181xxxxxx6666"
	userInfoKey := fmt.Sprintf(USERINFO_HASH_KEY, vccId, userId)

	pipeline := GetRedisClient().TxPipeline()
	pipeline.HSet(userInfoKey,"userName","guobin")
	pipeline.HSet(userInfoKey,"userId",userId)
	pipeline.HSet(userInfoKey,"nickName","gb")
	pipeline.HSet(userInfoKey,"age",28)
	pipeline.HSet(userInfoKey,"sex",1)
	_, exeErr := pipeline.Exec()
	if exeErr!=nil {
		logger.Warn("Fail to pipeline.Exec() for userInfoHash. err:",exeErr)
	}
}
