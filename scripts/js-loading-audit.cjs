#!/usr/bin/env node
'use strict';

// Audit JavaScript actually requested by pages, not total deployed asset size.
// Resource Timing omits/masks byte sizes when cross-origin TAO permission is absent.
const fs = require('node:fs');
const path = require('node:path');

function parseArgs(argv) {
  const opts = {urls: [], out: 'test-results/js-loading-audit.json', waitMs: 500};
  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i];
    if (arg === '--url' && argv[i + 1]) opts.urls.push(argv[++i]);
    else if (arg === '--out' && argv[i + 1]) opts.out = argv[++i];
    else if (arg === '--wait-ms' && argv[i + 1]) opts.waitMs = Number(argv[++i]);
    else if (arg === '--help' || arg === '-h') opts.help = true;
    else throw new Error('Unknown or incomplete option: ' + arg);
  }
  if (opts.help) return opts;
  if (!opts.urls.length) throw new Error('Provide at least one --url');
  if (!Number.isInteger(opts.waitMs) || opts.waitMs < 0 || opts.waitMs > 30000) {
    throw new Error('--wait-ms must be an integer between 0 and 30000');
  }
  opts.urls = opts.urls.map((item) => {
    const url = new URL(item);
    if (!['http:', 'https:'].includes(url.protocol)) throw new Error('Unsupported protocol: ' + item);
    return url.href;
  });
  return opts;
}

function isScriptResource(item) {
  if (!item || !item.name) return false;
  try { return /\.m?js$/.test(new URL(item.name).pathname); }
  catch (_) { return false; }
}

function summarize(snapshot) {
  const origin = new URL(snapshot.url).origin;
  const scripts = (snapshot.resources || []).filter(isScriptResource).map((item) => {
    const sameOrigin = new URL(item.name).origin === origin;
    const transferBytes = Number(item.transferSize) || 0;
    const decodedBytes = Number(item.decodedBodySize) || 0;
    const encodedBytes = Number(item.encodedBodySize) || 0;
    return {
      url: item.name,
      sameOrigin,
      initiatorType: item.initiatorType || '',
      transferBytes,
      encodedBytes,
      decodedBytes,
      durationMs: Math.round(Number(item.duration) * 10) / 10 || 0,
      // A zero transfer could mean cache reuse OR cross-origin timing restrictions.
      transferKnown: transferBytes > 0,
    };
  });
  const sum = (key) => scripts.reduce((total, item) => total + item[key], 0);
  return {
    url: snapshot.url,
    title: snapshot.title || '',
    scriptRequests: scripts.length,
    uniqueScriptUrls: new Set(scripts.map((item) => item.url)).size,
    sameOriginScriptRequests: scripts.filter((item) => item.sameOrigin).length,
    crossOriginScriptRequests: scripts.filter((item) => !item.sameOrigin).length,
    observedTransferBytes: sum('transferBytes'),
    observedEncodedBytes: sum('encodedBytes'),
    observedDecodedBytes: sum('decodedBytes'),
    zeroTransferEntries: scripts.filter((item) => !item.transferKnown).length,
    inlineExecutableScriptCount: snapshot.inlineExecutableScriptCount || 0,
    inlineExecutableScriptCharacters: snapshot.inlineExecutableScriptCharacters || 0,
    scripts,
  };
}

function playwrightChromium() {
  const candidates = [process.env.PLAYWRIGHT_MODULE, 'playwright', 'playwright-core'].filter(Boolean);
  try {
    const globalRoot = require('node:child_process').execFileSync('npm', ['root', '-g'], {encoding: 'utf8'}).trim();
    candidates.push(path.join(globalRoot, 'agent-browser/node_modules/playwright-core'));
  } catch (_) {}
  for (const candidate of candidates) {
    try { return require(candidate).chromium; }
    catch (error) { if (error.code !== 'MODULE_NOT_FOUND') throw error; }
  }
  throw new Error('Playwright unavailable. Install playwright or set PLAYWRIGHT_MODULE.');
}

