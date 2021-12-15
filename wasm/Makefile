GRT=`go env GOROOT`

all: build

vet:
	go vet .

js:
	cp -f ${GRT}/misc/wasm/wasm_exec.js .

build:
	GOOS=js GOARCH=wasm go build -o decode.wasm main.go
