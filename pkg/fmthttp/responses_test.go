package fmthttp_test

import (
	"testing"

	"github.com/janmalch/argus/pkg/fmthttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusLineString(t *testing.T) {
	sl := fmthttp.StatusLine{
		Proto:      "HTTP/1.1",
		StatusCode: 200,
		StatusText: "OK",
	}
	assert.Equal(t, "HTTP/1.1 200 OK", sl.String())
}

func TestResponseHeadString(t *testing.T) {
	expected := joinCrLf("HTTP/1.1 200 OK",
		"Host: example.com",
		"X-Foo: 1",
		"X-Foo: 2",
		"\r\n",
	)
	require.Equal(t, "\r\n\r\n", last4(expected))
	rh := fmthttp.ResponseHead{
		StatusLine: fmthttp.StatusLine{
			Proto:      "HTTP/1.1",
			StatusCode: 200,
			StatusText: "OK",
		},
		Headers: fmthttp.NewHeaders(
			"X-Foo", "1",
			"Host", "example.com",
			"X-Foo", "2",
		),
	}
	assert.Equal(t, expected, rh.String())
}
