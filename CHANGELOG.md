# Markata Changelog

## 0.9.1

- HUGE DOCS UPDATE
- Fix: feeds now update with any content change
- Fix: og tags should use property not name
- Feat: centralized `jinja_env` plugin that is used for template rendering accross all plugins
- Feat: added support for multiple templates per article https://markata.dev/multi-template/
- Feat: skip now supports environment variable `MARKATA_SKIP=true/false` to control skipping

## 0.9.0

- Feat: add skip plugin
- Fix: wikilinks uses a faster and more robust link lookup
- Feat: most plugins now only do work on posts not marked as skip
- Fix: wikilinks now give good error messages to show you which article is throwing the warning

### Performance Improvements

Here are a couple of build times ran back to back on the same content with 0.8.2 and 0.9.0

## Markata

63 posts

- 0.9.0 - 4.308s cold - 1.457s hot
- 0.8.2 - 9.063s cold - 1.542s hot

## waylonwalker.com

exactly 2k posts

- 0.9.0 - 69.200s cold - 2.860 hot
- 0.8.2 - 479s cold - 4.683 hot

#### Core Optimizations

- Added map cache with statistics tracking in Markata core
- Optimized template rendering with bytecode caching
- Parallelized markdown rendering and file loading
- Added configurable cache expiration times for different components

#### Component-specific Optimizations

- Updated auto_description to use markdown-it for better performance
- Improved date parsing with dateparser fallback
- Optimized post loading with bulk processing
- Modernized pydantic validators to use new field_validator syntax

## 0.8.2

- Fix: markata installs setuptools required by one dependency
- Fix: cleaup cli output
- Fix: speed up cli starup with some lazy imports
- Fix: all cache.adds were replaced with cache.set
- Fix: Updated to new typer format requiring name=
- Fix: teardown only runs if a the build process was started, i.e. some clis
  will not need to teardown

## 0.8.1

### Feeds have partials

The `feeds` plugin now has configurable `partial_template` that can be used to
render only the inside of the feeds page.  This is indended to allow you to
load small feeds into a page with htmx.

### Better Jinja Templates

Markata now fully supports jinja templates with a loader that will load from
your templates directory, the markata built-in templates, and from a
dynamically generated templates directory in your .markata.cache directory.

#### cli

You can list out your templates and configuration with the following command

``` bash
markata templates show
```

#### Variables

The following variables are available within jinja templates for post
templates.  This is now consistent accross all three built in plugins that
render jinja templates.

##### post_template

- `__version__` - the version of markata
- `markata` - the markata instance
- `config` - the markata config
- `body` - the body of the post
- `post` - the current post object

##### feeds

Similarly from within rendering feeds.

- `__version__` - the version of markata
- `markata` - the markata instance
- `config` - the markata config
- `posts` - the list of posts
- `post` - a pseudo post object with title, slug, description, and date for template consistency
- `feed` - the current feed object

##### jinja_md

Similar to posts from within jinja_md to render a markdown post as a template.

- `__version__` - the version of markata
- `markata` - the markata instance
- `body` - the body of the post
- `config` - the markata config
- `post` - the current post object

#### Feeds cli

The feeds cli will help show which templates each feed will be using.

