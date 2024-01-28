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
            <ul>
            <xsl:for-each select="/rss/channel/item">
                {% include 'rss_card.html' %}
            </xsl:for-each>
            </ul>
          </section>
        </div>
    </main>
{% endblock %}
