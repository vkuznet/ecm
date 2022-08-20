for Chrome V3 manifest please see:
https://groups.google.com/a/chromium.org/g/chromium-extensions/c/sJiaTnFMLHQ/m/tJQ9AE9vBQAJ

We may use
 "content_security_policy": "script-src 'self' 'unsafe-eval'; object-src 'self'"
but we should restrict it to sha-256 scheme:
 "content_security_policy": "script-src 'self' 'sha256-XXXX'; object-src 'self'"
where sha sum should be calculated as following:
shasum -a 256 -b <file.wasm>

store user options
https://developer.chrome.com/docs/extensions/mv3/options/

CORS:
https://groups.google.com/g/golang-nuts/c/Kz-14zEJ0Bg
https://stackoverflow.com/questions/40985920/making-golang-gorilla-cors-handler-work
https://github.com/rs/cors

XMLHttpRequest and xhr
https://developer.chrome.com/docs/extensions/mv3/xhr/
https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest/Using_XMLHttpRequest
https://www.w3schools.com/xml/xml_http.asp

Chrome security policy
https://newbedev.com/chrome-extension-content-security-policy-executing-inline-code
https://github.com/dteare/wasm-csp

Chrome capture user input
https://spin.atomicobject.com/2017/08/18/chrome-extension-form-data/
https://developer.chrome.com/docs/extensions/reference/webRequest/
https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/webRequest/onBeforeRequest
