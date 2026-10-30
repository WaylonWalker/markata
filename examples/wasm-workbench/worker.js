/* One persistent Go runtime. The UI owns debounce and coalesces pending edits. */
importScripts('wasm_exec.js');
// Keep the Go/WASM binary lean: hydrate optional inline results from the
// same-origin bundle already exported and cached by the workbench host.
const assetCache = new Map();
async function localAssetText(id) {
  if (!/^(css|js|assets\/vendor)\/[a-zA-Z0-9_./-]+$/.test(id) || id.includes('..')) {
    throw new Error('Invalid managed asset ID');
  }
  if (!assetCache.has(id)) {
    assetCache.set(id, fetch('appearance/' + id).then(response => {
      if (!response.ok) throw new Error('Cannot load local asset ' + id);
      return response.text();
    }).catch(error => {assetCache.delete(id);throw error;}));
  }
  return assetCache.get(id);
}
function inlineTag(asset) {
  const content = asset.content || '';
  if (asset.kind === 'style') return '<style>' + content.replace(/<\/style/gi, '<\\/style') + '</style>';
  return '<script' + (asset.module ? ' type="module"' : '') + '>' +
    content.replace(/<\/script/gi, '<\\/script') + '</script>';
}
async function inlineManagedAssets(result, offline) {
  for (const asset of result.assets) {
    if (!asset.id.startsWith('generated-') && !asset.content) {
      asset.content = await localAssetText(asset.id);
    }
    // Source CSS is controlled by the static host, not by the Go WASM embed.
    if (offline && asset.kind === 'style' &&
        /(?:url\s*\(\s*['"]?\s*|@import\s+(?:url\s*\(\s*)?['"]?\s*)(?:https?:)?\/\//i.test(asset.content || '')) {
      throw new Error('Offline interactive asset has a remote CSS dependency: ' + asset.id);
    }
  }
  result.head_html = result.assets.filter(asset => asset.placement === 'head').map(inlineTag).join('\n') + '\n';
  result.body_end_html = result.assets.filter(asset => asset.placement !== 'head').map(inlineTag).join('\n') + '\n';
  return result;
}
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
  onmessage = async ({data}) => {
    const started = performance.now();
    try {
      const profile = data.profile || 'web';
      const result = JSON.parse(markataRender(data.source, profile, Boolean(data.inline_assets)));
      if (!result.error && (data.inline_assets || profile === 'offline-interactive')) {
        await inlineManagedAssets(result, profile === 'offline-interactive');
      }
      postMessage({type:'result', revision:data.revision, result, render_ms:performance.now()-started});
    } catch (error) { postMessage({type:'result',revision:data.revision,result:{error:String(error)}}); }
  };
  postMessage({type:'ready'});
})().catch(error => postMessage({type:'fatal',error:String(error)}));
