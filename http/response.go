package http

import h "net/http"

// Resp http response
type Resp struct {
	header  h.Header
	content []byte
	code    int
}

// GetHeader 返回resp的header
func (r *Resp) GetHeader() h.Header {
	if r != nil {
		return r.header
	}
	return nil
}

// GetContent 返回response的[]byte格式内容
func (r *Resp) GetContent() []byte {
	if r != nil {
		return r.content
	}
	return nil
}

// GetStatusCode 返回response的http code
func (r *Resp) GetStatusCode() int {
	if r != nil {
		return r.code
	}
	return 0
}
