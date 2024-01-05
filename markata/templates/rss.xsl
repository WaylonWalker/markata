<?xml version="1.0" encoding="utf-8"?>
<!--

# Pretty Feed

Styles an RSS/Atom feed, making it friendly for humans viewers, and adds a link
to aboutfeeds.com for new user onboarding. See it in action:

   https://interconnected.org/home/feed


## How to use

1. Download this XML stylesheet from the following URL and host it on your own
   domain (this is a limitation of XSL in browsers):

   https://github.com/genmon/aboutfeeds/blob/main/tools/pretty-feed-v3.xsl

2. Include the XSL at the top of the RSS/Atom feed, like:

```
<?xml version="1.0" encoding="UTF-8"?>
<?xml-stylesheet href="/PATH-TO-YOUR-STYLES/pretty-feed-v3.xsl" type="text/xsl"?>
```

3. Serve the feed with the following HTTP headers:

```
Content-Type: application/xml; charset=utf-8  # not application/rss+xml
x-content-type-options: nosniff
```

(These headers are required to style feeds for users with Safari on iOS/Mac.)



## Limitations

- Styling the feed *prevents* the browser from automatically opening a
  newsreader application. This is a trade off, but it's a benefit to new users
  who won't have a newsreader installed, and they are saved from seeing or
  downloaded obscure XML content. For existing newsreader users, they will know
  to copy-and-paste the feed URL, and they get the benefit of an in-browser feed
  preview.
- Feed styling, for all browsers, is only available to site owners who control
  their own platform. The need to add both XML and HTTP headers makes this a
  limited solution.


## Credits

pretty-feed is based on work by lepture.com:

   https://lepture.com/en/2019/rss-style-with-xsl

This current version is maintained by aboutfeeds.com:

   https://github.com/genmon/aboutfeeds


## Feedback

This file is in BETA. Please test and contribute to the discussion:

     https://github.com/genmon/aboutfeeds/issues/8

-->
<xsl:stylesheet version="3.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                xmlns:atom="http://www.w3.org/2005/Atom" xmlns:dc="http://purl.org/dc/elements/1.1/"
                xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
  <xsl:output method="html" version="1.0" encoding="UTF-8" indent="yes"/>
  <xsl:template match="/">
    <html xmlns="http://www.w3.org/1999/xhtml" lang='en'>

