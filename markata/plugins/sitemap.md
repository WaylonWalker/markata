---
content: "None\n\n\n!! class <h2 id='SiteMapUrl' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>SiteMapUrl <em class='small'>class</em></h2>\n\n???+ source \"SiteMapUrl
  <em class='small'>source</em>\"\n\n```python\n\n        class SiteMapUrl(pydantic.BaseModel):\n
  \           slug: str = pydantic.Field(..., exclude=True)\n            loc: str
  = pydantic.Field(None, include=True)\n            changefreq: str = pydantic.Field(\"daily\",
  include=True)\n            priority: str = pydantic.Field(\"0.7\", include=True)\n\n
  \           @pydantic.validator(\"loc\")\n            def default_loc(cls, v, *,
  values):\n                if v is None:\n                    return cls.markata.config.url
  + \"/\" + values[\"slug\"] + \"/\"\n                return v\n\n            def
  dict(self, *args, **kwargs):\n                return {\"url\": {**super().dict(*args,
  **kwargs)}}\n```\n\n\n!! class <h2 id='SiteMapPost' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>SiteMapPost <em class='small'>class</em></h2>\n\n???+ source \"SiteMapPost
  <em class='small'>source</em>\"\n\n```python\n\n        class SiteMapPost(pydantic.BaseModel):\n
  \           slug: str = None\n            published: bool = True\n            sitemap_url:
  Optional[SiteMapUrl] = None\n\n            @pydantic.validator(\"sitemap_url\",
  pre=False, always=True)\n            def default_loc(cls, v, *, values):\n                if
  v is None:\n                    return SiteMapUrl(markata=cls.markata, slug=values[\"slug\"])\n
  \               if v.markata is None:\n                    return SiteMapUrl(markata=cls.markata,
  slug=values[\"slug\"])\n                return v\n```\n\n\n!! function <h2 id='post_model'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2>\n\n???+
  source \"post_model <em class='small'>source</em>\"\n\n```python\n\n        def
  post_model(markata: \"Markata\") -> None:\n            markata.post_models.append(SiteMapPost)\n```\n\n\n!!
  function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>render <em class='small'>function</em></h2>\n\n???+ source \"render <em class='small'>source</em>\"\n\n```python\n\n
  \       def render(markata: Markata) -> None:\n            sitemap = {\n                \"urlset\":
  markata.map(\"post.sitemap_url.dict()\", filter=\"post.published\")\n            }\n\n
  \           sitemap = (\n                anyconfig.dumps(sitemap, \"xml\")\n                .decode(\"utf-8\")\n
  \               .replace(\n                    \"<urlset>\",\n                    (\n
  \                       '<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\"\\n'\n
  \                       'xmlns:news=\"http://www.google.com/schemas/sitemap-news/0.9\"\\n'\n
  \                       'xmlns:xhtml=\"http://www.w3.org/1999/xhtml\"\\n'\n                        'xmlns:mobile=\"http://www.google.com/schemas/sitemap-mobile/1.0\"\\n'\n
  \                       'xmlns:image=\"http://www.google.com/schemas/sitemap-image/1.1\"\\n'\n
  \                       'xmlns:video=\"http://www.google.com/schemas/sitemap-video/1.1\">\\n'\n
  \                   ),\n                )\n                .replace(\"</url>\",
  \"</url>\\n\")\n            )\n            markata.sitemap = sitemap\n```\n\n\n!!
  function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save
  <em class='small'>function</em></h2>\n\n???+ source \"save <em class='small'>source</em>\"\n\n```python\n\n
  \       def save(markata: Markata) -> None:\n            with open(Path(markata.config.output_dir)
  / \"sitemap.xml\", \"w\") as f:\n                f.write(markata.sitemap)\n```\n\n\n!!
  method <h2 id='default_loc' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>default_loc <em class='small'>method</em></h2>\n\n???+ source \"default_loc
  <em class='small'>source</em>\"\n\n```python\n\n        def default_loc(cls, v,
  *, values):\n                if v is None:\n                    return cls.markata.config.url
  + \"/\" + values[\"slug\"] + \"/\"\n                return v\n```\n\n\n!! method
  <h2 id='dict' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dict
  <em class='small'>method</em></h2>\n\n???+ source \"dict <em class='small'>source</em>\"\n\n```python\n\n
  \       def dict(self, *args, **kwargs):\n                return {\"url\": {**super().dict(*args,
  **kwargs)}}\n```\n\n\n!! method <h2 id='default_loc' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>default_loc <em class='small'>method</em></h2>\n\n???+ source \"default_loc
  <em class='small'>source</em>\"\n\n```python\n\n        def default_loc(cls, v,
  *, values):\n                if v is None:\n                    return SiteMapUrl(markata=cls.markata,
  slug=values[\"slug\"])\n                if v.markata is None:\n                    return
  SiteMapUrl(markata=cls.markata, slug=values[\"slug\"])\n                return v\n```\n\n"
date: 0001-01-01
description: 'None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ???+ source  ! ???+ source  ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Sitemap.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \"
    />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\"
    href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script
    src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n<title>Sitemap.Py</title>\n<meta charset=\"UTF-8\" />\n<meta
    name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta name=\"description\"
    content=\"None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Sitemap.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='SiteMapUrl' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>SiteMapUrl <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">SiteMapUrl
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
    <span class=\"nc\">SiteMapUrl</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">slug</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span class=\"o\">.</span><span
    class=\"n\">Field</span><span class=\"p\">(</span><span class=\"o\">...</span><span
    class=\"p\">,</span> <span class=\"n\">exclude</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"n\">loc</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">Field</span><span
    class=\"p\">(</span><span class=\"kc\">None</span><span class=\"p\">,</span> <span
    class=\"n\">include</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"n\">changefreq</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;daily&quot;</span><span class=\"p\">,</span> <span class=\"n\">include</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">priority</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;0.7&quot;</span><span class=\"p\">,</span> <span class=\"n\">include</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;loc&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">default_loc</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">url</span> <span class=\"o\">+</span> <span
    class=\"s2\">&quot;/&quot;</span> <span class=\"o\">+</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">+</span> <span class=\"s2\">&quot;/&quot;</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">dict</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;url&quot;</span><span class=\"p\">:</span> <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"nb\">super</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">dict</span><span class=\"p\">(</span><span
    class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">,</span> <span
    class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">)}}</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='SiteMapPost' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>SiteMapPost <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">SiteMapPost
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
    <span class=\"nc\">SiteMapPost</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">slug</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n            <span class=\"n\">published</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n            <span class=\"n\">sitemap_url</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">SiteMapUrl</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;sitemap_url&quot;</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
    class=\"o\">=</span><span class=\"kc\">False</span><span class=\"p\">,</span>
    <span class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">default_loc</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">SiteMapUrl</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">slug</span><span class=\"o\">=</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">])</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span><span class=\"o\">.</span><span
    class=\"n\">markata</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">SiteMapUrl</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">slug</span><span
    class=\"o\">=</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">])</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class=\"n\">SiteMapPost</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>render <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render
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
    <span class=\"nf\">render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">sitemap</span> <span class=\"o\">=</span> <span
    class=\"p\">{</span>\n                <span class=\"s2\">&quot;urlset&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">map</span><span class=\"p\">(</span><span class=\"s2\">&quot;post.sitemap_url.dict()&quot;</span><span
    class=\"p\">,</span> <span class=\"nb\">filter</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;post.published&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"p\">}</span>\n\n            <span class=\"n\">sitemap</span> <span class=\"o\">=</span>
    <span class=\"p\">(</span>\n                <span class=\"n\">anyconfig</span><span
    class=\"o\">.</span><span class=\"n\">dumps</span><span class=\"p\">(</span><span
    class=\"n\">sitemap</span><span class=\"p\">,</span> <span class=\"s2\">&quot;xml&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"o\">.</span><span class=\"n\">decode</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;utf-8&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span>\n                    <span class=\"s2\">&quot;&lt;urlset&gt;&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">(</span>\n                        <span
    class=\"s1\">&#39;&lt;urlset xmlns=&quot;http://www.sitemaps.org/schemas/sitemap/0.9&quot;</span><span
    class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n                        <span
    class=\"s1\">&#39;xmlns:news=&quot;http://www.google.com/schemas/sitemap-news/0.9&quot;</span><span
    class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n                        <span
    class=\"s1\">&#39;xmlns:xhtml=&quot;http://www.w3.org/1999/xhtml&quot;</span><span
    class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n                        <span
    class=\"s1\">&#39;xmlns:mobile=&quot;http://www.google.com/schemas/sitemap-mobile/1.0&quot;</span><span
    class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n                        <span
    class=\"s1\">&#39;xmlns:image=&quot;http://www.google.com/schemas/sitemap-image/1.1&quot;</span><span
    class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n                        <span
    class=\"s1\">&#39;xmlns:video=&quot;http://www.google.com/schemas/sitemap-video/1.1&quot;&gt;</span><span
    class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n                    <span
    class=\"p\">),</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;&lt;/url&gt;&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;&lt;/url&gt;</span><span
    class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">sitemap</span> <span class=\"o\">=</span>
    <span class=\"n\">sitemap</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2
    id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
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
    class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">with</span> <span class=\"nb\">open</span><span
    class=\"p\">(</span><span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">output_dir</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;sitemap.xml&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;w&quot;</span><span class=\"p\">)</span>
    <span class=\"k\">as</span> <span class=\"n\">f</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">f</span><span class=\"o\">.</span><span class=\"n\">write</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">sitemap</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='default_loc' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>default_loc <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">default_loc
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
    <span class=\"nf\">default_loc</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">url</span>
    <span class=\"o\">+</span> <span class=\"s2\">&quot;/&quot;</span> <span class=\"o\">+</span>
    <span class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;/&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='dict' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dict
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">dict <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">dict</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;url&quot;</span><span class=\"p\">:</span> <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"nb\">super</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">dict</span><span class=\"p\">(</span><span
    class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">,</span> <span
    class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">)}}</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='default_loc' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>default_loc <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">default_loc
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
    <span class=\"nf\">default_loc</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">SiteMapUrl</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"bp\">cls</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">slug</span><span class=\"o\">=</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">])</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">markata</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">SiteMapUrl</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">slug</span><span class=\"o\">=</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">])</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Sitemap.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source \"
    />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\"
    href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script
    src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n<title>Sitemap.Py</title>\n<meta charset=\"UTF-8\" />\n<meta
    name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta name=\"description\"
    content=\"None ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Sitemap.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Sitemap.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    class <h2 id='SiteMapUrl' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>SiteMapUrl <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">SiteMapUrl
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
    <span class=\"nc\">SiteMapUrl</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">slug</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span class=\"o\">.</span><span
    class=\"n\">Field</span><span class=\"p\">(</span><span class=\"o\">...</span><span
    class=\"p\">,</span> <span class=\"n\">exclude</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n            <span class=\"n\">loc</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span> <span class=\"o\">=</span>
    <span class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">Field</span><span
    class=\"p\">(</span><span class=\"kc\">None</span><span class=\"p\">,</span> <span
    class=\"n\">include</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"n\">changefreq</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;daily&quot;</span><span class=\"p\">,</span> <span class=\"n\">include</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">priority</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">Field</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;0.7&quot;</span><span class=\"p\">,</span> <span class=\"n\">include</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;loc&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">default_loc</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">url</span> <span class=\"o\">+</span> <span
    class=\"s2\">&quot;/&quot;</span> <span class=\"o\">+</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">+</span> <span class=\"s2\">&quot;/&quot;</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">dict</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;url&quot;</span><span class=\"p\">:</span> <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"nb\">super</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">dict</span><span class=\"p\">(</span><span
    class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">,</span> <span
    class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">)}}</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='SiteMapPost' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>SiteMapPost <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">SiteMapPost
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
    <span class=\"nc\">SiteMapPost</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">slug</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"kc\">None</span>\n            <span class=\"n\">published</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n            <span class=\"n\">sitemap_url</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">SiteMapUrl</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span>\n\n            <span class=\"nd\">@pydantic</span><span
    class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;sitemap_url&quot;</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
    class=\"o\">=</span><span class=\"kc\">False</span><span class=\"p\">,</span>
    <span class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">default_loc</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">SiteMapUrl</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">slug</span><span class=\"o\">=</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">])</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span><span class=\"o\">.</span><span
    class=\"n\">markata</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">SiteMapUrl</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"bp\">cls</span><span class=\"o\">.</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">slug</span><span
    class=\"o\">=</span><span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;slug&quot;</span><span class=\"p\">])</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class=\"n\">SiteMapPost</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>render <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">render
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
    <span class=\"nf\">render</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">sitemap</span> <span class=\"o\">=</span> <span
    class=\"p\">{</span>\n                <span class=\"s2\">&quot;urlset&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">map</span><span class=\"p\">(</span><span class=\"s2\">&quot;post.sitemap_url.dict()&quot;</span><span
    class=\"p\">,</span> <span class=\"nb\">filter</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;post.published&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"p\">}</span>\n\n            <span class=\"n\">sitemap</span> <span class=\"o\">=</span>
    <span class=\"p\">(</span>\n                <span class=\"n\">anyconfig</span><span
    class=\"o\">.</span><span class=\"n\">dumps</span><span class=\"p\">(</span><span
    class=\"n\">sitemap</span><span class=\"p\">,</span> <span class=\"s2\">&quot;xml&quot;</span><span
    class=\"p\">)</span>\n                <span class=\"o\">.</span><span class=\"n\">decode</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;utf-8&quot;</span><span class=\"p\">)</span>\n
    \               <span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span>\n                    <span class=\"s2\">&quot;&lt;urlset&gt;&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">(</span>\n                        <span
    class=\"s1\">&#39;&lt;urlset xmlns=&quot;http://www.sitemaps.org/schemas/sitemap/0.9&quot;</span><span
    class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n                        <span
    class=\"s1\">&#39;xmlns:news=&quot;http://www.google.com/schemas/sitemap-news/0.9&quot;</span><span
    class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n                        <span
    class=\"s1\">&#39;xmlns:xhtml=&quot;http://www.w3.org/1999/xhtml&quot;</span><span
    class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n                        <span
    class=\"s1\">&#39;xmlns:mobile=&quot;http://www.google.com/schemas/sitemap-mobile/1.0&quot;</span><span
    class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n                        <span
    class=\"s1\">&#39;xmlns:image=&quot;http://www.google.com/schemas/sitemap-image/1.1&quot;</span><span
    class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n                        <span
    class=\"s1\">&#39;xmlns:video=&quot;http://www.google.com/schemas/sitemap-video/1.1&quot;&gt;</span><span
    class=\"se\">\\n</span><span class=\"s1\">&#39;</span>\n                    <span
    class=\"p\">),</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"o\">.</span><span class=\"n\">replace</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;&lt;/url&gt;&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;&lt;/url&gt;</span><span
    class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">sitemap</span> <span class=\"o\">=</span>
    <span class=\"n\">sitemap</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2
    id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
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
    class=\"p\">:</span> <span class=\"n\">Markata</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">with</span> <span class=\"nb\">open</span><span
    class=\"p\">(</span><span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">output_dir</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;sitemap.xml&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;w&quot;</span><span class=\"p\">)</span>
    <span class=\"k\">as</span> <span class=\"n\">f</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">f</span><span class=\"o\">.</span><span class=\"n\">write</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">sitemap</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='default_loc' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>default_loc <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">default_loc
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
    <span class=\"nf\">default_loc</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">url</span>
    <span class=\"o\">+</span> <span class=\"s2\">&quot;/&quot;</span> <span class=\"o\">+</span>
    <span class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">+</span> <span class=\"s2\">&quot;/&quot;</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='dict' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dict
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">dict <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">dict</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"n\">args</span><span
    class=\"p\">,</span> <span class=\"o\">**</span><span class=\"n\">kwargs</span><span
    class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"p\">{</span><span
    class=\"s2\">&quot;url&quot;</span><span class=\"p\">:</span> <span class=\"p\">{</span><span
    class=\"o\">**</span><span class=\"nb\">super</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">dict</span><span class=\"p\">(</span><span
    class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">,</span> <span
    class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">)}}</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='default_loc' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>default_loc <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">default_loc
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
    <span class=\"nf\">default_loc</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">SiteMapUrl</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"bp\">cls</span><span class=\"o\">.</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">slug</span><span class=\"o\">=</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">])</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">markata</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">SiteMapUrl</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">=</span><span class=\"bp\">cls</span><span
    class=\"o\">.</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">slug</span><span class=\"o\">=</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">])</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/sitemap
