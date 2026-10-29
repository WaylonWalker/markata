# Keyboard Shortcuts Testing Guide

This document provides comprehensive testing procedures for the generic keyboard shortcuts registry system in markata-go.

## Overview

The keyboard shortcuts system consists of:
- **Core Registry** (`shortcuts-registry.js`) - Centralized shortcut management
- **Feature Modules**:
  - `search-shortcuts.js` - Search and help modal
  - `scrolling-shortcuts.js` - Page and card scrolling
  - `navigation-shortcuts.js` - Feed card navigation
  - `history-shortcuts.js` - Browser history navigation
  - `palette-switcher.js` - Theme and palette switching

## Test Environment Setup

### Prerequisites
1. Access to a running markata-go instance
2. Browser DevTools open (F12 or Ctrl+Shift+I)
3. Console tab visible to check for errors
4. Multiple test pages:
   - Feed page (e.g., `/feed/`)
   - Post/Article page (e.g., `/posts/some-post/`)
   - Home page

### Initial Checks
1. Open DevTools console
2. Check that `window.shortcutsRegistry` is defined
3. Verify no JavaScript errors appear
4. Check all shortcut modules loaded (look for console logs if present)

---

## Test Categories

### 1. FEED PAGE SHORTCUTS

**Test Location:** `/feed/` or any page with `.card` or `[data-card]` elements

#### 1.1 Card Navigation - j/k
- [ ] Press `j` → Next card should highlight with 2px outline
- [ ] Press `j` multiple times → Should cycle through cards
- [ ] Press `k` → Previous card should highlight
- [ ] Last card + `j` → Should not wrap around
- [ ] First card + `k` → Should not wrap around
- [ ] Visual feedback: Card should have `.kb-highlighted` class with outline

#### 1.2 Open Card - o / Enter
- [ ] Select a card with `j`/`k`
- [ ] Press `o` → Card should navigate to post
- [ ] Repeat with `Enter` → Same behavior as `o`

#### 1.3 Open in New Tab - Shift+O
- [ ] Select a card with `j`/`k`
- [ ] Press `Shift+O` → Card should open in new tab
- [ ] Current tab should remain on feed page

#### 1.4 Pagination - [ / ]
- [ ] On multi-page feed, press `[` → Previous page should load
- [ ] Press `]` → Next page should load
- [ ] First page + `[` → Should do nothing (or keep on first page)
- [ ] Last page + `]` → Should do nothing (or keep on last page)

#### 1.5 Quick Navigation - g sequences
- [ ] Press `g h` → Should navigate to home page
- [ ] Press `g s` → Focus should move to search input (or search modal opens)

#### 1.6 Copy URL - yy
- [ ] Select a card with `j`/`k`
- [ ] Press `y y` → Card URL should copy to clipboard
- [ ] Verify clipboard by pasting (Ctrl+V) in text field

---

### 2. POST/ARTICLE PAGE SHORTCUTS

**Test Location:** Any page with `<article>` element or `[data-type="post"]` attribute

#### 2.1 Scrolling - j/k on Post Pages
- [ ] Press `j` → Page should scroll down ~100px (smooth)
- [ ] Press `j` repeatedly → Page should continue scrolling
- [ ] Press `k` → Page should scroll up ~100px (smooth)
- [ ] At top of page + `k` → Should not go above 0
- [ ] At bottom of page + `j` → Should not scroll past bottom
- [ ] **Critical:** j/k should NOT navigate cards on post page

#### 2.2 Half-Page Scrolling - d/u
- [ ] Press `d` → Page should scroll down ~50% of window height
- [ ] Press `u` → Page should scroll up ~50% of window height
- [ ] Multiple presses should accumulate scrolling

#### 2.3 Jump to Top/Bottom - gg / Shift+G
- [ ] Press `g g` (two 'g' presses within 500ms) → Should jump to page top
- [ ] Press `Shift+G` (hold Shift and press G) → Should jump to page bottom
- [ ] When at top, `g g` should have no visible effect
- [ ] When at bottom, `Shift+G` should have no visible effect