``` bash
❯ markata feeds show
                                          Feeds 6
┏━━━━━━━━━━━━━━━━━┳━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃            Feed ┃ posts ┃ config                                                        ┃
┡━━━━━━━━━━━━━━━━━╇━━━━━━━╇━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┩
│ project_gallery │ 2     │ DEFAULT_TITLE: All Posts                                      │
│                 │       │ title: Project Gallery                                        │
│                 │       │ slug: project-gallery                                         │
│                 │       │ name: project_gallery                                         │
│                 │       │ filter: 'project-gallery' in str(path)                        │
│                 │       │ sort: title                                                   │
│                 │       │ reverse: False                                                │
│                 │       │ rss: True                                                     │
│                 │       │ sitemap: True                                                 │
│                 │       │ card_template: card.html                                      │
│                 │       │ template: feed.html                                           │
│                 │       │ rss_template: rss.xml                                         │
│                 │       │ sitemap_template: sitemap.xml                                 │
│                 │       │ xsl_template: rss.xsl                                         │
│                 │       │                                                               │
│            docs │ 10    │ DEFAULT_TITLE: All Posts                                      │
│                 │       │ title: Documentation                                          │
│                 │       │ slug: docs                                                    │
│                 │       │ name: docs                                                    │
│                 │       │ filter: "markata" not in slug and "tests" not in slug and ... │
│                 │       │ sort: slug                                                    │
│                 │       │ reverse: False                                                │
│                 │       │ rss: True                                                     │
│                 │       │ sitemap: True                                                 │
│                 │       │ card_template: card.html                                      │
│                 │       │ template: feed.html                                           │
│                 │       │ rss_template: rss.xml                                         │
│                 │       │ sitemap_template: sitemap.xml                                 │
│                 │       │ xsl_template: rss.xsl                                         │
│                 │       │                                                               │
│         autodoc │ 17    │ DEFAULT_TITLE: All Posts                                      │
│                 │       │ title: AutoDoc Python Modules.                                │
│                 │       │ slug: autodoc                                                 │
│                 │       │ name: autodoc                                                 │
│                 │       │ filter: "markata" in slug and "plugin" not in slug and "te... │
│                 │       │ sort: slug                                                    │
│                 │       │ reverse: False                                                │
│                 │       │ rss: True                                                     │
│                 │       │ sitemap: True                                                 │
│                 │       │ card_template: card.html                                      │
│                 │       │ template: feed.html                                           │
│                 │       │ rss_template: rss.xml                                         │
│                 │       │ sitemap_template: sitemap.xml                                 │
│                 │       │ xsl_template: rss.xsl                                         │
│                 │       │                                                               │
│             all │ 73    │ DEFAULT_TITLE: All Posts                                      │
│                 │       │ title: All Markata Modules                                    │
│                 │       │ slug: all                                                     │
│                 │       │ name: all                                                     │
│                 │       │ filter: True                                                  │
│                 │       │ sort: date                                                    │
│                 │       │ reverse: False                                                │
│                 │       │ rss: True                                                     │
│                 │       │ sitemap: True                                                 │
│                 │       │ card_template: card.html                                      │
│                 │       │ template: feed.html                                           │
│                 │       │ rss_template: rss.xml                                         │
│                 │       │ sitemap_template: sitemap.xml                                 │
│                 │       │ xsl_template: rss.xsl                                         │
│                 │       │                                                               │
│    core_modules │ 17    │ DEFAULT_TITLE: All Posts                                      │
│                 │       │ title: Markata Core Modules                                   │
│                 │       │ slug: core_modules                                            │
│                 │       │ name: core_modules                                            │
│                 │       │ filter: 'plugin' not in slug and 'test' not in slug and ti... │
│                 │       │ sort: date                                                    │
│                 │       │ reverse: False                                                │
│                 │       │ rss: True                                                     │
│                 │       │ sitemap: True                                                 │
│                 │       │ card_template: card.html                                      │
│                 │       │ template: feed.html                                           │
│                 │       │ rss_template: rss.xml                                         │
│                 │       │ sitemap_template: sitemap.xml                                 │
│                 │       │ xsl_template: rss.xsl                                         │
│                 │       │                                                               │
│         plugins │ 42    │ DEFAULT_TITLE: All Posts                                      │
│                 │       │ title: Markata Plugins                                        │
│                 │       │ slug: plugins                                                 │
│                 │       │ name: plugins                                                 │
│                 │       │ filter: 'plugin' in slug and 'test' not in slug               │
│                 │       │ sort: date                                                    │
│                 │       │ reverse: False                                                │
│                 │       │ rss: True                                                     │
│                 │       │ sitemap: True                                                 │
│                 │       │ card_template: card.html                                      │
│                 │       │ template: feed.html                                           │
│                 │       │ rss_template: rss.xml                                         │
│                 │       │ sitemap_template: sitemap.xml                                 │
│                 │       │ xsl_template: rss.xsl                                         │
│                 │       │                                                               │
└─────────────────┴───────┴───────────────────────────────────────────────────────────────┘
```

## 0.8.0

- pydantic support

### Pydantic Support

Now plugins are configured through a pydantic Config object.

### breaking changes

There are a number of breaking changes going into 0.8.0. Use caution when
upgrading.

#### glob config is now under markata.glob

