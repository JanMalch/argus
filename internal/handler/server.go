package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/janmalch/argus/internal/config"
	"github.com/janmalch/argus/internal/handler/proxy"
	"github.com/janmalch/argus/pkg/fmthttp"
)

func NewServer(h Hooks, getServerConfig config.ServerProvider) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", handleEverything(h, getServerConfig))
	return mux
}

var ops atomic.Uint64

func handleEverything(h Hooks, getServerConfig config.ServerProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		id := ops.Add(1)
		defer func() {
			if r := recover(); r != nil {
				h.Log(id, fmt.Sprintf("panic: %+v\n%s", r, debug.Stack()))
			}
		}()

		conf := getServerConfig()

		overwrite := findMatching(conf.Response.Overwrites, r.RequestURI, r.Method)
		if overwrite != nil {
			if overwrite.File == "" {
				if overwrite.Status != 0 {
					// only status defined
					e := &Exchange{
						Id: id,
						Request: Request{
							Timestamp:   start,
							RequestHead: fmthttp.CopyRequestHead(r),
							Url:         proxy.PrepareUrl(id, r, conf.Upstream, conf.Request.Parameters).String(),
						},
					}
					err := h.OnRequestWithoutFurtherBodyUsage(e, r)
					if err != nil {
						h.Log(id, err.Error())
					}

					end := time.Now()
					e.Response = &Response{
						Timestamp:    end,
						ResponseHead: fmthttp.NewResponseHead(r.Proto, overwrite.Status, "", w.Header()),
					}
					w.WriteHeader(overwrite.Status)
					h.OnResponse(e, http.NoBody, "")
					return
				}
				// else: no file and no status isn't a case. fallthrough to upstream
			} else {
				// file is defined
				status := http.StatusOK
				if overwrite.Status > 0 {
					status = overwrite.Status
				}

				ofr, contentType, err := h.ReadFile(overwrite.File)
				if err != nil {
					h.Log(id, err.Error())
					// fallthrough to upstream
				} else {
					e := &Exchange{
						Id: id,
						Request: Request{
							Timestamp:   start,
							RequestHead: fmthttp.CopyRequestHead(r),
							Url:         proxy.PrepareUrl(id, r, conf.Upstream, conf.Request.Parameters).String(),
						},
					}
					err := h.OnRequestWithoutFurtherBodyUsage(e, r)
					if err != nil {
						h.Log(id, err.Error())
					}

					w.Header().Add("Content-Type", contentType)
					end := time.Now()
					e.Response = &Response{
						Timestamp:    end,
						ResponseHead: fmthttp.NewResponseHead(r.Proto, status, "", w.Header()),
					}

					// this will create a copy of the referenced file, but this way we have immutability
					downbodyFile, err := h.OnResponse(e, ofr, contentType)
					if err != nil {
						h.Log(id, err.Error())
						w.WriteHeader(http.StatusNotImplemented)
						return
					}
					downbody, err := os.Open(downbodyFile)
					if err != nil {
						h.Log(id, err.Error())
						w.WriteHeader(http.StatusNotImplemented)
						return
					}
					w.WriteHeader(status)

					_, err = io.Copy(w, downbody)
					if err != nil {
						// FIXME: all these error cases don't update the UI appropriately.
						//        h.OnError(id, err) which can also do the logging?
						h.Log(id, err.Error())
						return
					}
					if err = downbody.Close(); err != nil {
						h.Log(id, err.Error())
					}
					return
				}

			}
		}

		upreq, err := proxy.PrepareProxyRequest(id, r, conf.Upstream, conf.Request.Headers, conf.Request.Parameters)
		if err != nil {
			h.Log(id, err.Error())
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		e := &Exchange{
			Id: id,
			Request: Request{
				Timestamp:   start,
				RequestHead: upreq.RequestHead(),
				Url:         upreq.Url(),
			},
		}

		reqbody, err := h.OnRequest(e, r)
		if err != nil {
			h.Log(id, err.Error())
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		// reqbody will be closed by the transport layer
		upres, err := upreq.SetBody(reqbody).Do()
		end := time.Now()
		if err != nil {
			h.Log(id, err.Error())
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		downheaders := proxy.PrepareHeaders(id, &upres.Header, conf.Response.Headers)
		e.Response = &Response{
			Timestamp: end,
			ResponseHead: fmthttp.NewResponseHead(
				upres.Proto,
				upres.StatusCode,
				upres.Status,
				downheaders,
			),
		}
		for dhk, dhvv := range downheaders {
			for _, dhv := range dhvv {
				w.Header().Add(dhk, dhv)
			}
		}
		downbodyFile, err := h.OnResponse(e, upres.Body, upres.Header.Get("Content-Type"))
		if err != nil {
			h.Log(id, err.Error())
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
		downbody, err := os.Open(downbodyFile)
		if err != nil {
			h.Log(id, err.Error())
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
		w.WriteHeader(upres.StatusCode)
		_, err = io.Copy(w, downbody)
		if err != nil {
			h.Log(id, err.Error())
			return
		}
		err = downbody.Close()
		if err != nil {
			h.Log(id, err.Error())
		}
	}
}

func findMatching(overwrites []config.Overwrite, path string, method string) *config.Overwrite {
	for _, overwrite := range overwrites {
		if overwrite.Method != "" && overwrite.Method != method {
			continue
		}
		if overwrite.Exact != "" {
			if overwrite.Exact == path {
				return &overwrite
			}
		} else if overwrite.Regex != nil {
			if overwrite.Regex.MatchString(path) {
				return &overwrite
			}
		}
	}
	return nil
}
