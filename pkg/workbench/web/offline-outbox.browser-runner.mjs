// Native IndexedDB smoke via Chromium DevTools Protocol, no npm dependencies.
// Waits for the page's real asynchronous completion, not --dump-dom timing.
import {createServer} from 'node:http';
import {readFile, mkdtemp, rm} from 'node:fs/promises';
import {tmpdir} from 'node:os';
import {join, dirname} from 'node:path';
import {fileURLToPath} from 'node:url';
import {spawn, execFile as execFileCallback} from 'node:child_process';
import {promisify} from 'node:util';

const execFile = promisify(execFileCallback);
const pause = milliseconds => new Promise(resolve => setTimeout(resolve, milliseconds));
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
    response.writeHead(200, {'Content-Type': asset[1], 'Cache-Control': 'no-store'});
    response.end(await readFile(join(root, asset[0])));
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
    } catch { /* Try another Chrome/Chromium installation. */ }
  }
  throw new Error('Chromium/Chrome is not installed; native IndexedDB test cannot run');
}

async function debuggerTab(profile, expectedURL) {
  for (let attempt = 0; attempt < 140; attempt++) {
    try {
      const portFile = await readFile(join(profile, 'DevToolsActivePort'), 'utf8');
      const port = Number(portFile.split('\n')[0]);
      if (!Number.isSafeInteger(port) || port <= 0) throw new Error('Invalid Chrome debug port');
      const response = await fetch('http://127.0.0.1:' + port + '/json/list', {signal: AbortSignal.timeout(1500)});
      const pages = await response.json();
      const page = pages.find(item => item.type === 'page' && item.url === expectedURL);
      if (page?.webSocketDebuggerUrl) return page.webSocketDebuggerUrl;
    } catch { /* Browser may still be starting. */ }
    await pause(100);
  }
  throw new Error('Chrome did not open the local smoke page within 14 seconds');
}

async function connectDebugger(url) {
  const socket = new WebSocket(url);
  await new Promise((resolve, reject) => {
    socket.addEventListener('open', resolve, {once: true});
    socket.addEventListener('error', reject, {once: true});
  });
  let nextID = 0;
  const pending = new Map();
  socket.addEventListener('message', event => {
    const message = JSON.parse(event.data);
    const waiter = pending.get(message.id);
    if (!waiter) return;
    pending.delete(message.id);
    waiter(message);
  });
  return {
    socket,
    command(method, params) {
      return new Promise((resolve, reject) => {
        const id = ++nextID;
        pending.set(id, message => message.error ? reject(new Error(message.error.message)) : resolve(message.result));
        socket.send(JSON.stringify({id, method, params}));
      });
    },
  };
}

async function awaitBrowserCompletion(debuggerClient) {
  for (let attempt = 0; attempt < 90; attempt++) {
    const result = await debuggerClient.command('Runtime.evaluate', {
      expression: '(() => { const x = document.querySelector("#result"); return {status: x?.dataset.result || "pending", message: x?.textContent || document.body?.innerText?.slice(0, 400) || ""}; })()',
      returnByValue: true,
    });
    const value = result.result?.value;
    if (value?.status === 'pass') return;
    if (value?.status === 'fail') throw new Error('Native IndexedDB browser failure: ' + value.message);
    await pause(200);
  }
  const result = await debuggerClient.command('Runtime.evaluate', {
    expression: 'document.body?.innerText?.slice(0, 500) || "No body content"',
    returnByValue: true,
  });
  throw new Error('Browser never completed its IndexedDB operations: ' + result.result?.value);
}

let profile;
let chromeProcess;
let debuggerClient;
try {
  await new Promise((resolve, reject) => {
    server.once('error', reject);
    server.listen(0, '127.0.0.1', resolve);
  });
  profile = await mkdtemp(join(tmpdir(), 'markata-outbox-smoke-'));
  const chrome = await resolveChrome();
  const pageURL = 'http://127.0.0.1:' + server.address().port + '/offline-outbox.browser.html';
  chromeProcess = spawn(chrome, [
    '--headless=new', '--no-sandbox', '--disable-gpu', '--disable-dev-shm-usage',
    '--no-first-run', '--no-default-browser-check', '--no-proxy-server',
    '--remote-debugging-port=0', '--remote-allow-origins=*',
    '--user-data-dir=' + profile, pageURL,
  ], {stdio: 'ignore'});
  const socketURL = await debuggerTab(profile, pageURL);
  debuggerClient = await connectDebugger(socketURL);
  await awaitBrowserCompletion(debuggerClient);
  process.stdout.write('Native Chromium IndexedDB persistence and replay smoke passed\n');
} catch (error) {
  process.stderr.write('Native IndexedDB browser smoke failed: ' + error.message + '\n');
  process.exitCode = 1;
} finally {
  debuggerClient?.socket.close();
  chromeProcess?.kill();
  await pause(200);
  await new Promise(resolve => server.close(resolve));
  if (profile) await rm(profile, {recursive: true, force: true});
}
