/* Feed Studio: feed settings, authoritative filter matching and built-site preview. */
(() => {
  'use strict';

  const workbench = document.getElementById('workbench');
  const studio = document.getElementById('feed-studio');
  const select = document.getElementById('fs-feed');
  const form = document.getElementById('fs-form');
  const status = document.getElementById('fs-status');
  const count = document.getElementById('fs-count');
  const invalid = document.getElementById('fs-filter-error');
  const siteButton = document.getElementById('fs-preview-site');
  const matchesButton = document.getElementById('fs-preview-members');
  const iframe = document.getElementById('fs-rendered');
  const cards = document.getElementById('fs-cards');
  const siteEmpty = document.getElementById('fs-site-empty');
  const returnButton = document.getElementById('feed-studio-return');
  const fields = ['title', 'description', 'filter', 'sort', 'reverse', 'presentation', 'template'];
  const field = name => document.getElementById('fs-' + name);
  if (!studio || !select || !workbench) return;

  let feeds = [];
  let active = null;
  let original = null;
  let version = '';
  let timer;
  let requestNumber = 0;
  let aborter;
  let validFilter = true;
  let previewMode = 'site';
  let snapshotURL = '';
  let snapshotAvailable = false;
  let opened = false;

  const setStatus = (message, error = false) => {
    status.textContent = message;
    status.dataset.error = String(error);
  };
  const data = () => Object.fromEntries(fields.map(name => [
    name, name === 'reverse' ? field(name).checked : field(name).value,
  ]));
  const equal = (a, b) => JSON.stringify(a) === JSON.stringify(b);
  const dirty = () => original && !equal(data(), original);
  function updateSave() {
    const selected = feeds.find(feed => feed.slug === active);
    const writable = selected?.editable === true;
    document.getElementById('fs-save').disabled = !writable || !validFilter || !dirty();
    document.getElementById('fs-readonly').hidden = writable;

  }
  function setFields(feed) {
    version = feed.version || version;
    document.getElementById('fs-source').textContent = feed.source ? 'Editing ' + feed.source : 'Read-only feed source';
    for (const name of fields) {
      if (name === 'reverse') field(name).checked = Boolean(feed.reverse);
      else field(name).value = feed[name] ?? '';
      // Preserve arbitrary authored sort and presentation values.
      if (field(name).tagName === 'SELECT' && field(name).value !== String(feed[name] ?? '')) {
        const opt = document.createElement('option');
        opt.value = String(feed[name] ?? '');
        opt.textContent = opt.value || 'Default';
        field(name).append(opt);
        field(name).value = opt.value;
      }
    }
    original = data();
    validFilter = true;
    invalid.hidden = true;
    updateSave();
  }
  function feedLabel(feed) {
    const label = feed.title || feed.slug || 'Home';
    return label + ' · ' + Number(feed.count || 0).toLocaleString();
  }
  function renderSelector(nextFeeds, selected) {
    feeds = nextFeeds;
    select.replaceChildren();
    for (const feed of feeds) {
      const option = document.createElement('option');
      option.value = feed.slug;
      option.textContent = feedLabel(feed);
      select.append(option);
    }
    if (selected && feeds.some(feed => feed.slug === selected)) select.value = selected;
    else if (feeds.length) select.value = feeds[0].slug;
    active = select.value;
    document.getElementById('fs-empty').hidden = feeds.length > 0;
    form.hidden = feeds.length === 0;
  }
  async function fetchJSON(url, options = {}) {
    const response = await fetch(url, {cache: 'no-store', ...options});
    const value = await response.json();
    if (!response.ok) throw new Error(value.error || 'Could not load feed data');
    return value;
  }
  async function load(selected = active) {
    setStatus('Loading feed definitions…');
    try {
      const response = await fetchJSON('/workspace/feed-studio');
      version = response.version || '';
      renderSelector(response.feeds || [], selected);
      if (feeds.some(item => item.slug === active)) {
        const feed = feeds.find(item => item.slug === active);
        setFields(feed);
        try { localStorage.setItem('markata-feed-studio-selected', active); } catch {}
        await refreshPreview();
        if (validFilter) setStatus('Feed ready · edit fields to preview matching posts');
      } else setStatus('No configured feeds yet');
    } catch (error) {
      setStatus(error.message, true);
    }
  }
  function showMode(mode) {
    previewMode = mode;
    const rendered = mode === 'site';
    iframe.hidden = !rendered || !snapshotAvailable;
    siteEmpty.hidden = !rendered || snapshotAvailable;
    cards.hidden = rendered;
    siteButton.setAttribute('aria-pressed', String(rendered));
    matchesButton.setAttribute('aria-pressed', String(!rendered));
  }
  function setSnapshot(payload) {
    snapshotAvailable = Boolean(payload.rendered_available);
    snapshotURL = payload.rendered_url || '';
    siteButton.disabled = !snapshotAvailable;
    siteEmpty.textContent = snapshotAvailable
      ? '' : 'Build the site to enable the exact feed-template preview. Matching posts are available now.';
    if (snapshotAvailable && iframe.dataset.source !== snapshotURL) {
      iframe.dataset.source = snapshotURL;
      iframe.src = snapshotURL;
    }
    if (!snapshotAvailable || dirty()) showMode('members');
    else showMode(previewMode);
  }
  function text(tag, value, className) {
    const element = document.createElement(tag);
    element.textContent = value;
    if (className) element.className = className;
    return element;
  }
  async function openPost(path) {
    close();
    document.getElementById('workspace-tab-files')?.click();
    await window.MarkataSiteTree?.openFile(path);
    returnButton.hidden = false;
  }
  function renderPosts(payload) {
    cards.replaceChildren();
    const total = Number(payload.matched_count || 0);
    const shown = Number(payload.visible_count || 0);
    count.textContent = total.toLocaleString() + ' matching posts' +
      (shown !== total ? ' · ' + shown.toLocaleString() + ' after feed limits' : '') +
      (payload.partial ? ' · first 60 shown' : '');
    const posts = payload.posts || [];
    if (!posts.length) {
      cards.append(text('p', 'No matching published posts. Try a broader filter.', 'fs-no-results'));
      return;
    }
    for (const post of posts) {
      const article = document.createElement('article');
      article.className = 'fs-post-card';
      const main = document.createElement('div');
      main.append(text('h3', post.title || post.slug));
      main.append(text('p', post.date || post.path, 'fs-post-meta'));
      const button = text('button', 'Edit post', 'fs-edit-post');
      button.type = 'button';
      button.addEventListener('click', () => openPost(post.path));
      article.append(main, button);
      cards.append(article);
    }
  }
  async function refreshPreview() {
    if (!feeds.some(item => item.slug === active)) return;
    const sequence = ++requestNumber;
    aborter?.abort();
    aborter = new AbortController();
    const values = data();
    const params = new URLSearchParams({
      slug: active, filter: values.filter, sort: values.sort,
      reverse: String(values.reverse),
    });
    count.textContent = 'Checking feed membership…';
    try {
      const response = await fetchJSON('/workspace/feed-studio?' + params.toString(),
        {signal: aborter.signal});
      if (sequence !== requestNumber) return;
      validFilter = !response.filter_error;
      invalid.hidden = validFilter;
      invalid.textContent = response.filter_error || '';
      renderPosts(response);
      setSnapshot(response);
      updateSave();
      if (!validFilter) setStatus('Fix the filter expression before saving.', true);
    } catch (error) {
      if (error.name !== 'AbortError' && sequence === requestNumber) {
        setStatus(error.message, true);
        count.textContent = 'Preview unavailable';
      }
    }
  }
  function schedulePreview() {
    clearTimeout(timer);
    setStatus('Unsaved feed settings · previewing changes');
    updateSave();
    timer = setTimeout(refreshPreview, 260);
  }
  async function save(event) {
    event.preventDefault();
    if (!feeds.some(item => item.slug === active) || !validFilter || !dirty()) return;
    const selected = feeds.find(feed => feed.slug === active);
    if (!selected?.editable) return;
    const button = document.getElementById('fs-save');
    button.disabled = true;
    setStatus('Saving site configuration…');
    try {
      const response = await fetchJSON('/workspace/feed-studio', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({slug: active, expected_version: version, ...data()}),
      });
      version = response.version;
      original = data();
      setStatus('Saved to ' + response.file + ' · ready for review/publish');
      await refreshPreview();
    } catch (error) {
      setStatus(error.message + ' · your changes are still here', true);
    } finally {
      updateSave();
    }
  }
  function open() {
    studio.hidden = false;
    workbench.dataset.feedStudio = 'open';
    returnButton.hidden = true;
    if (opened) return;
    opened = true;
    let stored = '';
    try { stored = localStorage.getItem('markata-feed-studio-selected') || ''; } catch {}
    load(stored);
  }
  function close() {
    studio.hidden = true;
    delete workbench.dataset.feedStudio;
  }
  for (const name of fields) {
    field(name).addEventListener(name === 'reverse' || field(name).tagName === 'SELECT' ? 'change' : 'input', schedulePreview);
  }
  select.addEventListener('change', () => {
    if (dirty() && !window.confirm('Discard unsaved feed settings and switch feeds?')) {
      select.value = active;
      return;
    }
    active = select.value;
    const feed = feeds.find(item => item.slug === active);
    if (feed) {
      setFields(feed);
      try { localStorage.setItem('markata-feed-studio-selected', active); } catch {}
      previewMode = 'site';
      refreshPreview();
    }
  });
  siteButton.addEventListener('click', () => showMode('site'));
  matchesButton.addEventListener('click', () => showMode('members'));
  form.addEventListener('submit', save);
  document.getElementById('fs-refresh').addEventListener('click', async () => {
    if (dirty() && !window.confirm('Reload feed definitions and discard unsaved changes?')) return;
    await load(active);
  });
  document.getElementById('fs-close').addEventListener('click', close);
  returnButton.addEventListener('click', () => {
    document.getElementById('workspace-tab-feeds')?.click();
    open();
  });
  document.getElementById('workspace-tab-feeds')?.addEventListener('click', open);
  document.getElementById('site-view')?.addEventListener('change', () => {
    if (document.getElementById('site-view')?.value !== 'feeds') close();
  });
  window.MarkataFeedStudio = {open, close, selectedFeed: () => active};
})();