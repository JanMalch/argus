package fmthttp

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type StatusLine struct {
	Proto      string
	StatusCode int    // e.g. 200
	StatusText string // e.g. "OK"
}

// e.g. "200 OK"
func (s *StatusLine) Status() string {
	return fmt.Sprintf("%d %s", s.StatusCode, s.StatusText)
}

func (s *StatusLine) String() string {
	return fmt.Sprintf("%s %d %s", s.Proto, s.StatusCode, s.StatusText)
}

type ResponseHead struct {
	StatusLine
	Headers
}

func NewResponseHead(proto string, statusCode int, status string, header http.Header) ResponseHead {
	var statusText string
	if len(status) > 4 {
		statusText = status[4:]
	} else {
		statusText = http.StatusText(statusCode)
	}
	return ResponseHead{
		StatusLine: StatusLine{Proto: proto, StatusCode: statusCode, StatusText: statusText},
		Headers:    CopyToHeaders(header),
	}
}

func CopyResponseHead(r *http.Response) ResponseHead {
	var statusText string
	if len(r.Status) > 4 {
		statusText = r.Status[4:]
	} else {
		statusText = http.StatusText(r.StatusCode)
	}
	return ResponseHead{
		Headers: CopyToHeaders(r.Header),
		StatusLine: StatusLine{
			Proto:      r.Proto,
			StatusCode: r.StatusCode,
			StatusText: statusText,
		},
	}
}

func (r *ResponseHead) String() string {
	var sb strings.Builder
	sb.WriteString(r.StatusLine.String() + "\r\n")
	for _, e := range r.Headers {
		for _, v := range e.Values {
			sb.WriteString(fmt.Sprintf("%s: %s\r\n", e.Key, v))
		}
	}
	return sb.String() + "\r\n"
}

type Response struct {
	ResponseHead
	Body io.ReadCloser
}

func CopyResponse(r *http.Response) Response {
	return Response{
		ResponseHead: CopyResponseHead(r),
		Body:         r.Body,
	}
}
