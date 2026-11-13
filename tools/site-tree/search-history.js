import fuzzysort from 'fuzzysort';
const prefix = 'markata-picks-v1:';
const preferenceKey = 'markata-picks-enabled-v1';
const fallback = new Map();
const storage = () => { try { return globalThis.localStorage; } catch { return null; } };
export function enabled() { try { return storage()?.getItem(preferenceKey) !== 'false'; } catch { return true; } }
export function setEnabled(value) {
  try { storage()?.setItem(preferenceKey, String(value)); } catch { /* Optional persistence. */ }
  fallback.set(preferenceKey, value);
}
function isEnabled() { return fallback.has(preferenceKey) ? fallback.get(preferenceKey) : enabled(); }
function read(scope) {
  if (!isEnabled()) return [];
  let value = fallback.get(scope) || [];
  try { const raw = fallback.has(scope) ? null : storage()?.getItem(prefix + scope); if (raw) value = JSON.parse(raw); } catch { /* Keep memory history. */ }
  return Array.isArray(value) ? value.filter(x => x && typeof x.id === 'string' &&
    Number.isFinite(x.at) && Number.isFinite(x.count) && x.count > 0).slice(0, 100) : [];
}
export function remember(scope, id) {
  if (!isEnabled() || typeof id !== 'string') return;
  const old = read(scope), previous = old.find(x => x.id === id);
  const next = [{id, at:Date.now(), count:Math.min(100, (previous?.count || 0) + 1)},
    ...old.filter(x => x.id !== id)].slice(0,100);
  fallback.set(scope,next);
  try { storage()?.setItem(prefix + scope,JSON.stringify(next)); } catch { /* Works for this tab. */ }
}
export function clearHistory(scope) {
  fallback.set(scope,[]);
  try { storage()?.removeItem(prefix + scope); } catch { /* Optional persistence. */ }
}
export const fileScope = () => 'files:' + (globalThis.window?.MarkataSiteTree?.workspaceId?.() || globalThis.location?.origin || 'default');
export function createSearch(entries, keys, weights, scope, identify) {
  const snapshot = fuzzysort.snapshot(entries,{keys});
  return query => {
    const text = query.trim(), start = performance.now();
    const history = new Map(read(typeof scope === 'function' ? scope() : scope).map((x,i) => [x.id,{...x,index:i}]));
    const results = fuzzysort.go(text, snapshot, {threshold:0, limit:Infinity,
      scoreFn: r => Math.max(...weights.map((weight,i) => r[i].score * weight))});
    const rank = result => {
      const exact = String(result.obj[keys[0]] || '').toLowerCase().replace(keys[0] === 'name' ? /\.(md|markdown|mdown|mkd)$/ : /$^/, '') === text.toLowerCase();
      const use = history.get(identify(result.obj));
      // History breaks close match ties; it cannot promote a weak match over a strong one.
      const quality = text ? (exact ? 100 : Math.floor(result.score * 20)) : 0;
      const bonus = use ? 1 / (1 + (Math.max(0,Date.now() - use.at)) / 604800000) + Math.log1p(use.count) * .1 : 0;
      return {quality, bonus};
    };
    const ranks = new Map(results.map(r => [r,rank(r)]));
    results.sort((a,b) => ranks.get(b).quality - ranks.get(a).quality ||
      ranks.get(b).bonus - ranks.get(a).bonus || b.score-a.score || identify(a.obj).localeCompare(identify(b.obj)));
    const top = results.slice(0,60); top.total = results.total;
    globalThis.window?.MarkataDiagnostics?.record('search',performance.now()-start);
    return top;
  };
}
export const historyAPI = {createSearch, remember, clearHistory, enabled:isEnabled, setEnabled, fileScope};
if (typeof window !== 'undefined') window.MarkataSearch = historyAPI;
