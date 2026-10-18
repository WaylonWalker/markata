// Update relative editorial headings when a static page outlives its build date.
(function () {
  'use strict';
  if (window.__markataEditorialRelativeWeeks) return;
  window.__markataEditorialRelativeWeeks = true;

  function mondayOf(date) {
    var monday = new Date(date.getFullYear(), date.getMonth(), date.getDate());
    monday.setDate(monday.getDate() - ((monday.getDay() + 6) % 7));
    return monday;
  }

  function relativeLabel(startKey, now) {
    var match = /^(\d{4})-(\d{2})-(\d{2})$/.exec(startKey || '');
    if (!match) return null;
    var year = Number(match[1]);
    var month = Number(match[2]) - 1;
    var day = Number(match[3]);
    var start = new Date(year, month, day);
    if (start.getFullYear() !== year || start.getMonth() !== month || start.getDate() !== day) return null;

    // Use UTC calendar midnights to count whole weeks without DST skew.
    var monday = mondayOf(now);
    var weeksAgo = Math.round((Date.UTC(monday.getFullYear(), monday.getMonth(), monday.getDate()) -
      Date.UTC(year, month, day)) / (7 * 86400000));
    if (weeksAgo < 0) return weeksAgo === -1 ? 'Next week' : 'In ' + (-weeksAgo) + ' weeks';
    if (weeksAgo === 0) return 'This week';
    if (weeksAgo === 1) return 'Last week';
    if (weeksAgo < 5) return weeksAgo + ' weeks ago';

    var monthsAgo = (now.getFullYear() - year) * 12 + now.getMonth() - month;
    if (monthsAgo < 1) return weeksAgo + ' weeks ago';
    if (monthsAgo < 12) return monthsAgo + (monthsAgo === 1 ? ' month ago' : ' months ago');
    var yearsAgo = Math.floor(monthsAgo / 12);
    return yearsAgo + (yearsAgo === 1 ? ' year ago' : ' years ago');
  }

  function refresh(root) {
    var scope = root && root.querySelectorAll ? root : document;
    var selector = '[data-editorial-relative-week][data-week-start]';
    var headings = Array.from(scope.querySelectorAll(selector));
    if (scope.matches && scope.matches(selector)) headings.unshift(scope);
    var now = new Date();
    headings.forEach(function (heading) {
      var label = relativeLabel(heading.dataset.weekStart, now);
      if (label) heading.textContent = label;
    });
  }

  document.addEventListener('htmx:afterSwap', function (event) {
    refresh(event.detail && event.detail.target || event.target);
  });
  window.addEventListener('view-transition-complete', function () { refresh(document); });
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', function () { refresh(document); }, { once: true });
  } else {
    refresh(document);
  }
}());
