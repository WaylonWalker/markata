/* Shorts: static, finite feed with a virtualized 5+1+5 viewing window. */
(function () {
  'use strict';

  var viewer = document.getElementById('shorts-viewer');
  if (!viewer) return;

  var slides = document.getElementById('shorts-slides');
  var title = document.getElementById('shorts-title');
  var description = document.getElementById('shorts-description');
  var more = document.getElementById('shorts-more');
  var permalink = document.getElementById('shorts-permalink');
  var counter = document.getElementById('shorts-counter');
  var status = document.getElementById('shorts-message');
  var endNotice = document.getElementById('shorts-end');
  var up = document.getElementById('shorts-up');
  var down = document.getElementById('shorts-down');
  var playButton = document.getElementById('shorts-play');
  var soundButton = document.getElementById('shorts-sound');

  var indexURL = new URL(viewer.dataset.manifest, location.href);
  var manifest = null;
  var position = 0;
  var soundOn = false;
  var paused = false;
  var slots = new Map();
  var chunks = new Map();
  var inFlight = new Map();
  var lastWheel = 0;
  var touchStartY = null;
  var renderEpoch = 0;
  var activeVideo = null;

  var connection = navigator.connection || navigator.mozConnection || navigator.webkitConnection;
  var reducedMotion = window.matchMedia && matchMedia('(prefers-reduced-motion: reduce)').matches;

  function chunkName(n) { return String(n).padStart(4, '0') + '.json'; }
  function clamp(n) { return Math.max(0, Math.min((manifest ? manifest.ids.length : 1) - 1, n)); }
  function encodedHash(index) { return '#id=' + encodeURIComponent(manifest.ids[index]); }
  function idFromHash() {
    if (!location.hash.startsWith('#id=')) return '';
    try { return decodeURIComponent(location.hash.slice(4)); }
    catch (_) { return ''; }
  }

  function getChunk(n) {
    if (chunks.has(n)) return Promise.resolve(chunks.get(n));
    if (inFlight.has(n)) return inFlight.get(n);
    var req = fetch(new URL(chunkName(n), indexURL), { credentials: 'same-origin' })
      .then(function (response) {
        if (!response.ok) throw new Error('Failed to load shorts chunk ' + n);
        return response.json();
      }).then(function (items) {
        if (!Array.isArray(items)) throw new Error('Invalid shorts chunk ' + n);
        chunks.set(n, items);
        inFlight.delete(n);
        // Keep nearby chunks. Never accumulate the entire feed in memory.
        for (var key of chunks.keys()) {
          if (Math.abs(key - Math.floor(position / manifest.chunk_size)) > 3) {
            chunks.delete(key);
          }
        }
        return items;
      }).catch(function (err) { inFlight.delete(n); throw err; });
    inFlight.set(n, req);
    return req;
  }

  function getItem(n) {
    var chunkNumber = Math.floor(n / manifest.chunk_size);
    return getChunk(chunkNumber).then(function (items) {
      var item = items[n % manifest.chunk_size];
      if (!item || item.id !== manifest.ids[n]) throw new Error('Feed index out of sync');
      return item;
    });
  }

  function stopVideo() {
    if (activeVideo) {
      activeVideo.pause();
      activeVideo.removeAttribute('src');
      activeVideo.load();
      activeVideo = null;
    }
  }

  function previewBudget() {
    if (connection && (connection.saveData || connection.effectiveType === 'slow-2g' || connection.effectiveType === '2g')) return 0;
    return 1;
  }

  function canPrefetchSecondPreview() {
    return Boolean(connection && !connection.saveData && connection.effectiveType === '4g' &&
      (!connection.downlink || connection.downlink >= 5));
  }

  function shouldLoadPreview(n) {
    var ahead = n - position;
    return ahead > 0 && ahead <= previewBudget();
  }

  function loadImage(img) {
    var src = img.dataset.shortsSrc;
    if (!src) return;
    delete img.dataset.shortsSrc;
    img.src = src;
  }

  function loadImages(wrap) {
    wrap.querySelectorAll('img[data-shorts-src]:not([data-shorts-hold])').forEach(loadImage);
  }

  function prefetchSecondPreview() {
    if (!canPrefetchSecondPreview() || !manifest) return;
    var n = position + 2;
    var node = slots.get(n);
    if (!node || n >= manifest.ids.length || node.dataset.mediaMode === 'preview') return;
    var epoch = renderEpoch;
    getItem(n).then(function (item) {
      if (epoch !== renderEpoch || slots.get(n) !== node || !canPrefetchSecondPreview() || !item.thumb) return;
      releaseSlideMedia(node);
      var media = mediaFor(item, false, true);
      node.replaceChildren(media);
      node.dataset.role = 'preview';
      node.dataset.id = item.id;
      node.dataset.mediaMode = 'preview';
      node.classList.toggle('shorts-slide--video', item.kind === 'video');
      node.classList.toggle('shorts-slide--empty', !item.src);
      loadImages(media);
    }).catch(function () {});
  }

  function releaseSlideMedia(node) {
    if (activeVideo && node.contains(activeVideo)) stopVideo();
    node.querySelectorAll('img').forEach(function (img) {
      img.removeAttribute('src');
      img.removeAttribute('srcset');
      delete img.dataset.shortsSrc;
    });
  }

  function mediaFor(item, active, loadPreview, reusedPreview) {
    var wrap = document.createElement('div');
    wrap.className = 'shorts-media-wrap';
    wrap.style.cssText = 'position:absolute;inset:0;overflow:hidden;display:grid;place-items:center';
    if (!active && !loadPreview) return wrap;

    var placeholder = null;
    if (active && item.kind === 'image' && (reusedPreview || item.placeholder)) {
      // Keep the already-loaded 240px upcoming preview when it becomes active.
      // A 72px request would add bytes and visibly downgrade the image.
      placeholder = reusedPreview || document.createElement('img');
      placeholder.className = 'shorts-image-placeholder';
      placeholder.alt = '';
      placeholder.setAttribute('aria-hidden', 'true');
      placeholder.loading = 'eager';
      placeholder.decoding = 'async';
      placeholder.fetchPriority = 'high';
      if (!reusedPreview) placeholder.dataset.shortsSrc = item.placeholder;
      wrap.appendChild(placeholder);
    }
    if (item.thumb && active && item.kind === 'video') {
      var blur = document.createElement('img');
      blur.className = 'shorts-blur';
      blur.alt = '';
      blur.setAttribute('aria-hidden', 'true');
      blur.loading = 'eager';
      blur.decoding = 'async';
      blur.fetchPriority = 'low';
      blur.dataset.shortsSrc = item.thumb;
      wrap.appendChild(blur);
    }
    if (!item.src) return wrap;
    var el;
    if (item.kind === 'video' && active) {
      el = document.createElement('video');
      el.playsInline = true;
      el.muted = !soundOn;
      el.loop = true;
      el.autoplay = true;
      el.preload = 'metadata';
      if (item.poster) el.poster = item.poster;
      el.src = item.src;
      if (item.mime) el.setAttribute('type', item.mime);
      el.addEventListener('click', function () { togglePlay(); });
      activeVideo = el;
      if (!paused) {
        el.play().catch(function () {
          if (activeVideo === el) {
            paused = true;
            updatePlayback();
          }
        });
      }
    } else {
      var imageSource = active ? item.src : item.thumb;
      if (!imageSource) return wrap;
      el = document.createElement('img');
      el.dataset.shortsSrc = imageSource;
      el.loading = 'eager';
      el.decoding = 'async';
      el.alt = item.alt || item.title || 'Photograph';
      if (active) {
        el.className = 'shorts-media shorts-image--full';
        el.fetchPriority = 'high';
        if (placeholder && !reusedPreview) {
          // Give the tiny preview a brief head start without delaying the
          // sharp photo indefinitely if the placeholder fails or stalls.
          el.dataset.shortsHold = 'true';
          var started = false;
          var startSharp = function () {
            if (started) return;
            started = true;
            delete el.dataset.shortsHold;
            if (el.isConnected) loadImage(el);
          };
          placeholder.addEventListener('load', startSharp, { once: true });
          placeholder.addEventListener('error', function () {
            placeholder.remove();
            startSharp();
          }, { once: true });
          setTimeout(startSharp, 180);
        }
        el.addEventListener('load', function () {
          var decoded = typeof el.decode === 'function' ? el.decode() : Promise.resolve();
          decoded.catch(function () {}).then(function () {
            if (el.isConnected) {
              wrap.classList.add('shorts-media-wrap--sharp');
              prefetchSecondPreview();
            }
          });
        }, { once: true });
        el.addEventListener('error', function () {
          if (!el.isConnected) return;
          wrap.querySelector('.shorts-image-placeholder')?.remove();
          el.hidden = true;
          nodeFor(wrap)?.classList.add('shorts-slide--empty');
        }, { once: true });
      } else {
        el.fetchPriority = 'low';
      }
    }
    if (!(active && item.kind === 'image')) el.className = 'shorts-media';
    wrap.appendChild(el);
    return wrap;
  }

  function nodeFor(wrap) { return wrap.parentElement; }

  function updatePlayback() {
    playButton.textContent = paused ? '▶' : 'Ⅱ';
    playButton.setAttribute('aria-label', paused ? 'Play video' : 'Pause video');
    soundButton.textContent = soundOn ? '♫' : '♪';
    soundButton.setAttribute('aria-label', soundOn ? 'Mute video' : 'Unmute video');
  }

  function setHUD(item, currentIndex) {
    title.textContent = item.title || item.id;
    description.textContent = item.description || '';
    description.classList.remove('expanded');
    more.textContent = 'More';
    more.setAttribute('aria-expanded', 'false');
    more.hidden = !item.description || item.description.length < 135;
    permalink.href = item.href || '/';
    counter.textContent = (currentIndex + 1).toLocaleString() + ' / ' + manifest.total.toLocaleString();
    playButton.hidden = item.kind !== 'video' || !item.src;
    soundButton.hidden = item.kind !== 'video' || !item.src;
    updatePlayback();
    status.hidden = true;
  }

  function buildSlide(n) {
    var node = document.createElement('article');
    node.className = 'shorts-slide';
    node.dataset.index = String(n);
    node.setAttribute('aria-hidden', n === position ? 'false' : 'true');
    node.dataset.role = '';
    slides.appendChild(node);
    slots.set(n, node);
    return node;
  }

  function fillSlide(n, node, active, epoch, loadPreview) {
    getItem(n).then(function (item) {
      if (epoch !== renderEpoch || slots.get(n) !== node) return;
      var role = active ? 'active' : 'preview';
      var mediaMode = active ? 'active' : (loadPreview ? 'preview' : 'mounted');
      if (node.dataset.role !== role || node.dataset.id !== item.id || node.dataset.mediaMode !== mediaMode) {
        if (active) stopVideo();
        var reusedPreview = null;
        if (active && item.kind === 'image' &&
            node.dataset.mediaMode === 'preview' && node.dataset.id === item.id) {
          reusedPreview = node.querySelector('img.shorts-media');
          if (reusedPreview) reusedPreview.remove();
        }
        releaseSlideMedia(node);
        var media = mediaFor(item, active, loadPreview, reusedPreview);
        node.replaceChildren(media);
        loadImages(media);
        node.dataset.role = role;
        node.dataset.id = item.id;
        node.dataset.mediaMode = mediaMode;
        node.classList.toggle('shorts-slide--video', item.kind === 'video');
        node.classList.toggle('shorts-slide--empty', !item.src);
      }
      if (active) setHUD(item, n);
    }).catch(function (err) {
      if (epoch === renderEpoch && n === position) {
        status.hidden = false;
        status.textContent = 'Could not load this short. Swipe to another or open the gallery.';
        console.warn('[shorts]', err);
      }
    });
  }

  function render() {
    if (!manifest || !manifest.ids.length) return;
    renderEpoch += 1;
    var epoch = renderEpoch;
    for (var key of Array.from(slots.keys())) {
      if (Math.abs(key - position) > 5) {
        var stale = slots.get(key);
        releaseSlideMedia(stale);
        stale.remove();
        slots.delete(key);
      }
    }
    // Keep DOM bounded to eleven media slides regardless of feed length.
    for (var n = Math.max(0, position - 5); n <= Math.min(manifest.ids.length - 1, position + 5); n++) {
      var node = slots.get(n) || buildSlide(n);
      node.style.transform = 'translate3d(0,' + ((n - position) * 100) + '%,0)';
      node.setAttribute('aria-hidden', n === position ? 'false' : 'true');
      fillSlide(n, node, n === position, epoch, shouldLoadPreview(n));
    }
    up.disabled = position === 0;
    down.disabled = position === manifest.ids.length - 1;
    endNotice.hidden = position !== manifest.ids.length - 1;
  }

  function navigate(n, writeHash) {
    if (!manifest || !manifest.ids.length) return;
    n = clamp(n);
    if (n === position && status.hidden) return;
    if (activeVideo) stopVideo();
    position = n;
    paused = false;
    status.hidden = true;
    if (writeHash) history.replaceState(history.state, '', location.pathname + location.search + encodedHash(position));
    render();
  }

  function move(offset) {
    if (manifest) navigate(position + offset, true);
  }

  function togglePlay() {
    if (!activeVideo) return;
    paused = !activeVideo.paused;
    if (paused) activeVideo.pause();
    else activeVideo.play().catch(function () { paused = true; updatePlayback(); });
    updatePlayback();
  }
  function toggleSound() {
    soundOn = !soundOn;
    if (activeVideo) {
      activeVideo.muted = !soundOn;
      if (activeVideo.paused && !paused) activeVideo.play().catch(function () {});
    }
    updatePlayback();
  }

  up.addEventListener('click', function () { move(-1); });
  down.addEventListener('click', function () { move(1); });
  playButton.addEventListener('click', togglePlay);
  soundButton.addEventListener('click', toggleSound);
  document.getElementById('shorts-restart').addEventListener('click', function () { navigate(0, true); });
  more.addEventListener('click', function () {
    var expanded = description.classList.toggle('expanded');
    more.textContent = expanded ? 'Less' : 'More';
    more.setAttribute('aria-expanded', String(expanded));
  });

  viewer.addEventListener('wheel', function (event) {
    if (event.target.closest('.shorts-hud, .shorts-top, .shorts-end')) return;
    if (Math.abs(event.deltaY) < 8) return;
    event.preventDefault();
    var now = performance.now();
    if (now - lastWheel < (reducedMotion ? 140 : 320)) return;
    lastWheel = now;
    move(event.deltaY > 0 ? 1 : -1);
  }, { passive: false });

  viewer.addEventListener('touchstart', function (event) {
    if (event.touches.length !== 1 || event.target.closest('.shorts-hud, .shorts-top, .shorts-end')) {
      touchStartY = null;
      return;
    }
    touchStartY = event.touches[0].clientY;
  }, { passive: true });
  viewer.addEventListener('touchend', function (event) {
    if (touchStartY === null || !event.changedTouches.length) return;
    var delta = event.changedTouches[0].clientY - touchStartY;
    touchStartY = null;
    if (Math.abs(delta) > 45) move(delta < 0 ? 1 : -1);
  }, { passive: true });

  document.addEventListener('keydown', function (event) {
    if (event.altKey || event.ctrlKey || event.metaKey ||
        /^(INPUT|TEXTAREA|SELECT)$/.test(document.activeElement.tagName)) return;
    if (event.key === 'ArrowDown' || event.key === 'PageDown') { event.preventDefault(); move(1); }
    if (event.key === 'ArrowUp' || event.key === 'PageUp') { event.preventDefault(); move(-1); }
    if (event.key === ' ' && !event.target.closest('button, a')) { event.preventDefault(); togglePlay(); }
  });
  document.addEventListener('visibilitychange', function () {
    if (!activeVideo) return;
    if (document.hidden) activeVideo.pause();
    else if (!paused) activeVideo.play().catch(function () {});
  });
  function restoreHash() {
    if (!manifest) return;
    var target = manifest.ids.indexOf(idFromHash());
    if (target !== -1 && target !== position) navigate(target, false);
  }
  window.addEventListener('popstate', restoreHash);
  window.addEventListener('hashchange', restoreHash);

  fetch(indexURL, { credentials: 'same-origin' }).then(function (r) {
    if (!r.ok) throw new Error('Shorts index unavailable');
    return r.json();
  }).then(function (data) {
    if (!Array.isArray(data.ids) || !Number.isInteger(data.chunk_size) || data.chunk_size < 1) throw new Error('Invalid Shorts index');
    manifest = data;
    if (!manifest.ids.length) {
      status.textContent = 'No posts in this feed yet.';
      return;
    }
    var initial = manifest.ids.indexOf(idFromHash());
    navigate(initial >= 0 ? initial : 0, initial < 0);
  }).catch(function (err) {
    status.hidden = false;
    status.textContent = 'Shorts is unavailable. Return to Shots to browse posts.';
    console.warn('[shorts]', err);
  });

  // Diagnostic hook for performance and coverage tests; contains no private post data.
  window.__markataShorts = {
    get total() { return manifest ? manifest.ids.length : 0; },
    get position() { return position; },
    get mounted() { return slots.size; },
    goTo: function (i) { navigate(i, true); }
  };
})();
