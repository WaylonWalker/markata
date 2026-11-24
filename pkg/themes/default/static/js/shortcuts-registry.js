/**
 * Generic Shortcuts Registry for markata-go
 *
 * Provides a centralized system for registering and managing keyboard shortcuts.
 * Plugins can register their shortcuts with this system instead of managing their own
 * keyboard event listeners.
 *
 * Usage:
 *   window.shortcutsRegistry.register({
 *     key: 'j',
 *     modifiers: [],
 *     description: 'Scroll down',
 *     group: 'scrolling',
 *     handler: function(e) { ... },
 *     priority: 10
 *   });
 *
 * Accessibility: Shortcuts can be disabled via localStorage
 * WCAG 2.1.4 compliant - shortcuts are ignored when typing in inputs
 */

(function() {
  'use strict';

  // Storage keys
  var STORAGE_KEY_DISABLED = 'markata-shortcuts-disabled';
  var STORAGE_KEY_DISABLED_GROUPS = 'markata-shortcuts-disabled-groups';

  /**
   * Shortcut Registry
   * Manages registration, execution, and lifecycle of shortcuts
   */
  var ShortcutRegistry = function() {
    this.shortcuts = [];
    this.groups = new Map();
  };

  /**
   * Register a new shortcut
   * @param {Object} config - Shortcut configuration
   * @param {string} config.key - The key(s) to trigger (e.g., 'j', 'Enter', 'gg')
   * @param {Array<string>} [config.modifiers] - Modifier keys (e.g., ['Ctrl', 'Shift'])
   * @param {string} config.description - Human-readable description for help modal
   * @param {string} [config.group] - Grouping category (e.g., 'scrolling', 'navigation')
   * @param {Function} config.handler - Callback function(e, shortcut)
   * @param {number} [config.priority] - Priority for execution order (higher = first)
   * @returns {Object} The registered shortcut
   */
  ShortcutRegistry.prototype.register = function(config) {
    if (!config.key || !config.description || !config.handler) {
      console.error('[shortcuts-registry] Missing required config: key, description, handler');
      return null;
    }

    var shortcut = {
      key: config.key,
      modifiers: config.modifiers || [],
      description: config.description,
      group: config.group || 'other',
      handler: config.handler,
      priority: config.priority || 0,
      enabled: true
    };

    this.shortcuts.push(shortcut);

    // Track group
    if (!this.groups.has(shortcut.group)) {
      this.groups.set(shortcut.group, []);
    }
    this.groups.get(shortcut.group).push(shortcut);

    // Sort by priority (higher first)
    this.shortcuts.sort(function(a, b) {
      return b.priority - a.priority;
    });

    return shortcut;
  };

  /**
   * Find shortcuts matching a key combination
   * @param {KeyboardEvent} e
   * @returns {Array<Object>} Matching shortcuts
   */
  ShortcutRegistry.prototype.findMatches = function(e) {
    var matches = [];
    var key = e.key;

    for (var i = 0; i < this.shortcuts.length; i++) {
      var shortcut = this.shortcuts[i];
      if (!shortcut.enabled) continue;

      // Check if key matches
      if (shortcut.key !== key) continue;

      // Check modifiers
      var modifiersMatch = this._checkModifiers(e, shortcut.modifiers);
      if (modifiersMatch) {
        matches.push(shortcut);
      }
    }

    return matches;
  };

  /**
   * Check if keyboard event modifiers match expected modifiers
   * @private
   */
  ShortcutRegistry.prototype._checkModifiers = function(e, expectedModifiers) {
    var expectedSet = new Set(expectedModifiers.map(function(m) { return m.toLowerCase(); }));

    var actualModifiers = [];
    if (e.ctrlKey) actualModifiers.push('ctrl');
    if (e.altKey) actualModifiers.push('alt');
    if (e.shiftKey) actualModifiers.push('shift');
    if (e.metaKey) actualModifiers.push('meta');

    // If no modifiers expected, allow any printable character (even if shift is pressed)
    // This handles cases like ? which is Shift+/ on US keyboards
    if (expectedSet.size === 0) {
      // For unmodified shortcuts, we DON'T want Ctrl/Meta/Alt to be held
      // But we DO allow Shift since it's just the shift variant of a key
      var significantModifiers = actualModifiers.filter(function(m) {
        return m !== 'shift';
      });
      return significantModifiers.length === 0;
    }

    // For shortcuts with expected modifiers, check exact match
    if (expectedSet.size !== actualModifiers.length) {
      return false;
    }

    for (var i = 0; i < actualModifiers.length; i++) {
      if (!expectedSet.has(actualModifiers[i])) {
        return false;
      }
    }

    return true;
  };

  /**
   * Enable/disable a specific shortcut
   */
  ShortcutRegistry.prototype.setEnabled = function(key, enabled) {
    for (var i = 0; i < this.shortcuts.length; i++) {
      if (this.shortcuts[i].key === key) {
        this.shortcuts[i].enabled = enabled;
      }
    }
  };

  /**
   * Enable/disable all shortcuts in a group
   */
  ShortcutRegistry.prototype.setGroupEnabled = function(group, enabled) {
    if (this.groups.has(group)) {
      var shortcuts = this.groups.get(group);
      for (var i = 0; i < shortcuts.length; i++) {
        shortcuts[i].enabled = enabled;
      }
    }
    this._persistDisabledGroups();
  };

  /**
   * Check if a shortcut group is enabled
   */
  ShortcutRegistry.prototype.isGroupEnabled = function(group) {
    if (!this.groups.has(group)) return true;
    var shortcuts = this.groups.get(group);
    return shortcuts.length > 0 && shortcuts[0].enabled;
  };

  /**
   * Get all shortcuts organized by group
   */
  ShortcutRegistry.prototype.getShortcutsByGroup = function() {
    var byGroup = {};
    var self = this;

    this.groups.forEach(function(shortcuts, group) {
      byGroup[group] = shortcuts;
    });

    return byGroup;
  };

  /**
   * Get all registered shortcuts
   */
  ShortcutRegistry.prototype.getAll = function() {
    return this.shortcuts.slice();
  };

  /**
   * Persist disabled groups to localStorage
   * @private
   */
  ShortcutRegistry.prototype._persistDisabledGroups = function() {
    try {
      var disabledGroups = [];
      var self = this;

      this.groups.forEach(function(shortcuts, group) {
        if (shortcuts.length > 0 && !shortcuts[0].enabled) {
          disabledGroups.push(group);
        }
      });

      localStorage.setItem(STORAGE_KEY_DISABLED_GROUPS, JSON.stringify(disabledGroups));
    } catch (e) {
      // localStorage may be unavailable
    }
  };

  /**
   * Load disabled groups from localStorage
   * @private
   */
  ShortcutRegistry.prototype._loadDisabledGroups = function() {
    try {
      var stored = localStorage.getItem(STORAGE_KEY_DISABLED_GROUPS);
      if (stored) {
        var disabledGroups = JSON.parse(stored);
        var self = this;
        disabledGroups.forEach(function(group) {
          self.setGroupEnabled(group, false);
        });
      }
    } catch (e) {
      // localStorage may be unavailable or JSON invalid
    }
  };

  /**
   * Global instance
   */
  var registry = new ShortcutRegistry();

  /**
   * Check if shortcuts are disabled globally via localStorage
   * @returns {boolean}
   */
  function areShortcutsDisabled() {
    try {
      return localStorage.getItem(STORAGE_KEY_DISABLED) === 'true';
    } catch (e) {
      return false;
    }
  }

  /**
   * Toggle all shortcuts enabled/disabled state
   * @returns {boolean} New disabled state
   */
  function toggleAllShortcuts() {
    try {
      var currentlyDisabled = areShortcutsDisabled();
      localStorage.setItem(STORAGE_KEY_DISABLED, (!currentlyDisabled).toString());
      return !currentlyDisabled;
    } catch (e) {
      return false;
    }
  }

  /**
   * Check if the current element is a form field or interactive control that
   * must keep its native keyboard behavior.
   * @param {Element} element
   * @returns {boolean}
   */
  function deepActiveElement() {
    var active = document.activeElement;
    while (active && active.shadowRoot && active.shadowRoot.activeElement) {
      active = active.shadowRoot.activeElement;
    }
    return active;
  }

  function isInputElement(element) {
    // Events from inside a shadow root (e.g. the serve settings panel) are
    // retargeted to the shadow host, so also check the focused element inside
    // any open shadow roots.
    var elementsToCheck = [element];
    var active = deepActiveElement();
    if (active && active !== element) {
      elementsToCheck.push(active);
    }

    for (var i = 0; i < elementsToCheck.length; i++) {
      var el = elementsToCheck[i];
      if (!el) continue;

      // Global shortcuts must not consume Enter/Space (or other registered
      // keys) from a focused native or ARIA control. In particular, Enter on a
      // copy button must reach the button's native activation behavior.
      if (typeof el.closest === 'function' && el.closest(
        'button, a[href], summary, [role="button"], [role="link"], ' +
        '[role="menuitem"], [role="menuitemcheckbox"], [role="menuitemradio"], ' +
        '[role="option"], [role="tab"], [role="checkbox"], [role="radio"], ' +
        '[role="switch"], [contenteditable="true"]'
      )) {
        return true;
      }

      var tagName = el.tagName;

      // TEXTAREA always blocks shortcuts
      if (tagName === 'TEXTAREA') {
        return true;
      }

      // SELECT always blocks shortcuts (keyboard navigation)
      if (tagName === 'SELECT') {
        return true;
      }

      // INPUT: block for text-input types, allow for non-text types
      if (tagName === 'INPUT') {
        var inputType = (el.type || 'text').toLowerCase();
        // Non-text input types that should ALLOW shortcuts
        var nonTextTypes = [
          'submit',
          'button',
          'reset',
          'checkbox',
          'radio',
          'range',
          'color',
          'file',
          'hidden',
          'image'
        ];
        // If it's NOT a non-text type, block shortcuts
        if (nonTextTypes.indexOf(inputType) === -1) {
          return true;
        }
      }

      // Check for contenteditable
      if (el.isContentEditable) {
        return true;
      }

      // Check ARIA roles that indicate text input
      var role = el.getAttribute('role');
      if (role) {
        var textInputRoles = ['textbox', 'searchbox', 'combobox'];
        if (textInputRoles.indexOf(role.toLowerCase()) !== -1) {
          return true;
        }
      }
    }

    return false;
  }

  /**
   * Handle keyboard events and dispatch to registered shortcuts
   * @param {KeyboardEvent} e
   */
  function handleKeyDown(e) {
    // Always allow Escape to work, even in inputs
    if (e.key === 'Escape') {
      // Escape is always available for core system to handle
      return;
    }

    // Skip shortcuts when typing in input fields
    var origin = typeof e.composedPath === 'function' ? e.composedPath()[0] : e.target;
    if (isInputElement(origin || e.target)) {
      return;
    }

    // Check if shortcuts are disabled globally
    if (areShortcutsDisabled()) {
      return;
    }

    // Find and execute matching shortcuts
    var matches = registry.findMatches(e);
    for (var i = 0; i < matches.length; i++) {
      var shortcut = matches[i];
      try {
        shortcut.handler(e, shortcut);
        e.preventDefault();
        break; // Execute only the first (highest priority) match
      } catch (error) {
        console.error('[shortcuts-registry] Error executing shortcut:', error);
      }
    }
  }

  /**
   * Initialize the shortcuts registry system
   */
  function init() {
    // Load disabled groups from localStorage
    registry._loadDisabledGroups();

    // Add keyboard event listener
    document.addEventListener('keydown', handleKeyDown);
  }

  // Initialize when DOM is ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }

  /**
   * Expose the registry for plugins
   */
  window.shortcutsRegistry = {
    register: registry.register.bind(registry),
    setEnabled: registry.setEnabled.bind(registry),
    setGroupEnabled: registry.setGroupEnabled.bind(registry),
    isGroupEnabled: registry.isGroupEnabled.bind(registry),
    getShortcutsByGroup: registry.getShortcutsByGroup.bind(registry),
    getAll: registry.getAll.bind(registry),
    toggleAll: toggleAllShortcuts,
    areDisabled: areShortcutsDisabled,
    isInputElement: isInputElement
  };
})();
