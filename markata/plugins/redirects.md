---
content: "Creates redirects for times when your backend server can't.\n\n## Configuration\n\nEnable
  the redirect hook by adding it to your list of hooks.\n\n``` toml\n[markata]\n\n#
  Were you keep static assets to copy into the project, default is static\n# the assets_dir
  will set the default _redirects file directory\nassets_dir = \"static\"\n\n# You
  can override the default redirects file location\nredirects = static/_redirects\n\nhooks
  = [\n   # creates redirects from static/_redirects file\n   \"markata.plugins.redirects\",\n
  \  # copies your static assets into the output_dir (default: `markout`)\n   \"markata.plugins.copy_assets\",\n
  \ ...\n]\n```\n\n## Syntax\n\nYour `_redirects` file is a simplified version of
  what services like cloudflare\npages or netlify use.  In fact you can use the same
  redirects file!\n\nHere is an example that will redirect `/old` to `/new` and `/CHANGELOG`
  to\n`/changelog`\n\n```\n/old        /new\n/CHANGELOG  /changelog\n```\n\n## Limitations\n_no
  splats_\n\nSince it is static generated this plugin cannot cover *'s.  * or splat\nredirects
  need to be taken care of server side.  It also cannot change the http\ncode, this
  is only\n\n## Features\n\nThe features of markata.plugins.redirect is pretty limited
  since it is\nimplemented only as a static page.  Other features require server side\nimplementation.\n\n|
  Feature                             | Support | Example                                                         |
  Notes                                                                                             |\n|
  ----------------------------------- | ------- | ---------------------------------------------------------------
  | -------------------------------------------------------------------------------------------------
  |\n| Force                               | Yes     | `/pagethatexists /otherpage`
  \                                   | Creates an index.html with http-equiv and
  canonical                                               |\n| Redirects (301, 302,
  303, 307, 308) | No      | `/home / 301`                                                   |
  Ignored, requires server side implementation                                                      |\n|
  Rewrites (other status codes)       | No      | `/blog/* /blog/404.html 404`                                    |
  ...                                                                                               |\n|
  Splats                              | No      | `/blog/* /blog/:splat`                                          |
  ...                                                                                               |\n|
  Placeholders                        | No      | `/blog/:year/:month/:date/:slug
  /news/:year/:month/:date/:slug` | ...                                                                                               |\n|
  Query Parameters                    | No      | `/shop id=:id /blog/:id 301`                                    |
  ...                                                                                               |\n|
  Proxying                            | No      | `/blog/* https://blog.my.domain/:splat
  200`                     | ...                                                                                               |\n|
  Domain-level redirects              | No      | `workers.example.com/* workers.example.com/blog/:splat
  301`     | ...                                                                                               |\n|
  Redirect by country or language     | No      | `/ /us 302 Country=us`                                          |
  ...                                                                                               |\n|
  Redirect by cookie                  | No      | `/* /preview/:splat 302 Cookie=preview`
  \                       | ...                                                                                               |\n\n>
  Compare with\n> [cloudflare-pages](https://developers.cloudflare.com/pages/platform/redirects/)\n\n!!!
  tip\n    If you have a public site, pair this up with\n    [ahrefs](https://app.ahrefs.com/dashboard)
  to keep up with pages that have\n    moved without you realizing.\n\n\n!! class
  <h2 id='Redirect' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Redirect
  <em class='small'>class</em></h2>\n    DataClass to store the original and new url\n???+
  source \"Redirect <em class='small'>source</em>\"\n\n```python\n\n        class
  Redirect(pydantic.BaseModel):\n            \"DataClass to store the original and
  new url\"\n\n            original: str\n            new: str\n            markata:
  Markata\n            model_config = ConfigDict(validate_assignment=True, arbitrary_types_allowed=True)\n```\n\n\n!!
  class <h2 id='RedirectsConfig' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>RedirectsConfig <em class='small'>class</em></h2>\n\n???+ source \"RedirectsConfig
  <em class='small'>source</em>\"\n\n```python\n\n        class RedirectsConfig(pydantic.BaseModel):\n
  \           assets_dir: Path = Path(\"static\")\n            redirects_file: Optional[Path]
  = None\n\n            @pydantic.validator(\"redirects_file\", always=True)\n            def
  default_redirects_file(\n                cls: \"RedirectsConfig\", v: Path, *, values:
  Dict\n            ) -> Path:\n                if not v:\n                    return
  Path(values[\"assets_dir\"]) / \"_redirects\"\n                return v\n```\n\n\n!!
  class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config
  <em class='small'>class</em></h2>\n\n???+ source \"Config <em class='small'>source</em>\"\n\n```python\n\n
  \       class Config(pydantic.BaseModel):\n            redirects: RedirectsConfig
  = RedirectsConfig()\n```\n\n\n!! function <h2 id='config_model' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>\n\n???+
  source \"config_model <em class='small'>source</em>\"\n\n```python\n\n        def
  config_model(markata: \"Markata\") -> None:\n            markata.config_models.append(Config)\n```\n\n\n!!
  function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save
  <em class='small'>function</em></h2>\n    saves an index.html in the directory called
  out by the redirect.\n???+ source \"save <em class='small'>source</em>\"\n\n```python\n\n
  \       def save(markata: \"Markata\") -> None:\n            \"\"\"\n            saves
  an index.html in the directory called out by the redirect.\n            \"\"\"\n
  \           redirects_file = Path(markata.config.redirects.redirects_file)\n            if
  redirects_file.exists():\n                raw_redirects = redirects_file.read_text().split(\"\\n\")\n
  \           else:\n                raw_redirects = []\n\n            redirects =
  [\n                Redirect(original=s[0], new=s[1], markata=markata)\n                for
  r in raw_redirects\n                if \"*\" not in r and len(s := r.split()) ==
  2 and not r.strip().startswith(\"#\")\n            ]\n\n            if \"redirect_template\"
  in markata.config:\n                template_file = Path(str(markata.config.get(\"redirect_template\")))\n
  \           else:\n                template_file = DEFAULT_REDIRECT_TEMPLATE\n            template
  = Template(template_file.read_text())\n\n            for redirect in redirects:\n
  \               file = markata.config.output_dir / redirect.original.strip(\"/\")
  / \"index.html\"\n                file.parent.mkdir(parents=True, exist_ok=True)\n
  \               file.write_text(template.render(redirect.dict(), config=markata.config))\n```\n\n\n!!
  method <h2 id='default_redirects_file' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>default_redirects_file <em class='small'>method</em></h2>\n\n???+ source
  \"default_redirects_file <em class='small'>source</em>\"\n\n```python\n\n        def
  default_redirects_file(\n                cls: \"RedirectsConfig\", v: Path, *, values:
  Dict\n            ) -> Path:\n                if not v:\n                    return
  Path(values[\"assets_dir\"]) / \"_redirects\"\n                return v\n```\n\n"
date: 0001-01-01
description: Creates redirects for times when your backend server can Enable the redirect
  hook by adding it to your list of hooks. Your  Here is an example that will redirec
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Redirects.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Creates redirects for times when your
    backend server can Enable the redirect hook by adding it to your list of hooks.
    Your  Here is an example that will redirec\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Redirects.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Creates redirects for times when your
    backend server can Enable the redirect hook by adding it to your list of hooks.
    Your  Here is an example that will redirec\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<div class='container flex flex-row min-h-screen'>\n    <div>\n    </div>\n
    \   <div class='flex-grow px-8 mx-auto min-h-screen'>\n<header class='flex justify-center
    items-center p-8'>\n\n    <nav class='flex justify-center items-center my-8'>\n
    \       <a\n            href='/'>markata</a>\n        <a\n            href='https://github.com/WaylonWalker/markata'>GitHub</a>\n
    \       <a\n            href='https://markata.dev/docs/'>docs</a>\n        <a\n
    \           href='https://markata.dev/plugins/'>plugins</a>\n    </nav>\n\n    <div>\n
    \       <label id=\"theme-switch\" class=\"theme-switch\" for=\"checkbox-theme\"
    title=\"light/dark mode toggle\">\n            <input type=\"checkbox\" id=\"checkbox-theme\"
    />\n            <div class=\"slider round\"></div>\n        </label>\n    </div>\n</header><article
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Redirects.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Creates redirects
    for times when your backend server can't.</p>\n<h2 id=\"configuration\">Configuration
    <a class=\"header-anchor\" href=\"#configuration\"><svg class=\"heading-permalink\"
    aria-hidden=\"true\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
    viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
    d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Enable the redirect
    hook by adding it to your list of hooks.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
    class=\"c1\"># Were you keep static assets to copy into the project, default is
    static</span>\n<span class=\"c1\"># the assets_dir will set the default _redirects
    file directory</span>\n<span class=\"n\">assets_dir</span><span class=\"w\"> </span><span
    class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;static&quot;</span>\n\n<span
    class=\"c1\"># You can override the default redirects file location</span>\n<span
    class=\"n\">redirects</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"err\">static/</span><span class=\"mi\">_</span><span
    class=\"n\">redirects</span>\n\n<span class=\"n\">hooks</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[</span>\n<span
    class=\"w\">   </span><span class=\"c1\"># creates redirects from static/_redirects
    file</span>\n<span class=\"w\">   </span><span class=\"s2\">&quot;markata.plugins.redirects&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">   </span><span class=\"c1\"># copies
    your static assets into the output_dir (default: `markout`)</span>\n<span class=\"w\">
    \  </span><span class=\"s2\">&quot;markata.plugins.copy_assets&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">  </span><span class=\"err\">...</span>\n<span
    class=\"err\">]</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"syntax\">Syntax <a
    class=\"header-anchor\" href=\"#syntax\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
    fill=\"currentColor\" focusable=\"false\" height=\"1em\" viewBox=\"0 0 24 24\"
    width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199 13.599a5.99
    5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992 5.992 0
    0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242 6.003
    6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Your <code>_redirects</code>
    file is a simplified version of what services like cloudflare\npages or netlify
    use.  In fact you can use the same redirects file!</p>\n<p>Here is an example
    that will redirect <code>/old</code> to <code>/new</code> and <code>/CHANGELOG</code>
    to\n<code>/changelog</code></p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>/old        /new\n/CHANGELOG
    \ /changelog\n</pre></div>\n\n</pre>\n\n<h2 id=\"limitations\">Limitations <a
    class=\"header-anchor\" href=\"#limitations\"><svg class=\"heading-permalink\"
    aria-hidden=\"true\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
    viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
    d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p><em>no splats</em></p>\n<p>Since
    it is static generated this plugin cannot cover *'s.  * or splat\nredirects need
    to be taken care of server side.  It also cannot change the http\ncode, this is
    only</p>\n<h2 id=\"features\">Features <a class=\"header-anchor\" href=\"#features\"><svg
    class=\"heading-permalink\" aria-hidden=\"true\" fill=\"currentColor\" focusable=\"false\"
    height=\"1em\" viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
    d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The features of markata.plugins.redirect
    is pretty limited since it is\nimplemented only as a static page.  Other features
    require server side\nimplementation.</p>\n<table>\n<thead>\n<tr>\n<th>Feature</th>\n<th>Support</th>\n<th>Example</th>\n<th>Notes</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>Force</td>\n<td>Yes</td>\n<td><code>/pagethatexists
    /otherpage</code></td>\n<td>Creates an index.html with http-equiv and canonical</td>\n</tr>\n<tr>\n<td>Redirects
    (301, 302, 303, 307, 308)</td>\n<td>No</td>\n<td><code>/home / 301</code></td>\n<td>Ignored,
    requires server side implementation</td>\n</tr>\n<tr>\n<td>Rewrites (other status
    codes)</td>\n<td>No</td>\n<td><code>/blog/* /blog/404.html 404</code></td>\n<td>...</td>\n</tr>\n<tr>\n<td>Splats</td>\n<td>No</td>\n<td><code>/blog/*
    /blog/:splat</code></td>\n<td>...</td>\n</tr>\n<tr>\n<td>Placeholders</td>\n<td>No</td>\n<td><code>/blog/:year/:month/:date/:slug
    /news/:year/:month/:date/:slug</code></td>\n<td>...</td>\n</tr>\n<tr>\n<td>Query
    Parameters</td>\n<td>No</td>\n<td><code>/shop id=:id /blog/:id 301</code></td>\n<td>...</td>\n</tr>\n<tr>\n<td>Proxying</td>\n<td>No</td>\n<td><code>/blog/*
    https://blog.my.domain/:splat 200</code></td>\n<td>...</td>\n</tr>\n<tr>\n<td>Domain-level
    redirects</td>\n<td>No</td>\n<td><code>workers.example.com/* workers.example.com/blog/:splat
    301</code></td>\n<td>...</td>\n</tr>\n<tr>\n<td>Redirect by country or language</td>\n<td>No</td>\n<td><code>/
    /us 302 Country=us</code></td>\n<td>...</td>\n</tr>\n<tr>\n<td>Redirect by cookie</td>\n<td>No</td>\n<td><code>/*
    /preview/:splat 302 Cookie=preview</code></td>\n<td>...</td>\n</tr>\n</tbody>\n</table>\n<blockquote>\n<p>Compare
    with\n<a href=\"https://developers.cloudflare.com/pages/platform/redirects/\">cloudflare-pages</a></p>\n</blockquote>\n<div
    class=\"admonition tip\">\n<p class=\"admonition-title\">Tip</p>\n<p>If you have
    a public site, pair this up with\n<a href=\"https://app.ahrefs.com/dashboard\">ahrefs</a>
    to keep up with pages that have\nmoved without you realizing.</p>\n</div>\n<p>!!
    class <h2 id='Redirect' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>Redirect <em class='small'>class</em></h2>\nDataClass to store the original
    and new url</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">Redirect <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Redirect</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"s2\">&quot;DataClass to store the original and new url&quot;</span>\n\n
    \           <span class=\"n\">original</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span>\n            <span class=\"n\">new</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span>\n            <span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"n\">Markata</span>\n            <span class=\"n\">model_config</span>
    <span class=\"o\">=</span> <span class=\"n\">ConfigDict</span><span class=\"p\">(</span><span
    class=\"n\">validate_assignment</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">arbitrary_types_allowed</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='RedirectsConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>RedirectsConfig <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">RedirectsConfig
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">RedirectsConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">assets_dir</span><span class=\"p\">:</span> <span
    class=\"n\">Path</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;static&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">redirects_file</span><span class=\"p\">:</span>
    <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;redirects_file&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">always</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">default_redirects_file</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">cls</span><span class=\"p\">:</span> <span class=\"s2\">&quot;RedirectsConfig&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">:</span> <span
    class=\"n\">Path</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Path</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">v</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;assets_dir&quot;</span><span class=\"p\">])</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;_redirects&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Config <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">redirects</span><span class=\"p\">:</span> <span
    class=\"n\">RedirectsConfig</span> <span class=\"o\">=</span> <span class=\"n\">RedirectsConfig</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='config_model'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">config_model <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>save <em class='small'>function</em></h2>\nsaves an index.html in the directory
    called out by the redirect.</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">save <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            saves an index.html in the directory called out by the
    redirect.</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n            <span
    class=\"n\">redirects_file</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">redirects</span><span
    class=\"o\">.</span><span class=\"n\">redirects_file</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">redirects_file</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \               <span class=\"n\">raw_redirects</span> <span class=\"o\">=</span>
    <span class=\"n\">redirects_file</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n            <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                <span class=\"n\">raw_redirects</span> <span
    class=\"o\">=</span> <span class=\"p\">[]</span>\n\n            <span class=\"n\">redirects</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span>\n                <span class=\"n\">Redirect</span><span
    class=\"p\">(</span><span class=\"n\">original</span><span class=\"o\">=</span><span
    class=\"n\">s</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">],</span> <span class=\"n\">new</span><span class=\"o\">=</span><span
    class=\"n\">s</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">],</span> <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">)</span>\n                <span class=\"k\">for</span>
    <span class=\"n\">r</span> <span class=\"ow\">in</span> <span class=\"n\">raw_redirects</span>\n
    \               <span class=\"k\">if</span> <span class=\"s2\">&quot;*&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">r</span>
    <span class=\"ow\">and</span> <span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"n\">s</span> <span class=\"o\">:=</span> <span class=\"n\">r</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">())</span>
    <span class=\"o\">==</span> <span class=\"mi\">2</span> <span class=\"ow\">and</span>
    <span class=\"ow\">not</span> <span class=\"n\">r</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">()</span><span class=\"o\">.</span><span
    class=\"n\">startswith</span><span class=\"p\">(</span><span class=\"s2\">&quot;#&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"p\">]</span>\n\n            <span
    class=\"k\">if</span> <span class=\"s2\">&quot;redirect_template&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">:</span>\n                <span class=\"n\">template_file</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;redirect_template&quot;</span><span
    class=\"p\">)))</span>\n            <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">template_file</span> <span class=\"o\">=</span>
    <span class=\"n\">DEFAULT_REDIRECT_TEMPLATE</span>\n            <span class=\"n\">template</span>
    <span class=\"o\">=</span> <span class=\"n\">Template</span><span class=\"p\">(</span><span
    class=\"n\">template_file</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">())</span>\n\n            <span class=\"k\">for</span> <span class=\"n\">redirect</span>
    <span class=\"ow\">in</span> <span class=\"n\">redirects</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">file</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">output_dir</span> <span class=\"o\">/</span>
    <span class=\"n\">redirect</span><span class=\"o\">.</span><span class=\"n\">original</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;index.html&quot;</span>\n                <span class=\"n\">file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">parents</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">exist_ok</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n                <span class=\"n\">file</span><span class=\"o\">.</span><span
    class=\"n\">write_text</span><span class=\"p\">(</span><span class=\"n\">template</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span><span
    class=\"n\">redirect</span><span class=\"o\">.</span><span class=\"n\">dict</span><span
    class=\"p\">(),</span> <span class=\"n\">config</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='default_redirects_file'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_redirects_file
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">default_redirects_file <em
    class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">default_redirects_file</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">cls</span><span class=\"p\">:</span> <span class=\"s2\">&quot;RedirectsConfig&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">:</span> <span
    class=\"n\">Path</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Path</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">v</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;assets_dir&quot;</span><span class=\"p\">])</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;_redirects&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Redirects.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Creates redirects for times when your backend server
    can Enable the redirect hook by adding it to your list of hooks. Your  Here is
    an example that will redirec\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Redirects.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Creates redirects for times when your
    backend server can Enable the redirect hook by adding it to your list of hooks.
    Your  Here is an example that will redirec\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Redirects.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Redirects.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>Creates
    redirects for times when your backend server can't.</p>\n<h2 id=\"configuration\">Configuration
    <a class=\"header-anchor\" href=\"#configuration\"><svg class=\"heading-permalink\"
    aria-hidden=\"true\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
    viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
    d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Enable the redirect
    hook by adding it to your list of hooks.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
    class=\"c1\"># Were you keep static assets to copy into the project, default is
    static</span>\n<span class=\"c1\"># the assets_dir will set the default _redirects
    file directory</span>\n<span class=\"n\">assets_dir</span><span class=\"w\"> </span><span
    class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;static&quot;</span>\n\n<span
    class=\"c1\"># You can override the default redirects file location</span>\n<span
    class=\"n\">redirects</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"err\">static/</span><span class=\"mi\">_</span><span
    class=\"n\">redirects</span>\n\n<span class=\"n\">hooks</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"p\">[</span>\n<span
    class=\"w\">   </span><span class=\"c1\"># creates redirects from static/_redirects
    file</span>\n<span class=\"w\">   </span><span class=\"s2\">&quot;markata.plugins.redirects&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">   </span><span class=\"c1\"># copies
    your static assets into the output_dir (default: `markout`)</span>\n<span class=\"w\">
    \  </span><span class=\"s2\">&quot;markata.plugins.copy_assets&quot;</span><span
    class=\"p\">,</span>\n<span class=\"w\">  </span><span class=\"err\">...</span>\n<span
    class=\"err\">]</span>\n</pre></div>\n\n</pre>\n\n<h2 id=\"syntax\">Syntax <a
    class=\"header-anchor\" href=\"#syntax\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
    fill=\"currentColor\" focusable=\"false\" height=\"1em\" viewBox=\"0 0 24 24\"
    width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199 13.599a5.99
    5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992 5.992 0
    0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242 6.003
    6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Your <code>_redirects</code>
    file is a simplified version of what services like cloudflare\npages or netlify
    use.  In fact you can use the same redirects file!</p>\n<p>Here is an example
    that will redirect <code>/old</code> to <code>/new</code> and <code>/CHANGELOG</code>
    to\n<code>/changelog</code></p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>/old        /new\n/CHANGELOG
    \ /changelog\n</pre></div>\n\n</pre>\n\n<h2 id=\"limitations\">Limitations <a
    class=\"header-anchor\" href=\"#limitations\"><svg class=\"heading-permalink\"
    aria-hidden=\"true\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
    viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
    d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p><em>no splats</em></p>\n<p>Since
    it is static generated this plugin cannot cover *'s.  * or splat\nredirects need
    to be taken care of server side.  It also cannot change the http\ncode, this is
    only</p>\n<h2 id=\"features\">Features <a class=\"header-anchor\" href=\"#features\"><svg
    class=\"heading-permalink\" aria-hidden=\"true\" fill=\"currentColor\" focusable=\"false\"
    height=\"1em\" viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
    d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>The features of markata.plugins.redirect
    is pretty limited since it is\nimplemented only as a static page.  Other features
    require server side\nimplementation.</p>\n<table>\n<thead>\n<tr>\n<th>Feature</th>\n<th>Support</th>\n<th>Example</th>\n<th>Notes</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>Force</td>\n<td>Yes</td>\n<td><code>/pagethatexists
    /otherpage</code></td>\n<td>Creates an index.html with http-equiv and canonical</td>\n</tr>\n<tr>\n<td>Redirects
    (301, 302, 303, 307, 308)</td>\n<td>No</td>\n<td><code>/home / 301</code></td>\n<td>Ignored,
    requires server side implementation</td>\n</tr>\n<tr>\n<td>Rewrites (other status
    codes)</td>\n<td>No</td>\n<td><code>/blog/* /blog/404.html 404</code></td>\n<td>...</td>\n</tr>\n<tr>\n<td>Splats</td>\n<td>No</td>\n<td><code>/blog/*
    /blog/:splat</code></td>\n<td>...</td>\n</tr>\n<tr>\n<td>Placeholders</td>\n<td>No</td>\n<td><code>/blog/:year/:month/:date/:slug
    /news/:year/:month/:date/:slug</code></td>\n<td>...</td>\n</tr>\n<tr>\n<td>Query
    Parameters</td>\n<td>No</td>\n<td><code>/shop id=:id /blog/:id 301</code></td>\n<td>...</td>\n</tr>\n<tr>\n<td>Proxying</td>\n<td>No</td>\n<td><code>/blog/*
    https://blog.my.domain/:splat 200</code></td>\n<td>...</td>\n</tr>\n<tr>\n<td>Domain-level
    redirects</td>\n<td>No</td>\n<td><code>workers.example.com/* workers.example.com/blog/:splat
    301</code></td>\n<td>...</td>\n</tr>\n<tr>\n<td>Redirect by country or language</td>\n<td>No</td>\n<td><code>/
    /us 302 Country=us</code></td>\n<td>...</td>\n</tr>\n<tr>\n<td>Redirect by cookie</td>\n<td>No</td>\n<td><code>/*
    /preview/:splat 302 Cookie=preview</code></td>\n<td>...</td>\n</tr>\n</tbody>\n</table>\n<blockquote>\n<p>Compare
    with\n<a href=\"https://developers.cloudflare.com/pages/platform/redirects/\">cloudflare-pages</a></p>\n</blockquote>\n<div
    class=\"admonition tip\">\n<p class=\"admonition-title\">Tip</p>\n<p>If you have
    a public site, pair this up with\n<a href=\"https://app.ahrefs.com/dashboard\">ahrefs</a>
    to keep up with pages that have\nmoved without you realizing.</p>\n</div>\n<p>!!
    class <h2 id='Redirect' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>Redirect <em class='small'>class</em></h2>\nDataClass to store the original
    and new url</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">Redirect <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Redirect</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"s2\">&quot;DataClass to store the original and new url&quot;</span>\n\n
    \           <span class=\"n\">original</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span>\n            <span class=\"n\">new</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span>\n            <span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"n\">Markata</span>\n            <span class=\"n\">model_config</span>
    <span class=\"o\">=</span> <span class=\"n\">ConfigDict</span><span class=\"p\">(</span><span
    class=\"n\">validate_assignment</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">arbitrary_types_allowed</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='RedirectsConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>RedirectsConfig <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">RedirectsConfig
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">RedirectsConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">assets_dir</span><span class=\"p\">:</span> <span
    class=\"n\">Path</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;static&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">redirects_file</span><span class=\"p\">:</span>
    <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;redirects_file&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">always</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"k\">def</span>
    <span class=\"nf\">default_redirects_file</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">cls</span><span class=\"p\">:</span> <span class=\"s2\">&quot;RedirectsConfig&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">:</span> <span
    class=\"n\">Path</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Path</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">v</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;assets_dir&quot;</span><span class=\"p\">])</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;_redirects&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Config <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">redirects</span><span class=\"p\">:</span> <span
    class=\"n\">RedirectsConfig</span> <span class=\"o\">=</span> <span class=\"n\">RedirectsConfig</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='config_model'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">config_model <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>save <em class='small'>function</em></h2>\nsaves an index.html in the directory
    called out by the redirect.</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">save <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            saves an index.html in the directory called out by the
    redirect.</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n            <span
    class=\"n\">redirects_file</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">redirects</span><span
    class=\"o\">.</span><span class=\"n\">redirects_file</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">redirects_file</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \               <span class=\"n\">raw_redirects</span> <span class=\"o\">=</span>
    <span class=\"n\">redirects_file</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">split</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span
    class=\"s2\">&quot;</span><span class=\"p\">)</span>\n            <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                <span class=\"n\">raw_redirects</span> <span
    class=\"o\">=</span> <span class=\"p\">[]</span>\n\n            <span class=\"n\">redirects</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span>\n                <span class=\"n\">Redirect</span><span
    class=\"p\">(</span><span class=\"n\">original</span><span class=\"o\">=</span><span
    class=\"n\">s</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">],</span> <span class=\"n\">new</span><span class=\"o\">=</span><span
    class=\"n\">s</span><span class=\"p\">[</span><span class=\"mi\">1</span><span
    class=\"p\">],</span> <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">)</span>\n                <span class=\"k\">for</span>
    <span class=\"n\">r</span> <span class=\"ow\">in</span> <span class=\"n\">raw_redirects</span>\n
    \               <span class=\"k\">if</span> <span class=\"s2\">&quot;*&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">r</span>
    <span class=\"ow\">and</span> <span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"n\">s</span> <span class=\"o\">:=</span> <span class=\"n\">r</span><span
    class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">())</span>
    <span class=\"o\">==</span> <span class=\"mi\">2</span> <span class=\"ow\">and</span>
    <span class=\"ow\">not</span> <span class=\"n\">r</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">()</span><span class=\"o\">.</span><span
    class=\"n\">startswith</span><span class=\"p\">(</span><span class=\"s2\">&quot;#&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"p\">]</span>\n\n            <span
    class=\"k\">if</span> <span class=\"s2\">&quot;redirect_template&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">:</span>\n                <span class=\"n\">template_file</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">&quot;redirect_template&quot;</span><span
    class=\"p\">)))</span>\n            <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">template_file</span> <span class=\"o\">=</span>
    <span class=\"n\">DEFAULT_REDIRECT_TEMPLATE</span>\n            <span class=\"n\">template</span>
    <span class=\"o\">=</span> <span class=\"n\">Template</span><span class=\"p\">(</span><span
    class=\"n\">template_file</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">())</span>\n\n            <span class=\"k\">for</span> <span class=\"n\">redirect</span>
    <span class=\"ow\">in</span> <span class=\"n\">redirects</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">file</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">output_dir</span> <span class=\"o\">/</span>
    <span class=\"n\">redirect</span><span class=\"o\">.</span><span class=\"n\">original</span><span
    class=\"o\">.</span><span class=\"n\">strip</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;index.html&quot;</span>\n                <span class=\"n\">file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">parents</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">exist_ok</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n                <span class=\"n\">file</span><span class=\"o\">.</span><span
    class=\"n\">write_text</span><span class=\"p\">(</span><span class=\"n\">template</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span><span
    class=\"n\">redirect</span><span class=\"o\">.</span><span class=\"n\">dict</span><span
    class=\"p\">(),</span> <span class=\"n\">config</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='default_redirects_file'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_redirects_file
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">default_redirects_file <em
    class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">default_redirects_file</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">cls</span><span class=\"p\">:</span> <span class=\"s2\">&quot;RedirectsConfig&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">:</span> <span
    class=\"n\">Path</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Path</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">v</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;assets_dir&quot;</span><span class=\"p\">])</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;_redirects&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/redirects
title: Redirects.Py


---

Creates redirects for times when your backend server can't.

## Configuration

Enable the redirect hook by adding it to your list of hooks.

``` toml
[markata]

# Were you keep static assets to copy into the project, default is static
# the assets_dir will set the default _redirects file directory
assets_dir = "static"

# You can override the default redirects file location
redirects = static/_redirects

hooks = [
   # creates redirects from static/_redirects file
   "markata.plugins.redirects",
   # copies your static assets into the output_dir (default: `markout`)
   "markata.plugins.copy_assets",
  ...
]
```

## Syntax

Your `_redirects` file is a simplified version of what services like cloudflare
pages or netlify use.  In fact you can use the same redirects file!

Here is an example that will redirect `/old` to `/new` and `/CHANGELOG` to
`/changelog`

```
/old        /new
/CHANGELOG  /changelog
```

## Limitations
_no splats_

Since it is static generated this plugin cannot cover *'s.  * or splat
redirects need to be taken care of server side.  It also cannot change the http
code, this is only

## Features

The features of markata.plugins.redirect is pretty limited since it is
implemented only as a static page.  Other features require server side
implementation.

| Feature                             | Support | Example                                                         | Notes                                                                                             |
| ----------------------------------- | ------- | --------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| Force                               | Yes     | `/pagethatexists /otherpage`                                    | Creates an index.html with http-equiv and canonical                                               |
| Redirects (301, 302, 303, 307, 308) | No      | `/home / 301`                                                   | Ignored, requires server side implementation                                                      |
| Rewrites (other status codes)       | No      | `/blog/* /blog/404.html 404`                                    | ...                                                                                               |
| Splats                              | No      | `/blog/* /blog/:splat`                                          | ...                                                                                               |
| Placeholders                        | No      | `/blog/:year/:month/:date/:slug /news/:year/:month/:date/:slug` | ...                                                                                               |
| Query Parameters                    | No      | `/shop id=:id /blog/:id 301`                                    | ...                                                                                               |
| Proxying                            | No      | `/blog/* https://blog.my.domain/:splat 200`                     | ...                                                                                               |
| Domain-level redirects              | No      | `workers.example.com/* workers.example.com/blog/:splat 301`     | ...                                                                                               |
| Redirect by country or language     | No      | `/ /us 302 Country=us`                                          | ...                                                                                               |
| Redirect by cookie                  | No      | `/* /preview/:splat 302 Cookie=preview`                        | ...                                                                                               |

> Compare with
> [cloudflare-pages](https://developers.cloudflare.com/pages/platform/redirects/)

!!! tip
    If you have a public site, pair this up with
    [ahrefs](https://app.ahrefs.com/dashboard) to keep up with pages that have
    moved without you realizing.


!! class <h2 id='Redirect' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Redirect <em class='small'>class</em></h2>
    DataClass to store the original and new url
???+ source "Redirect <em class='small'>source</em>"

```python

        class Redirect(pydantic.BaseModel):
            "DataClass to store the original and new url"

            original: str
            new: str
            markata: Markata
            model_config = ConfigDict(validate_assignment=True, arbitrary_types_allowed=True)
```


!! class <h2 id='RedirectsConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>RedirectsConfig <em class='small'>class</em></h2>

???+ source "RedirectsConfig <em class='small'>source</em>"

```python

        class RedirectsConfig(pydantic.BaseModel):
            assets_dir: Path = Path("static")
            redirects_file: Optional[Path] = None

            @pydantic.validator("redirects_file", always=True)
            def default_redirects_file(
                cls: "RedirectsConfig", v: Path, *, values: Dict
            ) -> Path:
                if not v:
                    return Path(values["assets_dir"]) / "_redirects"
                return v
```


!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>

???+ source "Config <em class='small'>source</em>"

```python

        class Config(pydantic.BaseModel):
            redirects: RedirectsConfig = RedirectsConfig()
```


!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>

???+ source "config_model <em class='small'>source</em>"

```python

        def config_model(markata: "Markata") -> None:
            markata.config_models.append(Config)
```


!! function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>
    saves an index.html in the directory called out by the redirect.
???+ source "save <em class='small'>source</em>"

```python

        def save(markata: "Markata") -> None:
            """
            saves an index.html in the directory called out by the redirect.
            """
            redirects_file = Path(markata.config.redirects.redirects_file)
            if redirects_file.exists():
                raw_redirects = redirects_file.read_text().split("\n")
            else:
                raw_redirects = []

            redirects = [
                Redirect(original=s[0], new=s[1], markata=markata)
                for r in raw_redirects
                if "*" not in r and len(s := r.split()) == 2 and not r.strip().startswith("#")
            ]

            if "redirect_template" in markata.config:
                template_file = Path(str(markata.config.get("redirect_template")))
            else:
                template_file = DEFAULT_REDIRECT_TEMPLATE
            template = Template(template_file.read_text())

            for redirect in redirects:
                file = markata.config.output_dir / redirect.original.strip("/") / "index.html"
                file.parent.mkdir(parents=True, exist_ok=True)
                file.write_text(template.render(redirect.dict(), config=markata.config))
```


!! method <h2 id='default_redirects_file' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_redirects_file <em class='small'>method</em></h2>

???+ source "default_redirects_file <em class='small'>source</em>"

```python

        def default_redirects_file(
                cls: "RedirectsConfig", v: Path, *, values: Dict
            ) -> Path:
                if not v:
                    return Path(values["assets_dir"]) / "_redirects"
                return v
```

