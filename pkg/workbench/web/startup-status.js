/* Startup telemetry is local, read-only, and does not trigger indexing. */
(() => {
  'use strict';
  const startedAt = performance.now();
  const summary = document.querySelector('#startup-index-summary');
  const state = document.querySelector('#startup-index-state');
  const timings = document.querySelector('#startup-timings');
  if (!summary || !state || !timings) return;

  let workspaceAt = null, previewAt = null, indexStatus = null;
  let metadataRequestMS = null, feedsMS = null;
  let indexReadySignalled = false;
  let pollID = 0;
  const seconds = milliseconds => (milliseconds/1000).toFixed(2)+'s';
  const add = (key,value) => {
    const item = timings.querySelector('[data-startup-timing="'+key+'"]');
    if (item) item.textContent=value;
  };
  function renderStatus() {
    if (!indexStatus) return;
    const index = indexStatus;
    const count = index.files.toLocaleString();
    let message;
    switch (index.phase) {
      case 'frontmatter':
        message='Loading frontmatter for '+index.pending.toLocaleString()+' files…';
        break;
      case 'links':
        message='Building links and search metadata for '+count+' files…';
        break;
      case 'ready':
        message=metadataRequestMS===null
          ? 'Frontmatter and links ready. Preparing feed views…'
          : 'Site index ready · '+count+' files';
        break;
      case 'error':
        message='Site index encountered an error. Check the terminal log.';
        break;
      default:
        message='Files discovered ('+count+'). Waiting to load post details…';
    }
    state.textContent=message;
    summary.textContent=index.phase==='ready' ?
      (metadataRequestMS===null ? 'Preparing site views…' : 'Site index ready') :
      index.phase==='error' ? 'Site index needs attention' : 'Loading site index…';
    add('index',index.phase==='ready' ? seconds(index.last_load_ms) :
      (index.phase==='frontmatter'||index.phase==='links' ? seconds(index.elapsed_ms)+' so far' : 'Not started'));
    add('files',count);
    add('frontmatter',index.phase==='frontmatter' ? seconds(index.elapsed_ms)+' so far' :
      index.frontmatter_ms>=0 && (index.phase==='links'||index.phase==='ready') ? seconds(index.frontmatter_ms) : 'Waiting…');
    add('links',index.phase==='links' ? 'Building…' :
      index.phase==='ready' ? seconds(index.links_ms) : 'Waiting…');
    add('request',metadataRequestMS===null ? 'Loading…' : seconds(metadataRequestMS));
    add('feeds',feedsMS===null ? (metadataRequestMS===null ? 'Loading…' : 'Not reported') : seconds(feedsMS));
  }
  async function poll() {
    const id=++pollID;
    try {
      const response=await fetch('/workspace/index-status',{cache:'no-store'});
      if (!response.ok) throw new Error('status unavailable');
      const data=await response.json();
      if (id!==pollID) return;
      indexStatus=data;
      renderStatus();
      if (data.ready && !indexReadySignalled) {
        indexReadySignalled = true;
        document.dispatchEvent(new Event('markata-site-index-ready'));
      }
      // Keep the sidebar compact after the long-running operation finishes;
      // the measurements remain available by expanding the summary.
      if (data.ready && metadataRequestMS!==null) document.querySelector('#startup-index')?.removeAttribute('open');
      if (!data.ready && data.phase!=='error') setTimeout(poll,400);
    } catch {
      summary.textContent='Site index status unavailable';
      state.textContent='Site files may still load; check the editor terminal if this persists.';
    }
  }
  // The site index may finish before the HTTP response finishes assembling
  // configured feeds. Browser resource timing captures that second stage
  // without wrapping the site's fetch or slowing down the metadata response.
  try {
    const observer = new PerformanceObserver(list => {
      for (const entry of list.getEntries()) {
        const url = new URL(entry.name);
        if (url.pathname!=='/workspace' || url.searchParams.get('metadata')!=='true') continue;
        metadataRequestMS=entry.duration;
        feedsMS=entry.serverTiming?.find(item=>item.name==='site-feeds')?.duration ?? null;
        if (indexStatus) renderStatus();
        if (indexStatus?.ready) document.querySelector('#startup-index')?.removeAttribute('open');
      }
    });
    observer.observe({type:'resource', buffered:true});
    window.addEventListener('pagehide',()=>observer.disconnect(),{once:true});
  } catch { /* Older browsers still show local index phases and counts. */ }

  const api = {
    workspaceReady() {
      if (workspaceAt!==null) return;
      workspaceAt=performance.now()-startedAt;
      add('workspace',seconds(workspaceAt));
    },
    previewReady() {
      if (previewAt!==null) return;
      previewAt=performance.now()-startedAt;
      add('preview',seconds(previewAt));
    },
    snapshot() {
      return {workspace_ready_ms:workspaceAt, first_preview_ms:previewAt,
        metadata_request_ms:metadataRequestMS, site_feeds_ms:feedsMS,
        site_index:indexStatus ? {...indexStatus} : null};
    }
  };
  window.MarkataStartupStatus=api;
  const copy = document.querySelector('#startup-copy');
  copy?.addEventListener('click', async () => {
    try {
      await navigator.clipboard.writeText(JSON.stringify(api.snapshot(),null,2));
      copy.textContent='Copied diagnostics';
    } catch { copy.textContent='Copy unavailable'; }
  });
  poll();
})();
