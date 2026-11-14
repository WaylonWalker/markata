/* One settings surface; reuse existing controls and their persistence handlers. */
(() => {
  function attach(editor) {
    const dialog = document.querySelector('#settings-dialog');
    const options = dialog.querySelector('.settings-options');
    const controls = Object.fromEntries(['setting-wrap','setting-scroll-sync','setting-layout','keymap']
      .map(id => [id, document.getElementById(id)]));
    dialog.setAttribute('aria-labelledby','editor-settings-title');
    dialog.querySelector('h2').id = 'editor-settings-title';
    dialog.querySelector('h2').textContent = 'Settings';
    dialog.querySelector('header p').textContent = 'Make this writing space yours. Saved in this browser.';
    options.replaceChildren();
    const search = document.createElement('input');
    search.type = 'search'; search.id = 'settings-search'; search.placeholder = 'Search settings…';
    search.setAttribute('aria-label','Search settings');
    const body = document.createElement('div'); body.className = 'settings-body';
    const nav = document.createElement('nav'); nav.setAttribute('aria-label','Settings categories');
    const content = document.createElement('div'); content.className = 'settings-content';
    const empty = document.createElement('p'); empty.textContent = 'No matching settings.'; empty.hidden = true;
    content.append(empty); body.append(nav,content); options.append(search,body);
    const categories = new Map(); let current = 'writing', returnFocus;
    for (const [id,label,description] of [
      ['appearance','Appearance','Colors and text size for this writing space'],
      ['writing','Writing','Text editing and reading comfort'],
      ['preview','Preview','Arrange your editor and live preview'],
      ['shortcuts','Keyboard Shortcuts','Choose familiar keys and customize commands'],
      ['advanced','Advanced','Search history and performance diagnostics']]) {
      const button = document.createElement('button'); button.type = 'button'; button.textContent = label;
      button.dataset.category = id; nav.append(button);
      const section = document.createElement('section'); section.dataset.settingsCategory = id;
      const heading = document.createElement('h3'); heading.textContent = label;
      const subtitle = document.createElement('p'); subtitle.textContent = description;
      section.append(heading,subtitle); content.append(section); categories.set(id,{button,section});
      button.addEventListener('click',() => {current=id; search.value=''; filter();});
    }
    function row(category, title, description, control) {
      const node = document.createElement('div'); node.className = 'setting-row';
      const copy = document.createElement('div'); const label = document.createElement('label');
      label.textContent = title; label.htmlFor = control.id;
      const hint = document.createElement('p'); hint.textContent = description;
      copy.append(label,hint); node.append(copy,control); categories.get(category).section.append(node);
      return node;
    }
    row('writing','Wrap long lines','Keep long paragraphs within the source pane.',controls['setting-wrap']);
    row('writing','Editor keymap','Vim, Emacs and Sublime text editing. Application shortcuts have their own category.',controls.keymap);
    const font = document.createElement('input'); font.id='setting-font-size'; font.type='range'; font.min='12'; font.max='24'; font.step='1';
    try { font.value=localStorage.getItem('markata-editor-font-size') || '14'; } catch { font.value='14'; }
    const applyFont = () => {
      editor.getWrapperElement().style.fontSize=font.value+'px'; editor.refresh();
      font.setAttribute('aria-valuetext',font.value+' pixels');
    };
    font.addEventListener('input',() => {applyFont();try {localStorage.setItem('markata-editor-font-size',font.value);} catch {}});
    applyFont(); row('appearance','Source text size','Adjust the source text size from 12 to 24 pixels.',font);
    const appearance=window.markata?.paletteSwitcher;
    const pickerHost=document.createElement('div'); pickerHost.id='settings-theme-picker';
    const pickerNode=document.querySelector('.theme-picker');
    if(appearance && pickerNode) pickerHost.append(pickerNode);
    categories.get('appearance').section.append(pickerHost);
    // The picker closes on any document click outside its root. Keep clicks inside this dialog from reaching it.
    dialog.addEventListener('click',(event) => {if(pickerNode && !pickerNode.contains(event.target)) event.stopPropagation();});
    function openPicker() {
      if(!appearance || !pickerNode) return;
      if(pickerNode.querySelector('.theme-picker-panel')?.hidden) appearance.open('colors');
    }
    row('preview','Layout','Write beside the preview, or focus on one pane.',controls['setting-layout']);
    row('preview','Synchronize scrolling','Keep the preview aligned with the paragraph you are editing.',controls['setting-scroll-sync']);
    const shortcuts = document.createElement('div'); shortcuts.id='settings-shortcuts'; categories.get('shortcuts').section.append(shortcuts);
    const history = document.createElement('input'); history.type='checkbox'; history.id='setting-pick-history'; history.checked=window.MarkataSearch?.enabled() !== false;
    history.addEventListener('change',() => window.MarkataSearch?.setEnabled(history.checked));
    row('advanced','Remember recent picks','Use local command and workspace file history to rank similar matches.',history);
    const clear = document.createElement('button'); clear.type='button'; clear.id='settings-clear-history'; clear.textContent='Clear history';
    const feedback = document.createElement('p'); feedback.setAttribute('role','status');
    clear.addEventListener('click',() => {
      window.MarkataSearch?.clearHistory('commands');
      window.MarkataSearch?.clearHistory(window.MarkataSearch.fileScope());
      window.dispatchEvent(new Event('markata-history-cleared'));
      feedback.textContent='Command and current workspace file history cleared.';
    });
    row('advanced','Clear recent picks','Remove command history and this workspace’s recent files.',clear).append(feedback);
    const diagnostics = document.createElement('button'); diagnostics.id='settings-diagnostics'; diagnostics.type='button'; diagnostics.textContent='Open diagnostics';
    diagnostics.addEventListener('click',() => {dialog.close();window.MarkataDiagnostics?.open();});
    row('advanced','Performance diagnostics','Inspect timings and copy a report for troubleshooting.',diagnostics);
    function filter() {
      const query=search.value.trim().toLowerCase(); let count=0;
      for (const [id,{button,section}] of categories) {
        button.setAttribute('aria-current', !query && id===current ? 'page' : 'false');
        const categoryMatch=(section.querySelector('h3').textContent+' '+section.querySelector('p').textContent).toLowerCase().includes(query);
        let matches=0;
        for (const item of section.querySelectorAll('.setting-row')) {
          item.hidden=Boolean(query && !categoryMatch && !item.textContent.toLowerCase().includes(query));
          if (!item.hidden) matches++;
        }
        if (id==='shortcuts') matches=categoryMatch ? 1 : 0;
        section.hidden=query ? !matches : id!==current;
        if (!section.hidden) count++;
      }
      empty.hidden=count>0;
      if (!categories.get('shortcuts').section.hidden) window.MarkataShortcuts?.renderSettings?.();
      if (!categories.get('appearance').section.hidden) openPicker();
    }
    search.addEventListener('input',filter);
    dialog.addEventListener('close',() => {if(returnFocus?.isConnected)returnFocus.focus();else editor.focus();});
    function open(category='writing') {
      const other = document.querySelector('dialog[open]');
      if (other && other !== dialog) other.close();
      if (!dialog.open) returnFocus = document.activeElement;
      current=categories.has(category)?category:'writing'; search.value='';
      history.checked=window.MarkataSearch?.enabled() !== false;
      if(!dialog.open)dialog.showModal(); filter(); search.focus();
    }
    window.MarkataSettings={open};
    window.MarkataShortcuts?.mountSettings?.();
    const registry=window.MarkataShortcuts;
    registry.register({id:'workbench.settings',label:'Open Editor Settings',category:'Settings',bindings:{vscode:'Mod+,'},run:()=>open()});
    for(const [id,{button}] of categories) if (id !== 'shortcuts') registry.register({id:'workbench.settings.'+id,label:'Settings: '+button.textContent,category:'Settings',run:()=>open(id)});
    registry.register({id:'workbench.scrollSync',label:'Toggle Synchronized Scrolling',category:'Preview',aliases:'scroll sync',run:()=>{const c=controls['setting-scroll-sync'];c.checked=!c.checked;c.dispatchEvent(new Event('change'));}});
    for(const [value,label] of [['split','Editor and Preview'],['editor','Editor Only'],['preview','Preview Only']]) registry.register({id:'workbench.layout.'+value,label:'Layout: '+label,category:'Layout',run:()=>{controls['setting-layout'].value=value;controls['setting-layout'].dispatchEvent(new Event('change'));}});
    for(const [value,label] of [['files','Files'],['feeds','Feeds'],['tags','Tags']]) registry.register({id:'workbench.view.'+value,label:'Browse '+label,category:'Navigation',when:()=>window.MarkataSiteTree?.isReady?.() && !document.querySelector('#site-view option[value="'+value+'"]').disabled,run:()=>{registry.execute('workbench.focusExplorer');const select=document.querySelector('#site-view');select.value=value;select.dispatchEvent(new Event('change'));}});
    filter();
  }
  window.MarkataSettings={attach};
})();
