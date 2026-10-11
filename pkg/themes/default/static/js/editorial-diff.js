(function () {
  'use strict';
  var maxBytes = 512 * 1024;

  async function source(url, signal) {
    var parsed = new URL(url);
    if (parsed.protocol !== 'https:' || parsed.hostname !== 'raw.githubusercontent.com' || parsed.username || parsed.password) {
      throw new Error('Unsupported revision source');
    }
    var response = await fetch(url, { signal: signal, credentials: 'omit', referrerPolicy: 'no-referrer' });
    if (!response.ok || Number(response.headers.get('content-length')) > maxBytes) throw new Error('Revision unavailable');
    if (!response.body) {
      var text = await response.text();
      if (new TextEncoder().encode(text).length > maxBytes) throw new Error('Revision too large');
      return text;
    }
    var reader = response.body.getReader();
    var chunks = [];
    var size = 0;
    while (true) {
      var part = await reader.read();
      if (part.done) break;
      size += part.value.length;
      if (size > maxBytes) { await reader.cancel(); throw new Error('Revision too large'); }
      chunks.push(part.value);
    }
    var bytes = new Uint8Array(size);
    var offset = 0;
    chunks.forEach(function (chunk) { bytes.set(chunk, offset); offset += chunk.length; });
    return new TextDecoder().decode(bytes);
  }

  function body(text) {
    var normalized = text.replace(/\r\n?/g, '\n').replace(/^\uFEFF/, '');
    var lines = normalized.split('\n');
    if (lines[0].trim() !== '---') return normalized;
    for (var i = 1; i < lines.length; i++) {
      if (lines[i].trim() === '---') return lines.slice(i + 1).join('\n');
    }
    throw new Error('Unclosed revision frontmatter');
  }

  function tokens(text) { return text.match(/\s+|[\p{L}\p{N}_]+|[^\s]/gu) || []; }

  // Bound the middle comparison; common prefixes/suffixes keep small edits cheap.
  function compare(before, after) {
    var a = tokens(before), b = tokens(after);
    var start = 0;
    while (start < a.length && start < b.length && a[start] === b[start]) start++;
    var endA = a.length, endB = b.length;
    while (endA > start && endB > start && a[endA - 1] === b[endB - 1]) { endA--; endB--; }
    var pieces = [];
    function push(kind, text) {
      if (!text) return;
      var last = pieces[pieces.length - 1];
      if (last && last.kind === kind) last.text += text;
      else pieces.push({ kind: kind, text: text });
    }
    push('equal', a.slice(0, start).join(''));
    var n = endA - start, m = endB - start;
    if ((n + 1) * (m + 1) > 1000000) {
      push('removed', a.slice(start, endA).join(''));
      push('added', b.slice(start, endB).join(''));
    } else {
      var width = m + 1;
      var table = new Uint32Array((n + 1) * width);
      for (var i = n - 1; i >= 0; i--) {
        for (var j = m - 1; j >= 0; j--) {
          table[i * width + j] = a[start + i] === b[start + j]
            ? 1 + table[(i + 1) * width + j + 1]
            : Math.max(table[(i + 1) * width + j], table[i * width + j + 1]);
        }
      }
      var x = 0, y = 0;
      while (x < n || y < m) {
        if (x < n && y < m && a[start + x] === b[start + y]) { push('equal', a[start + x++]); y++; }
        else if (x < n && (y === m || table[(x + 1) * width + y] >= table[x * width + y + 1])) push('removed', a[start + x++]);
        else push('added', b[start + y++]);
      }
    }
    push('equal', a.slice(endA).join(''));
    return pieces;
  }

  async function show(target, beforeURL, afterURL) {
    var controller = new AbortController();
    var timer = setTimeout(function () { controller.abort(); }, 20000);
    try {
      var revisions = await Promise.all([source(beforeURL, controller.signal), source(afterURL, controller.signal)]);
      var before = body(revisions[0]), after = body(revisions[1]);
      target.replaceChildren();
      if (before === after) {
        var message = document.createElement('p');
        message.textContent = 'No Markdown body changes in this revision. The edit may only affect frontmatter.';
        target.appendChild(message);
        return;
      }
      var content = document.createElement('div');
      content.className = 'editorial-history__source';
      compare(before, after).forEach(function (piece) {
        if (piece.kind === 'equal') content.appendChild(document.createTextNode(piece.text));
        else {
          var change = document.createElement(piece.kind === 'added' ? 'ins' : 'del');
          change.textContent = piece.text;
          content.appendChild(change);
        }
      });
      target.appendChild(content);
    } finally { clearTimeout(timer); controller.abort(); }
  }
  window.MarkataEditorialDiff = { show: show, compare: compare, body: body };
}());