#### 2.4 Context Verification
- [ ] Confirm `j/k` on post pages DO scroll, NOT navigate
- [ ] Confirm navigation shortcuts like `o` do NOT work (or gracefully fail)

---

### 3. GLOBAL SHORTCUTS (All Pages)

#### 3.1 Search - / or Ctrl+K (or Cmd+K on Mac)
- [ ] Press `/` → Search input should focus
- [ ] Type in search → Should not trigger shortcuts
- [ ] Press `Escape` → Search modal should close
- [ ] Alt method: Press `Ctrl+K` (or `Cmd+K` on Mac) → Same as `/`

#### 3.2 Help Modal - ?
- [ ] Press `?` → Help modal should open showing all shortcuts
- [ ] Modal should display shortcuts organized by group
- [ ] Press `Escape` → Help modal should close
- [ ] **Note:** `?` is `Shift+/` on most keyboards

#### 3.3 Theme Toggle - \
- [ ] Press `\` (backslash) → Should toggle between light and dark mode
- [ ] Visual change should be immediate
- [ ] Should persist after page reload (if localStorage enabled)

#### 3.4 Palette Switching - [ / ]
- [ ] On pages with palette switcher: Press `[` → Previous palette loads
- [ ] Press `]` → Next palette loads
- [ ] Cycle through available palettes

#### 3.5 Aesthetic Switching - { / }
- [ ] On pages with palette switcher: Press `{` (Shift+[) → Previous aesthetic loads
- [ ] Press `}` (Shift+]) → Next aesthetic loads
- [ ] **Note:** These are Shift+[ and Shift+] on most keyboards

---

### 4. HISTORY NAVIGATION SHORTCUTS

#### 4.1 Go Back - h
- [ ] Visit page A
- [ ] Navigate to page B
- [ ] Press `h` → Should go back to page A
- [ ] Press `h` again → Should go back further in history
- [ ] At first page in history + `h` → Should do nothing

#### 4.2 Go Forward - l
- [ ] Visit pages in sequence: A → B → C
- [ ] Press `h` twice → Back to A
- [ ] Press `l` → Forward to B
- [ ] Press `l` → Forward to C
- [ ] At last page in history + `l` → Should do nothing

---

### 5. INPUT ELEMENT DETECTION

**Purpose:** Verify shortcuts don't interfere with form input

#### 5.1 Text Input Blocking
- [ ] Navigate to a page with search input
- [ ] Click in search input
- [ ] Type a letter (e.g., `j`) → Should type, not scroll
- [ ] Verify `window.shortcutsRegistry.isInputElement()` returns true

#### 5.2 Textarea Blocking
- [ ] On a page with a comment/message textarea
- [ ] Click in textarea
- [ ] Type characters → Should type normally
- [ ] No shortcuts should trigger

#### 5.3 Contenteditable Elements
- [ ] On pages with contenteditable elements
- [ ] Click in editable area
- [ ] Type characters → Should type normally

#### 5.4 Non-Text Inputs Allowed
- [ ] On a page with button inputs
- [ ] Click on a button → Shortcuts should work (if focused)
- [ ] On checkbox inputs → Shortcuts should work
- [ ] On radio buttons → Shortcuts should work

#### 5.5 Escape Always Works
- [ ] Open search modal
- [ ] Click in search input
- [ ] Press `Escape` → Modal should close (even though in input)

---

### 6. CONFLICT DETECTION

**Purpose:** Verify no shortcuts interfere with each other

#### 6.1 j/k Context Switching
- [ ] Feed page: `j/k` should navigate cards (priority 20)
- [ ] Post page: `j/k` should scroll page (priority 10)
- [ ] Verify no double actions occur
- [ ] Verify `navigation-shortcuts.js` priority (20) > `scrolling-shortcuts.js` (10)

#### 6.2 g-sequence Ordering
- [ ] `g h` should go home (not confused with `g g`)
- [ ] `g s` should focus search (not confused with `g g`)
- [ ] `g g` should go to top (two presses within 500ms)
- [ ] Wait >500ms between `g` presses → Should reset sequence

#### 6.3 Modifier Key Conflicts
- [ ] `Shift+O` should be different from `o`
- [ ] `Shift+G` should be different from `g g`
- [ ] `Ctrl+K` should work the same as `/`

#### 6.4 Priority-Based Execution
- [ ] If multiple shortcuts match a key, higher priority executes first
- [ ] No conflicting shortcuts should execute simultaneously

---

### 7. BROWSER COMPATIBILITY

Test on multiple browsers if possible:

#### 7.1 Chrome/Chromium
- [ ] All shortcuts work
- [ ] No console errors
- [ ] localStorage persistence works

#### 7.2 Firefox
- [ ] All shortcuts work
- [ ] No console errors
- [ ] localStorage persistence works

#### 7.3 Safari
- [ ] All shortcuts work
- [ ] No console errors
- [ ] Mac-specific: `Cmd+K` works (not just Ctrl+K)

#### 7.4 Edge
- [ ] All shortcuts work
- [ ] No console errors
- [ ] localStorage persistence works

---

### 8. ACCESSIBILITY TESTS

#### 8.1 Keyboard Navigation
- [ ] All shortcuts work with keyboard only (no mouse needed)
- [ ] Tab key still works for form navigation
- [ ] Focus management is correct

#### 8.2 Help Modal
- [ ] Help modal (`?`) displays all shortcuts
- [ ] Shortcuts are organized by group
- [ ] Descriptions are clear and helpful
- [ ] Modal is keyboard navigable

#### 8.3 Visual Feedback
- [ ] Highlighted cards have visible outline
- [ ] Scrolling animation is smooth
- [ ] No visual glitches or flickering

#### 8.4 Disable/Enable
- [ ] Shortcuts can be disabled via registry settings
- [ ] Disabled shortcuts don't respond
- [ ] Re-enabling works correctly

---

### 9. PERFORMANCE TESTS

#### 9.1 Rapid Key Presses
- [ ] Press keys rapidly → No lag or missed inputs
- [ ] Multiple sequences quickly → All register correctly
- [ ] Frame rate stable during scrolling

#### 9.2 Memory Leaks
- [ ] Open DevTools → Memory tab
- [ ] Use shortcuts extensively
- [ ] Take heap snapshot
- [ ] Compare heap snapshots → No growing memory

#### 9.3 CPU Usage
- [ ] Monitor CPU during shortcut use
- [ ] Smooth scrolling shouldn't spike CPU
- [ ] Registry events shouldn't consume excessive resources

---

## Test Report Template

```markdown
# Keyboard Shortcuts Test Report

