const go = new Go()

self.addEventListener('fetch', (e) => {
  const { pathname } = new URL(e.request.url);
  if (!pathname.startsWith('/swagger/')) return;

  e.respondWith(Promise.resolve(GoHandler(e.request)));
});

const registerWasmHTTPListener = async (wasm, { base, args = [] } = {}) => {
  go.argv = [wasm, ...args]
  let res;
  try {
    if (typeof wasm === 'string') {
      const fetchRes = fetch(wasm);
      res = fetchRes.then((response) => response.arrayBuffer());
    } else {
      res = Promise.resolve(wasm);
    }
    const buf = await res;
    const module = await WebAssembly.compile(buf);
    const instance = await WebAssembly.instantiate(module, go.importObject);
    go.run(instance);
  } catch (error) {
    console.error(error);
  }

  return true
}
