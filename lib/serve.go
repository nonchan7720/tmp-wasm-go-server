package lib

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"syscall/js"
)

func Serve(handler http.Handler) func() {
	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		var p, resolve, reject = Promise()
		go func() {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						reject(fmt.Sprintf("wasmhttp: panic: %+v\n", err))
					} else {
						reject(fmt.Sprintf("wasmhttp: panic: %v\n", r))
					}
				}
			}()
			w := NewResponseRecorder()
			handler.ServeHTTP(w, Request(args[0]))
			resolve(w.JSResponse())
		}()
		return p
	})

	js.Global().Set("GoHandler", cb)
	return cb.Release
}

func Request(r js.Value) *http.Request {
	jsBody := js.Global().Get("Uint8Array").New(Await(r.Call("arrayBuffer")))
	body := make([]byte, jsBody.Get("length").Int())
	js.CopyBytesToGo(body, jsBody)

	req := httptest.NewRequest(
		r.Get("method").String(),
		r.Get("url").String(),
		bytes.NewBuffer(body),
	)

	headersIt := r.Get("headers").Call("entries")
	for {
		e := headersIt.Call("next")
		if e.Get("done").Bool() {
			break
		}
		v := e.Get("value")
		req.Header.Set(v.Index(0).String(), v.Index(1).String())
	}

	return req
}

type ResponseRecorder struct {
	*httptest.ResponseRecorder
}

func NewResponseRecorder() ResponseRecorder {
	return ResponseRecorder{httptest.NewRecorder()}
}

func (rr ResponseRecorder) JSResponse() js.Value {
	var res = rr.Result()
	var body js.Value = js.Undefined()
	if res.ContentLength != 0 {
		var b, err = io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		body = js.Global().Get("Uint8Array").New(len(b))
		js.CopyBytesToJS(body, b)
	}

	var init = make(map[string]interface{}, 2)

	if res.StatusCode != 0 {
		init["status"] = res.StatusCode
	}
	header := rr.Header()
	if len(header) != 0 {
		var headers = make(map[string]interface{}, len(res.Header))
		for k := range res.Header {
			headers[k] = res.Header.Get(k)
		}
		init["headers"] = headers
	}

	return js.Global().Get("Response").New(body, init)
}
