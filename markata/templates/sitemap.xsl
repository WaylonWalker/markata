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
<xsl:stylesheet version="3.0"
    xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns:sm="http://www.sitemaps.org/schemas/sitemap/0.9">
  <xsl:output method="html" version="1.0" encoding="UTF-8" indent="yes"/>
  <xsl:template match="/">
      <html xmlns="http://www.w3.org/1999/xhtml" lang='{{ markata.config.lang }}'>

<head>
    <title>{{ markata.config.title }}</title>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="{{ markata.config.description }}" />
     <link href="{{ markata.config.icon.name }}" rel="icon" type="image/png" />
    <script>
        {% include "theme.js" %}
    </script>
    <style>
        {% include "post.css" %}
    </style>
</head>

        <nav class="container-md px-3 py-2 mt-2 mt-md-5 mb-5 markdown-body">
          <p class="bg-yellow-light ml-n1 px-1 py-1 mb-1">
            <strong>This is a web feed,</strong> also known as an RSS feed. <strong>Subscribe</strong> by copying the URL from the address bar into your newsreader.
          </p>
    <a
    href='/'>
        markata
    </a>
    <a
    href='https://github.com/WaylonWalker/markata'>
        GitHub
    </a>
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
                         <a href="" alt="edit post url" title="edit this post">
            </a>
        </h1>
    </section>

    <main>
        <div class="container">
          <header>
            <h1><xsl:value-of select="/urlset/title"/></h1>
            <p><xsl:value-of select="/urlset/description"/></p>
            <a class="head_link" target="_blank">
              <xsl:attribute name="href">
                <xsl:value-of select="/urlset/link"/>
              </xsl:attribute>
              Visit Website &#x2192;
            </a>
          </header>
          <section class="recent">
            <h2>Recent Items</h2>
            <ul>
            <xsl:for-each select="sm:urlset/sm:url">
                <li>
                    <a href="{sm:loc}">
                        <xsl:value-of select="sm:loc"/>
                    <xsl:value-of select="sm:lastmod"/>
                    <xsl:value-of select="sm:title"/>
                    <xsl:value-of select="sm:description"/>
                    </a>
                </li>
            </xsl:for-each>
            </ul>
          </section>
        </div>

    </main>
    <footer style='margin-top: 20rem;'>© {{ markata.config.today.year }}</footer>
</body>

</html>
  </xsl:template>
</xsl:stylesheet>
