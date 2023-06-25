const now = Date.now().toString();
importScripts(`js/wasm_exec.js?t=${now}`);
importScripts(`js/wasm_sw.js?t=${now}`);

self.started = false;

self.addEventListener('install', (event) => {
  event.waitUntil(skipWaiting());
});

self.addEventListener('activate', (event) => {
  event.waitUntil(self.clients.claim());
  console.log('activate');
  registerWasmHTTPListener('openapi.wasm').then((v) => {
    self.started = true;
  });
});

self.addEventListener('message', (e) => {
  switch (e.data) {
    case 'init':
      e.source.postMessage({ init: self.started });
  }
});
