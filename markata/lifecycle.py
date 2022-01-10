from enum import Enum, auto
from functools import total_ordering


@total_ordering
class LifeCycle(Enum):

    configure = auto()
    glob = auto()
    load = auto()
    pre_render = auto()
    render = auto()
    post_render = auto()
    save = auto()

    def __lt__(self, other: object) -> bool:
        if isinstance(other, LifeCycle):
            return self.value < other.value
        if isinstance(other, int):
            return self.value < other
        return NotImplemented

    def __eq__(self, other: object) -> bool:
        if isinstance(other, LifeCycle):
            return self.value == other.value
        if isinstance(other, int):
            return self.value == other
        return NotImplemented
