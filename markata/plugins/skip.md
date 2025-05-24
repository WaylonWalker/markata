---
date: 2025-05-24
description: "Post Skipping and Caching Plugin !!! function post_model  function ???+
  source \"post_model  source \" !!! function load  function ???+ source \"load  source
  \"\u2026"
published: false
slug: markata/plugins/skip
title: skip.py


---

---

Post Skipping and Caching Plugin

---

!!! function
    <h2 id="post_model" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">post_model <em class="small">function</em></h2>

    Add skip attribute to post models.

???+ source "post_model <em class='small'>source</em>"
    ```python
    def post_model(markata: "Markata") -> None:
        """Add skip attribute to post models."""
        markata.post_models.append(PostModel)
    ```
!!! function
    <h2 id="load" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">load <em class="small">function</em></h2>

    Runs after posts are loaded to check if they should be skipped.

???+ source "load <em class='small'>source</em>"
    ```python
    def load(markata: "Markata") -> None:
        """Runs after posts are loaded to check if they should be skipped."""
        raw_should_skip = os.environ.get("MARKATA_SKIP", "true")
        should_skip = raw_should_skip.lower() in ["true", "1", "t", "y", "yes", "on"]
        if not should_skip:
            return
        for post in markata.posts:
            if hasattr(post, "raw"):
                key = markata.make_hash("skip", post.raw)
                if markata.cache.get(key) == "done":
                    post.skip = True
        console.log(f"{len(markata.filter('skip'))}/{len(markata.posts)} posts skipped")
        console.log(
            f"{len(markata.filter('not skip'))}/{len(markata.posts)} posts not skipped"
        )
    ```
!!! function
    <h2 id="save" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">save <em class="small">function</em></h2>

    Save the 'done' status for processed posts.

???+ source "save <em class='small'>source</em>"
    ```python
    def save(markata: "Markata") -> None:
        """Save the 'done' status for processed posts."""
        # for post in markata.posts:
        raw_should_skip = os.environ.get("MARKATA_SKIP", "true")
        should_skip = raw_should_skip.lower() in ["true", "1", "t", "y", "yes", "on"]
        if not should_skip:
            return
        for post in markata.filter("not skip"):
            if hasattr(post, "raw"):
                if post.output_html.exists():
                    key = markata.make_hash("skip", post.raw)
                    markata.cache.set(key, "done")
    ```