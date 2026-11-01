/* Mirrors the Python appearance bridge; normalize linked Go asset URLs. */
(()=>{const allowed=__MARKATA_CHOICES__;

const setAppearance=attributes=>{for(const [key,values] of Object.entries(allowed)){const value=attributes?.[key];
if(value==null)document.documentElement.removeAttribute(key);
else if(values.includes(value))document.documentElement.setAttribute(key,value);
}};

let updates=Promise.resolve(), scriptSequence=0;
let sourceLines=[], sourceLineCount=0, programmed=false, scrollFrame=0;
const scrollState=user=>{
  const article=document.querySelector('.post-content');
  if(!article)return;
  const elements=[...article.children];
  const points=[{line:0,y:0}];
  for(let i=0;i<Math.min(elements.length,sourceLines.length);i++){
    const line=sourceLines[i];
    if(line<=points.at(-1).line)continue;
    points.push({line,y:elements[i].getBoundingClientRect().top+scrollY});
  }
  points.push({line:Math.max(sourceLineCount,(sourceLines.at(-1)??0)+1),y:document.documentElement.scrollHeight});
  parent.postMessage({type:'sync-scroll-state',user,top:scrollY,height:document.documentElement.scrollHeight,clientHeight:innerHeight,points},'*');
};
addEventListener('scroll',()=>{
  if(scrollFrame)cancelAnimationFrame(scrollFrame);
  scrollFrame=requestAnimationFrame(()=>{scrollFrame=0;scrollState(!programmed);programmed=false;});
},{passive:true});
addEventListener('resize',()=>scrollState(false));
const originalScripts=new WeakMap();
const scriptSource=s=>originalScripts.get(s) ?? s.textContent;

const renderMarkdown=async data=>{const article=document.querySelector('.post-content');
if(!article)return;
if(window.Chart)article.querySelectorAll('canvas').forEach(canvas=>window.Chart.getChart?.(canvas)?.destroy());

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
const inlineScripts=new Set(sources.filter(s=>!s.src&&s.textContent.includes('rerenderMermaid')).map(s=>s.textContent));
document.querySelectorAll('script[data-markata-generated]:not([src])').forEach(s=>{if(!inlineScripts.has(scriptSource(s))&&!s.textContent.includes('rerenderMermaid'))s.remove();
});
for(const source of sources){if(source.src&&[...document.scripts].some(s=>s.src===assetURL(source,'src')))continue;
if(!source.src&&inlineScripts.has(source.textContent)&&[...document.scripts].some(s=>!s.src&&scriptSource(s)===source.textContent))continue;
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
}else if(source.src){await new Promise(resolve=>{script.addEventListener('load',resolve,{once:true});
script.addEventListener('error',resolve,{once:true});
document.body.append(script);
});
}else document.body.append(script);
}dispatchEvent(new Event('view-transition-complete'));
setAppearance(data.attributes);
dispatchEvent(new Event('theme-change'));
requestAnimationFrame(()=>scrollState(false));
};

addEventListener('message',e=>{if(e.source!==parent)return;
if(e.data?.type==='markdown-render'){const article=document.querySelector('.post-content');
if(article)article.innerHTML=e.data.html;
sourceLines=e.data.source_lines||[];
updates=updates.then(()=>renderMarkdown(e.data)).catch(console.error);
return;
}if(e.data?.type==='sync-scroll-lines'){sourceLines=e.data.lines||[];sourceLineCount=e.data.lineCount||0;requestAnimationFrame(()=>scrollState(false));return;
}if(e.data?.type==='sync-scroll'){programmed=true;scrollTo({top:e.data.y,behavior:'instant'});setTimeout(()=>{programmed=false;},100);return;
}if(e.data?.type==='appearance'){updates=updates.then(()=>{setAppearance(e.data.attributes);
dispatchEvent(new Event('theme-change'));
}).catch(console.error);
}});
document.addEventListener('click',event=>{
const link=event.target.closest?.('a.wikilink');
if(!link)return;
const path=new URL(link.getAttribute('href'),location.href).pathname;
const slug=decodeURIComponent(path.replace(/^\/+|\/+$/g,''));
if(!slug)return;
event.preventDefault();
parent.postMessage({type:'open-post',slug},'*');
});
})();
