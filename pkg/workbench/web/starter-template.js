(function (root) {
  function render(source, title, slug, now = new Date()) {
    const date = [now.getFullYear(), String(now.getMonth() + 1).padStart(2, "0"), String(now.getDate()).padStart(2, "0")].join("-");
    const values = { title: String(title), slug: String(slug), date, datetime: now.toISOString() };
    const expand = text => text.replace(/\{\{(title|slug|date|datetime)\}\}/g, (_, key) => values[key]);
    const frontmatter = String(source).match(/^(---\r?\n)([\s\S]*?)(\r?\n---(?:\r?\n|$))/);
    if (!frontmatter) return expand(String(source));
    const lines = frontmatter[2].split(/\r?\n/).map(line => {
      const titleField = line.match(/^(\s*title:\s*).*/);
      if (titleField) return `${titleField[1]}${JSON.stringify(values.title)}`;
      const field = line.match(/^(\s*(?:[A-Za-z0-9_-]+:\s*|-\s*))(.*)$/);
      if (!field) return expand(line);
      const [, prefix, raw] = field;
      if (/\{\{title\}\}/.test(raw) && !/["']/.test(raw)) {
        return `${prefix}${JSON.stringify(expand(raw))}`;
      }
      const protectedValues = [];
      const protect = (_, inner) => {
        if (!/\{\{(?:title|slug|date|datetime)\}\}/.test(inner)) return _;
        const marker = `\u0000starter-value-${protectedValues.length}\u0000`;
        protectedValues.push(JSON.stringify(expand(inner)));
        return marker;
      };
      const protectedRaw = raw.replace(/"((?:\\.|[^"\\])*)"/g, protect)
        .replace(/'((?:''|[^'])*)'/g, (_, inner) => protect(_, inner.replaceAll("''", "'")));
      let expanded = expand(protectedRaw);
      protectedValues.forEach((replacement, index) => {
        expanded = expanded.replaceAll(`\u0000starter-value-${index}\u0000`, replacement);
      });
      return `${prefix}${expanded}`;
    });
    return `${frontmatter[1]}${lines.join("\n")}${frontmatter[3]}${expand(String(source).slice(frontmatter[0].length))}`;
  }

  root.MarkataContentStarters = { render };
  if (typeof module !== "undefined" && module.exports) module.exports = root.MarkataContentStarters;
})(globalThis);
