package http

import (
	"errors"
	"io"
	"io/ioutil"
	h "net/http"
)

// Get 获取url对应的文件内容
// 返回到content中
func Get(uri string) (content []byte, httpCode int, e error) {
	resp, err := h.Get(uri)
	if err != nil {
		return nil, 0, err
	}
	return getRespBody(resp)
}

// Post will post数据
// uri: post到的地址
// contentType: header Content-Type
// body: post的数据
func Post(uri, contentType string, body io.Reader) (content []byte, statusCode int, err error) {
	resp, err := h.Post(uri, contentType, body)
	if err != nil {
		return nil, 0, err
	}
	return getRespBody(resp)
}

func getRespBody(resp *h.Response) (content []byte, statusCode int, err error) {
	if resp == nil {
		return nil, 0, errors.New("response is nil")
	}
	content, err = ioutil.ReadAll(resp.Body)
	if resp.Body != nil {
		resp.Body.Close()
	}
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return content, resp.StatusCode, nil
}
