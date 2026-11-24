(function() {
  'use strict';

  var d3Loader = null;
  var initFrame = null;

  function loadD3() {
    if (window.d3) {
      return Promise.resolve(window.d3);
    }

    if (d3Loader) {
      return d3Loader;
    }

    d3Loader = new Promise(function(resolve, reject) {
      var script = document.createElement('script');
      script.src = window.MARKATA_GO_D3_URL || 'https://d3js.org/d3.v7.min.js';
      script.onload = function() { resolve(window.d3); };
      script.onerror = reject;
      document.body.appendChild(script);
    });

    return d3Loader;
  }

  function cleanupPostGraphs() {
    var cleanups = window.__postGraphCleanup || [];
    while (cleanups.length) {
      var cleanup = cleanups.pop();
      try {
        cleanup();
      } catch (_) {
      }
    }
  }

  function registerCleanup(cleanup) {
    if (!window.__postGraphCleanup) {
      window.__postGraphCleanup = [];
    }
    window.__postGraphCleanup.push(cleanup);
  }

  function initSection(section) {
    var graphUrl = section.dataset.graph;
    var postHref = section.dataset.post;
    var canvas = section.querySelector('.post-graph__canvas-el');
    var tooltip = section.querySelector('.post-graph__tooltip');
    var limitInput = section.querySelector('.post-graph__limit');
    var limitValue = section.querySelector('.post-graph__limit-value');
    var animateToggle = section.querySelector('.post-graph__animate');
    var labelsToggle = section.querySelector('.post-graph__labels');
    var ctx = canvas ? canvas.getContext('2d') : null;
    var graphData = null;
    var nodes = [];
    var edges = [];
    var simulation = null;
    var width = 0;
    var height = 0;
    var hoveredNode = null;
    var dragNode = null;
    var panX = 0;
    var panY = 0;
    var zoom = 1;
    var rebuildTimer = null;
    var destroyed = false;
    var prefersReducedMotion = window.matchMedia ? window.matchMedia('(prefers-reduced-motion: reduce)') : null;
    var styles = getComputedStyle(section);
    var accentColor = styles.getPropertyValue('--garden-accent').trim() || '#ffcd11';
    var mutedColor = styles.getPropertyValue('--garden-muted').trim() || '#999';
    var motionListener = null;

    if (!canvas || !ctx) {
      return;
    }

    function resizeCanvas() {
      if (destroyed || !canvas.parentElement) {
        return;
      }

      var rect = canvas.parentElement.getBoundingClientRect();
      width = rect.width;
      height = rect.height;
      var dpr = Math.min(window.devicePixelRatio || 1, 1.5);
      canvas.width = Math.max(1, Math.floor(width * dpr));
      canvas.height = Math.max(1, Math.floor(height * dpr));
      ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
      draw();
    }

    function filterToPost(graph) {
      var postId = null;
      var hrefMap = new Map();

      graph.nodes.forEach(function(node) {
        if (node.href) {
          hrefMap.set(node.href, node.id);
        }
      });

      if (hrefMap.has(postHref)) {
        postId = hrefMap.get(postHref);
      }

      if (!postId) {
        return { nodes: [], edges: [], postId: null };
      }

      var connected = new Set([postId]);
      graph.edges.forEach(function(edge) {
        if (edge.source === postId) {
          connected.add(edge.target);
        }
        if (edge.target === postId) {
          connected.add(edge.source);
        }
      });

      return {
        nodes: graph.nodes.filter(function(node) { return connected.has(node.id); }),
        edges: graph.edges.filter(function(edge) {
          return connected.has(edge.source) && connected.has(edge.target);
        }),
        postId: postId
      };
    }

    function buildPreview(limit) {
      if (destroyed || !graphData) {
        return;
      }

      var subset = filterToPost(graphData);
      if (!subset.postId) {
        section.style.display = 'none';
        return;
      }

      var currentPost = subset.nodes.find(function(node) { return node.id === subset.postId; });
      if (!currentPost) {
        section.style.display = 'none';
        return;
      }

      section.style.display = '';
      section.classList.remove('post-graph--pending');

      var connectedNodes = subset.nodes.filter(function(node) { return node.id !== subset.postId; });
      var totalConnections = connectedNodes.length;
      connectedNodes.sort(function(a, b) {
        var aScore = a.type === 'tag' ? (a.count || 0) * 1000 : (Date.parse(a.date || '') || 0);
        var bScore = b.type === 'tag' ? (b.count || 0) * 1000 : (Date.parse(b.date || '') || 0);
        return bScore - aScore;
      });

      var effectiveLimit = limit;
      if (limitInput) {
        var minValue = parseInt(limitInput.min, 10);
        var maxValue = Math.max(minValue, totalConnections);
        if (!limitInput.dataset.defaultSet) {
          var rawDefault = limitInput.dataset.default;
          if (rawDefault === 'auto') {
            effectiveLimit = totalConnections;
          } else if (rawDefault) {
            var parsedDefault = parseInt(rawDefault, 10);
            if (!isNaN(parsedDefault)) {
              effectiveLimit = parsedDefault;
            }
          }
          limitInput.dataset.defaultSet = 'true';
        }

        effectiveLimit = Math.max(minValue, Math.min(effectiveLimit, maxValue));
        limitInput.max = String(maxValue);
        if (parseInt(limitInput.value, 10) !== effectiveLimit) {
          limitInput.value = String(effectiveLimit);
        }
        if (limitValue) {
          limitValue.textContent = String(effectiveLimit);
        }

        var control = limitInput.closest('.post-graph__control') || limitInput.parentElement;
        if (control) {
          var shouldHide = totalConnections <= minValue;
          control.hidden = shouldHide;
          control.style.display = shouldHide ? 'none' : '';
        }
      }

      var cappedLimit = Math.min(effectiveLimit, totalConnections);
      var limitedNodes = connectedNodes.slice(0, cappedLimit);
      var keep = new Set(limitedNodes.map(function(node) { return node.id; }));
      keep.add(currentPost.id);

      nodes = subset.nodes.filter(function(node) { return keep.has(node.id); }).map(function(node) {
        return {
          id: node.id,
          label: node.label,
          count: node.count || 1,
          href: node.href,
          type: node.type,
          x: width / 2,
          y: height / 2,
          vx: 0,
          vy: 0,
          r: node.type === 'tag' ? 2 + Math.min(4, Math.sqrt(node.count || 1)) : 3
        };
      });

      var nodeIds = new Set(nodes.map(function(node) { return node.id; }));
      edges = subset.edges.filter(function(edge) {
        return nodeIds.has(edge.source) && nodeIds.has(edge.target);
      }).map(function(edge) {
        return {
          source: edge.source,
          target: edge.target,
          weight: edge.weight || 1,
          type: edge.type
        };
      });

      startSimulation();
    }

    function startSimulation() {
      if (destroyed || !window.d3) {
        return;
      }

      if (simulation) {
        simulation.stop();
      }

      simulation = window.d3.forceSimulation(nodes)
        .force('link', window.d3.forceLink(edges)
          .id(function(d) { return d.id; })
          .distance(function(d) { return d.type === 'co-occurrence' ? 120 : 90; })
          .strength(function(d) { return d.type === 'co-occurrence' ? 0.35 : 0.6; }))
        .force('charge', window.d3.forceManyBody().strength(-180).distanceMax(260))
        .force('center', window.d3.forceCenter(width / 2, height / 2))
        .force('collision', window.d3.forceCollide().radius(function(d) { return d.r + 6; }));

      if (animateToggle && animateToggle.checked) {
        simulation.on('tick', draw);
        simulation.alpha(1).restart();
      } else {
        simulation.on('tick', null);
        simulation.alpha(1);
        for (var i = 0; i < 160; i++) {
          simulation.tick();
        }
        simulation.stop();
        draw();
      }
    }

    function draw() {
      if (destroyed) {
        return;
      }

      ctx.clearRect(0, 0, width, height);
      ctx.save();
      ctx.translate(panX, panY);
      ctx.scale(zoom, zoom);

      edges.forEach(function(edge) {
        var weight = Math.min(2, Math.max(0.4, edge.weight / 4));
        if (edge.type === 'co-occurrence') {
          ctx.lineWidth = 0.35;
          ctx.strokeStyle = 'rgba(255, 255, 255, 0.05)';
        } else {
          ctx.lineWidth = weight;
          ctx.strokeStyle = 'rgba(255, 205, 17, 0.08)';
        }
        ctx.beginPath();
        ctx.moveTo(edge.source.x, edge.source.y);
        ctx.lineTo(edge.target.x, edge.target.y);
        ctx.stroke();
      });

      nodes.forEach(function(node) {
        ctx.beginPath();
        if (node.type === 'post') {
          ctx.strokeStyle = accentColor;
          ctx.lineWidth = 1;
          ctx.globalAlpha = 0.75;
          var size = node.r + 3;
          ctx.rect(node.x - size / 2, node.y - size / 2, size, size);
          ctx.stroke();
        } else {
          ctx.fillStyle = accentColor;
          ctx.globalAlpha = 0.5;
          ctx.arc(node.x, node.y, node.r, 0, Math.PI * 2);
          ctx.fill();
        }
        ctx.globalAlpha = 1;
      });

      if (labelsToggle && labelsToggle.checked) {
        ctx.fillStyle = mutedColor;
        ctx.font = '10px var(--font-display, "Space Grotesk"), ui-sans-serif, system-ui';
        nodes.forEach(function(node) {
          ctx.fillText(node.label, node.x + node.r + 4, node.y + 3);
        });
      }

      if (hoveredNode) {
        ctx.beginPath();
        ctx.strokeStyle = '#fff';
        ctx.lineWidth = 1.2;
        if (hoveredNode.type === 'post') {
          var highlightSize = hoveredNode.r + 6;
          ctx.rect(hoveredNode.x - highlightSize / 2, hoveredNode.y - highlightSize / 2, highlightSize, highlightSize);
          ctx.stroke();
        } else {
          ctx.arc(hoveredNode.x, hoveredNode.y, hoveredNode.r + 3, 0, Math.PI * 2);
          ctx.stroke();
        }
      }

      ctx.restore();
    }

    function setTooltip(node, x, y) {
      if (!tooltip) {
        return;
      }

      if (!node) {
        tooltip.style.opacity = '0';
        tooltip.textContent = '';
        return;
      }

      tooltip.textContent = node.label;
      tooltip.style.left = x + 'px';
      tooltip.style.top = y + 'px';
      tooltip.style.opacity = '1';
    }

    function findNodeAt(x, y) {
      for (var i = nodes.length - 1; i >= 0; i--) {
        var node = nodes[i];
        var dx = x - node.x;
        var dy = y - node.y;
        var radius = node.r + 4;
        if (Math.sqrt(dx * dx + dy * dy) <= radius) {
          return node;
        }
      }
      return null;
    }

    // Convert a mouse event to canvas-space CSS pixels. The drawing surface
    // is sized from the parent at resize time; if the layout has since
    // shifted (drawers, scrollbars, fonts) the canvas is displayed stretched,
    // so scale by the ratio between the logical size and the rendered box.
    function canvasPoint(event) {
      var rect = canvas.getBoundingClientRect();
      var sx = rect.width ? width / rect.width : 1;
      var sy = rect.height ? height / rect.height : 1;
      return {
        x: (event.clientX - rect.left) * sx,
        y: (event.clientY - rect.top) * sy,
        cssX: event.clientX - rect.left,
        cssY: event.clientY - rect.top
      };
    }

    function onPointerMove(event) {
      var point = canvasPoint(event);
      var x = (point.x - panX) / zoom;
      var y = (point.y - panY) / zoom;
      hoveredNode = findNodeAt(x, y);
      setTooltip(hoveredNode, point.cssX, point.cssY);
      draw();
    }

    function onPointerLeave() {
      hoveredNode = null;
      setTooltip(null);
      draw();
    }

    function onClick() {
      if (hoveredNode && hoveredNode.href) {
        window.location.href = hoveredNode.href;
      }
    }

    function onWheel(event) {
      event.preventDefault();
      var point = canvasPoint(event);
      var mx = point.x;
      var my = point.y;
      var delta = event.deltaY > 0 ? 0.92 : 1.08;
      var nextZoom = Math.max(0.4, Math.min(3, zoom * delta));
      panX = mx - (mx - panX) * (nextZoom / zoom);
      panY = my - (my - panY) * (nextZoom / zoom);
      zoom = nextZoom;
      draw();
    }

    function onPointerDown(event) {
      var point = canvasPoint(event);
      var x = (point.x - panX) / zoom;
      var y = (point.y - panY) / zoom;
      dragNode = findNodeAt(x, y);
      if (dragNode) {
        dragNode.fx = dragNode.x;
        dragNode.fy = dragNode.y;
        if (simulation) {
          simulation.alphaTarget(0.2).restart();
        }
      } else {
        dragNode = { pan: true, startX: event.clientX, startY: event.clientY, baseX: panX, baseY: panY };
      }
    }

    function onPointerMoveDrag(event) {
      if (!dragNode) {
        return;
      }

      if (dragNode.pan) {
        panX = dragNode.baseX + (event.clientX - dragNode.startX);
        panY = dragNode.baseY + (event.clientY - dragNode.startY);
      } else {
        var point = canvasPoint(event);
        dragNode.fx = (point.x - panX) / zoom;
        dragNode.fy = (point.y - panY) / zoom;
      }

      draw();
    }

    function onPointerUp() {
      if (dragNode && !dragNode.pan) {
        dragNode.fx = null;
        dragNode.fy = null;
        if (simulation) {
          simulation.alphaTarget(0);
        }
      }
      dragNode = null;
    }

    function onLimitInput() {
      if (limitValue) {
        limitValue.textContent = limitInput.value;
      }
      if (rebuildTimer) {
        window.clearTimeout(rebuildTimer);
      }
      rebuildTimer = window.setTimeout(rebuild, 150);
    }

    function onAnimateChange() {
      startSimulation();
    }

    function onReducedMotionChange(event) {
      if (animateToggle) {
        animateToggle.checked = !event.matches;
      }
      startSimulation();
    }

    function rebuild() {
      if (!graphData || destroyed) {
        return;
      }
      var limit = parseInt(limitInput ? limitInput.value : '18', 10);
      if (limitValue) {
        limitValue.textContent = String(limit);
      }
      buildPreview(limit);
    }

    function onResize() {
      resizeCanvas();
      rebuild();
    }

    canvas.addEventListener('mousemove', onPointerMove);
    canvas.addEventListener('mouseleave', onPointerLeave);
    canvas.addEventListener('click', onClick);
    canvas.addEventListener('wheel', onWheel, { passive: false });
    canvas.addEventListener('mousedown', onPointerDown);
    window.addEventListener('mousemove', onPointerMoveDrag);
    window.addEventListener('mouseup', onPointerUp);
    window.addEventListener('resize', onResize);

    // The window may not resize when the graph's container does (sidebar
    // drawers opening, fonts loading, scrollbars appearing). Track the
    // container itself so the drawing surface always matches the box.
    var resizeObserver = null;
    var resizeFrame = 0;
    if (typeof ResizeObserver === 'function' && canvas.parentElement) {
      resizeObserver = new ResizeObserver(function() {
        if (resizeFrame) {
          return;
        }
        resizeFrame = window.requestAnimationFrame(function() {
          resizeFrame = 0;
          if (destroyed) {
            return;
          }
          var rect = canvas.parentElement.getBoundingClientRect();
          if (Math.abs(rect.width - width) > 0.5 || Math.abs(rect.height - height) > 0.5) {
            onResize();
          }
        });
      });
      resizeObserver.observe(canvas.parentElement);
    }

    if (limitInput) {
      limitInput.addEventListener('input', onLimitInput);
    }

    if (animateToggle) {
      animateToggle.addEventListener('change', onAnimateChange);
    }

    if (animateToggle && prefersReducedMotion && prefersReducedMotion.matches) {
      animateToggle.checked = false;
    }

    if (prefersReducedMotion) {
      motionListener = onReducedMotionChange;
      if (typeof prefersReducedMotion.addEventListener === 'function') {
        prefersReducedMotion.addEventListener('change', motionListener);
      }
    }

    registerCleanup(function() {
      destroyed = true;
      if (rebuildTimer) {
        window.clearTimeout(rebuildTimer);
      }
      if (simulation) {
        simulation.stop();
      }

      canvas.removeEventListener('mousemove', onPointerMove);
      canvas.removeEventListener('mouseleave', onPointerLeave);
      canvas.removeEventListener('click', onClick);
      canvas.removeEventListener('wheel', onWheel);
      canvas.removeEventListener('mousedown', onPointerDown);
      window.removeEventListener('mousemove', onPointerMoveDrag);
      window.removeEventListener('mouseup', onPointerUp);
      window.removeEventListener('resize', onResize);
      if (resizeObserver) {
        resizeObserver.disconnect();
      }
      if (resizeFrame) {
        window.cancelAnimationFrame(resizeFrame);
      }

      if (limitInput) {
        limitInput.removeEventListener('input', onLimitInput);
      }
      if (animateToggle) {
        animateToggle.removeEventListener('change', onAnimateChange);
      }
      if (prefersReducedMotion && motionListener && typeof prefersReducedMotion.removeEventListener === 'function') {
        prefersReducedMotion.removeEventListener('change', motionListener);
      }
    });

    resizeCanvas();

    if (!graphUrl) {
      section.style.display = 'none';
      return;
    }

    fetch(graphUrl)
      .then(function(resp) { return resp.json(); })
      .then(function(data) {
        if (destroyed || !section.isConnected) {
          return;
        }
        graphData = data;
        rebuild();
      })
      .catch(function() {
        if (!destroyed && section.isConnected) {
          section.style.display = 'none';
        }
      });
  }

  function initPostGraph() {
    cleanupPostGraphs();

    var sections = document.querySelectorAll('.post-graph');
    if (!sections.length) {
      return;
    }

    loadD3()
      .then(function() {
        document.querySelectorAll('.post-graph').forEach(initSection);
      })
      .catch(function() {
        document.querySelectorAll('.post-graph').forEach(function(section) {
          section.style.display = 'none';
        });
      });
  }

  function scheduleInitPostGraph() {
    if (initFrame !== null) {
      (window.cancelAnimationFrame || window.clearTimeout)(initFrame);
    }

    var schedule = window.requestAnimationFrame || function(callback) {
      return window.setTimeout(callback, 16);
    };

    initFrame = schedule(function() {
      initFrame = null;
      initPostGraph();
    });
  }

  window.initPostGraph = scheduleInitPostGraph;

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', scheduleInitPostGraph);
  } else {
    scheduleInitPostGraph();
  }

  window.addEventListener('view-transition-complete', scheduleInitPostGraph);
})();
