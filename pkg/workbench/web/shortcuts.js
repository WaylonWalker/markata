/* App-wide keyboard commands. Editing keymaps remain owned by CodeMirror. */
(() => {
  'use strict';
  const commands = new Map();
  const storageKey = 'markata-app-shortcuts-v1';
  const presets = ['writer', 'vscode', 'vim', 'emacs'];
  const persisted = (() => {
    try { return JSON.parse(localStorage.getItem(storageKey) || '{}'); }
    catch { return {}; }
  })();
  let preset = presets.includes(persisted.preset) ? persisted.preset : 'vscode';
  let overridesByPreset = persisted.overridesByPreset && typeof persisted.overridesByPreset === 'object' ? persisted.overridesByPreset : {};
  let overrides = overridesByPreset[preset] || {};
  let editor, pending = '', pendingTimer, dialog, search, list, feedback, returnFocus;

  function persist() {
    try { localStorage.setItem(storageKey, JSON.stringify({preset, overridesByPreset})); }
    catch { /* Shortcut handling works with blocked storage. */ }
  }
  function register(command) {
    if (!command || !command.id || typeof command.run !== 'function' || commands.has(command.id))
      throw new Error('Invalid or duplicate Markata command: ' + (command && command.id));
    commands.set(command.id, command);
    return command;
  }
  function bindingFor(command) {
    const id = typeof command === 'string' ? command : command.id;
    const entry = commands.get(id);
    if (!entry) return '';
    if (Object.prototype.hasOwnProperty.call(overrides, id)) return overrides[id] || '';
    const bindings = entry.bindings || {};
    return bindings[preset] === undefined ? (bindings.vscode || '') : bindings[preset];
  }
  function conflict(binding, except) {
    if (!binding) return null;
    return [...commands.values()].find(command => command.id !== except && bindingFor(command) === binding) || null;
  }
  function setBinding(id, binding) {
    if (!commands.has(id) || (binding && !/^((Mod|Alt|Shift|Ctrl)\+)*[A-Za-z0-9\\]+(?: ((Mod|Alt|Shift|Ctrl)\+)*[A-Za-z0-9\\]+)*$/.test(binding)))
      return {ok:false, error:'Unrecognized shortcut.'};
    const other = conflict(binding, id);
    if (other) return {ok:false, error:'Already used by ' + other.label + '.'};
    overrides = {...overrides, [id]:binding};
    overridesByPreset = {...overridesByPreset, [preset]:overrides};
    persist();
    render();
    return {ok:true};
  }
  function reset() { overrides = {}; overridesByPreset = {...overridesByPreset,[preset]:{}}; persist(); render(); }
  function setPreset(value) {
    if (!presets.includes(value)) return;
    preset = value;
    overrides = overridesByPreset[preset] || {};
    clearPending();
    persist();
    render();
  }
  function execute(id) {
    const command = commands.get(id);
    if (!command || (command.when && !command.when())) return false;
    command.run();
    return true;
  }
  function stroke(event) {
    let key = event.key;
    if (key === ' ') key = 'Space';
    else if (key === 'Escape') key = 'Esc';
    else if (key === 'ArrowLeft') key = 'Left';
    else if (key === 'ArrowRight') key = 'Right';
    else if (key.length === 1) key = key.toUpperCase();
    const modifier = preset === 'emacs' && event.ctrlKey && !event.metaKey ? 'Ctrl' :
      event.ctrlKey || event.metaKey ? 'Mod' : '';
    const shifted = event.shiftKey && key !== 'Shift' && /^[A-Za-z0-9]$/.test(key);
    return [modifier, event.altKey ? 'Alt' : '', shifted ? 'Shift' : '', key].filter(Boolean).join('+');
  }
  function clearPending() {
    pending = '';
    if (pendingTimer) clearTimeout(pendingTimer);
    pendingTimer = undefined;
  }
  function dispatch(event) {
    if (event.defaultPrevented || event.isComposing || event.key === 'Process' || event.repeat) return false;
    if (document.querySelector('dialog[open]')) { clearPending(); return false; }
    const target = event.target;
    const textInput = target && (target.matches?.('input, textarea, select, [contenteditable="true"]') ||
      target.closest?.('[contenteditable="true"]'));
    const inEditor = Boolean(target && target.closest?.('.CodeMirror'));
    const isVimNormal = preset === 'vim' && inEditor && editor &&
      editor.getOption('keyMap') === 'vim' && !editor.state.vim?.insertMode;
    const key = stroke(event);
    if (textInput && !inEditor) { clearPending(); return false; }
    if (key === 'Space' && !isVimNormal) { clearPending(); return false; }
    let sequence = pending ? pending + ' ' + key : key;
    let options = [...commands.values()].filter(command =>
      (!command.when || command.when()) && bindingFor(command));
    let match = options.find(command => bindingFor(command) === sequence);
    if (!match && pending) {
      clearPending();
      sequence = key;
      options = options.filter(command => bindingFor(command).startsWith(sequence));
      match = options.find(command => bindingFor(command) === sequence);
    }
    if (match) {
      clearPending();
      event.preventDefault(); event.stopPropagation();
      execute(match.id);
      return true;
    }
    if (options.some(command => bindingFor(command).startsWith(sequence + ' '))) {
      event.preventDefault(); event.stopPropagation();
      pending = sequence;
      pendingTimer = setTimeout(clearPending, 1400);
      return true;
    }
    clearPending();
    return false;
  }
  function render() {
    if (!dialog || !dialog.open) return;
    dialog.querySelector('#shortcuts-preset').value = preset;
    const query = search.value.toLowerCase().trim();
    list.replaceChildren();
    const entries = [...commands.values()].filter(command =>
      (command.label + ' ' + command.id + ' ' + (command.category || '')).toLowerCase().includes(query));
    for (const command of entries) {
      const row = document.createElement('div'); row.className = 'shortcut-row';
      const name = document.createElement('div');
      const title = document.createElement('strong'); title.textContent = command.label;
      const category = document.createElement('small'); category.textContent = command.category || 'General';
      name.append(title, category);
      const record = document.createElement('button'); record.type = 'button';
      record.className = 'shortcut-binding';
      record.textContent = bindingFor(command) || 'Unbound';
      record.title = 'Change shortcut for ' + command.label;
      record.addEventListener('click', () => capture(command, record));
      row.append(name, record); list.append(row);
    }
    dialog.querySelector('#shortcuts-count').textContent = entries.length + ' commands';
  }
  function capture(command, button) {
    const input = document.createElement('input');
    input.className = 'shortcut-capture'; input.readOnly = true;
    input.setAttribute('aria-label', 'Press a shortcut for ' + command.label);
    input.placeholder = 'Press keys · Esc cancels';
    button.replaceWith(input); input.focus();
    let first = '';
    function finish() { if (input.isConnected) input.replaceWith(button); render(); }
    input.addEventListener('blur', finish);
    input.addEventListener('keydown', event => {
      if (event.key === 'Tab') { finish(); return; }
      event.preventDefault(); event.stopPropagation();
      if (event.key === 'Escape') { finish(); return; }
      if (event.key === 'Backspace') {
        const result = setBinding(command.id, '');
        feedback.textContent = result.ok ? 'Shortcut removed.' : result.error;
        finish(); return;
      }
      if (['Control','Meta','Alt','Shift','Tab'].includes(event.key)) return;
      const next = stroke(event);
      if (!first && ['Mod+K', 'Mod+X', 'Alt+X', 'Space'].includes(next)) {
        first = next; input.placeholder = next + ' … next key'; return;
      }
      const sequence = first ? first + ' ' + next : next;
      const result = setBinding(command.id, sequence);
      feedback.textContent = result.ok ? command.label + ': ' + sequence : result.error;
      if (result.ok) finish();
      else { first = ''; input.placeholder = 'Try another shortcut'; }
    });
  }
  function openSettings() {
    if (!dialog) buildSettings();
    const current = document.querySelector('dialog[open]');
    if (current && current !== dialog) current.close();
    returnFocus = document.activeElement;
    dialog.showModal(); search.value = ''; feedback.textContent = '';
    render(); search.focus();
  }
  function buildSettings() {
    dialog = document.createElement('dialog'); dialog.id = 'shortcuts-dialog';
    dialog.setAttribute('aria-labelledby', 'shortcuts-heading');
    dialog.innerHTML = '<div class="shortcuts-header"><div><h2 id="shortcuts-heading">Keyboard shortcuts</h2><p>Application keybindings are separate from your editor text keymap.</p></div><button id="shortcuts-close" type="button" aria-label="Close keyboard shortcuts">×</button></div>' +
      '<div class="shortcuts-controls"><label for="shortcuts-preset">Application keymap</label><select id="shortcuts-preset"><option value="writer">Writer</option><option value="vscode">VS Code</option><option value="vim">Vim</option><option value="emacs">Emacs</option></select>' +
      '<label for="shortcuts-search">Find a command</label><input id="shortcuts-search" type="search" placeholder="Search commands and categories"></div>' +
      '<div class="shortcuts-meta"><span id="shortcuts-count"></span><button id="shortcuts-reset" type="button">Reset overrides</button></div>' +
      '<div id="shortcuts-list"></div><p id="shortcuts-feedback" role="status" aria-live="polite"></p>';
    document.body.append(dialog);
    search = dialog.querySelector('#shortcuts-search'); list = dialog.querySelector('#shortcuts-list');
    feedback = dialog.querySelector('#shortcuts-feedback');
    search.addEventListener('input', render);
    dialog.querySelector('#shortcuts-preset').addEventListener('change', event => setPreset(event.target.value));
    dialog.querySelector('#shortcuts-reset').addEventListener('click', reset);
    dialog.querySelector('#shortcuts-close').addEventListener('click', () => dialog.close());
    dialog.addEventListener('close', () => {
      const closedParent = returnFocus?.closest?.('dialog');
      if (closedParent && !closedParent.open) document.querySelector('#settings-open')?.focus();
      else if (returnFocus?.isConnected) returnFocus.focus();
      else editor?.focus();
    });
  }
  function init(instance) {
    editor = instance;
    document.addEventListener('keydown', dispatch, true);
    const settings = document.querySelector('#settings-dialog .settings-options');
    if (settings) {
      const button = document.createElement('button'); button.type = 'button';
      button.className = 'shortcuts-settings-link'; button.textContent = 'Keyboard shortcuts…';
      button.addEventListener('click', openSettings); settings.append(button);
    }
  }
  window.MarkataShortcuts = {register, execute, bindingFor, setBinding, setPreset,
    getPreset:() => preset, reset, all:() => [...commands.values()], init, dispatch, openSettings};
})();
