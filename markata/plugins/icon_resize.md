---
content: "Icon Resize Plugin\n\nResized favicon to a set of common sizes.\n\n## markata.plugins.icon_resize
  configuration\n\n```toml title=markata.toml\n[markata]\noutput_dir = \"markout\"\nassets_dir
  = \"static\"\nicon = \"static/icon.png\"\n```\n\n\n!! class <h2 id='Config' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>\n\n???+
  source \"Config <em class='small'>source</em>\"\n\n```python\n\n        class Config(pydantic.BaseModel):\n
  \           output_dir: pydantic.DirectoryPath = \"markout\"\n            assets_dir:
  Path = pydantic.Field(\n                Path(\"static\"),\n                description=\"The
  directory to store static assets\",\n            )\n            icon: Optional[Path]
  = None\n            icon_out_file: Optional[Path] = None\n            icons: Optional[List[Dict[str,
  str]]] = []\n\n            @pydantic.validator(\"icon\")\n            def ensure_icon_exists(cls,
  v, *, values: Dict) -> Path:\n                if v is None:\n                    return\n
  \               if v.exists():\n                    return v\n\n                icon
  = Path(values[\"assets_dir\"]) / v\n\n                if icon.exists():\n                    return
  icon\n                else:\n                    raise FileNotFoundError(v)\n\n
  \           @pydantic.validator(\"icon_out_file\", pre=True, always=True)\n            def
  default_icon_out_file(cls, v, *, values: Dict) -> Path:\n                if v is
  None and values[\"icon\"] is not None:\n                    return Path(values[\"output_dir\"])
  / values[\"icon\"]\n                return v\n```\n\n\n!! function <h2 id='config_model'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>\n\n???+
  source \"config_model <em class='small'>source</em>\"\n\n```python\n\n        def
  config_model(markata: \"Markata\") -> None:\n            markata.config_models.append(Config)\n```\n\n\n!!
  function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>render <em class='small'>function</em></h2>\n\n???+ source \"render <em class='small'>source</em>\"\n\n```python\n\n
  \       def render(markata: \"Markata\") -> None:\n            if markata.config.icon
  is None:\n                return\n\n            with Image.open(markata.config.icon)
  as img:\n                for width in [48, 72, 96, 144, 192, 256, 384, 512]:\n                    height
  = int(float(img.size[1]) * float(width / float(img.size[0])))\n                    filename
  = Path(\n                        f\"{markata.config.icon_out_file.stem}_{width}x{height}{markata.config.icon_out_file.suffix}\",\n
  \                   )\n                    markata.config.icons.append(\n                        {\n
  \                           \"src\": str(filename),\n                            \"sizes\":
  f\"{width}x{width}\",\n                            \"type\": f\"image/{img.format}\".lower(),\n
  \                           \"purpose\": \"any maskable\",\n                        },\n
  \                   )\n```\n\n\n!! function <h2 id='save' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>\n\n???+
  source \"save <em class='small'>source</em>\"\n\n```python\n\n        def save(markata:
  \"Markata\") -> None:\n            if markata.config.icon is None:\n                return\n
  \           for width in [48, 72, 96, 144, 192, 256, 384, 512]:\n                with
  Image.open(markata.config.icon) as img:\n                    height = int(float(img.size[1])
  * float(width / float(img.size[0])))\n                    img = img.resize((width,
  height), Image.LANCZOS)\n                    filename = Path(\n                        f\"{markata.config.icon_out_file.stem}_{width}x{height}{markata.config.icon_out_file.suffix}\",\n
  \                   )\n                    out_file = Path(markata.config.output_dir)
  / filename\n                    if out_file.exists():\n                        continue\n
  \                   img = img.resize((width, height), Image.LANCZOS)\n                    img.save(out_file)\n```\n\n\n!!
  method <h2 id='ensure_icon_exists' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>ensure_icon_exists <em class='small'>method</em></h2>\n\n???+ source \"ensure_icon_exists
  <em class='small'>source</em>\"\n\n```python\n\n        def ensure_icon_exists(cls,
  v, *, values: Dict) -> Path:\n                if v is None:\n                    return\n
  \               if v.exists():\n                    return v\n\n                icon
  = Path(values[\"assets_dir\"]) / v\n\n                if icon.exists():\n                    return
  icon\n                else:\n                    raise FileNotFoundError(v)\n```\n\n\n!!
  method <h2 id='default_icon_out_file' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>default_icon_out_file <em class='small'>method</em></h2>\n\n???+ source \"default_icon_out_file
  <em class='small'>source</em>\"\n\n```python\n\n        def default_icon_out_file(cls,
  v, *, values: Dict) -> Path:\n                if v is None and values[\"icon\"]
  is not None:\n                    return Path(values[\"output_dir\"]) / values[\"icon\"]\n
  \               return v\n```\n\n"
date: 0001-01-01
description: 'Icon Resize Plugin Resized favicon to a set of common sizes. ! ???+
  source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Icon_Resize.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Icon Resize Plugin Resized favicon to
    a set of common sizes. ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Icon_Resize.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Icon Resize Plugin Resized favicon to
    a set of common sizes. ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Icon_Resize.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>Icon Resize Plugin</p>\n<p>Resized
    favicon to a set of common sizes.</p>\n<h2 id=\"markatapluginsicon_resize-configuration\">markata.plugins.icon_resize
    configuration <a class=\"header-anchor\" href=\"#markatapluginsicon_resize-configuration\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre class='wrapper'>\n\n<div
    class='filepath'>\n<p>markata.toml</p>\n\n<div class='right'>\n\n<button class='copy'
    title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"n\">output_dir</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;markout&quot;</span>\n<span class=\"n\">assets_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;static&quot;</span>\n<span class=\"n\">icon</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;static/icon.png&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    \           <span class=\"n\">output_dir</span><span class=\"p\">:</span> <span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">DirectoryPath</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;markout&quot;</span>\n            <span
    class=\"n\">assets_dir</span><span class=\"p\">:</span> <span class=\"n\">Path</span>
    <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span class=\"o\">.</span><span
    class=\"n\">Field</span><span class=\"p\">(</span>\n                <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;static&quot;</span><span class=\"p\">),</span>\n
    \               <span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;The directory to store static assets&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"n\">icon</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">Path</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">icon_out_file</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">Path</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">icons</span><span class=\"p\">:</span>
    <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"n\">Dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">]]]</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;icon&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">ensure_icon_exists</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Path</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n
    \               <span class=\"n\">icon</span> <span class=\"o\">=</span> <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;assets_dir&quot;</span><span class=\"p\">])</span>
    <span class=\"o\">/</span> <span class=\"n\">v</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">icon</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">icon</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"k\">raise</span>
    <span class=\"ne\">FileNotFoundError</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">)</span>\n\n            <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;icon_out_file&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">always</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">default_icon_out_file</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Path</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;icon&quot;</span><span class=\"p\">]</span>
    <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">])</span>
    <span class=\"o\">/</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;icon&quot;</span><span class=\"p\">]</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>config_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">config_model
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
    <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">icon</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span>\n\n            <span
    class=\"k\">with</span> <span class=\"n\">Image</span><span class=\"o\">.</span><span
    class=\"n\">open</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">icon</span><span class=\"p\">)</span> <span class=\"k\">as</span>
    <span class=\"n\">img</span><span class=\"p\">:</span>\n                <span
    class=\"k\">for</span> <span class=\"n\">width</span> <span class=\"ow\">in</span>
    <span class=\"p\">[</span><span class=\"mi\">48</span><span class=\"p\">,</span>
    <span class=\"mi\">72</span><span class=\"p\">,</span> <span class=\"mi\">96</span><span
    class=\"p\">,</span> <span class=\"mi\">144</span><span class=\"p\">,</span> <span
    class=\"mi\">192</span><span class=\"p\">,</span> <span class=\"mi\">256</span><span
    class=\"p\">,</span> <span class=\"mi\">384</span><span class=\"p\">,</span> <span
    class=\"mi\">512</span><span class=\"p\">]:</span>\n                    <span
    class=\"n\">height</span> <span class=\"o\">=</span> <span class=\"nb\">int</span><span
    class=\"p\">(</span><span class=\"nb\">float</span><span class=\"p\">(</span><span
    class=\"n\">img</span><span class=\"o\">.</span><span class=\"n\">size</span><span
    class=\"p\">[</span><span class=\"mi\">1</span><span class=\"p\">])</span> <span
    class=\"o\">*</span> <span class=\"nb\">float</span><span class=\"p\">(</span><span
    class=\"n\">width</span> <span class=\"o\">/</span> <span class=\"nb\">float</span><span
    class=\"p\">(</span><span class=\"n\">img</span><span class=\"o\">.</span><span
    class=\"n\">size</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">])))</span>\n                    <span class=\"n\">filename</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span>\n
    \                       <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">icon_out_file</span><span
    class=\"o\">.</span><span class=\"n\">stem</span><span class=\"si\">}</span><span
    class=\"s2\">_</span><span class=\"si\">{</span><span class=\"n\">width</span><span
    class=\"si\">}</span><span class=\"s2\">x</span><span class=\"si\">{</span><span
    class=\"n\">height</span><span class=\"si\">}{</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">icon_out_file</span><span class=\"o\">.</span><span class=\"n\">suffix</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">icons</span><span class=\"o\">.</span><span class=\"n\">append</span><span
    class=\"p\">(</span>\n                        <span class=\"p\">{</span>\n                            <span
    class=\"s2\">&quot;src&quot;</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">filename</span><span class=\"p\">),</span>\n
    \                           <span class=\"s2\">&quot;sizes&quot;</span><span class=\"p\">:</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">width</span><span class=\"si\">}</span><span class=\"s2\">x</span><span
    class=\"si\">{</span><span class=\"n\">width</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">,</span>\n                            <span
    class=\"s2\">&quot;type&quot;</span><span class=\"p\">:</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;image/</span><span class=\"si\">{</span><span class=\"n\">img</span><span
    class=\"o\">.</span><span class=\"n\">format</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">(),</span>\n                            <span class=\"s2\">&quot;purpose&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;any maskable&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"p\">},</span>\n                    <span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='save' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">save
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
    <span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">icon</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span>\n            <span
    class=\"k\">for</span> <span class=\"n\">width</span> <span class=\"ow\">in</span>
    <span class=\"p\">[</span><span class=\"mi\">48</span><span class=\"p\">,</span>
    <span class=\"mi\">72</span><span class=\"p\">,</span> <span class=\"mi\">96</span><span
    class=\"p\">,</span> <span class=\"mi\">144</span><span class=\"p\">,</span> <span
    class=\"mi\">192</span><span class=\"p\">,</span> <span class=\"mi\">256</span><span
    class=\"p\">,</span> <span class=\"mi\">384</span><span class=\"p\">,</span> <span
    class=\"mi\">512</span><span class=\"p\">]:</span>\n                <span class=\"k\">with</span>
    <span class=\"n\">Image</span><span class=\"o\">.</span><span class=\"n\">open</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">icon</span><span
    class=\"p\">)</span> <span class=\"k\">as</span> <span class=\"n\">img</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">height</span> <span
    class=\"o\">=</span> <span class=\"nb\">int</span><span class=\"p\">(</span><span
    class=\"nb\">float</span><span class=\"p\">(</span><span class=\"n\">img</span><span
    class=\"o\">.</span><span class=\"n\">size</span><span class=\"p\">[</span><span
    class=\"mi\">1</span><span class=\"p\">])</span> <span class=\"o\">*</span> <span
    class=\"nb\">float</span><span class=\"p\">(</span><span class=\"n\">width</span>
    <span class=\"o\">/</span> <span class=\"nb\">float</span><span class=\"p\">(</span><span
    class=\"n\">img</span><span class=\"o\">.</span><span class=\"n\">size</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">])))</span>\n
    \                   <span class=\"n\">img</span> <span class=\"o\">=</span> <span
    class=\"n\">img</span><span class=\"o\">.</span><span class=\"n\">resize</span><span
    class=\"p\">((</span><span class=\"n\">width</span><span class=\"p\">,</span>
    <span class=\"n\">height</span><span class=\"p\">),</span> <span class=\"n\">Image</span><span
    class=\"o\">.</span><span class=\"n\">LANCZOS</span><span class=\"p\">)</span>\n
    \                   <span class=\"n\">filename</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">icon_out_file</span><span class=\"o\">.</span><span
    class=\"n\">stem</span><span class=\"si\">}</span><span class=\"s2\">_</span><span
    class=\"si\">{</span><span class=\"n\">width</span><span class=\"si\">}</span><span
    class=\"s2\">x</span><span class=\"si\">{</span><span class=\"n\">height</span><span
    class=\"si\">}{</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">icon_out_file</span><span
    class=\"o\">.</span><span class=\"n\">suffix</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n                    <span class=\"n\">out_file</span> <span
    class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">output_dir</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"n\">filename</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">out_file</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n                        <span
    class=\"k\">continue</span>\n                    <span class=\"n\">img</span>
    <span class=\"o\">=</span> <span class=\"n\">img</span><span class=\"o\">.</span><span
    class=\"n\">resize</span><span class=\"p\">((</span><span class=\"n\">width</span><span
    class=\"p\">,</span> <span class=\"n\">height</span><span class=\"p\">),</span>
    <span class=\"n\">Image</span><span class=\"o\">.</span><span class=\"n\">LANCZOS</span><span
    class=\"p\">)</span>\n                    <span class=\"n\">img</span><span class=\"o\">.</span><span
    class=\"n\">save</span><span class=\"p\">(</span><span class=\"n\">out_file</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='ensure_icon_exists'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>ensure_icon_exists
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">ensure_icon_exists <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">ensure_icon_exists</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">:</span> <span class=\"n\">Dict</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n
    \               <span class=\"n\">icon</span> <span class=\"o\">=</span> <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;assets_dir&quot;</span><span class=\"p\">])</span>
    <span class=\"o\">/</span> <span class=\"n\">v</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">icon</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">icon</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"k\">raise</span>
    <span class=\"ne\">FileNotFoundError</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='default_icon_out_file'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_icon_out_file
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">default_icon_out_file <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">default_icon_out_file</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">:</span> <span class=\"n\">Dict</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">and</span>
    <span class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;icon&quot;</span><span
    class=\"p\">]</span> <span class=\"ow\">is</span> <span class=\"ow\">not</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;output_dir&quot;</span><span
    class=\"p\">])</span> <span class=\"o\">/</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;icon&quot;</span><span class=\"p\">]</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Icon_Resize.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Icon Resize Plugin Resized favicon to a set of
    common sizes. ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Icon_Resize.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Icon Resize Plugin Resized favicon to
    a set of common sizes. ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Icon_Resize.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Icon_Resize.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>Icon
    Resize Plugin</p>\n<p>Resized favicon to a set of common sizes.</p>\n<h2 id=\"markatapluginsicon_resize-configuration\">markata.plugins.icon_resize
    configuration <a class=\"header-anchor\" href=\"#markatapluginsicon_resize-configuration\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre class='wrapper'>\n\n<div
    class='filepath'>\n<p>markata.toml</p>\n\n<div class='right'>\n\n<button class='copy'
    title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n</div>\n\n<div
    class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
    class=\"n\">output_dir</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"s2\">&quot;markout&quot;</span>\n<span class=\"n\">assets_dir</span><span
    class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span
    class=\"s2\">&quot;static&quot;</span>\n<span class=\"n\">icon</span><span class=\"w\">
    </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">&quot;static/icon.png&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    \           <span class=\"n\">output_dir</span><span class=\"p\">:</span> <span
    class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">DirectoryPath</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;markout&quot;</span>\n            <span
    class=\"n\">assets_dir</span><span class=\"p\">:</span> <span class=\"n\">Path</span>
    <span class=\"o\">=</span> <span class=\"n\">pydantic</span><span class=\"o\">.</span><span
    class=\"n\">Field</span><span class=\"p\">(</span>\n                <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;static&quot;</span><span class=\"p\">),</span>\n
    \               <span class=\"n\">description</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;The directory to store static assets&quot;</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"n\">icon</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">Path</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">icon_out_file</span><span
    class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
    class=\"n\">Path</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
    class=\"kc\">None</span>\n            <span class=\"n\">icons</span><span class=\"p\">:</span>
    <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">List</span><span
    class=\"p\">[</span><span class=\"n\">Dict</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"nb\">str</span><span
    class=\"p\">]]]</span> <span class=\"o\">=</span> <span class=\"p\">[]</span>\n\n
    \           <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;icon&quot;</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">ensure_icon_exists</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Path</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n
    \               <span class=\"n\">icon</span> <span class=\"o\">=</span> <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;assets_dir&quot;</span><span class=\"p\">])</span>
    <span class=\"o\">/</span> <span class=\"n\">v</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">icon</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">icon</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"k\">raise</span>
    <span class=\"ne\">FileNotFoundError</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">)</span>\n\n            <span class=\"nd\">@pydantic</span><span class=\"o\">.</span><span
    class=\"n\">validator</span><span class=\"p\">(</span><span class=\"s2\">&quot;icon_out_file&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">pre</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span> <span class=\"n\">always</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">default_icon_out_file</span><span
    class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
    class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
    class=\"p\">,</span> <span class=\"n\">values</span><span class=\"p\">:</span>
    <span class=\"n\">Dict</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Path</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span> <span class=\"ow\">and</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;icon&quot;</span><span class=\"p\">]</span>
    <span class=\"ow\">is</span> <span class=\"ow\">not</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">])</span>
    <span class=\"o\">/</span> <span class=\"n\">values</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;icon&quot;</span><span class=\"p\">]</span>\n                <span
    class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>config_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">config_model
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
    <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">icon</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span>\n\n            <span
    class=\"k\">with</span> <span class=\"n\">Image</span><span class=\"o\">.</span><span
    class=\"n\">open</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">icon</span><span class=\"p\">)</span> <span class=\"k\">as</span>
    <span class=\"n\">img</span><span class=\"p\">:</span>\n                <span
    class=\"k\">for</span> <span class=\"n\">width</span> <span class=\"ow\">in</span>
    <span class=\"p\">[</span><span class=\"mi\">48</span><span class=\"p\">,</span>
    <span class=\"mi\">72</span><span class=\"p\">,</span> <span class=\"mi\">96</span><span
    class=\"p\">,</span> <span class=\"mi\">144</span><span class=\"p\">,</span> <span
    class=\"mi\">192</span><span class=\"p\">,</span> <span class=\"mi\">256</span><span
    class=\"p\">,</span> <span class=\"mi\">384</span><span class=\"p\">,</span> <span
    class=\"mi\">512</span><span class=\"p\">]:</span>\n                    <span
    class=\"n\">height</span> <span class=\"o\">=</span> <span class=\"nb\">int</span><span
    class=\"p\">(</span><span class=\"nb\">float</span><span class=\"p\">(</span><span
    class=\"n\">img</span><span class=\"o\">.</span><span class=\"n\">size</span><span
    class=\"p\">[</span><span class=\"mi\">1</span><span class=\"p\">])</span> <span
    class=\"o\">*</span> <span class=\"nb\">float</span><span class=\"p\">(</span><span
    class=\"n\">width</span> <span class=\"o\">/</span> <span class=\"nb\">float</span><span
    class=\"p\">(</span><span class=\"n\">img</span><span class=\"o\">.</span><span
    class=\"n\">size</span><span class=\"p\">[</span><span class=\"mi\">0</span><span
    class=\"p\">])))</span>\n                    <span class=\"n\">filename</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span>\n
    \                       <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span
    class=\"si\">{</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">icon_out_file</span><span
    class=\"o\">.</span><span class=\"n\">stem</span><span class=\"si\">}</span><span
    class=\"s2\">_</span><span class=\"si\">{</span><span class=\"n\">width</span><span
    class=\"si\">}</span><span class=\"s2\">x</span><span class=\"si\">{</span><span
    class=\"n\">height</span><span class=\"si\">}{</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">icon_out_file</span><span class=\"o\">.</span><span class=\"n\">suffix</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">icons</span><span class=\"o\">.</span><span class=\"n\">append</span><span
    class=\"p\">(</span>\n                        <span class=\"p\">{</span>\n                            <span
    class=\"s2\">&quot;src&quot;</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">filename</span><span class=\"p\">),</span>\n
    \                           <span class=\"s2\">&quot;sizes&quot;</span><span class=\"p\">:</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">width</span><span class=\"si\">}</span><span class=\"s2\">x</span><span
    class=\"si\">{</span><span class=\"n\">width</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">,</span>\n                            <span
    class=\"s2\">&quot;type&quot;</span><span class=\"p\">:</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;image/</span><span class=\"si\">{</span><span class=\"n\">img</span><span
    class=\"o\">.</span><span class=\"n\">format</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
    class=\"p\">(),</span>\n                            <span class=\"s2\">&quot;purpose&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;any maskable&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"p\">},</span>\n                    <span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='save' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">save
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
    <span class=\"nf\">save</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">icon</span> <span class=\"ow\">is</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span>\n            <span
    class=\"k\">for</span> <span class=\"n\">width</span> <span class=\"ow\">in</span>
    <span class=\"p\">[</span><span class=\"mi\">48</span><span class=\"p\">,</span>
    <span class=\"mi\">72</span><span class=\"p\">,</span> <span class=\"mi\">96</span><span
    class=\"p\">,</span> <span class=\"mi\">144</span><span class=\"p\">,</span> <span
    class=\"mi\">192</span><span class=\"p\">,</span> <span class=\"mi\">256</span><span
    class=\"p\">,</span> <span class=\"mi\">384</span><span class=\"p\">,</span> <span
    class=\"mi\">512</span><span class=\"p\">]:</span>\n                <span class=\"k\">with</span>
    <span class=\"n\">Image</span><span class=\"o\">.</span><span class=\"n\">open</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">icon</span><span
    class=\"p\">)</span> <span class=\"k\">as</span> <span class=\"n\">img</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">height</span> <span
    class=\"o\">=</span> <span class=\"nb\">int</span><span class=\"p\">(</span><span
    class=\"nb\">float</span><span class=\"p\">(</span><span class=\"n\">img</span><span
    class=\"o\">.</span><span class=\"n\">size</span><span class=\"p\">[</span><span
    class=\"mi\">1</span><span class=\"p\">])</span> <span class=\"o\">*</span> <span
    class=\"nb\">float</span><span class=\"p\">(</span><span class=\"n\">width</span>
    <span class=\"o\">/</span> <span class=\"nb\">float</span><span class=\"p\">(</span><span
    class=\"n\">img</span><span class=\"o\">.</span><span class=\"n\">size</span><span
    class=\"p\">[</span><span class=\"mi\">0</span><span class=\"p\">])))</span>\n
    \                   <span class=\"n\">img</span> <span class=\"o\">=</span> <span
    class=\"n\">img</span><span class=\"o\">.</span><span class=\"n\">resize</span><span
    class=\"p\">((</span><span class=\"n\">width</span><span class=\"p\">,</span>
    <span class=\"n\">height</span><span class=\"p\">),</span> <span class=\"n\">Image</span><span
    class=\"o\">.</span><span class=\"n\">LANCZOS</span><span class=\"p\">)</span>\n
    \                   <span class=\"n\">filename</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span>\n                        <span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">icon_out_file</span><span class=\"o\">.</span><span
    class=\"n\">stem</span><span class=\"si\">}</span><span class=\"s2\">_</span><span
    class=\"si\">{</span><span class=\"n\">width</span><span class=\"si\">}</span><span
    class=\"s2\">x</span><span class=\"si\">{</span><span class=\"n\">height</span><span
    class=\"si\">}{</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">icon_out_file</span><span
    class=\"o\">.</span><span class=\"n\">suffix</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n                    <span class=\"n\">out_file</span> <span
    class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">output_dir</span><span class=\"p\">)</span>
    <span class=\"o\">/</span> <span class=\"n\">filename</span>\n                    <span
    class=\"k\">if</span> <span class=\"n\">out_file</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n                        <span
    class=\"k\">continue</span>\n                    <span class=\"n\">img</span>
    <span class=\"o\">=</span> <span class=\"n\">img</span><span class=\"o\">.</span><span
    class=\"n\">resize</span><span class=\"p\">((</span><span class=\"n\">width</span><span
    class=\"p\">,</span> <span class=\"n\">height</span><span class=\"p\">),</span>
    <span class=\"n\">Image</span><span class=\"o\">.</span><span class=\"n\">LANCZOS</span><span
    class=\"p\">)</span>\n                    <span class=\"n\">img</span><span class=\"o\">.</span><span
    class=\"n\">save</span><span class=\"p\">(</span><span class=\"n\">out_file</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='ensure_icon_exists'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>ensure_icon_exists
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">ensure_icon_exists <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">ensure_icon_exists</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
    class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
    class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
    class=\"p\">:</span> <span class=\"n\">Dict</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">v</span>\n\n
    \               <span class=\"n\">icon</span> <span class=\"o\">=</span> <span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;assets_dir&quot;</span><span class=\"p\">])</span>
    <span class=\"o\">/</span> <span class=\"n\">v</span>\n\n                <span
    class=\"k\">if</span> <span class=\"n\">icon</span><span class=\"o\">.</span><span
    class=\"n\">exists</span><span class=\"p\">():</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">icon</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"k\">raise</span>
    <span class=\"ne\">FileNotFoundError</span><span class=\"p\">(</span><span class=\"n\">v</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='default_icon_out_file'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_icon_out_file
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">default_icon_out_file <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">default_icon_out_file</span><span class=\"p\">(</span><span
    class=\"bp\">cls</span><span class=\"p\">,</span> <span class=\"n\">v</span><span
    class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span> <span
    class=\"n\">values</span><span class=\"p\">:</span> <span class=\"n\">Dict</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Path</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span> <span class=\"ow\">and</span>
    <span class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;icon&quot;</span><span
    class=\"p\">]</span> <span class=\"ow\">is</span> <span class=\"ow\">not</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">&quot;output_dir&quot;</span><span
    class=\"p\">])</span> <span class=\"o\">/</span> <span class=\"n\">values</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;icon&quot;</span><span class=\"p\">]</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">v</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/icon-resize
title: Icon_Resize.Py


---

Icon Resize Plugin

Resized favicon to a set of common sizes.

## markata.plugins.icon_resize configuration

```toml title=markata.toml
[markata]
output_dir = "markout"
assets_dir = "static"
icon = "static/icon.png"
```


!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>

???+ source "Config <em class='small'>source</em>"

```python

        class Config(pydantic.BaseModel):
            output_dir: pydantic.DirectoryPath = "markout"
            assets_dir: Path = pydantic.Field(
                Path("static"),
                description="The directory to store static assets",
            )
            icon: Optional[Path] = None
            icon_out_file: Optional[Path] = None
            icons: Optional[List[Dict[str, str]]] = []

            @pydantic.validator("icon")
            def ensure_icon_exists(cls, v, *, values: Dict) -> Path:
                if v is None:
                    return
                if v.exists():
                    return v

                icon = Path(values["assets_dir"]) / v

                if icon.exists():
                    return icon
                else:
                    raise FileNotFoundError(v)

            @pydantic.validator("icon_out_file", pre=True, always=True)
            def default_icon_out_file(cls, v, *, values: Dict) -> Path:
                if v is None and values["icon"] is not None:
                    return Path(values["output_dir"]) / values["icon"]
                return v
```


!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>

???+ source "config_model <em class='small'>source</em>"

```python

        def config_model(markata: "Markata") -> None:
            markata.config_models.append(Config)
```


!! function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2>

???+ source "render <em class='small'>source</em>"

```python

        def render(markata: "Markata") -> None:
            if markata.config.icon is None:
                return

            with Image.open(markata.config.icon) as img:
                for width in [48, 72, 96, 144, 192, 256, 384, 512]:
                    height = int(float(img.size[1]) * float(width / float(img.size[0])))
                    filename = Path(
                        f"{markata.config.icon_out_file.stem}_{width}x{height}{markata.config.icon_out_file.suffix}",
                    )
                    markata.config.icons.append(
                        {
                            "src": str(filename),
                            "sizes": f"{width}x{width}",
                            "type": f"image/{img.format}".lower(),
                            "purpose": "any maskable",
                        },
                    )
```


!! function <h2 id='save' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>save <em class='small'>function</em></h2>

???+ source "save <em class='small'>source</em>"

```python

        def save(markata: "Markata") -> None:
            if markata.config.icon is None:
                return
            for width in [48, 72, 96, 144, 192, 256, 384, 512]:
                with Image.open(markata.config.icon) as img:
                    height = int(float(img.size[1]) * float(width / float(img.size[0])))
                    img = img.resize((width, height), Image.LANCZOS)
                    filename = Path(
                        f"{markata.config.icon_out_file.stem}_{width}x{height}{markata.config.icon_out_file.suffix}",
                    )
                    out_file = Path(markata.config.output_dir) / filename
                    if out_file.exists():
                        continue
                    img = img.resize((width, height), Image.LANCZOS)
                    img.save(out_file)
```


!! method <h2 id='ensure_icon_exists' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>ensure_icon_exists <em class='small'>method</em></h2>

???+ source "ensure_icon_exists <em class='small'>source</em>"

```python

        def ensure_icon_exists(cls, v, *, values: Dict) -> Path:
                if v is None:
                    return
                if v.exists():
                    return v

                icon = Path(values["assets_dir"]) / v

                if icon.exists():
                    return icon
                else:
                    raise FileNotFoundError(v)
```


!! method <h2 id='default_icon_out_file' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>default_icon_out_file <em class='small'>method</em></h2>

???+ source "default_icon_out_file <em class='small'>source</em>"

```python

        def default_icon_out_file(cls, v, *, values: Dict) -> Path:
                if v is None and values["icon"] is not None:
                    return Path(values["output_dir"]) / values["icon"]
                return v
```

