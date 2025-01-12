package fmthttp_test

import (
	"net/http"
	"testing"

	"github.com/janmalch/argus/pkg/fmthttp"
	"github.com/stretchr/testify/assert"
)

func TestNewHeadersEven(t *testing.T) {
	actual := fmthttp.NewHeaders(
		"X-Foo", "1",
		"Host", "example.com",
		"X-Foo", "2",
	)
	assert.EqualValues(t, []fmthttp.Header{
		{
			Key:    "Host",
			Values: []string{"example.com"},
		},
		{
			Key:    "X-Foo",
			Values: []string{"1", "2"},
		},
	}, actual)
}

func TestNewHeadersNone(t *testing.T) {
	actual := fmthttp.NewHeaders()
	assert.EqualValues(t, []fmthttp.Header{}, actual)
}

func TestNewHeadersOne(t *testing.T) {
	actual := fmthttp.NewHeaders("Orphan")
	assert.EqualValues(t, []fmthttp.Header{}, actual)
}

func TestNewHeadersOdd(t *testing.T) {
	actual := fmthttp.NewHeaders(
		"X-Foo", "1",
		"Host", "example.com",
		"X-Foo", "2",
		"Orphan",
	)
	assert.EqualValues(t, []fmthttp.Header{
		{
			Key:    "Host",
			Values: []string{"example.com"},
		},
		{
			Key:    "X-Foo",
			Values: []string{"1", "2"},
		},
	}, actual)
}

func TestCopyToHeaders(t *testing.T) {
	actual := fmthttp.CopyToHeaders(http.Header{
		"Host": {"example.com"},
		// can't test with more headers, because iteration order of maps is random
	})
	assert.EqualValues(t, []fmthttp.Header{
		{
			Key:    "Host",
			Values: []string{"example.com"},
		},
	}, actual)
}
