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
  const applyStatus = document.querySelector('#excalidraw-apply-status') || document.querySelector('#saved') || document.querySelector('#status');
  const artifactOptions = 'padding=20;fonts=inline-local-woff2;svg=1';
  const editableElementTypes = new Set(['rectangle','diamond','ellipse','line','arrow','freedraw','text','image','frame','magicframe','iframe','embeddable']);
  let cssReady;
  let rendererVersionPromise;

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
    if (scene.version !== 2) {
      throw new Error(`This Excalidraw scene uses unsupported version ${scene.version}. Its source has been kept unchanged.`);
    }
    for (let index = 0; index < scene.elements.length; index++) {
      const element = scene.elements[index];
      if (!element || typeof element !== 'object' || !editableElementTypes.has(element.type)) {
        throw new Error(`Excalidraw element ${index + 1} has an unsupported type. Its source has been kept unchanged.`);
      }
    }
    if (scene.appState != null && (typeof scene.appState !== 'object' || Array.isArray(scene.appState))) {
      throw new Error('This Excalidraw scene has invalid appState data. Its source has been kept unchanged.');
    }
    if (scene.files != null && (typeof scene.files !== 'object' || Array.isArray(scene.files))) {
      throw new Error('This Excalidraw scene has invalid image data. Its source has been kept unchanged.');
    }
    return scene;
  }

  function canonicalJSON(value) {
    if (Array.isArray(value)) return `[${value.map(canonicalJSON).join(',')}]`;
    if (value && typeof value === 'object') {
      return `{${Object.keys(value).sort().map(key => `${JSON.stringify(key)}:${canonicalJSON(value[key])}`).join(',')}}`;
    }
    return JSON.stringify(value);
  }

  async function sha256Hex(value) {
    const digest = await window.crypto.subtle.digest('SHA-256',new TextEncoder().encode(value));
    return [...new Uint8Array(digest)].map(byte => byte.toString(16).padStart(2,'0')).join('');
  }

  async function rendererVersion() {
    if (typeof window.MARKATA_EXCALIDRAW_RENDERER_VERSION === 'string' && window.MARKATA_EXCALIDRAW_RENDERER_VERSION) {
      return window.MARKATA_EXCALIDRAW_RENDERER_VERSION;
    }
    if (!rendererVersionPromise) rendererVersionPromise = (async () => {
      const manifestResponse = await fetch('/editor/excalidraw-fonts/manifest.json',{cache:'no-store'});
      if (!manifestResponse.ok) throw new Error('the bundled font manifest is unavailable');
      const manifest = await manifestResponse.json();
      if (typeof manifest.renderer_version !== 'string' || !Array.isArray(manifest.fonts) || manifest.fonts.length === 0) {
        throw new Error('the bundled font manifest is invalid');
      }
      const fonts = [...manifest.fonts].sort();
      const parts = [];
      let size = 0;
      for (const font of fonts) {
        if (typeof font !== 'string' || !/^[A-Za-z0-9_.-]+(?:\/[A-Za-z0-9_.-]+)*\.woff2$/.test(font) || font.split('/').some(part => part === '.' || part === '..')) {
          throw new Error('the bundled font manifest contains an invalid path');
        }
        const response = await fetch('/editor/excalidraw-fonts/' + font.split('/').map(encodeURIComponent).join('/'),{cache:'no-store'});
        if (!response.ok) throw new Error(`a bundled font could not be checked (${response.status})`);
        const nameBytes = new TextEncoder().encode(font);
        const data = new Uint8Array(await response.arrayBuffer());
        parts.push(nameBytes,new Uint8Array([0]),data);
        size += nameBytes.length + 1 + data.length;
      }
      // The Go validator includes the exported JS/CSS entrypoints as well as
      // fonts. The entrypoint imports content-hashed chunks, so exporter
      // rebuilds invalidate the cached SVG even when the Excalidraw package
      // version itself has not changed.
      for (const name of ['excalidraw-editor.js','excalidraw-editor.css']) {
        const response = await fetch('/editor/' + name,{cache:'no-store'});
        if (!response.ok) throw new Error(`a bundled renderer asset could not be checked (${response.status})`);
        const nameBytes = new TextEncoder().encode(name);
        const data = new Uint8Array(await response.arrayBuffer());
        parts.push(nameBytes,new Uint8Array([0]),data);
        size += nameBytes.length + 1 + data.length;
      }
      const input = new Uint8Array(size);
      let offset = 0;
      for (const part of parts) { input.set(part,offset); offset += part.length; }
      const digest = await window.crypto.subtle.digest('SHA-256',input);
      const hash = [...new Uint8Array(digest)].map(byte => byte.toString(16).padStart(2,'0')).join('');
      return `${manifest.renderer_version}:${hash}`;
    })();
    return rendererVersionPromise;
  }

  async function localFontDataURL(path) {
    if (path.split('/').some(part => !part || part === '.' || part === '..')) {
      throw new Error('the SVG contains an invalid font asset path');
    }
    const response = await fetch('/editor/excalidraw-fonts/' + path.split('/').map(encodeURIComponent).join('/'),{cache:'no-store'});
    if (!response.ok) throw new Error(`a bundled font could not be loaded (${response.status})`);
    const bytes = new Uint8Array(await response.arrayBuffer());
    let binary = '';
    for (let offset = 0; offset < bytes.length; offset += 0x8000) {
      binary += String.fromCharCode(...bytes.subarray(offset,offset+0x8000));
    }
    return 'data:font/woff2;base64,' + window.btoa(binary);
  }

  async function inlineLocalFonts(svg) {
    const fontURL = /https:\/\/esm\.sh\/@excalidraw\/excalidraw@[^/]+\/dist\/prod\/fonts\/([^)'"\s]+)/g;
    const matches = [...svg.matchAll(fontURL)];
    for (const match of matches) {
      const dataURL = await localFontDataURL(match[1]);
      svg = svg.split(match[0]).join(dataURL);
    }
    if (/url\(\s*["']?https?:\/\//i.test(svg)) {
      throw new Error('the SVG contains a remote resource that cannot be embedded offline');
    }
    return svg;
  }

  async function gzipBase64(value) {
    if (typeof CompressionStream !== 'function') throw new Error('this browser cannot compress SVG artifacts');
    const stream = new Blob([value]).stream().pipeThrough(new CompressionStream('gzip'));
    const bytes = new Uint8Array(await new Response(stream).arrayBuffer());
    let binary = '';
    for (let offset = 0; offset < bytes.length; offset += 0x8000) {
      binary += String.fromCharCode(...bytes.subarray(offset,offset+0x8000));
    }
    return window.btoa(binary);
  }

  function staticPublishBlockers(scene) {
    const blockers = new Set();
    for (const element of scene.elements || []) {
      if (element.isDeleted) continue;
      if (typeof element.link === 'string' && element.link.trim()) blockers.add('linked elements');
      if (element.type === 'image') blockers.add('embedded images');
      if (element.type === 'iframe' || element.type === 'embeddable') blockers.add('web embeds');
    }
    return [...blockers];
  }

  async function createSVGArtifact(module, serialized) {
    const currentRendererVersion = await rendererVersion();
    const theme = window.MARKATA_EXCALIDRAW_EXPORT_THEME;
    const exportTheme = ['dark','light'].includes(theme) ? theme : currentTheme();
    if (!currentRendererVersion || !['dark','light'].includes(exportTheme)) {
      throw new Error('the editor could not determine its renderer fingerprint or site theme');
    }
    const scene = parseScene(serialized);
    const blockers = staticPublishBlockers(scene);
    if (blockers.length) {
      throw new Error('static publishing does not yet support ' + blockers.join(', '));
    }
    delete scene._markata;
    const canonical = canonicalJSON(scene);
    const sourceSHA256 = await sha256Hex(canonical);
    const svg = await inlineLocalFonts(await module.renderExcalidrawSvg(scene,exportTheme === 'dark'));
    if (new TextEncoder().encode(svg).byteLength > 4 * 1024 * 1024) {
      throw new Error('the SVG exceeds Markata’s 4 MiB static output limit');
    }
    if (/<(?:image|a|foreignObject|script|iframe|object|embed|feImage)(?:\s|\/?>)/i.test(svg)) {
      throw new Error('the SVG uses an image, link, or active element that static publishing rejects');
    }
    const svgSHA256 = await sha256Hex(svg);
    return {
      format:'markata-excalidraw-svg-v1',
      source_sha256:sourceSHA256,
      renderer_version:currentRendererVersion,
      theme:exportTheme,
      options:artifactOptions,
      encoding:'gzip+base64',
      svg_sha256:svgSHA256,
      svg:await gzipBase64(svg),
    };
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
        if (!editor.getOption('readOnly')) openExisting(record.fence,record);
      });
      toggle.addEventListener('click', () => {
        record.collapsed = !record.collapsed;
        renderFence(record,record.fence);
        window.requestAnimationFrame(() => record.toggleButton?.focus());
      });
      // CodeMirror fence widgets can lose button focus between Space down and
      // Space up. Activate on keydown and cancel the native key action so no
      // space is inserted into the Markdown or a second click is dispatched.
      root.setAttribute('data-cm-ignore-events','true');
      toggle.addEventListener('keydown', event => {
        if (event.key !== ' ' && event.code !== 'Space') return;
        event.preventDefault();
        event.stopPropagation();
        if (event.repeat) return;
        toggle.click();
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

    function openDialog(mode, scene, target, launcherRecord = null) {
      // CodeMirror may replace the actual launcher while the modal is open.
      // Keep its fence record so the current button can be found on close.
      const opener = document.activeElement === document.body ? button : document.activeElement;
      active = {mode,target,opener,launcherRecord};
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
      dialog.returnValue = '';
      dialog.showModal();

      window.EXCALIDRAW_ASSET_PATH = '/editor/excalidraw-fonts/';
      Promise.all([ensureStyles(),import('/editor/excalidraw-editor.js')]).then(([,module]) => {
        if (!dialog.open || !active) return;
        active.module = module;
        canvas = module.mountExcalidraw(host,scene,source => {
          try {
            const unsupported = staticPublishBlockers(parseScene(source));
            hint.textContent = unsupported.length
              ? 'This drawing can be edited, but static publishing does not yet support ' + unsupported.join(', ') + '. Apply will save the drawing without a static image.'
              : mode === 'edit' ? 'Changes stay in this draft until you apply.' : 'The scene is added to this page when you apply it.';
          } catch {
            // Invalid intermediary canvas state is handled at Apply.
          }
        },currentTheme());
        host.hidden = false;
        loading.hidden = true;
        applyButton.disabled = false;
        applyAsButton.disabled = false;
      }).catch(cause => {
        if (!dialog.open) return;
        showError(`The drawing canvas could not start: ${cause.message}`);
      });
    }

    function openExisting(fence, launcherRecord = null) {
      if (editor.getOption('readOnly')) return;
      if (fence.closeLine == null) {
        showError('Close this Excalidraw fence before editing it. The source has not changed.');
        return;
      }
      const from = {line:fence.line+1,ch:0};
      const to = {line:fence.closeLine,ch:0};
      const body = editor.getRange(from,to).replace(/\r?\n$/,'');
      try {
        openDialog('edit',parseScene(body),fence,launcherRecord);
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

    async function apply(asNew = false) {
      if (!active || !canvas) return;
      const session = active;
      try {
        const serialized = canvas.serialize();
        parseScene(serialized);
        applyButton.disabled = true;
        applyAsButton.disabled = true;
        let sourceToApply = serialized;
        let warning = '';
        try {
          const artifact = await createSVGArtifact(session.module,serialized);
          const sceneWithArtifact = parseScene(serialized);
          sceneWithArtifact._markata = {...(sceneWithArtifact._markata || {}),excalidraw_svg:artifact};
          sourceToApply = JSON.stringify(sceneWithArtifact,null,2);
        } catch (cause) {
          warning = `No static SVG artifact was saved (${cause.message}).`;
        }
        if (!dialog.open || active !== session) return;
        if (active.mode === 'edit' && !asNew) {
          try {
            replaceSceneBody(editor,session.target,sourceToApply);
          } catch (cause) {
            if (sourceToApply === serialized || !cause.message.includes('exceed 200,000 bytes')) throw cause;
            replaceSceneBody(editor,session.target,serialized);
            warning = 'The editable scene was applied without an SVG artifact because the document limit would be exceeded.';
          }
        } else {
          const position = session.mode === 'edit' ? afterFence(session.target) : session.target;
          try {
            insertText(editor,position,sourceToApply);
          } catch (cause) {
            if (sourceToApply === serialized || !cause.message.includes('exceed 200,000 bytes')) throw cause;
            insertText(editor,position,serialized);
            warning = 'The editable scene was inserted without an SVG artifact because the document limit would be exceeded.';
          }
        }
        if (applyStatus) {
          applyStatus.textContent = warning ? `Scene applied. ${warning} Publishing will fail in the default mode until the unsupported content is removed or a compatible SVG artifact is available.` : '';
          applyStatus.dataset.warning = warning ? 'true' : 'false';
        }
        dialog.close('apply');
      } catch (cause) {
        error.textContent = cause.message;
        error.hidden = false;
        applyButton.disabled = false;
        applyAsButton.disabled = false;
      }
    }

    button.addEventListener('click',openAtCursor);
    document.querySelector('#excalidraw-cancel').addEventListener('click',() => dialog.close('cancel'));
    document.querySelector('#excalidraw-close').addEventListener('click',() => dialog.close('cancel'));
    applyButton.addEventListener('click',() => apply(false));
    applyAsButton.addEventListener('click',() => apply(true));
    dialog.addEventListener('close',() => {
      const opener = active?.opener;
      const launcherRecord = active?.launcherRecord;
      const applied = dialog.returnValue === 'apply';
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
      // On Cancel, return to the *current* fence control, not a detached
      // element retained before CodeMirror redrew the collapsed marker.
      // On Apply, return to the Markdown editor with the new drawing.
      if (applied) {
        editor.refresh();
        editor.focus();
      } else {
        window.requestAnimationFrame(() => {
          let launcher = opener;
          if (launcherRecord) {
            const currentRecord = fenceRecords.get(launcherRecord.id) ||
              [...fenceRecords.values()].find(record => record.fence.line === launcherRecord.fence.line);
            launcher = currentRecord?.editButton || opener;
          }
          if (launcher?.isConnected && !launcher.disabled) launcher.focus({preventScroll:true});
          else if (!button.disabled) button.focus({preventScroll:true});
        });
      }
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
