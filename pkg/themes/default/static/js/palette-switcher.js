/**
 * Theme picker and dark/light toggle.
 *
 * The inline <head> script in base.html restores the visitor's choice
 * (color-mode, theme-palette-<mode>, theme-aesthetic, theme-fontpack) before
 * CSS paints, so
 * this file only wires up the UI. Palette blocks in palette.css use bare
 * [data-palette="name"] selectors, so every picker card is scoped with its own
 * data-palette attribute and previews the theme with its real colors.
 */
(function() {
  'use strict';

  const MODE_KEY = 'color-mode';
  const LEGACY_MODE_KEY = 'theme';
  const PICK_PREFIX = 'theme-palette-';
  const AESTHETIC_KEY = 'theme-aesthetic';
  const FONT_KEY = 'theme-fontpack';
  // Stored in place of a palette name to follow the seasonal calendar.
  const SEASONAL = 'seasonal';
  // Keys written unconditionally by the previous switcher; they do not reflect
  // an explicit visitor choice, so they are discarded instead of migrated.
  const LEGACY_KEYS = ['selected-palette', 'selected-family', 'selected-aesthetic'];

  const root = document.documentElement;

  const store = {
    get(key) {
      try { return localStorage.getItem(key); } catch (_) { return null; }
    },
    set(key, value) {
      try { localStorage.setItem(key, value); } catch (_) { /* ignore */ }
    },
    remove(key) {
      try { localStorage.removeItem(key); } catch (_) { /* ignore */ }
    },
  };

  function readCSSString(prop) {
    let value = getComputedStyle(root).getPropertyValue(prop).trim();
    if ((value.startsWith('"') && value.endsWith('"')) || (value.startsWith("'") && value.endsWith("'"))) {
      value = value.slice(1, -1);
    }
    return value;
  }

  function readCSSJSON(prop) {
    const raw = readCSSString(prop).replace(/\\'/g, "'");
    if (!raw || raw === 'none') return [];
    try {
      const parsed = JSON.parse(raw);
      return Array.isArray(parsed) ? parsed : [];
    } catch (e) {
      console.warn('[theme-picker] Failed to parse ' + prop, e);
      return [];
    }
  }

  // ---------------------------------------------------------------------------
  // Data
  // ---------------------------------------------------------------------------

  let manifestCache = null;
  let manifestIndex = null;
  let manifestHasFullStyles = false;

  function getManifest() {
    const fullStyles = !!window.__markataPaletteStylesReady;
    if (manifestCache && manifestCache.length && manifestHasFullStyles === fullStyles) return manifestCache;
    manifestCache = readCSSJSON('--palette-manifest');
    manifestHasFullStyles = fullStyles;
    manifestIndex = new Map(manifestCache.map((entry) => [entry.name, entry]));
    return manifestCache;
  }

  function findEntry(name) {
    getManifest();
    return (name && manifestIndex && manifestIndex.get(name)) || null;
  }

  function themesForMode(mode) {
    return getManifest().filter((entry) => entry.variant === mode);
  }

  function isPickerEnabled() {
    const flag = readCSSString('--palette-switcher-enabled');
    return (flag === '1' || flag === 'true') && getManifest().length > 0;
  }

  function getFallbackMode() {
    return window.__markataThemeFallbackMode === 'light' ? 'light' : 'dark';
  }

  function getDefaults() {
    const inline = window.__markataThemeDefaults || {};
    return {
      light: inline.light || readCSSString('--palette-light'),
      dark: inline.dark || readCSSString('--palette-dark'),
    };
  }

  function otherMode(mode) {
    return mode === 'dark' ? 'light' : 'dark';
  }

  function getColorMode() {
    const stored = store.get(MODE_KEY) || store.get(LEGACY_MODE_KEY);
    if (stored === 'light' || stored === 'dark') return stored;
    if (root.dataset.theme === 'light' || root.dataset.theme === 'dark') return root.dataset.theme;
    return getFallbackMode();
  }

  function paletteForMode(mode) {
    if (isPickerEnabled()) {
      if (isSeasonal(mode)) {
        const seasonal = seasonalPick(mode);
        if (seasonal && findEntry(seasonal.name)) return seasonal.name;
      }
      const picked = store.get(PICK_PREFIX + mode);
      if (picked && findEntry(picked)) return picked;
    }
    return getDefaults()[mode] || '';
  }

  // ---------------------------------------------------------------------------
  // Seasonal
  // ---------------------------------------------------------------------------

  function hasSeasonal() {
    return typeof window.__markataSeasonalPick === 'function' && !!window.__markataSeasonal;
  }

  function seasonalIsDefault() {
    return hasSeasonal() && window.__markataSeasonalDefault === true;
  }

  /** Today's seasonal palette for mode, as { label, name }, or null. */
  function seasonalPick(mode) {
    if (!hasSeasonal()) return null;
    try {
      return window.__markataSeasonalPick(mode || getColorMode());
    } catch (_) {
      return null;
    }
  }

  function isSeasonal(mode) {
    if (!hasSeasonal()) return false;
    const picked = store.get(PICK_PREFIX + (mode || getColorMode()));
    return picked === SEASONAL || (!picked && seasonalIsDefault());
  }

  function selectSeasonal() {
    endPreview();
    const mode = getColorMode();
    const pick = seasonalPick(mode);
    if (!pick || !findEntry(pick.name)) return;
    if (seasonalIsDefault()) {
      store.remove(PICK_PREFIX + 'light');
      store.remove(PICK_PREFIX + 'dark');
    } else {
      store.set(PICK_PREFIX + 'light', SEASONAL);
      store.set(PICK_PREFIX + 'dark', SEASONAL);
    }
    applyPalette(pick.name);
  }

  function labelFor(name) {
    const entry = findEntry(name);
    return entry ? entry.displayName : name;
  }

  // ---------------------------------------------------------------------------
  // Applying state
  // ---------------------------------------------------------------------------

  // Swap colors in a single frame; per-element color transitions would make the
  // page smear through intermediate states at different speeds.
  let switchFrame = 0;
  function withoutTransitions(fn) {
    root.classList.add('theme-switching');
    fn();
    cancelAnimationFrame(switchFrame);
    switchFrame = requestAnimationFrame(() => {
      switchFrame = requestAnimationFrame(() => root.classList.remove('theme-switching'));
    });
  }

  function applyPalette(name) {
    if (!name) return;
    loadPaletteStyles(name);
    withoutTransitions(() => {
      root.dataset.palette = name;
    });
    syncUI();
    window.dispatchEvent(new CustomEvent('palette-change', {
      detail: { palette: name, family: labelFor(name) },
    }));
  }

  function loadPaletteStyles(name) {
    const defaults = getDefaults();
    if (!window.__markataLoadPaletteStyles || name === defaults.light || name === defaults.dark) return;
    window.__markataLoadPaletteStyles().catch((error) => showNotification(error.message));
  }

  function setModeAttributes(mode) {
    root.dataset.theme = mode;
    root.classList.toggle('dark', mode === 'dark');
  }

  // ---------------------------------------------------------------------------
  // Hover preview
  // ---------------------------------------------------------------------------

  // Hovering a card or a calendar period shows it on the page without storing
  // anything. The real attributes are saved on the first preview and put back
  // when the pointer leaves, or before a choice is committed.
  let previewSaved = null;
  let previewFrame = 0;
  let previewPending = null;
  let previewRestoreTimer = 0;

  function snapshotAttributes() {
    return {
      palette: root.dataset.palette,
      theme: root.dataset.theme,
      dark: root.classList.contains('dark'),
      aesthetic: root.dataset.aesthetic,
      fontpack: root.dataset.fontpack,
    };
  }

  function setOrRemove(key, value) {
    if (value === undefined) delete root.dataset[key];
    else root.dataset[key] = value;
  }

  function writePreview(target) {
    withoutTransitions(() => {
      const saved = previewSaved;
      setOrRemove('palette', saved.palette);
      setOrRemove('theme', saved.theme);
      root.classList.toggle('dark', saved.dark);
      setOrRemove('aesthetic', saved.aesthetic);
      setOrRemove('fontpack', saved.fontpack);
      if (target.palette) {
        root.dataset.palette = target.palette;
        if (target.mode) setModeAttributes(target.mode);
      }
      if (target.aesthetic) root.dataset.aesthetic = target.aesthetic;
      if (target.fontpack) root.dataset.fontpack = target.fontpack;
    });
  }

  // Say what the page is showing: the top label names the previewed choice,
  // and in the schedule the hovered period and its twin (strip segment or
  // list row) are highlighted with a readout of its dates and length.
  function showPreviewInfo(target) {
    if (!picker) return;
    if (picker.schedule) {
      picker.schedule.querySelectorAll('.is-linked').forEach((el) => el.classList.remove('is-linked'));
      if (target && target.key) {
        picker.schedule.querySelectorAll('[data-period-key="' + CSS.escape(target.key) + '"]')
          .forEach((el) => el.classList.add('is-linked'));
      }
      const readout = picker.schedule.querySelector('.ss-readout');
      if (readout) {
        readout.textContent = target && target.detail ? target.detail : readout.dataset.hint || '';
        readout.classList.toggle('is-active', !!(target && target.detail));
      }
    }
    if (!picker.current) return;
    if (target && target.label) {
      picker.current.textContent = 'Preview: ' + target.label;
      picker.current.classList.add('is-previewing');
    } else {
      picker.current.classList.remove('is-previewing');
      syncUI();
    }
  }

  function startPreview(target) {
    clearTimeout(previewRestoreTimer);
    showPreviewInfo(target);
    if (!previewSaved) previewSaved = snapshotAttributes();
    previewPending = target;
    // Sweeping across many cards costs at most one style change per frame.
    if (previewFrame) return;
    previewFrame = requestAnimationFrame(() => {
      previewFrame = 0;
      if (previewSaved && previewPending) writePreview(previewPending);
    });
  }

  function endPreview() {
    clearTimeout(previewRestoreTimer);
    cancelAnimationFrame(previewFrame);
    previewFrame = 0;
    previewPending = null;
    if (!previewSaved) return;
    const saved = previewSaved;
    previewSaved = null;
    withoutTransitions(() => {
      setOrRemove('palette', saved.palette);
      setOrRemove('theme', saved.theme);
      root.classList.toggle('dark', saved.dark);
      setOrRemove('aesthetic', saved.aesthetic);
      setOrRemove('fontpack', saved.fontpack);
    });
    showPreviewInfo(null);
  }

  // Leaving one card for the next should not flash the saved theme in between.
  function scheduleEndPreview() {
    clearTimeout(previewRestoreTimer);
    previewRestoreTimer = setTimeout(endPreview, 90);
  }

  function palettePreview(name) {
    const entry = findEntry(name);
    return entry ? { palette: entry.name, mode: entry.variant } : null;
  }

  function previewTargetFor(element) {
    if (!element || !element.closest) return null;
    const period = element.closest('[data-preview-palette]');
    if (period) {
      const target = palettePreview(period.dataset.previewPalette);
      if (target) {
        target.label = period.dataset.previewLabel || '';
        target.key = period.dataset.periodKey || '';
        target.detail = period.dataset.previewDetail || '';
      }
      return target;
    }
    const card = element.closest('.theme-card');
    if (!card) return null;
    const name = card.querySelector('.tc-label') || card.querySelector('.theme-card-name');
    const label = name ? name.textContent : '';
    let target;
    if (card.classList.contains('style-card')) target = { aesthetic: card.dataset.aesthetic };
    else if (card.classList.contains('font-card')) target = { fontpack: card.dataset.name };
    else target = palettePreview(card.dataset.palette);
    if (target) target.label = label;
    return target;
  }

  const PREVIEW_AREAS = '.theme-picker-grid, .ss-year, .ss-list';

  function bindHoverPreview(panel) {
    panel.addEventListener('pointerover', (event) => {
      if (event.pointerType !== 'mouse') return;
      const target = previewTargetFor(event.target);
      if (target) startPreview(target);
    });
    panel.addEventListener('pointerout', (event) => {
      if (event.pointerType !== 'mouse' || !previewSaved) return;
      const next = event.relatedTarget;
      if (next && panel.contains(next)) {
        if (previewTargetFor(next)) return;
        // Gaps between cards or calendar periods keep the last preview.
        const area = event.target.closest && event.target.closest(PREVIEW_AREAS);
        if (area && area.contains(next)) return;
      }
      scheduleEndPreview();
    });
  }

  function setColorMode(mode) {
    if (mode !== 'light' && mode !== 'dark') return;
    endPreview();
    store.set(MODE_KEY, mode);
    store.set(LEGACY_MODE_KEY, mode);
    const palette = paletteForMode(mode);
    loadPaletteStyles(palette);
    withoutTransitions(() => {
      setModeAttributes(mode);
      if (palette) root.dataset.palette = palette;
    });
    syncUI();
    window.dispatchEvent(new CustomEvent('color-mode-change', { detail: { mode: mode } }));
  }

  function toggleColorMode() {
    const next = otherMode(getColorMode());
    setColorMode(next);
    return next;
  }

  /**
   * Select a theme and remember it (and its counterpart) per color mode, so
   * toggling dark/light returns to the matching half of the same family.
   */
  function selectTheme(name) {
    const entry = findEntry(name);
    if (!entry) return;
    endPreview();
    const mode = entry.variant;
    const other = otherMode(mode);
    const defaults = getDefaults();
    // With a seasonal site default, a concrete pick must be stored explicitly.
    const isDefaultPair = !seasonalIsDefault() && name === defaults[mode] && (!entry.counterpart || entry.counterpart === defaults[other]);

    if (isDefaultPair) {
      store.remove(PICK_PREFIX + mode);
      store.remove(PICK_PREFIX + other);
    } else {
      store.set(PICK_PREFIX + mode, name);
      if (entry.counterpart) {
        store.set(PICK_PREFIX + other, entry.counterpart);
      } else {
        store.remove(PICK_PREFIX + other);
      }
    }

    if (getColorMode() !== mode) {
      store.set(MODE_KEY, mode);
      store.set(LEGACY_MODE_KEY, mode);
      withoutTransitions(() => setModeAttributes(mode));
    }
    applyPalette(name);
  }

  function resetTheme() {
    endPreview();
    store.remove(PICK_PREFIX + 'light');
    store.remove(PICK_PREFIX + 'dark');
    store.remove(AESTHETIC_KEY);
    const defaultAesthetic = getDefaultAesthetic();
    if (defaultAesthetic) applyAesthetic(defaultAesthetic);
    if (store.get(FONT_KEY)) {
      store.remove(FONT_KEY);
      applyFontpack(pageFontpack || getDefaultFontpack());
    }
    applyPalette(paletteForMode(getColorMode()));
  }

  function cycleTheme(direction) {
    const list = themesForMode(getColorMode());
    if (!list.length) return;
    const index = list.findIndex((entry) => entry.name === root.dataset.palette);
    const step = direction === 'prev' ? -1 : 1;
    const next = list[(index + step + list.length) % list.length];
    selectTheme(next.name);
    showNotification(next.displayName);
  }

  function randomTheme(candidates) {
    const list = (candidates && candidates.length ? candidates : themesForMode(getColorMode()))
      .filter((entry) => entry.name !== root.dataset.palette);
    if (!list.length) return null;
    const pick = list[Math.floor(Math.random() * list.length)];
    selectTheme(pick.name);
    return pick;
  }

  // ---------------------------------------------------------------------------
  // Aesthetic
  // ---------------------------------------------------------------------------

  function getAestheticManifest() {
    return readCSSJSON('--aesthetic-manifest');
  }

  function getDefaultAesthetic() {
    const names = getAestheticManifest().map((item) => item.name);
    const configured = readCSSString('--aesthetic-default');
    if (configured && names.includes(configured)) return configured;
    return names[0] || '';
  }

  function applyAesthetic(name) {
    root.dataset.aesthetic = name;
    syncUI();
    window.dispatchEvent(new CustomEvent('aesthetic-change', { detail: { aesthetic: name } }));
  }

  function getAesthetic() {
    const names = getAestheticManifest().map((item) => item.name);
    const stored = store.get(AESTHETIC_KEY);
    if (stored && names.includes(stored)) return stored;
    if (root.dataset.aesthetic && names.includes(root.dataset.aesthetic)) return root.dataset.aesthetic;
    return getDefaultAesthetic();
  }

  function setAesthetic(name) {
    const names = getAestheticManifest().map((item) => item.name);
    if (!names.includes(name)) return;
    endPreview();
    if (name === getDefaultAesthetic()) {
      store.remove(AESTHETIC_KEY);
    } else {
      store.set(AESTHETIC_KEY, name);
    }
    applyAesthetic(name);
  }

  function cycleAesthetic(direction) {
    const items = getAestheticManifest();
    if (!items.length) return;
    const names = items.map((item) => item.name);
    const index = names.indexOf(getAesthetic());
    const step = direction === 'prev' ? -1 : 1;
    const next = items[(index + step + items.length) % items.length];
    setAesthetic(next.name);
    showNotification('Style: ' + (next.displayName || next.name));
  }

  // ---------------------------------------------------------------------------
  // Fonts
  // ---------------------------------------------------------------------------

  // The fontpack the page shipped with (a post may override the site pack).
  // The head script records it before applying a visitor's stored choice.
  let pageFontpack = window.__markataPageFontpack || root.dataset.fontpack || '';

  function getFontManifest() {
    return readCSSJSON('--fontpack-manifest');
  }

  function getDefaultFontpack() {
    return readCSSString('--fontpack-default') || pageFontpack;
  }

  function fontpackNames() {
    return getFontManifest().map((item) => item.name);
  }

  function getFontpack() {
    const names = fontpackNames();
    const stored = store.get(FONT_KEY);
    if (stored && names.includes(stored)) return stored;
    if (root.dataset.fontpack && names.includes(root.dataset.fontpack)) return root.dataset.fontpack;
    return getDefaultFontpack();
  }

  function applyFontpack(name) {
    if (!name) return;
    root.dataset.fontpack = name;
    syncUI();
    window.dispatchEvent(new CustomEvent('fontpack-change', { detail: { fontpack: name } }));
  }

  function setFontpack(name) {
    if (!fontpackNames().includes(name)) return;
    endPreview();
    if (name === pageFontpack) {
      store.remove(FONT_KEY);
    } else {
      store.set(FONT_KEY, name);
    }
    applyFontpack(name);
  }

  function cycleFont(direction) {
    const items = getFontManifest();
    if (items.length < 2) return;
    const names = items.map((item) => item.name);
    const index = names.indexOf(getFontpack());
    const step = direction === 'prev' ? -1 : 1;
    const next = items[(index + step + items.length) % items.length];
    setFontpack(next.name);
    showNotification('Font: ' + (next.displayName || next.name));
  }

  function describeFontpack(item) {
    const parts = [];
    if (item.heading) parts.push('Headings: ' + item.heading);
    if (item.body) parts.push('Body: ' + item.body);
    if (item.code) parts.push('Code: ' + item.code);
    return parts.join(' \u00b7 ');
  }

  // ---------------------------------------------------------------------------
  // Bake into config (markata-go serve only)
  // ---------------------------------------------------------------------------

  /** The dev-server bake endpoint; only `markata-go serve` injects it. */
  function bakeEndpoint() {
    const endpoint = window.__markataThemeBakeEndpoint;
    return typeof endpoint === 'string' && endpoint ? endpoint : '';
  }

  /**
   * The visitor's current choices as [markata-go.theme] settings. Color mode
   * and text size are visitor preferences, so they are only included when the
   * visitor explicitly picked them in this browser.
   */
  function getBakeSettings() {
    const mode = getColorMode();
    const light = paletteForMode('light');
    const dark = paletteForMode('dark');
    const settings = {};
    if (store.get(MODE_KEY) || store.get(LEGACY_MODE_KEY)) settings.fallback_mode = mode;
    if (isSeasonal(mode)) {
      // Seasonal picks change by date; the site palettes stay as the fallback.
      settings.seasonal = true;
    } else {
      const current = root.dataset.palette || (mode === 'dark' ? dark : light);
      if (current) settings.palette = current;
      if (light && !isSeasonal('light')) settings.palette_light = light;
      if (dark && !isSeasonal('dark')) settings.palette_dark = dark;
    }
    const aesthetic = getAestheticManifest().length ? getAesthetic() : root.dataset.aesthetic;
    if (aesthetic) settings.aesthetic = aesthetic;
    const fontpack = getFontManifest().length ? getFontpack() : root.dataset.fontpack;
    if (fontpack) settings.fontpack = fontpack;
    if (store.get(TEXT_SIZE_KEY) && root.dataset.textSize) settings.text_size = root.dataset.textSize;
    return settings;
  }

  function bakeRequest(method, body) {
    const endpoint = bakeEndpoint();
    if (!endpoint) return Promise.reject(new Error('Bake is only available in markata-go serve'));
    const init = { method: method, headers: { Accept: 'application/json' }, credentials: 'same-origin' };
    if (body) {
      init.headers['Content-Type'] = 'application/json';
      init.body = JSON.stringify(body);
    }
    return fetch(endpoint, init).then((response) => response.json().catch(() => ({})).then((data) => {
      if (!response.ok || data.error) throw new Error(data.error || ('HTTP ' + response.status));
      return data;
    }));
  }

  // The bake control is an icon button; idle it shows only the icon, while
  // the confirm/progress/result states expand an inline label.
  function setBakeLabel(button, text, state) {
    if (!button) return;
    const label = button.querySelector('[data-picker-bake-label]');
    const idle = text === 'Bake' && !state;
    if (label) {
      label.textContent = idle ? '' : text;
    } else {
      button.textContent = text;
    }
    button.classList.toggle('has-label', !idle);
    button.setAttribute('aria-label', idle ? 'Bake: write these choices into your site config' : text);
    button.classList.toggle('is-armed', state === 'armed');
    button.classList.toggle('is-done', state === 'done');
    button.classList.toggle('is-error', state === 'error');
  }

  function disarmBake(button) {
    if (!button) return;
    clearTimeout(button._bakeTimeout);
    delete button.dataset.armed;
    setBakeLabel(button, 'Bake', '');
  }

  /** Write the current choices into the site config and clear stored picks. */
  function bakeConfig() {
    return bakeRequest('POST', getBakeSettings()).then((result) => {
      // The config now matches these picks; drop the per-browser overrides so
      // the rebuilt site's defaults take over after live reload.
      store.remove(PICK_PREFIX + 'light');
      store.remove(PICK_PREFIX + 'dark');
      store.remove(AESTHETIC_KEY);
      store.remove(FONT_KEY);
      if (result.keys && result.keys.indexOf('text_size') >= 0) store.remove(TEXT_SIZE_KEY);
      return result;
    });
  }

  /** First click names the target file; a second click within 4s bakes. */
  function onBakeClick(button) {
    if (button.dataset.busy) return;
    if (!button.dataset.armed) {
      button.dataset.busy = '1';
      bakeRequest('GET').then((info) => {
        const target = info.target || 'config';
        if (info.target_kind === 'global') {
          setBakeLabel(button, 'Bake unavailable', 'error');
          button.title = target + ' is your global config; add a site config to bake into.';
          showNotification('Bake unavailable: this site has no config file of its own (' + target + ' is global)');
          button._bakeTimeout = setTimeout(() => disarmBake(button), 4000);
          return;
        }
        button.dataset.armed = '1';
        const override = info.target_kind === 'override';
        setBakeLabel(button, 'Bake into ' + target + (override ? ' (override)' : '') + '?', 'armed');
        button.title = 'Click again to write these choices to ' + (info.path || target) +
          (override ? '\nThis is a --merge-config override file, not the base config.' : '') +
          (info.warnings && info.warnings.length ? '\n' + info.warnings.join('\n') : '');
        clearTimeout(button._bakeTimeout);
        button._bakeTimeout = setTimeout(() => disarmBake(button), 4000);
      }, (error) => {
        setBakeLabel(button, 'Bake failed', 'error');
        button.title = error.message;
        showNotification('Bake failed: ' + error.message);
      }).finally(() => { delete button.dataset.busy; });
      return;
    }
    clearTimeout(button._bakeTimeout);
    delete button.dataset.armed;
    button.dataset.busy = '1';
    setBakeLabel(button, 'Baking\u2026', '');
    bakeConfig().then((result) => {
      setBakeLabel(button, 'Baked', 'done');
      button.title = 'Saved to ' + (result.path || result.target);
      const keys = result.keys && result.keys.length ? ' (' + result.keys.join(', ') + ')' : '';
      showNotification('Baked' + keys + ' into ' + result.target + ' \u2014 rebuilding');
      if (result.warnings && result.warnings.length) console.warn('[theme-picker] ' + result.warnings.join('; '));
    }, (error) => {
      setBakeLabel(button, 'Bake failed', 'error');
      button.title = error.message;
      showNotification('Bake failed: ' + error.message);
    }).finally(() => {
      delete button.dataset.busy;
      button._bakeTimeout = setTimeout(() => disarmBake(button), 2500);
    });
  }

  // ---------------------------------------------------------------------------
  // Validation of restored state
  // ---------------------------------------------------------------------------

  function validateRestoredState() {
    LEGACY_KEYS.forEach((key) => store.remove(key));
    if (!isPickerEnabled()) return;

    const mode = getColorMode();
    const current = root.dataset.palette;
    const defaults = getDefaults();
    if (current && !findEntry(current) && current !== defaults[mode]) {
      // A stored palette no longer exists on this site.
      store.remove(PICK_PREFIX + 'light');
      store.remove(PICK_PREFIX + 'dark');
      withoutTransitions(() => { root.dataset.palette = defaults[mode]; });
    }

    const storedFont = store.get(FONT_KEY);
    if (storedFont && getFontManifest().length && !fontpackNames().includes(storedFont)) {
      store.remove(FONT_KEY);
      if (pageFontpack) root.dataset.fontpack = pageFontpack;
    }

    const storedAesthetic = store.get(AESTHETIC_KEY);
    if (storedAesthetic) {
      const names = getAestheticManifest().map((item) => item.name);
      if (!names.includes(storedAesthetic)) {
        store.remove(AESTHETIC_KEY);
        const fallback = getDefaultAesthetic();
        if (fallback) root.dataset.aesthetic = fallback;
      }
    }
  }

  // ---------------------------------------------------------------------------
  // UI
  // ---------------------------------------------------------------------------

  let picker = null; // state for the currently mounted picker

  const TABS = ['colors', 'style', 'font'];
  const TAB_KEY = 'theme-picker-tab';
  const TEXT_SIZE_KEY = 'text-size';
  const STYLE_NOTES = {
    minimal: 'Quiet, flat surfaces',
    balanced: 'Soft corners, light depth',
    elevated: 'Round cards that float',
    precision: 'Crisp, thin outlines',
    brutal: 'Hard edges, bold shadows',
  };

  function getTextSize() {
    if (window.MarkataTextSize) return window.MarkataTextSize.get();
    return root.dataset.textSize || 'large';
  }

  function setTextSize(size) {
    if (window.MarkataTextSize) {
      window.MarkataTextSize.set(size);
    } else {
      root.dataset.textSize = size;
      store.set(TEXT_SIZE_KEY, size);
    }
    syncUI();
    window.dispatchEvent(new CustomEvent('text-size-change', { detail: { size: size } }));
  }

  function tabGrid(tab) {
    if (!picker) return null;
    if (tab === 'style') return picker.styleGrid;
    if (tab === 'font') return picker.fontGrid;
    return picker.grid;
  }

  function currentName(tab) {
    if (tab === 'style') return getAesthetic();
    if (tab === 'font') return getFontpack();
    return isSeasonal() ? SEASONAL : root.dataset.palette;
  }

  function selectInTab(tab, name) {
    if (tab === 'style') setAesthetic(name);
    else if (tab === 'font') setFontpack(name);
    else if (name === SEASONAL) selectSeasonal();
    else selectTheme(name);
  }

  function currentLabel(tab) {
    const name = currentName(tab);
    if (tab === 'style') {
      const item = getAestheticManifest().find((entry) => entry.name === name);
      return item ? (item.displayName || item.name) : name;
    }
    if (tab === 'font') {
      const item = getFontManifest().find((entry) => entry.name === name);
      return item ? (item.displayName || item.name) : name;
    }
    if (name === SEASONAL) {
      const pick = seasonalPick();
      return 'Seasonal' + (pick ? ': ' + pick.label : '');
    }
    return labelFor(name);
  }

  function syncGrid(grid, current) {
    if (!grid) return;
    grid.querySelectorAll('.theme-card').forEach((card) => {
      const active = card.dataset.name === current;
      card.classList.toggle('is-current', active);
      card.setAttribute('aria-selected', String(active));
    });
  }

  function syncUI() {
    const mode = getColorMode();
    const label = mode === 'dark' ? 'Dark' : 'Light';
    document.querySelectorAll('.palette-mode-toggle').forEach((toggle) => {
      toggle.title = label + ' mode - click to toggle (\\)';
      toggle.setAttribute('aria-label', 'Current: ' + label + ' mode. Click to switch.');
    });

    if (!picker) return;
    picker.toggle.title = 'Theme: ' + currentLabel('colors') + ' (t)';
    picker.modeButtons.forEach((button) => {
      button.setAttribute('aria-pressed', String(button.dataset.pickerMode === picker.viewMode));
    });
    syncGrid(picker.grid, currentName('colors'));
    syncGrid(picker.styleGrid, getAesthetic());
    syncGrid(picker.fontGrid, getFontpack());
    const size = getTextSize();
    picker.sizeButtons.forEach((button) => {
      button.setAttribute('aria-checked', String(button.dataset.pickerSize === size));
    });
    if (picker.current) picker.current.textContent = currentLabel(picker.activeTab);
    const noun = picker.activeTab === 'colors' ? 'theme' : picker.activeTab;
    picker.stepButtons.forEach((button) => {
      const verb = button.dataset.pickerStep === 'prev' ? 'Previous ' : 'Next ';
      button.setAttribute('aria-label', verb + noun);
      button.title = verb + noun;
    });
  }

  function cardShell(className, name) {
    const card = document.createElement('button');
    card.type = 'button';
    card.className = 'theme-card ' + className;
    card.id = className + '-' + name.replace(/[^a-z0-9-]/gi, '_');
    card.setAttribute('role', 'option');
    card.setAttribute('aria-selected', 'false');
    card.tabIndex = -1;
    card.dataset.name = name;
    return card;
  }

  function buildCard(entry) {
    const card = cardShell('palette-card', entry.name);
    card.dataset.palette = entry.name;
    card.title = entry.displayName + (entry.derived ? ' (' + entry.variant + ' variant)' : '');
    card.innerHTML =
      '<span class="theme-card-preview" aria-hidden="true">' +
        '<span class="tc-title"></span>' +
        '<span class="tc-line"></span>' +
        '<span class="tc-line tc-line--short"></span>' +
        '<span class="tc-code"><i class="tc-k"></i><i class="tc-f"></i><i class="tc-s"></i></span>' +
        '<span class="tc-dots"><i class="tc-accent"></i><i class="tc-success"></i><i class="tc-warning"></i><i class="tc-error"></i></span>' +
      '</span>' +
      '<span class="theme-card-name"></span>';
    card.querySelector('.theme-card-name').textContent = entry.displayName;
    return card;
  }

  // The seasonal card previews today's pick for the mode being browsed.
  function buildSeasonalCard(mode) {
    const pick = seasonalPick(mode);
    if (!pick || !findEntry(pick.name)) return null;
    const card = buildCard({ name: pick.name, displayName: 'Seasonal' });
    card.id = 'palette-card-seasonal';
    card.dataset.name = SEASONAL;
    card.classList.add('seasonal-card');
    card.title = 'Seasonal - follows the seasons and celebrates holidays. Now: ' + pick.label + ' (' + labelFor(pick.name) + '). The calendar button shows the schedule.';
    const name = card.querySelector('.theme-card-name');
    name.textContent = '';
    const label = document.createElement('span');
    label.className = 'tc-label';
    label.textContent = 'Seasonal';
    const note = document.createElement('small');
    note.className = 'tc-note';
    note.textContent = pick.label;
    name.append(label, note);
    return card;
  }

  // Seasonal schedule
  //
  // Periods are derived by asking the same pick function the head script uses
  // for every day, so the schedule always matches what the site will do.
  function seasonalPeriods(mode, start, days) {
    const periods = [];
    if (!hasSeasonal()) return periods;
    for (let i = 0; i < days; i++) {
      const date = new Date(start.getFullYear(), start.getMonth(), start.getDate() + i, 12);
      let pick = null;
      try {
        pick = window.__markataSeasonalPick(mode, date);
      } catch (_) {
        return periods;
      }
      if (!pick || !pick.name) continue;
      const last = periods[periods.length - 1];
      if (last && last.label === pick.label && last.name === pick.name) {
        last.end = date;
        last.days++;
      } else {
        periods.push({ label: pick.label, name: pick.name, start: date, end: date, days: 1 });
      }
    }
    return periods;
  }

  function isSeasonLabel(label) {
    const s = window.__markataSeasonal;
    return !!(s && s.s && s.s.some((season) => season.l === label));
  }

  function formatDay(date) {
    return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
  }

  function formatDuration(days) {
    if (days < 14) return days + (days === 1 ? ' day' : ' days');
    const weeks = days / 7;
    return (Number.isInteger(weeks) ? '' : '~') + Math.round(weeks) + ' weeks';
  }

  // Full span of the season active on date (holidays inside it included).
  function seasonBounds(label, date) {
    const s = window.__markataSeasonal;
    if (!s || !s.s) return null;
    const seasons = s.s.slice().sort((a, b) => a.f - b.f);
    const index = seasons.findIndex((season) => season.l === label);
    if (index === -1) return null;
    const at = (f, year) => new Date(year, Math.floor(f / 100) - 1, f % 100, 12);
    let start = at(seasons[index].f, date.getFullYear());
    if (start > date) start = at(seasons[index].f, date.getFullYear() - 1);
    const next = seasons[(index + 1) % seasons.length];
    let end = at(next.f, start.getFullYear());
    if (end <= start) end = at(next.f, start.getFullYear() + 1);
    return { start: start, end: end, days: Math.round((end - start) / 864e5) };
  }

  function periodKey(period) {
    if (isSeasonLabel(period.label)) {
      const bounds = seasonBounds(period.label, period.start);
      if (bounds) return period.label + '|' + bounds.start.toDateString();
    }
    return period.label + '|' + period.start.toDateString();
  }

  function withPalette(text, period) {
    const name = labelFor(period.name);
    return name && name.toLowerCase() !== period.label.toLowerCase() ? text + ' · ' + name : text;
  }

  function markPeriod(el, period, detail) {
    el.dataset.previewPalette = period.name;
    el.dataset.periodKey = periodKey(period);
    el.dataset.previewLabel = period.label;
    el.dataset.previewDetail = detail;
  }

  function formatRange(period) {
    return period.days === 1 ? formatDay(period.start) : formatDay(period.start) + ' – ' + formatDay(period.end);
  }

  function scheduleSwatch(name) {
    const swatch = document.createElement('span');
    swatch.className = 'ss-swatch';
    swatch.dataset.palette = name;
    swatch.setAttribute('aria-hidden', 'true');
    swatch.innerHTML = '<i class="tc-accent"></i><i class="tc-success"></i><i class="tc-warning"></i>';
    return swatch;
  }

  function renderSchedule() {
    if (!picker || !picker.schedule) return;
    const mode = picker.viewMode;
    const view = picker.schedule;
    const today = new Date();
    const lead = (window.__markataSeasonal && window.__markataSeasonal.lead) || 0;

    const intro = document.createElement('div');
    intro.className = 'ss-intro';
    const text = document.createElement('p');
    text.textContent = 'Seasonal changes the colors with the seasons (northern hemisphere) and switches to a holiday theme ' +
      (lead ? lead + ' days before' : 'on') + ' each holiday. Showing ' + mode + ' themes.';
    const use = document.createElement('button');
    use.type = 'button';
    use.className = 'ss-use';
    use.dataset.scheduleUse = 'true';
    const active = isSeasonal(mode);
    use.textContent = active ? 'Using Seasonal' : 'Use Seasonal';
    use.disabled = active;
    intro.append(text, use);

    // Year strip for the current calendar year.
    const yearStart = new Date(today.getFullYear(), 0, 1);
    const yearDays = Math.round((new Date(today.getFullYear() + 1, 0, 1) - yearStart) / 864e5);
    const strip = document.createElement('div');
    strip.className = 'ss-year';
    strip.setAttribute('role', 'img');
    const yearPeriods = seasonalPeriods(mode, yearStart, yearDays);
    strip.setAttribute('aria-label', today.getFullYear() + ' seasonal calendar: ' +
      yearPeriods.map((p) => p.label + ' ' + formatRange(p)).join(', '));
    yearPeriods.forEach((period) => {
      const segment = document.createElement('span');
      segment.className = 'ss-seg' + (isSeasonLabel(period.label) ? '' : ' ss-seg--holiday');
      segment.dataset.palette = period.name;
      const detail = withPalette(period.label + ' · ' + formatRange(period) + ' · ' + formatDuration(period.days), period);
      markPeriod(segment, period, detail);
      segment.style.flexGrow = String(period.days);
      segment.title = detail;
      strip.appendChild(segment);
    });
    const dayOfYear = Math.floor((new Date(today.getFullYear(), today.getMonth(), today.getDate()) - yearStart) / 864e5);
    const marker = document.createElement('span');
    marker.className = 'ss-today';
    marker.style.left = ((dayOfYear + 0.5) / yearDays * 100).toFixed(2) + '%';
    marker.title = 'Today';
    strip.appendChild(marker);

    const months = document.createElement('div');
    months.className = 'ss-months';
    months.setAttribute('aria-hidden', 'true');
    for (let m = 0; m < 12; m++) {
      const label = document.createElement('span');
      label.textContent = new Date(today.getFullYear(), m, 1).toLocaleDateString(undefined, { month: 'narrow' });
      months.appendChild(label);
    }

    // Upcoming changes over the next year, starting with today's period.
    const heading = document.createElement('h3');
    heading.className = 'ss-heading';
    heading.textContent = 'Coming up';
    const list = document.createElement('ol');
    list.className = 'ss-list';
    // A season that resumes after a holiday is not a new change worth listing.
    const seen = new Set();
    seasonalPeriods(mode, today, 366).forEach((period, index) => {
      const season = isSeasonLabel(period.label);
      if (index > 0 && season && seen.has(period.label)) return;
      seen.add(period.label);
      const item = document.createElement('li');
      item.className = 'ss-item' + (index === 0 ? ' is-now' : '');
      const body = document.createElement('span');
      body.className = 'ss-item-body';
      const label = document.createElement('strong');
      label.textContent = period.label;
      body.appendChild(label);
      const paletteName = labelFor(period.name);
      if (paletteName && paletteName.toLowerCase() !== period.label.toLowerCase()) {
        const palette = document.createElement('small');
        palette.textContent = paletteName;
        body.appendChild(palette);
      }
      const when = document.createElement('span');
      when.className = 'ss-when';
      const date = document.createElement('span');
      const length = document.createElement('small');
      let detail;
      if (index === 0) {
        const todayNoon = new Date(today.getFullYear(), today.getMonth(), today.getDate(), 12);
        const left = Math.round((period.end - todayNoon) / 864e5) + 1;
        date.textContent = 'Now · until ' + formatDay(period.end);
        length.textContent = left === 1 ? 'last day' : formatDuration(left) + ' left';
        detail = period.label + ' · now until ' + formatDay(period.end) + ' · ' + length.textContent;
      } else if (season) {
        const bounds = seasonBounds(period.label, period.start);
        date.textContent = 'From ' + formatDay(period.start);
        length.textContent = bounds ? formatDuration(bounds.days) + ' season' : '';
        detail = period.label + ' · from ' + formatDay(period.start) +
          (bounds ? ' · ' + formatDuration(bounds.days) + ', holidays included' : '');
      } else {
        date.textContent = formatRange(period);
        length.textContent = formatDuration(period.days);
        detail = period.label + ' · ' + formatRange(period) + ' · ' + formatDuration(period.days) +
          (lead ? ', starting ' + lead + ' days before' : '');
      }
      when.append(date, length);
      markPeriod(item, period, withPalette(detail, period));
      item.append(scheduleSwatch(period.name), body, when);
      list.appendChild(item);
    });

    const readout = document.createElement('p');
    readout.className = 'ss-readout';
    readout.setAttribute('aria-live', 'polite');
    readout.dataset.hint = window.matchMedia('(hover: hover)').matches
      ? 'Hover the calendar to preview a period on this page.'
      : 'Tap a band or row for its dates. The line marks today.';
    readout.textContent = readout.dataset.hint;

    view.replaceChildren(intro, strip, months, readout, heading, list);
    picker.scheduleMode = mode;
  }

  function setScheduleOpen(open) {
    if (!picker || !picker.schedule || !picker.scheduleButton) return;
    endPreview();
    open = !!open && hasSeasonal();
    picker.scheduleOpen = open;
    picker.scheduleButton.setAttribute('aria-pressed', String(open));
    picker.schedule.hidden = !open;
    picker.grid.hidden = open;
    if (open) {
      picker.empty.hidden = true;
      renderSchedule();
      picker.schedule.scrollTop = 0;
    } else {
      renderGrid();
    }
  }

  // Style cards carry their own data-aesthetic, so the surface tokens
  // (--radius-*, --surface-border, --surface-shadow) resolve per card.
  function buildStyleCard(item) {
    const card = cardShell('style-card', item.name);
    card.dataset.aesthetic = item.name;
    const label = item.displayName || item.name;
    const note = STYLE_NOTES[item.name] || '';
    card.title = label + (note ? ' - ' + note : '');
    card.innerHTML =
      '<span class="theme-card-preview style-card-preview" aria-hidden="true">' +
        '<span class="sc-panel">' +
          '<span class="tc-title"></span>' +
          '<span class="tc-line"></span>' +
          '<span class="tc-line tc-line--short"></span>' +
        '</span>' +
        '<span class="sc-row"><span class="sc-code"><i class="tc-k"></i><i class="tc-s"></i></span><span class="sc-button"></span></span>' +
      '</span>' +
      '<span class="theme-card-name"><span class="tc-label"></span><small class="tc-note"></small></span>';
    card.querySelector('.tc-label').textContent = label;
    card.querySelector('.tc-note').textContent = note;
    return card;
  }

  function buildFontCard(item) {
    const card = cardShell('font-card', item.name);
    const label = item.displayName || item.name;
    card.title = label + ' - ' + describeFontpack(item);
    card.innerHTML =
      '<span class="theme-card-preview font-card-preview" aria-hidden="true">' +
        '<span class="fc-heading">Aa Heading</span>' +
        '<span class="fc-body">The quick brown fox jumps over the lazy dog.</span>' +
        '<span class="fc-code">let x = 42;</span>' +
      '</span>' +
      '<span class="theme-card-name"><span class="tc-label"></span><small class="tc-note"></small></span>';
    card.querySelector('.tc-label').textContent = label;
    card.querySelector('.tc-note').textContent = [item.heading, item.body].filter(Boolean)
      .filter((value, index, list) => list.indexOf(value) === index).join(' + ');
    card._fontItem = item;
    return card;
  }

  // Font previews apply each pack's families only once the card scrolls into
  // view, so opening the Font tab does not download every font on the site.
  function applyFontPreview(card) {
    const item = card._fontItem;
    if (!item || card.dataset.fontLoaded === 'true') return;
    card.dataset.fontLoaded = 'true';
    const heading = card.querySelector('.fc-heading');
    const body = card.querySelector('.fc-body');
    const code = card.querySelector('.fc-code');
    if (item.headingFont) heading.style.fontFamily = item.headingFont;
    if (item.headingWeight) heading.style.fontWeight = String(item.headingWeight);
    if (item.bodyFont) body.style.fontFamily = item.bodyFont;
    if (item.codeFont) code.style.fontFamily = item.codeFont;
  }

  function observeFontCards() {
    const cards = Array.from(picker.fontGrid.querySelectorAll('.font-card'));
    if (!('IntersectionObserver' in window)) {
      cards.forEach(applyFontPreview);
      return;
    }
    if (picker.fontObserver) picker.fontObserver.disconnect();
    picker.fontObserver = new IntersectionObserver((entries, observer) => {
      entries.forEach((entry) => {
        if (!entry.isIntersecting) return;
        applyFontPreview(entry.target);
        observer.unobserve(entry.target);
      });
    }, { root: picker.fontGrid, rootMargin: '120px 0px' });
    cards.forEach((card) => picker.fontObserver.observe(card));
  }

  function renderGrid() {
    if (!picker) return;
    const { grid, empty } = picker;
    const query = picker.search.value.trim().toLowerCase();
    const list = themesForMode(picker.viewMode).filter((entry) => {
      if (!query) return true;
      return entry.displayName.toLowerCase().includes(query) || entry.name.includes(query);
    });

    if (picker.renderedMode !== picker.viewMode) {
      const fragment = document.createDocumentFragment();
      const seasonalCard = buildSeasonalCard(picker.viewMode);
      if (seasonalCard) fragment.appendChild(seasonalCard);
      themesForMode(picker.viewMode).forEach((entry) => fragment.appendChild(buildCard(entry)));
      grid.replaceChildren(fragment);
      picker.renderedMode = picker.viewMode;
    }

    const visible = new Set(list.map((entry) => entry.name));
    const seasonal = seasonalPick(picker.viewMode);
    if (seasonal && (!query || ('seasonal holiday ' + seasonal.label).toLowerCase().includes(query))) visible.add(SEASONAL);
    grid.querySelectorAll('.theme-card').forEach((card) => {
      card.hidden = !visible.has(card.dataset.name);
    });
    empty.hidden = visible.size > 0;
    picker.search.setAttribute('placeholder', 'Search ' + themesForMode(picker.viewMode).length + ' ' + picker.viewMode + ' themes');
    syncUI();
  }

  function renderStyleGrid() {
    const items = getAestheticManifest();
    picker.tabs.style.hidden = items.length < 2;
    if (items.length < 2 || picker.styleGrid.childElementCount === items.length) return;
    picker.styleGrid.replaceChildren(...items.map(buildStyleCard));
  }

  function renderFontGrid() {
    const items = getFontManifest();
    const hasSizes = picker.sizeButtons.length > 0;
    picker.tabs.font.hidden = items.length < 2 && !hasSizes;
    if (!picker.fontGrid) return;
    picker.fontGrid.hidden = items.length < 2;
    if (items.length < 2 || picker.fontGrid.childElementCount === items.length) return;
    picker.fontGrid.replaceChildren(...items.map(buildFontCard));
  }

  function availableTabs() {
    return TABS.filter((tab) => picker.tabs[tab] && !picker.tabs[tab].hidden);
  }

  function visibleCards(tab) {
    const grid = tabGrid(tab || (picker && picker.activeTab));
    return grid && !grid.hidden ? Array.from(grid.querySelectorAll('.theme-card:not([hidden])')) : [];
  }

  function gridColumns(cards) {
    if (cards.length < 2) return 1;
    const top = cards[0].offsetTop;
    let columns = 0;
    for (const card of cards) {
      if (card.offsetTop !== top) break;
      columns++;
    }
    return Math.max(1, columns);
  }

  function revealCard(card, center) {
    if (!card || !picker) return;
    const grid = card.parentElement;
    const cardTop = card.offsetTop - grid.offsetTop;
    const cardBottom = cardTop + card.offsetHeight;
    if (center) {
      grid.scrollTop = cardTop - (grid.clientHeight - card.offsetHeight) / 2;
    } else if (cardTop < grid.scrollTop + 4) {
      grid.scrollTop = cardTop - 8;
    } else if (cardBottom > grid.scrollTop + grid.clientHeight - 4) {
      grid.scrollTop = cardBottom - grid.clientHeight + 8;
    }
  }

  function focusCard(card) {
    if (!card || !picker) return;
    if (picker.activeTab === 'colors') picker.search.setAttribute('aria-activedescendant', card.id);
    const active = document.activeElement;
    if (active && active.classList && (active.classList.contains('theme-card') || active.classList.contains('theme-picker-step'))) {
      if (active.classList.contains('theme-card')) card.focus({ preventScroll: true });
    }
    revealCard(card, false);
  }

  function moveSelection(key) {
    const tab = picker.activeTab;
    const cards = visibleCards(tab);
    if (!cards.length) return;
    const columns = gridColumns(cards);
    let index = cards.findIndex((card) => card.dataset.name === currentName(tab));
    if (index === -1) {
      index = key === 'ArrowUp' || key === 'ArrowLeft' || key === 'End' ? cards.length : -1;
    }
    let next = index;
    switch (key) {
      case 'ArrowRight': next = index + 1; break;
      case 'ArrowLeft': next = index - 1; break;
      case 'ArrowDown': next = index === -1 ? 0 : index + columns; break;
      case 'ArrowUp': next = index - columns; break;
      case 'Home': next = 0; break;
      case 'End': next = cards.length - 1; break;
      case 'PageDown': next = index + columns * 3; break;
      case 'PageUp': next = index - columns * 3; break;
      default: return;
    }
    if (key === 'ArrowRight' || key === 'ArrowLeft') {
      next = (next + cards.length) % cards.length;
    } else {
      next = Math.min(cards.length - 1, Math.max(0, next));
    }
    const card = cards[next];
    selectInTab(tab, card.dataset.name);
    focusCard(card);
  }

  function prefersSearchFocus() {
    return window.matchMedia('(hover: hover) and (pointer: fine)').matches;
  }

  function focusActivePane(focusGrid) {
    const tab = picker.activeTab;
    const current = tabGrid(tab) && tabGrid(tab).querySelector('.theme-card.is-current');
    if (tab === 'colors' && !focusGrid && prefersSearchFocus()) {
      picker.search.focus({ preventScroll: true });
    } else if (current) {
      current.focus({ preventScroll: true });
    } else {
      const first = visibleCards(tab)[0];
      (first || picker.panel).focus({ preventScroll: true });
    }
  }

  function showTab(tab, options) {
    if (!picker) return;
    endPreview();
    const tabs = availableTabs();
    if (!tabs.includes(tab)) tab = tabs[0] || 'colors';
    picker.activeTab = tab;
    store.set(TAB_KEY, tab);
    TABS.forEach((name) => {
      const button = picker.tabs[name];
      const pane = picker.panes[name];
      const active = name === tab;
      if (button) {
        button.setAttribute('aria-selected', String(active));
        button.tabIndex = active ? 0 : -1;
      }
      if (pane) pane.hidden = !active;
    });
    picker.panel.dataset.activeTab = tab;
    if (tab === 'font' && picker.fontGrid && !picker.fontGrid.hidden) observeFontCards();
    syncUI();
    const current = tabGrid(tab) && tabGrid(tab).querySelector('.theme-card.is-current');
    if (current) revealCard(current, true);
    if (options && options.focus) focusActivePane(options.focusGrid);
  }

  function stepActive(direction) {
    if (!picker) return;
    if (picker.activeTab === 'colors' && picker.scheduleOpen) setScheduleOpen(false);
    moveSelection(direction === 'prev' ? 'ArrowLeft' : 'ArrowRight');
  }

  function openPicker(options) {
    if (!picker || !picker.panel.hidden) return;
    if (window.__markataLoadPaletteStyles && !window.__markataPaletteStylesReady) {
      picker.toggle.setAttribute('aria-busy', 'true');
      window.__markataLoadPaletteStyles().then(() => {
        if (picker) picker.toggle.removeAttribute('aria-busy');
        openPicker(options);
      }).catch((error) => {
        if (picker) picker.toggle.removeAttribute('aria-busy');
        showNotification(error.message);
      });
      return;
    }
    picker.viewMode = getColorMode();
    renderStyleGrid();
    renderFontGrid();
    picker.panel.hidden = false;
    picker.toggle.setAttribute('aria-expanded', 'true');
    picker.root.classList.add('is-open');
    renderGrid();
    const current = picker.grid.querySelector('.theme-card.is-current');
    if (current) picker.search.setAttribute('aria-activedescendant', current.id);
    const tab = (options && options.tab) || store.get(TAB_KEY) || 'colors';
    showTab(tab, { focus: true, focusGrid: options && options.focusGrid });
  }

  function closePicker(returnFocus) {
    if (!picker || picker.panel.hidden) return;
    endPreview();
    picker.panel.hidden = true;
    picker.toggle.setAttribute('aria-expanded', 'false');
    picker.root.classList.remove('is-open');
    if (picker.fontObserver) {
      picker.fontObserver.disconnect();
      picker.fontObserver = null;
    }
    if (picker.search.value) {
      picker.search.value = '';
      renderGrid();
    }
    if (picker.scheduleOpen) setScheduleOpen(false);
    if (returnFocus) picker.toggle.focus({ preventScroll: true });
  }

  function mountPicker(container) {
    const toggle = container.querySelector('.theme-picker-toggle');
    const panel = container.querySelector('.theme-picker-panel');
    if (!toggle || !panel) return;

    const tabs = {};
    const panes = {};
    TABS.forEach((tab) => {
      tabs[tab] = panel.querySelector('[data-picker-tab="' + tab + '"]');
      panes[tab] = panel.querySelector('[data-picker-pane="' + tab + '"]');
    });

    picker = {
      root: container,
      toggle: toggle,
      panel: panel,
      tabs: tabs,
      panes: panes,
      activeTab: 'colors',
      search: panel.querySelector('.theme-picker-search'),
      grid: panel.querySelector('.theme-picker-grid'),
      styleGrid: panel.querySelector('[data-picker-style-grid]'),
      fontGrid: panel.querySelector('[data-picker-font-grid]'),
      empty: panel.querySelector('.theme-picker-empty'),
      current: panel.querySelector('[data-picker-current]'),
      modeButtons: Array.from(panel.querySelectorAll('[data-picker-mode]')),
      sizeButtons: Array.from(panel.querySelectorAll('[data-picker-size]')),
      stepButtons: Array.from(panel.querySelectorAll('[data-picker-step]')),
      schedule: panel.querySelector('[data-picker-schedule-view]'),
      scheduleButton: panel.querySelector('[data-picker-schedule]'),
      scheduleOpen: false,
      viewMode: getColorMode(),
      renderedMode: null,
      fontObserver: null,
    };
    panel.tabIndex = -1;

    if (container.dataset.themePickerBound === 'true') {
      syncUI();
      return;
    }
    container.dataset.themePickerBound = 'true';
    bindHoverPreview(panel);

    toggle.addEventListener('click', (event) => {
      event.stopPropagation();
      if (panel.hidden) {
        openPicker();
      } else {
        closePicker(false);
      }
    });

    TABS.forEach((tab) => {
      if (tabs[tab]) tabs[tab].addEventListener('click', () => showTab(tab));
    });

    picker.stepButtons.forEach((button) => {
      button.addEventListener('click', () => stepActive(button.dataset.pickerStep));
    });

    picker.sizeButtons.forEach((button) => {
      button.addEventListener('click', () => setTextSize(button.dataset.pickerSize));
    });

    picker.search.addEventListener('input', () => {
      if (picker.scheduleOpen) setScheduleOpen(false);
      renderGrid();
      const first = visibleCards('colors')[0];
      if (first) revealCard(first, false);
    });

    panel.addEventListener('click', (event) => {
      const card = event.target.closest('.theme-card');
      if (!card || !panel.contains(card)) return;
      const pane = card.closest('[data-picker-pane]');
      const tab = pane ? pane.dataset.pickerPane : 'colors';
      selectInTab(tab, card.dataset.name);
      if (tab === 'colors') picker.search.setAttribute('aria-activedescendant', card.id);
    });

    picker.modeButtons.forEach((button) => {
      button.addEventListener('click', () => {
        const mode = button.dataset.pickerMode;
        picker.viewMode = mode;
        if (getColorMode() !== mode) setColorMode(mode);
        renderGrid();
        if (picker.scheduleOpen) renderSchedule();
        else revealCard(picker.grid.querySelector('.theme-card.is-current'), true);
      });
    });

    if (picker.scheduleButton && picker.schedule && hasSeasonal()) {
      picker.scheduleButton.hidden = false;
      picker.scheduleButton.addEventListener('click', () => setScheduleOpen(!picker.scheduleOpen));
      picker.schedule.addEventListener('click', (event) => {
        if (event.target.closest('[data-schedule-use]')) {
          selectSeasonal();
          renderSchedule();
          return;
        }
        // Touch has no hover: a tap shows the period's details without re-theming.
        const period = event.target.closest('[data-period-key]');
        if (period && !previewSaved) {
          showPreviewInfo({ key: period.dataset.periodKey, detail: period.dataset.previewDetail });
        }
      });
    }

    const shuffle = panel.querySelector('[data-picker-shuffle]');
    if (shuffle) {
      shuffle.addEventListener('click', () => {
        if (picker.scheduleOpen) setScheduleOpen(false);
        const pool = visibleCards('colors').map((card) => findEntry(card.dataset.name)).filter(Boolean);
        const pick = randomTheme(pool);
        if (pick) {
          const card = picker.grid.querySelector('.theme-card[data-name="' + CSS.escape(pick.name) + '"]');
          if (card) {
            picker.search.setAttribute('aria-activedescendant', card.id);
            revealCard(card, true);
          }
        }
      });
    }

    const reset = panel.querySelector('[data-picker-reset]');
    if (reset) {
      reset.addEventListener('click', () => {
        resetTheme();
        const grid = tabGrid(picker.activeTab);
        revealCard(grid && grid.querySelector('.theme-card.is-current'), true);
      });
    }

    const bake = panel.querySelector('[data-picker-bake]');
    if (bake && bakeEndpoint()) {
      bake.hidden = false;
      bake.addEventListener('click', () => onBakeClick(bake));
    }

    const settings = panel.querySelector('[data-picker-settings]');
    if (settings && typeof window.__markataSettingsEndpoint === 'string' && window.__markataSettingsEndpoint) {
      settings.hidden = false;
      settings.addEventListener('click', () => {
        const devSettings = window.markataDevSettings;
        if (!devSettings) return;
        closePicker(false);
        devSettings.open('theme');
      });
    }

    panel.addEventListener('keydown', (event) => {
      const key = event.key;
      if (key === 'Escape') {
        event.preventDefault();
        event.stopPropagation();
        closePicker(true);
        return;
      }
      if (event.altKey || event.ctrlKey || event.metaKey) return;
      const target = event.target;
      const inSearch = target === picker.search;
      const onTab = target.getAttribute && target.getAttribute('role') === 'tab';
      const onCard = target.classList && target.classList.contains('theme-card');

      if (target.closest && target.closest('.adaptive-loading')) return;
      if (!inSearch && /^[1-3]$/.test(key)) {
        const tab = availableTabs()[Number(key) - 1];
        if (tab) {
          event.preventDefault();
          showTab(tab, { focus: true, focusGrid: true });
        }
        return;
      }

      if (onTab) {
        const list = availableTabs();
        const index = list.indexOf(picker.activeTab);
        if (key === 'ArrowRight' || key === 'ArrowLeft') {
          event.preventDefault();
          const step = key === 'ArrowRight' ? 1 : -1;
          const next = list[(index + step + list.length) % list.length];
          showTab(next);
          picker.tabs[next].focus();
        } else if (key === 'ArrowDown') {
          event.preventDefault();
          focusActivePane(true);
        }
        return;
      }

      if (!inSearch && !onCard) return;
      if (key === 'Enter') {
        event.preventDefault();
        closePicker(true);
        return;
      }
      const caretKeys = key === 'ArrowLeft' || key === 'ArrowRight' || key === 'Home' || key === 'End';
      if (inSearch && caretKeys && picker.search.value) return;
      if (['ArrowRight', 'ArrowLeft', 'ArrowDown', 'ArrowUp', 'Home', 'End', 'PageDown', 'PageUp'].includes(key)) {
        event.preventDefault();
        moveSelection(key);
      }
    });

    window.addEventListener('text-size-change', syncUI);
    syncUI();
  }

  function bindGlobalListeners() {
    if (window.__markataThemePickerGlobals) return;
    window.__markataThemePickerGlobals = true;

    document.addEventListener('click', (event) => {
      const toggle = event.target.closest && event.target.closest('.palette-mode-toggle');
      if (toggle) {
        toggleColorMode();
        if (picker && !picker.panel.hidden) {
          picker.viewMode = getColorMode();
          renderGrid();
          if (picker.scheduleOpen) renderSchedule();
        }
        return;
      }
      if (picker && !picker.panel.hidden && !picker.root.contains(event.target)) {
        closePicker(false);
      }
    });

    document.addEventListener('keydown', (event) => {
      if (event.key === 'Escape' && picker && !picker.panel.hidden) {
        closePicker(true);
      }
    });
  }

  function showNotification(text) {
    let notice = document.querySelector('.palette-notification');
    if (!notice) {
      notice = document.createElement('div');
      notice.className = 'palette-notification';
      notice.setAttribute('role', 'status');
      notice.setAttribute('aria-live', 'polite');
      document.body.appendChild(notice);
    }
    notice.textContent = text;
    notice.classList.add('visible');
    clearTimeout(notice._timeout);
    notice._timeout = setTimeout(() => notice.classList.remove('visible'), 1100);
  }

  function setupKeyboardShortcuts() {
    function waitForRegistry(callback, attempts) {
      if (window.shortcutsRegistry) {
        callback();
      } else if ((attempts || 0) < 50) {
        setTimeout(() => waitForRegistry(callback, (attempts || 0) + 1), 10);
      }
    }

    waitForRegistry(() => {
      const registry = window.shortcutsRegistry;
      if (isPickerEnabled()) {
        registry.register({
          key: 't', modifiers: [], description: 'Open theme picker', group: 'theme', priority: 50,
          handler: (event) => {
            event.preventDefault();
            if (picker && picker.panel.hidden) openPicker();
            else closePicker(true);
          },
        });
        registry.register({
          key: ',', modifiers: [], description: 'Previous theme', group: 'theme', priority: 50,
          handler: (event) => { event.preventDefault(); cycleTheme('prev'); },
        });
        registry.register({
          key: '.', modifiers: [], description: 'Next theme', group: 'theme', priority: 50,
          handler: (event) => { event.preventDefault(); cycleTheme('next'); },
        });
        registry.register({
          key: '<', modifiers: [], description: 'Previous style', group: 'theme', priority: 50,
          handler: (event) => { event.preventDefault(); cycleAesthetic('prev'); },
        });
        registry.register({
          key: '>', modifiers: [], description: 'Next style', group: 'theme', priority: 50,
          handler: (event) => { event.preventDefault(); cycleAesthetic('next'); },
        });
        if (getFontManifest().length > 1) {
          registry.register({
            key: 'F', modifiers: [], description: 'Previous font', group: 'theme', priority: 50,
            handler: (event) => { event.preventDefault(); cycleFont('prev'); },
          });
          registry.register({
            key: 'f', modifiers: [], description: 'Next font', group: 'theme', priority: 50,
            handler: (event) => { event.preventDefault(); cycleFont('next'); },
          });
        }
      }
      if (document.querySelector('.palette-mode-toggle')) {
        registry.register({
          key: '\\', modifiers: [], description: 'Toggle dark/light mode', group: 'theme', priority: 50,
          handler: (event) => {
            event.preventDefault();
            const next = toggleColorMode();
            if (picker && !picker.panel.hidden) {
              picker.viewMode = next;
              renderGrid();
              if (picker.scheduleOpen) renderSchedule();
            }
            showNotification(next === 'dark' ? 'Dark mode' : 'Light mode');
          },
        });
      }
    });
  }

  function init() {
    if (!store.get(FONT_KEY) && root.dataset.fontpack) pageFontpack = root.dataset.fontpack;
    validateRestoredState();
    bindGlobalListeners();
    const container = document.querySelector('[data-theme-picker]');
    if (container && isPickerEnabled()) {
      const panel = container.querySelector('.theme-picker-panel');
      const fallback = document.querySelector('.header-action-controls > [data-adaptive-loading]');
      if (panel && fallback) panel.insertBefore(fallback, panel.querySelector('.theme-picker-bar'));
      container.hidden = false;
      mountPicker(container);
    } else {
      // Keep loading preferences accessible without a palette manifest.
      if (container) {
        const loading = container.querySelector('[data-adaptive-loading]');
        const actions = container.closest('.header-action-controls');
        if (loading && actions) actions.appendChild(loading);
        container.hidden = true;
      }
      picker = null;
      syncUI();
    }
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init, { once: true });
  } else {
    init();
  }
  setupKeyboardShortcuts();

  window.initPaletteSwitcher = init;
  window.addEventListener('view-transition-complete', init);

  window.markata = window.markata || {};
  window.markata.paletteSwitcher = {
    getManifest: getManifest,
    getColorMode: getColorMode,
    setColorMode: setColorMode,
    toggleColorMode: toggleColorMode,
    applyPalette: selectTheme,
    selectTheme: selectTheme,
    resetTheme: resetTheme,
    cycleFamily: cycleTheme,
    nextFamily: () => cycleTheme('next'),
    prevFamily: () => cycleTheme('prev'),
    randomTheme: () => randomTheme(),
    open: (tab) => openPicker(tab ? { tab: tab } : undefined),
    close: () => closePicker(false),
    getAesthetic: getAesthetic,
    setAesthetic: setAesthetic,
    cycleAesthetic: cycleAesthetic,
    nextAesthetic: () => cycleAesthetic('next'),
    prevAesthetic: () => cycleAesthetic('prev'),
    getFontpack: getFontpack,
    setFontpack: setFontpack,
    cycleFont: cycleFont,
    nextFont: () => cycleFont('next'),
    prevFont: () => cycleFont('prev'),
    getTextSize: getTextSize,
    setTextSize: setTextSize,
    showTab: (tab) => { if (picker && picker.panel.hidden) openPicker({ tab: tab }); else showTab(tab); },
    selectSeasonal: selectSeasonal,
    isSeasonal: () => isSeasonal(),
    getSeasonal: (mode) => seasonalPick(mode),
    getBakeSettings: getBakeSettings,
    bake: bakeConfig,
  };
})();
