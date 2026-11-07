/* Compact recent workspace files; local-only and safe with blocked storage. */
(() => {
  'use strict';

  function attach() {
    const workspace = window.MarkataSiteTree;
    const tree = document.querySelector('#site-files');
    const label = document.querySelector('#active-file');
    if (!workspace?.files || !workspace?.openFile || !tree || !label) return;

    const storageKey = 'markata-recent-files-v1:' + (workspace.workspaceId?.() || location.origin);
    const read = () => {
      try {
        const saved = JSON.parse(localStorage.getItem(storageKey) || '[]');
        return Array.isArray(saved) ? saved.filter(path => typeof path === 'string').slice(0, 8) : [];
      } catch { return []; }
    };
    let recent = read();
    const details = document.createElement('details');
    details.id = 'recent-files';
    details.className = 'recent-files';
    details.hidden = true;
    details.open = true;
    const summary = document.createElement('summary');
    summary.textContent = 'Recent files';
    const list = document.createElement('div');
    list.id = 'recent-files-list';
    list.setAttribute('aria-label', 'Recently opened workspace files');
    details.append(summary, list);
    tree.before(details);

    let current = '';
    function render() {
      const files = new Map(workspace.files().filter(file => typeof file.path === 'string')
        .map(file => [file.path, file]));
      recent = [...new Set(recent)].filter(path => files.has(path)).slice(0, 8);
      details.hidden = !recent.length;
      list.replaceChildren();
      for (const path of recent) {
        const file = files.get(path);
        const button = document.createElement('button');
        button.type = 'button';
        button.className = 'recent-files-item';
        if (path === current) button.setAttribute('aria-current', 'page');
        const name = document.createElement('strong');
        name.textContent = file.title || path.split('/').pop();
        const location = document.createElement('span');
        location.textContent = path;
        button.title = path;
        button.append(name, location);
        button.addEventListener('click', () => workspace.openFile(path));
        list.append(button);
      }
    }
    function observeActive() {
      const path = workspace.activeFile?.();
      if (!path || path === current) return;
      current = path;
      if (!workspace.files().some(file => file.path === path)) {
        render();
        return; // Never persist the scratch buffer as a site document.
      }
      recent = [path, ...recent.filter(previous => previous !== path)].slice(0, 8);
      try { localStorage.setItem(storageKey, JSON.stringify(recent)); }
      catch { /* Still usable in the current tab when storage is unavailable. */ }
      render();
    }
    const observer = new MutationObserver(observeActive);
    observer.observe(label, {childList:true, characterData:true, subtree:true});
    observeActive();
    window.addEventListener('pagehide', () => observer.disconnect(), {once:true});
  }
  window.MarkataRecentFiles = {attach};
})();
