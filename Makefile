.PHONY: build
build:
	GOOS=js GOARCH=wasm go build -o openapi.wasm

DIR=
.PHONY: wasm_copy
wasm_copy: build
	cp openapi.wasm $(DIR)
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js $(DIR)/js
	cp sw.js $(DIR)/js/wasm_sw.js
