import test from 'node:test';
import assert from 'node:assert/strict';
import {mutationDigest, OfflineMutationSync} from './offline-sync.mjs';

const initial = Object.freeze({
  siteId: 'family-blog', subjectId: 'alice', documentId: 'post-001',
  mutationId: '8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafa',
  baseRevision: 0, source: '# First', deleted: false, visibility: 'private',
  localRevision: 1, state: 'queued',
});

class FakeOutbox {
  constructor(mutations = [initial]) {
    this.mutations = mutations.map(m => ({...m}));
    this.acks = [];
    this.conflicts = [];
  }
  async listUnsent(site, subject) {
    return this.mutations.filter(m => m.siteId === site && m.subjectId === subject && m.state === 'queued')
      .map(m => ({...m}));
  }
  async acknowledge(args) {
    const found = this.mutations.find(m => m.mutationId === args.mutationId);
    if (!found || found.state !== 'queued' || found.baseRevision + 1 !== args.acceptedRevision) {
      throw new Error('Unexpected receipt');
    }
    found.state = 'accepted';
    this.acks.push({...args});
  }
  async recordConflict(args) {
    const found = this.mutations.find(m => m.mutationId === args.mutationId);
    if (!found || found.state !== 'queued') throw new Error('Unexpected conflict');
    found.state = 'conflict';
    this.conflicts.push({...args});
  }
}

class ExclusiveLocks {
  held = false;
  async request(_name, options, callback) {
    assert.equal(options.mode, 'exclusive');
    assert.equal(options.ifAvailable, true);
    if (this.held) return callback(null);
    this.held = true;
    try { return await callback({name: 'exclusive'}); }
    finally { this.held = false; }
  }
}

const authorized = async () => ({siteId: initial.siteId, subjectId: initial.subjectId, verified: true});
function setup({store = new FakeOutbox(), locks = new ExclusiveLocks(), verifiedIdentity = authorized, fetchFn} = {}) {
  return {store, controller: new OfflineMutationSync({
    store, locks, verifiedIdentity, fetchFn: fetchFn || (async () => { throw new Error('offline'); }),
    siteId: initial.siteId, subjectId: initial.subjectId,
  })};
}
function response(status, json) {
  return {status, json: async () => json};
}

test('mutation fingerprint matches Go length-framed SHA-256 including revision and visibility', async () => {
  assert.equal(await mutationDigest(initial), '3564810fe245464ceb896332ead1040f39e5ba4b217e49059d2eccc35060a1a5');
  assert.notEqual(await mutationDigest({...initial, deleted: true}), await mutationDigest(initial));
  assert.notEqual(await mutationDigest({...initial, baseRevision: 1}), await mutationDigest(initial));
  assert.notEqual(await mutationDigest({...initial, source: '# Changed'}), await mutationDigest(initial));
});

test('lost HTTP response replays exact frozen intent and acknowledges only committed matching receipt', async () => {
  const sent = [];
  let attempts = 0;
  const {store, controller} = setup({fetchFn: async (url, request) => {
    assert.equal(url, '/api/v1/sites/family-blog/mutations');
    assert.equal(request.method, 'POST');
    assert.equal(request.mode, 'same-origin');
    assert.equal(request.credentials, 'same-origin');
    assert.equal(request.redirect, 'error');
    assert.equal(request.headers['X-Markata-Sync-Intent'], 'write');
    sent.push(request.body);
    if (++attempts === 1) throw new Error('response lost after commit');
    return response(200, {
      mutation_id: initial.mutationId,
      revision: 1,
      digest: await mutationDigest(initial),
    });
  }});
  assert.deepEqual(await controller.syncOnce(), {status: 'retry_network', synced: 0});
  assert.equal(store.acks.length, 0);
  assert.equal(store.mutations[0].state, 'queued');
  assert.deepEqual(await controller.syncOnce(), {status: 'synced_drafts', synced: 1});
  assert.equal(store.acks.length, 1);
  assert.equal(store.mutations[0].state, 'accepted');
  assert.deepEqual(sent[0], sent[1]);
  const body = JSON.parse(sent[0]);
  assert.deepEqual(Object.keys(body).sort(), ['base_revision', 'deleted', 'document_id', 'mutation_id', 'source'].sort());
  assert.equal(body.source, '# First');
});

