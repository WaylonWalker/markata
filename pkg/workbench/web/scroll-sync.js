// Source-line anchors keep documents with unequal rendered heights aligned.
window.MarkataScrollSync = function (editor, frame, setting, workbench) {
  let lines = [], preview = null, driver = 'editor', suppressEditor = false;
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
  const editorPoints = () => {
    const info = editor.getScrollInfo();
    const at = line => editor.heightAtLine(clamp(line, editor.lastLine()), 'local');
    return [{line:0,y:at(0)}, ...lines.map(line => ({line,y:at(line)})), {line:editor.lineCount(),y:info.height}]
      .filter((point, i, all) => i === 0 || point.line > all[i - 1].line);
  };
  function editorToPreview(followCursor = false) {
    if (!active() || !preview || !lines.length) return;
    const info = editor.getScrollInfo();
    const cursorY = editor.heightAtLine(editor.getCursor().line, 'local');
    const ratio = followCursor ? Math.max(.15, Math.min(.8, (cursorY - info.top) / info.clientHeight)) : .4;
    const line = followCursor ? editor.getCursor().line : position(editorPoints(), info.top + info.clientHeight * ratio);
    const y = height(preview.points, line);
    if (y != null) frame.contentWindow?.postMessage({type:'sync-scroll', y:clamp(y - preview.clientHeight * ratio, preview.height - preview.clientHeight)}, '*');
  }
  function previewToEditor() {
    if (!active() || !preview || !lines.length) return;
    const info = editor.getScrollInfo(), ratio = .4;
    const line = position(preview.points, preview.top + preview.clientHeight * ratio);
    const y = height(editorPoints(), line);
    if (y != null) {
      suppressEditor = true;
      editor.scrollTo(null, clamp(y - info.clientHeight * ratio, info.height - info.clientHeight));
      requestAnimationFrame(() => { suppressEditor = false; });
    }
  }
  editor.on('scroll', () => {
    if (suppressEditor || !active()) return;
    driver = 'editor';
    if (!editorPending) { editorPending = true; requestAnimationFrame(() => { editorPending = false; editorToPreview(); }); }
  });
  editor.on('change', () => { driver = 'editor'; requestAnimationFrame(() => editorToPreview(true)); });
  editor.on('cursorActivity', () => { if (editor.hasFocus()) { driver = 'editor'; requestAnimationFrame(() => editorToPreview(true)); } });
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
