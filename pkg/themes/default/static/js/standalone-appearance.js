    window.__markataSeasonalPick = function(forMode, date) {
      var s = window.__markataSeasonal, now = date || new Date(), y = now.getFullYear();
      var today = Date.UTC(y, now.getMonth(), now.getDate()) / 864e5;
      var best = null, bestDiff = 1e9;
      s.h.forEach(function(h) {
        h.d.forEach(function(d) {
          var p = d.split('-'), years = p.length === 3 ? [+p[0]] : [y - 1, y, y + 1];
          years.forEach(function(yy) {
            var diff = Date.UTC(yy, +p[p.length - 2] - 1, +p[p.length - 1]) / 864e5 - today;
            if (diff <= s.lead && diff > -(h.n || 1) && diff < bestDiff) { best = h; bestDiff = diff; }
          });
        });
      });
      if (!best) {
        var md = (now.getMonth() + 1) * 100 + now.getDate();
        best = s.s[s.s.length - 1];
        s.s.forEach(function(x) { if (md >= x.f) best = x; });
      }
      return { label: best.l, name: forMode === 'dark' ? best.dk : best.lt };
    };
/* Application bootstrap for the existing site appearance picker. */
(() => {
  const root = document.documentElement;
  const read = key => { try { return localStorage.getItem(key); } catch { return null; } };
  window.__markataThemeFallbackMode = 'dark';
  window.__markataThemeDefaults = {light:'default-light',dark:'default-dark'};
  window.__markataPaletteStylesReady = true;
  const mode = read('color-mode') || root.dataset.theme || 'dark';
  root.dataset.theme = mode === 'light' ? 'light' : 'dark';
  root.dataset.palette = read('theme-palette-' + root.dataset.theme) || window.__markataThemeDefaults[root.dataset.theme];
  root.dataset.aesthetic = read('theme-aesthetic') || root.dataset.aesthetic || 'minimal';
  root.dataset.fontpack = read('theme-fontpack') || root.dataset.fontpack || 'system';
  root.dataset.textSize = read('text-size') || root.dataset.textSize || 'large';
})();
