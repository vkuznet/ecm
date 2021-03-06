const go = new Go();
const wasmUrl = chrome.runtime.getURL("js/ecm.wasm");
WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject)
    .then(result => {
    go.run(result.instance);
});
