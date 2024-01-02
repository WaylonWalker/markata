---
content: "# Markata Changelog\n\n## 0.8.0\n\n- pydantic support\n\n### Pydantic Support\n\nNow
  plugins are configured through a pydantic Config object.\n\n### breaking changes\n\nThere
  are a number of breaking changes going into 0.8.0. Use caution when\nupgrading.\n\n####
  glob config is now under markata.glob\n\n```diff\n- [markata]\n- glob_patterns =
  \"pages/**/*.md\"\n+ [markata.glob]\n+ glob_patterns = \"pages/**/*.md\"\n```\n\n####
  Feeds are now a list\n\n```toml\n[markata.feeds.published]\ntemplate=\"pages/templates/archive_template.html\"\ncard_template
  = \"pages/templates/feed_card.html\"\nfilter=\"date<=today and templateKey in ['blog-post',
  'til'] and status.lower()=='published'\"\nsort=\"date\"\n```\n\n> old\n\n```toml\n[[markata.feeds.published]]\ntemplate=\"pages/templates/archive_template.html\"\ncard_template
  = \"pages/templates/feed_card.html\"\nfilter=\"date<=today and templateKey in ['blog-post',
  'til'] and status.lower()=='published'\"\nsort=\"date\"\n```\n\n> new\n\n### markata.summary.filter_count
  is now a list\n\nThe old way was to set up a dict, where the keys were the name,
  now its a list\nof Objects with an explicit name field.\n\n```toml\n[markata.summary.filter_count.drafts]\nfilter=\"published
  == 'False'\"\ncolor='red'\n```\n\n> Old\n\n```toml\n[[markata.summary.filter_count]]\nname='drafts'\nfilter=\"published
  == 'False'\"\ncolor='red'\n```\n\n## 0.7.4\n\n- Fix: Icon resize broken from PIL
  10.0.0 release\n\n## 0.7.3\n\n- Fix: broken default feed card template\n- Fix: broken
  links for index.md in feeds\n\n## 0.7.2\n\n- Fix: broken `markata new` command due
  to pydantic v2 compatability with copier.\n\n## 0.7.0\n\n- Adopt ruff linter 0.7.0.dev1
  #142\n- add support for markdown_it_py backend #145\n- add trogon tui support 0.7.0.dev6\n\n###
  trogon tui support\n\ninstall with pip\n\n```bash\npip install 'markata[tui]'\n```\n\nusage\n\n```bash\nmarkata
  tui\n```\n\n## 0.6.4\n\n- Fix: Icon resize broken from PIL 10.0.0 release\n\n##
  0.6.3\n\n- Fix: broken `markata new` command due to pydantic v2 compatability with
  copier.\n\n## 0.6.2\n\nUpdate License and Security files.\n\n## 0.6.1\n\n- Fix:
  allow feeds to be used from within Markdown.\n\n### Feeds in Markdown\n\n```markdown\n{%
  for post in markata.feeds.docs.posts %}\n[{{post.title}}](/{{post.slug}})\n{% endfor
  %}\n```\n\n## 0.6.0\n\n- Fix: article_html is not available to jinja 0.6.0.dev1
  #105\n- Fix: service worker is not upating client without hard refresh 0.6.0.dev2\n
  \ #106\n- Feat: create teardown lifecycle method 0.6.0.dev3 #110\n- Fix: implement
  teardown on all raises 0.6.0.dev4 #111\n- Fix: implement teardown in pyinstrument
  plugin 0.6.0.dev5 #112\n- Feat: Automatically call teardown without needing to remember
  it before raise\n  0.6.0.dev6 #113\n- Fix: only stop the profiler if it is running
  0.6.0.dev7 #114\n- Fix: map was giving inconsistent results 0.6.0.dev8 #116\n- Fix:
  tui continuously rebuilds if an input file exists in a parent directory\n  to the
  output directory 0.6.0.dev9 #118\n- Feat: report the cache stats for the current
  run 0.6.0.dev11 #121\n- Fix: prevent zerodivisionerror when reporting stats 0.6.0.dev12
  #122\n- Fix: properly set the pyinstrument profiler to prevent recurrsion errors\n
  \ 0.6.dev13 #123\n- Clean: cli attributes (`runner`, `summary`, `server`, `plugins`)
  are now\n  added as Markata properties through `register_attr` rather than directly
  to\n  the class 0.6.0.dev13 #107\n- Fix: Markata tui will remain running even when
  the runner fails 0.6.0.dev13\n  #107\n- Fix: Pinned to `textual<0.2.0` due to breaking
  changes 0.6.0.dev13 #107\n- Fix: snyk remove seuptools from requirements 0.6.0.dev14
  #130\n- Fix: Markata.filter was missing post and m 0.6.0.dev14 #133\n- Feat: add
  `path_prefix` config for gh-pages deploy 0.6.0.dev15 #132 fixes #57\n- Feat: created
  `markata plugin show` cli command 0.6.0.dev16 #109\n- Feat: use `path_prefix` in
  nav entries 0.6.0.dev17 #136\n- Feat: enable wikilinks extension by default 0.6.0.dev18
  #138\n- Feat: `markata` instance is available form inside card_template for feeds\n
  \ 0.6.0.dev19 #139\n\n### wikilinks\n\nwikilinks are now enabled by default ex:
  `[[home-page]]`. This will create a\nlink `<a class=\"wikilink\" href=\"/home-page/\">home-page</a>`.
  This will\nautomagically work if you leave `markata.plugins.flat_slug` plugin enabled\n(which
  is by default).\n\n> ProTip: this was highly inspired by the\n> [marksman-lsp](https://github.com/artempyanykh/marksman)
  by\n> [artempyanykh](https://github.com/artempyanykh/), which can autocomplete post\n>
  links in this style for you.\n\n[[home-page]]\n\n## 0.5.5\n\n- Fix: Icon resize
  broken from PIL 10.0.0 release\n\n## 0.5.4\n\n- Fix: broken `markata new` command
  due to pydantic v2 compatability with copier.\n\n## 0.5.2\n\n- clean up unnecessary
  images_url is missing warning #104 0.5.2.dev1\n\n## 0.5.1\n\n- fix: contrast ratio
  on admonitions was insufficient for A11y #103 0.5.1.dev0\n\n## 0.5.0\n\n- Create
  `new` cli command for creating new `blogs`, `posts`, and `plugins` #93\n  0.5.0.dev16
  [base_cli-docs](https://markata.dev/markata/plugins/base_cli/)\n- Remove unused
  function clif that was the original entrypoint #81 0.5.0.dev8\n- Allow template
  variables to be used in head config #88 0.5.0.dev12\n- Expose `markata.__version__`
  to templates as `__version__` #89 0.5.0.dev13\n- Fix, ignore post_template save
  on posts without an html attribute #92 0.5.0.dev13\n- Fix #33 sluggify paths #69
  **BREAKING CHANGE** 0.5.0.dev6\n- Configurable template #70 0.5.dev5, #85 0.5.0.dev11\n-
  Fix #40 Images overlfow outside of body #66 0.5.0.dev3\n- Created entrypoint hook
  allowing for users to extend marka with jinja\n  exensions #60 0.5.0.dev2\n- Moved
  to PEP 517 build #59 0.5.0.dev1\n- new `markata.plugins.redirects` will create redirect
  html files as a backup when\n  server-side redirects fail. #76 0.5.0.dev10\n  [redirects-docs](https://markata.dev/markata/plugins/redirects/)\n-
  create a slugify migration script #82\n- DeepMerge `config_overrides` with config
  in post render methods #91 0.5.0.dev13\n- Create ipython extension to automatically
  load markata #79 0.5.0.dev15\n- Fix: images wrapped in a link overflow outside the
  body #96\n- new `markata.plugins.service_worker` plugin to create service workers
  and\n  enable offline mode on sites #94 0.5.0.dev15\n  [service-worker-docs](https://markata.dev/markata/plugins/service-worker/)\n-
  Fix: icons were relatively linked, and were broken for any page other than\n  index,
  they are now absolutely linked to the root of the site. #97 0.5.0.dev16\n- Fix:
  auto_descriptions were not rendered on first pass in tui or at all in\n  build due
  to auto_description running after jinja_md. #100 0.5.0.dev18\n- Fix: give redirect
  pages a uniqe description and title #101 0.5.0.dev19\n\n### `new` cli command\n\nMore
  information in the [base_cli-docs](https://markata.dev/markata/plugins/base_cli/).\n\n```bash\n#
  create a new blog template\n# copier requires you to specify a directory\nmarkata
  new blog [directory]\n\n# create a new blog post\nmarkata new post\n\n# create a
  new plugin\nmarkata new plugin\n\nmarkata new --help\n\n Usage: markata new [OPTIONS]
  COMMAND [ARGS]...\n\n create new things from templates\n\n\u256D\u2500 Options \u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256E\n\u2502
  --help          Show this message and exit.                                                                                                                                       \u2502\n\u2570\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256F\n\u256D\u2500
  Commands \u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256E\n\u2502
  blog       Create a new blog from using the template from https://github.com/WaylonWalker/markata-blog-starter.
  \                                                                  \u2502\n\u2502
  plugin     Create a new plugin using the template at https://github.com/WaylonWalker/markata-plugin-template.
  \                                                                    \u2502\n\u2502
  post       Create new blog post in the pages directory from the template at  https://github.com/WaylonWalker/markata-post-template.
  \                                              \u2502\n\u2570\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256F\n```\n\n###
  sluggify paths\n\n`python-sluggify` was implemented to ensure good urls are in place
  despite the\nname of the original file.\n\nFor examples of how `python-slugify`
  will change your url's see the\n[project's home page](https://pypi.org/project/python-slugify/).
  One\ndifference is that `markata` will leave `/`'s for routing in the slugs.\n\n####
  OPTING OUT\n\nIf you have an existing site and do not want to implement redirects
  or if you\ndo not want to use slugify, you can opt out by setting `slugify=False`
  in your\n`markata.toml`.\n\n```toml\n[markata]\nslugify=false\n```\n\n#### Migrating
  to slugify\n\nFrom the command line with `markata>=0.5.0` installed run the\nmigration
  script from the command line to create a redirects file in the\ndefault location.
  This should avoid all 404's as it will create a redirects\nfile that many static
  hosting providers will issue a server-side 301 for, and\nfor those that don't, markata.plugins.redirects
  creates a redirect html page,\nthat will kick in as a backup.\n\n```bash\npython
  -m markata.scripts.migrate_to_slugify\n```\n\n### configurable page template\n\nNow
  injects seo into the default template through configuration. Here is an\nexample,
  by adding this to your `markata.toml` configuration.\n\n```toml\n[[markata.head.meta]]\nname
  = \"og:type\"\ncontent = \"article\"\n\n[[markata.head.meta]]\nname = \"og:author\"\ncontent
  = \"Waylon Walker\"\n\n[[markata.head.meta]]\nname = \"og:site_name\"\ncontent =
  \"Waylon Walker\"\n\n[[markata.head.meta]]\nname = \"theme-color\"\ncontent=\"#322D39\"\n\n[[markata.head.meta]]\nname
  = \"twitter:creator\"\ncontent=\"@_waylonwalker\"\n```\n\nYou will end up with these
  meta tags in your html.\n\n```html\n<meta name=\"og:type\" content=\"article\" />\n<meta
  name=\"og:author\" content=\"Waylon Walker\" />\n<meta name=\"og:site_name\" content=\"Waylon
  Walker\" />\n<meta name=\"theme-color\" content=\"#322D39\" />\n<meta name=\"twitter:creator\"
  content=\"@_waylonwalker\" />\n```\n\nYou can have an array of toml tables with
  a key of text. The text will be\nadded as plain text to the end of the head of each
  page.\n\n```toml\n[[markata.head]]\ntext = \"\"\"\n<style>\nimg {\nwidth: 100%;\nheight:
  auto;\n}\nul {\n  display: flex;\n  flex-wrap: wrap;\n}\n\nli {\n  flex: 1 2 400px;\n}\n</style>\n\n\"\"\"\n```\n\nDescriptions
  will now properly end up in each page.\n\n```html\n<meta name=\"description\" content=\"{{
  description }}\" />\n```\n\n## Config Overrides\n\nEach post can override config
  settings such as `head`. New meta tags can be\nadded to a single post, or anything
  that your template might reference from\nconfig.\n\n```yaml\nconfig_overrides:\n
  \ head:\n    meta:\n      - content: waylonwalker\n        name: author\n      -
  content: \"@_waylonwalker\"\n        name: \"twitter:creator\"\n    link:\n      -
  href: https://waylonwalker.com/that-special-post/\n        rel: canonical\n    text:\n
  \     - value: <link rel='stylesheet' href='/my-extra-styles.css' />\n```\n\n##
  ipython extension\n\nMarkata has an ipython extension if you want ipython to automatically
  load with\nan instance of `Markata` mapped to `m` and `markata` you can add the
  following\nto your `~/.ipython/profile_default/ipython_config.py`\n\n```python\nc.InteractiveShellApp.extensions.append('markata')\n```\n\n##
  0.4.1\n\n- Issue FutureWarning for upcoming change to slugify change that will change
  urls\n\nTo keep existing behavior add this to your `markata.toml`.\n\n```toml\n[markata]\nslugify=false\n```\n\n##
  0.4.0\n\n- feat: add html logging with [setup_logging](/markata/plugins/setup_logging/)\n
  \ plugin is all new closes #37\n- fix: remove HTML tidy as the site generator tag\n-
  feat: create configurable [navbar](https://markata.dev/nav)\n- perf: prevent double
  runs on pre-render and post-render #39\n- perf: prevent duplicate ruun from to_dict
  calling pre-render #53\n  - `to_dict` only runs up to `render` phase if necessary
  as directed by `register_attr`\n- perf: only prettify if configured #54\n- fix:
  pyinstrument will not create a second profiler causing it to end in\n  errors #50\n-
  fix: sites without feeds config do not create an index #55\n\n### Double Runs\n\nPreviously
  markata would catch AttributeError and run the previous step any\ntime you ran a
  step too early. The way this was implemented caused some steps\nsuch as pre-render
  and post-render to run twice with every single run.\n\nThis change will no longer
  catch attribute errors. If you run into any issues\nwith your plugins not running
  before asking for attributes created by your\nplugin make sure that you implement
  the\n[@register_attr](https://markata.dev/markata/hookspec/#register_attr-function)\ndecorator.\n\n###
  Prettify\n\nprettify html has been turned off by default as beautifulsoup4 prettify
  was\ntaking a significant time, and was often popping up as the slowest parts in
  my\npersonal `_profile`. If you want to continue running prettify throughout the\nbuild
  you can set a flag in your config to continue running prettify.\n\n```toml\n[markata]\nprettify_html
  = true\n```\n\n## 0.3.0\n\nSkipped from a bump2version misconfiguration.\n\n## 0.2.0\n\n-
  feat: [auto_description](/markata/plugins/auto_description/) plugin is all\n  new
  closes #13\n- deprecated: long_description has been deprecated by auto_description\n-
  fix: [covers](/markata/plugins/covers/) plugin would previously skip\n  every time.\n-
  feat: [`markata clean`](/markata/plugins/base_cli/#clean-function) cleans up\n  your
  cache and output from the command line\n- fix: [`publish_source`](/markata/plugins/publish_source/)
  plugin will now\n  ignore any non yaml serializable values\n- feat: Default template
  colors are now customizable\n- feat: Default template now has light and dark theme\n-
  feat: map now can map entire posts\n- feat: [prevnext](/markata/plugins/prevnext/)
  plugin was added to link between\n  posts closes #20\n- feat: [jinja_md](/markata/plugins/jinja_md/)
  plugins was added to incorporate\n  jinja into all the markdown\n- breaking: [feeds](/markata/plugins/feeds)
  config now has feeds and\n  feeds_config\n- feat: `output_html` can now be specified
  in the frontmatter\n  [see example](/markata/plugins/publish_html/#explicityly-set-the-output)\n-
  feat: edit link is now included in the default page template closes #21\n\n### breaking
  change to feeds config\n\nIf you are using the 0.1.0 version of feeds, and have
  configured custom\ntemplates in `markata.feeds.template` and `markata.feeds.card_template`
  they\nwill need to be moved to `markata.feeds_config`.\n\nHere is what you need
  to do to update your feeds_config.\n\n```diff\n+ [markata.feeds_config]\n+ template=\"pages/templates/archive_template.html\"\n+
  card_template=\"plugins/feed_card_template.html\"\n- [markata.feeds]\n- template=\"pages/templates/archive_template.html\"\n-
  card_template=\"plugins/feed_card_template.html\"\n```\n\n`markata.feeds` will only
  be used to configure feeds pages.\n\n### map entire posts\n\n`post` is now exposed
  to the `markata.map` object, allowing you to return a\nlist of posts.\n\n```python\nm
  = Markata()\n# 'post' will return the entire post\nm.map('post', filter='\"git\"
  in tags')\n```\n\n### Customizable colors\n\n```toml\n[markata]\n# default dark
  theme\ncolor_bg = '#1f2022'\ncolor_bg_code = '#1f2022'\ncolor_text = '#eefbfe'\ncolor_link
  = '#47cbff'\ncolor_accent = '#e1bd00c9'\noverlay_brightness = '.85'\n\n# pink and
  purple\ncolor_bg = 'deeppink'\ncolor_bg_code = 'rebeccapurple'\ncolor_text = 'white'\ncolor_link
  = 'aqua'\ncolor_accent = 'peachpuff'\noverlay_brightness = '1.2'\n\n# default light
  theme\ncolor_bg_light = '#eefbfe'\ncolor_bg_code_light = '#eefbfe'\ncolor_text_light
  = '#1f2022'\ncolor_link_light = '#47cbff'\ncolor_accent_light = '#ffeb00'\noverlay_brightness_light
  = '.95'\n```\n\n### All New auto_description plugin\n\n- Cache is busted on plugin
  change\n- plugin is configurable\n- plugin now has docs\n\n### auto_description
  Configuration\n\nOpen up your `markata.toml` file and add new entries for your\nauto_descriptions.
  You can have multiple desriptions, each one will be named\nafter the key you give
  it in your config.\n\n```toml\n[markata]\nhooks=[\n   \"markata.plugins.auto_description\",\n
  \  ]\n\n[markata.auto_description.description]\nlen=160\n[markata.auto_description.long_description]\nlen=250\n[markata.auto_description.super_description]\nlen=500\n```\n\nIn
  the above we will end up with three different descritpions,\n(`description`, `long_description`,
  and `super_description`) each will be the\nfirst number of characters from the document
  as specified in the config.\n\n### auto_descriptions are no longer duplicated\n\n4e299d6
  fixes the dedupe issue that was in develop for a while and closes #24\n\n## 0.1.0\n\n-
  fix: pyinstument plugin no longer overrides the main cli callback\n- feat: default
  is to run the profiler if pyinstrument is installed\n- fix: --profile is now under
  the build command\n- feat: --pretty/--no-pretty will make pretty tracebacks, and
  skip over\n  framework code (closes #4 )\n- fix: links are now absolute, so they
  work from GitHub, thanks\n  [MR Destructive](https://github.com/Mr-Destructive)\n-
  deprecate: `article['content_hash']` has been removed with preference for\n  hashing
  on `article['content']`\n- create `heading_link` plugin\n- doc: How to create your
  home page. [docs](/home-page/)\n\n### New cli help\n\nAfter the pyinstrument plugin
  was fixed --version and --to-json are back, and\n--profile is now under the build
  command.\n\n![image](https://user-images.githubusercontent.com/22648375/150662983-547aebbd-c18c-4c17-8985-a6dc01cd29c7.png)\n\n###
  New Heading Link Plugin\n\nThe new heading link plugin makes it easier to share
  the exact part of an\narticle you want with someone, by giving clickable links to
  the id of the\nheading it's next to.\n\n![image](https://user-images.githubusercontent.com/22648375/151718782-08a7cd26-41c1-4f00-a12c-0a208c593e9c.png)\n\n##
  0.0.1\n\nInitial Release \U0001F389\n"
date: 0001-01-01
description: 'pydantic support Now plugins are configured through a pydantic Config
  object. There are a number of breaking changes going into 0.8.0. Use caution when
  old new '
html: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>Changelog</title>\n<meta
  charset=\"utf-8\"/>\n<meta content=\"width=device-width, initial-scale=1\" name=\"viewport\"/>\n<meta
  content=\"pydantic support Now plugins are configured through a pydantic Config
  object. There are a number of breaking changes going into 0.8.0. Use caution when
  old new \" name=\"description\"/>\n<link href=\"/static/favicon.ico\" rel=\"icon\"
  type=\"image/png\"/>\n<script>\n        function setTheme(theme) {\n            document.documentElement.setAttribute(\"data-theme\",
  theme);\n        }\n\n        function detectColorSchemeOnLoad() {\n            //local
  storage is used to override OS theme settings\n            if (localStorage.getItem(\"theme\"))
  {\n                if (localStorage.getItem(\"theme\") == \"dark\") {\n                    setTheme(\"dark\");\n
  \               } else if (localStorage.getItem(\"theme\") == \"light\") {\n                    setTheme(\"light\");\n
  \               }\n            } else if (!window.matchMedia) {\n                //matchMedia
  method not supported\n                setTheme(\"light\");\n                return
  false;\n            } else if (window.matchMedia(\"(prefers-color-scheme: dark)\").matches)
  {\n                //OS theme setting detected as dark\n                setTheme(\"dark\");\n
  \           } else {\n                setTheme(\"light\");\n            }\n        }\n
  \       detectColorSchemeOnLoad();\n        document.addEventListener(\n            \"DOMContentLoaded\",\n
  \           function () {\n                //identify the toggle switch HTML element\n
  \               const toggleSwitch = document.querySelector(\n                    '#theme-switch
  input[type=\"checkbox\"]',\n                );\n\n                //function that
  changes the theme, and sets a localStorage variable to track the theme between page
  loads\n                function switchTheme(e) {\n                    if (e.target.checked)
  {\n                        localStorage.setItem(\"theme\", \"dark\");\n                        document.documentElement.setAttribute(\"data-theme\",
  \"dark\");\n                        toggleSwitch.checked = true;\n                    }
  else {\n                        localStorage.setItem(\"theme\", \"light\");\n                        document.documentElement.setAttribute(\"data-theme\",
  \"light\");\n                        toggleSwitch.checked = false;\n                    }\n
  \               }\n\n                //listener for changing themes\n                toggleSwitch.addEventListener(\"change\",
  switchTheme, false);\n\n                //pre-check the dark-theme checkbox if dark-theme
  is set\n                if (document.documentElement.getAttribute(\"data-theme\")
  == \"dark\") {\n                    toggleSwitch.checked = true;\n                }\n
  \           },\n            false,\n        );\n    </script>\n<style>\n      :root
  {\n        --color-bg: #1f2022;\n        --color-bg-2: ;\n        --color-bg-code:
  #1f2022;\n        --color-text: #eefbfe;\n        --color-link: #fb30c4; \n        --color-accent:
  #e1bd00c9;\n        --overlay-brightness: .85;\n        --body-width: 800px;\n      }\n
  \     [data-theme=\"dark\"] {\n        --color-bg: #1f2022;\n        --color-bg-2:
  ;\n        --color-bg-code: #1f2022;\n        --color-text: #eefbfe;\n        --color-link:
  #fb30c4; \n        --color-accent: #e1bd00c9;\n        --overlay-brightness: .85;\n
  \       --body-width: 800px;\n      }\n      [data-theme=\"light\"] {\n        --color-bg:
  #eefbfe;\n        --color-bg-2: ;\n        --color-bg-code: #eefbfe;\n        --color-text:
  #1f2022;\n        --color-link: #fb30c4; \n        --color-accent: #ffeb00;\n        --overlay-brightness:
  .95;\n      }\n\n        html {\n            font-family: \"Space Mono\", monospace;\n
  \           background: var(--color-bg);\n            color: var(--color-text);\n
  \       }\n\n        a {\n            color: var(--color-link);\n        }\n\n        main
  a {\n            max-width: 100%;\n        }\n\n        .heading-permalink {\n            font-size:
  .7em;\n        }\n\n        body {\n            max-width: var(--body-width);\n
  \           margin: 5rem auto;\n            padding: 0 .5rem;\n            font-size:
  1rem;\n            line-height: 1.56;\n        }\n\n        blockquote {\n            background:
  var(--color-bg);\n            filter: brightness(var(--overlay-brightness));\n            border-left:
  4px solid var(--color-accent);\n            border-radius: 4px;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #f1fa8c,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n            padding-left: 1rem;\n            margin: 1rem;\n
  \       }\n\n        li.post {\n            list-style-type: None;\n            padding:
  .2rem 0;\n        }\n\n        pre.wrapper {\n            padding: 0;\n            box-shadow:
  0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n            display: flex;\n            flex-direction:
  column;\n            position: relative;\n            margin: 2rem;\n        }\n\n
  \       pre {\n            margin: 0;\n            padding: 1rem;\n            min-width:
  -webkit-fill-available;\n            max-width: fit-content;\n            overflow-x:
  auto;\n        }\n\n        pre .filepath {\n            margin: 0;\n            padding-left:
  1rem;\n            border-radius: 4px 4px 0 0;\n            background: black;\n
  \           display: flex;\n            justify-content: space-between;\n            align-items:
  center;\n        }\n\n        pre .filepath p {\n            margin: 0\n        }\n\n
  \       pre .filepath .right {\n            display: flex;\n            gap: .2rem;\n
  \           align-items: center;\n        }\n\n        pre::-webkit-scrollbar {\n
  \           height: 4px;\n            background-color: transparent;\n        }\n\n
  \       pre::-webkit-scrollbar-thumb {\n            background-color: #d3d3d32e;\n
  \           border-radius: 2px;\n        }\n\n        pre::-webkit-scrollbar-track
  {\n            background-color: transparent;\n        }\n\n        .copy-wrapper
  {\n            background: none;\n            position: absolute;\n            width:
  100%;\n            z-index: 100;\n            display: flex;\n            justify-content:
  flex-end;\n        }\n\n        button.copy {\n            z-index: 100;\n            background:
  none;\n            fill: #ffffff45;\n            border: none;\n            width:
  32px;\n            align-self: flex-end;\n            top: 0;\n            right:
  0;\n            margin: 0.5rem 0.2rem;\n\n        }\n\n        button.copy:hover
  {\n            fill: white\n        }\n\n        a.help {\n            fill: #ffffff45;\n
  \       }\n\n        a.help:hover {\n            fill: white;\n        }\n\n        a.help
  svg {\n            height: 24px;\n            width: 24px;\n        }\n\n        .highlight
  {\n            background: var(--color-bg-code);\n            color: var(--color-text);\n
  \           filter: brightness(var(--overlay-brightness));\n            border-radius:
  0 0 4px 4px;\n        }\n\n        .highlight .c {\n            color: #8b8b8b\n
  \       }\n\n        /* Comment */\n        .highlight .err {\n            color:
  #960050;\n            background-color: #1e0010\n        }\n\n        /* Error */\n
  \       .highlight .k {\n            color: #c678dd\n        }\n\n        /* Keyword
  */\n        .highlight .l {\n            color: #ae81ff\n        }\n\n        /*
  Literal */\n        .highlight .n {\n            color: #abb2bf\n        }\n\n        /*
  Name */\n        .highlight .o {\n            color: #c678dd\n        }\n\n        /*
  Operator */\n        .highlight .p {\n            color: #abb2bf\n        }\n\n
  \       /* Punctuation */\n        .highlight .ch {\n            color: #8b8b8b\n
  \       }\n\n        /* Comment.Hashbang */\n        .highlight .cm {\n            color:
  #8b8b8b\n        }\n\n        /* Comment.Multiline */\n        .highlight .cp {\n
  \           color: #8b8b8b\n        }\n\n        /* Comment.Preproc */\n        .highlight
  .cpf {\n            color: #8b8b8b\n        }\n\n        /* Comment.PreprocFile
  */\n        .highlight .c1 {\n            color: #8b8b8b\n        }\n\n        /*
  Comment.Single */\n        .highlight .cs {\n            color: #8b8b8b\n        }\n\n
  \       /* Comment.Special */\n        .highlight .gd {\n            color: #c678dd\n
  \       }\n\n        /* Generic.Deleted */\n        .highlight .ge {\n            font-style:
  italic\n        }\n\n        /* Generic.Emph */\n        .highlight .gi {\n            color:
  #a6e22e\n        }\n\n        /* Generic.Inserted */\n        .highlight .gs {\n
  \           font-weight: bold\n        }\n\n        /* Generic.Strong */\n        .highlight
  .gu {\n            color: #8b8b8b\n        }\n\n        /* Generic.Subheading */\n
  \       .highlight .kc {\n            color: #c678dd\n        }\n\n        /* Keyword.Constant
  */\n        .highlight .kd {\n            color: #c678dd\n        }\n\n        /*
  Keyword.Declaration */\n        .highlight .kn {\n            color: #c678dd\n        }\n\n
  \       /* Keyword.Namespace */\n        .highlight .kp {\n            color: #c678dd\n
  \       }\n\n        /* Keyword.Pseudo */\n        .highlight .kr {\n            color:
  #c678dd\n        }\n\n        /* Keyword.Reserved */\n        .highlight .kt {\n
  \           color: #c678dd\n        }\n\n        /* Keyword.Type */\n        .highlight
  .ld {\n            color: #e6db74\n        }\n\n        /* Literal.Date */\n        .highlight
  .m {\n            color: #ae81ff\n        }\n\n        /* Literal.Number */\n        .highlight
  .s {\n            color: #e6db74\n        }\n\n        /* Literal.String */\n        .highlight
  .na {\n            color: #a6e22e\n        }\n\n        /* Name.Attribute */\n        .highlight
  .nb {\n            color: #98c379\n        }\n\n        /* Name.Builtin */\n        .highlight
  .nc {\n            color: #abb2bf\n        }\n\n        /* Name.Class */\n        .highlight
  .no {\n            color: #c678dd\n        }\n\n        /* Name.Constant */\n        .highlight
  .nd {\n            color: #abb2bf\n        }\n\n        /* Name.Decorator */\n        .highlight
  .ni {\n            color: #abb2bf\n        }\n\n        /* Name.Entity */\n        .highlight
  .ne {\n            color: #a6e22e\n        }\n\n        /* Name.Exception */\n        .highlight
  .nf {\n            color: #61afef\n        }\n\n        /* Name.Function */\n        .highlight
  .nl {\n            color: #abb2bf\n        }\n\n        /* Name.Label */\n        .highlight
  .nn {\n            color: #abb2bf\n        }\n\n        /* Name.Namespace */\n        .highlight
  .nx {\n            color: #a6e22e\n        }\n\n        /* Name.Other */\n        .highlight
  .py {\n            color: #abb2bf\n        }\n\n        /* Name.Property */\n        .highlight
  .nt {\n            color: #c678dd\n        }\n\n        /* Name.Tag */\n        .highlight
  .nv {\n            color: #abb2bf\n        }\n\n        /* Name.Variable */\n        .highlight
  .ow {\n            color: #c678dd\n        }\n\n        /* Operator.Word */\n        .highlight
  .w {\n            color: #abb2bf\n        }\n\n        /* Text.Whitespace */\n        .highlight
  .mb {\n            color: #ae81ff\n        }\n\n        /* Literal.Number.Bin */\n
  \       .highlight .mf {\n            color: #ae81ff\n        }\n\n        /* Literal.Number.Float
  */\n        .highlight .mh {\n            color: #ae81ff\n        }\n\n        /*
  Literal.Number.Hex */\n        .highlight .mi {\n            color: #ae81ff\n        }\n\n
  \       /* Literal.Number.Integer */\n        .highlight .mo {\n            color:
  #ae81ff\n        }\n\n        /* Literal.Number.Oct */\n        .highlight .sa {\n
  \           color: #e6db74\n        }\n\n        /* Literal.String.Affix */\n        .highlight
  .sb {\n            color: #e6db74\n        }\n\n        /* Literal.String.Backtick
  */\n        .highlight .sc {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Char */\n        .highlight .dl {\n            color: #e6db74\n        }\n\n
  \       /* Literal.String.Delimiter */\n        .highlight .sd {\n            color:
  #98c379\n        }\n\n        /* Literal.String.Doc */\n        .highlight .s2 {\n
  \           color: #98c379\n        }\n\n        /* Literal.String.Double */\n        .highlight
  .se {\n            color: #ae81ff\n        }\n\n        /* Literal.String.Escape
  */\n        .highlight .sh {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Heredoc */\n        .highlight .si {\n            color: #e6db74\n
  \       }\n\n        /* Literal.String.Interpol */\n        .highlight .sx {\n            color:
  #e6db74\n        }\n\n        /* Literal.String.Other */\n        .highlight .sr
  {\n            color: #e6db74\n        }\n\n        /* Literal.String.Regex */\n
  \       .highlight .s1 {\n            color: #e6db74\n        }\n\n        /* Literal.String.Single
  */\n        .highlight .ss {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Symbol */\n        .highlight .bp {\n            color: #abb2bf\n
  \       }\n\n        /* Name.Builtin.Pseudo */\n        .highlight .fm {\n            color:
  #61afef\n        }\n\n        /* Name.Function.Magic */\n        .highlight .vc
  {\n            color: #abb2bf\n        }\n\n        /* Name.Variable.Class */\n
  \       .highlight .vg {\n            color: #abb2bf\n        }\n\n        /* Name.Variable.Global
  */\n        .highlight .vi {\n            color: #abb2bf\n        }\n\n        /*
  Name.Variable.Instance */\n        .highlight .vm {\n            color: #abb2bf\n
  \       }\n\n        /* Name.Variable.Magic */\n        .highlight .il {\n            color:
  #ae81ff\n        }\n\n        /* Literal.Number.Integer.Long */\n\n        /* Tab
  style starts here */\n        .tabbed-set {\n            position: relative;\n            display:
  flex;\n            flex-wrap: wrap;\n            margin: 1em 0;\n            border-radius:
  0.1rem;\n        }\n\n        .tabbed-set>input {\n            display: none;\n
  \       }\n\n        .tabbed-set label {\n            width: auto;\n            padding:
  0.9375em 1.25em 0.78125em;\n            font-weight: 700;\n            font-size:
  0.84em;\n            white-space: nowrap;\n            border-bottom: 0.15rem solid
  transparent;\n            border-top-left-radius: 0.1rem;\n            border-top-right-radius:
  0.1rem;\n            cursor: pointer;\n            transition: background-color
  250ms, color 250ms;\n        }\n\n        .tabbed-set .tabbed-content {\n            width:
  100%;\n            display: none;\n            box-shadow: 0 -.05rem #ddd;\n        }\n\n
  \       .tabbed-set input {\n            position: absolute;\n            opacity:
  0;\n        }\n\n        /* fonts */\n        h1 {\n            font-weight: 700;\n
  \       }\n\n        h1#title a {\n            font-size: 16px;\n        }\n\n        h1,\n
  \       h2,\n        h3,\n        h4,\n        h5,\n        h6 {\n            margin-top:
  3rem;\n        }\n\n        h1 {\n            font-size: 2.5em;\n            margin-top:
  5rem;\n        }\n\n        h2 {\n            font-size: 1.63rem;\n            margin-top:
  5rem;\n        }\n\n\n\n        p {\n            font-size: 21px;\n            font-style:
  normal;\n            font-variant: normal;\n            font-weight: 400;\n            line-height:
  1.5;\n        }\n\n        @media only screen and (max-width: 700px) {\n            p
  {\n                font-size: 18px;\n            }\n        }\n\n        @media
  only screen and (max-width: 600px) {\n            p {\n                font-size:
  16px;\n            }\n        }\n\n        @media only screen and (max-width: 500px)
  {\n            p {\n                font-size: 14px;\n            }\n        }\n\n
  \       @media only screen and (max-width: 400px) {\n            p {\n                font-size:
  12px;\n            }\n        }\n\n\n        pre {\n            font-style: normal;\n
  \           font-variant: normal;\n            font-weight: 400;\n            line-height:
  18.5714px;\n            */\n        }\n\n        a {\n            font-weight: 600;\n
  \           text-decoration-color: var(--color-accent);\n            color: var(--color-link);\n
  \           padding: .3rem .5rem;\n            display: inline-block;\n        }\n\n
  \       .admonition,\n        details {\n            box-shadow: 0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n            margin: 5rem 0;\n            border: 1px solid transparent;\n
  \           border-radius: 4px;\n            text-align: left;\n            padding:
  0;\n            border: 0;\n\n        }\n\n        .admonition {\n            padding-bottom:
  1rem;\n        }\n\n        details[open] {\n            padding-bottom: .5rem;\n
  \       }\n\n        .admonition p {\n            padding: .2rem .6rem;\n        }\n\n
  \       .admonition-title,\n        .details-title,\n        summary {\n            background:
  var(--color-bg-2);\n            padding: 0;\n            margin: 0;\n            position:
  sticky;\n            top: 0;\n            z-index: 10;\n        }\n\n        summary:hover
  {\n            cursor: pointer;\n        }\n\n        summary.admonition-title,\n
  \       summary.details-title {\n            padding: .5rem;\n            padding-left:
  1rem;\n        }\n\n        .note {\n            border-left: 4px solid #f1fa8c;\n
  \           box-shadow:\n                -0.8rem 0rem 1rem -1rem #f1fa8c,\n                0.2rem
  0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .note>.admonition-title {\n            border-bottom:
  1px solid #3c3d2d;\n        }\n\n        .abstract {\n            border-left: 4px
  solid #8be9fd;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #8be9fd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .abstract>.admonition-title
  {\n            border-bottom: 1px solid #2c3a3f;\n        }\n\n        .info {\n
  \           border-left: 4px solid;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #8bb0fd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .info>.admonition-title {\n            border-bottom: 1px solid #2c313f;\n
  \       }\n\n        .tip {\n            border-left: 4px solid #008080;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #008080,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .tip>.admonition-title {\n            border-bottom:
  1px solid #1b2a2b;\n        }\n\n        .success {\n            border-left: 4px
  solid #50fa7b;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #50fa7b,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .success>.admonition-title
  {\n            border-bottom: 1px solid #263e2b;\n        }\n\n        .question
  {\n            border-left: 4px solid #a7fcbd;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #a7fcbd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .question>.admonition-title {\n            border-bottom: 1px solid #303e35;\n
  \       }\n\n        .warning {\n            border-left: 4px solid #ffb86c;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #ffb86c,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .warning>.admonition-title {\n            border-bottom:
  1px solid #3f3328;\n        }\n\n        .failure {\n            border-left: 4px
  solid #b23b3b;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #b23b3b,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .failure>.admonition-title
  {\n            border-bottom: 1px solid #34201f;\n        }\n\n        .danger {\n
  \           border-left: 4px solid #ff5555;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #ff5555,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .danger>.admonition-title {\n            border-bottom: 1px solid #402523;\n
  \       }\n\n        .bug {\n            border-left: 4px solid #b2548a;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #b2548a,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .bug>.admonition-title {\n            border-bottom:
  1px solid #32232c;\n        }\n\n        .example {\n            border-left: 4px
  solid #bd93f9;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #bd93f9,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .example>.admonition-title
  {\n            border-bottom: 1px solid #332d3e;\n        }\n\n        .source {\n
  \           border-left: 4px solid #bd93f9;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #bd93f9,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .source>.admonition-title {\n            border-bottom: 1px solid #332d3e;\n
  \       }\n\n        .quote {\n            border-left: 4px solid #999;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #999,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .quote>.admonition-title {\n            border-bottom:
  1px solid #2d2e2f;\n        }\n\n        table {\n            margin: 1rem 0;\n
  \           border-collapse: collapse;\n            border-spacing: 0;\n            display:
  block;\n            max-width: -moz-fit-content;\n            max-width: fit-content;\n
  \           overflow-x: auto;\n            white-space: nowrap;\n        }\n\n        table
  thead th {\n            border: solid 1px var(--color-text);\n            padding:
  10px;\n            text-align: left;\n        }\n\n        table tbody td {\n            border:
  solid 1px var(--color-text);\n            padding: 10px;\n        }\n\n        .theme-switch
  {\n            z-index: 10;\n            display: inline-block;\n            height:
  34px;\n            position: relative;\n            width: 60px;\n\n            display:
  flex;\n            justify-content: flex-end;\n            margin-right: 1rem;\n
  \           margin-left: auto;\n            position: fixed;\n            right:
  1rem;\n            top: 1rem;\n        }\n\n        .theme-switch input {\n            display:
  none;\n\n        }\n\n        .slider {\n            background-color: #ccc;\n            bottom:
  0;\n            cursor: pointer;\n            left: 0;\n            position: absolute;\n
  \           right: 0;\n            top: 0;\n            transition: .4s;\n        }\n\n
  \       .slider:before {\n            background-color: #fff;\n            bottom:
  4px;\n            content: \"\";\n            height: 26px;\n            left: 4px;\n
  \           position: absolute;\n            transition: .4s;\n            width:
  26px;\n        }\n\n        input:checked+.slider {\n            background-color:
  #343434;\n        }\n\n        input:checked+.slider:before {\n            background-color:
  #848484;\n        }\n\n        input:checked+.slider:before {\n            transform:
  translateX(26px);\n        }\n\n        .slider.round {\n            border-radius:
  34px;\n        }\n\n        .slider.round:before {\n            border-radius: 50%;\n
  \       }\n\n        main p img {\n            width: 100%;\n            width:
  -moz-available;\n            width: -webkit-fill-available;\n            width:
  fill-available;\n        }\n\n        details>* {\n            margin: 1rem;\n        }\n\n
  \       .admonition>* {\n            margin: 1rem;\n        }\n\n        p.admonition-title,\n
  \       summary {\n            margin: 0;\n            padding-left: 1.2rem;\n        }\n\n
  \       .small {\n            font-size: .9rem;\n            color: #888;\n        }\n\n
  \       admonition+admonition {\n            margin-top: 20rem;\n        }\n\n        ::-webkit-scrollbar
  {\n            height: 12px;\n            background-color: transparent;\n        }\n\n
  \       ::-webkit-scrollbar-thumb {\n            background-color: #d3d3d32e;\n
  \           border-radius: 6px;\n        }\n\n        ::-webkit-scrollbar-track
  {\n            background-color: transparent;\n        }\n    </style>\n<script>\n
  \       if (\"serviceWorker\" in navigator) {\n            navigator.serviceWorker.register(\"/service-worker.js\");\n
  \           navigator.serviceWorker.addEventListener(\"controllerchange\", () =>
  {\n                console.log(\"new worker\");\n                window.location.reload();\n
  \           });\n        }\n    </script>\n<meta content=\"waylon@waylonwalker.com\"
  name=\"og:author_email\"/>\n<meta content=\"waylon@waylonwalker.com\" name=\"og:author_email\"/>\n<meta
  content=\"Waylon Walker\" name=\"og:author\" property=\"og:author\"/><meta content=\"waylon@waylonwalaker.com\"
  name=\"og:author_email\" property=\"og:author_email\"/><meta content=\"website\"
  name=\"og:type\" property=\"og:type\"/><meta content=\"pydantic support Now plugins
  are configured through a pydantic Config object. There are a number of breaking
  changes going into 0.8.0. Use caution when old new \" name=\"description\" property=\"description\"/><meta
  content=\"pydantic support Now plugins are configured through a pydantic Config
  object. There are a number of breaking changes going into 0.8.0. Use caution when
  old new \" name=\"og:description\" property=\"og:description\"/><meta content=\"pydantic
  support Now plugins are configured through a pydantic Config object. There are a
  number of breaking changes going into 0.8.0. Use caution when old new \" name=\"twitter:description\"
  property=\"twitter:description\"/><meta content=\"Changelog | Markata\" name=\"og:title\"
  property=\"og:title\"/><meta content=\"Changelog | Markata\" name=\"twitter:title\"
  property=\"twitter:title\"/><meta content=\"https://markata.dev//changelog-og.png\"
  name=\"og:image\" property=\"og:image\"/><meta content=\"https://markata.dev//changelog-og.png\"
  name=\"twitter:image\" property=\"twitter:image\"/><meta content=\"1600\" name=\"og:image:width\"
  property=\"og:image:width\"/><meta content=\"900\" name=\"og:image:width\" property=\"og:image:width\"/><meta
  content=\"summary_large_image\" name=\"twitter:card\" property=\"twitter:card\"/><meta
  content=\"Markata\" name=\"og:site_name\" property=\"og:site_name\"/><meta content=\"@_waylonwalker\"
  name=\"twitter:creator\" property=\"twitter:creator\"/><meta content=\"Changelog\"
  name=\"title\" property=\"title\"/><meta content=\"markata 0.8.0.dev6\" name=\"generator\"
  property=\"generator\"/><link href=\"https://markata.dev//changelog/\" rel=\"canonical\"/><meta
  content=\"https://markata.dev//changelog/\" name=\"og:url\" property=\"og:url\"/></head>\n<body><nav>\n<a
  href=\"/\">markata</a>\n<a href=\"https://github.com/WaylonWalker/markata\">GitHub</a>\n</nav>\n<div>\n<label
  class=\"theme-switch\" for=\"checkbox-theme\" id=\"theme-switch\" title=\"light/dark
  mode toggle\">\n<input id=\"checkbox-theme\" type=\"checkbox\"/>\n<div class=\"slider
  round\"></div>\n</label>\n</div>\n<section class=\"title\">\n<h1 id=\"title\">\n
  \           Changelog \n            \n        </h1>\n</section>\n<main><h1 id=\"markata-changelog\">Markata
  Changelog <a class=\"header-anchor\" href=\"#markata-changelog\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<h2 id=\"080\">0.8.0 <a class=\"header-anchor\"
  href=\"#080\"><svg aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\"
  focusable=\"false\" height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>pydantic support</li>\n</ul>\n<h3>Pydantic
  Support</h3>\n<p>Now plugins are configured through a pydantic Config object.</p>\n<h3>breaking
  changes</h3>\n<p>There are a number of breaking changes going into 0.8.0. Use caution
  when\nupgrading.</p>\n<h4>glob config is now under markata.glob</h4>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"gd\">- [markata]</span>\n<span
  class=\"gd\">- glob_patterns = \"pages/**/*.md\"</span>\n<span class=\"gi\">+ [markata.glob]</span>\n<span
  class=\"gi\">+ glob_patterns = \"pages/**/*.md\"</span>\n</pre></div>\n\n</pre>\n<h4>Feeds
  are now a list</h4>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.feeds.published]</span>\n<span
  class=\"n\">template</span><span class=\"o\">=</span><span class=\"s2\">\"pages/templates/archive_template.html\"</span>\n<span
  class=\"n\">card_template</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"pages/templates/feed_card.html\"</span>\n<span
  class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">\"date&lt;=today
  and templateKey in ['blog-post', 'til'] and status.lower()=='published'\"</span>\n<span
  class=\"n\">sort</span><span class=\"o\">=</span><span class=\"s2\">\"date\"</span>\n</pre></div>\n\n</pre>\n<blockquote>\n<p>old</p>\n</blockquote>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.feeds.published]]</span>\n<span
  class=\"n\">template</span><span class=\"o\">=</span><span class=\"s2\">\"pages/templates/archive_template.html\"</span>\n<span
  class=\"n\">card_template</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"pages/templates/feed_card.html\"</span>\n<span
  class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">\"date&lt;=today
  and templateKey in ['blog-post', 'til'] and status.lower()=='published'\"</span>\n<span
  class=\"n\">sort</span><span class=\"o\">=</span><span class=\"s2\">\"date\"</span>\n</pre></div>\n\n</pre>\n<blockquote>\n<p>new</p>\n</blockquote>\n<h3>markata.summary.filter_count
  is now a list</h3>\n<p>The old way was to set up a dict, where the keys were the
  name, now its a list\nof Objects with an explicit name field.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.summary.filter_count.drafts]</span>\n<span
  class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">\"published
  == 'False'\"</span>\n<span class=\"n\">color</span><span class=\"o\">=</span><span
  class=\"s1\">'red'</span>\n</pre></div>\n\n</pre>\n<blockquote>\n<p>Old</p>\n</blockquote>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.summary.filter_count]]</span>\n<span
  class=\"n\">name</span><span class=\"o\">=</span><span class=\"s1\">'drafts'</span>\n<span
  class=\"n\">filter</span><span class=\"o\">=</span><span class=\"s2\">\"published
  == 'False'\"</span>\n<span class=\"n\">color</span><span class=\"o\">=</span><span
  class=\"s1\">'red'</span>\n</pre></div>\n\n</pre>\n<h2 id=\"074\">0.7.4 <a class=\"header-anchor\"
  href=\"#074\"><svg aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\"
  focusable=\"false\" height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>Fix: Icon resize
  broken from PIL 10.0.0 release</li>\n</ul>\n<h2 id=\"073\">0.7.3 <a class=\"header-anchor\"
  href=\"#073\"><svg aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\"
  focusable=\"false\" height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>Fix: broken default
  feed card template</li>\n<li>Fix: broken links for <a href=\"http://index.md\">index.md</a>
  in feeds</li>\n</ul>\n<h2 id=\"072\">0.7.2 <a class=\"header-anchor\" href=\"#072\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>Fix: broken <code>markata
  new</code> command due to pydantic v2 compatability with copier.</li>\n</ul>\n<h2
  id=\"070\">0.7.0 <a class=\"header-anchor\" href=\"#070\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>Adopt ruff linter
  0.7.0.dev1 #142</li>\n<li>add support for markdown_it_py backend #145</li>\n<li>add
  trogon tui support 0.7.0.dev6</li>\n</ul>\n<h3>trogon tui support</h3>\n<p>install
  with pip</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>pip<span class=\"w\"> </span>install<span
  class=\"w\"> </span><span class=\"s1\">'markata[tui]'</span>\n</pre></div>\n\n</pre>\n<p>usage</p>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>markata<span class=\"w\"> </span>tui\n</pre></div>\n\n</pre>\n<h2
  id=\"064\">0.6.4 <a class=\"header-anchor\" href=\"#064\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>Fix: Icon resize
  broken from PIL 10.0.0 release</li>\n</ul>\n<h2 id=\"063\">0.6.3 <a class=\"header-anchor\"
  href=\"#063\"><svg aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\"
  focusable=\"false\" height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>Fix: broken <code>markata
  new</code> command due to pydantic v2 compatability with copier.</li>\n</ul>\n<h2
  id=\"062\">0.6.2 <a class=\"header-anchor\" href=\"#062\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Update License and Security
  files.</p>\n<h2 id=\"061\">0.6.1 <a class=\"header-anchor\" href=\"#061\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>Fix: allow feeds
  to be used from within Markdown.</li>\n</ul>\n<h3>Feeds in Markdown</h3>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>{% for post in markata.feeds.docs.posts
  %}\n[<span class=\"nt\">{{post.title}}</span>](<span class=\"na\">/{{post.slug}}</span>)\n{%
  endfor %}\n</pre></div>\n\n</pre>\n<h2 id=\"060\">0.6.0 <a class=\"header-anchor\"
  href=\"#060\"><svg aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\"
  focusable=\"false\" height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>Fix: article_html
  is not available to jinja 0.6.0.dev1 #105</li>\n<li>Fix: service worker is not upating
  client without hard refresh 0.6.0.dev2\n#106</li>\n<li>Feat: create teardown lifecycle
  method 0.6.0.dev3 #110</li>\n<li>Fix: implement teardown on all raises 0.6.0.dev4
  #111</li>\n<li>Fix: implement teardown in pyinstrument plugin 0.6.0.dev5 #112</li>\n<li>Feat:
  Automatically call teardown without needing to remember it before raise\n0.6.0.dev6
  #113</li>\n<li>Fix: only stop the profiler if it is running 0.6.0.dev7 #114</li>\n<li>Fix:
  map was giving inconsistent results 0.6.0.dev8 #116</li>\n<li>Fix: tui continuously
  rebuilds if an input file exists in a parent directory\nto the output directory
  0.6.0.dev9 #118</li>\n<li>Feat: report the cache stats for the current run 0.6.0.dev11
  #121</li>\n<li>Fix: prevent zerodivisionerror when reporting stats 0.6.0.dev12 #122</li>\n<li>Fix:
  properly set the pyinstrument profiler to prevent recurrsion errors\n0.6.dev13 #123</li>\n<li>Clean:
  cli attributes (<code>runner</code>, <code>summary</code>, <code>server</code>,
  <code>plugins</code>) are now\nadded as Markata properties through <code>register_attr</code>
  rather than directly to\nthe class 0.6.0.dev13 #107</li>\n<li>Fix: Markata tui will
  remain running even when the runner fails 0.6.0.dev13\n#107</li>\n<li>Fix: Pinned
  to <code>textual&lt;0.2.0</code> due to breaking changes 0.6.0.dev13 #107</li>\n<li>Fix:
  snyk remove seuptools from requirements 0.6.0.dev14 #130</li>\n<li>Fix: Markata.filter
  was missing post and m 0.6.0.dev14 #133</li>\n<li>Feat: add <code>path_prefix</code>
  config for gh-pages deploy 0.6.0.dev15 #132 fixes #57</li>\n<li>Feat: created <code>markata
  plugin show</code> cli command 0.6.0.dev16 #109</li>\n<li>Feat: use <code>path_prefix</code>
  in nav entries 0.6.0.dev17 #136</li>\n<li>Feat: enable wikilinks extension by default
  0.6.0.dev18 #138</li>\n<li>Feat: <code>markata</code> instance is available form
  inside card_template for feeds\n0.6.0.dev19 #139</li>\n</ul>\n<h3>wikilinks</h3>\n<p>wikilinks
  are now enabled by default ex: <code>[[home-page]]</code>. This will create a\nlink
  <code>&lt;a class=\"wikilink\" href=\"/home-page/\"&gt;home-page&lt;/a&gt;</code>.
  This will\nautomagically work if you leave <code>markata.plugins.flat_slug</code>
  plugin enabled\n(which is by default).</p>\n<blockquote>\n<p>ProTip: this was highly
  inspired by the\n<a href=\"https://github.com/artempyanykh/marksman\">marksman-lsp</a>
  by\n<a href=\"https://github.com/artempyanykh/\">artempyanykh</a>, which can autocomplete
  post\nlinks in this style for you.</p>\n</blockquote>\n<p><a class=\"wikilink\"
  href=\"/home-page\">home-page</a></p>\n<h2 id=\"055\">0.5.5 <a class=\"header-anchor\"
  href=\"#055\"><svg aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\"
  focusable=\"false\" height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>Fix: Icon resize
  broken from PIL 10.0.0 release</li>\n</ul>\n<h2 id=\"054\">0.5.4 <a class=\"header-anchor\"
  href=\"#054\"><svg aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\"
  focusable=\"false\" height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>Fix: broken <code>markata
  new</code> command due to pydantic v2 compatability with copier.</li>\n</ul>\n<h2
  id=\"052\">0.5.2 <a class=\"header-anchor\" href=\"#052\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>clean up unnecessary
  images_url is missing warning #104 0.5.2.dev1</li>\n</ul>\n<h2 id=\"051\">0.5.1
  <a class=\"header-anchor\" href=\"#051\"><svg aria-hidden=\"true\" class=\"heading-permalink\"
  fill=\"currentColor\" focusable=\"false\" height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\"
  xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949
  2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992 5.992 0 0 0 1.695-4.285 5.976
  5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242 6.003 6.003 0 0 0-1.905 1.24l-1.731
  1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985 3.985 0 0 1 2.761-1.123 3.975
  3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005 3.006a3.982 3.982 0
  0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201 13.6zm5.602-3.198a5.99
  5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995 2.994a5.992 5.992 0
  0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0 0 6.431 1.242 6.003
  6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836 19.81a3.985
  3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0 0 1-.111-5.644l3.005-3.006a3.982
  3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632 1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>fix:
  contrast ratio on admonitions was insufficient for A11y #103 0.5.1.dev0</li>\n</ul>\n<h2
  id=\"050\">0.5.0 <a class=\"header-anchor\" href=\"#050\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>Create <code>new</code>
  cli command for creating new <code>blogs</code>, <code>posts</code>, and <code>plugins</code>
  #93\n0.5.0.dev16 <a href=\"https://markata.dev/markata/plugins/base_cli/\">base_cli-docs</a></li>\n<li>Remove
  unused function clif that was the original entrypoint #81 0.5.0.dev8</li>\n<li>Allow
  template variables to be used in head config #88 0.5.0.dev12</li>\n<li>Expose <code>markata.__version__</code>
  to templates as <code>__version__</code> #89 0.5.0.dev13</li>\n<li>Fix, ignore post_template
  save on posts without an html attribute #92 0.5.0.dev13</li>\n<li>Fix #33 sluggify
  paths #69 <strong>BREAKING CHANGE</strong> 0.5.0.dev6</li>\n<li>Configurable template
  #70 0.5.dev5, #85 0.5.0.dev11</li>\n<li>Fix #40 Images overlfow outside of body
  #66 0.5.0.dev3</li>\n<li>Created entrypoint hook allowing for users to extend marka
  with jinja\nexensions #60 0.5.0.dev2</li>\n<li>Moved to PEP 517 build #59 0.5.0.dev1</li>\n<li>new
  <code>markata.plugins.redirects</code> will create redirect html files as a backup
  when\nserver-side redirects fail. #76 0.5.0.dev10\n<a href=\"https://markata.dev/markata/plugins/redirects/\">redirects-docs</a></li>\n<li>create
  a slugify migration script #82</li>\n<li>DeepMerge <code>config_overrides</code>
  with config in post render methods #91 0.5.0.dev13</li>\n<li>Create ipython extension
  to automatically load markata #79 0.5.0.dev15</li>\n<li>Fix: images wrapped in a
  link overflow outside the body #96</li>\n<li>new <code>markata.plugins.service_worker</code>
  plugin to create service workers and\nenable offline mode on sites #94 0.5.0.dev15\n<a
  href=\"https://markata.dev/markata/plugins/service-worker/\">service-worker-docs</a></li>\n<li>Fix:
  icons were relatively linked, and were broken for any page other than\nindex, they
  are now absolutely linked to the root of the site. #97 0.5.0.dev16</li>\n<li>Fix:
  auto_descriptions were not rendered on first pass in tui or at all in\nbuild due
  to auto_description running after jinja_md. #100 0.5.0.dev18</li>\n<li>Fix: give
  redirect pages a uniqe description and title #101 0.5.0.dev19</li>\n</ul>\n<h3><code>new</code>
  cli command</h3>\n<p>More information in the <a href=\"https://markata.dev/markata/plugins/base_cli/\">base_cli-docs</a>.</p>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\"># create
  a new blog template</span>\n<span class=\"c1\"># copier requires you to specify
  a directory</span>\nmarkata<span class=\"w\"> </span>new<span class=\"w\"> </span>blog<span
  class=\"w\"> </span><span class=\"o\">[</span>directory<span class=\"o\">]</span>\n\n<span
  class=\"c1\"># create a new blog post</span>\nmarkata<span class=\"w\"> </span>new<span
  class=\"w\"> </span>post\n\n<span class=\"c1\"># create a new plugin</span>\nmarkata<span
  class=\"w\"> </span>new<span class=\"w\"> </span>plugin\n\nmarkata<span class=\"w\">
  </span>new<span class=\"w\"> </span>--help\n\n<span class=\"w\"> </span>Usage:<span
  class=\"w\"> </span>markata<span class=\"w\"> </span>new<span class=\"w\"> </span><span
  class=\"o\">[</span>OPTIONS<span class=\"o\">]</span><span class=\"w\"> </span>COMMAND<span
  class=\"w\"> </span><span class=\"o\">[</span>ARGS<span class=\"o\">]</span>...\n\n<span
  class=\"w\"> </span>create<span class=\"w\"> </span>new<span class=\"w\"> </span>things<span
  class=\"w\"> </span>from<span class=\"w\"> </span>templates\n\n\u256D\u2500<span
  class=\"w\"> </span>Options<span class=\"w\"> </span>\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256E\n\u2502<span
  class=\"w\"> </span>--help<span class=\"w\">          </span>Show<span class=\"w\">
  </span>this<span class=\"w\"> </span>message<span class=\"w\"> </span>and<span class=\"w\">
  </span>exit.<span class=\"w\">                                                                                                                                       </span>\u2502\n\u2570\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256F\n\u256D\u2500<span
  class=\"w\"> </span>Commands<span class=\"w\"> </span>\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256E\n\u2502<span
  class=\"w\"> </span>blog<span class=\"w\">       </span>Create<span class=\"w\">
  </span>a<span class=\"w\"> </span>new<span class=\"w\"> </span>blog<span class=\"w\">
  </span>from<span class=\"w\"> </span>using<span class=\"w\"> </span>the<span class=\"w\">
  </span>template<span class=\"w\"> </span>from<span class=\"w\"> </span>https://github.com/WaylonWalker/markata-blog-starter.<span
  class=\"w\">                                                                   </span>\u2502\n\u2502<span
  class=\"w\"> </span>plugin<span class=\"w\">     </span>Create<span class=\"w\">
  </span>a<span class=\"w\"> </span>new<span class=\"w\"> </span>plugin<span class=\"w\">
  </span>using<span class=\"w\"> </span>the<span class=\"w\"> </span>template<span
  class=\"w\"> </span>at<span class=\"w\"> </span>https://github.com/WaylonWalker/markata-plugin-template.<span
  class=\"w\">                                                                     </span>\u2502\n\u2502<span
  class=\"w\"> </span>post<span class=\"w\">       </span>Create<span class=\"w\">
  </span>new<span class=\"w\"> </span>blog<span class=\"w\"> </span>post<span class=\"w\">
  </span><span class=\"k\">in</span><span class=\"w\"> </span>the<span class=\"w\">
  </span>pages<span class=\"w\"> </span>directory<span class=\"w\"> </span>from<span
  class=\"w\"> </span>the<span class=\"w\"> </span>template<span class=\"w\"> </span>at<span
  class=\"w\">  </span>https://github.com/WaylonWalker/markata-post-template.<span
  class=\"w\">                                               </span>\u2502\n\u2570\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256F\n</pre></div>\n\n</pre>\n<h3>sluggify
  paths</h3>\n<p><code>python-sluggify</code> was implemented to ensure good urls
  are in place despite the\nname of the original file.</p>\n<p>For examples of how
  <code>python-slugify</code> will change your url's see the\n<a href=\"https://pypi.org/project/python-slugify/\">project's
  home page</a>. One\ndifference is that <code>markata</code> will leave <code>/</code>'s
  for routing in the slugs.</p>\n<h4>OPTING OUT</h4>\n<p>If you have an existing site
  and do not want to implement redirects or if you\ndo not want to use slugify, you
  can opt out by setting <code>slugify=False</code> in your\n<code>markata.toml</code>.</p>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
  class=\"n\">slugify</span><span class=\"o\">=</span><span class=\"kc\">false</span>\n</pre></div>\n\n</pre>\n<h4>Migrating
  to slugify</h4>\n<p>From the command line with <code>markata&gt;=0.5.0</code> installed
  run the\nmigration script from the command line to create a redirects file in the\ndefault
  location. This should avoid all 404's as it will create a redirects\nfile that many
  static hosting providers will issue a server-side 301 for, and\nfor those that don't,
  markata.plugins.redirects creates a redirect html page,\nthat will kick in as a
  backup.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>python<span class=\"w\"> </span>-m<span
  class=\"w\"> </span>markata.scripts.migrate_to_slugify\n</pre></div>\n\n</pre>\n<h3>configurable
  page template</h3>\n<p>Now injects seo into the default template through configuration.
  Here is an\nexample, by adding this to your <code>markata.toml</code> configuration.</p>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.head.meta]]</span>\n<span
  class=\"n\">name</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"og:type\"</span>\n<span class=\"n\">content</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">\"article\"</span>\n\n<span
  class=\"k\">[[markata.head.meta]]</span>\n<span class=\"n\">name</span><span class=\"w\">
  </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">\"og:author\"</span>\n<span
  class=\"n\">content</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"Waylon Walker\"</span>\n\n<span class=\"k\">[[markata.head.meta]]</span>\n<span
  class=\"n\">name</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"og:site_name\"</span>\n<span class=\"n\">content</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">\"Waylon
  Walker\"</span>\n\n<span class=\"k\">[[markata.head.meta]]</span>\n<span class=\"n\">name</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">\"theme-color\"</span>\n<span
  class=\"n\">content</span><span class=\"o\">=</span><span class=\"s2\">\"#322D39\"</span>\n\n<span
  class=\"k\">[[markata.head.meta]]</span>\n<span class=\"n\">name</span><span class=\"w\">
  </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s2\">\"twitter:creator\"</span>\n<span
  class=\"n\">content</span><span class=\"o\">=</span><span class=\"s2\">\"@_waylonwalker\"</span>\n</pre></div>\n\n</pre>\n<p>You
  will end up with these meta tags in your html.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"p\">&lt;</span><span
  class=\"nt\">meta</span> <span class=\"na\">name</span><span class=\"o\">=</span><span
  class=\"s\">\"og:type\"</span> <span class=\"na\">content</span><span class=\"o\">=</span><span
  class=\"s\">\"article\"</span> <span class=\"p\">/&gt;</span>\n<span class=\"p\">&lt;</span><span
  class=\"nt\">meta</span> <span class=\"na\">name</span><span class=\"o\">=</span><span
  class=\"s\">\"og:author\"</span> <span class=\"na\">content</span><span class=\"o\">=</span><span
  class=\"s\">\"Waylon Walker\"</span> <span class=\"p\">/&gt;</span>\n<span class=\"p\">&lt;</span><span
  class=\"nt\">meta</span> <span class=\"na\">name</span><span class=\"o\">=</span><span
  class=\"s\">\"og:site_name\"</span> <span class=\"na\">content</span><span class=\"o\">=</span><span
  class=\"s\">\"Waylon Walker\"</span> <span class=\"p\">/&gt;</span>\n<span class=\"p\">&lt;</span><span
  class=\"nt\">meta</span> <span class=\"na\">name</span><span class=\"o\">=</span><span
  class=\"s\">\"theme-color\"</span> <span class=\"na\">content</span><span class=\"o\">=</span><span
  class=\"s\">\"#322D39\"</span> <span class=\"p\">/&gt;</span>\n<span class=\"p\">&lt;</span><span
  class=\"nt\">meta</span> <span class=\"na\">name</span><span class=\"o\">=</span><span
  class=\"s\">\"twitter:creator\"</span> <span class=\"na\">content</span><span class=\"o\">=</span><span
  class=\"s\">\"@_waylonwalker\"</span> <span class=\"p\">/&gt;</span>\n</pre></div>\n\n</pre>\n<p>You
  can have an array of toml tables with a key of text. The text will be\nadded as
  plain text to the end of the head of each page.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[[markata.head]]</span>\n<span
  class=\"n\">text</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s2\">\"\"\"</span>\n<span class=\"s2\">&lt;style&gt;</span>\n<span
  class=\"s2\">img {</span>\n<span class=\"s2\">width: 100%;</span>\n<span class=\"s2\">height:
  auto;</span>\n<span class=\"s2\">}</span>\n<span class=\"s2\">ul {</span>\n<span
  class=\"s2\">  display: flex;</span>\n<span class=\"s2\">  flex-wrap: wrap;</span>\n<span
  class=\"s2\">}</span>\n\n<span class=\"s2\">li {</span>\n<span class=\"s2\">  flex:
  1 2 400px;</span>\n<span class=\"s2\">}</span>\n<span class=\"s2\">&lt;/style&gt;</span>\n\n<span
  class=\"s2\">\"\"\"</span>\n</pre></div>\n\n</pre>\n<p>Descriptions will now properly
  end up in each page.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"p\">&lt;</span><span
  class=\"nt\">meta</span> <span class=\"na\">name</span><span class=\"o\">=</span><span
  class=\"s\">\"description\"</span> <span class=\"na\">content</span><span class=\"o\">=</span><span
  class=\"s\">\"{{ description }}\"</span> <span class=\"p\">/&gt;</span>\n</pre></div>\n\n</pre>\n<h2
  id=\"config-overrides\">Config Overrides <a class=\"header-anchor\" href=\"#config-overrides\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Each post can override
  config settings such as <code>head</code>. New meta tags can be\nadded to a single
  post, or anything that your template might reference from\nconfig.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"nt\">config_overrides</span><span
  class=\"p\">:</span>\n<span class=\"w\">  </span><span class=\"nt\">head</span><span
  class=\"p\">:</span>\n<span class=\"w\">    </span><span class=\"nt\">meta</span><span
  class=\"p\">:</span>\n<span class=\"w\">      </span><span class=\"p p-Indicator\">-</span><span
  class=\"w\"> </span><span class=\"nt\">content</span><span class=\"p\">:</span><span
  class=\"w\"> </span><span class=\"l l-Scalar l-Scalar-Plain\">waylonwalker</span>\n<span
  class=\"w\">        </span><span class=\"nt\">name</span><span class=\"p\">:</span><span
  class=\"w\"> </span><span class=\"l l-Scalar l-Scalar-Plain\">author</span>\n<span
  class=\"w\">      </span><span class=\"p p-Indicator\">-</span><span class=\"w\">
  </span><span class=\"nt\">content</span><span class=\"p\">:</span><span class=\"w\">
  </span><span class=\"s\">\"@_waylonwalker\"</span>\n<span class=\"w\">        </span><span
  class=\"nt\">name</span><span class=\"p\">:</span><span class=\"w\"> </span><span
  class=\"s\">\"twitter:creator\"</span>\n<span class=\"w\">    </span><span class=\"nt\">link</span><span
  class=\"p\">:</span>\n<span class=\"w\">      </span><span class=\"p p-Indicator\">-</span><span
  class=\"w\"> </span><span class=\"nt\">href</span><span class=\"p\">:</span><span
  class=\"w\"> </span><span class=\"l l-Scalar l-Scalar-Plain\">https://waylonwalker.com/that-special-post/</span>\n<span
  class=\"w\">        </span><span class=\"nt\">rel</span><span class=\"p\">:</span><span
  class=\"w\"> </span><span class=\"l l-Scalar l-Scalar-Plain\">canonical</span>\n<span
  class=\"w\">    </span><span class=\"nt\">text</span><span class=\"p\">:</span>\n<span
  class=\"w\">      </span><span class=\"p p-Indicator\">-</span><span class=\"w\">
  </span><span class=\"nt\">value</span><span class=\"p\">:</span><span class=\"w\">
  </span><span class=\"l l-Scalar l-Scalar-Plain\">&lt;link rel='stylesheet' href='/my-extra-styles.css'
  /&gt;</span>\n</pre></div>\n\n</pre>\n<h2 id=\"ipython-extension\">ipython extension
  <a class=\"header-anchor\" href=\"#ipython-extension\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Markata has an ipython
  extension if you want ipython to automatically load with\nan instance of <code>Markata</code>
  mapped to <code>m</code> and <code>markata</code> you can add the following\nto
  your <code>~/.ipython/profile_default/ipython_config.py</code></p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"n\">c</span><span
  class=\"o\">.</span><span class=\"n\">InteractiveShellApp</span><span class=\"o\">.</span><span
  class=\"n\">extensions</span><span class=\"o\">.</span><span class=\"n\">append</span><span
  class=\"p\">(</span><span class=\"s1\">'markata'</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<h2
  id=\"041\">0.4.1 <a class=\"header-anchor\" href=\"#041\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>Issue FutureWarning
  for upcoming change to slugify change that will change urls</li>\n</ul>\n<p>To keep
  existing behavior add this to your <code>markata.toml</code>.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
  class=\"n\">slugify</span><span class=\"o\">=</span><span class=\"kc\">false</span>\n</pre></div>\n\n</pre>\n<h2
  id=\"040\">0.4.0 <a class=\"header-anchor\" href=\"#040\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>feat: add html
  logging with <a href=\"/markata/plugins/setup_logging/\">setup_logging</a>\nplugin
  is all new closes #37</li>\n<li>fix: remove HTML tidy as the site generator tag</li>\n<li>feat:
  create configurable <a href=\"https://markata.dev/nav\">navbar</a></li>\n<li>perf:
  prevent double runs on pre-render and post-render #39</li>\n<li>perf: prevent duplicate
  ruun from to_dict calling pre-render #53\n<ul>\n<li><code>to_dict</code> only runs
  up to <code>render</code> phase if necessary as directed by <code>register_attr</code></li>\n</ul>\n</li>\n<li>perf:
  only prettify if configured #54</li>\n<li>fix: pyinstrument will not create a second
  profiler causing it to end in\nerrors #50</li>\n<li>fix: sites without feeds config
  do not create an index #55</li>\n</ul>\n<h3>Double Runs</h3>\n<p>Previously markata
  would catch AttributeError and run the previous step any\ntime you ran a step too
  early. The way this was implemented caused some steps\nsuch as pre-render and post-render
  to run twice with every single run.</p>\n<p>This change will no longer catch attribute
  errors. If you run into any issues\nwith your plugins not running before asking
  for attributes created by your\nplugin make sure that you implement the\n<a href=\"https://markata.dev/markata/hookspec/#register_attr-function\">@register_attr</a>\ndecorator.</p>\n<h3>Prettify</h3>\n<p>prettify
  html has been turned off by default as beautifulsoup4 prettify was\ntaking a significant
  time, and was often popping up as the slowest parts in my\npersonal <code>_profile</code>.
  If you want to continue running prettify throughout the\nbuild you can set a flag
  in your config to continue running prettify.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
  class=\"n\">prettify_html</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"kc\">true</span>\n</pre></div>\n\n</pre>\n<h2
  id=\"030\">0.3.0 <a class=\"header-anchor\" href=\"#030\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Skipped from a bump2version
  misconfiguration.</p>\n<h2 id=\"020\">0.2.0 <a class=\"header-anchor\" href=\"#020\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>feat: <a href=\"/markata/plugins/auto_description/\">auto_description</a>
  plugin is all\nnew closes #13</li>\n<li>deprecated: long_description has been deprecated
  by auto_description</li>\n<li>fix: <a href=\"/markata/plugins/covers/\">covers</a>
  plugin would previously skip\nevery time.</li>\n<li>feat: <a href=\"/markata/plugins/base_cli/#clean-function\"><code>markata
  clean</code></a> cleans up\nyour cache and output from the command line</li>\n<li>fix:
  <a href=\"/markata/plugins/publish_source/\"><code>publish_source</code></a> plugin
  will now\nignore any non yaml serializable values</li>\n<li>feat: Default template
  colors are now customizable</li>\n<li>feat: Default template now has light and dark
  theme</li>\n<li>feat: map now can map entire posts</li>\n<li>feat: <a href=\"/markata/plugins/prevnext/\">prevnext</a>
  plugin was added to link between\nposts closes #20</li>\n<li>feat: <a href=\"/markata/plugins/jinja_md/\">jinja_md</a>
  plugins was added to incorporate\njinja into all the markdown</li>\n<li>breaking:
  <a href=\"/markata/plugins/feeds\">feeds</a> config now has feeds and\nfeeds_config</li>\n<li>feat:
  <code>output_html</code> can now be specified in the frontmatter\n<a href=\"/markata/plugins/publish_html/#explicityly-set-the-output\">see
  example</a></li>\n<li>feat: edit link is now included in the default page template
  closes #21</li>\n</ul>\n<h3>breaking change to feeds config</h3>\n<p>If you are
  using the 0.1.0 version of feeds, and have configured custom\ntemplates in <code>markata.feeds.template</code>
  and <code>markata.feeds.card_template</code> they\nwill need to be moved to <code>markata.feeds_config</code>.</p>\n<p>Here
  is what you need to do to update your feeds_config.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"gi\">+ [markata.feeds_config]</span>\n<span
  class=\"gi\">+ template=\"pages/templates/archive_template.html\"</span>\n<span
  class=\"gi\">+ card_template=\"plugins/feed_card_template.html\"</span>\n<span class=\"gd\">-
  [markata.feeds]</span>\n<span class=\"gd\">- template=\"pages/templates/archive_template.html\"</span>\n<span
  class=\"gd\">- card_template=\"plugins/feed_card_template.html\"</span>\n</pre></div>\n\n</pre>\n<p><code>markata.feeds</code>
  will only be used to configure feeds pages.</p>\n<h3>map entire posts</h3>\n<p><code>post</code>
  is now exposed to the <code>markata.map</code> object, allowing you to return a\nlist
  of posts.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"n\">m</span>
  <span class=\"o\">=</span> <span class=\"n\">Markata</span><span class=\"p\">()</span>\n<span
  class=\"c1\"># 'post' will return the entire post</span>\n<span class=\"n\">m</span><span
  class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
  class=\"s1\">'post'</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
  class=\"o\">=</span><span class=\"s1\">'\"git\" in tags'</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<h3>Customizable
  colors</h3>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
  class=\"c1\"># default dark theme</span>\n<span class=\"n\">color_bg</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">'#1f2022'</span>\n<span
  class=\"n\">color_bg_code</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s1\">'#1f2022'</span>\n<span class=\"n\">color_text</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">'#eefbfe'</span>\n<span
  class=\"n\">color_link</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s1\">'#47cbff'</span>\n<span class=\"n\">color_accent</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">'#e1bd00c9'</span>\n<span
  class=\"n\">overlay_brightness</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s1\">'.85'</span>\n\n<span class=\"c1\"># pink
  and purple</span>\n<span class=\"n\">color_bg</span><span class=\"w\"> </span><span
  class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">'deeppink'</span>\n<span
  class=\"n\">color_bg_code</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s1\">'rebeccapurple'</span>\n<span class=\"n\">color_text</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">'white'</span>\n<span
  class=\"n\">color_link</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s1\">'aqua'</span>\n<span class=\"n\">color_accent</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">'peachpuff'</span>\n<span
  class=\"n\">overlay_brightness</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s1\">'1.2'</span>\n\n<span class=\"c1\"># default
  light theme</span>\n<span class=\"n\">color_bg_light</span><span class=\"w\"> </span><span
  class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">'#eefbfe'</span>\n<span
  class=\"n\">color_bg_code_light</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s1\">'#eefbfe'</span>\n<span class=\"n\">color_text_light</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">'#1f2022'</span>\n<span
  class=\"n\">color_link_light</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s1\">'#47cbff'</span>\n<span class=\"n\">color_accent_light</span><span
  class=\"w\"> </span><span class=\"o\">=</span><span class=\"w\"> </span><span class=\"s1\">'#ffeb00'</span>\n<span
  class=\"n\">overlay_brightness_light</span><span class=\"w\"> </span><span class=\"o\">=</span><span
  class=\"w\"> </span><span class=\"s1\">'.95'</span>\n</pre></div>\n\n</pre>\n<h3>All
  New auto_description plugin</h3>\n<ul>\n<li>Cache is busted on plugin change</li>\n<li>plugin
  is configurable</li>\n<li>plugin now has docs</li>\n</ul>\n<h3>auto_description
  Configuration</h3>\n<p>Open up your <code>markata.toml</code> file and add new entries
  for your\nauto_descriptions. You can have multiple desriptions, each one will be
  named\nafter the key you give it in your config.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n<span
  class=\"n\">hooks</span><span class=\"o\">=</span><span class=\"p\">[</span>\n<span
  class=\"w\">   </span><span class=\"s2\">\"markata.plugins.auto_description\"</span><span
  class=\"p\">,</span>\n<span class=\"w\">   </span><span class=\"p\">]</span>\n\n<span
  class=\"k\">[markata.auto_description.description]</span>\n<span class=\"n\">len</span><span
  class=\"o\">=</span><span class=\"mi\">160</span>\n<span class=\"k\">[markata.auto_description.long_description]</span>\n<span
  class=\"n\">len</span><span class=\"o\">=</span><span class=\"mi\">250</span>\n<span
  class=\"k\">[markata.auto_description.super_description]</span>\n<span class=\"n\">len</span><span
  class=\"o\">=</span><span class=\"mi\">500</span>\n</pre></div>\n\n</pre>\n<p>In
  the above we will end up with three different descritpions,\n(<code>description</code>,
  <code>long_description</code>, and <code>super_description</code>) each will be
  the\nfirst number of characters from the document as specified in the config.</p>\n<h3>auto_descriptions
  are no longer duplicated</h3>\n<p>4e299d6 fixes the dedupe issue that was in develop
  for a while and closes #24</p>\n<h2 id=\"010\">0.1.0 <a class=\"header-anchor\"
  href=\"#010\"><svg aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\"
  focusable=\"false\" height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<ul>\n<li>fix: pyinstument
  plugin no longer overrides the main cli callback</li>\n<li>feat: default is to run
  the profiler if pyinstrument is installed</li>\n<li>fix: --profile is now under
  the build command</li>\n<li>feat: --pretty/--no-pretty will make pretty tracebacks,
  and skip over\nframework code (closes #4 )</li>\n<li>fix: links are now absolute,
  so they work from GitHub, thanks\n<a href=\"https://github.com/Mr-Destructive\">MR
  Destructive</a></li>\n<li>deprecate: <code>article['content_hash']</code> has been
  removed with preference for\nhashing on <code>article['content']</code></li>\n<li>create
  <code>heading_link</code> plugin</li>\n<li>doc: How to create your home page. <a
  href=\"/home-page/\">docs</a></li>\n</ul>\n<h3>New cli help</h3>\n<p>After the pyinstrument
  plugin was fixed --version and --to-json are back, and\n--profile is now under the
  build command.</p>\n<p>!<a href=\"https://user-images.githubusercontent.com/22648375/150662983-547aebbd-c18c-4c17-8985-a6dc01cd29c7.png\">image</a></p>\n<h3>New
  Heading Link Plugin</h3>\n<p>The new heading link plugin makes it easier to share
  the exact part of an\narticle you want with someone, by giving clickable links to
  the id of the\nheading it's next to.</p>\n<p>!<a href=\"https://user-images.githubusercontent.com/22648375/151718782-08a7cd26-41c1-4f00-a12c-0a208c593e9c.png\">image</a></p>\n<h2
  id=\"001\">0.0.1 <a class=\"header-anchor\" href=\"#001\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>Initial Release \U0001F389</p>\n</main>\n<footer>\xA9
  2024</footer>\n</body></html>"
published: true
slug: changelog
title: Changelog


---

# Markata Changelog

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

```toml
[markata.feeds.published]
template="pages/templates/archive_template.html"
card_template = "pages/templates/feed_card.html"
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

Initial Release 🎉
