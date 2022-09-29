# Markata Changelog

## 0.5.0-dev

* Remove unused function clif that was the original entrypoint #81 0.5.0.dev8
* Allow template variables to be used in head config #88 0.5.0.dev12
* Expose `markata.__version__` to templates as `__version__` #89 0.5.0.dev13
* Fix #33 sluggify paths #69 **BREAKING CHANGE** 0.5.0.dev6
* Configurable template #70 0.5.dev5, #85 0.5.0.dev11
* Fix #40 Images overlfow outside of body #66 0.5.0.dev3
* Created entrypoint hook allowing for users to extend marka with jinja
  exensions #60 0.5.0.dev2
* Moved to PEP 517 build #59 0.5.0.dev1
* markata.plugins.redirects will create redirect html files as a backup when
  server-side redirects fail, or the user does not have the ability to issue. #76
* create a slugify migration script #82
* DeepMerge `config_overrides` with config in post render methods #91 0.5.0.dev13

### sluggify paths

`python-sluggify` was implemented to ensure good urls are in place despite the
name of the original file.

For examples of how `python-slugify` will change your url's see the
[project's home page](https://pypi.org/project/python-slugify/).  One
difference is that `markata` will leave `/`'s for routing in the slugs.

#### OPTING OUT

do not want to use slugify, you can opt out by setting `slugify=False` in your
`markata.toml`.
If you have an existing site and do not want to implement redirects, or simply

``` toml
[markata]
slugify=false
```

#### Migrating to slugify

From the command line with `markata>=0.5.0` installed run the
migration script from the command line to create a redirects file in the
default location.  This should avoid all 404's as it will create a redirects
file that many static hosting providers will issue a server-side 301 for, and
for those that don't, markata.plugins.redirects creates a redirect html page,
that will kick in as a backup.

``` bash
python -m markata.scripts.migrate_to_slugify
```

### configurable page template

Now injects seo into the default template through configuration.  Here is an
example, by adding this to your `markata.toml` configuration.

``` toml
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

``` html
<meta name="og:type" content="article">
<meta name="og:author" content="Waylon Walker">
<meta name="og:site_name" content="Waylon Walker">
<meta name="theme-color" content="#322D39">
<meta name="twitter:creator" content="@_waylonwalker">
```

You can have an array of toml tables with a key of text.  The text will be
added as plain text to the end of the head of each page.

``` toml
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

``` html
  <meta name="description" content="{{ description }}">
```

## Config Overrides

Each post can override config settings such as `head`.  New meta tags can be
added to a single post, or anything that your template might reference from
config.

``` yaml
config_overrides:
  head:
    meta:
    - content: waylonwalker
      name: author
    - content: '@_waylonwalker'
      name: 'twitter:creator'
    link:
    - href: https://waylonwalker.com/that-special-post/
      rel: canonical
    text:
    - value: <link rel='stylesheet' href='/my-extra-styles.css' />
```

## 0.4.1

* Issue FutureWarning for upcoming change to slugify change that will change urls

To keep existing behavior add this to your `markata.toml`.

``` toml
[markata]
slugify=false
```

## 0.4.0

* feat: add html logging with [setup_logging](/markata/plugins/setup_logging/)
  plugin is all new closes #37
* fix: remove HTML tidy as the site generator tag
* feat: create configurable [navbar](https://markata.dev/nav)
* perf: prevent double runs on pre-render and post-render #39
* perf: prevent duplicate ruun from to_dict calling pre-render #53
   * `to_dict` only runs up to `render` phase if necessary as directed by `register_attr`
* perf: only prettify if configured #54
* fix: pyinstrument will not create a second profiler causing it to end in errors #50
* fix: sites without feeds config do not create an index #55

### Double Runs

Previously markata would catch AttributeError and run the previous step any
time you ran a step too early.  The way this was implemented caused some steps
such as pre-render and post-render to run twice with every single run.

This change will no longer catch attribute errors. If you run into any issues
with your plugins not running before asking for attributes created by your
plugin make sure that you implement the
[@register_attr](https://markata.dev/markata/hookspec/#register_attr-function)
decorator.

### Prettify

prettify html has been turned off by default as beautifulsoup4 prettify was
taking a significant time, and was often popping up as the slowest parts in my
personal `_profile`.  If you want to continue running prettify throughout the
build you can set a flag in your config to continue running prettify.

``` toml
[markata]
prettify_html = true
```

## 0.3.0

Skipped from a bump2version misconfiguration.

## 0.2.0

* feat: [auto_description](/markata/plugins/auto_description/) plugin is all
  new closes #13
* deprecated: long_description has been deprecated by auto_description
* fix: [covers](/markata/plugins/covers/) plugin which would previously skip
  every time.
* feat: [`markata clean`](/markata/plugins/base_cli/#clean-function) cleans up
  your cache and output from the command line
* fix: [`publish_source`](/markata/plugins/publish_source/) plugin will now
  ignore any non yaml serializable values 
* feat: Default template colors are now customizable
* feat: Default template now has light and dark theme
* feat: map now has the ability to map entire posts
* feat: [prevnext](/markata/plugins/prevnext/) plugin was added to link between
  posts closes #20
* feat: [jinja_md](/markata/plugins/jinja_md/) plugins was added to incorporate
  jinja into all the markdown
* breaking: [feeds](/markata/plugins/feeds) config now has feeds and
  feeds_config
* feat: `output_html` can now be specified in the frontmatter [see
  example](/markata/plugins/publish_html/#explicityly-set-the-output)
* feat: edit link is now included in the default page template closes #21

### breaking change to feeds config

If you are using the 0.1.0 version of feeds, and have configured custom
templates in `markata.feeds.template` and `markata.feeds.card_template` they
will need to be moved to `markata.feeds_config`.

Here is what you need to do to update your feeds_config.

``` diff
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

``` python
m = Markata()
# 'post' will return the entire post
m.map('post', filter='"git" in tags')
```

### Customizable colors

``` toml 
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

* Cache is busted on plugin change 
* plugin is configurable
* plugin now has docs

### auto_description Configuration

Open up your `markata.toml` file and add new entries for your
auto_descriptions.  You can have multiple desriptions, each one will be named
after the key you give it in your config.

``` toml
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

4e299d6 fixes the dedupe issue that was in develop for awhile and closes #24

## 0.1.0

* fix: pyinstument plugin no longer overrides the main cli callback
* feat: default is to run the profiler if pyinstrument is installed
* fix: --profile is now under the build command
* feat: --pretty/--no-pretty will make pretty tracebacks, and skip over framework code (closes #4 )
* fix: links are now absolute so they work from github, thanks [MR Destructive](https://github.com/Mr-Destructive)
* deprecate: `article['content_hash']` has been removed with preference for simply hashing on `article['content']`
* create `heading_link` plugin
* doc: How to create your home page. [docs](/home-page/)

### New cli help

After the pyinstrument plugin was fixed --version and  --to-json are back, and --profile is now under the build command.

![image](https://user-images.githubusercontent.com/22648375/150662983-547aebbd-c18c-4c17-8985-a6dc01cd29c7.png)

### New Heading Link Plugin

The new heading link plugin makes it easier to share the exact part of an article you want with someone, by giving clickable links to the id of the heading it's next to.

![image](https://user-images.githubusercontent.com/22648375/151718782-08a7cd26-41c1-4f00-a12c-0a208c593e9c.png)

## 0.0.1

Initial Release 🎉

