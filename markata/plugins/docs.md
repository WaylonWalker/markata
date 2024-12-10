---
content: "leading docstring\n\n\n!! function <h2 id='add_parents' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>add_parents <em class='small'>function</em></h2>\n\n???+
  source \"add_parents <em class='small'>source</em>\"\n\n```python\n\n        def
  add_parents(tree: ast.AST) -> None:\n            for node in ast.walk(tree):\n                for
  child in ast.iter_child_nodes(node):\n                    child.parent = node\n
  \                   if not hasattr(child, \"parents\"):\n                        child.parents
  = [node]\n                    child.parents.append(node)\n                    if
  isinstance(node, ast.ClassDef) and isinstance(child, ast.FunctionDef):\n                        child.type
  = \"method\"\n                    elif isinstance(child, ast.FunctionDef):\n                        child.type
  = \"function\"\n                    elif isinstance(child, ast.ClassDef):\n                        child.type
  = \"class\"\n```\n\n\n!! function <h2 id='glob' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>glob <em class='small'>function</em></h2>\n    finds k\n\n    ## Parameters\n\n
  \   `markata` the markata object\n???+ source \"glob <em class='small'>source</em>\"\n\n```python\n\n
  \       def glob(markata: \"MarkataDocs\") -> None:\n            \"\"\"\n            finds
  k\n\n            ## Parameters\n\n            `markata` the markata object\n\n            \"\"\"\n\n
  \           import glob\n\n            markata.py_files = [Path(f) for f in glob.glob(\"**/*.py\",
  recursive=True)]\n\n            content_directories = list({f.parent for f in markata.py_files})\n
  \           if \"content_directories\" in markata.__dict__:\n                markata.content_directories.extend(content_directories)\n
  \           else:\n                markata.content_directories = content_directories\n\n
  \           try:\n                ignore = True\n            except KeyError:\n
  \               ignore = True\n\n            if ignore and (Path(\".gitignore\").exists()
  or Path(\".markataignore\").exists()):\n                import pathspec\n\n                lines
  = []\n\n                if Path(\".gitignore\").exists():\n                    lines.extend(Path(\".gitignore\").read_text().splitlines())\n\n
  \               if Path(\".markataignore\").exists():\n                    lines.extend(Path(\".markataignore\").read_text().splitlines())\n\n
  \           spec = pathspec.PathSpec.from_lines(\"gitwildmatch\", lines)\n\n            markata.py_files
  = [\n                file for file in markata.py_files if not spec.match_file(str(file))\n
  \           ]\n```\n\n\n!! function <h2 id='get_template' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>get_template <em class='small'>function</em></h2>\n\n???+
  source \"get_template <em class='small'>source</em>\"\n\n```python\n\n        def
  get_template():\n            jinja_env = jinja2.Environment()\n            template
  = jinja_env.from_string(\n                (Path(__file__).parent / \"default_doc_template.md\").read_text(),\n
  \           )\n            return template\n```\n\n\n!! function <h2 id='make_article'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>make_article <em class='small'>function</em></h2>\n\n???+
  source \"make_article <em class='small'>source</em>\"\n\n```python\n\n        def
  make_article(markata: \"Markata\", file: Path, cache) -> frontmatter.Post:\n            with
  open(file) as f:\n                raw_source = f.read()\n            key = markata.make_hash(\"docs\",
  \"file\", raw_source)\n            slug = f\"{file.parent}/{slugify(file.stem)}\".lstrip(\"/\").lstrip(\"./\")\n
  \           output_html = markata.config.output_dir / slug / \"index.html\"\n            edit_link
  = (\n                str(markata.config.get(\"repo_url\", \"https://github.com/\"))\n
  \               + \"edit/\"\n                + str(markata.config.get(\"repo_branch\",
  \"main\"))\n                + \"/\"\n                + str(file)\n            )\n
  \           article_from_cache = markata.precache.get(key)\n            if article_from_cache
  is not None:\n                article = article_from_cache\n            else:\n
  \               tree = ast.parse(raw_source)\n                add_parents(tree)\n
  \               nodes = [\n                    n for n in ast.walk(tree) if isinstance(n,
  (ast.FunctionDef, ast.ClassDef))\n                ]\n\n                article =
  get_template().render(\n                    ast=ast,\n                    file=file,\n
  \                   slug=slug,\n                    edit_link=edit_link,\n                    tree=tree,\n
  \                   datetime=datetime,\n                    nodes=nodes,\n                    raw_source=raw_source,\n
  \                   indent=textwrap.indent,\n                )\n                cache.add(\n
  \                   key,\n                    article,\n                    expire=markata.config.default_cache_expire,\n
  \               )\n\n            try:\n                article = markata.Post(\n
  \                   markata=markata,\n                    path=str(file).replace(\".py\",
  \".md\"),\n                    output_html=output_html,\n                    title=file.name,\n
  \                   content=article,\n                    file=file,\n                    slug=slug,\n
  \                   edit_link=edit_link,\n                    # enable sidebar in
  the future\n                    # sidebar=\"plugins\",\n                )\n\n            except
  pydantic.ValidationError as e:\n                from markata.plugins.load import
  ValidationError, get_models\n\n                models = get_models(markata=markata,
  error=e)\n                models = list(models.values())\n                models
  = \"\\n\".join(models)\n                raise ValidationError(f\"{e}\\n\\n{models}\\nfailed
  to load {path}\") from e\n\n            return article\n```\n\n\n!! function <h2
  id='load' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>load <em
  class='small'>function</em></h2>\n    similar to [glob](../glob)\n???+ source \"load
  <em class='small'>source</em>\"\n\n```python\n\n        def load(markata: \"MarkataDocs\")
  -> None:\n            \"\"\"\n            similar to [glob](../glob)\n            \"\"\"\n
  \           if \"articles\" not in markata.__dict__:\n                markata.articles
  = []\n            for py_file in markata.py_files:\n                with markata.cache
  as cache:\n                    markata.articles.append(make_article(markata, py_file,
  cache))\n```\n\n\n!! class <h2 id='MarkataDocs' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>MarkataDocs <em class='small'>class</em></h2>\n\n???+ source \"MarkataDocs
  <em class='small'>source</em>\"\n\n```python\n\n        class MarkataDocs(Markata):\n
  \               py_files: List = []\n                content_directories: List =
  []\n```\n\n"
date: 0001-01-01
description: 'leading docstring ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Docs.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"leading docstring ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Docs.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"leading docstring ! ???+ source  ! ???+
    source  ! ???+ source  ! ???+ source  ! ! ???+ source \" />\n <link href=\"/favicon.ico\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Docs.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>leading docstring</p>\n<p>!!
    function <h2 id='add_parents' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>add_parents <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">add_parents
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
    <span class=\"nf\">add_parents</span><span class=\"p\">(</span><span class=\"n\">tree</span><span
    class=\"p\">:</span> <span class=\"n\">ast</span><span class=\"o\">.</span><span
    class=\"n\">AST</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">node</span> <span class=\"ow\">in</span> <span class=\"n\">ast</span><span
    class=\"o\">.</span><span class=\"n\">walk</span><span class=\"p\">(</span><span
    class=\"n\">tree</span><span class=\"p\">):</span>\n                <span class=\"k\">for</span>
    <span class=\"n\">child</span> <span class=\"ow\">in</span> <span class=\"n\">ast</span><span
    class=\"o\">.</span><span class=\"n\">iter_child_nodes</span><span class=\"p\">(</span><span
    class=\"n\">node</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">child</span><span class=\"o\">.</span><span class=\"n\">parent</span>
    <span class=\"o\">=</span> <span class=\"n\">node</span>\n                    <span
    class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"nb\">hasattr</span><span
    class=\"p\">(</span><span class=\"n\">child</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;parents&quot;</span><span class=\"p\">):</span>\n                        <span
    class=\"n\">child</span><span class=\"o\">.</span><span class=\"n\">parents</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">node</span><span
    class=\"p\">]</span>\n                    <span class=\"n\">child</span><span
    class=\"o\">.</span><span class=\"n\">parents</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">node</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">node</span><span class=\"p\">,</span> <span
    class=\"n\">ast</span><span class=\"o\">.</span><span class=\"n\">ClassDef</span><span
    class=\"p\">)</span> <span class=\"ow\">and</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">child</span><span class=\"p\">,</span> <span
    class=\"n\">ast</span><span class=\"o\">.</span><span class=\"n\">FunctionDef</span><span
    class=\"p\">):</span>\n                        <span class=\"n\">child</span><span
    class=\"o\">.</span><span class=\"n\">type</span> <span class=\"o\">=</span> <span
    class=\"s2\">&quot;method&quot;</span>\n                    <span class=\"k\">elif</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">child</span><span
    class=\"p\">,</span> <span class=\"n\">ast</span><span class=\"o\">.</span><span
    class=\"n\">FunctionDef</span><span class=\"p\">):</span>\n                        <span
    class=\"n\">child</span><span class=\"o\">.</span><span class=\"n\">type</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;function&quot;</span>\n                    <span
    class=\"k\">elif</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">child</span><span class=\"p\">,</span> <span class=\"n\">ast</span><span
    class=\"o\">.</span><span class=\"n\">ClassDef</span><span class=\"p\">):</span>\n
    \                       <span class=\"n\">child</span><span class=\"o\">.</span><span
    class=\"n\">type</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;class&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='glob' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>glob <em class='small'>function</em></h2>\nfinds k</p>\n<pre><code>## Parameters\n\n`markata`
    the markata object\n</code></pre>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">glob <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">glob</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;MarkataDocs&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            finds k</span>\n\n<span class=\"sd\">            ## Parameters</span>\n\n<span
    class=\"sd\">            `markata` the markata object</span>\n\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n\n            <span class=\"kn\">import</span>
    <span class=\"nn\">glob</span>\n\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">py_files</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span><span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">f</span><span class=\"p\">)</span> <span class=\"k\">for</span> <span
    class=\"n\">f</span> <span class=\"ow\">in</span> <span class=\"n\">glob</span><span
    class=\"o\">.</span><span class=\"n\">glob</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;**/*.py&quot;</span><span class=\"p\">,</span> <span class=\"n\">recursive</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)]</span>\n\n
    \           <span class=\"n\">content_directories</span> <span class=\"o\">=</span>
    <span class=\"nb\">list</span><span class=\"p\">({</span><span class=\"n\">f</span><span
    class=\"o\">.</span><span class=\"n\">parent</span> <span class=\"k\">for</span>
    <span class=\"n\">f</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">py_files</span><span class=\"p\">})</span>\n
    \           <span class=\"k\">if</span> <span class=\"s2\">&quot;content_directories&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"vm\">__dict__</span><span class=\"p\">:</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">content_directories</span><span
    class=\"o\">.</span><span class=\"n\">extend</span><span class=\"p\">(</span><span
    class=\"n\">content_directories</span><span class=\"p\">)</span>\n            <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">content_directories</span> <span class=\"o\">=</span>
    <span class=\"n\">content_directories</span>\n\n            <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                <span class=\"n\">ignore</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n            <span class=\"k\">except</span> <span
    class=\"ne\">KeyError</span><span class=\"p\">:</span>\n                <span
    class=\"n\">ignore</span> <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n
    \           <span class=\"k\">if</span> <span class=\"n\">ignore</span> <span
    class=\"ow\">and</span> <span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.gitignore&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()</span>
    <span class=\"ow\">or</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.markataignore&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()):</span>\n
    \               <span class=\"kn\">import</span> <span class=\"nn\">pathspec</span>\n\n
    \               <span class=\"n\">lines</span> <span class=\"o\">=</span> <span
    class=\"p\">[]</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.gitignore&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">lines</span><span class=\"o\">.</span><span
    class=\"n\">extend</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.gitignore&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">splitlines</span><span class=\"p\">())</span>\n\n
    \               <span class=\"k\">if</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.markataignore&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">():</span>\n                    <span class=\"n\">lines</span><span
    class=\"o\">.</span><span class=\"n\">extend</span><span class=\"p\">(</span><span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"s2\">&quot;.markataignore&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">splitlines</span><span
    class=\"p\">())</span>\n\n            <span class=\"n\">spec</span> <span class=\"o\">=</span>
    <span class=\"n\">pathspec</span><span class=\"o\">.</span><span class=\"n\">PathSpec</span><span
    class=\"o\">.</span><span class=\"n\">from_lines</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;gitwildmatch&quot;</span><span class=\"p\">,</span> <span class=\"n\">lines</span><span
    class=\"p\">)</span>\n\n            <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">py_files</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \               <span class=\"n\">file</span> <span class=\"k\">for</span> <span
    class=\"n\">file</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">py_files</span> <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"n\">spec</span><span class=\"o\">.</span><span
    class=\"n\">match_file</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">file</span><span class=\"p\">))</span>\n
    \           <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!! function
    <h2 id='get_template' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_template
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">get_template <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">get_template</span><span class=\"p\">():</span>\n            <span
    class=\"n\">jinja_env</span> <span class=\"o\">=</span> <span class=\"n\">jinja2</span><span
    class=\"o\">.</span><span class=\"n\">Environment</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">template</span> <span class=\"o\">=</span> <span
    class=\"n\">jinja_env</span><span class=\"o\">.</span><span class=\"n\">from_string</span><span
    class=\"p\">(</span>\n                <span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"vm\">__file__</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">parent</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;default_doc_template.md&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">(),</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">template</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2
    id='make_article' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>make_article
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">make_article <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">make_article</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">file</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">,</span> <span class=\"n\">cache</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">frontmatter</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"p\">:</span>\n            <span class=\"k\">with</span>
    <span class=\"nb\">open</span><span class=\"p\">(</span><span class=\"n\">file</span><span
    class=\"p\">)</span> <span class=\"k\">as</span> <span class=\"n\">f</span><span
    class=\"p\">:</span>\n                <span class=\"n\">raw_source</span> <span
    class=\"o\">=</span> <span class=\"n\">f</span><span class=\"o\">.</span><span
    class=\"n\">read</span><span class=\"p\">()</span>\n            <span class=\"n\">key</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">make_hash</span><span class=\"p\">(</span><span class=\"s2\">&quot;docs&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;file&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">raw_source</span><span class=\"p\">)</span>\n            <span
    class=\"n\">slug</span> <span class=\"o\">=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"si\">}</span><span
    class=\"s2\">/</span><span class=\"si\">{</span><span class=\"n\">slugify</span><span
    class=\"p\">(</span><span class=\"n\">file</span><span class=\"o\">.</span><span
    class=\"n\">stem</span><span class=\"p\">)</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"o\">.</span><span class=\"n\">lstrip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">lstrip</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;./&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">output_html</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">output_dir</span> <span class=\"o\">/</span> <span class=\"n\">slug</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \           <span class=\"n\">edit_link</span> <span class=\"o\">=</span> <span
    class=\"p\">(</span>\n                <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;repo_url&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;https://github.com/&quot;</span><span
    class=\"p\">))</span>\n                <span class=\"o\">+</span> <span class=\"s2\">&quot;edit/&quot;</span>\n
    \               <span class=\"o\">+</span> <span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;repo_branch&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;main&quot;</span><span class=\"p\">))</span>\n                <span
    class=\"o\">+</span> <span class=\"s2\">&quot;/&quot;</span>\n                <span
    class=\"o\">+</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">file</span><span class=\"p\">)</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">article_from_cache</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">precache</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">article_from_cache</span> <span class=\"ow\">is</span> <span
    class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">article</span> <span class=\"o\">=</span> <span
    class=\"n\">article_from_cache</span>\n            <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                <span class=\"n\">tree</span> <span class=\"o\">=</span>
    <span class=\"n\">ast</span><span class=\"o\">.</span><span class=\"n\">parse</span><span
    class=\"p\">(</span><span class=\"n\">raw_source</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">add_parents</span><span class=\"p\">(</span><span
    class=\"n\">tree</span><span class=\"p\">)</span>\n                <span class=\"n\">nodes</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span>\n                    <span
    class=\"n\">n</span> <span class=\"k\">for</span> <span class=\"n\">n</span> <span
    class=\"ow\">in</span> <span class=\"n\">ast</span><span class=\"o\">.</span><span
    class=\"n\">walk</span><span class=\"p\">(</span><span class=\"n\">tree</span><span
    class=\"p\">)</span> <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">n</span><span class=\"p\">,</span> <span
    class=\"p\">(</span><span class=\"n\">ast</span><span class=\"o\">.</span><span
    class=\"n\">FunctionDef</span><span class=\"p\">,</span> <span class=\"n\">ast</span><span
    class=\"o\">.</span><span class=\"n\">ClassDef</span><span class=\"p\">))</span>\n
    \               <span class=\"p\">]</span>\n\n                <span class=\"n\">article</span>
    <span class=\"o\">=</span> <span class=\"n\">get_template</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span>\n
    \                   <span class=\"n\">ast</span><span class=\"o\">=</span><span
    class=\"n\">ast</span><span class=\"p\">,</span>\n                    <span class=\"n\">file</span><span
    class=\"o\">=</span><span class=\"n\">file</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">slug</span><span class=\"o\">=</span><span
    class=\"n\">slug</span><span class=\"p\">,</span>\n                    <span class=\"n\">edit_link</span><span
    class=\"o\">=</span><span class=\"n\">edit_link</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">tree</span><span class=\"o\">=</span><span
    class=\"n\">tree</span><span class=\"p\">,</span>\n                    <span class=\"n\">datetime</span><span
    class=\"o\">=</span><span class=\"n\">datetime</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">nodes</span><span class=\"o\">=</span><span
    class=\"n\">nodes</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">raw_source</span><span class=\"o\">=</span><span class=\"n\">raw_source</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">indent</span><span
    class=\"o\">=</span><span class=\"n\">textwrap</span><span class=\"o\">.</span><span
    class=\"n\">indent</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n
    \               <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">add</span><span class=\"p\">(</span>\n                    <span class=\"n\">key</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">article</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">expire</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">default_cache_expire</span><span
    class=\"p\">,</span>\n                <span class=\"p\">)</span>\n\n            <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                <span class=\"n\">article</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"p\">(</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">path</span><span class=\"o\">=</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">file</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.py&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;.md&quot;</span><span class=\"p\">),</span>\n                    <span
    class=\"n\">output_html</span><span class=\"o\">=</span><span class=\"n\">output_html</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">title</span><span
    class=\"o\">=</span><span class=\"n\">file</span><span class=\"o\">.</span><span
    class=\"n\">name</span><span class=\"p\">,</span>\n                    <span class=\"n\">content</span><span
    class=\"o\">=</span><span class=\"n\">article</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">file</span><span class=\"o\">=</span><span
    class=\"n\">file</span><span class=\"p\">,</span>\n                    <span class=\"n\">slug</span><span
    class=\"o\">=</span><span class=\"n\">slug</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">edit_link</span><span class=\"o\">=</span><span
    class=\"n\">edit_link</span><span class=\"p\">,</span>\n                    <span
    class=\"c1\"># enable sidebar in the future</span>\n                    <span
    class=\"c1\"># sidebar=&quot;plugins&quot;,</span>\n                <span class=\"p\">)</span>\n\n
    \           <span class=\"k\">except</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">ValidationError</span> <span class=\"k\">as</span>
    <span class=\"n\">e</span><span class=\"p\">:</span>\n                <span class=\"kn\">from</span>
    <span class=\"nn\">markata.plugins.load</span> <span class=\"kn\">import</span>
    <span class=\"n\">ValidationError</span><span class=\"p\">,</span> <span class=\"n\">get_models</span>\n\n
    \               <span class=\"n\">models</span> <span class=\"o\">=</span> <span
    class=\"n\">get_models</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">error</span><span class=\"o\">=</span><span class=\"n\">e</span><span
    class=\"p\">)</span>\n                <span class=\"n\">models</span> <span class=\"o\">=</span>
    <span class=\"nb\">list</span><span class=\"p\">(</span><span class=\"n\">models</span><span
    class=\"o\">.</span><span class=\"n\">values</span><span class=\"p\">())</span>\n
    \               <span class=\"n\">models</span> <span class=\"o\">=</span> <span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">models</span><span class=\"p\">)</span>\n                <span class=\"k\">raise</span>
    <span class=\"n\">ValidationError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">e</span><span
    class=\"si\">}</span><span class=\"se\">\\n\\n</span><span class=\"si\">{</span><span
    class=\"n\">models</span><span class=\"si\">}</span><span class=\"se\">\\n</span><span
    class=\"s2\">failed to load </span><span class=\"si\">{</span><span class=\"n\">path</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>
    <span class=\"kn\">from</span> <span class=\"nn\">e</span>\n\n            <span
    class=\"k\">return</span> <span class=\"n\">article</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='load' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>load <em class='small'>function</em></h2>\nsimilar to <a href=\"../glob\">glob</a></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">load
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
    <span class=\"nf\">load</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;MarkataDocs&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            similar to [glob](../glob)</span>\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n            <span class=\"k\">if</span>
    <span class=\"s2\">&quot;articles&quot;</span> <span class=\"ow\">not</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"vm\">__dict__</span><span class=\"p\">:</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">articles</span>
    <span class=\"o\">=</span> <span class=\"p\">[]</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">py_file</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">py_files</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">with</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">articles</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">make_article</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">py_file</span><span class=\"p\">,</span>
    <span class=\"n\">cache</span><span class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='MarkataDocs' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>MarkataDocs <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">MarkataDocs
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
    <span class=\"nc\">MarkataDocs</span><span class=\"p\">(</span><span class=\"n\">Markata</span><span
    class=\"p\">):</span>\n                <span class=\"n\">py_files</span><span
    class=\"p\">:</span> <span class=\"n\">List</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n                <span class=\"n\">content_directories</span><span
    class=\"p\">:</span> <span class=\"n\">List</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Docs.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"leading docstring ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Docs.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"leading docstring ! ???+ source  ! ???+
    source  ! ???+ source  ! ???+ source  ! ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Docs.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Docs.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>leading
    docstring</p>\n<p>!! function <h2 id='add_parents' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>add_parents <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">add_parents
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
    <span class=\"nf\">add_parents</span><span class=\"p\">(</span><span class=\"n\">tree</span><span
    class=\"p\">:</span> <span class=\"n\">ast</span><span class=\"o\">.</span><span
    class=\"n\">AST</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">node</span> <span class=\"ow\">in</span> <span class=\"n\">ast</span><span
    class=\"o\">.</span><span class=\"n\">walk</span><span class=\"p\">(</span><span
    class=\"n\">tree</span><span class=\"p\">):</span>\n                <span class=\"k\">for</span>
    <span class=\"n\">child</span> <span class=\"ow\">in</span> <span class=\"n\">ast</span><span
    class=\"o\">.</span><span class=\"n\">iter_child_nodes</span><span class=\"p\">(</span><span
    class=\"n\">node</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">child</span><span class=\"o\">.</span><span class=\"n\">parent</span>
    <span class=\"o\">=</span> <span class=\"n\">node</span>\n                    <span
    class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"nb\">hasattr</span><span
    class=\"p\">(</span><span class=\"n\">child</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;parents&quot;</span><span class=\"p\">):</span>\n                        <span
    class=\"n\">child</span><span class=\"o\">.</span><span class=\"n\">parents</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">node</span><span
    class=\"p\">]</span>\n                    <span class=\"n\">child</span><span
    class=\"o\">.</span><span class=\"n\">parents</span><span class=\"o\">.</span><span
    class=\"n\">append</span><span class=\"p\">(</span><span class=\"n\">node</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">node</span><span class=\"p\">,</span> <span
    class=\"n\">ast</span><span class=\"o\">.</span><span class=\"n\">ClassDef</span><span
    class=\"p\">)</span> <span class=\"ow\">and</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">child</span><span class=\"p\">,</span> <span
    class=\"n\">ast</span><span class=\"o\">.</span><span class=\"n\">FunctionDef</span><span
    class=\"p\">):</span>\n                        <span class=\"n\">child</span><span
    class=\"o\">.</span><span class=\"n\">type</span> <span class=\"o\">=</span> <span
    class=\"s2\">&quot;method&quot;</span>\n                    <span class=\"k\">elif</span>
    <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">child</span><span
    class=\"p\">,</span> <span class=\"n\">ast</span><span class=\"o\">.</span><span
    class=\"n\">FunctionDef</span><span class=\"p\">):</span>\n                        <span
    class=\"n\">child</span><span class=\"o\">.</span><span class=\"n\">type</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;function&quot;</span>\n                    <span
    class=\"k\">elif</span> <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span
    class=\"n\">child</span><span class=\"p\">,</span> <span class=\"n\">ast</span><span
    class=\"o\">.</span><span class=\"n\">ClassDef</span><span class=\"p\">):</span>\n
    \                       <span class=\"n\">child</span><span class=\"o\">.</span><span
    class=\"n\">type</span> <span class=\"o\">=</span> <span class=\"s2\">&quot;class&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='glob' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>glob <em class='small'>function</em></h2>\nfinds k</p>\n<pre><code>## Parameters\n\n`markata`
    the markata object\n</code></pre>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">glob <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">glob</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;MarkataDocs&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            finds k</span>\n\n<span class=\"sd\">            ## Parameters</span>\n\n<span
    class=\"sd\">            `markata` the markata object</span>\n\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n\n            <span class=\"kn\">import</span>
    <span class=\"nn\">glob</span>\n\n            <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">py_files</span> <span class=\"o\">=</span>
    <span class=\"p\">[</span><span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"n\">f</span><span class=\"p\">)</span> <span class=\"k\">for</span> <span
    class=\"n\">f</span> <span class=\"ow\">in</span> <span class=\"n\">glob</span><span
    class=\"o\">.</span><span class=\"n\">glob</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;**/*.py&quot;</span><span class=\"p\">,</span> <span class=\"n\">recursive</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)]</span>\n\n
    \           <span class=\"n\">content_directories</span> <span class=\"o\">=</span>
    <span class=\"nb\">list</span><span class=\"p\">({</span><span class=\"n\">f</span><span
    class=\"o\">.</span><span class=\"n\">parent</span> <span class=\"k\">for</span>
    <span class=\"n\">f</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">py_files</span><span class=\"p\">})</span>\n
    \           <span class=\"k\">if</span> <span class=\"s2\">&quot;content_directories&quot;</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"vm\">__dict__</span><span class=\"p\">:</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">content_directories</span><span
    class=\"o\">.</span><span class=\"n\">extend</span><span class=\"p\">(</span><span
    class=\"n\">content_directories</span><span class=\"p\">)</span>\n            <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">content_directories</span> <span class=\"o\">=</span>
    <span class=\"n\">content_directories</span>\n\n            <span class=\"k\">try</span><span
    class=\"p\">:</span>\n                <span class=\"n\">ignore</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n            <span class=\"k\">except</span> <span
    class=\"ne\">KeyError</span><span class=\"p\">:</span>\n                <span
    class=\"n\">ignore</span> <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n
    \           <span class=\"k\">if</span> <span class=\"n\">ignore</span> <span
    class=\"ow\">and</span> <span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.gitignore&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()</span>
    <span class=\"ow\">or</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;.markataignore&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">()):</span>\n
    \               <span class=\"kn\">import</span> <span class=\"nn\">pathspec</span>\n\n
    \               <span class=\"n\">lines</span> <span class=\"o\">=</span> <span
    class=\"p\">[]</span>\n\n                <span class=\"k\">if</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.gitignore&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">exists</span><span class=\"p\">():</span>\n
    \                   <span class=\"n\">lines</span><span class=\"o\">.</span><span
    class=\"n\">extend</span><span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.gitignore&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">splitlines</span><span class=\"p\">())</span>\n\n
    \               <span class=\"k\">if</span> <span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.markataignore&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">exists</span><span
    class=\"p\">():</span>\n                    <span class=\"n\">lines</span><span
    class=\"o\">.</span><span class=\"n\">extend</span><span class=\"p\">(</span><span
    class=\"n\">Path</span><span class=\"p\">(</span><span class=\"s2\">&quot;.markataignore&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">read_text</span><span
    class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">splitlines</span><span
    class=\"p\">())</span>\n\n            <span class=\"n\">spec</span> <span class=\"o\">=</span>
    <span class=\"n\">pathspec</span><span class=\"o\">.</span><span class=\"n\">PathSpec</span><span
    class=\"o\">.</span><span class=\"n\">from_lines</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;gitwildmatch&quot;</span><span class=\"p\">,</span> <span class=\"n\">lines</span><span
    class=\"p\">)</span>\n\n            <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">py_files</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \               <span class=\"n\">file</span> <span class=\"k\">for</span> <span
    class=\"n\">file</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">py_files</span> <span class=\"k\">if</span>
    <span class=\"ow\">not</span> <span class=\"n\">spec</span><span class=\"o\">.</span><span
    class=\"n\">match_file</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">file</span><span class=\"p\">))</span>\n
    \           <span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>!! function
    <h2 id='get_template' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_template
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">get_template <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">get_template</span><span class=\"p\">():</span>\n            <span
    class=\"n\">jinja_env</span> <span class=\"o\">=</span> <span class=\"n\">jinja2</span><span
    class=\"o\">.</span><span class=\"n\">Environment</span><span class=\"p\">()</span>\n
    \           <span class=\"n\">template</span> <span class=\"o\">=</span> <span
    class=\"n\">jinja_env</span><span class=\"o\">.</span><span class=\"n\">from_string</span><span
    class=\"p\">(</span>\n                <span class=\"p\">(</span><span class=\"n\">Path</span><span
    class=\"p\">(</span><span class=\"vm\">__file__</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">parent</span> <span class=\"o\">/</span>
    <span class=\"s2\">&quot;default_doc_template.md&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">read_text</span><span class=\"p\">(),</span>\n
    \           <span class=\"p\">)</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">template</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2
    id='make_article' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>make_article
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">make_article <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">make_article</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">file</span><span class=\"p\">:</span> <span class=\"n\">Path</span><span
    class=\"p\">,</span> <span class=\"n\">cache</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">frontmatter</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"p\">:</span>\n            <span class=\"k\">with</span>
    <span class=\"nb\">open</span><span class=\"p\">(</span><span class=\"n\">file</span><span
    class=\"p\">)</span> <span class=\"k\">as</span> <span class=\"n\">f</span><span
    class=\"p\">:</span>\n                <span class=\"n\">raw_source</span> <span
    class=\"o\">=</span> <span class=\"n\">f</span><span class=\"o\">.</span><span
    class=\"n\">read</span><span class=\"p\">()</span>\n            <span class=\"n\">key</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">make_hash</span><span class=\"p\">(</span><span class=\"s2\">&quot;docs&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;file&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">raw_source</span><span class=\"p\">)</span>\n            <span
    class=\"n\">slug</span> <span class=\"o\">=</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">file</span><span
    class=\"o\">.</span><span class=\"n\">parent</span><span class=\"si\">}</span><span
    class=\"s2\">/</span><span class=\"si\">{</span><span class=\"n\">slugify</span><span
    class=\"p\">(</span><span class=\"n\">file</span><span class=\"o\">.</span><span
    class=\"n\">stem</span><span class=\"p\">)</span><span class=\"si\">}</span><span
    class=\"s2\">&quot;</span><span class=\"o\">.</span><span class=\"n\">lstrip</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;/&quot;</span><span class=\"p\">)</span><span
    class=\"o\">.</span><span class=\"n\">lstrip</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;./&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">output_html</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">output_dir</span> <span class=\"o\">/</span> <span class=\"n\">slug</span>
    <span class=\"o\">/</span> <span class=\"s2\">&quot;index.html&quot;</span>\n
    \           <span class=\"n\">edit_link</span> <span class=\"o\">=</span> <span
    class=\"p\">(</span>\n                <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;repo_url&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;https://github.com/&quot;</span><span
    class=\"p\">))</span>\n                <span class=\"o\">+</span> <span class=\"s2\">&quot;edit/&quot;</span>\n
    \               <span class=\"o\">+</span> <span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;repo_branch&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;main&quot;</span><span class=\"p\">))</span>\n                <span
    class=\"o\">+</span> <span class=\"s2\">&quot;/&quot;</span>\n                <span
    class=\"o\">+</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">file</span><span class=\"p\">)</span>\n            <span class=\"p\">)</span>\n
    \           <span class=\"n\">article_from_cache</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">precache</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">)</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">article_from_cache</span> <span class=\"ow\">is</span> <span
    class=\"ow\">not</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"n\">article</span> <span class=\"o\">=</span> <span
    class=\"n\">article_from_cache</span>\n            <span class=\"k\">else</span><span
    class=\"p\">:</span>\n                <span class=\"n\">tree</span> <span class=\"o\">=</span>
    <span class=\"n\">ast</span><span class=\"o\">.</span><span class=\"n\">parse</span><span
    class=\"p\">(</span><span class=\"n\">raw_source</span><span class=\"p\">)</span>\n
    \               <span class=\"n\">add_parents</span><span class=\"p\">(</span><span
    class=\"n\">tree</span><span class=\"p\">)</span>\n                <span class=\"n\">nodes</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span>\n                    <span
    class=\"n\">n</span> <span class=\"k\">for</span> <span class=\"n\">n</span> <span
    class=\"ow\">in</span> <span class=\"n\">ast</span><span class=\"o\">.</span><span
    class=\"n\">walk</span><span class=\"p\">(</span><span class=\"n\">tree</span><span
    class=\"p\">)</span> <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">n</span><span class=\"p\">,</span> <span
    class=\"p\">(</span><span class=\"n\">ast</span><span class=\"o\">.</span><span
    class=\"n\">FunctionDef</span><span class=\"p\">,</span> <span class=\"n\">ast</span><span
    class=\"o\">.</span><span class=\"n\">ClassDef</span><span class=\"p\">))</span>\n
    \               <span class=\"p\">]</span>\n\n                <span class=\"n\">article</span>
    <span class=\"o\">=</span> <span class=\"n\">get_template</span><span class=\"p\">()</span><span
    class=\"o\">.</span><span class=\"n\">render</span><span class=\"p\">(</span>\n
    \                   <span class=\"n\">ast</span><span class=\"o\">=</span><span
    class=\"n\">ast</span><span class=\"p\">,</span>\n                    <span class=\"n\">file</span><span
    class=\"o\">=</span><span class=\"n\">file</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">slug</span><span class=\"o\">=</span><span
    class=\"n\">slug</span><span class=\"p\">,</span>\n                    <span class=\"n\">edit_link</span><span
    class=\"o\">=</span><span class=\"n\">edit_link</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">tree</span><span class=\"o\">=</span><span
    class=\"n\">tree</span><span class=\"p\">,</span>\n                    <span class=\"n\">datetime</span><span
    class=\"o\">=</span><span class=\"n\">datetime</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">nodes</span><span class=\"o\">=</span><span
    class=\"n\">nodes</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">raw_source</span><span class=\"o\">=</span><span class=\"n\">raw_source</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">indent</span><span
    class=\"o\">=</span><span class=\"n\">textwrap</span><span class=\"o\">.</span><span
    class=\"n\">indent</span><span class=\"p\">,</span>\n                <span class=\"p\">)</span>\n
    \               <span class=\"n\">cache</span><span class=\"o\">.</span><span
    class=\"n\">add</span><span class=\"p\">(</span>\n                    <span class=\"n\">key</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">article</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">expire</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">default_cache_expire</span><span
    class=\"p\">,</span>\n                <span class=\"p\">)</span>\n\n            <span
    class=\"k\">try</span><span class=\"p\">:</span>\n                <span class=\"n\">article</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">Post</span><span class=\"p\">(</span>\n                    <span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">path</span><span class=\"o\">=</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">file</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">replace</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;.py&quot;</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;.md&quot;</span><span class=\"p\">),</span>\n                    <span
    class=\"n\">output_html</span><span class=\"o\">=</span><span class=\"n\">output_html</span><span
    class=\"p\">,</span>\n                    <span class=\"n\">title</span><span
    class=\"o\">=</span><span class=\"n\">file</span><span class=\"o\">.</span><span
    class=\"n\">name</span><span class=\"p\">,</span>\n                    <span class=\"n\">content</span><span
    class=\"o\">=</span><span class=\"n\">article</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">file</span><span class=\"o\">=</span><span
    class=\"n\">file</span><span class=\"p\">,</span>\n                    <span class=\"n\">slug</span><span
    class=\"o\">=</span><span class=\"n\">slug</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">edit_link</span><span class=\"o\">=</span><span
    class=\"n\">edit_link</span><span class=\"p\">,</span>\n                    <span
    class=\"c1\"># enable sidebar in the future</span>\n                    <span
    class=\"c1\"># sidebar=&quot;plugins&quot;,</span>\n                <span class=\"p\">)</span>\n\n
    \           <span class=\"k\">except</span> <span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">ValidationError</span> <span class=\"k\">as</span>
    <span class=\"n\">e</span><span class=\"p\">:</span>\n                <span class=\"kn\">from</span>
    <span class=\"nn\">markata.plugins.load</span> <span class=\"kn\">import</span>
    <span class=\"n\">ValidationError</span><span class=\"p\">,</span> <span class=\"n\">get_models</span>\n\n
    \               <span class=\"n\">models</span> <span class=\"o\">=</span> <span
    class=\"n\">get_models</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">=</span><span class=\"n\">markata</span><span class=\"p\">,</span>
    <span class=\"n\">error</span><span class=\"o\">=</span><span class=\"n\">e</span><span
    class=\"p\">)</span>\n                <span class=\"n\">models</span> <span class=\"o\">=</span>
    <span class=\"nb\">list</span><span class=\"p\">(</span><span class=\"n\">models</span><span
    class=\"o\">.</span><span class=\"n\">values</span><span class=\"p\">())</span>\n
    \               <span class=\"n\">models</span> <span class=\"o\">=</span> <span
    class=\"s2\">&quot;</span><span class=\"se\">\\n</span><span class=\"s2\">&quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span><span
    class=\"n\">models</span><span class=\"p\">)</span>\n                <span class=\"k\">raise</span>
    <span class=\"n\">ValidationError</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;</span><span class=\"si\">{</span><span class=\"n\">e</span><span
    class=\"si\">}</span><span class=\"se\">\\n\\n</span><span class=\"si\">{</span><span
    class=\"n\">models</span><span class=\"si\">}</span><span class=\"se\">\\n</span><span
    class=\"s2\">failed to load </span><span class=\"si\">{</span><span class=\"n\">path</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">)</span>
    <span class=\"kn\">from</span> <span class=\"nn\">e</span>\n\n            <span
    class=\"k\">return</span> <span class=\"n\">article</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='load' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>load <em class='small'>function</em></h2>\nsimilar to <a href=\"../glob\">glob</a></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">load
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
    <span class=\"nf\">load</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;MarkataDocs&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            similar to [glob](../glob)</span>\n<span class=\"sd\">
    \           &quot;&quot;&quot;</span>\n            <span class=\"k\">if</span>
    <span class=\"s2\">&quot;articles&quot;</span> <span class=\"ow\">not</span> <span
    class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"vm\">__dict__</span><span class=\"p\">:</span>\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">articles</span>
    <span class=\"o\">=</span> <span class=\"p\">[]</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">py_file</span> <span class=\"ow\">in</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">py_files</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">with</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">cache</span> <span class=\"k\">as</span>
    <span class=\"n\">cache</span><span class=\"p\">:</span>\n                    <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">articles</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">make_article</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">py_file</span><span class=\"p\">,</span>
    <span class=\"n\">cache</span><span class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='MarkataDocs' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>MarkataDocs <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">MarkataDocs
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
    <span class=\"nc\">MarkataDocs</span><span class=\"p\">(</span><span class=\"n\">Markata</span><span
    class=\"p\">):</span>\n                <span class=\"n\">py_files</span><span
    class=\"p\">:</span> <span class=\"n\">List</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n                <span class=\"n\">content_directories</span><span
    class=\"p\">:</span> <span class=\"n\">List</span> <span class=\"o\">=</span>
    <span class=\"p\">[]</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/docs
title: Docs.Py


---

leading docstring


!! function <h2 id='add_parents' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>add_parents <em class='small'>function</em></h2>

???+ source "add_parents <em class='small'>source</em>"

```python

        def add_parents(tree: ast.AST) -> None:
            for node in ast.walk(tree):
                for child in ast.iter_child_nodes(node):
                    child.parent = node
                    if not hasattr(child, "parents"):
                        child.parents = [node]
                    child.parents.append(node)
                    if isinstance(node, ast.ClassDef) and isinstance(child, ast.FunctionDef):
                        child.type = "method"
                    elif isinstance(child, ast.FunctionDef):
                        child.type = "function"
                    elif isinstance(child, ast.ClassDef):
                        child.type = "class"
```


!! function <h2 id='glob' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>glob <em class='small'>function</em></h2>
    finds k

    ## Parameters

    `markata` the markata object
???+ source "glob <em class='small'>source</em>"

```python

        def glob(markata: "MarkataDocs") -> None:
            """
            finds k

            ## Parameters

            `markata` the markata object

            """

            import glob

            markata.py_files = [Path(f) for f in glob.glob("**/*.py", recursive=True)]

            content_directories = list({f.parent for f in markata.py_files})
            if "content_directories" in markata.__dict__:
                markata.content_directories.extend(content_directories)
            else:
                markata.content_directories = content_directories

            try:
                ignore = True
            except KeyError:
                ignore = True

            if ignore and (Path(".gitignore").exists() or Path(".markataignore").exists()):
                import pathspec

                lines = []

                if Path(".gitignore").exists():
                    lines.extend(Path(".gitignore").read_text().splitlines())

                if Path(".markataignore").exists():
                    lines.extend(Path(".markataignore").read_text().splitlines())

            spec = pathspec.PathSpec.from_lines("gitwildmatch", lines)

            markata.py_files = [
                file for file in markata.py_files if not spec.match_file(str(file))
            ]
```


!! function <h2 id='get_template' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_template <em class='small'>function</em></h2>

???+ source "get_template <em class='small'>source</em>"

```python

        def get_template():
            jinja_env = jinja2.Environment()
            template = jinja_env.from_string(
                (Path(__file__).parent / "default_doc_template.md").read_text(),
            )
            return template
```


!! function <h2 id='make_article' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>make_article <em class='small'>function</em></h2>

???+ source "make_article <em class='small'>source</em>"

```python

        def make_article(markata: "Markata", file: Path, cache) -> frontmatter.Post:
            with open(file) as f:
                raw_source = f.read()
            key = markata.make_hash("docs", "file", raw_source)
            slug = f"{file.parent}/{slugify(file.stem)}".lstrip("/").lstrip("./")
            output_html = markata.config.output_dir / slug / "index.html"
            edit_link = (
                str(markata.config.get("repo_url", "https://github.com/"))
                + "edit/"
                + str(markata.config.get("repo_branch", "main"))
                + "/"
                + str(file)
            )
            article_from_cache = markata.precache.get(key)
            if article_from_cache is not None:
                article = article_from_cache
            else:
                tree = ast.parse(raw_source)
                add_parents(tree)
                nodes = [
                    n for n in ast.walk(tree) if isinstance(n, (ast.FunctionDef, ast.ClassDef))
                ]

                article = get_template().render(
                    ast=ast,
                    file=file,
                    slug=slug,
                    edit_link=edit_link,
                    tree=tree,
                    datetime=datetime,
                    nodes=nodes,
                    raw_source=raw_source,
                    indent=textwrap.indent,
                )
                cache.add(
                    key,
                    article,
                    expire=markata.config.default_cache_expire,
                )

            try:
                article = markata.Post(
                    markata=markata,
                    path=str(file).replace(".py", ".md"),
                    output_html=output_html,
                    title=file.name,
                    content=article,
                    file=file,
                    slug=slug,
                    edit_link=edit_link,
                    # enable sidebar in the future
                    # sidebar="plugins",
                )

            except pydantic.ValidationError as e:
                from markata.plugins.load import ValidationError, get_models

                models = get_models(markata=markata, error=e)
                models = list(models.values())
                models = "\n".join(models)
                raise ValidationError(f"{e}\n\n{models}\nfailed to load {path}") from e

            return article
```


!! function <h2 id='load' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>load <em class='small'>function</em></h2>
    similar to [glob](../glob)
???+ source "load <em class='small'>source</em>"

```python

        def load(markata: "MarkataDocs") -> None:
            """
            similar to [glob](../glob)
            """
            if "articles" not in markata.__dict__:
                markata.articles = []
            for py_file in markata.py_files:
                with markata.cache as cache:
                    markata.articles.append(make_article(markata, py_file, cache))
```


!! class <h2 id='MarkataDocs' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>MarkataDocs <em class='small'>class</em></h2>

???+ source "MarkataDocs <em class='small'>source</em>"

```python

        class MarkataDocs(Markata):
                py_files: List = []
                content_directories: List = []
```

