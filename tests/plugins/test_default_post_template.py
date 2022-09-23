from pathlib import Path

import toml
from conftest import set_directory

from markata import Markata
from markata.plugins import post_template


def test_head_config_text_str(tmp_path):
    with set_directory(tmp_path):
        m = Markata()
        m.config["head"] = {}
        m.config["head"]["text"] = "here"
        post_template.configure(m)
        assert m.config["head"]["text"] == "here"


def test_head_config_text_dict(tmp_path):
    with set_directory(tmp_path):
        m = Markata()
        m.config["head"] = {}
        m.config["head"]["text"] = [{"value": "one"}, {"value": "two"}]
        post_template.configure(m)
        assert m.config["head"]["text"] == "one\ntwo"


def test_head_config_text_str_toml(tmp_path):
    with set_directory(tmp_path):
        Path("markata.toml").write_text(
            toml.dumps({"markata": {"head": {"text": "here"}}})
        )
        m = Markata()
        post_template.configure(m)
        assert m.config["head"]["text"] == "here"


def test_head_config_text_list_toml(tmp_path):
    with set_directory(tmp_path):
        Path("markata.toml").write_text(
            toml.dumps(
                {"markata": {"head": {"text": [{"value": "one"}, {"value": "two"}]}}}
            )
        )
        m = Markata()
        post_template.configure(m)
        assert m.config["head"]["text"] == "one\ntwo"
