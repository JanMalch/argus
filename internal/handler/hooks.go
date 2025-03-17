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
	AddRequest(id uint64, req *http.Request, timestamp time.Time) (string, error)
	AddResponse(id uint64, res *fmthttp.Response, timestamp time.Time) (string, error)
	ReadFile(file string) (io.ReadCloser, string, error)
	Log(id uint64, msg string)
}
