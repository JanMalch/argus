package fmthttp

import (
	"fmt"
	"net/http"
	"strings"
)

type RequestLine struct {
	Proto         string
	Method        string
	RequestTarget string
}

func (r *RequestLine) String() string {
	return fmt.Sprintf("%s %s %s", r.Method, r.RequestTarget, r.Proto)
}

type RequestHead struct {
	RequestLine
	Headers
}

func NewRequestHead(proto, method, requestTarget string, header http.Header) RequestHead {
	return RequestHead{
		Headers: CopyToHeaders(header),
		RequestLine: RequestLine{
			Proto:         proto,
			Method:        method,
			RequestTarget: requestTarget,
		},
	}
}

func CopyRequestHead(r *http.Request) RequestHead {
	return RequestHead{
		Headers: CopyToHeaders(r.Header),
		RequestLine: RequestLine{
			Proto:         r.Proto,
			Method:        r.Method,
			RequestTarget: r.RequestURI,
		},
	}
}

func (r *RequestHead) String() string {
	var sb strings.Builder
	sb.WriteString(r.RequestLine.String() + "\r\n")
	for _, line := range r.Headers {
		for _, v := range line.Values {
			sb.WriteString(fmt.Sprintf("%s: %s\r\n", line.Key, v))
		}
	}
	return sb.String() + "\r\n"
}
