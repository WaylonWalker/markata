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
    xmlns:atom="http://www.w3.org/2005/Atom"
    xmlns:dc="http://purl.org/dc/elements/1.1/"
    xmlns:sm="http://www.sitemaps.org/schemas/sitemap/0.9">
  <xsl:output method="html" version="1.0" encoding="UTF-8" indent="yes"/>
  <xsl:template match="/">
    <html xmlns="http://www.w3.org/1999/xhtml" lang='en'>  
    <head>
<title>Markata's Docs</title>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<meta name="description" content="Documentation for using the Markata static site generator" />
 <link href="/favicon.ico" rel="icon" type="image/png" />

<link rel="stylesheet" href="/post.css" />
<link rel="stylesheet" href="/app.css" />
<script src="/theme.js"></script>


        <meta property="og:author_email" content="waylon@waylonwalker.com" />

    </head>
    <body>
<header class='flex justify-center items-center p-8'>

    <nav class='flex justify-center items-center my-8'>
        <a
            href='/'>markata</a>
        <a
            href='https://github.com/WaylonWalker/markata'>GitHub</a>
        <a
            href='https://markata.dev/docs/'>docs</a>
        <a
            href='https://markata.dev/plugins/'>plugins</a>
    </nav>

    <div>
        <label id="theme-switch" class="theme-switch" for="checkbox-theme" title="light/dark mode toggle">
            <input type="checkbox" id="checkbox-theme" />
            <div class="slider round"></div>
        </label>
    </div>
</header>    <main>
        <div class="container">
          <header>
            <h1>
              <xsl:choose>
                <!-- RSS title -->
                <xsl:when test="/rss/channel/title">
                  <xsl:value-of select="/rss/channel/title"/>
                </xsl:when>
                <!-- Atom title -->
                <xsl:when test="/atom:feed/atom:title">
                  <xsl:value-of select="/atom:feed/atom:title"/>
                </xsl:when>
              </xsl:choose>
            </h1>

            <p>
              <xsl:choose>
                <!-- RSS description -->
                <xsl:when test="/rss/channel/description">
                  <xsl:value-of select="/rss/channel/description"/>
                </xsl:when>
                <!-- Atom subtitle (rough equivalent) -->
                <xsl:when test="/atom:feed/atom:subtitle">
                  <xsl:value-of select="/atom:feed/atom:subtitle"/>
                </xsl:when>
              </xsl:choose>
            </p>

            <a class="head_link" target="_blank">
              <xsl:attribute name="href">
                <xsl:choose>
                  <!-- RSS site link -->
                  <xsl:when test="/rss/channel/link">
                    <xsl:value-of select="/rss/channel/link"/>
                  </xsl:when>
                  <!-- Atom: prefer rel='alternate' link, otherwise first link -->
                  <xsl:when test="/atom:feed/atom:link[@rel='alternate']">
                    <xsl:value-of select="/atom:feed/atom:link[@rel='alternate'][1]/@href"/>
                  </xsl:when>
                  <xsl:when test="/atom:feed/atom:link">
                    <xsl:value-of select="/atom:feed/atom:link[1]/@href"/>
                  </xsl:when>
                  <!-- Fallback: Atom id (not ideal as a URL, but better than nothing) -->
                  <xsl:otherwise>
                    <xsl:value-of select="/atom:feed/atom:id"/>
                  </xsl:otherwise>
                </xsl:choose>
              </xsl:attribute>
              Visit Website &#x2192;
            </a>
          </header>

          <section class="recent">
            <h2>Recent Items</h2>
            <ul>
              <!-- Union: RSS items OR Atom entries -->
              <xsl:for-each select="/rss/channel/item | /atom:feed/atom:entry">
<li class='my-6 list-none'>
    <a class='w-full text-white no-underline rounded bg-neutral-800 shadow-lgc shadow-neutral-900 hover:bg-neutral-700'>
        <!-- href: RSS <link> text OR Atom <atom:link>@href -->
        <xsl:attribute name="href">
            <xsl:choose>
                <!-- RSS link -->
                <xsl:when test="link">
                    <xsl:value-of select="link" />
                </xsl:when>
                <!-- Atom: prefer rel='alternate' -->
                <xsl:when test="atom:link[@rel='alternate']">
                    <xsl:value-of select="atom:link[@rel='alternate'][1]/@href" />
                </xsl:when>
                <!-- Atom: fallback to first link -->
                <xsl:when test="atom:link">
                    <xsl:value-of select="atom:link[1]/@href" />
                </xsl:when>
            </xsl:choose>
        </xsl:attribute>

        <h2 class='mx-4 mt-1 text-3xl font-bold'>
            <!-- RSS <title> OR Atom <atom:title> -->
            <xsl:value-of select="title | atom:title" />
        </h2>

        <p class='mx-4 text-neutral-200'>
            <!-- RSS <description> OR Atom <atom:summary> OR Atom <atom:content> -->
            <xsl:choose>
                <xsl:when test="description">
                    <xsl:value-of select="description" />
                </xsl:when>
                <xsl:when test="atom:summary">
                    <xsl:value-of select="atom:summary" />
                </xsl:when>
                <xsl:when test="atom:content">
                    <xsl:value-of select="atom:content" />
                </xsl:when>
            </xsl:choose>
        </p>

        <time class='mx-4 text-neutral-600'>
            <!-- RSS pubDate OR Atom published/updated -->
            <xsl:value-of select="pubDate | atom:published | atom:updated" />
        </time>
    </a>
</li>
              </xsl:for-each>
            </ul>
          </section>
        </div>
    </main>
            <footer style='margin-top: 20rem;'>© 2025</footer>

    </body>
    </html>
  </xsl:template>
</xsl:stylesheet>