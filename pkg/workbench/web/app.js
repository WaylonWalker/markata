
const source = document.querySelector('#source');
const frame = document.querySelector('#preview');
window.addEventListener('message',event=>{
  if (event.source !== frame.contentWindow) return;
  if (event.data?.type === 'open-post') window.MarkataSiteTree?.openPost(event.data.slug);
  if (event.data?.type === 'preview-updated' && event.data.version === revision) {
    finishPreview(event.data.version);
    scrollSync.previewUpdated();
  }
});
const status = document.querySelector('#status');
const errorBox = document.querySelector('#error');
const saved = document.querySelector('#saved');
const initial = '# A place to think\n\nEdit this side. Markata renders the other.\n\n!!! note\n    Your draft stays in this browser.\n\n## A small example\n\n```go\nresult, err := plugins.RenderStandalone("# Hello", options)\n```\n\n## Follow the idea\n\n```mermaid\ngraph LR\n  Markdown --> Go --> Preview\n```\n\n- [x] Start writing\n- [ ] Make something useful\n';
try { source.value = localStorage.getItem('markata-workbench-draft') ?? initial; }
catch { source.value = initial; storageWarning(); }
const editor = CodeMirror.fromTextArea(source, {
  mode:'gfm', lineNumbers:true, lineWrapping:true, matchBrackets:true,
  indentUnit:2, tabSize:4, viewportMargin:10, screenReaderLabel:'Markdown source'
});
const settingsDialog = document.querySelector('#settings-dialog');
document.querySelector('#settings-open').addEventListener('click',() => settingsDialog.showModal());
document.querySelector('#settings-close').addEventListener('click',() => settingsDialog.close());
const wrapSetting = document.querySelector('#setting-wrap');
try { wrapSetting.checked = localStorage.getItem('markata-editor-wrap') !== 'false'; } catch { storageWarning(); }
editor.setOption('lineWrapping',wrapSetting.checked);
wrapSetting.addEventListener('change',() => {editor.setOption('lineWrapping',wrapSetting.checked);try {localStorage.setItem('markata-editor-wrap',String(wrapSetting.checked));} catch {storageWarning();}});
const workbench = document.querySelector('#workbench');
const layoutSetting = document.querySelector('#setting-layout');
try { layoutSetting.value = localStorage.getItem('markata-editor-layout') || 'split'; } catch { storageWarning(); }
function setLayout(layout) {
  workbench.dataset.layout = ['split','editor','preview'].includes(layout) ? layout : 'split';
  layoutSetting.value = workbench.dataset.layout;
  try { localStorage.setItem('markata-editor-layout',workbench.dataset.layout); } catch { storageWarning(); }
  requestAnimationFrame(() => editor.refresh());
}
setLayout(layoutSetting.value);
const scrollSync = window.MarkataScrollSync(editor, frame, document.querySelector('#setting-scroll-sync'), workbench);
layoutSetting.addEventListener('change',() => setLayout(layoutSetting.value));
const paneDivider = document.querySelector('#pane-divider');
const editorPane = document.querySelector('#editor-pane');
function resizeEditorPane(clientX) {
  const bounds = workbench.getBoundingClientRect();
  const available = bounds.width - document.querySelector('#site-sidebar').getBoundingClientRect().width - paneDivider.offsetWidth;
  const width = Math.max(280,Math.min(available - 280,clientX - editorPane.getBoundingClientRect().left));
  workbench.style.setProperty('--editor-width',`${width}px`);
  paneDivider.setAttribute('aria-valuenow',String(Math.round(width / available * 100)));
  editor.refresh();
}
paneDivider.addEventListener('pointerdown',event => {if(workbench.dataset.layout !== 'split')return;paneDivider.setPointerCapture(event.pointerId);resizeEditorPane(event.clientX);});
paneDivider.addEventListener('pointermove',event => {if(paneDivider.hasPointerCapture(event.pointerId))resizeEditorPane(event.clientX);});
paneDivider.addEventListener('pointerup',event => {if(paneDivider.hasPointerCapture(event.pointerId))paneDivider.releasePointerCapture(event.pointerId);});
paneDivider.addEventListener('dblclick',() => {workbench.style.removeProperty('--editor-width');paneDivider.setAttribute('aria-valuenow','50');editor.refresh();});
paneDivider.addEventListener('keydown',event => {
  if (!['ArrowLeft','ArrowRight'].includes(event.key) || workbench.dataset.layout !== 'split') return;
  event.preventDefault(); resizeEditorPane(editorPane.getBoundingClientRect().right + (event.key === 'ArrowRight' ? 24 : -24));
});
const problemsPanel = document.querySelector('#problems-panel');
const problemsList = document.querySelector('#problems-list');
const problemToggle = document.querySelector('#problems-toggle');
function showProblems(diagnostics = []) {
  const errors = diagnostics.filter(item => item.severity === 1).length;
  const warnings = diagnostics.filter(item => item.severity === 2).length;
  document.querySelector('#error-count').textContent = errors;
  document.querySelector('#warning-count').textContent = warnings;
  problemToggle.setAttribute('aria-label', `${errors} errors and ${warnings} warnings`);
  problemsList.replaceChildren();
  for (const diagnostic of diagnostics) {
    const row = document.createElement('li');
    const range = diagnostic.range?.start || {line:0,character:0};
    const level = document.createElement('span');
    level.className = diagnostic.severity === 1 ? 'problem-error' : 'problem-warning';
    level.textContent = diagnostic.severity === 1 ? 'Error' : 'Warning';
    const message = document.createElement('span'); message.textContent = diagnostic.message;
    const location = document.createElement('span'); location.className = 'problem-location';
    location.textContent = `Ln ${range.line + 1}, Col ${range.character + 1}`;
    row.append(level, message, location);
    row.addEventListener('click', () => {problemsPanel.hidden = true; problemToggle.setAttribute('aria-expanded','false'); editor.focus(); editor.setCursor(range.line,range.character);});
    problemsList.append(row);
  }
}
window.MarkataWorkbenchDiagnostics = {update:showProblems};
const toggleProblems = () => {problemsPanel.hidden = !problemsPanel.hidden; problemToggle.setAttribute('aria-expanded',String(!problemsPanel.hidden));};
problemToggle.addEventListener('click',toggleProblems);
document.querySelector('#problems-close').addEventListener('click',toggleProblems);
const createDialog = document.querySelector('#create-dialog');
const createError = document.querySelector('#create-error');
const templatePicker = document.querySelector('#template-picker');
const templateQuery = document.querySelector('#template-picker-query');
const templateResults = document.querySelector('#template-picker-results');
const templateCount = document.querySelector('#template-picker-count');
const createTitle = document.querySelector('#create-title');
const createSlug = document.querySelector('#create-slug');
const createDirectory = document.querySelector('#create-directory');
const createPathPreview = document.querySelector('#create-path-preview');
let contentTemplates = [];
async function showCreateDialog() {
  try {
    const response = await fetch('/workspace/templates',{cache:'no-store'});
    const data = await response.json();
    if (!response.ok) throw new Error(data.error || 'Templates are unavailable.');
    contentTemplates = data.templates || [];
    templateQuery.value = '';
    renderTemplateChoices();
    templatePicker.showModal();
    templateQuery.focus();
  } catch (error) {
    createError.textContent = error.message; createError.hidden = false;
    createDialog.showModal();
  }
}
let templateMatches = [], activeTemplate = 0;
function templateScore(template, query) {
  const haystack = `${template.name} ${template.description} ${template.id}`.toLowerCase();
  if (!query) return 0;
  const exact = haystack.indexOf(query);
  if (exact >= 0) return exact;
  let cursor = 0, gaps = 0;
  for (const char of query) {
    const position = haystack.indexOf(char, cursor);
    if (position < 0) return Infinity;
    gaps += position - cursor; cursor = position + 1;
  }
  return 100 + gaps;
}
function chooseTemplate(template) {
  if (!template) return;
  templatePicker.dataset.selectedTemplate = template.id;
  templatePicker.close();
  createError.hidden = true;
  document.querySelector('#create-description').textContent = template.description || template.name;
  const titleMatch = template.source.match(/^title:\s*(.*)$/m);
  const title = titleMatch?.[1]?.trim().replace(/^(?:"([\s\S]*)"|'([\s\S]*)')$/, '$1$2') || template.name;
  createTitle.value = title;
  createSlug.value = '';
  createSlug.placeholder = slugify(title) || 'post-title';
  createDirectory.value = defaultContentDirectory(template);
  updateCreatePathPreview();
  createDialog.showModal();
  createTitle.focus();
  createTitle.select();
}
function renderTemplateChoices() {
  const query = templateQuery.value.trim().toLowerCase();
  templateMatches = contentTemplates.map(template => ({template,score:templateScore(template,query)}))
    .filter(item => Number.isFinite(item.score))
    .sort((a,b) => a.score-b.score || a.template.name.localeCompare(b.template.name))
    .slice(0,40).map(item => item.template);
  templateCount.textContent = templateMatches.length ? `${templateMatches.length} template${templateMatches.length===1?'':'s'}` : 'No matching templates';
  templateResults.replaceChildren(...templateMatches.map((template,index) => {
    const button = document.createElement('button');
    button.type = 'button'; button.className = 'template-choice'; button.setAttribute('role','option');
    button.setAttribute('aria-selected',String(index===activeTemplate)); button.id = `template-choice-${index}`;
    const name = document.createElement('strong'); name.textContent = template.name;
    const description = document.createElement('span'); description.textContent = template.description || template.id;
    button.append(name,description);
    button.addEventListener('pointermove',() => selectTemplateChoice(index));
    button.addEventListener('click',() => chooseTemplate(template));
    return button;
  }));
  activeTemplate = Math.min(activeTemplate,Math.max(0,templateMatches.length-1));
  selectTemplateChoice(activeTemplate);
}
function selectTemplateChoice(index) {
  activeTemplate = index;
  [...templateResults.children].forEach((item,i) => item.setAttribute('aria-selected',String(i===index)));
  const row = templateResults.children[index];
  if (row) {templateQuery.setAttribute('aria-activedescendant',row.id);row.scrollIntoView({block:'nearest'});}
  else templateQuery.removeAttribute('aria-activedescendant');
}
function slugify(value) {
  return value.normalize('NFKD').replace(/[\u0300-\u036f]/g,'').toLowerCase()
    .replace(/[^a-z0-9]+/g,'-').replace(/^-+|-+$/g,'');
}
function defaultContentDirectory(template) {
  const key = template.id.replace(/^builtin:/,'');
  return ({thought:'pages/thoughts',note:'pages/notes',photo:'pages/photos'})[key] || 'pages/blog';
}
function updateCreatePathPreview() {
  const slug = createSlug.value.trim() || slugify(createTitle.value) || 'new-post';
  const directory = createDirectory.value.trim().replace(/\\/g,'/').replace(/^\/+|\/+$/g,'');
  createSlug.placeholder = slugify(createTitle.value) || 'post-title';
  createPathPreview.textContent = `${directory ? `${directory}/` : ''}${slug}.md`;
}
templateQuery.addEventListener('input',() => {activeTemplate=0;renderTemplateChoices();});
templatePicker.querySelector('#template-picker-close').addEventListener('click',() => templatePicker.close());
templatePicker.addEventListener('keydown',event => {
  if (event.isComposing) return;
  if (event.key === 'Escape') {event.preventDefault();event.stopPropagation();templatePicker.close();return;}
  if (!['ArrowDown','ArrowUp','Enter'].includes(event.key)) return;
  event.preventDefault();
  if (event.key === 'Enter') chooseTemplate(templateMatches[activeTemplate]);
  else if (templateMatches.length) selectTemplateChoice((activeTemplate+(event.key==='ArrowDown'?1:-1)+templateMatches.length)%templateMatches.length);
});
templatePicker.addEventListener('close',() => {if(!createDialog.open)editor.focus();});
createDialog.addEventListener('close',() => {if(!templatePicker.open)editor.focus();});
document.querySelector('#create-back').addEventListener('click',() => {createDialog.close();templateQuery.value='';renderTemplateChoices();templatePicker.showModal();templateQuery.focus();});
document.querySelector('#create-close').addEventListener('click',() => createDialog.close());
document.querySelector('#create-cancel').addEventListener('click',() => createDialog.close());
createTitle.addEventListener('input',updateCreatePathPreview);
createSlug.addEventListener('input',updateCreatePathPreview);
createDirectory.addEventListener('input',updateCreatePathPreview);
document.querySelector('#create-form').addEventListener('submit',async event => {
  event.preventDefault(); createError.hidden = true;
  const template = contentTemplates.find(item => item.id === templatePicker.dataset.selectedTemplate);
  const title = createTitle.value.trim();
  const slug = createSlug.value.trim() || slugify(title);
  const directory = createDirectory.value.trim().replace(/\\/g,'/').replace(/^\/+|\/+$/g,'');
  const path = `${directory ? `${directory}/` : ''}${slug}.md`;
  if (!template || !title || !slug || !path) return;
  const source = template.source.replace(/^title:\s*.*$/m,`title: ${JSON.stringify(title)}`)
    .replace(/^(#{1,6})\s+.*$/m,(heading,level) => `${level} ${title}`);
  try {
    const response = await fetch('/workspace/file?path=' + encodeURIComponent(path),{method:'POST',headers:{'Content-Type':'text/plain;charset=utf-8'},body:source});
    const data = await response.json();
    if (!response.ok) throw new Error(data.error || 'Cannot create this file.');
    window.MarkataSiteTree?.activate(path);
    createDialog.close(); location.reload();
  } catch (error) { createError.textContent = error.message; createError.hidden = false; }
});
window.MarkataCommandPalette.attach(editor,{actions:[
  {id:'new-from-template',label:'New Post',description:'Choose a template, then set its title and path.',aliases:'new create post page template',run:showCreateDialog},
  {id:'open-markdown',label:'Open Markdown file',description:'Import a local Markdown file into an unsaved draft.',aliases:'file import',run:()=>document.querySelector('#open').click()},
  {id:'download-markdown',label:'Download current file',description:'Save the current source as a Markdown file.',aliases:'export save',run:()=>document.querySelector('#download').click()}
]});
window.MarkataEditorLSP.connect(editor, document.querySelector('#language-status'));
if (new URLSearchParams(location.search).has('example')) {
  fetch(document.body.dataset.exampleUrl).then(response => {
    if (!response.ok) throw new Error('Cannot load the command palette example.');
    return response.text();
  }).then(async markdown => {
    await workspaceReady;
    editor.setValue(markdown);
    editor.focus();
    setTimeout(() => document.querySelector('#commands').click(), 250);
  }).catch(error => { errorBox.textContent = error.message; errorBox.hidden = false; });
}
const keymap = document.querySelector('#keymap');
try { keymap.value = localStorage.getItem('markata-editor-keymap') || 'default'; } catch { storageWarning(); }
if (!keymap.value) keymap.value = 'default';
editor.setOption('keyMap', keymap.value);
function focusEditorForWriting() {
  editor.focus();
  if (keymap.value === 'vim' && window.CodeMirror?.Vim && !editor.state.vim?.insertMode) {
    window.CodeMirror.Vim.handleKey(editor,'i');
  }
}
keymap.addEventListener('change', () => {
  editor.setOption('keyMap', keymap.value);
  try { localStorage.setItem('markata-editor-keymap', keymap.value); } catch { storageWarning(); }
  focusEditorForWriting();
});
let revision = 0, timer, busy = false, ready = false, lastEditAt = null;
let previewStarted = false, previewReady = false, pendingPreview = null, activeRender = null;
function finishPreview(version) {
  if (!activeRender || activeRender.version !== version || version !== revision) return;
  const {data} = activeRender;
  const liveTime = Math.round(performance.now() - (activeRender.editAt ?? activeRender.startedAt));
  status.textContent = `Go ${data.render_ms.toFixed(1)} ms${data.cached ? ' (cached)' : ''} · preview ${liveTime} ms · input ${(data.input_bytes/1024).toFixed(1)} KB · fragment ${((data.html_bytes)/1024).toFixed(1)} KB`;
}
function storageWarning() { saved.textContent = 'Draft not saved · storage unavailable or full · download a copy'; }
function saveDraft() {
 if (window.MarkataSiteTree?.save()) return;
  try { localStorage.setItem('markata-workbench-draft', source.value); saved.textContent = 'Draft saved on this device'; } catch { storageWarning(); }
}
function updateFrontmatterSummary(markdown) {
  const panel = document.querySelector('#preview-meta');
  const match = markdown.match(/^---\r?\n([\s\S]*?)\r?\n---(?:\r?\n|$)/);
  if (!match) {panel.hidden = true;return;}
  const lines = match[1].split(/\r?\n/), fields = new Map();
  for (let i=0;i<lines.length;i++) {
    const field = lines[i].match(/^([A-Za-z][\w-]*):\s*(.*)$/);
    if (!field) continue;
    const key = field[1], value = field[2].trim();
    if (key === 'tags' && !value) {
      const tags=[];
      for (let j=i+1;j<lines.length && /^\s+-\s+/.test(lines[j]);j++) tags.push(lines[j].replace(/^\s+-\s+/,''));
      fields.set(key,tags.join(', '));
    } else fields.set(key,value.replace(/^(?:"([^"]*)"|'([^']*)')$/,'$1$2'));
  }
  const title = fields.get('title') || fields.get('templateKey') || 'Untitled document';
  const details = [];
  if (fields.has('date')) details.push(fields.get('date'));
  if (fields.get('published') === 'true') details.push('Published');
  else if (fields.get('draft') === 'true' || fields.get('published') === 'false') details.push('Draft');
  if (fields.has('tags')) details.push(`Tags: ${fields.get('tags')}`);
  document.querySelector('#preview-title').textContent = title;
  document.querySelector('#preview-details').textContent = details.join(' · ');
  panel.hidden = false;
}
async function render() {
  if (!ready || busy) return;
  const version = revision, originalMarkdown = source.value;
  const markdown = window.MarkataSiteTree?.markdown(originalMarkdown) ?? originalMarkdown;
  const frontmatter = originalMarkdown.match(/^---\r?\n[\s\S]*?\r?\n---(?:\r?\n|$)/);
  const sourceLineOffset = frontmatter && markdown !== originalMarkdown
    ? (frontmatter[0].match(/\n/g) || []).length
    : 0;
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
      body:JSON.stringify({markdown, include_document:!previewStarted, appearance:Object.fromEntries(Object.entries(attributes()).filter(([,v])=>v))}),
      signal:controller.signal
    });
    const data = await response.json();
    if (version !== revision) return;
    if (!response.ok) throw new Error(data.error || 'Rendering failed.');
    if (sourceLineOffset) data.source_lines = data.source_lines.map(line => line + sourceLineOffset);
    scrollSync.update(data.source_lines, false);
    activeRender = {data, version, editAt:lastEditAt, startedAt:started};
    if (!previewStarted) {
      previewStarted = true;
      frame.onload = () => {
        previewReady = true;
        syncAppearance();
        scrollSync.ready();
        if (pendingPreview) {
          const latest = pendingPreview;
          pendingPreview = null;
          if (latest.version === revision) updatePreview(latest.data, latest.version);
        } else {
          scrollSync.previewUpdated();
          finishPreview(activeRender?.version);
        }
      };
      frame.srcdoc = data.document;
    } else if (!previewReady) {
      pendingPreview = {data, version};
    } else {
      updatePreview(data, version);
    }
    errorBox.hidden = true;
    status.textContent = 'Updating preview…';
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
editor.on('change', (_instance, change) => {
  editor.save(); revision++; clearTimeout(timer);
  lastEditAt = performance.now();
  updateFrontmatterSummary(editor.getValue());
  if (change.origin !== 'setValue') saveDraft();
  status.textContent = 'Waiting for edits…';
  timer = setTimeout(render, 100);
});
document.querySelector('#download').addEventListener('click', () => {
  const url = URL.createObjectURL(new Blob([source.value], {type:'text/markdown;charset=utf-8'}));
  const link = document.createElement('a');
  link.href = url; link.download = window.MarkataSiteTree?.filename() || 'draft.md'; link.click();
  setTimeout(() => URL.revokeObjectURL(url), 1000);
});
const appearanceKeys = ['data-palette','data-theme','data-aesthetic','data-fontpack','data-text-size'];
function updatePreview(data, version) {
  frame.contentWindow?.postMessage({type:'markdown-render', html:data.html,
    source_lines:data.source_lines,
    line_count:editor.lineCount(),
    feature_head_html:data.feature_head_html, body_end_html:data.body_end_html, version,
    attributes:attributes()}, '*');
}
function attributes() { return Object.fromEntries(appearanceKeys.map(k => [k, document.documentElement.getAttribute(k)])); }
function syncAppearance() { frame.contentWindow?.postMessage({type:'appearance',attributes:attributes()}, '*'); }
new MutationObserver(syncAppearance).observe(document.documentElement, {attributes:true,attributeFilter:appearanceKeys});
window.addEventListener('load', async () => { await workspaceReady; ready = true; render(); focusEditorForWriting(); }, {once:true});

document.querySelector('#open').addEventListener('click', () => document.querySelector('#import').click());
document.querySelector('#import').addEventListener('change', async event => {
  const file = event.target.files[0];
  event.target.value = '';
  if (!file) return;
  if (file.size > 200000) {
    errorBox.textContent = 'Document exceeds 200,000 UTF-8 bytes. Your current draft has been kept.';
    errorBox.hidden = false; return;
  }
  try {
    const text = await file.text();
    if (new TextEncoder().encode(text).length > 200000) throw new Error('Document exceeds 200,000 UTF-8 bytes.');
    window.MarkataSiteTree.import(text);
  } catch (error) { errorBox.textContent = 'Cannot open file: ' + error.message; errorBox.hidden = false; }
});
if (!saved.textContent) saveDraft();

const workspaceReady = window.MarkataSiteTree.attach(editor, storageWarning);
workspaceReady.then(() => window.MarkataFileFinder?.attach(editor));
