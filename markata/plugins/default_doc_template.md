---

{{ ast.get_docstring(tree) }}

---

{% for node in nodes %}
{% if ast.get_docstring(node) %}
!!! {{ node.type }}
    <h2 id="{{ node.name }}" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">{{ node.name }} <em class="small">{{ node.type }}</em></h2>

    {{ ast.get_docstring(node) | replace('\n\n\n', '\n\n') | indent(4) }}

???+ source "{{ node.name }} <em class='small'>source</em>"
    ```python
    {{ ast.get_source_segment(raw_source, node) | replace('\n\n\n', '\n\n') | indent(4) }}
    ```
{% endif %}
{% endfor %}
