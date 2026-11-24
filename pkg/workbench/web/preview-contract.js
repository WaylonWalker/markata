/* The hosted editor and local WASM renderer share one fail-closed preview contract.
 * Keep this dependency-free: the production editor does not need Node/npm. */
(() => {
  'use strict';
  const SCHEMA_VERSION = 1;
  const knownNumber = value => Number.isFinite(value) && value >= 0;
  const lineList = (values, field) => {
    if (!Array.isArray(values) || !values.every(value => Number.isSafeInteger(value) && value >= 0))
      throw new Error('Invalid preview ' + field);
    return values;
  };
  function fromGo(response, sourceOffset = 0, requireDocument = false) {
    if (!response || response.schema_version !== SCHEMA_VERSION)
      throw new Error('Incompatible Go preview contract; reload the editor');
    if (!Number.isSafeInteger(sourceOffset) || sourceOffset < 0)
      throw new Error('Invalid preview source offset');
    if (typeof response.renderer_fingerprint !== 'string' ||
        typeof response.html !== 'string' ||
        typeof response.feature_head_html !== 'string' ||
        typeof response.body_end_html !== 'string')
      throw new Error('Incomplete Go preview payload');
    if (requireDocument && typeof response.document !== 'string')
      throw new Error('Initial Go preview document is missing');
    if (response.document != null && typeof response.document !== 'string')
      throw new Error('Invalid preview document');
    const sourceLines = lineList(response.source_lines, 'source lines');
    const sourceAnchors = lineList(response.source_anchors, 'source anchors');
    const metrics = ['render_ms', 'queue_ms', 'input_bytes', 'html_bytes'];
    if (metrics.some(field => !knownNumber(response[field])))
      throw new Error('Invalid Go preview timings');
    return {
      ...response,
      source_lines: sourceLines.map(line => line + sourceOffset),
      source_anchors: sourceAnchors.map(line => line + sourceOffset),
      source_offset: sourceOffset,
    };
  }
  /* This only accepts base-style, parity-certified fragments. Never paste
   * arbitrary standalone tags into the authenticated workbench preview. */
  function fromWasm(result, timing = {}) {
    if (!result || result.schema_version !== SCHEMA_VERSION || result.error)
      throw new Error(result?.error || 'Incompatible WASM preview result');
    if (typeof result.html !== 'string')
      throw new Error('WASM preview has no HTML');
    const sourceLines = lineList(result.source_lines, 'source lines');
    // The opaque preview shell already owns these three base stylesheets.
    // Feature CSS, generated scripts, and other dependencies remain Go-only
    // until the asset bridge explicitly supports them.
    const baseStyles = new Set(['css/variables.css', 'css/main.css', 'css/components.css']);
    if (!Array.isArray(result.assets) ||
        result.assets.some(asset => !baseStyles.has(asset.id) || asset.kind !== 'style' || asset.content) ||
        Object.values(result.features || {}).some(Boolean))
      throw new Error('WASM preview requires Go-managed assets or features');
    if (!knownNumber(timing.render_ms || 0) || !knownNumber(timing.input_bytes || 0))
      throw new Error('Invalid WASM preview timings');
    return {
      schema_version: SCHEMA_VERSION, document: null, html: result.html,
      source_lines: [...sourceLines], source_anchors: [], source_offset: 0,
      feature_head_html: '', body_end_html: '',
      render_ms: timing.render_ms || 0, queue_ms: 0,
      input_bytes: timing.input_bytes || 0,
      html_bytes: new TextEncoder().encode(result.html).length,
      document_bytes: 0, cached: false,
    };
  }
  /* No implicit guessing from Markdown syntax is sufficient to certify parity.
   * An explicit matched site renderer fingerprint must authorize the worker. */
  function chooseRenderer({workerReady = false, previewFrameReady = false,
                           siteFingerprint = '', workerFingerprint = '',
                           parityCertified = false, requiresSiteRender = false} = {}) {
    if (!workerReady) return {renderer: 'go', reason: 'worker-not-ready'};
    if (!previewFrameReady) return {renderer: 'go', reason: 'preview-shell-not-ready'};
    if (!siteFingerprint || !workerFingerprint || siteFingerprint !== workerFingerprint)
      return {renderer: 'go', reason: 'renderer-fingerprint-mismatch'};
    if (!parityCertified) return {renderer: 'go', reason: 'parity-not-certified'};
    if (requiresSiteRender) return {renderer: 'go', reason: 'site-dependent-content'};
    return {renderer: 'wasm', reason: 'certified-local-preview'};
  }
  function wasmEligibility({workerSupported = false, wasmSupported = false,
                            siteFingerprint = '', workerFingerprint = '',
                            deviceMemory, hardwareConcurrency, saveData = false,
                            effectiveType = ''} = {}) {
    if (!workerSupported || !wasmSupported) return {allowed: false, reason: 'browser-unsupported'};
    if (!siteFingerprint || siteFingerprint !== workerFingerprint)
      return {allowed: false, reason: 'renderer-fingerprint-mismatch'};
    if (!Number.isFinite(deviceMemory) || deviceMemory <= 4)
      return {allowed: false, reason: 'device-memory-limited'};
    if (!Number.isFinite(hardwareConcurrency) || hardwareConcurrency <= 4)
      return {allowed: false, reason: 'processor-capacity-limited'};
    if (saveData || effectiveType === 'slow-2g' || effectiveType === '2g')
      return {allowed: false, reason: 'reduced-data'};
    return {allowed: true, reason: 'eligible'};
  }
  window.MarkataPreviewContract = Object.freeze({
    schemaVersion: SCHEMA_VERSION, fromGo, fromWasm, chooseRenderer, wasmEligibility,
  });
})();
