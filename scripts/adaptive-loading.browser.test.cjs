// Real browser timing integration. Uses the Playwright installation bundled with
// agent-browser unless PLAYWRIGHT_MODULE points at another playwright(-core).
const {test} = require('node:test');
const assert = require('node:assert/strict');
const http = require('node:http');
const fs = require('node:fs');
const path = require('node:path');
const {execFileSync} = require('node:child_process');
const playwright = require(process.env.PLAYWRIGHT_MODULE || path.join(execFileSync('npm', ['root', '-g'], {encoding:'utf8'}).trim(), 'agent-browser/node_modules/playwright-core'));
const controller = fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/js/adaptive-loading.js'));
// Loading controls now live in a shared partial, included either in the
// appearance dialog or directly in the header when the picker is disabled.
const control = fs.readFileSync(path.join(__dirname, '../pkg/themes/default/templates/partials/adaptive-loading.html'), 'utf8');
const svg = Buffer.from('<svg xmlns="http://www.w3.org/2000/svg" width="160" height="90"><rect width="160" height="90" fill="blue"/><!--' + 'x'.repeat(80000) + '--></svg>');

function serve() {
  return http.createServer((req, res) => {
    if (req.url === '/clip.webm') {res.setHeader('Content-Type','video/webm'); return res.end(fs.readFileSync(path.join(__dirname,'fixtures/adaptive-video.webm')));}
    if (req.url === '/adaptive.js') return res.end(controller);
    if (req.url === '/view.js') return res.end(fs.readFileSync(path.join(__dirname,'../pkg/themes/default/static/js/view-transitions.js')));
    if (req.url === '/startup.js') return res.end('window.controllerReadyState = document.readyState;');
    if (req.url.startsWith('/sample/')) {
      res.setHeader('Content-Type','image/svg+xml');
      res.setHeader('Cache-Control', req.url.includes('cached') ? 'max-age=3600' : 'no-store');
      const poor = req.url.includes('poor');
      let offset = 0;
      const timer = setInterval(() => {
        res.write(svg.subarray(offset, offset + 8000));
        offset += 8000;
        if (offset >= svg.length) {clearInterval(timer); res.end();}
      }, poor ? 90 : 5);
      res.on('close', () => clearInterval(timer));
      return;
    }
    res.setHeader('Content-Type','text/html');
    res.end('<!doctype html><html><body>' + control + '<div id="view-transition-page"><main><p>Useful content</p>' +
      '<video width="160" height="90" muted controls preload="none" data-authored-autoplay="true"><source src="/clip.webm" type="video/webm"></video>' +
      '<video width="160" height="90" muted controls preload="none" id="no-intent"><source src="/clip.webm" type="video/webm"></video>' +
      [0,1,2,3].map(i => '<img width="160" height="90" src="/sample/' + (req.url.includes('good') ? 'good' : 'poor') + i + '">').join('') +
      '<a id="next-route" href="/good-native-next">Next route</a></main></div><script src="/startup.js" defer></script><script src="/adaptive.js" defer></script><script src="/view.js" defer></script></body></html>');
  });
}
async function more(page, kind, count) {
  await page.evaluate(({kind,count}) => Promise.all(Array.from({length:count}, (_,i) => new Promise(resolve => {
    const img = new Image(); img.onload = resolve; img.onerror = resolve;
    img.src = '/sample/' + kind + '-' + (kind === 'good-cached' ? '' : Date.now() + '-') + i;
    document.querySelector('main').appendChild(img);
  }))), {kind,count});
}
async function recommendation(page, value) {
  await page.waitForFunction(value => document.documentElement.dataset.loadingRecommendation === value, value);
}

