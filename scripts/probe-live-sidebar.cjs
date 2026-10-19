// Verify current production CSS after main-branch site hotfix was published.
// Production-only browser probe. Reports the actual deployed markup and script
// failures before and after a physical click on the sidebar handle.
const { chromium } = require('playwright');
const path = require('node:path');
const fs = require('node:fs');
(async () => {
  const browser = await chromium.launch({headless: true,args:['--no-sandbox']});
  const reports = [];
  try {
    for (const width of [412, 1440]) {
      const context = await browser.newContext({viewport:{width,height:850},ignoreHTTPSErrors:false});
      for (const route of ['/archive/', '/have-you-taken-the-agentic-pill-yet/']) {
        const page = await context.newPage();
        const errors = [], failed = [], scripts = [];
        page.on('pageerror',e=>errors.push(String(e).slice(0,600)));
        page.on('console',m=>{ if(m.type()==='error')errors.push('console: '+m.text().slice(0,400));});
        page.on('requestfailed',r=>failed.push(r.url()+' '+r.failure()?.errorText));
        page.on('response',r=>{ if(r.request().resourceType()==='script')scripts.push({url:r.url(),status:r.status()}); });
        const report={route,width};
        try {
          const response=await page.goto('https://waylonwalker.com'+route,{waitUntil:'domcontentloaded',timeout:45000});
          await page.waitForTimeout(2000);
          report.status=response.status();
          report.url=page.url();
          const selector='.feed-sidebar > .sidebar-toggle';
          report.stateBefore=await page.evaluate(sel=>{
            const b=document.querySelector(sel);
            if(!b)return {missing:sel,available:[...document.querySelectorAll('.sidebar-toggle')].map(x=>x.outerHTML.slice(0,300))};
            const a=b.closest('aside'),cs=getComputedStyle(b),as=getComputedStyle(a),r=b.getBoundingClientRect();
            const x=Math.min(innerWidth-1,Math.max(1,r.left+r.width/2)),y=Math.min(innerHeight-1,Math.max(1,r.top+r.height/2));
            return {button:b.outerHTML.slice(0,400),asideClass:a.className,asideHTMLLength:a.outerHTML.length, rect:{x:r.x,y:r.y,w:r.width,h:r.height}, buttonVisibility:cs.visibility,buttonPointerEvents:cs.pointerEvents,asideVisibility:as.visibility,asideTransform:as.transform,hit:document.elementFromPoint(x,y)?.outerHTML.slice(0,260),initPin:typeof window.initSidebarPinButtons, initRescue:typeof window.toggleSidebarPinned, scripts:[...document.querySelectorAll('script[src]')].filter(s=>/sidebar|view-transitions|feed-cycling/.test(s.src)).map(s=>s.src)};
          },selector);
          if(!report.stateBefore.missing){
            try{await page.locator(selector).click({timeout:5000}); report.click='passed';}
            catch(e){report.click='failed: '+String(e).slice(0,700); }
            await page.waitForTimeout(700);
            report.stateAfter=await page.evaluate(sel=>{
              const b=document.querySelector(sel), a=b?.closest('aside');
              const pinnedStyles=[],stylesheets=[];
              function checkRules(rules,href){
                if(!rules)return;
                for (const rule of rules) {
                  if(rule.selectorText && rule.selectorText.includes('sidebar--pinned')) {
                    pinnedStyles.push({href,selector:rule.selectorText.slice(0,220),visibility:rule.style?.visibility,transform:rule.style?.transform});
                  }
                  if(rule.cssRules)checkRules(rule.cssRules,href);
                }
              }
              for(const s of document.styleSheets){
                stylesheets.push(s.href);
                try{checkRules(s.cssRules,s.href)}catch(e){}
              }
              return {expanded:b?.getAttribute('aria-expanded'),pinned:a?.classList.contains('sidebar--pinned'),
                visibility:a&&getComputedStyle(a).visibility,transform:a&&getComputedStyle(a).transform,
                rect:a&&{x:a.getBoundingClientRect().x,width:a.getBoundingClientRect().width},
                styleSheets:stylesheets.filter(x=>x&&/components|main|palette/.test(x)),pinnedRuleCount:pinnedStyles.length,
                pinnedRules:pinnedStyles.slice(0,16)};
            },selector);
          }
        }catch(e){report.exception=String(e).slice(0,900);}
        report.errors=errors.slice(0,15);report.failedRequests=failed.slice(0,20);report.scriptResponses=scripts.slice(0,40);
        reports.push(report);console.log('SIDEBAR_PROBE '+JSON.stringify(report));await page.close();
      }
      await context.close();
    }
  } finally {await browser.close();}
  if (reports.some(r=>{
    if(r.exception || r.stateBefore?.missing || r.click?.startsWith('failed'))return true;
    const shouldOpen=!r.stateBefore.asideClass.includes('sidebar--pinned');
    if(r.stateAfter?.pinned!==shouldOpen)return true;
    if(shouldOpen && (r.stateAfter.visibility!=='visible' || r.stateAfter.transform!=='none'))return true;
    return false;
  })) process.exitCode=1;
})().catch(e=>{console.error(e);process.exitCode=1});
