const test = require("node:test");
const assert = require("node:assert/strict");
const { render, generatedLinkSlug, setFrontmatterValue, suggestLinkSlug } = require("./starter-template.js");

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
  const slug = generatedLinkSlug(source, "Link", "title", "", "my-link");
  assert.equal(slug, "link/my-link");
  assert.match(setFrontmatterValue(source, "slug", slug), /^slug: "link\/my-link"$/m);
});

test("a site-local starter named Link uses the same generated slug", () => {
  const source = '---\ntitle: "{{title}}"\n---\n';
  assert.equal(generatedLinkSlug(source, "link", "title", "", "my-link"), "link/my-link");
});

test("an explicit Link starter slug is preserved", () => {
  const source = '---\ntitle: "{{title}}"\nslug: custom/existing\n---\n';
  assert.equal(generatedLinkSlug(source, "Link", "title", "", "my-link"), "");
});

test("a nested YAML slug does not block the post-level generated slug", () => {
  const source = '---\nmetadata:\n  slug: nested/value\n---\n';
  const slug = generatedLinkSlug(source, "Link", "title", "", "my-link");
  assert.equal(slug, "link/my-link");
  const result = setFrontmatterValue(source, "slug", slug);
  assert.match(result, /^slug: "link\/my-link"$/m);
  assert.match(result, /^  slug: nested\/value$/m);
});

test("manual and numbered Link filenames do not get generated URL slugs", () => {
  const source = '---\ntitle: "{{title}}"\n---\n';
  assert.equal(generatedLinkSlug(source, "Link", "title", "chosen-name", "my-link"), "");
  assert.equal(generatedLinkSlug(source, "Link", "numbered", "", "link-1"), "");
});

test("Link slug suggestions include stable source post IDs", () => {
  const title = "typecraft (@typecraft_dev) on X";
  const first = suggestLinkSlug(title, "https://x.com/typecraft_dev/status/2123191397247578186");
  const second = suggestLinkSlug(title, "https://x.com/typecraft_dev/status/2123191397247578187");

  assert.equal(first, "typecraft-typecraft-dev-on-x-2123191397247578186");
  assert.notEqual(first, second);
});

test("Link slug suggestions use social post IDs and leave ordinary URLs alone", () => {
  assert.equal(suggestLinkSlug("A Post", "https://threads.net/@person/post/123456789"), "a-post-123456789");
  assert.equal(suggestLinkSlug("Same Title", "https://example.com/article"), "same-title");
});

test("starter without frontmatter receives creation defaults", () => {
  const result = render("# {{title}}", "A Guide", "a-guide", new Date(2024, 1, 3, 4, 5, 6));
  assert.match(result, /^date: 2024-02-03$/m);
  assert.match(result, /^published: true$/m);
  assert.match(result, /# A Guide/);
});
