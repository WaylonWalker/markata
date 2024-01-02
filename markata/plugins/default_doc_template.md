{{ ast.get_docstring(tree) }}

{% for node in nodes %}
!! {{node.type}} <h2 id='{{node.name}}' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>{{node.name}} <em class='small'>{{node.type}}</em></h2>
{{ indent(ast.get_docstring(node) or '', '    ') }}
???+ source "{{node.name}} <em class='small'>source</em>"

```python

{{ indent(ast.get_source_segment(raw_source, node), '        ') }}
```

{% endfor %}
