package fmthttp

import (
	"fmt"
	"net/http"
	"net/textproto"
	"sort"
	"strings"
)

type Header struct {
	Key    string
	Values []string
}

// A collection of HTTP headers with a stable order.
type Headers []Header

func (h Headers) String() string {
	var sb strings.Builder
	for _, e := range h {
		for _, v := range e.Values {
			sb.WriteString(fmt.Sprintf("%s: %s\r\n", e.Key, v))
		}
	}
	return sb.String()
}

func (h Headers) Get(key string) string {
	canon := textproto.CanonicalMIMEHeaderKey(key)
	for _, e := range h {
		if e.Key == canon {
			if len(e.Values) == 0 {
				return ""
			}
			return e.Values[0]
		}
	}
	return ""
}

func (h Headers) LongestName() string {
	maxLen := 0
	maxStr := ""
	for _, e := range h {
		eLen := len(e.Key)
		if eLen > maxLen {
			maxLen = eLen
			maxStr = e.Key
		}
	}
	return maxStr
}

// Creates new Headers from the given key-value pairs.
//
//	h := fmthttp.NewHeaders(
//		"Host", "example.com",
//		"X-Foo", "1",
//		"X-Foo", "2",
//	)
//
// If the number of arguments is odd, the last key is ignored.
func NewHeaders(kv ...string) Headers {
	inLen := len(kv)
	if inLen < 2 {
		return []Header{}
	}
	h := http.Header{}
	i := 0
	for {
		h.Add(kv[i], kv[i+1])
		i += 2
		if i+1 >= inLen {
			break
		}
	}
	return CopyToHeaders(h)
}

func CopyToHeaders(header http.Header) Headers {
	size := len(header)

	i := 0
	keys := make([]string, size)
	for k := range header {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	i = 0
	h := make([]Header, size)
	for _, k := range keys {
		vv := header.Values(k)
		h[i] = Header{
			Key:    k,
			Values: append(vv[:0:0], vv...),
		}
		i++
	}
	return h
}
