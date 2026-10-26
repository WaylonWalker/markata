(function () {
  "use strict";

  function ready(callback) {
    if (document.readyState === "loading") {
      document.addEventListener("DOMContentLoaded", callback);
    } else {
      callback();
    }
  }

  ready(function () {
    var root = document.querySelector("[data-image-library]");
    if (!root) return;

    var grid = root.querySelector("[data-image-grid]");
    var cards = Array.prototype.slice.call(root.querySelectorAll("[data-image-card]"));
    var baseHref = (root.getAttribute ? root.getAttribute("data-image-base") : null) || "/images/";
    var currentPage = Number((root.getAttribute ? root.getAttribute("data-image-page") : null) || 1);
    var totalPages = Number((root.getAttribute ? root.getAttribute("data-image-total-pages") : null) || 1);
    var totalCount = Number((root.getAttribute ? root.getAttribute("data-image-total-count") : null) || cards.length);
    var allPagesLoaded = totalPages <= 1;
    var loadPromise = null;
    var search = root.querySelector("[data-image-search]");
    var sort = root.querySelector("[data-image-sort]");
    var count = root.querySelector("[data-image-count]");
    var empty = root.querySelector("[data-image-empty]");
    var status = root.querySelector("[data-image-status]");
    var filterButtons = Array.prototype.slice.call(root.querySelectorAll("[data-image-filter]"));
    var activeFilter = "all";
    var frame = 0;
    var recentWindowSeconds = 30 * 24 * 60 * 60;

    function text(value) {
      return (value || "").toLowerCase().trim();
    }

    function cardMatches(card, query) {
      if (!query) return true;
      return text(card.getAttribute("data-search")).indexOf(query) !== -1;
    }

    function recentlyAdded(card) {
      var addedAt = Number(card.getAttribute("data-added") || 0);
      var now = Math.floor(Date.now() / 1000);
      if (!(addedAt > 0)) return false;
      return addedAt >= now - recentWindowSeconds && addedAt <= now;
    }

    function compareNames(left, right) {
      return text(left.getAttribute("data-name")).localeCompare(text(right.getAttribute("data-name"))) || text(left.getAttribute("data-src")).localeCompare(text(right.getAttribute("data-src")));
    }

    function filterMatches(card) {
      if (activeFilter === "used") return card.getAttribute("data-used") === "true";
      if (activeFilter === "unused") return card.getAttribute("data-used") !== "true";
      if (activeFilter === "cover") return card.getAttribute("data-cover") === "true";
      if (activeFilter === "recent") return recentlyAdded(card);
      return true;
    }

    function compareCards(left, right) {
      var mode = sort ? sort.value : "name";
      var leftValue;
      var rightValue;
      if (mode === "recent") {
        leftValue = Number(left.getAttribute("data-added") || 0);
        rightValue = Number(right.getAttribute("data-added") || 0);
        return rightValue - leftValue || compareNames(left, right);
      }
      if (mode === "used") {
        leftValue = Number(left.getAttribute("data-uses") || 0);
        rightValue = Number(right.getAttribute("data-uses") || 0);
        return rightValue - leftValue || compareNames(left, right);
      }
      if (mode === "latest-used") {
        leftValue = Number(left.getAttribute("data-last-used") || 0);
        rightValue = Number(right.getAttribute("data-last-used") || 0);
        return rightValue - leftValue || compareNames(left, right);
      }
      return compareNames(left, right);
    }

    function loadVideo(video) {
      if (video.getAttribute("data-video-loaded") === "true") return;
      var sourceURL = video.getAttribute("data-video-src");
      if (!sourceURL) return;

      var source = document.createElement("source");
      source.src = sourceURL;
      var type = video.getAttribute("data-video-type");
      if (type) source.type = type;
      video.appendChild(source);
      video.setAttribute("data-video-loaded", "true");
      video.load();
    }

    function setupLazyVideos() {
      var videos = Array.prototype.slice.call(root.querySelectorAll("[data-video-src]"));
      if (!videos.length) return;

      if (typeof window.IntersectionObserver === "function") {
        var observer = new IntersectionObserver(function (entries) {
          entries.forEach(function (entry) {
            if (!entry.isIntersecting && !(entry.intersectionRatio > 0)) return;
            loadVideo(entry.target);
            observer.unobserve(entry.target);
          });
        });
        videos.forEach(function (video) {
          observer.observe(video);
        });
        return;
      }

      videos.forEach(function (video) {
        var card = video.closest("[data-image-card]") || video;
        ["pointerenter", "focusin", "touchstart"].forEach(function (eventName) {
          card.addEventListener(eventName, function () {
            loadVideo(video);
          }, { once: true, passive: eventName === "touchstart" });
        });
      });
    }

    function apply() {
      var query = text(search && search.value);
      var visible = [];
      cards.forEach(function (card) {
        var isVisible = filterMatches(card) && cardMatches(card, query);
        card.hidden = !isVisible;
        if (isVisible) visible.push(card);
      });

      visible.sort(compareCards);
      // Moving each card independently makes the large library expensive to
      // sort, and can restart image loads. Skip unchanged order and batch
      // actual reordering into a single DOM insertion.
      // Hidden cards still occupy grid.children positions. Compare only the
      // visible DOM order or an unchanged filter would repeatedly detach and
      // reinsert every matching image on each search/sort interaction.
      var visibleInDOM = Array.prototype.filter.call(grid.children, function (card) {
        return !card.hidden;
      });
      var needsReorder = visible.some(function (card, index) {
        return visibleInDOM[index] !== card;
      });
      if (needsReorder) {
        var fragment = document.createDocumentFragment();
        visible.forEach(function (card) {
          fragment.appendChild(card);
        });
        grid.appendChild(fragment);
      }

      if (count) {
        count.textContent = visible.length + " of " + (allPagesLoaded ? cards.length : totalCount) + " media items" +
          (allPagesLoaded ? "" : " (page " + currentPage + " of " + totalPages + ")");
      }
      if (empty) {
        empty.hidden = visible.length !== 0;
      }
    }

    // Keep the initial document bounded. Only if the visitor searches, sorts,
    // or filters do we hydrate the remaining, already server-rendered pages.
    // This retains global authoring controls without a JS dependency for basic
    // paging or navigation.
    function hydrateAllPages() {
      if (allPagesLoaded) return Promise.resolve(true);
      if (loadPromise) return loadPromise;
      if (typeof window.fetch !== "function" || typeof window.DOMParser !== "function") {
        announce("Global filtering needs JavaScript fetch; browse the page links instead.", true);
        return Promise.resolve(false);
      }
      if (status) status.textContent = "Loading the rest of the image library...";
      var requests = [];
      for (var number = 1; number <= totalPages; number++) {
        if (number === currentPage) continue;
        (function (pageNumber) {
          var href = pageNumber === 1 ? baseHref : baseHref + "page/" + pageNumber + "/";
          requests.push(fetch(href, { credentials: "same-origin" }).then(function (response) {
            if (!response.ok) throw new Error("Image library page " + pageNumber + " returned " + response.status);
            return response.text();
          }).then(function (markup) {
            var page = new DOMParser().parseFromString(markup, "text/html");
            var pageCards = page.querySelectorAll("[data-image-card]");
            if (!pageCards.length) throw new Error("Image library page " + pageNumber + " has no cards");
            return Array.prototype.slice.call(pageCards);
          }));
        }(number));
      }
      loadPromise = Promise.all(requests).then(function (pageCards) {
        var fragment = document.createDocumentFragment();
        pageCards.forEach(function (entry) {
          entry.forEach(function (card) { fragment.appendChild(card); });
        });
        grid.appendChild(fragment);
        cards = Array.prototype.slice.call(root.querySelectorAll("[data-image-card]"));
        allPagesLoaded = true;
        var pagination = root.querySelector("[data-image-pagination]");
        if (pagination) pagination.hidden = true;
        setupLazyVideos();
        if (status) status.textContent = "";
        return true;
      }).catch(function () {
        announce("Could not load the whole library. Use the numbered page links to continue.", true);
        return false;
      }).then(function (success) {
        loadPromise = null;
        return success;
      });
      return loadPromise;
    }

    function scheduleApply() {
      if (!allPagesLoaded) {
        hydrateAllPages().then(function (success) {
          if (success) scheduleApply();
        });
        return;
      }
      if (frame) window.cancelAnimationFrame(frame);
      frame = window.requestAnimationFrame(function () {
        frame = 0;
        apply();
      });
    }

    filterButtons.forEach(function (button) {
      button.addEventListener("click", function () {
        activeFilter = button.getAttribute("data-image-filter") || "all";
        if (activeFilter === "recent" && sort) sort.value = "recent";
        filterButtons.forEach(function (item) {
          item.setAttribute("aria-pressed", item === button ? "true" : "false");
        });
        scheduleApply();
      });
    });

    if (search) search.addEventListener("input", scheduleApply);
    if (sort) sort.addEventListener("change", scheduleApply);

    function announce(message, failed) {
      if (!status) return;
      status.textContent = message;
      status.dataset.state = failed ? "error" : "success";
      window.setTimeout(function () {
        status.textContent = "";
        status.removeAttribute("data-state");
      }, 2600);
    }

    function fallbackCopy(value) {
      var area = document.createElement("textarea");
      area.value = value;
      area.setAttribute("readonly", "");
      area.style.position = "fixed";
      area.style.opacity = "0";
      document.body.appendChild(area);
      area.select();
      var copied = false;
      try {
        copied = document.execCommand("copy");
      } catch (_) {
        copied = false;
      }
      document.body.removeChild(area);
      return copied;
    }

    function copyValue(value) {
      if (navigator.clipboard && navigator.clipboard.writeText) {
        return navigator.clipboard.writeText(value).then(function () {
          return true;
        }).catch(function () {
          return fallbackCopy(value);
        });
      }
      return Promise.resolve(fallbackCopy(value));
    }

    root.addEventListener("click", function (event) {
      var button = event.target.closest("[data-copy-value]");
      if (!button) return;
      var value = button.getAttribute("data-copy-value") || "";
      var original = button.textContent;
      button.disabled = true;
      copyValue(value).then(function (copied) {
        if (copied) {
          button.textContent = "Copied";
          announce("Copied to clipboard.", false);
        } else {
          announce("Copy failed. Select the value manually.", true);
        }
      }).catch(function () {
        announce("Copy failed. Select the value manually.", true);
      }).finally(function () {
        window.setTimeout(function () {
          button.disabled = false;
          button.textContent = original;
        }, 1800);
      });
    });

    setupLazyVideos();
    apply();
  });
}());
