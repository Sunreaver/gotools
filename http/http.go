package http

import (
	"io/ioutil"
	h "net/http"
)

// Get 获取url对应的文件内容
// 返回到content中
func Get(url string) (content string, err error) {
	resp, err := h.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		err = err2
		return
	}
	content = string(data)
	return
}
