/**
 * View Transitions API - Global Navigation Interceptor
 *
 * Provides smooth animated transitions for all internal navigation.
 * Gracefully degrades on unsupported browsers.
 *
 * Handles:
 * - Wikilinks, card links, nav/footer links, post navigation
 * - Back/forward browser buttons
 * - Script re-initialization after transitions
 * - Error handling and fallbacks
 *
 * Configuration (set window.VIEW_TRANSITIONS_CONFIG in your template):
 * {
 *   enabled: true,              // Enable/disable view transitions
 *   debug: false,               // Log debug messages
 *   skipClasses: [],            // Additional CSS classes to skip
 *   skipSelectors: [],          // Additional selectors to skip (e.g., '[data-no-transition]')
 *   transitionDuration: 300,    // Default transition duration in ms (for prefetching timeout)
 *   updateMeta: true,           // Update meta tags on navigation
 *   scrollToTop: true,          // Scroll to top on navigation (unless hash present)
 * }
 */

(function() {
  'use strict';

  // Default configuration
  const DEFAULT_CONFIG = {
    enabled: true,
    debug: false,
    skipClasses: [],
    skipSelectors: [],
    transitionDuration: 120,
    updateMeta: true,
    scrollToTop: true,
  };

  const HTML_EXTENSIONS = new Set(['.html', '.htm', '.xhtml']);
  const MAX_PREFETCHED_DOCUMENTS = 8;
  const PREFETCH_DEBOUNCE_MS = 120;
  // Keep hover/focus work useful for an ordinary click, but revalidate a
  // document left in the in-memory cache for longer than this window.
  const PREFETCH_FRESHNESS_MS = 5000;
  const STAGED_STYLESHEET_ATTR = 'data-markata-staged-stylesheet';
  const STYLESHEET_READY_TIMEOUT_MS = 3500;
  // Visitor-controlled theme state must survive in-site navigation; the new
  // document only carries server defaults.
  const RUNTIME_HTML_ATTRIBUTES = new Set(['data-theme', 'data-text-size', 'data-palette', 'data-aesthetic']);
  const RUNTIME_HTML_CLASS_NAMES = new Set(['dark']);
  const RUNTIME_HTML_ATTRIBUTE_PREFIXES = [
    'data-loading-',
    'data-adaptive-',
    'data-shared-transition-',
    'data-post-transition-',
  ];
  const prefetchedDocuments = new Map();
  let navigationInFlight = false;
  let lastPrefetchUrl = '';
  let lastPrefetchAt = 0;

  function getNavigationPath(url) {
    if (!url) return '';
    return url.pathname || '';
  }

  function getNavigationKey(url) {
    if (!url) return '';
    return `${url.pathname || ''}${url.search || ''}`;
  }

  // Path+query of the document currently rendered. Fragment navigation (TOC
  // and heading anchors) fires popstate without changing it, so popstate uses
  // this to leave in-page jumps to the browser instead of refetching the page.
  let renderedNavigationKey = getNavigationKey(window.location);

  function markRenderedLocation() {
    renderedNavigationKey = getNavigationKey(window.location);
  }

  function createSharedTransitionToken(path) {
    const token = String(path || '')
      .toLowerCase()
      .replace(/[^a-z0-9_-]+/g, '-')
      .replace(/^-+|-+$/g, '');

    return token || 'post';
  }

  function createSharedTransitionNames(path) {
    const token = createSharedTransitionToken(path);
    return {
      shell: `post-shell-${token}`,
      title: `post-title-${token}`,
    };
  }

  function createSidebarTransitionNames(path) {
    const token = createSharedTransitionToken(path);
    return {
      shell: `sidebar-header-${token}`,
      title: `sidebar-title-${token}`,
    };
  }

  function getNow() {
    if (window.performance && typeof window.performance.now === 'function') {
      return window.performance.now();
    }

    return Date.now();
  }

  function roundTiming(value) {
    return Math.round(value * 10) / 10;
  }

  function createNavigationMetrics(url, source) {
    return {
      url: url,
      source: source,
      prefetched: false,
      fetchMs: 0,
      parseMs: 0,
      swapMs: 0,
      reinitMs: 0,
      totalMs: 0,
      startedAt: getNow(),
    };
  }

  function finalizeNavigationMetrics(metrics) {
    if (!metrics) return;

    metrics.totalMs = roundTiming(getNow() - metrics.startedAt);

    const finalized = {
      url: metrics.url,
      source: metrics.source,
      prefetched: metrics.prefetched,
      fetchMs: roundTiming(metrics.fetchMs),
      parseMs: roundTiming(metrics.parseMs),
      swapMs: roundTiming(metrics.swapMs),
      reinitMs: roundTiming(metrics.reinitMs),
      totalMs: metrics.totalMs,
    };

    window.__lastViewTransitionMetrics = finalized;
    window.dispatchEvent(new CustomEvent('view-transition-timing', { detail: finalized }));

    if (config.debug) {
      console.log('[view-transitions] timing', finalized);
    }

    return finalized;
  }

  function normalizePathname(p) {
    if (typeof p !== 'string') return '';
    let s = p.trim();
    if (!s) return '';
    // Drop origin if an absolute URL was provided.
    try {
      if (s.startsWith('http://') || s.startsWith('https://')) {
        s = new URL(s).pathname;
      }
    } catch (_) {
      // ignore
    }
    if (!s.startsWith('/')) s = '/' + s;
    // Collapse multiple slashes.
    s = s.replace(/\/{2,}/g, '/');
    return s;
  }

  function shouldSkipSpecialRoutes(url) {
    if (!url || !url.pathname) return false;

    const pathname = url.pathname;

    // Random post endpoint performs a client-side redirect.
    // View-transition navigation swaps DOM without executing inline scripts,
    // which prevents the redirect script from running.
    const randomCandidates = new Set(['/random', '/random/']);
    const configured = normalizePathname(window.MARKATA_GO_RANDOM_POST_PATH);
    if (configured) {
      randomCandidates.add(configured);
      if (configured.endsWith('/')) {
        randomCandidates.add(configured.slice(0, -1));
      } else {
        randomCandidates.add(configured + '/');
      }
    }

    if (randomCandidates.has(pathname)) {
      if (config.debug) console.log('Skipping special route navigation:', pathname);
      return true;
    }

    return false;
  }

  function getLowercaseExtension(pathname) {
    const lastSegment = (pathname || '').split('/').pop() || '';
    const dotIndex = lastSegment.lastIndexOf('.');
    if (dotIndex <= 0) return '';
    return lastSegment.slice(dotIndex).toLowerCase();
  }

  function isNonHTMLDocumentURL(url) {
    const ext = getLowercaseExtension(url.pathname);
    if (!ext) return false; // Pretty URLs and extensionless routes are assumed HTML.
    return !HTML_EXTENSIONS.has(ext);
  }

  function shouldTransitionToURL(url) {
    if (!url) return false;

    if (url.origin !== window.location.origin) return false;
    if (shouldSkipSpecialRoutes(url)) return false;
    if (isNonHTMLDocumentURL(url)) return false;

    if (url.pathname === window.location.pathname &&
        url.search === window.location.search) {
      return false;
    }

    return true;
  }

  // Merge with user config
  const config = Object.assign({}, DEFAULT_CONFIG, window.VIEW_TRANSITIONS_CONFIG || {});

  // Exit early if disabled
  if (!config.enabled) {
    if (config.debug) console.log('View Transitions disabled via config');
    return;
  }

  // Feature detection - exit early if not supported
  if (!document.startViewTransition) {
    if (config.debug) console.log('View Transitions API not supported, using standard navigation');
    return;
  }

  if (config.debug) console.log('View Transitions API enabled', config);

  /**
   * Check if a link should use view transitions
   */
  function shouldTransition(link) {
    // Must be an <a> element with href
    if (!link || !link.href) return false;

    // Parse URL
    let url;
    try {
      url = new URL(link.href);
    } catch (e) {
      return false;
    }

    if (!shouldTransitionToURL(url)) {
      if (config.debug && isNonHTMLDocumentURL(url)) {
        console.log('Skipping non-HTML navigation:', url.href);
      }
      return false;
    }

    // Skip external links (target=_blank)
    if (link.target === '_blank') return false;

    // Skip non-default navigation targets
    if (link.target && link.target !== '_self') return false;

    // Skip download links
    if (link.hasAttribute('download')) return false;

    // Skip links with rel=external
    if (link.rel && link.rel.includes('external')) return false;

    // Skip TOC links (they have their own smooth scroll handling)
    if (link.classList.contains('toc-link')) return false;

    // Skip HTMX-managed links (they handle their own updates)
    if (link.hasAttribute('hx-get') || link.hasAttribute('hx-post')) return false;

    // Skip GLightbox links (inline lightbox should not trigger navigation)
    if (link.classList.contains('glightbox-mermaid') || link.hasAttribute('data-glightbox')) return false;

    // Skip links that explicitly opt-out
    if (link.dataset.noTransition) return false;

    // Skip user-configured classes
    for (const className of config.skipClasses) {
      if (link.classList.contains(className)) {
        if (config.debug) console.log('Skipping link with class:', className);
        return false;
      }
    }

    // Skip user-configured selectors
    for (const selector of config.skipSelectors) {
      if (link.matches(selector)) {
        if (config.debug) console.log('Skipping link matching selector:', selector);
        return false;
      }
    }

    return true;
  }

  function resolveCardRoot(link) {
    if (!link) return null;
    return link.closest('.card, .photo-figure');
  }

  function resolvePrimaryCardLink(cardRoot) {
    if (!cardRoot) return null;

    return cardRoot.querySelector([
      'a.u-url[href]',
      '.card-title a[href]',
      'a.card-cover[href]',
      'a.card-video-link[href]',
      'a.card-image-link[href]',
      'a.card-link[href]',
      'a[href]'
    ].join(', '));
  }

  function resolveCardTitleElement(cardRoot) {
    if (!cardRoot) return null;

    return cardRoot.querySelector([
      '.card-title',
      '.card-caption',
      'figcaption'
    ].join(', '));
  }

  function getSharedTransitionContext(link, targetURL) {
    const cardRoot = resolveCardRoot(link);
    if (!cardRoot) return null;

    const primaryLink = resolvePrimaryCardLink(cardRoot);
    if (!primaryLink) return null;

    let primaryURL;
    try {
      primaryURL = new URL(primaryLink.href, window.location.href);
    } catch (_) {
      return null;
    }

    if (getNavigationKey(primaryURL) !== getNavigationKey(targetURL)) {
      return null;
    }

    return {
      path: getNavigationPath(targetURL),
      names: createSharedTransitionNames(getNavigationPath(targetURL)),
      shell: cardRoot,
      title: resolveCardTitleElement(cardRoot),
      incomingShellSelector: '[data-shared-transition-path]',
      incomingTitleSelector: '[data-shared-transition-title]',
    };
  }

  function getSidebarSharedTransitionContext(link, targetURL) {
    if (!link || !link.closest('.feed-sidebar')) return null;

    const row = link.closest('[data-sidebar-transition-path]');
    if (!row) return null;

    if (row.getAttribute('data-sidebar-transition-path') !== getNavigationPath(targetURL)) {
      return null;
    }

    return {
      path: getNavigationPath(targetURL),
      names: createSidebarTransitionNames(getNavigationPath(targetURL)),
      shell: row,
      title: row.querySelector('[data-sidebar-transition-title]'),
      incomingShellSelector: '[data-sidebar-transition-header]',
      incomingTitleSelector: '[data-shared-transition-title]',
    };
  }

  function markSharedTransitionElement(element, name) {
    if (!element || !name) return;

    element.style.setProperty('view-transition-name', name);
    element.setAttribute('data-shared-transition-active', 'true');
  }

  function setSharedTransitionState(active) {
    if (active) {
      document.documentElement.setAttribute('data-shared-transition-active', 'true');
      return;
    }

    document.documentElement.removeAttribute('data-shared-transition-active');
  }

  function setPostNavigationState(context) {
    if (!context) return;

    document.documentElement.setAttribute('data-post-transition-active', 'true');
    document.documentElement.setAttribute('data-post-transition-source', context.source);
    document.documentElement.setAttribute('data-post-transition-direction', context.direction);
  }

  function clearPostNavigationState() {
    document.documentElement.removeAttribute('data-post-transition-active');
    document.documentElement.removeAttribute('data-post-transition-source');
    document.documentElement.removeAttribute('data-post-transition-direction');
  }

  function resolveSidebarFeedLinks() {
    return Array.from(document.querySelectorAll('.feed-sidebar a[href]')).filter((link) => {
      try {
        const url = new URL(link.href, window.location.href);
        return shouldTransitionToURL(url);
      } catch (_) {
        return false;
      }
    });
  }

  function getSidebarDirection(link) {
    if (!link || !link.closest('.feed-sidebar')) return null;

    const links = resolveSidebarFeedLinks();
    if (!links.length) return null;

    const currentPath = window.location.pathname;
    const activeIndex = links.findIndex((candidate) => {
      try {
        return new URL(candidate.href, window.location.href).pathname === currentPath;
      } catch (_) {
        return false;
      }
    });
    const targetIndex = links.indexOf(link);

    if (activeIndex === -1 || targetIndex === -1 || activeIndex === targetIndex) {
      return null;
    }

    return targetIndex > activeIndex ? 'next' : 'prev';
  }

  function getPostNavigationContext(link) {
    if (!link) return null;

    if (link.closest('.post-nav')) {
      if (link.matches('.next, [data-action="next"]')) {
        return { source: 'post-nav', direction: 'next' };
      }

      if (link.matches('.prev, [data-action="prev"]')) {
        return { source: 'post-nav', direction: 'prev' };
      }
    }

    const sidebarDirection = getSidebarDirection(link);
    if (sidebarDirection) {
      return { source: 'sidebar', direction: sidebarDirection };
    }

    return null;
  }

  function clearSharedTransitionElements(root) {
    if (!root || !root.querySelectorAll) return;

    root.querySelectorAll('[data-shared-transition-active]').forEach((element) => {
      element.style.removeProperty('view-transition-name');
      element.removeAttribute('data-shared-transition-active');
    });
  }

  function findIncomingSharedTransitionTarget(newDoc, context) {
    const selector = context && context.incomingShellSelector ? context.incomingShellSelector : '[data-shared-transition-path]';
    const candidates = newDoc.querySelectorAll(selector);

    for (const candidate of candidates) {
      if (!context || !context.path || !candidate.hasAttribute('data-shared-transition-path')) {
        return candidate;
      }

      if (candidate.getAttribute('data-sidebar-transition-path') === context.path ||
          candidate.getAttribute('data-shared-transition-path') === context.path) {
        return candidate;
      }
    }

    return null;
  }

  function activateSharedTransitionContext(context) {
    if (!context) return;

    setSharedTransitionState(true);
    markSharedTransitionElement(context.shell, context.names.shell);
    markSharedTransitionElement(context.title, context.names.title);
  }

  function prepareIncomingSharedTransition(newDoc, context) {
    if (!context) return;

    const target = findIncomingSharedTransitionTarget(newDoc, context);
    if (!target) return;

    markSharedTransitionElement(target, context.names.shell);
    markSharedTransitionElement(target.querySelector(context.incomingTitleSelector || '[data-shared-transition-title]'), context.names.title);
  }

  /**
   * Fetch document for transition
   */
  function buildFetchOptions(extraOptions) {
    const options = Object.assign({}, extraOptions || {});
    options.headers = Object.assign({
      'Accept': 'text/html,application/xhtml+xml',
    }, options.headers || {});

    return options;
  }

  async function fetchDocument(url, options) {
    const metrics = options && options.metrics;
    const fetchOptions = options && options.fetchOptions;
    const fetchStartedAt = getNow();

    // Fetch new page content.
    // Allow default caching so prefetching and back/forward navigation is instantaneous.
    const response = await fetch(url, buildFetchOptions(fetchOptions));

    if (metrics) {
      metrics.fetchMs += getNow() - fetchStartedAt;
    }

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const contentType = (response.headers.get('content-type') || '').toLowerCase();
    if (!contentType.includes('text/html') && !contentType.includes('application/xhtml+xml')) {
      throw new Error(`Non-HTML navigation (content-type: ${contentType || 'unknown'})`);
    }

    const html = await response.text();

    // Parse new document
    const parseStartedAt = getNow();
    const parser = new DOMParser();
    const newDoc = parser.parseFromString(html, 'text/html');

    if (metrics) {
      metrics.parseMs += getNow() - parseStartedAt;
    }

    return newDoc;
  }

  async function resolveDocument(url, metrics) {
    const prefetched = prefetchedDocuments.get(url);
    if (prefetched) {
      prefetchedDocuments.delete(url);

      if (getNow() - prefetched.prefetchedAt <= PREFETCH_FRESHNESS_MS) {
        // Failed speculative requests are cache misses, not navigation failures.
        const document = await prefetched.promise;
        if (document) {
          metrics.prefetched = true;
          return document;
        }
        return fetchDocument(url, {
          metrics: metrics,
          fetchOptions: { cache: 'no-cache' },
        });
      }

      // Abort any old in-flight speculative request and ask the HTTP cache to
      // validate with ETag/Last-Modified when it retained those validators.
      prefetched.controller.abort();
      return fetchDocument(url, {
        metrics: metrics,
        fetchOptions: { cache: 'no-cache' },
      });
    }

    return fetchDocument(url, { metrics: metrics });
  }

  function syncBodyAttributes(newDoc) {
    if (!newDoc || !newDoc.body) return;

    const currentAttributes = Array.from(document.body.attributes);
    for (const attribute of currentAttributes) {
      if (!newDoc.body.hasAttribute(attribute.name)) {
        document.body.removeAttribute(attribute.name);
      }
    }

    Array.from(newDoc.body.attributes).forEach((attribute) => {
      document.body.setAttribute(attribute.name, attribute.value);
    });
  }

  function syncDocumentElementAttributes(newDoc) {
    if (!newDoc || !newDoc.documentElement) return;

    const preservedAttributes = new Map();
    const preservedClasses = new Set();

    Array.from(document.documentElement.attributes).forEach((attribute) => {
      if (RUNTIME_HTML_ATTRIBUTES.has(attribute.name) ||
          RUNTIME_HTML_ATTRIBUTE_PREFIXES.some((prefix) => attribute.name.startsWith(prefix))) {
        preservedAttributes.set(attribute.name, attribute.value);
      }
    });

    // A visitor-picked font follows navigation; otherwise each page keeps the
    // fontpack its author assigned.
    let pickedFontpack = null;
    try { pickedFontpack = localStorage.getItem('theme-fontpack'); } catch (_) { /* ignore */ }
    if (pickedFontpack && document.documentElement.hasAttribute('data-fontpack')) {
      preservedAttributes.set('data-fontpack', document.documentElement.getAttribute('data-fontpack'));
    }

    RUNTIME_HTML_CLASS_NAMES.forEach((className) => {
      if (document.documentElement.classList.contains(className)) {
        preservedClasses.add(className);
      }
    });

    const currentAttributes = Array.from(document.documentElement.attributes);
    for (const attribute of currentAttributes) {
      if (preservedAttributes.has(attribute.name)) {
        continue;
      }
      if (!newDoc.documentElement.hasAttribute(attribute.name)) {
        document.documentElement.removeAttribute(attribute.name);
      }
    }

    Array.from(newDoc.documentElement.attributes).forEach((attribute) => {
      document.documentElement.setAttribute(attribute.name, attribute.value);
    });

    preservedAttributes.forEach((value, name) => {
      document.documentElement.setAttribute(name, value);
    });

    preservedClasses.forEach((className) => {
      document.documentElement.classList.add(className);
    });
  }

  function stylesheetKey(node) {
    if (!node) return '';
    if (node.tagName === 'LINK') {
      return ['link', node.getAttribute('rel') || '', node.getAttribute('href') || '', node.getAttribute('media') || ''].join('::');
    }
    if (node.tagName === 'STYLE') {
      return ['style', node.textContent || ''].join('::');
    }
    return '';
  }

  function headScriptKey(node) {
    if (!node || node.tagName !== 'SCRIPT') return '';
    const src = node.getAttribute('src') || '';
    if (src) {
      return ['script', node.getAttribute('type') || '', src].join('::');
    }
    return ['script', node.getAttribute('type') || '', 'inline', node.textContent || ''].join('::');
  }

  function isManagedHeadScript(node) {
    if (!node || node.tagName !== 'SCRIPT') return false;
    if (node.hasAttribute('data-markata-persist')) return false;

    const type = (node.getAttribute('type') || '').toLowerCase();
    const src = node.getAttribute('src') || '';

    if (type === 'application/ld+json') return true;
    if (src) return true;
    if (!src) return true;
    return false;
  }

  function shouldRemoveManagedHeadScript(node) {
    if (!isManagedHeadScript(node)) return false;
    return true;
  }

  /**
   * Load destination-only styles INSIDE the View Transition update callback.
   * The browser has already captured the outgoing page at that point, so
   * enabling the incoming CSS cannot flash or restyle what visitors see.
   * Wait for real CSSOM readiness before capturing the destination frame.
   * A failed/stalled stylesheet falls back to ordinary navigation.
   */
  async function prepareHeadStyles(newDoc) {
    if (!newDoc || !newDoc.head) return;

    const alreadyLoaded = new Set(
      Array.from(document.head.querySelectorAll('link[rel="stylesheet"]'))
        .filter(node => !node.hasAttribute(STAGED_STYLESHEET_ATTR))
        .map(stylesheetKey)
    );
    const staged = [];
    const loads = [];

    const nextStyles = Array.from(newDoc.head.querySelectorAll('link[rel="stylesheet"], style'));
    nextStyles.forEach((node, index) => {
      if (node.tagName !== 'LINK') return;
      const key = stylesheetKey(node);
      if (!key || alreadyLoaded.has(key)) return;

      // Print/on-demand styles must retain their intentional lazy loading.
      const media = node.getAttribute('media');
      if (media && window.matchMedia && !window.matchMedia(media).matches) return;

      const link = node.cloneNode(true);
      link.removeAttribute('onload');
      link.setAttribute(STAGED_STYLESHEET_ATTR, key);
      const loaded = new Promise((resolve, reject) => {
        link.addEventListener('load', resolve, { once: true });
        link.addEventListener('error', () => reject(new Error('Stylesheet failed: ' + node.href)), { once: true });
      });
      staged.push(link);
      loads.push(loaded);

      // A stylesheet link can lose its CSSOM when moved after load. Put
      // destination-only links into final cascade order BEFORE loading them.
      const existingFor = desired => Array.from(document.head.querySelectorAll('link[rel="stylesheet"], style'))
        .find(candidate => stylesheetKey(candidate) === stylesheetKey(desired) ||
          candidate.getAttribute(STAGED_STYLESHEET_ATTR) === stylesheetKey(desired));
      const following = nextStyles.slice(index + 1).map(existingFor).find(Boolean);
      const preceding = nextStyles.slice(0, index).reverse().map(existingFor).find(Boolean);
      if (following) {
        document.head.insertBefore(link, following);
      } else if (preceding) {
        preceding.insertAdjacentElement('afterend', link);
      } else {
        document.head.appendChild(link);
      }
    });

    if (!loads.length) return;
    let timeout;
    try {
      await Promise.race([
        Promise.all(loads),
        new Promise((_, reject) => {
          timeout = setTimeout(() => reject(new Error('Stylesheet preparation timed out')), STYLESHEET_READY_TIMEOUT_MS);
        })
      ]);
    } catch (error) {
      staged.forEach(node => node.remove());
      throw error;
    } finally {
      clearTimeout(timeout);
    }
  }

  function syncHeadStyles(newDoc) {
    if (!newDoc || !newDoc.head) return;

    const selector = 'link[rel="stylesheet"], style';
    const head = document.head;
    const currentNodes = Array.from(head.querySelectorAll(selector));
    const nextNodes = Array.from(newDoc.head.querySelectorAll(selector));
    const nextKeys = new Set(nextNodes.map(stylesheetKey).filter(Boolean));

    // Activate prepared styles before removing outgoing ones. Preserve the
    // destination cascade order within the synchronous transition callback.
    const firstCurrent = currentNodes.find(node => !node.hasAttribute(STAGED_STYLESHEET_ATTR));
    let previousInserted = null;
    nextNodes.forEach(node => {
      const key = stylesheetKey(node);
      if (!key) return;

      const existing = Array.from(head.querySelectorAll(selector)).find(candidate =>
        !candidate.hasAttribute(STAGED_STYLESHEET_ATTR) && stylesheetKey(candidate) === key
      );
      let replacement = existing;
      if (!replacement) {
        replacement = currentNodes.find(candidate =>
          candidate.getAttribute(STAGED_STYLESHEET_ATTR) === key
        );
        if (replacement) {
          if (node.hasAttribute('media')) {
            replacement.setAttribute('media', node.getAttribute('media'));
          } else {
            replacement.removeAttribute('media');
          }
          if (node.hasAttribute('onload')) {
            replacement.setAttribute('onload', node.getAttribute('onload'));
          }
          replacement.removeAttribute(STAGED_STYLESHEET_ATTR);
        } else {
          // Non-blocking media (such as print) can still load on demand.
          replacement = node.cloneNode(true);
        }
      }

      // Never detach/reinsert an existing <link>: Chromium unloads its
      // CSSOM during that move, causing a one-frame flash. Prepared links
      // were already inserted at their destination cascade position.
      if (!replacement.parentNode) {
        if (previousInserted && previousInserted.parentNode === head) {
          previousInserted.insertAdjacentElement('afterend', replacement);
        } else if (firstCurrent && firstCurrent.parentNode === head) {
          head.insertBefore(replacement, firstCurrent);
        } else {
          head.appendChild(replacement);
        }
      }
      previousInserted = replacement;
    });

    currentNodes.forEach(node => {
      if (node.hasAttribute('data-markata-persist')) return;
      const key = stylesheetKey(node);
      if (key && !nextKeys.has(key)) node.remove();
    });
  }

  /**
   * Sync page-managed <script> tags from the new document's <head>.
   *
   * Persistent loaders can opt out with data-markata-persist so they are not
   * touched across navigations. Other external head scripts are synchronized
   * with the destination document so a transitioned page has the same managed
   * head state as a full reload.
   */
  function syncHeadScripts(newDoc) {
    if (!newDoc || !newDoc.head) return;

    const head = document.head;
    const currentNodes = Array.from(head.querySelectorAll('script')).filter(isManagedHeadScript);
    const nextNodes = Array.from(newDoc.head.querySelectorAll('script')).filter(isManagedHeadScript);
    const nextKeys = new Set(nextNodes.map(headScriptKey).filter(Boolean));

    currentNodes.forEach((node) => {
      const key = headScriptKey(node);
      if (key && !nextKeys.has(key) && shouldRemoveManagedHeadScript(node)) {
        node.remove();
      }
    });

    const existingKeys = new Set(
      Array.from(head.querySelectorAll('script')).map(headScriptKey).filter(Boolean)
    );

    nextNodes.forEach((node) => {
      const key = headScriptKey(node);
      if (!key || existingKeys.has(key)) return;

      const fresh = document.createElement('script');
      Array.from(node.attributes).forEach((attr) => {
        fresh.setAttribute(attr.name, attr.value);
      });
      if (!node.hasAttribute('src')) {
        fresh.textContent = node.textContent || '';
      }
      head.appendChild(fresh);
      existingKeys.add(key);
    });
  }

  function replaceElementContents(selector, newDoc) {
    const currentElement = document.querySelector(selector);
    const nextElement = newDoc.querySelector(selector);

    if (!currentElement || !nextElement) {
      return false;
    }

    currentElement.innerHTML = nextElement.innerHTML;
    return true;
  }

  function replaceElement(selector, newDoc) {
    const currentElement = document.querySelector(selector);
    const nextElement = newDoc.querySelector(selector);

    if (!currentElement || !nextElement) {
      return false;
    }

    currentElement.replaceWith(nextElement.cloneNode(true));
    return true;
  }

  function updateLayoutRegions(newDoc) {
    // Preserve feed sidebar scroll position across DOM swap so
    // the list doesn't jump to top before scrollIntoView runs.
    var feedList = document.querySelector('#feed-nav-collapsible');
    var savedScrollTop = feedList ? feedList.scrollTop : 0;

    // Preserve docked sidebar pinned state across DOM swap. Overlay drawers
    // (<= 1200px) close on navigation.
    var dockedSidebars = !(typeof window.isSidebarOverlayMode === 'function' && window.isSidebarOverlayMode());
    var pinnedLeft = dockedSidebars && !!document.querySelector('.feed-sidebar--left.sidebar--pinned, .doc-sidebar--left.sidebar--pinned, .content-sidebar--left.sidebar--pinned');
    var pinnedRight = dockedSidebars && !!document.querySelector('.feed-sidebar--right.sidebar--pinned, .doc-sidebar--right.sidebar--pinned, .content-sidebar--right.sidebar--pinned');

    // Replace the full wrapper so page-specific layout classes do not leak
    // across feed -> post and post -> feed navigations.
    const replacedPage = replaceElement('#view-transition-page', newDoc);
    if (!replacedPage) {
      document.body.innerHTML = newDoc.body.innerHTML;
      return false;
    }

    // Mark new sidebars as hydrating so they stay hidden until scripts reinit
    document.querySelectorAll('.feed-sidebar, .doc-sidebar, .content-sidebar').forEach(function(el) {
      el.setAttribute('data-sidebar-hydrating', '');
    });

    // Restore sidebar pinned state on new elements
    if (pinnedLeft) {
      document.querySelectorAll('.feed-sidebar--left, .doc-sidebar--left, .content-sidebar--left').forEach(function(el) {
        el.classList.add('sidebar--pinned');
      });
    }
    if (pinnedRight) {
      document.querySelectorAll('.feed-sidebar--right, .doc-sidebar--right, .content-sidebar--right').forEach(function(el) {
        el.classList.add('sidebar--pinned');
      });
    }

    // Restore feed sidebar scroll position after content swap
    var newFeedList = document.querySelector('#feed-nav-collapsible');
    if (newFeedList && savedScrollTop > 0) {
      newFeedList.scrollTop = savedScrollTop;
    }

    replaceElement('#view-transition-progress', newDoc);
    replaceElementContents('.site-footer', newDoc);
    return true;
  }

  /**
   * Update the current document with content from new document
   */
  function updateDocument(newDoc, metrics, options) {
    const updateOptions = Object.assign({ reinitialize: true, hydrateCritical: false }, options || {});
    const swapStartedAt = getNow();

    syncDocumentElementAttributes(newDoc);
    syncBodyAttributes(newDoc);
    syncHeadStyles(newDoc);
    syncHeadScripts(newDoc);

    // Update title
    document.title = newDoc.title;

    // Update meta tags (description, og tags, etc.)
    if (config.updateMeta) {
      updateMetaTags(newDoc);
    }

    // Replace the main layout region instead of the entire body when possible.
    updateLayoutRegions(newDoc);

    if (updateOptions.hydrateCritical) {
      hydrateCriticalLayoutScripts();
    }

    if (metrics) {
      metrics.swapMs += getNow() - swapStartedAt;
    }

    if (updateOptions.reinitialize) {
      reinitializeDocument(metrics, undefined, newDoc);
    }

    // Scroll to top (or to hash if present)
    if (config.scrollToTop) {
      if (window.location.hash) {
        // New page: land on the anchor directly. getElementById avoids
        // selector errors for ids such as "2024-notes".
        let targetId = window.location.hash.slice(1);
        try {
          targetId = decodeURIComponent(targetId);
        } catch (_) {
          // keep the raw fragment
        }
        const target = document.getElementById(targetId);
        if (target) {
          target.scrollIntoView({ block: 'start', behavior: 'instant' });
        } else {
          window.scrollTo(0, 0);
        }
      } else {
        window.scrollTo(0, 0);
      }
    }
  }

  function reinitializeDocument(metrics, options, newDoc) {
    const reinitOptions = Object.assign({ skipCritical: false }, options || {});
    const reinitStartedAt = getNow();
    reexecuteInlineModuleScripts();
    reinitializeScripts(reinitOptions, newDoc);

    if (metrics) {
      metrics.reinitMs += getNow() - reinitStartedAt;
    }
  }

  function reexecuteInlineModuleScripts() {
    // Replaced fragments do not run inline module scripts, so evaluate them
    // after the transition snapshot instead of during the critical swap.
    document.body.querySelectorAll('script[type="module"]').forEach(old => {
      const fresh = document.createElement('script');
      fresh.type = 'module';
      fresh.textContent = old.textContent;
      old.replaceWith(fresh);
    });
  }

  function callOptionalInit(name, fn) {
    if (typeof fn !== 'function') return;
    try {
      fn();
    } catch (error) {
      if (config.debug) console.warn('View transition init failed:', name, error);
    }
  }

  function hydrateCriticalLayoutScripts() {
    // The destination snapshot must contain hydrated, visible sidebars.
    // Revealing them only after transition.finished flashes empty drawers.
    callOptionalInit('initFeedCycling', window.initFeedCycling);
    document.querySelectorAll('[data-sidebar-hydrating]').forEach(el => {
      el.removeAttribute('data-sidebar-hydrating');
    });
  }

  /**
   * Update meta tags from new document
   */
  function updateMetaTags(newDoc) {
    const metaSelectors = [
      'meta[name="description"]',
      'meta[property^="og:"]',
      'meta[name^="twitter:"]',
      'link[rel="canonical"]'
    ];

    metaSelectors.forEach(selector => {
      const oldMeta = document.querySelector(selector);
      const newMeta = newDoc.querySelector(selector);

      if (newMeta) {
        if (oldMeta) {
          oldMeta.replaceWith(newMeta.cloneNode(true));
        } else {
          document.head.appendChild(newMeta.cloneNode(true));
        }
      } else if (oldMeta) {
        oldMeta.remove();
      }
    });
  }

  function preserveFeedContext(rawUrl, triggerElement) {
    try {
      const currentURL = new URL(window.location.href);
      const targetURL = new URL(rawUrl, window.location.href);
      const currentFeed = currentURL.searchParams.get('feed');

      if (!currentFeed) return targetURL.href;
      if (targetURL.searchParams.get('feed')) return targetURL.href;
      if (!triggerElement) return targetURL.href;

      if (triggerElement.closest('.feed-sidebar, .post-nav')) {
        targetURL.searchParams.set('feed', currentFeed);
      }

      return targetURL.href;
    } catch (_) {
      return rawUrl;
    }
  }

  /**
   * Re-initialize scripts after content replacement
   */
  function reinitializeScripts(options, newDoc) {
    const reinitOptions = Object.assign({ skipCritical: false }, options || {});

    // Dispatch custom event for other scripts to listen to
    window.dispatchEvent(new CustomEvent('view-transition-complete'));

    // Pagefind Search (navbar) - needs re-init after DOM swap
    callOptionalInit('initPagefindSearch', window.initPagefindSearch);
    callOptionalInit('initBleveSearch', window.initBleveSearch);

    // Re-initialize feed cycling first so sidebar chrome is hydrated before
    // other scripts inspect feed-aware DOM state.
    if (!reinitOptions.skipCritical) {
      callOptionalInit('initFeedCycling', window.initFeedCycling);
    }

    // If feed cycling script was never loaded (e.g. navigating from a page
    // without a feed sidebar to one that has it), load it dynamically.
    if (typeof window.initFeedCycling !== 'function' && document.getElementById('feed-sidebar-data')) {
      var existing = document.querySelector('script[src*="feed-cycling"]');
      if (!existing) {
        var searchDoc = newDoc || document;
        var inlines = searchDoc.querySelectorAll('script:not([src])');
        for (var si = 0; si < inlines.length; si++) {
          var m = inlines[si].textContent.match(/['"]([^'"]*feed-cycling[^'"]*\.js)['"]/);
          if (m) {
            var s = document.createElement('script');
            s.src = m[1];
            s.defer = true;
            document.body.appendChild(s);
            break;
          }
        }
      }
    }

    // Re-initialize common scripts if they exist
    callOptionalInit('initScrollSpy', window.initScrollSpy);
    callOptionalInit('initTooltips', window.initTooltips);
    callOptionalInit('initMentionCards', window.initMentionCards);
    callOptionalInit('initPagination', window.initPagination);
    callOptionalInit('initNavigationShortcuts', window.initNavigationShortcuts);

    // Re-scroll feed sidebar active item into view (inline scripts don't re-run after DOM swap)
    callOptionalInit('initFeedSidebarScroll', window.initFeedSidebarScroll);

    // Re-bind feed sidebar collapse toggle (tablet/mobile)
    callOptionalInit('initSidebarToggle', window.initSidebarToggle);

    // Re-bind sidebar drawer handles and re-sync their open state/labels
    callOptionalInit('initSidebarPinButtons', window.initSidebarPinButtons);
    callOptionalInit('initSidebarTop', window.initSidebarTop);
    callOptionalInit('restoreSidebarState', window.restoreSidebarState);

    // Close hamburger menu after navigation (header is outside #view-transition-page so it persists)
    var openHamburger = document.querySelector('.hamburger-toggle--open');
    if (openHamburger) {
      openHamburger.classList.remove('hamburger-toggle--open');
      openHamburger.setAttribute('aria-expanded', 'false');
      var navGroup = document.querySelector('.mobile-nav-group--open');
      if (navGroup) navGroup.classList.remove('mobile-nav-group--open');
    }

    // Re-initialize mermaid diagrams (module script won't re-execute after DOM swap)
    callOptionalInit('initMermaid', window.initMermaid);

    // Re-attach event listeners
    initNavigationInterceptor();

    // Reveal sidebars now that all scripts have rehydrated
    document.querySelectorAll('[data-sidebar-hydrating]').forEach(function(el) {
      el.removeAttribute('data-sidebar-hydrating');
    });
  }

  async function navigateWithViewTransition(url, options) {
    const navOptions = Object.assign({
      source: 'script',
      pushState: true,
      triggerElement: null,
      bypassTransition: false,
    }, options || {});

    if (navigationInFlight) {
      if (config.debug) console.log('Skipping navigation while transition is active:', url);
      return false;
    }

    let targetURL;
    try {
      targetURL = new URL(preserveFeedContext(url, navOptions.triggerElement), window.location.href);
    } catch (_) {
      window.location.href = url;
      return false;
    }

    if (!shouldTransitionToURL(targetURL)) {
      window.location.href = targetURL.href;
      return false;
    }

    if (navOptions.bypassTransition) {
      const metrics = createNavigationMetrics(targetURL.href, navOptions.source);
      const newDoc = await resolveDocument(targetURL.href, metrics);
      // This fast path does not capture an old snapshot. If new render-
      // blocking CSS is required, prefer an ordinary page load to a FOUC.
      const currentStyles = new Set(Array.from(document.head.querySelectorAll('link[rel="stylesheet"]')).map(stylesheetKey));
      const hasNewStyles = Array.from(newDoc.head.querySelectorAll('link[rel="stylesheet"]'))
        .some(node => (!node.media || !window.matchMedia || window.matchMedia(node.media).matches) &&
          !currentStyles.has(stylesheetKey(node)));
      if (hasNewStyles) {
        window.location.href = targetURL.href;
        return false;
      }
      updateDocument(newDoc, metrics);

      if (navOptions.pushState) {
        history.pushState(null, '', targetURL.href);
      }
      markRenderedLocation();

      finalizeNavigationMetrics(metrics);
      return true;
    }

    navigationInFlight = true;

    const metrics = createNavigationMetrics(targetURL.href, navOptions.source);
    try {
      const newDoc = await resolveDocument(targetURL.href, metrics);
      const sharedContext = getSidebarSharedTransitionContext(navOptions.triggerElement, targetURL) || getSharedTransitionContext(navOptions.triggerElement, targetURL);
      const postNavigationContext = getPostNavigationContext(navOptions.triggerElement);

      if (sharedContext) {
        activateSharedTransitionContext(sharedContext);
      }

      if (postNavigationContext) {
        setPostNavigationState(postNavigationContext);
      }

      const transition = document.startViewTransition(async () => {
        await prepareHeadStyles(newDoc);
        prepareIncomingSharedTransition(newDoc, sharedContext);
        updateDocument(newDoc, metrics, { reinitialize: false, hydrateCritical: true });
        if (navOptions.pushState) {
          history.pushState(null, '', targetURL.href);
        }
        markRenderedLocation();
      });

      await transition.finished;
      reinitializeDocument(metrics, { skipCritical: true }, newDoc);
      finalizeNavigationMetrics(metrics);
      return true;
    } finally {
      clearSharedTransitionElements(document);
      setSharedTransitionState(false);
      clearPostNavigationState();
      navigationInFlight = false;
    }
  }

  /**
   * Handle navigation click with view transition
   */
  async function handleNavigationClick(event) {
    // Respect other handlers and normal browser behaviors
    if (event.defaultPrevented) return;
    if (event.button !== 0) return; // Only left clicks
    if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;

    const link = event.target.closest('a');

    if (!shouldTransition(link)) return;

    // Prevent default navigation
    event.preventDefault();

    const url = preserveFeedContext(link.href, link);

    if (config.debug) console.log('Starting view transition to:', url);

    try {
      await navigateWithViewTransition(url, {
        source: 'click',
        pushState: true,
        triggerElement: link,
      });

      if (config.debug) console.log('View transition completed');
    } catch (error) {
      console.error('View transition failed:', error);
      clearSharedTransitionElements(document);
      setSharedTransitionState(false);
      clearPostNavigationState();
      // Fallback to normal navigation
      window.location.href = url;
    }
  }

  /**
   * Handle back/forward button navigation
   */
  async function handlePopState(event) {
    if (config.debug) console.log('Handling popstate to:', window.location.href);

    // Same document, different fragment: the browser already scrolled.
    if (getNavigationKey(window.location) === renderedNavigationKey) {
      return;
    }

    // If we land on a special route, force a full reload so its scripts run.
    try {
      const currentURL = new URL(window.location.href);
      if (shouldSkipSpecialRoutes(currentURL)) {
        window.location.reload();
        return;
      }
    } catch (_) {
      // ignore
    }

    try {
      const url = window.location.href;
      const metrics = createNavigationMetrics(url, 'popstate');
      const newDoc = await resolveDocument(url, metrics);

      const transition = document.startViewTransition(async () => {
        await prepareHeadStyles(newDoc);
        updateDocument(newDoc, metrics, { reinitialize: false, hydrateCritical: true });
        markRenderedLocation();
      });

      await transition.finished;
      reinitializeDocument(metrics, { skipCritical: true }, newDoc);

      finalizeNavigationMetrics(metrics);
    } catch (error) {
      console.error('View transition failed on popstate:', error);
      // Fallback to reload
      window.location.reload();
    }
  }

  /**
   * Prefetch a URL to make subsequent navigation instant
   */
  function evictStalePrefetches(now) {
    for (const [url, prefetched] of prefetchedDocuments) {
      if (now - prefetched.prefetchedAt <= PREFETCH_FRESHNESS_MS) continue;

      prefetched.controller.abort();
      prefetchedDocuments.delete(url);
    }
  }

  function prefetchUrl(url) {
    // Only prefetch if we're not heavily resource constrained
    if (navigator.connection && navigator.connection.saveData) {
      return;
    }

    if (!url || prefetchedDocuments.has(url)) return;

    const now = getNow();
    evictStalePrefetches(now);
    if (url === lastPrefetchUrl && now - lastPrefetchAt < PREFETCH_DEBOUNCE_MS) {
      return;
    }
    lastPrefetchUrl = url;
    lastPrefetchAt = now;

    if (prefetchedDocuments.size >= MAX_PREFETCHED_DOCUMENTS) {
      if (config.debug) console.log('[view-transitions] prefetch cache full, skipping:', url);
      return;
    }

    if (config.debug) console.log('Prefetching URL:', url);

    const controller = new AbortController();
    const prefetched = {
      prefetchedAt: now,
      controller: controller,
      promise: null,
    };
    prefetched.promise = fetchDocument(url, {
      fetchOptions: { priority: 'low', signal: controller.signal },
    }).then((document) => {
      // Age the cache from when the document is ready to consume, not from
      // when a slow speculative request started.
      prefetched.prefetchedAt = getNow();
      return document;
    }).catch((error) => {
      if (prefetchedDocuments.get(url) === prefetched) {
        prefetchedDocuments.delete(url);
      }

      if (controller.signal.aborted) {
        return null;
      }

      if (config.debug) {
        console.warn('[view-transitions] prefetch failed', url, error);
      }

      // A hover is not a navigation: never leave an unobserved rejection.
      // Clicking this link will retry through the ordinary navigation path.
      return null;
    });

    prefetchedDocuments.set(url, prefetched);
  }

  /**
   * Handle hovering over links to prefetch them
   */
  function handleLinkHover(event) {
    const link = event.target.closest('a');
    if (!link || !shouldTransition(link)) return;

    prefetchUrl(link.href);
  }

  /**
   * Initialize click interceptor for navigation links
   */
  function initNavigationInterceptor() {
    // Use event delegation on document for better performance
    document.removeEventListener('click', handleNavigationClick);
    document.addEventListener('click', handleNavigationClick);

    // Add prefetch on hover/focus for instant transitions
    document.removeEventListener('mouseover', handleLinkHover);
    document.addEventListener('mouseover', handleLinkHover, { passive: true });

    document.removeEventListener('focusin', handleLinkHover);
    document.addEventListener('focusin', handleLinkHover, { passive: true });
  }

  /**
   * Initialize view transitions on page load
   */
  function init() {
    if (config.debug) console.log('View Transitions API initialized');

    // Intercept navigation clicks
    initNavigationInterceptor();

    // Handle browser back/forward buttons
    window.addEventListener('popstate', handlePopState);

    // Expose reinitialization function for other scripts
    window.reinitViewTransitions = initNavigationInterceptor;
    window.prefetchViewTransitionUrl = prefetchUrl;
    window.navigateWithViewTransition = navigateWithViewTransition;

    // Expose config for runtime inspection
    window.VIEW_TRANSITIONS_CONFIG = config;
  }

  // ── Feed Sidebar: scroll active item into view ──
  // Defined outside init() so it's available even if view transitions are disabled.
  // Uses manual scrollTop math so that only the sidebar container scrolls --
  // scrollIntoView() would scroll every ancestor including the page viewport.
  window.initFeedSidebarScroll = function() {
    var active = document.querySelector('.feed-nav-item--active');
    if (!active) return;
    var drawer = active.closest('.feed-sidebar');
    if (drawer && !drawer.classList.contains('sidebar--pinned')) return;
    var container = document.getElementById('feed-nav-collapsible');
    if (!container) return;
    var activeTop = active.offsetTop - container.offsetTop;
    var activeHeight = active.offsetHeight;
    var containerHeight = container.clientHeight;
    // Center the active item within the scrollable container
    container.scrollTop = activeTop - (containerHeight / 2) + (activeHeight / 2);
  };

  // ── Feed Sidebar: mobile/tablet collapse toggle ──
  // The feed-nav-header has role="button" and aria-controls="feed-nav-collapsible".
  // On narrow viewports (≤1200px) the CSS hides #feed-nav-collapsible and shows
  // it when .feed-sidebar has the --expanded modifier.  This wires up the click.
  window.initSidebarToggle = function() {
    var header = document.querySelector('.feed-nav-header[aria-controls="feed-nav-collapsible"]');
    if (!header) return;
    // Avoid double-binding
    if (header._sidebarToggleBound) return;
    header._sidebarToggleBound = true;

    function toggle() {
      var sidebar = header.closest('.feed-sidebar');
      if (!sidebar) return;
      var expanded = sidebar.classList.toggle('feed-sidebar--expanded');
      header.setAttribute('aria-expanded', String(expanded));
    }

    header.addEventListener('click', toggle);
    header.addEventListener('keydown', function(e) {
      if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        toggle();
      }
    });
  };

  // ── Sidebar drawers ──
  // Drawers open only on explicit action: the edge handle (.sidebar-toggle)
  // or the b / B shortcuts. Docked mode (>= 1201px) pushes content aside and
  // remembers open state per side in localStorage. Overlay mode (<= 1200px)
  // slides a single drawer over a scrim; Escape, the scrim, or following a
  // link closes it, and nothing is persisted.
  var SIDEBAR_STORAGE_PREFIX = 'markata-sidebar:';
  var DRAWER_SELECTOR = '.feed-sidebar, .doc-sidebar, .content-sidebar';
  var OVERLAY_DRAWER_SELECTOR = '.feed-sidebar, .doc-sidebar';
  var sidebarOverlayQuery = window.matchMedia ? window.matchMedia('(max-width: 1200px)') : null;

  function isSidebarOverlayMode() {
    return !!(sidebarOverlayQuery && sidebarOverlayQuery.matches);
  }

  function isActiveDrawer(sidebar) {
    // Content sidebars flow below the article in overlay mode.
    return !(isSidebarOverlayMode() && sidebar.classList.contains('content-sidebar'));
  }

  function sidebarSide(el) {
    return /--left\b/.test(el.className) ? 'left' : 'right';
  }

  function sidebarStorageKey(side) {
    return SIDEBAR_STORAGE_PREFIX + side;
  }

  function readSidebarState(side) {
    try {
      return localStorage.getItem(sidebarStorageKey(side));
    } catch (e) {
      return null;
    }
  }

  function writeSidebarState(side, open) {
    try {
      localStorage.setItem(sidebarStorageKey(side), open ? 'open' : 'closed');
    } catch (e) {
      /* storage unavailable: state is per-page only */
    }
  }

  function syncSidebarToggle(sidebar) {
    var btn = sidebar.querySelector(':scope > .sidebar-toggle');
    if (!btn) return;
    var open = sidebar.classList.contains('sidebar--pinned');
    var label = btn.getAttribute('data-sidebar-label') || 'sidebar';
    var hotkey = btn.getAttribute('data-sidebar-hotkey');
    var action = (open ? 'Close ' : 'Open ') + label + (label === 'sidebar' ? '' : ' sidebar');
    btn.setAttribute('aria-expanded', String(open));
    btn.setAttribute('aria-label', action);
    btn.title = hotkey ? action + ' (' + hotkey + ')' : action;
  }

  function withoutSidebarMotion(sidebar, fn) {
    sidebar.classList.add('sidebar--no-motion');
    fn();
    // Two frames so the class change lands before transitions resume.
    requestAnimationFrame(function() {
      requestAnimationFrame(function() {
        sidebar.classList.remove('sidebar--no-motion');
      });
    });
  }

  function focusDrawerStart(sidebar) {
    var target = sidebar.querySelector('.toc-link--active, .feed-nav-item--active a') ||
      sidebar.querySelector('.toc a[href], .feed-nav a[href], a[href], button:not(.sidebar-toggle)');
    if (target && typeof target.focus === 'function') {
      target.focus({ preventScroll: true });
    }
  }

  function setSidebarOpen(sidebar, open, persist) {
    var overlay = isSidebarOverlayMode();
    if (open && overlay) {
      document.querySelectorAll(OVERLAY_DRAWER_SELECTOR).forEach(function(other) {
        if (other !== sidebar && other.classList.contains('sidebar--pinned')) {
          setSidebarOpen(other, false, false);
        }
      });
    }
    sidebar.classList.toggle('sidebar--pinned', open);
    syncSidebarToggle(sidebar);
    if (open) {
      updateSidebarTop();
      window.initFeedSidebarScroll();
    }
    if (open && overlay) {
      // Wait for visibility to flip before moving focus into the drawer.
      requestAnimationFrame(function() { focusDrawerStart(sidebar); });
    }
    if (!open) {
      var active = document.activeElement;
      var btn = sidebar.querySelector(':scope > .sidebar-toggle');
      if (active && sidebar.contains(active) && active !== btn) {
        if (overlay && btn) {
          btn.focus({ preventScroll: true });
        } else if (typeof active.blur === 'function') {
          active.blur();
        }
      }
    }
    if (persist && !overlay) writeSidebarState(sidebarSide(sidebar), open);
  }

  function closeOverlayDrawers() {
    var closed = null;
    document.querySelectorAll(OVERLAY_DRAWER_SELECTOR).forEach(function(sidebar) {
      if (sidebar.classList.contains('sidebar--pinned')) {
        setSidebarOpen(sidebar, false, false);
        closed = sidebar;
      }
    });
    return closed;
  }

  function openOverlayDrawer() {
    return document.querySelector('.feed-sidebar.sidebar--pinned, .doc-sidebar.sidebar--pinned');
  }

  // Keep fixed drawers flush with the bottom edge of the header while it is
  // on screen, then let them grow to the top once the header scrolls away.
  // Overlay drawers fill the viewport and need no header geometry. Scope the
  // offset to drawers so updates do not invalidate the article subtree.
  var sidebarTopFrame = null;
  var sidebarTopValue = null;
  function updateSidebarTop() {
    sidebarTopFrame = null;
    if (isSidebarOverlayMode()) return;
    var drawers = document.querySelectorAll(DRAWER_SELECTOR);
    if (!drawers.length) return;
    var header = document.querySelector('.site-header');
    var bottom = header ? Math.max(0, Math.round(header.getBoundingClientRect().bottom)) : 0;
    if (bottom === sidebarTopValue) return;
    sidebarTopValue = bottom;
    drawers.forEach(function(drawer) {
      drawer.style.setProperty('--sidebar-top', bottom + 'px');
    });
  }
  function scheduleSidebarTop() {
    if (sidebarTopFrame === null) {
      sidebarTopFrame = requestAnimationFrame(updateSidebarTop);
    }
  }
  window.initSidebarTop = function() {
    sidebarTopValue = null;
    if (window._sidebarTopBound) {
      updateSidebarTop();
      return;
    }
    window._sidebarTopBound = true;
    window.addEventListener('scroll', scheduleSidebarTop, { passive: true });
    window.addEventListener('resize', scheduleSidebarTop);
    updateSidebarTop();
  };

  // Restore persisted drawer state without animating the initial slide.
  // Overlay drawers always start closed.
  window.restoreSidebarState = function() {
    var overlay = isSidebarOverlayMode();
    document.querySelectorAll(DRAWER_SELECTOR).forEach(function(sidebar) {
      var wantOpen = !overlay && readSidebarState(sidebarSide(sidebar)) === 'open';
      if (wantOpen !== sidebar.classList.contains('sidebar--pinned')) {
        withoutSidebarMotion(sidebar, function() {
          sidebar.classList.toggle('sidebar--pinned', wantOpen);
        });
      }
      syncSidebarToggle(sidebar);
    });
    updateSidebarTop();
    window.initFeedSidebarScroll();
  };

  window.initSidebarPinButtons = function() {
    document.querySelectorAll('.sidebar-toggle').forEach(function(btn) {
      if (btn._sidebarToggleBound) return;
      btn._sidebarToggleBound = true;
      btn.addEventListener('click', function(e) {
        e.stopPropagation();
        var sidebar = btn.closest(DRAWER_SELECTOR);
        if (!sidebar) return;
        setSidebarOpen(sidebar, !sidebar.classList.contains('sidebar--pinned'), true);
      });
    });

    if (window._sidebarDocumentBound) return;
    window._sidebarDocumentBound = true;

    // Scrim tap: any click outside the open overlay drawer closes it.
    document.addEventListener('click', function(e) {
      if (!isSidebarOverlayMode() || !openOverlayDrawer()) return;
      var target = e.target;
      var inDrawer = target && target.closest && target.closest(OVERLAY_DRAWER_SELECTOR);
      if (!inDrawer) {
        e.preventDefault();
        e.stopPropagation();
        closeOverlayDrawers();
        return;
      }
      // Following a link (TOC jump or feed navigation) dismisses the overlay.
      if (target.closest('a[href]')) {
        setSidebarOpen(inDrawer, false, false);
      }
    }, true);

    document.addEventListener('keydown', function(e) {
      if (e.key !== 'Escape' || !isSidebarOverlayMode()) return;
      var closed = closeOverlayDrawers();
      if (closed) {
        e.preventDefault();
        var btn = closed.querySelector(':scope > .sidebar-toggle');
        if (btn) btn.focus({ preventScroll: true });
      }
    });

    if (sidebarOverlayQuery) {
      var onModeChange = function() {
        if (isSidebarOverlayMode()) {
          document.querySelectorAll(DRAWER_SELECTOR).forEach(function(sidebar) {
            if (sidebar.classList.contains('sidebar--pinned')) {
              withoutSidebarMotion(sidebar, function() {
                sidebar.classList.remove('sidebar--pinned');
              });
              syncSidebarToggle(sidebar);
            }
          });
        } else {
          window.restoreSidebarState();
        }
      };
      if (typeof sidebarOverlayQuery.addEventListener === 'function') {
        sidebarOverlayQuery.addEventListener('change', onModeChange);
      } else if (typeof sidebarOverlayQuery.addListener === 'function') {
        sidebarOverlayQuery.addListener(onModeChange);
      }
    }
  };

  // Used by updateLayoutRegions to decide whether open drawers survive SPA
  // navigation (docked only).
  window.isSidebarOverlayMode = isSidebarOverlayMode;

  // ── Keyboard shortcuts to open/close sidebars ──
  // b = toggle left drawer, B = toggle right drawer
  // Registered via shortcuts registry in navigation-shortcuts.js.
  window.toggleSidebarPinned = function(side) {
    var selector = side === 'left'
      ? '.feed-sidebar--left, .doc-sidebar--left, .content-sidebar--left'
      : '.feed-sidebar--right, .doc-sidebar--right, .content-sidebar--right';
    var drawers = Array.prototype.filter.call(document.querySelectorAll(selector), isActiveDrawer);
    if (!drawers.length) return;
    var anyOpen = drawers.some(function(el) { return el.classList.contains('sidebar--pinned'); });
    if (isSidebarOverlayMode()) {
      setSidebarOpen(drawers[0], !anyOpen, false);
      return;
    }
    drawers.forEach(function(el) {
      setSidebarOpen(el, !anyOpen, true);
    });
  };

  // Initialize when DOM is ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', function() {
      init();
      window.initFeedSidebarScroll();
      window.initSidebarToggle();
      window.initSidebarPinButtons();
      window.initSidebarTop();
      window.restoreSidebarState();
    });
  } else {
    init();
    window.initFeedSidebarScroll();
    window.initSidebarToggle();
    window.initSidebarPinButtons();
    window.initSidebarTop();
    window.restoreSidebarState();
  }
})();
