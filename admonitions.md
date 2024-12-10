---
content: "\n???+ note open by default\nyou can open a details tab with '+'\n\n<!--
  prettier-ignore -->\n??? note closed by default\n    you can open a details tab
  with '+'\n\n## all of the admonitions\n\n<!-- prettier-ignore -->\n!!! note\n    a
  note\n\n<!-- prettier-ignore -->\n!!! abstract\n    an abstract\n\n<!-- prettier-ignore
  -->\n!!! info\n\n    admonitions\n\n<!-- prettier-ignore -->\n!!! tip\n\n    You
  should think about using admonitions\n\n<!-- prettier-ignore -->\n!!! success\n\n
  \   Run Successful!\n\n<!-- prettier-ignore -->\n!!! question\n\n    What do you
  think of this?\n\n<!-- prettier-ignore -->\n!!! source\n    Add some source code.\n
  \   ```python\n    print('hello world')\n\n````\n\n<!-- prettier-ignore -->\n!!!
  warning\n    a warning\n\n<!-- prettier-ignore -->\n!!! failure\n    a failure\n\n<!--
  prettier-ignore -->\n!!! danger\n    some danger\n\n<!-- prettier-ignore -->\n!!!
  bug\n    a bug\n\n<!-- prettier-ignore -->\n!!! example\n    an example\n\n    ```
  python\n    print('hello world')\n    ```\n\n<!-- prettier-ignore -->\n!!! quote\n\n
  \   a quote\n\n    > include a nice quote\n\n<!-- Background #282a36 -->\n<!-- Current
  Line #44475a -->\n<!-- Foreground #f8f8f2 -->\n<!-- Comment #6272a4 -->\n<!-- Cyan
  #8be9fd -->\n<!-- Green #50fa7b -->\n<!-- Orange #ffb86c -->\n<!-- Pink #ff79c6
  -->\n<!-- Purple #bd93f9 -->\n<!-- Red #ff5555 -->\n<!-- Yellow #f1fa8c -->\n\n<!--
  note: $drac-dark-yellow, -->\n<!-- abstract: $drac-cyan, -->\n<!-- info: $drac-light-blue,
  -->\n<!-- tip: $drac-teal, -->\n<!-- success: $drac-green, -->\n<!-- question: $drac-light-green,
  -->\n<!-- warning: $drac-orange, -->\n<!-- failure: $drac-dark-red, -->\n<!-- danger:
  $drac-red, -->\n<!-- bug: $drac-dark-pink, -->\n<!-- example: $drac-purple, -->\n<!--
  quote: $drac-grey -->\n````\n"
date: 0001-01-01
description: This is what the default admonition styles look like and how to create
  them.
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Admonitions</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"This is what the default admonition styles look
    like and how to create them.\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Admonitions</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"This is what the default admonition styles
    look like and how to create them.\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
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
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Admonitions\n
    \   </h1>\n</section>    <section class=\"body\">\n        <div class=\"admonition
    note is-collapsible collapsible-open\">\n<p class=\"admonition-title\">open by
    default</p>\n</div>\n<p>you can open a details tab with '+'</p>\n<!-- prettier-ignore
    -->\n<div class=\"admonition note is-collapsible collapsible-closed\">\n<p class=\"admonition-title\">closed
    by default</p>\n<p>you can open a details tab with '+'</p>\n</div>\n<h2 id=\"all-of-the-admonitions\">all
    of the admonitions <a class=\"header-anchor\" href=\"#all-of-the-admonitions\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<!-- prettier-ignore -->\n<div
    class=\"admonition note\">\n<p class=\"admonition-title\">Note</p>\n<p>a note</p>\n</div>\n<!--
    prettier-ignore -->\n<div class=\"admonition abstract\">\n<p class=\"admonition-title\">Abstract</p>\n<p>an
    abstract</p>\n</div>\n<!-- prettier-ignore -->\n<div class=\"admonition info\">\n<p
    class=\"admonition-title\">Info</p>\n<p>admonitions</p>\n</div>\n<!-- prettier-ignore
    -->\n<div class=\"admonition tip\">\n<p class=\"admonition-title\">Tip</p>\n<p>You
    should think about using admonitions</p>\n</div>\n<!-- prettier-ignore -->\n<div
    class=\"admonition success\">\n<p class=\"admonition-title\">Success</p>\n<p>Run
    Successful!</p>\n</div>\n<!-- prettier-ignore -->\n<div class=\"admonition question\">\n<p
    class=\"admonition-title\">Question</p>\n<p>What do you think of this?</p>\n</div>\n<!--
    prettier-ignore -->\n<div class=\"admonition source\">\n<p class=\"admonition-title\">Source</p>\n<p>Add
    some source code.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"nb\">print</span><span
    class=\"p\">(</span><span class=\"s1\">&#39;hello world&#39;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span>&lt;!-- prettier-ignore --&gt;\n!!!
    warning\n    a warning\n\n&lt;!-- prettier-ignore --&gt;\n!!! failure\n    a failure\n\n&lt;!--
    prettier-ignore --&gt;\n!!! danger\n    some danger\n\n&lt;!-- prettier-ignore
    --&gt;\n!!! bug\n    a bug\n\n&lt;!-- prettier-ignore --&gt;\n!!! example\n    an
    example\n\n    ``` python\n    print(&#39;hello world&#39;)\n    ```\n\n&lt;!--
    prettier-ignore --&gt;\n!!! quote\n\n    a quote\n\n    &gt; include a nice quote\n\n&lt;!--
    Background #282a36 --&gt;\n&lt;!-- Current Line #44475a --&gt;\n&lt;!-- Foreground
    #f8f8f2 --&gt;\n&lt;!-- Comment #6272a4 --&gt;\n&lt;!-- Cyan #8be9fd --&gt;\n&lt;!--
    Green #50fa7b --&gt;\n&lt;!-- Orange #ffb86c --&gt;\n&lt;!-- Pink #ff79c6 --&gt;\n&lt;!--
    Purple #bd93f9 --&gt;\n&lt;!-- Red #ff5555 --&gt;\n&lt;!-- Yellow #f1fa8c --&gt;\n\n&lt;!--
    note: $drac-dark-yellow, --&gt;\n&lt;!-- abstract: $drac-cyan, --&gt;\n&lt;!--
    info: $drac-light-blue, --&gt;\n&lt;!-- tip: $drac-teal, --&gt;\n&lt;!-- success:
    $drac-green, --&gt;\n&lt;!-- question: $drac-light-green, --&gt;\n&lt;!-- warning:
    $drac-orange, --&gt;\n&lt;!-- failure: $drac-dark-red, --&gt;\n&lt;!-- danger:
    $drac-red, --&gt;\n&lt;!-- bug: $drac-dark-pink, --&gt;\n&lt;!-- example: $drac-purple,
    --&gt;\n&lt;!-- quote: $drac-grey --&gt;\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>
    \   </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Admonitions</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"This is what the default admonition styles look
    like and how to create them.\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Admonitions</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"This is what the default admonition styles
    look like and how to create them.\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Admonitions\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Admonitions\n    </h1>\n</section>    <section class=\"body\">\n        <div
    class=\"admonition note is-collapsible collapsible-open\">\n<p class=\"admonition-title\">open
    by default</p>\n</div>\n<p>you can open a details tab with '+'</p>\n<!-- prettier-ignore
    -->\n<div class=\"admonition note is-collapsible collapsible-closed\">\n<p class=\"admonition-title\">closed
    by default</p>\n<p>you can open a details tab with '+'</p>\n</div>\n<h2 id=\"all-of-the-admonitions\">all
    of the admonitions <a class=\"header-anchor\" href=\"#all-of-the-admonitions\"><svg
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
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<!-- prettier-ignore -->\n<div
    class=\"admonition note\">\n<p class=\"admonition-title\">Note</p>\n<p>a note</p>\n</div>\n<!--
    prettier-ignore -->\n<div class=\"admonition abstract\">\n<p class=\"admonition-title\">Abstract</p>\n<p>an
    abstract</p>\n</div>\n<!-- prettier-ignore -->\n<div class=\"admonition info\">\n<p
    class=\"admonition-title\">Info</p>\n<p>admonitions</p>\n</div>\n<!-- prettier-ignore
    -->\n<div class=\"admonition tip\">\n<p class=\"admonition-title\">Tip</p>\n<p>You
    should think about using admonitions</p>\n</div>\n<!-- prettier-ignore -->\n<div
    class=\"admonition success\">\n<p class=\"admonition-title\">Success</p>\n<p>Run
    Successful!</p>\n</div>\n<!-- prettier-ignore -->\n<div class=\"admonition question\">\n<p
    class=\"admonition-title\">Question</p>\n<p>What do you think of this?</p>\n</div>\n<!--
    prettier-ignore -->\n<div class=\"admonition source\">\n<p class=\"admonition-title\">Source</p>\n<p>Add
    some source code.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
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
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"nb\">print</span><span
    class=\"p\">(</span><span class=\"s1\">&#39;hello world&#39;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n</div>\n<pre
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
    \       \n<div class=\"highlight\"><pre><span></span>&lt;!-- prettier-ignore --&gt;\n!!!
    warning\n    a warning\n\n&lt;!-- prettier-ignore --&gt;\n!!! failure\n    a failure\n\n&lt;!--
    prettier-ignore --&gt;\n!!! danger\n    some danger\n\n&lt;!-- prettier-ignore
    --&gt;\n!!! bug\n    a bug\n\n&lt;!-- prettier-ignore --&gt;\n!!! example\n    an
    example\n\n    ``` python\n    print(&#39;hello world&#39;)\n    ```\n\n&lt;!--
    prettier-ignore --&gt;\n!!! quote\n\n    a quote\n\n    &gt; include a nice quote\n\n&lt;!--
    Background #282a36 --&gt;\n&lt;!-- Current Line #44475a --&gt;\n&lt;!-- Foreground
    #f8f8f2 --&gt;\n&lt;!-- Comment #6272a4 --&gt;\n&lt;!-- Cyan #8be9fd --&gt;\n&lt;!--
    Green #50fa7b --&gt;\n&lt;!-- Orange #ffb86c --&gt;\n&lt;!-- Pink #ff79c6 --&gt;\n&lt;!--
    Purple #bd93f9 --&gt;\n&lt;!-- Red #ff5555 --&gt;\n&lt;!-- Yellow #f1fa8c --&gt;\n\n&lt;!--
    note: $drac-dark-yellow, --&gt;\n&lt;!-- abstract: $drac-cyan, --&gt;\n&lt;!--
    info: $drac-light-blue, --&gt;\n&lt;!-- tip: $drac-teal, --&gt;\n&lt;!-- success:
    $drac-green, --&gt;\n&lt;!-- question: $drac-light-green, --&gt;\n&lt;!-- warning:
    $drac-orange, --&gt;\n&lt;!-- failure: $drac-dark-red, --&gt;\n&lt;!-- danger:
    $drac-red, --&gt;\n&lt;!-- bug: $drac-dark-pink, --&gt;\n&lt;!-- example: $drac-purple,
    --&gt;\n&lt;!-- quote: $drac-grey --&gt;\n</pre></div>\n\n</pre>\n\n\n    </section>\n</article>"
  raw.md: "---\ntitle: Admonitions\ndescription: This is what the default admonition
    styles look like and how to create them.\n---\n\n???+ note open by default\nyou
    can open a details tab with '+'\n\n<!-- prettier-ignore -->\n??? note closed by
    default\n    you can open a details tab with '+'\n\n## all of the admonitions\n\n<!--
    prettier-ignore -->\n!!! note\n    a note\n\n<!-- prettier-ignore -->\n!!! abstract\n
    \   an abstract\n\n<!-- prettier-ignore -->\n!!! info\n\n    admonitions\n\n<!--
    prettier-ignore -->\n!!! tip\n\n    You should think about using admonitions\n\n<!--
    prettier-ignore -->\n!!! success\n\n    Run Successful!\n\n<!-- prettier-ignore
    -->\n!!! question\n\n    What do you think of this?\n\n<!-- prettier-ignore -->\n!!!
    source\n    Add some source code.\n    ```python\n    print('hello world')\n\n````\n\n<!--
    prettier-ignore -->\n!!! warning\n    a warning\n\n<!-- prettier-ignore -->\n!!!
    failure\n    a failure\n\n<!-- prettier-ignore -->\n!!! danger\n    some danger\n\n<!--
    prettier-ignore -->\n!!! bug\n    a bug\n\n<!-- prettier-ignore -->\n!!! example\n
    \   an example\n\n    ``` python\n    print('hello world')\n    ```\n\n<!-- prettier-ignore
    -->\n!!! quote\n\n    a quote\n\n    > include a nice quote\n\n<!-- Background
    #282a36 -->\n<!-- Current Line #44475a -->\n<!-- Foreground #f8f8f2 -->\n<!--
    Comment #6272a4 -->\n<!-- Cyan #8be9fd -->\n<!-- Green #50fa7b -->\n<!-- Orange
    #ffb86c -->\n<!-- Pink #ff79c6 -->\n<!-- Purple #bd93f9 -->\n<!-- Red #ff5555
    -->\n<!-- Yellow #f1fa8c -->\n\n<!-- note: $drac-dark-yellow, -->\n<!-- abstract:
    $drac-cyan, -->\n<!-- info: $drac-light-blue, -->\n<!-- tip: $drac-teal, -->\n<!--
    success: $drac-green, -->\n<!-- question: $drac-light-green, -->\n<!-- warning:
    $drac-orange, -->\n<!-- failure: $drac-dark-red, -->\n<!-- danger: $drac-red,
    -->\n<!-- bug: $drac-dark-pink, -->\n<!-- example: $drac-purple, -->\n<!-- quote:
    $drac-grey -->\n````\n"
published: true
slug: admonitions
title: Admonitions


---


???+ note open by default
you can open a details tab with '+'

<!-- prettier-ignore -->
??? note closed by default
    you can open a details tab with '+'

## all of the admonitions

<!-- prettier-ignore -->
!!! note
    a note

<!-- prettier-ignore -->
!!! abstract
    an abstract

<!-- prettier-ignore -->
!!! info

    admonitions

<!-- prettier-ignore -->
!!! tip

    You should think about using admonitions

<!-- prettier-ignore -->
!!! success

    Run Successful!

<!-- prettier-ignore -->
!!! question

    What do you think of this?

<!-- prettier-ignore -->
!!! source
    Add some source code.
    ```python
    print('hello world')

````

<!-- prettier-ignore -->
!!! warning
    a warning

<!-- prettier-ignore -->
!!! failure
    a failure

<!-- prettier-ignore -->
!!! danger
    some danger

<!-- prettier-ignore -->
!!! bug
    a bug

<!-- prettier-ignore -->
!!! example
    an example

    ``` python
    print('hello world')
    ```

<!-- prettier-ignore -->
!!! quote

    a quote

    > include a nice quote

<!-- Background #282a36 -->
<!-- Current Line #44475a -->
<!-- Foreground #f8f8f2 -->
<!-- Comment #6272a4 -->
<!-- Cyan #8be9fd -->
<!-- Green #50fa7b -->
<!-- Orange #ffb86c -->
<!-- Pink #ff79c6 -->
<!-- Purple #bd93f9 -->
<!-- Red #ff5555 -->
<!-- Yellow #f1fa8c -->

<!-- note: $drac-dark-yellow, -->
<!-- abstract: $drac-cyan, -->
<!-- info: $drac-light-blue, -->
<!-- tip: $drac-teal, -->
<!-- success: $drac-green, -->
<!-- question: $drac-light-green, -->
<!-- warning: $drac-orange, -->
<!-- failure: $drac-dark-red, -->
<!-- danger: $drac-red, -->
<!-- bug: $drac-dark-pink, -->
<!-- example: $drac-purple, -->
<!-- quote: $drac-grey -->
````
