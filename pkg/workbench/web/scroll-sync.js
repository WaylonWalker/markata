// Source-line anchors keep documents with unequal rendered heights aligned.
window.MarkataScrollSync = function (editor, frame, setting, workbench) {
  let lines = [], preview = null, driver = 'editor', expectedEditorTop = null, markedLine = -1;
  let framePending = false, editorPending = false;
  try { setting.checked = localStorage.getItem('markata-editor-scroll-sync') !== 'false'; } catch { /* storage is optional */ }
  const active = () => setting.checked && workbench.dataset.layout === 'split';
  const clamp = (n, max) => Math.max(0, Math.min(max, n));
  const position = (points, value) => {
    if (points.length < 2) return null;
    if (value <= points[0].y) return points[0].line;
    for (let i = 1; i < points.length; i++) {
      if (value <= points[i].y) {
        const a = points[i - 1], b = points[i];
        return a.line + (b.line - a.line) * (value - a.y) / Math.max(1, b.y - a.y);
      }
    }
    return points.at(-1).line;
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
  const editorPoints = () => {
    const info = editor.getScrollInfo();
    const at = line => editor.heightAtLine(clamp(line, editor.lastLine()), 'local');
    const raw = [{line:0,y:at(0)}, ...lines.map(line => ({line,y:at(line)})), {line:editor.lineCount(),y:info.height}];
    return raw.reduce((points, point) => {
      const previous = points.at(-1);
      if (Number.isFinite(point.line) && Number.isFinite(point.y) && (!previous || (point.line > previous.line && point.y > previous.y))) points.push(point);
      return points;
    }, []);
  };
  function mark(line) {
    line = clamp(Math.round(line || 0), editor.lastLine());
    if (markedLine >= 0) editor.removeLineClass(markedLine, 'background', 'markata-sync-cursor');
    markedLine = line;
    editor.addLineClass(markedLine, 'background', 'markata-sync-cursor');
    frame.contentWindow?.postMessage({type:'sync-cursor', line}, '*');
  }
  function editorToPreview(followCursor = false) {
    if (!active() || !preview || !lines.length) return;
    const info = editor.getScrollInfo();
    const cursorY = editor.heightAtLine(editor.getCursor().line, 'local');
    const preferred = followCursor ? Math.max(.15, Math.min(.8, (cursorY - info.top) / info.clientHeight)) : .4;
    const sourcePoints = editorPoints();
    const line = followCursor ? editor.getCursor().line : position(sourcePoints, info.top + info.clientHeight * preferred);
    if (line != null) mark(line);
    const sourceY = height(sourcePoints, line), previewY = height(preview.points, line);
    if (sourceY != null && previewY != null) {
      const ratio = sharedRatio(sourceY, info.height - info.clientHeight, info.clientHeight, previewY, preview.height - preview.clientHeight, preview.clientHeight, preferred);
      placeEditor(sourceY, ratio, info);
      frame.contentWindow?.postMessage({type:'sync-scroll', y:clamp(previewY - preview.clientHeight * ratio, preview.height - preview.clientHeight)}, '*');
    }
  }
  function previewToEditor() {
    if (!active() || !preview || !lines.length) return;
    const info = editor.getScrollInfo(), preferred = .4;
    const line = position(preview.points, preview.top + preview.clientHeight * preferred);
    if (line != null) mark(line);
    const editorY = height(editorPoints(), line), previewY = height(preview.points, line);
    if (editorY != null && previewY != null) {
      const ratio = sharedRatio(editorY, info.height - info.clientHeight, info.clientHeight, previewY, preview.height - preview.clientHeight, preview.clientHeight, preferred);
      placeEditor(editorY, ratio, info);
      frame.contentWindow?.postMessage({type:'sync-scroll', y:clamp(previewY - preview.clientHeight * ratio, preview.height - preview.clientHeight)}, '*');
    }
  }
  editor.on('scroll', () => {
    if (expectedEditorTop != null) {
      const actual = editor.getScrollInfo().top;
      if (Math.abs(actual - expectedEditorTop) < 2) { expectedEditorTop = null; return; }
      expectedEditorTop = null;
    }
    if (!active()) return;
    driver = 'editor';
    if (!editorPending) { editorPending = true; requestAnimationFrame(() => { editorPending = false; editorToPreview(); }); }
  });
  editor.on('change', () => { driver = 'editor'; requestAnimationFrame(() => editorToPreview(true)); });
  editor.on('cursorActivity', () => { if (editor.hasFocus()) { driver = 'editor'; mark(editor.getCursor().line); requestAnimationFrame(() => editorToPreview(true)); } });
  window.addEventListener('message', event => {
    if (event.source !== frame.contentWindow || event.data?.type !== 'sync-scroll-state') return;
    preview = event.data;
    if (event.data.user) driver = 'preview';
    if (event.data.user && !framePending) {
      framePending = true;
      requestAnimationFrame(() => { framePending = false; previewToEditor(); });
    } else if (!event.data.user && driver === 'editor') editorToPreview();
  });
  setting.addEventListener('change', () => {
    try { localStorage.setItem('markata-editor-scroll-sync', String(setting.checked)); } catch { /* storage is optional */ }
    if (setting.checked) { driver = 'editor'; editorToPreview(); }
  });
  return {
    update(sourceLines) { lines = sourceLines || []; frame.contentWindow?.postMessage({type:'sync-scroll-lines', lines, lineCount:editor.lineCount()}, '*'); },
    ready() { frame.contentWindow?.postMessage({type:'sync-scroll-lines', lines, lineCount:editor.lineCount()}, '*'); },
    refresh() { if (driver === 'preview') previewToEditor(); else editorToPreview(); }
  };
};
