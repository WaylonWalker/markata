/* Local WASM preview worker transport. No credentials, site mutation or DOM. */
(() => {
  'use strict';
  function create(version) {
    if (!version || !window.Worker) throw new Error('WASM preview unavailable');
    const worker = new Worker('/editor/wasm/worker.js?v=' + encodeURIComponent(version),
      {name: 'markata-preview-wasm'});
    const pending = new Map();
    let counter = 0, running = true, ready = false;
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
        clearTimeout(initializationTimeout);
        completeReady();
        return;
      }
      if (data?.type !== 'result' || !Number.isSafeInteger(data.revision)) return;
      const entry = pending.get(data.revision);
      if (!entry) return; // Response from an obsolete request.
      pending.delete(data.revision);
      clearTimeout(entry.timeout);
      if (data.result?.error) entry.reject(new Error(data.result.error));
      else entry.resolve({result: data.result, render_ms: data.render_ms});
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
  window.MarkataWasmPreview = Object.freeze({create});
})();