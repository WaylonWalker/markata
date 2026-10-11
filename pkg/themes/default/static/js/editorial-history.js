(function () {
  'use strict';
  if (window.__markataEditorialHistory) return;
  window.__markataEditorialHistory = true;
  var viewerPromise = null;

  function initialize() {
    document.querySelectorAll('[data-editorial-open]').forEach(function (button) { button.hidden = false; });
    document.querySelectorAll('[data-editorial-events]').forEach(function (list) {
      list.querySelectorAll('.editorial-history__event').forEach(function (entry, index) { entry.hidden = index !== 0; });
    });
  }

  function loadViewer(src) {
    if (window.MarkataEditorialDiff) return Promise.resolve(window.MarkataEditorialDiff);
    if (viewerPromise) return viewerPromise;
    viewerPromise = new Promise(function (resolve, reject) {
      var script = document.createElement('script');
      var timer = setTimeout(fail, 15000);
      function fail() {
        clearTimeout(timer);
        script.remove();
        viewerPromise = null;
        reject(new Error('The comparison viewer could not load.'));
      }
      script.src = src;
      script.onload = function () {
        clearTimeout(timer);
        if (window.MarkataEditorialDiff) resolve(window.MarkataEditorialDiff);
        else fail();
      };
      script.onerror = fail;
      document.head.appendChild(script);
    });
    return viewerPromise;
  }

  function loadComparison(entry) {
    var panel = entry && entry.querySelector('[data-editorial-comparison]');
    if (!panel || panel.dataset.loaded === 'true' || panel.dataset.loading === 'true') return;
    panel.dataset.loading = 'true';
    panel.setAttribute('aria-busy', 'true');
    var status = panel.querySelector('[data-editorial-status]');
    status.textContent = 'Loading recorded changes…';
    loadViewer(panel.dataset.viewer)
      .then(function (viewer) {
        return viewer.show(panel.querySelector('[data-editorial-result]'), panel.dataset.before, panel.dataset.after);
      })
      .then(function () {
        panel.dataset.loaded = 'true';
        status.textContent = 'Added text is underlined; removed text is struck through. Unchanged text is collapsed around the edits.';
      })
      .catch(function () {
        status.textContent = 'Changes could not load. Open the diff in GitHub to view the revision.';
      })
      .finally(function () {
        delete panel.dataset.loading;
        panel.removeAttribute('aria-busy');
      });
  }

  document.addEventListener('click', function (event) {
    if (event.target.matches('dialog.editorial-history')) {
      event.target.close();
      return;
    }
    var opener = event.target.closest('[data-editorial-open]');
    if (opener) {
      var dialog = document.getElementById(opener.getAttribute('aria-controls'));
      if (dialog && !dialog.open) dialog.showModal();
      if (dialog) loadComparison(dialog.querySelector('.editorial-history__event:not([hidden])'));
      return;
    }
    var closer = event.target.closest('[data-editorial-close]');
    if (closer) {
      closer.closest('dialog').close();
      return;
    }
    var navigation = event.target.closest('[data-editorial-previous], [data-editorial-next]');
    if (navigation) {
      var historyDialog = navigation.closest('dialog');
      var entries = Array.from(historyDialog.querySelectorAll('[data-editorial-events] > .editorial-history__event'));
      var active = entries.findIndex(function (entry) { return !entry.hidden; });
      active += navigation.hasAttribute('data-editorial-next') ? 1 : -1;
      if (active < 0 || active >= entries.length) return;
      entries.forEach(function (entry, index) { entry.hidden = index !== active; });
      var previous = historyDialog.querySelector('[data-editorial-previous]');
      var next = historyDialog.querySelector('[data-editorial-next]');
      if (previous) previous.disabled = active === 0;
      if (next) next.disabled = active === entries.length - 1;
      historyDialog.querySelector('[data-editorial-position]').textContent = (active + 1) + ' of ' + entries.length;
      historyDialog.scrollTop = 0;
      loadComparison(entries[active]);
      return;
    }
  });
  window.addEventListener('view-transition-complete', initialize);
  if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', initialize, { once: true });
  else initialize();
}());
