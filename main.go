package main

import (
	"net/http"

	"github.com/n-creativesystem/short-url/pkg/cmd/wasm/lib"
)

func main() {
	c := make(chan struct{})
	mux := http.NewServeMux()
	mux.HandleFunc("/swagger/api", func(w http.ResponseWriter, r *http.Request) {
		buf, err := openAPI.ReadFile("openapi/swagger.yaml")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(buf)
		}
	})
	r := lib.Serve(mux)
	defer r()
	<-c
}
