import {createSearch, fileScope} from './search-history.js';
export function createFileSearch(files) {
  const entries = files.map(file => ({...file, name:file.path.split('/').pop(), title:file.title || ''}));
  return createSearch(entries, ['name','path','title'], [1,.9,.8], fileScope, file => file.path);
}
