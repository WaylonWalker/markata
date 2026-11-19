
const source = document.querySelector('#source');
const frame = document.querySelector('#preview');
const sourceVisibleText = value => String(value||'').replace(/^\s{0,3}(?:#{1,6}\s+|>\s?|[-*+]\s+|\d+\.\s+)/,'').replace(/!\[([^\]]*)\]\([^)]+\)/g,'$1').replace(/\[([^\]]+)\]\([^)]+\)/g,'$1').replace(/\[([^\]]+)\]\[[^\]]*\]/g,'$1').replace(/<[^>]*>/g,'').replace(/[`*_~]/g,'').replace(/[‘’]/g,"'").replace(/[“”]/g,'"').replace(/[–—]/g,'-').replace(/\s+/g,' ').trim();
const placeEditorCursor = position => {
  const lines = editor.getValue().split(/\r?\n/);
  const start = Math.max(0,Math.min(lines.length-1,Number(position.line)||0));
  const end = Math.max(start,Math.min(lines.length-1,Number(position.endLine)||start));
  const clicked = String(position.text||'').replace(/[‘’]/g,"'").replace(/[“”]/g,'"').replace(/[–—]/g,'-');
  const findSourceMatch = (from,to) => {
    let match=null;
    for(let i=from;i<=to;i++){
      const parts=[];
      for(let j=i;j<=to&&j<i+16;j++){
        if(!lines[j].trim())break;
        const visible=sourceVisibleText(lines[j]);
        if(visible)parts.push({line:j,text:visible,raw:lines[j]});
        const joined=parts.map(part=>part.text).join(' ');
        const index=joined?clicked.indexOf(joined):-1;
        if(index>=0&&(!match||joined.length>match.length))match={parts,index,length:joined.length};
      }
    }
    return match;
  };
  const best=findSourceMatch(start,end)||findSourceMatch(0,lines.length-1);
  let line=start,ch=0;
  if(best){
    const within=Math.max(0,Math.min(best.length,Number(position.offset)||0)-best.index);
    let consumed=0;
    for(const part of best.parts){
      if(within<=consumed+part.text.length){
        line=part.line;
        const leading=part.raw.length-part.raw.trimStart().length;
        const syntax=part.raw.slice(leading).match(/^(?:#{1,6}\s+|>\s?|[-*+]\s+|\d+[.)]\s+)/)?.[0].length||0;
        ch=Math.min(part.raw.length,leading+syntax+Math.max(0,within-consumed));
        break;
      }
      consumed+=part.text.length+1;
    }
  }
  ch=Math.max(0,Math.min(lines[line]?.length||0,ch));
  editor.focus();
  editor.setCursor({line,ch});
};
window.addEventListener('message',event=>{
  if (event.source !== frame.contentWindow) return;
  if (event.data?.type === 'preview-content-ready') markPreviewContentReady();
  if (event.data?.type === 'preview-features-ready') markPreviewReady();
  if (event.data?.type === 'open-post') window.MarkataSiteTree?.openPost(event.data.slug);
  if (event.data?.type === 'place-editor-cursor') placeEditorCursor(event.data);
  if (event.data?.type === 'workbench-preview-focus-editor')
    window.MarkataShortcuts?.execute('workbench.focusEditor');
  if (event.data?.type === 'preview-updated' && event.data.version === revision) {
    finishPreview(event.data.version,event.data.patch?.timing);
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
  indentUnit:2, tabSize:4, viewportMargin:10, cursorScrollMargin:110, screenReaderLabel:'Markdown source'
});
const excalidrawActions = window.MarkataExcalidraw.attach(editor);
const settingsDialog = document.querySelector('#settings-dialog');
document.querySelector('#settings-open').addEventListener('click',() => window.MarkataSettings.open());
document.querySelector('#settings-close').addEventListener('click',() => settingsDialog.close());
const wrapSetting = document.querySelector('#setting-wrap');
try { wrapSetting.checked = localStorage.getItem('markata-editor-wrap') !== 'false'; } catch { storageWarning(); }
editor.setOption('lineWrapping',wrapSetting.checked);
wrapSetting.addEventListener('change',() => {editor.setOption('lineWrapping',wrapSetting.checked);try {localStorage.setItem('markata-editor-wrap',String(wrapSetting.checked));} catch {storageWarning();}});
const workbench = document.querySelector('#workbench');
const layoutSetting = document.querySelector('#setting-layout');
try { layoutSetting.value = localStorage.getItem('markata-editor-layout') || 'split'; } catch { storageWarning(); }
function setLayout(layout) {
  const previouslyFocusedPane = document.activeElement?.closest?.('#editor-pane, #preview-pane')?.id;
  workbench.dataset.layout = ['split','editor','preview'].includes(layout) ? layout : 'split';
  layoutSetting.value = workbench.dataset.layout;
  try { localStorage.setItem('markata-editor-layout',workbench.dataset.layout); } catch { storageWarning(); }
  if (previouslyFocusedPane === 'editor-pane' && workbench.dataset.layout === 'preview')
    document.querySelector('#preview-pane')?.focus();
  if (previouslyFocusedPane === 'preview-pane' && workbench.dataset.layout === 'editor')
    editor.focus();
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
const createLinkField = document.querySelector('#create-link-field');
const createLinkURL = document.querySelector('#create-link-url');
const createLinkMetadata = document.querySelector('#create-link-metadata');
const createLinkMetadataStatus = document.querySelector('#create-link-metadata-status');
const createLinkDescription = document.querySelector('#create-link-description');
const createLinkSiteName = document.querySelector('#create-link-site-name');
const createLinkMediaType = document.querySelector('#create-link-media-type');
const createLinkProvider = document.querySelector('#create-link-provider');
const createLinkAuthor = document.querySelector('#create-link-author');
const createLinkAuthorURL = document.querySelector('#create-link-author-url');
const createImageField = document.querySelector('#create-image-field');
const createImageURL = document.querySelector('#create-image-url');
const linkMetadataControls = {
  title: createTitle,
  description: createLinkDescription,
  image: createImageURL,
  site_name: createLinkSiteName,
  type: createLinkMediaType,
  provider_name: createLinkProvider,
  author_name: createLinkAuthor,
  author_url: createLinkAuthorURL
};
const createImageFile = document.querySelector('#create-image-file');
const createImagePaste = document.querySelector('#create-image-paste');
const createImagePreview = document.querySelector('#create-image-preview');
const createImageName = document.querySelector('#create-image-name');
const createImageStorage = document.querySelector('#create-image-storage');
const createNumbering = document.querySelector('#create-numbering');
const createFilenameMode = document.querySelector('#create-filename-mode');
const createRememberDefaults = document.querySelector('#create-remember-defaults');
const createDate = document.querySelector('#create-date');
const createSlug = document.querySelector('#create-slug');
const createDirectory = document.querySelector('#create-directory');
const createDirectoryHelp = document.querySelector('#create-directory-help');
const createPathPreview = document.querySelector('#create-path-preview');
let contentTemplates = [];
let imageUploadProvider = 'local';
let selectedContentTemplate = null;
let activeMobileShare = null;
let selectedImageFile = null;
let imagePreviewURL = '';
let numberLookupRevision = 0;
let linkMetadataTimer = 0;
let linkMetadataRevision = 0;
let linkMetadataResolvedURL = '';
let linkMetadataApplied = {};
let linkMetadataTouched = new Set();
let applyingLinkMetadata = false;
let linkMetadataPromise = Promise.resolve();
async function showCreateDialog(preferredTemplateID = '') {
  try {
    const response = await fetch('/workspace/templates',{cache:'no-store'});
    const data = await response.json();
    if (!response.ok) throw new Error(data.error || 'Templates are unavailable.');
    contentTemplates = data.templates || [];
    imageUploadProvider = data.image_upload || 'local';
    templateQuery.value = '';
    renderTemplateChoices();
    if (typeof preferredTemplateID === 'string') {
      const chosen = contentTemplates.find(item => item.id === preferredTemplateID);
      if (chosen) { chooseTemplate(chosen); return; }
    }
    templatePicker.showModal();
    templateQuery.focus();
  } catch (error) {
    createError.textContent = error.message; createError.hidden = false;
    createDialog.showModal();
  }
}
let templateMatches = [], activeTemplate = 0;
function templateScore(template, query) {
  const haystack = `${template.name} ${(template.aliases || []).join(' ')} ${template.description} ${template.id}`.toLowerCase();
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
  numberLookupRevision++;
  selectedContentTemplate = template;
  templatePicker.dataset.selectedTemplate = template.id;
  templatePicker.close();
  createError.hidden = true;
  document.querySelector('#create-description').textContent = template.description || template.name;
  const titleMatch = template.source.match(/^title:\s*(.*)$/m);
  const parsedTitle = titleMatch?.[1]?.trim().replace(/^(?:"([\s\S]*)"|'([\s\S]*)')$/, '$1$2');
  const title = !parsedTitle || parsedTitle.includes('{{') ? template.name : parsedTitle;
  const isLinkPost = template.id === 'builtin:link';
  const isImagePost = ['builtin:photo','builtin:shot'].includes(template.id) || /^image:\s*/m.test(template.source);
  createLinkField.hidden = !isLinkPost;
  createLinkMetadata.hidden = !isLinkPost;
  createLinkURL.required = isLinkPost;
  createLinkURL.value = '';
  createTitle.required = !isLinkPost;
  createImageField.hidden = !isImagePost && !isLinkPost;
  clearTimeout(linkMetadataTimer);
  linkMetadataRevision++;
  linkMetadataResolvedURL = '';
  linkMetadataApplied = {};
  linkMetadataTouched.clear();
  createLinkMetadataStatus.textContent = '';
  createLinkDescription.value = '';
  createLinkSiteName.value = '';
  createLinkMediaType.value = '';
  createLinkProvider.value = '';
  createLinkAuthor.value = '';
  createLinkAuthorURL.value = '';
  createImageStorage.textContent = imageUploadProvider === 'dropper'
    ? 'Uploaded images go to Dropper.'
    : 'Uploaded images are copied to static/images/.';
  createImageURL.value = '';
  createImageFile.value = '';
  clearSelectedImage();
  createTitle.value = isLinkPost ? '' : title;
  createTitle.placeholder = isLinkPost ? 'Filled from the source URL, or enter a title' : '';
  createSlug.value = '';
  createSlug.readOnly = false;
  createSlug.placeholder = slugify(title) || 'post-title';
  createDirectory.value = defaultContentDirectory(template);
  createDirectoryHelp.textContent = `Current site default for ${template.name}: ${template.directory || createDirectory.value}.`;
  createDate.value = '';
  createNumbering.hidden = false;
  createFilenameMode.value = template.filename_mode || (template.numbered_by_default ? 'numbered' : 'title');
  createRememberDefaults.checked = true;
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
    if (template.aliases?.length) {
      const aliases = document.createElement('span');
      aliases.className = 'template-aliases';
      aliases.textContent = `Also known as: ${template.aliases.join(', ')}`;
      button.append(aliases);
    }
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
function clearSelectedImage() {
  selectedImageFile = null;
  if (imagePreviewURL) URL.revokeObjectURL(imagePreviewURL);
  imagePreviewURL = '';
  createImagePreview.removeAttribute('src');
  createImagePreview.hidden = true;
  createImageName.hidden = true;
  createImageName.textContent = '';
}
function selectImageFile(file) {
  if (!file || !/^image\/(?:png|jpeg|gif|webp)$/.test(file.type)) {
    createError.textContent = 'Choose a PNG, JPEG, GIF, or WebP image.';
    createError.hidden = false;
    return;
  }
  if (file.size > 15 * 1024 * 1024) {
    createError.textContent = 'Image must be 15 MiB or smaller.';
    createError.hidden = false;
    return;
  }
  createError.hidden = true;
  createImageURL.value = '';
  clearSelectedImage();
  selectedImageFile = file;
  imagePreviewURL = URL.createObjectURL(file);
  createImagePreview.src = imagePreviewURL;
  createImagePreview.hidden = false;
  createImageName.textContent = `${file.name || 'Pasted image'} · ${(file.size / 1024 / 1024).toFixed(1)} MiB`;
  createImageName.hidden = false;
}
function previewRemoteImage(value) {
  const raw = value.trim();
  if (!raw) return;
  let parsed;
  try { parsed = new URL(raw); } catch { return; }
  if (!['http:','https:'].includes(parsed.protocol) || !parsed.hostname || parsed.username || parsed.password) return;
  createImagePreview.src = parsed.href;
  createImagePreview.hidden = false;
  createImageName.textContent = 'Previewing remote image';
  createImageName.hidden = false;
}
function defaultContentDirectory(template) {
  if (template.directory) return template.directory;
  const key = template.id.replace(/^builtin:/,'');
  return ({thought:'pages/thoughts',note:'pages/notes',photo:'pages/photos'})[key] || 'pages/blog';
}
function updateCreatePathPreview() {
  const template = selectedContentTemplate;
  if (createFilenameMode.value === 'numbered') {
    const number = template.next_number || 1;
    const prefix = template.number_prefix || slugify(template.name);
    createSlug.value = `${prefix}-${number}`;
    createSlug.readOnly = true;
  } else {
    if (createSlug.readOnly) createSlug.value = '';
    createSlug.readOnly = false;
  }
  const titleSlug = template?.id === 'builtin:link'
    ? window.MarkataContentStarters.suggestLinkSlug(createTitle.value, createLinkURL.value)
    : slugify(createTitle.value);
  const slug = createSlug.value.trim() || titleSlug || 'new-post';
  const directory = createDirectory.value.trim().replace(/\\/g,'/').replace(/^\/+|\/+$/g,'');
  createSlug.placeholder = titleSlug || 'post-title';
  const path = `${directory ? `${directory}/` : ''}${slug}.md`;
  const generatedLinkSlug = window.MarkataContentStarters.generatedLinkSlug(
    template?.source, template?.name, createFilenameMode.value, createSlug.value, slug);
  createPathPreview.textContent = generatedLinkSlug ? `${path} · URL /${generatedLinkSlug}/` : path;
}
async function refreshNextContentNumber() {
  if (!selectedContentTemplate || createFilenameMode.value !== 'numbered') return;
  const revision = ++numberLookupRevision;
  const directory = createDirectory.value.trim().replace(/\\/g,'/').replace(/^\/+|\/+$/g,'');
  const prefix = selectedContentTemplate.number_prefix || slugify(selectedContentTemplate.name);
  try {
    const query = new URLSearchParams({directory,prefix});
    const response = await fetch('/workspace/content-number?'+query.toString(),{cache:'no-store'});
    const data = await response.json();
    if (response.ok && revision === numberLookupRevision && Number.isInteger(data.next_number)) {
      selectedContentTemplate.next_number = data.next_number;
      updateCreatePathPreview();
    }
  } catch {}
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
for (const [key, control] of Object.entries(linkMetadataControls)) {
  control.addEventListener('input',() => {
    if (applyingLinkMetadata) return;
    linkMetadataTouched.add(key);
    if (linkMetadataApplied[key] && control.value.trim() !== linkMetadataApplied[key]) delete linkMetadataApplied[key];
  });
}
createLinkURL.addEventListener('input',() => {scheduleLinkMetadata();updateCreatePathPreview();});
createSlug.addEventListener('input',updateCreatePathPreview);
createDirectory.addEventListener('input',updateCreatePathPreview);
createDirectory.addEventListener('change',refreshNextContentNumber);
createFilenameMode.addEventListener('change',() => {updateCreatePathPreview();refreshNextContentNumber();});
createImageURL.addEventListener('input',() => {
  if(createImageURL.value.trim()){
    createImageFile.value='';
    clearSelectedImage();
    previewRemoteImage(createImageURL.value);
  }else clearSelectedImage();
});
createImageFile.addEventListener('change',() => selectImageFile(createImageFile.files?.[0]));
const createCameraFile = document.querySelector('#create-camera-file');
createCameraFile.addEventListener('change',() => {
  const file = createCameraFile.files?.[0];
  createCameraFile.value = '';
  if (file) selectImageFile(file);
});
createImagePaste.addEventListener('click',() => createImageFile.click());
createImagePaste.addEventListener('keydown',event => {
  if(event.key==='Enter'||event.key===' '){event.preventDefault();createImageFile.click();}
});
createImagePaste.addEventListener('dragover',event => {event.preventDefault();createImagePaste.dataset.dragging='true';});
createImagePaste.addEventListener('dragleave',() => delete createImagePaste.dataset.dragging);
createImagePaste.addEventListener('drop',event => {
  event.preventDefault();delete createImagePaste.dataset.dragging;
  selectImageFile([...event.dataTransfer.files].find(file=>file.type.startsWith('image/')));
});
createDialog.addEventListener('paste',event => {
  if(!createDialog.open)return;
  const item=[...(event.clipboardData?.items||[])].find(candidate=>candidate.type.startsWith('image/'));
  const file=item?.getAsFile();
  if(file){event.preventDefault();selectImageFile(file);}
});
createDialog.addEventListener('close',()=>{clearSelectedImage();activeMobileShare=null;});
async function loadLinkMetadata() {
  const rawURL=createLinkURL.value.trim();
  if(!rawURL)return;
  const revision=++linkMetadataRevision;
  linkMetadataResolvedURL=rawURL;
  let parsedURL;
  try{parsedURL=new URL(rawURL);}catch{
    createLinkMetadataStatus.textContent='Enter a complete HTTP(S) URL to look up metadata.';
    return;
  }
  if(!['http:','https:'].includes(parsedURL.protocol)||!parsedURL.hostname||parsedURL.username||parsedURL.password){
    createLinkMetadataStatus.textContent='Enter a complete HTTP(S) URL without credentials.';
    return;
  }
  createLinkMetadataStatus.textContent='Looking up page metadata…';
  try{
    const response=await fetch('/workspace/link-metadata',{
      method:'POST',headers:{'Content-Type':'application/json'},
      body:JSON.stringify({url:rawURL}),cache:'no-store'
    });
    const data=await response.json();
    if(!response.ok)throw new Error(data.error||'Metadata lookup failed.');
    if(revision!==linkMetadataRevision||rawURL!==createLinkURL.value.trim())return;
    // User edits win even when the field was subsequently cleared. Do not
    // mistake synthetic image-preview input events for manual edits.
    applyingLinkMetadata=true;
    try{
      for(const [key,control] of Object.entries(linkMetadataControls)){
        const incoming=String(data[key]||'');
        const current=control.value.trim();
        const previous=linkMetadataApplied[key]||'';
        if(!linkMetadataTouched.has(key)&&(!current||(previous&&current===previous))){
          control.value=incoming;
          if(incoming)linkMetadataApplied[key]=incoming;
          else delete linkMetadataApplied[key];
          if(key==='image')control.dispatchEvent(new Event('input',{bubbles:true}));
        }
      }
    }finally{applyingLinkMetadata=false;}
    const provider=data.provider_name||data.site_name;
    const source=data.source?` via ${data.source}`:'';
    if(data.available){
      createLinkMetadataStatus.textContent=`Loaded metadata${provider?` from ${provider}`:''}${source}. Review or edit any field.`;
      if(!data.title)createLinkMetadataStatus.textContent+=' Add a title to create the post.';
    }else createLinkMetadataStatus.textContent='No embed metadata found. Enter a title to create the post.';
    updateCreatePathPreview();
  }catch(error){
    if(revision===linkMetadataRevision&&rawURL===createLinkURL.value.trim())
      createLinkMetadataStatus.textContent='Could not fetch metadata. Enter a title and any other details manually.';
  }
}
function scheduleLinkMetadata() {
  clearTimeout(linkMetadataTimer);
  linkMetadataRevision++;
  linkMetadataResolvedURL='';
  const rawURL=createLinkURL.value.trim();
  if(!rawURL){createLinkMetadataStatus.textContent='';return;}
  createLinkMetadataStatus.textContent='Metadata lookup starts when the URL is complete.';
  linkMetadataTimer=setTimeout(()=>{linkMetadataPromise=loadLinkMetadata();},450);
}
// Home-screen quick actions reuse the same authoring templates as the desktop editor.
// A failed remote template request is shown by showCreateDialog, never silently saved.
window.addEventListener('markata-compose', event => {
  const kind = event.detail;
  if (!['note','photo','link'].includes(kind) || createDialog.open || templatePicker.open) return;
  void showCreateDialog('builtin:' + kind);
});
window.addEventListener('markata-share-select', async event => {
  const shared = event.detail;
  if (!shared || typeof shared.id !== 'string') return;
  // Never overwrite an open create dialog or publish without an explicit user action.
  if (createDialog.open || templatePicker.open) return;
  const id = shared.image ? 'builtin:photo' : shared.url ? 'builtin:link' : 'builtin:note';
  await showCreateDialog(id);
  if (!createDialog.open) return;
  activeMobileShare = shared;
  const name = shared.image?.name?.replace(/\.[a-z0-9]+$/i,'') || '';
  const fallback = shared.image ? (name || 'Shared photo') : shared.url ? '' : 'Shared text';
  const sharedTitle = String(shared.title || fallback).trim();
  if (sharedTitle) {
    createTitle.value = sharedTitle;
    createTitle.dispatchEvent(new Event('input', {bubbles:true}));
  }
  if (shared.url && id === 'builtin:link') {
    createLinkURL.value = shared.url;
    createLinkURL.dispatchEvent(new Event('input', {bubbles:true}));
  }
  if (shared.image) selectImageFile(shared.image);
  updateCreatePathPreview();
});

document.querySelector('#create-form').addEventListener('submit',async event => {
  event.preventDefault(); createError.hidden = true;
  const submitButton=document.querySelector('#create-submit');
  submitButton.disabled=true;
  const template = contentTemplates.find(item => item.id === templatePicker.dataset.selectedTemplate);
  const filenameMode = createFilenameMode.value;
  const directory = createDirectory.value.trim().replace(/\\/g,'/').replace(/^\/+|\/+$/g,'');
  try {
    if (!template) throw new Error('Choose a content template first.');
    if (template.id === 'builtin:link') {
      const rawLinkURL=createLinkURL.value.trim();
      let sourceURL;
      try { sourceURL = new URL(rawLinkURL); }
      catch { throw new Error('Enter a valid source URL beginning with https:// or http://.'); }
      if (!['https:','http:'].includes(sourceURL.protocol) ||
          !sourceURL.hostname || sourceURL.username || sourceURL.password) {
        throw new Error('Use a public-style http(s) URL without embedded credentials.');
      }
      clearTimeout(linkMetadataTimer);
      if (linkMetadataResolvedURL !== rawLinkURL) linkMetadataPromise=loadLinkMetadata();
      await linkMetadataPromise;
    }
    if (filenameMode === 'numbered') await refreshNextContentNumber();
    const title = createTitle.value.trim();
    if (!title) throw new Error(template.id === 'builtin:link'
      ? 'No title was found for this URL. Enter a title to create the link post.'
      : 'Enter a title to create the post.');
    const enteredFilename = createSlug.value.trim();
    const titleSlug = template.id === 'builtin:link' && filenameMode === 'title'
      ? window.MarkataContentStarters.suggestLinkSlug(title, createLinkURL.value)
      : slugify(title);
    const slug = enteredFilename || titleSlug;
    if (!slug) throw new Error('Enter a title or file name for this post.');
    const path = `${directory && directory !== '.' ? `${directory}/` : ''}${slug}.md`;
    let source = window.MarkataContentStarters.render(template.source,title,slug,new Date(),createDate.value);
    const generatedLinkSlug = window.MarkataContentStarters.generatedLinkSlug(
      source, template.name, filenameMode, enteredFilename, slug);
    if (generatedLinkSlug)
      source=window.MarkataContentStarters.setFrontmatterValue(source,'slug',generatedLinkSlug);
    if (template.id === 'builtin:link') {
      source=window.MarkataContentStarters.setFrontmatterValue(source,'url',createLinkURL.value.trim());
      for (const key of ['description','site_name','type','provider_name','author_name','author_url']) {
        const value=linkMetadataControls[key].value.trim();
        if(value)source=window.MarkataContentStarters.setFrontmatterValue(source,key,value);
      }
    }
    if (activeMobileShare) {
      // Shared content is always a new unpublished post, not a public publication.
      source = source.replace(/^published:[^\r\n]*$/m, 'published: false');
      const sharedText=String(activeMobileShare.text || '').trim();
      const caption = activeMobileShare.url && sharedText === activeMobileShare.url ? '' : sharedText;
      if (caption) source = source.trimEnd() + '\n\n' + caption + '\n';
      if (activeMobileShare.url && template.id !== 'builtin:link')
        source += '\n' + activeMobileShare.url + '\n';
    }
    const imageURL=createImageURL.value.trim();
    if(imageURL){
      let parsedImageURL;
      try{parsedImageURL=new URL(imageURL);}catch{throw new Error('Enter a valid image URL beginning with https:// or http://.');}
      if(!['https:','http:'].includes(parsedImageURL.protocol)||!parsedImageURL.hostname||parsedImageURL.username||parsedImageURL.password)
        throw new Error('Use a public-style http(s) image URL without embedded credentials.');
      source=setFrontmatterValue(source,'image',imageURL);
    }
    const changedDefaults = directory !== template.directory || filenameMode !== (template.filename_mode || 'title');
    if (createRememberDefaults.checked && changedDefaults) {
      const settingsResponse = await fetch('/workspace/template-settings', {
        method:'POST',headers:{'Content-Type':'application/json'},
        body:JSON.stringify({id:template.id,directory,filename_mode:filenameMode})
      });
      const settingsData = await settingsResponse.json();
      if (!settingsResponse.ok) throw new Error(settingsData.error || 'Cannot save these site defaults.');
    }
    let response;
    if(selectedImageFile){
      const form=new FormData();form.set('path',path);form.set('source',source);form.set('image_name',slug);form.set('image',selectedImageFile,selectedImageFile.name||'clipboard-image');
      response=await fetch('/workspace/new-post',{method:'POST',body:form});
    }else{
      response=await fetch('/workspace/file?path='+encodeURIComponent(path),{method:'POST',headers:{'Content-Type':'text/plain;charset=utf-8'},body:source});
    }
    const data = await response.json();
    if (!response.ok) throw new Error(data.error || 'Cannot create this file.');
    if (activeMobileShare && window.MarkataMobileShare) {
      try { await window.MarkataMobileShare.consume(activeMobileShare.id); }
      catch (error) { console.warn('Share inbox item needs manual cleanup:',error); }
    }
    window.MarkataSiteTree?.activate(path);
    createDialog.close(); location.reload();
  } catch (error) { createError.textContent = error.message; createError.hidden = false; }
  finally {submitButton.disabled=false;}
});
window.MarkataFrontmatter.attach(editor);
const appShortcuts = window.MarkataShortcuts;
appShortcuts.init(editor);
function toggleSidebar(visible = workbench.classList.contains('sidebar-hidden')) {
  const sidebar = document.querySelector('#site-sidebar');
  const focusedInside = sidebar.contains(document.activeElement);
  if (!visible && focusedInside) {
    if (workbench.dataset.layout === 'preview') document.querySelector('#preview-pane').focus();
    else editor.focus();
  }
  workbench.classList.toggle('sidebar-hidden', !visible);
  sidebar.inert = !visible;
  sidebar.setAttribute('aria-hidden', String(!visible));
  editor.refresh();
}
function focusExplorer() {
  toggleSidebar(true);
  const tree = document.querySelector('.site-tree:not([hidden])');
  (tree?.querySelector('[tabindex="0"], button, [role="treeitem"]') ||
    document.querySelector('#workspace-view-tabs [aria-selected="true"]'))?.focus();
}
const registerApp = (id,label,category,bindings,run,description='') =>
  appShortcuts.register({id,label,category,bindings,run,description});
// Make editor motion keymaps discoverable without opening Settings. The
// palette executes the same change event as the existing settings control.
for (const [value, name] of [['default','Default'],['vim','Vim'],['emacs','Emacs'],['sublime','Sublime']]) {
  registerApp('workbench.editorKeymap.'+value,'Editor keymap: '+name,'Settings',{},
    () => {
      if (keymap.value !== value) {
        keymap.value = value;
        keymap.dispatchEvent(new Event('change'));
      } else focusEditorForWriting();
    },
    'Switch Markdown text editing keys to '+name+' (separate from application shortcuts).');
}
registerApp('workbench.toggleSidebar','Toggle File Explorer','Navigation',
  {vscode:'Mod+B',writer:'Mod+\\',vim:'Space E',emacs:'Ctrl+X Ctrl+B'},toggleSidebar);
registerApp('workbench.focusExplorer','Focus File Explorer','Navigation',
  {vscode:'Mod+Shift+E',writer:'Mod+Shift+E',vim:'Space F E',emacs:'Ctrl+X Ctrl+E'},focusExplorer);
registerApp('workbench.focusEditor','Focus Editor','Navigation',
  {vscode:'Mod+Alt+1',writer:'Mod+Alt+1',vim:'Space 1',emacs:'Alt+1'},
  () => {if(workbench.dataset.layout==='preview')setLayout('split');focusEditorForWriting();});
registerApp('workbench.focusPreview','Focus Preview','Navigation',
  {vscode:'Mod+Alt+2',writer:'Mod+Alt+2',vim:'Space 2',emacs:'Alt+2'},
  () => {if(workbench.dataset.layout==='editor')setLayout('split');document.querySelector('#preview-pane').focus();});
registerApp('workbench.preview','Toggle Preview Only','Layout',
  {vscode:'Mod+Alt+V',writer:'Mod+Alt+V',vim:'Space P V',emacs:'Alt+V'},
  () => setLayout(workbench.dataset.layout==='preview'?'split':'preview'));
registerApp('workbench.lineWrap','Toggle Line Wrapping','Writing',
  {vscode:'Alt+Z',writer:'Alt+Z',vim:'Alt+Z',emacs:''},
  () => {wrapSetting.checked=!wrapSetting.checked;wrapSetting.dispatchEvent(new Event('change'));});
registerApp('workbench.problems','Toggle Problems','Diagnostics',
  {vscode:'Mod+Shift+M',writer:'Mod+Shift+M',vim:'Space X',emacs:''},toggleProblems);
registerApp('workbench.save','Publish Current Post','Writing',
  {vscode:'Mod+S',writer:'Mod+S',vim:'Mod+S',emacs:'Ctrl+X Ctrl+S'},
  () => window.MarkataPublishReview?.quickPublish(),
  'Publish only the active site post. Other changes remain untouched.');
registerApp('workbench.newPost','New Post','Files',
  {vscode:'Mod+Alt+N',writer:'Mod+Alt+N',vim:'Space F N',emacs:''},
  showCreateDialog,'Create a workspace post from a template.');
registerApp('workbench.openMarkdown','Open Markdown File','Files',
  {vscode:'',writer:'',vim:'',emacs:''},
  () => document.querySelector('#open').click(),'Import a local Markdown file as a draft.');
registerApp('workbench.downloadMarkdown','Download Markdown File','Files',
  {vscode:'',writer:'',vim:'',emacs:''},
  () => document.querySelector('#download').click(),'Export the current source as Markdown.');
registerApp('workbench.shortcutSettings','Keyboard Shortcuts','Settings',
  {vscode:'Mod+K Mod+S',writer:'Mod+Alt+K',vim:'Space ?',emacs:'Alt+Shift+K'},
  () => appShortcuts.openSettings(),'View, change, or reset application keybindings.');
window.MarkataSettings.attach(editor);
window.MarkataDiagnostics.attach();
window.MarkataPublishReview.attach(editor);
window.MarkataCommandPalette.attach(editor,{actions:[
  ...excalidrawActions,
  {id:'edit-properties',label:'Page: Edit Properties',description:'Open the visual frontmatter object editor.',aliases:'frontmatter metadata tags date title publish fields',run:()=>window.MarkataFrontmatter.open()},
  {id:'add-property',label:'Page: Add a Property',description:'Search and add a frontmatter property.',aliases:'frontmatter metadata field key',run:()=>window.MarkataFrontmatter.openAddProperty()},
  {id:'back-to-markdown',label:'Page: Back to Markdown',description:'Close page details and return to the article source.',aliases:'frontmatter source article escape close',run:()=>window.MarkataFrontmatter.close()},
  {id:'edit-frontmatter-yaml',label:'Page: Edit YAML Frontmatter',description:'Hand edit the raw YAML header.',aliases:'frontmatter source advanced raw metadata',run:()=>window.MarkataFrontmatter.editYAML()},

]});
document.querySelector('#preview-focus-editor').addEventListener('click', () =>
  appShortcuts.execute('workbench.focusEditor'));
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
  scrollSync.vimModeChanged();
});
// Rejected imports must not be dismissed by an already-running preview response.
let rejectedImportRevision = -1;
let revision = 0, timer, busy = false, ready = false, lastEditAt = null, pendingSince = null;
let pendingRender = false;
let activeRenderController = null;
let previewStarted = false, previewContentReady = false, previewReady = false, pendingPreview = null, activeRender = null;
let previewFrameStartedAt = null, previewFrameContentReadyMS = null;
function markPreviewContentReady() {
  if (!previewStarted || previewContentReady) return;
  previewContentReady = true;
  if (previewFrameStartedAt !== null) {
    previewFrameContentReadyMS = performance.now() - previewFrameStartedAt;
    window.MarkataDiagnostics?.record('preview_frame_content_ready', previewFrameContentReadyMS);
  }
  scrollSync.ready();
  scrollSync.previewUpdated();
  if (activeRender?.initialDocument) activeRender.frameContentReadyMS = previewFrameContentReadyMS;
  finishPreview(activeRender?.version);
}
function markPreviewReady() {
  if (!previewStarted || previewReady) return;
  previewReady = true;
  syncAppearance();
  scrollSync.ready();
  if (previewFrameStartedAt !== null) {
    window.MarkataDiagnostics?.record('preview_feature_startup', performance.now() - previewFrameStartedAt);
  }
  if (pendingPreview) {
    const latest = pendingPreview;
    pendingPreview = null;
    if (latest.version === revision) updatePreview(latest.data, latest.version);
  } else {
    scrollSync.previewUpdated();
    finishPreview(activeRender?.version);
  }
}
function markPreviewLoadComplete() {
  if (previewFrameStartedAt !== null) {
    window.MarkataDiagnostics?.record('preview_iframe_load', performance.now() - previewFrameStartedAt);
  }
  markPreviewReady();
}
function finishPreview(version, patchTiming) {
  if (!activeRender || activeRender.pipelineRecorded || activeRender.version !== version || version !== revision) return;
  window.MarkataStartupStatus?.previewReady();
  const {data} = activeRender;
  const finishedAt = performance.now();
  const editAt = activeRender.editAt ?? activeRender.startedAt;
  const responseAt = activeRender.responseAt ?? activeRender.startedAt;
  const requestMS = Math.max(0,responseAt-activeRender.startedAt);
  const queueMS = Math.max(0,data.queue_ms ?? 0);
  const renderMS = Math.max(0,data.render_ms ?? 0);
  const stages = {
    waiting_ms:Math.max(0,activeRender.startedAt-editAt),
    rendering_ms:queueMS+renderMS,
    network_ms:Math.max(0,requestMS-queueMS-renderMS),
    display_ms:Math.max(0,finishedAt-responseAt),
  };
  const liveTime = Math.round(finishedAt-editAt);
  window.MarkataDiagnostics?.record('edit_to_preview',liveTime);
  status.textContent = `Preview ${liveTime} ms`;
  status.title = `${activeRender.mode} mode · Go ${data.render_ms.toFixed(1)} ms${data.cached ? ' (cached)' : ''} · queue ${(data.queue_ms ?? 0).toFixed(1)} ms · input ${(data.input_bytes/1024).toFixed(1)} KB · fragment ${((data.html_bytes)/1024).toFixed(1)} KB`;
  if (!activeRender.pipelineRecorded) {
    activeRender.pipelineRecorded = true;
    try {
      window.MarkataDiagnostics?.recordPreviewPipeline(stages,{
        request_ms:requestMS,queue_ms:queueMS,go_render_ms:renderMS,cached:data.cached,
        input_bytes:data.input_bytes,html_bytes:data.html_bytes,mode:activeRender.mode,iframe:patchTiming||null,
        frame_content_ready_ms:activeRender.frameContentReadyMS ?? previewFrameContentReadyMS,
      });
    } catch(error) { console.error('Could not record preview timings',error); }
  }
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
function previewDelay() {
  const now = performance.now();
  const mode = activePreviewMode();
  const cadence = mode === 'fast' ? {idle:0,burst:0} : mode === 'slow' ? {idle:300,burst:1000} : {idle:50,burst:250};
  const idle = lastEditAt === null ? 0 : Math.max(0, cadence.idle - (now - lastEditAt));
  return pendingSince === null ? idle : Math.min(idle, Math.max(0, cadence.burst - (now - pendingSince)));
}
let fastRenderQueued=false;
function schedulePreviewRender() {
  clearTimeout(timer);
  if(activePreviewMode()==='fast') {
    if(fastRenderQueued)return;
    fastRenderQueued=true;
    queueMicrotask(()=>{fastRenderQueued=false;void render();});
    return;
  }
  timer=setTimeout(render,previewDelay());
}
const previewModeStorageKey='markata-preview-mode';
let previewModePreference='normal', batterySaverActive=false;
try {
  const stored=localStorage.getItem(previewModeStorageKey);
  if(['auto','fast','normal','slow'].includes(stored))previewModePreference=stored;
} catch {}
function autoPreviewShouldSlow() {
  const connection=navigator.connection||navigator.mozConnection||navigator.webkitConnection;
  const effectiveType=connection?.effectiveType;
  return batterySaverActive||connection?.saveData===true||effectiveType==='slow-2g'||effectiveType==='2g'||
    (Number.isFinite(navigator.deviceMemory)&&navigator.deviceMemory<=4)||
    (Number.isFinite(navigator.hardwareConcurrency)&&navigator.hardwareConcurrency<=4);
}
function activePreviewMode() {
  return previewModePreference==='auto'?(autoPreviewShouldSlow()?'slow':'normal'):previewModePreference;
}
function previewModeChanged(mode) {
  if(['auto','fast','normal','slow'].includes(mode))previewModePreference=mode;
  if(pendingSince!==null)schedulePreviewRender();
}
window.addEventListener('markata-preview-mode-change',event=>previewModeChanged(event.detail?.mode));
const connectionStatus=navigator.connection||navigator.mozConnection||navigator.webkitConnection;
connectionStatus?.addEventListener?.('change',()=>previewModeChanged(previewModePreference));
try {
  if(typeof navigator.getBattery==='function')navigator.getBattery().then(battery=>{
    const update=()=>{
      const next=!battery.charging&&battery.level<=0.2;
      if(next===batterySaverActive)return;
      batterySaverActive=next;previewModeChanged(previewModePreference);
    };
    update();battery.addEventListener('chargingchange',update);battery.addEventListener('levelchange',update);
  }).catch(()=>{});
} catch {}
async function render() {
  if (!ready) return;
  if (busy) {
    // Keep the newest request pending rather than losing changes made while
    // Go is still rendering an earlier revision.
    pendingRender = true;
    return;
  }
  pendingRender = false;
  pendingSince = null;
  // The CodeMirror document is the source of truth. The textarea is only
  // a form mirror and may lag a programmatic document replacement.
  const version = revision, originalMarkdown = editor.getValue();
  const initialDocument = !previewStarted;
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
  activeRenderController = controller;
  const timeout = setTimeout(() => controller.abort(), 30000);
  const started = performance.now();
  const editAt = lastEditAt;
  const renderMode = activePreviewMode();
  try {
    const response = await fetch('/render', {
      method:'POST', headers:{'Content-Type':'application/json'},
      body:JSON.stringify({markdown, frontmatter_source:frontmatter?.[0] || '', include_document:!previewStarted, appearance:Object.fromEntries(Object.entries(attributes()).filter(([,v])=>v))}),
      signal:controller.signal
    });
    const data = await response.json();
    const responseAt = performance.now();
    if (version !== revision) return;
    if (!response.ok) throw new Error(data.error || 'Rendering failed.');
    if (sourceLineOffset) {
      data.source_lines = data.source_lines.map(line => line + sourceLineOffset);
      data.source_anchors = data.source_anchors.map(line => line + sourceLineOffset);
    }
    data.source_offset = sourceLineOffset;
    scrollSync.update(data.source_lines, false, sourceLineOffset, data.source_anchors);
    activeRender = {data, version, editAt, startedAt:started, responseAt, mode:renderMode, initialDocument};
    if (!previewStarted) {
      previewStarted = true;
      frame.onload = markPreviewLoadComplete;
      previewFrameStartedAt = performance.now();
      frame.srcdoc = data.document;
    } else if (!previewReady) {
      pendingPreview = {data, version};
    } else {
      updatePreview(data, version);
    }
    if (rejectedImportRevision !== version) errorBox.hidden = true;
    status.textContent = 'Updating preview…';
  } catch (error) {
    if (version !== revision) return;
    errorBox.textContent = error.name === 'AbortError' ? 'Preview timed out. Edit the Markdown or reload to retry.' : error.message;
    errorBox.hidden = false;
    status.textContent = 'Preview failed. Edit the Markdown to try again.';
  } finally {
    clearTimeout(timeout);
    if (activeRenderController === controller) activeRenderController = null;
    busy = false;
    // Stale or skipped renders must always yield to the latest editor text.
    // Keep one pending revision; its selected preview mode controls when it runs.
    if (pendingRender || version !== revision) {
      pendingRender = false;
      schedulePreviewRender();
    }
  }
}
editor.on('change', (_instance, change) => {
  // Discard and file-open replace the entire CodeMirror document. An old
  // preview scroll map must not be reused for the replacement's cursor.
  if (change.origin === 'setValue') scrollSync.reset();
  editor.save(); revision++; clearTimeout(timer);
  // A render already posted into the iframe may still be queued or executing
  // async feature setup. Invalidate it the moment the source changes.
  frame.contentWindow?.postMessage({type:'preview-invalidate',version:revision},'*');
  lastEditAt = performance.now();
  if (pendingSince === null) pendingSince = lastEditAt;
  if (busy && activeRenderController) {
    pendingRender = true;
    activeRenderController.abort();
  }
  updateFrontmatterSummary(editor.getValue());
  if (change.origin !== 'setValue') saveDraft();
  status.textContent = 'Waiting for edits…';
  schedulePreviewRender();
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
    source_anchors:data.source_anchors,
    source_offset:data.source_offset || 0,
    line_count:editor.lineCount(),
    feature_head_html:data.feature_head_html, body_end_html:data.body_end_html, version,
    attributes:attributes()}, '*');
}
function attributes() { return Object.fromEntries(appearanceKeys.map(k => [k, document.documentElement.getAttribute(k)])); }
function syncAppearance() { frame.contentWindow?.postMessage({type:'appearance',attributes:attributes()}, '*'); }
new MutationObserver(syncAppearance).observe(document.documentElement, {attributes:true,attributeFilter:appearanceKeys});
document.querySelector('#open').addEventListener('click', () => document.querySelector('#import').click());
document.querySelector('#import').addEventListener('change', async event => {
  const file = event.target.files[0];
  event.target.value = '';
  if (!file) return;
  if (file.size > 200000) {
    rejectedImportRevision = revision;
    errorBox.textContent = 'Document exceeds 200,000 UTF-8 bytes. Your current draft has been kept.';
    errorBox.hidden = false; return;
  }
  try {
    const text = await file.text();
    if (new TextEncoder().encode(text).length > 200000) throw new Error('Document exceeds 200,000 UTF-8 bytes.');
    window.MarkataSiteTree.import(text);
  } catch (error) { rejectedImportRevision = revision; errorBox.textContent = 'Cannot open file: ' + error.message; errorBox.hidden = false; }
});
if (!saved.textContent) saveDraft();

// Typing and restoring a local draft must not await workspace discovery,
// metadata indexing, LSP, or Go rendering.
requestAnimationFrame(() => {
  focusEditorForWriting();
  window.MarkataStartupStatus?.editorReady?.();
});
const workspaceReady = window.MarkataSiteTree.attach(editor, storageWarning, {saveToSite:true});
window.addEventListener('markata-site-index-ready', () => {
  if (!ready) return;
  revision++;
  frame.contentWindow?.postMessage({type:'preview-invalidate',version:revision},'*');
  activeRenderController?.abort();
  pendingRender = busy;
  clearTimeout(timer);
  timer = setTimeout(render, 0);
});
// Once the selected site file is attached, display *its* latest preview.
workspaceReady.then(() => {
  ready = true;
  window.MarkataStartupStatus?.workspaceReady();
  clearTimeout(timer);
  timer = setTimeout(render, 0);
}).catch(error => {
  errorBox.textContent='Unable to initialize workspace: '+error.message;
  errorBox.hidden=false;
});


workspaceReady.then(() => {
  window.MarkataFileFinder?.attach(editor);
  window.MarkataRecentFiles?.attach();
  window.MarkataConflictCompare?.attach(editor);
});
