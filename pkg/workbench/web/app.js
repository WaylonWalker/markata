
const source = document.querySelector('#source');
const frame = document.querySelector('#preview');
const status = document.querySelector('#status');
const errorBox = document.querySelector('#error');
const initial = '# A place to think\n\nEdit this side. Markata renders the other.\n\n!!! note\n    Your draft stays in this browser.\n\n## A small example\n\n```go\nresult, err := plugins.RenderStandalone("# Hello", options)\n```\n\n## Follow the idea\n\n```mermaid\ngraph LR\n  Markdown --> Go --> Preview\n```\n\n- [x] Start writing\n- [ ] Make something useful\n';
try { source.value = localStorage.getItem('markata-workbench-draft') ?? initial; }
catch { source.value = initial; storageWarning(); }
const editor = CodeMirror.fromTextArea(source, {
  mode:'gfm', lineNumbers:true, lineWrapping:true, matchBrackets:true,
  indentUnit:2, tabSize:4, viewportMargin:10, screenReaderLabel:'Markdown source'
});
const keymap = document.querySelector('#keymap');
try { keymap.value = localStorage.getItem('markata-editor-keymap') || 'default'; } catch { storageWarning(); }
if (!keymap.value) keymap.value = 'default';
editor.setOption('keyMap', keymap.value);
keymap.addEventListener('change', () => {
  editor.setOption('keyMap', keymap.value);
  try { localStorage.setItem('markata-editor-keymap', keymap.value); } catch { storageWarning(); }
  editor.focus();
});
let revision = 0, timer, busy = false, ready = false;
function storageWarning() { document.querySelector('#storage').textContent = 'Browser storage unavailable; download your draft to keep it. '; }
async function render() {
  if (!ready || busy) return;
  const version = revision, markdown = source.value;
  const bytes = new TextEncoder().encode(markdown).length;
  if (bytes > 200000) {
    errorBox.textContent = 'Document exceeds 200,000 UTF-8 bytes. Download your draft or shorten it.';
    errorBox.hidden = false; status.textContent = 'Preview paused.'; return;
  }
  busy = true;
  status.textContent = 'Rendering…';
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 30000);
  const started = performance.now();
  try {
    const response = await fetch('/render', {
      method:'POST', headers:{'Content-Type':'application/json'},
      body:JSON.stringify({markdown, appearance:Object.fromEntries(Object.entries(attributes()).filter(([,v])=>v))}),
      signal:controller.signal
    });
    const data = await response.json();
    if (version !== revision) return;
    if (!response.ok) throw new Error(data.error || 'Rendering failed.');
    frame.onload = syncAppearance;
    frame.srcdoc = data.document;
    errorBox.hidden = true;
    status.textContent = `Go ${data.render_ms.toFixed(1)} ms${data.cached ? ' (cached)' : ''} · round trip ${Math.round(performance.now()-started)} ms · input ${(data.input_bytes/1024).toFixed(1)} KB · preview ${(data.document_bytes/1024).toFixed(1)} KB`;
  } catch (error) {
    if (version !== revision) return;
    errorBox.textContent = error.name === 'AbortError' ? 'Preview timed out. Edit the Markdown or reload to retry.' : error.message;
    errorBox.hidden = false;
    status.textContent = 'Preview failed. Edit the Markdown to try again.';
  } finally {
    clearTimeout(timeout); busy = false;
    if (version !== revision) { clearTimeout(timer); timer = setTimeout(render, 0); }
  }
}
editor.on('change', () => {
  editor.save(); revision++; clearTimeout(timer);
  try { localStorage.setItem('markata-workbench-draft', source.value); } catch { storageWarning(); }
  status.textContent = 'Waiting for edits…';
  timer = setTimeout(render, 250);
});
document.querySelector('#download').addEventListener('click', () => {
  const url = URL.createObjectURL(new Blob([source.value], {type:'text/markdown;charset=utf-8'}));
  const link = document.createElement('a');
  link.href = url; link.download = 'draft.md'; link.click();
  setTimeout(() => URL.revokeObjectURL(url), 1000);
});
const appearanceKeys = ['data-palette','data-theme','data-aesthetic','data-fontpack','data-text-size'];
function attributes() { return Object.fromEntries(appearanceKeys.map(k => [k, document.documentElement.getAttribute(k)])); }
function syncAppearance() { frame.contentWindow?.postMessage({type:'appearance',attributes:attributes()}, '*'); }
new MutationObserver(syncAppearance).observe(document.documentElement, {attributes:true,attributeFilter:appearanceKeys});
window.addEventListener('load', () => { ready = true; render(); }, {once:true});
