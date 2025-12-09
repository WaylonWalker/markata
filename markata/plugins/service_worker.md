---
date: 2025-12-09
description: "Adds a service_worker to your site. This will make it installable on
  mobile, viewable offline, and potentially more responsive as the user goes between
  good\u2026"
published: false
slug: markata/plugins/service-worker
title: service_worker.py


---

---

Adds a service_worker to your site.  This will make it installable on mobile,
viewable offline, and potentially more responsive as the user goes between good
and bad connections.

## Configuration

Enable this plugin by adding it to your `markata.toml` hooks list.

``` toml
[markata]
hooks=[
  # your hooks
  "markata.plugins.service_worker",
]
```

If you have any content that you want to  precache, add it to the list of
precache.  You can use devtools, change your network to offline, and see what
files send 404's to the console.  These files likely need precache.

``` toml
[markata]
precache_urls = ['archive-styles.css', 'scroll.css', 'manifest.json']
```

# cache busting

Markata uses the checksum.dirhash of your output directory as the cache key.
This is likely to change and bust the cache on every build.

# pre-caching feeds

You can add and entire feed to your precache, this will automatically load
these posts into the cache anytime someone visits your site and their browser
installs the service worker.

Be nice to your users and don't try to install everything possible in their
cache, but maybe a few that they are most likely to click on.

``` toml
[markata.feeds.recent]
filter="date<today and date>today-timedelta(days=30) and published"
sort="slug"
precache=true
```

> note this assumes that the blog implements a published boolean in each posts
frontmatter.

---

!!! function
    <h2 id="render" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">render <em class="small">function</em></h2>

    sets precache_urls in markata.config to be used in
    `markata.plugins.service_worker.save`.

???+ source "render <em class='small'>source</em>"
    ```python
    def render(markata: "Markata") -> None:
        """
        sets precache_urls in markata.config to be used in
        `markata.plugins.service_worker.save`.
        """

        config = markata.config.service_worker

        if config.precache_feeds:
            for feed, config in markata.config.feeds:
                config.precache_urls.append(f"/{feed}/")

        if config.precache_posts:
            with markata.console.status("pre-caching posts...") as status:
                for post in markata.map("post", **config):
                    status.update(f"pre-caching {post.get('slug', '')}...")
                    config.precache_urls.append(f"/{post.get('slug', '')}/")

        config.precache_urls = list(set(config.precache_urls))
    ```
!!! function
    <h2 id="save" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">save <em class="small">function</em></h2>

    Renders the service-worker.js file with your precache urls, and dirhash.

???+ source "save <em class='small'>source</em>"
    ```python
    def save(markata: "Markata") -> None:
        """
        Renders the service-worker.js file with your precache urls, and dirhash.
        """

        template = Template(markata.config.service_worker.template_file.read_text())
        service_worker_js = template.render(
            __version__=__version__,
            config=copy.deepcopy(markata.config),
            output_dirhash=dirhash(markata.config.output_dir),
        )

        output_file = markata.config.output_dir / "service-worker.js"
        current_content = output_file.read_text() if output_file.exists() else ""
        if current_content != service_worker_js:
            output_file.write_text(service_worker_js)
    ```