test('another browser tab holding the lock prevents a second upload', async () => {
  const locks = new ExclusiveLocks();
  let finish;
  let posts = 0;
  const delayed = new Promise(resolve => { finish = resolve; });
  const store = new FakeOutbox();
  const send = async () => {
    posts++;
    await delayed;
    return response(200, {
      mutation_id: initial.mutationId, revision: 1, digest: await mutationDigest(initial),
    });
  };
  const a = setup({store, locks, fetchFn: send}).controller;
  const b = setup({store, locks, fetchFn: send}).controller;
  const running = a.syncOnce();
  // Give the first tab enough microtasks to hold its lock during fetch.
  await new Promise(resolve => setImmediate(resolve));
  const other = await b.syncOnce();
  assert.deepEqual(other, {status: 'busy_other_tab', synced: 0});
  finish();
  assert.deepEqual(await running, {status: 'synced_drafts', synced: 1});
  assert.equal(posts, 1);
});

test('authentication loss and account switching preserve all queued bytes', async () => {
  const store = new FakeOutbox();
  for (const identity of [null, {siteId: initial.siteId, subjectId: 'bob', verified: true},
    {siteId: initial.siteId, subjectId: initial.subjectId, verified: false}]) {
    const {controller} = setup({
      store, verifiedIdentity: async () => identity,
      fetchFn: async () => { throw new Error('Should never call network'); },
    });
    assert.deepEqual(await controller.syncOnce(), {status: 'paused_auth', synced: 0});
  }
  assert.equal(store.mutations[0].state, 'queued');
  const blocked = setup({store, fetchFn: async () => response(401, {})}).controller;
  assert.deepEqual(await blocked.syncOnce(), {status: 'paused_auth', synced: 0});
  assert.equal(store.mutations[0].source, '# First');
});

test('409 preserved conflict is kept for reconciliation and never retried automatically', async () => {
  let posts = 0;
  const {store, controller} = setup({fetchFn: async () => {
    posts++;
    return response(409, {code: 'revision_conflict', current_revision: 5, preserved: true});
  }});
  assert.deepEqual(await controller.syncOnce(), {status: 'conflict', synced: 0});
  assert.equal(store.conflicts.length, 1);
  assert.equal(store.conflicts[0].currentRevision, 5);
  assert.equal(store.mutations[0].source, '# First');
  assert.deepEqual(await controller.syncOnce(), {status: 'synced_drafts', synced: 0});
  assert.equal(posts, 1);
});

test('mismatched digest, ID or revision cannot discard unsynced source', async () => {
  for (const receipt of [
    {mutation_id: initial.mutationId, revision: 1, digest: 'incorrect'},
    {mutation_id: 'different-uuid', revision: 1, digest: await mutationDigest(initial)},
    {mutation_id: initial.mutationId, revision: 7, digest: await mutationDigest(initial)},
  ]) {
    const {store, controller} = setup({fetchFn: async () => response(200, receipt)});
    assert.deepEqual(await controller.syncOnce(), {status: 'invalid_receipt', synced: 0});
    assert.equal(store.mutations[0].state, 'queued');
    assert.equal(store.acks.length, 0);
  }
});

test('network/server errors and missing browser Web Locks do not drop the queue', async () => {
  for (const code of [429, 500, 503]) {
    const {store, controller} = setup({fetchFn: async () => response(code, {})});
    assert.deepEqual(await controller.syncOnce(), {status: 'retry_server', synced: 0});
    assert.equal(store.mutations[0].state, 'queued');
  }
  const noLock = setup({locks: null, fetchFn: async () => { throw new Error('Unexpected upload'); }});
  assert.deepEqual(await noLock.controller.syncOnce(), {status: 'paused_no_web_locks', synced: 0});
  assert.equal(noLock.store.mutations[0].source, '# First');
});
