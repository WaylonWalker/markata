---
content: "Sets the articles `output_html` path, and saves the article's `html` to
  the\n`output_html` file.\n\n##  Ouptut Directory\n\nOutput will always be written
  inside of the configured `output_dir`\n\n```toml\n[markata]\n# markout is the default,
  but you can override it in your markata.toml file\noutput_dir = \"markout\"\n```\n\n##
  Explicityly set the output\n\nmarkata will save the articles `html` to the `output_html`
  specified in the\narticles metadata, loaded from frontmatter.\n\n## 404 example
  use case\n\nHere is an example use case of explicitly setting the output_html.  By
  default\nmarkata will turn `pages/404.md` into `markout/404/index.html`, but many\nhosting
  providers look for a 404.html to redirect the user to when a page is\nnot found.\n\n```markdown\n---\ntitle:
  Whoops that page was not found\ndescription: 404, looks like we can't find the page
  you are looking for\noutput_html: 404.html\n\n---\n\n404, looks like we can't find
  the page you are looking for.  Try one of these\npages.\n\n<ul>\n{% for post in\n
  \   markata.map(\n        'post',\n        filter='\"markata\" not in slug and \"tests\"
  not in slug and \"404\" not in slug'\n        )\n %}\n    <li><a href=\"{{ post.slug
  }}\">{{ post.title or \"CHANGELOG\" }}</a></li>\n{% endfor %}\n</ul>\n```\n\n##
  Index.md is the one special case\n\nIf you have a file `pages/index.md` it will
  become `markout/index.html` rather\nthan `markout/index/inject.html` This is one
  of the primary ways that markata\nlets you [make your home page](https://markata.dev/home-page/)\n\n\n!!
  class <h2 id='OutputHTML' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>OutputHTML <em class='small'>class</em></h2>\n\n???+ source \"OutputHTML
  <em class='small'>source</em>\"\n\n```python\n\n        class OutputHTML(pydantic.BaseModel):\n
  \           markata: Any = Field(None, exclude=True)\n            path: Path\n            slug:
  str = None\n            output_html: Path = None\n\n            @pydantic.validator(\"slug\",
  pre=True, always=True)\n            @classmethod\n            def default_slug(cls,
  v, *, values):\n                if v is None:\n                    return slugify(str(values[\"path\"].stem))\n
  \               return v\n\n            @pydantic.validator(\"output_html\", pre=True,
  always=True)\n            def default_output_html(\n                cls: \"OutputHTML\",
  v: Optional[Path], *, values: Dict\n            ) -> Path:\n                if isinstance(v,
  str):\n                    v = Path(v)\n                if v is not None:\n                    return
  v\n                if \"slug\" not in values:\n                    for validator
  in cls.__validators__[\"slug\"]:\n                        values[\"slug\"] = validator.func(cls,
  v, values=values)\n\n                if values[\"slug\"] == \"index\":\n                    return
  cls.markata.config.output_dir / \"index.html\"\n                return cls.markata.config.output_dir
  / values[\"slug\"] / \"index.html\"\n\n            @pydantic.validator(\"output_html\")\n
  \           def output_html_relative(\n                cls: \"OutputHTML\", v: Optional[Path],
  *, values: Dict\n            ) -> Path:\n                if isinstance(v, str):\n
  \                   v = Path(v)\n                if cls.markata.config.output_dir.absolute()
  not in v.absolute().parents:\n                    return cls.markata.config.output_dir
  / v\n                return v\n\n            @pydantic.validator(\"output_html\")\n
  \           def output_html_exists(\n                cls: \"OutputHTML\", v: Optional[Path],
  *, values: Dict\n            ) -> Path:\n                if isinstance(v, str):\n
  \                   v = Path(v)\n                if not v.parent.exists():\n                    v.parent.mkdir(parents=True,
  exist_ok=True)\n                return v\n```\n\n\n!! function <h2 id='post_model'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2>\n\n???+
  source \"post_model <em class='small'>source</em>\"\n\n```python\n\n        def
  post_model(markata: \"Markata\") -> None:\n            markata.post_models.append(OutputHTML)\n```\n\n\n!!
  function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save
  <em class='small'>function</em></h2>\n    Saves all the articles to their set `output_html`
  location if that location\n    is relative to the specified `output_dir`.  If its
  not relative to the\n    `output_dir` it will log an error and move on.\n???+ source
  \"save <em class='small'>source</em>\"\n\n```python\n\n        def save(markata:
  \"Markata\") -> None:\n            \"\"\"\n            Saves all the articles to
  their set `output_html` location if that location\n            is relative to the
  specified `output_dir`.  If its not relative to the\n            `output_dir` it
  will log an error and move on.\n            \"\"\"\n\n            for article in
  markata.articles:\n                if article.html is None:\n                    continue\n
  \               if isinstance(article.html, str):\n                    article.output_html.write_text(article.html)\n
  \               if isinstance(article.html, Dict):\n                    for slug,
  html in article.html.items():\n                        if slug == \"index\":\n                            slug
  = \"\"\n                            output_html = article.output_html\n                        elif
  \".\" in slug:\n                            output_html = article.output_html.parent
  / slug\n                        else:\n                            slug = slugify(slug)\n
  \                           output_html = article.output_html.parent / slug / \"index.html\"\n\n
  \                       output_html.parent.mkdir(parents=True, exist_ok=True)\n
  \                       output_html.write_text(html)\n```\n\n\n!! method <h2 id='default_slug'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_slug <em class='small'>method</em></h2>\n\n???+
  source \"default_slug <em class='small'>source</em>\"\n\n```python\n\n        def
  default_slug(cls, v, *, values):\n                if v is None:\n                    return
  slugify(str(values[\"path\"].stem))\n                return v\n```\n\n\n!! method
  <h2 id='default_output_html' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>default_output_html <em class='small'>method</em></h2>\n\n???+ source \"default_output_html
  <em class='small'>source</em>\"\n\n```python\n\n        def default_output_html(\n
  \               cls: \"OutputHTML\", v: Optional[Path], *, values: Dict\n            )
  -> Path:\n                if isinstance(v, str):\n                    v = Path(v)\n
  \               if v is not None:\n                    return v\n                if
  \"slug\" not in values:\n                    for validator in cls.__validators__[\"slug\"]:\n
  \                       values[\"slug\"] = validator.func(cls, v, values=values)\n\n
  \               if values[\"slug\"] == \"index\":\n                    return cls.markata.config.output_dir
  / \"index.html\"\n                return cls.markata.config.output_dir / values[\"slug\"]
  / \"index.html\"\n```\n\n\n!! method <h2 id='output_html_relative' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>output_html_relative <em class='small'>method</em></h2>\n\n???+
  source \"output_html_relative <em class='small'>source</em>\"\n\n```python\n\n        def
  output_html_relative(\n                cls: \"OutputHTML\", v: Optional[Path], *,
  values: Dict\n            ) -> Path:\n                if isinstance(v, str):\n                    v
  = Path(v)\n                if cls.markata.config.output_dir.absolute() not in v.absolute().parents:\n
  \                   return cls.markata.config.output_dir / v\n                return
  v\n```\n\n\n!! method <h2 id='output_html_exists' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>output_html_exists <em class='small'>method</em></h2>\n\n???+ source \"output_html_exists
  <em class='small'>source</em>\"\n\n```python\n\n        def output_html_exists(\n
  \               cls: \"OutputHTML\", v: Optional[Path], *, values: Dict\n            )
  -> Path:\n                if isinstance(v, str):\n                    v = Path(v)\n
  \               if not v.parent.exists():\n                    v.parent.mkdir(parents=True,
  exist_ok=True)\n                return v\n```\n\n"
date: 0001-01-01
description: Sets the articles  Output will always be written inside of the configured  markata
  will save the articles  Here is an example use case of explicitly setting the
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Publish_Html.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Sets the articles  Output will always
    be written inside of the configured  markata will save the articles  Here is an
    example use case of explicitly setting the\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Publish_Html.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Sets the articles  Output will always
    be written inside of the configured  markata will save the articles  Here is an
    example use case of explicitly setting the\" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Publish_Html.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Sets the articles
    <code>output_html</code> path, and saves the article's <code>html</code> to the\n<code>output_html</code>
    file.</p>\n<h2 id=\"ouptut-directory\">Ouptut Directory <a class=\"header-anchor\"
    href=\"#ouptut-directory\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Output will always be
    written inside of the configured <code>output_dir</code></p>\n<pre class='wrapper'>\n\n<div
    class='copy-wrapper'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"c1\"># markout is the default, but you can override it in your markata.toml
    file</span>\n<span class=\"n\">output_dir</span><span class=\"w\"> </span><span
    class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;markout&quot;</span>\n</pre></div>\n\n</pre>\n\n<h2
    id=\"explicityly-set-the-output\">Explicityly set the output <a class=\"header-anchor\"
    href=\"#explicityly-set-the-output\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>markata will save the
    articles <code>html</code> to the <code>output_html</code> specified in the\narticles
    metadata, loaded from frontmatter.</p>\n<h2 id=\"404-example-use-case\">404 example
    use case <a class=\"header-anchor\" href=\"#404-example-use-case\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Here is an example use
    case of explicitly setting the output_html.  By default\nmarkata will turn <code>pages/404.md</code>
    into <code>markout/404/index.html</code>, but many\nhosting providers look for
    a 404.html to redirect the user to when a page is\nnot found.</p>\n<pre class='wrapper'>\n\n<div
    class='copy-wrapper'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
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
    \       \n<div class=\"highlight\"><pre><span></span>---\ntitle: Whoops that page
    was not found\ndescription: 404, looks like we can&#39;t find the page you are
    looking for\noutput_html: 404.html\n\n---\n\n404, looks like we can&#39;t find
    the page you are looking for.  Try one of these\npages.\n\n&lt;ul&gt;\n{% for
    post in\n    markata.map(\n        &#39;post&#39;,\n        filter=&#39;&quot;markata&quot;
    not in slug and &quot;tests&quot; not in slug and &quot;404&quot; not in slug&#39;\n
    \       )\n %}\n    &lt;li&gt;&lt;a href=&quot;{{ post.slug }}&quot;&gt;{{ post.title
    or &quot;CHANGELOG&quot; }}&lt;/a&gt;&lt;/li&gt;\n{% endfor %}\n&lt;/ul&gt;\n</pre></div>\n\n</pre>\n\n<h2
    id=\"indexmd-is-the-one-special-case\"><a href=\"http://Index.md\">Index.md</a>
    is the one special case <a class=\"header-anchor\" href=\"#indexmd-is-the-one-special-case\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>If you have a file <code>pages/index.md</code>
    it will become <code>markout/index.html</code> rather\nthan <code>markout/index/inject.html</code>
    This is one of the primary ways that markata\nlets you <a href=\"https://markata.dev/home-page/\">make
    your home page</a></p>\n<p>!! class <h2 id='OutputHTML' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>OutputHTML <em class='small'>class</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">OutputHTML
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
    <span class=\"nc\">OutputHTML</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Any</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"kc\">None</span><span class=\"p\">,</span> <span class=\"n\">exclude</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">path</span><span class=\"p\">:</span> <span class=\"n\">Path</span>\n
    \           <span class=\"n\">slug</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n            <span class=\"n\">output_html</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"nd\">@classmethod</span>\n            <span
    class=\"k\">def</span> <span class=\"nf\">default_slug</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">v</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">slugify</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;path&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">stem</span><span class=\"p\">))</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">v</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;output_html&quot;</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">default_output_html</span><span
    class=\"p\">(</span>\n                <span class=\"bp\">cls</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;OutputHTML&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Path</span><span class=\"p\">],</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">:</span> <span class=\"n\">Dict</span>\n            <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">v</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">v</span>\n
    \               <span class=\"k\">if</span> <span class=\"s2\">&quot;slug&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">values</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">for</span> <span class=\"n\">validator</span>
    <span class=\"ow\">in</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">__validators__</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">]:</span>\n                        <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"n\">validator</span><span class=\"o\">.</span><span
    class=\"n\">func</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"o\">=</span><span class=\"n\">values</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;index&quot;</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span>
    <span class=\"o\">/</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;index.html&quot;</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;output_html&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"k\">def</span> <span class=\"nf\">output_html_relative</span><span class=\"p\">(</span>\n
    \               <span class=\"bp\">cls</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;OutputHTML&quot;</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">Path</span><span class=\"p\">],</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Path</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                    <span class=\"n\">v</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">)</span>\n                <span class=\"k\">if</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"o\">.</span><span class=\"n\">absolute</span><span class=\"p\">()</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">absolute</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">parents</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span>
    <span class=\"o\">/</span> <span class=\"n\">v</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">v</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;output_html&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"k\">def</span> <span class=\"nf\">output_html_exists</span><span class=\"p\">(</span>\n
    \               <span class=\"bp\">cls</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;OutputHTML&quot;</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">Path</span><span class=\"p\">],</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Path</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                    <span class=\"n\">v</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">)</span>\n                <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"n\">v</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">v</span><span class=\"o\">.</span><span
    class=\"n\">parent</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
    class=\"p\">(</span><span class=\"n\">parents</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">exist_ok</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>post_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">post_model
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">post_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">post_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">OutputHTML</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>save <em class='small'>function</em></h2>\nSaves all the articles to their
    set <code>output_html</code> location if that location\nis relative to the specified
    <code>output_dir</code>.  If its not relative to the\n<code>output_dir</code>
    it will log an error and move on.</p>\n<div class=\"admonition source is-collapsible
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
    class=\"sd\">            Saves all the articles to their set `output_html` location
    if that location</span>\n<span class=\"sd\">            is relative to the specified
    `output_dir`.  If its not relative to the</span>\n<span class=\"sd\">            `output_dir`
    it will log an error and move on.</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n\n
    \           <span class=\"k\">for</span> <span class=\"n\">article</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">:</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">html</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">continue</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">html</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">output_html</span><span
    class=\"o\">.</span><span class=\"n\">write_text</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">html</span><span
    class=\"p\">)</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">html</span><span class=\"p\">,</span> <span class=\"n\">Dict</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">for</span> <span
    class=\"n\">slug</span><span class=\"p\">,</span> <span class=\"n\">html</span>
    <span class=\"ow\">in</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">html</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n                        <span class=\"k\">if</span> <span
    class=\"n\">slug</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;index&quot;</span><span
    class=\"p\">:</span>\n                            <span class=\"n\">slug</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n                            <span
    class=\"n\">output_html</span> <span class=\"o\">=</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">output_html</span>\n                        <span
    class=\"k\">elif</span> <span class=\"s2\">&quot;.&quot;</span> <span class=\"ow\">in</span>
    <span class=\"n\">slug</span><span class=\"p\">:</span>\n                            <span
    class=\"n\">output_html</span> <span class=\"o\">=</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">output_html</span><span class=\"o\">.</span><span
    class=\"n\">parent</span> <span class=\"o\">/</span> <span class=\"n\">slug</span>\n
    \                       <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                           <span class=\"n\">slug</span> <span class=\"o\">=</span>
    <span class=\"n\">slugify</span><span class=\"p\">(</span><span class=\"n\">slug</span><span
    class=\"p\">)</span>\n                            <span class=\"n\">output_html</span>
    <span class=\"o\">=</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">output_html</span><span class=\"o\">.</span><span class=\"n\">parent</span>
    <span class=\"o\">/</span> <span class=\"n\">slug</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;index.html&quot;</span>\n\n                        <span
    class=\"n\">output_html</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
    class=\"o\">.</span><span class=\"n\">mkdir</span><span class=\"p\">(</span><span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n                        <span
    class=\"n\">output_html</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
    class=\"p\">(</span><span class=\"n\">html</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='default_slug' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>default_slug <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">default_slug
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">default_slug</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">slugify</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">stem</span><span
    class=\"p\">))</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='default_output_html' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>default_output_html <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">default_output_html
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">default_output_html</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">cls</span><span class=\"p\">:</span> <span class=\"s2\">&quot;OutputHTML&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
    class=\"p\">],</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">:</span> <span class=\"n\">Dict</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">v</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">v</span>\n
    \               <span class=\"k\">if</span> <span class=\"s2\">&quot;slug&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">values</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">for</span> <span class=\"n\">validator</span>
    <span class=\"ow\">in</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">__validators__</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">]:</span>\n                        <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"n\">validator</span><span class=\"o\">.</span><span
    class=\"n\">func</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"o\">=</span><span class=\"n\">values</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;index&quot;</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span>
    <span class=\"o\">/</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;index.html&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='output_html_relative' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>output_html_relative <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">output_html_relative
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">output_html_relative</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">cls</span><span class=\"p\">:</span> <span class=\"s2\">&quot;OutputHTML&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
    class=\"p\">],</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">:</span> <span class=\"n\">Dict</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">v</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                <span
    class=\"k\">if</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">output_dir</span><span class=\"o\">.</span><span
    class=\"n\">absolute</span><span class=\"p\">()</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"n\">v</span><span class=\"o\">.</span><span
    class=\"n\">absolute</span><span class=\"p\">()</span><span class=\"o\">.</span><span
    class=\"n\">parents</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">output_dir</span> <span class=\"o\">/</span>
    <span class=\"n\">v</span>\n                <span class=\"k\">return</span> <span
    class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='output_html_exists'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>output_html_exists
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">output_html_exists <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">output_html_exists</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">cls</span><span class=\"p\">:</span> <span class=\"s2\">&quot;OutputHTML&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
    class=\"p\">],</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">:</span> <span class=\"n\">Dict</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">v</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                <span
    class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n                    <span
    class=\"n\">v</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
    class=\"o\">.</span><span class=\"n\">mkdir</span><span class=\"p\">(</span><span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Publish_Html.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Sets the articles  Output will always
    be written inside of the configured  markata will save the articles  Here is an
    example use case of explicitly setting the\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Publish_Html.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Sets the articles  Output will always
    be written inside of the configured  markata will save the articles  Here is an
    example use case of explicitly setting the\" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Publish_Html.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Publish_Html.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>Sets
    the articles <code>output_html</code> path, and saves the article's <code>html</code>
    to the\n<code>output_html</code> file.</p>\n<h2 id=\"ouptut-directory\">Ouptut
    Directory <a class=\"header-anchor\" href=\"#ouptut-directory\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Output will always be
    written inside of the configured <code>output_dir</code></p>\n<pre class='wrapper'>\n\n<div
    class='copy-wrapper'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"c1\"># markout is the default, but you can override it in your markata.toml
    file</span>\n<span class=\"n\">output_dir</span><span class=\"w\"> </span><span
    class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;markout&quot;</span>\n</pre></div>\n\n</pre>\n\n<h2
    id=\"explicityly-set-the-output\">Explicityly set the output <a class=\"header-anchor\"
    href=\"#explicityly-set-the-output\"><svg class=\"heading-permalink\" aria-hidden=\"true\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>markata will save the
    articles <code>html</code> to the <code>output_html</code> specified in the\narticles
    metadata, loaded from frontmatter.</p>\n<h2 id=\"404-example-use-case\">404 example
    use case <a class=\"header-anchor\" href=\"#404-example-use-case\"><svg class=\"heading-permalink\"
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Here is an example use
    case of explicitly setting the output_html.  By default\nmarkata will turn <code>pages/404.md</code>
    into <code>markout/404/index.html</code>, but many\nhosting providers look for
    a 404.html to redirect the user to when a page is\nnot found.</p>\n<pre class='wrapper'>\n\n<div
    class='copy-wrapper'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
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
    \       \n<div class=\"highlight\"><pre><span></span>---\ntitle: Whoops that page
    was not found\ndescription: 404, looks like we can&#39;t find the page you are
    looking for\noutput_html: 404.html\n\n---\n\n404, looks like we can&#39;t find
    the page you are looking for.  Try one of these\npages.\n\n&lt;ul&gt;\n{% for
    post in\n    markata.map(\n        &#39;post&#39;,\n        filter=&#39;&quot;markata&quot;
    not in slug and &quot;tests&quot; not in slug and &quot;404&quot; not in slug&#39;\n
    \       )\n %}\n    &lt;li&gt;&lt;a href=&quot;{{ post.slug }}&quot;&gt;{{ post.title
    or &quot;CHANGELOG&quot; }}&lt;/a&gt;&lt;/li&gt;\n{% endfor %}\n&lt;/ul&gt;\n</pre></div>\n\n</pre>\n\n<h2
    id=\"indexmd-is-the-one-special-case\"><a href=\"http://Index.md\">Index.md</a>
    is the one special case <a class=\"header-anchor\" href=\"#indexmd-is-the-one-special-case\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>If you have a file <code>pages/index.md</code>
    it will become <code>markout/index.html</code> rather\nthan <code>markout/index/inject.html</code>
    This is one of the primary ways that markata\nlets you <a href=\"https://markata.dev/home-page/\">make
    your home page</a></p>\n<p>!! class <h2 id='OutputHTML' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>OutputHTML <em class='small'>class</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">OutputHTML
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
    <span class=\"nc\">OutputHTML</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Any</span>
    <span class=\"o\">=</span> <span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"kc\">None</span><span class=\"p\">,</span> <span class=\"n\">exclude</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">path</span><span class=\"p\">:</span> <span class=\"n\">Path</span>\n
    \           <span class=\"n\">slug</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n            <span class=\"n\">output_html</span><span
    class=\"p\">:</span> <span class=\"n\">Path</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"nd\">@classmethod</span>\n            <span
    class=\"k\">def</span> <span class=\"nf\">default_slug</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">):</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">v</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">slugify</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;path&quot;</span><span class=\"p\">]</span><span class=\"o\">.</span><span
    class=\"n\">stem</span><span class=\"p\">))</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">v</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;output_html&quot;</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
    class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">default_output_html</span><span
    class=\"p\">(</span>\n                <span class=\"bp\">cls</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;OutputHTML&quot;</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">:</span> <span class=\"n\">Optional</span><span
    class=\"p\">[</span><span class=\"n\">Path</span><span class=\"p\">],</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">:</span> <span class=\"n\">Dict</span>\n            <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">v</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">v</span>\n
    \               <span class=\"k\">if</span> <span class=\"s2\">&quot;slug&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">values</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">for</span> <span class=\"n\">validator</span>
    <span class=\"ow\">in</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">__validators__</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">]:</span>\n                        <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"n\">validator</span><span class=\"o\">.</span><span
    class=\"n\">func</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"o\">=</span><span class=\"n\">values</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;index&quot;</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span>
    <span class=\"o\">/</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;index.html&quot;</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;output_html&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"k\">def</span> <span class=\"nf\">output_html_relative</span><span class=\"p\">(</span>\n
    \               <span class=\"bp\">cls</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;OutputHTML&quot;</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">Path</span><span class=\"p\">],</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Path</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                    <span class=\"n\">v</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">)</span>\n                <span class=\"k\">if</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
    class=\"o\">.</span><span class=\"n\">absolute</span><span class=\"p\">()</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">absolute</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">parents</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span>
    <span class=\"o\">/</span> <span class=\"n\">v</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">v</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;output_html&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"k\">def</span> <span class=\"nf\">output_html_exists</span><span class=\"p\">(</span>\n
    \               <span class=\"bp\">cls</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;OutputHTML&quot;</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">Path</span><span class=\"p\">],</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Path</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">):</span>\n                    <span class=\"n\">v</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">)</span>\n                <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"n\">v</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">v</span><span class=\"o\">.</span><span
    class=\"n\">parent</span><span class=\"o\">.</span><span class=\"n\">mkdir</span><span
    class=\"p\">(</span><span class=\"n\">parents</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">exist_ok</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>post_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">post_model
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">post_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">post_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">OutputHTML</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>save <em class='small'>function</em></h2>\nSaves all the articles to their
    set <code>output_html</code> location if that location\nis relative to the specified
    <code>output_dir</code>.  If its not relative to the\n<code>output_dir</code>
    it will log an error and move on.</p>\n<div class=\"admonition source is-collapsible
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
    class=\"sd\">            Saves all the articles to their set `output_html` location
    if that location</span>\n<span class=\"sd\">            is relative to the specified
    `output_dir`.  If its not relative to the</span>\n<span class=\"sd\">            `output_dir`
    it will log an error and move on.</span>\n<span class=\"sd\">            &quot;&quot;&quot;</span>\n\n
    \           <span class=\"k\">for</span> <span class=\"n\">article</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">articles</span><span class=\"p\">:</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">html</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">continue</span>\n                <span class=\"k\">if</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">html</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">output_html</span><span
    class=\"o\">.</span><span class=\"n\">write_text</span><span class=\"p\">(</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">html</span><span
    class=\"p\">)</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">html</span><span class=\"p\">,</span> <span class=\"n\">Dict</span><span
    class=\"p\">):</span>\n                    <span class=\"k\">for</span> <span
    class=\"n\">slug</span><span class=\"p\">,</span> <span class=\"n\">html</span>
    <span class=\"ow\">in</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">html</span><span class=\"o\">.</span><span class=\"n\">items</span><span
    class=\"p\">():</span>\n                        <span class=\"k\">if</span> <span
    class=\"n\">slug</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;index&quot;</span><span
    class=\"p\">:</span>\n                            <span class=\"n\">slug</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;&quot;</span>\n                            <span
    class=\"n\">output_html</span> <span class=\"o\">=</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">output_html</span>\n                        <span
    class=\"k\">elif</span> <span class=\"s2\">&quot;.&quot;</span> <span class=\"ow\">in</span>
    <span class=\"n\">slug</span><span class=\"p\">:</span>\n                            <span
    class=\"n\">output_html</span> <span class=\"o\">=</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">output_html</span><span class=\"o\">.</span><span
    class=\"n\">parent</span> <span class=\"o\">/</span> <span class=\"n\">slug</span>\n
    \                       <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                           <span class=\"n\">slug</span> <span class=\"o\">=</span>
    <span class=\"n\">slugify</span><span class=\"p\">(</span><span class=\"n\">slug</span><span
    class=\"p\">)</span>\n                            <span class=\"n\">output_html</span>
    <span class=\"o\">=</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">output_html</span><span class=\"o\">.</span><span class=\"n\">parent</span>
    <span class=\"o\">/</span> <span class=\"n\">slug</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;index.html&quot;</span>\n\n                        <span
    class=\"n\">output_html</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
    class=\"o\">.</span><span class=\"n\">mkdir</span><span class=\"p\">(</span><span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n                        <span
    class=\"n\">output_html</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
    class=\"p\">(</span><span class=\"n\">html</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='default_slug' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>default_slug <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">default_slug
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">default_slug</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">slugify</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;path&quot;</span><span
    class=\"p\">]</span><span class=\"o\">.</span><span class=\"n\">stem</span><span
    class=\"p\">))</span>\n                <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='default_output_html' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>default_output_html <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">default_output_html
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">default_output_html</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">cls</span><span class=\"p\">:</span> <span class=\"s2\">&quot;OutputHTML&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
    class=\"p\">],</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">:</span> <span class=\"n\">Dict</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">v</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">v</span>\n
    \               <span class=\"k\">if</span> <span class=\"s2\">&quot;slug&quot;</span>
    <span class=\"ow\">not</span> <span class=\"ow\">in</span> <span class=\"n\">values</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">for</span> <span class=\"n\">validator</span>
    <span class=\"ow\">in</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">__validators__</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">]:</span>\n                        <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"n\">validator</span><span class=\"o\">.</span><span
    class=\"n\">func</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"o\">=</span><span class=\"n\">values</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">==</span> <span class=\"s2\">&quot;index&quot;</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span>
    <span class=\"o\">/</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;index.html&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='output_html_relative' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>output_html_relative <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">output_html_relative
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
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">output_html_relative</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">cls</span><span class=\"p\">:</span> <span class=\"s2\">&quot;OutputHTML&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
    class=\"p\">],</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">:</span> <span class=\"n\">Dict</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">v</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                <span
    class=\"k\">if</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">output_dir</span><span class=\"o\">.</span><span
    class=\"n\">absolute</span><span class=\"p\">()</span> <span class=\"ow\">not</span>
    <span class=\"ow\">in</span> <span class=\"n\">v</span><span class=\"o\">.</span><span
    class=\"n\">absolute</span><span class=\"p\">()</span><span class=\"o\">.</span><span
    class=\"n\">parents</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">output_dir</span> <span class=\"o\">/</span>
    <span class=\"n\">v</span>\n                <span class=\"k\">return</span> <span
    class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='output_html_exists'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>output_html_exists
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">output_html_exists <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">output_html_exists</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">cls</span><span class=\"p\">:</span> <span class=\"s2\">&quot;OutputHTML&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">:</span> <span
    class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
    class=\"p\">],</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">:</span> <span class=\"n\">Dict</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"nb\">str</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">v</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n                <span
    class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n                    <span
    class=\"n\">v</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
    class=\"o\">.</span><span class=\"n\">mkdir</span><span class=\"p\">(</span><span
    class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">,</span> <span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n                <span class=\"k\">return</span>
    <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/publish-html
title: Publish_Html.Py


---

Sets the articles `output_html` path, and saves the article's `html` to the
`output_html` file.

##  Ouptut Directory

Output will always be written inside of the configured `output_dir`

```toml
[markata]
# markout is the default, but you can override it in your markata.toml file
output_dir = "markout"
```

## Explicityly set the output

markata will save the articles `html` to the `output_html` specified in the
articles metadata, loaded from frontmatter.

## 404 example use case

Here is an example use case of explicitly setting the output_html.  By default
markata will turn `pages/404.md` into `markout/404/index.html`, but many
hosting providers look for a 404.html to redirect the user to when a page is
not found.

```markdown
---
title: Whoops that page was not found
description: 404, looks like we can't find the page you are looking for
output_html: 404.html

---

404, looks like we can't find the page you are looking for.  Try one of these
pages.

<ul>
{% for post in
    markata.map(
        'post',
        filter='"markata" not in slug and "tests" not in slug and "404" not in slug'
        )
 %}
    <li><a href="{{ post.slug }}">{{ post.title or "CHANGELOG" }}</a></li>
{% endfor %}
</ul>
```

## Index.md is the one special case

If you have a file `pages/index.md` it will become `markout/index.html` rather
than `markout/index/inject.html` This is one of the primary ways that markata
lets you [make your home page](https://markata.dev/home-page/)


!! class <h2 id='OutputHTML' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>OutputHTML <em class='small'>class</em></h2>

???+ source "OutputHTML <em class='small'>source</em>"

```python

        class OutputHTML(pydantic.BaseModel):
            markata: Any = Field(None, exclude=True)
            path: Path
            slug: str = None
            output_html: Path = None

            @pydantic.validator("slug", pre=True, always=True)
            @classmethod
            def default_slug(cls, v, *, values):
                if v is None:
                    return slugify(str(values["path"].stem))
                return v

            @pydantic.validator("output_html", pre=True, always=True)
            def default_output_html(
                cls: "OutputHTML", v: Optional[Path], *, values: Dict
            ) -> Path:
                if isinstance(v, str):
                    v = Path(v)
                if v is not None:
                    return v
                if "slug" not in values:
                    for validator in cls.__validators__["slug"]:
                        values["slug"] = validator.func(cls, v, values=values)

                if values["slug"] == "index":
                    return cls.markata.config.output_dir / "index.html"
                return cls.markata.config.output_dir / values["slug"] / "index.html"

            @pydantic.validator("output_html")
            def output_html_relative(
                cls: "OutputHTML", v: Optional[Path], *, values: Dict
            ) -> Path:
                if isinstance(v, str):
                    v = Path(v)
                if cls.markata.config.output_dir.absolute() not in v.absolute().parents:
                    return cls.markata.config.output_dir / v
                return v

            @pydantic.validator("output_html")
            def output_html_exists(
                cls: "OutputHTML", v: Optional[Path], *, values: Dict
            ) -> Path:
                if isinstance(v, str):
                    v = Path(v)
                if not v.parent.exists():
                    v.parent.mkdir(parents=True, exist_ok=True)
                return v
```


!! function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2>

???+ source "post_model <em class='small'>source</em>"

```python

        def post_model(markata: "Markata") -> None:
            markata.post_models.append(OutputHTML)
```


!! function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>
    Saves all the articles to their set `output_html` location if that location
    is relative to the specified `output_dir`.  If its not relative to the
    `output_dir` it will log an error and move on.
???+ source "save <em class='small'>source</em>"

```python

        def save(markata: "Markata") -> None:
            """
            Saves all the articles to their set `output_html` location if that location
            is relative to the specified `output_dir`.  If its not relative to the
            `output_dir` it will log an error and move on.
            """

            for article in markata.articles:
                if article.html is None:
                    continue
                if isinstance(article.html, str):
                    article.output_html.write_text(article.html)
                if isinstance(article.html, Dict):
                    for slug, html in article.html.items():
                        if slug == "index":
                            slug = ""
                            output_html = article.output_html
                        elif "." in slug:
                            output_html = article.output_html.parent / slug
                        else:
                            slug = slugify(slug)
                            output_html = article.output_html.parent / slug / "index.html"

                        output_html.parent.mkdir(parents=True, exist_ok=True)
                        output_html.write_text(html)
```


!! method <h2 id='default_slug' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_slug <em class='small'>method</em></h2>

???+ source "default_slug <em class='small'>source</em>"

```python

        def default_slug(cls, v, *, values):
                if v is None:
                    return slugify(str(values["path"].stem))
                return v
```


!! method <h2 id='default_output_html' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_output_html <em class='small'>method</em></h2>

???+ source "default_output_html <em class='small'>source</em>"

```python

        def default_output_html(
                cls: "OutputHTML", v: Optional[Path], *, values: Dict
            ) -> Path:
                if isinstance(v, str):
                    v = Path(v)
                if v is not None:
                    return v
                if "slug" not in values:
                    for validator in cls.__validators__["slug"]:
                        values["slug"] = validator.func(cls, v, values=values)

                if values["slug"] == "index":
                    return cls.markata.config.output_dir / "index.html"
                return cls.markata.config.output_dir / values["slug"] / "index.html"
```


!! method <h2 id='output_html_relative' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>output_html_relative <em class='small'>method</em></h2>

???+ source "output_html_relative <em class='small'>source</em>"

```python

        def output_html_relative(
                cls: "OutputHTML", v: Optional[Path], *, values: Dict
            ) -> Path:
                if isinstance(v, str):
                    v = Path(v)
                if cls.markata.config.output_dir.absolute() not in v.absolute().parents:
                    return cls.markata.config.output_dir / v
                return v
```


!! method <h2 id='output_html_exists' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>output_html_exists <em class='small'>method</em></h2>

???+ source "output_html_exists <em class='small'>source</em>"

```python

        def output_html_exists(
                cls: "OutputHTML", v: Optional[Path], *, values: Dict
            ) -> Path:
                if isinstance(v, str):
                    v = Path(v)
                if not v.parent.exists():
                    v.parent.mkdir(parents=True, exist_ok=True)
                return v
```

