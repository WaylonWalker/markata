(function (root) {
  function render(source, title, slug, now = new Date(), dateOverride = "") {
    const today = [now.getFullYear(), String(now.getMonth() + 1).padStart(2, "0"), String(now.getDate()).padStart(2, "0")].join("-");
    const date = /^\d{4}-\d{2}-\d{2}$/.test(dateOverride) ? dateOverride : today;
    const values = { title: String(title), slug: String(slug), date, datetime: now.toISOString() };
    const expand = text => text.replace(/\{\{(title|slug|date|datetime)\}\}/g, (_, key) => values[key]);
    const frontmatter = String(source).match(/^(---\r?\n)([\s\S]*?)(\r?\n---(?:\r?\n|$))/);
    if (!frontmatter) return `---\ndate: ${date}\npublished: true\n---\n\n${expand(String(source))}`;
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
    let hasDate = false;
    let hasPublished = false;
    const hasDraft = lines.some(line => /^draft:\s*true\s*(?:#.*)?$/i.test(line));
    for (let index = 0; index < lines.length; index += 1) {
      const dateField = lines[index].match(/^(date:\s*)(.*)$/i);
      if (dateField) {
        hasDate = true;
        const value = dateField[2].replace(/\s+#.*$/, "").trim();
        if (dateOverride || !value || /^(?:null|~)$/i.test(value)) lines[index] = `${dateField[1]}${date}`;
      }
      const publishedField = lines[index].match(/^(published:\s*)(.*)$/i);
      if (publishedField) {
        hasPublished = true;
        const value = publishedField[2].replace(/\s+#.*$/, "").trim();
        if (!value || /^(?:null|~)$/i.test(value)) lines[index] = `${publishedField[1]}${hasDraft ? "false" : "true"}`;
      }
    }
    if (!hasDate) lines.push(`date: ${date}`);
    if (!hasPublished) lines.push(`published: ${hasDraft ? "false" : "true"}`);
    return `${frontmatter[1]}${lines.join("\n")}${frontmatter[3]}${expand(String(source).slice(frontmatter[0].length))}`;
  }

  function hasFrontmatterKey(source, key) {
    const match = String(source).match(/^(---\r?\n)([\s\S]*?)(\r?\n---(?:\r?\n|$))/);
    if (!match) return false;
    const escapedKey = String(key).replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
    const field = new RegExp(`^\\s*${escapedKey}:\\s*`);
    return match[2].split(/\r?\n/).some(line => field.test(line));
  }

  function setFrontmatterValue(source, key, value) {
    const text = String(source);
    const match = text.match(/^(---\r?\n)([\s\S]*?)(\r?\n---(?:\r?\n|$))/);
    const quoted = `${key}: ${JSON.stringify(String(value))}`;
    if (!match) return `---\n${quoted}\n---\n\n${text}`;
    const lines = match[2].split(/\r?\n/);
    const escapedKey = String(key).replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
    const field = new RegExp(`^${escapedKey}:\\s*`);
    const index = lines.findIndex(line => field.test(line));
    if (index >= 0) lines[index] = quoted;
    else lines.push(quoted);
    return `${match[1]}${lines.join("\n")}${match[3]}${text.slice(match[0].length)}`;
  }

  function generatedLinkSlug(source, templateID, filenameMode, enteredFilename, slug) {
    if (templateID !== "builtin:link" || filenameMode !== "title" || String(enteredFilename).trim() || !String(slug).trim()) {
      return "";
    }
    if (hasFrontmatterKey(source, "slug")) return "";
    return `link/${slug}`;
  }

  root.MarkataContentStarters = { render, hasFrontmatterKey, setFrontmatterValue, generatedLinkSlug };
  if (typeof module !== "undefined" && module.exports) module.exports = root.MarkataContentStarters;
})(globalThis);
