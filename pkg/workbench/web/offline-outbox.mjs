// Private-draft local persistence foundation. Not connected to the hosted UI.
// An authenticated caller must supply the verified site and stable subject IDs.
// Merely partitioning IndexedDB is NOT encryption or a shared-device security boundary.

const DB_NAME = 'markata-editor-local-v1';
const SITE_RE = /^[a-z0-9](?:[a-z0-9-]{0,62}[a-z0-9])?$/;
const UUID_RE = /^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$/;
const MAX_SOURCE_BYTES = 2 << 20;

function scope(siteId, subjectId) {
  if (!SITE_RE.test(siteId) || typeof subjectId !== 'string' ||
      !subjectId || subjectId.length > 256 ||
      subjectId.trim() !== subjectId || /[\x00\r\n]/.test(subjectId)) {
    throw new TypeError('A verified site ID and stable subject ID are required');
  }
  return {siteId, subjectId};
}

function documentID(value) {
  if (typeof value !== 'string' || !value || value.length > 256 ||
      value.trim() !== value || /[\x00\r\n]/.test(value)) {
    throw new TypeError('Invalid document ID');
  }
  return value;
}

function requireRevision(value) {
  if (!Number.isSafeInteger(value) || value < 0) {
    throw new TypeError('Revision must be a nonnegative safe integer');
  }
  return value;
}

function sourceText(value) {
  if (typeof value !== 'string' ||
      new TextEncoder().encode(value).byteLength > MAX_SOURCE_BYTES) {
    throw new TypeError('Draft source must be a UTF-8 string of at most 2 MiB');
  }
  return value;
}

function requestValue(request) {
  return new Promise((resolve, reject) => {
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error || new Error('IndexedDB request failed'));
  });
}

// A write resolves only after transaction.oncomplete. No UI may show
// "Saved on this device" before the returned promise resolves.
export class BrowserIndexedDBDriver {
  constructor(db) { this.db = db; }

  async write(stores, callback) {
    const transaction = this.db.transaction(stores, 'readwrite');
    const completed = new Promise((resolve, reject) => {
      transaction.oncomplete = resolve;
      transaction.onerror = () => reject(transaction.error || new Error('IndexedDB transaction error'));
      transaction.onabort = () => reject(transaction.error || new Error('IndexedDB transaction aborted'));
    });
    // Avoid an unhandled rejection if an individual request fails first.
    void completed.catch(() => {});
    const tx = {
      get: (name, key) => requestValue(transaction.objectStore(name).get(key)),
      put: (name, value) => requestValue(transaction.objectStore(name).put(value)),
      forScope: (name, siteId, subjectId) =>
        requestValue(transaction.objectStore(name).index('byScope').getAll([siteId, subjectId])),
    };
    try {
      const value = await callback(tx);
      await completed;
      return value;
    } catch (error) {
      try { transaction.abort(); } catch { /* Already committed or aborted. */ }
      await completed.catch(() => {});
      throw error;
    }
  }
}

export function openBrowserDraftStore(indexedDB = globalThis.indexedDB) {
  if (!indexedDB || typeof indexedDB.open !== 'function') {
    return Promise.reject(new Error('IndexedDB is unavailable; offline saving is disabled'));
  }
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, 1);
    request.onupgradeneeded = () => {
      const db = request.result;
      for (const [name, key] of [
        ['drafts', ['siteId', 'subjectId', 'documentId']],
        ['mutations', ['siteId', 'subjectId', 'mutationId']],
      ]) {
        if (db.objectStoreNames.contains(name)) continue;
        const store = db.createObjectStore(name, {keyPath: key});
        store.createIndex('byScope', ['siteId', 'subjectId'], {unique: false});
      }
    };
    request.onerror = () => reject(request.error || new Error('Cannot open IndexedDB'));
    request.onblocked = () => reject(new Error('Close other editor tabs to upgrade local storage'));
    request.onsuccess = () => {
      const db = request.result;
      db.onversionchange = () => db.close();
      resolve(new LocalDraftOutbox(new BrowserIndexedDBDriver(db)));
    };
  });
}

export class LocalDraftOutbox {
  constructor(driver) {
    if (!driver || typeof driver.write !== 'function') {
      throw new TypeError('A transactional draft driver is required');
    }
    this.driver = driver;
  }

  // Each keystroke/batch can be locally committed without creating a new
  // server mutation. The server outbox snapshots a revision only at enqueue.
  async saveDraft({siteId, subjectId, documentId, source, deleted = false}) {
    scope(siteId, subjectId);
    documentID(documentId);
    sourceText(source);
    if (typeof deleted !== 'boolean') throw new TypeError('Invalid deletion intent');
    return this.driver.write(['drafts'], async tx => {
      const key = [siteId, subjectId, documentId];
      const prior = await tx.get('drafts', key);
      const draft = {
        siteId, subjectId, documentId, source, deleted,
        localRevision: (prior?.localRevision || 0) + 1,
        serverRevision: prior?.serverRevision ?? 0,
        syncedLocalRevision: prior?.syncedLocalRevision ?? 0,
      };
      await tx.put('drafts', draft);
      return draft;
    });
  }

