/* Opaque-origin preview: only the parent can supply renderer results. */
(() => {
  const loaded = new Map();
  const received = new Map();
  const modules = new Set();
  let pending, updating = false, moduleSequence = 0;
  async function scriptAsset(asset) {
    if (asset.id.startsWith('generated-')) {
      if (asset.module && modules.has(asset.content)) return;
      const script = document.createElement('script');
      if (asset.module) script.type='module';
      script.textContent=asset.content;
      if (asset.module) await new Promise((resolve,reject) => {
        const event='markata-module-'+(++moduleSequence);
        const timer=setTimeout(() => {removeEventListener(event,done);reject(new Error('Preview module did not initialize'));},10000);
        function done() {clearTimeout(timer);modules.add(asset.content);resolve();}
        addEventListener(event,done,{once:true});
        script.textContent += '\nwindow.dispatchEvent(new Event('+JSON.stringify(event)+'));';
        script.onerror=() => {clearTimeout(timer);removeEventListener(event,done);reject(new Error('Preview module failed'));};
        document.body.append(script);
      });
      else document.body.append(script);
      script.remove();
      return;
    }
    if (!loaded.has(asset.id)) {
      const script=document.createElement('script');
      script.textContent=asset.content;
      document.body.append(script);script.remove();
      loaded.set(asset.id,true);
    }
  }
  async function update() {
    if (updating) return;
    updating=true;
    try {
      while(pending) {
        const {result} = pending; pending=null;
        // Destroy old chart observers before replacing their canvases.
        if (window.Chart) for (const chart of Object.values(Chart.instances)) chart.destroy();
        document.querySelector('.post-content').innerHTML=result.html;
        for (const asset of result.assets) {
          if (asset.kind !== 'style') continue;
          if (!document.querySelector(`style[data-render-id="${CSS.escape(asset.id)}"]`)) {
            const style=document.createElement('style');style.dataset.renderId=asset.id;style.textContent=asset.content;
            document.head.append(style);
          }
        }
        for (const asset of result.assets) if (asset.kind === 'script') await scriptAsset(asset);
        document.dispatchEvent(new CustomEvent('markata:content-decrypted',{detail:{container:document.querySelector('.post-content')}}));
        dispatchEvent(new Event('view-transition-complete'));
      }
    } catch (error) { parent.postMessage({type:'preview-error',error:String(error)},'*'); }
    finally { updating=false; }
  }
  addEventListener('message', e => {
    if (e.source !== parent) return;
    if (e.data?.type === 'appearance') {
      for (const [key,value] of Object.entries(e.data.attributes || {})) {
        if (!['data-theme','data-palette','data-aesthetic','data-fontpack','data-text-size'].includes(key)) continue;
        if (typeof value === 'string') document.documentElement.setAttribute(key,value);
      }
      dispatchEvent(new Event('theme-change'));
    }
    if (e.data?.type === 'render') {
      // Cache bytes before coalescing pending updates, so a skipped revision
      // cannot discard the only transfer of a newly introduced dependency.
      for (const asset of e.data.result.assets) {
        if (asset.id.startsWith('generated-')) continue;
        if (asset.content) received.set(asset.id,asset.content);
        else asset.content=received.get(asset.id) || '';
      }
      pending=e.data;update();
    }
  });
  parent.postMessage({type:'preview-ready'},'*');
})();
