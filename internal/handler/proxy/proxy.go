package proxy

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/janmalch/argus/pkg/fmthttp"
)

type ProxyRequest struct {
	upurl *url.URL
	upreq *http.Request
}

func PrepareProxyRequest(
	id uint64,
	r *http.Request,
	upstream *url.URL,
	header map[string]string,
	query map[string]string,
) (*ProxyRequest, error) {
	upurl := PrepareUrl(id, r, upstream, query)
	upreq, err := http.NewRequest(r.Method, upurl.String(), http.NoBody)
	if err != nil {
		return nil, err
	}
	upreq.Header = PrepareHeaders(id, &r.Header, header)
	// https://stackoverflow.com/a/13131806
	upreq.Header.Del("Accept-Encoding")

	return &ProxyRequest{upurl, upreq}, nil
}

func (p *ProxyRequest) SetBody(body io.ReadCloser) *ProxyRequest {
	p.upreq.Body = body
	return p
}

func (p *ProxyRequest) Url() string {
	return p.upurl.String()
}

func (p *ProxyRequest) Do() (*http.Response, error) {
	return http.DefaultClient.Do(p.upreq)
}

func (p *ProxyRequest) RequestHead() fmthttp.RequestHead {
	h := fmthttp.CopyRequestHead(p.upreq)
	if h.RequestTarget == "" {
		h.RequestTarget = p.upurl.RequestURI()
	}
	return h
}

func PrepareHeaders(id uint64, h *http.Header, c map[string]string) http.Header {
	r := h.Clone()
	for k, v := range c {
		if v != "" {
			r.Set(k, fillPlaceholders(v, id))
		} else {
			r.Del(k)
		}
	}
	return r
}

func fillPlaceholders(p string, id uint64) string {
	s := p
	if strings.Contains(s, "{{id}}") {
		s = strings.ReplaceAll(p, "{{id}}", strconv.FormatUint(id, 10))
	}
	for strings.Contains(s, "{{rng.uuid}}") {
		s = strings.Replace(p, "{{rng.uuid}}", uuid.New().String(), 1)
	}
	return s
}

func PrepareUrl(
	id uint64,
	r *http.Request,
	upstream *url.URL,
	query map[string]string,
) *url.URL {
	// create copy
	up := *r.URL
	up.Host = upstream.Host
	up.Scheme = upstream.Scheme
	q := up.Query()
	for k, v := range query {
		q.Add(k, fillPlaceholders(v, id))
	}
	up.RawQuery = q.Encode()
	return &up
}
