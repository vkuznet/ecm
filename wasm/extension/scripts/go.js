const go = new Go();
const wasmUrl = chrome.runtime.getURL("../scripts/decode.wasm");
WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject)
    .then(result => {
    go.run(result.instance);
});
