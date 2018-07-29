package encrypt

import (
	"testing"
	"baselib/logger"
	"time"
)

func TestGetMd5(t *testing.T) {
	logger.Info(GetSha1("123"))
	logger.Info(GetMd5("123"))
	logger.Info(time.Now())
	logger.Info(TimeFormat(time.Now()))
	logger.Info(GetTimestamp())
}

