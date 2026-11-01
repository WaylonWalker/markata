// Vendored at build time; never loaded until the user chooses Full diff.
import { FileDiff } from '@pierre/diffs';

const instances = new WeakMap();
const palettes = { light: 'pierre-light', dark: 'pierre-dark' };

function release(host) {
  const old = instances.get(host);
  if (old) {
    old.cleanUp();
    instances.delete(host);
  }
  host.replaceChildren();
}

function render(host, before, after, options = {}) {
  if (!host || typeof before !== 'string' || typeof after !== 'string')
    throw new Error('Invalid source comparison');
  release(host);
  const file = options.filename || 'post.md';
  const view = new FileDiff({
    diffStyle: options.layout === 'split' ? 'split' : 'stacked',
    theme: palettes,
    diffIndicators: 'bars',
        preferredHighlighter: 'highlights'
  });
  view.render({
    fileContainer: host,
    oldFile: options.newFile ? null : { name:file, contents:before },
    newFile: { name:file, contents:after }
  });
  instances.set(host, view);
  return view;
}

window.MarkataPierreDiff = {render,release};
