// Stateless, same-origin identity bootstrap for the offline mutation drain.
// Do not persist a positive authentication decision, cookie or bearer token.
// Each invocation must contact the trusted hosted editor session endpoint.
const SITE_RE = /^[a-z0-9](?:[a-z0-9-]{0,62}[a-z0-9])?$/;

function validSubject(subject) {
  return typeof subject === 'string' && subject.length > 0 && subject.length <= 256 &&
    subject.trim() === subject && !/[\u0000-\u001f\u007f]/.test(subject);
}

// Inject as verifiedIdentity: () => fetchVerifiedSyncIdentity(siteId).
// A null response intentionally pauses OfflineMutationSync without deleting
// anything from IndexedDB on logout, unavailable session or network errors.
export async function fetchVerifiedSyncIdentity(siteId, {fetchFn = globalThis.fetch} = {}) {
  if (typeof siteId !== 'string' || !SITE_RE.test(siteId) || typeof fetchFn !== 'function') {
    throw new TypeError('Verified sync identity needs a site ID and same-origin fetch');
  }
  try {
    const response = await fetchFn('/api/v1/session', {
      method: 'GET',
      mode: 'same-origin',
      credentials: 'same-origin',
      redirect: 'error',
      cache: 'no-store',
      headers: {Accept: 'application/json'},
    });
    if (!response || response.status !== 200 || response.redirected === true ||
        response.type === 'opaque' || response.type === 'opaqueredirect') {
      return null;
    }
    const contentType = response.headers?.get?.('content-type');
    if (typeof contentType !== 'string' || !/^application\/json(?:\s*;|\s*$)/i.test(contentType)) {
      return null;
    }
    const payload = await response.json();
    if (!payload || !validSubject(payload.subject_id) || !Array.isArray(payload.sites) ||
        payload.sites.length > 1000) {
      return null;
    }
    const granted = payload.sites.filter(site => site && site.id === siteId);
    if (granted.length !== 1 || granted[0].permissions?.read !== true ||
        granted[0].permissions?.write !== true) {
      return null;
    }
    return {siteId, subjectId: payload.subject_id, verified: true};
  } catch {
    // Includes offline, session rejection, JSON failures and redirect errors.
    // The immutable outbox still owns the unsent source; never clear it here.
    return null;
  }
}
