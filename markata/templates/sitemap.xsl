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
{% endblock %}
