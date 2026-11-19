/* Share target only: do not cache private workspace source, auth or editor HTML. */
importScripts('/editor/share-target-utils.js');
const DB_NAME = 'markata-share-inbox';
const STORE = 'shares';
const ROOT = new URL('/', self.registration.scope).href;
function openInbox() {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, 1);
    request.onupgradeneeded = () => {
      if (!request.result.objectStoreNames.contains(STORE))
        request.result.createObjectStore(STORE, {keyPath:'id'});
    };
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
}
async function storeShares(shares) {
  const db = await openInbox();
  try {
    await new Promise((resolve, reject) => {
      const tx = db.transaction(STORE, 'readwrite');
      const object = tx.objectStore(STORE);
      for (const share of shares) object.add(share);
      tx.oncomplete = resolve;
      tx.onabort = () => reject(tx.error || new Error('Cannot save incoming share'));
      tx.onerror = () => reject(tx.error || new Error('Cannot save incoming share'));
    });
  } finally {db.close();}
}
self.addEventListener('install', event => event.waitUntil(self.skipWaiting()));
self.addEventListener('activate', event => event.waitUntil(self.clients.claim()));
self.addEventListener('fetch', event => {
  const url = new URL(event.request.url);
  if (url.origin !== self.location.origin || url.pathname !== '/share-target' || event.request.method !== 'POST') return;
  event.respondWith((async () => {
    try {
      const form = await event.request.formData();
      const details = self.MarkataShareUtils.normalizeSharedText({
        title:form.get('title'),text:form.get('text'),url:form.get('url')
      });
      const images = form.getAll('images').filter(value => value instanceof File);
      self.MarkataShareUtils.validateImages(images);
      if (!details.title && !details.text && !details.url && images.length === 0)
        throw new Error('No share content received');
      const batch = crypto.randomUUID();
      // Queue each photo independently; do not silently lose multi-photo shares.
      const shares = (images.length ? images : [null]).map((image,index) => ({
        id:batch+'-'+index,receivedAt:Date.now(),
        title:details.title,text:details.text,url:details.url,image,
        imageCount:images.length,imageIndex:index+1
      }));
      await storeShares(shares);
      return Response.redirect(ROOT+'?shared=1',303);
    } catch (error) {
      console.warn('Markata could not retain shared content:',error);
      return Response.redirect(ROOT+'?share_error=1',303);
    }
  })());
});
