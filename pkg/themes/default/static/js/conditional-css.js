/**
 * Conditional CSS Runtime Loader
 *
 * Handles CSS injection after view-transition SPA navigations.
 * When markata-go uses view transitions, only document.body.innerHTML is
 * swapped - <head> content persists. This means:
 *   - CSS loaded on page A stays available on page B (harmless, cached).
 *   - CSS NOT loaded on page A will be missing on page B if page B needs it.
 *
 * This script detects DOM elements requiring conditional CSS and injects
 * the appropriate <link> tags if they are not already present in <head>.
 */
(function() {
  'use strict';

  // Map of CSS base names to DOM selectors that require them.
  var conditionalCSS = {
    'admonitions': '.admonition',
    'code':        'pre > code, .highlight, code[class*="language-"]',
    'chroma':      '.chroma',
    'cards':       '.card, .posts-list',
    'webmentions': '.webmentions',
    'encryption':  '.encrypted-content, [data-encrypted]',
    'home':        '.home-hero, .home-feeds',
    'feeds':       '.feeds-page, .feed.h-feed, .feed.feed-simple',
    'roads':       '.road, [data-road]'
  };

  var conditionalJS = {
    'feed-sparklines-v2': '.feed-sparkline-wrap, .feed-header-sparkline',
    'code-chrome': '.post-content pre > code, .post-content .chroma pre, .post-content .highlight pre',
    'roads': '.road, [data-road]'
  };

  var navHoverSuppressed = false;
  var navHoverReleaseArmed = false;
  var NAV_HOVER_ATTRIBUTE = 'data-nav-hover-suppressed';
  var NAV_HOVER_STYLE_ID = 'markata-nav-hover-reset-style';
  var deferredCSS = { admonitions: true, code: true, chroma: true };
  var deferredCSSObservers = {};

  /**
   * Check if a stylesheet containing the given base name is already loaded.
   */
  function hasCSS(baseName) {
    var links = document.querySelectorAll('link[rel="stylesheet"]');
    for (var i = 0; i < links.length; i++) {
      var sources = links[i].getAttribute('data-css-sources') || '';
      if ((links[i].href && links[i].href.indexOf(baseName) !== -1) || sources.indexOf('/' + baseName + '.') !== -1) {
        return true;
      }
    }
    return false;
  }

  /**
   * Inject a CSS file by base name. Looks for a matching link already in
   * the page to derive the full hashed path; if not found, falls back to
   * /css/{baseName}.css.
   */
  function injectCSS(baseName) {
    var href = '/css/' + baseName + '.css';
    var link = document.createElement('link');
    link.rel = 'stylesheet';
    link.href = href;
    document.head.appendChild(link);
  }

  function injectCSSWhenNearViewport(baseName, selector) {
    var previousObserver = deferredCSSObservers[baseName];
    if (previousObserver) previousObserver.disconnect();

    var elements = document.querySelectorAll(selector);
    if (!elements.length) return;

    if (!window.IntersectionObserver) {
      injectCSS(baseName);
      return;
    }

    var observer = new IntersectionObserver(function(entries) {
      for (var i = 0; i < entries.length; i++) {
        if (!entries[i].isIntersecting) continue;
        observer.disconnect();
        if (deferredCSSObservers[baseName] === observer) delete deferredCSSObservers[baseName];
        if (!hasCSS(baseName)) injectCSS(baseName);
        return;
      }
    }, { rootMargin: '2400px 0px' });

    deferredCSSObservers[baseName] = observer;
    for (var i = 0; i < elements.length; i++) observer.observe(elements[i]);
  }

  function isFeedCardContent(selector) {
    var elements = document.querySelectorAll(selector);
    if (!elements.length) return false;
    for (var i = 0; i < elements.length; i++) {
      if (!elements[i].closest || !elements[i].closest('.card')) return false;
    }
    return true;
  }

  function hasScript(baseName) {
    var scripts = document.querySelectorAll('script[src]');
    for (var i = 0; i < scripts.length; i++) {
      if (scripts[i].src && scripts[i].src.indexOf(baseName) !== -1) {
        return true;
      }
    }
    return false;
  }

  function injectJS(baseName) {
    var script = document.createElement('script');
    script.src = '/js/' + baseName + '.js';
    script.defer = true;
    document.body.appendChild(script);
  }

  /**
   * Lite YouTube's heavyweight assets are page-scoped. The normal inline
   * loader does not execute when View Transitions swap the page wrapper, so
   * inspect the newly swapped wrapper and load missing assets here as well.
   * Asset URLs travel with the wrapper to preserve site-specific CDN overrides.
   */
  function ensureLiteYouTubeAssets() {
    if (!document.querySelector('lite-youtube')) return;
    var page = document.getElementById('view-transition-page');
    if (!page) return;

    var cssURL = page.getAttribute('data-markata-lite-youtube-css');
    var jsURL = page.getAttribute('data-markata-lite-youtube-js');
    if (!cssURL || !jsURL) return;

    function absoluteURL(value) {
      try {
        return new URL(value, document.baseURI).href;
      } catch (_) {
        return '';
      }
    }

    cssURL = absoluteURL(cssURL);
    jsURL = absoluteURL(jsURL);
    if (!cssURL || !jsURL) return;

    var styles = document.querySelectorAll('link[rel="stylesheet"]');
    var hasStyle = Array.prototype.some.call(styles, function (link) {
      return link.href === cssURL;
    });
    if (!hasStyle) {
      var style = document.createElement('link');
      style.rel = 'stylesheet';
      style.href = cssURL;
      style.onerror = function () { style.remove(); };
      document.head.appendChild(style);
    }

    // Custom element registrations persist across document transitions, even
    // if the script node itself was removed during an earlier body swap.
    if (window.customElements && window.customElements.get('lite-youtube')) return;
    var scripts = document.querySelectorAll('script[src]');
    var hasScript = Array.prototype.some.call(scripts, function (script) {
      return script.src === jsURL;
    });
    if (!hasScript) {
      var loader = document.createElement('script');
      loader.src = jsURL;
      loader.onerror = function () { loader.remove(); };
      document.body.appendChild(loader);
    }
  }

  /**
   * Scan the current DOM and inject any missing conditional CSS.
   */
  function loadConditionalCSS() {
    for (var baseName in conditionalCSS) {
      if (!conditionalCSS.hasOwnProperty(baseName)) continue;
      var selector = conditionalCSS[baseName];
      var hasMatchingElements = !!document.querySelector(selector);
      if (deferredCSSObservers[baseName] && (!hasMatchingElements || hasCSS(baseName))) {
        deferredCSSObservers[baseName].disconnect();
        delete deferredCSSObservers[baseName];
      }
      if (hasMatchingElements && !hasCSS(baseName)) {
        if (deferredCSS[baseName] && isFeedCardContent(selector)) {
          injectCSSWhenNearViewport(baseName, selector);
        } else {
          injectCSS(baseName);
        }
      }
    }

    for (var scriptName in conditionalJS) {
      if (!conditionalJS.hasOwnProperty(scriptName)) continue;
      var scriptSelector = conditionalJS[scriptName];
      if (!document.querySelector(scriptSelector)) continue;
      // SPA body swaps can remove the script tag while retaining the event
      // handler. Re-run the existing initializer instead of loading a second
      // copy of the same code-block controls.
      if (scriptName === 'code-chrome' && typeof window.markataInitCodeChrome === 'function') {
        window.markataInitCodeChrome();
        continue;
      }
      if (!hasScript(scriptName)) {
        injectJS(scriptName);
      }
    }

    ensureLiteYouTubeAssets();

    // Handle decryption.js for encrypted content
    if (document.querySelector('.encrypted-content, [data-encrypted]')) {
      var scripts = document.querySelectorAll('script[src]');
      var hasDecryption = false;
      for (var j = 0; j < scripts.length; j++) {
        if (scripts[j].src.indexOf('decryption') !== -1) {
          hasDecryption = true;
          break;
        }
      }
      if (!hasDecryption) {
        var script = document.createElement('script');
        script.src = '/js/decryption.js';
        script.defer = true;
        document.body.appendChild(script);
      }
    }

    // Handle GLightbox re-initialization after navigation
    if (window.initGLightbox && document.querySelector('.glightbox')) {
      window.initGLightbox();
    }
  }

  function applyNavHoverSuppression() {
    ensureNavHoverResetStyle();
    document.documentElement.setAttribute(NAV_HOVER_ATTRIBUTE, 'true');
  }

  function clearNavHoverSuppression() {
    navHoverSuppressed = false;
    navHoverReleaseArmed = false;
    document.documentElement.removeAttribute(NAV_HOVER_ATTRIBUTE);
  }

  function ensureNavHoverResetStyle() {
    if (document.getElementById(NAV_HOVER_STYLE_ID)) return;

    var style = document.createElement('style');
    style.id = NAV_HOVER_STYLE_ID;
    style.textContent = [
      'html[' + NAV_HOVER_ATTRIBUTE + '] .nav-entry--preview:hover:not(:focus-within) .nav-preview {',
      '  opacity: 0 !important;',
      '  visibility: hidden !important;',
      '  pointer-events: none !important;',
      '  transform: translateY(-0.35rem) !important;',
      '}',
      'html[' + NAV_HOVER_ATTRIBUTE + '] .nav-entry--preview:hover:not(:focus-within) > .nav-link::after {',
      '  opacity: 0.4 !important;',
      '  transform: scaleX(0.35) !important;',
      '}',
      'html[' + NAV_HOVER_ATTRIBUTE + '] .nav-entry--preview:hover:not(:focus-within) {',
      '  z-index: auto !important;',
      '}'
    ].join('\n');
    style.setAttribute('data-markata-persist', 'true');
    document.head.appendChild(style);
  }

  function suppressNavHoverUntilPointerMoves() {
    navHoverSuppressed = true;
    navHoverReleaseArmed = false;
    applyNavHoverSuppression();
  }

  function shouldResetNavHover(event, link) {
    if (!document.startViewTransition) return false;
    if (window.VIEW_TRANSITIONS_CONFIG && window.VIEW_TRANSITIONS_CONFIG.enabled === false) return false;
    if (!link || event.defaultPrevented || event.button !== 0) return false;
    if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return false;
    if (link.target && link.target !== '_self') return false;
    if (link.hasAttribute('download')) return false;

    var target;
    try {
      target = new URL(link.href, window.location.href);
    } catch (_) {
      return false;
    }

    if (target.origin !== window.location.origin) return false;
    return target.pathname !== window.location.pathname || target.search !== window.location.search;
  }

  // A SPA navigation swaps in a new nav underneath an unmoved pointer. Without
  // an explicit reset, the new element immediately matches :hover and makes a
  // preview look like it survived the page transition. Keep pointer-only hover
  // suppressed through the swap and only release it on pointer movement after
  // the transition is complete; :focus-within still works throughout.
  document.addEventListener('click', function(event) {
    var target = event.target;
    if (!target || !target.closest) return;

    var link = target.closest('.nav-entry--preview a[href]');
    if (!shouldResetNavHover(event, link)) return;

    suppressNavHoverUntilPointerMoves();
  }, true);

  window.addEventListener('pointermove', function() {
    if (!navHoverSuppressed || !navHoverReleaseArmed) return;
    clearNavHoverSuppression();
  }, true);

  // Feed embeds and generated card excerpts can require conditional styles
  // that are not visible to the post-level template markers.
  document.addEventListener('DOMContentLoaded', function() {
    loadConditionalCSS();
    ensureNavHoverResetStyle();
  });
  window.addEventListener('view-transition-complete', function() {
    loadConditionalCSS();
    ensureNavHoverResetStyle();
    if (navHoverSuppressed) {
      applyNavHoverSuppression();
      navHoverReleaseArmed = true;
    }
  });
})();
