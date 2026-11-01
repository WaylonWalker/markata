import { build } from 'esbuild';
import { readdir, readFile, rm, writeFile } from 'node:fs/promises';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

const root = fileURLToPath(new URL('../../', import.meta.url));
const web = path.join(root, 'pkg/workbench/web');
await rm(path.join(web, 'pierre-diffs'), { recursive: true, force: true });
await rm(path.join(web, 'pierre-diffs.bundle.js'), { force: true });
await build({
  entryPoints: [path.join(root, 'tools/publish-diffs/pierre.js')],
  outdir: path.join(root, 'pkg/workbench/web'),
  entryNames: 'pierre-diffs.bundle',
  chunkNames: 'pierre-diffs/[name]-[hash]',
  platform: 'browser',
  format: 'esm',
  target: 'es2022',
  bundle: true,
  splitting: true,
  minify: true,
  legalComments: 'eof',
  define: { 'process.env.NODE_ENV': '"production"' },
  plugins: [{
    name: 'markdown-only-highlighting',
    setup(builder) {
      builder.onResolve({ filter: /^shiki\/langs$/ }, () => ({
        path: 'markdown-languages', namespace: 'markata-publish-diff'
      }));
      builder.onLoad({ filter: /^markdown-languages$/, namespace: 'markata-publish-diff' }, () => ({
        contents: 'export const bundledLanguages = { markdown: () => import("@shikijs/langs/markdown"), md: () => import("@shikijs/langs/markdown") };',
        resolveDir: path.join(root, 'tools/publish-diffs'),
        loader: 'js'
      }));
    }
  }]
});

// Some bundled dependencies include literal tabs at the end of template
// literal lines. They are data used by the Markdown highlighter, but fail
// `git diff --check`; escape them without changing the runtime strings.
const chunks = path.join(web, 'pierre-diffs');
for (const file of await readdir(chunks)) {
  const output = path.join(chunks, file);
  const contents = await readFile(output, 'utf8');
  await writeFile(output, contents.replace(/\t(?=\r?$)/gm, '\\t'));
}
