### WabAssembly plugin
This area contains wasm (WebAssembly) plugin for ECM. This plugin can be used
in web browser extensions.

#### Build procedue
Please use the following instructions:
```
# build wasm executable
make

# copy exec_wasm.js file
make js

# copy wasm files to extension folder
cp ecm.wasm exec_wasm.js extension
```
Now, you can load your extension into a browser. The following steps
are tested with Chrome based browser(s) (Chrome, Brave, etc.):
- open Extention menu of your browser
- click on *Load unpacked*
- navigate to your extension folder and load it (the folder name)
- the new extention will appear in your browser

In order to test wasm code you may user provided server. Please follow
these steps:
- build your `ecm.wasm` and `exec_wasm.js` files
- copy these files to server folder, build and run the server code, e.g.
```
# copy necessary files to the stand-alone server for testing purposes
cp ecm.wasm wasm_exec.js server

# run go server to use wasm code
cd server
go run server.go

# visit localhost:9090 to see it in action
```

### References
[Go in the browser using WebAssembly](https://dev.bitolog.com/go-in-the-browser-using-webassembly/)
and
[Go WebAssembly Wiki](https://zchee.github.io/golang-wiki/WebAssembly/)
and
[Go WebAssembly handling HTTP requests](https://withblue.ink/2020/10/03/go-webassembly-http-requests-and-promises.html)
