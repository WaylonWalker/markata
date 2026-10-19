/**
 * Tiny fallback for sidebar drawer controls.
 *
 * Keep navigation usable if the optional view-transitions bundle is disabled,
 * missing, or fails to initialize. The normal controller remains authoritative:
 * its _sidebarToggleBound marker prevents double handling.
 *
 * This module does not fetch feed data, read layout for closed drawers, or
 * eagerly initialize optional sidebar features.
 */
(function () {
  'use strict';

  var DRAWERS = '.feed-sidebar, .doc-sidebar, .content-sidebar';
  var OVERLAYS = '.feed-sidebar, .doc-sidebar';
  var STORAGE_PREFIX = 'markata-sidebar:';
  var media = window.matchMedia && window.matchMedia('(max-width: 1200px)');

  function overlayMode() { return !!(media && media.matches); }
  function isDrawer(drawer) {
    return !!drawer && !(overlayMode() && drawer.classList.contains('content-sidebar'));
  }
  function side(drawer) { return drawer.classList.contains('feed-sidebar--left') ||
    drawer.classList.contains('doc-sidebar--left') ||
    drawer.classList.contains('content-sidebar--left') ? 'left' : 'right'; }
  function button(drawer) { return drawer && drawer.querySelector(':scope > .sidebar-toggle'); }
  function needsFallback(drawer) {
    var control = button(drawer);
    return !!control && !control._sidebarToggleBound && isDrawer(drawer);
  }
  function savedState(drawer) {
    try { return localStorage.getItem(STORAGE_PREFIX + side(drawer)); }
    catch (_) { return null; }
  }
  function saveState(drawer, open) {
    try { localStorage.setItem(STORAGE_PREFIX + side(drawer), open ? 'open' : 'closed'); }
    catch (_) { /* storage is optional */ }
  }
  function sync(drawer) {
    var control = button(drawer);
    if (!control) return;
    var open = drawer.classList.contains('sidebar--pinned');
    var label = control.dataset.sidebarLabel || 'sidebar';
    var action = (open ? 'Close ' : 'Open ') + label + (label === 'sidebar' ? '' : ' sidebar');
    control.setAttribute('aria-expanded', String(open));
    control.setAttribute('aria-label', action);
    control.title = control.dataset.sidebarHotkey ? action + ' (' + control.dataset.sidebarHotkey + ')' : action;
  }
  function updateTop(drawer) {
    if (overlayMode()) return;
    var header = document.querySelector('.site-header');
    var bottom = header ? Math.max(0, Math.round(header.getBoundingClientRect().bottom)) : 0;
    drawer.style.setProperty('--sidebar-top', bottom + 'px');
  }
  function setOpen(drawer, open, persist) {
    if (!needsFallback(drawer)) return;
    if (open && overlayMode()) {
      document.querySelectorAll(OVERLAYS).forEach(function (other) {
        if (other !== drawer && needsFallback(other)) {
          other.classList.remove('sidebar--pinned');
          sync(other);
        }
      });
    }
    drawer.classList.toggle('sidebar--pinned', open);
    sync(drawer);
    if (open) {
      updateTop(drawer); // never force layout for a closed drawer
      if (typeof window.initFeedSidebarScroll === 'function') window.initFeedSidebarScroll();
      if (overlayMode()) {
        requestAnimationFrame(function () {
          var target = drawer.querySelector('.feed-nav-item--active a, .toc-link--active, .feed-nav a[href], .toc a[href]');
          if (target && drawer.classList.contains('sidebar--pinned')) target.focus({ preventScroll: true });
        });
      }
    }
    if (persist && !overlayMode()) saveState(drawer, open);
  }
  function findDrawer(selector) {
    return Array.prototype.find.call(document.querySelectorAll(selector), needsFallback);
  }
  function toggleSide(which) {
    var selector = which === 'left'
      ? '.feed-sidebar--left, .doc-sidebar--left, .content-sidebar--left'
      : '.feed-sidebar--right, .doc-sidebar--right, .content-sidebar--right';
    var drawer = findDrawer(selector);
    if (drawer) setOpen(drawer, !drawer.classList.contains('sidebar--pinned'), true);
  }
  function restore() {
    document.querySelectorAll(DRAWERS).forEach(function (drawer) {
      if (!needsFallback(drawer)) return;
      if (!overlayMode() && savedState(drawer) === 'open') {
        drawer.classList.add('sidebar--pinned');
        updateTop(drawer);
      } else if (overlayMode()) {
        drawer.classList.remove('sidebar--pinned');
      }
      sync(drawer);
    });
  }

  document.addEventListener('click', function (event) {
    var target = event.target;
    if (!target || typeof target.closest !== 'function') return;
    var control = target.closest('.sidebar-toggle');
    var drawer = control && control.closest(DRAWERS);
    if (drawer && control === button(drawer) && needsFallback(drawer)) {
      event.preventDefault();
      setOpen(drawer, !drawer.classList.contains('sidebar--pinned'), true);
      return;
    }
    if (!overlayMode()) return;
    var open = document.querySelector(OVERLAYS.split(', ').map(function (s) { return s + '.sidebar--pinned'; }).join(', '));
    if (!open || !needsFallback(open)) return;
    if (!open.contains(target)) setOpen(open, false, false);
    else if (target.closest('a[href]')) setOpen(open, false, false);
  });

  document.addEventListener('keydown', function (event) {
    if (event.key !== 'Escape' || !overlayMode()) return;
    var open = document.querySelector('.feed-sidebar.sidebar--pinned, .doc-sidebar.sidebar--pinned');
    if (!open || !needsFallback(open)) return;
    setOpen(open, false, false);
    event.preventDefault();
    var control = button(open);
    if (control) control.focus({ preventScroll: true });
  });

  // navigation-shortcuts.js invokes this for b / Shift+B even without SPA.
  if (typeof window.toggleSidebarPinned !== 'function') window.toggleSidebarPinned = toggleSide;
  if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', restore);
  else restore();
  window.addEventListener('view-transition-complete', restore);
  if (media && media.addEventListener) media.addEventListener('change', restore);
})();
