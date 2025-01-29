"""The LifeCycle is a core component for the internal workings of Markata.  It
sets fourth the hooks available, the methods to run them on the Markata
instance, and the order they run in.

# Build Lifecycle

The build process follows these stages in order. Each stage runs all registered hooks
for that stage before proceeding to the next stage.

## 1. Configuration Stage

Sets up the build environment:


```python
# First: Define configuration models
# Example plugin:
from markata.hookspec import hook_impl

@hook_impl
def config_model(markata):
    markata.config_models.append(MyConfig)
```


``` python
# Second: Configure plugins
# Example plugin:
from markata.hookspec import hook_impl

@hook_impl
def configure(markata):
    markata.my_data = setup_resources()
```

## 2. Model Creation Stage

Defines content structure:

First: Define post model fragments

```python
from pydantic import BaseModel
from markata.hookspec import hook_impl

class MyPostFields(BaseModel):
    title: str
    tags: List[str]

# Example plugin:
@hook_impl
def post_model(markata):
    return MyPostFields
```

## 3. Content Discovery Stage

Finds and loads content:

```python
# First: Find content files
# Example plugin:
@hook_impl
def glob(markata):
    return list(Path("content").glob("**/*.md"))

# Second: Load content
# Example plugin:
@hook_impl
def load(markata):
    for path in markata.paths:
        content = path.read_text()
        markata.articles.append(parse_content(content))
```

## 4. Content Processing Stage

Transforms content:

```python
# Second: Convert markdown
# Example plugin:
@hook_impl
def render_markdown(markata):
    for article in markata.articles:
        article.html = convert_markdown(article.content)

# Third: Process content
# Example plugin:
@hook_impl
def render(markata):
    for article in markata.articles:
        article.html = apply_template(article.html)
```

## 5. Output Generation Stage

Saves processed content:

```python
# First: Save content
# Example plugin:
@hook_impl
def save(markata):
    for article in markata.articles:
        save_article(article)

# Finally: Clean up
# Example plugin:
@hook_impl
def teardown(markata):
    cleanup_resources()
```

# Hook Execution Order

Within each stage, hooks are executed in this order:
1. Hooks with tryfirst=True (earliest)
2. Hooks with no ordering specified
3. Hooks with trylast=True (latest)

Example ordering:
```python
@hook_impl(tryfirst=True)
def configure(markata):
    \"\"\"Runs first in configure stage\"\"\"
    setup_required_resources()

@hook_impl
def configure(markata):
    \"\"\"Runs in middle of configure stage\"\"\"
    setup_optional_features()

@hook_impl(trylast=True)
def configure(markata):
    \"\"\"Runs last in configure stage\"\"\"
    finalize_configuration()
```

# Error Handling

The lifecycle manager handles errors in hooks:
1. Logs errors with traceback
2. Continues execution if possible
3. Raises fatal errors that prevent build

Example error handling:
```python
@hook_impl
def render(markata):
    try:
        process_content()
    except NonFatalError:
        # Log and continue
        markata.logger.warning("Non-fatal error in render")
    except FatalError:
        # Stop build
        raise
```

# Parallel Processing

Some stages support parallel execution:
- render_markdown: Parallel markdown conversion
- render: Parallel template rendering
- save: Parallel file writing

Example parallel hook:
```python
@hook_impl
def render_markdown(markata):
    with ThreadPoolExecutor() as executor:
        futures = []
        for article in markata.articles:
            future = executor.submit(convert_markdown, article.content)
            futures.append((article, future))

        for article, future in futures:
            article.html = future.result()
```

See [[ hookspec ]] for detailed hook specifications and standard_config.py for
configuration options.

### Usage

``` python
from markata import Lifecycle

step = Lifecycle.glob
```

"""

from enum import Enum, auto
from functools import total_ordering


@total_ordering
class LifeCycle(Enum):
    """
    LifeCycle currently supports the following steps.

    * config_model - load configuration models from plugins
    * post_model - load post models from plugins
    * create_models - merge models from all plugins into markata.Post and markata.Plugin
    * load_config - load configuration
    * configure - load and fix configuration
    * validate_config - validate configuration
    * glob - find files
    * load - load files
    * validate_posts
    * pre_render - clean up files/metadata before render
    * render - render content
    * post_render - clean up rendered content
    * save - store results to disk
    * teardown - runs on exit

    """

    config_model = auto()
    post_model = auto()
    create_models = auto()
    load_config = auto()
    configure = auto()
    validate_config = auto()
    glob = auto()
    load = auto()
    pre_render = auto()
    render = auto()
    post_render = auto()
    save = auto()
    teardown = auto()

    def __lt__(self, other: object) -> bool:
        """
        Determine whether other is less than this instance.
        """
        if isinstance(other, LifeCycle):
            return self.value < other.value
        if isinstance(other, int):
            return self.value < other
        return NotImplemented

    def __eq__(self, other: object) -> bool:
        """
        Determine whether other is equal to this instance.
        """
        if isinstance(other, LifeCycle):
            return self.value == other.value
        if isinstance(other, int):
            return self.value == other
        return NotImplemented
