(function(root, factory) {
  if (typeof module === 'object' && module.exports) module.exports = factory();
  else root.MarkataAdaptiveLoading = factory();
})(typeof window === 'undefined' ? this : window, function() {
  'use strict';

  var KEY = 'markata-loading-mode-v1';
  var SESSION_KEY = 'markata-loading-auto-v1';
  var DISMISS_KEY = 'markata-loading-dismissed-v1';
  var SESSION_TTL = 10 * 60 * 1000;
  var MAX_SAMPLES = 12;

  function state() {
    return { recommendation: 'unknown', poor: 0, good: 0, samples: 0, offline: false };
  }

  function vote(current, kind) {
    var next = Object.assign(state(), current || {});
    if (kind === 'offline') {
      next.offline = true;
      next.good = 0;
      next.recommendation = 'constrained';
      return next;
    }
    if (kind === 'online') {
      next.offline = false;
      return next;
    }
    if (kind === 'poor') {
      next.poor = Math.min(8, next.poor + 1);
      next.good = 0;
      next.samples = Math.min(MAX_SAMPLES, next.samples + 1);
    } else if (kind === 'good') {
      next.good = Math.min(8, next.good + 1);
      next.poor = 0;
      next.samples = Math.min(MAX_SAMPLES, next.samples + 1);
    } else if (kind === 'strong-poor') {
      next.poor = Math.max(2, next.poor + 2);
      next.good = 0;
    }
    if (next.offline || next.poor >= 2) next.recommendation = 'constrained';
    else if (next.good >= 4) next.recommendation = 'full-quality';
    else if (next.recommendation === 'unknown') next.recommendation = 'normal';
    return next;
  }

  function effectiveMode(userMode, automatic) {
    if (userMode === 'save-data') return 'constrained';
    if (userMode === 'full-quality') return 'full-quality';
    return automatic || 'unknown';
  }

  // Article videos are allowed to start when on screen, even before Auto has
  // collected enough evidence to classify the connection as Full Quality.
  function canAutoplayInViewport(effective, reducedMotion, saveData) {
    return effective !== 'constrained' && !reducedMotion && !saveData;
  }

  function readMode(storage, key) {
    try {
      var mode = storage && storage.getItem(key);
      return mode === 'save-data' || mode === 'full-quality' ? mode : 'auto';
    } catch (_) { return 'auto'; }
  }

  function persistMode(storage, key, mode) {
    try {
      if (!storage) return;
      if (mode === 'auto') storage.removeItem(key);
      else if (mode === 'save-data' || mode === 'full-quality') storage.setItem(key, mode);
    } catch (_) { /* storage may be blocked */ }
  }

  function mismatch(userMode, current) {
    if (current.offline) return userMode === 'full-quality' ? 'constrained' : '';
    if (userMode === 'full-quality' && current.poor >= 4) return 'constrained';
    if (userMode === 'save-data' && current.good >= 5) return 'full-quality';
    return '';
  }

  function networkHint(current, info) {
    if (!info) return current || state();
    if (info.saveData) return vote(current, 'strong-poor');
    if (current && (current.samples || current.offline || current.recommendation === 'constrained')) return current;
    var next = Object.assign(state(), current || {});
    var type = String(info.effectiveType || '').toLowerCase();
    if (type === 'slow-2g' || type === '2g') next.recommendation = 'constrained';
    else if (type === '4g' || Number(info.downlink) >= 1.5) next.recommendation = 'normal';
    return next;
  }

  function restore(serialized, now) {
    try {
      var saved = JSON.parse(serialized);
      if (!saved || now - saved.at > SESSION_TTL || now < saved.at) return state();
      var restored = state();
      restored.poor = Math.max(0, Math.min(8, saved.poor | 0));
      restored.good = Math.max(0, Math.min(8, saved.good | 0));
      restored.samples = Math.max(0, Math.min(MAX_SAMPLES, saved.samples | 0));
      restored.offline = !!saved.offline;
      if (restored.offline) restored.recommendation = 'constrained';
      else if (saved.recommendation === 'constrained' || restored.poor >= 2) restored.recommendation = 'constrained';
      else if (restored.good >= 4) restored.recommendation = 'full-quality';
      else if (restored.poor || restored.good) restored.recommendation = 'normal';
      return restored;
    } catch (_) { return state(); }
  }

  function savedState(current, now) {
    return JSON.stringify({ poor: current.poor, good: current.good, samples: current.samples, offline: current.offline, recommendation: current.recommendation, at: now });
  }

  var api = { state: state, vote: vote, networkHint: networkHint, effectiveMode: effectiveMode, readMode: readMode, persistMode: persistMode, mismatch: mismatch, restore: restore, savedState: savedState, canAutoplayInViewport: canAutoplayInViewport };
  if (typeof document === 'undefined') return api;

  var rootElement = document.documentElement;
  var control = document.querySelector('[data-adaptive-loading]');
  if (!control) return api;

  var radios = control.querySelectorAll('input[name="loading-mode"]');
  var status = control.querySelector('[data-loading-status]');
  var message = control.querySelector('[data-loading-mismatch]');
  var actions = control.querySelector('[data-loading-actions]');
  var useAuto = control.querySelector('[data-loading-use-auto]');
  var keepMode = control.querySelector('[data-loading-keep]');
  var details = control.querySelector('details');
  var sessionStore = null;
  var localStore = null;
  try { sessionStore = window.sessionStorage; } catch (_) { /* storage may be blocked */ }
  try { localStore = window.localStorage; } catch (_) { /* storage may be blocked */ }
  var current = restore(readStorage(sessionStore, SESSION_KEY), Date.now());
  var userMode = readMode(localStore, KEY);
  var samples = 0;
  var seen = new Set();
  var observer;
  var stopTimer;
  if (navigator.onLine && current.offline) current = vote(current, 'online');

  function readStorage(storage, key) {
    try { return storage ? storage.getItem(key) : null; } catch (_) { return null; }
  }

  function writeStorage(storage, key, value) {
    try { if (storage) storage.setItem(key, value); } catch (_) { /* storage may be blocked */ }
  }

  function clearStorage(storage, key) {
    try { if (storage) storage.removeItem(key); } catch (_) { /* storage may be blocked */ }
  }

  function setUserMode(mode) {
    userMode = mode;
    persistMode(localStore, KEY, mode);
    update();
    syncAutoplay();
  }

  function update() {
    if (!navigator.onLine) current = vote(current, 'offline');
    var connection = navigator.connection || navigator.mozConnection || navigator.webkitConnection;
    var automatic = current.offline || (connection && connection.saveData) ? 'constrained' : current.recommendation;
    var effective = effectiveMode(userMode, automatic);
    rootElement.dataset.loadingMode = userMode;
    rootElement.dataset.loadingPolicy = effective;
    rootElement.dataset.loadingEffective = effective;
    rootElement.dataset.loadingRecommendation = automatic;
    rootElement.dataset.adaptiveReady = 'true';
    // Autoplay of authored article media is handled by a viewport observer.
    // Preview-card videos have no authored intent and remain click/hover-only.
    for (var i = 0; i < radios.length; i++) radios[i].checked = radios[i].value === userMode;
    var labels = { auto: 'Auto', 'save-data': 'Save Data', 'full-quality': 'Full Quality' };
    var summary = control.querySelector('[data-loading-label]');
    if (summary) summary.textContent = 'Loading: ' + labels[userMode];
    var suggestion = userMode === 'save-data' && automatic === 'constrained' ? '' : mismatch(userMode, current);
    var dismissed = readStorage(sessionStore, DISMISS_KEY);
    if (actions) actions.hidden = !suggestion || dismissed === suggestion;
    if (keepMode) keepMode.textContent = 'Keep ' + labels[userMode];
    if (suggestion && dismissed !== suggestion) {
      if (status) status.hidden = false;
      if (message) {
        message.hidden = false;
        message.textContent = userMode === 'save-data'
          ? 'Save Data is on. Your connection appears fast now.'
          : 'Full Quality is on. Your connection appears constrained.';
      }
    } else {
      if (status) status.hidden = true;
      if (message) message.hidden = true;
    }
  }

  var autoplayObserver = null;
  var autoplayStates = new WeakMap();

  function playbackAllowed() {
    var connection = navigator.connection || navigator.mozConnection || navigator.webkitConnection;
    var reduced = window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    return navigator.onLine && !document.hidden &&
      canAutoplayInViewport(rootElement.dataset.loadingPolicy, reduced, !!(connection && connection.saveData));
  }

  // Pauses only playback started by this controller; manual playback is left alone.
  function pauseManaged(video, state) {
    if (!state || !state.started || video.paused) return;
    state.pausing = true;
    video.pause();
    window.setTimeout(function() { state.pausing = false; }, 0);
  }

  function playVisible(video, state) {
    if (!state || !state.visible || !video.isConnected || state.userPaused || !playbackAllowed()) return;
    // Never hijack a video the reader already started or scrubbed manually.
    if (!state.started && (video.currentTime || !video.paused)) return;
    if (!state.started) {
      state.started = true;
      video.dataset.autoplayStarted = 'true';
    }
    // Deliberately do not add a native autoplay attribute: preload remains none
    // until the video is visible and the browser accepts the play() request.
    var playing = video.play();
    if (playing && playing.catch) playing.catch(function() {
      // Browser policy can reject autoplay; native controls remain available.
    });
  }

  function watchAutoplay() {
    if (!window.IntersectionObserver) return; // No-JS/old browsers use controls.
    if (!autoplayObserver) {
      autoplayObserver = new IntersectionObserver(function(entries) {
        entries.forEach(function(entry) {
          var video = entry.target;
          var state = autoplayStates.get(video);
          if (!state) return;
          state.visible = entry.isIntersecting && entry.intersectionRatio >= 0.25;
          if (state.visible) playVisible(video, state);
          else pauseManaged(video, state);
        });
      }, { threshold: [0, 0.25] });
    }
    document.querySelectorAll('video[data-authored-autoplay="true"]:not([data-adaptive-preview])').forEach(function(video) {
      if (autoplayStates.has(video)) return;
      var state = { visible: false, started: false, userPaused: false, pausing: false };
      autoplayStates.set(video, state);
      video.addEventListener('pause', function() {
        if (state.started && state.visible && !state.pausing) state.userPaused = true;
      });
      autoplayObserver.observe(video);
    });
  }

  function syncAutoplay() {
    document.querySelectorAll('video[data-authored-autoplay="true"]').forEach(function(video) {
      var state = autoplayStates.get(video);
      if (!state) return;
      if (!playbackAllowed()) pauseManaged(video, state);
      else if (state.visible) playVisible(video, state);
    });
  }

  function resetAutoplay() {
    if (autoplayObserver) autoplayObserver.disconnect();
    autoplayStates = new WeakMap();
    watchAutoplay();
  }

  document.addEventListener('visibilitychange', syncAutoplay);
  if (window.matchMedia) {
    var reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
    if (reducedMotion.addEventListener) reducedMotion.addEventListener('change', syncAutoplay);
  }

  function record(kind) {
    if (samples >= MAX_SAMPLES) return;
    current = vote(current, kind);
    if (kind === 'good' || kind === 'poor') samples++;
    if (samples >= MAX_SAMPLES && observer) observer.disconnect();
    writeStorage(sessionStore, SESSION_KEY, savedState(current, Date.now()));
    update();
  }

  function connectionHint() {
    var connection = navigator.connection || navigator.mozConnection || navigator.webkitConnection;
    if (!connection) return;
    current = networkHint(current, connection);
    writeStorage(sessionStore, SESSION_KEY, savedState(current, Date.now()));
  }

  function inspect(entry) {
    if (samples >= MAX_SAMPLES || !entry || (entry.entryType !== 'navigation' && entry.initiatorType !== 'img' && entry.initiatorType !== 'video' && entry.initiatorType !== 'audio')) return;
    var identity = entry.entryType + ':' + entry.name + ':' + entry.startTime;
    if (seen.has(identity)) return;
    seen.add(identity);
    var bytes = Number(entry.transferSize || 0);
    var elapsed = Number(entry.responseEnd - entry.responseStart);
    if (bytes < 24000 || elapsed <= 0 || elapsed > 30000) return;
    var rate = bytes * 8 / (elapsed / 1000);
    if (rate < 2000000) record('poor');
    else if (rate > 4000000 && entry.responseStart - entry.startTime < 500) record('good');
  }

  for (var r = 0; r < radios.length; r++) radios[r].addEventListener('change', function() {
    if (this.checked) setUserMode(this.value);
  });
  if (useAuto) useAuto.addEventListener('click', function() { clearStorage(sessionStore, DISMISS_KEY); setUserMode('auto'); });
  if (keepMode) keepMode.addEventListener('click', function() {
    var recommendation = mismatch(userMode, current);
    if (recommendation) writeStorage(sessionStore, DISMISS_KEY, recommendation);
    if (details) details.open = false;
    update();
  });
  if (status && details) status.addEventListener('click', function() { details.open = true; });
  window.addEventListener('offline', function() {
    current = vote(current, 'offline');
    writeStorage(sessionStore, SESSION_KEY, savedState(current, Date.now()));
    update();
  });
  window.addEventListener('online', function() {
    current = vote(current, 'online');
    writeStorage(sessionStore, SESSION_KEY, savedState(current, Date.now()));
    update();
  });

  connectionHint();
  update();
  function collectTimings() {
    if (!performance.getEntriesByType) return;
    var navigation = performance.getEntriesByType('navigation');
    if (navigation.length) inspect(navigation[0]);
    var entries = performance.getEntriesByType('resource');
    for (var e = 0; e < entries.length && samples < MAX_SAMPLES; e++) inspect(entries[e]);
  }
  function observeTimings(buffered) {
    if (!window.PerformanceObserver) return;
    try {
      if (!observer) observer = new PerformanceObserver(function(list) { list.getEntries().forEach(inspect); });
      observer.observe({ type: 'resource', buffered: buffered });
      window.clearTimeout(stopTimer);
      stopTimer = window.setTimeout(function() { observer.disconnect(); }, 30000);
    } catch (_) { /* load snapshot remains available */ }
  }
  window.addEventListener('view-transition-complete', function() {
    var saved = readStorage(sessionStore, SESSION_KEY);
    if (saved) current = restore(saved, Date.now());
    samples = 0;
    seen = new Set();
    update();
    resetAutoplay();
    observeTimings(false);
  });
  observeTimings(true);
  // Wait for page resources to finish before starting initial autoplay, so
  // measured poor-network evidence can suppress unnecessary video downloads.
  function loaded() {
    collectTimings();
    watchAutoplay();
  }
  if (document.readyState === 'complete') loaded();
  else window.addEventListener('load', loaded, { once: true });
  return api;
});
