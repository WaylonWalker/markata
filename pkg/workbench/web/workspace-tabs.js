/* Stable tab navigation for the workspace sidebar. */
(() => {
  'use strict';

  const select = document.querySelector('#site-view');
  const tablist = document.querySelector('#workspace-view-tabs');
  const views = ['recent', 'files', 'feeds', 'tags', 'changes'];
  if (!select || !tablist) return;

  const tabs = new Map(views.map(view => [view,
    tablist.querySelector(`[data-workspace-view="${view}"]`)]));
  const panels = new Map(views.map(view => [view,
    document.querySelector('#site-' + view)]));
  const preferenceKey = 'markata-workspace-view-v1';
  let pendingInitialView = null;

  function persist(view) {
    try {
      localStorage.setItem(preferenceKey, view);
      localStorage.setItem('markata-workspace-changes-open', String(view === 'changes'));
    } catch { /* Keep navigation usable when browser storage is unavailable. */ }
  }

  function sync(view = select.value) {
    if (tabs.get(view)?.disabled) {
      if (pendingInitialView === view) view = 'files';
      else {
        view = 'files';
        select.value = view;
        persist(view);
        select.dispatchEvent(new Event('change', {bubbles:true}));
      }
    }
    for (const [name, tab] of tabs) {
      const selected = name === view;
      tab.setAttribute('aria-selected', String(selected));
      tab.tabIndex = selected ? 0 : -1;
      const option = select.querySelector(`option[value="${name}"]`);
      tab.disabled = Boolean(option?.disabled);
      const panel = panels.get(name);
      if (panel) panel.hidden = !selected;
    }
    const changes = tabs.get('changes');
    const count = changes.querySelector('.workspace-changes-count')?.textContent;
    changes.setAttribute('aria-label', count ? `Changes, ${count} files` : 'Changes');
  }

  function activate(view, focus = false) {
    const tab = tabs.get(view);
    if (!tab || tab.disabled) return;
    pendingInitialView = null;
    select.value = view;
    persist(view);
    sync(view);
    select.dispatchEvent(new Event('change', {bubbles:true}));
    if (focus) tab.focus();
  }

  select.hidden = true;
  select.addEventListener('change', () => {
    if (pendingInitialView && select.value === 'files') {
      sync('files');
      return;
    }
    pendingInitialView = null;
    persist(select.value);
    sync(select.value);
  });
  new MutationObserver(() => {
    if (pendingInitialView) {
      const option = select.querySelector(`option[value="${pendingInitialView}"]`);
      const state = option?.textContent?.toLowerCase() || '';
      if (/\(\d+\)|unavailable/.test(state)) {
        const view = pendingInitialView;
        pendingInitialView = null;
        if (tabs.get(view)?.disabled) persist('files');
        else activate(view);
        return;
      }
    }
    sync();
  }).observe(select, {
    subtree:true, childList:true, characterData:true, attributes:true,
    attributeFilter:['disabled'],
  });
  for (const [view, tab] of tabs) {
    tab.addEventListener('click', () => activate(view));
    tab.addEventListener('keydown', event => {
      if (!['ArrowLeft', 'ArrowRight', 'Home', 'End'].includes(event.key)) return;
      event.preventDefault();
      const enabled = views.filter(name => !tabs.get(name).disabled);
      const index = enabled.indexOf(view);
      const next = event.key === 'Home' ? enabled[0]
        : event.key === 'End' ? enabled.at(-1)
          : enabled[(index + (event.key === 'ArrowRight' ? 1 : -1) + enabled.length) % enabled.length];
      if (next) activate(next, true);
    });
  }

  // Files is the stable first view; the hidden compatibility select's first
  // option is Recent, so do not use its implicit browser default here.
  let initial = 'files';
  try {
    const saved = localStorage.getItem(preferenceKey);
    if (views.includes(saved)) {
      if (saved === 'feeds' || saved === 'tags') {
        pendingInitialView = saved;
        initial = 'files';
      } else initial = saved;
    }
    else if (localStorage.getItem('markata-workspace-changes-open') === 'true') initial = 'changes';
  } catch { /* Use the Files tab when browser storage is blocked. */ }
  select.value = initial;
  sync(initial);
})();
