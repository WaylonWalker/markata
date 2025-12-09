---
date: 2025-12-09
description: Markata is a tool for handling directories of markdown.
published: false
slug: markata/init
title: __init__.py


---

---

Markata is a tool for handling directories of markdown.

---

!!! method
    <h2 id="teardown" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">teardown <em class="small">method</em></h2>

    Cleanup and print statistics when Markata is done.

???+ source "teardown <em class='small'>source</em>"
    ```python
    def teardown(self: "Markata"):
            """Cleanup and print statistics when Markata is done."""
            # Print map cache statistics if they exist
            if hasattr(self, "_map_cache_stats"):
                stats = self._map_cache_stats
                total = stats["total"]
                if total > 0:
                    hit_rate = (stats["hits"] / total) * 100
                    self.console.print("\n[yellow]Map Cache Statistics:[/yellow]")
                    self.console.print(f"Total calls: {total}")
                    self.console.print(f"Cache hits: {stats['hits']}")
                    self.console.print(f"Cache misses: {stats['misses']}")
                    self.console.print(f"Hit rate: {hit_rate:.1f}%")
                    self.console.print(
                        f"Cache size: {len(getattr(self, '_filtered_cache', {}))}"
                    )
            if self.stages_ran:
                self._pm.hook.teardown(markata=self)
            return self
    ```
!!! method
    <h2 id="_compile_sort_key" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_compile_sort_key <em class="small">method</em></h2>

    Compile a sort key function for better performance

???+ source "_compile_sort_key <em class='small'>source</em>"
    ```python
    def _compile_sort_key(self, sort: str):
            """Compile a sort key function for better performance"""
            if "datetime" in sort.lower():
                return lambda a: a.get(sort, datetime.datetime(1970, 1, 1))
            if "date" in sort.lower():
                return lambda a: a.get(sort, datetime.date(1970, 1, 1))

            # Create a compiled function for complex sort expressions
            try:
                code = compile(sort, "<string>", "eval")

                def sort_key(a):
                    try:
                        value = eval(code, a.to_dict(), {})
                        if isinstance(value, (int, float)):
                            return value
                        if hasattr(value, "timestamp"):
                            return value.timestamp()
                        if isinstance(value, datetime.date):
                            return datetime.datetime.combine(
                                value,
                                datetime.datetime.min.time(),
                            ).timestamp()
                        return sum(ord(c) for c in str(value))
                    except Exception:
                        return -1

                return sort_key
            except Exception:
                return lambda _: -1
    ```
!!! method
    <h2 id="_get_sort_key" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_get_sort_key <em class="small">method</em></h2>

    Cache compiled sort key functions

???+ source "_get_sort_key <em class='small'>source</em>"
    ```python
    def _get_sort_key(self, sort: str):
            """Cache compiled sort key functions"""
            return self._compile_sort_key(sort)
    ```
!!! method
    <h2 id="_get_eval_globals" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_get_eval_globals <em class="small">method</em></h2>

    Get common globals used in eval operations

???+ source "_get_eval_globals <em class='small'>source</em>"
    ```python
    def _get_eval_globals(self):
            """Get common globals used in eval operations"""
            if not hasattr(self, "_eval_globals"):
                self._eval_globals = {"timedelta": timedelta}
            return self._eval_globals
    ```
!!! method
    <h2 id="_eval_with_article" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">_eval_with_article <em class="small">method</em></h2>

    Evaluate code with article context, reusing dict where possible

???+ source "_eval_with_article <em class='small'>source</em>"
    ```python
    def _eval_with_article(self, code, article, extra_globals=None):
            """Evaluate code with article context, reusing dict where possible"""
            if not hasattr(article, "_eval_dict"):
                article._eval_dict = article.to_dict()
                article._eval_dict.update({"post": article, "m": self})

            globals_dict = self._get_eval_globals()
            if extra_globals:
                globals_dict.update(extra_globals)

            try:
                return eval(code, article._eval_dict, globals_dict)
            except Exception:
                return None
    ```