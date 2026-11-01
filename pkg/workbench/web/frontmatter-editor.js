/* Source-backed frontmatter properties. No second persisted metadata model. */
(() => {
  'use strict';
  const definitions = [
    ['title','Title','text','Page basics','Name shown on the page and in search results'],
    ['description','Description','textarea','Page basics','Short summary for search and social previews'],
    ['slug','URL slug','text','Page basics','Leave blank to use the file name'],
    ['date','Publication date','text','Page basics','Use YYYY-MM-DD or a supported timestamp'],
    ['published','Published','boolean','Page basics','Show this page in public feeds and listings'],
    ['draft','Draft','boolean','Page basics','Keep this page out of site output'],
    ['tags','Tags','list','Page basics','Help readers find related pages'],
    ['image','Featured image','text','Presentation','Social and Open Graph image'],
    ['video','Featured video','text','Presentation','Media URL for video cards'],
    ['cover','Cover','text','Presentation','Explicit cover media'],
    ['cover_image','Cover image','text','Presentation','Image-library cover override'],
    ['content_width','Content width','text','Presentation','wide opts into story layout'],
    ['author','Author','text','Presentation','Single author name'],
    ['authors','Authors','list','Presentation','Multiple authors or structured author references'],
    ['template','Template','text','Presentation','Choose a content template'],
    ['templateKey','Template key','text','Presentation','Site or built-in template identifier'],
    ['layout','Layout','text','Presentation','Override the page layout'],
    ['toc','Table of contents','boolean','Presentation','Show a table of contents'],
    ['sidebar','Sidebar','boolean','Presentation','Show a sidebar'],
    ['modified','Modified date','text','Publishing & navigation','Date of most recent update'],
    ['publishdate','Publish date alias','text','Publishing & navigation','Alternative publication date'],
    ['pubdate','Pubdate alias','text','Publishing & navigation','Alternative publication date'],
    ['lastmod','Last modified alias','text','Publishing & navigation','Alternative modified date'],
    ['updated','Updated alias','text','Publishing & navigation','Alternative modified date'],
    ['updated_at','Updated-at alias','text','Publishing & navigation','Alternative modified date'],
    ['last_modified','Last-modified alias','text','Publishing & navigation','Alternative modified date'],
    ['by','By (author alias)','text','Publishing & navigation','Alternative single-author field'],
    ['writer','Writer (author alias)','text','Publishing & navigation','Alternative single-author field'],
    ['series','Series','text','Publishing & navigation','Connect this page to a series'],
    ['prevnext_feed','Previous/next feed','text','Publishing & navigation','Navigation feed or series slug'],
    ['aliases','Aliases','list','Publishing & navigation','Other names and handles'],
    ['canonical_url','Canonical URL','text','Publishing & navigation','Original source for republished content'],
    ['private','Private','boolean','Advanced','Private publishing controls'],
    ['skip','Skip processing','boolean','Advanced','Exclude this content from processing'],
    ['templates','Templates','object','Advanced','Per-output template mapping'],
    ['secret_key','Secret key','text','Advanced','Sensitive: only edit in a trusted workspace']
  ];
  const common = new Set(['title','description','slug','date','published','draft','tags','image']);
  const known = new Map(definitions.map(row => [row[0],{key:row[0],label:row[1],type:row[2],group:row[3],help:row[4]}]));
  const panel = document.querySelector('#frontmatter-panel');
  const fieldsRoot = document.querySelector('#fm-fields');
  const message = document.querySelector('#fm-message');
  const addSelect = document.querySelector('#fm-add-select');
  const addSearch = document.querySelector('#fm-add-search');
  const newKey = document.querySelector('#fm-new-key');
  const toggle = document.querySelector('#frontmatter-toggle');
  const rawButton = document.querySelector('#frontmatter-raw');
  const moreButton = document.querySelector('#fm-more');
  const addToggle = document.querySelector('#fm-add-toggle');
  const addFlow = document.querySelector('#fm-add-flow');
  const addKnown = document.querySelector('#fm-add');
  const addCustom = document.querySelector('#fm-add-custom');
  const customField = document.querySelector('#fm-custom-property');
  const editorPane = document.querySelector('#editor-pane');
  let editor, raw = false, marker, expanded = false, refreshing = 0, refreshTimer, busy = false, seenFields = [];
  let showMore = false;
  let fieldId = 0;

  function frontmatterEnd(source) {
    const match = /^---\r?\n(?:[\s\S]*?\r?\n)?---(?:\r?\n|$)/.exec(source);
    return match ? match[0].length : 0;
  }

  function showMessage(text, isError = false) {
    message.textContent = text;
    message.hidden = !text;
    message.dataset.error = String(isError);
  }

  function collapseYAML() {
    if (!editor) return;
    if (marker) { marker.clear(); marker = null; }
    if (raw) return;
    const end = frontmatterEnd(editor.getValue());
    if (!end) return;
    const stub = document.createElement('span');
    stub.className = 'fm-collapsed';
    stub.textContent = 'Frontmatter hidden · Edit properties';
    stub.setAttribute('role','button');
    stub.tabIndex = 0;
    stub.addEventListener('click', openProperties);
    stub.addEventListener('keydown', e => {if(e.key==='Enter'||e.key===' '){e.preventDefault();openProperties();}});
    marker = editor.markText(editor.posFromIndex(0), editor.posFromIndex(end), {collapsed:true,replacedWith:stub});
  }

  function setRaw(next) {
    raw = next;
    rawButton.textContent = raw ? 'Hide YAML' : 'Edit YAML';
    rawButton.setAttribute('aria-pressed',String(raw));
    if (raw) {
      if (marker) {marker.clear();marker=null;}
      if (expanded) closeProperties();
      editor.focus();
      editor.setCursor(1,0);
    } else {
      collapseYAML();
      editor.focus();
    }
  }

  function closeProperties() {
    expanded=false;
    panel.hidden=true;
    editorPane.removeAttribute('data-frontmatter-open');
    toggle.setAttribute('aria-expanded','false');
    toggle.textContent='Properties';
    editor.refresh();
  }
  async function openProperties() {
    if (!editor) return;
    const layout = document.querySelector('#setting-layout');
    if (document.querySelector('#workbench').dataset.layout === 'preview') {
      layout.value='split';
      layout.dispatchEvent(new Event('change'));
    }
    if (raw) {raw=false;rawButton.textContent='Edit YAML';rawButton.setAttribute('aria-pressed','false');collapseYAML();}
    expanded=true;panel.hidden=false;editorPane.dataset.frontmatterOpen='true';
    toggle.setAttribute('aria-expanded','true');
    toggle.textContent='Back to Markdown';
    editor.refresh();
    await refresh();
  }
  function toggleProperties() { return expanded ? closeProperties() : openProperties(); }

  async function request(operation,payload) {
    const response = await fetch('/workspace/frontmatter/'+operation,{
      method:'POST',headers:{'Content-Type':'application/json'},
      body:JSON.stringify(payload),cache:'no-store'
    });
    const data = await response.json();
    if (!response.ok) throw new Error(data.error || 'Could not update frontmatter');
    return data;
  }
  function nodeButton(text,run,title) {
    const button=document.createElement('button');
    button.type='button';button.textContent=text;
    if(title) button.title=title;
    button.addEventListener('click',run);
    return button;
  }
  function inputFor(value,kind,commit) {
    let control;
    if (kind==='boolean') {
      control=document.createElement('label');control.className='fm-switch';
      const checkbox=document.createElement('input');checkbox.type='checkbox';checkbox.checked=Boolean(value);
      const state=document.createElement('span');state.textContent=checkbox.checked?'On':'Off';
      checkbox.addEventListener('change',()=>{state.textContent=checkbox.checked?'On':'Off';commit(checkbox.checked);});
      control.append(checkbox,state);
      return control;
    } else if (kind==='null') {
      const label=document.createElement('span');label.textContent='null';return label;
    } else {
      control=document.createElement(kind==='textarea'?'textarea':'input');
      if (control.tagName==='INPUT') control.type=kind==='number'?'number':'text';
      if (kind==='number') control.step='any';
      control.value=value==null?'':String(value);
      control.addEventListener('change',()=>{
        if (kind==='number') {
          if (!control.value.trim() || !Number.isFinite(Number(control.value))) {showMessage('Enter a valid number.',true);return;}
          commit(Number(control.value));
        } else commit(control.value);
      });
    }
    control.className='fm-input';
    return control;
  }
  function valueType(v) {
    if (v===null) return 'null';
    if (Array.isArray(v)) return 'list';
    if (typeof v==='boolean') return 'boolean';
    if (typeof v==='number') return 'number';
    if (typeof v==='object') return 'object';
    return 'text';
  }
  function emptyFor(kind) {
    const values={list:[],object:{},boolean:false,number:0,null:null};
    return Object.hasOwn(values,kind) ? values[kind] : '';
  }

  // Recursive editor for arbitrary nested YAML maps, lists, and mixed author entries.
  // Edits replace only their containing top-level YAML property on the server.
  function structured(value,commit,depth=0) {
    const wrapper=document.createElement('div');wrapper.className='fm-structured';
    if (depth>8) {
      wrapper.append('Too deeply nested for the form. Use Edit YAML.');
      return wrapper;
    }
    const type=valueType(value);
    const typePicker=document.createElement('select');typePicker.className='fm-type';
    typePicker.setAttribute('aria-label','Value type');
    for(const kind of ['text','number','boolean','null','list','object']) {
      const option=document.createElement('option');option.value=kind;option.textContent=kind==='text'?'Text':kind[0].toUpperCase()+kind.slice(1);
      typePicker.append(option);
    }
    typePicker.value=type;
    typePicker.addEventListener('change',()=>commit(emptyFor(typePicker.value)));
    wrapper.append(typePicker);
    if (type==='list') {
      const list=document.createElement('div');list.className='fm-items';
      value.forEach((item,index)=>{
        const line=document.createElement('div');line.className='fm-nested-item';
        const heading=document.createElement('span');heading.textContent=String(index+1)+'.';
        line.append(heading,structured(item,next=>{const copy=value.slice();copy[index]=next;commit(copy);},depth+1),
          nodeButton('Remove',()=>commit(value.filter((_,i)=>i!==index)),'Remove item'));
        list.append(line);
      });
      wrapper.append(list,nodeButton('+ Add item',()=>commit([...value,''])));
    } else if (type==='object') {
      const list=document.createElement('div');list.className='fm-items';
      for(const [key,item] of Object.entries(value)){
        const line=document.createElement('div');line.className='fm-nested-item';
        const heading=document.createElement('strong');heading.textContent=key;
        line.append(heading,structured(item,next=>commit({...value,[key]:next}),depth+1),
          nodeButton('Remove',()=>{const copy={...value};delete copy[key];commit(copy);}));
        list.append(line);
      }
      wrapper.append(list);
      const add=document.createElement('div');add.className='fm-add-nested';
      const name=document.createElement('input');name.placeholder='New property name';name.setAttribute('aria-label','Nested property name');
      add.append(name,nodeButton('+ Add property',()=>{
        const key=name.value.trim();
        if(!key||Object.hasOwn(value,key)){showMessage('Enter a unique property name.',true);return;}
        commit({...value,[key]:''});
      }));
      wrapper.append(add);
    } else if (type!=='null') {
      wrapper.append(inputFor(value,type,commit));
    }
    return wrapper;
  }

  function makeField(field) {
    const def=known.get(field.name);
    const type=def?.type || field.kind;
    const row=document.createElement('div');row.className='fm-field';row.dataset.field=field.name;
    const head=document.createElement('div');head.className='fm-field-head';
    const label=document.createElement('label');
    label.textContent=def?.label||field.name;
    label.id='fm-label-'+(++fieldId);
    const key=document.createElement('code');key.textContent=field.name;
    const remove=nodeButton('Remove',()=>setField(field.name,null,true),'Remove '+(def?.label||field.name)+' from this page');
    remove.setAttribute('aria-label','Remove '+(def?.label||field.name));
    remove.className='fm-remove';
    head.append(label,key,remove);
    row.append(head);
    if (def?.help) {
      const hint=document.createElement('small');hint.textContent=def.help;row.append(hint);
    }
    if (field.kind==='alias') {
      const warning=document.createElement('p');warning.textContent='YAML alias: edit this value in source mode to preserve anchors.';
      row.append(warning,nodeButton('Edit YAML',()=>setRaw(true)));
    } else if (type==='boolean' && field.kind==='boolean') {
      row.append(inputFor(field.value,'boolean',value=>setField(field.name,value)));
    } else if (type==='list' && (typeof field.value==='string' || Array.isArray(field.value) && field.value.every(item=>typeof item==='string'))) {
      const items=Array.isArray(field.value)?field.value:[field.value];
      const tags=inputFor(items.join(', '),'text',text=>setField(field.name,text.split(',').map(t=>t.trim()).filter(Boolean)));
      tags.placeholder='Separate with commas';row.append(tags);
    } else if (type==='object'||type==='list'||field.kind==='object'||field.kind==='list') {
      row.append(structured(field.value,value=>setField(field.name,value)));
    } else {
      row.append(inputFor(field.value,field.kind==='number'?'number':type==='textarea'?'textarea':'text',value=>setField(field.name,value)));
    }
    row.setAttribute('role','group');
    row.setAttribute('aria-labelledby',label.id);
    const primaryControl=row.querySelector('input,select,textarea');
    if(primaryControl&&!row.querySelector('.fm-structured')){
      primaryControl.id='fm-control-'+fieldId;
      label.htmlFor=primaryControl.id;
    }
    return row;
  }

  function renderFields(fields) {
    fieldsRoot.replaceChildren();
    const byName=new Map(fields.map(f=>[f.name,f]));
    const shown=new Set();
    const groups=['Page basics','Presentation','Publishing & navigation','Advanced','Custom'];
    const advancedCount=definitions.filter(def=>def[3]==='Advanced'&&byName.has(def[0])).length;
    for(const group of groups){
      const choices=[];
      for(const def of definitions) {
        if(def[3]!==group) continue;
        if(byName.has(def[0]) || (common.has(def[0]) && group==='Page basics')) {
          const field=byName.get(def[0]);
          // Unset fields are provided through the add menu, not silently created.
          if(field) {choices.push(field);shown.add(def[0]);}
        }
      }
      if(group==='Custom') {
        for(const field of fields) if(!shown.has(field.name)) choices.push(field);
      }
      if(!choices.length) continue;
      if(group==='Advanced'&&!showMore) {
        continue;
      }
      const groupEl=document.createElement('section');groupEl.className='fm-group';
      const heading=document.createElement('h3');heading.textContent=group;groupEl.append(heading);
      for(const field of choices) groupEl.append(makeField(field));
      fieldsRoot.append(groupEl);
    }
    if(!fields.length) {
      const empty=document.createElement('p');empty.className='fm-empty';empty.textContent='No properties yet. Add a title, tags, or any custom field below.';
      fieldsRoot.append(empty);
    }
    const published=byName.get('published')?.value, draft=byName.get('draft')?.value;
    const flag=document.querySelector('#fm-status');
    flag.hidden=false;
    flag.textContent=draft===true&&published===true?'Conflicting status':
      draft===true?'Draft':published===true?'Published':
      published===false?'Not listed':'No publish status';
    flag.dataset.warning=String(draft===true&&published===true);
    document.querySelector('#fm-summary-title').textContent=String(byName.get('title')?.value||'Untitled page');
    const slug=byName.get('slug')?.value;
    const activeFile=document.querySelector('#active-file').textContent;
    document.querySelector('#fm-summary-location').textContent=slug?'URL slug · '+String(slug):'File · '+activeFile;
    document.querySelector('#fm-save-hint').textContent=activeFile==='draft.md'?'Scratch drafts stay in this browser · Ctrl/Cmd+Z to undo':'Site pages save as you edit · Ctrl/Cmd+Z to undo';
    renderAddOptions(fields);
    moreButton.hidden=advancedCount===0;
    moreButton.textContent=showMore?'Hide advanced properties':'Show advanced properties'+(advancedCount?' ('+advancedCount+')':'');
  }

  function renderAddOptions(fields) {
    const byName=new Set(fields.map(field=>field.name));
    const term=addSearch.value.trim().toLowerCase();
    const selected=addSelect.value;
    addSelect.replaceChildren();
    const placeholder=document.createElement('option');placeholder.value='';placeholder.textContent='Choose a property…';addSelect.append(placeholder);
    for(const group of ['Page basics','Presentation','Publishing & navigation','Advanced']) {
      const available=definitions.filter(def=>def[3]===group&&!byName.has(def[0])&&
        (!term||[def[0],def[1],def[3]].some(value=>value.toLowerCase().includes(term))));
      if(!available.length)continue;
      const optionGroup=document.createElement('optgroup');optionGroup.label=group;
      for(const def of available){
        const option=document.createElement('option');option.value=def[0];option.textContent=def[1]+' · '+def[0];optionGroup.append(option);
      }
      addSelect.append(optionGroup);
    }
    if(!addSelect.querySelector('optgroup')){
      placeholder.textContent=term?'No matching properties':'All known properties are already on this page';
      placeholder.disabled=true;
    }else if(selected&&addSelect.querySelector('option[value="'+selected+'"]'))addSelect.value=selected;
  }

  async function refresh() {
    if (!expanded||!editor||busy) return;
    const source=editor.getValue(), serial=++refreshing;
    try {
      const data=await request('inspect',{source});
      if(serial!==refreshing||source!==editor.getValue()||!expanded) return;
      seenFields=data.fields || [];
      renderFields(seenFields);
      showMessage('');
    } catch(error) {
      if(serial!==refreshing||!expanded) return;
      showMessage(error.message+' · Use Edit YAML to repair it.',true);
      fieldsRoot.replaceChildren();
    }
  }
  function scheduleRefresh() {
    if(!expanded||busy) return;
    clearTimeout(refreshTimer);
    refreshTimer=setTimeout(refresh,200);
  }
  function applySource(oldSource,newSource) {
    if(oldSource===newSource)return;
    let left=0,right=0;
    while(left<oldSource.length&&left<newSource.length&&oldSource[left]===newSource[left]) left++;
    while(right<oldSource.length-left&&right<newSource.length-left&&
      oldSource[oldSource.length-1-right]===newSource[newSource.length-1-right]) right++;
    editor.replaceRange(newSource.slice(left,newSource.length-right),
      editor.posFromIndex(left),editor.posFromIndex(oldSource.length-right),'+frontmatter');
  }
  async function setField(name,value,remove=false) {
    if(busy)return;
    const source=editor.getValue();
    busy=true;
    panel.setAttribute('aria-busy','true');
    fieldsRoot.querySelectorAll('input,button,select,textarea').forEach(el=>el.disabled=true);
    showMessage('Updating page details…');
    let failure='';
    try {
      const data=await request('edit',{source,field:name,value,remove});
      if(source!==editor.getValue())throw new Error('The document changed while editing. Your newer edits were kept; try again.');
      applySource(source,data.source);
      showMessage('');
    } catch(error) {failure=error.message;}
    finally {
      busy=false;panel.removeAttribute('aria-busy');await refresh();
      if(failure)showMessage(failure,true);
    }
  }

  function addField(name,type) {
    if(!/^[A-Za-z_][A-Za-z0-9_-]{0,127}$/.test(name)) {showMessage('Use letters, numbers, _ or - for a field name.',true);return;}
    if(seenFields.some(field=>field.name===name)) {showMessage('This property already exists.',true);return;}
    newKey.value='';addSelect.value='';addSearch.value='';customField.open=false;addFlow.hidden=true;
    addToggle.setAttribute('aria-expanded','false');
    addToggle.textContent='+ Add property';
    if(known.get(name)?.group==='Advanced')showMore=true;
    setField(name,emptyFor(type==='textarea'?'text':type)).then(()=>{
      const row=[...fieldsRoot.querySelectorAll('.fm-field')].find(field=>field.dataset.field===name);
      row?.querySelector('input,select,textarea,button')?.focus();
    });
  }

  function attach(instance) {
    editor=instance;
    collapseYAML();
    editor.on('change',()=>{
      // Collapse after CodeMirror has finished recording the change.
      setTimeout(collapseYAML,0);
      scheduleRefresh();
    });
    toggle.addEventListener('click',toggleProperties);
    rawButton.addEventListener('click',()=>setRaw(!raw));
    moreButton.addEventListener('click',()=>{
      showMore=!showMore;
      renderFields(seenFields);
    });
    addToggle.addEventListener('click',()=>{
      const opening=addFlow.hidden;
      addFlow.hidden=!opening;
      addToggle.setAttribute('aria-expanded',String(opening));
      addToggle.textContent=opening?'Close property picker':'+ Add property';
      if(!addFlow.hidden)addSearch.focus();
    });
    addSearch.addEventListener('input',()=>renderAddOptions(seenFields));
    addKnown.addEventListener('click',()=>{
      const name=addSelect.value;
      if(!name){showMessage('Choose a property to add.',true);addSelect.focus();return;}
      addField(name,known.get(name)?.type||'text');
    });
    addCustom.addEventListener('click',()=>{
      const name=newKey.value.trim();
      if(!name){showMessage('Enter a name for the custom property.',true);newKey.focus();return;}
      addField(name,document.querySelector('#fm-new-type').value);
    });
    document.addEventListener('keydown',event=>{
      if(event.isComposing||event.repeat||event.altKey===false||!event.shiftKey||(event.key.toLowerCase()!=='m')||!(event.metaKey||event.ctrlKey))return;
      event.preventDefault();toggleProperties();
    });
  }

  window.MarkataFrontmatter = {attach,open:openProperties,editYAML:()=>setRaw(true),close:closeProperties};
})();
