(() => {
  const collator = new Intl.Collator(undefined, { numeric: true, sensitivity: 'base' });
  const numberPattern = /^[+-]?(?:\d+\.?\d*|\.\d+)(?:e[+-]?\d+)?%?$/i;

  document.querySelectorAll('.post-content table thead').forEach((head) => {
    const table = head.closest('table');
    const body = table?.tBodies[0];
    const headerRow = head.rows[head.rows.length - 1];
    if (!table || !body || !headerRow) return;

    Array.from(headerRow.cells).forEach((cell, column) => {
      const label = cell.textContent.trim();
      // Sorting controls may wrap formatted text but must not swallow
      // links, nested controls, or incorrectly indexed spanning headers.
      if (!label || cell.tagName !== 'TH' || cell.colSpan !== 1 ||
          cell.querySelector('a, button, input, select, textarea, [tabindex]')) return;

      const button = document.createElement('button');
      button.type = 'button';
      button.className = 'table-sort-button';
      button.setAttribute('aria-label', `Sort by ${label}`);
      // Preserve <code>, <em>, <abbr>, and other authored heading markup.
      // Replacing the header with its text lost these details and links.
      while (cell.firstChild) button.appendChild(cell.firstChild);
      cell.appendChild(button);
      cell.setAttribute('aria-sort', 'none');

      button.addEventListener('click', () => {
        const ascending = cell.getAttribute('aria-sort') !== 'ascending';
        headerRow.querySelectorAll('th').forEach((header) => {
          if (header !== cell) header.setAttribute('aria-sort', 'none');
        });
        cell.setAttribute('aria-sort', ascending ? 'ascending' : 'descending');

        const rows = Array.from(body.rows).map((row, index) => ({
          row,
          index,
          value: row.cells[column]?.textContent.trim() ?? '',
        }));
        const values = rows.map(({ value }) => value.replace(/[,\s%]/g, ''));
        const numeric = values.length > 0 && values.every((value) => value === '' || numberPattern.test(value));

        rows.sort((a, b) => {
          if (!a.value) return b.value ? 1 : a.index - b.index;
          if (!b.value) return -1;
          const result = numeric
            ? Number(a.value.replace(/[,\s%]/g, '')) - Number(b.value.replace(/[,\s%]/g, ''))
            : collator.compare(a.value, b.value);
          return (ascending ? result : -result) || a.index - b.index;
        });
        rows.forEach(({ row }) => body.append(row));
      });
    });
  });
})();
