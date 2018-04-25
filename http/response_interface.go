package http

import h "net/http"

// Response define Response interface
type Response interface {
	GetHeader() h.Header
	GetContent() []byte
	GetStatusCode() int
}
