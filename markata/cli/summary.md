---
content: "## Run it\n\n``` bash\npython -m markata.cli.summary\n```\n\n## Configuration\n\nThere
  are two main things currently supported by summary, It can count the\nnumber of
  posts based on a filter (`filter_count'), and it can automatically\nlist all the
  values of an attribute and the number of posts that have that\nattribute (`grid_attr`).\n\n###
  grid_attr\n\n`grid_attr` will map over all posts, find all values for each attribute\nconfigured,
  then report the number of posts for each value.\n\n``` toml\n[markata.summary]\ngrid_attr
  = ['tags', 'series']\n```\n\nExample output that will be shown in the summary, counting
  all posts with each\ntag value.\n\n```\nTAGS\n247  python\n90   linux\n68   cli\n49
  \  kedro\n46   bash\n```\n\n### filter_counts\n\n`filter_count` will pass a given
  filter into `markata.map` and return the number\nof posts.\n\n```\n[[markata.summary.filter_count]]\nname='drafts'\nfilter=\"published
  == 'False'\"\ncolor='red'\n\n[[markata.summary.filter_count]]\nname='articles'\ncolor='dark_orange'\n\n[[markata.summary.filter_count]]\nname='py_modules'\nfilter='\"plugin\"
  not in slug and \"docs\" not in str(path)'\ncolor=\"yellow1\"\n\n[markata.summary.filter_count.published]\nfilter=\"published
  == 'True'\"\ncolor='green1'\n\n[markata.summary.filter_count.plugins]\nfilter='\"plugin\"
  in slug and \"docs\" not in str(path)'\ncolor=\"blue\"\n\n[markata.summary.filter_count.docs]\nfilter=\"'docs'
  in str(path)\"\ncolor='purple'\n```\n\nExample output might look like this, showing
  the number of posts that contained\nwithin each filter specified.\n\n```\n8 drafts\n66
  articles\n20 py_modules\n58 published\n38 plugins\n8 docs\n```\n\n\n!! class <h2
  id='FilterCount' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>FilterCount
  <em class='small'>class</em></h2>\n\n???+ source \"FilterCount <em class='small'>source</em>\"\n\n```python\n\n
  \       class FilterCount(pydantic.BaseModel):\n            name: str\n            filter:
  str = \"True\"\n            color: str = \"white\"\n```\n\n\n!! class <h2 id='SummaryConfig'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>SummaryConfig <em
  class='small'>class</em></h2>\n\n???+ source \"SummaryConfig <em class='small'>source</em>\"\n\n```python\n\n
  \       class SummaryConfig(pydantic.BaseModel):\n            grid_attr: List[str]
  = [\"tags\", \"series\"]\n            filter_count: List[FilterCount] = FilterCount(\n
  \               name=\"drafts\", filter=\"published == 'False'\", color=\"red\"\n
  \           )\n```\n\n\n!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>Config <em class='small'>class</em></h2>\n\n???+ source \"Config <em class='small'>source</em>\"\n\n```python\n\n
  \       class Config(pydantic.BaseModel):\n            summary: SummaryConfig =
  SummaryConfig()\n```\n\n\n!! function <h2 id='config_model' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>\n\n???+
  source \"config_model <em class='small'>source</em>\"\n\n```python\n\n        def
  config_model(markata: \"Markata\") -> None:\n            markata.config_models.append(Config)\n```\n\n\n!!
  class <h2 id='Summary' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Summary
  <em class='small'>class</em></h2>\n\n???+ source \"Summary <em class='small'>source</em>\"\n\n```python\n\n
  \       class Summary:\n            def __init__(self, m: \"Markata\", simple: bool
  = False) -> None:\n                self.m = m\n                self.simple = simple\n\n
  \           def get_grid(self) -> None:\n                \"create a rich grid to
  display the summary\"\n                self.grid = Table.grid(expand=True)\n\n                for
  filter_count in self.m.config.summary.filter_count:\n                    self.filter_count(filter_count)\n\n
  \               for attr in self.m.config.summary.grid_attr:\n                    self.grid_attr(attr)\n\n
  \               return self.grid\n\n            def filter_count(\n                self,\n
  \               fc: FilterCount,\n            ) -> None:\n                \"add
  a row in the grid for the number of items in a filter config\"\n                self.grid.add_row(\n
  \                   f\"[{fc.color}]{len(self.m.map(filter=fc.filter))}[/] {fc.name}\"\n
  \               )\n\n            def grid_attr(self, attr: str) -> None:\n                \"add
  attribute the the object grid\"\n                posts = list(\n                    flatten(\n
  \                       [\n                            tags if isinstance(tags,
  list) else [tags]\n                            for a in self.m.posts\n                            if
  (tags := a.get(attr, None)) is not None\n                        ],\n                    ),\n
  \               )\n                if len(posts) > 0:\n                    self.grid.add_row()\n
  \                   self.grid.add_row(f\"[bold gold1]{attr.upper()}[/]\")\n                    for
  post, count in Counter(posts).most_common():\n                        self.grid.add_row(f'{count}
  {\" \"*(3-len(str(count)))} {post}')\n\n            def __rich__(self) -> Union[Panel,
  Table]:\n                grid = self.get_grid()\n\n                if self.simple:\n
  \                   return grid\n                else:\n                    return
  Panel(\n                        grid,\n                        title=\"[gold1]summary[/]\",\n
  \                       border_style=\"magenta\",\n                        expand=False,\n
  \                   )\n```\n\n\n!! function <h2 id='configure' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>configure <em class='small'>function</em></h2>\n\n???+
  source \"configure <em class='small'>source</em>\"\n\n```python\n\n        def configure(markata:
  \"Markata\") -> None:\n            def get_summary(self):\n                try:\n
  \                   return self._summary\n                except AttributeError:\n
  \                   self._summary: Summary = Summary(self)\n                    return
  self._summary\n\n            from markata import Markata\n\n            Markata.summary
  = property(get_summary)\n```\n\n\n!! function <h2 id='cli' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2>\n
  \   Markata hook to implement base cli commands.\n???+ source \"cli <em class='small'>source</em>\"\n\n```python\n\n
  \       def cli(app: typer.Typer, markata: \"Markata\") -> None:\n            \"\"\"\n
  \           Markata hook to implement base cli commands.\n            \"\"\"\n            summary_app
  = typer.Typer()\n            app.add_typer(summary_app, name=\"summary\")\n\n            @summary_app.callback(invoke_without_command=True)\n
  \           def summary():\n                \"show the application summary\"\n                from
  rich import print\n\n                markata.console.quiet = True\n\n                print(Summary(markata,
  simple=True))\n```\n\n\n!! method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>__init__ <em class='small'>method</em></h2>\n\n???+ source \"__init__ <em
  class='small'>source</em>\"\n\n```python\n\n        def __init__(self, m: \"Markata\",
  simple: bool = False) -> None:\n                self.m = m\n                self.simple
  = simple\n```\n\n\n!! method <h2 id='get_grid' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>get_grid <em class='small'>method</em></h2>\n    create a rich grid to display
  the summary\n???+ source \"get_grid <em class='small'>source</em>\"\n\n```python\n\n
  \       def get_grid(self) -> None:\n                \"create a rich grid to display
  the summary\"\n                self.grid = Table.grid(expand=True)\n\n                for
  filter_count in self.m.config.summary.filter_count:\n                    self.filter_count(filter_count)\n\n
  \               for attr in self.m.config.summary.grid_attr:\n                    self.grid_attr(attr)\n\n
  \               return self.grid\n```\n\n\n!! method <h2 id='filter_count' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>filter_count <em class='small'>method</em></h2>\n
  \   add a row in the grid for the number of items in a filter config\n???+ source
  \"filter_count <em class='small'>source</em>\"\n\n```python\n\n        def filter_count(\n
  \               self,\n                fc: FilterCount,\n            ) -> None:\n
  \               \"add a row in the grid for the number of items in a filter config\"\n
  \               self.grid.add_row(\n                    f\"[{fc.color}]{len(self.m.map(filter=fc.filter))}[/]
  {fc.name}\"\n                )\n```\n\n\n!! method <h2 id='grid_attr' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>grid_attr <em class='small'>method</em></h2>\n
  \   add attribute the the object grid\n???+ source \"grid_attr <em class='small'>source</em>\"\n\n```python\n\n
  \       def grid_attr(self, attr: str) -> None:\n                \"add attribute
  the the object grid\"\n                posts = list(\n                    flatten(\n
  \                       [\n                            tags if isinstance(tags,
  list) else [tags]\n                            for a in self.m.posts\n                            if
  (tags := a.get(attr, None)) is not None\n                        ],\n                    ),\n
  \               )\n                if len(posts) > 0:\n                    self.grid.add_row()\n
  \                   self.grid.add_row(f\"[bold gold1]{attr.upper()}[/]\")\n                    for
  post, count in Counter(posts).most_common():\n                        self.grid.add_row(f'{count}
  {\" \"*(3-len(str(count)))} {post}')\n```\n\n\n!! method <h2 id='__rich__' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>__rich__ <em class='small'>method</em></h2>\n\n???+
  source \"__rich__ <em class='small'>source</em>\"\n\n```python\n\n        def __rich__(self)
  -> Union[Panel, Table]:\n                grid = self.get_grid()\n\n                if
  self.simple:\n                    return grid\n                else:\n                    return
  Panel(\n                        grid,\n                        title=\"[gold1]summary[/]\",\n
  \                       border_style=\"magenta\",\n                        expand=False,\n
  \                   )\n```\n\n\n!! function <h2 id='get_summary' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>get_summary <em class='small'>function</em></h2>\n\n???+
  source \"get_summary <em class='small'>source</em>\"\n\n```python\n\n        def
  get_summary(self):\n                try:\n                    return self._summary\n
  \               except AttributeError:\n                    self._summary: Summary
  = Summary(self)\n                    return self._summary\n```\n\n\n!! function
  <h2 id='summary' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>summary
  <em class='small'>function</em></h2>\n    show the application summary\n???+ source
  \"summary <em class='small'>source</em>\"\n\n```python\n\n        def summary():\n
  \               \"show the application summary\"\n                from rich import
  print\n\n                markata.console.quiet = True\n\n                print(Summary(markata,
  simple=True))\n```\n\n"
date: 0001-01-01
description: There are two main things currently supported by summary, It can count
  the grid_attr Example output that will be shown in the summary, counting all posts
  with e
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Summary.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"There are two main things currently supported by
    summary, It can count the grid_attr Example output that will be shown in the summary,
    counting all posts with e\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Summary.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"There are two main things currently supported
    by summary, It can count the grid_attr Example output that will be shown in the
    summary, counting all posts with e\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<div class='container flex flex-row min-h-screen'>\n    <div>\n    </div>\n
    \   <div class='flex-grow px-8 mx-auto min-h-screen'>\n<header class='flex justify-center
    items-center p-8'>\n\n    <nav class='flex justify-center items-center my-8'>\n
    \       <a\n            href='/'>markata</a>\n        <a\n            href='https://github.com/WaylonWalker/markata'>GitHub</a>\n
    \       <a\n            href='https://markata.dev/docs/'>docs</a>\n        <a\n
    \           href='https://markata.dev/plugins/'>plugins</a>\n    </nav>\n\n    <div>\n
    \       <label id=\"theme-switch\" class=\"theme-switch\" for=\"checkbox-theme\"
    title=\"light/dark mode toggle\">\n            <input type=\"checkbox\" id=\"checkbox-theme\"
    />\n            <div class=\"slider round\"></div>\n        </label>\n    </div>\n</header><article
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Summary.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <h2 id=\"run-it\">Run
    it <a class=\"header-anchor\" href=\"#run-it\"><svg class=\"heading-permalink\"
    aria-hidden=\"true\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
    viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
    d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre class='wrapper'>\n\n<div
    class='copy-wrapper'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>python<span class=\"w\">
    </span>-m<span class=\"w\"> </span>markata.cli.summary\n</pre></div>\n\n</pre>\n\n<h2
    id=\"configuration\">Configuration <a class=\"header-anchor\" href=\"#configuration\"><svg
    class=\"heading-permalink\" aria-hidden=\"true\" fill=\"currentColor\" focusable=\"false\"
    height=\"1em\" viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
    d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>There are two main things
    currently supported by summary, It can count the\nnumber of posts based on a filter
    (<code>filter_count'), and it can automatically list all the values of an attribute
    and the number of posts that have that attribute (</code>grid_attr`).</p>\n<h3>grid_attr</h3>\n<p><code>grid_attr</code>
    will map over all posts, find all values for each attribute\nconfigured, then
    report the number of posts for each value.</p>\n<pre class='wrapper'>\n\n<div
    class='copy-wrapper'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.summary]</span>\n<span
    class=\"n\">grid_attr</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span><span class=\"s1\">&#39;tags&#39;</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">&#39;series&#39;</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>Example output that will be
    shown in the summary, counting all posts with each\ntag value.</p>\n<pre class='wrapper'>\n\n<div
    class='copy-wrapper'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>TAGS\n247  python\n90   linux\n68
    \  cli\n49   kedro\n46   bash\n</pre></div>\n\n</pre>\n\n<h3>filter_counts</h3>\n<p><code>filter_count</code>
    will pass a given filter into <code>markata.map</code> and return the number\nof
    posts.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy'
    title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>[[markata.summary.filter_count]]\nname=&#39;drafts&#39;\nfilter=&quot;published
    == &#39;False&#39;&quot;\ncolor=&#39;red&#39;\n\n[[markata.summary.filter_count]]\nname=&#39;articles&#39;\ncolor=&#39;dark_orange&#39;\n\n[[markata.summary.filter_count]]\nname=&#39;py_modules&#39;\nfilter=&#39;&quot;plugin&quot;
    not in slug and &quot;docs&quot; not in str(path)&#39;\ncolor=&quot;yellow1&quot;\n\n[markata.summary.filter_count.published]\nfilter=&quot;published
    == &#39;True&#39;&quot;\ncolor=&#39;green1&#39;\n\n[markata.summary.filter_count.plugins]\nfilter=&#39;&quot;plugin&quot;
    in slug and &quot;docs&quot; not in str(path)&#39;\ncolor=&quot;blue&quot;\n\n[markata.summary.filter_count.docs]\nfilter=&quot;&#39;docs&#39;
    in str(path)&quot;\ncolor=&#39;purple&#39;\n</pre></div>\n\n</pre>\n\n<p>Example
    output might look like this, showing the number of posts that contained\nwithin
    each filter specified.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>8 drafts\n66 articles\n20
    py_modules\n58 published\n38 plugins\n8 docs\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='FilterCount' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>FilterCount <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">FilterCount
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">FilterCount</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">name</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n
    \           <span class=\"nb\">filter</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span>\n            <span
    class=\"n\">color</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;white&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='SummaryConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>SummaryConfig <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">SummaryConfig
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">SummaryConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">grid_attr</span><span class=\"p\">:</span> <span
    class=\"n\">List</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
    class=\"s2\">&quot;tags&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;series&quot;</span><span
    class=\"p\">]</span>\n            <span class=\"n\">filter_count</span><span class=\"p\">:</span>
    <span class=\"n\">List</span><span class=\"p\">[</span><span class=\"n\">FilterCount</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">FilterCount</span><span
    class=\"p\">(</span>\n                <span class=\"n\">name</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;drafts&quot;</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;published == &#39;False&#39;&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">color</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;red&quot;</span>\n            <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Config <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">summary</span><span class=\"p\">:</span> <span class=\"n\">SummaryConfig</span>
    <span class=\"o\">=</span> <span class=\"n\">SummaryConfig</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>config_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">config_model
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Summary' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>Summary <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Summary
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Summary</span><span class=\"p\">:</span>\n            <span
    class=\"k\">def</span> <span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">m</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">simple</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
    <span class=\"o\">=</span> <span class=\"kc\">False</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">m</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">simple</span>
    <span class=\"o\">=</span> <span class=\"n\">simple</span>\n\n            <span
    class=\"k\">def</span> <span class=\"nf\">get_grid</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"s2\">&quot;create a rich grid to display the summary&quot;</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid</span>
    <span class=\"o\">=</span> <span class=\"n\">Table</span><span class=\"o\">.</span><span
    class=\"n\">grid</span><span class=\"p\">(</span><span class=\"n\">expand</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">for</span> <span class=\"n\">filter_count</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">summary</span><span class=\"o\">.</span><span
    class=\"n\">filter_count</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">filter_count</span><span
    class=\"p\">(</span><span class=\"n\">filter_count</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">for</span> <span class=\"n\">attr</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">summary</span><span class=\"o\">.</span><span
    class=\"n\">grid_attr</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid_attr</span><span
    class=\"p\">(</span><span class=\"n\">attr</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">grid</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">filter_count</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">,</span>\n                <span class=\"n\">fc</span><span
    class=\"p\">:</span> <span class=\"n\">FilterCount</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"s2\">&quot;add a row in the
    grid for the number of items in a filter config&quot;</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
    class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span>\n
    \                   <span class=\"sa\">f</span><span class=\"s2\">&quot;[</span><span
    class=\"si\">{</span><span class=\"n\">fc</span><span class=\"o\">.</span><span
    class=\"n\">color</span><span class=\"si\">}</span><span class=\"s2\">]</span><span
    class=\"si\">{</span><span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"nb\">filter</span><span class=\"o\">=</span><span class=\"n\">fc</span><span
    class=\"o\">.</span><span class=\"n\">filter</span><span class=\"p\">))</span><span
    class=\"si\">}</span><span class=\"s2\">[/] </span><span class=\"si\">{</span><span
    class=\"n\">fc</span><span class=\"o\">.</span><span class=\"n\">name</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                <span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">grid_attr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">attr</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"s2\">&quot;add attribute the
    the object grid&quot;</span>\n                <span class=\"n\">posts</span> <span
    class=\"o\">=</span> <span class=\"nb\">list</span><span class=\"p\">(</span>\n
    \                   <span class=\"n\">flatten</span><span class=\"p\">(</span>\n
    \                       <span class=\"p\">[</span>\n                            <span
    class=\"n\">tags</span> <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">tags</span><span class=\"p\">,</span> <span
    class=\"nb\">list</span><span class=\"p\">)</span> <span class=\"k\">else</span>
    <span class=\"p\">[</span><span class=\"n\">tags</span><span class=\"p\">]</span>\n
    \                           <span class=\"k\">for</span> <span class=\"n\">a</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">posts</span>\n
    \                           <span class=\"k\">if</span> <span class=\"p\">(</span><span
    class=\"n\">tags</span> <span class=\"o\">:=</span> <span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">attr</span><span class=\"p\">,</span> <span class=\"kc\">None</span><span
    class=\"p\">))</span> <span class=\"ow\">is</span> <span class=\"ow\">not</span>
    <span class=\"kc\">None</span>\n                        <span class=\"p\">],</span>\n
    \                   <span class=\"p\">),</span>\n                <span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">)</span> <span
    class=\"o\">&gt;</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">grid</span><span class=\"o\">.</span><span class=\"n\">add_row</span><span
    class=\"p\">()</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_row</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;[bold gold1]</span><span class=\"si\">{</span><span class=\"n\">attr</span><span
    class=\"o\">.</span><span class=\"n\">upper</span><span class=\"p\">()</span><span
    class=\"si\">}</span><span class=\"s2\">[/]&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">for</span> <span class=\"n\">post</span><span
    class=\"p\">,</span> <span class=\"n\">count</span> <span class=\"ow\">in</span>
    <span class=\"n\">Counter</span><span class=\"p\">(</span><span class=\"n\">posts</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">most_common</span><span
    class=\"p\">():</span>\n                        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_row</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s1\">&#39;</span><span class=\"si\">{</span><span class=\"n\">count</span><span
    class=\"si\">}</span><span class=\"s1\"> </span><span class=\"si\">{</span><span
    class=\"s2\">&quot; &quot;</span><span class=\"o\">*</span><span class=\"p\">(</span><span
    class=\"mi\">3</span><span class=\"o\">-</span><span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">count</span><span class=\"p\">)))</span><span class=\"si\">}</span><span
    class=\"s1\"> </span><span class=\"si\">{</span><span class=\"n\">post</span><span
    class=\"si\">}</span><span class=\"s1\">&#39;</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">__rich__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Union</span><span class=\"p\">[</span><span
    class=\"n\">Panel</span><span class=\"p\">,</span> <span class=\"n\">Table</span><span
    class=\"p\">]:</span>\n                <span class=\"n\">grid</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">get_grid</span><span
    class=\"p\">()</span>\n\n                <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">simple</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">grid</span>\n
    \               <span class=\"k\">else</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">Panel</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">grid</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;[gold1]summary[/]&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">border_style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;magenta&quot;</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">expand</span><span class=\"o\">=</span><span class=\"kc\">False</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>configure <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">configure
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">get_summary</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_summary</span>\n                <span class=\"k\">except</span> <span
    class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_summary</span><span
    class=\"p\">:</span> <span class=\"n\">Summary</span> <span class=\"o\">=</span>
    <span class=\"n\">Summary</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_summary</span>\n\n
    \           <span class=\"kn\">from</span> <span class=\"nn\">markata</span> <span
    class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n            <span
    class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">summary</span>
    <span class=\"o\">=</span> <span class=\"nb\">property</span><span class=\"p\">(</span><span
    class=\"n\">get_summary</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli
    <em class='small'>function</em></h2>\nMarkata hook to implement base cli commands.</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">cli
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">cli</span><span class=\"p\">(</span><span class=\"n\">app</span><span
    class=\"p\">:</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Typer</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Markata hook to implement base cli commands.</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n            <span class=\"n\">summary_app</span>
    <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Typer</span><span class=\"p\">()</span>\n            <span class=\"n\">app</span><span
    class=\"o\">.</span><span class=\"n\">add_typer</span><span class=\"p\">(</span><span
    class=\"n\">summary_app</span><span class=\"p\">,</span> <span class=\"n\">name</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;summary&quot;</span><span class=\"p\">)</span>\n\n
    \           <span class=\"nd\">@summary_app</span><span class=\"o\">.</span><span
    class=\"n\">callback</span><span class=\"p\">(</span><span class=\"n\">invoke_without_command</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">summary</span><span
    class=\"p\">():</span>\n                <span class=\"s2\">&quot;show the application
    summary&quot;</span>\n                <span class=\"kn\">from</span> <span class=\"nn\">rich</span>
    <span class=\"kn\">import</span> <span class=\"nb\">print</span>\n\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n\n                <span class=\"nb\">print</span><span
    class=\"p\">(</span><span class=\"n\">Summary</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">simple</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>init</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>init</strong>
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">m</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span class=\"n\">simple</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"kc\">False</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span> <span
    class=\"o\">=</span> <span class=\"n\">m</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">simple</span> <span class=\"o\">=</span>
    <span class=\"n\">simple</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='get_grid'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_grid <em class='small'>method</em></h2>\ncreate
    a rich grid to display the summary</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">get_grid <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">get_grid</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"s2\">&quot;create a rich grid
    to display the summary&quot;</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">grid</span> <span class=\"o\">=</span> <span
    class=\"n\">Table</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
    class=\"p\">(</span><span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n\n                <span class=\"k\">for</span>
    <span class=\"n\">filter_count</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">summary</span><span
    class=\"o\">.</span><span class=\"n\">filter_count</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">filter_count</span><span class=\"p\">(</span><span class=\"n\">filter_count</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">for</span> <span class=\"n\">attr</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">summary</span><span class=\"o\">.</span><span
    class=\"n\">grid_attr</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid_attr</span><span
    class=\"p\">(</span><span class=\"n\">attr</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">grid</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='filter_count' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>filter_count <em class='small'>method</em></h2>\nadd a row in the grid
    for the number of items in a filter config</p>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">filter_count
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">filter_count</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">,</span>\n                <span class=\"n\">fc</span><span
    class=\"p\">:</span> <span class=\"n\">FilterCount</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"s2\">&quot;add a row in the
    grid for the number of items in a filter config&quot;</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
    class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span>\n
    \                   <span class=\"sa\">f</span><span class=\"s2\">&quot;[</span><span
    class=\"si\">{</span><span class=\"n\">fc</span><span class=\"o\">.</span><span
    class=\"n\">color</span><span class=\"si\">}</span><span class=\"s2\">]</span><span
    class=\"si\">{</span><span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"nb\">filter</span><span class=\"o\">=</span><span class=\"n\">fc</span><span
    class=\"o\">.</span><span class=\"n\">filter</span><span class=\"p\">))</span><span
    class=\"si\">}</span><span class=\"s2\">[/] </span><span class=\"si\">{</span><span
    class=\"n\">fc</span><span class=\"o\">.</span><span class=\"n\">name</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='grid_attr' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>grid_attr <em class='small'>method</em></h2>\nadd attribute the the object
    grid</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">grid_attr <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">grid_attr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">attr</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"s2\">&quot;add attribute the the object grid&quot;</span>\n                <span
    class=\"n\">posts</span> <span class=\"o\">=</span> <span class=\"nb\">list</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">flatten</span><span
    class=\"p\">(</span>\n                        <span class=\"p\">[</span>\n                            <span
    class=\"n\">tags</span> <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">tags</span><span class=\"p\">,</span> <span
    class=\"nb\">list</span><span class=\"p\">)</span> <span class=\"k\">else</span>
    <span class=\"p\">[</span><span class=\"n\">tags</span><span class=\"p\">]</span>\n
    \                           <span class=\"k\">for</span> <span class=\"n\">a</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">posts</span>\n
    \                           <span class=\"k\">if</span> <span class=\"p\">(</span><span
    class=\"n\">tags</span> <span class=\"o\">:=</span> <span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">attr</span><span class=\"p\">,</span> <span class=\"kc\">None</span><span
    class=\"p\">))</span> <span class=\"ow\">is</span> <span class=\"ow\">not</span>
    <span class=\"kc\">None</span>\n                        <span class=\"p\">],</span>\n
    \                   <span class=\"p\">),</span>\n                <span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">)</span> <span
    class=\"o\">&gt;</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">grid</span><span class=\"o\">.</span><span class=\"n\">add_row</span><span
    class=\"p\">()</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_row</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;[bold gold1]</span><span class=\"si\">{</span><span class=\"n\">attr</span><span
    class=\"o\">.</span><span class=\"n\">upper</span><span class=\"p\">()</span><span
    class=\"si\">}</span><span class=\"s2\">[/]&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">for</span> <span class=\"n\">post</span><span
    class=\"p\">,</span> <span class=\"n\">count</span> <span class=\"ow\">in</span>
    <span class=\"n\">Counter</span><span class=\"p\">(</span><span class=\"n\">posts</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">most_common</span><span
    class=\"p\">():</span>\n                        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_row</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s1\">&#39;</span><span class=\"si\">{</span><span class=\"n\">count</span><span
    class=\"si\">}</span><span class=\"s1\"> </span><span class=\"si\">{</span><span
    class=\"s2\">&quot; &quot;</span><span class=\"o\">*</span><span class=\"p\">(</span><span
    class=\"mi\">3</span><span class=\"o\">-</span><span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">count</span><span class=\"p\">)))</span><span class=\"si\">}</span><span
    class=\"s1\"> </span><span class=\"si\">{</span><span class=\"n\">post</span><span
    class=\"si\">}</span><span class=\"s1\">&#39;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>rich</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>rich</strong>
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Union</span><span
    class=\"p\">[</span><span class=\"n\">Panel</span><span class=\"p\">,</span> <span
    class=\"n\">Table</span><span class=\"p\">]:</span>\n                <span class=\"n\">grid</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">get_grid</span><span class=\"p\">()</span>\n\n                <span
    class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">simple</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">grid</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"k\">return</span>
    <span class=\"n\">Panel</span><span class=\"p\">(</span>\n                        <span
    class=\"n\">grid</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s2\">&quot;[gold1]summary[/]&quot;</span><span
    class=\"p\">,</span>\n                        <span class=\"n\">border_style</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;magenta&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">False</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='get_summary'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_summary <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">get_summary <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">get_summary</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_summary</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_summary</span><span
    class=\"p\">:</span> <span class=\"n\">Summary</span> <span class=\"o\">=</span>
    <span class=\"n\">Summary</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_summary</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='summary' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>summary <em class='small'>function</em></h2>\nshow the application summary</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">summary
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">summary</span><span class=\"p\">():</span>\n                <span
    class=\"s2\">&quot;show the application summary&quot;</span>\n                <span
    class=\"kn\">from</span> <span class=\"nn\">rich</span> <span class=\"kn\">import</span>
    <span class=\"nb\">print</span>\n\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">quiet</span> <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n
    \               <span class=\"nb\">print</span><span class=\"p\">(</span><span
    class=\"n\">Summary</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">simple</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Summary.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"There are two main things currently supported by
    summary, It can count the grid_attr Example output that will be shown in the summary,
    counting all posts with e\" />\n <link href=\"/favicon.ico\" rel=\"icon\" type=\"image/png\"
    />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link rel=\"stylesheet\"
    href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta name=\"og:author_email\"
    content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Summary.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"There are two main things currently supported
    by summary, It can count the grid_attr Example output that will be shown in the
    summary, counting all posts with e\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Summary.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Summary.Py\n    </h1>\n</section>    <section class=\"body\">\n        <h2
    id=\"run-it\">Run it <a class=\"header-anchor\" href=\"#run-it\"><svg class=\"heading-permalink\"
    aria-hidden=\"true\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
    viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
    d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<pre class='wrapper'>\n\n<div
    class='copy-wrapper'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>python<span class=\"w\">
    </span>-m<span class=\"w\"> </span>markata.cli.summary\n</pre></div>\n\n</pre>\n\n<h2
    id=\"configuration\">Configuration <a class=\"header-anchor\" href=\"#configuration\"><svg
    class=\"heading-permalink\" aria-hidden=\"true\" fill=\"currentColor\" focusable=\"false\"
    height=\"1em\" viewBox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
    d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
    5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
    6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
    3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
    3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
    13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
    2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0
    0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
    19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997
    0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
    1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h2>\n<p>There are two main things
    currently supported by summary, It can count the\nnumber of posts based on a filter
    (<code>filter_count'), and it can automatically list all the values of an attribute
    and the number of posts that have that attribute (</code>grid_attr`).</p>\n<h3>grid_attr</h3>\n<p><code>grid_attr</code>
    will map over all posts, find all values for each attribute\nconfigured, then
    report the number of posts for each value.</p>\n<pre class='wrapper'>\n\n<div
    class='copy-wrapper'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata.summary]</span>\n<span
    class=\"n\">grid_attr</span><span class=\"w\"> </span><span class=\"o\">=</span><span
    class=\"w\"> </span><span class=\"p\">[</span><span class=\"s1\">&#39;tags&#39;</span><span
    class=\"p\">,</span><span class=\"w\"> </span><span class=\"s1\">&#39;series&#39;</span><span
    class=\"p\">]</span>\n</pre></div>\n\n</pre>\n\n<p>Example output that will be
    shown in the summary, counting all posts with each\ntag value.</p>\n<pre class='wrapper'>\n\n<div
    class='copy-wrapper'>\n\n<button class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>TAGS\n247  python\n90   linux\n68
    \  cli\n49   kedro\n46   bash\n</pre></div>\n\n</pre>\n\n<h3>filter_counts</h3>\n<p><code>filter_count</code>
    will pass a given filter into <code>markata.map</code> and return the number\nof
    posts.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy'
    title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>[[markata.summary.filter_count]]\nname=&#39;drafts&#39;\nfilter=&quot;published
    == &#39;False&#39;&quot;\ncolor=&#39;red&#39;\n\n[[markata.summary.filter_count]]\nname=&#39;articles&#39;\ncolor=&#39;dark_orange&#39;\n\n[[markata.summary.filter_count]]\nname=&#39;py_modules&#39;\nfilter=&#39;&quot;plugin&quot;
    not in slug and &quot;docs&quot; not in str(path)&#39;\ncolor=&quot;yellow1&quot;\n\n[markata.summary.filter_count.published]\nfilter=&quot;published
    == &#39;True&#39;&quot;\ncolor=&#39;green1&#39;\n\n[markata.summary.filter_count.plugins]\nfilter=&#39;&quot;plugin&quot;
    in slug and &quot;docs&quot; not in str(path)&#39;\ncolor=&quot;blue&quot;\n\n[markata.summary.filter_count.docs]\nfilter=&quot;&#39;docs&#39;
    in str(path)&quot;\ncolor=&#39;purple&#39;\n</pre></div>\n\n</pre>\n\n<p>Example
    output might look like this, showing the number of posts that contained\nwithin
    each filter specified.</p>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>8 drafts\n66 articles\n20
    py_modules\n58 published\n38 plugins\n8 docs\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='FilterCount' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>FilterCount <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">FilterCount
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">FilterCount</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">name</span><span class=\"p\">:</span> <span class=\"nb\">str</span>\n
    \           <span class=\"nb\">filter</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;True&quot;</span>\n            <span
    class=\"n\">color</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
    <span class=\"o\">=</span> <span class=\"s2\">&quot;white&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='SummaryConfig' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>SummaryConfig <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">SummaryConfig
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">SummaryConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">grid_attr</span><span class=\"p\">:</span> <span
    class=\"n\">List</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"p\">[</span><span
    class=\"s2\">&quot;tags&quot;</span><span class=\"p\">,</span> <span class=\"s2\">&quot;series&quot;</span><span
    class=\"p\">]</span>\n            <span class=\"n\">filter_count</span><span class=\"p\">:</span>
    <span class=\"n\">List</span><span class=\"p\">[</span><span class=\"n\">FilterCount</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">FilterCount</span><span
    class=\"p\">(</span>\n                <span class=\"n\">name</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;drafts&quot;</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;published == &#39;False&#39;&quot;</span><span
    class=\"p\">,</span> <span class=\"n\">color</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;red&quot;</span>\n            <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Config <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
    class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
    \           <span class=\"n\">summary</span><span class=\"p\">:</span> <span class=\"n\">SummaryConfig</span>
    <span class=\"o\">=</span> <span class=\"n\">SummaryConfig</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>config_model <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">config_model
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
    class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
    class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Summary' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>Summary <em class='small'>class</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">Summary
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Summary</span><span class=\"p\">:</span>\n            <span
    class=\"k\">def</span> <span class=\"fm\">__init__</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"n\">m</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">simple</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
    <span class=\"o\">=</span> <span class=\"kc\">False</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">m</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">simple</span>
    <span class=\"o\">=</span> <span class=\"n\">simple</span>\n\n            <span
    class=\"k\">def</span> <span class=\"nf\">get_grid</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"s2\">&quot;create a rich grid to display the summary&quot;</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid</span>
    <span class=\"o\">=</span> <span class=\"n\">Table</span><span class=\"o\">.</span><span
    class=\"n\">grid</span><span class=\"p\">(</span><span class=\"n\">expand</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">for</span> <span class=\"n\">filter_count</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">summary</span><span class=\"o\">.</span><span
    class=\"n\">filter_count</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">filter_count</span><span
    class=\"p\">(</span><span class=\"n\">filter_count</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">for</span> <span class=\"n\">attr</span> <span
    class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">summary</span><span class=\"o\">.</span><span
    class=\"n\">grid_attr</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid_attr</span><span
    class=\"p\">(</span><span class=\"n\">attr</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">grid</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">filter_count</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">,</span>\n                <span class=\"n\">fc</span><span
    class=\"p\">:</span> <span class=\"n\">FilterCount</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"s2\">&quot;add a row in the
    grid for the number of items in a filter config&quot;</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
    class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span>\n
    \                   <span class=\"sa\">f</span><span class=\"s2\">&quot;[</span><span
    class=\"si\">{</span><span class=\"n\">fc</span><span class=\"o\">.</span><span
    class=\"n\">color</span><span class=\"si\">}</span><span class=\"s2\">]</span><span
    class=\"si\">{</span><span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"nb\">filter</span><span class=\"o\">=</span><span class=\"n\">fc</span><span
    class=\"o\">.</span><span class=\"n\">filter</span><span class=\"p\">))</span><span
    class=\"si\">}</span><span class=\"s2\">[/] </span><span class=\"si\">{</span><span
    class=\"n\">fc</span><span class=\"o\">.</span><span class=\"n\">name</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                <span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">grid_attr</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
    class=\"n\">attr</span><span class=\"p\">:</span> <span class=\"nb\">str</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"s2\">&quot;add attribute the
    the object grid&quot;</span>\n                <span class=\"n\">posts</span> <span
    class=\"o\">=</span> <span class=\"nb\">list</span><span class=\"p\">(</span>\n
    \                   <span class=\"n\">flatten</span><span class=\"p\">(</span>\n
    \                       <span class=\"p\">[</span>\n                            <span
    class=\"n\">tags</span> <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">tags</span><span class=\"p\">,</span> <span
    class=\"nb\">list</span><span class=\"p\">)</span> <span class=\"k\">else</span>
    <span class=\"p\">[</span><span class=\"n\">tags</span><span class=\"p\">]</span>\n
    \                           <span class=\"k\">for</span> <span class=\"n\">a</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">posts</span>\n
    \                           <span class=\"k\">if</span> <span class=\"p\">(</span><span
    class=\"n\">tags</span> <span class=\"o\">:=</span> <span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">attr</span><span class=\"p\">,</span> <span class=\"kc\">None</span><span
    class=\"p\">))</span> <span class=\"ow\">is</span> <span class=\"ow\">not</span>
    <span class=\"kc\">None</span>\n                        <span class=\"p\">],</span>\n
    \                   <span class=\"p\">),</span>\n                <span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">)</span> <span
    class=\"o\">&gt;</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">grid</span><span class=\"o\">.</span><span class=\"n\">add_row</span><span
    class=\"p\">()</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_row</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;[bold gold1]</span><span class=\"si\">{</span><span class=\"n\">attr</span><span
    class=\"o\">.</span><span class=\"n\">upper</span><span class=\"p\">()</span><span
    class=\"si\">}</span><span class=\"s2\">[/]&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">for</span> <span class=\"n\">post</span><span
    class=\"p\">,</span> <span class=\"n\">count</span> <span class=\"ow\">in</span>
    <span class=\"n\">Counter</span><span class=\"p\">(</span><span class=\"n\">posts</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">most_common</span><span
    class=\"p\">():</span>\n                        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_row</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s1\">&#39;</span><span class=\"si\">{</span><span class=\"n\">count</span><span
    class=\"si\">}</span><span class=\"s1\"> </span><span class=\"si\">{</span><span
    class=\"s2\">&quot; &quot;</span><span class=\"o\">*</span><span class=\"p\">(</span><span
    class=\"mi\">3</span><span class=\"o\">-</span><span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">count</span><span class=\"p\">)))</span><span class=\"si\">}</span><span
    class=\"s1\"> </span><span class=\"si\">{</span><span class=\"n\">post</span><span
    class=\"si\">}</span><span class=\"s1\">&#39;</span><span class=\"p\">)</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">__rich__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Union</span><span class=\"p\">[</span><span
    class=\"n\">Panel</span><span class=\"p\">,</span> <span class=\"n\">Table</span><span
    class=\"p\">]:</span>\n                <span class=\"n\">grid</span> <span class=\"o\">=</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">get_grid</span><span
    class=\"p\">()</span>\n\n                <span class=\"k\">if</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">simple</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">grid</span>\n
    \               <span class=\"k\">else</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">Panel</span><span class=\"p\">(</span>\n
    \                       <span class=\"n\">grid</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">title</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;[gold1]summary[/]&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">border_style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;magenta&quot;</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">expand</span><span class=\"o\">=</span><span class=\"kc\">False</span><span
    class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>configure <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">configure
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">get_summary</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_summary</span>\n                <span class=\"k\">except</span> <span
    class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_summary</span><span
    class=\"p\">:</span> <span class=\"n\">Summary</span> <span class=\"o\">=</span>
    <span class=\"n\">Summary</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_summary</span>\n\n
    \           <span class=\"kn\">from</span> <span class=\"nn\">markata</span> <span
    class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n            <span
    class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">summary</span>
    <span class=\"o\">=</span> <span class=\"nb\">property</span><span class=\"p\">(</span><span
    class=\"n\">get_summary</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli
    <em class='small'>function</em></h2>\nMarkata hook to implement base cli commands.</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">cli
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">cli</span><span class=\"p\">(</span><span class=\"n\">app</span><span
    class=\"p\">:</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Typer</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
    class=\"w\">            </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">            Markata hook to implement base cli commands.</span>\n<span
    class=\"sd\">            &quot;&quot;&quot;</span>\n            <span class=\"n\">summary_app</span>
    <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Typer</span><span class=\"p\">()</span>\n            <span class=\"n\">app</span><span
    class=\"o\">.</span><span class=\"n\">add_typer</span><span class=\"p\">(</span><span
    class=\"n\">summary_app</span><span class=\"p\">,</span> <span class=\"n\">name</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;summary&quot;</span><span class=\"p\">)</span>\n\n
    \           <span class=\"nd\">@summary_app</span><span class=\"o\">.</span><span
    class=\"n\">callback</span><span class=\"p\">(</span><span class=\"n\">invoke_without_command</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">summary</span><span
    class=\"p\">():</span>\n                <span class=\"s2\">&quot;show the application
    summary&quot;</span>\n                <span class=\"kn\">from</span> <span class=\"nn\">rich</span>
    <span class=\"kn\">import</span> <span class=\"nb\">print</span>\n\n                <span
    class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
    class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span>\n\n                <span class=\"nb\">print</span><span
    class=\"p\">(</span><span class=\"n\">Summary</span><span class=\"p\">(</span><span
    class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">simple</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>init</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>init</strong>
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"fm\">__init__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">m</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">,</span> <span class=\"n\">simple</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"kc\">False</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span> <span
    class=\"o\">=</span> <span class=\"n\">m</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">simple</span> <span class=\"o\">=</span>
    <span class=\"n\">simple</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='get_grid'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_grid <em class='small'>method</em></h2>\ncreate
    a rich grid to display the summary</p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">get_grid <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">get_grid</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"s2\">&quot;create a rich grid
    to display the summary&quot;</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">grid</span> <span class=\"o\">=</span> <span
    class=\"n\">Table</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
    class=\"p\">(</span><span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">)</span>\n\n                <span class=\"k\">for</span>
    <span class=\"n\">filter_count</span> <span class=\"ow\">in</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">summary</span><span
    class=\"o\">.</span><span class=\"n\">filter_count</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">filter_count</span><span class=\"p\">(</span><span class=\"n\">filter_count</span><span
    class=\"p\">)</span>\n\n                <span class=\"k\">for</span> <span class=\"n\">attr</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"o\">.</span><span class=\"n\">summary</span><span class=\"o\">.</span><span
    class=\"n\">grid_attr</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid_attr</span><span
    class=\"p\">(</span><span class=\"n\">attr</span><span class=\"p\">)</span>\n\n
    \               <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">grid</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='filter_count' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>filter_count <em class='small'>method</em></h2>\nadd a row in the grid
    for the number of items in a filter config</p>\n<div class=\"admonition source
    is-collapsible collapsible-open\">\n<p class=\"admonition-title\">filter_count
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">filter_count</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">,</span>\n                <span class=\"n\">fc</span><span
    class=\"p\">:</span> <span class=\"n\">FilterCount</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"s2\">&quot;add a row in the
    grid for the number of items in a filter config&quot;</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">grid</span><span
    class=\"o\">.</span><span class=\"n\">add_row</span><span class=\"p\">(</span>\n
    \                   <span class=\"sa\">f</span><span class=\"s2\">&quot;[</span><span
    class=\"si\">{</span><span class=\"n\">fc</span><span class=\"o\">.</span><span
    class=\"n\">color</span><span class=\"si\">}</span><span class=\"s2\">]</span><span
    class=\"si\">{</span><span class=\"nb\">len</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
    class=\"nb\">filter</span><span class=\"o\">=</span><span class=\"n\">fc</span><span
    class=\"o\">.</span><span class=\"n\">filter</span><span class=\"p\">))</span><span
    class=\"si\">}</span><span class=\"s2\">[/] </span><span class=\"si\">{</span><span
    class=\"n\">fc</span><span class=\"o\">.</span><span class=\"n\">name</span><span
    class=\"si\">}</span><span class=\"s2\">&quot;</span>\n                <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='grid_attr' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>grid_attr <em class='small'>method</em></h2>\nadd attribute the the object
    grid</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">grid_attr <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">grid_attr</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">,</span> <span class=\"n\">attr</span><span class=\"p\">:</span> <span
    class=\"nb\">str</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"s2\">&quot;add attribute the the object grid&quot;</span>\n                <span
    class=\"n\">posts</span> <span class=\"o\">=</span> <span class=\"nb\">list</span><span
    class=\"p\">(</span>\n                    <span class=\"n\">flatten</span><span
    class=\"p\">(</span>\n                        <span class=\"p\">[</span>\n                            <span
    class=\"n\">tags</span> <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
    class=\"p\">(</span><span class=\"n\">tags</span><span class=\"p\">,</span> <span
    class=\"nb\">list</span><span class=\"p\">)</span> <span class=\"k\">else</span>
    <span class=\"p\">[</span><span class=\"n\">tags</span><span class=\"p\">]</span>\n
    \                           <span class=\"k\">for</span> <span class=\"n\">a</span>
    <span class=\"ow\">in</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">m</span><span class=\"o\">.</span><span class=\"n\">posts</span>\n
    \                           <span class=\"k\">if</span> <span class=\"p\">(</span><span
    class=\"n\">tags</span> <span class=\"o\">:=</span> <span class=\"n\">a</span><span
    class=\"o\">.</span><span class=\"n\">get</span><span class=\"p\">(</span><span
    class=\"n\">attr</span><span class=\"p\">,</span> <span class=\"kc\">None</span><span
    class=\"p\">))</span> <span class=\"ow\">is</span> <span class=\"ow\">not</span>
    <span class=\"kc\">None</span>\n                        <span class=\"p\">],</span>\n
    \                   <span class=\"p\">),</span>\n                <span class=\"p\">)</span>\n
    \               <span class=\"k\">if</span> <span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"n\">posts</span><span class=\"p\">)</span> <span
    class=\"o\">&gt;</span> <span class=\"mi\">0</span><span class=\"p\">:</span>\n
    \                   <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">grid</span><span class=\"o\">.</span><span class=\"n\">add_row</span><span
    class=\"p\">()</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_row</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s2\">&quot;[bold gold1]</span><span class=\"si\">{</span><span class=\"n\">attr</span><span
    class=\"o\">.</span><span class=\"n\">upper</span><span class=\"p\">()</span><span
    class=\"si\">}</span><span class=\"s2\">[/]&quot;</span><span class=\"p\">)</span>\n
    \                   <span class=\"k\">for</span> <span class=\"n\">post</span><span
    class=\"p\">,</span> <span class=\"n\">count</span> <span class=\"ow\">in</span>
    <span class=\"n\">Counter</span><span class=\"p\">(</span><span class=\"n\">posts</span><span
    class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">most_common</span><span
    class=\"p\">():</span>\n                        <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">grid</span><span class=\"o\">.</span><span
    class=\"n\">add_row</span><span class=\"p\">(</span><span class=\"sa\">f</span><span
    class=\"s1\">&#39;</span><span class=\"si\">{</span><span class=\"n\">count</span><span
    class=\"si\">}</span><span class=\"s1\"> </span><span class=\"si\">{</span><span
    class=\"s2\">&quot; &quot;</span><span class=\"o\">*</span><span class=\"p\">(</span><span
    class=\"mi\">3</span><span class=\"o\">-</span><span class=\"nb\">len</span><span
    class=\"p\">(</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"n\">count</span><span class=\"p\">)))</span><span class=\"si\">}</span><span
    class=\"s1\"> </span><span class=\"si\">{</span><span class=\"n\">post</span><span
    class=\"si\">}</span><span class=\"s1\">&#39;</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'><strong>rich</strong> <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\"><strong>rich</strong>
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Union</span><span
    class=\"p\">[</span><span class=\"n\">Panel</span><span class=\"p\">,</span> <span
    class=\"n\">Table</span><span class=\"p\">]:</span>\n                <span class=\"n\">grid</span>
    <span class=\"o\">=</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">get_grid</span><span class=\"p\">()</span>\n\n                <span
    class=\"k\">if</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">simple</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">grid</span>\n                <span
    class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"k\">return</span>
    <span class=\"n\">Panel</span><span class=\"p\">(</span>\n                        <span
    class=\"n\">grid</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"s2\">&quot;[gold1]summary[/]&quot;</span><span
    class=\"p\">,</span>\n                        <span class=\"n\">border_style</span><span
    class=\"o\">=</span><span class=\"s2\">&quot;magenta&quot;</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">False</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='get_summary'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_summary <em
    class='small'>function</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">get_summary <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">get_summary</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_summary</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_summary</span><span
    class=\"p\">:</span> <span class=\"n\">Summary</span> <span class=\"o\">=</span>
    <span class=\"n\">Summary</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span>\n                    <span class=\"k\">return</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_summary</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='summary' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>summary <em class='small'>function</em></h2>\nshow the application summary</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">summary
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">summary</span><span class=\"p\">():</span>\n                <span
    class=\"s2\">&quot;show the application summary&quot;</span>\n                <span
    class=\"kn\">from</span> <span class=\"nn\">rich</span> <span class=\"kn\">import</span>
    <span class=\"nb\">print</span>\n\n                <span class=\"n\">markata</span><span
    class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
    class=\"n\">quiet</span> <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n
    \               <span class=\"nb\">print</span><span class=\"p\">(</span><span
    class=\"n\">Summary</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">,</span> <span class=\"n\">simple</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">))</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/cli/summary
title: Summary.Py


---

## Run it

``` bash
python -m markata.cli.summary
```

## Configuration

There are two main things currently supported by summary, It can count the
number of posts based on a filter (`filter_count'), and it can automatically
list all the values of an attribute and the number of posts that have that
attribute (`grid_attr`).

### grid_attr

`grid_attr` will map over all posts, find all values for each attribute
configured, then report the number of posts for each value.

``` toml
[markata.summary]
grid_attr = ['tags', 'series']
```

Example output that will be shown in the summary, counting all posts with each
tag value.

```
TAGS
247  python
90   linux
68   cli
49   kedro
46   bash
```

### filter_counts

`filter_count` will pass a given filter into `markata.map` and return the number
of posts.

```
[[markata.summary.filter_count]]
name='drafts'
filter="published == 'False'"
color='red'

[[markata.summary.filter_count]]
name='articles'
color='dark_orange'

[[markata.summary.filter_count]]
name='py_modules'
filter='"plugin" not in slug and "docs" not in str(path)'
color="yellow1"

[markata.summary.filter_count.published]
filter="published == 'True'"
color='green1'

[markata.summary.filter_count.plugins]
filter='"plugin" in slug and "docs" not in str(path)'
color="blue"

[markata.summary.filter_count.docs]
filter="'docs' in str(path)"
color='purple'
```

Example output might look like this, showing the number of posts that contained
within each filter specified.

```
8 drafts
66 articles
20 py_modules
58 published
38 plugins
8 docs
```


!! class <h2 id='FilterCount' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>FilterCount <em class='small'>class</em></h2>

???+ source "FilterCount <em class='small'>source</em>"

```python

        class FilterCount(pydantic.BaseModel):
            name: str
            filter: str = "True"
            color: str = "white"
```


!! class <h2 id='SummaryConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>SummaryConfig <em class='small'>class</em></h2>

???+ source "SummaryConfig <em class='small'>source</em>"

```python

        class SummaryConfig(pydantic.BaseModel):
            grid_attr: List[str] = ["tags", "series"]
            filter_count: List[FilterCount] = FilterCount(
                name="drafts", filter="published == 'False'", color="red"
            )
```


!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>

???+ source "Config <em class='small'>source</em>"

```python

        class Config(pydantic.BaseModel):
            summary: SummaryConfig = SummaryConfig()
```


!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>

???+ source "config_model <em class='small'>source</em>"

```python

        def config_model(markata: "Markata") -> None:
            markata.config_models.append(Config)
```


!! class <h2 id='Summary' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Summary <em class='small'>class</em></h2>

???+ source "Summary <em class='small'>source</em>"

```python

        class Summary:
            def __init__(self, m: "Markata", simple: bool = False) -> None:
                self.m = m
                self.simple = simple

            def get_grid(self) -> None:
                "create a rich grid to display the summary"
                self.grid = Table.grid(expand=True)

                for filter_count in self.m.config.summary.filter_count:
                    self.filter_count(filter_count)

                for attr in self.m.config.summary.grid_attr:
                    self.grid_attr(attr)

                return self.grid

            def filter_count(
                self,
                fc: FilterCount,
            ) -> None:
                "add a row in the grid for the number of items in a filter config"
                self.grid.add_row(
                    f"[{fc.color}]{len(self.m.map(filter=fc.filter))}[/] {fc.name}"
                )

            def grid_attr(self, attr: str) -> None:
                "add attribute the the object grid"
                posts = list(
                    flatten(
                        [
                            tags if isinstance(tags, list) else [tags]
                            for a in self.m.posts
                            if (tags := a.get(attr, None)) is not None
                        ],
                    ),
                )
                if len(posts) > 0:
                    self.grid.add_row()
                    self.grid.add_row(f"[bold gold1]{attr.upper()}[/]")
                    for post, count in Counter(posts).most_common():
                        self.grid.add_row(f'{count} {" "*(3-len(str(count)))} {post}')

            def __rich__(self) -> Union[Panel, Table]:
                grid = self.get_grid()

                if self.simple:
                    return grid
                else:
                    return Panel(
                        grid,
                        title="[gold1]summary[/]",
                        border_style="magenta",
                        expand=False,
                    )
```


!! function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>configure <em class='small'>function</em></h2>

???+ source "configure <em class='small'>source</em>"

```python

        def configure(markata: "Markata") -> None:
            def get_summary(self):
                try:
                    return self._summary
                except AttributeError:
                    self._summary: Summary = Summary(self)
                    return self._summary

            from markata import Markata

            Markata.summary = property(get_summary)
```


!! function <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2>
    Markata hook to implement base cli commands.
???+ source "cli <em class='small'>source</em>"

```python

        def cli(app: typer.Typer, markata: "Markata") -> None:
            """
            Markata hook to implement base cli commands.
            """
            summary_app = typer.Typer()
            app.add_typer(summary_app, name="summary")

            @summary_app.callback(invoke_without_command=True)
            def summary():
                "show the application summary"
                from rich import print

                markata.console.quiet = True

                print(Summary(markata, simple=True))
```


!! method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>

???+ source "__init__ <em class='small'>source</em>"

```python

        def __init__(self, m: "Markata", simple: bool = False) -> None:
                self.m = m
                self.simple = simple
```


!! method <h2 id='get_grid' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_grid <em class='small'>method</em></h2>
    create a rich grid to display the summary
???+ source "get_grid <em class='small'>source</em>"

```python

        def get_grid(self) -> None:
                "create a rich grid to display the summary"
                self.grid = Table.grid(expand=True)

                for filter_count in self.m.config.summary.filter_count:
                    self.filter_count(filter_count)

                for attr in self.m.config.summary.grid_attr:
                    self.grid_attr(attr)

                return self.grid
```


!! method <h2 id='filter_count' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>filter_count <em class='small'>method</em></h2>
    add a row in the grid for the number of items in a filter config
???+ source "filter_count <em class='small'>source</em>"

```python

        def filter_count(
                self,
                fc: FilterCount,
            ) -> None:
                "add a row in the grid for the number of items in a filter config"
                self.grid.add_row(
                    f"[{fc.color}]{len(self.m.map(filter=fc.filter))}[/] {fc.name}"
                )
```


!! method <h2 id='grid_attr' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>grid_attr <em class='small'>method</em></h2>
    add attribute the the object grid
???+ source "grid_attr <em class='small'>source</em>"

```python

        def grid_attr(self, attr: str) -> None:
                "add attribute the the object grid"
                posts = list(
                    flatten(
                        [
                            tags if isinstance(tags, list) else [tags]
                            for a in self.m.posts
                            if (tags := a.get(attr, None)) is not None
                        ],
                    ),
                )
                if len(posts) > 0:
                    self.grid.add_row()
                    self.grid.add_row(f"[bold gold1]{attr.upper()}[/]")
                    for post, count in Counter(posts).most_common():
                        self.grid.add_row(f'{count} {" "*(3-len(str(count)))} {post}')
```


!! method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__rich__ <em class='small'>method</em></h2>

???+ source "__rich__ <em class='small'>source</em>"

```python

        def __rich__(self) -> Union[Panel, Table]:
                grid = self.get_grid()

                if self.simple:
                    return grid
                else:
                    return Panel(
                        grid,
                        title="[gold1]summary[/]",
                        border_style="magenta",
                        expand=False,
                    )
```


!! function <h2 id='get_summary' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_summary <em class='small'>function</em></h2>

???+ source "get_summary <em class='small'>source</em>"

```python

        def get_summary(self):
                try:
                    return self._summary
                except AttributeError:
                    self._summary: Summary = Summary(self)
                    return self._summary
```


!! function <h2 id='summary' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>summary <em class='small'>function</em></h2>
    show the application summary
???+ source "summary <em class='small'>source</em>"

```python

        def summary():
                "show the application summary"
                from rich import print

                markata.console.quiet = True

                print(Summary(markata, simple=True))
```

