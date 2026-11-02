/* Enhance semantic article code blocks with a copy control. */
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
    return Promise.reject(new Error('Clipboard API unavailable'));
  }

  function codeText(pre, code) {
    const lines = code.querySelectorAll('.cl');
    const value = lines.length
      ? Array.from(lines, (line) => line.textContent).join('')
      : code.textContent;
    return value.replace(/\n+$/, '');
  }

  function addCopyButton(pre, code, header) {
    const button = document.createElement('button');
    button.type = 'button';
    button.className = 'code-block__copy';
    button.setAttribute('aria-label', 'Copy code');
    button.setAttribute('data-markata-ui', 'code-copy');
    button.innerHTML = '<svg viewBox="0 0 24 24" aria-hidden="true" focusable="false"><path d="M9 9V4a1 1 0 0 1 1-1h9a1 1 0 0 1 1 1v11a1 1 0 0 1-1 1h-5v-2h4V5h-7v4H9Zm-4 3h9a1 1 0 0 1 1 1v7a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1v-7a1 1 0 0 1 1-1Zm1 2v5h7v-5H6Z" fill="currentColor"/></svg><span>Copy</span>';

    button.addEventListener('click', () => {
      copyText(codeText(pre, code)).then(() => {
        button.classList.add('code-block__copy--done');
        button.querySelector('span').textContent = 'Copied';
        button.setAttribute('aria-label', 'Code copied');
        window.setTimeout(() => {
          button.classList.remove('code-block__copy--done');
          button.querySelector('span').textContent = 'Copy';
          button.setAttribute('aria-label', 'Copy code');
        }, 1600);
      }).catch(() => {
        button.querySelector('span').textContent = 'Select code';
        button.setAttribute('aria-label', 'Select code to copy manually');
        try {
          const range = document.createRange();
          range.selectNodeContents(code);
          const selection = window.getSelection();
          selection.removeAllRanges();
          selection.addRange(range);
        } catch (_) {
          // Keep the code visible when selection is unavailable.
        }
      });
    });

    header.appendChild(button);
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
          pre.closest('.mermaid, .chartjs-container, .csv-table, .post-copy, .sidenote, [data-no-code-chrome]') ||
          pre.classList.contains('mermaid')) return;

      const code = pre.querySelector(':scope > code') || pre.querySelector('code');
      if (!code && !pre.matches('.chroma, .highlight') &&
          !pre.closest('.chroma, .highlight')) return;

      const codeClass = (code && code.className) || '';
      if (/(?:^|\s)language-(?:mermaid|chartjs|csv|contribution-graph)(?:\s|$)/i.test(codeClass)) return;

      const text = (code || pre).textContent;
      if (!text || !text.trim()) return;
      pre.dataset.markataCodeEnhanced = 'true';

      const language = codeClass.match(/(?:^|\s)language-([a-z0-9_+.-]+)/i);
      const wrapper = pre.closest('.code-block') || (() => {
        const created = document.createElement('div');
        created.className = 'code-block';
        pre.parentNode.insertBefore(created, pre);
        created.appendChild(pre);
        return created;
      })();

      if (language && !wrapper.dataset.lang) wrapper.dataset.lang = language[1];

      let header = wrapper.querySelector(':scope > .code-block__header');
      if (!header) {
        header = document.createElement('div');
        header.className = 'code-block__header';
        const meta = document.createElement('div');
        meta.className = 'code-block__meta';
        if (wrapper.dataset.title) {
          const title = document.createElement('span');
          title.className = 'code-block__title';
          title.textContent = wrapper.dataset.title;
          meta.appendChild(title);
        }
        if (wrapper.dataset.lang) {
          const lang = document.createElement('span');
          lang.className = 'code-block__lang';
          lang.textContent = wrapper.dataset.lang;
          meta.appendChild(lang);
        }
        header.appendChild(meta);
        wrapper.insertBefore(header, pre);
      }

      if (!header.querySelector('.code-block__copy')) addCopyButton(pre, code || pre, header);
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
