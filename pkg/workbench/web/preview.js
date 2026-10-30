/* Mirrors the Python appearance bridge; normalize linked Go asset URLs. */
(()=>{const allowed=__MARKATA_CHOICES__;

const setAppearance=attributes=>{for(const [key,values] of Object.entries(allowed)){const value=attributes?.[key];
if(value==null)document.documentElement.removeAttribute(key);
else if(values.includes(value))document.documentElement.setAttribute(key,value);
}};

let updates=Promise.resolve();

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
document.querySelectorAll('script[data-markata-generated]:not([src])').forEach(s=>{if(!inlineScripts.has(s.textContent)&&!s.textContent.includes('rerenderMermaid'))s.remove();
});
for(const source of sources){if(source.src&&[...document.scripts].some(s=>s.src===assetURL(source,'src')))continue;
if(!source.src&&inlineScripts.has(source.textContent)&&[...document.scripts].some(s=>!s.src&&s.textContent===source.textContent))continue;
const script=document.createElement('script');
for(const attr of source.attributes)script.setAttribute(attr.name,attr.value);
script.dataset.markataGenerated='true';
script.textContent=source.textContent;
if(source.src||source.type==='module'){await new Promise(resolve=>{script.addEventListener('load',resolve,{once:true});
script.addEventListener('error',resolve,{once:true});
document.body.append(script);
});
}else document.body.append(script);
}dispatchEvent(new Event('view-transition-complete'));
setAppearance(data.attributes);
dispatchEvent(new Event('theme-change'));
};

addEventListener('message',e=>{if(e.source!==parent)return;
if(e.data?.type==='markdown-render'){const article=document.querySelector('.post-content');
if(article)article.innerHTML=e.data.html;
updates=updates.then(()=>renderMarkdown(e.data)).catch(console.error);
return;
}if(e.data?.type==='appearance'){updates=updates.then(()=>{setAppearance(e.data.attributes);
dispatchEvent(new Event('theme-change'));
}).catch(console.error);
}});
})();
