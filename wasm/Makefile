ECM=`go env GOROOT`
flags=-ldflags="-s -w -extldflags -static"

all: build

vet:
	go vet .

js:
	cp -f ${ECM}/misc/wasm/wasm_exec.js .

build:
	GOOS=js GOARCH=wasm go build -o ecm.wasm ${flags}

tinygo:
	GOROOT=/usr/local/go tinygo build -o ecm.wasm -target=wasm
