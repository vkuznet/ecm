ECM=`go env GOROOT`
# flags=-ldflags="-s -w -extldflags -static -X main.RootCA=$(shell pwd)/rootCA.crt"
flags=-ldflags="-s -w -extldflags -static"

all: build

vet:
	go vet .

js:
	cp -f ${ECM}/misc/wasm/wasm_exec.js .

build: js
	GOOS=js GOARCH=wasm go build -o ecm.wasm ${flags}
	/bin/cp -f ecm.wasm wasm_exec.js extension/js

tinygo:
	GOROOT=/usr/local/go tinygo build -o ecm.wasm -target=wasm