```diff
- [markata]
- glob_patterns = "pages/**/*.md"
+ [markata.glob]
+ glob_patterns = "pages/**/*.md"
```

#### Feeds are now a list

Feeds are now a list of Objects within the configuration that you choose from
whether its toml or yaml.  Also templates_dir is now configurable, and once you
have a templates dir it is better to specify templates by name relative to your
templates_dir.

```toml
[markata]
templates_dir = "pages/templates"

[markata.feeds.published]
template="archive_template.html"
card_template = "feed_card.html"
filter="date<=today and templateKey in ['blog-post', 'til'] and status.lower()=='published'"
sort="date"
```

> old

```toml
[[markata.feeds.published]]
template="pages/templates/archive_template.html"
card_template = "pages/templates/feed_card.html"
filter="date<=today and templateKey in ['blog-post', 'til'] and status.lower()=='published'"
sort="date"
```

> new

### markata.summary.filter_count is now a list

The old way was to set up a dict, where the keys were the name, now its a list
of Objects with an explicit name field.

```toml
[markata.summary.filter_count.drafts]
filter="published == 'False'"
color='red'
```

> Old

```toml
[[markata.summary.filter_count]]
name='drafts'
filter="published == 'False'"
color='red'
```

## 0.7.4

- Fix: Icon resize broken from PIL 10.0.0 release

## 0.7.3

- Fix: broken default feed card template
- Fix: broken links for index.md in feeds

## 0.7.2

- Fix: broken `markata new` command due to pydantic v2 compatability with copier.

## 0.7.0

- Adopt ruff linter 0.7.0.dev1 #142
- add support for markdown_it_py backend #145
- add trogon tui support 0.7.0.dev6

### trogon tui support

install with pip

```bash
pip install 'markata[tui]'
```

usage

```bash
markata tui
```

## 0.6.4

- Fix: Icon resize broken from PIL 10.0.0 release

## 0.6.3

- Fix: broken `markata new` command due to pydantic v2 compatability with copier.

## 0.6.2

Update License and Security files.

## 0.6.1

- Fix: allow feeds to be used from within Markdown.

### Feeds in Markdown

```markdown
{% for post in markata.feeds.docs.posts %}
[{{post.title}}](/{{post.slug}})
{% endfor %}
```

## 0.6.0

- Fix: article_html is not available to jinja 0.6.0.dev1 #105
- Fix: service worker is not upating client without hard refresh 0.6.0.dev2
  #106
- Feat: create teardown lifecycle method 0.6.0.dev3 #110
- Fix: implement teardown on all raises 0.6.0.dev4 #111
- Fix: implement teardown in pyinstrument plugin 0.6.0.dev5 #112
- Feat: Automatically call teardown without needing to remember it before raise
  0.6.0.dev6 #113
- Fix: only stop the profiler if it is running 0.6.0.dev7 #114
- Fix: map was giving inconsistent results 0.6.0.dev8 #116
- Fix: tui continuously rebuilds if an input file exists in a parent directory
  to the output directory 0.6.0.dev9 #118
- Feat: report the cache stats for the current run 0.6.0.dev11 #121
- Fix: prevent zerodivisionerror when reporting stats 0.6.0.dev12 #122
- Fix: properly set the pyinstrument profiler to prevent recurrsion errors
  0.6.dev13 #123
- Clean: cli attributes (`runner`, `summary`, `server`, `plugins`) are now
  added as Markata properties through `register_attr` rather than directly to
  the class 0.6.0.dev13 #107
- Fix: Markata tui will remain running even when the runner fails 0.6.0.dev13
  #107
- Fix: Pinned to `textual<0.2.0` due to breaking changes 0.6.0.dev13 #107
- Fix: snyk remove seuptools from requirements 0.6.0.dev14 #130
- Fix: Markata.filter was missing post and m 0.6.0.dev14 #133
- Feat: add `path_prefix` config for gh-pages deploy 0.6.0.dev15 #132 fixes #57
- Feat: created `markata plugin show` cli command 0.6.0.dev16 #109
- Feat: use `path_prefix` in nav entries 0.6.0.dev17 #136
- Feat: enable wikilinks extension by default 0.6.0.dev18 #138
- Feat: `markata` instance is available form inside card_template for feeds
  0.6.0.dev19 #139

