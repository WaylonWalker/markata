/* Modern CMS shell interactions: desktop mode switcher, live document header,
   sidebar action grouping, and full-height mobile CMS navigation. */
(() => {
  'use strict';

  const workbench = document.querySelector('#workbench');
  const layoutSelect = document.querySelector('#setting-layout');
  const modeButtons = [...document.querySelectorAll('.header-modes [data-layout-mode]')];
  const sidebarToggle = document.querySelector('#sidebar-toggle');
  const headerNewPost = document.querySelector('#header-new-post');
  const headerCommands = document.querySelector('#header-commands');
  const mobileNewPost = document.querySelector('#mobile-new-post');
  const mobileEditingToolbar = document.querySelector('#mobile-editing-toolbar');
  const mobilePullCue = document.querySelector('#mobile-pull-cue');
  const headerDoc = document.querySelector('#header-doc');
  const headerDocTitle = document.querySelector('#header-doc-title');
  const headerDocMeta = document.querySelector('#header-doc-meta');
  const mobileNav = document.querySelector('#mobile-nav');
  const mobileTabs = [...document.querySelectorAll('#mobile-nav [data-mobile-tab]')];
  const mobileQuery = window.matchMedia('(max-width: 768px)');

  if (!workbench) return;

  const MOBILE_VIEWS = ['content', 'write', 'preview'];
  let userHasInteractedWithFiles = false;

  function isMobile() {
    return mobileQuery.matches;
  }

  function refreshEditor() {
    requestAnimationFrame(() => {
      document.querySelector('.CodeMirror')?.CodeMirror?.refresh();
    });
  }

  function setMobileView(view, { focusEditor = false } = {}) {
    const next = MOBILE_VIEWS.includes(view) ? view : 'content';
    workbench.dataset.mobileView = next;
    for (const button of mobileTabs) {
      const active = button.dataset.mobileTab === next;
      button.setAttribute('aria-current', active ? 'page' : 'false');
    }
    if (sidebarToggle && isMobile()) {
      sidebarToggle.setAttribute('aria-expanded', String(next === 'content'));
    }
    if (next === 'write') {
      refreshEditor();
      if (focusEditor) {
        requestAnimationFrame(() => document.querySelector('.CodeMirror')?.CodeMirror?.focus());
      }
    }
  }

  // Default to the CMS Content view on mobile so browsing/changes work immediately.
  setMobileView('content');

  for (const button of mobileTabs) {
    button.addEventListener('click', () => {
      setMobileView(button.dataset.mobileTab);
    });
  }

  // Sync desktop segmented layout buttons with #setting-layout and workbench[data-layout].
  function syncLayoutModes() {
    const current = workbench.dataset.layout || layoutSelect?.value || 'split';
    for (const button of modeButtons) {
      button.setAttribute('aria-pressed', String(button.dataset.layoutMode === current));
    }
  }

  for (const button of modeButtons) {
    button.addEventListener('click', () => {
      const mode = button.dataset.layoutMode;
      if (!mode || !layoutSelect) return;
      layoutSelect.value = mode;
      layoutSelect.dispatchEvent(new Event('change', { bubbles: true }));
      syncLayoutModes();
    });
  }

  new MutationObserver(() => {
    syncLayoutModes();
    if (sidebarToggle && !isMobile()) {
      sidebarToggle.setAttribute('aria-expanded', String(!workbench.classList.contains('sidebar-hidden')));
    }
  }).observe(workbench, { attributes: true, attributeFilter: ['data-layout', 'class'] });

  layoutSelect?.addEventListener('change', syncLayoutModes);
  syncLayoutModes();

  // Sidebar toggle: collapses rail on desktop, switches to/from Content on mobile.
  if (sidebarToggle) {
    sidebarToggle.setAttribute('aria-expanded', String(!workbench.classList.contains('sidebar-hidden')));
    sidebarToggle.addEventListener('click', () => {
      if (isMobile()) {
        setMobileView(workbench.dataset.mobileView === 'content' ? 'write' : 'content');
        return;
      }
      if (window.MarkataShortcuts?.execute) {
        window.MarkataShortcuts.execute('workbench.toggleSidebar');
      } else {
        workbench.classList.toggle('sidebar-hidden');
        refreshEditor();
      }
      sidebarToggle.setAttribute('aria-expanded', String(!workbench.classList.contains('sidebar-hidden')));
    });
  }

  // + New post & Command Palette triggers
  function openNewPost() {
    window.MarkataShortcuts?.execute?.('workbench.newPost');
  }
  function openCommandPalette() {
    const cmdBtn = document.querySelector('#commands');
    if (cmdBtn) {
      cmdBtn.click();
      return;
    }
    window.MarkataShortcuts?.execute?.('workbench.commandPalette');
  }
  headerNewPost?.addEventListener('click', openNewPost);
  mobileNewPost?.addEventListener('click', openNewPost);
  headerCommands?.addEventListener('click', openCommandPalette);

  // Clicking the header document chip opens Quick Open on desktop.
  // On mobile, tapping it opens Page Properties (like md.waylonwalker.com's document bar)
  // when in Write/Preview, or switches to Write when in Content.
  headerDoc?.addEventListener('click', () => {
    if (isMobile()) {
      const current = workbench.dataset.mobileView || 'content';
      if (current === 'content') {
        setMobileView('write');
        return;
      }
      if (current === 'preview') {
        setMobileView('write');
      }
      document.querySelector('#frontmatter-toggle')?.click();
      return;
    }
    const quickOpen = document.querySelector('#site-sidebar .files-open');
    if (quickOpen && !quickOpen.disabled) quickOpen.click();
    else window.MarkataShortcuts?.execute?.('workbench.openFile');
  });

  // Keep the header document title and save state synchronized with active file,
  // frontmatter preview title, and #saved status.
  const activeFileEl = document.querySelector('#active-file');
  const previewTitleEl = document.querySelector('#preview-title');
  const savedEl = document.querySelector('#saved');

  function updateHeaderDoc() {
    if (!headerDocTitle || !headerDocMeta) return;
    const rawPath = (activeFileEl?.textContent || 'draft.md').trim();
    const previewTitle = (previewTitleEl?.textContent || '').trim();
    const fileName = rawPath.split('/').pop() || rawPath || 'draft.md';
    headerDocTitle.textContent = previewTitle && previewTitle !== 'Untitled document' ? previewTitle : fileName;
    headerDocTitle.title = rawPath;

    const savedText = (savedEl?.textContent || '').trim();
    let dotState = 'saved';
    let shortStatus = 'Saved';
    if (!savedText) {
      dotState = 'local';
      shortStatus = rawPath;
    } else if (/failed|conflict|changed on disk|not saved/i.test(savedText)) {
      dotState = 'error';
      shortStatus = savedText;
    } else if (/saving|unsaved|dirty/i.test(savedText)) {
      dotState = 'dirty';
      shortStatus = savedText;
    } else if (/this device|scratch/i.test(savedText)) {
      dotState = 'local';
      shortStatus = 'Local draft';
    } else {
      dotState = 'saved';
      shortStatus = 'Saved';
    }

    const parts = [];
    if (rawPath && rawPath !== headerDocTitle.textContent) parts.push(rawPath);
    if (shortStatus && shortStatus !== rawPath) parts.push(shortStatus);

    headerDocMeta.replaceChildren();
    const dot = document.createElement('span');
    dot.className = 'cms-save-dot';
    dot.dataset.state = dotState;
    dot.setAttribute('aria-hidden', 'true');
    const text = document.createElement('span');
    text.textContent = parts.join(' · ') || 'Draft';
    headerDocMeta.append(dot, text);
  }

  for (const target of [activeFileEl, previewTitleEl, savedEl]) {
    if (!target) continue;
    new MutationObserver(updateHeaderDoc).observe(target, {
      childList: true,
      characterData: true,
      subtree: true,
    });
  }
  updateHeaderDoc();

  // Group Open file + Commands into one sleek 2-column action row in the sidebar
  // once site-tree.js and command-palette.js have mounted them.
  const sidebar = document.querySelector('#site-sidebar');
  function arrangeSidebarActions() {
    if (!sidebar) return;
    const filesOpen = sidebar.querySelector('.files-open');
    const commands = sidebar.querySelector('#commands');
    if (!filesOpen && !commands) return;
    let bar = sidebar.querySelector('.cms-sidebar-actions');
    if (!bar) {
      bar = document.createElement('div');
      bar.className = 'cms-sidebar-actions';
      const anchor = sidebar.querySelector('.pane-label');
      if (anchor) anchor.after(bar);
      else sidebar.prepend(bar);
    }
    if (filesOpen && filesOpen.parentElement !== bar) bar.prepend(filesOpen);
    if (commands && commands.parentElement !== bar) bar.append(commands);
    if (filesOpen && commands && bar.firstElementChild !== filesOpen) bar.prepend(filesOpen);
  }

  if (sidebar) {
    arrangeSidebarActions();
    new MutationObserver(arrangeSidebarActions).observe(sidebar, { childList: true });
  }

  // Track user interactions inside file browsers so tapping a post on mobile
  // smoothly navigates into the Write view without hijacking initial page load.
  const fileBrowserSelectors = '#site-files, #site-recent, #site-feeds, #site-tags, #files-results, #file-finder-results';
  document.addEventListener('click', (event) => {
    if (event.target?.closest?.(fileBrowserSelectors)) {
      userHasInteractedWithFiles = true;
    }
  }, true);

  let lastObservedFile = (activeFileEl?.textContent || '').trim();
  if (activeFileEl) {
    new MutationObserver(() => {
      const next = (activeFileEl.textContent || '').trim();
      if (next && next !== lastObservedFile) {
        lastObservedFile = next;
        if (isMobile() && userHasInteractedWithFiles) {
          userHasInteractedWithFiles = false;
          setMobileView('write');
        }
      }
    }).observe(activeFileEl, { childList: true, characterData: true, subtree: true });
  }

  // Returning to editor from preview or creating a new post switches to Write on mobile.
  document.querySelector('#preview-focus-editor')?.addEventListener('click', () => {
    if (isMobile()) setMobileView('write', { focusEditor: true });
  });

  document.querySelector('#create-form')?.addEventListener('submit', () => {
    if (isMobile()) userHasInteractedWithFiles = true;
  });

  // Mobile Markdown formatting toolbar (md.waylonwalker.com pattern)
  function getCM() {
    return document.querySelector('.CodeMirror')?.CodeMirror || null;
  }

  function applyMobileTool(tool) {
    if (tool === 'commands') {
      openCommandPalette();
      return;
    }
    if (tool === 'done') {
      document.activeElement?.blur?.();
      document.body.classList.remove('mobile-writing-active');
      refreshEditor();
      return;
    }
    const cm = getCM();
    if (!cm) return;
    cm.focus();
    const sel = cm.getSelection();
    const cursor = cm.getCursor();
    switch (tool) {
      case 'undo':
        cm.undo();
        break;
      case 'redo':
        cm.redo();
        break;
      case 'h2': {
        const line = cm.getLine(cursor.line) || '';
        if (/^##\s+/.test(line)) {
          cm.replaceRange(line.replace(/^##\s+/, ''), { line: cursor.line, ch: 0 }, { line: cursor.line, ch: line.length });
        } else {
          cm.replaceRange('## ' + line.replace(/^#+\s*/, ''), { line: cursor.line, ch: 0 }, { line: cursor.line, ch: line.length });
        }
        break;
      }
      case 'list': {
        const line = cm.getLine(cursor.line) || '';
        if (/^-\s+/.test(line)) {
          cm.replaceRange(line.replace(/^-\s+/, ''), { line: cursor.line, ch: 0 }, { line: cursor.line, ch: line.length });
        } else {
          cm.replaceRange('- ' + line, { line: cursor.line, ch: 0 }, { line: cursor.line, ch: line.length });
        }
        break;
      }
      case 'code':
        if (sel) {
          cm.replaceSelection(sel.includes('\n') ? '```\n' + sel + '\n```' : '`' + sel + '`');
        } else {
          cm.replaceSelection('``');
          cm.setCursor({ line: cursor.line, ch: cursor.ch + 1 });
        }
        break;
      case 'link':
        if (sel) {
          cm.replaceSelection('[' + sel + '](https://)');
        } else {
          cm.replaceSelection('[](https://)');
          cm.setCursor({ line: cursor.line, ch: cursor.ch + 1 });
        }
        break;
      default:
        break;
    }
  }

  if (mobileEditingToolbar) {
    // Keep virtual keyboard open when tapping formatting buttons
    mobileEditingToolbar.addEventListener('mousedown', (event) => {
      if (event.target.closest('button[data-mobile-tool]:not([data-mobile-tool="done"]):not([data-mobile-tool="commands"])')) {
        event.preventDefault();
      }
    });
    mobileEditingToolbar.addEventListener('click', (event) => {
      const btn = event.target.closest('button[data-mobile-tool]');
      if (!btn) return;
      applyMobileTool(btn.dataset.mobileTool);
    });
  }

  // Virtual-keyboard-aware viewport height + writing mode
  function syncVisualViewport() {
    if (!isMobile()) {
      document.documentElement.style.removeProperty('--cms-viewport-height');
      document.documentElement.style.removeProperty('--cms-viewport-top');
      document.body.classList.remove('mobile-writing-active');
      return;
    }
    const vv = window.visualViewport;
    if (vv) {
      document.documentElement.style.setProperty('--cms-viewport-height', `${Math.floor(vv.height)}px`);
      document.documentElement.style.setProperty('--cms-viewport-top', `${Math.max(0, Math.floor(vv.offsetTop))}px`);
      const keyboardLikelyOpen = window.innerHeight - vv.height > 120;
      const cmFocused = Boolean(document.activeElement?.closest?.('.CodeMirror'));
      document.body.classList.toggle('mobile-writing-active', keyboardLikelyOpen && cmFocused);
      refreshEditor();
    }
  }

  window.visualViewport?.addEventListener('resize', syncVisualViewport);
  window.visualViewport?.addEventListener('scroll', syncVisualViewport);
  document.addEventListener('focusin', (event) => {
    if (isMobile() && (event.target?.closest?.('.CodeMirror') || event.target?.closest?.('dialog'))) {
      syncVisualViewport();
      setTimeout(syncVisualViewport, 120);
    }
  });
  document.addEventListener('focusout', () => {
    if (isMobile()) {
      setTimeout(() => {
        if (!document.activeElement?.closest?.('.CodeMirror')) {
          document.body.classList.remove('mobile-writing-active');
        }
      }, 120);
    }
  });

  // Mobile touch gestures (ported from md.waylonwalker.com):
  // - Swipe left/right between Write and Preview (avoiding 32px screen edges)
  // - Pull down at top of Write/Preview (>= 72px) to open Command Palette
  let touchState = null;
  const EDGE_INSET = 32;
  const SWIPE_X_THRESHOLD = 52;
  const PULL_Y_THRESHOLD = 72;

  function hidePullCue() {
    if (!mobilePullCue) return;
    mobilePullCue.hidden = true;
    mobilePullCue.dataset.ready = 'false';
  }

  function isPaneAtTop(view) {
    if (view === 'write') {
      const scroller = document.querySelector('#editor-pane .CodeMirror-scroll');
      return !scroller || scroller.scrollTop <= 2;
    }
    if (view === 'preview') {
      try {
        const win = document.querySelector('#preview')?.contentWindow;
        return !win || win.scrollY <= 2;
      } catch {
        return true;
      }
    }
    return false;
  }

  workbench.addEventListener('touchstart', (event) => {
    if (!isMobile() || event.touches.length !== 1) return;
    const view = workbench.dataset.mobileView || 'content';
    if (view !== 'write' && view !== 'preview') return;
    const touch = event.touches[0];
    touchState = {
      x: touch.clientX,
      y: touch.clientY,
      view,
      atTop: isPaneAtTop(view),
      peakPull: 0,
    };
  }, { passive: true });

  workbench.addEventListener('touchmove', (event) => {
    if (!touchState || event.touches.length !== 1) return;
    const touch = event.touches[0];
    const dx = touch.clientX - touchState.x;
    const dy = touch.clientY - touchState.y;
    if (touchState.atTop && dy > 24 && Math.abs(dx) < 42) {
      touchState.peakPull = Math.max(touchState.peakPull, dy);
      if (mobilePullCue) {
        mobilePullCue.hidden = false;
        const ready = dy >= PULL_Y_THRESHOLD;
        mobilePullCue.dataset.ready = String(ready);
        mobilePullCue.textContent = ready ? 'Release for Commands' : 'Pull for Commands';
      }
    } else if (dy < 16) {
      hidePullCue();
    }
  }, { passive: true });

  workbench.addEventListener('touchend', (event) => {
    if (!touchState) return;
    const state = touchState;
    touchState = null;
    hidePullCue();
    const touch = event.changedTouches[0];
    if (!touch) return;
    const dx = touch.clientX - state.x;
    const dy = touch.clientY - state.y;

    if (state.atTop && dy >= PULL_Y_THRESHOLD && Math.abs(dx) <= 48) {
      openCommandPalette();
      return;
    }
    if (state.x <= EDGE_INSET || state.x >= window.innerWidth - EDGE_INSET) return;
    if (Math.abs(dx) >= SWIPE_X_THRESHOLD && Math.abs(dx) > Math.abs(dy) * 1.4) {
      if (state.view === 'write' && dx < 0) setMobileView('preview');
      else if (state.view === 'preview' && dx > 0) setMobileView('write');
    }
  }, { passive: true });

  workbench.addEventListener('touchcancel', () => {
    touchState = null;
    hidePullCue();
  }, { passive: true });

  // Ensure primary buttons (Publish post, Create file) always have crisp contrast
  // even when a site palette uses a bright amber/gold --color-primary without
  // setting --color-on-primary.
  function ensurePrimaryContrast() {
    const probe = document.createElement('span');
    probe.style.cssText = 'position:fixed;visibility:hidden;pointer-events:none;color:var(--color-primary)';
    document.body.append(probe);
    const rgb = getComputedStyle(probe).color;
    probe.remove();
    const match = rgb.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)/);
    if (!match) return;
    const [r, g, b] = match.slice(1, 4).map((v) => {
      const c = Number(v) / 255;
      return c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4;
    });
    const luminance = 0.2126 * r + 0.7152 * g + 0.0722 * b;
    // Only write when needed. The observer also watches the root style
    // attribute, so unconditional setProperty creates a self-triggering loop.
    const ink = luminance > 0.179 ? '#11151d' : '#ffffff';
    for (const token of ['--btn-on-accent', '--cms-on-accent']) {
      if (document.documentElement.style.getPropertyValue(token) !== ink) {
        document.documentElement.style.setProperty(token, ink);
      }
    }
  }
  ensurePrimaryContrast();
  new MutationObserver(ensurePrimaryContrast).observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-theme', 'data-palette', 'class', 'style'],
  });

  // Split Publish menu (simple 1-click Publish + progressive disclosure for
  // advanced review workflows) and floating Publish Toast visible on desktop & mobile.
  function openWorkspaceChanges() {
    if (isMobile()) {
      setMobileView('content');
    } else if (workbench.classList.contains('sidebar-hidden')) {
      // Reuse the existing navigation action so inert/aria-hidden are restored.
      if (window.MarkataShortcuts?.execute) {
        window.MarkataShortcuts.execute('workbench.toggleSidebar');
      } else {
        sidebarToggle?.click();
      }
    }
    const changesTab = document.querySelector('#workspace-tab-changes');
    changesTab?.click();
  }

  function initPublishWorkflow() {
    const publishBtn = document.querySelector('#publish-open');
    if (!publishBtn || publishBtn.closest('.cms-publish-split')) return;

    const split = document.createElement('div');
    split.className = 'cms-publish-split';
    publishBtn.parentElement.insertBefore(split, publishBtn);
    split.append(publishBtn);

    const toggle = document.createElement('button');
    toggle.type = 'button';
    toggle.id = 'publish-menu-toggle';
    toggle.className = 'cms-publish-menu-toggle';
    toggle.setAttribute('aria-label', 'More publish and review options');
    toggle.setAttribute('aria-haspopup', 'menu');
    toggle.setAttribute('aria-expanded', 'false');
    toggle.title = 'Review changes, save as draft, or inspect workspace changes';
    toggle.innerHTML = '<svg viewBox="0 0 24 24" aria-hidden="true"><path d="m6 9 6 6 6-6"/></svg>';
    split.append(toggle);

    const menu = document.createElement('div');
    menu.id = 'publish-menu';
    menu.className = 'cms-publish-menu';
    menu.setAttribute('role', 'menu');
    menu.setAttribute('aria-label', 'Publish and review options');
    menu.hidden = true;
    menu.innerHTML = `
      <div class="cms-publish-menu-head">
        <strong>Publishing</strong>
        <button type="button" class="cms-publish-menu-close" aria-label="Close publish menu">✕</button>
      </div>
      <button type="button" role="menuitem" data-publish-action="quick">
        <span class="cms-pm-title">Publish now <kbd>Ctrl+S</kbd></span>
        <span class="cms-pm-desc">Mark public, commit, and push this post immediately</span>
      </button>
      <button type="button" role="menuitem" data-publish-action="review">
        <span class="cms-pm-title">Review before publishing… <kbd>Ctrl+Shift+↵</kbd></span>
        <span class="cms-pm-desc">Compare diff, rendered preview, edit commit message, or commit only</span>
      </button>
      <button type="button" role="menuitem" data-publish-action="draft">
        <span class="cms-pm-title">Save as draft</span>
        <span class="cms-pm-desc">Set published: false and draft: true locally without pushing</span>
      </button>
      <hr class="cms-publish-menu-sep">
      <button type="button" role="menuitem" data-publish-action="changes">
        <span class="cms-pm-title">All workspace changes <span class="cms-pm-badge" id="publish-menu-changes-count" hidden>0</span></span>
        <span class="cms-pm-desc">Inspect every uncommitted file across your site</span>
      </button>
      <button type="button" role="menuitem" data-publish-action="discard" class="cms-pm-danger">
        <span class="cms-pm-title">Discard unpublished changes…</span>
        <span class="cms-pm-desc">Open review to safely restore this file to Git HEAD</span>
      </button>
    `;
    document.body.append(menu);

    function syncMenuChangesCount() {
      const badge = menu.querySelector('#publish-menu-changes-count');
      const source = document.querySelector('#workspace-changes-count, .workspace-changes-count');
      if (!badge) return;
      const count = parseInt((source?.textContent || '').replace(/\D+/g, ''), 10) || 0;
      badge.textContent = String(count);
      badge.hidden = count <= 0;
    }

    function positionMenu() {
      if (isMobile()) {
        menu.style.top = '';
        menu.style.right = '';
        return;
      }
      const rect = split.getBoundingClientRect();
      menu.style.top = `${Math.round(rect.bottom + 8)}px`;
      menu.style.right = `${Math.max(12, Math.round(window.innerWidth - rect.right))}px`;
    }

    function setMenuOpen(open) {
      menu.hidden = !open;
      toggle.setAttribute('aria-expanded', String(open));
      if (open) {
        syncMenuChangesCount();
        positionMenu();
        menu.querySelector('[role="menuitem"]')?.focus();
      }
    }

    toggle.addEventListener('click', (event) => {
      event.stopPropagation();
      setMenuOpen(menu.hidden);
    });
    menu.querySelector('.cms-publish-menu-close')?.addEventListener('click', () => setMenuOpen(false));

    menu.addEventListener('click', (event) => {
      const item = event.target.closest('[data-publish-action]');
      if (!item) return;
      const action = item.dataset.publishAction;
      setMenuOpen(false);
      if (action === 'quick') window.MarkataPublishReview?.quickPublish?.();
      else if (action === 'review' || action === 'discard') window.MarkataPublishReview?.open?.();
      else if (action === 'draft') window.MarkataPublishReview?.saveAsDraft?.();
      else if (action === 'changes') openWorkspaceChanges();
    });

    document.addEventListener('click', (event) => {
      if (!menu.hidden && !menu.contains(event.target) && !split.contains(event.target)) {
        setMenuOpen(false);
      }
    });
    document.addEventListener('keydown', (event) => {
      if (event.key === 'Escape' && !menu.hidden) {
        event.preventDefault();
        setMenuOpen(false);
        toggle.focus();
      }
    });
    window.addEventListener('resize', () => {
      if (!menu.hidden) positionMenu();
    });
  }

  // Wait for publish-review.js to inject #publish-open, then upgrade it to a split button.
  if (document.querySelector('#publish-open')) {
    initPublishWorkflow();
  } else {
    const publishObs = new MutationObserver(() => {
      if (document.querySelector('#publish-open')) {
        publishObs.disconnect();
        initPublishWorkflow();
      }
    });
    publishObs.observe(document.querySelector('#workbench-header') || document.body, {
      childList: true,
      subtree: true,
    });
  }

  // Floating Publish Toast (visible on mobile & desktop)
  const toast = document.createElement('div');
  toast.id = 'cms-publish-toast';
  toast.className = 'cms-publish-toast';
  toast.setAttribute('role', 'status');
  toast.setAttribute('aria-live', 'polite');
  toast.hidden = true;
  toast.innerHTML = `
    <span class="cms-toast-dot" aria-hidden="true"></span>
    <span class="cms-toast-message"></span>
    <button type="button" class="cms-toast-action" hidden></button>
    <button type="button" class="cms-toast-close" aria-label="Dismiss status">✕</button>
  `;
  document.body.append(toast);

  let toastTimer = null;
  const toastMessage = toast.querySelector('.cms-toast-message');
  const toastAction = toast.querySelector('.cms-toast-action');
  const toastClose = toast.querySelector('.cms-toast-close');

  function hideToast() {
    clearTimeout(toastTimer);
    toast.classList.remove('is-dragging');
    toast.style.transform = '';
    toast.style.opacity = '';
    toast.hidden = true;
  }
  toastClose.addEventListener('click', hideToast);

  // 1:1 finger-tracked upward swipe-to-dismiss on the floating publish toast (matches entrance axis).
  let toastDrag = null;
  toast.addEventListener('pointerdown', (event) => {
    if (event.target.closest('button')) return;
    toastDrag = {
      pointerId: event.pointerId,
      startY: event.clientY,
      startTime: performance.now(),
      dy: 0,
    };
    toast.classList.add('is-dragging');
    toast.setPointerCapture?.(event.pointerId);
  });
  toast.addEventListener('pointermove', (event) => {
    if (!toastDrag || event.pointerId !== toastDrag.pointerId) return;
    const rawDy = event.clientY - toastDrag.startY;
    // Allow 1:1 upward drag; rubber-band downward drag gently.
    const dy = rawDy < 0 ? rawDy : rawDy * 0.22;
    toastDrag.dy = dy;
    toast.style.transform = `translateY(${dy.toFixed(1)}px)`;
    toast.style.opacity = String(Math.max(0.25, 1 + Math.min(0, dy) / 110));
  });
  function endToastDrag(event) {
    if (!toastDrag || event.pointerId !== toastDrag.pointerId) return;
    const { dy, startTime } = toastDrag;
    const dt = Math.max(1, performance.now() - startTime);
    const velocity = dy / dt; // px/ms (negative = upward flick)
    toastDrag = null;
    toast.classList.remove('is-dragging');
    if (dy < -36 || velocity < -0.45) {
      toast.style.transform = 'translateY(-110%)';
      toast.style.opacity = '0';
      setTimeout(hideToast, 180);
    } else {
      toast.style.transform = '';
      toast.style.opacity = '';
    }
  }
  toast.addEventListener('pointerup', endToastDrag);
  toast.addEventListener('pointercancel', endToastDrag);

  // Keep mobile browser/PWA status bar color matched to the active header surface.
  function syncThemeColorMeta() {
    const metas = document.querySelectorAll('meta[name="theme-color"]');
    const topHeader = document.querySelector('body > header');
    if (!metas.length || !topHeader) return;
    const headerBg = getComputedStyle(topHeader).backgroundColor;
    if (headerBg && headerBg !== 'rgba(0, 0, 0, 0)' && headerBg !== 'transparent') {
      metas.forEach((meta) => meta.setAttribute('content', headerBg));
    }
  }
  syncThemeColorMeta();
  requestAnimationFrame(syncThemeColorMeta);
  window.addEventListener('load', syncThemeColorMeta);
  new MutationObserver(() => requestAnimationFrame(syncThemeColorMeta)).observe(
    document.documentElement,
    { attributes: true, attributeFilter: ['data-theme', 'data-palette', 'style', 'class'] }
  );

  window.addEventListener('markata:publish-status', (event) => {
    const { message, type = 'info', pending = false } = event.detail || {};
    if (!message) return;
    clearTimeout(toastTimer);
    toast.dataset.state = type;
    toastMessage.textContent = message;
    if (pending) {
      toastAction.textContent = 'Retry';
      toastAction.hidden = false;
      toastAction.onclick = () => {
        hideToast();
        window.MarkataPublishReview?.quickPublish?.();
      };
    } else if (type === 'warning' || /No unpublished/i.test(message)) {
      toastAction.textContent = 'Open Review';
      toastAction.hidden = false;
      toastAction.onclick = () => {
        hideToast();
        window.MarkataPublishReview?.open?.();
      };
    } else {
      toastAction.hidden = true;
      toastAction.onclick = null;
    }
    toast.hidden = false;
    if (type !== 'loading' && type !== 'warning') {
      toastTimer = setTimeout(hideToast, 5500);
    }
  });

  // Close any open <dialog> modal when the user clicks its ::backdrop outside the dialog box.
  // Require both pointerdown and click outside bounds so dragging text selection out doesn't dismiss.
  let backdropDownDialog = null;
  function isOutsideDialogBounds(dialog, event) {
    if (event.target !== dialog) return false;
    if (event.clientX === 0 && event.clientY === 0) return false;
    const rect = dialog.getBoundingClientRect();
    return (
      event.clientX < rect.left ||
      event.clientX > rect.right ||
      event.clientY < rect.top ||
      event.clientY > rect.bottom
    );
  }
  document.addEventListener('pointerdown', (event) => {
    const dialog = event.target instanceof HTMLDialogElement ? event.target : null;
    backdropDownDialog =
      dialog && dialog.open && !dialog.hasAttribute('data-no-backdrop-close') && isOutsideDialogBounds(dialog, event)
        ? dialog
        : null;
  });
  document.addEventListener('click', (event) => {
    const dialog = backdropDownDialog;
    backdropDownDialog = null;
    if (!dialog || !dialog.open || dialog.hasAttribute('data-no-backdrop-close')) return;
    if (isOutsideDialogBounds(dialog, event)) {
      dialog.close();
    }
  });

  mobileQuery.addEventListener('change', () => {
    syncVisualViewport();
    if (isMobile()) {
      setMobileView(workbench.dataset.mobileView || 'content');
    } else {
      refreshEditor();
    }
  });

  window.MarkataCMSShell = { setMobileView, openCommandPalette, openWorkspaceChanges };
})();
