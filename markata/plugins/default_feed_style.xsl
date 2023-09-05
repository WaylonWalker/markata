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
    <html xmlns="http://www.w3.org/1999/xhtml">
      <head>
        <title><xsl:value-of select="/rss/channel/title"/> Web Feed</title>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
        <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1"/>
        <style type="text/css">
          * { box-sizing: border-box; }

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
          img { max-width: 100%; }
            body { --gap: 5vw; margin: 0; font-family: system-ui; line-height: 1.7;
                  background-color: #1f2022;
                  color: #eefbfe;
                font-family: "Space Mono", monospace;
            }
            body a { color: inherit; }

          h1,h2,h3 { margin-block-start: 0; margin-block-end: 0; }
          .pb-5 { padding-bottom: calc(var(--gap) / 2); }
          .meta { color: #676767; }
          .container {
            display: grid;
            gap: var(--gap);
            max-width: 46rem;
            width: 95%;
            margin: auto;
          }
          .intro {
            background-color: #EEDD82;
            background-color: #FAFACC;
            margin-block-end: var(--gap);
            padding-block: calc(var(--gap) / 2);
            color: black;
          }
          .intro .container {
            gap: 1rem;
            grid-template-columns:  4fr 2fr;
            align-items: top;
          }
          @media (min-width: 40rem) {
            .intro .container {
              grid-template-columns:  4fr 1fr;
              align-items: center;
            }
          }
          .recent {
            padding-block-end: var(--gap);
          }
        </style>
      </head>
      <body>
        <nav class="container-md px-3 py-2 mt-2 mt-md-5 mb-5 markdown-body">
          <p class="bg-yellow-light ml-n1 px-1 py-1 mb-1">
            <strong>This is a web feed,</strong> also known as an RSS feed. <strong>Subscribe</strong> by copying the URL from the address bar into your newsreader.
          </p>
          <p class="text-gray">
            Visit <a href="https://aboutfeeds.com">About Feeds</a> to get started with newsreaders and subscribing. It’s free.
          </p>
        </nav>
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
      </body>
    </html>
  </xsl:template>
</xsl:stylesheet>