<head>
    {% if post.title or config.title %}
    <title>{{ post.title or config.title }}</title>
    {% endif %}
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    {% if post.description or config.description %}
    <meta name="description" content="{{ post.description or config.description }}" />
    {% endif %} {% if config.icon %}
    <link href="/{{ config.icon }}" rel="icon" type="image/png" />
    {% endif %}

    <script>
        function setTheme(theme) {
            document.documentElement.setAttribute("data-theme", theme);
        }

        function detectColorSchemeOnLoad() {
            //local storage is used to override OS theme settings
            if (localStorage.getItem("theme")) {
                if (localStorage.getItem("theme") == "dark") {
                    setTheme("dark");
                } else if (localStorage.getItem("theme") == "light") {
                    setTheme("light");
                }
            } else if (!window.matchMedia) {
                //matchMedia method not supported
                setTheme("light");
                return false;
            } else if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
                //OS theme setting detected as dark
                setTheme("dark");
            } else {
                setTheme("light");
            }
        }
        detectColorSchemeOnLoad();
        document.addEventListener(
            "DOMContentLoaded",
            function () {
                //identify the toggle switch HTML element
                const toggleSwitch = document.querySelector(
                    '#theme-switch input[type="checkbox"]',
                );

                //function that changes the theme, and sets a localStorage variable to track the theme between page loads
                function switchTheme(e) {
                    if (e.target.checked) {
                        localStorage.setItem("theme", "dark");
                        document.documentElement.setAttribute("data-theme", "dark");
                        toggleSwitch.checked = true;
                    } else {
                        localStorage.setItem("theme", "light");
                        document.documentElement.setAttribute("data-theme", "light");
                        toggleSwitch.checked = false;
                    }
                }

                //listener for changing themes
                toggleSwitch.addEventListener("change", switchTheme, false);

                //pre-check the dark-theme checkbox if dark-theme is set
                if (document.documentElement.getAttribute("data-theme") == "dark") {
                    toggleSwitch.checked = true;
                }
            },
            false,
        );
    </script>
    <style>
        :root {
            --color-bg: {{ config.style.color_bg }} ;
            --color-bg-code: {{ config.style.color_bg_code }};
            --color-text: {{ config.style.color_text }};
            --color-link: {{ config.style.color_link }};
            --color-accent: {{ config.style.color_accent }};
            --overlay-brightness: {{ config.style.overlay_brightness }};
            --body-width: {{ config.style.body_width }};
        }

        [data-theme="dark"] {
            --color-bg: {{ config.style.color_bg }};
            --color-bg-code: {{ config.style.color_bg_code }};
            --color-text: {{ config.style.color_text }};
            --color-link: {{ config.style.color_link }};
            --color-accent: {{ config.style.color_accent }};
            --overlay-brightness: {{ config.style.overlay_brightness }};
            --body-width: {{ config.style.body_width }};
        }

        [data-theme="light"] {
            --color-bg: {{ config.style.color_bg_light }};
            --color-bg-2: {{ config.style.color_bg_light_2 }};
            --color-bg-code: {{ config.style.color_bg_code_light }};
            --color-text: {{ config.style.color_text_light }};
            --color-link: {{ config.style.color_link_light }};
            --color-accent: {{ config.style.color_accent_light }};
            --overlay-brightness: {{ config.style.overlay_brightness_light }};
        }

        html {
            font-family: "Space Mono", monospace;
            background: var(--color-bg);
            color: var(--color-text);
        }

        a {
            color: var(--color-link);
        }

        main a {
            max-width: 100%;
        }

        .heading-permalink {
            font-size: .7em;
        }

        body {
            max-width: var(--body-width);
            margin: 5rem auto;
            padding: 0 .5rem;
            font-size: 1rem;
            line-height: 1.56;
        }

        blockquote {
            background: var(--color-bg);
            filter: brightness(var(--overlay-brightness));
            border-left: 4px solid var(--color-accent);
            border-radius: 4px;
            box-shadow:
                -0.8rem 0rem 1rem -1rem #f1fa8c,
                0.2rem 0rem 1rem rgb(0, 0, 0, .4);
            padding-left: 1rem;
            margin: 1rem;
        }

        li.post {
            list-style-type: None;
            padding: .2rem 0;
        }

        pre.wrapper {
            padding: 0;
            box-shadow: 0.2rem 0rem 1rem rgb(0, 0, 0, .4);
            display: flex;
            flex-direction: column;
            position: relative;
            margin: 2rem;
        }

        pre {
            margin: 0;
            padding: 1rem;
            min-width: -webkit-fill-available;
            max-width: fit-content;
            overflow-x: auto;
        }

        pre .filepath {
            margin: 0;
            padding-left: 1rem;
            border-radius: 4px 4px 0 0;
            background: black;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        pre .filepath p {
            margin: 0
        }

        pre .filepath .right {
            display: flex;
            gap: .2rem;
            align-items: center;
        }

        pre::-webkit-scrollbar {
            height: 4px;
            background-color: transparent;
        }

        pre::-webkit-scrollbar-thumb {
            background-color: #d3d3d32e;
            border-radius: 2px;
        }

        pre::-webkit-scrollbar-track {
            background-color: transparent;
        }

        .copy-wrapper {
            background: none;
            position: absolute;
            width: 100%;
            z-index: 100;
            display: flex;
            justify-content: flex-end;
        }

        button.copy {
            z-index: 100;
            background: none;
            fill: #ffffff45;
            border: none;
            width: 32px;
            align-self: flex-end;
            top: 0;
            right: 0;
            margin: 0.5rem 0.2rem;

        }

        button.copy:hover {
            fill: white
        }

        a.help {
            fill: #ffffff45;
        }

        a.help:hover {
            fill: white;
        }

        a.help svg {
            height: 24px;
            width: 24px;
        }

        .highlight {
            background: var(--color-bg-code);
            color: var(--color-text);
            filter: brightness(var(--overlay-brightness));
            border-radius: 0 0 4px 4px;
        }

        .highlight .c {
            color: #8b8b8b
        }

        /* Comment */
        .highlight .err {
            color: #960050;
            background-color: #1e0010
        }

        /* Error */
        .highlight .k {
            color: #c678dd
        }

        /* Keyword */
        .highlight .l {
            color: #ae81ff
        }

        /* Literal */
        .highlight .n {
            color: #abb2bf
        }

        /* Name */
        .highlight .o {
            color: #c678dd
        }

        /* Operator */
        .highlight .p {
            color: #abb2bf
        }

        /* Punctuation */
        .highlight .ch {
            color: #8b8b8b
        }

        /* Comment.Hashbang */
        .highlight .cm {
            color: #8b8b8b
        }

        /* Comment.Multiline */
        .highlight .cp {
            color: #8b8b8b
        }

        /* Comment.Preproc */
        .highlight .cpf {
            color: #8b8b8b
        }

        /* Comment.PreprocFile */
        .highlight .c1 {
            color: #8b8b8b
        }

        /* Comment.Single */
        .highlight .cs {
            color: #8b8b8b
        }

        /* Comment.Special */
        .highlight .gd {
            color: #c678dd
        }

        /* Generic.Deleted */
        .highlight .ge {
            font-style: italic
        }

        /* Generic.Emph */
        .highlight .gi {
            color: #a6e22e
        }

        /* Generic.Inserted */
        .highlight .gs {
            font-weight: bold
        }

        /* Generic.Strong */
        .highlight .gu {
            color: #8b8b8b
        }

        /* Generic.Subheading */
        .highlight .kc {
            color: #c678dd
        }

        /* Keyword.Constant */
        .highlight .kd {
            color: #c678dd
        }

        /* Keyword.Declaration */
        .highlight .kn {
            color: #c678dd
        }

        /* Keyword.Namespace */
        .highlight .kp {
            color: #c678dd
        }

        /* Keyword.Pseudo */
        .highlight .kr {
            color: #c678dd
        }

        /* Keyword.Reserved */
        .highlight .kt {
            color: #c678dd
        }

        /* Keyword.Type */
        .highlight .ld {
            color: #e6db74
        }

        /* Literal.Date */
        .highlight .m {
            color: #ae81ff
        }

        /* Literal.Number */
        .highlight .s {
            color: #e6db74
        }

        /* Literal.String */
        .highlight .na {
            color: #a6e22e
        }

        /* Name.Attribute */
        .highlight .nb {
            color: #98c379
        }

        /* Name.Builtin */
        .highlight .nc {
            color: #abb2bf
        }

        /* Name.Class */
        .highlight .no {
            color: #c678dd
        }

        /* Name.Constant */
        .highlight .nd {
            color: #abb2bf
        }

        /* Name.Decorator */
        .highlight .ni {
            color: #abb2bf
        }

        /* Name.Entity */
        .highlight .ne {
            color: #a6e22e
        }

        /* Name.Exception */
        .highlight .nf {
            color: #61afef
        }

        /* Name.Function */
        .highlight .nl {
            color: #abb2bf
        }

        /* Name.Label */
        .highlight .nn {
            color: #abb2bf
        }

        /* Name.Namespace */
        .highlight .nx {
            color: #a6e22e
        }

        /* Name.Other */
        .highlight .py {
            color: #abb2bf
        }

        /* Name.Property */
        .highlight .nt {
            color: #c678dd
        }

        /* Name.Tag */
        .highlight .nv {
            color: #abb2bf
        }

        /* Name.Variable */
        .highlight .ow {
            color: #c678dd
        }

        /* Operator.Word */
        .highlight .w {
            color: #abb2bf
        }

        /* Text.Whitespace */
        .highlight .mb {
            color: #ae81ff
        }

        /* Literal.Number.Bin */
        .highlight .mf {
            color: #ae81ff
        }

        /* Literal.Number.Float */
        .highlight .mh {
            color: #ae81ff
        }

        /* Literal.Number.Hex */
        .highlight .mi {
            color: #ae81ff
        }

        /* Literal.Number.Integer */
        .highlight .mo {
            color: #ae81ff
        }

        /* Literal.Number.Oct */
        .highlight .sa {
            color: #e6db74
        }

        /* Literal.String.Affix */
        .highlight .sb {
            color: #e6db74
        }

        /* Literal.String.Backtick */
        .highlight .sc {
            color: #e6db74
        }

        /* Literal.String.Char */
        .highlight .dl {
            color: #e6db74
        }

        /* Literal.String.Delimiter */
        .highlight .sd {
            color: #98c379
        }

        /* Literal.String.Doc */
        .highlight .s2 {
            color: #98c379
        }

        /* Literal.String.Double */
        .highlight .se {
            color: #ae81ff
        }

        /* Literal.String.Escape */
        .highlight .sh {
            color: #e6db74
        }

        /* Literal.String.Heredoc */
        .highlight .si {
            color: #e6db74
        }

        /* Literal.String.Interpol */
        .highlight .sx {
            color: #e6db74
        }

        /* Literal.String.Other */
        .highlight .sr {
            color: #e6db74
        }

        /* Literal.String.Regex */
        .highlight .s1 {
            color: #e6db74
        }

        /* Literal.String.Single */
        .highlight .ss {
            color: #e6db74
        }

        /* Literal.String.Symbol */
        .highlight .bp {
            color: #abb2bf
        }

        /* Name.Builtin.Pseudo */
        .highlight .fm {
            color: #61afef
        }

        /* Name.Function.Magic */
        .highlight .vc {
            color: #abb2bf
        }

        /* Name.Variable.Class */
        .highlight .vg {
            color: #abb2bf
        }

        /* Name.Variable.Global */
        .highlight .vi {
            color: #abb2bf
        }

        /* Name.Variable.Instance */
        .highlight .vm {
            color: #abb2bf
        }

        /* Name.Variable.Magic */
        .highlight .il {
            color: #ae81ff
        }

        /* Literal.Number.Integer.Long */

        /* Tab style starts here */
        .tabbed-set {
            position: relative;
            display: flex;
            flex-wrap: wrap;
            margin: 1em 0;
            border-radius: 0.1rem;
        }

        .tabbed-set>input {
            display: none;
        }

        .tabbed-set label {
            width: auto;
            padding: 0.9375em 1.25em 0.78125em;
            font-weight: 700;
            font-size: 0.84em;
            white-space: nowrap;
            border-bottom: 0.15rem solid transparent;
            border-top-left-radius: 0.1rem;
            border-top-right-radius: 0.1rem;
            cursor: pointer;
            transition: background-color 250ms, color 250ms;
        }

        .tabbed-set .tabbed-content {
            width: 100%;
            display: none;
            box-shadow: 0 -.05rem #ddd;
        }

        .tabbed-set input {
            position: absolute;
            opacity: 0;
        }

        /* fonts */
        h1 {
            font-weight: 700;
        }

        h1#title a {
            font-size: 16px;
        }

        h1,
        h2,
        h3,
        h4,
        h5,
        h6 {
            margin-top: 3rem;
        }

        h1 {
            font-size: 2.5em;
            margin-top: 5rem;
        }

        h2 {
            font-size: 1.63rem;
            margin-top: 5rem;
        }



        p {
            font-size: 21px;
            font-style: normal;
            font-variant: normal;
            font-weight: 400;
            line-height: 1.5;
        }

        @media only screen and (max-width: 700px) {
            p {
                font-size: 18px;
            }
        }

        @media only screen and (max-width: 600px) {
            p {
                font-size: 16px;
            }
        }

        @media only screen and (max-width: 500px) {
            p {
                font-size: 14px;
            }
        }

        @media only screen and (max-width: 400px) {
            p {
                font-size: 12px;
            }
        }


        pre {
            font-style: normal;
            font-variant: normal;
            font-weight: 400;
            line-height: 18.5714px;
            */
        }

        a {
            font-weight: 600;
            text-decoration-color: var(--color-accent);
            color: var(--color-link);
            padding: .3rem .5rem;
            display: inline-block;
        }

        .admonition,
        details {
            box-shadow: 0.2rem 0rem 1rem rgb(0, 0, 0, .4);
            margin: 5rem 0;
            border: 1px solid transparent;
            border-radius: 4px;
            text-align: left;
            padding: 0;
            border: 0;

        }

        .admonition {
            padding-bottom: 1rem;
        }

        details[open] {
            padding-bottom: .5rem;
        }

        .admonition p {
            padding: .2rem .6rem;
        }

        .admonition-title,
        .details-title,
        summary {
            background: var(--color-bg-2);
            padding: 0;
            margin: 0;
            position: sticky;
            top: 0;
            z-index: 10;
        }

        summary:hover {
            cursor: pointer;
        }

        summary.admonition-title,
        summary.details-title {
            padding: .5rem;
            padding-left: 1rem;
        }

        .note {
            border-left: 4px solid #f1fa8c;
            box-shadow:
                -0.8rem 0rem 1rem -1rem #f1fa8c,
                0.2rem 0rem 1rem rgb(0, 0, 0, .4);
        }

        .note>.admonition-title {
            border-bottom: 1px solid #3c3d2d;
        }

        .abstract {
            border-left: 4px solid #8be9fd;
            box-shadow:
                -0.8rem 0rem 1rem -1rem #8be9fd,
                0.2rem 0rem 1rem rgb(0, 0, 0, .4);
        }

        .abstract>.admonition-title {
            border-bottom: 1px solid #2c3a3f;
        }

        .info {
            border-left: 4px solid;
            box-shadow:
                -0.8rem 0rem 1rem -1rem #8bb0fd,
                0.2rem 0rem 1rem rgb(0, 0, 0, .4);
        }

        .info>.admonition-title {
            border-bottom: 1px solid #2c313f;
        }

        .tip {
            border-left: 4px solid #008080;
            box-shadow:
                -0.8rem 0rem 1rem -1rem #008080,
                0.2rem 0rem 1rem rgb(0, 0, 0, .4);
        }

        .tip>.admonition-title {
            border-bottom: 1px solid #1b2a2b;
        }

        .success {
            border-left: 4px solid #50fa7b;
            box-shadow:
                -0.8rem 0rem 1rem -1rem #50fa7b,
                0.2rem 0rem 1rem rgb(0, 0, 0, .4);
        }

        .success>.admonition-title {
            border-bottom: 1px solid #263e2b;
        }

        .question {
            border-left: 4px solid #a7fcbd;
            box-shadow:
                -0.8rem 0rem 1rem -1rem #a7fcbd,
                0.2rem 0rem 1rem rgb(0, 0, 0, .4);
        }

        .question>.admonition-title {
            border-bottom: 1px solid #303e35;
        }

        .warning {
            border-left: 4px solid #ffb86c;
            box-shadow:
                -0.8rem 0rem 1rem -1rem #ffb86c,
                0.2rem 0rem 1rem rgb(0, 0, 0, .4);
        }

        .warning>.admonition-title {
            border-bottom: 1px solid #3f3328;
        }

        .failure {
            border-left: 4px solid #b23b3b;
            box-shadow:
                -0.8rem 0rem 1rem -1rem #b23b3b,
                0.2rem 0rem 1rem rgb(0, 0, 0, .4);
        }

        .failure>.admonition-title {
            border-bottom: 1px solid #34201f;
        }

        .danger {
            border-left: 4px solid #ff5555;
            box-shadow:
                -0.8rem 0rem 1rem -1rem #ff5555,
                0.2rem 0rem 1rem rgb(0, 0, 0, .4);
        }

        .danger>.admonition-title {
            border-bottom: 1px solid #402523;
        }

        .bug {
            border-left: 4px solid #b2548a;
            box-shadow:
                -0.8rem 0rem 1rem -1rem #b2548a,
                0.2rem 0rem 1rem rgb(0, 0, 0, .4);
        }

        .bug>.admonition-title {
            border-bottom: 1px solid #32232c;
        }

        .example {
            border-left: 4px solid #bd93f9;
            box-shadow:
                -0.8rem 0rem 1rem -1rem #bd93f9,
                0.2rem 0rem 1rem rgb(0, 0, 0, .4);
        }

        .example>.admonition-title {
            border-bottom: 1px solid #332d3e;
        }

        .source {
            border-left: 4px solid #bd93f9;
            box-shadow:
                -0.8rem 0rem 1rem -1rem #bd93f9,
                0.2rem 0rem 1rem rgb(0, 0, 0, .4);
        }

        .source>.admonition-title {
            border-bottom: 1px solid #332d3e;
        }

        .quote {
            border-left: 4px solid #999;
            box-shadow:
                -0.8rem 0rem 1rem -1rem #999,
                0.2rem 0rem 1rem rgb(0, 0, 0, .4);
        }

        .quote>.admonition-title {
            border-bottom: 1px solid #2d2e2f;
        }

        table {
            margin: 1rem 0;
            border-collapse: collapse;
            border-spacing: 0;
            display: block;
            max-width: -moz-fit-content;
            max-width: fit-content;
            overflow-x: auto;
            white-space: nowrap;
        }

        table thead th {
            border: solid 1px var(--color-text);
            padding: 10px;
            text-align: left;
        }

        table tbody td {
            border: solid 1px var(--color-text);
            padding: 10px;
        }

        .theme-switch {
            z-index: 10;
            display: inline-block;
            height: 34px;
            position: relative;
            width: 60px;

            display: flex;
            justify-content: flex-end;
            margin-right: 1rem;
            margin-left: auto;
            position: fixed;
            right: 1rem;
            top: 1rem;
        }

        .theme-switch input {
            display: none;

        }

        .slider {
            background-color: #ccc;
            bottom: 0;
            cursor: pointer;
            left: 0;
            position: absolute;
            right: 0;
            top: 0;
            transition: .4s;
        }

        .slider:before {
            background-color: #fff;
            bottom: 4px;
            content: "";
            height: 26px;
            left: 4px;
            position: absolute;
            transition: .4s;
            width: 26px;
        }

        input:checked+.slider {
            background-color: #343434;
        }

        input:checked+.slider:before {
            background-color: #848484;
        }

        input:checked+.slider:before {
            transform: translateX(26px);
        }

        .slider.round {
            border-radius: 34px;
        }

        .slider.round:before {
            border-radius: 50%;
        }

        main p img {
            width: 100%;
            width: -moz-available;
            width: -webkit-fill-available;
            width: fill-available;
        }

        details>* {
            margin: 1rem;
        }

        .admonition>* {
            margin: 1rem;
        }

        p.admonition-title,
        summary {
            margin: 0;
            padding-left: 1.2rem;
        }

        .small {
            font-size: .9rem;
            color: #888;
        }

        admonition+admonition {
            margin-top: 20rem;
        }

        ::-webkit-scrollbar {
            height: 12px;
            background-color: transparent;
        }

        ::-webkit-scrollbar-thumb {
            background-color: #d3d3d32e;
            border-radius: 6px;
        }

        ::-webkit-scrollbar-track {
            background-color: transparent;
        }
    </style>

    {% if 'markata.plugins.service_worker' in config.hooks %}
    <script>
        if ("serviceWorker" in navigator) {
            navigator.serviceWorker.register("/service-worker.js");
            navigator.serviceWorker.addEventListener("controllerchange", () => {
                console.log("new worker");
                window.location.reload();
            });
        }
    </script>
    {% endif %} {% for text in config.head.text %} {{ text.value }}{% endfor %}
    {% for meta in config.head.meta %}{% for _meta in meta %}
    <meta content="{{ meta.content }}" name="{{ meta.name }}" />
    {% endfor %}{% endfor %}
