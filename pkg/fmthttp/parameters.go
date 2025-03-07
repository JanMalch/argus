package fmthttp

import (
	"net/url"
	"sort"
)

type Parameter struct {
	Key    string
	Values []string
}

// A collection of query parameters with a stable order.
type Parameters []Parameter

// Creates new Headers from the given key-value pairs.
//
//	h := fmthttp.NewParameters(
//		"a", "1",
//		"b", "Hello%2C%20World!",
//		"c", "2",
//		"c", "3",
//	)
//
// If the number of arguments is odd, the last key is ignored.
func NewParameters(kv ...string) Parameters {
	inLen := len(kv)
	if inLen < 2 {
		return Parameters{}
	}
	lut := make(map[string][]string, 0)
	keys := make([]string, 0)
	i := 0
	for {
		key := kv[i]
		value := kv[i+1]
		group, ok := lut[key]
		if ok {
			lut[key] = append(group, value)
		} else {
			lut[key] = []string{value}
			keys = append(keys, key)
		}
		i += 2
		if i+1 >= inLen {
			break
		}
	}
	sort.Strings(keys)

	p := make([]Parameter, len(lut))
	for li, lk := range keys {
		lv := lut[lk]
		p[li] = Parameter{Key: lk, Values: lv}
	}
	return p
}

func CopyToParameters(url *url.URL) Parameters {
	q := url.Query()

	keys := make([]string, 0)
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := make([]Parameter, len(keys))
	for i, k := range keys {
		vv := q[k]
		h[i] = Parameter{
			Key:    k,
			Values: append(vv[:0:0], vv...),
		}
	}
	return h
}
