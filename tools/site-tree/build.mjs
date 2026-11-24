import {build} from 'esbuild';
import {cp, mkdir, readFile, rm, writeFile, readdir} from 'node:fs/promises';
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
await build({entryPoints:[path.join(root,'tools/site-tree/search-history.js')], bundle:true, minify:true,
  format:'iife', target:'es2020', outfile:path.join(root,'pkg/workbench/web/search-history.js'), legalComments:'eof'});
const workbenchWeb = path.join(root, 'pkg/workbench/web');
await mkdir(workbenchWeb, {recursive:true});
const sharedExcalidraw = path.join(root, 'pkg/visualizations/excalidraw');
const excalidrawRendererVersion = 'excalidraw-0.18.1-static-svg-v2';
await mkdir(sharedExcalidraw, {recursive:true});
await rm(path.join(sharedExcalidraw, 'excalidraw-editor.js'), {force:true});
await rm(path.join(sharedExcalidraw, 'excalidraw-editor.css'), {force:true});
await rm(path.join(sharedExcalidraw, 'excalidraw-editor-chunks'), {recursive:true, force:true});
await build({entryPoints:[path.join(root,'tools/site-tree/excalidraw-editor.js')], bundle:true, splitting:true, minify:true,
  format:'esm', target:'es2022', outdir:sharedExcalidraw, entryNames:'excalidraw-editor',
  chunkNames:'excalidraw-editor-chunks/[name]-[hash]', legalComments:'eof'});
const excalidrawFonts = path.join(sharedExcalidraw, 'excalidraw-fonts');
await rm(excalidrawFonts, {recursive:true, force:true});
await mkdir(excalidrawFonts, {recursive:true});
const excalidrawFontSource = path.join(root,'tools/site-tree/node_modules/@excalidraw/excalidraw/dist/prod/fonts');
for (const font of ['Assistant','Cascadia','ComicShanns','Excalifont','Liberation','Lilita','Nunito','Virgil']) {
  await cp(path.join(excalidrawFontSource,font), path.join(excalidrawFonts,font), {recursive:true});
}
const fontAssets = [];
async function collectFontAssets(dir, prefix = '') {
  for (const entry of await readdir(dir, {withFileTypes:true})) {
    const name = prefix ? `${prefix}/${entry.name}` : entry.name;
    if (entry.isDirectory()) await collectFontAssets(path.join(dir,entry.name),name);
    else if (entry.isFile() && entry.name.endsWith('.woff2')) fontAssets.push(name);
  }
}
await collectFontAssets(excalidrawFonts);
fontAssets.sort();
await writeFile(path.join(excalidrawFonts,'manifest.json'), JSON.stringify({renderer_version:excalidrawRendererVersion,fonts:fontAssets},null,2)+'\n');
await build({entryPoints:[path.join(root,'tools/site-tree/excalidraw-editor.css')], bundle:true, minify:true,
  format:'esm', target:'es2022', outfile:path.join(sharedExcalidraw,'excalidraw-editor.css'),
  assetNames:'excalidraw-fonts/[name]', loader:{'.woff2':'file'}, legalComments:'eof', conditions:['production']});
async function stripGeneratedJavaScriptWhitespace(dir) {
  for (const entry of await readdir(dir, {withFileTypes:true})) {
    const file = path.join(dir, entry.name);
    if (entry.isDirectory()) await stripGeneratedJavaScriptWhitespace(file);
    else if (entry.isFile() && entry.name.endsWith('.js')) {
      const source = await readFile(file, 'utf8');
      const normalized = source.replace(/[ \t]+$/gm, '');
      if (normalized !== source) await writeFile(file, normalized);
    }
  }
}
// Third-party bundles contain trailing spaces in embedded shader and license
// text. Keep the checked-in generated assets clean for git diff --check.
await stripGeneratedJavaScriptWhitespace(sharedExcalidraw);
await rm(path.join(workbenchWeb, 'excalidraw-editor.js'), {force:true});
await rm(path.join(workbenchWeb, 'excalidraw-editor.css'), {force:true});
await rm(path.join(workbenchWeb, 'excalidraw-fonts'), {recursive:true, force:true});
// Retain notices for the browser runtime and its bundled dependencies.
const licenses = [
  ['fuzzysort','node_modules/fuzzysort/LICENSE'],
  ['@pierre/trees','node_modules/@pierre/trees/LICENSE.md'],
  ['@pierre/trees notice','node_modules/@pierre/trees/NOTICE.md'],
  ['@pierre/theming','node_modules/@pierre/theming/LICENSE.md'],
  ['preact','node_modules/preact/LICENSE'],
  ['preact-render-to-string','node_modules/preact-render-to-string/LICENSE'],
  ['@excalidraw/excalidraw','excalidraw.LICENSE'],
  ['react','node_modules/react/LICENSE'],
  ['react-dom','node_modules/react-dom/LICENSE']
];
const notices = await Promise.all(licenses.map(async ([name,file]) =>
  name + '\n\n' + await readFile(new URL(file, import.meta.url), 'utf8')));
await writeFile(path.join(root,'pkg/workbench/web/site-tree.LICENSE'), notices.join('\n\n').replace(/\r\n/g, '\n'));
const {execFileSync} = await import('node:child_process');
execFileSync('python3', [path.join(root,'scripts/package-python-editor-assets.py')], {stdio:'inherit'});
