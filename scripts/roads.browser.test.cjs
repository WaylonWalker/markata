const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const asset = (name) => fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static', name), 'utf8');
const icon = '<span class="icon" data-icon="lucide/key-round"><svg viewBox="0 0 24 24"><path d="M2 12h20" stroke="currentColor"/></svg></span>';
const html = `<!doctype html><meta name="viewport" content="width=device-width, initial-scale=1">
<style>:root{--color-primary:teal;--color-border:gray;--color-text:black;--color-text-muted:gray;--color-surface:white;--color-background:white}
body{font:16px system-ui;margin:1rem} .road{max-width:70rem;margin-inline:auto} h3{font-size:1.5rem}</style>
<style>${asset('css/roads.css')}</style>
<div class="road editorial story reveal" id="story">
<h3 id="first">${icon} First milestone</h3><p><em>2026 Q3</em></p>
<p>A lead paragraph.</p><ul><li>One</li><li>Two</li><li>Three</li><li>Four</li></ul>
<ul class="tasks"><li><input type="checkbox"> Task</li></ul>
<div style="height:1200px"></div>
<h3 id="last">Last milestone</h3><p><em>2027 Q1</em></p><p>Last lead.</p></div>
<div class="road editorial" id="plain">
<h3>${icon} <a href="/title">Linked title</a> Existing label</h3><p><em>Keep this paragraph</em></p>
<h3>Plain title</h3><p><em>Not alone</em> extra text</p></div><div style="height:900px"></div>
<div class="road editorial" id="empty"><p>No entry headings.</p></div>
<div class="road timeline" id="no-list"><p>No milestone list.</p></div>`;

for (const width of [1440, 1024, 390]) {
  test(`editorial badges, labels, story lists and anchor reveal at ${width}px`, async () => {
    const browser = await chromium.launch({args:['--no-sandbox'], executablePath:process.env.CHROMIUM_EXECUTABLE || undefined});
    try {
      const page = await browser.newPage({viewport:{width, height:900}});
      const errors = [];
      page.on('pageerror', (error) => errors.push(error.message));
      await page.setContent(html);
      await page.addScriptTag({content:asset('js/roads.js')});
      await page.waitForTimeout(750);
      assert.equal(await page.locator('#first > .road-entry-label').textContent(), '2026 Q3');
      assert.equal(await page.locator('#first > .icon').count(), 0);
      assert.equal(await page.locator('#story .road-entry-icon > .icon').count(), 1);
      assert.equal(await page.locator('#plain .road-entry-label').first().textContent(), 'Existing label');
      assert.equal(await page.locator('#plain .road-entry-body p').first().textContent(), 'Keep this paragraph');
      assert.equal(await page.locator('#plain .road-entry-label').count(), 1);
      assert.equal(await page.locator('#story .road-entry-dot:not(.road-entry-icon)').count(), 1);
      const metrics = await page.evaluate(() => {
        const rail = document.querySelector('#story .road-rail').getBoundingClientRect();
        const badge = document.querySelector('#story .road-entry-icon').getBoundingClientRect();
        const list = document.querySelector('#story .road-entry-body > ul');
        return {rail:rail.x + rail.width/2, badge:badge.x + badge.width/2,
          columns:getComputedStyle(list).gridTemplateColumns.split(' ').length,
          tasks:getComputedStyle(document.querySelector('.tasks')).display,
          active:getComputedStyle(document.querySelector('.road-entry-icon')).color,
          muted:getComputedStyle(document.querySelector('#last').parentElement.querySelector('.road-entry-dot')).borderColor,
          overflow:document.documentElement.scrollWidth > innerWidth};
      });
      assert.ok(Math.abs(metrics.rail - metrics.badge) < 1, 'badge must be centered on rail');
      assert.equal(metrics.columns, width >= 1100 ? 2 : 1);
      assert.notEqual(metrics.tasks, 'grid', 'task lists must keep existing layout');
      assert.equal(metrics.active, 'rgb(0, 128, 128)');
      assert.equal(metrics.muted, 'rgb(85, 85, 85)');
      assert.equal(metrics.overflow, false);
      const laterBadge = page.locator('#plain .road-entry-icon');
      assert.equal(await laterBadge.evaluate((el) => getComputedStyle(el).color), 'rgb(128, 128, 128)');
      await laterBadge.evaluate((el) => el.scrollIntoView());
      await page.waitForTimeout(500);
      assert.equal(await laterBadge.evaluate((el) => getComputedStyle(el).color), 'rgb(0, 128, 128)');
      await page.evaluate(() => { location.hash = 'last'; });
      await page.waitForTimeout(100);
      assert.equal(await page.locator('#last').evaluate((el) => getComputedStyle(el.parentElement).opacity), '1');
      assert.equal(await page.locator('#last').evaluate((el) => getComputedStyle(el.parentElement).animationName), 'none');
      await page.emulateMedia({media:'print'});
      assert.equal(await page.locator('#first').evaluate((el) => getComputedStyle(el.parentElement).opacity), '1');
      assert.equal(await page.locator('#first').evaluate((el) => getComputedStyle(el.parentElement).animationName), 'none');
      assert.deepEqual(errors, []);
    } finally { await browser.close(); }
  });
}

