/* Opt-in interactive Markdown checklists.
 * Standard Goldmark task checkboxes remain disabled/read-only unless wrapped
 * in ::: md-checklist. No JavaScript or storage is needed to read a list.
 */
(() => {
  'use strict';

  if (typeof window.markataInitChecklists === 'function') {
    window.markataInitChecklists();
    return;
  }

  const prefix = 'markata:checklist:v1:';

  function initChecklists() {
    document.querySelectorAll('.post-content .md-checklist').forEach((container) => {
      const boxes = Array.from(container.querySelectorAll('li.task-list-item > input[type="checkbox"]'))
        .filter((box) => box.disabled);
      if (!boxes.length) return;

      // Explicit IDs are required for persistence. Lists without IDs stay
      // interactive for the current page only, with no storage side effects.
      const id = (container.getAttribute('data-id') || '').trim();
      const key = id ? prefix + location.pathname + ':' + id : null;
      let stored = null;
      if (key) {
        try {
          const parsed = JSON.parse(localStorage.getItem(key));
          if (Array.isArray(parsed) && parsed.length === boxes.length &&
              parsed.every((value) => typeof value === 'boolean')) {
            stored = parsed;
          }
        } catch (_) {
          // Storage can be blocked in private browsing or restricted embeds.
        }
      }

      boxes.forEach((box, index) => {
        if (stored) box.checked = stored[index];
        box.disabled = false;
        const item = box.closest('li.task-list-item');
        if (item && !box.hasAttribute('aria-label')) {
          const label = (item.textContent || '').trim().replace(/\s+/g, ' ').slice(0, 160);
          box.setAttribute('aria-label', label || 'Checklist item ' + (index + 1));
        }
        box.addEventListener('change', () => {
          if (!key) return;
          try {
            localStorage.setItem(key, JSON.stringify(boxes.map((current) => current.checked)));
          } catch (_) {
            // Continue to work transiently if persistence is unavailable.
          }
        });
      });
    });
  }

  window.markataInitChecklists = initChecklists;
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initChecklists, {once: true});
  } else {
    initChecklists();
  }
  window.addEventListener('view-transition-complete', initChecklists);
})();
