ECM=`go env GOROOT`

all: build

vet:
	go vet .

js:
	cp -f ${ECM}/misc/wasm/wasm_exec.js .

build:
	GOOS=js GOARCH=wasm go build -o ecm.wasm main.go

tinygo:
	GOROOT=/usr/local/go tinygo build -o ecm.wasm -target=wasm main.go
