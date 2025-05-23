.PHONY: build clean test

build: 
	GOOS=js GOARCH=wasm go build -o ./client/public/wasm/crdt_conn.wasm ./document/wasm/main.go

clean: 
	rm -rf ./client/public/wasm/crdt_conn.wasm

test: go test ./...