test('reveal is progressive, motion-safe, idempotent, and optional', async () => {
  const browser = await chromium.launch({args:['--no-sandbox'], executablePath:process.env.CHROMIUM_EXECUTABLE || undefined});
  try {
    for (const mode of ['normal', 'reduced', 'no-observer', 'no-js', 'script-failure']) {
      const context = await browser.newContext({reducedMotion:mode === 'reduced' ? 'reduce' : 'no-preference', javaScriptEnabled:mode !== 'no-js'});
      const page = await context.newPage();
      await page.setContent(html);
      if (mode !== 'no-js' && mode !== 'script-failure') {
        if (mode === 'no-observer') await page.evaluate(() => { delete window.IntersectionObserver; });
        await page.addScriptTag({content:asset('js/roads.js')});
        await page.evaluate(() => window.MarkataRoads.init());
        await page.waitForTimeout(100);
        assert.equal(await page.locator('#story .road-entry').count(), 2);
        assert.equal(await page.locator('#plain .road-entry-revealed').count(), 0);
        const entry = page.locator('#story .road-entry').first();
        if (mode === 'normal') {
          for (const id of ['empty', 'no-list']) {
            const root = page.locator(`#${id}`);
            assert.equal(await root.getAttribute('data-road-ready'), '');
            assert.equal(await root.evaluate((el) => getComputedStyle(el).visibility), 'visible');
            await root.evaluate((el) => { el.innerHTML = '<h3>Later heading</h3>'; });
          }
          await page.evaluate(() => window.MarkataRoads.init());
          assert.equal(await page.locator('#empty .road-entry, #no-list .road-entry').count(), 0,
            'roots whose build returns null must not be reprocessed');
          assert.equal(await entry.evaluate((el) => getComputedStyle(el).animationName), 'road-entry-reveal');
          await entry.locator('h3').evaluate((el) => { const a = document.createElement('a'); a.href = '#last'; a.textContent = 'Focus'; el.append(a); a.focus(); });
          assert.equal(await entry.evaluate((el) => getComputedStyle(el).animationName), 'none');
          await page.locator('#last').evaluate((el) => el.scrollIntoView());
          await page.waitForTimeout(100);
          assert.equal(await page.locator('#last').evaluate((el) => getComputedStyle(el.parentElement).animationName), 'road-entry-reveal');
        } else {
          assert.equal(await page.locator('#story .road-entry-revealed').count(), 0);
          assert.equal(await entry.evaluate((el) => getComputedStyle(el).opacity), '1');
          if (mode === 'reduced') {
            assert.equal(await page.locator('.road-entry-icon').first().evaluate((el) => getComputedStyle(el).transitionDuration), '0s');
          }
        }
      } else {
        if (mode === 'script-failure') await page.waitForTimeout(2600);
        assert.equal(await page.locator('#story .road-entry').count(), 0);
        assert.equal(await page.locator('#story').evaluate((el) => getComputedStyle(el).visibility), 'visible');
        assert.equal(await page.locator('#first').textContent(), ' First milestone');
      }
      await context.close();
    }
  } finally { await browser.close(); }
});

