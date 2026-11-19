// Mobile share inbox; incoming items remain drafts until Create file is confirmed.
(function() {
  'use strict';
  const NAME='markata-share-inbox', STORE='shares';
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
  async function refresh() {
    const root=ensureBanner();
    try {
      const rows=(await pending()).sort((a,b)=>a.receivedAt-b.receivedAt);
      if(!rows.length){root.hidden=true;return;}
      const first=rows[0], summary=document.createElement('span');
      summary.textContent=(rows.length>1?rows.length+' shared items saved. ':'Shared item saved. ')+
        (first.image?'Photo':first.url?'Link':'Text')+' ready to make into a post.';
      root.replaceChildren(summary,
        button('Create draft',()=>window.dispatchEvent(new CustomEvent('markata-share-select',{detail:first}))),
        button('Discard',async()=>{await remove(first.id);await refresh();}));
      root.hidden=false;
    } catch {
      root.textContent='Could not read the local share inbox; nothing has been published.';
      root.hidden=false;
    }
  }
  window.MarkataMobileShare={consume:async id=>{await remove(id);await refresh();},refresh};
  // Chrome exposes a native app-install prompt only when installability is met.
  // Safari uses the Add to Home Screen action in its Share menu instead.
  let installPrompt;
  const installButton=document.createElement('button');
  installButton.id='mobile-install';
  installButton.type='button';
  installButton.textContent='Install app';
  installButton.hidden=true;
  document.querySelector('body > header')?.append(installButton);
  window.addEventListener('beforeinstallprompt',event=>{
    event.preventDefault();
    installPrompt=event;
    installButton.hidden=false;
  });
  window.addEventListener('appinstalled',()=>{installButton.hidden=true;installPrompt=null;});
  installButton.addEventListener('click',async()=>{
    if(installPrompt){
      const current=installPrompt;
      installPrompt=null;
      installButton.hidden=true;
      await current.prompt();
    } else if(/iPhone|iPad|iPod/.test(navigator.userAgent)) {
      alert('In Safari, open Share and choose Add to Home Screen.');
    }
  });
  if(/iPhone|iPad|iPod/.test(navigator.userAgent) &&
     !window.navigator.standalone && !window.matchMedia('(display-mode: standalone)').matches) {
    installButton.hidden=false;
  }
  if('serviceWorker' in navigator && window.isSecureContext)
    navigator.serviceWorker.register('/editor/share-sw.js',{scope:'/'}).catch(err=>console.warn('Markata share target registration failed:',err));
  if(new URLSearchParams(location.search).has('share_error')) {
    const root=ensureBanner();
    root.textContent='Could not save the shared item. Try sharing again; nothing was published.';
    root.hidden=false;
  } else refresh();
})();
