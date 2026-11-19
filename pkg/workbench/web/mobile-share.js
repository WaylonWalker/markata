// Mobile share inbox and quick capture. Nothing is published without the editor's
// explicit create/publish actions. Store only per-device, never in a shared cache.
(function() {
  'use strict';
  const NAME='markata-share-inbox', STORE='shares';
  const params=new URLSearchParams(location.search);
  let banner;
  function openInbox() {
    return new Promise((resolve,reject) => {
      const req=indexedDB.open(NAME,1);
      req.onupgradeneeded=()=>{
        if(!req.result.objectStoreNames.contains(STORE))
          req.result.createObjectStore(STORE,{keyPath:'id'});
      };
      req.onsuccess=()=>resolve(req.result);
      req.onerror=()=>reject(req.error);
    });
  }
  async function transact(mode,operation) {
    const db=await openInbox();
    try {
      return await new Promise((resolve,reject)=>{
        const tx=db.transaction(STORE,mode), req=operation(tx.objectStore(STORE));
        let value;
        req.onsuccess=()=>{value=req.result;};
        tx.oncomplete=()=>resolve(value);
        tx.onabort=()=>reject(tx.error||new Error('Share inbox unavailable'));
        tx.onerror=()=>reject(tx.error||new Error('Share inbox unavailable'));
      });
    } finally {db.close();}
  }
  const pending=()=>transact('readonly',store=>store.getAll());
  const remove=id=>transact('readwrite',store=>store.delete(id));
  function ensureBanner() {
    if(banner)return banner;
    banner=document.createElement('aside');
    banner.id='mobile-share-inbox';
    banner.setAttribute('aria-live','polite');
    const workbench=document.querySelector('#workbench');
    workbench.parentNode.insertBefore(banner,workbench);
    return banner;
  }
  function button(label,handler) {
    const element=document.createElement('button');
    element.type='button';
    element.textContent=label;
    element.addEventListener('click',handler);
    return element;
  }
  function announce(message) {
    const root=ensureBanner();
    root.textContent=message;
    root.hidden=false;
  }
  async function refresh() {
    const root=ensureBanner();
    try {
      const rows=(await pending()).sort((a,b)=>a.receivedAt-b.receivedAt);
      if(!rows.length){
        if(params.has('share_inbox')) {
          root.textContent='Share inbox is empty. Share a link or take a photo to get started.';
          root.hidden=false;
        } else root.hidden=true;
        return;
      }
      const first=rows[0], summary=document.createElement('span');
      summary.textContent=(rows.length>1?rows.length+' shared items saved. ':'Shared item saved. ')+
        (first.image?'Photo':first.url?'Link':'Text')+' ready to make into a post.';
      const actions=[
        summary,
        button('Create draft',()=>window.dispatchEvent(new CustomEvent('markata-share-select',{detail:first}))),
        button('Discard',async()=>{
          try {await remove(first.id);await refresh();}
          catch {announce('Could not discard this saved item. Try again.');}
        })
      ];
      if(navigator.storage && navigator.storage.persist && navigator.storage.persisted) {
        let protectedStorage=false;
        try {protectedStorage=await navigator.storage.persisted();} catch {}
        if(!protectedStorage) actions.push(button('Protect storage',async()=>{
          try {
            const granted=await navigator.storage.persist();
            if(!granted) announce('The browser cannot guarantee persistent storage. Shared items may be evicted; create or back up important drafts.');
            else await refresh();
          } catch {announce('Storage protection is unavailable on this device.');}
        }));
      }
      root.replaceChildren(...actions);
      root.hidden=false;
    } catch {
      announce('Could not read the local share inbox; nothing has been published.');
    }
  }
  window.MarkataMobileShare={consume:async id=>{await remove(id);await refresh();},refresh};
  const header=document.querySelector('body > header');

  // Capture to IndexedDB before asking the remote workspace to create a file.
  // This also works when the phone is offline, unlike the server-backed dialog.
  const cameraInput=document.createElement('input');
  cameraInput.type='file';
  cameraInput.accept='image/png,image/jpeg,image/gif,image/webp';
  cameraInput.setAttribute('capture','environment');
  cameraInput.hidden=true;
  cameraInput.setAttribute('aria-label','Take a photo for a new draft');
  document.body.append(cameraInput);
  const cameraButton=button('Take photo',()=>cameraInput.click());
  cameraButton.id='mobile-camera';
  header?.append(cameraButton);
  cameraInput.addEventListener('change',async()=>{
    const file=cameraInput.files?.[0];
    cameraInput.value='';
    if(!file)return;
    try {
      window.MarkataShareUtils.validateImages([file]);
      if(navigator.storage?.estimate) {
        const estimate=await navigator.storage.estimate();
        if(Number.isFinite(estimate.quota) && Number.isFinite(estimate.usage) &&
           estimate.quota-estimate.usage < file.size+1024*1024)
          throw new Error('Not enough device storage for this photo.');
      }
      const entry={id:crypto.randomUUID()+'-0',receivedAt:Date.now(),
        title:'',text:'',url:'',image:file,imageCount:1,imageIndex:1};
      await transact('readwrite',store=>store.add(entry));
      await refresh();
    } catch(error) {
      announce('Photo was not saved: '+(error?.message||'Local storage unavailable')+'. Nothing was published.');
    }
  });

  // Chrome exposes an install prompt only when installability is met.
  // Safari uses Add to Home Screen in its Share menu instead.
  let installPrompt;
  const installButton=button('Install app',async()=>{
    if(installPrompt){
      const current=installPrompt;
      installPrompt=null;
      installButton.hidden=true;
      await current.prompt();
    } else if(/iPhone|iPad|iPod/.test(navigator.userAgent)) {
      alert('In Safari, open Share and choose Add to Home Screen.');
    }
  });
  installButton.id='mobile-install';
  installButton.hidden=true;
  header?.append(installButton);
  window.addEventListener('beforeinstallprompt',event=>{
    event.preventDefault();
    installPrompt=event;
    installButton.hidden=false;
  });
  window.addEventListener('appinstalled',()=>{installButton.hidden=true;installPrompt=null;});
  if(/iPhone|iPad|iPod/.test(navigator.userAgent) &&
     !window.navigator.standalone && !window.matchMedia('(display-mode: standalone)').matches) {
    installButton.hidden=false;
  }

  // Never reload the editor on a background service-worker update.
  // Only the explicit Update app action activates a waiting worker.
  let updateButton;
  function offerUpdate(registration) {
    if(!registration.waiting || !navigator.serviceWorker.controller || updateButton)return;
    updateButton=button('Update app',()=>{
      if(!window.confirm('Save your current edits before reloading Markata. Update now?'))return;
      updateButton.disabled=true;
      navigator.serviceWorker.addEventListener('controllerchange',()=>location.reload(),{once:true});
      registration.waiting.postMessage({type:'MARKATA_ACTIVATE_UPDATE'});
    });
    updateButton.id='mobile-update';
    header?.append(updateButton);
  }
  if('serviceWorker' in navigator && window.isSecureContext) {
    navigator.serviceWorker.register('/editor/share-sw.js',{scope:'/'}).then(registration=>{
      offerUpdate(registration);
      registration.addEventListener('updatefound',()=>{
        const worker=registration.installing;
        worker?.addEventListener('statechange',()=>{
          if(worker.state==='installed')offerUpdate(registration);
        });
      });
    }).catch(err=>console.warn('Markata share target registration failed:',err));
  }

  if(params.has('share_error')) {
    announce('Could not save the shared item. Try sharing again; nothing was published.');
  } else {
    void refresh();
  }
  const compose=params.get('compose');
  if(['note','link','photo'].includes(compose))
    window.dispatchEvent(new CustomEvent('markata-compose',{detail:compose}));
  if(params.has('compose') || params.has('share_inbox')) {
    const clean=new URL(location.href);
    clean.searchParams.delete('compose');
    clean.searchParams.delete('share_inbox');
    history.replaceState(history.state,'',clean.pathname+clean.search+clean.hash);
  }
})();
