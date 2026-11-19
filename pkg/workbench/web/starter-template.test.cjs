const test = require("node:test");
const assert = require("node:assert/strict");
const { render, generatedLinkSlug, setFrontmatterValue } = require("./starter-template.js");

test("starter title substitutions remain valid quoted YAML scalars", () => {
  const source = [
    "---",
    "title: '{{title}}'",
    'subtitle: "About {{title}}"',
    "description: From {{title}}",
    'tags: ["{{title}}"]',
    "tags:",
    "  - {{title}}",
    "---",
    "Body: {{title}}",
  ].join("\n");
  const result = render(source, 'A "quoted"\nTitle', "a-quoted-title", new Date("2024-02-03T04:05:06Z"));

  assert.match(result, /title: "A \\\"quoted\\\"\\nTitle"/);
  assert.match(result, /subtitle: "About A \\\"quoted\\\"\\nTitle"/);
  assert.match(result, /description: "From A \\\"quoted\\\"\\nTitle"/);
  assert.match(result, /  - "A \\\"quoted\\\"\\nTitle"/);
  assert.match(result, /tags: \["A \\\"quoted\\\"\\nTitle"\]/);
  assert.match(result, /Body: A "quoted"\nTitle/);
});

test("literal placeholder text entered as a title is not expanded again", () => {
  const result = render('---\ntitle: "{{title}}"\nsubtitle: "{{title}}"\n---\n', "{{slug}}", "actual-slug", new Date("2024-02-03T04:05:06Z"));
  assert.match(result, /title: "\{\{slug\}\}"/);
  assert.match(result, /subtitle: "\{\{slug\}\}"/);
});

test("minimal starter receives date and published creation defaults", () => {
  const source = '---\ntitle: "{{title}}"\n---\n\n## {{title}}\n';
  const result = render(source, "A Guide", "a-guide", new Date(2024, 1, 3, 4, 5, 6));
  assert.match(result, /title: "A Guide"/);
  assert.match(result, /^date: 2024-02-03$/m);
  assert.match(result, /^published: true$/m);
  assert.match(result, /## A Guide/);
  for (const key of ["slug", "draft", "tags", "description"]) {
    assert.doesNotMatch(result, new RegExp(`^${key}:`, "m"));
  }
});

test("frontmatter without a title remains sparse", () => {
  const result = render("---\npublished: false\n---\n## {{title}}", "A Guide", "a-guide");
  assert.match(result, /published: false/);
  assert.match(result, /^date: \d{4}-\d{2}-\d{2}$/m);
  assert.match(result, /## A Guide/);
  assert.doesNotMatch(result, /^title:/m);
  assert.doesNotMatch(result, /^slug:/m);
});

test("missing publication follows an explicit draft field", () => {
  const result = render("---\ndraft: true\ndate:\npublished: null\n---\nBody", "Draft", "draft", new Date(2024, 1, 3, 4, 5, 6));
  assert.match(result, /^date:\s*\d{4}-\d{2}-\d{2}$/m);
  assert.match(result, /^published: false$/m);
});

test("generated Link filenames get a namespaced URL slug", () => {
  const source = '---\ntitle: "{{title}}"\nurl: ""\n---\n';
  const slug = generatedLinkSlug(source, "builtin:link", "title", "", "my-link");
  assert.equal(slug, "link/my-link");
  assert.match(setFrontmatterValue(source, "slug", slug), /^slug: "link\/my-link"$/m);
});

test("an explicit Link starter slug is preserved", () => {
  const source = '---\ntitle: "{{title}}"\nslug: custom/existing\n---\n';
  assert.equal(generatedLinkSlug(source, "builtin:link", "title", "", "my-link"), "");
});

test("manual and numbered Link filenames do not get generated URL slugs", () => {
  const source = '---\ntitle: "{{title}}"\n---\n';
  assert.equal(generatedLinkSlug(source, "builtin:link", "title", "chosen-name", "my-link"), "");
  assert.equal(generatedLinkSlug(source, "builtin:link", "numbered", "", "link-1"), "");
});

test("starter without frontmatter receives creation defaults", () => {
  const result = render("# {{title}}", "A Guide", "a-guide", new Date(2024, 1, 3, 4, 5, 6));
  assert.match(result, /^date: 2024-02-03$/m);
  assert.match(result, /^published: true$/m);
  assert.match(result, /# A Guide/);
});
