/**
 * Search Shortcuts Module for markata-go
 *
 * Registers search-related keyboard shortcuts with the shortcuts registry.
 * - `/` or `Ctrl/Cmd+K` - Focus search input
 * - `?` - Show shortcuts help modal
 * - `Escape` - Close search/modals
 *
 * Also supports deep-linking into search with `?q=<query>` and bridges the
 * generated 404 recovery form into the site's real search UI.
 */

(function() {
  'use strict';

  // Wait for registry to be available
  function waitForRegistry(callback, attempts = 0) {
    if (window.shortcutsRegistry) {
      callback();
    } else if (attempts < 50) {
      setTimeout(function() {
        waitForRegistry(callback, attempts + 1);
      }, 10);
    }
  }

  /**
   * Detect if the user is on a Mac platform
   */
  function isMacPlatform() {
    if (navigator.userAgentData && navigator.userAgentData.platform) {
      return navigator.userAgentData.platform.toUpperCase().indexOf('MAC') >= 0;
    }
    return navigator.platform.toUpperCase().indexOf('MAC') >= 0;
  }

  /**
   * Focus the search input (Pagefind).
   * If pagefind hasn't loaded yet, trigger lazy-load and focus after init.
   */
  function focusSearch() {
    // Try Pagefind's input first
    var pagefindInput = document.querySelector('.pagefind-ui__search-input');
    if (pagefindInput) {
      pagefindInput.focus();
      return true;
    }

    // Fallback to any search input
    var searchInput = document.querySelector('#pagefind-search input, #search input, [type="search"]');
    if (searchInput) {
      searchInput.focus();
      return true;
    }

    // Pagefind not loaded yet -- trigger lazy-load and focus after init
    if (window.loadPagefind) {
      window.loadPagefind(function() {
        var input = document.querySelector('.pagefind-ui__search-input');
        if (input) input.focus();
      });
      return true;
    }

    return false;
  }

  /**
   * Put a query into the site's real search control and fire the same input
   * event Pagefind/bleve receives during normal typing.
   */
  function seedSearchInput(query) {
    var input = document.querySelector('.pagefind-ui__search-input, #pagefind-search input, #search input');
    if (!input) return false;

    input.value = query;
    input.dispatchEvent(new Event('input', { bubbles: true }));
    input.focus();
    return true;
  }

  /**
   * Open the site search for a query, accounting for Pagefind's lazy loader.
   * The shortcuts bundle can execute before base.html has registered
   * window.loadPagefind, so a short retry window avoids script-order races.
   */
  function openSiteSearch(query) {
    query = (query || '').trim();
    if (!query) return false;

    if (seedSearchInput(query)) {
      return true;
    }

    if (window.loadPagefind) {
      window.loadPagefind(function() {
        seedSearchInput(query);
      });
      return true;
    }

    var attempts = 0;
    var timer = window.setInterval(function() {
      attempts += 1;
      if (seedSearchInput(query)) {
        window.clearInterval(timer);
        return;
      }
      if (window.loadPagefind) {
        window.clearInterval(timer);
        window.loadPagefind(function() {
          seedSearchInput(query);
        });
        return;
      }
      if (attempts >= 50) {
        window.clearInterval(timer);
      }
    }, 20);

    return true;
  }

  /**
   * Seed the site search from a `?q=` query parameter.
   *
   * The generated 404 historically falls back to `/?q=<query>`. Treating that
   * URL as a first-class search deep link makes the fallback useful and also
   * gives sites a simple shareable search URL.
   */
  function applySearchQueryParam() {
    var params = new URLSearchParams(window.location.search);
    var query = (params.get('q') || '').trim();
    if (!query) return false;
    return openSiteSearch(query);
  }

  /**
   * Keep the generated 404 recovery form on the 404 page when the real site
   * search is available. The 404 template's own submit listener is attached at
   * the target/bubble phase; listening on document in capture phase lets us
   * route the query into Pagefind before that legacy `/?q=` redirect runs.
   */
  function handle404RecoverySubmit(event) {
    var form = event.target;
    if (!form || form.id !== 'search-form') return;
    if (!form.closest || !form.closest('.error-404')) return;

    var input = form.querySelector('#search-input, [name="q"]');
    var query = input ? input.value.trim() : '';
    if (!query) return;

    event.preventDefault();
    event.stopImmediatePropagation();
    openSiteSearch(query);
  }

  // Remember the opener so keyboard users return to their previous position.
  var shortcutsModalOpener = null;

  /**
   * Show the shortcuts help modal
   */
  function showShortcutsModal() {
    var modal = document.getElementById('shortcuts-modal');
    if (modal) {
      shortcutsModalOpener = document.activeElement;
      // Restore the dialog to the document before moving focus inside.
      modal.removeAttribute('hidden');
      modal.removeAttribute('inert');
      modal.classList.add('shortcuts-modal--open');
      modal.setAttribute('aria-hidden', 'false');
      // Focus the close button for accessibility
      var closeBtn = modal.querySelector('.shortcuts-modal-close');
      if (closeBtn) {
        closeBtn.focus();
      }
      // Prevent body scrolling while modal is open
      document.body.style.overflow = 'hidden';
    }
  }

  /**
   * Hide the shortcuts help modal
   */
  function hideShortcutsModal() {
    var modal = document.getElementById('shortcuts-modal');
    if (modal && modal.classList.contains('shortcuts-modal--open')) {
      modal.classList.remove('shortcuts-modal--open');
      // Move focus out before making its subtree inert.
      if (shortcutsModalOpener && shortcutsModalOpener.isConnected && typeof shortcutsModalOpener.focus === 'function') {
        shortcutsModalOpener.focus();
      } else if (modal.contains(document.activeElement)) {
        document.activeElement.blur();
      }
      shortcutsModalOpener = null;
      modal.setAttribute('inert', '');
      modal.setAttribute('aria-hidden', 'true');
      document.body.style.overflow = '';
      return true;
    }
    return false;
  }

  /**
   * Close all open modals and clear search focus
   */
  function closeModals() {
    var closed = false;

    // Close shortcuts modal
    if (hideShortcutsModal()) {
      closed = true;
    }

    // Blur search input if focused
    var activeElement = document.activeElement;
    if (activeElement && (activeElement.tagName === 'INPUT' || activeElement.tagName === 'TEXTAREA')) {
      if (window.dismissBleveSearch && activeElement.id === 'pagefind-search-input') {
        window.dismissBleveSearch();
      }
      activeElement.blur();
      closed = true;
    }

    // Clear Pagefind results if present
    var pagefindResults = document.querySelector('.pagefind-ui__results-area');
    if (pagefindResults && pagefindResults.children.length > 0) {
      var pagefindInput = document.querySelector('.pagefind-ui__search-input');
      if (pagefindInput) {
        pagefindInput.value = '';
        pagefindInput.dispatchEvent(new Event('input', { bubbles: true }));
        closed = true;
      }
    }

    return closed;
  }

  /**
   * Update the toggle button state in the modal
   */
  function updateToggleButton() {
    var toggleBtn = document.getElementById('shortcuts-toggle');
    if (toggleBtn) {
      var disabled = window.shortcutsRegistry.areDisabled();
      toggleBtn.textContent = disabled ? 'Enable Shortcuts' : 'Disable Shortcuts';
      toggleBtn.setAttribute('aria-pressed', (!disabled).toString());
    }
  }

  /**
   * Update the modifier key display in the modal based on platform
   */
  function updateModifierKeyDisplay() {
    var isMac = isMacPlatform();
    var macKeys = document.querySelectorAll('.kbd-mac');
    var winKeys = document.querySelectorAll('.kbd-win');

    macKeys.forEach(function(el) {
      el.style.display = isMac ? 'inline-block' : 'none';
    });

    winKeys.forEach(function(el) {
      el.style.display = isMac ? 'none' : 'inline-block';
    });
  }

  /**
   * Check if the shortcuts modal is currently open
   */
  function isModalOpen() {
    var modal = document.getElementById('shortcuts-modal');
    return modal && modal.classList.contains('shortcuts-modal--open');
  }

  /**
   * Handle click outside modal to close it
   */
  function handleModalBackdropClick(e) {
    var modal = document.getElementById('shortcuts-modal');
    if (!modal) return;

    // Close if clicking the backdrop (the modal overlay itself) or
    // any area outside the modal content
    var content = modal.querySelector('.shortcuts-modal-content');
    if (content && !content.contains(e.target)) {
      hideShortcutsModal();
    }
  }

  /**
   * Initialize search shortcuts
   */
  function init() {
    // Register / shortcut for search
    window.shortcutsRegistry.register({
      key: '/',
      modifiers: [],
      description: 'Focus search input',
      group: 'search',
      handler: function(e) {
        e.preventDefault();
        focusSearch();
      },
      priority: 100
    });

    // Register Ctrl/Cmd+K for search (need to handle modifiers specially)
    document.addEventListener('keydown', function(e) {
      if (window.shortcutsRegistry.areDisabled()) return;
      if (window.shortcutsRegistry.isInputElement(e.target)) return;

      var modifier = isMacPlatform() ? e.metaKey : e.ctrlKey;
      if (modifier && e.key === 'k') {
        e.preventDefault();
        focusSearch();
      }
    });

    // Register ? for help modal (toggles open/close)
    window.shortcutsRegistry.register({
      key: '?',
      modifiers: [],
      description: 'Toggle shortcuts help',
      group: 'help',
      handler: function(e) {
        e.preventDefault();
        if (isModalOpen()) {
          hideShortcutsModal();
        } else {
          showShortcutsModal();
        }
      },
      priority: 100
    });

    // Handle Escape and x to close modals (always works, even in inputs)
    document.addEventListener('keydown', function(e) {
      if (e.key === 'Escape') {
        if (closeModals()) {
          e.preventDefault();
        }
      } else if (e.key === 'x' && isModalOpen()) {
        // x closes the shortcuts modal when it is open
        e.preventDefault();
        hideShortcutsModal();
      }
    });

    // Setup modal close button
    var closeBtn = document.querySelector('.shortcuts-modal-close');
    if (closeBtn) {
      closeBtn.addEventListener('click', hideShortcutsModal);
    }

    // Setup toggle button
    var toggleBtn = document.getElementById('shortcuts-toggle');
    if (toggleBtn) {
      toggleBtn.addEventListener('click', function() {
        window.shortcutsRegistry.toggleAll();
        updateToggleButton();
      });
      updateToggleButton();
    }

    // Close modal on backdrop click
    var modal = document.getElementById('shortcuts-modal');
    if (modal) {
      modal.addEventListener('click', handleModalBackdropClick);
    }

    // Update modifier key display
    updateModifierKeyDisplay();
  }

  // Capture the 404 recovery submit before the generated template's legacy
  // redirect handler. This keeps the user on the recovery page and opens the
  // same Pagefind/bleve search used by the site's normal header search.
  document.addEventListener('submit', handle404RecoverySubmit, true);

  // Apply a deep-linked search independently of shortcut-registry readiness.
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', applySearchQueryParam, { once: true });
  } else {
    applySearchQueryParam();
  }

  // Initialize when DOM is ready (or registry is available)
  waitForRegistry(function() {
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', init);
    } else {
      init();
    }
  });

  // Expose for backward compatibility
  window.markataShortcuts = {
    focusSearch: focusSearch,
    showShortcutsModal: showShortcutsModal,
    hideShortcutsModal: hideShortcutsModal,
    toggleShortcuts: function() {
      return window.shortcutsRegistry.toggleAll();
    },
    areShortcutsDisabled: function() {
      return window.shortcutsRegistry.areDisabled();
    }
  };
})();
