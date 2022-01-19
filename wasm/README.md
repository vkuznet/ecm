### WabAssembly plugin
This area contains wasm (WebAssembly) plugin for ECM.
Please use the following instructions:
```
# build wasm executable
make
# copy exec_wasm.js file
make js
# adjust and copy index.html to use in your browser/server
# for instance, we'll copy index.html to server arear and
# run Go server from there
cp decode.wasm index.html wasm_exec.js server

# run go server to use wasm code
cd server
go run server.go

# visit localhost:9090 to see it in action
```

The idea is borrowed from:
[Go in the browser using WebAssembly](https://dev.bitolog.com/go-in-the-browser-using-webassembly/)
and
[Go WebAssembly Wiki](https://zchee.github.io/golang-wiki/WebAssembly/)
and
[Go WebAssembly handling HTTP requests](https://withblue.ink/2020/10/03/go-webassembly-http-requests-and-promises.html)
