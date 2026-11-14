// Applies the workbench button system (buttons.css) to every button, including
// buttons that scripts create later. Segmented tabs, tag chips, and list-style
// choices keep their own component styling and are skipped.
(() => {
  const SKIP = [
    '.publish-mode-tabs', '.publish-side-tabs', '.publish-pierre-actions',
    '.fm-field', '.fm-tag-chip', '.fm-nested-item', '.fm-add-nested',
    '.template-choice', '.settings-body nav', '.recent-files-item', '.CodeMirror',
    '[role="tab"]', '[role="option"]',
  ].join(',');
  const PRIMARY = '#publish-open,#publish-confirm,#fm-add,#fm-add-custom,#create-form [type="submit"]';
  const DANGER = '#publish-discard-accept';
  const SUBTLE = '.editor-toolbar > button,#frontmatter-raw,#fm-add-toggle,#fm-more,#preview-focus-editor,#site-sidebar .files-open';
  const GHOST = '#commands,#publish-review-open,#workspace-changes-toggle,#settings-open,#theme-picker-toggle,#palette-mode-toggle,#diagnostics-toggle,#problems-toggle,.theme-switcher button';
  const SMALL = '#preview-focus-editor,#diagnostics-toggle,#problems-toggle,#workspace-conflict-review,#fm-add-toggle';
  const KBD = {
    'publish-open': 'Ctrl+S',
    'publish-review-open': 'Ctrl+Shift+Enter',
    commands: 'Ctrl+Shift+P',
  };

  const isClose = (button) => /(^|-)close$/.test(button.id)
    || /(^|\s)[\w-]*-close(\s|$)/.test(button.className.toString())
    || /^close\b/i.test(button.getAttribute('aria-label') || '');
  const isIcon = (button) => isClose(button) || button.id === 'settings-open';

  function style(button) {
    if (button.matches(SKIP) || button.closest(SKIP)) return;
    button.classList.add('btn');
    if (button.matches(PRIMARY)) button.classList.add('btn-primary');
    else if (button.matches(DANGER)) button.classList.add('btn-danger');
    else if (isIcon(button)) button.classList.add('btn-icon', 'btn-ghost');
    else if (button.matches(SUBTLE)) button.classList.add('btn-subtle');
    else if (button.matches(GHOST)) button.classList.add('btn-ghost');
    if (button.matches(SMALL)) button.classList.add('btn-sm');
    if (KBD[button.id] && button.dataset.kbd !== KBD[button.id]) button.dataset.kbd = KBD[button.id];
  }

  function styleAll(root) {
    if (root.nodeType !== Node.ELEMENT_NODE) return;
    if (root.tagName === 'BUTTON') style(root);
    root.querySelectorAll?.('button').forEach(style);
  }

  styleAll(document.body);
  new MutationObserver((records) => {
    for (const record of records) {
      record.addedNodes.forEach(styleAll);
      if (record.type === 'attributes') styleAll(record.target);
    }
  }).observe(document.body, {
    childList: true,
    subtree: true,
    attributes: true,
    attributeFilter: ['id', 'aria-label', 'type'],
  });
})();
