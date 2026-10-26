/**
 * Navigation Shortcuts Module for markata-go
 *
 * Registers navigation-related keyboard shortcuts with the shortcuts registry.
 * - `j` or `↓` - Next post (in feeds) / Highlight next card
 * - `k` or `↑` - Previous post (in feeds) / Highlight previous card
 * - `Enter` or `o` - Open highlighted post
 * - `Shift+O` - Open in new tab
 * - `g h` - Go to home
 * - `g s` - Focus search
 * - `s` - Toggle simple/rich feed view (feed pages) or reader mode (post pages)
 * - `[` - Previous page
 * - `]` - Next page
 * - `y y` - Copy URL to clipboard
 *
 * When feed cards are present, j/k will navigate between cards (with visual highlight).
 * Press o/Enter to open the selected card.
 */

(function() {
  'use strict';

  // Wait for registry to be available
  function waitForRegistry(callback, attempts) {
    attempts = attempts || 0;
    if (window.shortcutsRegistry) {
      callback();
    } else if (attempts < 50) {
      setTimeout(function() {
        waitForRegistry(callback, attempts + 1);
      }, 10);
    }
  }

   var state = {
     selectedCard: null,
     cards: [],
     lastKeyTime: 0,
     lastKey: null,
     keySequenceTimeout: 500, // ms
     jKeyDown: false,
     kKeyDown: false,
     navRepeatTimer: null,
     navRepeatDelay: 80, // ms between navigations when key held (faster)
     initialized: false,  // Track if already initialized
     jkRegistered: false  // Track if j/k shortcuts are registered
   };

   /**
    * Check if we're on a dedicated feed page (not a post page with an embedded feed widget).
    * Feed pages have <div class="feed"> but NOT <article class="post">.
    */
   function isDedicatedFeedPage() {
     return document.querySelector('div.feed') !== null &&
            document.querySelector('article.post') === null;
   }

   /**
    * Get all post cards on the page
    */
   function getCards() {
     return Array.from(document.querySelectorAll('.card, [data-card]'));
   }

  /**
   * Highlight a card
   */
  function highlightCard(card) {
    // Remove previous highlight
    if (state.selectedCard) {
      state.selectedCard.classList.remove('kb-highlighted');
    }

    // Highlight new card
    state.selectedCard = card;
    card.classList.add('kb-highlighted');
    card.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
  }

  /**
   * Navigate to next post
   */
  function nextPost() {
    state.cards = getCards();
    if (state.cards.length === 0) return;

    var currentIndex = state.cards.indexOf(state.selectedCard);
    var nextIndex = currentIndex + 1;

    if (nextIndex >= state.cards.length) {
      // Try to navigate to next page
      var nextBtn = document.querySelector('[data-action="next"]');
      if (nextBtn && !nextBtn.disabled) {
        nextBtn.click();
      }
      return;
    }

    highlightCard(state.cards[nextIndex]);
  }

  /**
   * Navigate to previous post
   */
  function previousPost() {
    state.cards = getCards();
    if (state.cards.length === 0) return;

    var currentIndex = state.selectedCard ? state.cards.indexOf(state.selectedCard) : -1;
    var prevIndex = currentIndex - 1;

    if (prevIndex < 0) {
      // Try to navigate to previous page
      var prevBtn = document.querySelector('[data-action="prev"]');
      if (prevBtn && !prevBtn.disabled) {
        prevBtn.click();
      }
      return;
    }

    highlightCard(state.cards[prevIndex]);
  }

  /**
   * Open highlighted post
   */
  function openPost(newTab) {
    if (!state.selectedCard) {
      // Highlight first card if none selected
      state.cards = getCards();
      if (state.cards.length > 0) {
        highlightCard(state.cards[0]);
        return;
      }
      return;
    }

    var link = state.selectedCard.querySelector('a');
    if (!link) return;

    if (newTab) {
      window.open(link.href, '_blank');
    } else {
      navigateToURL(link.href, { triggerElement: link });
    }
  }

  /**
   * Go to home page
   */
  function goHome() {
    var homeLink = document.querySelector('[data-nav="home"], a[href="/"], nav a:first-child');
    if (homeLink) {
      navigateToURL(homeLink.href, { triggerElement: homeLink });
    } else {
      navigateToURL('/');
    }
  }

  function navigateToURL(url, options) {
    if (!url) return;

    if (window.navigateWithViewTransition && typeof window.navigateWithViewTransition === 'function') {
      window.navigateWithViewTransition(url, Object.assign({ source: 'shortcut' }, options || {}));
      return;
    }

    window.location.href = url;
  }

  function getSidebarAdjacentLink(direction) {
    var activeItem = document.querySelector('.feed-nav-item--active');
    if (!activeItem) {
      var links = Array.from(document.querySelectorAll('.feed-nav-list a.feed-nav-link[href]'));
      if (!links.length) return null;
      return direction === 'next' ? links[0] : links[links.length - 1];
    }

    var sibling = direction === 'next' ? activeItem.nextElementSibling : activeItem.previousElementSibling;
    if (!sibling) return null;

    return sibling.querySelector('a.feed-nav-link[href]');
  }

  /**
   * Focus search
   */
  function focusSearch() {
    var pagefindInput = document.querySelector('.pagefind-ui__search-input');
    if (pagefindInput) {
      pagefindInput.focus();
      return;
    }

    var searchInput = document.querySelector('#pagefind-search input, #search input, [type="search"]');
    if (searchInput) {
      searchInput.focus();
    }
  }

    /**
     * Toggle between simple and rich feed views.
     * On feed pages: toggles between /feed/ and /feed/simple/.
     * On post pages: toggles reader mode (hides chrome for distraction-free reading).
     */
    function toggleFeedView() {
      // Case 1: On a dedicated feed page -- toggle between simple and rich
      if (isDedicatedFeedPage()) {
        var path = window.location.pathname;
        var feedEl = document.querySelector('div.feed');
        var isSimple = feedEl.classList.contains('feed-simple');

        if (isSimple) {
          // On simple feed -> go to rich feed: remove /simple/ from path
          var richPath = path.replace(/\/simple\/(page\/\d+\/)?$/, '/');
          // Preserve page number if present
          var pageMatch = path.match(/\/simple\/page\/(\d+)\//);
          if (pageMatch) {
            richPath = path.replace(/\/simple\/page\/(\d+)\//, '/page/' + pageMatch[1] + '/');
          }
          navigateToURL(richPath);
        } else {
          // On rich feed -> go to simple feed: insert /simple/ before any /page/ or at end
          var simplePath;
          var richPageMatch = path.match(/\/page\/(\d+)\/$/);
          if (richPageMatch) {
            simplePath = path.replace(/\/page\/(\d+)\/$/, '/simple/page/' + richPageMatch[1] + '/');
          } else {
            // Ensure trailing slash
            simplePath = path.replace(/\/?$/, '/') + 'simple/';
          }
          navigateToURL(simplePath);
        }
        return;
      }

      // Case 2: On a post page -- toggle reader mode
      if (document.querySelector('article.post')) {
        toggleReaderMode();
      }
    }

    /**
     * Toggle reader mode on post pages.
     * Strips chrome (sidebar, nav, footer) for distraction-free reading.
     * Press 's' again to restore normal view.
     * Persists via ?reader=1 URL parameter across navigation.
     */
    function toggleReaderMode() {
      var body = document.body;
      var isReaderMode = body.classList.toggle('reader-mode');

      // Update URL parameter to persist across navigation
      var url = new URL(window.location.href);
      if (isReaderMode) {
        url.searchParams.set('reader', '1');
      } else {
        url.searchParams.delete('reader');
      }
      history.replaceState(null, '', url.toString());

      // Show a brief notification
      if (isReaderMode) {
        showNotification('Reader mode -- press s to exit');
      } else {
        showNotification('Reader mode off');
      }
    }

    /**
     * Check for ?reader=1 URL param and apply reader-mode class.
     * Called on init (page load and after view transitions).
     */
    function applyReaderModeFromURL() {
      var params = new URLSearchParams(window.location.search);
      if (params.get('reader') === '1' && document.querySelector('article.post')) {
        document.body.classList.add('reader-mode');
      } else {
        document.body.classList.remove('reader-mode');
      }
    }

   /**
    * Get the active feed slug when on a post page.
    * Checks feed-cycling state first, then falls back to the sidebar title link.
    */
   function getPostPageFeedSlug() {
     // Try feed-cycling JS state (if feed cycling is active)
     try {
       var dataEl = document.getElementById('feed-sidebar-data');
       if (dataEl) {
         var data = JSON.parse(dataEl.textContent);
         // Check URL for ?feed= param (user may have cycled)
         var params = new URLSearchParams(window.location.search);
         var feedParam = params.get('feed');
         if (feedParam) return feedParam;
         // Fall back to the default feed from JSON
         if (data.feeds && data.feeds[data.currentFeedIndex]) {
           return data.feeds[data.currentFeedIndex].slug;
         }
       }
     } catch (e) { /* ignore parse errors */ }

     // Fall back to sidebar title link href
     var titleLink = document.querySelector('.feed-nav-title a[href]');
     if (titleLink) {
       // Extract slug from href like "/blog/" -> "blog"
       var href = titleLink.getAttribute('href');
       var match = href.match(/^\/([^?#]+?)\/?$/);
       if (match) return match[1];
     }

     return null;
   }

  /**
   * Navigate to next page.
   * Checks for JS pagination buttons (data-action="next") and also
   * manual/htmx pagination links (.pagination-next).
   */
  function nextPage() {
    var sidebarNextLink = getSidebarAdjacentLink('next');
    if (sidebarNextLink) {
      navigateToURL(sidebarNextLink.href, { triggerElement: sidebarNextLink, bypassTransition: true, source: 'shortcut' });
      return;
    }

    // Try JS pagination button first
    var nextBtn = document.querySelector('[data-action="next"]:not(:disabled)');
    if (nextBtn) {
      nextBtn.click();
      return;
    }
    // Try manual/htmx pagination link
    var nextLink = document.querySelector('a.pagination-next');
    if (nextLink) {
      nextLink.click();
    }
  }

  /**
   * Navigate to previous page.
   * Checks for JS pagination buttons (data-action="prev") and also
   * manual/htmx pagination links (.pagination-prev).
   */
  function previousPage() {
    var sidebarPrevLink = getSidebarAdjacentLink('prev');
    if (sidebarPrevLink) {
      navigateToURL(sidebarPrevLink.href, { triggerElement: sidebarPrevLink, bypassTransition: true, source: 'shortcut' });
      return;
    }

    // Try JS pagination button first
    var prevBtn = document.querySelector('[data-action="prev"]:not(:disabled)');
    if (prevBtn) {
      prevBtn.click();
      return;
    }
    // Try manual/htmx pagination link
    var prevLink = document.querySelector('a.pagination-prev');
    if (prevLink) {
      prevLink.click();
    }
  }

  /**
   * Copy current page URL to clipboard
   */
  function copyUrl() {
    var url = window.location.href;
    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(url).then(function() {
        showNotification('URL copied to clipboard');
      }).catch(function(err) {
        console.error('Failed to copy URL:', err);
      });
    } else {
      // Fallback for older browsers
      var textarea = document.createElement('textarea');
      textarea.value = url;
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand('copy');
      document.body.removeChild(textarea);
      showNotification('URL copied to clipboard');
    }
  }

  /**
   * Show a brief notification
   */
  function showNotification(text) {
    var notif = document.querySelector('.shortcuts-notification');
    if (!notif) {
      notif = document.createElement('div');
      notif.className = 'shortcuts-notification';
      notif.style.cssText = 'position: fixed; bottom: 20px; right: 20px; ' +
        'background: var(--color-bg); color: var(--color-text); ' +
        'padding: 12px 16px; border-radius: 4px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); ' +
        'font-size: 14px; z-index: 10000; opacity: 0; transition: opacity 0.3s;';
      document.body.appendChild(notif);
    }

    notif.textContent = text;
    notif.style.opacity = '1';

    clearTimeout(notif._timeout);
    notif._timeout = setTimeout(function() {
      notif.style.opacity = '0';
    }, 2000);
  }

  /**
   * Clean up state for re-initialization (for view transitions)
   */
  function cleanup() {
    // Clear any timers
    if (state.navRepeatTimer) {
      clearTimeout(state.navRepeatTimer);
      state.navRepeatTimer = null;
    }

    // Reset key states
    state.jKeyDown = false;
    state.kKeyDown = false;

    // Remove highlight from old card (if it still exists in DOM)
    if (state.selectedCard && state.selectedCard.parentNode) {
      state.selectedCard.classList.remove('kb-highlighted');
    }

    // Clear card references
    state.selectedCard = null;
    state.cards = [];
  }

  /**
   * Preload adjacent pages for faster pagination
   */
  function preloadAdjacentPages() {
    // Only preload if we're not heavily resource constrained
    if (navigator.connection && navigator.connection.saveData) {
      return;
    }

    if (typeof window.prefetchViewTransitionUrl !== 'function') {
      return;
    }

    var nextLink = document.querySelector('[data-action="next"], a.pagination-next');
    if (nextLink && nextLink.href) {
      window.prefetchViewTransitionUrl(nextLink.href);
    }

    var prevLink = document.querySelector('[data-action="prev"], a.pagination-prev');
    if (prevLink && prevLink.href) {
      window.prefetchViewTransitionUrl(prevLink.href);
    }
  }

  window.preloadNavigationShortcuts = preloadAdjacentPages;

  /**
   * Initialize navigation shortcuts
   */
  function init() {
    // Clean up previous state
    cleanup();

    // Apply reader mode from URL param (persists across view transitions)
    applyReaderModeFromURL();

    // Adjacent pages are fetched on link hover/focus or explicit feed cycling.
    // Initializing keyboard controls is not navigation intent.

    // Initialize cards -- but do NOT auto-highlight or scroll on page load.
    // The first card will be highlighted when the user presses j/k.
    state.cards = getCards();
    state.selectedCard = null;

    // Only register these shortcuts once (they use document-level listeners internally)
    if (!state.initialized) {
      // Enter - Open highlighted post
      window.shortcutsRegistry.register({
        key: 'Enter',
        modifiers: [],
        description: 'Open highlighted post',
        group: 'navigation',
        handler: function(e) {
          e.preventDefault();
          openPost(false);
        },
        priority: 15
      });

      // o - Open highlighted post
      window.shortcutsRegistry.register({
        key: 'o',
        modifiers: [],
        description: 'Open highlighted card',
        group: 'navigation',
        handler: function(e) {
          e.preventDefault();
          openPost(false);
        },
        priority: 15
      });

      // Shift+O - Open in new tab
      window.shortcutsRegistry.register({
        key: 'O',
        modifiers: [],
        description: 'Open highlighted card in new tab',
        group: 'navigation',
        handler: function(e) {
          e.preventDefault();
          openPost(true);
        },
        priority: 15
      });

      // [ and ] - Previous/Next page
      // Priority 55: higher than palette-switcher (50) so pagination wins.
      // Registered unconditionally -- handlers gracefully no-op when prev/next
      // elements aren't on the current page.
      window.shortcutsRegistry.register({
        key: '[',
        modifiers: [],
        description: 'Previous page',
        group: 'navigation',
        handler: function(e) {
          e.preventDefault();
          previousPage();
        },
        priority: 55
      });

      window.shortcutsRegistry.register({
        key: ']',
        modifiers: [],
        description: 'Next page',
        group: 'navigation',
        handler: function(e) {
          e.preventDefault();
          nextPage();
        },
        priority: 55
      });

      // b and B - Toggle sidebar pinned open/closed
      window.shortcutsRegistry.register({
        key: 'b',
        modifiers: [],
        description: 'Toggle left sidebar',
        group: 'navigation',
        handler: function(e) {
          e.preventDefault();
          if (window.toggleSidebarPinned) window.toggleSidebarPinned('left');
        },
        priority: 55
      });

      window.shortcutsRegistry.register({
        key: 'B',
        modifiers: [],
        description: 'Toggle right sidebar',
        group: 'navigation',
        handler: function(e) {
          e.preventDefault();
          if (window.toggleSidebarPinned) window.toggleSidebarPinned('right');
        },
        priority: 55
      });

      // J and K (Shift+j/k) - Previous/Next post in feed
      // Vim-like list navigation: j/k move within the view (scroll or card
      // select), J/K move to the adjacent item (prev/next post).
      // Same handlers as [/] -- an additional binding for muscle-memory.
      window.shortcutsRegistry.register({
        key: 'K',
        modifiers: [],
        description: 'Previous post in feed',
        group: 'navigation',
        handler: function(e) {
          e.preventDefault();
          previousPage();
        },
        priority: 55
      });

      window.shortcutsRegistry.register({
        key: 'J',
        modifiers: [],
        description: 'Next post in feed',
        group: 'navigation',
        handler: function(e) {
          e.preventDefault();
          nextPage();
        },
        priority: 55
      });

      // Handle multi-key sequences: g h, g s, and standalone s for feed toggle
      document.addEventListener('keydown', function(e) {
        if (window.shortcutsRegistry.areDisabled()) return;
        if (window.shortcutsRegistry.isInputElement(e.target)) return;

        var now = Date.now();
        var timeSinceLastKey = now - state.lastKeyTime;

        if (e.key === 'g') {
          state.lastKey = 'g';
          state.lastKeyTime = now;
        } else if (e.key === 'h' && state.lastKey === 'g' && timeSinceLastKey < state.keySequenceTimeout) {
          // g h - go to home
          e.preventDefault();
          goHome();
          state.lastKey = null;
          state.lastKeyTime = 0;
        } else if (e.key === 's' && state.lastKey === 'g' && timeSinceLastKey < state.keySequenceTimeout) {
          // g s - focus search
          e.preventDefault();
          focusSearch();
          state.lastKey = null;
          state.lastKeyTime = 0;
         } else if (e.key === 's' && state.lastKey !== 'g') {
           // Standalone s - toggle simple/rich feed view (feed pages) or reader mode (post pages)
           e.preventDefault();
           toggleFeedView();
           state.lastKey = null;
           state.lastKeyTime = 0;
        } else {
          state.lastKey = null;
          state.lastKeyTime = 0;
        }
      });

      // y y - Copy URL (special handling for repeated key)
      var yKeyTime = 0;
      document.addEventListener('keydown', function(e) {
        if (window.shortcutsRegistry.areDisabled()) return;
        if (window.shortcutsRegistry.isInputElement(e.target)) return;

        if (e.key === 'y') {
          var now = Date.now();
          var timeSinceLastY = now - yKeyTime;

          if (timeSinceLastY < state.keySequenceTimeout) {
            // y y - copy URL
            e.preventDefault();
            copyUrl();
            yKeyTime = 0;
          } else {
            yKeyTime = now;
          }
        } else {
          yKeyTime = 0;
        }
      });

      // Propagate ?reader=1 to navigation links when reader mode is active.
      // Uses capture phase so it runs before the view transition click handler.
      document.addEventListener('click', function(e) {
        if (!document.body.classList.contains('reader-mode')) return;

        var link = e.target.closest('a');
        if (!link || !link.href) return;

        try {
          var linkUrl = new URL(link.href);
          // Only modify same-origin links
          if (linkUrl.origin !== window.location.origin) return;
          // Don't add if already present
          if (linkUrl.searchParams.get('reader') === '1') return;

          linkUrl.searchParams.set('reader', '1');
          link.href = linkUrl.toString();
        } catch (ex) { /* ignore invalid URLs */ }
      }, true); // capture phase

      state.initialized = true;
    }

    // Register j/k shortcuts if we have cards (these use state.cards which gets refreshed)
    // Only register once - the handlers will use the refreshed state.cards
    if (state.cards.length > 0 && !state.jkRegistered) {
      // j - Next card in feed
      window.shortcutsRegistry.register({
        key: 'j',
        modifiers: [],
        description: 'Select next card in feed',
        group: 'navigation',
        handler: function(e) {
          // Refresh cards on each call to handle view transitions
          state.cards = getCards();
          if (state.cards.length === 0) return;

          // Initialize selection if needed
          if (!state.selectedCard || !state.selectedCard.parentNode) {
            highlightCard(state.cards[0]);
          } else {
            e.preventDefault();
            nextPost();
          }
        },
        priority: 20
      });

      // k - Previous card in feed
      window.shortcutsRegistry.register({
        key: 'k',
        modifiers: [],
        description: 'Select previous card in feed',
        group: 'navigation',
        handler: function(e) {
          // Refresh cards on each call to handle view transitions
          state.cards = getCards();
          if (state.cards.length === 0) return;

          // Initialize selection if needed
          if (!state.selectedCard || !state.selectedCard.parentNode) {
            highlightCard(state.cards[0]);
          } else {
            e.preventDefault();
            previousPost();
          }
        },
        priority: 20
      });

      // Handle held j/k for continuous navigation
      document.addEventListener('keydown', function(e) {
        if (window.shortcutsRegistry.areDisabled()) return;
        if (window.shortcutsRegistry.isInputElement(e.target)) return;

        // Refresh cards to handle view transitions
        state.cards = getCards();
        if (state.cards.length === 0) return;

        if (e.key === 'j' && !state.jKeyDown) {
          state.jKeyDown = true;
          // Start repeat timer after initial delay
          state.navRepeatTimer = setTimeout(function repeatJ() {
            if (state.jKeyDown) {
              nextPost();
              state.navRepeatTimer = setTimeout(repeatJ, state.navRepeatDelay);
            }
          }, state.navRepeatDelay);
          e.preventDefault();
        } else if (e.key === 'k' && !state.kKeyDown) {
          state.kKeyDown = true;
          // Start repeat timer after initial delay
          state.navRepeatTimer = setTimeout(function repeatK() {
            if (state.kKeyDown) {
              previousPost();
              state.navRepeatTimer = setTimeout(repeatK, state.navRepeatDelay);
            }
          }, state.navRepeatDelay);
          e.preventDefault();
        }
      });

      // Clear repeat timer when keys are released
      document.addEventListener('keyup', function(e) {
        if (e.key === 'j') {
          state.jKeyDown = false;
          if (state.navRepeatTimer && !state.kKeyDown) {
            clearTimeout(state.navRepeatTimer);
            state.navRepeatTimer = null;
          }
        } else if (e.key === 'k') {
          state.kKeyDown = false;
          if (state.navRepeatTimer && !state.jKeyDown) {
            clearTimeout(state.navRepeatTimer);
            state.navRepeatTimer = null;
          }
        }
      });

      state.jkRegistered = true;
    }
  }

  // Initialize when registry is available
  waitForRegistry(function() {
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', init);
    } else {
      init();
    }
  });

  // Expose init for view transitions
  window.initNavigationShortcuts = init;

  // Re-initialize after view transitions
  window.addEventListener('view-transition-complete', init);
})();
