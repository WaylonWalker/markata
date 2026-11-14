/* Publish review entry points for the site-backed Go editor.
 * This stage previews last-committed vs saved content. It never deploys.
 */
(() => {
  'use strict';
  let dialog, editor, button, notice, panels, selected = 'after';
  const sourceLimit = 200000;
  const sleep = ms => new Promise(resolve => setTimeout(resolve, ms));
  const esc = value => String(value).replace(/[&<>"']/g, char => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[char]));
  const state = {path:null, data:null, mode:'rich', loading:false, commit:null, busy:false};
  let reviewGeneration = 0;
  let reviewRenderPromise = null;
  let reviewRenderController = null;
  const comparisonMaps = {before:[],after:[]};
  const comparisonMax = {before:0,after:0};

  function titleFromSource(source) {
    const front = source.match(/^---\r?\n([\s\S]*?)\r?\n---/);
    return front?.[1].match(/^title:\s*["']?(.+?)["']?\s*$/m)?.[1] || null;
  }
  function frontmatter(source) {
    const front = source.match(/^---\r?\n([\s\S]*?)\r?\n---/);
    const fields = {};
    if (front) for (const line of front[1].split(/\r?\n/)) {
      const match = line.match(/^([a-z][\w-]*):\s*(.*?)\s*$/i);
      if (match) {
        const key=match[1].toLowerCase(),value=match[2].replace(/^["']|["']$/g,'');
        fields[key]=['true','false'].includes(value.toLowerCase())?value.toLowerCase():value;
      }
    }
    return fields;
  }
  function contentBody(source) {
    return source.replace(/^---\r?\n[\s\S]*?\r?\n---(?:\r?\n|$)/,'');
  }
  function canOpen() {
    const path = window.MarkataSiteTree?.activeFile?.();
    if (!path || !/\.md$/i.test(path)) return null;
    // The site tree is compiled into a separate browser bundle. Do not
    // require a newly added source-only API to recognize existing files.
    // Only the known scratch aliases need special treatment here. The
    // server authoritatively validates paths on the review request.
    if (/^Unsaved draft(?: \d+)?\.md$/i.test(path)) return null;
    if (path === 'draft.md' && !window.MarkataSiteTree?.files?.()?.some(file => file.path === path)) return null;
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
      script.type='module';
      script.async=true;
      script.onload=()=>window.MarkataPierreDiff?resolve():reject(new Error('The full diff renderer did not initialize.'));
      script.onerror=()=>reject(new Error('Full diff assets could not load. Use Source diff, or rebuild tools/publish-diffs assets.'));
      document.head.append(script);
    }).catch(error => {
      // A single interrupted connection must not permanently poison the
      // optional renderer until the user reloads the entire editor.
      pierreRequest=null;
      throw error;
    });
    return pierreRequest;
  }
  async function showPierre() {
    const error=dialog.querySelector('#publish-pierre-error');
    const host=dialog.querySelector('#publish-pierre-host');
    if(!state.data || state.mode!=='pierre')return;
    error.hidden=true;
    host.dataset.state='loading';
    try{
      await loadPierre();
      if(state.mode!=='pierre' || !state.data)return;
      await window.MarkataPierreDiff.render(host,state.data.baseline,state.data.current,{
        filename:state.path,newFile:!state.data.baseline_exists,
        layout:dialog.querySelector('[data-pierre-layout].active')?.dataset.pierreLayout||'stacked'
      });
      if(state.mode!=='pierre'||!dialog.open)return;
      host.dataset.state='ready';
    }catch(e){
      host.replaceChildren();
      host.dataset.state='error';
      error.textContent='Full diff is unavailable: '+e.message+'. Source diff still works.';
      error.hidden=false;
    }
  }
  function setMode(mode) {
    state.mode = mode;
    for (const tab of dialog.querySelectorAll('[data-publish-mode]')) {
      const on = tab.dataset.publishMode === mode;
      tab.classList.toggle('active',on);
      tab.setAttribute('aria-selected',String(on));
      tab.tabIndex=on?0:-1;
    }
    for (const section of dialog.querySelectorAll('[data-publish-view]'))
      section.hidden = section.dataset.publishView !== mode;
    if(mode==='pierre')showPierre();
    if((mode==='rich'||mode==='side') && state.data?.changed) ensureRendered();
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
    const before=oldText.split('\n'), after=newText.split('\n');
    let start=0,oldEnd=before.length,newEnd=after.length;
    while(start<oldEnd && start<newEnd && before[start]===after[start]) start++;
    while(oldEnd>start && newEnd>start && before[oldEnd-1]===after[newEnd-1]){
      oldEnd--;newEnd--;
    }
    if(start===oldEnd && start===newEnd)return [];

    const oldMiddle=before.slice(start,oldEnd),newMiddle=after.slice(start,newEnd);
    const output=[];
    const addContext=(from,to,lines)=>{
      if(from>0)output.push({type:'note',text:`… ${from} unchanged lines …`});
      for(let i=from;i<to;i++)output.push({type:'same',text:lines[i]});
    };
    addContext(Math.max(0,start-3),start,before);
    // Compare only the changed middle; a one-line edit in a 2,000-line
    // document must not hit a quadratic whole-document guard.
    if(!oldMiddle.length) {
      for(const line of newMiddle)output.push({type:'add',text:line});
    }else if(!newMiddle.length){
      for(const line of oldMiddle)output.push({type:'del',text:line});
    }else if(oldMiddle.length*newMiddle.length<=180000){
      const dp=Array.from({length:oldMiddle.length+1},()=>new Uint16Array(newMiddle.length+1));
      for(let i=oldMiddle.length-1;i>=0;i--)
        for(let j=newMiddle.length-1;j>=0;j--)
          dp[i][j]=oldMiddle[i]===newMiddle[j]?1+dp[i+1][j+1]:Math.max(dp[i+1][j],dp[i][j+1]);
      let i=0,j=0;
      while(i<oldMiddle.length||j<newMiddle.length){
        if(i<oldMiddle.length&&j<newMiddle.length&&oldMiddle[i]===newMiddle[j]){
          output.push({type:'same',text:oldMiddle[i++]});j++;
        }else if(j<newMiddle.length&&(i===oldMiddle.length||dp[i][j+1]>=dp[i+1][j])){
          output.push({type:'add',text:newMiddle[j++]});
        }else output.push({type:'del',text:oldMiddle[i++]});
      }
    }else{
      // Do not block the UI or claim the entire document is too large.
      // Show bounded source excerpts with an explicit omitted-lines count.
      const limit=80;
      for(const line of oldMiddle.slice(0,limit))output.push({type:'del',text:line});
      if(oldMiddle.length>limit)output.push({type:'note',text:`… ${oldMiddle.length-limit} more removed lines not shown …`});
      for(const line of newMiddle.slice(0,limit))output.push({type:'add',text:line});
      if(newMiddle.length>limit)output.push({type:'note',text:`… ${newMiddle.length-limit} more added lines not shown …`});
    }
    const endContext=Math.min(after.length,newEnd+3);
    for(let i=newEnd;i<endContext;i++)output.push({type:'same',text:after[i]});
    if(endContext<after.length)output.push({type:'note',text:`… ${after.length-endContext} unchanged lines …`});
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
  function comparisonScrollBridge(side) {
    return "<script>(()=>{" +
      "const side=" + JSON.stringify(side) + ";let intentUntil=0,frame=0;" +
      "const content=document.querySelector('.post-content')||document.querySelector('main')||document.body;" +
      "const anchors=()=>{const seen=new Map();return [...content.querySelectorAll('h1,h2,h3,h4,h5,h6,p,li,blockquote,pre,table,figure')].map(node=>{const text=(node.textContent||'').replace(/\\s+/g,' ').trim().slice(0,180);if(!text)return null;const base=node.tagName.toLowerCase()+':'+text,count=seen.get(base)||0;seen.set(base,count+1);return{key:base+'#'+count,top:Math.round(node.getBoundingClientRect().top+scrollY)}}).filter(Boolean)};" +
      "let points=[];const sendMap=()=>{points=anchors();parent.postMessage({type:'publish-diff-map',side,points,max:Math.max(0,document.documentElement.scrollHeight-innerHeight),viewport:innerHeight},'*')};" +
      "const sendScroll=()=>parent.postMessage({type:'publish-diff-scroll',side,top:scrollY,max:Math.max(0,document.documentElement.scrollHeight-innerHeight),viewport:innerHeight,user:performance.now()<intentUntil},'*');" +
      "const scheduleMap=()=>{if(frame)return;frame=requestAnimationFrame(()=>{frame=0;sendMap();sendScroll()})};" +
      "for(const name of ['wheel','touchstart','pointerdown'])addEventListener(name,()=>{intentUntil=performance.now()+700},{passive:true});" +
      "addEventListener('keydown',event=>{if(['ArrowUp','ArrowDown','PageUp','PageDown','Home','End',' '].includes(event.key))intentUntil=performance.now()+700},true);" +
      "addEventListener('scroll',sendScroll,{passive:true});addEventListener('resize',scheduleMap);addEventListener('load',scheduleMap,true);" +
      "new ResizeObserver(scheduleMap).observe(content);" +
      "addEventListener('message',event=>{if(event.source!==parent||event.data?.type!=='publish-diff-scroll-to')return;scrollTo(0,Math.max(0,event.data.top||0))});scheduleMap()})()</script>";
  }
  function updateCommitMessage() {
    const input=dialog.querySelector('#publish-commit-message');
    const bytes=new TextEncoder().encode(input.value).length;
    dialog.querySelector('#publish-message-count').textContent=bytes+' / 4000 bytes';
    const fields=state.data?frontmatter(state.data.current):{};
    const canCommit=Boolean(state.data?.changed&&fields.published==='true'&&fields.draft!=='true');
    dialog.querySelector('#publish-confirm').disabled=!canCommit||!input.value.trim()||bytes>4000;
  }
  function comparisonScrollSync(event) {
    const before=dialog.querySelector('#publish-before-frame');
    const after=dialog.querySelector('#publish-after-frame');
    const source=event.source===before.contentWindow?'before':event.source===after.contentWindow?'after':null;
    if(!source)return;
    if(event.data?.type==='publish-diff-map'&&event.data.side===source) {
      comparisonMaps[source]=Array.isArray(event.data.points)?event.data.points:[];
      comparisonMax[source]=event.data.max||0;
      return;
    }
    if(event.data?.type!=='publish-diff-scroll'||event.data.side!==source)return;
    comparisonMax[source]=event.data.max||0;
    if(!event.data.user||state.mode!=='side')return;
    const target=source==='before'?'after':'before';
    const targetFrame=target==='before'?before:after;
    const targetPoints=comparisonMaps[target],byKey=new Map(targetPoints.map(point=>[point.key,point]));
    const reference=event.data.top+event.data.viewport*.3;
    let match=null,distance=Infinity;
    for(const point of comparisonMaps[source]) {
      const candidate=byKey.get(point.key);
      if(candidate&&Math.abs(point.top-reference)<distance){match={point,candidate};distance=Math.abs(point.top-reference);}
    }
    const targetTop=match
      ? match.candidate.top-(match.point.top-event.data.top)
      : (event.data.max?event.data.top/event.data.max:0)*comparisonMax[target];
    targetFrame.contentWindow?.postMessage({type:'publish-diff-scroll-to',top:targetTop},'*');
  }
  async function renderVersion(source, frame, side, signal) {
    const appearance={};
    for (const key of ['data-palette','data-theme','data-aesthetic','data-fontpack','data-text-size']) {
      const value=document.documentElement.getAttribute(key);
      if(value) appearance[key]=value;
    }
    const response=await fetch('/render',{method:'POST',headers:{'Content-Type':'application/json'},
      body:JSON.stringify({markdown:contentBody(source),include_document:true,appearance}),signal});
    const result=await response.json();
    if(!response.ok||!result.document) throw new Error(result.error||'Could not render the comparison');
    // Do not mutate an iframe until the caller verifies this is still the
    // active review. Responses from closed dialogs must not resurrect old HTML.
    return {result, document:result.document.replace('</body>',`${comparisonScrollBridge(side)}</body>`)};
  }
  function ensureRendered() {
    if(reviewRenderPromise || !state.data?.changed)return;
    const generation=reviewGeneration, data=state.data;
    reviewRenderController=new AbortController();
    const signal=reviewRenderController.signal;
    setNotice('Source diff is ready. Rendering the comparison on demand…');
    const loadSide=(source,id,side)=>renderVersion(source,dialog.querySelector('#'+id),side,signal)
      .then(version=>{
        if(generation===reviewGeneration && dialog.open)
          dialog.querySelector('#'+id).srcdoc=version.document;
        return version;
      });
    // Let Side by side show whichever version finishes first. The richer
    // annotations wait for both, but a slow Go render no longer blanks
    // the entire comparison meanwhile.
    reviewRenderPromise=Promise.all([
      loadSide(data.baseline,'publish-before-frame','before'),
      loadSide(data.current,'publish-after-frame','after')
    ]).then(([before,after])=>{
      if(generation!==reviewGeneration || !dialog.open)return;
      try {
        window.MarkataRichPublishDiff.render(before.result,after.result,dialog.querySelector('#publish-rich-frame'));
        setNotice('Compared with the last Git commit. Rendered and source views are ready; nothing has been pushed or deployed.');
      }catch(error) {
        const tab=dialog.querySelector('[data-publish-mode="rich"]');
        tab.disabled=true;tab.title=error.message;
        if(state.mode==='rich')setMode('side');
        setNotice('Rendered annotations unavailable: '+error.message+'. Source diff and side-by-side remain available.','warning');
      }
    }).catch(error=>{
      if(generation!==reviewGeneration || !dialog.open)return;
      dialog.querySelector('[data-publish-mode="rich"]').disabled=true;
      dialog.querySelector('[data-publish-mode="side"]').disabled=true;
      if(state.mode==='rich'||state.mode==='side')setMode('source');
      setNotice('Rendered preview unavailable: '+error.message+'. Source diff is still ready.','warning');
    });
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
      dialog.dataset.unavailable='true';
      dialog.querySelector('#publish-confirm').disabled=true;
      const selected=window.MarkataSiteTree?.activeFile?.();
      const message=!selected
        ? 'No site document is selected. Open a Markdown file from the sidebar before reviewing changes.'
        : !/\.md$/i.test(selected)
        ? 'Review changes only supports Markdown (.md) files. Open a Markdown page from the sidebar.'
        : 'This draft is only in your browser. Use New Post to save it as a site Markdown file, then review its changes.';
      dialog.showModal();setNotice(message,'warning');
      dialog.querySelector('#publish-title').textContent='Choose a site Markdown file';
      dialog.querySelector('#publish-body').hidden=true;
      return;
    }
    dialog.dataset.unavailable='false';
    const generation=++reviewGeneration;
    state.loading=true;state.path=path;state.data=null;state.commit=null;
    reviewRenderController?.abort();
    reviewRenderController=null;
    reviewRenderPromise=null;
    for (const id of ['publish-before-frame','publish-after-frame','publish-rich-frame']) {
      dialog.querySelector('#'+id).srcdoc='<!doctype html><html><body><p>Preparing current comparison…</p></body></html>';
    }
    dialog.querySelector('[data-publish-mode="rich"]').disabled=false;
    dialog.querySelector('#publish-confirm').hidden=false;
    dialog.querySelector('#publish-confirm').disabled=true;
    dialog.querySelector('#publish-discard').hidden=true;
    dialog.querySelector('#publish-discard-warning').hidden=true;
    dialog.querySelector('#publish-discard-confirmation').hidden=true;
    dialog.querySelector('#publish-no-changes').hidden=true;
    dialog.querySelector('.publish-mode-tabs').hidden=false;
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
      // A closed review must not finish populating an obsolete dialog.
      if(generation!==reviewGeneration || !dialog.open) return;
      const response=await fetch('/workspace/publish/review?path='+encodeURIComponent(path),{cache:'no-store'});
      const data=await response.json();
      if(generation!==reviewGeneration || !dialog.open) return;
      if(!response.ok) throw new Error(data.error||'Review unavailable. This workspace may not be a Git checkout.');
      if(data.current!==snapshot || window.MarkataSiteTree?.activeFile?.()!==path || editor.getValue()!==snapshot)
        throw new Error('The document changed while preparing the review. Close and reopen to compare the latest text.');
      state.data=data;
      dialog.querySelector('#publish-discard').hidden=!data.changed;
      dialog.querySelector('#publish-discard').disabled=false;
      dialog.querySelector('#publish-discard-warning').hidden=true;
      dialog.querySelector('#publish-title').textContent=titleFromSource(data.current)||path.split('/').pop();
      dialog.querySelector('#publish-path').textContent=path;
      dialog.querySelector('#publish-before-name').textContent=data.baseline_exists?'Last committed':'New page';
      renderChips(data);
      renderDiff(data);
      const fields=frontmatter(data.current);
      dialog.querySelector('#publish-mark-public').hidden=fields.published==='true' && fields.draft!=='true';
      dialog.querySelector('#publish-confirm').disabled=!data.changed || fields.published!=='true' || fields.draft==='true';
      dialog.querySelector('#publish-commit-message').value=`Publish: ${path}`;
      updateCommitMessage();
      comparisonMaps.before=[];comparisonMaps.after=[];
      setSide('after');
      setMode('source');
      dialog.querySelector('#publish-body').hidden=false;
      const noChanges=!data.changed;
      dialog.querySelector('#publish-no-changes').hidden=!noChanges;
      dialog.querySelector('.publish-mode-tabs').hidden=noChanges;
      for(const panel of dialog.querySelectorAll('[data-publish-view]'))if(noChanges)panel.hidden=true;
      if(noChanges){
        setNotice('Your saved file matches the last Git commit. No changes need review or commit.','info');
      }else{
        setNotice('Source diff ready. Choose Rendered diff or Side by side to load the full page render. A commit alone does not deploy.','info');
      }
    }catch(error){
      if(generation!==reviewGeneration || !dialog.open) return;
      setNotice(error.message,'warning');
      dialog.querySelector('#publish-body').hidden=true;
      dialog.dataset.unavailable='true';
      dialog.querySelector('#publish-title').textContent='Review unavailable';
    }finally {
      // A previous request must not clear the loading flag of a newer review.
      if(generation===reviewGeneration) state.loading=false;
    }
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
      if(/^published:\s*(true|false)\s*$/im.test(changes))
        changes=changes.replace(/^published:\s*(true|false)\s*$/im,'published: true');
      else changes+=newline+'published: true';
      if(/^draft:\s*(true|false)\s*$/im.test(changes))
        changes=changes.replace(/^draft:\s*(true|false)\s*$/im,'draft: false');
      next=changes+rest;
    }
    if(next===source)return;
    editor.operation(()=>editor.replaceRange(next,{line:0,ch:0},{line:editor.lastLine(),ch:editor.getLine(editor.lastLine()).length},'+publish-meta'));
    dialog.close();
    setTimeout(open,0);
  }
  async function discardReviewed() {
    const data=state.data;
    if (!data || state.busy) return;
    if (state.path!==window.MarkataSiteTree?.activeFile?.() || editor.getValue()!==data.current) {
      setNotice('Draft changed since review. Reopen Review changes before discarding.','warning');return;
    }
    const button=dialog.querySelector('#publish-discard-accept');
    state.busy=true;
    button.disabled=true;
    button.textContent='Restoring…';
    try {
      // Save timers may still be pending: read exact source before requesting
      // restoration; the backend also rejects a stale source/version.
      const response=await fetch('/workspace/publish/discard',{
        method:'POST', headers:{'Content-Type':'application/json'},
        body:JSON.stringify({path:data.path,head:data.head,version:data.version})
      });
      const result=await response.json();
      if (!response.ok)throw new Error(result.error||'Could not discard reviewed changes.');
      if (result.status === 'removed-local') {
        // The server verified Git's untracked status and exact source version.
        // The browser must not re-create the deleted draft via recovery or autosave.
        if (!window.MarkataSiteTree?.forgetDeleted?.(data.path))
          throw new Error('New file deleted on disk, but browser recovery could not be cleared. Reload to reconcile.');
        refreshWorkspaceChanges();
        location.reload();
        return;
      }
      if (result.status !== 'restored-local')
        throw new Error('Unexpected discard response. Reload to reconcile.');
      // Set the live editor and its recovery copy to the exact server-acknowledged
      // baseline. Reloading here would restore an older dirty browser buffer.
      if (window.MarkataSiteTree?.activeFile?.()!==data.path || editor.getValue()!==data.current) {
        setNotice('The server restored the reviewed version, but this editor changed during the request. Reload the file to reconcile before editing.','warning');
        return;
      }
      editor.setValue(data.baseline);
      if (!window.MarkataSiteTree?.reconcileRestored?.(data.path,data.baseline,result.version))
        throw new Error('Restore succeeded on disk but browser recovery could not be reconciled. Reload this file before editing.');
      // Verify the real file after the server acknowledges it, rather than
      // calling the ordinary Save flow with the pre-discard source version.
      const check=await fetch('/workspace/file?path='+encodeURIComponent(data.path),{cache:'no-store'});
      const saved=await check.json();
      if (!check.ok || saved.source!==data.baseline || saved.version!==result.version)
        throw new Error('Disk changed again after discard. Reopen the file and resolve the new version.');
      refreshChangeGutter();
      refreshWorkspaceChanges();
      dialog.querySelector('#publish-discard-confirmation').hidden=true;
      dialog.close();
    } catch(error) {
      setNotice(error.message,'warning');
      dialog.querySelector('#publish-discard-confirmation').hidden=true;
    } finally {
      state.busy=false;button.disabled=false;button.textContent='Discard reviewed changes';
    }
  }
  async function commitReview() {
    const data=state.data;
    if(!data||state.busy||!data.changed)return;
    const messageField=dialog.querySelector('#publish-commit-message');
    const message=messageField.value.trim();
    if(messageField.value!==message){messageField.value=message;updateCommitMessage();}
    if(!message||new TextEncoder().encode(message).length>4000) {
      setNotice('Enter a commit message of 1 to 4000 bytes before committing.','warning');return;
    }
    if(state.path!==window.MarkataSiteTree?.activeFile?.()||editor.getValue()!==data.current) {
      setNotice('The document changed. Close and reopen review before committing.','warning');return;
    }
    state.busy=true;
    const btn=dialog.querySelector('#publish-confirm');
    btn.disabled=true;btn.textContent='Creating commit…';
    try{
      const response=await fetch('/workspace/publish/commit',{
        method:'POST',headers:{'Content-Type':'application/json'},
        body:JSON.stringify({path:data.path,head:data.head,version:data.version,message})
      });
      const result=await response.json();
      if(!response.ok)throw new Error(result.error||'Cannot commit this revision.');
      state.commit=result;
      refreshChangeGutter();
      refreshWorkspaceChanges();
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
      refreshWorkspaceChanges();
      dialog.querySelector('#publish-live-check').hidden=false;
      setNotice('Pushed to origin. Click Check live status to verify the current builder release, if configured.','info');
    }catch(err){setNotice(err.message,'warning');btn.disabled=false;}
    finally{state.busy=false;btn.textContent='Push to origin';}
  }

  // Only an explicit Publish request changes public metadata; ordinary
  // autosave and draft editing never change a post's publication state.
  function publicPostSource(source) {
    const newline = source.includes('\r\n') ? '\r\n' : '\n';
    if (!source.startsWith('---' + newline))
      return '---' + newline + 'published: true' + newline + 'draft: false' +
        newline + '---' + newline + newline + source;
    const end = source.indexOf(newline + '---', 3 + newline.length);
    if (end < 0) throw new Error('Fix malformed YAML frontmatter before publishing.');
    let head = source.slice(0,end);
    const tail = source.slice(end);
    if (/^published:\s*(true|false)\s*$/im.test(head))
      head = head.replace(/^published:\s*(true|false)\s*$/im, 'published: true');
    else head += newline + 'published: true';
    if (/^draft:\s*(true|false)\s*$/im.test(head))
      head = head.replace(/^draft:\s*(true|false)\s*$/im, 'draft: false');
    return head + tail;
  }

  // Writer-first publication keeps the Git transaction behind one explicit
  // action, while preserving the guarded review/commit/push endpoints.
  let quickBusy = false;
  let pendingQuick = null;
  let publishStatus = null;
  function reportQuick(message, type='info') {
    if (!publishStatus) return;
    publishStatus.textContent = message;
    publishStatus.dataset.state = type;
    publishStatus.title = message;
  }
  async function finishQuickPush(pending) {
    reportQuick('Publishing to your site…', 'loading');
    const response = await fetch('/workspace/publish/push?path=' +
      encodeURIComponent(pending.path), {
      method:'POST', headers:{'Content-Type':'application/json'},
      body:JSON.stringify({commit:pending.commit,branch:pending.branch})
    });
    const result = await response.json();
    if (!response.ok) throw new Error(result.error || 'Could not push. Your local commit is retained.');
    pendingQuick = null;
    reportQuick('Published to Git · site build not yet verified', 'info');
    refreshWorkspaceChanges();
    refreshChangeGutter();
  }
  async function quickPublish() {
    if (quickBusy || state.busy) return;
    const path = canOpen();
    if (!path) {
      reportQuick('Choose a site Markdown post before publishing.', 'warning');
      return;
    }
    quickBusy = true;
    button.disabled = true;
    try {
      let snapshot = editor.getValue();
      if (new TextEncoder().encode(snapshot).length > sourceLimit)
        throw new Error('Post is too large to publish.');
      if (pendingQuick) {
        if (pendingQuick.path !== path || pendingQuick.source !== snapshot)
          throw new Error('An earlier publish is pending. Return to that post without changing it and retry before publishing new edits.');
        await finishQuickPush(pendingQuick);
        return;
      }
      const meta = frontmatter(snapshot);
      if (meta.published !== 'true' || meta.draft === 'true') {
        // Publishing is the writer's explicit request to mark this post public.
        // The commit endpoint will still reject invalid or duplicate YAML.
        const next = publicPostSource(snapshot);
        if (new TextEncoder().encode(next).length > sourceLimit)
          throw new Error('Post is too large to publish.');
        editor.operation(() => editor.replaceRange(next, {line:0,ch:0},
          {line:editor.lastLine(),ch:editor.getLine(editor.lastLine()).length}, '+publish-meta'));
        snapshot = editor.getValue();
      }
      reportQuick('Saving current post…', 'loading');
      await awaitSaved(path, snapshot);
      if (window.MarkataSiteTree?.activeFile?.() !== path || editor.getValue() !== snapshot)
        throw new Error('Post changed during save. Try Publish again.');
      const reviewResponse = await fetch('/workspace/publish/review?path=' + encodeURIComponent(path),
        {cache:'no-store'});
      const reviewData = await reviewResponse.json();
      if (!reviewResponse.ok)
        throw new Error(reviewData.error || 'Could not check current post.');
      if (!reviewData.changed) {
        reportQuick('No unpublished content changes. Git does not confirm your site is live.');
        refreshWorkspaceChanges();
        return;
      }
      if (reviewData.current !== snapshot)
        throw new Error('Saved post changed during publication. Retry after reviewing the latest version.');
      reportQuick('Checking publishing safety…', 'loading');
      const readyResponse = await fetch('/workspace/publish/readiness?path=' +
        encodeURIComponent(path), {cache:'no-store'});
      const ready = await readyResponse.json();
      if (!readyResponse.ok || !ready.ready)
        throw new Error(ready.error || 'Publishing safety check failed.');
      if (ready.head !== reviewData.head)
        throw new Error('Git changed during publishing. Please retry after syncing your branch.');
      reportQuick('Publishing this post…', 'loading');
      const commitResponse = await fetch('/workspace/publish/commit', {
        method:'POST', headers:{'Content-Type':'application/json'},
        body:JSON.stringify({
          path,head:reviewData.head,version:reviewData.version,message:'Publish: '+path
        })
      });
      const commit = await commitResponse.json();
      if (!commitResponse.ok)
        throw new Error(commit.error || 'Could not prepare publication.');
      pendingQuick = {path,source:snapshot,commit:commit.commit,branch:commit.branch};
      refreshWorkspaceChanges();
      await finishQuickPush(pendingQuick);
    } catch (error) {
      reportQuick('Not published: ' + error.message, 'warning');
    } finally {
      quickBusy = false;
      button.disabled = false;
      updatePublishButton();
    }
  }

  function updatePublishButton() {
    if (!button || !editor) return;
    const path = canOpen();
    const label = pendingQuick ? 'Retry publish' : (!path ? 'Choose post' : 'Publish post');
    const icon=document.createElementNS('http://www.w3.org/2000/svg','svg');
    icon.setAttribute('viewBox','0 0 24 24');icon.setAttribute('aria-hidden','true');icon.setAttribute('focusable','false');
    const pathNode=document.createElementNS('http://www.w3.org/2000/svg','path');
    pathNode.setAttribute('d','M12 15V3m0 0 4 4m-4-4L8 7M4 14v6h16v-6');icon.append(pathNode);
    button.replaceChildren(icon,document.createTextNode(label));
    button.title = 'Publish this post only (Ctrl/Cmd+S). Review changes with Ctrl/Cmd+Shift+Enter.';
  }

  // Git-style gutter. Compare the visible editor to Git HEAD, not autosave.
  let gutterPath=null, gutterBase=null, gutterRequest=0, gutterTimer=null, gutterStatus=null;
  function reportGutter(message, kind='info') {
    if (!gutterStatus) return;
    gutterStatus.textContent=message;
    gutterStatus.dataset.state=kind;
    gutterStatus.title=message;
  }
  function clearChangeGutter() {
    editor?.clearGutter('markata-change-gutter');
  }
  function markChange(line, type) {
    if (line<0 || line>=editor.lineCount()) return;
    const marker=document.createElement('span');
    marker.className='markata-change-marker '+type;
    marker.setAttribute('aria-label',type==='add'?'Added line':type==='del'?'Deleted lines above':'Modified line');
    marker.title=marker.getAttribute('aria-label')+' since last Git commit';
    editor.setGutterMarker(line,'markata-change-gutter',marker);
  }
  // Trim unchanged leading/trailing lines before invoking the bounded LCS.
  // Large posts with a small edit should not silently lose their markers.
  // For a huge changed middle region, mark it conservatively as one region
  // rather than blocking the editor with quadratic diff work.
  function gutterChanges(before, after) {
    const oldLines=before.split('\n'), newLines=after.split('\n');
    let prefix=0;
    while(prefix<oldLines.length && prefix<newLines.length &&
          oldLines[prefix]===newLines[prefix])prefix++;
    let suffix=0;
    while(suffix<oldLines.length-prefix && suffix<newLines.length-prefix &&
          oldLines[oldLines.length-1-suffix]===newLines[newLines.length-1-suffix])suffix++;
    const oldMiddle=oldLines.slice(prefix,oldLines.length-suffix);
    const newMiddle=newLines.slice(prefix,newLines.length-suffix);
    if(!oldMiddle.length && !newMiddle.length)return [];
    if(oldMiddle.length*newMiddle.length>180000) {
      const type=oldMiddle.length?(newMiddle.length?'modify':'del'):'add';
      return [{line:prefix,added:newMiddle.length,deleted:oldMiddle.length,type}];
    }
    const diff=diffLines(oldMiddle.join('\n'),newMiddle.join('\n'));
    const changes=[];let offset=prefix;
    // Empty middles are represented as zero lines by the diff algorithm.
    if(!oldMiddle.length)return [{line:prefix,added:newMiddle.length,deleted:0,type:'add'}];
    if(!newMiddle.length)return [{line:prefix,added:0,deleted:oldMiddle.length,type:'del'}];
    for(let i=0;i<diff.length;) {
      if(diff[i].type==='same'){offset++;i++;continue;}
      const line=offset;let added=0,deleted=0;
      while(i<diff.length&&diff[i].type!=='same'){
        if(diff[i].type==='add'){added++;offset++;}
        if(diff[i].type==='del')deleted++;
        i++;
      }
      changes.push({line,added,deleted,type:deleted?(added?'modify':'del'):'add'});
    }
    return changes;
  }
  function drawChangeGutter() {
    clearChangeGutter();
    if(gutterPath!==canOpen() || typeof gutterBase!=='string')return;
    const changes=gutterChanges(gutterBase,editor.getValue());
    let marked=0;
    editor.operation(()=>{
      for(const change of changes) {
        // CodeMirror only renders visible lines; cap large replaced regions.
        // Their leading and trailing markers still identify the change.
        const limit=Math.min(change.added,600);
        for(let n=0;n<limit;n++){markChange(change.line+n,change.type);marked++;}
        if(change.added>limit)markChange(change.line+change.added-1,change.type);
        if(change.deleted&&!change.added){markChange(Math.min(change.line,editor.lineCount()-1),'del');marked++;}
      }
    });
    reportGutter(changes.length?(changes.length+' Git change region'+(changes.length===1?'':'s')):'No Git changes',changes.length?'changed':'clean');
  }
  async function refreshChangeGutter() {
    const path=canOpen(), request=++gutterRequest;
    gutterPath=path;gutterBase=null;clearChangeGutter();
    if(!path){reportGutter('Git gutter: select a site Markdown file','muted');return;}
    reportGutter('Loading Git change gutter…','loading');
    try {
      const response=await fetch('/workspace/publish/review?path='+encodeURIComponent(path),{cache:'no-store'});
      if(!response.ok) {
        reportGutter('Git gutter unavailable: review returned HTTP '+response.status,'warning');
        return;
      }
      const data=await response.json();
      if(request!==gutterRequest||canOpen()!==path)return;
      if(typeof data.baseline!=='string')throw new Error('missing committed baseline');
      gutterBase=data.baseline_exists?data.baseline:'';
      drawChangeGutter();
    } catch(error) {
      if(request===gutterRequest)reportGutter('Git gutter unavailable: '+error.message,'warning');
    }
  }
  function scheduleChangeGutter() {
    clearTimeout(gutterTimer);
    const path=canOpen();
    if(path!==gutterPath){refreshChangeGutter();return;}
    gutterTimer=setTimeout(drawChangeGutter,160);
  }


  // The Changes panel is intentionally separate from per-page Review. Git
  // status must be the inventory source, never just the set of opened files.
  let changesPanel = null, changesButton = null, changesList = null;
  let changesGeneration = 0, changesTimer = null;
  const changesExpanded = new Set();
  function changesVisible(open) {
    if (!changesPanel) return;
    changesPanel.hidden = !open;
    changesButton.setAttribute('aria-expanded', String(open));
    clearInterval(changesTimer);
    changesTimer = null;
    try { localStorage.setItem('markata-workspace-changes-open', open ? 'true' : 'false'); }
    catch { /* Keep the panel usable without storage. */ }
    if (open) {
      refreshWorkspaceChanges();
      changesTimer = setInterval(() => {
        if (document.visibilityState === 'visible') refreshWorkspaceChanges();
      }, 15000);
    } else {
      changesGeneration++;
    }
  }
  function changeLabel(entry) {
    if (entry.untracked) return 'New / untracked';
    if (entry.kind === 'conflict') return 'Git conflict';
    if (entry.kind === 'deleted') return 'Deleted';
    if (entry.kind === 'renamed') return 'Renamed';
    if (entry.kind === 'copied') return 'Copied';
    if (entry.kind === 'added') return 'Added';
    return 'Modified';
  }
  function showWorkspaceChanges(data) {
    changesList.replaceChildren();
    changesButton.textContent = 'Changes (' + data.total + ')';
    const heading = changesPanel.querySelector('#workspace-changes-summary');
    heading.textContent = data.total
      ? data.total + ' files with changes'
      : 'All caught up — no working-tree changes. Unpushed commits may still exist.';
    for (const change of data.entries) {
      const row = document.createElement('div');
      row.className = 'workspace-change';
      row.setAttribute('role', 'listitem');
      const description = document.createElement('div');
      description.className = 'workspace-change-description';
      const name = document.createElement('strong');
      name.textContent = change.path;
      name.title = change.path;
      const detail = document.createElement('small');
      detail.textContent = changeLabel(change) +
        (change.staged ? ' · staged' : '') +
        (change.unstaged ? ' · unstaged' : '') +
        (change.previous_path ? ' · from ' + change.previous_path : '') +
        (!change.workspace ? ' · outside this site' : !change.markdown ? ' · not Markdown' : '');
      description.append(name, detail);
      row.append(description);
      const canInspect = change.workspace && change.markdown && change.editor_path &&
        change.kind !== 'deleted' && change.kind !== 'conflict' &&
        window.MarkataSiteTree?.files?.()?.some(file => file.path === change.editor_path);
      if (canInspect) {
        const disclosure = document.createElement('button');
        disclosure.type = 'button';
        disclosure.textContent = 'View diff';
        const inline = document.createElement('div');
        inline.className = 'workspace-change-inline';
        inline.hidden = true;
        disclosure.setAttribute('aria-expanded', 'false');
        disclosure.addEventListener('click', async () => {
          const show = inline.hidden;
          inline.hidden = !show;
          disclosure.setAttribute('aria-expanded', String(show));
          disclosure.textContent = show ? 'Hide diff' : 'View diff';
          if (!show) return;
          inline.textContent = 'Loading comparison…';
          disclosure.disabled = true;
          try {
            const response = await fetch('/workspace/publish/review?path=' +
              encodeURIComponent(change.editor_path), {cache:'no-store'});
            const comparison = await response.json();
            if (!response.ok || typeof comparison.current !== 'string' ||
                typeof comparison.baseline !== 'string')
              throw new Error(comparison.error || 'Comparison unavailable.');
            inline.replaceChildren();
            const diff = document.createElement('div');
            diff.className = 'workspace-change-diff';
            diff.setAttribute('role', 'region');
            diff.setAttribute('aria-label', 'Diff for ' + change.path);
            const lines = diffLines(comparison.baseline, comparison.current);
            for (const line of lines.slice(0,160)) {
              const item = document.createElement('div');
              item.className = 'publish-diff-line ' + line.type;
              item.textContent = (line.type === 'add' ? '+' : line.type === 'del' ? '-' : ' ') + ' ' + line.text;
              diff.append(item);
            }
            if (!lines.length) diff.textContent = 'No source differences from the last commit.';
            if (lines.length > 160) {
              const more = document.createElement('p');
              more.textContent = 'Showing 160 of ' + lines.length + ' changed and context lines.';
              diff.append(more);
            }
            inline.append(diff);
            const edit = document.createElement('button');
            edit.type = 'button';
            edit.textContent = 'Edit post';
            edit.addEventListener('click', () => window.MarkataSiteTree.openFile(change.editor_path));
            inline.append(edit);
          } catch (error) {
            inline.textContent = 'Could not load diff: ' + error.message;
          } finally {
            disclosure.disabled = false;
          }
        });
        row.append(disclosure);
        changesList.append(row, inline);
      } else {
        const hint = document.createElement('small');
        hint.className = 'workspace-change-unavailable';
        hint.textContent = change.kind === 'deleted' ? 'Deleted' :
          change.kind === 'conflict' ? 'Resolve conflict' : 'Outside editor';
        row.append(hint);
        changesList.append(row);
      }
    }
  }
  async function refreshWorkspaceChanges() {
    if (!changesPanel || changesPanel.hidden) return;
    const generation = ++changesGeneration;
    const summary = changesPanel.querySelector('#workspace-changes-summary');
    if (!changesList.children.length) summary.textContent = 'Checking for changes…';
    try {
      const response = await fetch('/workspace/changes', {cache:'no-store'});
      const data = await response.json();
      if (generation !== changesGeneration || changesPanel.hidden) return;
      if (!response.ok || !Array.isArray(data.entries)) throw new Error(data.error || 'Change inventory unavailable');
      showWorkspaceChanges(data);
    } catch (error) {
      if (generation !== changesGeneration || changesPanel.hidden) return;
      summary.textContent = 'Cannot list all changes: ' + error.message;
      changesButton.textContent = 'Changes (?)';
      changesList.replaceChildren();
    }
  }
  function attachWorkspaceChanges() {
    const settings = document.querySelector('#settings-open');
    changesButton = document.createElement('button');
    changesButton.type = 'button';
    changesButton.id = 'workspace-changes-toggle';
    changesButton.textContent = 'Changes';
    changesButton.setAttribute('aria-controls', 'workspace-changes-panel');
    changesButton.setAttribute('aria-expanded', 'false');
    settings.parentElement.insertBefore(changesButton, settings);
    changesPanel = document.createElement('section');
    changesPanel.id = 'workspace-changes-panel';
    changesPanel.hidden = true;
    changesPanel.setAttribute('aria-label', 'Workspace changes');
    changesPanel.innerHTML =
      '<header><h2>Changes</h2><button type="button" id="workspace-changes-refresh">Refresh</button>' +
      '<button type="button" id="workspace-changes-close" aria-label="Close changes panel">Close</button></header>' +
      '<p id="workspace-changes-summary" role="status">Checking all Git changes…</p>' +
      '<p class="workspace-changes-help">See what changed without leaving your post. ' +
      'Saving a draft does not publish it.</p>' +
      '<div id="workspace-changes-list" role="list"></div>';
    document.body.append(changesPanel);
    changesList = changesPanel.querySelector('#workspace-changes-list');
    changesButton.addEventListener('click', () => changesVisible(changesPanel.hidden));
    changesPanel.querySelector('#workspace-changes-refresh')
      .addEventListener('click', refreshWorkspaceChanges);
    changesPanel.querySelector('#workspace-changes-close')
      .addEventListener('click', () => changesVisible(false));
    window.addEventListener('focus', () => { if (!changesPanel.hidden) refreshWorkspaceChanges(); });
    let persisted = false;
    try { persisted = localStorage.getItem('markata-workspace-changes-open') === 'true'; }
    catch { /* Default closed when storage is blocked. */ }
    if (persisted) changesVisible(true);
  }

  function attach(instance) {
    editor=instance;
    editor.setOption('gutters',['CodeMirror-linenumbers','markata-change-gutter']);
    gutterStatus=document.createElement('span');
    gutterStatus.id='markata-gutter-status';
    gutterStatus.setAttribute('role','status');
    gutterStatus.textContent='Loading Git change gutter…';
    document.querySelector('#saved')?.insertAdjacentElement('afterend',gutterStatus);
    editor.on('change',scheduleChangeGutter);
    const activeLabel=document.querySelector('#active-file');
    if(activeLabel)new MutationObserver(refreshChangeGutter).observe(activeLabel,{childList:true,characterData:true,subtree:true});
    setTimeout(refreshChangeGutter,0);
    const startPublish=()=>{quickPublish();return true;};
    editor.addKeyMap({'Ctrl-S':startPublish,'Cmd-S':startPublish,'Ctrl-X Ctrl-S':startPublish});
    window.CodeMirror?.Vim?.defineEx('write','w',()=>quickPublish());
    button=document.createElement('button');
    button.type='button';button.id='publish-open';
    button.textContent='Publish post';button.title='Publish this post only (Ctrl/Cmd+S, Vim :w)';
    button.setAttribute('aria-label','Publish current post');
    const settings=document.querySelector('#settings-open');
    settings.parentElement.insertBefore(button,settings);
    const doc=document.createElement('dialog');
    doc.id='publish-review';doc.setAttribute('aria-labelledby','publish-title');
    doc.innerHTML='<header class="publish-review-head"><div><small>REVIEW CHANGES</small><h2 id="publish-title">Review changes</h2><p id="publish-path"></p></div><button type="button" id="publish-close" aria-label="Close publish review">✕</button></header>'+
    '<div id="publish-notice" role="status"></div><div id="publish-body" hidden><div id="publish-chips" aria-label="Change summary"></div><section id="publish-no-changes" hidden aria-label="No changes"><h3>All caught up</h3><p>This file matches its last Git commit. You have no new content changes to review.</p><p>If the page is still unpublished, you can mark it as public to prepare a new review.</p></section>'+
    '<nav class="publish-mode-tabs" role="tablist" aria-label="Comparison view"><button id="publish-tab-side" type="button" role="tab" aria-controls="publish-view-side" data-publish-mode="side">Side by side</button><button id="publish-tab-rich" type="button" role="tab" aria-controls="publish-view-rich" data-publish-mode="rich">Rendered diff</button><button id="publish-tab-source" type="button" role="tab" aria-controls="publish-view-source" data-publish-mode="source">Source diff</button><button id="publish-tab-pierre" type="button" role="tab" aria-controls="publish-view-pierre" data-publish-mode="pierre">Full diff</button><button id="publish-tab-message" type="button" role="tab" aria-controls="publish-view-message" data-publish-mode="message">Commit message</button></nav>'+
    '<div id="publish-view-side" role="tabpanel" aria-labelledby="publish-tab-side" data-publish-view="side"><div class="publish-side-tabs" role="group" aria-label="Version to read"><button data-publish-side="before" type="button">Last committed</button><button data-publish-side="after" type="button">Current draft</button></div>'+
    '<div class="publish-comparison"><section data-publish-pane="before"><h3 id="publish-before-name">Last committed</h3><iframe id="publish-before-frame" title="Last committed rendering" sandbox="allow-scripts"></iframe></section>'+
    '<section data-publish-pane="after"><h3>Current saved draft</h3><iframe id="publish-after-frame" title="Current draft rendering" sandbox="allow-scripts"></iframe></section></div></div>'+
    '<div id="publish-view-rich" role="tabpanel" aria-labelledby="publish-tab-rich" data-publish-view="rich" hidden><p class="publish-rich-help">Structural comparison: inline changes for text, before/after blocks for complex elements. Use Side by side for exact rendering.</p><iframe id="publish-rich-frame" title="Unified rendered Markdown change annotations" sandbox="allow-scripts"></iframe></div>'+
    '<div id="publish-view-pierre" role="tabpanel" aria-labelledby="publish-tab-pierre" data-publish-view="pierre" hidden><div class="publish-pierre-actions"><span>Pierre source diff</span><button type="button" data-pierre-layout="stacked" class="active">Unified</button><button type="button" data-pierre-layout="split">Split</button></div><p id="publish-pierre-error" role="status" hidden></p><div id="publish-pierre-host"></div></div>'+
    '<div id="publish-view-source" role="tabpanel" aria-labelledby="publish-tab-source" data-publish-view="source" hidden><div id="publish-source-diff" role="region" tabindex="0" aria-label="Unified Markdown source diff"></div></div>'+
    '<div id="publish-view-message" role="tabpanel" aria-labelledby="publish-tab-message" data-publish-view="message" hidden><label for="publish-commit-message">Commit message</label><textarea id="publish-commit-message" rows="4" maxlength="4000" spellcheck="true"></textarea><p class="publish-message-help"><span id="publish-message-count">0 / 4000 bytes</span> Used for the local commit. Nothing is pushed from this screen.</p></div></div>'+
    '<div id="publish-discard-warning" hidden role="status"></div>'+
    '<div id="publish-discard-confirmation" hidden role="alertdialog" aria-label="Confirm discard"><p id="publish-discard-description">Restore this document to its last committed version? This removes your reviewed unpublished changes, without changing Git history.</p><button type="button" id="publish-discard-cancel">Cancel</button><button type="button" id="publish-discard-accept">Discard reviewed changes</button></div>'+
    '<footer class="publish-actions"><span>A Git commit is not a verified deployment.</span><button type="button" id="publish-discard" hidden>Discard changes…</button><button type="button" id="publish-back">Back to editing</button><button type="button" id="publish-mark-public" hidden>Mark as public</button><button type="button" id="publish-confirm" disabled>Commit selected page</button><button type="button" id="publish-push" hidden>Push to origin</button><button type="button" id="publish-live-check" hidden>Check live status</button></footer>';
    document.body.append(doc);dialog=doc;notice=doc.querySelector('#publish-notice');panels=doc.querySelector('#publish-body');
    button.addEventListener('click',quickPublish);
    // Keep detailed review discoverable without making it a publish prerequisite.
    const reviewButton=document.createElement('button');
    reviewButton.type='button';
    reviewButton.id='publish-review-open';
    reviewButton.textContent='Review changes';
    reviewButton.title='Advanced review for the active post (Ctrl/Cmd+Shift+Enter)';
    button.insertAdjacentElement('afterend',reviewButton);
    reviewButton.addEventListener('click',open);
    publishStatus=document.createElement('span');
    publishStatus.id='markata-publish-status';
    publishStatus.setAttribute('role','status');
    publishStatus.setAttribute('aria-live','polite');
    publishStatus.textContent='Not published';
    document.querySelector('#saved')?.insertAdjacentElement('afterend',publishStatus);
    attachWorkspaceChanges();
    editor.on('change',updatePublishButton);
    const fileLabel = document.querySelector('#active-file');
    if (fileLabel) new MutationObserver(updatePublishButton).observe(fileLabel,{subtree:true,characterData:true,childList:true});
    updatePublishButton();
    doc.querySelector('#publish-close').addEventListener('click',()=>doc.close());
    doc.querySelector('#publish-back').addEventListener('click',()=>doc.close());
    doc.querySelector('#publish-discard').addEventListener('click',()=>{
      if (state.busy || !state.data) return;
      doc.querySelector('#publish-discard-description').textContent=state.data.baseline_exists
        ? 'Restore '+state.path+' to the last committed Git version? The reviewed unpublished edits will be removed from the active workspace.'
        : 'Permanently delete the new untracked file '+state.path+' and its browser recovery copy? This cannot be undone; other files and Git history will not be changed.';
      doc.querySelector('#publish-discard-confirmation').hidden=false;
      doc.querySelector('#publish-discard-cancel').focus();
    });
    doc.querySelector('#publish-discard-cancel').addEventListener('click',()=>{doc.querySelector('#publish-discard-confirmation').hidden=true;doc.querySelector('#publish-discard').focus();});
    doc.querySelector('#publish-discard-accept').addEventListener('click',discardReviewed);
    doc.querySelector('#publish-mark-public').addEventListener('click',markPublicAndReviewAgain);
    doc.querySelector('#publish-confirm').addEventListener('click',commitReview);
    doc.querySelector('#publish-push').addEventListener('click',pushReview);
    doc.querySelector('#publish-live-check').addEventListener('click',checkLiveDeployment);
    doc.addEventListener('close',()=>{
      reviewGeneration++;
      // Reopening must not wait on a closed dialog's network request.
      state.loading=false;
      state.data=null;
      reviewRenderController?.abort();
      reviewRenderController=null;
      window.MarkataPierreDiff?.release?.(doc.querySelector('#publish-pierre-host'));
      editor.focus();
    });
    doc.querySelectorAll('[data-publish-mode]').forEach(b=>b.addEventListener('click',()=>setMode(b.dataset.publishMode)));
    doc.querySelector('.publish-mode-tabs').addEventListener('keydown',event=>{
      if(!['ArrowLeft','ArrowRight','Home','End'].includes(event.key))return;
      const tabs=[...doc.querySelectorAll('[data-publish-mode]:not(:disabled)')];
      const current=tabs.indexOf(event.target.closest('[data-publish-mode]'));
      if(current<0)return;
      event.preventDefault();
      const next=event.key==='Home'?0:event.key==='End'?tabs.length-1:(current+(event.key==='ArrowRight'?1:-1)+tabs.length)%tabs.length;
      tabs[next].click();tabs[next].focus();
    });
    doc.querySelector('#publish-commit-message').addEventListener('input',updateCommitMessage);
    doc.querySelector('#publish-commit-message').addEventListener('blur',event=>{event.currentTarget.value=event.currentTarget.value.trim();updateCommitMessage();});
    window.addEventListener('message',comparisonScrollSync);
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
    if (window.MarkataShortcuts) {
      window.MarkataShortcuts.register({
        id:'publishing.review', label:'Publishing: Review changes', category:'Publishing',
        description:'Review saved and committed changes before publishing.',
        aliases:'review changes publish go live deploy release git',
        bindings:{vscode:'Mod+Shift+Enter',writer:'Mod+Shift+Enter',vim:'Mod+Shift+Enter',emacs:'Ctrl+Shift+Enter'},
        run:open
      });
    } else {
      document.addEventListener('keydown',event=>{
        if(event.isComposing||event.altKey||!event.shiftKey||!(event.ctrlKey||event.metaKey)||event.key!=='Enter')return;
        if(document.querySelector('dialog[open]'))return;
        event.preventDefault();event.stopPropagation();open();
      },true);
    }
  }
  window.MarkataPublishReview={attach,open,quickPublish};
})();
