package http_request

import (
	"testing"
	"fmt"
	"baselib/logger"
)

type PostBody struct {
	Body string		`json:"content"`
}

func TestPost(t *testing.T) {
	logger.Info(fmt.Sprintf("will post msg to url..."))

	//造数据
	var(
		rsp interface{}
		userId = "user001"
		url = "www.google.com/xxx"
		req = PostBody{
			Body:"hello world ^~^",
		}
	)

	postErr := postReq(url,userId,req,rsp)
	if postErr != nil {
		log := fmt.Sprintf("post msg to url failed, url:%s,err:%v",url,postErr)
		logger.Error(log)
	}
}

func postReq(url,userId string, body,rsp interface{}) error {
	holder := NewRequest{
		headerInfo:make(map[string]string),
	}

	var headerInfo = make(map[string]string)
	headerInfo["Content-Type"] = "application/json"
	headerInfo["userId"] = userId

	holder.SetUrl(url)
	holder.SetHeader(headerInfo)
	holder.SetBody(body)
	holder.SetResponse(rsp)

	postErr := holder.Post()
	if postErr != nil {
		return postErr
	}

	return nil
}