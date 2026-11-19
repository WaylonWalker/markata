// Shared by the service worker and editor; no HTML interpretation.
(function(root) {
  'use strict';
  const allowed = new Set(['image/png','image/jpeg','image/gif','image/webp']);
  function safeURL(value) {
    if (typeof value !== 'string') return '';
    try {
      const parsed = new URL(value.trim());
      return ['http:','https:'].includes(parsed.protocol) && parsed.hostname && !parsed.username && !parsed.password ? parsed.href : '';
    } catch { return ''; }
  }
  function normalizeSharedText(share) {
    const title = String(share.title || '').trim().slice(0, 500);
    const text = String(share.text || '').trim().slice(0, 8192);
    const direct = safeURL(String(share.url || ''));
    const match = text.match(/https?:\/\/[^\s<>"']+/i);
    const inferred = match ? safeURL(match[0].replace(/[),.;!?]+$/, '')) : '';
    return {title, text, url: direct || inferred};
  }
  function validateImages(images) {
    if (!Array.isArray(images) || images.length > 6) throw new Error('Share at most six images.');
    let total = 0;
    for (const image of images) {
      if (!image || !allowed.has(image.type) || !Number.isFinite(image.size))
        throw new Error('Only PNG, JPEG, GIF and WebP images are supported.');
      if (image.size > 15*1024*1024) throw new Error('Each image must be 15 MiB or smaller.');
      total += image.size;
    }
    if (total > 40*1024*1024) throw new Error('Shared images must total 40 MiB or less.');
    return images;
  }
  const api = {safeURL, normalizeSharedText, validateImages};
  root.MarkataShareUtils = api;
  if (typeof module !== 'undefined' && module.exports) module.exports = api;
})(typeof globalThis !== 'undefined' ? globalThis : this);