**Date:** [DATE]
**Tester:** [NAME]
**Browser:** [BROWSER + VERSION]
**Platform:** [OS + VERSION]

## Results Summary

### Feed Page Shortcuts
- [ ] Card navigation (j/k) - **PASS** / **FAIL**
- [ ] Open card (o/Enter) - **PASS** / **FAIL**
- [ ] Open in new tab (Shift+O) - **PASS** / **FAIL**
- [ ] Pagination ([/]) - **PASS** / **FAIL**
- [ ] Quick navigation (g h, g s) - **PASS** / **FAIL**
- [ ] Copy URL (yy) - **PASS** / **FAIL**

### Post Page Shortcuts
- [ ] Scrolling (j/k) - **PASS** / **FAIL**
- [ ] Half-page scroll (d/u) - **PASS** / **FAIL**
- [ ] Jump to top/bottom (gg/Shift+G) - **PASS** / **FAIL**

### Global Shortcuts
- [ ] Search (/ or Ctrl+K) - **PASS** / **FAIL**
- [ ] Help (?) - **PASS** / **FAIL**
- [ ] Theme toggle (\) - **PASS** / **FAIL**
- [ ] Palette switching ([/]) - **PASS** / **FAIL**
- [ ] Aesthetic switching ({/}) - **PASS** / **FAIL**