for (const engine of (process.env.ADAPTIVE_BROWSERS || 'chromium,firefox,brave').split(',')) {
  test(engine + ': deferred startup, actual timing, recovery, manual modes and UI', {timeout:90000}, async () => {
    const server = serve(); await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
    const browser = await (engine === 'firefox' ? playwright.firefox : playwright.chromium).launch({headless:true,
      ...(engine === 'firefox' ? (process.env.ADAPTIVE_FIREFOX_EXECUTABLE ? {executablePath:process.env.ADAPTIVE_FIREFOX_EXECUTABLE} : {}) : {executablePath:engine === 'brave' ? '/usr/bin/brave' : '/usr/bin/chromium', args:['--no-sandbox']})}).catch(async error => {await new Promise(resolve => server.close(resolve)); throw error;});
    try {
      const context = await browser.newContext({viewport:{width:390,height:844}});
      await context.addInitScript(() => Object.defineProperty(navigator,'connection',{value:undefined, configurable:true}));
      const page = await context.newPage();
      const origin = 'http://127.0.0.1:' + server.address().port;
      await page.goto(origin + '/poor', {waitUntil:'load'});
      assert.equal(await page.evaluate(() => controllerReadyState), 'interactive');
      await recommendation(page, 'constrained');
      let evidence = await page.evaluate(() => JSON.parse(sessionStorage.getItem('markata-loading-auto-v1')));
      assert.ok(evidence.poor >= 4);
      assert.equal(await page.evaluate(() => localStorage.length), 0);
      assert.equal(await page.getAttribute('html','data-loading-policy'),'constrained');
      assert.equal(await page.locator('video').first().getAttribute('autoplay'),null);
      // A cache hit is not evidence that the network improved. Preload via
      // fetch (not a measured media initiator), then request cached images.
      await page.evaluate(async () => {
        for (let i=0;i<4;i++) await (await fetch('/sample/good-cached-'+i)).arrayBuffer();
      });
      await more(page,'good-cached',4);
      assert.equal(await page.getAttribute('html','data-loading-policy'),'constrained');
      const geometry = await page.locator('main p').boundingBox();
      const videoGeometry = await page.locator('video').first().boundingBox();
      const sources = await page.locator('main img').evaluateAll(imgs => imgs.map(img => img.currentSrc));
      // Native controls remain usable while Auto is constrained.
      await page.locator('#no-intent').evaluate(v => v.play());
      await page.waitForFunction(() => !document.querySelector('#no-intent').paused);
      await more(page,'good',1);
      assert.equal(await page.getAttribute('html','data-loading-policy'),'constrained');
      await more(page,'good',3);
      await recommendation(page,'full-quality');
      assert.equal(await page.locator('video').first().getAttribute('data-autoplay-started'),null);
      assert.equal(await page.locator('#no-intent').getAttribute('autoplay'),null);
      assert.deepEqual((await page.locator('main img').evaluateAll(imgs => imgs.map(img => img.currentSrc))).slice(0,sources.length), sources);
      assert.deepEqual(await page.locator('main p').boundingBox(), geometry);
      await page.waitForFunction(() => !document.querySelector('#no-intent').paused && document.querySelector('#no-intent').currentTime > 0);
      assert.deepEqual(await page.locator('video').first().boundingBox(), videoGeometry);
      const playback = await page.locator('#no-intent').evaluate(v => v.currentTime);
      // Offline/online cannot restore the old healthy recommendation.
      await context.setOffline(true); await recommendation(page,'constrained');
      assert.ok(await page.locator('#no-intent').evaluate((v,t) => !v.paused && v.currentTime >= t, playback));
      await context.setOffline(false);
      assert.equal(await page.getAttribute('html','data-loading-policy'),'constrained');
      await more(page,'good',4);
      await recommendation(page,'full-quality');
      await more(page,'poor',2);
      assert.equal(await page.getAttribute('html','data-loading-policy'),'full-quality');
      assert.equal(await page.evaluate(() => JSON.parse(sessionStorage.getItem('markata-loading-auto-v1')).samples),12);
      await page.goto(origin + '/good-again');
      assert.equal(await page.locator('video').first().getAttribute('data-autoplay-started'),'true');
      await page.waitForFunction(() => !document.querySelector('video').paused);
      // Exhausting one page's budget cannot freeze session confidence forever.
      await page.goto(origin + '/poor-again');
      await recommendation(page,'constrained');
      assert.ok(await page.evaluate(() => JSON.parse(sessionStorage.getItem('markata-loading-auto-v1')).poor >= 4));
      await page.close(); await context.close();

      for (const mode of ['save-data','full-quality']) {
        const ctx = await browser.newContext();
        await ctx.addInitScript(mode => {Object.defineProperty(navigator,'connection',{value:undefined}); localStorage.setItem('markata-loading-mode-v1',mode);}, mode);
        const p = await ctx.newPage();
        await p.goto(origin + (mode === 'save-data' ? '/good' : '/poor'));
        if (mode === 'save-data') await more(p,'good',1);
        await p.waitForFunction(() => !document.querySelector('[data-loading-actions]').hidden);
        assert.equal(await p.getAttribute('html','data-loading-policy'), mode === 'save-data' ? 'constrained' : 'full-quality');
        await p.locator('[data-adaptive-loading] summary').click();
        await p.locator('[data-loading-keep]').click();
        assert.equal(await p.evaluate(() => localStorage.getItem('markata-loading-mode-v1')),mode);
        assert.ok(await p.locator('[data-loading-actions]').isHidden());
        await p.reload();
        assert.ok(await p.locator('[data-loading-actions]').isHidden());
        // Choosing Auto via radio is always available after dismissing.
        await p.locator('[data-adaptive-loading] summary').click();
        await p.locator('input[value="auto"]').check();
        assert.equal(await p.getAttribute('html','data-loading-mode'),'auto');
        await ctx.close();
      }
      const intent = await browser.newContext();
      await intent.addInitScript(() => {
        Object.defineProperty(navigator,'connection',{value:{saveData:true}});
        localStorage.setItem('markata-loading-mode-v1','save-data');
      });
      const intentPage = await intent.newPage();
      await intentPage.goto(origin + '/good'); await more(intentPage,'good',1);
      assert.equal(await intentPage.getAttribute('html','data-loading-recommendation'),'constrained');
      assert.ok(await intentPage.locator('[data-loading-actions]').isHidden());
      await intent.close();

      if (engine === 'chromium') {
        for (const mode of ['auto','save-data','full-quality']) {
          const nav = await browser.newContext();
          await nav.addInitScript(mode => {
            Object.defineProperty(navigator,'connection',{value:undefined});
            if (mode !== 'auto') localStorage.setItem('markata-loading-mode-v1',mode);
            window.routeSwaps = 0;
            window.entryRecommendations = [];
            window.addEventListener('view-transition-complete', () => {
              entryRecommendations.push(document.documentElement.dataset.loadingRecommendation);
              routeSwaps++;
            });
          },mode);
          const n = await nav.newPage(); await n.goto(origin + '/poor-native');
          await n.locator('#next-route').click();
          await n.waitForFunction(() => routeSwaps > 0);
          await recommendation(n,'full-quality');
          assert.equal(await n.getAttribute('html','data-loading-mode'),mode);
          assert.equal(await n.getAttribute('html','data-loading-policy'),mode === 'save-data' ? 'constrained' : 'full-quality');
          // Auto learns fast on this destination but keeps it still. The next
          // user-requested route may honor that already-measured confidence.
          if (mode === 'auto') {
            const atEntry = await n.evaluate(() => entryRecommendations[0]);
            assert.equal(await n.locator('video').first().getAttribute('data-autoplay-started'),atEntry === 'full-quality' ? 'true' : null);
            await n.evaluate(() => window.navigateWithViewTransition('/good-native-last'));
            await n.waitForFunction(() => routeSwaps >= 2);
            assert.equal(await n.locator('video').first().getAttribute('data-autoplay-started'),'true');
          }
          await nav.close();
        }
      }
      // With the observer unavailable, only the deferred load listener can
      // collect these real Resource Timing entries.
      for (const connection of [undefined, {effectiveType:'4g', downlink:10}, {saveData:true}]) {
        const fallback = await browser.newContext();
        await fallback.addInitScript(connection => {
          Object.defineProperty(navigator,'connection',{value:connection});
          window.PerformanceObserver = undefined;
        }, connection);
        const f = await fallback.newPage();
        await f.goto(origin + (connection && connection.saveData ? '/good' : '/poor'));
        assert.equal(await f.evaluate(() => controllerReadyState), 'interactive');
        assert.equal(await f.getAttribute('html','data-loading-mode'),'auto');
        assert.equal(await f.getAttribute('html','data-loading-policy'),'constrained');
        const saved = await f.evaluate(() => JSON.parse(sessionStorage.getItem('markata-loading-auto-v1')));
        assert.ok(saved.samples >= 4);
        await fallback.close();
      }
      const ctx = await browser.newContext(); const p = await ctx.newPage();
      await ctx.addInitScript(() => {Object.defineProperty(navigator,'connection',{value:undefined}); localStorage.setItem('markata-loading-mode-v1','full-quality');});
      await p.goto(origin + '/poor');
      await p.locator('[data-adaptive-loading] summary').click();
      await p.locator('[data-loading-use-auto]').click();
      assert.equal(await p.getAttribute('html','data-loading-policy'),'constrained');
      assert.equal(await p.evaluate(() => localStorage.getItem('markata-loading-mode-v1')),null);
      await ctx.close();
    } finally {await browser.close(); await new Promise(resolve => server.close(resolve));}
  });
}
