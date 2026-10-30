/* One persistent Go runtime. The UI owns debounce and coalesces pending edits. */
importScripts('wasm_exec.js');
(async () => {
  const go = new Go();
  const response = await fetch('markata.wasm');
  if (!response.ok) throw new Error(`Cannot load renderer (${response.status})`);
  // instantiateStreaming needs application/wasm; static hosts can omit it.
  let instance;
  try { ({instance} = await WebAssembly.instantiateStreaming(response.clone(), go.importObject)); }
  catch { ({instance} = await WebAssembly.instantiate(await response.arrayBuffer(), go.importObject)); }
  go.run(instance).catch(error => postMessage({type:'fatal', error:String(error)}));
  if (typeof markataRender !== 'function') throw new Error('Renderer did not initialize');
  onmessage = ({data}) => {
    const started = performance.now();
    try {
      const result = JSON.parse(markataRender(data.source, data.profile || 'web', Boolean(data.inline_assets)));
      postMessage({type:'result', revision:data.revision, result, render_ms:performance.now()-started});
    } catch (error) { postMessage({type:'result',revision:data.revision,result:{error:String(error)}}); }
  };
  postMessage({type:'ready'});
})().catch(error => postMessage({type:'fatal',error:String(error)}));
