""" The LifeCycle is a core component for the internal workings of Markata.  It
sets fourth the hooks available, the methods to run them on the Markata
instance, and the order they run in.

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


    * configure - load and fix configuration
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
