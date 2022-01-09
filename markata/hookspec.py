"""Define hook specs."""
import functools
from typing import TYPE_CHECKING

import pluggy

from markata.lifecycle import LifeCycle

if TYPE_CHECKING:
    import typer

    from markata import Markata

HOOK_NAMESPACE = "markata"
hook_spec = pluggy.HookspecMarker(HOOK_NAMESPACE)
hook_impl = pluggy.HookimplMarker(HOOK_NAMESPACE)


class MarkataSpecs:
    """
    Namespace that defines all specifications for Load hooks.

    configure -> glob -> load -> render -> save
    """


@hook_spec
def generic_lifecycle_method(
    markata: "Markata",
):
    ...


@hook_spec
def cli_lifecycle_method(markata: "Markata", app: "typer.Typer"):
    "A Markata lifecycle methos that includes a typer app used for cli's"


for method in LifeCycle._member_map_:
    if "cli" in method:
        setattr(MarkataSpecs, method, cli_lifecycle_method)
    else:
        setattr(MarkataSpecs, method, generic_lifecycle_method)

registered_attrs = {}


def register_attr(*attrs):
    def decorator_register(func):

        for attr in attrs:
            if attr not in registered_attrs:
                registered_attrs[attr] = []
            registered_attrs[attr].append(
                {
                    "func": func,
                    "funcname": func.__code__.co_name,
                    "lifecycle": getattr(LifeCycle, func.__code__.co_name),
                }
            )

        @functools.wraps(func)
        def wrapper_register(markata, *args, **kwargs):
            return func(markata, *args, **kwargs)

        return wrapper_register

    return decorator_register
