/* Hosted editor WASM renderer. Dedicated worker; never handles authentication
 * or persistence. The shell is authenticated, but Markdown stays in this worker. */
importScripts('wasm_exec.js' + (new URL(location.href).searchParams.has('v')
  ? '?v=' + encodeURIComponent(new URL(location.href).searchParams.get('v')) : ''));

(async () => {
  const url = new URL(location.href);
  const suffix = url.searchParams.has('v') ? '?v=' + encodeURIComponent(url.searchParams.get('v')) : '';
  const go = new Go();
  const response = await fetch('markata.wasm' + suffix, {credentials: 'same-origin'});
  if (!response.ok) throw new Error('Cannot load editor WASM (' + response.status + ')');
  let instance;
  try {
    ({instance} = await WebAssembly.instantiateStreaming(response.clone(), go.importObject));
  } catch {
    ({instance} = await WebAssembly.instantiate(await response.arrayBuffer(), go.importObject));
  }
  go.run(instance).catch(error => postMessage({type: 'fatal', error: String(error)}));
  if (typeof markataRender !== 'function') throw new Error('WASM renderer unavailable');
  onmessage = ({data}) => {
    if (data?.type !== 'render' || !Number.isSafeInteger(data.revision) ||
        typeof data.source !== 'string') return;
    const started = performance.now();
    try {
      const result = JSON.parse(markataRender(data.source));
      postMessage({type: 'result', revision: data.revision, result,
        render_ms: performance.now() - started});
    } catch (error) {
      postMessage({type: 'result', revision: data.revision,
        result: {error: String(error)}, render_ms: performance.now() - started});
    }
  };
  postMessage({type: 'ready', schema_version: 1, version: url.searchParams.get('v') || ''});
})().catch(error => postMessage({type: 'fatal', error: String(error)}));
