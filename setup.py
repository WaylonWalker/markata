"""
markata uses setup tools for packaging.

To Build markata as a Python package

    $ python setup.py sdist bdist_wheel --bdist-dir ~/temp/bdistwheel

Regular install

    $ pip install -e .

To setup local Development

    $ pip install -e .

"""
from pathlib import Path

from setuptools import find_packages, setup

NAME = "markata"

README = (Path(__file__).parent / "README.md").read_text(encoding="utf-8")

with open("requirements.txt", "r", encoding="utf-8") as f:
    requires = [x.strip() for x in f if x.strip()]

with open("requirements_dev.txt", "r", encoding="utf-8") as f:
    dev_requires = [x.strip() for x in f if x.strip()]

setup(
    name=NAME,
    version="0.0.1",
    packages=find_packages(),
    entry_points={"console_scripts": ["markata=markata.cli.cli:cli"]},
    platforms="any",
    license="MIT",
    install_requires=requires,
    extras_require={"dev": dev_requires},
    include_package_data=True,
)
