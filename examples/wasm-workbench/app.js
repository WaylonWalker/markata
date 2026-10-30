/* The only draft copy lives in browser storage or the user's exported file. */
(async () => {
  const source = document.querySelector('#source');
  const frame = document.querySelector('#preview');
  const status = document.querySelector('#status');
  const saved = document.querySelector('#saved');
  const errorBox = document.querySelector('#error');
  const offline = document.querySelector('#offline');
  const renderProfile = document.querySelector('#render-profile');
  const renderAssets = document.querySelector('#render-assets');
  const draftKey = 'markata-wasm-draft-v1';
  const initial = '# A place to think\n\nMarkata, running entirely in your browser.\n\n!!! note\n    Your draft stays on this device. Export a copy to keep it elsewhere.\n\n## Follow the idea\n\n```mermaid\ngraph LR\n  Markdown --> WASM --> Preview\n```\n\n## Small things, clearly\n\n```go\nmessage := "Hello from WebAssembly"\n```\n\n- [x] Keep rendering local\n- [ ] Make something useful\n';
  function showError(message) { errorBox.textContent = message; errorBox.hidden = false; }
  try { source.value = localStorage.getItem(draftKey) ?? localStorage.getItem('markata-workbench-draft') ?? initial; }
  catch { source.value = initial; saved.textContent = 'Storage unavailable · export to keep your draft'; }
  const editor = CodeMirror.fromTextArea(source, {mode:'gfm',lineNumbers:true,lineWrapping:true,matchBrackets:true,indentUnit:2,tabSize:4,viewportMargin:10,screenReaderLabel:'Markdown source'});
  // Exposed on the element by CodeMirror for host integrations and browser tests.
  window.MarkataCommandPalette.attach(editor);
if (new URLSearchParams(location.search).has('example')) {
  fetch(document.body.dataset.exampleUrl).then(response => {
    if (!response.ok) throw new Error('Cannot load the command palette example.');
    return response.text();
  }).then(markdown => {
    editor.setValue(markdown);
    editor.focus();
    setTimeout(() => document.querySelector('#commands').click(), 250);
  }).catch(error => { errorBox.textContent = error.message; errorBox.hidden = false; });
}
  const keymap = document.querySelector('#keymap');
  try { keymap.value = localStorage.getItem('markata-editor-keymap') || 'default'; } catch {}
  if (!keymap.value) keymap.value = 'default';
  editor.setOption('keyMap',keymap.value);
  keymap.addEventListener('change',() => {
    editor.setOption('keyMap',keymap.value);
    try { localStorage.setItem('markata-editor-keymap',keymap.value); } catch {}
    editor.focus();
  });
  let revision = 0, busy = false, ready = false, pending = false, timer, worker;
  let previewReady = false, latestResult;
  const transferred = new Set();
  const attributes = () => Object.fromEntries(['data-theme','data-palette','data-aesthetic','data-fontpack','data-text-size'].map(key => [key,document.documentElement.getAttribute(key)]));
  const syncAppearance = () => frame.contentWindow?.postMessage({type:'appearance',attributes:attributes()},'*');
  new MutationObserver(syncAppearance).observe(document.documentElement,{attributes:true,attributeFilter:Object.keys(attributes())});
  const sendPreview = () => {
    if (!previewReady || !latestResult) return;
    syncAppearance();
    const result={...latestResult,assets:latestResult.assets.map(asset =>
      !asset.id.startsWith('generated-') && transferred.has(asset.id) ? {...asset,content:''} : asset)};
    frame.contentWindow.postMessage({type:'render',revision,result},'*');
    for (const asset of result.assets) if (!asset.id.startsWith('generated-')) transferred.add(asset.id);
  };
  window.addEventListener('message',e => {
    if (e.source !== frame.contentWindow) return;
    if (e.data?.type === 'preview-ready') { transferred.clear();previewReady = true; sendPreview(); }
    if (e.data?.type === 'preview-error') showError(e.data.error);
  });
  const assets = new Map();
  async function assetText(id) {
    if (!assets.has(id)) assets.set(id,fetch('appearance/'+id).then(response => {
      if (!response.ok) throw new Error('Cannot load '+id);
      return response.text();
    }).catch(error => {assets.delete(id);throw error;}));
    return assets.get(id);
  }
  async function preparePreview() {
    const response=await fetch('preview.html');
    if (!response.ok) throw new Error('Cannot load preview shell');
    const shell=new DOMParser().parseFromString(await response.text(),'text/html');
    for (const link of shell.querySelectorAll('link[rel=stylesheet]')) {
      const id=link.getAttribute('href').replace(/^appearance\//,'');
      const style=shell.createElement('style');
      style.textContent=await assetText(id==='css/fonts.css' ? 'css/preview-fonts.css' : id);
      link.replaceWith(style);
    }
    const bridge=shell.querySelector('script[src="preview.js"]');
    const script=await fetch('preview.js');
    if (!script.ok) throw new Error('Cannot load preview bridge');
    bridge.removeAttribute('src');bridge.removeAttribute('defer');
    bridge.textContent=await script.text();
    frame.srcdoc='<!doctype html>'+shell.documentElement.outerHTML;
  }
  const previewPromise=preparePreview().catch(error => showError(String(error)));
  async function populateAssets(result) {
    return {...result,assets:await Promise.all(result.assets.map(async asset => {
      if (asset.id.startsWith('generated-') || asset.content) return asset;
      return {...asset,content:await assetText(asset.id)};
    }))};
  }
  function schedule() {
    if (!ready || busy || !pending) return;
    pending = false;
    const bytes = new TextEncoder().encode(source.value).length;
    if (bytes > 200000) { showError('Document exceeds 200,000 UTF-8 bytes. Export it before shortening it.'); status.textContent='Preview paused · document too large'; return; }
    busy = true;
    status.textContent = 'Rendering on this device…';
    worker.postMessage({source:source.value,revision,profile:renderProfile.value,inline_assets:renderAssets.value==='inline'});
  }
  function save() {
    try { localStorage.setItem(draftKey,source.value); saved.textContent='Draft saved on this device'; }
    catch { saved.textContent='Draft not saved · storage unavailable or full · export a copy'; }
  }
  const rendererSettingsChanged = () => {
    if (renderProfile.value === 'offline-interactive') {
      renderAssets.value = 'inline';
      renderAssets.disabled = true;
    } else {
      renderAssets.disabled = false;
    }
    revision++;
    pending = true;
    clearTimeout(timer);
    timer = setTimeout(schedule, 0);
  };
  renderProfile.addEventListener('change', rendererSettingsChanged);
  renderAssets.addEventListener('change', rendererSettingsChanged);
  rendererSettingsChanged();

  editor.on('change',() => {
    editor.save(); revision++; pending = true;
    clearTimeout(timer);
    // Save first, including when the renderer is still loading or has failed.
    save();
    status.textContent = 'Waiting for edits…';
    timer = setTimeout(schedule,250);
  });
  document.querySelector('#download').addEventListener('click',() => {
    const url = URL.createObjectURL(new Blob([source.value],{type:'text/markdown;charset=utf-8'}));
    const link = document.createElement('a'); link.href=url; link.download='draft.md'; link.click();
    setTimeout(() => URL.revokeObjectURL(url),1000);
  });
  document.querySelector('#open').addEventListener('click',() => document.querySelector('#import').click());
  document.querySelector('#import').addEventListener('change',async e => {
    const file = e.target.files[0];
    if (!file) return;
    e.target.value='';
    if (file.size > 200000) { showError('Document exceeds 200,000 UTF-8 bytes. Your current draft has been kept.'); return; }
    try {
      const text = await file.text();
      if (new TextEncoder().encode(text).length > 200000) throw new Error('Document exceeds 200,000 UTF-8 bytes.');
      editor.setValue(text);
    } catch (error) { showError(`Cannot read file: ${error.message}`); }
    e.target.value='';
  });
  try {
    worker = new Worker('worker.js');
    const fatal = message => { ready=false; busy=false; showError(message); status.textContent='Renderer unavailable · your draft can still be saved and exported'; };
    worker.onerror = e => fatal(e.message || 'Worker failed. Reload to retry.');
    worker.onmessage = async ({data}) => {
      if (data.type === 'fatal') { fatal(data.error); return; }
      if (data.type === 'ready') { ready=true;pending=true;schedule();return; }
      if (data.type !== 'result') return;
      if (data.revision === revision) {
        if (data.result.error) { showError(data.result.error);status.textContent='Preview failed · edit to retry'; }
        else {
          errorBox.hidden=true;
          try {
            const populated=await populateAssets(data.result);
            await previewPromise;
            if (data.revision === revision) {
              latestResult=populated;sendPreview();
              status.textContent=`WASM ${data.render_ms.toFixed(1)} ms · ${(new TextEncoder().encode(source.value).length/1024).toFixed(1)} KB · rendered locally`;
            }
          } catch (error) {showError(String(error));}

        }
      }
      busy=false;
      // If edits arrived during rendering, send only the newest revision.
      if (pending) { clearTimeout(timer); timer=setTimeout(schedule,0); }
    };
  } catch (error) { showError(error.message); }
  if (!saved.textContent) save();
  if ('serviceWorker' in navigator) {
    try {
      await navigator.serviceWorker.register('sw.js');
      await navigator.serviceWorker.ready;
      offline.textContent='Ready offline';
    } catch { offline.textContent='Offline copy unavailable · keep this tab open'; }
  } else offline.textContent='Offline reload unavailable in this browser';
})();
