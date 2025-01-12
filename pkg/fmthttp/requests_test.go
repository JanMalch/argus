package fmthttp_test

import (
	"testing"

	"github.com/janmalch/argus/pkg/fmthttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineString(t *testing.T) {
	rl := fmthttp.RequestLine{
		Method:        "POST",
		RequestTarget: "/users",
		Proto:         "HTTP/1.1",
	}
	assert.Equal(t, reqLine, rl.String())
}

func TestRequestHeadString(t *testing.T) {
	expected := joinCrLf(reqLine,
		"Host: example.com",
		"X-Foo: 1",
		"X-Foo: 2",
		"\r\n",
	)
	require.Equal(t, "\r\n\r\n", last4(expected))
	rh := fmthttp.RequestHead{
		RequestLine: fmthttp.RequestLine{
			Method:        "POST",
			RequestTarget: "/users",
			Proto:         "HTTP/1.1",
		},
		Headers: fmthttp.NewHeaders(
			"X-Foo", "1",
			"Host", "example.com",
			"X-Foo", "2",
		),
	}
	assert.Equal(t, expected, rh.String())
}