title: Sitemap.Py


---

None


!! class <h2 id='SiteMapUrl' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>SiteMapUrl <em class='small'>class</em></h2>

???+ source "SiteMapUrl <em class='small'>source</em>"

```python

        class SiteMapUrl(pydantic.BaseModel):
            slug: str = pydantic.Field(..., exclude=True)
            loc: str = pydantic.Field(None, include=True)
            changefreq: str = pydantic.Field("daily", include=True)
            priority: str = pydantic.Field("0.7", include=True)

            @pydantic.validator("loc")
            def default_loc(cls, v, *, values):
                if v is None:
                    return cls.markata.config.url + "/" + values["slug"] + "/"
                return v

            def dict(self, *args, **kwargs):
                return {"url": {**super().dict(*args, **kwargs)}}
```


!! class <h2 id='SiteMapPost' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>SiteMapPost <em class='small'>class</em></h2>

???+ source "SiteMapPost <em class='small'>source</em>"

```python

        class SiteMapPost(pydantic.BaseModel):
            slug: str = None
            published: bool = True
            sitemap_url: Optional[SiteMapUrl] = None

            @pydantic.validator("sitemap_url", pre=False, always=True)
            def default_loc(cls, v, *, values):
                if v is None:
                    return SiteMapUrl(markata=cls.markata, slug=values["slug"])
                if v.markata is None:
                    return SiteMapUrl(markata=cls.markata, slug=values["slug"])
                return v
```


