{% extends "base.xsl" %}
{% block content %}
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
        <h2>The Items</h2>
        <ul>
        <xsl:for-each select="sm:urlset/sm:url">
            {% include 'sitemap_card.html' %}
        </xsl:for-each>
        </ul>
        </section>
    </div>

</main>
{% endblock %}
