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
  let lastPushedCommit = null, liveCheckButton = null;
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
  function setPagePublishStatus(message, level='info') {
    if (publishStatus) {
      publishStatus.textContent = message;
      publishStatus.dataset.state = level;
      publishStatus.title = message;
    }
    window.dispatchEvent(new CustomEvent('markata:publish-status', {
      detail: {message, type:level, pending:false}
    }));
  }
  function finishPush() {
    if (!state.commit?.commit) return;
    lastPushedCommit = {...state.commit, path:state.path};
    setPagePublishStatus('Pushed to origin · live status not verified.');
    if (liveCheckButton) liveCheckButton.hidden = false;
    dialog.close();
    schedulePublishStateRefresh();
  }
  function isRemoteBranchRace(status, message) {
    return status === 409 && /remote branch changed|sync and review again/i.test(message || '');
  }
  function showRemoteBranchRace() {
    const retry = dialog.querySelector('#publish-push');
    retry.hidden = true;
    retry.disabled = true;
    dialog.querySelector('#publish-back').textContent = 'Close and keep local commit';
    setNotice('Origin changed after this commit. Your local commit is retained. This editor cannot sync the branch yet; close this review before making another push.', 'warning');
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
      // Pierre measures the host while mounting. Let the tab's hidden state
      // reach layout before creating the renderer, or it can mount blank.
      await new Promise(resolve => requestAnimationFrame(() => requestAnimationFrame(resolve)));
      if(state.mode!=='pierre' || !state.data || !dialog.open)return;
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
    const disabled=!canCommit||!input.value.trim()||bytes>4000;
    dialog.querySelector('#publish-confirm').disabled=disabled;
    const commitPush=dialog.querySelector('#publish-commit-and-push');
    if(commitPush)commitPush.disabled=disabled;
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
    const frontmatter=source.match(/^---\r?\n[\s\S]*?\r?\n---(?:\r?\n|$)/);
    const response=await fetch('/render',{method:'POST',headers:{'Content-Type':'application/json'},
      body:JSON.stringify({markdown:contentBody(source),frontmatter_source:frontmatter?.[0] || '',include_document:true,appearance}),signal});
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
    dialog.querySelector('#publish-back').textContent='Back to editing';
    const commitPushBtn=dialog.querySelector('#publish-commit-and-push');
    if(commitPushBtn){commitPushBtn.hidden=false;commitPushBtn.disabled=true;}
    dialog.querySelector('#publish-discard').hidden=true;
    dialog.querySelector('#publish-discard-warning').hidden=true;
    dialog.querySelector('#publish-discard-confirmation').hidden=true;
    dialog.querySelector('#publish-no-changes').hidden=true;
    dialog.querySelector('.publish-mode-tabs').hidden=false;
    dialog.querySelector('#publish-push').hidden=true;
    dialog.querySelector('#publish-push').disabled=false;
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
      const commitDisabled=!data.changed || fields.published!=='true' || fields.draft==='true';
      dialog.querySelector('#publish-confirm').disabled=commitDisabled;
      const commitPushBtn2=dialog.querySelector('#publish-commit-and-push');
      if(commitPushBtn2)commitPushBtn2.disabled=commitDisabled;
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
      const commitPush=dialog.querySelector('#publish-commit-and-push');
      if(commitPush)commitPush.hidden=true;
      dialog.querySelector('#publish-mark-public').hidden=true;
      dialog.querySelector('#publish-push').hidden=false;
      setNotice('Committed only this page to Git. Other staged files were excluded. Push separately to trigger your configured deployment pipeline.','info');
    }catch(err){setNotice(err.message,'warning');btn.disabled=false;}
    finally{state.busy=false;btn.textContent='Commit selected page';}
  }
  async function commitAndPushReview() {
    if(state.busy||!state.data)return;
    const data=state.data;
    const messageField=dialog.querySelector('#publish-commit-message');
    const message=messageField.value.trim();
    if(messageField.value!==message){messageField.value=message;updateCommitMessage();}
    if(!message||new TextEncoder().encode(message).length>4000) {
      setNotice('Enter a commit message of 1 to 4000 bytes before publishing.','warning');return;
    }
    if(state.path!==window.MarkataSiteTree?.activeFile?.()||editor.getValue()!==data.current) {
      setNotice('The document changed. Close and reopen review before publishing.','warning');return;
    }
    state.busy=true;
    const pushBtn=dialog.querySelector('#publish-commit-and-push');
    const commitBtn=dialog.querySelector('#publish-confirm');
    pushBtn.disabled=true;commitBtn.disabled=true;
    pushBtn.textContent='Publishing…';
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
      commitBtn.hidden=true;
      pushBtn.hidden=true;
      dialog.querySelector('#publish-mark-public').hidden=true;
      const pushResp=await fetch('/workspace/publish/push?path='+encodeURIComponent(state.path),{
        method:'POST',headers:{'Content-Type':'application/json'},
        body:JSON.stringify({commit:state.commit.commit,branch:state.commit.branch})
      });
      const pushResult=await pushResp.json();
      if(!pushResp.ok){
        if (isRemoteBranchRace(pushResp.status, pushResult.error)) {
          showRemoteBranchRace();
          return;
        }
        const retry = dialog.querySelector('#publish-push');
        retry.hidden=false;
        retry.disabled=false;
        throw new Error(pushResult.error||'Committed locally, but push failed. Click Push to origin to retry.');
      }
      refreshWorkspaceChanges();
      finishPush();
    }catch(err){
      setNotice(err.message,'warning');
      if(!state.commit){pushBtn.disabled=false;commitBtn.disabled=false;}
      else {
        const retry = dialog.querySelector('#publish-push');
        retry.hidden=false;
        retry.disabled=false;
        setNotice(err.message+' The local commit is retained; retrying checks whether this exact commit already reached origin.', 'warning');
      }
    }finally{
      state.busy=false;
      pushBtn.textContent='Commit & Push';
    }
  }
  async function checkLiveDeployment() {
    const receipt = state.commit?.commit ? state.commit : lastPushedCommit;
    if(!receipt || state.busy)return;
    const btn=dialog.querySelector('#publish-live-check');
    const pageBtn=liveCheckButton;
    state.busy=true;
    if (btn) { btn.disabled=true; btn.textContent='Checking…'; }
    if (pageBtn) { pageBtn.disabled=true; pageBtn.textContent='Checking…'; }
    const report = (message, level='info') => {
      if (dialog.open) setNotice(message, level);
      else setPagePublishStatus(message, level);
    };
    try{
      const response=await fetch('/workspace/publish/deployment?commit='+encodeURIComponent(receipt.commit),{cache:'no-store'});
      const result=await response.json();
      if(!response.ok)throw new Error(result.error||'Cannot verify deployment.');
      if(result.verified){
        if (btn) btn.hidden=true;
        if (pageBtn) pageBtn.hidden=true;
        report('Live · Verified source revision '+result.live_source_commit.slice(0,10)+(result.release_id?' · release '+result.release_id:'')+'.','info');
      } else {
        report(result.message||'Deployment has not been verified.','warning');
      }
    }catch(error){report(error.message,'warning');}
    finally{
      state.busy=false;
      if (btn) { btn.disabled=false; btn.textContent='Check live status'; }
      if (pageBtn) { pageBtn.disabled=false; pageBtn.textContent='Check live status'; }
    }
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
      if(!response.ok) {
        if (isRemoteBranchRace(response.status, result.error)) {
          showRemoteBranchRace();
          return;
        }
        throw new Error(result.error||'Push failed. The local commit is retained.');
      }
      refreshWorkspaceChanges();
      finishPush();
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

  function draftPostSource(source) {
    const newline = source.includes('\r\n') ? '\r\n' : '\n';
    if (!source.startsWith('---' + newline))
      return '---' + newline + 'published: false' + newline + 'draft: true' +
        newline + '---' + newline + newline + source;
    const end = source.indexOf(newline + '---', 3 + newline.length);
    if (end < 0) throw new Error('Fix malformed YAML frontmatter before saving as draft.');
    let head = source.slice(0,end);
    const tail = source.slice(end);
    if (/^published:\s*(true|false)\s*$/im.test(head))
      head = head.replace(/^published:\s*(true|false)\s*$/im, 'published: false');
    else head += newline + 'published: false';
    if (/^draft:\s*(true|false)\s*$/im.test(head))
      head = head.replace(/^draft:\s*(true|false)\s*$/im, 'draft: true');
    else head += newline + 'draft: true';
    return head + tail;
  }

  // Writer-first publication keeps the Git transaction behind one explicit
  // action, while preserving the guarded review/commit/push endpoints.
  let quickBusy = false;
  let pendingQuick = null;
  let publishStatus = null;
  let publishSnapshot = null;
  let publishRefreshTimer = null;
  function reportQuick(message, type='info') {
    if (publishStatus) {
      publishStatus.textContent = message;
      publishStatus.dataset.state = type;
      publishStatus.title = message;
    }
    window.dispatchEvent(new CustomEvent('markata:publish-status', {
      detail: { message, type, pending: Boolean(pendingQuick && !pendingQuick.blocked) }
    }));
  }
  async function saveAsDraft() {
    if (quickBusy || state.busy) return;
    const path = canOpen();
    if (!path) {
      reportQuick('Choose a site Markdown post first.', 'warning');
      return;
    }
    quickBusy = true;
    try {
      const current = editor.getValue();
      const next = draftPostSource(current);
      if (next !== current) {
        editor.operation(() => editor.replaceRange(next, {line:0,ch:0},
          {line:editor.lastLine(),ch:editor.getLine(editor.lastLine()).length}, '+draft-meta'));
      }
      reportQuick('Saving as draft…', 'loading');
      await awaitSaved(path, editor.getValue());
      reportQuick('Saved as draft (not published)', 'info');
      refreshWorkspaceChanges();
      refreshChangeGutter();
    } catch (error) {
      reportQuick('Could not save draft: ' + error.message, 'warning');
    } finally {
      quickBusy = false;
      updatePublishButton();
    }
  }
  async function finishQuickPush(pending) {
    reportQuick('Publishing to your site…', 'loading');
    const response = await fetch('/workspace/publish/push?path=' +
      encodeURIComponent(pending.path), {
      method:'POST', headers:{'Content-Type':'application/json'},
      body:JSON.stringify({commit:pending.commit,branch:pending.branch})
    });
    const result = await response.json();
    if (!response.ok) {
      if (isRemoteBranchRace(response.status, result.error)) {
        pending.blocked = true;
        reportQuick('Origin changed. Your local commit is retained; review the post before syncing. This editor cannot sync the branch yet.', 'warning');
        updatePublishButton();
        return;
      }
      throw new Error(result.error || 'Could not push. Your local commit is retained.');
    }
    pendingQuick = null;
    lastPushedCommit = {commit:pending.commit, branch:pending.branch, path:pending.path};
    if (liveCheckButton) liveCheckButton.hidden = false;
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
      updatePublishButton();
    }
  }

  function updatePublishButton() {
    if (!button || !editor) return;
    const path = canOpen();
    const snapshot = publishSnapshot?.path === path ? publishSnapshot : null;
    const source = editor.getValue();
    const fields = frontmatter(source);
    const upToDate = Boolean(snapshot && snapshot.baselineExists && !snapshot.changed &&
      snapshot.current === source && fields.published === 'true' && fields.draft !== 'true');
    const publishState = quickBusy ? 'publishing' : pendingQuick?.blocked ? 'sync-needed' : pendingQuick ? 'retry' :
      !path ? 'no-post' : upToDate ? 'up-to-date' : 'ready';
    const label = {publishing:'Publishing…',retry:'Retry publish','sync-needed':'Remote changed','no-post':'Choose post',
      'up-to-date':'Up to date',ready:'Publish post'}[publishState];
    if (button.dataset.publishState !== publishState) {
      const icon=document.createElementNS('http://www.w3.org/2000/svg','svg');
      icon.setAttribute('viewBox','0 0 24 24');icon.setAttribute('aria-hidden','true');icon.setAttribute('focusable','false');
      const pathNode=document.createElementNS('http://www.w3.org/2000/svg','path');
      pathNode.setAttribute('d',publishState==='up-to-date'?'M5 12l4 4L19 6':'M12 15V3m0 0 4 4m-4-4L8 7M4 14v6h16v-6');icon.append(pathNode);
      button.replaceChildren(icon,document.createTextNode(label));
    }
    button.dataset.publishState = publishState;
    button.disabled = quickBusy || publishState === 'up-to-date' || publishState === 'sync-needed' || !path;
    button.setAttribute('aria-busy',String(quickBusy));
    button.setAttribute('aria-label',publishState==='up-to-date'?'Up to date in Git; deployment not verified':
      publishState==='retry'?'Retry pending publish':publishState==='sync-needed'?'Remote branch changed; review before publishing again':publishState==='no-post'?'Choose a site post to publish':
      publishState==='publishing'?'Publishing current post':'Publish current post');
    button.title = publishState==='up-to-date'?
      'This saved post matches its local Git commit. Remote push and site deployment are not verified.':
      publishState==='sync-needed'?'Origin changed. The local commit is retained; review and sync the branch before another push.':
      'Publish this post only (Ctrl/Cmd+S). Review changes with Ctrl/Cmd+Shift+Enter.';
  }

  function schedulePublishStateRefresh() {
    // A saved file or regained focus invalidates the previous proof immediately.
    // Reuse the gutter's review request; never poll or fetch a second full source.
    publishSnapshot = null;
    gutterRequest++;
    updatePublishButton();
    clearTimeout(publishRefreshTimer);
    publishRefreshTimer = setTimeout(refreshChangeGutter,180);
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
    clearTimeout(publishRefreshTimer);
    publishSnapshot=null;updatePublishButton();
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
      if (typeof data.current==='string' && typeof data.changed==='boolean') {
        publishSnapshot={path,current:data.current,baselineExists:data.baseline_exists===true,changed:data.changed};
      }
      updatePublishButton();
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
  let changesPanel = null, changesButton = null, changesList = null, changesHost = null;
  let changesReviewDialog = null, changesReviewQueue = null, changesReviewDetail = null;
  let changesSnapshot = null, changesReviewPath = null;
  let changesGeneration = 0, changesTimer = null, changesReviewGeneration = 0;
  let changesFilter = 'all', changesReviewFilter = 'all';
  const changesReviewed = new Set();
  function changesOpen() {
    return document.querySelector('#site-view')?.value === 'changes';
  }
  // The sidebar view selector is the single source of truth for whether the
  // Changes list is shown; this only requests a view change.
  function changesVisible(open) {
    const select = document.querySelector('#site-view');
    if (!changesPanel || !select) return;
    const next = open ? 'changes' : (select.value === 'changes' ? 'files' : select.value);
    if (select.value === next) return;
    select.value = next;
    select.dispatchEvent(new Event('change'));
  }
  function syncChangesView() {
    const open = changesOpen();
    changesHost.hidden = !open;
    changesPanel.hidden = !open;
    clearInterval(changesTimer);
    changesTimer = null;
    try { localStorage.setItem('markata-workspace-changes-open', open ? 'true' : 'false'); }
    catch { /* Keep the panel usable without storage. */ }
    if (open) {
      if (changesSnapshot) showWorkspaceChanges(changesSnapshot);
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
  function changeCode(change) {
    if (change.untracked) return 'U';
    if (change.kind === 'conflict') return '!';
    if (change.kind === 'deleted') return 'D';
    if (change.kind === 'renamed') return 'R';
    if (change.kind === 'copied') return 'C';
    if (change.kind === 'added') return 'A';
    return 'M';
  }
  function setChangesCount(label, description) {
    const badge = changesButton.querySelector('.workspace-changes-count');
    badge.textContent = label;
    badge.hidden = !label;
    changesButton.setAttribute('aria-label', 'Changes' + (description ? ': ' + description : ''));
    changesButton.title = description ? 'Changes · ' + description : 'Changes';
  }
  function changeInventorySignature(data) {
    return JSON.stringify({head:data.head,entries:data.entries});
  }
  function filteredChanges(entries, filter, query) {
    const needle = String(query || '').trim().toLocaleLowerCase();
    return entries.filter(change => {
      if (filter === 'posts' && !change.markdown) return false;
      if (filter === 'other' && change.markdown) return false;
      return !needle || (change.path + ' ' + change.status + ' ' + changeLabel(change))
        .toLocaleLowerCase().includes(needle);
    }).sort((a, b) => Number(!a.markdown) - Number(!b.markdown) || a.path.localeCompare(b.path));
  }
  function changeStateText(change) {
    const states = [];
    if (change.staged) states.push('staged');
    if (change.unstaged) states.push('unstaged');
    if (change.untracked) states.push('untracked');
    if (!change.workspace) states.push('outside this site');
    if (!states.length) states.push('working copy');
    return states.join(' · ');
  }
  function renderInventoryRows() {
    if (!changesSnapshot) return;
    const entries = filteredChanges(changesSnapshot.entries, changesFilter,
      changesPanel.querySelector('#workspace-changes-search').value);
    changesList.replaceChildren();
    if (!entries.length) {
      const empty = document.createElement('p');
      empty.className = 'workspace-changes-empty';
      empty.textContent = changesSnapshot.total
        ? 'No files match this filter.'
        : 'All caught up. No working-tree changes. Unpushed commits may still exist.';
      changesList.append(empty);
      return;
    }
    let group = null;
    for (const change of entries) {
      const nextGroup = change.markdown ? 'posts' : 'other';
      if (nextGroup !== group) {
        group = nextGroup;
        const heading = document.createElement('div');
        heading.className = 'workspace-change-group';
        heading.setAttribute('role', 'presentation');
        heading.textContent = group === 'posts' ? 'Posts' : 'Other files';
        changesList.append(heading);
      }
      const row = document.createElement('button');
      row.type = 'button';
      row.className = 'workspace-change';
      row.dataset.kind = change.untracked ? 'untracked' : change.kind || 'modified';
      row.setAttribute('aria-label', change.path + ', ' + changeLabel(change) + ', ' + changeStateText(change));
      const code = document.createElement('span');
      code.className = 'workspace-change-code';
      code.textContent = changeCode(change);
      code.setAttribute('aria-hidden', 'true');
      const description = document.createElement('span');
      description.className = 'workspace-change-description';
      const name = document.createElement('strong');
      name.textContent = change.path;
      name.title = change.path;
      const status = document.createElement('span');
      status.className = 'workspace-change-state';
      status.textContent = changeLabel(change) + ' · ' + changeStateText(change);
      description.append(name, status);
      const arrow = document.createElement('span');
      arrow.className = 'workspace-change-open-mark';
      arrow.setAttribute('aria-hidden', 'true');
      arrow.textContent = '›';
      row.append(code, description, arrow);
      row.addEventListener('click', () => openChangesReview(change.path));
      changesList.append(row);
    }
  }
  function showWorkspaceChanges(data, force = false) {
    const changed = !changesSnapshot || changeInventorySignature(data) !== changeInventorySignature(changesSnapshot);
    changesSnapshot = data;
    setChangesCount(data.total ? String(data.total) : '', data.total ? data.total + ' with changes' : '');
    const heading = changesPanel.querySelector('#workspace-changes-summary');
    heading.textContent = data.total
      ? data.total + (data.total === 1 ? ' file' : ' files') + ' with changes · ' +
        data.staged + ' staged · ' + data.unstaged + ' unstaged'
      : 'Nothing to review';
    changesPanel.querySelector('#workspace-changes-review-open').disabled = !data.total;
    changesPanel.querySelector('#workspace-changes-review-open').textContent = 'Review queue';
    changesPanel.querySelector('#workspace-changes-staged').textContent = data.staged;
    changesPanel.querySelector('#workspace-changes-unstaged').textContent = data.unstaged;
    changesPanel.querySelector('#workspace-changes-untracked').textContent = data.untracked;
    renderInventoryRows();
    if (changesReviewDialog?.open) {
      renderChangesReviewQueue();
      if (changesReviewPath && (changed || force)) selectWorkspaceChange(changesReviewPath);
    }
  }
  function changesReviewEntries() {
    if (!changesSnapshot) return [];
    return filteredChanges(changesSnapshot.entries, changesReviewFilter,
      changesReviewDialog.querySelector('#workspace-review-search').value);
  }
  function renderChangesReviewQueue() {
    if (!changesReviewQueue || !changesSnapshot) return;
    const entries = changesReviewEntries();
    if (!entries.some(change => change.path === changesReviewPath))
      changesReviewPath = entries[0]?.path || null;
    changesReviewQueue.replaceChildren();
    if (!entries.length) {
      const empty = document.createElement('p');
      empty.className = 'workspace-review-queue-empty';
      empty.textContent = changesSnapshot.total ? 'No files match this filter.' : 'No working-tree changes.';
      changesReviewQueue.append(empty);
    }
    let group = null;
    for (const change of entries) {
      const nextGroup = change.markdown ? 'posts' : 'other';
      if (nextGroup !== group) {
        group = nextGroup;
        const heading = document.createElement('div');
        heading.className = 'workspace-review-group';
        heading.textContent = group === 'posts' ? 'Posts' : 'Other files';
        changesReviewQueue.append(heading);
      }
      const row = document.createElement('button');
      row.type = 'button';
      row.className = 'workspace-review-item';
      row.dataset.kind = change.untracked ? 'untracked' : change.kind || 'modified';
      row.dataset.path = change.path;
      row.setAttribute('aria-current', String(change.path === changesReviewPath));
      row.setAttribute('aria-label', change.path + ', ' + changeLabel(change) + ', ' + changeStateText(change));
      const mark = document.createElement('span');
      mark.className = 'workspace-review-item-mark';
      mark.textContent = changeCode(change);
      mark.setAttribute('aria-hidden', 'true');
      const copy = document.createElement('span');
      copy.className = 'workspace-review-item-copy';
      const path = document.createElement('strong');
      path.textContent = change.path;
      path.title = change.path;
      const stateText = document.createElement('small');
      stateText.textContent = changeLabel(change) + ' · ' + changeStateText(change);
      copy.append(path, stateText);
      const flag = document.createElement('span');
      flag.className = changesReviewed.has(change.path) ? 'workspace-review-kept' : 'workspace-review-pending';
      flag.textContent = changesReviewed.has(change.path) ? 'Kept' : 'Review';
      row.append(mark, copy, flag);
      row.addEventListener('click', () => selectWorkspaceChange(change.path));
      changesReviewQueue.append(row);
    }
    const position = entries.findIndex(change => change.path === changesReviewPath);
    changesReviewDialog.querySelector('#workspace-review-position').textContent =
      position < 0 ? '0 / ' + entries.length : (position + 1) + ' / ' + entries.length;
    changesReviewDialog.querySelector('#workspace-review-progress').textContent =
      changesReviewed.size + ' kept · ' + Math.max(0, changesSnapshot.total - changesReviewed.size) + ' to review';
    const active = changesReviewQueue.querySelector('[aria-current="true"]');
    active?.scrollIntoView({block:'nearest'});
  }
  function resetUnavailableReason(change) {
    if (!change.workspace) return 'This file is outside the selected site.';
    if (change.staged) return 'Reset is unavailable while this file has staged changes.';
    if (change.kind === 'deleted') return 'Deleted files cannot be restored from this review queue.';
    if (change.kind === 'conflict') return 'Resolve this Git conflict before resetting the file.';
    if (!change.markdown) return 'Reset is available from the workbench for Markdown posts only.';
    if (!isOpenableWorkspacePost(change))
      return 'This Markdown file is not available in the editor workspace.';
    return '';
  }
  function isOpenableWorkspacePost(change) {
    return Boolean(change.workspace && change.markdown && change.editor_path &&
      !['deleted', 'conflict'].includes(change.kind) &&
      window.MarkataSiteTree?.files?.()?.some(file => file.path === change.editor_path));
  }
  function renderWorkspaceChangeDiff(change, data) {
    const detail = changesReviewDetail;
    const diff = detail.querySelector('#workspace-review-diff');
    diff.replaceChildren();
    const lines = diffLines(data.baseline, data.current);
    detail.querySelector('#workspace-review-diff-count').textContent = lines.length ? 'Source changes' : 'No source differences';
    if (!lines.length) {
      const empty = document.createElement('p');
      empty.className = 'workspace-review-diff-empty';
      empty.textContent = 'The saved file matches the last committed version.';
      diff.append(empty);
    }
    detail.querySelector('#workspace-review-baseline').textContent = data.baseline_exists
      ? 'Last commit · ' + data.head.slice(0, 8)
      : 'No committed baseline · new file';
    detail.querySelector('#workspace-review-working').textContent = 'Saved workspace version';
    detail.querySelector('#workspace-review-meta').textContent = change.previous_path
      ? 'Renamed from ' + change.previous_path + ' · ' + changeStateText(change)
      : changeStateText(change);
    detail.querySelector('#workspace-review-unavailable').hidden = true;
    detail.querySelector('#workspace-review-diff-wrap').hidden = false;
    diff.setAttribute('aria-label', 'Source diff for ' + change.path);
    for (const line of lines) {
      const row = document.createElement('div');
      row.className = 'workspace-review-diff-line ' + line.type;
      const prefix = document.createElement('span');
      prefix.className = 'workspace-review-diff-prefix';
      prefix.textContent = line.type === 'add' ? '+' : line.type === 'del' ? '−' : ' ';
      const content = document.createElement('span');
      content.textContent = line.text || ' ';
      row.append(prefix, content);
      diff.append(row);
    }
  }
  function showWorkspaceChangeUnavailable(message) {
    const detail = changesReviewDetail;
    changesReviewGeneration++;
    window.MarkataPierreDiff?.release?.(detail.querySelector('#workspace-review-diff'));
    const load = detail.querySelector('#workspace-review-load');
    if (load) load.hidden = true;
    detail.querySelector('#workspace-review-diff-wrap').hidden = true;
    const unavailable = detail.querySelector('#workspace-review-unavailable');
    unavailable.textContent = message;
    unavailable.hidden = false;
    detail.querySelector('#workspace-review-diff-count').textContent = '';
    detail.querySelector('#workspace-review-baseline').textContent = 'Baseline unavailable';
    detail.querySelector('#workspace-review-working').textContent = '';
  }
  async function selectWorkspaceChange(path) {
    if (!changesReviewDialog?.open) return;
    const change = changesSnapshot?.entries.find(entry => entry.path === path);
    if (!change) return;
    window.MarkataPierreDiff?.release?.(changesReviewDetail.querySelector('#workspace-review-diff'));
    changesReviewPath = path;
    const generation = ++changesReviewGeneration;
    renderChangesReviewQueue();
    changesReviewDetail.innerHTML =
      '<header class="workspace-review-detail-head"><div><small id="workspace-review-kind"></small><h3 id="workspace-review-path"></h3><p id="workspace-review-meta"></p></div><button type="button" id="workspace-review-formats">Full review…</button></header>' +
      '<div class="workspace-review-baselines"><span id="workspace-review-baseline">Loading baseline…</span><span id="workspace-review-working"></span><span id="workspace-review-diff-count"></span></div>' +
      '<p class="workspace-review-load" id="workspace-review-load" role="status">Loading saved source comparison…</p>' +
      '<div class="workspace-review-diff-wrap" id="workspace-review-diff-wrap" hidden><div class="workspace-review-diff" id="workspace-review-diff" role="region" tabindex="0" aria-label="Source diff"></div></div>' +
      '<p class="workspace-review-unavailable" id="workspace-review-unavailable" role="status" hidden></p>' +
      '<footer class="workspace-review-actions"><div><strong>Keep leaves this saved file unchanged.</strong><span id="workspace-review-reset-note">Full review has rendered, side-by-side, source, and full diff modes; reset needs confirmation.</span></div><button type="button" id="workspace-review-keep">Keep in workspace</button><button type="button" id="workspace-review-open-file">Open in editor</button><button type="button" id="workspace-review-reset">Review and reset…</button></footer>';
    const detail = changesReviewDetail;
    detail.querySelector('#workspace-review-kind').textContent = changeLabel(change) + ' · ' + change.status;
    detail.querySelector('#workspace-review-path').textContent = change.path;
    detail.querySelector('#workspace-review-meta').textContent = change.previous_path
      ? 'Renamed from ' + change.previous_path + ' · ' + changeStateText(change)
      : changeStateText(change);
    const keep = detail.querySelector('#workspace-review-keep');
    keep.disabled = changesReviewed.has(change.path);
    keep.textContent = changesReviewed.has(change.path) ? 'Kept for this review' : 'Keep in workspace';
    keep.addEventListener('click', () => {
      changesReviewed.add(change.path);
      renderChangesReviewQueue();
      keep.disabled = true;
      keep.textContent = 'Kept for this review';
      const next = changesReviewEntries().find(entry => !changesReviewed.has(entry.path));
      if (next) {
        selectWorkspaceChange(next.path);
        changesReviewQueue.querySelector('[aria-current="true"]')?.focus();
      }
    });
    const reset = detail.querySelector('#workspace-review-reset');
    const formats = detail.querySelector('#workspace-review-formats');
    const openFile = detail.querySelector('#workspace-review-open-file');
    const resetReason = resetUnavailableReason(change);
    reset.disabled = Boolean(resetReason);
    detail.querySelector('#workspace-review-reset-note').textContent = resetReason || 'Full review has rendered, side-by-side, source, and full diff modes; reset needs confirmation.';
    reset.title = resetReason || 'Open the saved Markdown post and review its reset safeguards.';
    formats.disabled = !isOpenableWorkspacePost(change);
    formats.title = formats.disabled
      ? 'Full review modes are available for Markdown posts in this workspace.'
      : 'Choose rendered, side-by-side, source, or full diff.';
    formats.addEventListener('click', async () => {
      if (formats.disabled || !window.MarkataSiteTree?.openFile || !window.MarkataPublishReview?.open) return;
      formats.disabled = true;
      formats.textContent = 'Opening review…';
      try {
        await window.MarkataSiteTree.openFile(change.editor_path);
        if (window.MarkataSiteTree.activeFile?.() !== change.editor_path)
          throw new Error('The selected file is not available in the editor workspace.');
        changesReviewDialog.close();
        await window.MarkataPublishReview.open();
      } catch (error) {
        formats.disabled = false;
        formats.textContent = 'Full review…';
        showWorkspaceChangeUnavailable('Could not open the full post review: ' + error.message);
      }
    });
    openFile.disabled = !isOpenableWorkspacePost(change);
    openFile.title = openFile.disabled ? resetReason || 'This file cannot be opened in the editor.' : 'Open this Markdown post in the editor.';
    openFile.addEventListener('click', async () => {
      if (openFile.disabled) return;
      openFile.disabled = true;
      try {
        await window.MarkataSiteTree.openFile(change.editor_path);
        if (window.MarkataSiteTree.activeFile?.() !== change.editor_path)
          throw new Error('The selected file is not available in the editor workspace.');
        changesReviewDialog.close();
        changesVisible(false);
      } catch (error) {
        openFile.disabled = false;
        showWorkspaceChangeUnavailable('Could not open this file: ' + error.message);
      }
    });
    reset.addEventListener('click', async () => {
      if (reset.disabled || !window.MarkataSiteTree?.openFile || !window.MarkataPublishReview?.open) return;
      reset.disabled = true;
      reset.textContent = 'Opening review…';
      try {
        await window.MarkataSiteTree.openFile(change.editor_path);
        if (window.MarkataSiteTree.activeFile?.() !== change.editor_path)
          throw new Error('The selected file is not available in the editor workspace.');
        changesReviewDialog.close();
        await window.MarkataPublishReview.open();
      } catch (error) {
        reset.disabled = false;
        reset.textContent = 'Review and reset…';
        showWorkspaceChangeUnavailable('Could not open the post review: ' + error.message);
      }
    });
    if (resetReason) reset.title = resetReason;
    if (!change.workspace) {
      showWorkspaceChangeUnavailable('This change is outside the selected site. Its path and Git status are listed here, but source content is not opened from this workspace.');
      return;
    }
    if (change.kind === 'deleted') {
      showWorkspaceChangeUnavailable('This file was deleted from the working tree. The inventory shows its Git status, but there is no saved working copy to compare here.');
      return;
    }
    if (change.kind === 'conflict') {
      showWorkspaceChangeUnavailable('This file has an unresolved Git conflict. Resolve the conflict before reviewing or resetting its source.');
      return;
    }
    if (!change.editor_path) {
      showWorkspaceChangeUnavailable('The current file path is not available in this workspace.');
      return;
    }
    try {
      const response = await fetch('/workspace/publish/review?path=' + encodeURIComponent(change.editor_path), {cache:'no-store'});
      const data = await response.json();
      if (generation !== changesReviewGeneration || changesReviewPath !== path) return;
      if (!response.ok || typeof data.current !== 'string' || typeof data.baseline !== 'string')
        throw new Error(data.error || 'Comparison unavailable.');
      detail.querySelector('#workspace-review-load').hidden = true;
      if (generation !== changesReviewGeneration || !changesReviewDialog?.open) return;
      renderWorkspaceChangeDiff(change, data);
    } catch (error) {
      if (generation !== changesReviewGeneration || changesReviewPath !== path) return;
      detail.querySelector('#workspace-review-load').hidden = true;
      reset.disabled = true;
      reset.title = 'Reset is unavailable because this file could not be compared safely.';
      detail.querySelector('#workspace-review-reset-note').textContent = 'Reset is unavailable because the saved source could not be compared safely.';
      showWorkspaceChangeUnavailable('Could not compare this file: ' + error.message);
    }
  }
  async function openChangesReview(path) {
    if (!changesSnapshot) {
      await refreshWorkspaceChanges();
      if (!changesSnapshot) return;
    }
    changesReviewed.clear();
    changesReviewFilter = 'all';
    changesReviewDialog.querySelector('#workspace-review-search').value = '';
    for (const button of changesReviewDialog.querySelectorAll('[data-review-filter]')) {
      const selected = button.dataset.reviewFilter === 'all';
      button.setAttribute('aria-pressed', String(selected));
    }
    changesReviewPath = path || changesSnapshot.entries[0]?.path || null;
    renderChangesReviewQueue();
    changesReviewDialog.showModal();
    if (changesReviewPath) {
      selectWorkspaceChange(changesReviewPath);
      changesReviewQueue.querySelector('[aria-current="true"]')?.focus();
    }
    else showWorkspaceChangeUnavailable('There are no workspace changes to review. Unpushed commits may still exist.');
  }
  async function refreshWorkspaceChanges(force = false) {
    if (!changesPanel) return;
    const generation = ++changesGeneration;
    const summary = changesPanel.querySelector('#workspace-changes-summary');
    if (changesOpen() && !changesList.children.length) summary.textContent = 'Checking for changes…';
    try {
      const response = await fetch('/workspace/changes', {cache:'no-store'});
      const data = await response.json();
      if (generation !== changesGeneration) return;
      if (!response.ok || !Array.isArray(data.entries)) throw new Error(data.error || 'Change inventory unavailable');
      const changed = !changesSnapshot || changeInventorySignature(data) !== changeInventorySignature(changesSnapshot);
      if (changed || force) changesReviewed.clear();
      if (changesOpen()) {
        showWorkspaceChanges(data, force);
      } else {
        changesSnapshot = data;
        setChangesCount(data.total ? String(data.total) : '', data.total ? data.total + ' with changes' : '');
        if (changesReviewDialog?.open) {
          renderChangesReviewQueue();
          if (changesReviewPath && (changed || force)) selectWorkspaceChange(changesReviewPath);
        }
      }
    } catch (error) {
      if (generation !== changesGeneration) return;
      if (changesOpen()) summary.textContent = 'Cannot list all changes: ' + error.message;
      setChangesCount('!', 'cannot list changes');
      if (changesOpen()) changesList.replaceChildren();
      if (changesReviewDialog?.open)
        changesReviewDialog.querySelector('#workspace-review-progress').textContent = 'Refresh failed · ' + error.message;
    }
  }
  function attachWorkspaceChanges() {
    changesHost = document.querySelector('#site-changes');
    if (!changesHost) return;
    changesButton = document.querySelector('#workspace-tab-changes');
    if (!changesButton) return;
    changesPanel = document.createElement('section');
    changesPanel.id = 'workspace-changes-panel';
    changesPanel.hidden = true;
    changesPanel.setAttribute('aria-label', 'Workspace changes');
    changesPanel.innerHTML =
      '<header class="workspace-changes-header"><div class="workspace-changes-heading"><small>WORKTREE</small><h2>Changes</h2></div>' +
      '<div class="workspace-changes-tools"><button type="button" id="workspace-changes-review-open">Review changes</button>' +
      '<button type="button" id="workspace-changes-refresh" aria-label="Refresh changes" title="Refresh changes">↻</button>' +
      '<button type="button" id="workspace-changes-close" aria-label="Close changes panel" title="Close">×</button></div></header>' +
      '<div class="workspace-changes-intro"><p id="workspace-changes-summary" role="status">Checking all Git changes…</p>' +
      '<p class="workspace-changes-help">Review saved files against the last commit. Keeping leaves them in the worktree.</p>' +
      '<div class="workspace-changes-stats" aria-label="Change status counts">' +
      '<div><strong id="workspace-changes-staged">0</strong><span>Staged</span></div>' +
      '<div><strong id="workspace-changes-unstaged">0</strong><span>Unstaged</span></div>' +
      '<div><strong id="workspace-changes-untracked">0</strong><span>Untracked</span></div></div>' +
      '<label class="workspace-changes-search"><span>Find a file</span><input id="workspace-changes-search" type="search" placeholder="Filter by path" autocomplete="off"></label>' +
      '<div class="workspace-changes-filters" role="group" aria-label="Filter change list">' +
      '<button type="button" data-inventory-filter="all" aria-pressed="true">All</button>' +
      '<button type="button" data-inventory-filter="posts" aria-pressed="false">Posts</button>' +
      '<button type="button" data-inventory-filter="other" aria-pressed="false">Other</button></div></div>' +
      '<div id="workspace-changes-list" role="group" aria-label="Changed files"></div>';
    changesHost.append(changesPanel);
    changesList = changesPanel.querySelector('#workspace-changes-list');
    document.querySelector('#site-view').addEventListener('change', syncChangesView);
    changesPanel.querySelector('#workspace-changes-review-open')
      .addEventListener('click', () => openChangesReview());
    changesPanel.querySelector('#workspace-changes-refresh')
      .addEventListener('click', () => refreshWorkspaceChanges(true));
    changesPanel.querySelector('#workspace-changes-close')
      .addEventListener('click', () => changesVisible(false));
    changesPanel.querySelector('#workspace-changes-search').addEventListener('input', renderInventoryRows);
    changesPanel.querySelectorAll('[data-inventory-filter]').forEach(button => button.addEventListener('click', () => {
      changesFilter = button.dataset.inventoryFilter;
      changesPanel.querySelectorAll('[data-inventory-filter]').forEach(other =>
        other.setAttribute('aria-pressed', String(other === button)));
      renderInventoryRows();
    }));
    changesReviewDialog = document.createElement('dialog');
    changesReviewDialog.id = 'workspace-changes-review';
    changesReviewDialog.setAttribute('aria-labelledby', 'workspace-review-title');
    changesReviewDialog.innerHTML =
      '<header class="workspace-review-header"><div><small>WORKSPACE REVIEW</small><h2 id="workspace-review-title">Review changes</h2>' +
      '<p>Compare saved files with the last committed version. This review does not publish anything.</p></div>' +
      '<div class="workspace-review-header-tools"><span id="workspace-review-position">0 / 0</span><span id="workspace-review-progress">0 kept · 0 to review</span>' +
      '<button type="button" id="workspace-review-refresh" title="Refresh inventory">Refresh</button>' +
      '<button type="button" id="workspace-review-close" aria-label="Close review">×</button></div></header>' +
      '<div class="workspace-review-layout"><aside class="workspace-review-sidebar" aria-label="Changed files">' +
      '<label class="workspace-review-search"><span>Find a file</span><input id="workspace-review-search" type="search" placeholder="Filter by path" autocomplete="off"></label>' +
      '<div class="workspace-review-filters" role="group" aria-label="Filter review queue">' +
      '<button type="button" data-review-filter="all" aria-pressed="true">All</button>' +
      '<button type="button" data-review-filter="posts" aria-pressed="false">Posts</button>' +
      '<button type="button" data-review-filter="other" aria-pressed="false">Other</button></div>' +
      '<div id="workspace-review-queue" role="group" aria-label="Changes to review"></div></aside>' +
      '<section class="workspace-review-detail" aria-label="Selected change">' +
      '<header class="workspace-review-detail-head"><div><small>SELECT A FILE</small><h3>Review workspace changes</h3><p>Choose a file to compare it with the last commit.</p></div></header>' +
      '<div class="workspace-review-baselines"><span id="workspace-review-baseline">Last committed version</span><span id="workspace-review-working">Saved workspace version</span><span id="workspace-review-diff-count"></span></div>' +
      '<div class="workspace-review-diff-wrap" id="workspace-review-diff-wrap" hidden><div class="workspace-review-diff" id="workspace-review-diff" role="region" tabindex="0" aria-label="Source diff"></div></div>' +
      '<p class="workspace-review-unavailable" id="workspace-review-unavailable" role="status">No file selected.</p>' +
      '<footer class="workspace-review-actions"><div><strong>Keep leaves this saved file unchanged.</strong><span>Reset opens the full post review and still asks you to confirm.</span></div>' +
      '<button type="button" id="workspace-review-keep" disabled>Keep in workspace</button>' +
      '<button type="button" id="workspace-review-reset" disabled>Review and reset…</button></footer></section></div>';
    document.body.append(changesReviewDialog);
    changesReviewQueue = changesReviewDialog.querySelector('#workspace-review-queue');
    changesReviewDetail = changesReviewDialog.querySelector('.workspace-review-detail');
    changesReviewDialog.querySelector('#workspace-review-close').addEventListener('click', () => changesReviewDialog.close());
    changesReviewDialog.querySelector('#workspace-review-refresh').addEventListener('click', () => refreshWorkspaceChanges(true));
    changesReviewDialog.querySelector('#workspace-review-search').addEventListener('input', () => {
      changesReviewPath = null;
      renderChangesReviewQueue();
      if (changesReviewPath) selectWorkspaceChange(changesReviewPath);
      else showWorkspaceChangeUnavailable('No files match this filter.');
    });
    changesReviewDialog.querySelectorAll('[data-review-filter]').forEach(button => button.addEventListener('click', () => {
      changesReviewFilter = button.dataset.reviewFilter;
      changesReviewDialog.querySelectorAll('[data-review-filter]').forEach(other =>
        other.setAttribute('aria-pressed', String(other === button)));
      changesReviewPath = null;
      renderChangesReviewQueue();
      if (changesReviewPath) selectWorkspaceChange(changesReviewPath);
      else showWorkspaceChangeUnavailable('No files match this filter.');
    }));
    changesReviewQueue.addEventListener('keydown', event => {
      if (!['ArrowUp', 'ArrowDown'].includes(event.key)) return;
      const rows = [...changesReviewQueue.querySelectorAll('.workspace-review-item')];
      const index = rows.indexOf(event.target.closest('.workspace-review-item'));
      if (index < 0) return;
      event.preventDefault();
      const next = rows[Math.max(0, Math.min(rows.length - 1, index + (event.key === 'ArrowDown' ? 1 : -1)))];
      next.focus();
      selectWorkspaceChange(next.dataset.path);
    });
    changesReviewDialog.addEventListener('close', () => {
      changesReviewGeneration++;
      changesReviewPath = null;
      changesReviewed.clear();
      window.MarkataPierreDiff?.release?.(changesReviewDetail.querySelector('#workspace-review-diff'));
      refreshWorkspaceChanges();
    });
    window.addEventListener('focus', () => {
      refreshWorkspaceChanges(true);
      schedulePublishStateRefresh();
    });
    window.addEventListener('markata-workspace-file-saved', () => {
      refreshWorkspaceChanges(true);
      schedulePublishStateRefresh();
    });
    let persisted = false;
    try { persisted = localStorage.getItem('markata-workspace-changes-open') === 'true'; }
    catch { /* Default closed when storage is blocked. */ }
    if (persisted) changesVisible(true);
    else {
      syncChangesView();
      refreshWorkspaceChanges();
    }
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
    '<div id="publish-view-pierre" role="tabpanel" aria-labelledby="publish-tab-pierre" data-publish-view="pierre" hidden><div class="publish-pierre-actions"><span>Pierre source diff</span><button type="button" data-pierre-layout="unified" class="active">Unified</button><button type="button" data-pierre-layout="split">Split</button></div><p id="publish-pierre-error" role="status" hidden></p><div id="publish-pierre-host"></div></div>'+
    '<div id="publish-view-source" role="tabpanel" aria-labelledby="publish-tab-source" data-publish-view="source" hidden><div id="publish-source-diff" role="region" tabindex="0" aria-label="Unified Markdown source diff"></div></div>'+
    '<div id="publish-view-message" role="tabpanel" aria-labelledby="publish-tab-message" data-publish-view="message" hidden><label for="publish-commit-message">Commit message</label><textarea id="publish-commit-message" rows="4" maxlength="4000" spellcheck="true"></textarea><p class="publish-message-help"><span id="publish-message-count">0 / 4000 bytes</span> Used for the local commit. Nothing is pushed from this screen.</p></div></div>'+
    '<div id="publish-discard-warning" hidden role="status"></div>'+
    '<div id="publish-discard-confirmation" hidden role="alertdialog" aria-label="Confirm discard"><p id="publish-discard-description">Restore this document to its last committed version? This removes your reviewed unpublished changes, without changing Git history.</p><button type="button" id="publish-discard-cancel">Cancel</button><button type="button" id="publish-discard-accept">Discard reviewed changes</button></div>'+
    '<footer class="publish-actions"><span>A Git commit is not a verified deployment.</span><button type="button" id="publish-discard" hidden>Discard changes…</button><button type="button" id="publish-back">Back to editing</button><button type="button" id="publish-mark-public" hidden>Mark as public</button><button type="button" id="publish-confirm" disabled>Commit selected page</button><button type="button" id="publish-commit-and-push" disabled>Commit &amp; Push</button><button type="button" id="publish-push" hidden>Push to origin</button><button type="button" id="publish-live-check" hidden>Check live status</button></footer>';
    document.body.append(doc);dialog=doc;notice=doc.querySelector('#publish-notice');panels=doc.querySelector('#publish-body');
    button.addEventListener('click',quickPublish);
    // Keep detailed review discoverable without making it a publish prerequisite.
    const reviewButton=document.createElement('button');
    reviewButton.type='button';
    reviewButton.id='publish-review-open';
    reviewButton.className='workspace-active-review';
    reviewButton.textContent='Review';
    reviewButton.setAttribute('aria-label','Review changes for the active post');
    reviewButton.title='Advanced review for the active post (Ctrl/Cmd+Shift+Enter)';
    reviewButton.addEventListener('click',open);
    publishStatus=document.createElement('span');
    publishStatus.id='markata-publish-status';
    publishStatus.setAttribute('role','status');
    publishStatus.setAttribute('aria-live','polite');
    publishStatus.textContent='Not published';
    document.querySelector('#saved')?.insertAdjacentElement('afterend',publishStatus);
    liveCheckButton=document.createElement('button');
    liveCheckButton.type='button';
    liveCheckButton.id='markata-live-check';
    liveCheckButton.textContent='Check live status';
    liveCheckButton.setAttribute('aria-label','Check the last pushed revision for deployment status');
    liveCheckButton.hidden=true;
    liveCheckButton.addEventListener('click',checkLiveDeployment);
    publishStatus.insertAdjacentElement('afterend',liveCheckButton);
    attachWorkspaceChanges();
    const changesHeader=changesPanel?.querySelector('header');
    const changesTools=changesPanel?.querySelector('.workspace-changes-tools');
    if(changesTools)changesTools.insertBefore(reviewButton,changesTools.querySelector('#workspace-changes-refresh'));
    else button.insertAdjacentElement('afterend',reviewButton);
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
    doc.querySelector('#publish-commit-and-push').addEventListener('click',commitAndPushReview);
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
  window.MarkataPublishReview={attach,open,quickPublish,saveAsDraft};
})();
