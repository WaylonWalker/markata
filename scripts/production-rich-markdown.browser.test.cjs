// Public, read-only acceptance check for the deployed rich Markdown release.
// Does not access Builder Admin, require credentials, or assert an unknown deploy revision.
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require('playwright');

const base = process.env.LIVE_SAMPLE_URL || 'https://waylonwalker.com/sample/';
const output = process.env.SCREENSHOT_DIR || 'test-results/live-rich-markdown';
const widths = [320, 360, 412, 768, 1440];
const themes = [
  {name:'light', mode:'light', palette:'ayu-light'},
  {name:'dark', mode:'dark', palette:'ayu-dark'},
  {name:'custom', mode:'dark', palette:'nord-dark'},
];
const report = {url:base, timestamp:new Date().toISOString(), checks:[], failures:[], observations:[]};
fs.mkdirSync(output,{recursive:true});

function check(condition, message) {
  report.checks.push({message,pass:Boolean(condition)});
  if (!condition) throw new Error(message);
}

async function snapshot(page) {
  return page.evaluate(() => {
    const all = selector => [...document.querySelectorAll(selector)];
    const diagrams = all('pre.mermaid');
    const svgs = all('pre.mermaid svg');
    const glossary = all('a.glossary-term').map(a => ({
      name:a.textContent.trim(), href:a.getAttribute('href'),
      preview:a.getAttribute('data-local-preview'),
    }));
    const text = document.querySelector('article')?.innerText || '';
    const edge = document.querySelector('pre.mermaid svg path.flowchart-link');
    const resolveColor = input => {
      if (!input || input === 'none') return null;
      const probe = document.createElement('span');
      probe.style.color = input;
      document.body.appendChild(probe);
      const value = getComputedStyle(probe).color;
      probe.remove();
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');
      ctx.fillStyle = value;
      ctx.fillRect(0,0,1,1);
      return [...ctx.getImageData(0,0,1,1).data].slice(0,3).join(',');
    };
    return {
      width:innerWidth, scrollWidth:document.documentElement.scrollWidth,
      theme:document.documentElement.dataset.theme,
      palette:document.documentElement.dataset.palette,
      diagrams:diagrams.length, rendered:svgs.length,
      errors:all('pre.mermaid svg .error-icon, pre.mermaid svg .error-text').length,
      arrows:all('pre.mermaid svg path[marker-end], pre.mermaid svg path[marker-start]').length,
      clusters:all('pre.mermaid svg .cluster').length,
      glossary, emoji:['😂','🚀','👍','😀','🔥'].map(e=>text.includes(e)),
      copyButtons:all('button.md-code-copy').length,
      nested:all('.sample-nested').length,
      tables:all('article table').length,
      edge:edge ? resolveColor(getComputedStyle(edge).stroke) : null,
      expectedText:resolveColor('var(--color-text)'),
    };
  });
}

