import { build } from 'esbuild';
import { rm } from 'node:fs/promises';
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
    name: 'highlights-only',
    setup(builder) {
      builder.onResolve({ filter: /^\.\/backends\/shiki\.js$/ }, args => {
        if (args.importer.endsWith('/shared_highlighter.js'))
          return { path: 'unused-shiki-backend', namespace: 'markata-publish-diff' };
      });
      builder.onResolve({ filter: /^shiki\/langs$/ }, () => ({
        path: 'unused-shiki-languages', namespace: 'markata-publish-diff'
      }));
      builder.onLoad({ filter: /.*/, namespace: 'markata-publish-diff' }, args => ({
        contents: args.path === 'unused-shiki-backend'
          ? 'export class ShikiHighlighter { static async create() { throw new Error("Only the bundled Highlights backend is available."); } }'
          : 'export const bundledLanguages = {};',
        loader: 'js'
      }));
    }
  }]
});
