/* Shared CodeMirror 5 adapter for the Go and FastAPI examples. */
(() => {
  'use strict';
  window.MarkataEditorLSP = {connect};
  async function connect(editor, status) {
    let socket, uri, version = 1, sequence = 0, ready = false, disposed = false;
    let timer, autoTimer, generation = 0;
    const pending = new Map(), marks = [];
    const pos = p => ({line:p.line, ch:p.character});
    const rpcPos = p => ({line:p.line, character:p.ch});
    const label = message => { status.textContent = message; };
    function plainSnippet(text) {
      return text.replace(/\$\{\d+\|([^,|}]+)(?:,[^|}]*)?\|\}/g,'$1')
        .replace(/\$\{\d+:([^}]*)\}/g,'$1').replace(/\$\{\d+\}/g,'').replace(/\$\d+/g,'');
    }
    function clearMarks() { marks.splice(0).forEach(mark => mark.clear()); }
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
    editor.addKeyMap({'Ctrl-Space':complete,'Cmd-Space':complete});
    editor.on('change', () => {
      generation++;
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
    function dispose() {
      if (disposed) return;
      disposed = true;
      try { if (ready) send('textDocument/didClose',{textDocument:{uri}}); } catch {}
      unavailable();
      socket?.close();
    }
    window.addEventListener('pagehide',dispose,{once:true});
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
          clearMarks();
          window.MarkataWorkbenchDiagnostics?.update(message.params.diagnostics);
          for (const diagnostic of message.params.diagnostics) {
            marks.push(editor.markText(pos(diagnostic.range.start),pos(diagnostic.range.end),
              {className:'markata-lsp-diagnostic',title:diagnostic.message}));
          }
        }
      });
      socket.addEventListener('close',unavailable);
      socket.addEventListener('error',unavailable);
      await new Promise((resolve,reject) => {
        const timeout = setTimeout(() => {socket.close();reject(new Error('Connection timed out'));},10000);
        socket.addEventListener('open',() => {clearTimeout(timeout);resolve();},{once:true});
        socket.addEventListener('error',() => {clearTimeout(timeout);reject(new Error('Connection failed'));},{once:true});
      });
      await request('initialize',{processId:null,rootUri:config.rootUri,capabilities:{}});
      send('initialized',{});
      send('textDocument/didOpen',{textDocument:{uri,languageId:'markdown',version,text:editor.getValue()}});
      ready = true;
      label('LSP connected · Ctrl+Space for links');
    } catch { unavailable();socket?.close(); }
  }
})();