  async getDraft(siteId, subjectId, documentId) {
    scope(siteId, subjectId);
    documentID(documentId);
    return this.driver.write(['drafts'], tx => tx.get('drafts', [siteId, subjectId, documentId]));
  }

  // A mutation is frozen at enqueue: retries must send this identical intent,
  // even if the user has subsequently edited the locally saved draft.
  async enqueue({siteId, subjectId, documentId, mutationId, baseRevision}) {
    scope(siteId, subjectId);
    documentID(documentId);
    if (!UUID_RE.test(mutationId)) throw new TypeError('Invalid mutation UUID');
    requireRevision(baseRevision);
    return this.driver.write(['drafts', 'mutations'], async tx => {
      const existing = await tx.get('mutations', [siteId, subjectId, mutationId]);
      if (existing) {
        if (existing.documentId !== documentId || existing.baseRevision !== baseRevision) {
          throw new Error('Mutation UUID already belongs to a different intent');
        }
        return existing;
      }
      const pending = await tx.forScope('mutations', siteId, subjectId);
      if (pending.some(m => m.documentId === documentId && m.state !== 'accepted')) {
        throw new Error('Resolve or acknowledge the existing mutation before queuing another');
      }
      const draft = await tx.get('drafts', [siteId, subjectId, documentId]);
      if (!draft) throw new Error('Save the local draft before enqueuing it');
      if (draft.serverRevision !== baseRevision) {
        throw new Error('Base revision does not match the locally confirmed server revision');
      }
      const mutation = {
        siteId, subjectId, documentId, mutationId,
        baseRevision, visibility: 'private', source: draft.source,
        deleted: draft.deleted, localRevision: draft.localRevision, state: 'queued',
      };
      await tx.put('mutations', mutation);
      return mutation;
    });
  }

  async getMutation(siteId, subjectId, mutationId) {
    scope(siteId, subjectId);
    if (!UUID_RE.test(mutationId)) throw new TypeError('Invalid mutation UUID');
    return this.driver.write(['mutations'], tx => tx.get('mutations', [siteId, subjectId, mutationId]));
  }

  async listUnsent(siteId, subjectId) {
    scope(siteId, subjectId);
    return this.driver.write(['mutations'], async tx =>
      (await tx.forScope('mutations', siteId, subjectId))
        .filter(m => m.state === 'queued' || m.state === 'conflict')
        .sort((a, b) => a.documentId.localeCompare(b.documentId)));
  }

  // A receipt applies only to its original immutable mutation. A newer local
  // edit stays dirty and is never replaced by the server's older source.
  async acknowledge({siteId, subjectId, mutationId, acceptedRevision}) {
    scope(siteId, subjectId);
    if (!UUID_RE.test(mutationId)) throw new TypeError('Invalid mutation UUID');
    requireRevision(acceptedRevision);
    return this.driver.write(['mutations', 'drafts'], async tx => {
      const mutation = await tx.get('mutations', [siteId, subjectId, mutationId]);
      if (!mutation) throw new Error('Unknown local mutation');
      if (mutation.state === 'accepted') {
        if (mutation.acceptedRevision !== acceptedRevision) throw new Error('Receipt mismatch');
        return mutation;
      }
      if (mutation.state !== 'queued' || acceptedRevision !== mutation.baseRevision + 1) {
        throw new Error('Receipt does not match the queued mutation');
      }
      const accepted = {...mutation, state: 'accepted', acceptedRevision};
      await tx.put('mutations', accepted);
      const draft = await tx.get('drafts', [siteId, subjectId, mutation.documentId]);
      if (draft) {
        draft.serverRevision = acceptedRevision;
        if (draft.localRevision === mutation.localRevision) {
          draft.syncedLocalRevision = mutation.localRevision;
        }
        await tx.put('drafts', draft);
      }
      return accepted;
    });
  }

  async recordConflict({siteId, subjectId, mutationId, currentRevision}) {
    scope(siteId, subjectId);
    if (!UUID_RE.test(mutationId)) throw new TypeError('Invalid mutation UUID');
    requireRevision(currentRevision);
    return this.driver.write(['mutations'], async tx => {
      const mutation = await tx.get('mutations', [siteId, subjectId, mutationId]);
      if (!mutation || mutation.state === 'accepted') throw new Error('Unknown pending mutation');
      const conflicted = {...mutation, state: 'conflict', currentRevision};
      await tx.put('mutations', conflicted);
      return conflicted;
    });
  }
}
