/* Native Markdown tabs: progressive enhancement only. Without JS every panel
 * stays readable; print always expands all panels. */
(() => {
  function initTabs(root) {
    if (root.dataset.mdTabsReady) return;
    const nav = root.querySelector(':scope > .md-tabs-nav');
    const panelsWrap = root.querySelector(':scope > .md-tabs-panels');
    if (!nav || !panelsWrap) return;
    const tabs = Array.from(nav.querySelectorAll(':scope > [role="tab"]'));
    const panels = Array.from(panelsWrap.querySelectorAll(':scope > .md-tab'));
    if (!tabs.length || tabs.length !== panels.length) return;
    if (!tabs.every((tab, i) => tab.getAttribute('aria-controls') === panels[i].id)) return;

    function activate(index, focus) {
      tabs.forEach((tab, i) => {
        const active = i === index;
        tab.setAttribute('aria-selected', String(active));
        tab.tabIndex = active ? 0 : -1;
        panels[i].hidden = !active;
      });
      if (focus) tabs[index].focus();
    }
    tabs.forEach((tab, index) => {
      tab.addEventListener('click', () => activate(index, false));
      tab.addEventListener('keydown', event => {
        let next = index;
        if (event.key === 'ArrowRight') next = (index + 1) % tabs.length;
        else if (event.key === 'ArrowLeft') next = (index + tabs.length - 1) % tabs.length;
        else if (event.key === 'Home') next = 0;
        else if (event.key === 'End') next = tabs.length - 1;
        else return;
        event.preventDefault();
        activate(next, true);
      });
    });
    activate(0, false);
    root.dataset.mdTabsReady = 'true';
    nav.hidden = false;
  }
  function initAll() {
    document.querySelectorAll('[data-md-tabs]').forEach(initTabs);
  }
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initAll, { once: true });
  } else {
    initAll();
  }
  window.addEventListener('view-transition-complete', initAll);
})();