test('fixed progress stays authored across scroll, resize, and view transitions', async () => {
  const browser = await chromium.launch({args:['--no-sandbox'], executablePath:process.env.CHROMIUM_EXECUTABLE || undefined});
  try {
    const page = await browser.newPage({viewport:{width:900, height:700}});
    const errors = [];
    page.on('pageerror', (error) => errors.push(error.message));
    await page.setContent(`<!doctype html><style>.road{margin-top:80px}</style>
      <div class="road timeline" data-road-progress-mode="static" data-road-progress="0" id="zero"><ul><li><a href="#a">A</a></li><li><a href="#b">B</a></li></ul></div>
      <div class="road timeline" data-road-progress-mode="static" data-road-progress="50" id="half"><ul><li><a href="#a">A</a></li><li><a href="#b">B</a></li></ul></div>
      <div class="road timeline" data-road-progress-mode="static" data-road-progress="100" id="full"><ul><li><a href="#a">A</a></li><li><a href="#b">B</a></li></ul></div>
      <div class="road git" data-road-progress-mode="milestones" data-road-current="1" id="git"><ul><li><a href="#main">main</a></li><li><a href="#feature">feature from main</a></li></ul></div>
      <div class="road editorial" data-road-progress-mode="milestones" data-road-current="1" id="editorial"><h3><a href="#first">First</a></h3><p>First body</p><h3><a href="#second">Second</a></h3><p>Second body</p></div>
      <div style="height:1200px"></div>`);
    await page.addScriptTag({content:asset('js/roads.js')});
    const state = async (selector) => page.locator(selector).evaluate((root) => ({
      active: [...root.querySelectorAll('.road-point, .road-entry')].map((point) => point.classList.contains('active')),
      progress: root.querySelector('.road-progress')?.style.strokeDashoffset || '',
      fill: root.querySelector('.road-fill')?.style.height || '',
    }));
    const initial = {
      zero: await state('#zero'), half: await state('#half'),
      full: await state('#full'), git: await state('#git'), editorial: await state('#editorial'),
    };
    assert.deepEqual(initial.zero.active, [false, false]);
    assert.deepEqual(initial.half.active, [true, false]);
    assert.deepEqual(initial.full.active, [true, true]);
    assert.deepEqual(initial.git.active, [true, false]);
    assert.deepEqual(initial.editorial.active, [true, false]);
    await page.evaluate(async () => {
      window.scrollTo(0, document.body.scrollHeight);
      await new Promise((resolve) => requestAnimationFrame(() => requestAnimationFrame(resolve)));
      window.dispatchEvent(new Event('view-transition-complete'));
      await new Promise((resolve) => requestAnimationFrame(() => requestAnimationFrame(resolve)));
    });
    await page.setViewportSize({width:390, height:700});
    assert.deepEqual(await state('#zero'), initial.zero);
    assert.deepEqual(await state('#half'), initial.half);
    assert.deepEqual(await state('#full'), initial.full);
    assert.deepEqual(await state('#git'), initial.git);
    assert.deepEqual(await state('#editorial'), initial.editorial);
    assert.deepEqual(errors, []);
  } finally { await browser.close(); }
});

test('fixed progress markup remains a readable list without JavaScript', async () => {
  const browser = await chromium.launch({args:['--no-sandbox'], executablePath:process.env.CHROMIUM_EXECUTABLE || undefined});
  try {
    const context = await browser.newContext({javaScriptEnabled:false});
    const page = await context.newPage();
    await page.setContent('<div class="road timeline" data-road-progress-mode="static" data-road-progress="65"><ul><li><a href="/one/">One</a></li><li><a href="/two/">Two</a></li></ul></div>');
    assert.deepEqual(Array.from(await page.locator('.road a').allTextContents()), ['One', 'Two']);
    assert.equal(await page.locator('.road a').evaluateAll((links) => links.every((link) => link.getAttribute('href'))), true);
    await context.close();
  } finally { await browser.close(); }
});
