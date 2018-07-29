package http_request

import (
	"net/http"
	"encoding/json"
	"time"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"baselib/logger"
)

type NewRequest struct {
	headerInfo map[string]string
	req interface{}
	rsp interface{}
	url string
}

func (r *NewRequest) SetHeader(header map[string]string) {
	r.headerInfo = header
}

func (r *NewRequest) GetHeader() map[string]string {
	return r.headerInfo
}

func (r *NewRequest) SetBody(body interface{}) {
	r.req = body
}

func (r *NewRequest) GetBody() interface{} {
	return r.req
}

func (r *NewRequest) SetResponse(response interface{}) {
	r.rsp = response
}

func (r *NewRequest) GetResponse() interface{} {
	return r.rsp
}

func (r *NewRequest) SetUrl(url string) {
	r.url = url
}

func (r *NewRequest) GetUrl() string {
	return r.url
}

// http get 请求
func (r *NewRequest) Get() ([]byte, error) {
	client := &http.Client{}
	logger.Info(r.url)
	request, err := http.NewRequest(http.MethodGet, r.url, nil)
	for k, v := range r.headerInfo {
		request.Header.Add(k, v)
	}

	if err != nil {
		errInfo := errors.New(fmt.Sprintf("http.NewRequest error %v: url:%s", err, r.url))
		logger.Error(errInfo)
		return nil, errInfo
	}

	response, err := client.Do(request)
	if err != nil {
		errInfo := errors.New(fmt.Sprintf("do request error %v: url:%s", err, r.url))
		logger.Error(errInfo)
		return nil, errInfo
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errInfo := errors.New(fmt.Sprintf("get response body error %v: url:%s", err, r.url))
		logger.Error(errInfo)
		return nil, errInfo
	}

	if response.StatusCode == 200 {
		return body, nil
	}

	errInfo := errors.New(fmt.Sprintf("response code error, body:%s, url:%s", string(body),r.url))
	logger.Error(errInfo)
	return nil, errInfo
}

//http post 请求
func (r *NewRequest) Post() error {
	http.DefaultClient.Timeout = time.Second * 2
	bts, err := json.Marshal(r.req)
	if err != nil {
		return err
	}
	logger.Info("req is %s", string(bts))

	request, err := http.NewRequest(http.MethodPost, r.url, bytes.NewBuffer(bts))
	if err != nil {
		return err
	}

	for k, v := range r.headerInfo {
		request.Header.Add(k, v)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	logger.Info("success to get response body:%s", string(all))
	err = json.Unmarshal(all, r.rsp)
	if err != nil {
		return err
	}

	return nil
}