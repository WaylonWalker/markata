/* Mirrors the Python appearance bridge; normalize linked Go asset URLs. */
(()=>{const allowed=__MARKATA_CHOICES__;

const setAppearance=attributes=>{for(const [key,values] of Object.entries(allowed)){const value=attributes?.[key];
if(value==null)document.documentElement.removeAttribute(key);
else if(values.includes(value))document.documentElement.setAttribute(key,value);
}};

let updates=Promise.resolve(), scriptSequence=0, latestVersion=0;
let sourceLines=[], sourceLineCount=0, sourceLineOffset=0, sourcePositionAnchors=[], programmedTarget=null, scrollFrame=0, positionLine=-1, positionNode=null, scrollIntentUntil=0;
let articleObserver;
const executedInlineScripts=new Set();
let appearanceSignature=JSON.stringify(Object.fromEntries(Object.keys(allowed).map(key=>[key,document.documentElement.getAttribute(key)])));
const markerStyle=document.createElement('style');markerStyle.textContent='.markata-sync-position{outline:1px solid color-mix(in srgb,var(--color-primary,#d39b32) 38%,transparent);outline-offset:2px;border-radius:3px}';document.head.append(markerStyle);
const blockSignature=node=>{
  if(node.matches('pre.mermaid'))return `mermaid:${node.dataset.markataMermaidSource??node.textContent}`;
  const generated=node.matches('.markata-sync-position,[data-markata-code-enhanced],.chartjs-container,[data-markata-source-line],[data-markata-source-end-line]')||node.getAttribute('class')===''||node.querySelector('[data-markata-ui],[data-markata-code-enhanced],.chartjs-container canvas,[data-markata-source-line],[data-markata-source-end-line]');
  if(!generated)return node.outerHTML;
  const clone=node.cloneNode(true);
  if(clone.hasAttribute('data-markata-source-line'))clone.removeAttribute('data-markata-source-line');
  clone.querySelectorAll('[data-markata-source-line],[data-markata-source-end-line]').forEach(child=>{
    child.removeAttribute('data-markata-source-line');child.removeAttribute('data-markata-source-end-line');
  });
  clone.removeAttribute('data-markata-source-end-line');
  clone.querySelectorAll('[data-markata-ui]').forEach(child=>child.remove());
  clone.querySelectorAll('.markata-sync-position').forEach(child=>child.classList.remove('markata-sync-position'));
  const normalizeCodeChrome=pre=>{
    const generatedLanguage=pre.dataset.markataGeneratedLanguage==='true';
    pre.removeAttribute('data-markata-code-enhanced');
    pre.removeAttribute('data-markata-generated-language');
    pre.classList.remove('md-code-enhanced','markata-sync-position');
    if(generatedLanguage)pre.removeAttribute('data-language');
    if(!pre.classList.length)pre.removeAttribute('class');
  };
  if(clone.hasAttribute('data-markata-code-enhanced')||clone.classList.contains('markata-sync-position'))normalizeCodeChrome(clone);
  clone.querySelectorAll('[data-markata-code-enhanced]').forEach(normalizeCodeChrome);
  if(clone.getAttribute('class')==='')clone.removeAttribute('class');
  clone.querySelectorAll('.chartjs-container canvas').forEach(canvas=>canvas.remove());
  clone.querySelectorAll('.chartjs-container [id]').forEach(element=>element.removeAttribute('id'));
  clone.querySelectorAll('.chartjs-container details').forEach(details=>details.removeAttribute('open'));
  return clone.outerHTML;
};
const patchArticle=(article,html)=>{
  const parseStarted=performance.now();
  const template=document.createElement('template');template.innerHTML=html;
  const parseMs=performance.now()-parseStarted;
  const reconcileStarted=performance.now();
  for(const child of [...article.childNodes])if(child.nodeType===Node.TEXT_NODE&&!child.textContent.trim())child.remove();
  const previous=[...article.children],available=new Map();
  for(const node of previous){const key=blockSignature(node);if(!available.has(key))available.set(key,[]);available.get(key).push(node);}
  const changedNodes=[],retained=new Set(),next=[];
  for(const fresh of [...template.content.children]){
    const matches=available.get(blockSignature(fresh));
    const existing=matches?.shift();
    if(existing){syncSourceLines(existing,fresh);next.push(existing);retained.add(existing);}
    else{next.push(fresh);changedNodes.push(fresh);}
  }
  for(const node of previous)if(!retained.has(node))node.remove();
  let cursor=article.firstElementChild;
  for(const node of next){if(node===cursor)cursor=cursor.nextElementSibling;else article.insertBefore(node,cursor);}
  return {changedNodes,total:next.length,retained:retained.size,removed:previous.length-retained.size,parseMs,reconcileMs:performance.now()-reconcileStarted};
};
const syncSourceLines=(existing,fresh)=>{
  const selector='[data-markata-source-line],[data-markata-source-end-line]';
  const markers=node=>[...(node.matches(selector)?[node]:[]),...node.querySelectorAll(selector)];
  const current=markers(existing),next=markers(fresh);
  if(current.length!==next.length)return;
  for(let i=0;i<current.length;i++)for(const attribute of ['data-markata-source-line','data-markata-source-end-line']){
    const value=next[i].getAttribute(attribute);
    if(value==null)current[i].removeAttribute(attribute);else current[i].setAttribute(attribute,value);
  }
};
const syncPositionMarker=()=>{
  const article=document.querySelector('.post-content');
  if(!article)return;
  let block=null,container=null;
  for(const anchor of sourcePositionAnchors){
    if(!anchor.nested){if(anchor.line<=positionLine)block=anchor;continue;}
    if(anchor.line>positionLine||anchor.endLine<positionLine)continue;
    if(!container||(anchor.endLine-anchor.line<container.endLine-container.line)||
      (anchor.endLine-anchor.line===container.endLine-container.line&&anchor.depth>container.depth))container=anchor;
  }
  const point=container||block||sourcePositionAnchors[0];
  const next=point?.node||(point&&point.index>=0&&point.index<article.children.length?article.children[point.index]:null);
  if(positionNode===next)return;
  positionNode?.classList.remove('markata-sync-position');
  if(positionNode&&!positionNode.classList.length)positionNode.removeAttribute('class');
  positionNode=next;positionNode?.classList.add('markata-sync-position');
};
const scrollState=(user,syncID=null)=>{
  const article=document.querySelector('.post-content');
  if(!article)return;
  const elements=[...article.children];
  const rect=article.getBoundingClientRect();
  const firstLine=Number.isFinite(sourceLines[0])?sourceLines[0]:0;
  const points=[{line:firstLine,y:rect.top+scrollY,index:0}];
  const positionAnchors=[];
  for(let i=0;i<Math.min(elements.length,sourceLines.length);i++){
    if(elements[i].getClientRects().length===0)continue;
    const line=sourceLines[i];
    const y=elements[i].getBoundingClientRect().top+scrollY;
    if(Number.isFinite(line)&&Number.isFinite(y))positionAnchors.push({line,index:i,node:elements[i],nested:false});
    if(!Number.isFinite(line)||!Number.isFinite(y))continue;
    points.push({line,y,index:i});
  }
  const selector='[data-markata-source-line],[data-markata-source-end-line]';
  for(const element of article.querySelectorAll(selector)){
    if(element.getClientRects().length===0)continue;
    const rect=element.getBoundingClientRect();
    let root=element,depth=0;
    while(root.parentElement&&root.parentElement!==article){root=root.parentElement;depth++;}
    const index=elements.indexOf(root);
    if(index<0)continue;
    const startValue=element.getAttribute('data-markata-source-line');
    const endValue=element.getAttribute('data-markata-source-end-line');
    const start=startValue==null?NaN:Number(startValue);
    const end=endValue==null?NaN:Number(endValue);
    let startLine=null,endLine=null;
    if(Number.isFinite(start)){
      startLine=start+sourceLineOffset;
      endLine=Number.isFinite(end)?end+sourceLineOffset:Math.max(startLine,sourceLineCount);
      const y=rect.top+scrollY;
      points.push({line:startLine,y,index});
      positionAnchors.push({line:startLine,endLine,index,node:element,nested:true,depth});
    }
    if(Number.isFinite(end)&&startLine!=null){
      const line=end+sourceLineOffset,y=rect.bottom+scrollY;
      points.push({line,y,index});
    }
  }
  sourcePositionAnchors=positionAnchors.sort((a,b)=>a.line-b.line);
  points.sort((a,b)=>a.line-b.line);
  const ordered=[];
  for(const point of points){
    if(!Number.isFinite(point.line)||!Number.isFinite(point.y))continue;
    if(ordered.at(-1)?.line===point.line)ordered[ordered.length-1]=point;
    else if(!ordered.length||point.line>ordered.at(-1).line)ordered.push(point);
  }
  const endLine=Math.max(sourceLineCount,(sourceLines.at(-1)??0)+1);
  const bottom=rect.bottom+scrollY;
  if(endLine>ordered.at(-1).line&&bottom>=ordered.at(-1).y)ordered.push({line:endLine,y:bottom,index:elements.length});
  const visualPoints=[...ordered].sort((a,b)=>a.y-b.y||a.line-b.line);
  syncPositionMarker();
  parent.postMessage({type:'sync-scroll-state',user,top:scrollY,height:document.documentElement.scrollHeight,clientHeight:innerHeight,points:ordered,visual_points:visualPoints,sync_id:syncID},'*');
};
const noteScrollIntent=duration=>{scrollIntentUntil=performance.now()+duration;};
document.addEventListener('wheel',()=>noteScrollIntent(250),{passive:true,capture:true});
document.addEventListener('touchstart',()=>noteScrollIntent(500),{passive:true,capture:true});
document.addEventListener('touchmove',()=>noteScrollIntent(250),{passive:true,capture:true});
document.addEventListener('pointerdown',()=>noteScrollIntent(100),{passive:true,capture:true});
document.addEventListener('pointermove',event=>{if(event.buttons)noteScrollIntent(250);},{passive:true,capture:true});
document.addEventListener('pointerup',()=>noteScrollIntent(100),{passive:true,capture:true});
document.addEventListener('keydown',event=>{const scrollKey=['ArrowUp','ArrowDown','PageUp','PageDown','Home','End'].includes(event.key)||(event.key===' '&&!event.target.closest?.('input,textarea,button,[contenteditable="true"]'));if(scrollKey)noteScrollIntent(250);},true);
addEventListener('scroll',()=>{
  if(scrollFrame)cancelAnimationFrame(scrollFrame);
  scrollFrame=requestAnimationFrame(()=>{scrollFrame=0;const request=programmedTarget;const isProgrammed=request!=null&&Math.abs(scrollY-request.y)<2;if(isProgrammed)programmedTarget=null;scrollState(!isProgrammed&&performance.now()<scrollIntentUntil,isProgrammed?request.id:null);});
},{passive:true});
addEventListener('resize',()=>scrollState(false));
const observeArticle=article=>{
  if(!article)return;
  if(!articleObserver){
    articleObserver=new ResizeObserver(()=>requestAnimationFrame(()=>scrollState(false)));
    article.addEventListener('load',()=>requestAnimationFrame(()=>scrollState(false)),true);
  }
  articleObserver.disconnect();articleObserver.observe(article);
};
if(document.fonts?.ready)document.fonts.ready.then(()=>scrollState(false));
const originalScripts=new WeakMap();
const scriptSource=s=>originalScripts.get(s) ?? s.textContent;

const renderMarkdown=async(data,changedNodes)=>{const article=document.querySelector('.post-content');
if(!article)return;
if(window.Chart)changedNodes.forEach(node=>[...(node.matches('canvas')?[node]:[]),...node.querySelectorAll('canvas')].forEach(canvas=>window.Chart.getChart?.(canvas)?.destroy()));

const head=document.createElement('template');
head.innerHTML=data.feature_head_html;
const styles=[...head.content.querySelectorAll('style,link[rel=stylesheet]')];
const assetURL=(node,attribute)=>new URL(node.getAttribute(attribute),document.baseURI).href;
const key=n=>n.tagName==='STYLE'?n.textContent:assetURL(n,'href');
const desired=new Set(styles.map(key));
document.querySelectorAll('[data-markata-render-head]').forEach(n=>{if(!desired.has(key(n)))n.remove();
});
for(const source of styles){const existing=[...document.head.querySelectorAll('style,link[rel=stylesheet]')];
if(existing.some(n=>key(n)===key(source)))continue;
const style=source.cloneNode(true);
style.dataset.markataRenderHead='true';
document.head.append(style);
}
const body=document.createElement('template');
body.innerHTML=data.body_end_html;
const sources=[...body.content.querySelectorAll('script')];
const chartChanged=changedNodes.some(node=>node.matches('.chartjs-container')||node.querySelector('.chartjs-container'));
const changedDiagrams=changedNodes.flatMap(node=>[...(node.matches('pre.mermaid')?[node]:[]),...node.querySelectorAll('pre.mermaid')]);
for(const source of sources){if(source.src&&[...document.scripts].some(s=>s.src===assetURL(source,'src')))continue;
if(!source.src){
  const scriptText=source.textContent;
  const chartID=scriptText.match(/chartjs-\d+/)?.[0];
  if(chartID){
    const canvas=document.getElementById(chartID);
    if(!canvas||(!chartChanged&&window.Chart?.getChart?.(canvas))){executedInlineScripts.add(scriptText);continue;}
  }else if(executedInlineScripts.has(scriptText)||[...document.scripts].some(s=>!s.src&&scriptSource(s)===scriptText)){executedInlineScripts.add(scriptText);continue;}
  executedInlineScripts.add(scriptText);
}
const script=document.createElement('script');
for(const attr of source.attributes)script.setAttribute(attr.name,attr.value);
script.dataset.markataGenerated='true';
script.textContent=source.textContent;
originalScripts.set(script,source.textContent);
if(!source.src&&source.type==='module'){
  await new Promise(resolve=>{
    const event='markata-module-ready-'+(++scriptSequence);
    addEventListener(event,resolve,{once:true});
    script.addEventListener('error',resolve,{once:true});
    script.textContent+='\n;window.dispatchEvent(new Event('+JSON.stringify(event)+'));';
    document.body.append(script);
  });
  script.remove();
}else if(source.src){await new Promise(resolve=>{script.addEventListener('load',resolve,{once:true});
script.addEventListener('error',resolve,{once:true});
document.body.append(script);
});
}else{document.body.append(script);script.remove();}
}
window.markataInitCodeChrome?.(changedNodes);
if(changedDiagrams.length)window.initMermaid?.(changedDiagrams);
const nextAppearance=JSON.stringify(data.attributes||{});
if(nextAppearance!==appearanceSignature){setAppearance(data.attributes);dispatchEvent(new Event('theme-change'));appearanceSignature=nextAppearance;}
observeArticle(article);
};

addEventListener('message',e=>{if(e.source!==parent)return;
if(e.data?.type==='markdown-render'){
  latestVersion=Math.max(latestVersion,e.data.version||0);
  updates=updates.then(async()=>{
    if(e.data.version<latestVersion)return;
    const article=document.querySelector('.post-content');if(!article)return;
    const updateStarted=performance.now();
    sourceLines=e.data.source_lines||[];sourceLineCount=e.data.line_count||sourceLineCount;sourceLineOffset=e.data.source_offset||0;
    const patch=patchArticle(article,e.data.html);
    const patchFinished=performance.now();
    await renderMarkdown(e.data,patch.changedNodes);
    const featuresFinished=performance.now();
    if(e.data.version!==latestVersion)return;
    scrollState(false);
    const updateFinished=performance.now();
    const diagrams=patch.changedNodes.reduce((count,node)=>count+(node.matches('pre.mermaid')?1:0)+node.querySelectorAll('pre.mermaid').length,0);
    parent.postMessage({type:'preview-updated',version:e.data.version,patch:{blocks:patch.total,changed:patch.changedNodes.length,reused:patch.retained,removed:patch.removed,diagrams,kinds:patch.changedNodes.map(node=>node.tagName.toLowerCase()+(node.classList.length?'.'+[...node.classList].join('.'):'')),timing:{html_parse_ms:patch.parseMs,dom_patch_ms:patch.reconcileMs,features_ms:featuresFinished-patchFinished,total_ms:updateFinished-updateStarted}}},'*');
  }).catch(console.error);
return;
}if(e.data?.type==='sync-scroll-lines'){sourceLines=e.data.lines||[];sourceLineCount=e.data.lineCount||0;sourceLineOffset=e.data.sourceOffset||0;requestAnimationFrame(()=>scrollState(false));return;
}if(e.data?.type==='sync-scroll'){programmedTarget={id:e.data.id,y:Math.max(0,Math.min(e.data.y,document.documentElement.scrollHeight-innerHeight))};document.scrollingElement.scrollTop=programmedTarget.y;const id=e.data.id;setTimeout(()=>{if(programmedTarget?.id===id)programmedTarget=null;},300);return;
}if(e.data?.type==='sync-position'){positionLine=e.data.line;syncPositionMarker();return;
}if(e.data?.type==='appearance'){updates=updates.then(()=>{setAppearance(e.data.attributes);
dispatchEvent(new Event('theme-change'));appearanceSignature=JSON.stringify(e.data.attributes||{});
}).catch(console.error);
}});
document.addEventListener('click',event=>{
const link=event.target.closest?.('a');
if(!link)return;
const href=link.getAttribute('href')||'';
if(!link.classList.contains('wikilink')&&(!href.startsWith('/')||href.startsWith('//')))return;
const path=new URL(href,location.href).pathname;
const slug=decodeURIComponent(path.replace(/^\/+|\/+$/g,''));
if(!slug)return;
event.preventDefault();
parent.postMessage({type:'open-post',slug},'*');
},true);
})();
