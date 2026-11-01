// Source-line anchors keep documents with unequal rendered heights aligned.
window.MarkataScrollSync = function (editor, frame, setting, workbench) {
  let lines = [], anchors = [], sourceOffset = 0, preview = null, driver = 'editor', expectedEditorTop = null;
  let editorFollowCursor = true, lastMarkedLine = -1, lastCursorLine = -1;
  let geometryKey = '', requestedPreview = null, previewRequestID = 0, editorScrollIntentUntil = 0;
  let contentUnrendered = false, deferredCursorFollow = false;
  let framePending = false, editorPending = false;
  const editorWrapper = editor.getWrapperElement();
  const noteEditorScrollIntent = duration => { editorScrollIntentUntil = performance.now() + duration; };
  editorWrapper.addEventListener('wheel', () => noteEditorScrollIntent(250), {passive:true});
  editorWrapper.addEventListener('touchstart', () => noteEditorScrollIntent(500), {passive:true});
  editorWrapper.addEventListener('touchmove', () => noteEditorScrollIntent(250), {passive:true});
  editorWrapper.addEventListener('pointerdown', () => noteEditorScrollIntent(100), {passive:true});
  editorWrapper.addEventListener('pointermove', event => { if (event.buttons) noteEditorScrollIntent(250); }, {passive:true});
  editorWrapper.addEventListener('pointerup', () => noteEditorScrollIntent(100), {passive:true});
  window.addEventListener('keydown', event => {
    if (editor.hasFocus() && ['ArrowUp','ArrowDown','PageUp','PageDown','Home','End'].includes(event.key)) noteEditorScrollIntent(250);
  }, true);
  try { setting.checked = localStorage.getItem('markata-editor-scroll-sync') !== 'false'; } catch { /* storage is optional */ }
  const active = () => setting.checked && workbench.dataset.layout === 'split';
  const clamp = (n, max) => Math.max(0, Math.min(max, n));
  const position = (points, value, preferredLine = null) => {
    if (points.length < 2) return null;
    const lineAt = (start, end) => {
      let selected = end;
      if (Number.isFinite(preferredLine)) {
        let distance = Infinity;
        for (let i = start; i <= end; i++) {
          const nextDistance = Math.abs(points[i].line - preferredLine);
          if (nextDistance < distance || (nextDistance === distance && points[i].line > points[selected].line)) {
            selected = i;
            distance = nextDistance;
          }
        }
      }
      return points[selected].line;
    };
    let start = 0, end = 0;
    while (end + 1 < points.length && points[end + 1].y === points[start].y) end++;
    let previousY = points[start].y, previousLine = lineAt(start, end);
    if (value <= previousY) return previousLine;
    start = end + 1;
    while (start < points.length) {
      end = start;
      while (end + 1 < points.length && points[end + 1].y === points[start].y) end++;
      const nextY = points[start].y, nextLine = lineAt(start, end);
      if (value <= nextY) {
        return previousLine + (nextLine - previousLine) * (value - previousY) / Math.max(1, nextY - previousY);
      }
      previousY = nextY;
      previousLine = nextLine;
      start = end + 1;
    }
    return previousLine;
  };
  const height = (points, line) => {
    if (points.length < 2) return null;
    if (line <= points[0].line) return points[0].y;
    for (let i = 1; i < points.length; i++) {
      if (line <= points[i].line) {
        const a = points[i - 1], b = points[i];
        return a.y + (b.y - a.y) * (line - a.line) / Math.max(1, b.line - a.line);
      }
    }
    return points.at(-1).y;
  };
  const sharedRatio = (a, aMax, aHeight, b, bMax, bHeight, preferred) => {
    const low = Math.max(0, (a - aMax) / aHeight, (b - bMax) / bHeight);
    const high = Math.min(1, a / aHeight, b / bHeight);
    return low <= high ? Math.max(low, Math.min(high, preferred)) : preferred;
  };
  const placeEditor = (y, ratio, info) => {
    const top = clamp(y - info.clientHeight * ratio, info.height - info.clientHeight);
    if (Math.abs(info.top - top) < 1) return;
    expectedEditorTop = top;
    editor.scrollTo(null, top);
    setTimeout(() => { if (expectedEditorTop === top) expectedEditorTop = null; }, 300);
  };
  const placePreview = y => {
    const top = clamp(y, preview.height - preview.clientHeight);
    if (Math.abs(preview.top - top) < 1 || (requestedPreview != null && Math.abs(requestedPreview.top - top) < 1)) return;
    const id = ++previewRequestID;
    requestedPreview = {id, top};
    frame.contentWindow?.postMessage({type:'sync-scroll', y:top, id}, '*');
    setTimeout(() => { if (requestedPreview?.id === id) requestedPreview = null; }, 300);
  };
  const editorPoints = () => {
    const info = editor.getScrollInfo();
    const at = line => editor.heightAtLine(clamp(line, editor.lastLine()), 'local');
    const sourceLines = [...new Set([...lines, ...anchors])].sort((a,b) => a-b);
    const raw = [{line:0,y:at(0)}, ...sourceLines.map(line => ({line,y:at(line)})), {line:editor.lineCount(),y:info.height}];
    return raw.reduce((points, point) => {
      const previous = points.at(-1);
      if (Number.isFinite(point.line) && Number.isFinite(point.y) && (!previous || (point.line > previous.line && point.y > previous.y))) points.push(point);
      return points;
    }, []);
  };
  function markPreviewPosition(line) {
    line = clamp(Math.round(line || 0), editor.lastLine());
    if (line === lastMarkedLine) return;
    lastMarkedLine = line;
    frame.contentWindow?.postMessage({type:'sync-position', line}, '*');
  }
  function editorToPreview(followCursor = false) {
    if (!active() || !preview || !lines.length) return;
    const info = editor.getScrollInfo();
    const cursor = editor.getCursor();
    const cursorY = editor.cursorCoords(cursor, 'local').top;
    const preferred = followCursor ? Math.max(.15, Math.min(.8, (cursorY - info.top) / info.clientHeight)) : .4;
    const sourcePoints = editorPoints();
    const line = followCursor ? cursor.line : position(sourcePoints, info.top + info.clientHeight * preferred);
    if (line != null) markPreviewPosition(line);
    const sourceY = height(sourcePoints, line), previewY = height(preview.points, line);
    if (sourceY != null && previewY != null) {
      const ratio = sharedRatio(sourceY, info.height - info.clientHeight, info.clientHeight, previewY, preview.height - preview.clientHeight, preview.clientHeight, preferred);
      placePreview(previewY - preview.clientHeight * ratio);
    }
  }
  function previewToEditor() {
    if (!active() || !preview || !lines.length) return;
    const info = editor.getScrollInfo(), preferred = .4;
    const line = position(preview.visual_points || preview.points, preview.top + preview.clientHeight * preferred, lastMarkedLine >= 0 ? lastMarkedLine : null);
    if (line != null) markPreviewPosition(line);
    const editorY = height(editorPoints(), line), previewY = height(preview.points, line);
    if (editorY != null && previewY != null) {
      const ratio = sharedRatio(editorY, info.height - info.clientHeight, info.clientHeight, previewY, preview.height - preview.clientHeight, preview.clientHeight, preferred);
      placeEditor(editorY, ratio, info);
    }
  }
  editor.on('scroll', () => {
    if (expectedEditorTop != null) {
      const actual = editor.getScrollInfo().top;
      if (Math.abs(actual - expectedEditorTop) < 2) { expectedEditorTop = null; return; }
      expectedEditorTop = null;
    }
    if (!active()) return;
    if (performance.now() > editorScrollIntentUntil) return;
    driver = 'editor';
    editorFollowCursor = false;
    deferredCursorFollow = false;
    if (!editorPending) { editorPending = true; requestAnimationFrame(() => { editorPending = false; editorToPreview(); }); }
  });
  editor.on('change', () => {
    contentUnrendered = true;
    if (!editor.hasFocus()) return;
    driver = 'editor';
    editorFollowCursor = true;
    deferredCursorFollow = true;
  });
  editor.on('cursorActivity', () => {
    if (!editor.hasFocus()) return;
    const cursor = editor.getCursor();
    if (cursor.line === lastCursorLine) return;
    driver = 'editor'; editorFollowCursor = true;
    lastCursorLine = cursor.line;
    markPreviewPosition(cursor.line);
    if (contentUnrendered) { deferredCursorFollow = true; return; }
    requestAnimationFrame(() => editorToPreview(true));
  });
  window.addEventListener('message', event => {
    if (event.source !== frame.contentWindow || event.data?.type !== 'sync-scroll-state') return;
    preview = event.data;
    if (requestedPreview && event.data.sync_id === requestedPreview.id) requestedPreview = null;
    const nextGeometryKey = `${preview.height}:${preview.clientHeight}:` + preview.points.map(point => `${point.line},${Math.round(point.y)}`).join(';');
    const geometryChanged = nextGeometryKey !== geometryKey;
    geometryKey = nextGeometryKey;
    if (event.data.user) driver = 'preview';
    if (event.data.user && !framePending) {
      framePending = true;
      requestAnimationFrame(() => { framePending = false; previewToEditor(); });
    } else if (!event.data.user && geometryChanged && !contentUnrendered) {
      requestAnimationFrame(() => driver === 'preview' ? previewToEditor() : editorToPreview(editorFollowCursor));
    }
  });
  setting.addEventListener('change', () => {
    try { localStorage.setItem('markata-editor-scroll-sync', String(setting.checked)); } catch { /* storage is optional */ }
    if (setting.checked) { driver = 'editor'; editorFollowCursor = false; editorToPreview(); }
  });
  return {
    update(sourceLines, notifyPreview = true, offset = 0, sourceAnchors = []) {
      lines = sourceLines || [];
      anchors = sourceAnchors || [];
      sourceOffset = offset;
      if (notifyPreview) frame.contentWindow?.postMessage({type:'sync-scroll-lines', lines, anchors, sourceOffset, lineCount:editor.lineCount()}, '*');
    },
    ready() { frame.contentWindow?.postMessage({type:'sync-scroll-lines', lines, anchors, sourceOffset, lineCount:editor.lineCount()}, '*'); },
    previewUpdated() {
      contentUnrendered = false;
      const followCursor = deferredCursorFollow || editorFollowCursor;
      deferredCursorFollow = false;
      requestAnimationFrame(() => driver === 'preview' ? previewToEditor() : editorToPreview(followCursor));
    },
    refresh() { if (driver === 'preview') previewToEditor(); else editorToPreview(); }
  };
};
