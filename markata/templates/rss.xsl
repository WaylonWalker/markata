{% extends "base.xsl" %}
{% block content %}
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
    </main>
{% endblock %}
