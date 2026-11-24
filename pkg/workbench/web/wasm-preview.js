/* Local WASM preview worker transport. No credentials, site mutation or DOM. */
(() => {
  'use strict';
  function create(version) {
    if (!version || !window.Worker) throw new Error('WASM preview unavailable');
    const worker = new Worker('/editor/wasm/worker.js?v=' + encodeURIComponent(version),
      {name: 'markata-preview-wasm'});
    const pending = new Map();
    let counter = 0, running = true, ready = false, memoryBytes = 0;
    let completeReady, rejectReady;
    const initialized = new Promise((resolve, reject) => {
      completeReady = resolve; rejectReady = reject;
    });
    const fail = error => {
      if (!running) return;
      running = false;
      ready = false;
      const failure = error instanceof Error ? error : new Error(String(error));
      rejectReady(failure);
      for (const entry of pending.values()) {
        clearTimeout(entry.timeout);
        entry.reject(failure);
      }
      pending.clear();
      worker.terminate();
    };
    const initializationTimeout = setTimeout(() => fail(new Error('WASM worker initialization timed out')), 45000);
    worker.onerror = event => {
      event.preventDefault?.();
      fail(new Error('WASM worker failed to start'));
    };
    worker.onmessage = ({data}) => {
      if (data?.type === 'fatal') { fail(data.error); return; }
      if (data?.type === 'ready') {
        if (data.schema_version !== 1 || data.version !== version) {
          fail(new Error('WASM worker version mismatch'));
          return;
        }
        ready = true;
        if (Number.isSafeInteger(data.memory_bytes) && data.memory_bytes > 0) memoryBytes = data.memory_bytes;
        clearTimeout(initializationTimeout);
        completeReady({memory_bytes: memoryBytes});
        return;
      }
      if (data?.type !== 'result' || !Number.isSafeInteger(data.revision)) return;
      const entry = pending.get(data.revision);
      if (!entry) return; // Response from an obsolete request.
      pending.delete(data.revision);
      clearTimeout(entry.timeout);
      if (Number.isSafeInteger(data.memory_bytes) && data.memory_bytes > 0) memoryBytes = data.memory_bytes;
      if (data.result?.error) entry.reject(new Error(data.result.error));
      else entry.resolve({result: data.result, render_ms: data.render_ms, memory_bytes: memoryBytes});
    };
    initialized.catch(() => {}); // Worker may fail before anyone requests it.
    return {
      ready: initialized,
      render(source) {
        if (!ready || !running) return Promise.reject(new Error('WASM worker is not ready'));
        if (typeof source !== 'string') return Promise.reject(new Error('Invalid preview source'));
        const revision = ++counter;
        return new Promise((resolve, reject) => {
          const timeout = setTimeout(() => {
            pending.delete(revision);
            reject(new Error('WASM preview timed out'));
          }, 15000);
          pending.set(revision, {resolve, reject, timeout});
          worker.postMessage({type: 'render', revision, source});
        });
      },
      close() {
        clearTimeout(initializationTimeout);
        fail(new Error('WASM preview worker closed'));
      },
    };
  }
  // A deliberately small fast-path grammar. Site-specific syntax, assets and
  // formatting stay on Go. This is not a general Markdown feature detector.
  function simpleTextOnly(source) {
    if (typeof source !== 'string') return false;
    return !/(?:^---[ \t]*$|^[ \t]*(?:\x60{3}|~{3}|:{3}|!{3}|@@)|^ {4}\S|\[\[|!\[|\[[^\]]+\]\(|:lucide-|<|\{\{|\{%|^[ \t]*\||^[ \t]*[-*+] \[[xX ]\])/m.test(source);
  }
  function compareParity(goPayload, wasmPayload, sourceOffset = 0) {
    if (!goPayload || !wasmPayload || !Array.isArray(goPayload.source_lines) ||
        !Array.isArray(wasmPayload.source_lines) ||
        (goPayload.source_anchors || []).length !== 0) return null;
    if (!goPayload.html.endsWith(wasmPayload.html)) return null;
    if (goPayload.source_lines.length !== wasmPayload.source_lines.length) return null;
    for (let i = 0; i < wasmPayload.source_lines.length; i++) {
      if (goPayload.source_lines[i] !== wasmPayload.source_lines[i] + sourceOffset) return null;
    }
    return {prefix: goPayload.html.slice(0, goPayload.html.length - wasmPayload.html.length)};
  }
  window.MarkataWasmPreview = Object.freeze({create, simpleTextOnly, compareParity});
})();
