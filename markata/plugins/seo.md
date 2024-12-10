---
content: "manifest plugin\n\n\n!! function <h2 id='_create_seo' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>_create_seo <em class='small'>function</em></h2>\n\n???+
  source \"_create_seo <em class='small'>source</em>\"\n\n```python\n\n        def
  _create_seo(\n            markata: Markata,\n            soup: BeautifulSoup,\n
  \           article: \"frontmatter.Post\",\n            site_name: str,\n            author_name:
  str,\n            author_email: str,\n            twitter_card: str,\n            twitter_creator:
  str,\n            config_seo: Dict,\n            images_url: str,\n        ) ->
  List:\n            if article.metadata[\"description\"] == \"\" or None:\n                try:\n
  \                   article.metadata[\"description\"] = \" \".join(\n                        [p.text
  for p in soup.find(id=\"post-body\").find_all(\"p\")],\n                    ).strip()[:120]\n
  \               except AttributeError:\n                    article.metadata[\"description\"]
  = \"\"\n\n            seo = [\n                *config_seo,\n                {\n
  \                   \"name\": \"og:author\",\n                    \"property\":
  \"og:author\",\n                    \"content\": author_name,\n                },\n
  \               {\n                    \"name\": \"og:author_email\",\n                    \"property\":
  \"og:author_email\",\n                    \"content\": author_email,\n                },\n
  \               {\n                    \"name\": \"og:type\",\n                    \"property\":
  \"og:type\",\n                    \"content\": \"website\",\n                },\n
  \               {\n                    \"name\": \"description\",\n                    \"property\":
  \"description\",\n                    \"content\": article.metadata[\"description\"],\n
  \               },\n                {\n                    \"name\": \"og:description\",\n
  \                   \"property\": \"og:description\",\n                    \"content\":
  article.metadata[\"description\"],\n                },\n                {\n                    \"name\":
  \"twitter:description\",\n                    \"property\": \"twitter:description\",\n
  \                   \"content\": article.metadata[\"description\"],\n                },\n
  \               {\n                    \"name\": \"og:title\",\n                    \"property\":
  \"og:title\",\n                    \"content\": f'{article.metadata[\"title\"]}
  | {site_name}'[:60],\n                },\n                {\n                    \"name\":
  \"twitter:title\",\n                    \"property\": \"twitter:title\",\n                    \"content\":
  f'{article.metadata[\"title\"]} | {site_name}'[:60],\n                },\n                {\n
  \                   \"name\": \"og:image\",\n                    \"property\": \"og:image\",\n
  \                   \"content\": f'{images_url}/{article.metadata[\"slug\"]}-og.png',\n
  \               },\n                {\n                    \"name\": \"twitter:image\",\n
  \                   \"property\": \"twitter:image\",\n                    \"content\":
  f'{images_url}/{article.metadata[\"slug\"]}-og.png',\n                },\n                {\n
  \                   \"name\": \"og:image:width\",\n                    \"property\":
  \"og:image:width\",\n                    \"content\": \"1600\",\n                },\n
  \               {\n                    \"name\": \"og:image:width\",\n                    \"property\":
  \"og:image:width\",\n                    \"content\": \"900\",\n                },\n
  \               {\n                    \"name\": \"twitter:card\",\n                    \"property\":
  \"twitter:card\",\n                    \"content\": twitter_card,\n                },\n
  \               {\n                    \"name\": \"og:site_name\",\n                    \"property\":
  \"og:site_name\",\n                    \"content\": site_name,\n                },\n
  \               {\n                    \"name\": \"twitter:creator\",\n                    \"property\":
  \"twitter:creator\",\n                    \"content\": twitter_creator,\n                },\n
  \               {\n                    \"name\": \"title\",\n                    \"property\":
  \"title\",\n                    \"content\": article.metadata[\"title\"],\n                },\n
  \               {\n                    \"name\": \"generator\",\n                    \"property\":
  \"generator\",\n                    \"content\": f\"markata {__version__}\",\n                },\n
  \           ]\n            return seo\n```\n\n\n!! function <h2 id='_add_seo_tags'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_add_seo_tags <em
  class='small'>function</em></h2>\n\n???+ source \"_add_seo_tags <em class='small'>source</em>\"\n\n```python\n\n
  \       def _add_seo_tags(seo: List, article: \"frontmatter.Post\", soup: BeautifulSoup)
  -> None:\n            for meta in seo:\n                soup.head.append(_create_seo_tag(meta,
  soup))\n```\n\n\n!! function <h2 id='_create_seo_tag' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>_create_seo_tag <em class='small'>function</em></h2>\n\n???+ source \"_create_seo_tag
  <em class='small'>source</em>\"\n\n```python\n\n        def _create_seo_tag(meta:
  dict, soup: BeautifulSoup) -> \"Tag\":\n            tag = soup.new_tag(\"meta\")\n
  \           for k in meta:\n                tag.attrs[k] = meta[k]\n            return
  tag\n```\n\n\n!! function <h2 id='_get_or_warn' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>_get_or_warn <em class='small'>function</em></h2>\n\n???+ source \"_get_or_warn
  <em class='small'>source</em>\"\n\n```python\n\n        def _get_or_warn(config:
  Dict, key: str, default: str) -> Any:\n            if key not in config.keys():\n
  \               logger.warning(\n                    f\"{key} is missing from markata.toml
  config, using default value {default}\",\n                )\n            return
  config.get(key, default)\n```\n\n\n!! function <h2 id='render' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2>\n\n???+
  source \"render <em class='small'>source</em>\"\n\n```python\n\n        def render(markata:
  Markata) -> None:\n            url = _get_or_warn(markata.config, \"url\", \"\")\n
  \           images_url = markata.config.get(\"images_url\", url)\n            site_name
  = _get_or_warn(markata.config, \"site_name\", \"\")\n            author_name = _get_or_warn(markata.config,
  \"author_name\", \"\")\n            author_email = _get_or_warn(markata.config,
  \"author_email\", \"\")\n            twitter_creator = _get_or_warn(markata.config,
  \"twitter_creator\", \"\")\n            twitter_card = _get_or_warn(markata.config,
  \"twitter_card\", \"summary_large_image\")\n            config_seo = markata.config.get(\"seo\",
  {})\n            should_prettify = markata.config.get(\"prettify_html\", False)\n\n
  \           with markata.cache as cache:\n                for article in markata.iter_articles(\"add
  seo tags from seo.py\"):\n                    key = markata.make_hash(\n                        \"seo\",\n
  \                       \"render\",\n                        article.html,\n                        site_name,\n
  \                       url,\n                        article.metadata[\"slug\"],\n
  \                       twitter_card,\n                        article.metadata[\"title\"],\n
  \                       str(config_seo),\n                    )\n\n                    html_from_cache
  = markata.precache.get(key)\n\n                    if html_from_cache is None:\n
  \                       soup = BeautifulSoup(article.html, features=\"lxml\")\n
  \                       seo = _create_seo(\n                            markata=markata,\n
  \                           soup=soup,\n                            article=article,\n
  \                           site_name=site_name,\n                            author_name=author_name,\n
  \                           author_email=author_email,\n                            twitter_card=twitter_card,\n
  \                           twitter_creator=twitter_creator,\n                            config_seo=config_seo,\n
  \                           images_url=images_url,\n                        )\n
  \                       _add_seo_tags(seo, article, soup)\n                        canonical_link
  = soup.new_tag(\"link\")\n                        canonical_link.attrs[\"rel\"]
  = \"canonical\"\n                        if article.metadata[\"slug\"] == \"index\":\n
  \                           canonical_link.attrs[\"href\"] = f\"{url}/\"\n                        else:\n
  \                           canonical_link.attrs[\"href\"] = f'{url}/{article.metadata[\"slug\"]}/'\n
  \                       soup.head.append(canonical_link)\n\n                        meta_url
  = soup.new_tag(\"meta\")\n                        meta_url.attrs[\"name\"] = \"og:url\"\n
  \                       meta_url.attrs[\"property\"] = \"og:url\"\n                        if
  article.metadata[\"slug\"] == \"index\":\n                            meta_url.attrs[\"content\"]
  = f\"{url}/\"\n                        else:\n                            meta_url.attrs[\"content\"]
  = f'{url}/{article.metadata[\"slug\"]}/'\n                        soup.head.append(meta_url)\n\n
  \                       html = soup.prettify() if should_prettify else str(soup)\n
  \                       cache.add(key, html, expire=markata.config.default_cache_expire)\n\n
  \                   else:\n                        html = html_from_cache\n                    article.html
  = html\n```\n\n"
date: 0001-01-01
description: 'manifest plugin ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ???+ source '
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Seo.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"manifest plugin ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Seo.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"manifest plugin ! ???+ source  ! ???+
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Seo.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>manifest plugin</p>\n<p>!!
    function <h2 id='_create_seo' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_create_seo <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_create_seo
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
    <span class=\"nf\">_create_seo</span><span class=\"p\">(</span>\n            <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">,</span>\n            <span class=\"n\">soup</span><span class=\"p\">:</span>
    <span class=\"n\">BeautifulSoup</span><span class=\"p\">,</span>\n            <span
    class=\"n\">article</span><span class=\"p\">:</span> <span class=\"s2\">&quot;frontmatter.Post&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"n\">site_name</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">,</span>\n            <span class=\"n\">author_name</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">author_email</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">,</span>\n            <span class=\"n\">twitter_card</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">twitter_creator</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">,</span>\n            <span class=\"n\">config_seo</span><span
    class=\"p\">:</span> <span class=\"n\">Dict</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">images_url</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">,</span>\n        <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;description&quot;</span><span class=\"p\">]</span> <span class=\"o\">==</span>
    <span class=\"s2\">&quot;&quot;</span> <span class=\"ow\">or</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"s2\">&quot; &quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span>\n
    \                       <span class=\"p\">[</span><span class=\"n\">p</span><span
    class=\"o\">.</span><span class=\"n\">text</span> <span class=\"k\">for</span>
    <span class=\"n\">p</span> <span class=\"ow\">in</span> <span class=\"n\">soup</span><span
    class=\"o\">.</span><span class=\"n\">find</span><span class=\"p\">(</span><span
    class=\"nb\">id</span><span class=\"o\">=</span><span class=\"s2\">&quot;post-body&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">find_all</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;p&quot;</span><span class=\"p\">)],</span>\n
    \                   <span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">()[:</span><span class=\"mi\">120</span><span
    class=\"p\">]</span>\n                <span class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;description&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;&quot;</span>\n\n            <span class=\"n\">seo</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span>\n                <span class=\"o\">*</span><span
    class=\"n\">config_seo</span><span class=\"p\">,</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;name&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;og:author&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:author&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span> <span class=\"n\">author_name</span><span
    class=\"p\">,</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;name&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;og:author_email&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;property&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;og:author_email&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">author_email</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:type&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;og:type&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;website&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;description&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;description&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:description&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:description&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;twitter:description&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;twitter:description&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:title&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;og:title&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">:</span> <span class=\"sa\">f</span><span class=\"s1\">&#39;</span><span
    class=\"si\">{</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">]</span><span class=\"si\">}</span><span class=\"s1\"> | </span><span
    class=\"si\">{</span><span class=\"n\">site_name</span><span class=\"si\">}</span><span
    class=\"s1\">&#39;</span><span class=\"p\">[:</span><span class=\"mi\">60</span><span
    class=\"p\">],</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;name&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;twitter:title&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;property&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;twitter:title&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">:</span> <span class=\"sa\">f</span><span class=\"s1\">&#39;</span><span
    class=\"si\">{</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">]</span><span class=\"si\">}</span><span class=\"s1\"> | </span><span
    class=\"si\">{</span><span class=\"n\">site_name</span><span class=\"si\">}</span><span
    class=\"s1\">&#39;</span><span class=\"p\">[:</span><span class=\"mi\">60</span><span
    class=\"p\">],</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;name&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;og:image&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:image&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span> <span class=\"sa\">f</span><span
    class=\"s1\">&#39;</span><span class=\"si\">{</span><span class=\"n\">images_url</span><span
    class=\"si\">}</span><span class=\"s1\">/</span><span class=\"si\">{</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s1\">-og.png&#39;</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;twitter:image&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;twitter:image&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"sa\">f</span><span class=\"s1\">&#39;</span><span class=\"si\">{</span><span
    class=\"n\">images_url</span><span class=\"si\">}</span><span class=\"s1\">/</span><span
    class=\"si\">{</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">]</span><span class=\"si\">}</span><span class=\"s1\">-og.png&#39;</span><span
    class=\"p\">,</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;name&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;og:image:width&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;property&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;og:image:width&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;1600&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:image:width&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:image:width&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;900&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"p\">},</span>\n                <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;twitter:card&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;property&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;twitter:card&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">twitter_card</span><span class=\"p\">,</span>\n                <span
    class=\"p\">},</span>\n                <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;og:site_name&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;property&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;og:site_name&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">site_name</span><span class=\"p\">,</span>\n                <span
    class=\"p\">},</span>\n                <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;twitter:creator&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;property&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;twitter:creator&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">twitter_creator</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;title&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">],</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;name&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;generator&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;generator&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;markata </span><span class=\"si\">{</span><span class=\"n\">__version__</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n            <span class=\"p\">]</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">seo</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='_add_seo_tags' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_add_seo_tags <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_add_seo_tags
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
    <span class=\"nf\">_add_seo_tags</span><span class=\"p\">(</span><span class=\"n\">seo</span><span
    class=\"p\">:</span> <span class=\"n\">List</span><span class=\"p\">,</span> <span
    class=\"n\">article</span><span class=\"p\">:</span> <span class=\"s2\">&quot;frontmatter.Post&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">soup</span><span class=\"p\">:</span> <span
    class=\"n\">BeautifulSoup</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">meta</span> <span class=\"ow\">in</span> <span class=\"n\">seo</span><span
    class=\"p\">:</span>\n                <span class=\"n\">soup</span><span class=\"o\">.</span><span
    class=\"n\">head</span><span class=\"o\">.</span><span class=\"n\">append</span><span
    class=\"p\">(</span><span class=\"n\">_create_seo_tag</span><span class=\"p\">(</span><span
    class=\"n\">meta</span><span class=\"p\">,</span> <span class=\"n\">soup</span><span
    class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='_create_seo_tag'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_create_seo_tag
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_create_seo_tag <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">_create_seo_tag</span><span class=\"p\">(</span><span class=\"n\">meta</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>
    <span class=\"n\">soup</span><span class=\"p\">:</span> <span class=\"n\">BeautifulSoup</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"s2\">&quot;Tag&quot;</span><span
    class=\"p\">:</span>\n            <span class=\"n\">tag</span> <span class=\"o\">=</span>
    <span class=\"n\">soup</span><span class=\"o\">.</span><span class=\"n\">new_tag</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;meta&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">for</span> <span class=\"n\">k</span> <span class=\"ow\">in</span>
    <span class=\"n\">meta</span><span class=\"p\">:</span>\n                <span
    class=\"n\">tag</span><span class=\"o\">.</span><span class=\"n\">attrs</span><span
    class=\"p\">[</span><span class=\"n\">k</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"n\">meta</span><span class=\"p\">[</span><span
    class=\"n\">k</span><span class=\"p\">]</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">tag</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='_get_or_warn'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_get_or_warn <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_get_or_warn <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">_get_or_warn</span><span class=\"p\">(</span><span class=\"n\">config</span><span
    class=\"p\">:</span> <span class=\"n\">Dict</span><span class=\"p\">,</span> <span
    class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">key</span> <span class=\"ow\">not</span> <span class=\"ow\">in</span>
    <span class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
    class=\"p\">():</span>\n                <span class=\"n\">logger</span><span class=\"o\">.</span><span
    class=\"n\">warning</span><span class=\"p\">(</span>\n                    <span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">key</span><span class=\"si\">}</span><span class=\"s2\"> is missing
    from markata.toml config, using default value </span><span class=\"si\">{</span><span
    class=\"n\">default</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"p\">)</span>\n            <span
    class=\"k\">return</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    \           <span class=\"n\">url</span> <span class=\"o\">=</span> <span class=\"n\">_get_or_warn</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">,</span> <span class=\"s2\">&quot;url&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">images_url</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;images_url&quot;</span><span class=\"p\">,</span> <span class=\"n\">url</span><span
    class=\"p\">)</span>\n            <span class=\"n\">site_name</span> <span class=\"o\">=</span>
    <span class=\"n\">_get_or_warn</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;site_name&quot;</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">author_name</span> <span class=\"o\">=</span> <span class=\"n\">_get_or_warn</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">,</span> <span class=\"s2\">&quot;author_name&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">author_email</span> <span class=\"o\">=</span> <span
    class=\"n\">_get_or_warn</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;author_email&quot;</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">twitter_creator</span> <span class=\"o\">=</span> <span class=\"n\">_get_or_warn</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">,</span> <span class=\"s2\">&quot;twitter_creator&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">twitter_card</span> <span class=\"o\">=</span> <span
    class=\"n\">_get_or_warn</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;twitter_card&quot;</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;summary_large_image&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">config_seo</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;seo&quot;</span><span class=\"p\">,</span> <span class=\"p\">{})</span>\n
    \           <span class=\"n\">should_prettify</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;prettify_html&quot;</span><span class=\"p\">,</span> <span
    class=\"kc\">False</span><span class=\"p\">)</span>\n\n            <span class=\"k\">with</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span>
    <span class=\"k\">as</span> <span class=\"n\">cache</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">article</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">iter_articles</span><span class=\"p\">(</span><span class=\"s2\">&quot;add
    seo tags from seo.py&quot;</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span>\n
    \                       <span class=\"s2\">&quot;seo&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"s2\">&quot;render&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">html</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">site_name</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">url</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">],</span>\n
    \                       <span class=\"n\">twitter_card</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">],</span>\n                        <span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">config_seo</span><span class=\"p\">),</span>\n
    \                   <span class=\"p\">)</span>\n\n                    <span class=\"n\">html_from_cache</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">)</span>\n\n
    \                   <span class=\"k\">if</span> <span class=\"n\">html_from_cache</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">soup</span> <span class=\"o\">=</span>
    <span class=\"n\">BeautifulSoup</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">html</span><span class=\"p\">,</span> <span
    class=\"n\">features</span><span class=\"o\">=</span><span class=\"s2\">&quot;lxml&quot;</span><span
    class=\"p\">)</span>\n                        <span class=\"n\">seo</span> <span
    class=\"o\">=</span> <span class=\"n\">_create_seo</span><span class=\"p\">(</span>\n
    \                           <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">soup</span><span class=\"o\">=</span><span class=\"n\">soup</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">article</span><span
    class=\"o\">=</span><span class=\"n\">article</span><span class=\"p\">,</span>\n
    \                           <span class=\"n\">site_name</span><span class=\"o\">=</span><span
    class=\"n\">site_name</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">author_name</span><span class=\"o\">=</span><span class=\"n\">author_name</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">author_email</span><span
    class=\"o\">=</span><span class=\"n\">author_email</span><span class=\"p\">,</span>\n
    \                           <span class=\"n\">twitter_card</span><span class=\"o\">=</span><span
    class=\"n\">twitter_card</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">twitter_creator</span><span class=\"o\">=</span><span class=\"n\">twitter_creator</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">config_seo</span><span
    class=\"o\">=</span><span class=\"n\">config_seo</span><span class=\"p\">,</span>\n
    \                           <span class=\"n\">images_url</span><span class=\"o\">=</span><span
    class=\"n\">images_url</span><span class=\"p\">,</span>\n                        <span
    class=\"p\">)</span>\n                        <span class=\"n\">_add_seo_tags</span><span
    class=\"p\">(</span><span class=\"n\">seo</span><span class=\"p\">,</span> <span
    class=\"n\">article</span><span class=\"p\">,</span> <span class=\"n\">soup</span><span
    class=\"p\">)</span>\n                        <span class=\"n\">canonical_link</span>
    <span class=\"o\">=</span> <span class=\"n\">soup</span><span class=\"o\">.</span><span
    class=\"n\">new_tag</span><span class=\"p\">(</span><span class=\"s2\">&quot;link&quot;</span><span
    class=\"p\">)</span>\n                        <span class=\"n\">canonical_link</span><span
    class=\"o\">.</span><span class=\"n\">attrs</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;rel&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;canonical&quot;</span>\n                        <span
    class=\"k\">if</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;index&quot;</span><span
    class=\"p\">:</span>\n                            <span class=\"n\">canonical_link</span><span
    class=\"o\">.</span><span class=\"n\">attrs</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;href&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">url</span><span class=\"si\">}</span><span class=\"s2\">/&quot;</span>\n
    \                       <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                           <span class=\"n\">canonical_link</span><span class=\"o\">.</span><span
    class=\"n\">attrs</span><span class=\"p\">[</span><span class=\"s2\">&quot;href&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"sa\">f</span><span
    class=\"s1\">&#39;</span><span class=\"si\">{</span><span class=\"n\">url</span><span
    class=\"si\">}</span><span class=\"s1\">/</span><span class=\"si\">{</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s1\">/&#39;</span>\n                        <span
    class=\"n\">soup</span><span class=\"o\">.</span><span class=\"n\">head</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">canonical_link</span><span class=\"p\">)</span>\n\n                        <span
    class=\"n\">meta_url</span> <span class=\"o\">=</span> <span class=\"n\">soup</span><span
    class=\"o\">.</span><span class=\"n\">new_tag</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;meta&quot;</span><span class=\"p\">)</span>\n                        <span
    class=\"n\">meta_url</span><span class=\"o\">.</span><span class=\"n\">attrs</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;name&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;og:url&quot;</span>\n                        <span
    class=\"n\">meta_url</span><span class=\"o\">.</span><span class=\"n\">attrs</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;property&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;og:url&quot;</span>\n                        <span
    class=\"k\">if</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;index&quot;</span><span
    class=\"p\">:</span>\n                            <span class=\"n\">meta_url</span><span
    class=\"o\">.</span><span class=\"n\">attrs</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;content&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">url</span><span class=\"si\">}</span><span class=\"s2\">/&quot;</span>\n
    \                       <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                           <span class=\"n\">meta_url</span><span class=\"o\">.</span><span
    class=\"n\">attrs</span><span class=\"p\">[</span><span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"sa\">f</span><span
    class=\"s1\">&#39;</span><span class=\"si\">{</span><span class=\"n\">url</span><span
    class=\"si\">}</span><span class=\"s1\">/</span><span class=\"si\">{</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s1\">/&#39;</span>\n                        <span
    class=\"n\">soup</span><span class=\"o\">.</span><span class=\"n\">head</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">meta_url</span><span class=\"p\">)</span>\n\n                        <span
    class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">soup</span><span
    class=\"o\">.</span><span class=\"n\">prettify</span><span class=\"p\">()</span>
    <span class=\"k\">if</span> <span class=\"n\">should_prettify</span> <span class=\"k\">else</span>
    <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">soup</span><span
    class=\"p\">)</span>\n                        <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">html</span><span
    class=\"p\">,</span> <span class=\"n\">expire</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">default_cache_expire</span><span class=\"p\">)</span>\n\n
    \                   <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">html</span> <span class=\"o\">=</span>
    <span class=\"n\">html_from_cache</span>\n                    <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">html</span> <span class=\"o\">=</span> <span
    class=\"n\">html</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Seo.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"manifest plugin ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Seo.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"manifest plugin ! ???+ source  ! ???+
    source  ! ???+ source  ! ???+ source  ! ???+ source \" />\n <link href=\"/favicon.ico\"
    rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\"
    />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Seo.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Seo.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>manifest
    plugin</p>\n<p>!! function <h2 id='_create_seo' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_create_seo <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_create_seo
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
    <span class=\"nf\">_create_seo</span><span class=\"p\">(</span>\n            <span
    class=\"n\">markata</span><span class=\"p\">:</span> <span class=\"n\">Markata</span><span
    class=\"p\">,</span>\n            <span class=\"n\">soup</span><span class=\"p\">:</span>
    <span class=\"n\">BeautifulSoup</span><span class=\"p\">,</span>\n            <span
    class=\"n\">article</span><span class=\"p\">:</span> <span class=\"s2\">&quot;frontmatter.Post&quot;</span><span
    class=\"p\">,</span>\n            <span class=\"n\">site_name</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">,</span>\n            <span class=\"n\">author_name</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">author_email</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">,</span>\n            <span class=\"n\">twitter_card</span><span
    class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">twitter_creator</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">,</span>\n            <span class=\"n\">config_seo</span><span
    class=\"p\">:</span> <span class=\"n\">Dict</span><span class=\"p\">,</span>\n
    \           <span class=\"n\">images_url</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">,</span>\n        <span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"n\">List</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">if</span> <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;description&quot;</span><span class=\"p\">]</span> <span class=\"o\">==</span>
    <span class=\"s2\">&quot;&quot;</span> <span class=\"ow\">or</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"s2\">&quot; &quot;</span><span
    class=\"o\">.</span><span class=\"n\">join</span><span class=\"p\">(</span>\n
    \                       <span class=\"p\">[</span><span class=\"n\">p</span><span
    class=\"o\">.</span><span class=\"n\">text</span> <span class=\"k\">for</span>
    <span class=\"n\">p</span> <span class=\"ow\">in</span> <span class=\"n\">soup</span><span
    class=\"o\">.</span><span class=\"n\">find</span><span class=\"p\">(</span><span
    class=\"nb\">id</span><span class=\"o\">=</span><span class=\"s2\">&quot;post-body&quot;</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">find_all</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;p&quot;</span><span class=\"p\">)],</span>\n
    \                   <span class=\"p\">)</span><span class=\"o\">.</span><span
    class=\"n\">strip</span><span class=\"p\">()[:</span><span class=\"mi\">120</span><span
    class=\"p\">]</span>\n                <span class=\"k\">except</span> <span class=\"ne\">AttributeError</span><span
    class=\"p\">:</span>\n                    <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">metadata</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;description&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;&quot;</span>\n\n            <span class=\"n\">seo</span>
    <span class=\"o\">=</span> <span class=\"p\">[</span>\n                <span class=\"o\">*</span><span
    class=\"n\">config_seo</span><span class=\"p\">,</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;name&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;og:author&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:author&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span> <span class=\"n\">author_name</span><span
    class=\"p\">,</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;name&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;og:author_email&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;property&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;og:author_email&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">author_email</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:type&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;og:type&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;website&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;description&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;description&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:description&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:description&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;twitter:description&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;twitter:description&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;description&quot;</span><span class=\"p\">],</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:title&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;og:title&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">:</span> <span class=\"sa\">f</span><span class=\"s1\">&#39;</span><span
    class=\"si\">{</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">]</span><span class=\"si\">}</span><span class=\"s1\"> | </span><span
    class=\"si\">{</span><span class=\"n\">site_name</span><span class=\"si\">}</span><span
    class=\"s1\">&#39;</span><span class=\"p\">[:</span><span class=\"mi\">60</span><span
    class=\"p\">],</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;name&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;twitter:title&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;property&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;twitter:title&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">:</span> <span class=\"sa\">f</span><span class=\"s1\">&#39;</span><span
    class=\"si\">{</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">]</span><span class=\"si\">}</span><span class=\"s1\"> | </span><span
    class=\"si\">{</span><span class=\"n\">site_name</span><span class=\"si\">}</span><span
    class=\"s1\">&#39;</span><span class=\"p\">[:</span><span class=\"mi\">60</span><span
    class=\"p\">],</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;name&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;og:image&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:image&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span> <span class=\"sa\">f</span><span
    class=\"s1\">&#39;</span><span class=\"si\">{</span><span class=\"n\">images_url</span><span
    class=\"si\">}</span><span class=\"s1\">/</span><span class=\"si\">{</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s1\">-og.png&#39;</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;twitter:image&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;twitter:image&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"sa\">f</span><span class=\"s1\">&#39;</span><span class=\"si\">{</span><span
    class=\"n\">images_url</span><span class=\"si\">}</span><span class=\"s1\">/</span><span
    class=\"si\">{</span><span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">]</span><span class=\"si\">}</span><span class=\"s1\">-og.png&#39;</span><span
    class=\"p\">,</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;name&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;og:image:width&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;property&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;og:image:width&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;1600&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:image:width&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;og:image:width&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;900&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"p\">},</span>\n                <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;twitter:card&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;property&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;twitter:card&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">twitter_card</span><span class=\"p\">,</span>\n                <span
    class=\"p\">},</span>\n                <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;og:site_name&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;property&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;og:site_name&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span>
    <span class=\"n\">site_name</span><span class=\"p\">,</span>\n                <span
    class=\"p\">},</span>\n                <span class=\"p\">{</span>\n                    <span
    class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;twitter:creator&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;property&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;twitter:creator&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">twitter_creator</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n                <span class=\"p\">{</span>\n
    \                   <span class=\"s2\">&quot;name&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;title&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span> <span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">:</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">],</span>\n                <span class=\"p\">},</span>\n                <span
    class=\"p\">{</span>\n                    <span class=\"s2\">&quot;name&quot;</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;generator&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;property&quot;</span><span class=\"p\">:</span>
    <span class=\"s2\">&quot;generator&quot;</span><span class=\"p\">,</span>\n                    <span
    class=\"s2\">&quot;content&quot;</span><span class=\"p\">:</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;markata </span><span class=\"si\">{</span><span class=\"n\">__version__</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">},</span>\n            <span class=\"p\">]</span>\n
    \           <span class=\"k\">return</span> <span class=\"n\">seo</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='_add_seo_tags' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>_add_seo_tags <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">_add_seo_tags
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
    <span class=\"nf\">_add_seo_tags</span><span class=\"p\">(</span><span class=\"n\">seo</span><span
    class=\"p\">:</span> <span class=\"n\">List</span><span class=\"p\">,</span> <span
    class=\"n\">article</span><span class=\"p\">:</span> <span class=\"s2\">&quot;frontmatter.Post&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">soup</span><span class=\"p\">:</span> <span
    class=\"n\">BeautifulSoup</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"k\">for</span>
    <span class=\"n\">meta</span> <span class=\"ow\">in</span> <span class=\"n\">seo</span><span
    class=\"p\">:</span>\n                <span class=\"n\">soup</span><span class=\"o\">.</span><span
    class=\"n\">head</span><span class=\"o\">.</span><span class=\"n\">append</span><span
    class=\"p\">(</span><span class=\"n\">_create_seo_tag</span><span class=\"p\">(</span><span
    class=\"n\">meta</span><span class=\"p\">,</span> <span class=\"n\">soup</span><span
    class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='_create_seo_tag'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_create_seo_tag
    <em class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_create_seo_tag <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">_create_seo_tag</span><span class=\"p\">(</span><span class=\"n\">meta</span><span
    class=\"p\">:</span> <span class=\"nb\">dict</span><span class=\"p\">,</span>
    <span class=\"n\">soup</span><span class=\"p\">:</span> <span class=\"n\">BeautifulSoup</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"s2\">&quot;Tag&quot;</span><span
    class=\"p\">:</span>\n            <span class=\"n\">tag</span> <span class=\"o\">=</span>
    <span class=\"n\">soup</span><span class=\"o\">.</span><span class=\"n\">new_tag</span><span
    class=\"p\">(</span><span class=\"s2\">&quot;meta&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">for</span> <span class=\"n\">k</span> <span class=\"ow\">in</span>
    <span class=\"n\">meta</span><span class=\"p\">:</span>\n                <span
    class=\"n\">tag</span><span class=\"o\">.</span><span class=\"n\">attrs</span><span
    class=\"p\">[</span><span class=\"n\">k</span><span class=\"p\">]</span> <span
    class=\"o\">=</span> <span class=\"n\">meta</span><span class=\"p\">[</span><span
    class=\"n\">k</span><span class=\"p\">]</span>\n            <span class=\"k\">return</span>
    <span class=\"n\">tag</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='_get_or_warn'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_get_or_warn <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">_get_or_warn <em class='small'>source</em></p>\n</div>\n<pre
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
    <span class=\"nf\">_get_or_warn</span><span class=\"p\">(</span><span class=\"n\">config</span><span
    class=\"p\">:</span> <span class=\"n\">Dict</span><span class=\"p\">,</span> <span
    class=\"n\">key</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">:</span>
    <span class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"n\">Any</span><span class=\"p\">:</span>\n            <span class=\"k\">if</span>
    <span class=\"n\">key</span> <span class=\"ow\">not</span> <span class=\"ow\">in</span>
    <span class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">keys</span><span
    class=\"p\">():</span>\n                <span class=\"n\">logger</span><span class=\"o\">.</span><span
    class=\"n\">warning</span><span class=\"p\">(</span>\n                    <span
    class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">key</span><span class=\"si\">}</span><span class=\"s2\"> is missing
    from markata.toml config, using default value </span><span class=\"si\">{</span><span
    class=\"n\">default</span><span class=\"si\">}</span><span class=\"s2\">&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"p\">)</span>\n            <span
    class=\"k\">return</span> <span class=\"n\">config</span><span class=\"o\">.</span><span
    class=\"n\">get</span><span class=\"p\">(</span><span class=\"n\">key</span><span
    class=\"p\">,</span> <span class=\"n\">default</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
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
    \           <span class=\"n\">url</span> <span class=\"o\">=</span> <span class=\"n\">_get_or_warn</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">,</span> <span class=\"s2\">&quot;url&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">images_url</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;images_url&quot;</span><span class=\"p\">,</span> <span class=\"n\">url</span><span
    class=\"p\">)</span>\n            <span class=\"n\">site_name</span> <span class=\"o\">=</span>
    <span class=\"n\">_get_or_warn</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;site_name&quot;</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">author_name</span> <span class=\"o\">=</span> <span class=\"n\">_get_or_warn</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">,</span> <span class=\"s2\">&quot;author_name&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">author_email</span> <span class=\"o\">=</span> <span
    class=\"n\">_get_or_warn</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;author_email&quot;</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n            <span
    class=\"n\">twitter_creator</span> <span class=\"o\">=</span> <span class=\"n\">_get_or_warn</span><span
    class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">,</span> <span class=\"s2\">&quot;twitter_creator&quot;</span><span
    class=\"p\">,</span> <span class=\"s2\">&quot;&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">twitter_card</span> <span class=\"o\">=</span> <span
    class=\"n\">_get_or_warn</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">,</span>
    <span class=\"s2\">&quot;twitter_card&quot;</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;summary_large_image&quot;</span><span class=\"p\">)</span>\n
    \           <span class=\"n\">config_seo</span> <span class=\"o\">=</span> <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;seo&quot;</span><span class=\"p\">,</span> <span class=\"p\">{})</span>\n
    \           <span class=\"n\">should_prettify</span> <span class=\"o\">=</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;prettify_html&quot;</span><span class=\"p\">,</span> <span
    class=\"kc\">False</span><span class=\"p\">)</span>\n\n            <span class=\"k\">with</span>
    <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">cache</span>
    <span class=\"k\">as</span> <span class=\"n\">cache</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">for</span> <span class=\"n\">article</span>
    <span class=\"ow\">in</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">iter_articles</span><span class=\"p\">(</span><span class=\"s2\">&quot;add
    seo tags from seo.py&quot;</span><span class=\"p\">):</span>\n                    <span
    class=\"n\">key</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">make_hash</span><span class=\"p\">(</span>\n
    \                       <span class=\"s2\">&quot;seo&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"s2\">&quot;render&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">html</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">site_name</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">url</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">],</span>\n
    \                       <span class=\"n\">twitter_card</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;title&quot;</span><span
    class=\"p\">],</span>\n                        <span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">config_seo</span><span class=\"p\">),</span>\n
    \                   <span class=\"p\">)</span>\n\n                    <span class=\"n\">html_from_cache</span>
    <span class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
    class=\"n\">precache</span><span class=\"o\">.</span><span class=\"n\">get</span><span
    class=\"p\">(</span><span class=\"n\">key</span><span class=\"p\">)</span>\n\n
    \                   <span class=\"k\">if</span> <span class=\"n\">html_from_cache</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">soup</span> <span class=\"o\">=</span>
    <span class=\"n\">BeautifulSoup</span><span class=\"p\">(</span><span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">html</span><span class=\"p\">,</span> <span
    class=\"n\">features</span><span class=\"o\">=</span><span class=\"s2\">&quot;lxml&quot;</span><span
    class=\"p\">)</span>\n                        <span class=\"n\">seo</span> <span
    class=\"o\">=</span> <span class=\"n\">_create_seo</span><span class=\"p\">(</span>\n
    \                           <span class=\"n\">markata</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">soup</span><span class=\"o\">=</span><span class=\"n\">soup</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">article</span><span
    class=\"o\">=</span><span class=\"n\">article</span><span class=\"p\">,</span>\n
    \                           <span class=\"n\">site_name</span><span class=\"o\">=</span><span
    class=\"n\">site_name</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">author_name</span><span class=\"o\">=</span><span class=\"n\">author_name</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">author_email</span><span
    class=\"o\">=</span><span class=\"n\">author_email</span><span class=\"p\">,</span>\n
    \                           <span class=\"n\">twitter_card</span><span class=\"o\">=</span><span
    class=\"n\">twitter_card</span><span class=\"p\">,</span>\n                            <span
    class=\"n\">twitter_creator</span><span class=\"o\">=</span><span class=\"n\">twitter_creator</span><span
    class=\"p\">,</span>\n                            <span class=\"n\">config_seo</span><span
    class=\"o\">=</span><span class=\"n\">config_seo</span><span class=\"p\">,</span>\n
    \                           <span class=\"n\">images_url</span><span class=\"o\">=</span><span
    class=\"n\">images_url</span><span class=\"p\">,</span>\n                        <span
    class=\"p\">)</span>\n                        <span class=\"n\">_add_seo_tags</span><span
    class=\"p\">(</span><span class=\"n\">seo</span><span class=\"p\">,</span> <span
    class=\"n\">article</span><span class=\"p\">,</span> <span class=\"n\">soup</span><span
    class=\"p\">)</span>\n                        <span class=\"n\">canonical_link</span>
    <span class=\"o\">=</span> <span class=\"n\">soup</span><span class=\"o\">.</span><span
    class=\"n\">new_tag</span><span class=\"p\">(</span><span class=\"s2\">&quot;link&quot;</span><span
    class=\"p\">)</span>\n                        <span class=\"n\">canonical_link</span><span
    class=\"o\">.</span><span class=\"n\">attrs</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;rel&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"s2\">&quot;canonical&quot;</span>\n                        <span
    class=\"k\">if</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;index&quot;</span><span
    class=\"p\">:</span>\n                            <span class=\"n\">canonical_link</span><span
    class=\"o\">.</span><span class=\"n\">attrs</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;href&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">url</span><span class=\"si\">}</span><span class=\"s2\">/&quot;</span>\n
    \                       <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                           <span class=\"n\">canonical_link</span><span class=\"o\">.</span><span
    class=\"n\">attrs</span><span class=\"p\">[</span><span class=\"s2\">&quot;href&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"sa\">f</span><span
    class=\"s1\">&#39;</span><span class=\"si\">{</span><span class=\"n\">url</span><span
    class=\"si\">}</span><span class=\"s1\">/</span><span class=\"si\">{</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s1\">/&#39;</span>\n                        <span
    class=\"n\">soup</span><span class=\"o\">.</span><span class=\"n\">head</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">canonical_link</span><span class=\"p\">)</span>\n\n                        <span
    class=\"n\">meta_url</span> <span class=\"o\">=</span> <span class=\"n\">soup</span><span
    class=\"o\">.</span><span class=\"n\">new_tag</span><span class=\"p\">(</span><span
    class=\"s2\">&quot;meta&quot;</span><span class=\"p\">)</span>\n                        <span
    class=\"n\">meta_url</span><span class=\"o\">.</span><span class=\"n\">attrs</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;name&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;og:url&quot;</span>\n                        <span
    class=\"n\">meta_url</span><span class=\"o\">.</span><span class=\"n\">attrs</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;property&quot;</span><span class=\"p\">]</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;og:url&quot;</span>\n                        <span
    class=\"k\">if</span> <span class=\"n\">article</span><span class=\"o\">.</span><span
    class=\"n\">metadata</span><span class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">==</span> <span class=\"s2\">&quot;index&quot;</span><span
    class=\"p\">:</span>\n                            <span class=\"n\">meta_url</span><span
    class=\"o\">.</span><span class=\"n\">attrs</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;content&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"sa\">f</span><span class=\"s2\">&quot;</span><span class=\"si\">{</span><span
    class=\"n\">url</span><span class=\"si\">}</span><span class=\"s2\">/&quot;</span>\n
    \                       <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                           <span class=\"n\">meta_url</span><span class=\"o\">.</span><span
    class=\"n\">attrs</span><span class=\"p\">[</span><span class=\"s2\">&quot;content&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"sa\">f</span><span
    class=\"s1\">&#39;</span><span class=\"si\">{</span><span class=\"n\">url</span><span
    class=\"si\">}</span><span class=\"s1\">/</span><span class=\"si\">{</span><span
    class=\"n\">article</span><span class=\"o\">.</span><span class=\"n\">metadata</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;slug&quot;</span><span class=\"p\">]</span><span
    class=\"si\">}</span><span class=\"s1\">/&#39;</span>\n                        <span
    class=\"n\">soup</span><span class=\"o\">.</span><span class=\"n\">head</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">meta_url</span><span class=\"p\">)</span>\n\n                        <span
    class=\"n\">html</span> <span class=\"o\">=</span> <span class=\"n\">soup</span><span
    class=\"o\">.</span><span class=\"n\">prettify</span><span class=\"p\">()</span>
    <span class=\"k\">if</span> <span class=\"n\">should_prettify</span> <span class=\"k\">else</span>
    <span class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">soup</span><span
    class=\"p\">)</span>\n                        <span class=\"n\">cache</span><span
    class=\"o\">.</span><span class=\"n\">add</span><span class=\"p\">(</span><span
    class=\"n\">key</span><span class=\"p\">,</span> <span class=\"n\">html</span><span
    class=\"p\">,</span> <span class=\"n\">expire</span><span class=\"o\">=</span><span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">default_cache_expire</span><span class=\"p\">)</span>\n\n
    \                   <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                       <span class=\"n\">html</span> <span class=\"o\">=</span>
    <span class=\"n\">html_from_cache</span>\n                    <span class=\"n\">article</span><span
    class=\"o\">.</span><span class=\"n\">html</span> <span class=\"o\">=</span> <span
    class=\"n\">html</span>\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: ''
published: true
slug: markata/plugins/seo
title: Seo.Py


---

manifest plugin


!! function <h2 id='_create_seo' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_create_seo <em class='small'>function</em></h2>

???+ source "_create_seo <em class='small'>source</em>"

```python

        def _create_seo(
            markata: Markata,
            soup: BeautifulSoup,
            article: "frontmatter.Post",
            site_name: str,
            author_name: str,
            author_email: str,
            twitter_card: str,
            twitter_creator: str,
            config_seo: Dict,
            images_url: str,
        ) -> List:
            if article.metadata["description"] == "" or None:
                try:
                    article.metadata["description"] = " ".join(
                        [p.text for p in soup.find(id="post-body").find_all("p")],
                    ).strip()[:120]
                except AttributeError:
                    article.metadata["description"] = ""

            seo = [
                *config_seo,
                {
                    "name": "og:author",
                    "property": "og:author",
                    "content": author_name,
                },
                {
                    "name": "og:author_email",
                    "property": "og:author_email",
                    "content": author_email,
                },
                {
                    "name": "og:type",
                    "property": "og:type",
                    "content": "website",
                },
                {
                    "name": "description",
                    "property": "description",
                    "content": article.metadata["description"],
                },
                {
                    "name": "og:description",
                    "property": "og:description",
                    "content": article.metadata["description"],
                },
                {
                    "name": "twitter:description",
                    "property": "twitter:description",
                    "content": article.metadata["description"],
                },
                {
                    "name": "og:title",
                    "property": "og:title",
                    "content": f'{article.metadata["title"]} | {site_name}'[:60],
                },
                {
                    "name": "twitter:title",
                    "property": "twitter:title",
                    "content": f'{article.metadata["title"]} | {site_name}'[:60],
                },
                {
                    "name": "og:image",
                    "property": "og:image",
                    "content": f'{images_url}/{article.metadata["slug"]}-og.png',
                },
                {
                    "name": "twitter:image",
                    "property": "twitter:image",
                    "content": f'{images_url}/{article.metadata["slug"]}-og.png',
                },
                {
                    "name": "og:image:width",
                    "property": "og:image:width",
                    "content": "1600",
                },
                {
                    "name": "og:image:width",
                    "property": "og:image:width",
                    "content": "900",
                },
                {
                    "name": "twitter:card",
                    "property": "twitter:card",
                    "content": twitter_card,
                },
                {
                    "name": "og:site_name",
                    "property": "og:site_name",
                    "content": site_name,
                },
                {
                    "name": "twitter:creator",
                    "property": "twitter:creator",
                    "content": twitter_creator,
                },
                {
                    "name": "title",
                    "property": "title",
                    "content": article.metadata["title"],
                },
                {
                    "name": "generator",
                    "property": "generator",
                    "content": f"markata {__version__}",
                },
            ]
            return seo
```


!! function <h2 id='_add_seo_tags' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_add_seo_tags <em class='small'>function</em></h2>

???+ source "_add_seo_tags <em class='small'>source</em>"

```python

        def _add_seo_tags(seo: List, article: "frontmatter.Post", soup: BeautifulSoup) -> None:
            for meta in seo:
                soup.head.append(_create_seo_tag(meta, soup))
```


!! function <h2 id='_create_seo_tag' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_create_seo_tag <em class='small'>function</em></h2>

???+ source "_create_seo_tag <em class='small'>source</em>"

```python

        def _create_seo_tag(meta: dict, soup: BeautifulSoup) -> "Tag":
            tag = soup.new_tag("meta")
            for k in meta:
                tag.attrs[k] = meta[k]
            return tag
```


!! function <h2 id='_get_or_warn' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_get_or_warn <em class='small'>function</em></h2>

???+ source "_get_or_warn <em class='small'>source</em>"

```python

        def _get_or_warn(config: Dict, key: str, default: str) -> Any:
            if key not in config.keys():
                logger.warning(
                    f"{key} is missing from markata.toml config, using default value {default}",
                )
            return config.get(key, default)
```


!! function <h2 id='render' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>render <em class='small'>function</em></h2>

???+ source "render <em class='small'>source</em>"

```python

        def render(markata: Markata) -> None:
            url = _get_or_warn(markata.config, "url", "")
            images_url = markata.config.get("images_url", url)
            site_name = _get_or_warn(markata.config, "site_name", "")
            author_name = _get_or_warn(markata.config, "author_name", "")
            author_email = _get_or_warn(markata.config, "author_email", "")
            twitter_creator = _get_or_warn(markata.config, "twitter_creator", "")
            twitter_card = _get_or_warn(markata.config, "twitter_card", "summary_large_image")
            config_seo = markata.config.get("seo", {})
            should_prettify = markata.config.get("prettify_html", False)

            with markata.cache as cache:
                for article in markata.iter_articles("add seo tags from seo.py"):
                    key = markata.make_hash(
                        "seo",
                        "render",
                        article.html,
                        site_name,
                        url,
                        article.metadata["slug"],
                        twitter_card,
                        article.metadata["title"],
                        str(config_seo),
                    )

                    html_from_cache = markata.precache.get(key)

                    if html_from_cache is None:
                        soup = BeautifulSoup(article.html, features="lxml")
                        seo = _create_seo(
                            markata=markata,
                            soup=soup,
                            article=article,
                            site_name=site_name,
                            author_name=author_name,
                            author_email=author_email,
                            twitter_card=twitter_card,
                            twitter_creator=twitter_creator,
                            config_seo=config_seo,
                            images_url=images_url,
                        )
                        _add_seo_tags(seo, article, soup)
                        canonical_link = soup.new_tag("link")
                        canonical_link.attrs["rel"] = "canonical"
                        if article.metadata["slug"] == "index":
                            canonical_link.attrs["href"] = f"{url}/"
                        else:
                            canonical_link.attrs["href"] = f'{url}/{article.metadata["slug"]}/'
                        soup.head.append(canonical_link)

                        meta_url = soup.new_tag("meta")
                        meta_url.attrs["name"] = "og:url"
                        meta_url.attrs["property"] = "og:url"
                        if article.metadata["slug"] == "index":
                            meta_url.attrs["content"] = f"{url}/"
                        else:
                            meta_url.attrs["content"] = f'{url}/{article.metadata["slug"]}/'
                        soup.head.append(meta_url)

                        html = soup.prettify() if should_prettify else str(soup)
                        cache.add(key, html, expire=markata.config.default_cache_expire)

                    else:
                        html = html_from_cache
                    article.html = html
```

