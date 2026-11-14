import {createSearch, fileScope} from './search-history.js';

// Operators narrow the list before fuzzy matching: #tag or tag:name, status:draft, in:pages/blog.
const OPERATOR = /(^|\s)(?:#|tag:)([^\s#]+)|(^|\s)status:(draft|published)\b|(^|\s)in:([^\s]+)|(^|\s)(?:feed:|\$)([^\s]+)/gi;

export function parseFileQuery(raw) {
  const filters = [];
  const text = raw.replace(OPERATOR, (match, _a, tag, _b, status, _c, folder, _d, feed) => {
    if (tag) filters.push({kind:'tag', value:tag.toLowerCase()});
    else if (status) filters.push({kind:'status', value:status.toLowerCase()});
    else if (folder) filters.push({kind:'folder', value:folder.replace(/^\/+|\/+$/g,'').toLowerCase()});
    else if (feed) filters.push({kind:'feed', value:feed.toLowerCase()});
    return ' ';
  }).replace(/\s+/g, ' ').trim();
  return {text, filters};
}

function matches(entry, filter, members) {
  if (filter.kind === 'tag') return (entry.tags || []).some(tag => tag.toLowerCase() === filter.value);
  if (filter.kind === 'status') return entry.status === filter.value;
  if (filter.kind === 'feed') return members(filter.value).has(entry.path);
  return entry.path.toLowerCase().startsWith(filter.value + '/');
}

// feeds() returns [{slug, paths}] for the workspace's configured feeds; it is read at search time.
export function createFileSearch(files, feeds = () => []) {
  const entries = files.map(file => ({...file, name:file.path.split('/').pop(), title:file.title || ''}));
  const make = pool => createSearch(pool, ['name','path','title'], [1,.9,.8], fileScope, file => file.path);
  const everything = make(entries);
  const filtered = new Map();
  const members = slug => {
    const feed = feeds().find(item => String(item.slug).toLowerCase() === slug);
    return new Set(feed?.paths || []);
  };
  return query => {
    const {text, filters} = parseFileQuery(query);
    if (!filters.length) return everything(text);
    // Feed membership can change while the picker is open, so feed searches are not cached.
    if (filters.some(filter => filter.kind === 'feed')) {
      return make(entries.filter(entry => filters.every(filter => matches(entry, filter, members))))(text);
    }
    const key = JSON.stringify(filters);
    if (!filtered.has(key)) {
      if (filtered.size > 20) filtered.clear();
      filtered.set(key, make(entries.filter(entry => filters.every(filter => matches(entry, filter, members)))));
    }
    return filtered.get(key)(text);
  };
}

// Tags and folders present in the workspace, for completing operators as they are typed.
export function fileFacets(files) {
  const tagCounts = new Map(), folders = new Set();
  for (const file of files) {
    for (const tag of file.tags || []) tagCounts.set(tag, (tagCounts.get(tag) || 0) + 1);
    const parts = String(file.path || '').split('/').slice(0, -1);
    for (let depth = 1; depth <= Math.min(parts.length, 3); depth++) folders.add(parts.slice(0, depth).join('/'));
  }
  const tags = [...tagCounts].sort((a, b) => b[1] - a[1] || a[0].localeCompare(b[0])).map(([tag]) => tag);
  return {tags, folders: [...folders].sort()};
}

// Completions for the word being typed: #tag, tag:, status:, in:, feed:, and $feed.
export function suggestFor(raw, facets) {
  const token = raw.match(/(?:^|\s)(\S*)$/)?.[1] ?? '';
  const pick = (values, prefix, make) => values
    .filter(value => value.toLowerCase().startsWith(prefix.toLowerCase()))
    .slice(0, 8).map(value => ({label: make(value), value: make(value)}));
  let items = [];
  if (token.startsWith('#')) items = pick(facets.tags, token.slice(1), value => '#' + value);
  else if (/^tag:/i.test(token)) items = pick(facets.tags, token.slice(4), value => 'tag:' + value);
  else if (/^status:/i.test(token)) items = pick(['draft', 'published'], token.slice(7), value => 'status:' + value);
  else if (/^in:/i.test(token)) items = pick(facets.folders, token.slice(3), value => 'in:' + value);
  else if (/^feed:/i.test(token)) items = pick(facets.feeds || [], token.slice(5), value => 'feed:' + value);
  else if (token.startsWith('$')) items = pick(facets.feeds || [], token.slice(1), value => '$' + value);
  return {token, items};
}

export function acceptSuggestion(raw, completed) {
  return raw.replace(/\S*$/, completed) + ' ';
}