### History Navigation
- [ ] Go back (h) - **PASS** / **FAIL**
- [ ] Go forward (l) - **PASS** / **FAIL**

### Input Handling
- [ ] Text input blocking - **PASS** / **FAIL**
- [ ] Textarea blocking - **PASS** / **FAIL**
- [ ] Escape always works - **PASS** / **FAIL**

### Conflict Detection
- [ ] No j/k conflicts - **PASS** / **FAIL**
- [ ] No g-sequence conflicts - **PASS** / **FAIL**
- [ ] Priority-based execution - **PASS** / **FAIL**

## Issues Found

### Critical Issues
1. [Issue 1]
2. [Issue 2]

### Minor Issues
1. [Issue 1]
2. [Issue 2]

## Notes

[Additional observations or testing notes]
```

---

## Debugging Tips

### Console Checks
```javascript
// Check if registry is loaded
window.shortcutsRegistry

// Get all registered shortcuts
window.shortcutsRegistry.getAll()

// Get shortcuts by group
window.shortcutsRegistry.getShortcutsByGroup()

// Check if a specific group is enabled
window.shortcutsRegistry.isGroupEnabled('navigation')

// Check if shortcuts are globally disabled
window.shortcutsRegistry.areDisabled()

// Check input element detection
window.shortcutsRegistry.isInputElement(document.activeElement)
```

### Developer Tools
1. **Event Listener Debugging:**
   - DevTools → Sources → Event Listener Breakpoints
   - Enable "Keyboard" → "keydown"
   - Shortcuts will pause on keydown

2. **Performance Profiling:**
   - DevTools → Performance
   - Record → Use shortcuts → Stop
   - Analyze flame chart for bottlenecks

3. **Network Tab:**
   - Verify all `.js` files load
   - Check for 404s on shortcut files

---

## Common Issues & Solutions

### Issue: Shortcuts Not Working
- **Check:** Console for JavaScript errors
- **Check:** `window.shortcutsRegistry` exists
- **Check:** Page has correct elements (`.card`, `<article>`, etc.)
- **Solution:** Refresh page, clear cache (Ctrl+Shift+Del)

### Issue: j/k Scrolling on Feed (Should Navigate)
- **Cause:** `navigation-shortcuts.js` didn't load
- **Check:** DevTools → Network → Verify `navigation-shortcuts.js` loaded
- **Solution:** Check if `.card` elements exist on page

### Issue: Shortcuts Fire in Text Input
- **Cause:** `isInputElement()` not detecting input
- **Check:** `window.shortcutsRegistry.isInputElement(el)` returns true for input
- **Solution:** Ensure element matches one of the detection rules

### Issue: Escape Not Closing Modal
- **Cause:** Modal has its own escape handler that runs first
- **Check:** Modal escape handler isn't `stopPropagation()`
- **Solution:** Ensure modal escape handler doesn't block registry

### Issue: History Shortcuts Don't Work
- **Cause:** Browser security prevents history changes in some contexts
- **Check:** Test on different page navigation (not same-page SPAs)
- **Solution:** Works best with traditional multi-page sites

---

## Regression Testing Checklist

When making changes to shortcuts, verify:
- [ ] All previous shortcuts still work
- [ ] New shortcuts don't break existing ones
- [ ] Input detection still works correctly
- [ ] No console errors
- [ ] Performance is acceptable
- [ ] Help modal shows all shortcuts
- [ ] LocalStorage persistence works
- [ ] No memory leaks

---

## Future Testing Scenarios

1. **Touch/Mobile Testing:** Test on touch devices
2. **Screen Reader Testing:** NVDA/JAWS with shortcuts
3. **VoiceOver Testing:** MacOS accessibility
4. **Multiple Window Testing:** Same site in multiple tabs
5. **Offline Testing:** Service worker + shortcuts
6. **PWA Testing:** Installed app mode + shortcuts

---

**Document Version:** 1.0
**Last Updated:** [DATE]
**Maintained By:** [TEAM]
