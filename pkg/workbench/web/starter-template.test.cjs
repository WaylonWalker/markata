const test = require("node:test");
const assert = require("node:assert/strict");
const { render } = require("./starter-template.js");

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

test("minimal starter does not materialize inherited defaults", () => {
  const source = '---\ntitle: "{{title}}"\n---\n\n## {{title}}\n';
  const result = render(source, "A Guide", "a-guide", new Date("2024-02-03T04:05:06Z"));
  assert.match(result, /title: "A Guide"/);
  assert.match(result, /## A Guide/);
  for (const key of ["slug", "date", "draft", "published", "tags", "description"]) {
    assert.doesNotMatch(result, new RegExp(`^${key}:`, "m"));
  }
});

test("frontmatter without a title remains sparse", () => {
  const result = render("---\npublished: false\n---\n## {{title}}", "A Guide", "a-guide");
  assert.match(result, /published: false/);
  assert.match(result, /## A Guide/);
  assert.doesNotMatch(result, /^title:/m);
  assert.doesNotMatch(result, /^slug:/m);
});