### wikilinks

wikilinks are now enabled by default ex: `[[home-page]]`. This will create a
link `<a class="wikilink" href="/home-page/">home-page</a>`. This will
automagically work if you leave `markata.plugins.flat_slug` plugin enabled
(which is by default).

> ProTip: this was highly inspired by the
> [marksman-lsp](https://github.com/artempyanykh/marksman) by
> [artempyanykh](https://github.com/artempyanykh/), which can autocomplete post
> links in this style for you.

[[home-page]]

## 0.5.5

- Fix: Icon resize broken from PIL 10.0.0 release

## 0.5.4

- Fix: broken `markata new` command due to pydantic v2 compatability with copier.

## 0.5.2

- clean up unnecessary images_url is missing warning #104 0.5.2.dev1

## 0.5.1

- fix: contrast ratio on admonitions was insufficient for A11y #103 0.5.1.dev0

## 0.5.0

- Create `new` cli command for creating new `blogs`, `posts`, and `plugins` #93
  0.5.0.dev16 [base_cli-docs](https://markata.dev/markata/plugins/base_cli/)
- Remove unused function clif that was the original entrypoint #81 0.5.0.dev8
- Allow template variables to be used in head config #88 0.5.0.dev12
- Expose `markata.__version__` to templates as `__version__` #89 0.5.0.dev13
- Fix, ignore post_template save on posts without an html attribute #92 0.5.0.dev13
- Fix #33 sluggify paths #69 **BREAKING CHANGE** 0.5.0.dev6
- Configurable template #70 0.5.dev5, #85 0.5.0.dev11
- Fix #40 Images overlfow outside of body #66 0.5.0.dev3
- Created entrypoint hook allowing for users to extend marka with jinja
  exensions #60 0.5.0.dev2
- Moved to PEP 517 build #59 0.5.0.dev1
- new `markata.plugins.redirects` will create redirect html files as a backup when
  server-side redirects fail. #76 0.5.0.dev10
  [redirects-docs](https://markata.dev/markata/plugins/redirects/)
- create a slugify migration script #82
- DeepMerge `config_overrides` with config in post render methods #91 0.5.0.dev13
- Create ipython extension to automatically load markata #79 0.5.0.dev15
- Fix: images wrapped in a link overflow outside the body #96
- new `markata.plugins.service_worker` plugin to create service workers and
  enable offline mode on sites #94 0.5.0.dev15
  [service-worker-docs](https://markata.dev/markata/plugins/service-worker/)
- Fix: icons were relatively linked, and were broken for any page other than
  index, they are now absolutely linked to the root of the site. #97 0.5.0.dev16
- Fix: auto_descriptions were not rendered on first pass in tui or at all in
  build due to auto_description running after jinja_md. #100 0.5.0.dev18
- Fix: give redirect pages a uniqe description and title #101 0.5.0.dev19

### `new` cli command

More information in the [base_cli-docs](https://markata.dev/markata/plugins/base_cli/).

```bash
# create a new blog template
# copier requires you to specify a directory
markata new blog [directory]

# create a new blog post
markata new post

# create a new plugin
markata new plugin

markata new --help

 Usage: markata new [OPTIONS] COMMAND [ARGS]...

 create new things from templates

╭─ Options ─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ --help          Show this message and exit.                                                                                                                                       │
╰───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
╭─ Commands ────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ blog       Create a new blog from using the template from https://github.com/WaylonWalker/markata-blog-starter.                                                                   │
│ plugin     Create a new plugin using the template at https://github.com/WaylonWalker/markata-plugin-template.                                                                     │
│ post       Create new blog post in the pages directory from the template at  https://github.com/WaylonWalker/markata-post-template.                                               │
╰───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
```

### sluggify paths

`python-sluggify` was implemented to ensure good urls are in place despite the
name of the original file.

For examples of how `python-slugify` will change your url's see the
[project's home page](https://pypi.org/project/python-slugify/). One
difference is that `markata` will leave `/`'s for routing in the slugs.

#### OPTING OUT

If you have an existing site and do not want to implement redirects or if you
do not want to use slugify, you can opt out by setting `slugify=False` in your
`markata.toml`.

```toml
[markata]
slugify=false
```

#### Migrating to slugify

From the command line with `markata>=0.5.0` installed run the
migration script from the command line to create a redirects file in the
default location. This should avoid all 404's as it will create a redirects
file that many static hosting providers will issue a server-side 301 for, and
for those that don't, markata.plugins.redirects creates a redirect html page,
that will kick in as a backup.

```bash
python -m markata.scripts.migrate_to_slugify
```

### configurable page template

Now injects seo into the default template through configuration. Here is an
example, by adding this to your `markata.toml` configuration.

```toml
[[markata.head.meta]]
name = "og:type"
content = "article"

[[markata.head.meta]]
name = "og:author"
content = "Waylon Walker"

[[markata.head.meta]]
name = "og:site_name"
content = "Waylon Walker"

[[markata.head.meta]]
name = "theme-color"
content="#322D39"

[[markata.head.meta]]
name = "twitter:creator"
content="@_waylonwalker"
```

You will end up with these meta tags in your html.

```html
<meta name="og:type" content="article" />
<meta name="og:author" content="Waylon Walker" />
<meta name="og:site_name" content="Waylon Walker" />
<meta name="theme-color" content="#322D39" />
<meta name="twitter:creator" content="@_waylonwalker" />
```

You can have an array of toml tables with a key of text. The text will be
added as plain text to the end of the head of each page.

```toml
[[markata.head]]
text = """
<style>
img {
width: 100%;
height: auto;
}
ul {
  display: flex;
  flex-wrap: wrap;
}

li {
  flex: 1 2 400px;
}
</style>

"""
```

Descriptions will now properly end up in each page.

```html
<meta name="description" content="{{ description }}" />
```

## Config Overrides

Each post can override config settings such as `head`. New meta tags can be
added to a single post, or anything that your template might reference from
config.

```yaml
config_overrides:
  head:
    meta:
      - content: waylonwalker
        name: author
      - content: "@_waylonwalker"
        name: "twitter:creator"
    link:
      - href: https://waylonwalker.com/that-special-post/
        rel: canonical
    text:
      - value: <link rel='stylesheet' href='/my-extra-styles.css' />
```

## ipython extension

Markata has an ipython extension if you want ipython to automatically load with
an instance of `Markata` mapped to `m` and `markata` you can add the following
to your `~/.ipython/profile_default/ipython_config.py`

```python
c.InteractiveShellApp.extensions.append('markata')
```

## 0.4.1

- Issue FutureWarning for upcoming change to slugify change that will change urls

To keep existing behavior add this to your `markata.toml`.

```toml
[markata]
slugify=false
```

## 0.4.0

- feat: add html logging with [setup_logging](/markata/plugins/setup_logging/)
  plugin is all new closes #37
- fix: remove HTML tidy as the site generator tag
- feat: create configurable [navbar](https://markata.dev/nav)
- perf: prevent double runs on pre-render and post-render #39
- perf: prevent duplicate ruun from to_dict calling pre-render #53
  - `to_dict` only runs up to `render` phase if necessary as directed by `register_attr`
- perf: only prettify if configured #54
- fix: pyinstrument will not create a second profiler causing it to end in
  errors #50
- fix: sites without feeds config do not create an index #55

### Double Runs

Previously markata would catch AttributeError and run the previous step any
time you ran a step too early. The way this was implemented caused some steps
such as pre-render and post-render to run twice with every single run.

This change will no longer catch attribute errors. If you run into any issues
with your plugins not running before asking for attributes created by your
plugin make sure that you implement the
[@register_attr](https://markata.dev/markata/hookspec/#register_attr-function)
decorator.

### Prettify

prettify html has been turned off by default as beautifulsoup4 prettify was
taking a significant time, and was often popping up as the slowest parts in my
personal `_profile`. If you want to continue running prettify throughout the
build you can set a flag in your config to continue running prettify.

```toml
[markata]
prettify_html = true
```

## 0.3.0

Skipped from a bump2version misconfiguration.

## 0.2.0

- feat: [auto_description](/markata/plugins/auto_description/) plugin is all
  new closes #13
- deprecated: long_description has been deprecated by auto_description
- fix: [covers](/markata/plugins/covers/) plugin would previously skip
  every time.
- feat: [`markata clean`](/markata/plugins/base_cli/#clean-function) cleans up
  your cache and output from the command line
- fix: [`publish_source`](/markata/plugins/publish_source/) plugin will now
  ignore any non yaml serializable values
- feat: Default template colors are now customizable
- feat: Default template now has light and dark theme
- feat: map now can map entire posts
- feat: [prevnext](/markata/plugins/prevnext/) plugin was added to link between
  posts closes #20
- feat: [jinja_md](/markata/plugins/jinja_md/) plugins was added to incorporate
  jinja into all the markdown
- breaking: [feeds](/markata/plugins/feeds) config now has feeds and
  feeds_config
- feat: `output_html` can now be specified in the frontmatter
  [see example](/markata/plugins/publish_html/#explicityly-set-the-output)
- feat: edit link is now included in the default page template closes #21

### breaking change to feeds config

If you are using the 0.1.0 version of feeds, and have configured custom
templates in `markata.feeds.template` and `markata.feeds.card_template` they
will need to be moved to `markata.feeds_config`.

Here is what you need to do to update your feeds_config.

```diff
+ [markata.feeds_config]
+ template="pages/templates/archive_template.html"
+ card_template="plugins/feed_card_template.html"
- [markata.feeds]
- template="pages/templates/archive_template.html"
- card_template="plugins/feed_card_template.html"
```

`markata.feeds` will only be used to configure feeds pages.

### map entire posts

`post` is now exposed to the `markata.map` object, allowing you to return a
list of posts.

```python
m = Markata()
# 'post' will return the entire post
m.map('post', filter='"git" in tags')
```

### Customizable colors

```toml
[markata]
# default dark theme
color_bg = '#1f2022'
color_bg_code = '#1f2022'
color_text = '#eefbfe'
color_link = '#47cbff'
color_accent = '#e1bd00c9'
overlay_brightness = '.85'

# pink and purple
color_bg = 'deeppink'
color_bg_code = 'rebeccapurple'
color_text = 'white'
color_link = 'aqua'
color_accent = 'peachpuff'
overlay_brightness = '1.2'

# default light theme
color_bg_light = '#eefbfe'
color_bg_code_light = '#eefbfe'
color_text_light = '#1f2022'
color_link_light = '#47cbff'
color_accent_light = '#ffeb00'
overlay_brightness_light = '.95'
```

### All New auto_description plugin

- Cache is busted on plugin change
- plugin is configurable
- plugin now has docs

### auto_description Configuration

Open up your `markata.toml` file and add new entries for your
auto_descriptions. You can have multiple desriptions, each one will be named
after the key you give it in your config.

```toml
[markata]
hooks=[
   "markata.plugins.auto_description",
   ]

[markata.auto_description.description]
len=160
[markata.auto_description.long_description]
len=250
[markata.auto_description.super_description]
len=500
```

In the above we will end up with three different descritpions,
(`description`, `long_description`, and `super_description`) each will be the
first number of characters from the document as specified in the config.

### auto_descriptions are no longer duplicated

4e299d6 fixes the dedupe issue that was in develop for a while and closes #24

## 0.1.0

- fix: pyinstument plugin no longer overrides the main cli callback
- feat: default is to run the profiler if pyinstrument is installed
- fix: --profile is now under the build command
- feat: --pretty/--no-pretty will make pretty tracebacks, and skip over
  framework code (closes #4 )
- fix: links are now absolute, so they work from GitHub, thanks
  [MR Destructive](https://github.com/Mr-Destructive)
- deprecate: `article['content_hash']` has been removed with preference for
  hashing on `article['content']`
- create `heading_link` plugin
- doc: How to create your home page. [docs](/home-page/)

### New cli help

After the pyinstrument plugin was fixed --version and --to-json are back, and
--profile is now under the build command.

![image](https://user-images.githubusercontent.com/22648375/150662983-547aebbd-c18c-4c17-8985-a6dc01cd29c7.png)

### New Heading Link Plugin

The new heading link plugin makes it easier to share the exact part of an
article you want with someone, by giving clickable links to the id of the
heading it's next to.

![image](https://user-images.githubusercontent.com/22648375/151718782-08a7cd26-41c1-4f00-a12c-0a208c593e9c.png)

## 0.0.1

Initial Release 
