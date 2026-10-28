# Visualize: Create, Serve, Open

## 1. Select the evidence

For a session post, use the relevant conversation, supplied artifacts, and observed results.
Separate completed work, decisions, and unresolved work.
For a supplied topic, use that material instead of unrelated session history.
Exclude secrets and unrequested private material.
Label uncertainty instead of inventing facts, measurements, or success.

**Done:** The post scope and the evidence for its claims are clear.

## 2. Create the visual post

Read `visual-authoring.md` and `rich-markdown.md` in this topic directory.
For numerical charts, also read `chart-authoring.md`.
Use their composition rules and native Markdown primitives.

In an existing site, inspect its configuration and content conventions.
Use `writing-frontmatter.md` for the post metadata.
Preserve existing publication flags.
Keep new session posts out of public publication unless the user authorizes publication.

Outside a Markata site, keep preview scaffolding in an isolated temporary directory.
Inspect the available Markata binary and its `--help` before creating that preview site.
Keep its configuration, output, and cache outside the unrelated repository.
Retain the authored Markdown at the agreed delivery location.

If the user requests source only or exactly one file, deliver that artifact without preview scaffolding or a server.
Explain that rich styling requires Markata.
Stop here for those requests.

**Done:** The Markdown artifact exists and follows the requested scope.
A composed visual post contains more than a Mermaid diagram followed by an ordinary table.

## 3. Serve a local preview

Read `build-deployment.md` for local preview guidance.
Use an existing server only after you establish that it serves the intended site and output.
Otherwise, start `markata-go serve --fast --host 127.0.0.1 --port <available-port>` from the preview site.
Use the installed binary's equivalent command name if it differs.
Select an available port without stopping another service.

Use the agent's supported persistent-process mechanism so the server remains available after the response.
Retain the process handle, site directory, and log location.
If persistent processes require approval, obtain that approval through the available permission mechanism.
If persistence is unavailable, report the limitation instead of promising a lasting preview.

Keep preview-only draft inclusion in isolated configuration or an isolated preview copy.
Inspect current CLI help and site filters before selecting that method.
Keep the source publication flags unchanged.
Local preview does not authorize public publishing, deployment, or a public tunnel.

**Done:** The intended server is running on loopback, with a known process handle and lifetime.

## 4. Check the post URL

Resolve the exact post URL from the site's output and routing.
Poll that URL with a bounded timeout until the server returns a successful response.
Check the response for the intended post's title or distinctive content.
A successful homepage response alone does not establish that the post is available.
If the build fails or the post remains unavailable, report the error and relevant logs.
Correct the cause before attempting the browser open.

If browser inspection is available, inspect the rendered layout and narrow-screen behavior.
Keep HTTP availability distinct from visual inspection.

**Done:** The exact post URL returns the intended post.

## 5. Open and hand off

Open the checked URL with the available browser tool or platform opener.
For desktop environments, use an available opener such as `xdg-open` on Linux or `open` on macOS.
For remote environments, use an available client-browser integration.
If no opener is available, report that limitation and provide the URL.
An opener's successful exit proves an open request, not that the user saw the page.

Report the Markdown path, exact preview URL, browser-open result, and server handle or stop instruction.
Keep the preview site and server available for the user.
After the user finishes, stop only the owned server and remove only its disposable preview artifacts.
Preserve the authored Markdown.

**Done:** The artifact exists, its preview responds, and browser opening was attempted.
If any stage is blocked, name that stage without claiming full completion.