</head>

        <nav class="container-md px-3 py-2 mt-2 mt-md-5 mb-5 markdown-body">
          <p class="bg-yellow-light ml-n1 px-1 py-1 mb-1">
            <strong>This is a web feed,</strong> also known as an RSS feed. <strong>Subscribe</strong> by copying the URL from the address bar into your newsreader.
          </p>
    {% for text, link in config.nav.items() %}
    <a
    href='{{"/" if link.startswith("/")}}{{"" if "://" in link else config["path_prefix"]}}{{link.lstrip("/")}}'>
        {{text}}
    </a>
    {% endfor %}
</nav>

<body>
    <div>
        <label id="theme-switch" class="theme-switch" for="checkbox-theme" title="light/dark mode toggle">
            <input type="checkbox" id="checkbox-theme" />
            <div class="slider round"></div>
        </label>
    </div>
    <section class="title">
        <h1 id="title">
            {{ title }} {% if config.get %}
            <a href="{{ edit_link }}" alt="edit post url" title="edit this post">
            </a>
            {% endif %}
        </h1>
    </section>

    <main>
        <div class="container">
          <header>
            <h1><xsl:value-of select="/rss/channel/title"/></h1>
            <p><xsl:value-of select="/rss/channel/description"/></p>
            <a class="head_link" target="_blank">
              <xsl:attribute name="href">
                <xsl:value-of select="/rss/channel/link"/>
              </xsl:attribute>
              Visit Website &#x2192;
            </a>
          </header>
          <section class="recent">
            <h2>Recent Items</h2>
            <xsl:for-each select="/rss/channel/item">
              <div class="pb-5">
                <h3>
                  <a target="_blank">
                    <xsl:attribute name="href">
                      <xsl:value-of select="link"/>
                    </xsl:attribute>
                    <xsl:value-of select="title"/>
                  </a>
                </h3>
                <p>
                <xsl:value-of select="description"/>
                </p>
                <small class="meta">
                  Published: <xsl:value-of select="pubDate" />
                </small>
              </div>
            </xsl:for-each>
          </section>
        </div>

      {{ body }}
    </main>
    <footer>© {{ today.year }}</footer>
</body>

</html>
  </xsl:template>
</xsl:stylesheet>
