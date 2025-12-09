{% extends "base.xsl" %}
{% block content %}
    <main>
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
                {% include 'rss_card.html' %}
              </xsl:for-each>
            </ul>
          </section>
        </div>
    </main>
{% endblock %}
