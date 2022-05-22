const go = new Go();
WebAssembly.instantiateStreaming(fetch("static/js/ecm.wasm"), go.importObject)
    .then(result => {
    go.run(result.instance);
});
