# Keyboard Shortcuts System - Implementation Guide

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Core Registry](#core-registry)
3. [Feature Modules](#feature-modules)
4. [Adding New Shortcuts](#adding-new-shortcuts)
5. [API Reference](#api-reference)
6. [Best Practices](#best-practices)
7. [Troubleshooting](#troubleshooting)

---

## Architecture Overview

### System Diagram

```
┌─────────────────────────────────────────────────────────┐
│ Browser Keyboard Event (keydown)                        │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│ shortcuts-registry.js (Central Event Listener)          │
│  - Validates input element status                       │
│  - Checks if shortcuts are disabled                     │
│  - Finds matching shortcuts by priority                 │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│ Execute Highest Priority Match                          │
│  - Call handler function                               │
│  - Prevent default if needed                           │
└─────────────────┬───────────────────────────────────────┘
                  │
        ┌─────────┼─────────┬─────────┬──────────┐
        ▼         ▼         ▼         ▼          ▼
    search   scrolling navigation history  palette
    module   module    module      module   module

    ↓         ↓         ↓         ↓         ↓
    /,?      d,u,gg   j,k,o,g   h,l      [,],{,}\
           (posts)   (feeds)
```

### Module Loading Order

**base.html** (template) loads modules in this order:

```html
1. shortcuts-registry.js         (FIRST - creates window.shortcutsRegistry)
   ↓
2. search-shortcuts.js           (conditional - if search/modal exists)
   ↓
3. scrolling-shortcuts.js        (always loaded)
   ↓
4. navigation-shortcuts.js       (conditional - if .card elements exist)
   ↓
5. history-shortcuts.js          (always loaded)
   ↓
6. palette-switcher.js           (conditional - if switcher enabled)
```

**Important:** Later modules cannot load before earlier ones due to dependencies.

---

## Core Registry

### Location
`pkg/themes/default/static/js/shortcuts-registry.js` (~400 lines)

### Key Responsibilities

1. **Global Event Listener**
   - Single `keydown` listener prevents multiple listeners
   - Fires before individual module listeners

2. **Input Detection**
   - Detects when user is typing in form inputs
   - Blocks shortcuts to prevent interference
   - Exception: Escape always works

3. **Priority Queue**
   - Executes shortcuts by priority (highest first)
   - Prevents lower-priority shortcuts from executing

4. **Enable/Disable Management**
   - Global disable via localStorage
   - Group-based disable
   - Individual shortcut disable

5. **Modifier Key Handling**
   - Ctrl/Alt/Meta when unexpected = no match
   - Shift allowed as "variant" (e.g., `?` = Shift+/)

### Core Functions

```javascript
// Register a shortcut
window.shortcutsRegistry.register({
  key: 'j',
  modifiers: [],
  description: 'Scroll down',
  group: 'scrolling',
  handler: function(e) { scroll(100); },
  priority: 10
});

// Get all shortcuts
window.shortcutsRegistry.getAll()

// Get shortcuts by group
window.shortcutsRegistry.getShortcutsByGroup()

// Enable/disable
window.shortcutsRegistry.setEnabled('j', false)
window.shortcutsRegistry.setGroupEnabled('scrolling', true)

// Check state
window.shortcutsRegistry.areDisabled()
window.shortcutsRegistry.isGroupEnabled('navigation')
window.shortcutsRegistry.isInputElement(element)
```

---

## Feature Modules

### 1. search-shortcuts.js

**Purpose:** Search and help modal

**Shortcuts:**
- `/` - Focus search input
- `Ctrl+K` / `Cmd+K` - Focus search (alternative)
- `?` - Show help modal
- `Escape` - Close modals (always works)

**Loading:** Conditional - only if `#pagefind-search`, `#shortcuts-modal`, or `[type="search"]` exist

**Module Structure:**
```javascript
(function() {
  'use strict';

  // Utility functions
  function openSearchModal() { ... }
  function closeModals() { ... }

  // Wait for registry
  waitForRegistry(init);

  // Registration
  function init() {
    window.shortcutsRegistry.register({...});
  }
})();
```

---

### 2. scrolling-shortcuts.js

**Purpose:** Page scrolling and content navigation

**Shortcuts:**
- `d` - Scroll half-page down (always)
- `u` - Scroll half-page up (always)
- `g g` - Scroll to top (always)
- `Shift+G` - Scroll to bottom (always)
- `j` - Scroll down 100px (post pages only)
- `k` - Scroll up 100px (post pages only)

**Smart Detection:**
- `isPostPage()` checks for `<article>`, `[data-type="post"]`, or `.post-content`
- j/k only register if on post page
- Prevents conflicts with feed navigation

**Two-Key Sequences:**
- `g g` requires two 'g' presses within 500ms
- Implemented with custom keydown listener (not registry)

**Loading:** Always

---

### 3. navigation-shortcuts.js

**Purpose:** Feed card and page navigation

**Shortcuts:**
- `j` / `k` - Navigate cards (highlight with outline)
- `o` / `Enter` - Open highlighted card
- `Shift+O` - Open card in new tab
- `[` / `]` - Previous/next page
- `g h` - Go home
- `g s` - Focus search
- `y y` - Copy card URL to clipboard

**Context:**
- Only works when `.card` or `[data-card]` elements exist
- Priority 20 (higher than scrolling j/k)
- Highlights selected card with `.kb-highlighted` class

**Visual Feedback:**
- CSS class: `.kb-highlighted`
- Defined in `components.css`
- 2px outline + background tint + smooth transition

**Loading:** Conditional - only if `.card` or `[data-card]` exist

---

### 4. history-shortcuts.js

**Purpose:** Browser history navigation

**Shortcuts:**
- `h` - Go back (window.history.back())
- `l` - Go forward (window.history.forward())

**Priority:** 15 (navigation group)

**Considerations:**
- Works with traditional multi-page navigation
- May not work with SPA (Single Page Application) routing
- Browser security may prevent programmatic history changes

**Loading:** Always

---

### 5. palette-switcher.js (Refactored)

**Purpose:** Theme and palette switching

**Shortcuts:**
- `[` - Previous palette
- `]` - Next palette
- `{` - Previous aesthetic (Shift+[)
- `}` - Next aesthetic (Shift+])
- `\` - Toggle dark/light mode

**Loading:** Conditional - if palette switcher UI exists

**Refactoring:** Updated to use registry instead of direct listeners

---

## Adding New Shortcuts

### Step 1: Create Module File

Create `pkg/themes/default/static/js/your-module.js`:

```javascript
/**
 * Your Shortcut Module
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

  // Your handler functions
  function yourHandler() {
    // Implementation
  }

  // Initialize shortcuts
  function init() {
    window.shortcutsRegistry.register({
      key: 'x',
      modifiers: [],
      description: 'Your shortcut description',
      group: 'your-group',
      handler: function(e) {
        e.preventDefault();
        yourHandler();
      },
      priority: 10
    });
  }

  // Initialize when registry is ready
  waitForRegistry(function() {
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', init);
    } else {
      init();
    }
  });
})();
```

### Step 2: Register in Template

Edit `pkg/themes/default/templates/base.html` (around line 200-220):

```html
<!-- Your Module Shortcuts -->
<!-- (conditional or always) -->
const yourScript = document.createElement('script');
yourScript.src = '{{ "js/your-module.js" | theme_asset }}';
yourScript.defer = true;
document.body.appendChild(yourScript);
```

### Step 3: Decide on Loading Strategy

**Always Load:**
- Core functionality
- Shortcuts used on most pages
- Low overhead

**Conditional Load:**
- Feature-specific shortcuts
- Depends on page elements
- Improves performance

Example conditional:
```javascript
if (document.querySelector('.your-element')) {
  // Load your module
}
```

### Step 4: Testing

1. Verify registry loads
2. Test shortcuts work
3. Check for console errors
4. Verify no conflicts
5. Test input detection
6. Check help modal displays shortcuts

---

## API Reference

### `window.shortcutsRegistry`

#### `register(config)`

Register a new shortcut.

**Parameters:**
```javascript
config = {
  key: string,              // The key to trigger (required)
  modifiers: string[],      // ['Ctrl', 'Alt', 'Shift', 'Meta'] (optional, default: [])
  description: string,      // Help text (required)
  group: string,            // Category (optional, default: 'other')
  handler: function(e),     // Callback (required)
  priority: number          // Execution order (optional, default: 0, higher first)
}
```

**Returns:** The registered shortcut object

**Example:**
```javascript
window.shortcutsRegistry.register({
  key: 'j',
  modifiers: [],
  description: 'Jump to next',
  group: 'navigation',
  handler: function(e) {
    e.preventDefault();
    jumpToNext();
  },
  priority: 10
});
```

---

#### `getAll()`

Get array of all registered shortcuts.

**Returns:** Array of shortcut objects

**Example:**
```javascript
const shortcuts = window.shortcutsRegistry.getAll();
shortcuts.forEach(s => {
  console.log(`${s.key}: ${s.description}`);
});
```

---

#### `getShortcutsByGroup()`

Get shortcuts organized by group.

**Returns:** Object with groups as keys, arrays of shortcuts as values

**Example:**
```javascript
const byGroup = window.shortcutsRegistry.getShortcutsByGroup();
// {
//   navigation: [...],
//   scrolling: [...],
//   search: [...],
//   ...
// }
```

---

#### `setEnabled(key, enabled)`

Enable or disable a specific shortcut.

**Parameters:**
- `key` (string) - The key to enable/disable
- `enabled` (boolean) - true to enable, false to disable

**Example:**
```javascript
// Disable 'j' shortcut
window.shortcutsRegistry.setEnabled('j', false);

// Re-enable
window.shortcutsRegistry.setEnabled('j', true);
```

---

#### `setGroupEnabled(group, enabled)`

Enable or disable all shortcuts in a group.

**Parameters:**
- `group` (string) - The group name
- `enabled` (boolean) - true to enable, false to disable

**Persists to localStorage automatically**

**Example:**
```javascript
// Disable all navigation shortcuts
window.shortcutsRegistry.setGroupEnabled('navigation', false);

// Re-enable
window.shortcutsRegistry.setGroupEnabled('navigation', true);
```

---

#### `isGroupEnabled(group)`

Check if a shortcut group is enabled.

**Parameters:**
- `group` (string) - The group name

**Returns:** boolean

**Example:**
```javascript
if (window.shortcutsRegistry.isGroupEnabled('navigation')) {
  console.log('Navigation shortcuts are enabled');
}
```

---

#### `areDisabled()`

Check if all shortcuts are globally disabled.

**Returns:** boolean

**Example:**
```javascript
if (window.shortcutsRegistry.areDisabled()) {
  console.log('All shortcuts are disabled');
}
```

---

#### `toggleDisabled()`

Toggle global disabled state.

**Returns:** boolean (new state)

**Example:**
```javascript
// Toggle shortcuts on/off
const newState = window.shortcutsRegistry.toggleDisabled();
console.log('Shortcuts ' + (newState ? 'enabled' : 'disabled'));
```

---

#### `isInputElement(element)`

Check if element is an input where shortcuts should be blocked.

**Parameters:**
- `element` (Element) - DOM element to check

**Returns:** boolean

**Detects:**
- `<textarea>`
- `<input type="text">`
- `<input type="search">`
- `contenteditable` elements
- ARIA roles: textbox, searchbox, combobox

**Allows shortcuts on:**
- `<input type="button">`
- `<input type="checkbox">`
- `<input type="radio">`
- `<input type="range">`
- `<input type="color">`
- `<input type="file">`

**Example:**
```javascript
if (window.shortcutsRegistry.isInputElement(document.activeElement)) {
  console.log('User is typing in an input');
}
```

---

## Best Practices

### 1. Use Meaningful Groups

Organize shortcuts by function:
- `navigation` - Movement and page selection
- `scrolling` - Viewport movement
- `search` - Search and discovery
- `theme` - Visual settings
- `editing` - Content manipulation

### 2. Set Appropriate Priorities

**Priority Guidelines:**
- `30+` - System/critical (Escape, etc.)
- `20-29` - Context-specific navigation
- `10-19` - Global utilities
- `0-9` - Low-priority alternatives

**Example:**
```
Navigation j/k on feed: priority 20
Scrolling j/k on post:  priority 10
d/u scrolling:          priority 10
gg/Shift+G:             priority 10
```

### 3. Provide Clear Descriptions

Descriptions appear in help modal:
- Use active verbs: "Scroll down", "Go home"
- Be concise: under 50 characters
- Include modifiers: "Open in new tab (Shift+O)"

### 4. Wait for Registry

Always use `waitForRegistry()` pattern:

```javascript
function waitForRegistry(callback, attempts = 0) {
  if (window.shortcutsRegistry) {
    callback();
  } else if (attempts < 50) {
    setTimeout(function() {
      waitForRegistry(callback, attempts + 1);
    }, 10);
  }
}
```

### 5. Handle DOMContentLoaded

Account for different loading scenarios:

```javascript
waitForRegistry(function() {
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
});
```

### 6. Prevent Default Appropriately

Only call `e.preventDefault()` if you handle the event:

```javascript
handler: function(e) {
  if (conditions.met()) {
    e.preventDefault();  // Only if you'll handle it
    doSomething();
  }
  // If not prevented, browser default occurs
}
```

### 7. Handle Edge Cases

```javascript
function handleCardNavigation(direction) {
  const cards = document.querySelectorAll('.card');
  if (cards.length === 0) return;  // No cards

  const highlighted = document.querySelector('.kb-highlighted');
  if (!highlighted) {
    // First card
    cards[0].classList.add('kb-highlighted');
  } else {
    // Next/prev logic with bounds checking
  }
}
```

### 8. Use LocalStorage Responsibly

Wrap in try-catch for privacy mode:

```javascript
try {
  localStorage.setItem('key', 'value');
} catch (e) {
  // localStorage unavailable
}
```

---

## Troubleshooting

### Issue: Shortcut Doesn't Work

**Diagnosis:**
```javascript
// 1. Check registry exists
console.log(window.shortcutsRegistry);

// 2. Check shortcut registered
const shortcuts = window.shortcutsRegistry.getAll();
console.log(shortcuts.find(s => s.key === 'j'));

// 3. Check group enabled
console.log(window.shortcutsRegistry.isGroupEnabled('scrolling'));

// 4. Check globally disabled
console.log(window.shortcutsRegistry.areDisabled());
```

**Solutions:**
1. Verify module loaded (Network tab)
2. Check for console errors
3. Verify registry exists before registration
4. Check DOM elements exist (for conditional loading)
5. Clear browser cache

---

### Issue: Shortcut Triggers in Text Input

**Diagnosis:**
```javascript
// Test input detection
document.querySelector('input').focus();
window.shortcutsRegistry.isInputElement(document.activeElement)
// Should return: true
```

**Solution:**
Verify element matches detection criteria in `isInputElement()`:
- Is it a TEXTAREA?
- Is it INPUT with text type?
- Does it have `contenteditable`?
- Does it have text input ARIA role?

---

### Issue: Multiple Shortcuts Execute

**Cause:** Lower priority shortcut executing after higher priority

**Diagnosis:**
```javascript
// Check priorities
window.shortcutsRegistry.getAll()
  .filter(s => s.key === 'j')
  .forEach(s => console.log(s.priority));
```

**Solution:**
Adjust priorities so only one matches:
- Feed: navigation j/k priority 20
- Post: scrolling j/k priority 10
- Use conditional registration

---

### Issue: Escape Doesn't Close Modal

**Cause:** Modal's Escape handler runs first (not registry's)

**Solution:**
Ensure modal's Escape handler doesn't have `stopPropagation()`:

```javascript
// BAD - blocks registry
document.addEventListener('keydown', function(e) {
  if (e.key === 'Escape') {
    closeModal();
    e.stopPropagation();  // ← PROBLEM
  }
});

// GOOD - registry still gets event
document.addEventListener('keydown', function(e) {
  if (e.key === 'Escape') {
    closeModal();
    // No stopPropagation - registry can also handle
  }
});
```

---

### Issue: Performance Lag on Key Press

**Diagnosis:**
1. DevTools → Performance → Record shortcut usage
2. Check flame chart for bottlenecks
3. Look for event listener execution time

**Solutions:**
1. Debounce handler if expensive:
   ```javascript
   function debounce(fn, delay) {
     let timeout;
     return function() {
       clearTimeout(timeout);
       timeout = setTimeout(fn, delay);
     };
   }
   ```

2. Move DOM queries outside handler:
   ```javascript
   // BAD - queries every time
   handler: function(e) {
     document.querySelectorAll('.card').forEach(...)
   }

   // GOOD - cache query
   const cards = document.querySelectorAll('.card');
   handler: function(e) {
     cards.forEach(...)
   }
   ```

3. Use event delegation instead of multiple listeners

---

### Issue: Module Doesn't Load

**Diagnosis:**
```javascript
// Check if registry initialized
console.log(window.shortcutsRegistry);

// Check if module's init function ran
// (Add console.log in init)

// Check network tab for module file
// Check for CORS errors
```

**Solutions:**
1. Ensure registry loads first
2. Check conditional loading criteria met
3. Verify file path correct in template
4. Check CORS headers on server
5. Look for JavaScript errors preventing initialization

---

## Performance Considerations

### Memory
- Registry stores only shortcut metadata (not handlers)
- Single global event listener (not per shortcut)
- Disabled shortcuts don't execute (no memory overhead)

### CPU
- keydown fires on every key press (cannot avoid)
- findMatches() is O(n) where n = number of shortcuts (~20 typical)
- Handlers should be < 16ms for 60fps

### Network
- Conditional loading reduces initial JS size
- Only load modules when needed
- Minify/bundle in production

---

## Future Enhancements

1. **Auto-generate Help Modal**
   - Extract from registered shortcuts
   - Dynamic shortcut documentation

2. **Conflict Detection**
   - Warn on conflicting registrations
   - Suggest priority adjustments

3. **Custom Keybindings**
   - Allow user to remap shortcuts
   - Persist to localStorage

4. **Analytics**
   - Track shortcut usage
   - Identify underused shortcuts

5. **Debugging Mode**
   - Visual indicator when shortcut executes
   - Performance timing
   - Event flow visualization

6. **Mobile/Touch**
   - Gesture alternatives
   - Touch-based shortcuts

---

**Document Version:** 1.0
**Last Updated:** [DATE]
**Maintained By:** [TEAM]
