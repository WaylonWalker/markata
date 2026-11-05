// Run the native browser smoke page without npm/browser-test dependencies.
// GitHub Ubuntu runners ship Chrome; local Linux may use chromium instead.
import {createServer} from 'node:http';
import {readFile, mkdtemp, rm} from 'node:fs/promises';
import {tmpdir} from 'node:os';
import {join, dirname} from 'node:path';
import {fileURLToPath} from 'node:url';
import {execFile as execFileCallback} from 'node:child_process';
import {promisify} from 'node:util';

const execFile = promisify(execFileCallback);
const root = dirname(fileURLToPath(import.meta.url));
const assets = new Map([
  ['/offline-outbox.browser.html', ['offline-outbox.browser.html', 'text/html; charset=utf-8']],
  ['/offline-outbox.mjs', ['offline-outbox.mjs', 'text/javascript; charset=utf-8']],
]);
const server = createServer(async (request, response) => {
  const asset = assets.get(request.url);
  if (!asset) {
    response.writeHead(404);
    response.end('Not found');
    return;
  }
  try {
    const contents = await readFile(join(root, asset[0]));
    response.writeHead(200, {'Content-Type': asset[1], 'Cache-Control': 'no-store'});
    response.end(contents);
  } catch {
    response.writeHead(500);
    response.end('Read failed');
  }
});

async function resolveChrome() {
  for (const name of [process.env.CHROME_BIN, 'google-chrome', 'chromium', 'google-chrome-stable', 'chromium-browser']) {
    if (!name) continue;
    try {
      await execFile(name, ['--version'], {timeout: 5000});
      return name;
    } catch { /* Try the next installed browser. */ }
  }
  throw new Error('Chromium/Chrome is not installed; native IndexedDB smoke cannot run');
}

let profile;
try {
  await new Promise((resolve, reject) => {
    server.once('error', reject);
    server.listen(0, '127.0.0.1', resolve);
  });
  profile = await mkdtemp(join(tmpdir(), 'markata-outbox-smoke-'));
  const chrome = await resolveChrome();
  const address = server.address();
  const url = 'http://127.0.0.1:' + address.port + '/offline-outbox.browser.html';
  const args = [
    '--headless=new', '--no-sandbox', '--disable-gpu', '--disable-dev-shm-usage',
    '--no-first-run', '--no-default-browser-check', '--disable-background-networking',
    '--virtual-time-budget=12000', '--user-data-dir=' + profile, '--dump-dom', url,
  ];
  const {stdout, stderr} = await execFile(chrome, args, {timeout: 45000, maxBuffer: 10 * 1024 * 1024});
  if (!stdout.includes('data-result="pass"') ||
      !stdout.includes('Offline outbox native IndexedDB smoke: PASS')) {
    throw new Error('Browser smoke did not pass. DOM: ' + stdout.slice(-2200) +
      '\nBrowser stderr: ' + stderr.slice(-800));
  }
  process.stdout.write('Native Chromium IndexedDB persistence and replay smoke passed\n');
} catch (error) {
  process.stderr.write('Native IndexedDB browser smoke failed: ' + error.message + '\n');
  process.exitCode = 1;
} finally {
  await new Promise(resolve => server.close(resolve));
  if (profile) await rm(profile, {recursive: true, force: true});
}
