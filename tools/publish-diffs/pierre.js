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

async function render(host, before, after, options = {}) {
  if (!host || typeof before !== 'string' || typeof after !== 'string')
    throw new Error('Invalid source comparison');
  release(host);
  const file = options.filename || 'post.md';
  let mounted;
  let resolveMounted;
  let rejectMounted;
  mounted = new Promise((resolve, reject) => {
    resolveMounted = resolve;
    rejectMounted = reject;
  });
  const view = new FileDiff({
    diffStyle: options.layout === 'split' ? 'split' : 'unified',
    theme: palettes,
    diffIndicators: 'bars',
    preferredHighlighter: 'shiki-js',
    // One standalone file in a nested dialog scroll container should not
    // allocate off-screen virtualization spacers for collapsed hunks.
    disableVirtualizationBuffers: true,
    hunkSeparators: 'line-info-basic',
    // The caller must not report success until Pierre has committed its DOM.
    disableErrorHandling: true,
    onPostRender(_node, _instance, phase) {
      if (phase === 'mount') resolveMounted();
    }
  });
  const timeout = setTimeout(() => rejectMounted(new Error('Pierre did not finish rendering the diff.')), 15000);
  try {
    view.render({
      fileContainer: host,
      oldFile: options.newFile ? null : { name:file, contents:before },
      newFile: { name:file, contents:after }
    });
    await mounted;
  } finally {
    clearTimeout(timeout);
  }
  instances.set(host, view);
  return view;
}

window.MarkataPierreDiff = {render,release};
