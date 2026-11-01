/* Semantic best-effort rendered Markdown diff, kept inside the sandboxed preview. */
(() => {
  'use strict';
  const escape = text => String(text).replace(/[&<>"']/g, ch => ({
    '&':'&amp;', '<':'&lt;', '>':'&gt;', '"':'&quot;', "'":'&#39;'
  }[ch]));
  const tokens = text => text.match(/\s+|[^\s]+/g) || [];
  const safeTextTags = new Set(['P','H1','H2','H3','H4','H5','H6']);
  const styles = [
    '.mk-rich-add,.mk-rich-remove{padding-left:12px;border-left:3px solid;margin:12px 0}',
    '.mk-rich-add{border-color:#50a487;background:color-mix(in srgb,#50a487 9%,transparent)}',
    '.mk-rich-remove{border-color:#c87a83;background:color-mix(in srgb,#c87a83 9%,transparent)}',
    'ins.mk-rich-insert{background:color-mix(in srgb,#51a487 23%,transparent);text-decoration:none;border-radius:3px}',
    'del.mk-rich-delete{background:color-mix(in srgb,#c87a83 21%,transparent);text-decoration:line-through;border-radius:3px}',
    '.mk-rich-legend{font:12px/1.6 system-ui,sans-serif;opacity:.8;padding:12px 0 16px;border-bottom:1px solid #8885;margin-bottom:16px}',
    '.mk-rich-changed{border-left:2px solid #8b8a8a66;padding-left:9px}'
  ].join('');
  function lcsOps(a,b,limit) {
    if(a.length*b.length>limit)return null;
    const width=b.length+1;
    const dp=new Uint16Array((a.length+1)*width);
    for(let i=a.length-1;i>=0;i--)for(let j=b.length-1;j>=0;j--){
      const pos=i*width+j;
      dp[pos]=a[i]===b[j] ? 1+dp[pos+width+1] : Math.max(dp[pos+width],dp[pos+1]);
    }
    const ops=[];let i=0,j=0;
    while(i<a.length||j<b.length){
      if(i<a.length&&j<b.length&&a[i]===b[j]){ops.push('same');i++;j++;}
      else if(j<b.length&&(i===a.length||dp[i*width+j+1]>=dp[(i+1)*width+j])){ops.push('add');j++;}
      else{ops.push('del');i++;}
    }
    return ops;
  }
  function parseBlocks(html) {
    const doc=new DOMParser().parseFromString('<main id="mk-snippet">'+html+'</main>','text/html');
    return [...doc.querySelector('#mk-snippet').children].map(el=>({
      tag:el.tagName,text:el.textContent||'',markup:el.outerHTML,key:el.outerHTML
    }));
  }
  function whole(block,kind) {
    return block?'<div class="mk-rich-'+(kind==='add'?'add':'remove')+'">'+block.markup+'</div>':'';
  }
  function changed(before,after) {
    if(before.tag!==after.tag||!safeTextTags.has(after.tag))
      return whole(before,'del')+whole(after,'add');
    const old=tokens(before.text),fresh=tokens(after.text);
    const ops=lcsOps(old,fresh,65000);
    if(!ops)return whole(before,'del')+whole(after,'add');
    let i=0,j=0,html='';
    for(const op of ops){
      if(op==='same'){html+=escape(fresh[j]);i++;j++;}
      else if(op==='add')html+='<ins class="mk-rich-insert">'+escape(fresh[j++])+'</ins>';
      else html+='<del class="mk-rich-delete">'+escape(old[i++])+'</del>';
    }
    const tag=after.tag.toLowerCase();
    return '<'+tag+' class="mk-rich-changed">'+html+'</'+tag+'>';
  }
  function richHTML(oldHTML,newHTML) {
    const old=parseBlocks(oldHTML),fresh=parseBlocks(newHTML);
    const ops=lcsOps(old.map(x=>x.key),fresh.map(x=>x.key),50000);
    if(!ops)throw new Error('Too many blocks for a rich comparison. Use side-by-side or source diff.');
    let i=0,j=0,result='',removed=[],added=[];
    function flush(){
      const paired=Math.min(removed.length,added.length);
      for(let k=0;k<paired;k++)result+=changed(removed[k],added[k]);
      for(let k=paired;k<removed.length;k++)result+=whole(removed[k],'del');
      for(let k=paired;k<added.length;k++)result+=whole(added[k],'add');
      removed=[];added=[];
    }
    for(const op of ops){
      if(op==='same'){flush();result+=fresh[j].markup;i++;j++;}
      else if(op==='del')removed.push(old[i++]);
      else added.push(fresh[j++]);
    }
    flush();
    const legend='<div class="mk-rich-legend"><del class="mk-rich-delete">Removed</del> &nbsp; <ins class="mk-rich-insert">Added</ins> &nbsp; Unchanged content is unmarked.</div>';
    return legend+(result||'<p>No visible body changes. Check the frontmatter changes above.</p>');
  }
  function render(before,after,frame) {
    const rich=richHTML(before.html,after.html);
    if(!after.html||!after.document||!after.document.includes(after.html))
      throw new Error('Cannot locate the article content in its rendered document.');
    let page=after.document.replace(after.html,rich);
    page=page.replace('</head>','<style data-markata-rich-diff>'+styles+'</style></head>');
    frame.srcdoc=page;
  }
  window.MarkataRichPublishDiff={render,richHTML};
})();