package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/janmalch/argus/pkg/fmthttp"
)

type Request struct {
	fmthttp.RequestHead
	Timestamp time.Time
	Url       string
}

type Response struct {
	fmthttp.ResponseHead
	Timestamp time.Time
}

type Exchange struct {
	Id       uint64
	Request  Request
	Response *Response
}

type Hooks interface {
	OnRequest(e *Exchange, r *http.Request) (io.ReadCloser, error)
	OnRequestWithoutFurtherBodyUsage(e *Exchange, r *http.Request) error
	OnResponse(e *Exchange, body io.ReadCloser, contentType string) (string, error)
	ReadFile(file string) (io.ReadCloser, string, error)
	Log(id uint64, msg string)
}