!! function <h2 id='post_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post_model <em class='small'>function</em></h2>

???+ source "post_model <em class='small'>source</em>"

```python

        def post_model(markata: "Markata") -> None:
            markata.post_models.append(SiteMapPost)
```


!! function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2>

???+ source "render <em class='small'>source</em>"

```python

        def render(markata: Markata) -> None:
            sitemap = {
                "urlset": markata.map("post.sitemap_url.dict()", filter="post.published")
            }

            sitemap = (
                anyconfig.dumps(sitemap, "xml")
                .decode("utf-8")
                .replace(
                    "<urlset>",
                    (
                        '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"\n'
                        'xmlns:news="http://www.google.com/schemas/sitemap-news/0.9"\n'
                        'xmlns:xhtml="http://www.w3.org/1999/xhtml"\n'
                        'xmlns:mobile="http://www.google.com/schemas/sitemap-mobile/1.0"\n'
                        'xmlns:image="http://www.google.com/schemas/sitemap-image/1.1"\n'
                        'xmlns:video="http://www.google.com/schemas/sitemap-video/1.1">\n'
                    ),
                )
                .replace("</url>", "</url>\n")
            )
            markata.sitemap = sitemap
```


!! function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>

???+ source "save <em class='small'>source</em>"

```python

        def save(markata: Markata) -> None:
            with open(Path(markata.config.output_dir) / "sitemap.xml", "w") as f:
                f.write(markata.sitemap)
```


!! method <h2 id='default_loc' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_loc <em class='small'>method</em></h2>

???+ source "default_loc <em class='small'>source</em>"

```python

        def default_loc(cls, v, *, values):
                if v is None:
                    return cls.markata.config.url + "/" + values["slug"] + "/"
                return v
```


!! method <h2 id='dict' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>dict <em class='small'>method</em></h2>

???+ source "dict <em class='small'>source</em>"

```python

        def dict(self, *args, **kwargs):
                return {"url": {**super().dict(*args, **kwargs)}}
```


!! method <h2 id='default_loc' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_loc <em class='small'>method</em></h2>

???+ source "default_loc <em class='small'>source</em>"

```python

        def default_loc(cls, v, *, values):
                if v is None:
                    return SiteMapUrl(markata=cls.markata, slug=values["slug"])
                if v.markata is None:
                    return SiteMapUrl(markata=cls.markata, slug=values["slug"])
                return v
```

