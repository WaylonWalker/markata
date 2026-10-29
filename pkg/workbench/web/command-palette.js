/* Shared by the Go, Python and offline WASM workbenches. */
(() => {
  'use strict';
  // {{text}} is replaced with the selection, or the command's editable default.
  const commands = [
    ['note', 'Note box', 'Add a titled note.', '!!! note', 'callout admonition', '!!! note "Note"\n    {{text}}', 'Write your note here.', true],
    ['tip', 'Tip box', 'Add a helpful suggestion.', '!!! tip', 'callout admonition', '!!! tip "Tip"\n    {{text}}', 'Write your tip here.', true],
    ['warning', 'Warning box', 'Draw attention to a warning.', '!!! warning', 'callout admonition caution', '!!! warning "Warning"\n    {{text}}', 'Explain what needs attention.', true],
    ['aside', 'Aside', 'Add a side note.', '!!! aside', 'callout sidenote', '!!! aside\n    {{text}}', 'A related thought.', true],
    ['collapse', 'Collapsible note', 'Add a note readers can expand.', '??? note', 'details disclosure callout', '??? note "Click to expand"\n    {{text}}', 'Hidden content.', true],
    ['highlight', 'Highlighted text', 'Highlight the selected text.', '==text==', 'mark emphasis', '=={{text}}==', 'highlighted text'],
    ['keys', 'Keyboard keys', 'Display a keyboard shortcut.', '++Ctrl+Alt+Del++', 'kbd shortcut', '++{{text}}++', 'Ctrl+Alt+Del'],
    ['definition', 'Definition list', 'Explain a term.', 'Term / : Definition', 'glossary meaning', 'Term\n:   {{text}}', 'Definition', true],
    ['mermaid', 'Mermaid diagram', 'Start a flow diagram.', '```mermaid', 'flowchart graph', '```mermaid\n{{text}}\n```', 'graph LR\n  Idea --> Draft --> Publish', true],
    ['chart', 'Chart.js chart', 'Start a bar chart with sample data.', '```chartjs', 'chartjs chart.js plot bar data', '```chartjs\n{{text}}\n```', '{\n  "type": "bar",\n  "data": {\n    "labels": ["Draft", "Review", "Done"],\n    "datasets": [{"label": "Posts", "data": [3, 2, 5]}]\n  }\n}', true],
    ['table', 'Table', 'Insert a two-column Markdown table.', '| Header | Header |', 'grid rows columns', '| Column 1 | Column 2 |\n| --- | --- |\n| {{text}} | Value |', 'Text', true],
    ['task', 'Task list', 'Add tasks to check off.', '- [ ] Task', 'checklist todo', '- [ ] {{text}}', 'Task', true],
    ['code', 'Code block', 'Insert a fenced code example.', '```text', 'fence syntax source', '```text\n{{text}}\n```', 'Your code here', true],
    ['bold', 'Bold', 'Emphasize the selected text.', '**text**', 'strong formatting', '**{{text}}**', 'bold text'],
    ['italic', 'Italic', 'Italicize the selected text.', '*text*', 'emphasis formatting', '*{{text}}*', 'italic text'],
    ['link', 'Link', 'Link the selected text to a URL.', '[text](URL)', 'anchor hyperlink', '[{{text}}](https://example.com)', 'Link text'],
    ['image', 'Image URL', 'Insert an image with descriptive alt text.', '![alt](URL)', 'photo picture', '![{{text}}](https://example.com/image.jpg)', 'Describe the image']
  ].map(([id, label, description, syntax, aliases, template, placeholder, block]) =>
    ({id, label, description, syntax, aliases, template, placeholder, block}));

  window.MarkataCommandPalette = {attach};
  function attach(editor) {
    const button = document.createElement('button');
    button.type = 'button'; button.id = 'commands';
    button.textContent = 'Commands'; button.title = 'Commands (Ctrl/Cmd+Shift+P)';
    button.setAttribute('aria-haspopup', 'dialog');
    document.querySelector('.editor-toolbar').append(button);
    const dialog = document.createElement('dialog');
    dialog.className = 'markdown-commands'; dialog.setAttribute('aria-labelledby', 'commands-title');
    dialog.innerHTML = '<div class="commands-heading"><h2 id="commands-title">Insert Markdown</h2><button type="button" class="commands-close" aria-label="Close commands">Close</button></div><label for="commands-query">Find a command</label><input id="commands-query" type="search" placeholder="Try callout, diagram, or highlight" autocomplete="off" role="combobox" aria-autocomplete="list" aria-expanded="true" aria-controls="commands-results"><p id="commands-count" role="status"></p><div id="commands-results" role="listbox" aria-label="Markdown commands"></div><p class="commands-help">↑ ↓ Choose · Enter Insert · Esc Close</p>';
    document.body.append(dialog);
    const query = dialog.querySelector('input');
    const results = dialog.querySelector('[role=listbox]');
    const count = dialog.querySelector('[role=status]');
    let matches = commands, active = 0, selection;
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
      matches = commands.filter(c => words.every(word => `${c.label} ${c.description} ${c.syntax} ${c.aliases}`.toLowerCase().includes(word)));
      results.replaceChildren();
      count.textContent = matches.length ? `${matches.length} commands` : 'No matching commands. Try another word or clear the search.';
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
      const text = selection.text || command.placeholder;
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
      if (!event.isComposing && (event.ctrlKey || event.metaKey) && event.shiftKey && !event.altKey && event.key.toLowerCase() === 'p') {
        event.preventDefault(); event.stopPropagation(); open();
      }
    }, true);
  }
})();
