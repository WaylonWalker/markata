import test from 'node:test';
import assert from 'node:assert/strict';
import {fetchVerifiedSyncIdentity} from './offline-session.mjs';
import {OfflineMutationSync, mutationDigest} from './offline-sync.mjs';

const writer = {
  subject_id: 'alice-opaque',
  sites: [
    {id: 'family-blog', name: 'Family Blog', permissions: {read: true, write: true, publish: false}},
    {id: 'another-blog', name: 'Another Blog', permissions: {read: true, write: false, publish: false}},
  ],
};
function response(status, body, {contentType = 'application/json; charset=utf-8', redirected = false} = {}) {
  return {
    status, redirected, type: 'basic',
    headers: {get: name => name.toLowerCase() === 'content-type' ? contentType : null},
    json: async () => body,
  };
}
function checkedFetch(payload = writer) {
  let calls = 0;
  const fetchFn = async (route, options) => {
    calls++;
    assert.equal(route, '/api/v1/session');
    assert.equal(options.method, 'GET');
    assert.equal(options.mode, 'same-origin');
    assert.equal(options.credentials, 'same-origin');
    assert.equal(options.redirect, 'error');
    assert.equal(options.cache, 'no-store');
    assert.equal(options.headers.Accept, 'application/json');
    return response(200, payload);
  };
  return {fetchFn, calls: () => calls};
}

test('fresh server-confirmed subject with a site write grant is returned without credentials', async () => {
  const fixture = checkedFetch();
  assert.deepEqual(await fetchVerifiedSyncIdentity('family-blog', fixture), {
    siteId: 'family-blog', subjectId: 'alice-opaque', verified: true,
  });
  assert.deepEqual(await fetchVerifiedSyncIdentity('family-blog', fixture), {
    siteId: 'family-blog', subjectId: 'alice-opaque', verified: true,
  });
  assert.equal(fixture.calls(), 2, 'session identity must be revalidated rather than cached');
  assert.equal(await fetchVerifiedSyncIdentity('another-blog', fixture), null,
    'read-only site must not authorize sync POSTs');
  assert.equal(await fetchVerifiedSyncIdentity('hidden-site', fixture), null);
});

test('logout, account switch, session expiry and bad JSON fail closed', async () => {
  const switched = [
    response(200, writer),
    response(200, {...writer, subject_id: 'bob-opaque'}),
    response(401, null),
  ];
  const fetchFn = async () => switched.shift();
  assert.equal((await fetchVerifiedSyncIdentity('family-blog', {fetchFn})).subjectId, 'alice-opaque');
  assert.equal((await fetchVerifiedSyncIdentity('family-blog', {fetchFn})).subjectId, 'bob-opaque');
  assert.equal(await fetchVerifiedSyncIdentity('family-blog', {fetchFn}), null);
  for (const status of [403, 404, 500]) {
    assert.equal(await fetchVerifiedSyncIdentity('family-blog', {
      fetchFn: async () => response(status, writer),
    }), null);
  }
  assert.equal(await fetchVerifiedSyncIdentity('family-blog', {
    fetchFn: async () => { throw new Error('offline'); },
  }), null);
  assert.equal(await fetchVerifiedSyncIdentity('family-blog', {
    fetchFn: async () => ({status: 200, json: async () => { throw new Error('invalid JSON'); }}),
  }), null);
});

test('cross-origin or malformed session responses never yield a verified identity', async () => {
  const invalid = [
    response(200, writer, {redirected: true}),
    response(200, writer, {contentType: 'text/html'}),
    response(200, writer, {contentType: ''}),
    {...response(200, writer), type: 'opaqueredirect'},
    response(200, {...writer, subject_id: 'alice\n'}),
    response(200, {...writer, subject_id: ''}),
    response(200, {...writer, sites: []}),
    response(200, {...writer, sites: [writer.sites[0], writer.sites[0]]}),
    response(200, {...writer, sites: [undefined]}),
    response(200, {...writer, sites: [{id: 'family-blog', permissions: {read: true, write: 'yes'}}]}),
    response(200, {...writer, sites: Array(1001).fill(writer.sites[0])}),
  ];
  for (const resp of invalid) {
    assert.equal(await fetchVerifiedSyncIdentity('family-blog', {fetchFn: async () => resp}), null);
  }
  assert.equal(await fetchVerifiedSyncIdentity('family-blog', {
    fetchFn: async () => ({...response(200, writer), type: 'opaque'}),
  }), null);
  await assert.rejects(fetchVerifiedSyncIdentity('BAD SITE', {fetchFn: async () => response(200, writer)}), /site ID/);
  await assert.rejects(fetchVerifiedSyncIdentity(undefined, {fetchFn: async () => response(200, writer)}), /site ID/);
});

test('offline mutation sync can use a fresh session bootstrap without localStorage identity', async () => {
  const mutation = {
    siteId: 'family-blog', subjectId: 'alice-opaque', documentId: 'post-a',
    mutationId: '8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafa',
    baseRevision: 0, source: '# Changed while offline', deleted: false,
    state: 'queued', localRevision: 1,
  };
  let identityCalls = 0;
  let uploads = 0;
  let accepted = false;
  const fetchFn = async (route, options) => {
    if (route === '/api/v1/session') {
      identityCalls++;
      assert.equal(options.cache, 'no-store');
      return response(200, writer);
    }
    assert.equal(route, '/api/v1/sites/family-blog/mutations');
    assert.equal(options.credentials, 'same-origin');
    assert.equal(JSON.parse(options.body).source, mutation.source);
    uploads++;
    return response(200, {
      mutation_id: mutation.mutationId, revision: 1, digest: await mutationDigest(mutation),
    });
  };
  const store = {
    listUnsent: async () => accepted ? [] : [mutation],
    acknowledge: async ({acceptedRevision}) => {
      assert.equal(acceptedRevision, 1);
      accepted = true;
    },
    recordConflict: async () => { throw new Error('unexpected conflict'); },
  };
  const locks = {request: async (_name, _options, fn) => fn({name: 'exclusive'})};
  const sync = new OfflineMutationSync({
    siteId: mutation.siteId, subjectId: mutation.subjectId, store, locks, fetchFn,
    verifiedIdentity: () => fetchVerifiedSyncIdentity(mutation.siteId, {fetchFn}),
  });
  assert.deepEqual(await sync.syncOnce(), {status: 'synced_drafts', synced: 1});
  assert.equal(uploads, 1);
  assert.equal(accepted, true);
  assert.ok(identityCalls >= 3, 'must recheck session at async sync boundaries');
});
