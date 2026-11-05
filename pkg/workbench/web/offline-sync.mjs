// Opt-in, identity-scoped drain of frozen IndexedDB mutations.
// This module is not wired to the workbench or to any public editor origin.

const encoder = new TextEncoder();

function frame(value) {
  const length = new Uint8Array(8);
  new DataView(length.buffer).setBigUint64(0, BigInt(value.length), false);
  return [length, value];
}

function uint64(value) {
  if (!Number.isSafeInteger(value) || value < 0) throw new TypeError('Unsafe base revision');
  const bytes = new Uint8Array(8);
  new DataView(bytes.buffer).setBigUint64(0, BigInt(value), false);
  return bytes;
}

// Matches the Go editorsync.Mutation.Digest length-framed SHA-256 contract.
// A response with the wrong digest is never allowed to acknowledge an outbox item.
export async function mutationDigest(mutation, cryptoAPI = globalThis.crypto) {
  if (!cryptoAPI?.subtle) throw new Error('WebCrypto unavailable; cannot authenticate mutation receipt');
  const pieces = [
    encoder.encode(mutation.siteId), encoder.encode(mutation.subjectId),
    encoder.encode(mutation.documentId), encoder.encode(mutation.mutationId),
    encoder.encode('private'), uint64(mutation.baseRevision),
    Uint8Array.of(mutation.deleted ? 1 : 0), encoder.encode(mutation.source),
  ];
  const framed = pieces.flatMap(frame);
  const size = framed.reduce((sum, bytes) => sum + bytes.length, 0);
  const bytes = new Uint8Array(size);
  let offset = 0;
  for (const part of framed) {
    bytes.set(part, offset);
    offset += part.length;
  }
  const hash = new Uint8Array(await cryptoAPI.subtle.digest('SHA-256', bytes));
  return [...hash].map(b => b.toString(16).padStart(2, '0')).join('');
}

function matchesScope(identity, siteId, subjectId) {
  return identity && identity.siteId === siteId && identity.subjectId === subjectId &&
    identity.verified === true;
}

export class OfflineMutationSync {
  constructor({store, siteId, subjectId, verifiedIdentity, locks = globalThis.navigator?.locks,
    fetchFn = globalThis.fetch, cryptoAPI = globalThis.crypto}) {
    if (!store || typeof store.listUnsent !== 'function' || typeof store.acknowledge !== 'function' ||
        typeof store.recordConflict !== 'function' || typeof verifiedIdentity !== 'function' ||
        typeof fetchFn !== 'function' || !cryptoAPI?.subtle ||
        !/^[a-z0-9](?:[a-z0-9-]{0,62}[a-z0-9])?$/.test(siteId) ||
        typeof subjectId !== 'string' || !subjectId || subjectId.length > 256) {
      throw new TypeError('Offline sync requires a scoped outbox, verified identity callback and WebCrypto');
    }
    this.store = store;
    this.siteId = siteId;
    this.subjectId = subjectId;
    this.verifiedIdentity = verifiedIdentity;
    this.locks = locks;
    this.fetchFn = fetchFn;
    this.cryptoAPI = cryptoAPI;
  }

  // Never send without exclusive Web Locks support. SQLite UUID receipts
  // remain idempotent if a tab crashes after server commit but before local ACK.
  async syncOnce() {
    if (typeof this.locks?.request !== 'function') {
      return {status: 'paused_no_web_locks', synced: 0};
    }
    const key = 'markata-sync-v1:' + this.siteId + ':' + this.subjectId;
    return this.locks.request(key, {mode: 'exclusive', ifAvailable: true}, async lock => {
      if (!lock) return {status: 'busy_other_tab', synced: 0};
      return this.#drainLocked();
    });
  }

  async #drainLocked() {
    let synced = 0;
    if (!matchesScope(await this.verifiedIdentity(), this.siteId, this.subjectId)) {
      return {status: 'paused_auth', synced};
    }
    const pending = await this.store.listUnsent(this.siteId, this.subjectId);
    for (const mutation of pending) {
      // The session may expire/change while the device processes many posts.
      if (!matchesScope(await this.verifiedIdentity(), this.siteId, this.subjectId)) {
        return {status: 'paused_auth', synced};
      }
      const digest = await mutationDigest(mutation, this.cryptoAPI);
      const body = JSON.stringify({
        document_id: mutation.documentId, mutation_id: mutation.mutationId,
        base_revision: mutation.baseRevision, source: mutation.source, deleted: mutation.deleted,
      });
      let response;
      try {
        response = await this.fetchFn('/api/v1/sites/' + encodeURIComponent(this.siteId) + '/mutations', {
          method: 'POST', mode: 'same-origin', credentials: 'same-origin',
          redirect: 'error', cache: 'no-store',
          headers: {'Content-Type': 'application/json', 'X-Markata-Sync-Intent': 'write'},
          body,
        });
      } catch {
        return {status: 'retry_network', synced};
      }
      if (response.status === 401 || response.status === 403) {
        return {status: 'paused_auth', synced};
      }
      if (response.status === 409) {
        let rejected;
        try { rejected = await response.json(); } catch { return {status: 'retry_server', synced}; }
        if (rejected?.code !== 'revision_conflict' || !Number.isSafeInteger(rejected.current_revision) ||
            rejected.current_revision < 0 || rejected.preserved !== true) {
          return {status: 'retry_server', synced};
        }
        await this.store.recordConflict({
          siteId: this.siteId, subjectId: this.subjectId, mutationId: mutation.mutationId,
          currentRevision: rejected.current_revision,
        });
        return {status: 'conflict', synced};
      }
      if (response.status !== 200) {
        return {status: 'retry_server', synced};
      }
      let receipt;
      try { receipt = await response.json(); } catch { return {status: 'retry_server', synced}; }
      if (receipt?.mutation_id !== mutation.mutationId || receipt.digest !== digest ||
          !Number.isSafeInteger(receipt.revision) ||
          receipt.revision !== mutation.baseRevision + 1) {
        return {status: 'invalid_receipt', synced};
      }
      // This is the only success transition: the local transaction must commit.
      await this.store.acknowledge({
        siteId: this.siteId, subjectId: this.subjectId,
        mutationId: mutation.mutationId, acceptedRevision: receipt.revision,
      });
      synced++;
    }
    return {status: 'synced_drafts', synced};
  }
}
