/* Publish review entry points for the site-backed Go editor.
 * This stage previews last-committed vs saved content. It never deploys.
 */
(() => {
  'use strict';
  let dialog, editor, button, notice, panels, selected = 'after';
  const sourceLimit = 200000;
  const sleep = ms => new Promise(resolve => setTimeout(resolve, ms));
  const esc = value => String(value).replace(/[&<>"']/g, char => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[char]));
  const state = {path:null, data:null, mode:'side', loading:false, commit:null, busy:false};

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
    const files = window.MarkataSiteTree?.files?.();
    if (!path || !/\.md$/i.test(path) || !Array.isArray(files) || !files.some(file => file.path === path)) return null;
    return path;
  }
  function setNotice(message, level='info') {
    notice.hidden = !message;
    notice.textContent = message || '';
    notice.dataset.level = level;
  }
  let pierreRequest;
  function loadPierre() {
    if(window.MarkataPierreDiff) return Promise.resolve();
    if(pierreRequest) return pierreRequest;
    pierreRequest=new Promise((resolve,reject)=>{
      const script=document.createElement('script');
      script.src='/editor/pierre-diffs.bundle.js';
      script.async=true;
      script.onload=()=>window.MarkataPierreDiff?resolve():reject(new Error('Pierre diff bundle did not initialize'));
      script.onerror=()=>reject(new Error('Optional Pierre bundle is not installed. Run the tools/publish-diffs build to enable Full diff.'));
      document.head.append(script);
    });
    return pierreRequest;
  }
  async function showPierre() {
    const error=dialog.querySelector('#publish-pierre-error');
    const host=dialog.querySelector('#publish-pierre-host');
    if(!state.data || state.mode!=='pierre')return;
    error.hidden=true;
    try{
      await loadPierre();
      if(state.mode!=='pierre' || !state.data)return;
      window.MarkataPierreDiff.render(host,state.data.baseline,state.data.current,{
        filename:state.path,newFile:!state.data.baseline_exists,
        layout:dialog.querySelector('[data-pierre-layout].active')?.dataset.pierreLayout||'stacked'
      });
    }catch(e){error.textContent=e.message;error.hidden=false;}
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
    if(mode==='pierre')showPierre();
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
    return result;
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
    state.loading=true;state.path=path;state.data=null;state.commit=null;
    dialog.querySelector('[data-publish-mode="rich"]').disabled=false;
    dialog.querySelector('#publish-confirm').hidden=false;
    dialog.querySelector('#publish-confirm').disabled=true;
    dialog.querySelector('#publish-push').hidden=true;
    dialog.querySelector('#publish-live-check').hidden=true;
    dialog.querySelector('#publish-mark-public').hidden=true;
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
      const fields=frontmatter(data.current);
      dialog.querySelector('#publish-mark-public').hidden=fields.published==='true' && fields.draft!=='true';
      dialog.querySelector('#publish-confirm').disabled=!data.changed || fields.published!=='true' || fields.draft==='true';
      setSide('after');setMode('side');
      dialog.querySelector('#publish-body').hidden=false;
      setNotice('Comparison against last Git commit, not necessarily the deployed site. Committing does not push or deploy.','info');
      const [before,after]=await Promise.all([
        renderVersion(data.baseline,dialog.querySelector('#publish-before-frame')),
        renderVersion(data.current,dialog.querySelector('#publish-after-frame'))
      ]);
      try {
        window.MarkataRichPublishDiff.render(before,after,dialog.querySelector('#publish-rich-frame'));
      } catch(error) {
        const tab=dialog.querySelector('[data-publish-mode="rich"]');
        tab.disabled=true;tab.title=error.message;
      }
    }catch(error){
      setNotice(error.message,'warning');
      dialog.querySelector('#publish-body').hidden=true;
      dialog.querySelector('#publish-title').textContent='Review unavailable';
    }finally {state.loading=false;}
  }
  function markPublicAndReviewAgain() {
    const source=editor.getValue();
    const newline=source.includes('\r\n')?'\r\n':'\n';
    let next=source;
    if (!source.startsWith('---'+newline)) {
      next='---'+newline+'published: true'+newline+'draft: false'+newline+'---'+newline+newline+source;
    } else {
      const end=source.indexOf(newline+'---',3+newline.length);
      if(end<0){setNotice('Fix malformed YAML frontmatter before publishing.','warning');return;}
      const before=source.slice(0,end);
      const rest=source.slice(end);
      let changes=before;
      if(/^published:\s*(true|false)\s*$/m.test(changes))
        changes=changes.replace(/^published:\s*(true|false)\s*$/m,'published: true');
      else changes+=newline+'published: true';
      if(/^draft:\s*(true|false)\s*$/m.test(changes))
        changes=changes.replace(/^draft:\s*(true|false)\s*$/m,'draft: false');
      next=changes+rest;
    }
    if(next===source)return;
    editor.operation(()=>editor.replaceRange(next,{line:0,ch:0},{line:editor.lastLine(),ch:editor.getLine(editor.lastLine()).length},'+publish-meta'));
    dialog.close();
    setTimeout(open,0);
  }
  async function commitReview() {
    const data=state.data;
    if(!data||state.busy||!data.changed)return;
    if(state.path!==window.MarkataSiteTree?.activeFile?.()||editor.getValue()!==data.current) {
      setNotice('The document changed. Close and reopen review before committing.','warning');return;
    }
    state.busy=true;
    const btn=dialog.querySelector('#publish-confirm');
    btn.disabled=true;btn.textContent='Creating commit…';
    try{
      const response=await fetch('/workspace/publish/commit',{
        method:'POST',headers:{'Content-Type':'application/json'},
        body:JSON.stringify({path:data.path,head:data.head,version:data.version})
      });
      const result=await response.json();
      if(!response.ok)throw new Error(result.error||'Cannot commit this revision.');
      state.commit=result;
      btn.hidden=true;
      dialog.querySelector('#publish-mark-public').hidden=true;
      dialog.querySelector('#publish-push').hidden=false;
      setNotice('Committed only this page to Git. Other staged files were excluded. Push separately to trigger your configured deployment pipeline.','info');
    }catch(err){setNotice(err.message,'warning');btn.disabled=false;}
    finally{state.busy=false;btn.textContent='Commit selected page';}
  }
  async function checkLiveDeployment() {
    if(!state.commit || state.busy)return;
    const btn=dialog.querySelector('#publish-live-check');
    state.busy=true;btn.disabled=true;btn.textContent='Checking…';
    try{
      const response=await fetch('/workspace/publish/deployment?commit='+encodeURIComponent(state.commit.commit),{cache:'no-store'});
      const result=await response.json();
      if(!response.ok)throw new Error(result.error||'Cannot verify deployment.');
      if(result.verified){
        btn.hidden=true;
        setNotice('Live · Verified source revision '+result.live_source_commit.slice(0,10)+(result.release_id?' · release '+result.release_id:'')+'.','info');
      } else {
        setNotice(result.message||'Deployment has not been verified.','warning');
      }
    }catch(error){setNotice(error.message,'warning');}
    finally{state.busy=false;btn.disabled=false;btn.textContent='Check live status';}
  }
  async function pushReview() {
    if(state.busy||!state.commit)return;
    const btn=dialog.querySelector('#publish-push');
    state.busy=true;btn.disabled=true;btn.textContent='Pushing…';
    try {
      const response=await fetch('/workspace/publish/push?path='+encodeURIComponent(state.path),{
        method:'POST',headers:{'Content-Type':'application/json'},
        body:JSON.stringify({commit:state.commit.commit,branch:state.commit.branch})
      });
      const result=await response.json();
      if(!response.ok)throw new Error(result.error||'Push failed. The local commit is retained.');
      btn.hidden=true;
      dialog.querySelector('#publish-live-check').hidden=false;
      setNotice('Pushed to origin. Click Check live status to verify the current builder release, if configured.','info');
    }catch(err){setNotice(err.message,'warning');btn.disabled=false;}
    finally{state.busy=false;btn.textContent='Push to origin';}
  }
  function updatePublishButton() {
    if (!button || !editor) return;
    const path = canOpen();
    const fields = frontmatter(editor.getValue());
    const published = fields.published === 'true' && fields.draft !== 'true';
    const conflicting = fields.draft === 'true' && fields.published === 'true';
    button.textContent = !path ? 'Save to site' : conflicting ? 'Review changes' : published ? 'Publish updates' : 'Publish';
    button.title = 'Review changes (Ctrl/Cmd+Shift+Enter)';
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
    '<nav class="publish-mode-tabs" aria-label="Comparison view"><button type="button" data-publish-mode="side">Side by side</button><button type="button" data-publish-mode="rich">Rendered diff</button><button type="button" data-publish-mode="source">Source diff</button><button type="button" data-publish-mode="pierre">Full diff</button></nav>'+
    '<div data-publish-view="side"><div class="publish-side-tabs" role="group" aria-label="Version to read"><button data-publish-side="before" type="button">Last committed</button><button data-publish-side="after" type="button">Current draft</button></div>'+
    '<div class="publish-comparison"><section data-publish-pane="before"><h3 id="publish-before-name">Last committed</h3><iframe id="publish-before-frame" title="Last committed rendering" sandbox="allow-scripts"></iframe></section>'+
    '<section data-publish-pane="after"><h3>Current saved draft</h3><iframe id="publish-after-frame" title="Current draft rendering" sandbox="allow-scripts"></iframe></section></div></div>'+
    '<div data-publish-view="rich" hidden><p class="publish-rich-help">Structural comparison: inline changes for text, before/after blocks for complex elements. Use Side by side for exact rendering.</p><iframe id="publish-rich-frame" title="Unified rendered Markdown change annotations" sandbox="allow-scripts"></iframe></div>'+
    '<div data-publish-view="pierre" hidden><div class="publish-pierre-actions"><span>Pierre source diff</span><button type="button" data-pierre-layout="stacked" class="active">Unified</button><button type="button" data-pierre-layout="split">Split</button></div><p id="publish-pierre-error" role="status" hidden></p><div id="publish-pierre-host"></div></div>'+
    '<div data-publish-view="source" hidden><div id="publish-source-diff" role="region" tabindex="0" aria-label="Unified Markdown source diff"></div></div></div>'+
    '<footer class="publish-actions"><span>A Git commit is not a verified deployment.</span><button type="button" id="publish-back">Back to editing</button><button type="button" id="publish-mark-public" hidden>Mark as public</button><button type="button" id="publish-confirm" disabled>Commit selected page</button><button type="button" id="publish-push" hidden>Push to origin</button><button type="button" id="publish-live-check" hidden>Check live status</button></footer>';
    document.body.append(doc);dialog=doc;notice=doc.querySelector('#publish-notice');panels=doc.querySelector('#publish-body');
    button.addEventListener('click',open);
    editor.on('change',updatePublishButton);
    const fileLabel = document.querySelector('#active-file');
    if (fileLabel) new MutationObserver(updatePublishButton).observe(fileLabel,{subtree:true,characterData:true,childList:true});
    updatePublishButton();
    doc.querySelector('#publish-close').addEventListener('click',()=>doc.close());
    doc.querySelector('#publish-back').addEventListener('click',()=>doc.close());
    doc.querySelector('#publish-mark-public').addEventListener('click',markPublicAndReviewAgain);
    doc.querySelector('#publish-confirm').addEventListener('click',commitReview);
    doc.querySelector('#publish-push').addEventListener('click',pushReview);
    doc.querySelector('#publish-live-check').addEventListener('click',checkLiveDeployment);
    doc.addEventListener('close',()=>{window.MarkataPierreDiff?.release?.(doc.querySelector('#publish-pierre-host'));editor.focus();});
    doc.querySelectorAll('[data-publish-mode]').forEach(b=>b.addEventListener('click',()=>setMode(b.dataset.publishMode)));
    doc.querySelectorAll('[data-pierre-layout]').forEach(b=>b.addEventListener('click',()=>{
      doc.querySelectorAll('[data-pierre-layout]').forEach(other=>other.classList.toggle('active',other===b));
      if(state.mode==='pierre')showPierre();
    }));
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