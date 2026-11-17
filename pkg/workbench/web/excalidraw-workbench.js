/* Excalidraw scene editing for the Go and FastAPI workbenches. */
(() => {
  'use strict';

  const currentTheme = () => document.documentElement.dataset.theme === 'light' ? 'light' : 'dark';
  const fenceMarkerLine = /^ {0,3}(?:`{3,}|~{3,})/;
  const emptyScene = () => ({
    type: 'excalidraw',
    version: 2,
    source: 'https://excalidraw.com',
    elements: [],
    appState: {viewBackgroundColor:'#ffffff',gridSize:null},
    files: {},
  });
  const dialog = document.querySelector('#excalidraw-dialog');
  const title = document.querySelector('#excalidraw-title');
  const hint = document.querySelector('#excalidraw-hint');
  const host = document.querySelector('#excalidraw-canvas');
  const loading = document.querySelector('#excalidraw-loading');
  const error = document.querySelector('#excalidraw-error');
  const applyButton = document.querySelector('#excalidraw-apply');
  const applyAsButton = document.querySelector('#excalidraw-apply-as');
  let cssReady;

  function openingFence(line) {
    const match = /^ {0,3}(`{3,}|~{3,})(.*)$/.exec(line);
    if (!match) return null;
    return {marker:match[1][0],length:match[1].length,info:match[2].trim()};
  }

  function closingFence(line, open) {
    const match = /^ {0,3}(`+|~+)\s*$/.exec(line);
    return !!match && match[1][0] === open.marker && match[1].length >= open.length;
  }

  function fenceAtCursor(editor) {
    const cursor = editor.getCursor();
    let open = null;
    for (let line = 0; line < editor.lineCount(); line++) {
      const text = editor.getLine(line);
      if (!open) {
        const candidate = openingFence(text);
        if (candidate) open = {...candidate,line};
        continue;
      }
      if (!closingFence(text,open)) continue;
      if (cursor.line >= open.line && cursor.line <= line) {
        const language = open.info.split(/\s+/,1)[0].toLowerCase();
        return {...open,closeLine:line,language};
      }
      open = null;
    }
    if (open && cursor.line >= open.line) {
      const language = open.info.split(/\s+/,1)[0].toLowerCase();
      return {...open,closeLine:null,language};
    }
    return null;
  }

  function ensureStyles() {
    if (cssReady) return cssReady;
    cssReady = new Promise((resolve,reject) => {
      const link = document.createElement('link');
      link.rel = 'stylesheet';
      link.href = '/editor/excalidraw-editor.css';
      link.dataset.excalidrawStyles = 'true';
      link.onload = resolve;
      link.onerror = () => reject(new Error('Could not load the local Excalidraw styles.'));
      document.head.append(link);
    });
    return cssReady;
  }

  function showError(message) {
    error.textContent = message;
    error.hidden = false;
    loading.hidden = true;
    host.hidden = true;
    applyButton.hidden = true;
    applyAsButton.hidden = true;
    if (!dialog.open) dialog.showModal();
  }

  function parseScene(source) {
    let scene;
    try {
      scene = JSON.parse(source);
    } catch (cause) {
      throw new Error(`This Excalidraw block is not valid JSON: ${cause.message}`);
    }
    if (!scene || typeof scene !== 'object' || scene.type !== 'excalidraw' || !Array.isArray(scene.elements)) {
      throw new Error('This block is not an Excalidraw scene. Its source has been kept unchanged.');
    }
    if (scene.appState != null && (typeof scene.appState !== 'object' || Array.isArray(scene.appState))) {
      throw new Error('This Excalidraw scene has invalid appState data. Its source has been kept unchanged.');
    }
    if (scene.files != null && (typeof scene.files !== 'object' || Array.isArray(scene.files))) {
      throw new Error('This Excalidraw scene has invalid image data. Its source has been kept unchanged.');
    }
    return scene;
  }

  function documentSize(source) {
    return new TextEncoder().encode(source).length;
  }

  function insertText(editor, position, text) {
    const source = editor.getValue();
    const offset = editor.indexFromPos(position);
    const before = source.slice(0,offset);
    const after = source.slice(offset);
    const prefix = !before ? '' : before.endsWith('\n\n') ? '' : before.endsWith('\n') ? '\n' : '\n\n';
    const suffix = !after ? '' : after.startsWith('\n\n') ? '' : after.startsWith('\n') ? '\n' : '\n\n';
    const block = '```excalidraw\n' + text + '\n```';
    const insertion = prefix + block + suffix;
    const next = before + insertion + after;
    if (documentSize(next) > 200000) {
      throw new Error('This drawing would make the Markdown document exceed 200,000 bytes. The document was not changed.');
    }
    editor.operation(() => {
      editor.changeGeneration(true);
      editor.replaceRange(insertion,position,position,'+excalidraw');
      editor.changeGeneration(true);
    });
  }

  function replaceSceneBody(editor, fence, text) {
    const from = {line:fence.line+1,ch:0};
    const to = {line:fence.closeLine,ch:0};
    const source = editor.getValue();
    const start = editor.indexFromPos(from);
    const end = editor.indexFromPos(to);
    const body = `${text}\n`;
    const next = source.slice(0,start) + body + source.slice(end);
    if (documentSize(next) > 200000) {
      throw new Error('This drawing would make the Markdown document exceed 200,000 bytes. The document was not changed.');
    }
    editor.operation(() => {
      editor.changeGeneration(true);
      editor.replaceRange(body,from,to,'+excalidraw');
      editor.changeGeneration(true);
    });
  }

  function attach(editor) {
    const button = document.querySelector('#drawing-open');
    let active = null;
    let canvas = null;
    const fenceRecords = new Map();
    let nextFenceID = 1;
    let refreshScheduled = false;

    function excalidrawFences() {
      const fences = [];
      let open = null;
      for (let line = 0; line < editor.lineCount(); line++) {
        const text = editor.getLine(line);
        if (!open) {
          const candidate = openingFence(text);
          if (candidate) open = {...candidate,line};
          continue;
        }
        if (!closingFence(text,open)) continue;
        if (open.info.split(/\s+/,1)[0].toLowerCase() === 'excalidraw') {
          fences.push({...open,closeLine:line,language:'excalidraw'});
        }
        open = null;
      }
      return fences;
    }

    function sceneSummary(fence) {
      const from = {line:fence.line+1,ch:0};
      const to = {line:fence.closeLine,ch:0};
      const body = editor.getRange(from,to).replace(/\r?\n$/,'');
      try {
        const scene = parseScene(body);
        const count = scene.elements.filter(element => !element.isDeleted).length;
        return `${count} ${count === 1 ? 'element' : 'elements'}`;
      } catch {
        return 'Scene JSON needs attention';
      }
    }

    function wireFenceControls(root, record, expanded) {
      const label = document.createElement('span');
      label.className = 'excalidraw-fence-controls__label';
      label.textContent = `Excalidraw · ${record.summary}`;
      const actions = document.createElement('span');
      actions.className = 'excalidraw-fence-controls__actions';
      const toggle = document.createElement('button');
      toggle.type = 'button';
      toggle.className = 'excalidraw-fence-controls__button';
      toggle.textContent = expanded ? 'Hide JSON' : 'Show JSON';
      toggle.setAttribute('aria-label', `${expanded ? 'Hide' : 'Show'} Excalidraw JSON`);
      toggle.setAttribute('aria-expanded', String(expanded));
      const edit = document.createElement('button');
      edit.type = 'button';
      edit.className = 'excalidraw-fence-controls__button excalidraw-fence-controls__button--edit';
      edit.textContent = 'Edit drawing';
      edit.disabled = !!editor.getOption('readOnly');
      edit.addEventListener('click', () => {
        if (!editor.getOption('readOnly')) openExisting(record.fence);
      });
      toggle.addEventListener('click', () => {
        record.collapsed = !record.collapsed;
        renderFence(record,record.fence);
        window.requestAnimationFrame(() => record.toggleButton?.focus());
      });
      root.addEventListener('mousedown',event => event.stopPropagation());
      root.addEventListener('click',event => event.stopPropagation());
      root.addEventListener('keydown',event => event.stopPropagation());
      actions.append(toggle,edit);
      root.append(label,actions);
      record.label = label;
      record.toggleButton = toggle;
      record.editButton = edit;
    }

    function renderFence(record, fence) {
      record.fence = fence;
      record.fold?.clear();
      record.fold = null;
      record.widget?.clear();
      record.widget = null;
      record.summary = sceneSummary(fence);

      if (record.collapsed) {
        const placeholder = document.createElement('span');
        placeholder.className = 'excalidraw-fence-controls excalidraw-fence-controls--folded';
        wireFenceControls(placeholder,record,false);
        record.fold = editor.markText(
          {line:fence.line+1,ch:0},
          {line:fence.closeLine,ch:0},
          {collapsed:true,clearOnEnter:false,inclusiveLeft:true,inclusiveRight:true,replacedWith:placeholder},
        );
      } else {
        const controls = document.createElement('div');
        controls.className = 'excalidraw-fence-controls excalidraw-fence-controls--expanded';
        wireFenceControls(controls,record,true);
        record.widget = editor.addLineWidget(fence.line,controls,{noHScroll:true,showIfHidden:true});
      }
    }

    function removeFence(record) {
      record.fold?.clear();
      record.widget?.clear();
      record.anchor.clear();
      fenceRecords.delete(record.id);
    }

    function refreshFences() {
      refreshScheduled = false;
      const fences = excalidrawFences();
      const remaining = new Map(fences.map(fence => [fence.line,fence]));
      for (const record of fenceRecords.values()) {
        const position = record.anchor.find();
        const fence = position && remaining.get(position.line);
        if (!fence) {
          removeFence(record);
          continue;
        }
        remaining.delete(position.line);
        record.fence = fence;
        const foldRange = record.fold?.find();
        const foldIsCurrent = foldRange?.from.line === fence.line+1 &&
          foldRange.from.ch === 0 && foldRange.to.line === fence.closeLine && foldRange.to.ch === 0;
        if ((record.collapsed && !foldIsCurrent) || (!record.collapsed && !record.widget)) {
          renderFence(record,fence);
        } else if (record.dirty) {
          record.summary = sceneSummary(fence);
          if (record.label) record.label.textContent = `Excalidraw · ${record.summary}`;
          if (record.editButton) record.editButton.disabled = !!editor.getOption('readOnly');
        }
        record.dirty = false;
      }
      for (const fence of remaining.values()) {
        const record = {id:nextFenceID++,anchor:editor.setBookmark({line:fence.line,ch:0},{insertLeft:true}),collapsed:true,dirty:false};
        fenceRecords.set(record.id,record);
        renderFence(record,fence);
      }
    }

    function scheduleFenceRefresh(change) {
      let shouldRefresh = false;
      for (let item = change; item; item = item.next) {
        if (item.text?.some(line => fenceMarkerLine.test(line))) shouldRefresh = true;
        const firstLine = Math.max(0,item.from.line-1);
        const lastLine = Math.min(editor.lineCount()-1,item.to.line+1);
        for (let line = firstLine; line <= lastLine; line++) {
          if (fenceMarkerLine.test(editor.getLine(line))) shouldRefresh = true;
        }
      }
      for (const record of fenceRecords.values()) {
        const fence = record.fence;
        for (let item = change; item; item = item.next) {
          if (item.from.line <= fence.closeLine) shouldRefresh = true;
          if (item.from.line <= fence.closeLine && item.to.line >= fence.line+1) record.dirty = true;
        }
      }
      if (!shouldRefresh) return;
      if (refreshScheduled) return;
      refreshScheduled = true;
      window.requestAnimationFrame(refreshFences);
    }

    function afterFence(fence) {
      const line = fence.closeLine + 1;
      return line < editor.lineCount()
        ? {line,ch:0}
        : {line:fence.closeLine,ch:editor.getLine(fence.closeLine).length};
    }

    function updateButton() {
      const readOnly = editor.getOption('readOnly');
      const fence = fenceAtCursor(editor);
      button.disabled = !!readOnly;
      button.textContent = fence?.language === 'excalidraw' ? 'Edit drawing' : 'Drawing';
      button.title = fence?.language === 'excalidraw'
        ? 'Edit the Excalidraw scene at the cursor'
        : 'Create an Excalidraw scene at the cursor';
    }

    function openDialog(mode, scene, target) {
      active = {mode,target};
      title.textContent = mode === 'edit' ? 'Edit drawing' : 'New drawing';
      hint.textContent = mode === 'edit'
        ? 'Changes stay in this draft until you apply.'
        : 'The scene is added to this page when you apply it.';
      error.hidden = true;
      error.textContent = '';
      loading.hidden = false;
      loading.textContent = 'Loading the local drawing canvas…';
      host.hidden = true;
      host.replaceChildren();
      applyButton.hidden = false;
      applyButton.disabled = true;
      applyAsButton.hidden = mode !== 'edit';
      applyAsButton.disabled = true;
      dialog.showModal();

      window.EXCALIDRAW_ASSET_PATH = '/editor/excalidraw-fonts/';
      Promise.all([ensureStyles(),import('/editor/excalidraw-editor.js')]).then(([,module]) => {
        if (!dialog.open || !active) return;
        canvas = module.mountExcalidraw(host,scene,undefined,currentTheme());
        host.hidden = false;
        loading.hidden = true;
        applyButton.disabled = false;
        applyAsButton.disabled = false;
      }).catch(cause => {
        if (!dialog.open) return;
        showError(`The drawing canvas could not start: ${cause.message}`);
      });
    }

    function openExisting(fence) {
      if (editor.getOption('readOnly')) return;
      if (fence.closeLine == null) {
        showError('Close this Excalidraw fence before editing it. The source has not changed.');
        return;
      }
      const from = {line:fence.line+1,ch:0};
      const to = {line:fence.closeLine,ch:0};
      const body = editor.getRange(from,to).replace(/\r?\n$/,'');
      try {
        openDialog('edit',parseScene(body),fence);
      } catch (cause) {
        showError(cause.message);
      }
    }

    function insertionPosition() {
      const fence = fenceAtCursor(editor);
      if (fence && fence.closeLine == null) {
        throw new Error('Close the surrounding code fence before adding a drawing.');
      }
      if (fence && fence.language !== 'excalidraw') {
        return afterFence(fence);
      }
      const cursor = editor.getCursor();
      return {line:cursor.line,ch:editor.getLine(cursor.line).length};
    }

    function openNew() {
      if (editor.getOption('readOnly')) return;
      try {
        openDialog('new',emptyScene(),insertionPosition());
      } catch (cause) {
        showError(cause.message);
      }
    }

    function openAtCursor() {
      if (editor.getOption('readOnly')) return;
      const fence = fenceAtCursor(editor);
      if (fence?.language === 'excalidraw') openExisting(fence);
      else openNew();
    }

    function apply(asNew = false) {
      if (!active || !canvas) return;
      try {
        const serialized = canvas.serialize();
        parseScene(serialized);
        if (active.mode === 'edit' && !asNew) {
          replaceSceneBody(editor,active.target,serialized);
        } else {
          const position = active.mode === 'edit' ?
            afterFence(active.target) : active.target;
          insertText(editor,position,serialized);
        }
        dialog.close('apply');
      } catch (cause) {
        error.textContent = cause.message;
        error.hidden = false;
      }
    }

    button.addEventListener('click',openAtCursor);
    document.querySelector('#excalidraw-cancel').addEventListener('click',() => dialog.close('cancel'));
    document.querySelector('#excalidraw-close').addEventListener('click',() => dialog.close('cancel'));
    applyButton.addEventListener('click',() => apply(false));
    applyAsButton.addEventListener('click',() => apply(true));
    dialog.addEventListener('close',() => {
      canvas?.destroy();
      canvas = null;
      active = null;
      host.replaceChildren();
      host.hidden = true;
      loading.hidden = false;
      error.hidden = true;
      applyButton.hidden = false;
      applyAsButton.hidden = true;
      applyButton.disabled = false;
      applyAsButton.disabled = false;
      editor.focus();
      editor.refresh();
      updateButton();
    });
    editor.on('cursorActivity',updateButton);
    editor.on('change',(_editor,change) => {updateButton();scheduleFenceRefresh(change);});
    editor.on('optionChange',(_editor,option) => {
      if (option === 'readOnly') {
        for (const record of fenceRecords.values()) {
          if (record.editButton) record.editButton.disabled = !!editor.getOption('readOnly');
        }
      }
    });
    refreshFences();
    updateButton();

    return [
      {id:'insert-excalidraw-drawing',label:'Insert Excalidraw drawing',
        description:'Create a sketch canvas and add its scene to this Markdown page.',
        aliases:'draw sketch whiteboard canvas diagram',
        when:() => !editor.getOption('readOnly'),run:openNew},
      {id:'edit-excalidraw-drawing',label:'Edit Excalidraw drawing',
        description:'Open the Excalidraw scene at the current cursor.',
        aliases:'draw sketch whiteboard canvas diagram',
        when:() => !editor.getOption('readOnly') && fenceAtCursor(editor)?.language === 'excalidraw',
        run:() => {const fence=fenceAtCursor(editor);if(fence)openExisting(fence);}},
    ];
  }

  window.MarkataExcalidraw = {attach};
})();
