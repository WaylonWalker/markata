---
content: "\nThis page will be rendered multiple times.  This is the normal page, check\nout
  [og](./og). It uses the og.html.  This is achieved through the use of the\ntemplate
  option in the frontmatter.\n\n``` md\n---\ntitle: Using Multiple Templates\ndescription:
  Guide to using multiple templates in markata\njinja: false\ntemplate:\n  index:
  post.html\n  title: og.html\n\n\n---\n```\n"
date: 0001-01-01
description: Guide to using multiple templates in markata
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Using Multiple Templates</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Guide to using multiple templates in
    markata\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link
    rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\"
    />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n<title>Using Multiple Templates</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Guide to using multiple templates in markata\"
    />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\"
    href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script
    src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    </head>\n    <body>\n<div class='container flex flex-row min-h-screen'>\n
    \   <div>\n    </div>\n    <div class='flex-grow px-8 mx-auto min-h-screen'>\n<header
    class='flex justify-center items-center p-8'>\n\n    <nav class='flex justify-center
    items-center my-8'>\n        <a\n            href='/'>markata</a>\n        <a\n
    \           href='https://github.com/WaylonWalker/markata'>GitHub</a>\n        <a\n
    \           href='https://markata.dev/docs/'>docs</a>\n        <a\n            href='https://markata.dev/plugins/'>plugins</a>\n
    \   </nav>\n\n    <div>\n        <label id=\"theme-switch\" class=\"theme-switch\"
    for=\"checkbox-theme\" title=\"light/dark mode toggle\">\n            <input type=\"checkbox\"
    id=\"checkbox-theme\" />\n            <div class=\"slider round\"></div>\n        </label>\n
    \   </div>\n</header><article class='w-full'>\n<section class=\"title\">\n    <h1
    id=\"title\">\n        Using Multiple Templates\n    </h1>\n</section>    <section
    class=\"body\">\n        <p>This page will be rendered multiple times.  This is
    the normal page, check\nout <a href=\"./og\">og</a>. It uses the og.html.  This
    is achieved through the use of the\ntemplate option in the frontmatter.</p>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span>---\ntitle: Using Multiple
    Templates\ndescription: Guide to using multiple templates in markata\njinja: false\ntemplate:\n
    \ index: post.html\n  title: og.html\n\n\n---\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Using Multiple Templates</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Guide to using multiple templates in
    markata\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link
    rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\"
    />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n<title>Using Multiple Templates</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Guide to using multiple templates in markata\"
    />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\"
    href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script
    src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    </head>\n    <body>\n<article style=\"text-align: center;\">\n    <style>\n
    \       section {\n            font-size: 200%;\n        }\n\n\n        .edit
    {\n            display: none;\n        }\n    </style>\n<section class=\"title\">\n
    \   <h1 id=\"title\">\n        Using Multiple Templates\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Using Multiple Templates\n    </h1>\n</section>    <section class=\"body\">\n
    \       <p>This page will be rendered multiple times.  This is the normal page,
    check\nout <a href=\"./og\">og</a>. It uses the og.html.  This is achieved through
    the use of the\ntemplate option in the frontmatter.</p>\n<pre class='wrapper'>\n\n<div
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
    \       \n<div class=\"highlight\"><pre><span></span>---\ntitle: Using Multiple
    Templates\ndescription: Guide to using multiple templates in markata\njinja: false\ntemplate:\n
    \ index: post.html\n  title: og.html\n\n\n---\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: "---\ntitle: Using Multiple Templates\ndescription: Guide to using multiple
    templates in markata\njinja: false\ntemplate:\n  index: post.html\n  title: og.html\n\n\n---\n\nThis
    page will be rendered multiple times.  This is the normal page, check\nout [og](./og).
    It uses the og.html.  This is achieved through the use of the\ntemplate option
    in the frontmatter.\n\n``` md\n---\ntitle: Using Multiple Templates\ndescription:
    Guide to using multiple templates in markata\njinja: false\ntemplate:\n  index:
    post.html\n  title: og.html\n\n\n---\n```\n"
  title: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Using Multiple Templates</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"Guide to using multiple templates in
    markata\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link
    rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\"
    />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    <head>\n<title>Using Multiple Templates</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"Guide to using multiple templates in markata\"
    />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\" />\n\n<link rel=\"stylesheet\"
    href=\"/post.css\" />\n<link rel=\"stylesheet\" href=\"/app.css\" />\n<script
    src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\" content=\"waylon@waylonwalker.com\"
    />\n\n    </head>\n    <body>\n<article style=\"text-align: center;\">\n    <style>\n
    \       section {\n            font-size: 200%;\n        }\n\n\n        .edit
    {\n            display: none;\n        }\n    </style>\n<section class=\"title\">\n
    \   <h1 id=\"title\">\n        Using Multiple Templates\n    </h1>\n</section></article>\n
    \    </body>\n</html>"
published: true
slug: multi-template
title: Using Multiple Templates


---


This page will be rendered multiple times.  This is the normal page, check
out [og](./og). It uses the og.html.  This is achieved through the use of the
template option in the frontmatter.

``` md
---
title: Using Multiple Templates
description: Guide to using multiple templates in markata
jinja: false
template:
  index: post.html
  title: og.html


---
```
