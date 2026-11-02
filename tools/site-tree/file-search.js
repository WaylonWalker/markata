import fuzzysort from 'fuzzysort';
export function createFileSearch(files) {
  const entries = files.map(file => ({...file, name:file.path.split('/').pop(), title:file.title || ''}));
  const targets = fuzzysort.snapshot(entries, {keys:['name','path','title']});
  return query => fuzzysort.go(query.trim(), targets, {limit:60, threshold:0,
    scoreFn:r => Math.max(r[0].score, r[1].score * .9, r[2].score * .8)});
}
