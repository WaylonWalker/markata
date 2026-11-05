import test from 'node:test';
import assert from 'node:assert/strict';
import {LocalDraftOutbox, openBrowserDraftStore} from './offline-outbox.mjs';

const ids = [
  '8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafa',
  '8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafb',
  '8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafc',
];
const shared = {siteId: 'family-blog', subjectId: 'alice', documentId: 'post-a'};

// In-memory transactional driver follows the same injected driver contract as
// the browser adapter. Fault injection tests app-level commit-before-ack logic;
// true browser IndexedDB commit/upgrade tests are a separate Playwright gate.
class MemoryDriver {
  constructor() {
    this.stores = {drafts: new Map(), mutations: new Map()};
    this.failCommit = false;
    this.queue = Promise.resolve();
  }

  write(names, callback) {
    const job = this.queue.then(async () => {
      const staged = Object.fromEntries(Object.entries(this.stores).map(([name, rows]) =>
        [name, new Map([...rows].map(([key, value]) => [key, structuredClone(value)]))]));
      const recordKey = (name, value) => JSON.stringify(name === 'drafts'
        ? [value.siteId, value.subjectId, value.documentId]
        : [value.siteId, value.subjectId, value.mutationId]);
      const tx = {
        get: async (name, key) => structuredClone(staged[name].get(JSON.stringify(key))),
        put: async (name, value) => staged[name].set(recordKey(name, value), structuredClone(value)),
        forScope: async (name, siteId, subjectId) => [...staged[name].values()]
          .filter(value => value.siteId === siteId && value.subjectId === subjectId)
          .map(value => structuredClone(value)),
      };
      const result = await callback(tx);
      if (this.failCommit) throw new Error('simulated IDB transaction abort/quota error');
      this.stores = staged;
      return structuredClone(result);
    });
    this.queue = job.then(() => {}, () => {});
    return job;
  }
}

test('100 disconnected edits survive a new store instance, with no sync claim', async () => {
  const driver = new MemoryDriver();
  const first = new LocalDraftOutbox(driver);
  let last;
  for (let n = 1; n <= 100; n++) {
    last = await first.saveDraft({...shared, source: '# Draft ' + n, expectedLocalRevision: n - 1});
  }
  assert.equal(last.localRevision, 100);
  assert.equal(last.syncedLocalRevision, 0);
  const reopened = new LocalDraftOutbox(driver);
  const recovered = await reopened.getDraft(shared.siteId, shared.subjectId, shared.documentId);
  assert.equal(recovered.source, '# Draft 100');
  assert.equal(recovered.localRevision, 100);
  assert.deepEqual(await reopened.listUnsent(shared.siteId, shared.subjectId), []);
});

test('queued intent is immutable through further local editing and lost ack retry', async () => {
  const store = new LocalDraftOutbox(new MemoryDriver());
  await store.saveDraft({...shared, source: '# Initial', expectedLocalRevision: 0});
  const staged = await store.enqueue({...shared, mutationId: ids[0], baseRevision: 0});
  assert.equal(staged.source, '# Initial');
  assert.equal(staged.visibility, 'private');
  assert.equal(staged.state, 'queued');
  await store.saveDraft({...shared, source: '# Newer unsynced edit', expectedLocalRevision: 1});
  assert.deepEqual(await store.enqueue({...shared, mutationId: ids[0], baseRevision: 0}), staged);
  await assert.rejects(store.enqueue({...shared, mutationId: ids[0], baseRevision: 1}), /different intent/);
  await assert.rejects(store.enqueue({...shared, mutationId: ids[1], baseRevision: 0}), /existing mutation/);
  const ack = await store.acknowledge({siteId: shared.siteId, subjectId: shared.subjectId,
    mutationId: ids[0], acceptedRevision: 1});
  assert.equal(ack.state, 'accepted');
  assert.equal(ack.acceptedRevision, 1);
  const draft = await store.getDraft(shared.siteId, shared.subjectId, shared.documentId);
  assert.equal(draft.source, '# Newer unsynced edit');
  assert.equal(draft.localRevision, 2);
  assert.equal(draft.serverRevision, 1);
  assert.equal(draft.syncedLocalRevision, 0);
  assert.deepEqual(await store.listUnsent(shared.siteId, shared.subjectId), []);
  await assert.rejects(store.acknowledge({siteId: shared.siteId, subjectId: shared.subjectId,
    mutationId: ids[0], acceptedRevision: 2}), /Receipt mismatch/);
  const later = await store.enqueue({...shared, mutationId: ids[1], baseRevision: 1});
  assert.equal(later.source, '# Newer unsynced edit');
  assert.equal(later.localRevision, 2);
});

test('acknowledged content becomes synced only after committed receipt', async () => {
  const driver = new MemoryDriver();
  const store = new LocalDraftOutbox(driver);
  await store.saveDraft({...shared, source: 'hello', expectedLocalRevision: 0});
  await store.enqueue({...shared, mutationId: ids[0], baseRevision: 0});
  driver.failCommit = true;
  await assert.rejects(store.acknowledge({...shared, mutationId: ids[0], acceptedRevision: 1}), /abort/);
  driver.failCommit = false;
  assert.equal((await store.getMutation(shared.siteId, shared.subjectId, ids[0])).state, 'queued');
  assert.equal((await store.getDraft(shared.siteId, shared.subjectId, shared.documentId)).serverRevision, 0);
  const receipt = await store.acknowledge({...shared, mutationId: ids[0], acceptedRevision: 1});
  assert.equal(receipt.state, 'accepted');
  const draft = await store.getDraft(shared.siteId, shared.subjectId, shared.documentId);
  assert.equal(draft.syncedLocalRevision, draft.localRevision);
  await assert.rejects(store.enqueue({...shared, mutationId: ids[1], baseRevision: 1}), /No new local edits/);
});