async function measure(page, url, waitMs) {
  await page.goto(url, {waitUntil: 'load', timeout: 45000});
  if (waitMs) await page.waitForTimeout(waitMs);
  const snapshot = await page.evaluate(() => {
    const inline = Array.from(document.scripts)
      .filter((script) => !script.src && (!script.type || script.type === 'module' || /javascript/.test(script.type)));
    return {
      url: location.href,
      title: document.title,
      resources: performance.getEntriesByType('resource').map((r) => ({
        name: r.name,
        initiatorType: r.initiatorType,
        transferSize: r.transferSize,
        encodedBodySize: r.encodedBodySize,
        decodedBodySize: r.decodedBodySize,
        duration: r.duration,
      })),
      inlineExecutableScriptCount: inline.length,
      inlineExecutableScriptCharacters: inline.reduce((sum, script) => sum + (script.textContent || '').length, 0),
    };
  });
  return summarize(snapshot);
}

async function run(opts) {
  const chromium = playwrightChromium();
  const launchOpts = {headless: true};
  if (process.env.CHROMIUM_EXECUTABLE) launchOpts.executablePath = process.env.CHROMIUM_EXECUTABLE;
  const browser = await chromium.launch(launchOpts);
  const reports = [];
  try {
    for (const url of opts.urls) {
      // Fresh context for cold navigation; same context on repeat navigation.
      const context = await browser.newContext();
      try {
        const page = await context.newPage();
        const cold = await measure(page, url, opts.waitMs);
        const repeat = await measure(page, url, opts.waitMs);
        reports.push({requestedUrl: url, cold, repeat});
      } finally { await context.close(); }
    }
  } finally { await browser.close(); }

  const report = {
    schemaVersion: 1,
    generatedAt: new Date().toISOString(),
    notes: [
      'Resource Timing bytes are observations, not complete transfer accounting.',
      'Zero transferSize may mean cache reuse or blocked cross-origin timing; it is not classified as a cache hit.',
      'Inline script characters are not compressed bytes or execution timings.',
      'Repeat navigation is a full-page navigation, not a View Transition SPA click.',
      'Interaction-driven JS is excluded unless it loads in the settle window.',
    ],
    pages: reports,
  };
  fs.mkdirSync(path.dirname(opts.out), {recursive: true});
  fs.writeFileSync(opts.out, JSON.stringify(report, null, 2) + '\n');
  for (const entry of reports) {
    const kb = (num) => (num / 1000).toFixed(1);
    process.stdout.write(entry.requestedUrl + '\n  cold: ' +
      entry.cold.scriptRequests + ' requests, ' + kb(entry.cold.observedTransferBytes) +
      ' KB observed transfer, ' + entry.cold.zeroTransferEntries + ' zero-transfer entries\n  repeat: ' +
      entry.repeat.scriptRequests + ' requests, ' + kb(entry.repeat.observedTransferBytes) +
      ' KB observed transfer, ' + entry.repeat.zeroTransferEntries + ' zero-transfer entries\n');
  }
  process.stdout.write('Report: ' + opts.out + '\n');
}

if (require.main === module) {
  let opts;
  try { opts = parseArgs(process.argv.slice(2)); }
  catch (error) {
    process.stderr.write(error.message + '\nTry: node scripts/js-loading-audit.cjs --url http://localhost:8000/post/\n');
    process.exitCode = 2;
  }
  if (opts && opts.help) {
    process.stdout.write('Usage: node scripts/js-loading-audit.cjs --url URL [--url URL...] [--out FILE] [--wait-ms 500]\nSet PLAYWRIGHT_MODULE / CHROMIUM_EXECUTABLE if needed.\n');
  } else if (opts) {
    run(opts).catch((error) => { process.stderr.write(String(error.stack || error) + '\n'); process.exitCode = 1; });
  }
}

module.exports = {parseArgs, isScriptResource, summarize};
