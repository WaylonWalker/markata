---
title: {{file.name}}
published: True
slug: {{slug}}
edit_link: {{edit_link}}
path: {{file.stem}}.md
today: {{datetime.datetime.today()}}
description: Docs for {{file.stem}}

---

{{ ast.get_docstring(tree) }}

{% for node in nodes %}
!!! {{node.type}} {{node.name}} <em class='small'>{{node.type}}</em>

{{ indent(ast.get_docstring(node) or '', '    ') }}

    ???+ source "{{node.name}} <em class='small'>source</em>"
        ``` python

{{ indent(ast.get_source_segment(raw_source, node), '        ') }}
        ```

{% endfor %}
