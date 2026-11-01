/* Progressive enhancement for article code blocks.
 * A fenced code block remains ordinary semantic <pre><code> without JS.
 * Loaded only on pages containing article code via conditional-css.js.
 */
(() => {
  'use strict';

  if (typeof window.markataInitCodeChrome === 'function') {
    window.markataInitCodeChrome();
    return;
  }

  function copyText(value) {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      return navigator.clipboard.writeText(value);
    }
    // Clipboard API is not available on some non-HTTPS local sites.
    // Do not use execCommand or a hidden focused textarea: leave a
    // clean manual-copy fallback and preserve existing text selection.
    return Promise.reject(new Error('Clipboard API unavailable'));
  }

  function initCodeChrome(root = document) {
    const roots = Array.isArray(root) ? root : [root];
    const blocks = roots.flatMap((node) => {
      if (node.matches?.('pre')) return [node];
      if (node === document) return [...document.querySelectorAll('.post-content pre')];
      return [...(node.querySelectorAll?.('pre') || [])];
    });
    blocks.forEach((pre) => {
      if (pre.dataset.markataCodeEnhanced === 'true' ||
          pre.closest('.mermaid, [data-no-code-chrome]') ||
          pre.classList.contains('mermaid')) return;

      const code = pre.querySelector(':scope > code');
      if (!code && !pre.matches('.chroma, .highlight') &&
          !pre.closest('.chroma, .highlight')) return;

      const raw = (code || pre).textContent;
      if (!raw || !raw.trim()) return;

      pre.dataset.markataCodeEnhanced = 'true';
      pre.classList.add('md-code-enhanced');

      const codeClass = (code && code.className) || '';
      const language = codeClass.match(/(?:^|\s)language-([a-z0-9_+.-]+)/i);
      if (language && !pre.hasAttribute('data-language')) {
        pre.setAttribute('data-language', language[1]);
        pre.dataset.markataGeneratedLanguage = 'true';
      }

      const button = document.createElement('button');
      button.type = 'button';
      button.className = 'md-code-copy';
      button.textContent = 'Copy';
      button.setAttribute('aria-label', 'Copy code block');
      button.setAttribute('data-markata-ui', 'code-copy');

      button.addEventListener('click', () => {
        copyText(raw.replace(/\n+$/, '')).then(() => {
          button.textContent = 'Copied';
          button.setAttribute('aria-label', 'Code copied');
        }).catch(() => {
          button.textContent = 'Select code';
          button.setAttribute('aria-label', 'Select code to copy manually');
          try {
            const range = document.createRange();
            range.selectNodeContents(code || pre);
            const selection = window.getSelection();
            selection.removeAllRanges();
            selection.addRange(range);
          } catch (_) {
            // Selection is optional; text stays visible either way.
          }
        });
      });
      pre.appendChild(button);
    });
  }

  window.markataInitCodeChrome = initCodeChrome;
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initCodeChrome, {once: true});
  } else {
    initCodeChrome();
  }
  window.addEventListener('view-transition-complete', initCodeChrome);
})();
