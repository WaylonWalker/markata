import os

from markata import Markata
from markata.cli import cli

if __name__ == "__main__":
    m = Markata()

    if "ipython" not in os.environ.get("_", ""):
        cli()
