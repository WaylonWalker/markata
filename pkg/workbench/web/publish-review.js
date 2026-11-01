/* Publish review entry points for the site-backed Go editor.
 * This stage previews last-committed vs saved content. It never deploys.
 */
(() => {
  'use strict';
  let dialog, editor, button, notice, panels, selected = 'after';
  const sourceLimit = 200000;
  const sleep = ms => new Promise(resolve => setTimeout(resolve, ms));
  const esc = value => String(value).replace(/[&<>"']/g, char => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[char]));
  const state = {path:null, data:null, mode:'side', loading:false};

  function titleFromSource(source) {
    const front = source.match(/^---\r?\n([\s\S]*?)\r?\n---/);
    return front?.[1].match(/^title:\s*["']?(.+?)["']?\s*$/m)?.[1] || null;
  }
  function frontmatter(source) {
    const front = source.match(/^---\r?\n([\s\S]*?)\r?\n---/);
    const fields = {};
    if (front) for (const line of front[1].split(/\r?\n/)) {
      const match = line.match(/^([a-z][\w-]*):\s*(.*?)\s*$/i);
      if (match) fields[match[1].toLowerCase()] = match[2].replace(/^["']|["']$/g,'');
    }
    return fields;
  }
  function contentBody(source) {
    return source.replace(/^---\r?\n[\s\S]*?\r?\n---(?:\r?\n|$)/,'');
  }
  function canOpen() {
    const path = window.MarkataSiteTree?.activeFile?.();
    if (!path || !/\.md$/i.test(path) || path === 'draft.md') return null;
    return path;
  }
  function setNotice(message, level='info') {
    notice.hidden = !message;
    notice.textContent = message || '';
    notice.dataset.level = level;
  }
  function setMode(mode) {
    state.mode = mode;
    for (const tab of dialog.querySelectorAll('[data-publish-mode]')) {
      const on = tab.dataset.publishMode === mode;
      tab.classList.toggle('active',on);
      tab.setAttribute('aria-pressed',String(on));
    }
    for (const section of dialog.querySelectorAll('[data-publish-view]'))
      section.hidden = section.dataset.publishView !== mode;
  }
  function setSide(side) {
    selected = side;
    for (const tab of dialog.querySelectorAll('[data-publish-side]')) {
      const on = tab.dataset.publishSide === side;
      tab.classList.toggle('active',on);
      tab.setAttribute('aria-pressed',String(on));
    }
    for (const pane of dialog.querySelectorAll('[data-publish-pane]'))
      pane.classList.toggle('mobile-visible',pane.dataset.publishPane === side);
  }
  function renderChips(data) {
    const before = frontmatter(data.baseline), after = frontmatter(data.current);
    const chips = dialog.querySelector('#publish-chips');
    chips.replaceChildren();
    function chip(message, risk='info') {
      const span = document.createElement('span'); span.className='publish-chip';
      span.dataset.level=risk; span.textContent=message; chips.append(span);
    }
    if (!data.baseline_exists) chip('New uncommitted page','warning');
    if (before.draft !== after.draft) chip('draft: '+(before.draft||'unset')+' → '+(after.draft||'unset'),'warning');
    if (before.published !== after.published) chip('published: '+(before.published||'unset')+' → '+(after.published||'unset'),'warning');
    if (before.slug !== after.slug) chip('Slug changed','warning');
    if ((after.draft==='true' && after.published==='true') || (after.draft==='false' && after.published==='false'))
      chip('Conflicting publication flags','warning');
    chip(data.changed ? 'Content changes detected' : 'No changes','info');
    chip('Compared with last Git commit','muted');
  }
  function diffLines(oldText,newText) {
    const a=oldText.split('\n'),b=newText.split('\n');
    if (a.length*b.length>180000) return [{type:'note',text:'Diff too large for an inline comparison; use the source files.'}];
    const dp=Array.from({length:a.length+1},()=>new Uint16Array(b.length+1));
    for(let i=a.length-1;i>=0;i--) for(let j=b.length-1;j>=0;j--)
      dp[i][j]=a[i]===b[j] ? 1+dp[i+1][j+1] : Math.max(dp[i+1][j],dp[i][j+1]);
    let i=0,j=0;const output=[];
    while(i<a.length||j<b.length) {
      if(i<a.length && j<b.length && a[i]===b[j]) {output.push({type:'same',text:a[i++]});j++;}
      else if(j<b.length && (i===a.length||dp[i][j+1]>=dp[i+1][j])) output.push({type:'add',text:b[j++]});
      else output.push({type:'del',text:a[i++]});
    }
    return output;
  }
  function renderDiff(data) {
    const box=dialog.querySelector('#publish-source-diff');
    box.replaceChildren();
    for (const line of diffLines(data.baseline,data.current)) {
      const row=document.createElement('div');
      row.className='publish-diff-line '+line.type;
      const prefix=document.createElement('span');
      prefix.className='publish-diff-prefix';
      prefix.textContent=line.type==='add'?'+':line.type==='del'?'-':' ';
      const content=document.createElement('span');
      content.textContent=line.text;
      row.append(prefix,content); box.append(row);
    }
  }
  async function renderVersion(source, frame) {
    const appearance={};
    for (const key of ['data-palette','data-theme','data-aesthetic','data-fontpack','data-text-size']) {
      const value=document.documentElement.getAttribute(key);
      if(value) appearance[key]=value;
    }
    const response=await fetch('/render',{method:'POST',headers:{'Content-Type':'application/json'},
      body:JSON.stringify({markdown:contentBody(source),include_document:true,appearance})});
    const result=await response.json();
    if(!response.ok||!result.document) throw new Error(result.error||'Could not render the comparison');
    frame.srcdoc=result.document;
  }
  async function awaitSaved(path, snapshot) {
    window.MarkataSiteTree.save();
    // The workbench writes edits asynchronously; never compare an older disk
    // version with a newer editor buffer and call it a publishable review.
    for(let n=0;n<18;n++) {
      const response=await fetch('/workspace/file?path='+encodeURIComponent(path),{cache:'no-store'});
      if(!response.ok) throw new Error('Could not read the saved workspace file.');
      const file=await response.json();
      if(file.source===snapshot) return;
      await sleep(200);
    }
    throw new Error('Saving is still in progress. Check the save status and try again.');
  }
  async function open() {
    if(state.loading||dialog.open)return;
    const path=canOpen();
    if(!path) {
      dialog.showModal();setNotice('This is a browser-only or non-Markdown file. Create or open a site Markdown post before reviewing publication.','warning');
      dialog.querySelector('#publish-title').textContent='Save to site first';
      dialog.querySelector('#publish-body').hidden=true;
      return;
    }
    state.loading=true;state.path=path;
    dialog.showModal();
    dialog.querySelector('#publish-body').hidden=true;
    dialog.querySelector('#publish-title').textContent='Preparing review…';
    setNotice('Saving the current editor text and reading the committed baseline.');
    try {
      const snapshot=editor.getValue();
      if(new TextEncoder().encode(snapshot).length>sourceLimit) throw new Error('Document is too large to review.');
      await awaitSaved(path,snapshot);
      const response=await fetch('/workspace/publish/review?path='+encodeURIComponent(path),{cache:'no-store'});
      const data=await response.json();
      if(!response.ok) throw new Error(data.error||'Review unavailable. This workspace may not be a Git checkout.');
      if(data.current!==snapshot || window.MarkataSiteTree?.activeFile?.()!==path || editor.getValue()!==snapshot)
        throw new Error('The document changed while preparing the review. Close and reopen to compare the latest text.');
      state.data=data;
      dialog.querySelector('#publish-title').textContent=titleFromSource(data.current)||path.split('/').pop();
      dialog.querySelector('#publish-path').textContent=path;
      dialog.querySelector('#publish-before-name').textContent=data.baseline_exists?'Last committed':'New page';
      renderChips(data);
      renderDiff(data);
      setSide('after');setMode('side');
      dialog.querySelector('#publish-body').hidden=false;
      setNotice('Review only. Last committed is not necessarily the version currently deployed; publishing and deployment are not connected yet.','info');
      await Promise.all([
        renderVersion(data.baseline,dialog.querySelector('#publish-before-frame')),
        renderVersion(data.current,dialog.querySelector('#publish-after-frame'))
      ]);
    }catch(error){
      setNotice(error.message,'warning');
      dialog.querySelector('#publish-body').hidden=true;
      dialog.querySelector('#publish-title').textContent='Review unavailable';
    }finally {state.loading=false;}
  }
  function attach(instance) {
    editor=instance;
    button=document.createElement('button');
    button.type='button';button.id='publish-open';
    button.textContent='Publish';button.title='Review & Publish (Ctrl/Cmd+Shift+Enter)';
    button.setAttribute('aria-haspopup','dialog');
    const settings=document.querySelector('#settings-open');
    settings.parentElement.insertBefore(button,settings);
    const doc=document.createElement('dialog');
    doc.id='publish-review';doc.setAttribute('aria-labelledby','publish-title');
    doc.innerHTML='<header class="publish-review-head"><div><small>REVIEW &amp; PUBLISH</small><h2 id="publish-title">Review changes</h2><p id="publish-path"></p></div><button type="button" id="publish-close" aria-label="Close publish review">✕</button></header>'+
    '<div id="publish-notice" role="status"></div><div id="publish-body" hidden><div id="publish-chips" aria-label="Change summary"></div>'+
    '<nav class="publish-mode-tabs" aria-label="Comparison view"><button type="button" data-publish-mode="side">Side by side</button><button type="button" data-publish-mode="source">Source diff</button></nav>'+
    '<div data-publish-view="side"><div class="publish-side-tabs" role="group" aria-label="Version to read"><button data-publish-side="before" type="button">Last committed</button><button data-publish-side="after" type="button">Current draft</button></div>'+
    '<div class="publish-comparison"><section data-publish-pane="before"><h3 id="publish-before-name">Last committed</h3><iframe id="publish-before-frame" title="Last committed rendering" sandbox="allow-scripts"></iframe></section>'+
    '<section data-publish-pane="after"><h3>Current saved draft</h3><iframe id="publish-after-frame" title="Current draft rendering" sandbox="allow-scripts"></iframe></section></div></div>'+
    '<div data-publish-view="source" hidden><div id="publish-source-diff" role="region" tabindex="0" aria-label="Unified Markdown source diff"></div></div></div>'+
    '<footer class="publish-actions"><span>Nothing is published by opening this review.</span><button type="button" id="publish-back">Back to editing</button><button type="button" id="publish-confirm" disabled title="Deployment backend is not yet configured">Publish · coming next</button></footer>';
    document.body.append(doc);dialog=doc;notice=doc.querySelector('#publish-notice');panels=doc.querySelector('#publish-body');
    button.addEventListener('click',open);
    doc.querySelector('#publish-close').addEventListener('click',()=>doc.close());
    doc.querySelector('#publish-back').addEventListener('click',()=>doc.close());
    doc.addEventListener('close',()=>editor.focus());
    doc.querySelectorAll('[data-publish-mode]').forEach(b=>b.addEventListener('click',()=>setMode(b.dataset.publishMode)));
    doc.querySelectorAll('[data-publish-side]').forEach(b=>b.addEventListener('click',()=>setSide(b.dataset.publishSide)));
    const gesture=doc.querySelector('.publish-comparison');
    let start=null;
    gesture.addEventListener('touchstart',e=>{start=[e.touches[0].clientX,e.touches[0].clientY]},{passive:true});
    gesture.addEventListener('touchend',e=>{
      if(!start)return;
      const dx=e.changedTouches[0].clientX-start[0],dy=e.changedTouches[0].clientY-start[1];
      if(Math.abs(dx)>55 && Math.abs(dx)>Math.abs(dy)*1.6)setSide(dx<0?'after':'before');
      start=null;
    },{passive:true});
    document.addEventListener('keydown',event=>{
      if(event.isComposing||event.altKey||!event.shiftKey||!(event.ctrlKey||event.metaKey)||event.key!=='Enter')return;
      if(document.querySelector('dialog[open]'))return;
      event.preventDefault();event.stopPropagation();open();
    },true);
  }
  window.MarkataPublishReview={attach,open};
})();