/* Bounded, local measurements. Never store source text, paths, or search queries. */
(() => {
  const samples = new Map();
  function record(name, milliseconds) {
    if (!Number.isFinite(milliseconds) || milliseconds < 0) return;
    const values=samples.get(name)||[]; values.push(milliseconds); if(values.length>100)values.shift(); samples.set(name,values);
  }
  function snapshot() {
    const measurements={};
    for(const [name,values] of samples) {
      const sorted=[...values].sort((a,b)=>a-b);
      measurements[name]={count:values.length,last_ms:values.at(-1),p50_ms:sorted[Math.floor((sorted.length-1)*.5)],p95_ms:sorted[Math.ceil((sorted.length-1)*.95)]};
    }
    return {schema_version:1, startup:window.MarkataStartupStatus?.snapshot() || null, measurements};
  }
  function attach() {
    const panel=document.createElement('aside'); panel.id='diagnostics-panel';panel.hidden=true;panel.setAttribute('aria-label','Performance diagnostics');
    panel.innerHTML='<div class="diagnostics-heading"><h2>Performance diagnostics</h2><button type="button" aria-label="Close diagnostics">Close</button></div><p>Local measurements from this session. No content or search terms are collected.</p><pre id="diagnostics-live"></pre><button id="diagnostics-copy" type="button">Copy full report</button><span id="diagnostics-copy-status" role="status"></span>';
    document.body.append(panel);
    const startup=document.querySelector('#startup-index');if(startup)panel.insertBefore(startup,panel.querySelector('pre'));
    const toggle=document.createElement('button');toggle.id='diagnostics-toggle';toggle.type='button';toggle.textContent='Diagnostics';toggle.setAttribute('aria-expanded','false');toggle.setAttribute('aria-controls',panel.id);document.querySelector('footer').append(toggle);
    function render(){panel.querySelector('pre').textContent=JSON.stringify(snapshot().measurements,null,2);}
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
  window.MarkataDiagnostics={record,snapshot,attach};
})();
