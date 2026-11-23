// Mobile-first CMS flows: Daily Notes, Quick Notes, Mobile PWA Share Target,
// and Cross-Link ([[wikilink]]) Picker.
(() => {
  'use strict';

  // Register installable PWA service worker (enables mobile home-screen install
  // and OS Share Target without interfering with live editing).
  if ('serviceWorker' in navigator && (location.protocol === 'https:' || location.hostname === 'localhost' || location.hostname === '127.0.0.1')) {
    navigator.serviceWorker.register('/sw.js', { scope: '/' }).catch(() => {});
  }

  function getEditor() {
    return document.querySelector('#editor-pane .CodeMirror')?.CodeMirror || null;
  }

  function todayISO() {
    const now = new Date();
    const y = now.getFullYear();
    const m = String(now.getMonth() + 1).padStart(2, '0');
    const d = String(now.getDate()).padStart(2, '0');
    return `${y}-${m}-${d}`;
  }

  function nowStamp() {
    const now = new Date();
    const hh = String(now.getHours()).padStart(2, '0');
    const mm = String(now.getMinutes()).padStart(2, '0');
    const ss = String(now.getSeconds()).padStart(2, '0');
    return `${hh}${mm}${ss}`;
  }

  function nowTime() {
    const now = new Date();
    return `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`;
  }

  function toast(message, type = 'info') {
    window.dispatchEvent(new CustomEvent('markata:publish-status', {
      detail: { message, type, pending: false },
    }));
  }

  function detectDirectory(candidates, fallback) {
    const files = window.MarkataSiteTree?.files?.() || [];
    for (const dir of candidates) {
      const prefix = dir.replace(/\/+$/, '') + '/';
      if (files.some((f) => String(f.path || '').startsWith(prefix))) {
        return dir.replace(/\/+$/, '');
      }
    }
    return fallback;
  }

  async function openOrCreateDailyNote(appendBlock = '') {
    const date = todayISO();
    const dir = detectDirectory(['daily', 'pages/daily', 'journal', 'posts/daily'], 'daily');
    const path = `${dir}/${date}.md`;
    try {
      const existing = await fetch('/workspace/file?path=' + encodeURIComponent(path), { cache: 'no-store' });
      if (existing.ok) {
        if (appendBlock) {
          const data = await existing.json();
          const current = data.source || '';
          const next = current.replace(/\s*$/, '') + '\n\n' + appendBlock.trim() + '\n';
          const putResp = await fetch('/workspace/file?path=' + encodeURIComponent(path), {
            method: 'PUT',
            headers: {
              'Content-Type': 'text/plain; charset=utf-8',
              'X-Markata-Expected-Version': data.version || '',
            },
            body: next,
          });
          if (!putResp.ok) {
            const err = await putResp.json().catch(() => ({}));
            throw new Error(err.error || 'Could not append to daily note');
          }
        }
      } else {
        const body = [
          '---',
          `title: "${date}"`,
          `date: ${date}`,
          'published: false',
          'draft: true',
          'tags:',
          '  - daily',
          '---',
          '',
          `# ${date}`,
          '',
          appendBlock ? appendBlock.trim() + '\n' : '',
        ].join('\n');
        const createResp = await fetch('/workspace/file?path=' + encodeURIComponent(path), {
          method: 'POST',
          headers: { 'Content-Type': 'text/plain; charset=utf-8' },
          body,
        });
        if (!createResp.ok) {
          const err = await createResp.json().catch(() => ({}));
          throw new Error(err.error || 'Could not create daily note');
        }
      }
      await window.MarkataSiteTree?.refresh?.();
      await window.MarkataSiteTree?.openFile?.(path);
      if (window.matchMedia('(max-width: 768px)').matches) {
        window.MarkataCMSShell?.setMobileView?.('write');
      }
      const cm = getEditor();
      if (cm) {
        const last = cm.lastLine();
        cm.setCursor({ line: last, ch: cm.getLine(last).length });
        cm.focus();
      }
      toast(appendBlock ? `Appended to ${path}` : `Opened daily note ${path}`, 'info');
      return path;
    } catch (err) {
      toast(err.message || 'Could not open daily note', 'warning');
      return null;
    }
  }

  // ── Quick Note & Mobile PWA Share Capture Dialog ─────────────────────────
  const quickDialog = document.createElement('dialog');
  quickDialog.id = 'cms-quick-note-dialog';
  quickDialog.className = 'cms-flow-dialog';
  quickDialog.setAttribute('aria-labelledby', 'cms-qn-title');
  quickDialog.innerHTML = `
    <header class="cms-flow-head">
      <div>
        <small>QUICK CAPTURE</small>
        <h2 id="cms-qn-title">Quick Note</h2>
      </div>
      <button type="button" class="cms-flow-close" aria-label="Close quick note">✕</button>
    </header>
    <form id="cms-qn-form" class="cms-flow-body">
      <div class="cms-flow-segmented" role="radiogroup" aria-label="Destination">
        <button type="button" data-qn-dest="note" aria-pressed="true">New Quick Note</button>
        <button type="button" data-qn-dest="daily" aria-pressed="false">Append to Today's Daily</button>
      </div>
      <label class="cms-flow-field">
        <span>Title <small>(optional)</small></span>
        <input id="cms-qn-heading" type="text" placeholder="Quick thought or bookmark title" autocomplete="off">
      </label>
      <label class="cms-flow-field">
        <span>URL <small>(optional — auto-fetches metadata)</small></span>
        <input id="cms-qn-url" type="url" placeholder="https://..." autocomplete="off">
      </label>
      <label class="cms-flow-field">
        <span>Note</span>
        <textarea id="cms-qn-body" rows="5" placeholder="Jot a thought, quote, or link notes…"></textarea>
      </label>
      <label class="cms-flow-field">
        <span>Tags <small>(comma-separated)</small></span>
        <input id="cms-qn-tags" type="text" placeholder="note, link" autocomplete="off">
      </label>
      <label class="cms-flow-field">
        <span>Photos <small>(up to 8 · 15 MiB each, 48 MiB total)</small></span>
        <input id="cms-qn-images" type="file" accept="image/png,image/jpeg,image/gif,image/webp" multiple>
      </label>
      <div id="cms-qn-photo-list" class="cms-qn-photo-list" aria-live="polite"></div>
      <label class="cms-flow-check">
        <input id="cms-qn-publish" type="checkbox">
        <span>Mark public &amp; publish immediately after saving</span>
      </label>
      <p id="cms-qn-status" class="cms-flow-status" role="status" hidden></p>
      <footer class="cms-flow-actions">
        <button type="button" class="cms-flow-cancel">Cancel</button>
        <button type="submit" id="cms-qn-submit" class="cms-flow-primary">Save Note</button>
      </footer>
    </form>
  `;
  document.body.append(quickDialog);

  let qnDest = 'note';
  const qnHeading = quickDialog.querySelector('#cms-qn-heading');
  const qnUrl = quickDialog.querySelector('#cms-qn-url');
  const qnBody = quickDialog.querySelector('#cms-qn-body');
  const qnTags = quickDialog.querySelector('#cms-qn-tags');
  const qnPublish = quickDialog.querySelector('#cms-qn-publish');
  const qnStatus = quickDialog.querySelector('#cms-qn-status');
  const qnSubmit = quickDialog.querySelector('#cms-qn-submit');
  const qnImageInput = quickDialog.querySelector('#cms-qn-images');
  const qnPhotoList = quickDialog.querySelector('#cms-qn-photo-list');
  let qnPhotos = [];
  let qnPhotoURLs = [];
  let pendingShareId = null;
  const SHARE_DB = 'markata-share-captures-v1';
  const SHARE_STORE = 'captures';

  function shareDatabase() {
    return new Promise((resolve, reject) => {
      const request = indexedDB.open(SHARE_DB, 1);
      request.onupgradeneeded = () => request.result.createObjectStore(SHARE_STORE, { keyPath: 'id' });
      request.onsuccess = () => resolve(request.result);
      request.onerror = () => reject(request.error);
    });
  }
  async function readStagedShare(id) {
    const db = await shareDatabase();
    try {
      return await new Promise((resolve, reject) => {
        const tx = db.transaction(SHARE_STORE, 'readonly');
        const req = tx.objectStore(SHARE_STORE).get(id);
        req.onsuccess = () => resolve(req.result || null);
        req.onerror = () => reject(req.error);
      });
    } finally { db.close(); }
  }
  async function removeStagedShare(id) {
    const db = await shareDatabase();
    try {
      await new Promise((resolve, reject) => {
        const tx = db.transaction(SHARE_STORE, 'readwrite');
        tx.objectStore(SHARE_STORE).delete(id);
        tx.oncomplete = resolve;
        tx.onerror = () => reject(tx.error);
        tx.onabort = () => reject(tx.error);
      });
    } finally { db.close(); }
  }

  function photoType(file) {
    const types = { png: 'image/png', jpg: 'image/jpeg', jpeg: 'image/jpeg',
      gif: 'image/gif', webp: 'image/webp' };
    return file.type || types[(file.name || '').split('.').pop().toLowerCase()] || '';
  }

  function updatePhotos(files) {
    const next = Array.from(files || []);
    const sum = next.reduce((total, file) => total + file.size, 0);
    if (next.length > 8 || sum > 48 * 1024 * 1024 ||
        next.some((file) => !file.size || file.size > 15 * 1024 * 1024 ||
          !['image/png','image/jpeg','image/gif','image/webp'].includes(photoType(file)))) {
      throw new Error('Choose up to 8 PNG/JPEG/GIF/WebP photos (15 MiB each, 48 MiB total).');
    }
    qnPhotoURLs.forEach((url) => URL.revokeObjectURL(url));
    qnPhotoURLs = [];
    qnPhotos = next;
    qnImageInput.value = '';
    qnPhotoList.replaceChildren();
    for (const [index, file] of qnPhotos.entries()) {
      const wrapper = document.createElement('div');
      wrapper.className = 'cms-qn-photo';
      const preview = document.createElement('img');
      const objectURL = URL.createObjectURL(file);
      qnPhotoURLs.push(objectURL);
      preview.src = objectURL;
      preview.alt = file.name || ('Shared photo ' + (index + 1));
      const remove = document.createElement('button');
      remove.type = 'button';
      remove.textContent = 'Remove';
      remove.setAttribute('aria-label', 'Remove photo ' + (index + 1));
      remove.addEventListener('click', () => {
        updatePhotos(qnPhotos.filter((_, position) => position !== index));
      });
      wrapper.append(preview, remove);
      qnPhotoList.append(wrapper);
    }
    const daily = quickDialog.querySelector('[data-qn-dest="daily"]');
    daily.disabled = qnPhotos.length > 0;
    if (qnPhotos.length && qnDest === 'daily') setQuickDest('note');
  }
  qnImageInput.addEventListener('change', () => {
    try {
      updatePhotos([...qnPhotos, ...qnImageInput.files]);
      qnStatus.hidden = true;
    } catch (error) {
      qnImageInput.value = '';
      qnStatus.hidden = false;
      qnStatus.textContent = error.message;
    }
  });

  function setQuickDest(dest) {
    qnDest = dest;
    quickDialog.querySelectorAll('[data-qn-dest]').forEach((btn) => {
      btn.setAttribute('aria-pressed', String(btn.dataset.qnDest === dest));
    });
    qnSubmit.textContent = dest === 'daily' ? 'Append to Daily' : 'Save Note';
  }
  quickDialog.querySelectorAll('[data-qn-dest]').forEach((btn) => {
    btn.addEventListener('click', () => setQuickDest(btn.dataset.qnDest));
  });
  quickDialog.querySelector('.cms-flow-close').addEventListener('click', () => quickDialog.close());
  quickDialog.querySelector('.cms-flow-cancel').addEventListener('click', () => quickDialog.close());

  async function enrichURL(url) {
    if (!url || !/^https?:\/\//i.test(url)) return;
    qnStatus.hidden = false;
    qnStatus.textContent = 'Fetching link details…';
    try {
      // Keep this in sync with app.js: the workspace metadata endpoint
      // intentionally accepts POST JSON rather than a query-string GET.
      const resp = await fetch('/workspace/link-metadata', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
        cache: 'no-store',
      });
      if (!resp.ok) {
        qnStatus.hidden = true;
        return;
      }
      const meta = await resp.json();
      if (!qnHeading.value.trim() && meta.title) qnHeading.value = meta.title;
      if (!qnBody.value.trim() && meta.description) qnBody.value = meta.description;
      qnStatus.hidden = true;
    } catch {
      qnStatus.hidden = true;
    }
  }
  qnUrl.addEventListener('blur', () => {
    if (qnUrl.value.trim() && !qnHeading.value.trim()) {
      enrichURL(qnUrl.value.trim());
    }
  });

  function slugify(text) {
    return String(text || '')
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '')
      .slice(0, 48);
  }

  function openQuickNote(prefill = {}) {
    qnHeading.value = prefill.title || '';
    qnUrl.value = prefill.url || '';
    qnBody.value = prefill.text || '';
    qnTags.value = prefill.tags || (prefill.url ? 'link' : 'note');
    qnPublish.checked = false;
    qnStatus.hidden = true;
    pendingShareId = prefill.shareId || null;
    updatePhotos(prefill.files || []);
    setQuickDest(prefill.dest || 'note');
    quickDialog.showModal();
    if (prefill.url && !prefill.title) {
      enrichURL(prefill.url);
    }
    (prefill.title || prefill.url ? qnBody : qnHeading).focus();
  }

  quickDialog.querySelector('#cms-qn-form').addEventListener('submit', async (event) => {
    event.preventDefault();
    const title = qnHeading.value.trim();
    const url = qnUrl.value.trim();
    const bodyText = qnBody.value.trim();
    const tags = qnTags.value.split(',').map((t) => t.trim()).filter(Boolean);
    const shouldPublish = qnPublish.checked;

    if (!title && !url && !bodyText && !qnPhotos.length) {
      qnStatus.hidden = false;
      qnStatus.textContent = 'Add a title, URL, or note before saving.';
      return;
    }

    qnSubmit.disabled = true;
    try {
      if (qnDest === 'daily') {
        if (qnPhotos.length) throw new Error('Photos must be saved in a new note, not appended to a daily note.');
        const heading = url
          ? `- **${nowTime()}** — [${title || url}](${url})`
          : `- **${nowTime()}**${title ? ` — **${title}**` : ''}`;
        const block = bodyText ? `${heading}\n  ${bodyText.replace(/\n/g, '\n  ')}` : heading;
        quickDialog.close();
        const savedPath = await openOrCreateDailyNote(block);
        if (!savedPath) {
          // Never publish the previously active document if daily capture fails.
          throw new Error('Could not save the daily note; nothing was published.');
        }
        if (shouldPublish) await window.MarkataPublishReview?.quickPublish?.();
        return;
      }

      const date = todayISO();
      const dir = detectDirectory(['thoughts', 'notes', 'posts', 'pages'], 'notes');
      const slug = slugify(title) || (qnPhotos.length ? 'shared-photos' : 'quick-note');
      // Titles can repeat, and share targets can arrive seconds apart; keep
      // independent captures from colliding with existing filenames.
      const unique = `${date}-${nowStamp()}-${Math.random().toString(36).slice(2, 10)}`;
      const path = `${dir}/${slug}-${unique}.md`;
      const yamlTags = tags.length ? ['tags:', ...tags.map((t) => `  - ${JSON.stringify(t)}`)].join('\n') : 'tags:\n  - note';
      const lines = [
        '---',
        `title: ${JSON.stringify(title || `Quick Note ${date}`)}`,
        `date: ${date}`,
        `published: ${shouldPublish ? 'true' : 'false'}`,
        `draft: ${shouldPublish ? 'false' : 'true'}`,
        yamlTags,
      ];
      if (url) lines.push(`link: ${JSON.stringify(url)}`);
      lines.push('---', '');
      if (url) lines.push(`[${title || url}](${url})`, '');
      if (bodyText) lines.push(bodyText, '');

      let resp;
      if (qnPhotos.length) {
        const form = new FormData();
        form.set('path', path);
        form.set('source', lines.join('\n'));
        form.set('image_name', slug + '-' + unique);
        for (const file of qnPhotos) {
          // Some Android providers omit File.type. The server still verifies
          // the actual image bytes independently.
          const upload = file.type ? file : new File([file], file.name || 'shared-photo', {
            type: photoType(file),
          });
          form.append('image', upload, upload.name || 'shared-photo');
        }
        resp = await fetch('/workspace/new-post', { method: 'POST', body: form });
      } else {
        resp = await fetch('/workspace/file?path=' + encodeURIComponent(path), {
          method: 'POST',
          headers: { 'Content-Type': 'text/plain; charset=utf-8' },
          body: lines.join('\n'),
        });
      }
      if (!resp.ok) {
        const err = await resp.json().catch(() => ({}));
        throw new Error(err.error || 'Could not create quick note');
      }
      quickDialog.close();
      if (pendingShareId) {
        try {
          await removeStagedShare(pendingShareId);
          pendingShareId = null;
          history.replaceState(null, '', location.pathname + location.hash);
        } catch {
          // The document is already saved; keeping the staged copy is safe.
          toast('Note saved, but the temporary share could not be cleared.', 'warning');
        }
      }
      await window.MarkataSiteTree?.refresh?.();
      await window.MarkataSiteTree?.openFile?.(path);
      if (window.matchMedia('(max-width: 768px)').matches) {
        window.MarkataCMSShell?.setMobileView?.('write');
      }
      toast(`Created ${path}`, 'info');
      if (shouldPublish) {
        setTimeout(() => window.MarkataPublishReview?.quickPublish?.(), 250);
      }
    } catch (err) {
      qnStatus.hidden = false;
      qnStatus.textContent = err.message || 'Could not save quick note';
    } finally {
      qnSubmit.disabled = false;
    }
  });

  // ── Cross-Link ([[wikilink]]) Fuzzy Picker ───────────────────────────────
  const linkDialog = document.createElement('dialog');
  linkDialog.id = 'cms-wikilink-dialog';
  linkDialog.className = 'cms-flow-dialog';
  linkDialog.setAttribute('aria-labelledby', 'cms-wl-title');
  linkDialog.innerHTML = `
    <header class="cms-flow-head">
      <div>
        <small>CROSS-LINK</small>
        <h2 id="cms-wl-title">Insert Wikilink</h2>
      </div>
      <button type="button" class="cms-flow-close" aria-label="Close wikilink picker">✕</button>
    </header>
    <div class="cms-flow-body">
      <input id="cms-wl-query" type="search" placeholder="Search posts by title, slug, or tag…" autocomplete="off">
      <div id="cms-wl-results" class="cms-wl-results" role="listbox" aria-label="Matching posts"></div>
    </div>
  `;
  document.body.append(linkDialog);

  const wlQuery = linkDialog.querySelector('#cms-wl-query');
  const wlResults = linkDialog.querySelector('#cms-wl-results');
  linkDialog.querySelector('.cms-flow-close').addEventListener('click', () => linkDialog.close());

  function slugFromFile(file) {
    if (file.slug) return String(file.slug).replace(/^\/+|\/+$/g, '');
    return String(file.path || '').replace(/^.*\//, '').replace(/\.md$/i, '');
  }

  function insertWikilink(file) {
    const cm = getEditor();
    if (!cm) return;
    const slug = slugFromFile(file);
    const title = (file.title || '').trim();
    const snippet = title && title.toLowerCase() !== slug.toLowerCase()
      ? `[[${slug}|${title}]]`
      : `[[${slug}]]`;
    linkDialog.close();
    if (window.matchMedia('(max-width: 768px)').matches) {
      window.MarkataCMSShell?.setMobileView?.('write');
    }
    cm.replaceSelection(snippet, 'end');
    cm.focus();
  }

  function renderWikilinkResults() {
    const q = wlQuery.value.trim().toLowerCase();
    const terms = q ? q.split(/\s+/).filter(Boolean) : [];
    const files = (window.MarkataSiteTree?.files?.() || []).filter((f) => /\.md$/i.test(f.path || ''));
    const scored = [];
    for (const file of files) {
      const hay = `${file.title || ''} ${file.slug || ''} ${file.path || ''} ${(file.tags || []).join(' ')}`.toLowerCase();
      if (terms.every((t) => hay.includes(t))) {
        const titleMatch = q && (file.title || '').toLowerCase().includes(q) ? 0 : 1;
        scored.push({ file, rank: titleMatch });
      }
      if (scored.length >= 60 && !q) break;
    }
    scored.sort((a, b) => a.rank - b.rank);
    const top = scored.slice(0, 24);
    wlResults.replaceChildren();
    if (!top.length) {
      const empty = document.createElement('p');
      empty.className = 'cms-wl-empty';
      empty.textContent = 'No matching posts found.';
      wlResults.append(empty);
      return;
    }
    top.forEach(({ file }, idx) => {
      const btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'cms-wl-item';
      btn.setAttribute('role', 'option');
      btn.setAttribute('aria-selected', idx === 0 ? 'true' : 'false');
      const title = document.createElement('strong');
      title.textContent = file.title || slugFromFile(file);
      const meta = document.createElement('span');
      meta.textContent = `[[${slugFromFile(file)}]] · ${file.path}`;
      btn.append(title, meta);
      btn.addEventListener('click', () => insertWikilink(file));
      wlResults.append(btn);
    });
  }

  function openWikilinkPicker() {
    wlQuery.value = '';
    renderWikilinkResults();
    linkDialog.showModal();
    wlQuery.focus();
  }

  wlQuery.addEventListener('input', renderWikilinkResults);
  wlQuery.addEventListener('keydown', (event) => {
    if (event.key === 'Enter') {
      const first = wlResults.querySelector('.cms-wl-item');
      if (first) {
        event.preventDefault();
        first.click();
      }
    }
  });

  // ── Quick Actions Strip in Content Sidebar ───────────────────────────────
  function injectQuickActions() {
    const sidebar = document.querySelector('#site-sidebar');
    if (!sidebar || sidebar.querySelector('#cms-quick-actions')) return;
    const strip = document.createElement('div');
    strip.id = 'cms-quick-actions';
    strip.className = 'cms-quick-actions';
    strip.setAttribute('aria-label', 'Quick writing flows');
    strip.innerHTML = `
      <button type="button" data-cms-flow="daily" title="Open or create today's daily note">Daily</button>
      <button type="button" data-cms-flow="quick" title="Capture a quick note or link">Quick Note</button>
      <button type="button" data-cms-flow="wikilink" title="Insert a [[wikilink]] to another post">[[Link]]</button>
    `;
    strip.addEventListener('click', (event) => {
      const btn = event.target.closest('[data-cms-flow]');
      if (!btn) return;
      const flow = btn.dataset.cmsFlow;
      if (flow === 'daily') openOrCreateDailyNote();
      else if (flow === 'quick') openQuickNote();
      else if (flow === 'wikilink') openWikilinkPicker();
    });
    const tabs = sidebar.querySelector('#workspace-view-tabs');
    if (tabs) tabs.insertAdjacentElement('afterend', strip);
    else sidebar.prepend(strip);
  }
  injectQuickActions();

  // Wire mobile editing toolbar [[ button if present
  document.querySelector('#mobile-editing-toolbar')?.addEventListener('click', (event) => {
    const btn = event.target.closest('button[data-mobile-tool="wikilink"]');
    if (!btn) return;
    event.stopPropagation();
    openWikilinkPicker();
  });

  // Register flows in Command Palette
  if (window.MarkataShortcuts) {
    window.MarkataShortcuts.register({
      id: 'flows.dailyNote',
      label: "Flow: Open Today's Daily Note",
      category: 'Writing',
      description: "Open or create today's daily note and jump to the editor.",
      aliases: 'daily note journal today log',
      bindings: { vscode: 'Mod+Shift+D', writer: 'Mod+Shift+D', vim: 'Mod+Shift+D', emacs: 'Ctrl+Shift+D' },
      run: () => openOrCreateDailyNote(),
    });
    window.MarkataShortcuts.register({
      id: 'flows.quickNote',
      label: 'Flow: New Quick Note / Share Capture',
      category: 'Writing',
      description: 'Capture a quick note, thought, or bookmark URL.',
      aliases: 'quick note thought bookmark share capture',
      bindings: { vscode: 'Mod+Shift+N', writer: 'Mod+Shift+N', vim: 'Mod+Shift+N', emacs: 'Ctrl+Shift+N' },
      run: () => openQuickNote(),
    });
    window.MarkataShortcuts.register({
      id: 'flows.wikilink',
      label: 'Flow: Insert Wikilink ([[)',
      category: 'Writing',
      description: 'Search site posts and insert a [[slug|Title]] cross-link.',
      aliases: 'wikilink cross link backlink internal link',
      bindings: { vscode: 'Mod+Shift+K', writer: 'Mod+Shift+K', vim: 'Mod+Shift+K', emacs: 'Ctrl+Shift+K' },
      run: () => openWikilinkPicker(),
    });
  }

  // Handle PWA shortcuts (?flow=daily|quick|new) and OS Share Target (?share_title=...&share_text=...&share_url=...)
  function handleIncomingURL() {
    const params = new URLSearchParams(location.search);
    const flow = params.get('flow');
    const shareTitle = params.get('share_title') || '';
    const shareText = params.get('share_text') || '';
    const shareUrl = params.get('share_url') || '';
    const shareId = params.get('share_id');
    const shareError = params.get('share_error');
    if (!flow && !shareTitle && !shareText && !shareUrl && !shareId && !shareError) return;

    const cleanURL = location.pathname + location.hash;
    // Keep a photo share's ID in the URL until it is successfully saved:
    // cancel/reload should not silently discard the staged attachment.
    if (!shareId) history.replaceState(null, '', cleanURL);

    // Wait briefly for site tree to initialize
    setTimeout(async () => {
      if (shareError) {
        toast('Share failed: ' + shareError, 'warning');
      } else if (shareId) {
        try {
          const capture = await readStagedShare(shareId);
          if (!capture || Date.now() - capture.createdAt > 24 * 60 * 60 * 1000) {
            throw new Error('Shared photos expired or are no longer available.');
          }
          const text = capture.text || '';
          const sharedURL = capture.url || (/^https?:\/\/\S+$/i.test(text.trim()) ? text.trim() : '');
          openQuickNote({
            shareId, title: capture.title || '', text: sharedURL === text.trim() ? '' : text,
            url: sharedURL, files: capture.files || [],
          });
        } catch (error) {
          toast('Could not open shared photos: ' + error.message, 'warning');
        }
      } else if (shareTitle || shareText || shareUrl) {
        // Android often puts the shared URL inside share_text when share_url is empty
        let url = shareUrl;
        let text = shareText;
        if (!url && /^https?:\/\/\S+$/i.test(text.trim())) {
          url = text.trim();
          text = '';
        }
        openQuickNote({ title: shareTitle, text, url });
      } else if (flow === 'daily') {
        openOrCreateDailyNote();
      } else if (flow === 'quick') {
        openQuickNote();
      } else if (flow === 'new') {
        document.querySelector('#header-new-post')?.click();
      }
    }, 350);
  }
  handleIncomingURL();

  window.MarkataFlows = {
    openDailyNote: openOrCreateDailyNote,
    openQuickNote,
    openWikilinkPicker,
  };
})();
