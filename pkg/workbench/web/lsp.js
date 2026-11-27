/* Shared CodeMirror 5 adapter for the Go and FastAPI examples. */
(() => {
  'use strict';
  window.MarkataEditorLSP = {connect};
  async function connect(editor, status) {
    let socket, uri, version = 1, sequence = 0, ready = false, disposed = false;
    let timer, autoTimer, reconnectTimer, retryDelay = 500, generation = 0;
    const pending = new Map(), marks = [];
    const onceIgnored = new Set(), mutedRules = new Set();
    let proseDiagnostics = [], allDiagnostics = [], prosePopup, proseTicket = 0;
    const pos = p => ({line:p.line, ch:p.character});
    const rpcPos = p => ({line:p.line, character:p.ch});
    const label = message => { status.textContent = message; };
    function plainSnippet(text) {
      return text.replace(/\$\{\d+\|([^,|}]+)(?:,[^|}]*)?\|\}/g,'$1')
        .replace(/\$\{\d+:([^}]*)\}/g,'$1').replace(/\$\{\d+\}/g,'').replace(/\$\d+/g,'');
    }
    function clearMarks() { marks.splice(0).forEach(mark => mark.clear()); }
    const comparePos = (a,b) => a.line - b.line || a.ch - b.ch;
    const isWithin = (cursor,range) => comparePos(cursor,pos(range.start)) >= 0 && comparePos(cursor,pos(range.end)) <= 0;
    function proseKey(diagnostic) {
      const from = pos(diagnostic.range.start), to = pos(diagnostic.range.end);
      return JSON.stringify([diagnostic.code,from.line,from.ch,to.line,to.ch,editor.getRange(from,to)]);
    }
    function closeProse() { proseTicket++; prosePopup?.remove(); prosePopup = undefined; }
    function refreshMarks() {
      clearMarks();
      const visible = allDiagnostics.filter(d => d.source !== 'markata-prose' ||
        (!mutedRules.has(String(d.code)) && !onceIgnored.has(proseKey(d))));
      proseDiagnostics = visible.filter(d => d.source === 'markata-prose');
      window.MarkataWorkbenchDiagnostics?.update(visible.filter(d => d.source !== 'markata-prose'));
      for (const diagnostic of visible) {
        const proseHint = diagnostic.source === 'markata-prose';
        marks.push(editor.markText(pos(diagnostic.range.start),pos(diagnostic.range.end),
          {className:proseHint ? 'markata-lsp-prose-hint' : 'markata-lsp-diagnostic',
            ...(proseHint ? {} : {title:diagnostic.message})}));
      }
    }
    function addProseButton(container,label,handler,secondary = false) {
      const button = document.createElement('button');
      button.type = 'button';
      button.className = secondary ? 'markata-prose-button markata-prose-secondary' : 'markata-prose-button';
      button.textContent = label;
      button.addEventListener('click',event => { event.stopPropagation(); void handler(); });
      container.appendChild(button);
      return button;
    }
    function positionProsePopup(x,y) {
      if (!prosePopup) return;
      const padding = 12, bounds = prosePopup.getBoundingClientRect();
      prosePopup.style.left = Math.max(padding,Math.min(x,window.innerWidth-bounds.width-padding))+'px';
      prosePopup.style.top = Math.max(padding,
        y+bounds.height+padding > window.innerHeight ? y-bounds.height-padding : y+padding)+'px';
    }
    async function openProse(diagnostic,x,y) {
      if (!ready || !diagnostic || mutedRules.has(String(diagnostic.code))) return;
      closeProse();
      if (timer) sync();
      const ticket = ++proseTicket, snapshot = editor.getValue(), documentVersion = version;
      const from = pos(diagnostic.range.start), to = pos(diagnostic.range.end);
      if (editor.getRange(from,to) === '') return;
      let actions = [];
      try {
        actions = await request('textDocument/codeAction',{
          textDocument:{uri},range:{start:diagnostic.range.start,end:diagnostic.range.end},
          context:{diagnostics:[diagnostic]}
        });
      } catch (error) {
        if (ready) label('Writing suggestions unavailable · you can keep editing');
      }
      if (ticket !== proseTicket || !ready || snapshot !== editor.getValue() || documentVersion !== version) return;
      const popup = document.createElement('div');
      popup.className = 'markata-prose-popover';
      popup.setAttribute('role','dialog');
      popup.setAttribute('aria-label','Writing suggestion');
      const title = document.createElement('strong');
      title.textContent = 'Writing suggestion';
      const message = document.createElement('p');
      message.textContent = diagnostic.message;
      const buttons = document.createElement('div');
      buttons.className = 'markata-prose-actions';
      popup.append(title,message,buttons);
      const edits = (Array.isArray(actions) ? actions : []).flatMap(action =>
        action.edit?.documentChanges?.flatMap(change =>
          change.textDocument?.uri === uri && change.textDocument.version === documentVersion
            ? change.edits.map(edit => ({title:action.title,edit})) : []) || []);
      if (edits.length) {
        const selected = edits[0];
        addProseButton(buttons,selected.title,() => {
          if (version !== documentVersion || editor.getValue() !== snapshot) return closeProse();
          editor.operation(() => editor.replaceRange(selected.edit.newText,
            pos(selected.edit.range.start),pos(selected.edit.range.end),'+prose'));
          closeProse();
          editor.focus();
        });
      }
      addProseButton(buttons,'Ignore once',() => {
        onceIgnored.add(proseKey(diagnostic));
        closeProse();
        refreshMarks();
        editor.focus();
      },true);
      addProseButton(buttons,'Always ignore rule',async () => {
        const button = buttons.lastElementChild;
        button.disabled = true;
        try {
          const response = await fetch('/workspace/prose/ignore-rule',{
            method:'POST',credentials:'same-origin',
            headers:{'Content-Type':'application/json'},
            body:JSON.stringify({rule:String(diagnostic.code)})
          });
          if (!response.ok) throw new Error('Could not save writing preference');
          mutedRules.add(String(diagnostic.code));
          closeProse();
          refreshMarks();
          sync();
          editor.focus();
        } catch {
          button.disabled = false;
          label('Could not save writing preference · try again');
        }
      },true);
      document.body.appendChild(popup);
      prosePopup = popup;
      positionProsePopup(x,y);
    }
    function openProseAtCursor() {
      const cursor = editor.getCursor(), diagnostic = proseDiagnostics.find(d => isWithin(cursor,d.range));
      if (!diagnostic) return;
      const coords = editor.cursorCoords(cursor,'window');
      void openProse(diagnostic,coords.left,coords.bottom);
    }
    const wrapper = editor.getWrapperElement();
    wrapper.addEventListener('click',event => {
      if (!(event.target instanceof Element) || !event.target.closest('.markata-lsp-prose-hint')) return;
      const cursor = editor.coordsChar({left:event.clientX,top:event.clientY},'window');
      const diagnostic = proseDiagnostics.find(d => isWithin(cursor,d.range));
      if (diagnostic) void openProse(diagnostic,event.clientX,event.clientY);
    });
    editor.on('scroll',closeProse);
    document.addEventListener('pointerdown',event => {
      if (prosePopup && !prosePopup.contains(event.target) &&
        !(event.target instanceof Element && event.target.closest('.markata-lsp-prose-hint'))) closeProse();
    });
    document.addEventListener('keydown',event => {
      if (event.key === 'Escape' && prosePopup) { closeProse(); event.preventDefault(); }
    });
    function send(method, params, id) {
      if (socket?.readyState !== WebSocket.OPEN) throw new Error('Language service disconnected');
      socket.send(JSON.stringify({jsonrpc:'2.0', ...(id === undefined ? {} : {id}), method, params}));
    }
    function request(method, params) {
      const id = ++sequence;
      return new Promise((resolve,reject) => {
        const timeout = setTimeout(() => {pending.delete(id); reject(new Error('Language service timed out'));},10000);
        pending.set(id, {resolve,reject,timeout});
        try { send(method,params,id); } catch(error) { clearTimeout(timeout); pending.delete(id); reject(error); }
      });
    }
    function sync() {
      clearTimeout(timer); timer = undefined;
      if (ready) send('textDocument/didChange', {textDocument:{uri,version:++version},contentChanges:[{text:editor.getValue()}]});
    }
    function unavailable() {
      ready = false;
      generation++;
      clearTimeout(timer);
      clearTimeout(autoTimer);
      closeProse();
      allDiagnostics = [];
      proseDiagnostics = [];
      clearMarks();
      editor.closeHint?.();
      window.MarkataWorkbenchDiagnostics?.update([]);
      pending.forEach(({reject,timeout}) => {clearTimeout(timeout); reject(new Error('Language service disconnected'));});
      pending.clear();
      label('Language service unavailable · editing and preview remain available');
    }
    async function complete(syncDocument = true) {
      if (!ready || editor.somethingSelected()) return;
      const ticket = ++generation, text = editor.getValue(), cursor = editor.getCursor();
      try {
        if (syncDocument) sync();
        const result = await request('textDocument/completion', {textDocument:{uri},position:rpcPos(cursor)});
        const current = editor.getCursor();
        if (!ready || ticket !== generation || text !== editor.getValue() || current.line !== cursor.line || current.ch !== cursor.ch) return;
        const list = (Array.isArray(result) ? result : result?.items || []).map(item => {
          const edit = item.textEdit;
          const inserted = edit?.newText ?? item.insertText ?? item.label;
          return {text:item.insertTextFormat === 2 ? plainSnippet(inserted) : inserted,
            displayText:item.detail ? `${item.label} — ${item.detail}` : item.label,
            from:edit ? pos(edit.range.start) : cursor, to:edit ? pos(edit.range.end) : cursor};
        });
        editor.closeHint();
        if (list.length) editor.showHint({hint:() => ({list,from:cursor,to:cursor}),completeSingle:false});
      } catch(error) { if (ready) label('Language assistance error · Ctrl+Space to retry'); }
    }
    editor.addKeyMap({'Ctrl-Space':complete,'Cmd-Space':complete,'Alt-Enter':openProseAtCursor});
    editor.on('change', () => {
      generation++;
      closeProse();
      proseDiagnostics = [];
      clearMarks();
      editor.closeHint?.();
      clearTimeout(timer);
      if (ready) timer = setTimeout(() => {try {sync();} catch {unavailable();}},150);
    });
    editor.on('cursorActivity', () => {generation++;});
    editor.on('inputRead', (_cm,change) => {
      if (ready && change.origin === '+input') {
        clearTimeout(autoTimer);
        // CodeMirror emits cursorActivity after inputRead in the same operation.
        autoTimer = setTimeout(() => complete(false),180);
      }
    });
    function scheduleReconnect() {
      if (disposed || reconnectTimer) return;
      label('Language assistance reconnecting… · writing remains available');
      reconnectTimer = setTimeout(() => {
        reconnectTimer = undefined;
        startConnection();
      }, retryDelay);
      retryDelay = Math.min(retryDelay * 2, 8000);
    }
    function dispose() {
      if (disposed) return;
      disposed = true;
      clearTimeout(reconnectTimer);
      try { if (ready) send('textDocument/didClose',{textDocument:{uri}}); } catch {}
      unavailable();
      socket?.close();
    }
    window.addEventListener('pagehide',dispose,{once:true});
    async function startConnection() {
    if (disposed) return;
    try {
      const response = await fetch('/lsp/config',{cache:'no-store'});
      if (!response.ok) throw new Error('Language service configuration unavailable');
      const config = await response.json();
      if (disposed) return;
      uri = config.documentUri;
      const endpoint = new URL('/lsp',location.href);
      endpoint.protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
      socket = new WebSocket(endpoint);
      socket.addEventListener('message', event => {
        let message;
        try { message = JSON.parse(event.data); } catch {return;}
        if (message.id !== undefined) {
          const item = pending.get(message.id);
          if (!item) return;
          clearTimeout(item.timeout); pending.delete(message.id);
          if (message.error) item.reject(new Error(message.error.message)); else item.resolve(message.result);
        } else if (message.method === 'textDocument/publishDiagnostics' && message.params.uri === uri) {
          // Full-sync notifications arrive in wire order. Ignore diagnostics while a newer edit is pending.
          if (timer || (message.params.version !== undefined && message.params.version !== version)) return;
          allDiagnostics = message.params.diagnostics || [];
          refreshMarks();
        }
      });
      socket.addEventListener('close', () => {unavailable();scheduleReconnect();});
      socket.addEventListener('error', () => {unavailable();scheduleReconnect();});
      await new Promise((resolve,reject) => {
        const timeout = setTimeout(() => {socket.close();reject(new Error('Connection timed out'));},10000);
        socket.addEventListener('open',() => {clearTimeout(timeout);resolve();},{once:true});
        socket.addEventListener('error',() => {clearTimeout(timeout);reject(new Error('Connection failed'));},{once:true});
      });
      await request('initialize',{processId:null,rootUri:config.rootUri,capabilities:{}});
      send('initialized',{});
      send('textDocument/didOpen',{textDocument:{uri,languageId:'markdown',version,text:editor.getValue()}});
      ready = true;
      retryDelay = 500;
      label('LSP connected · Ctrl+Space for links');
    } catch { unavailable();socket?.close();scheduleReconnect(); }
    }
    startConnection();
  }
})();