test('quota/transaction errors never imply a saved local draft or queued mutation', async () => {
  const driver = new MemoryDriver();
  const store = new LocalDraftOutbox(driver);
  driver.failCommit = true;
  await assert.rejects(store.saveDraft({...shared, source: 'Do not lose this', expectedLocalRevision: 0}), /abort/);
  driver.failCommit = false;
  assert.equal(await store.getDraft(shared.siteId, shared.subjectId, shared.documentId), undefined);
  await store.saveDraft({...shared, source: 'Still here', expectedLocalRevision: 0});
  driver.failCommit = true;
  await assert.rejects(store.enqueue({...shared, mutationId: ids[0], baseRevision: 0}), /abort/);
  driver.failCommit = false;
  assert.equal(await store.getMutation(shared.siteId, shared.subjectId, ids[0]), undefined);
  assert.equal((await store.getDraft(shared.siteId, shared.subjectId, shared.documentId)).source, 'Still here');
});

test('private data is partitioned by site and stable subject in all read methods', async () => {
  const store = new LocalDraftOutbox(new MemoryDriver());
  await store.saveDraft({...shared, source: 'Alice private', expectedLocalRevision: 0});
  await store.saveDraft({...shared, subjectId: 'bob', source: 'Bob private', expectedLocalRevision: 0});
  await store.saveDraft({...shared, siteId: 'second-blog', source: 'Alice second site', expectedLocalRevision: 0});
  await store.enqueue({...shared, mutationId: ids[0], baseRevision: 0});
  await store.enqueue({...shared, subjectId: 'bob', mutationId: ids[0], baseRevision: 0});
  assert.equal((await store.getDraft(shared.siteId, 'alice', shared.documentId)).source, 'Alice private');
  assert.equal((await store.getDraft(shared.siteId, 'bob', shared.documentId)).source, 'Bob private');
  assert.equal((await store.getDraft('second-blog', 'alice', shared.documentId)).source, 'Alice second site');
  assert.equal((await store.listUnsent('second-blog', 'alice')).length, 0);
  assert.equal((await store.listUnsent('family-blog', 'bob')).length, 1);
  assert.equal(await store.getMutation('second-blog', 'alice', ids[0]), undefined);
});

test('conflicts retain the original source and cannot be turned into fake success', async () => {
  const store = new LocalDraftOutbox(new MemoryDriver());
  await store.saveDraft({...shared, source: 'offline change', deleted: true, expectedLocalRevision: 0});
  await store.enqueue({...shared, mutationId: ids[0], baseRevision: 0});
  const original = await store.recordConflict({...shared, mutationId: ids[0], currentRevision: 3});
  assert.equal(original.deleted, true);
  assert.equal(original.state, 'conflict');
  await store.saveDraft({...shared, source: 'another local edit', deleted: false, expectedLocalRevision: 1});
  assert.deepEqual(await store.recordConflict({...shared, mutationId: ids[0], currentRevision: 4}), original);
  assert.equal((await store.getMutation(shared.siteId, shared.subjectId, ids[0])).source, 'offline change');
  await assert.rejects(store.acknowledge({...shared, mutationId: ids[0], acceptedRevision: 1}), /Receipt does not match/);
  assert.equal((await store.listUnsent(shared.siteId, shared.subjectId)).length, 0);
  assert.equal((await store.listConflicts(shared.siteId, shared.subjectId)).length, 1);
});

test('bad identities, unknown UUID, unsafe revision, and unavailable IDB fail closed', async () => {
  const store = new LocalDraftOutbox(new MemoryDriver());
  await assert.rejects(store.saveDraft({...shared, subjectId: '', source: 'private', expectedLocalRevision: 0}), /verified site/);
  await assert.rejects(store.saveDraft({...shared, documentId: '\ninvalid', source: 'private', expectedLocalRevision: 0}), /Invalid document/);
  await assert.rejects(store.saveDraft({...shared, source: 123, expectedLocalRevision: 0}), /UTF-8/);
  await store.saveDraft({...shared, source: '# A', expectedLocalRevision: 0});
  await assert.rejects(store.enqueue({...shared, mutationId: 'not-uuid', baseRevision: 0}), /Invalid mutation/);
  await assert.rejects(store.enqueue({...shared, mutationId: ids[0], baseRevision: 9007199254740992}), /safe integer/);
  await assert.rejects(openBrowserDraftStore(null), /unavailable/);
});

test('optimistic local revisions prevent stale tabs from silently overwriting work', async () => {
  const driver = new MemoryDriver();
  const tabA = new LocalDraftOutbox(driver);
  const tabB = new LocalDraftOutbox(driver);
  const initial = await tabA.saveDraft({...shared, source: '# First', expectedLocalRevision: 0});
  assert.equal(initial.localRevision, 1);
  const competing = await Promise.allSettled([
    tabA.saveDraft({...shared, source: '# Tab A', expectedLocalRevision: 1}),
    tabB.saveDraft({...shared, source: '# Tab B', expectedLocalRevision: 1}),
  ]);
  assert.equal(competing.filter(x => x.status === 'fulfilled').length, 1);
  assert.equal(competing.filter(x => x.status === 'rejected').length, 1);
  assert.match(competing.find(x => x.status === 'rejected').reason.message, /another tab/);
  const final = await tabB.getDraft(shared.siteId, shared.subjectId, shared.documentId);
  assert.equal(final.localRevision, 2);
  assert.ok(['# Tab A', '# Tab B'].includes(final.source));
  assert.equal(final.syncedLocalRevision, 0);
});
