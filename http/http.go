package http

import (
	"errors"
	"io"
	"io/ioutil"
	h "net/http"
	"net/url"
)

// Get 获取url对应的文件内容
// header: header 请求体的header map,例如Content-Type:application/json
// 返回到content中
func Get(uri string) (response *Resp, e error) {
	client := &h.Client{}
	resp, err := client.Get(uri)
	if err != nil {
		return nil, err
	}
	return getRespBody(resp)
}

// Post will post数据
// uri: post到的地址
// header: header 请求体的header map,例如Content-Type:application/json
// body: post的数据
func Post(uri string, header map[string]string, body io.Reader) (response *Resp, err error) {
	req, err := h.NewRequest(h.MethodPost, uri, body)
	for key, item := range header {
		req.Header.Add(key, item)
	}
	client := &h.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return getRespBody(resp)
}

// PostForm will post form数据
// uri: post到的地址
// body: form数据
func PostForm(uri string, data url.Values) (response *Resp, err error) {
	client := &h.Client{}
	resp, err := client.PostForm(uri, data)
	if err != nil {
		return nil, err
	}
	return getRespBody(resp)
}

// Put will put数据
// uri: post到的地址
// header: header 请求体的header map,例如Content-Type:application/json
// body: post的数据
func Put(uri string, header map[string]string, body io.Reader) (response *Resp, err error) {
	req, err := h.NewRequest(h.MethodPut, uri, body)
	for key, item := range header {
		req.Header.Add(key, item)
	}
	client := &h.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return getRespBody(resp)
}

// Delete will Delete数据
// uri: post到的地址
// header: header 请求体的header map,例如Content-Type:application/json
// body: post的数据
func Delete(uri string, header map[string]string, body io.Reader) (response *Resp, err error) {
	req, err := h.NewRequest(h.MethodDelete, uri, body)
	for key, item := range header {
		req.Header.Add(key, item)
	}
	client := &h.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return getRespBody(resp)
}

func getRespBody(resp *h.Response) (response *Resp, err error) {
	if resp == nil {
		return nil, errors.New("response is nil")
	}
	var content []byte
	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Resp{
		header:  resp.Header,
		content: content,
		code:    resp.StatusCode,
	}, nil
}
