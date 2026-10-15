// Reader previews are optional: a blocked hotlink or a stale cached URL
// should never leave a broken <img> icon behind.
(() => {
  function handleBrokenPreview(img) {
    if (!img.matches('img.reader-entry-image')) return;
    const link = img.closest('.reader-entry-image-link');
    const article = img.closest('.reader-entry');
    if (!link || !article || !img.isConnected) return;

    const sourceLink = article.querySelector('.reader-entry-source-link');
    const sourceName = article.querySelector('.reader-entry-source')?.textContent?.trim() || 'Source';
    const tile = document.createElement('span');
    tile.className = 'reader-entry-source-tile';
    tile.setAttribute('aria-hidden', 'true');

    const top = document.createElement('span');
    top.className = 'reader-entry-source-tile-top';
    const sourceIcon = sourceLink?.querySelector('img.reader-entry-source-icon');
    if (sourceIcon && sourceIcon.complete && sourceIcon.naturalWidth > 0) {
      const icon = sourceIcon.cloneNode(false);
      icon.className = 'reader-entry-source-tile-icon';
      top.appendChild(icon);
    } else {
      const badge = document.createElement('span');
      badge.className = 'reader-entry-source-tile-badge';
      badge.textContent = sourceName.charAt(0).toUpperCase();
      top.appendChild(badge);
    }

    const label = document.createElement('span');
    label.className = 'reader-entry-source-tile-label';
    label.textContent = sourceName;
    top.appendChild(label);
    tile.appendChild(top);

    const domain = document.createElement('span');
    domain.className = 'reader-entry-source-tile-domain';
    try {
      domain.textContent = new URL(sourceLink?.href || link.href).hostname;
    } catch {
      domain.textContent = '';
    }
    tile.appendChild(domain);

    link.className = 'reader-entry-image-link reader-entry-image-link--source-tile';
    article.classList.remove('has-image');
    img.replaceWith(tile);
  }

  // Capture image errors because they do not bubble. This also covers lazy
  // images and entries appended by HTMX or other pagination modes.
  document.addEventListener('error', (event) => {
    if (event.target instanceof HTMLImageElement) {
      handleBrokenPreview(event.target);
    }
  }, true);

  // Images may fail before this deferred script starts.
  document.querySelectorAll('img.reader-entry-image').forEach((img) => {
    if (img.complete && img.naturalWidth === 0) handleBrokenPreview(img);
  });
})();
