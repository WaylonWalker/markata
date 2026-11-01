/* Mirrors the Python appearance bridge; normalize linked Go asset URLs. */
(()=>{const allowed=__MARKATA_CHOICES__;

const setAppearance=attributes=>{for(const [key,values] of Object.entries(allowed)){const value=attributes?.[key];
if(value==null)document.documentElement.removeAttribute(key);
else if(values.includes(value))document.documentElement.setAttribute(key,value);
}};

let updates=Promise.resolve(), scriptSequence=0, latestVersion=0;
let sourceLines=[], sourceLineCount=0, sourceLineOffset=0, sourcePositionAnchors=[], programmedTarget=null, scrollFrame=0, positionLine=-1, positionColumn=null, positionSourceText=null, positionPending=false, positionVimMode=null, positionNode=null, lastExactPosition=null, scrollIntentUntil=0;
let articleObserver;
const executedInlineScripts=new Set();
let appearanceSignature=JSON.stringify(Object.fromEntries(Object.keys(allowed).map(key=>[key,document.documentElement.getAttribute(key)])));
const markerStyle=document.createElement('style');markerStyle.textContent='.markata-sync-position{outline:1px solid color-mix(in srgb,var(--color-primary,#d39b32) 18%,transparent);outline-offset:2px;border-radius:3px}.markata-presence-layer{position:fixed;inset:0;z-index:2147483646;overflow:hidden;pointer-events:none}.markata-presence-cursor{--presence-color:var(--color-primary,#d39b32);position:fixed;width:1px;height:28px;background:var(--presence-color);opacity:.85;filter:drop-shadow(0 0 2px color-mix(in srgb,var(--presence-color) 35%,transparent))}.markata-presence-cursor[data-vim-mode="normal"],.markata-presence-cursor[data-vim-mode="visual"]{width:var(--presence-glyph-width,8px);height:var(--presence-line-height,20px);opacity:.48;filter:none;border-radius:2px}.markata-presence-cursor[data-vim-mode="visual"]{outline:1px solid color-mix(in srgb,var(--presence-color) 70%,transparent)}.markata-presence-cursor[data-vim-mode="replace"]{height:2px!important;opacity:.95;filter:none}.markata-presence-cursor[hidden]{display:none}.markata-presence-flag{position:absolute;bottom:calc(100% + 3px);left:-1px;padding:2px 6px;border-radius:3px 3px 3px 0;background:var(--presence-color);color:var(--color-on-primary,#11151d);font:600 10px/1.4 system-ui,sans-serif;white-space:nowrap;letter-spacing:.01em}';document.head.append(markerStyle);
const presenceLayer=document.createElement('div');presenceLayer.className='markata-presence-layer';presenceLayer.setAttribute('aria-hidden','true');document.body.append(presenceLayer);
const presenceParticipants=[{id:'local',label:'You',color:'var(--color-primary,#d39b32)'}];
const presenceCursors=new Map();
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
const visibleSourceText=source=>{
  const decoded=document.createElement('textarea');
  decoded.innerHTML=String(source||'').replace(/^\s{0,3}(?:#{1,6}\s+|>\s?|[-*+]\s+|\d+\.\s+)/,'').replace(/!\[([^\]]*)\]\([^)]+\)/g,'$1').replace(/\[([^\]]+)\]\([^)]+\)/g,'$1').replace(/\[([^\]]+)\]\[[^\]]*\]/g,'$1').replace(/<[^>]*>/g,'').replace(/[`*_~]/g,'');
  return decoded.value.replace(/[‘’]/g,"'").replace(/[“”]/g,'"').replace(/[–—]/g,'-').replace(/\s+/g,' ').trim();
};
const visibleSourcePrefix=source=>{
  const decoded=document.createElement('textarea');
  decoded.innerHTML=String(source||'').replace(/^\s{0,3}(?:#{1,6}\s+|>\s?|[-*+]\s+|\d+\.\s+)/,'').replace(/!\[([^\]]*)\]\([^)]+\)/g,'$1').replace(/\[([^\]]+)\]\([^)]+\)/g,'$1').replace(/\[([^\]]+)\]\[[^\]]*\]/g,'$1').replace(/<[^>]*>/g,'').replace(/[`*_~]/g,'');
  return decoded.value.replace(/[‘’]/g,"'").replace(/[“”]/g,'"').replace(/[–—]/g,'-').replace(/\s+/g,' ').replace(/^\s+/,'');
};
const textPosition=(root,source,ch)=>{
  const wanted=visibleSourceText(source);
  if(!wanted)return null;
  const walker=document.createTreeWalker(root,NodeFilter.SHOW_TEXT);
  const chars=[];let node,lastWasSpace=true;
  while((node=walker.nextNode()))for(let offset=0;offset<node.nodeValue.length;offset++){
    const value=node.nodeValue[offset].replace(/[‘’]/g,"'").replace(/[“”]/g,'"').replace(/[–—]/g,'-');
    if(/\s/.test(value)){if(!lastWasSpace&&chars.length){chars.push({value:' ',node,offset});lastWasSpace=true;}continue;}
    chars.push({value,node,offset});lastWasSpace=false;
  }
  while(chars.at(-1)?.value===' ')chars.pop();
  const rendered=chars.map(char=>char.value).join('');
  const start=rendered.indexOf(wanted);
  if(start<0)return null;
  const prefix=visibleSourcePrefix(String(source).slice(0,Math.max(0,ch??0)));
  const index=start+Math.min(prefix.length,wanted.length);
  const target=chars[index];
  const range=document.createRange();
  if(target)range.setStart(target.node,target.offset);
  else{const last=chars.at(-1);if(!last)return null;range.setStart(last.node,last.offset+1);}
  range.collapse(true);
  const rect=range.getBoundingClientRect();
  const glyph=target||chars.at(-1);
  const glyphRange=document.createRange();
  if(target){glyphRange.setStart(glyph.node,glyph.offset);glyphRange.setEnd(glyph.node,Math.min(glyph.offset+1,glyph.node.nodeValue.length));}
  else{glyphRange.setStart(glyph.node,Math.max(0,glyph.offset));glyphRange.setEnd(glyph.node,glyph.offset+1);}
  const glyphRect=glyphRange.getBoundingClientRect();
  return rect.width||rect.height?{left:rect.left,top:rect.top,height:rect.height,width:glyphRect.width}:null;
};
const syncPositionMarker=()=>{
  const article=document.querySelector('.post-content');
  if(!article)return;
  const scrollOff=140;
  article.style.paddingBottom=`${scrollOff}px`;
  let block=null,container=null;
  for(const anchor of sourcePositionAnchors){
    if(!anchor.nested){if(anchor.line<=positionLine)block=anchor;continue;}
    if(anchor.line>positionLine||anchor.endLine<positionLine)continue;
    if(!container||(anchor.endLine-anchor.line<container.endLine-container.line)||
      (anchor.endLine-anchor.line===container.endLine-container.line&&anchor.depth>container.depth))container=anchor;
  }
  const point=container||block||sourcePositionAnchors[0];
  const next=point?.node||(point&&point.index>=0&&point.index<article.children.length?article.children[point.index]:null);
  if(positionNode!==next){
    positionNode?.classList.remove('markata-sync-position');
    if(positionNode&&!positionNode.classList.length)positionNode.removeAttribute('class');
    positionNode=next;positionNode?.classList.add('markata-sync-position');
  }
  for(const participant of presenceParticipants){
    let cursor=presenceCursors.get(participant.id);
    if(!cursor){cursor=document.createElement('div');cursor.className='markata-presence-cursor';cursor.hidden=true;cursor.dataset.presenceId=participant.id;const flag=document.createElement('span');flag.className='markata-presence-flag';cursor.append(flag);presenceLayer.append(cursor);presenceCursors.set(participant.id,cursor);}
    cursor.querySelector('.markata-presence-flag').textContent=participant.label;
    cursor.style.setProperty('--presence-color',participant.color);
    if(!next){cursor.hidden=true;continue;}
    const rect=next.getBoundingClientRect();
    const visible=rect.bottom>0&&rect.top<innerHeight&&rect.right>0&&rect.left<innerWidth;
    cursor.hidden=!visible;
    if(visible){
      const span=point?.nested?point.endLine-point.line:0;
      const progress=span>0?Math.max(0,Math.min(1,(positionLine-point.line)/span)):0;
      const lineHeight=Math.max(16,Math.min(20,parseFloat(getComputedStyle(next).lineHeight)||18));
      const top=rect.top+Math.max(0,rect.height-lineHeight)*progress;
      const character=textPosition(next,positionSourceText,positionColumn);
      if(character){lastExactPosition={node:next,line:positionLine,left:character.left,top:character.top,height:lineHeight};}
      const keepExact=positionPending&&lastExactPosition?.node===next&&lastExactPosition.line===positionLine;
      const left=character?.left??(keepExact?lastExactPosition.left:null);
      const markerTop=character?.top??(keepExact?lastExactPosition.top:top);
      cursor.style.left=`${Math.max(2,left==null?rect.left+(point?.nested?10:-7):left-1)}px`;
      cursor.style.top=`${Math.max(0,Math.min(innerHeight-lineHeight,markerTop))}px`;
      cursor.dataset.vimMode=positionVimMode||'';
      cursor.style.setProperty('--presence-glyph-width',`${Math.max(6,character?.width||8)}px`);
      cursor.style.setProperty('--presence-line-height',`${lineHeight}px`);
      cursor.style.height=`${character?.height||lineHeight}px`;
      if(positionVimMode==='replace')cursor.style.top=`${Math.max(0,Math.min(innerHeight-2,markerTop+lineHeight-2))}px`;
      if(!positionPending&&!character)lastExactPosition=null;
    }
  }
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
  const bottom=rect.bottom+scrollY-scrollOff;
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
}if(e.data?.type==='sync-position'){positionLine=e.data.line;positionColumn=Number.isFinite(e.data.ch)?e.data.ch:null;positionSourceText=typeof e.data.text==='string'?e.data.text:null;positionPending=e.data.pending===true;positionVimMode=['normal','insert','replace','visual'].includes(e.data.vimMode)?e.data.vimMode:null;if(positionColumn==null)lastExactPosition=null;syncPositionMarker();return;
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
document.addEventListener('click',event=>{
  if(event.button!==0||event.target.closest?.('a')||!window.getSelection()?.isCollapsed)return;
  const article=document.querySelector('.post-content');
  if(!article)return;
  let range=document.caretRangeFromPoint?.(event.clientX,event.clientY)||null;
  if(!range&&document.caretPositionFromPoint){const caret=document.caretPositionFromPoint(event.clientX,event.clientY);if(caret){range=document.createRange();range.setStart(caret.offsetNode,caret.offset);range.collapse(true);}}
  if(!range||!article.contains(range.startContainer))return;
  let textBlock=range.startContainer.nodeType===Node.TEXT_NODE?range.startContainer.parentElement:range.startContainer;
  textBlock=textBlock?.closest?.('p,h1,h2,h3,h4,h5,h6,li,td,th,blockquote,pre')||null;
  if(!textBlock||!article.contains(textBlock))return;
  const text=String(textBlock.innerText||textBlock.textContent||'').replace(/\s+/g,' ').trim();
  if(!text)return;
  const before=document.createRange();before.selectNodeContents(textBlock);before.setEnd(range.startContainer,range.startOffset);
  const offset=String(before.toString()).replace(/\s+/g,' ').trimStart().length;
  let root=textBlock;while(root.parentElement&&root.parentElement!==article)root=root.parentElement;
  const index=[...article.children].indexOf(root);if(index<0)return;
  const nested=sourcePositionAnchors.filter(anchor=>anchor.nested&&anchor.index===index&&anchor.node.contains(textBlock)).sort((a,b)=>b.depth-a.depth)[0];
  const start=(nested?.line??((sourceLines[index]??0)+sourceLineOffset))-1;
  const nextTop=sourceLines.slice(index+1).find(Number.isFinite);
  const end=(nested?.endLine??(nextTop==null?sourceLineOffset+sourceLineCount:nextTop+sourceLineOffset-1))-1;
  parent.postMessage({type:'place-editor-cursor',line:start,endLine:Math.max(start,end),text,offset},'*');
},true);
})();
