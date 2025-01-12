package fmthttp_test

import (
	"strings"
)

func last4(s string) string {
	return s[len(s)-4:]
}

func joinCrLf(lines ...string) string {
	return strings.Join(lines, "\r\n")
}

var jsonData = `{
    "data": "ABC123"
}`

var reqLine = "POST /users HTTP/1.1"
var reqHead = joinCrLf(
	reqLine,
	"Host: example.com",
	"\r\n",
)
var jsonRequest = joinCrLf(
	reqLine,
	"Host: example.com",
	"",
	jsonData,
)

var statusLine = "HTTP/1.1 200 OK"
var resHead = joinCrLf(
	statusLine,
	"Host: example.com",
	"\r\n",
)

var jsonResponse = joinCrLf(
	statusLine,
	"Host: example.com",
	"",
	jsonData,
)
