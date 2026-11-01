import { build } from 'esbuild';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

const root = fileURLToPath(new URL('../../', import.meta.url));
await build({
  entryPoints: [path.join(root, 'tools/publish-diffs/pierre.js')],
  outfile: path.join(root, 'pkg/workbench/web/pierre-diffs.bundle.js'),
  platform: 'browser',
  format: 'iife',
  target: 'es2022',
  bundle: true,
  minify: true,
  legalComments: 'eof',
  define: { 'process.env.NODE_ENV': '"production"' }
});