async function main() {
  const browser = await chromium.launch({headless:true,args:['--no-sandbox']});
  let context, noJS;
  try {
    context = await browser.newContext({
      viewport:{width:1440,height:1000}, colorScheme:'dark',
      reducedMotion:'reduce',
    });
    const page = await context.newPage();
    const errors = [];
    page.on('pageerror',error => errors.push(error.message));
    const url = new URL(base);
    url.searchParams.set('rich-markdown-audit',String(Date.now()));
    const response = await page.goto(url.toString(),{waitUntil:'domcontentloaded',timeout:90000});
    check(response && response.status() === 200,'Public /sample/ returns HTTP 200');
    await page.locator('pre.mermaid svg').first().waitFor({timeout:60000});
    await page.waitForFunction(() => {
      const source=document.querySelectorAll('pre.mermaid').length;
      return source >= 12 && document.querySelectorAll('pre.mermaid svg').length === source;
    },null,{timeout:90000});

    let state = await snapshot(page);
    report.initial = state;
    check(state.diagrams===12 && state.rendered===12,'All 12 Mermaid examples render SVGs');
    check(state.errors===0,'No Mermaid syntax-error SVGs');
    check(state.arrows>=20,'Mermaid diagrams contain directional arrows');
    check(state.glossary.length>=3,'Server-rendered glossary links present');
    check(state.glossary.every(a => {const href=a.href||''; return (href.startsWith('/') && !href.includes('<') && !href.includes('&lt;')) || href.startsWith('https://');}),'Glossary hrefs are valid local or HTTPS paths');
    check(state.glossary.every(a => !a.preview || (()=>{try{JSON.parse(a.preview);return true}catch{return false}})()),'Glossary local previews contain valid JSON');
    check(state.emoji.every(Boolean),'Five supported emoji shortcodes render as emoji');
    check(state.nested>=1 && state.tables>=2,'Nested containers and semantic tables are present');
    check(state.copyButtons>=1,'Code-copy controls are present');

    const button=page.locator('button.md-code-copy').first();
    await button.focus();
    report.copyGeometry=await page.evaluate(()=>{
      const button=document.querySelector('button.md-code-copy');
      const pre=button?.closest('pre');
      const box=button?.getBoundingClientRect();
      const parent=pre?.getBoundingClientRect();
      const style=button&&getComputedStyle(button);
      return {
        position:style?.position, insetInlineEnd:style?.insetInlineEnd,
        insetBlockStart:style?.insetBlockStart,
        prePosition:pre&&getComputedStyle(pre).position,
        preClass:pre?.className,
        rect:box&&{top:box.top,right:box.right,left:box.left,bottom:box.bottom},
        parent:parent&&{top:parent.top,right:parent.right,left:parent.left,bottom:parent.bottom},
        codeCSSLinks:[...document.querySelectorAll('link[rel="stylesheet"]')].map(l=>l.href).filter(h=>/code\.|chroma\./.test(h)),
      };
    });
    report.codeCSS=await page.evaluate(async()=>{
      const links=[...document.querySelectorAll('link[rel="stylesheet"]')]
        .filter(l=>/\\/code[.]/.test(l.href));
      return await Promise.all(links.map(async link=>{
        const response=await fetch(link.href,{cache:'no-cache'});
        const css=await response.text();
        return {url:link.href,status:response.status,bytes:css.length,
          hasCopySelector:css.includes('md-code-copy'),
          hasPositionRule:css.includes('position: absolute')||css.includes('position:absolute'),
          finalCSS:css.slice(-1200)};
      }));
    });
    await page.screenshot({path:path.join(output,'code-copy-placement.png')});
    check(report.copyGeometry.position==='absolute',
      'Copy button must be anchored above code, not flow after it: '+JSON.stringify(report.copyGeometry));
    check(report.copyGeometry.rect.top<report.copyGeometry.parent.top+80 &&
      report.copyGeometry.rect.right>report.copyGeometry.parent.right-40,
      'Copy button must sit at the top-right of its code block: '+JSON.stringify(report.copyGeometry));
    await page.evaluate(() => {
      window.__auditCopied=null;
      Object.defineProperty(navigator,'clipboard',{configurable:true,value:{
        writeText:async text=>{window.__auditCopied=text;},
      }});
    });
    await page.keyboard.press('Enter');
    await page.waitForFunction(()=>typeof window.__auditCopied==='string',{timeout:7000});
    check((await page.evaluate(()=>window.__auditCopied)).length>0,
      'Focused code-copy button activates with Enter');

    for (const width of widths) {
      await page.setViewportSize({width,height:1000});
      for(const theme of themes) {
        await page.evaluate(async ({mode,palette,name})=>{
          const root=document.documentElement;
          root.dataset.theme=mode;
          root.dataset.palette=palette;
          localStorage.setItem('color-mode',mode);
          localStorage.setItem('theme-palette-'+mode,palette);
          if(window.__markataLoadPaletteStyles) await window.__markataLoadPaletteStyles();
          const eventName=name==='light'?'color-mode-change':name==='custom'?'palette-change':'theme-change';
          window.dispatchEvent(new CustomEvent(eventName,{detail:{theme:mode,mode,palette}}));
        },theme);
        await page.waitForFunction(() => {
          const root=document.documentElement;
          const edge=document.querySelector('pre.mermaid svg path.flowchart-link');
          if (!edge) return false;
          const color=getComputedStyle(edge).stroke;
          return color!=='none' && color!=='' && getComputedStyle(root).getPropertyValue('--color-text').trim()!=='';
        },null,{timeout:20000});
        // Allow palette CSS and the asynchronous Mermaid rerender to settle.
        await page.waitForTimeout(700);
        state=await snapshot(page);
        report.observations.push({viewport:width,theme:theme.name,...state});
        check(state.rendered===12 && state.errors===0,
          `All Mermaid diagrams render at ${width}px/${theme.name}`);
        check(state.scrollWidth<=state.width+2,
          `No horizontal page overflow at ${width}px/${theme.name}`);
        check(Boolean(state.edge) && state.edge===state.expectedText,
          `Flowchart edges use readable theme text color at ${width}px/${theme.name}`);
        await page.screenshot({
          path:path.join(output,`${width}-${theme.name}.png`),animations:'disabled',
        });
      }
    }

    await page.emulateMedia({media:'print'});
    check((await page.locator('pre.mermaid svg').count())===12,'Print keeps Mermaid SVGs');

    noJS=await browser.newContext({
      viewport:{width:320,height:900},javaScriptEnabled:false,
    });
    const noPage=await noJS.newPage();
    const responseNoJS=await noPage.goto(base,{waitUntil:'domcontentloaded',timeout:90000});
    check(responseNoJS && responseNoJS.status()===200,'No-JS page responds successfully');
    check((await noPage.locator('article').innerText()).length>1200,
      'No-JS article remains readable');
    check((await noPage.locator('pre.mermaid').count())>=12,
      'No-JS Mermaid source is retained as readable text');
    check((await noPage.locator('article table').count())>=2,
      'No-JS semantic tables remain visible');
    await noPage.screenshot({path:path.join(output,'320-no-js.png')});
    report.pageErrors=errors;
  } catch(error) {
    report.failures.push(error.stack||String(error));
    process.exitCode=1;
  } finally {
    fs.writeFileSync(path.join(output,'report.json'),JSON.stringify(report,null,2)+'\n');
    await noJS?.close();
    await context?.close();
    await browser.close();
    console.log(JSON.stringify({
      successful:report.failures.length===0,
      checks:report.checks.length,
      lastError:report.failures[0]||null,
      screenshots:report.observations.length,
    },null,2));
  }
}
main().catch(error=>{console.error(error);process.exitCode=1;});
