(function () {
  'use strict';

  function bindSparklineReadout(wrap, svg, value) {
    if (!wrap || !svg || !value || svg.dataset.sparklineBound === 'true') return;

    var counts = (svg.dataset.values || '').split(',').map(Number);
    if (counts.length < 2 || counts.some(function (count) { return !Number.isFinite(count); })) return;

    var start = (svg.dataset.startMonth || '').split('-').map(Number);
    if (start.length !== 2 || !start[0] || !start[1]) return;

    var defaultValue = value.textContent;
    var activeIndex = counts.length - 1;
    var guide = svg.querySelector('line');
    var dot = svg.querySelector('circle');
    if (!guide || !dot) return;

    function show(index) {
      activeIndex = Math.max(0, Math.min(counts.length - 1, index));
      var x = activeIndex * 96 / (counts.length - 1);
      var max = Math.max.apply(null, counts);
      var y = 28 - (counts[activeIndex] / max) * 24;
      guide.setAttribute('x1', x);
      guide.setAttribute('x2', x);
      guide.setAttribute('y1', 0);
      guide.setAttribute('y2', 28);
      dot.setAttribute('cx', x);
      dot.setAttribute('cy', y);

      var date = new Date(Date.UTC(start[0], start[1] - 1 + activeIndex, 1));
      var month = new Intl.DateTimeFormat(undefined, {
        month: 'short', year: 'numeric', timeZone: 'UTC'
      }).format(date);
      value.textContent = month + ' | ' + counts[activeIndex] + ' posts';
      svg.classList.add('is-active');
    }

    function reset() {
      svg.classList.remove('is-active');
      value.textContent = defaultValue;
    }

    svg.addEventListener('pointermove', function (event) {
      var point = svg.createSVGPoint();
      point.x = event.clientX;
      point.y = event.clientY;
      var matrix = svg.getScreenCTM();
      if (!matrix) return;
      var chartPoint = point.matrixTransform(matrix.inverse());
      show(Math.round(chartPoint.x / 96 * (counts.length - 1)));
    });

    svg.addEventListener('pointerleave', function () {
      if (document.activeElement !== svg) reset();
    });

    svg.addEventListener('focus', function () { show(activeIndex); });
    svg.addEventListener('blur', reset);
    svg.addEventListener('keydown', function (event) {
      var next = activeIndex;
      if (event.key === 'ArrowLeft' || event.key === 'ArrowDown') next -= 1;
      else if (event.key === 'ArrowRight' || event.key === 'ArrowUp') next += 1;
      else if (event.key === 'Home') next = 0;
      else if (event.key === 'End') next = counts.length - 1;
      else return;
      event.preventDefault();
      show(next);
    });

    svg.dataset.sparklineBound = 'true';
  }

  function init() {
    document.querySelectorAll('.feed-sparkline-wrap').forEach(function (wrap) {
      bindSparklineReadout(wrap, wrap.querySelector('.feed-sparkline'), wrap.querySelector('.feed-sparkline-value'));
    });
    document.querySelectorAll('.feed-header-sparkline').forEach(function (wrap) {
      bindSparklineReadout(wrap, wrap.querySelector('.feed-header-sparkline-chart'), wrap.querySelector('.feed-header-sparkline-value'));
    });
  }

  document.addEventListener('DOMContentLoaded', init);
  window.addEventListener('view-transition-complete', init);
})();
