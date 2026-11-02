import {build} from 'esbuild';
import {readFile, writeFile, readdir} from 'node:fs/promises';
import {fileURLToPath} from 'node:url';
import path from 'node:path';
const root = fileURLToPath(new URL('../../', import.meta.url));
const files = {};
async function collect(dir) {
  for (const entry of await readdir(path.join(root, 'examples/editor-site', dir), {withFileTypes:true})) {
    const name = path.posix.join(dir, entry.name);
    if (entry.isDirectory()) await collect(name);
    else files[name] = await readFile(path.join(root, 'examples/editor-site', name), 'utf8');
  }
}
await collect('');
await build({entryPoints:[path.join(root,'tools/site-tree/site-tree.js')], bundle:true, minify:true,
  format:'iife', target:'es2020', outfile:path.join(root,'pkg/workbench/web/site-tree.js'),
  define:{SITE_FILES:JSON.stringify(files)}, legalComments:'eof'});
// Retain notices for the browser runtime and its bundled dependencies.
const licenses = [
  'fuzzysort/LICENSE', '@pierre/trees/LICENSE.md', '@pierre/trees/NOTICE.md',
  '@pierre/theming/LICENSE.md', 'preact/LICENSE', 'preact-render-to-string/LICENSE'
];
const notices = await Promise.all(licenses.map(async name =>
  name + '\n\n' + await readFile(new URL('node_modules/' + name, import.meta.url), 'utf8')));
await writeFile(path.join(root,'pkg/workbench/web/site-tree.LICENSE'), notices.join('\n\n').replace(/\r\n/g, '\n'));
const {execFileSync} = await import('node:child_process');
execFileSync('python3', [path.join(root,'scripts/package-python-editor-assets.py')], {stdio:'inherit'});
