// The visible navbar must keep its geometry while deferred search styles arrive.
const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');
const cssRoot = path.join(__dirname, '../pkg/themes/default/static/css');
const readCSS = name => fs.readFileSync(path.join(cssRoot, name + '.css'), 'utf8');
const coreCSS = ['variables', 'main', 'components', 'text-size', 'feeds'].map(readCSS).join('\n');
const searchCSS = readCSS('search');
const calendarCSS = readCSS('calendar-feed');

for (const width of [412, 1440]) {
  test(`navbar search reserves layout at ${width}px`, async () => {
    const browser = await chromium.launch({headless: true, executablePath: process.env.CHROMIUM_EXECUTABLE || undefined, args: ['--no-sandbox']});
    try {
      const page = await browser.newPage({viewport: {width, height: 900}});
      await page.setContent(`<!doctype html><html data-text-size="large"><head><style>${coreCSS}</style></head>
        <body><header class="site-header"><div class="container">
        <a class="site-title" href="/">Waylon Walker</a>
        <nav class="site-nav site-nav--header site-nav--horizontal"><a href="/archive/">Archive</a></nav>
        <div class="header-controls"><div class="search-container search-container--navbar">
        <div id="pagefind-search"></div></div><div class="header-action-controls"><button>Theme</button></div></div>
        </div></header><div class="page-wrapper page-wrapper--feed-only"><aside class="feed-sidebar feed-sidebar--left"><nav class="feed-nav">Series</nav></aside><main class="main-content" id="main-content"><div class="feed h-feed calendar-feed feed--editorial"><header class="feed-header"><h1>Archive</h1></header><div id="stories">Stories</div></div></main></div></body></html>`);
      const geometry = () => page.evaluate(() => Object.fromEntries(
        ['.site-header', '.site-nav', '.header-controls', '.search-container', 'main', '.feed', '#stories'].map(selector => {
          const r = document.querySelector(selector).getBoundingClientRect();
          return [selector, {x:r.x, y:r.y, width:r.width, height: selector === '.site-header' ? r.height : undefined}];
        })
      ));
      const before = await geometry();
      await page.addStyleTag({content: calendarCSS});
      assert.deepEqual(await geometry(), before, 'calendar CSS must not shift the header or content');
      for (const markup of [
        '<div class="bleve-search"><input class="bleve-search__input" aria-label="Search"></div>',
        '<div class="pagefind-ui"><form class="pagefind-ui__form"><input class="pagefind-ui__search-input" aria-label="Search"></form></div>',
      ]) {
        await page.locator('#pagefind-search').evaluate((e, html) => {e.innerHTML = html;}, markup);
        assert.deepEqual(await geometry(), before, 'search hydration before CSS must not shift content');
      }
      await page.addStyleTag({content: searchCSS});
      assert.deepEqual(await geometry(), before, 'deferred search CSS must not shift content');
    } finally {
      await browser.close();
    }
  });
}

const conditionalJS = fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/js/conditional-css.js'), 'utf8');
for (const width of [412, 1440]) for (const provider of ['loadPagefind', 'initPagefindSearch']) {
  const input = '<div class="pagefind-ui"><form class="pagefind-ui__form"><input class="pagefind-ui__search-input" aria-label="Search"></form></div>';
  test(`Pagefind ${provider} initializes without forced navbar resize at ${width}px`, async () => {
    const browser = await chromium.launch({headless:true, executablePath:process.env.CHROMIUM_EXECUTABLE || undefined, args:['--no-sandbox']});
    try {
      const page = await browser.newPage({viewport:{width,height:900}});
      await page.setContent(`<!doctype html><html><head><style>${coreCSS}</style>
        <link rel="stylesheet" href="data:text/css," data-css-sources="/feeds.fixture.css /cards.fixture.css">
        <script>
          window.layoutReads=0;window.resizes=[];window.searchCalls=0;window.fallbackCalls=0;
          if('${provider}'==='loadPagefind')window.initPagefindSearch=()=>window.fallbackCalls++;
          const widthGetter=Object.getOwnPropertyDescriptor(HTMLElement.prototype,'offsetWidth');
          Object.defineProperty(HTMLElement.prototype,'offsetWidth',{configurable:true,get(){
            if(this.matches('.search-container--navbar'))window.layoutReads++;
            return widthGetter.get.call(this);
          }});
          addEventListener('resize',e=>window.resizes.push(e.isTrusted));
          window.${provider}=()=>{window.searchCalls++;document.querySelector('#pagefind-search').innerHTML=
            ${JSON.stringify(input)}};
        </script></head><body>
        <header class="site-header"><div class="container"><a class="site-title" href="/">Waylon Walker</a>
        <nav class="site-nav site-nav--header site-nav--horizontal"><a href="/archive/">Archive</a></nav>
        <div class="header-controls"><div class="search-container search-container--navbar"><div id="pagefind-search"></div></div>
        <div class="header-action-controls"><button>Theme</button></div></div></div></header>
        <main><div class="feed h-feed">Archive</div></main><script>${conditionalJS}</script></body></html>`);
      const settle = () => page.evaluate(() => new Promise(resolve => requestAnimationFrame(() => requestAnimationFrame(resolve))));
      await settle();
      assert.deepEqual(await page.evaluate(() => ({reads:layoutReads,resizes,searchCalls,fallbackCalls})), {reads:0,resizes:[],searchCalls:0,fallbackCalls:0}, 'cold startup must leave lazy search unloaded');
      const before = await page.locator('.search-container').boundingBox();
      await page.addStyleTag({content:searchCSS});
      await page.evaluate(() => dispatchEvent(new Event('view-transition-complete')));
      await settle();
      assert.deepEqual(await page.evaluate(() => ({reads:layoutReads,resizes,searchCalls,fallbackCalls})), {reads:0,resizes:[],searchCalls:0,fallbackCalls:0}, 'conditional CSS transition handling must preserve lazy search');
      await page.evaluate(provider => window[provider](), provider);
      await settle();
      assert.deepEqual(await page.evaluate(() => ({reads:layoutReads,resizes,searchCalls,fallbackCalls})), {reads:0,resizes:[],searchCalls:1,fallbackCalls:0}, 'SPA navigation must initialize search once');
      assert.deepEqual(await page.locator('.search-container').boundingBox(), before, 'SPA search hydration must keep its dimensions');
      await page.setViewportSize({width:width+24,height:900});
      await page.waitForFunction(() => window.resizes.some(trusted => trusted));
      assert.ok((await page.locator('.search-container input').boundingBox()).width <= (await page.locator('.search-container').boundingBox()).width + 1);
    } finally {
      await browser.close();
    }
  });
}
