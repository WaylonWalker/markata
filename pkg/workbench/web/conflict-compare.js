/* Compare a conflicting browser draft with the latest file and pick one copy. */
(() => {
  'use strict';

  function attach(editor) {
    const status = document.querySelector('#saved');
    const workspace = window.MarkataSiteTree;
    if (!editor || !status || !workspace?.activeFile) return;

    const button = document.createElement('button');
    button.id = 'workspace-conflict-review';
    button.type = 'button';
    button.textContent = 'Compare copies…';
    button.hidden = true;
    status.insertAdjacentElement('afterend', button);

    const dialog = document.createElement('dialog');
    dialog.id = 'workspace-conflict-dialog';
    dialog.setAttribute('aria-labelledby', 'workspace-conflict-title');
    dialog.innerHTML = [
      '<div class="workspace-conflict-heading"><div><h2 id="workspace-conflict-title">Resolve conflicting copies</h2>',
      '<p>Pick one copy. The other copy is replaced.</p></div>',
      '<button type="button" id="workspace-conflict-close" aria-label="Close comparison">Close</button></div>',
      '<p id="workspace-conflict-path"></p><p id="workspace-conflict-state" role="status"></p>',
      '<div class="workspace-conflict-columns">',
      '<section><h3>Browser draft</h3><pre id="workspace-conflict-browser"></pre></section>',
      '<section><h3>Current file on disk</h3><pre id="workspace-conflict-disk"></pre></section>',
      '</div><footer><button type="button" id="workspace-conflict-keep-browser" disabled>Keep browser draft</button>',
      '<button type="button" id="workspace-conflict-use-disk" disabled>Use disk copy</button>',
      '<button type="button" id="workspace-conflict-done">Keep editing</button></footer>'
    ].join('');
    document.body.append(dialog);

    const keepBrowserButton = dialog.querySelector('#workspace-conflict-keep-browser');
    const useDiskButton = dialog.querySelector('#workspace-conflict-use-disk');
    let currentComparison = null;
    let requestNumber = 0;
    let conflictPath = null;
    const pathIsMarkdown = path => typeof path === 'string' && /\.(?:md|markdown|mdown|mkd)$/i.test(path);
    const hasConflict = () => /(?:Site save failed|File changed on disk)/i.test(status.textContent || '');
    function update() {
      const path = workspace.activeFile();
      button.hidden = !(pathIsMarkdown(path) && path === conflictPath && hasConflict());
    }
    conflictPath = hasConflict() ? workspace.activeFile() : null;
    const observer = new MutationObserver(() => {
      conflictPath = hasConflict() ? workspace.activeFile() : null;
      update();
    });
    observer.observe(status, {childList:true, characterData:true, subtree:true});
    const activeObserver = new MutationObserver(update);
    const activeLabel = document.querySelector('#active-file');
    if (activeLabel) activeObserver.observe(activeLabel, {childList:true, characterData:true, subtree:true});
    update();

    function setChoicesEnabled(enabled) {
      keepBrowserButton.disabled = !enabled;
      useDiskButton.disabled = !enabled;
    }

    button.addEventListener('click', async () => {
      const path = workspace.activeFile();
      if (!pathIsMarkdown(path) || !hasConflict()) return;
      const browserSource = editor.getValue();
      const request = ++requestNumber;
      currentComparison = {path, browserSource, diskSource:null, diskVersion:null};
      dialog.querySelector('#workspace-conflict-path').textContent = path;
      dialog.querySelector('#workspace-conflict-browser').textContent = browserSource;
      dialog.querySelector('#workspace-conflict-disk').textContent = 'Loading current file…';
      dialog.querySelector('#workspace-conflict-state').textContent = 'Reading latest file from the workspace server…';
      setChoicesEnabled(false);
      dialog.showModal();
      try {
        const response = await fetch('/workspace/file?path=' + encodeURIComponent(path), {cache:'no-store'});
        const data = await response.json();
        if (!response.ok || typeof data.source !== 'string' || typeof data.version !== 'string')
          throw new Error(data.error || 'Unable to read the disk copy.');
        if (request !== requestNumber || !currentComparison || currentComparison.path !== path || !dialog.open) return;
        currentComparison.diskSource = data.source;
        currentComparison.diskVersion = data.version;
        dialog.querySelector('#workspace-conflict-disk').textContent = data.source;
        dialog.querySelector('#workspace-conflict-state').textContent =
          'Pick the copy to keep. The other copy is replaced.';
        setChoicesEnabled(true);
      } catch (error) {
        if (request !== requestNumber || !dialog.open) return;
        dialog.querySelector('#workspace-conflict-disk').textContent = '';
        dialog.querySelector('#workspace-conflict-state').textContent =
          'Could not read the current disk copy: ' + error.message;
      }
    });

    keepBrowserButton.addEventListener('click', () => {
      const comparison = currentComparison;
      if (!comparison || typeof comparison.diskVersion !== 'string') return;
      if (!window.confirm('Replace the file on disk with your browser draft?')) return;
      if (workspace.keepBrowserDraft(comparison.path, comparison.diskVersion)) dialog.close();
    });
    useDiskButton.addEventListener('click', () => {
      const comparison = currentComparison;
      if (!comparison || typeof comparison.diskSource !== 'string') return;
      if (workspace.useDiskCopy(comparison.path, comparison.diskSource, comparison.diskVersion)) dialog.close();
    });
    function close() { dialog.close(); }
    dialog.querySelector('#workspace-conflict-close').addEventListener('click', close);
    dialog.querySelector('#workspace-conflict-done').addEventListener('click', close);
    dialog.addEventListener('close', () => { requestNumber++; currentComparison = null; update(); button.focus(); });
    window.addEventListener('pagehide', () => {
      observer.disconnect();
      activeObserver.disconnect();
    }, {once:true});
  }
  window.MarkataConflictCompare = {attach};
})();
