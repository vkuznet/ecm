We may use
 "content_security_policy": "script-src 'self' 'unsafe-eval'; object-src 'self'"
but we should restrict it to sha-256 scheme:
 "content_security_policy": "script-src 'self' 'sha256-XXXX'; object-src 'self'"
where sha sum should be calculated as following:
shasum -a 256 -b <file.wasm>
