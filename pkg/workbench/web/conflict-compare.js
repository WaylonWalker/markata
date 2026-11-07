/* Read-only comparison of a conflicting browser draft with the latest file. */
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
      '<div class="workspace-conflict-heading"><div><h2 id="workspace-conflict-title">Compare unsynced copies</h2>',
      '<p>No changes will be saved, restored or discarded here.</p></div>',
      '<button type="button" id="workspace-conflict-close" aria-label="Close comparison">Close</button></div>',
      '<p id="workspace-conflict-path"></p><p id="workspace-conflict-state" role="status"></p>',
      '<div class="workspace-conflict-columns">',
      '<section><h3>Browser draft</h3><pre id="workspace-conflict-browser"></pre></section>',
      '<section><h3>Current file on disk</h3><pre id="workspace-conflict-disk"></pre></section>',
      '</div><footer><button type="button" id="workspace-conflict-download-browser">Download browser copy</button>',
      '<button type="button" id="workspace-conflict-download-disk" disabled>Download disk copy</button>',
      '<button type="button" id="workspace-conflict-done">Keep editing</button></footer>'
    ].join('');
    document.body.append(dialog);

    let currentComparison = null;
    const pathIsMarkdown = path => typeof path === 'string' && /\.(?:md|markdown|mdown|mkd)$/i.test(path);
    const hasConflict = () => /(?:Site save failed|File changed on disk)/i.test(status.textContent || '');
    function update() {
      const path = workspace.activeFile();
      button.hidden = !(pathIsMarkdown(path) && hasConflict());
    }
    const observer = new MutationObserver(update);
    observer.observe(status, {childList:true, characterData:true, subtree:true});
    update();

    function download(source, path, suffix) {
      if (typeof source !== 'string' || typeof path !== 'string') return;
      const basename = path.split('/').pop().replace(/\.(?:md|markdown|mdown|mkd)$/i, '');
      const blob = new Blob([source], {type:'text/markdown;charset=utf-8'});
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = basename + '.' + suffix + '.md';
      document.body.append(link);
      link.click();
      link.remove();
      setTimeout(() => URL.revokeObjectURL(url), 1000);
    }

    button.addEventListener('click', async () => {
      const path = workspace.activeFile();
      if (!pathIsMarkdown(path) || !hasConflict()) return;
      const browserSource = editor.getValue();
      currentComparison = {path, browserSource, diskSource:null};
      dialog.querySelector('#workspace-conflict-path').textContent = path;
      dialog.querySelector('#workspace-conflict-browser').textContent = browserSource;
      dialog.querySelector('#workspace-conflict-disk').textContent = 'Loading current file…';
      dialog.querySelector('#workspace-conflict-state').textContent = 'Reading latest file from the workspace server…';
      dialog.querySelector('#workspace-conflict-download-disk').disabled = true;
      dialog.showModal();
      try {
        const response = await fetch('/workspace/file?path=' + encodeURIComponent(path), {cache:'no-store'});
        const data = await response.json();
        if (!response.ok || typeof data.source !== 'string')
          throw new Error(data.error || 'Unable to read the disk copy.');
        if (!currentComparison || currentComparison.path !== path || !dialog.open) return;
        currentComparison.diskSource = data.source;
        dialog.querySelector('#workspace-conflict-disk').textContent = data.source;
        dialog.querySelector('#workspace-conflict-state').textContent =
          'Read-only comparison. Download either copy before resolving the conflict manually.';
        dialog.querySelector('#workspace-conflict-download-disk').disabled = false;
      } catch (error) {
        if (!dialog.open) return;
        dialog.querySelector('#workspace-conflict-disk').textContent = '';
        dialog.querySelector('#workspace-conflict-state').textContent =
          'Could not read the current disk copy: ' + error.message;
      }
    });

    dialog.querySelector('#workspace-conflict-download-browser').addEventListener('click', () => {
      if (currentComparison) download(currentComparison.browserSource, currentComparison.path, 'browser');
    });
    dialog.querySelector('#workspace-conflict-download-disk').addEventListener('click', () => {
      if (currentComparison && typeof currentComparison.diskSource === 'string')
        download(currentComparison.diskSource, currentComparison.path, 'disk');
    });
    function close() { dialog.close(); }
    dialog.querySelector('#workspace-conflict-close').addEventListener('click', close);
    dialog.querySelector('#workspace-conflict-done').addEventListener('click', close);
    dialog.addEventListener('close', () => { currentComparison = null; update(); button.focus(); });
    window.addEventListener('pagehide', () => observer.disconnect(), {once:true});
  }
  window.MarkataConflictCompare = {attach};
})();
