/**
 * Keep <meta name="theme-color"> matched to the rendered page background.
 *
 * Mobile status bars and installed-app title bars take their color from
 * theme-color. One static value cannot match every palette and mode, so this
 * resolves the painted background color and follows palette, mode, and
 * view-transition changes.
 */
(function () {
  'use strict';

  if (window.__markataThemeColor) return;
  window.__markataThemeColor = true;

  var root = document.documentElement;
  var meta = document.querySelector('meta[name="theme-color"]');
  if (!meta) {
    meta = document.createElement('meta');
    meta.name = 'theme-color';
    document.head.appendChild(meta);
  }

  function painted(el) {
    var bg = el && getComputedStyle(el).backgroundColor;
    return bg && bg !== 'transparent' && bg !== 'rgba(0, 0, 0, 0)' ? bg : '';
  }

  function sync() {
    var bg = painted(document.body) || painted(root);
    if (bg && meta.content !== bg) meta.content = bg;
  }

  function schedule() { requestAnimationFrame(sync); }

  sync();
  new MutationObserver(schedule).observe(root, {
    attributes: true,
    attributeFilter: ['data-theme', 'data-palette', 'class'],
  });
  window.addEventListener('view-transition-complete', schedule);
  // Palette stylesheets can arrive after the attribute change.
  document.addEventListener('load', function (event) {
    if (event.target && event.target.tagName === 'LINK') schedule();
  }, true);
})();
