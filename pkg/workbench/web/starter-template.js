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
    const field = new RegExp(`^${escapedKey}:\\s*`);
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

  function slugify(value) {
    return String(value).normalize("NFKD").replace(/[\u0300-\u036f]/g, "").toLowerCase()
      .replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "");
  }

  function linkURLIdentifier(rawURL) {
    try {
      const source = new URL(rawURL);
      const pathSegments = source.pathname.split("/").filter(Boolean);
      const segments = pathSegments.map(segment => decodeURIComponent(segment));
      for (let index = 0; index + 1 < segments.length; index++) {
        if (!["status", "statuses", "post", "posts", "p", "reel", "reels", "comment", "comments"].includes(segments[index].toLowerCase())) continue;
        const identifier = segments[index + 1];
        if (identifier.length < 6 || !/[0-9]/.test(identifier)) continue;
        const slug = slugify(identifier);
        if (slug.length >= 6) return slug;
      }
    } catch {}
    return "";
  }

  function suggestLinkSlug(title, rawURL) {
    const slug = slugify(title);
    const identifier = linkURLIdentifier(rawURL);
    return identifier && slug && !slug.endsWith("-" + identifier) ? slug + "-" + identifier : slug;
  }

  function generatedLinkSlug(source, templateName, filenameMode, enteredFilename, slug) {
    const isLinkStarter = String(templateName).trim().toLowerCase() === "link";
    if (!isLinkStarter || filenameMode !== "title" || String(enteredFilename).trim() || !String(slug).trim()) {
      return "";
    }
    if (hasFrontmatterKey(source, "slug")) return "";
    return `link/${slug}`;
  }

  root.MarkataContentStarters = { render, hasFrontmatterKey, setFrontmatterValue, generatedLinkSlug, suggestLinkSlug };
  if (typeof module !== "undefined" && module.exports) module.exports = root.MarkataContentStarters;
})(globalThis);
