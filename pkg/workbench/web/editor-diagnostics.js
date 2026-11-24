/* Bounded, local measurements. Never store source text, paths, or search queries. */
(() => {
  const samples = new Map();
  let latestPreviewPipeline = null;
  let peakWasmMemoryBytes = 0;
  function record(name, milliseconds) {
    if (!Number.isFinite(milliseconds) || milliseconds < 0) return;
    const values=samples.get(name)||[]; values.push(milliseconds); if(values.length>100)values.shift(); samples.set(name,values);
  }
  function recordWasmMemory(bytes) {
    if (!Number.isSafeInteger(bytes) || bytes <= 0) return;
    peakWasmMemoryBytes = Math.max(peakWasmMemoryBytes, bytes);
    renderPipeline();
  }
  function snapshot() {
    const measurements={};
    for(const [name,values] of samples) {
      const sorted=[...values].sort((a,b)=>a-b);
      measurements[name]={count:values.length,last_ms:values.at(-1),p50_ms:sorted[Math.floor((sorted.length-1)*.5)],p95_ms:sorted[Math.ceil((sorted.length-1)*.95)]};
    }
    return {schema_version:1, startup:window.MarkataStartupStatus?.snapshot() || null,
      measurements, preview_pipeline:latestPreviewPipeline,
      wasm_worker_memory_bytes:peakWasmMemoryBytes || null};
  }
  function recordPreviewPipeline(stages, details) {
    latestPreviewPipeline={stages,details,recorded_at:new Date().toISOString()};
    for(const [name,value] of Object.entries(stages)) record(`preview_${name.replace('_ms','')}`,value);
    renderPipeline();
  }
  function renderPipeline() {
    const host=document.querySelector('#preview-pipeline');
    if(!host||!latestPreviewPipeline)return;
    const {stages,details}=latestPreviewPipeline;
    const transportLabel=details.engine==='wasm'?'Worker message / response':'Network / response';
    const entries=[['Waiting',stages.waiting_ms,'waiting'],['Rendering',stages.rendering_ms,'rendering'],[transportLabel,stages.network_ms,'network'],['Display',stages.display_ms,'display']];
    const total=entries.reduce((sum,[,value])=>sum+Math.max(0,value),0);
    const bar=host.querySelector('.pipeline-bar');
    bar.replaceChildren(...entries.map(([label,value,key])=>{
      const segment=document.createElement('span');segment.className=`pipeline-segment is-${key}`;
      const share=total>0?Math.max(0,value)/total*100:0;
      segment.style.width=`${share}%`;segment.setAttribute('aria-hidden','true');
      return segment;
    }));
    bar.setAttribute('aria-label',entries.map(([label,value])=>`${label}: ${value.toFixed(1)} milliseconds`).join(', '));
    host.querySelector('.pipeline-total').textContent=`${Math.round(total)} ms`;
    const rows=entries.map(([label,value,key])=>{
      const item=document.createElement('div');item.className='pipeline-row';
      const swatch=document.createElement('span');swatch.className=`pipeline-swatch is-${key}`;swatch.setAttribute('aria-hidden','true');
      const name=document.createElement('span');name.textContent=label;
      const amount=document.createElement('strong');amount.textContent=`${value.toFixed(1)} ms`;
      item.append(swatch,name,amount);return item;
    });
    const rendering=host.querySelector('.pipeline-render-detail');
    const workerMemory=details.wasm_memory_bytes||peakWasmMemoryBytes;
    const memoryText=workerMemory?` · WASM linear memory ${(workerMemory/1048576).toFixed(1)} MiB`:'';
    rendering.textContent=details.engine==='wasm' ? `${details.mode||'Normal'} mode · local WASM ${details.wasm_render_ms.toFixed(1)} ms · no render HTTP request${memoryText}` : `${details.mode||'Normal'} mode · server queue ${details.queue_ms.toFixed(1)} ms · Go render ${details.go_render_ms.toFixed(1)} ms${details.cached?' · cached':''}`;
    const memoryDetail=host.querySelector('.pipeline-memory-detail');
    memoryDetail.textContent=peakWasmMemoryBytes?`Worker memory high-water mark: ${(peakWasmMemoryBytes/1048576).toFixed(1)} MiB`:'Worker memory is not in use';
    const iframe=details.iframe;
    const display=host.querySelector('.pipeline-display-detail');
    display.textContent=iframe?`Iframe ${iframe.total_ms.toFixed(1)} ms · HTML parse ${iframe.html_parse_ms.toFixed(1)} ms · DOM patch ${iframe.dom_patch_ms.toFixed(1)} ms · features ${iframe.features_ms.toFixed(1)} ms`:Number.isFinite(details.frame_content_ready_ms)?`Initial frame content ${details.frame_content_ready_ms.toFixed(1)} ms · feature startup continues separately`: 'Iframe stage detail unavailable';
    host.querySelector('.pipeline-note').textContent=details.engine==='wasm' ? 'Local mode uses a worker message round trip, not an HTTP render request. Worker initialization is measured separately.' : 'Waiting includes edit debounce. Network / response is the HTTP round trip after measured server queue and Go rendering; it includes response encoding and browser JSON decoding.';
    const scroll=iframe?.scroll_state;
    host.querySelector('.pipeline-scroll-detail').textContent=scroll?`Scroll sync ${scroll.total_ms.toFixed(1)} ms · layout ${scroll.layout_flush_ms.toFixed(1)} ms · blocks ${scroll.block_geometry_ms.toFixed(1)} ms · anchors ${scroll.nested_anchor_ms.toFixed(1)} ms · map ${scroll.mapping_ms.toFixed(1)} ms · caret ${scroll.marker_ms.toFixed(1)} ms (rect ${scroll.caret_rect_ms.toFixed(1)} · line ${scroll.line_height_ms.toFixed(1)} · text ${scroll.text_position_ms.toFixed(1)} · style ${scroll.caret_style_ms.toFixed(1)} ms)`:'';
    host.querySelector('.pipeline-rows').replaceChildren(...rows);
  }
  function attach() {
    const panel=document.createElement('aside'); panel.id='diagnostics-panel';panel.hidden=true;panel.setAttribute('aria-label','Performance diagnostics');
    panel.innerHTML='<div class="diagnostics-heading"><h2>Performance diagnostics</h2><button type="button" aria-label="Close diagnostics">Close</button></div><p>Local measurements from this session. No content or search terms are collected.</p><div class="diagnostics-mode"><label for="diagnostics-preview-mode">Preview speed</label><select id="diagnostics-preview-mode"><option value="auto">Auto</option><option value="fast">Fast</option><option value="normal">Normal</option><option value="slow">Slow</option></select></div><section id="preview-pipeline" aria-labelledby="preview-pipeline-title"><div class="pipeline-heading"><h3 id="preview-pipeline-title">Latest preview pipeline</h3><strong class="pipeline-total">Waiting for a preview</strong></div><div class="pipeline-bar" role="img" aria-label="Preview timing stages"></div><div class="pipeline-rows"></div><p class="pipeline-render-detail"></p><p class="pipeline-memory-detail" aria-live="polite"></p><p class="pipeline-display-detail"></p><p class="pipeline-scroll-detail"></p><p class="pipeline-note">Waiting includes the edit debounce. Network / response is the HTTP round trip after measured server queue and Go rendering; it includes response encoding and browser JSON decoding.</p></section><pre id="diagnostics-live"></pre><button id="diagnostics-copy" type="button">Copy full report</button><span id="diagnostics-copy-status" role="status"></span>';
    document.body.append(panel);
    const startup=document.querySelector('#startup-index');if(startup)panel.insertBefore(startup,panel.querySelector('pre'));
    const toggle=document.createElement('button');toggle.id='diagnostics-toggle';toggle.type='button';toggle.textContent='Diagnostics';toggle.setAttribute('aria-expanded','false');toggle.setAttribute('aria-controls',panel.id);document.querySelector('#workbench-footer').append(toggle);
    function render(){renderPipeline();panel.querySelector('pre').textContent=JSON.stringify(snapshot().measurements,null,2);}
    const modeSelect=panel.querySelector('#diagnostics-preview-mode');
    const settingsMode=document.querySelector('#setting-preview-mode');
    if(settingsMode){modeSelect.value=settingsMode.value;settingsMode.addEventListener('change',()=>{modeSelect.value=settingsMode.value;});}
    modeSelect.addEventListener('change',()=>{
      if(!settingsMode)return;
      settingsMode.value=modeSelect.value;
      settingsMode.dispatchEvent(new Event('change',{bubbles:true}));
    });
    let refresh;
    function open(){panel.hidden=false;toggle.setAttribute('aria-expanded','true');render();clearInterval(refresh);refresh=setInterval(render,1000);panel.querySelector('button').focus();}
    function close(){panel.hidden=true;toggle.setAttribute('aria-expanded','false');clearInterval(refresh);toggle.focus();}
    toggle.addEventListener('click',()=>panel.hidden?open():close());panel.querySelector('button').addEventListener('click',close);
    panel.addEventListener('keydown',e=>{if(e.key==='Escape'){e.preventDefault();close();}});
    panel.querySelector('#diagnostics-copy').addEventListener('click',async()=>{
      try{await navigator.clipboard.writeText(JSON.stringify(snapshot(),null,2));panel.querySelector('#diagnostics-copy-status').textContent='Report copied';}
      catch{panel.querySelector('#diagnostics-copy-status').textContent='Clipboard unavailable. Use MarkataDiagnostics.snapshot() in developer tools.';}
    });
    window.MarkataDiagnostics.open=open;
    window.MarkataShortcuts.register({id:'workbench.diagnostics',label:'Toggle Performance Diagnostics',category:'Diagnostics',run:()=>panel.hidden?open():close()});
    try {
      const observer=new PerformanceObserver(list=>{
        for(const entry of list.getEntries()) {
          const url=new URL(entry.name);
          // Resource timings measure network requests, not acknowledgement of durable saves.
          if(url.pathname.startsWith('/workspace/publish/'))record('publish_request',entry.duration);
          if(url.pathname==='/workspace/file')record('file_request',entry.duration);
        }
      });observer.observe({type:'resource',buffered:true});window.addEventListener('pagehide',()=>{observer.disconnect();clearInterval(refresh);},{once:true});
    }catch{/* Resource timing is optional. */}
  }
  window.MarkataDiagnostics={record,recordPreviewPipeline,recordWasmMemory,snapshot,attach};
})();
