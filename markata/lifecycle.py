from enum import Enum, auto


class LifeCycle(Enum):
    configure = auto()
    glob = auto()
    pre_render = auto()
    render = auto()
    post_render = auto()
    save = auto()

    def __gt__(self, other):
        try:
            return self.value > other.value
        except AttributeError:
            return self.value > other

    def __lt__(self, other):
        try:
            return self.value < other.value
        except AttributeError:
            return self.value < other

    def __eq__(self, other):
        try:
            return self.value == other.value
        except AttributeError:
            return self.value == other
