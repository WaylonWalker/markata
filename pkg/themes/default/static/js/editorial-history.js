(function () {
  'use strict';
  if (window.__markataEditorialHistory) return;
  window.__markataEditorialHistory = true;
  var viewerPromise = null;
  var nextId = 0;

  function initialize() {
    document.querySelectorAll('[data-editorial-open]').forEach(function (button) { button.hidden = false; });
    document.querySelectorAll('[data-editorial-events]').forEach(function (list) {
      list.querySelectorAll('.editorial-history__event').forEach(function (entry, index) { entry.hidden = index !== 0; });
    });
    document.querySelectorAll('[data-editorial-show]').forEach(function (button) {
      button.hidden = false;
      if (!button.hasAttribute('aria-controls')) {
        var panel = button.closest('.editorial-history__event').querySelector('.editorial-history__comparison');
        panel.id = 'editorial-comparison-' + (++nextId);
        button.setAttribute('aria-controls', panel.id);
      }
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

  document.addEventListener('click', async function (event) {
    var opener = event.target.closest('[data-editorial-open]');
    if (opener) {
      var dialog = document.getElementById(opener.getAttribute('aria-controls'));
      if (dialog && !dialog.open) dialog.showModal();
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
      return;
    }
    var button = event.target.closest('[data-editorial-show]');
    if (!button || button.disabled) return;
    var panel = document.getElementById(button.getAttribute('aria-controls'));
    if (!panel) return;
    if (button.dataset.loaded) {
      panel.hidden = !panel.hidden;
      button.setAttribute('aria-expanded', String(!panel.hidden));
      button.textContent = panel.hidden ? 'Show changes here' : 'Hide changes';
      return;
    }
    panel.hidden = false;
    button.setAttribute('aria-expanded', 'true');
    button.disabled = true;
    panel.setAttribute('aria-busy', 'true');
    var status = panel.querySelector('[data-editorial-status]');
    status.textContent = 'Loading the recorded revisions…';
    try {
      var viewer = await loadViewer(button.dataset.viewer);
      await viewer.show(panel.querySelector('[data-editorial-result]'), button.dataset.before, button.dataset.after);
      status.textContent = 'Recorded Markdown changes. Added text is underlined; removed text is struck through. Unchanged text is collapsed around the edits.';
      var firstChange = panel.querySelector('[data-editorial-result] ins, [data-editorial-result] del');
      if (firstChange) firstChange.scrollIntoView({ block: 'center', behavior: 'smooth' });
      button.dataset.loaded = 'true';
      button.textContent = 'Hide changes';
    } catch (_) {
      status.textContent = 'Changes could not load. The repository may be private or the revision unavailable. Open the diff in GitHub, or retry.';
      button.textContent = 'Retry changes';
    } finally {
      button.disabled = false;
      panel.removeAttribute('aria-busy');
    }
  });
  window.addEventListener('view-transition-complete', initialize);
  if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', initialize, { once: true });
  else initialize();
}());
