/* Shared by the Go, Python and offline WASM workbenches. */
(() => {
  'use strict';
  const commands = window.MarkataMarkdownCommands;

  window.MarkataCommandPalette = {attach};
  function attach(editor) {
    const button = document.createElement('button');
    button.type = 'button'; button.id = 'commands';
    button.textContent = 'Commands'; button.title = 'Commands (Ctrl/Cmd+K)';
    button.setAttribute('aria-haspopup', 'dialog');
    document.querySelector('.editor-toolbar').append(button);
    const dialog = document.createElement('dialog');
    dialog.className = 'markdown-commands'; dialog.setAttribute('aria-labelledby', 'commands-title');
    dialog.innerHTML = '<div class="commands-heading"><h2 id="commands-title">Insert Markdown</h2><button type="button" class="commands-close" aria-label="Close commands">Close</button></div><label for="commands-query">Find a command</label><input id="commands-query" type="search" placeholder="Try callout, diagram, or highlight" autocomplete="off" role="combobox" aria-autocomplete="list" aria-expanded="true" aria-controls="commands-results"><p id="commands-count" role="status"></p><div id="commands-results" role="listbox" aria-label="Markdown commands"></div><p class="commands-help">↑ ↓ Choose · Enter Insert · Esc Close</p>';
    document.body.append(dialog);
    const query = dialog.querySelector('input');
    const results = dialog.querySelector('[role=listbox]');
    const count = dialog.querySelector('[role=status]');
    const available = window.MarkataSiteTree ? [{id:'open-file', label:'Open file', description:'Fuzzy search workspace files', syntax:'Ctrl/Cmd+P', aliases:'find quick open picker', action:() => window.MarkataSiteTree.openPicker()}, ...commands] : commands;
    let matches = available, active = 0, selection;
    function select(index) {
      active = index;
      [...results.children].forEach((row, i) => row.setAttribute('aria-selected', String(i === active)));
      const row = results.children[active];
      if (row) {
        query.setAttribute('aria-activedescendant', row.id);
        row.scrollIntoView({block:'nearest'});
      } else query.removeAttribute('aria-activedescendant');
    }
    function refresh() {
      const words = query.value.toLowerCase().trim().split(/\s+/).filter(Boolean);
      const filtered = available.filter(c => words.every(word => `${c.label} ${c.description} ${c.syntax} ${c.aliases}`.toLowerCase().includes(word)));
      matches = filtered.slice(0, 60);
      results.replaceChildren();
      count.textContent = filtered.length > 60 ? `Showing 60 of ${filtered.length} commands. Refine your search to see more.` : filtered.length ? `${filtered.length} commands` : 'No matching commands. Try another word or clear the search.';
      matches.forEach((command, index) => {
        const row = document.createElement('div');
        row.id = `markdown-command-${command.id}`; row.setAttribute('role','option');
        const label = document.createElement('strong'); label.textContent = command.label;
        const detail = document.createElement('span'); detail.textContent = command.description;
        const syntax = document.createElement('code'); syntax.textContent = command.syntax;
        row.append(label, detail, syntax);
        row.addEventListener('pointermove', () => select(index));
        row.addEventListener('click', () => insert(command));
        results.append(row);
      });
      select(0);
    }
    function open() {
      if (dialog.open) return;
      selection = {from:editor.getCursor('from'), to:editor.getCursor('to'), text:editor.getSelection()};
      query.value = ''; dialog.showModal(); refresh(); query.focus();
    }
    function close() { dialog.close(); }
    function insert(command) {
      if (command.action) { close(); command.action(); return; }
      if (editor.getOption('readOnly')) return;
      const text = command.useSelection === false ? command.placeholder : selection.text || command.placeholder;
      // Keep all lines of a selected passage inside an indented block.
      const indent = command.template.slice(0, command.template.indexOf('{{text}}')).match(/\n( +)$/)?.[1] || '';
      const content = text.replace(/\n/g, '\n' + indent);
      const marker = command.template.indexOf('{{text}}');
      let snippet = command.template.replace('{{text}}', () => content);
      let prefix = '', suffix = '';
      const from = editor.indexFromPos(selection.from), to = editor.indexFromPos(selection.to);
      const source = editor.getValue();
      if (command.block) {
        const before = source.slice(0, from), after = source.slice(to);
        if (before) prefix = before.endsWith('\n\n') ? '' : before.endsWith('\n') ? '\n' : '\n\n';
        if (after) suffix = after.startsWith('\n\n') ? '' : after.startsWith('\n') ? '\n' : '\n\n';
      }
      snippet = prefix + snippet + suffix;
      editor.operation(() => {
        // Separate palette insertion from the user's preceding/following typing.
        editor.changeGeneration(true);
        editor.replaceRange(snippet, selection.from, selection.to, '+command');
        editor.setSelection(editor.posFromIndex(from + prefix.length + marker), editor.posFromIndex(from + prefix.length + marker + content.length));
        editor.changeGeneration(true);
      });
      close();
    }
    button.addEventListener('click', open);
    query.addEventListener('input', refresh);
    dialog.querySelector('.commands-close').addEventListener('click', close);
    dialog.addEventListener('click', event => {
      const box = dialog.getBoundingClientRect();
      if (event.target === dialog && (event.clientX < box.left || event.clientX > box.right || event.clientY < box.top || event.clientY > box.bottom)) close();
    });
    dialog.addEventListener('close', () => editor.focus());
    dialog.addEventListener('keydown', event => {
      if (event.isComposing) return;
      if (event.key === 'Escape') { event.preventDefault(); event.stopPropagation(); close(); return; }
      if (['ArrowDown', 'ArrowUp', 'Enter'].includes(event.key) && event.target === query) {
        event.preventDefault();
        if (event.key === 'Enter') { if (matches[active]) insert(matches[active]); }
        else if (matches.length) select((active + (event.key === 'ArrowDown' ? 1 : -1) + matches.length) % matches.length);
      }
    });
    // Capture before Vim, Emacs or browser shortcuts can consume the event.
    document.addEventListener('keydown', event => {
      const key = event.key.toLowerCase();
      const modifier = event.ctrlKey || event.metaKey;
      const shortcut = (key === 'k' && !event.shiftKey) || (key === 'p' && event.shiftKey);
      if (!event.isComposing && modifier && shortcut && !event.altKey) {
        event.preventDefault(); event.stopPropagation(); open();
      }
    }, true);
  }
})();
