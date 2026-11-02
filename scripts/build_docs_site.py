#!/usr/bin/env python3
"""Build a small, JavaScript-optional static documentation site.

Source: docs/**/*.md. Output: docs-site/_site/.
Requires the Python-Markdown package for authoring, not for running Markata.
"""
from __future__ import annotations

import html
import os
import re
import shutil
from pathlib import Path
from urllib.parse import urlsplit, urlunsplit

from markdown import markdown

ROOT = Path(__file__).resolve().parent.parent
DOCS = ROOT / "docs"
SITE = ROOT / "docs-site"
OUT = SITE / "_site"
BASE = os.environ.get("DOCS_BASE_PATH", "/markata/")
if not (BASE.startswith("/") and BASE.endswith("/")):
    raise SystemExit("DOCS_BASE_PATH must start and end with /")
if "//" in BASE:
    raise SystemExit("DOCS_BASE_PATH must not contain //")


def page_url(source: Path) -> str:
    relative = source.relative_to(DOCS).with_suffix("")
    if relative.name == "index":
        relative = relative.parent
    suffix = relative.as_posix().strip(".")
    return BASE + "docs/" + (suffix + "/" if suffix else "")


def page_output(source: Path) -> Path:
    relative = source.relative_to(DOCS).with_suffix("")
    if relative.name == "index":
        relative = relative.parent
    return OUT / "docs" / relative / "index.html"


def metadata(markdown_text: str) -> tuple[dict[str, str], str]:
    if not markdown_text.startswith("---\n"):
        return {}, markdown_text
    marker = markdown_text.find("\n---\n", 4)
    if marker == -1:
        return {}, markdown_text
    front, content = markdown_text[4:marker], markdown_text[marker + 5:]
    fields = {}
    for line in front.splitlines():
        match = re.match(r"^([a-z_]+):\s*(.*)$", line)
        if match:
            fields[match[1]] = match[2].strip().strip("'\"")
    return fields, content


def rewrite_url(href: str, source: Path, known: dict[str, Path]) -> str:
    if href.startswith(("#", "mailto:", "tel:", "data:", "//")):
        return href
    parsed = urlsplit(href)
    if parsed.scheme or parsed.netloc:
        return href
    target = parsed.path
    if target.startswith("/docs/"):
        replacement = BASE + target.lstrip("/")
    elif target.startswith("/"):
        # Other site-root links can refer to user sites; leave them unchanged.
        return href
    elif target.endswith(".md"):
        match = (source.parent / target).resolve()
        if match.is_relative_to(DOCS.resolve()) and match.exists():
            replacement = page_url(match)
        else:
            return href
    else:
        return href
    return urlunsplit(("", "", replacement, parsed.query, parsed.fragment))


def render(source: Path, known: dict[str, Path]) -> tuple[str, str]:
    info, body = metadata(source.read_text(encoding="utf-8"))
    title = info.get("title", source.stem.replace("-", " ").title())
    description = info.get("description", "Markata documentation")
    def wiki(match: re.Match[str]) -> str:
        parts = match.group(1).split("|", 1)
        name, label = parts[0].strip(), (parts[-1] if len(parts) > 1 else parts[0]).strip()
        target = known.get(name.lower())
        if target is None:
            return label
        return f"[{label}]({page_url(target)})"
    body = re.sub(r"\[\[([^\[\]]+)\]\]", wiki, body)
    rendered = markdown(body, extensions=["fenced_code", "tables", "toc", "sane_lists"])
    def fix_href(match: re.Match[str]) -> str:
        quote, old = match.group(1), match.group(2)
        new = rewrite_url(html.unescape(old), source, known)
        return 'href=' + quote + html.escape(new, quote=True) + quote
    rendered = re.sub(r'href=(["\'])(.*?)\1', fix_href, rendered)
    return title, description, rendered


def shell(title: str, description: str, content: str) -> str:
    nav = [
        ("Quickstart", "docs/quickstart/"),
        ("Installation", "docs/installation/"),
        ("Agents", "docs/guides/agent-skills/"),
        ("Authoring", "docs/guides/authoring/"),
        ("All docs", "docs/"),
        ("GitHub ↗", "https://github.com/WaylonWalker/markata"),
    ]
    links = "".join(
        f'<a href="{url if url.startswith("https:") else BASE + url}">{html.escape(label)}</a>'
        for label, url in nav
    )
    return f"""<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="description" content="{html.escape(description, quote=True)}">
<title>{html.escape(title)} — Markata</title><link rel="stylesheet" href="{BASE}assets/site.css"></head>
<body><header class="site-header shell"><a href="{BASE}" class="wordmark">markata<span class="period">.</span></a><nav aria-label="Main navigation">{links}</nav></header>
<main class="shell docs-layout"><aside class="docs-nav" aria-label="Documentation sections">
<a href="{BASE}">← Overview</a>{links}</aside>
<article class="docs-article">{content}</article></main>
<footer class="site-footer shell"><span>Markata · Markdown to website</span><a href="{BASE}docs/guides/name-migration/">Naming migration</a></footer>
</body></html>"""


def main() -> None:
    if OUT.exists():
        shutil.rmtree(OUT)
    (OUT / "assets").mkdir(parents=True)
    shutil.copyfile(SITE / "site.css", OUT / "assets" / "site.css")
    # Keep the docs homepage on the same native tab controller as Markata sites.
    shutil.copyfile(
        ROOT / "pkg" / "themes" / "default" / "static" / "js" / "md-tabs.js",
        OUT / "assets" / "md-tabs.js",
    )
    (OUT / "index.html").write_text(
        (SITE / "home.html").read_text(encoding="utf-8").replace("@@BASE@@", BASE),
        encoding="utf-8",
    )
    sources = sorted(DOCS.rglob("*.md"))
    known: dict[str, Path] = {}
    for source in sources:
        known.setdefault(source.stem.lower(), source)
    for source in sources:
        title, description, content = render(source, known)
        dest = page_output(source)
        dest.parent.mkdir(parents=True, exist_ok=True)
        dest.write_text(shell(title, description, content), encoding="utf-8")
    if (DOCS / "images").exists():
        shutil.copytree(DOCS / "images", OUT / "docs" / "images", dirs_exist_ok=True)
    (OUT / ".nojekyll").touch()
    print(f"Built {len(sources)} documentation pages in {OUT}")


if __name__ == "__main__":
    main()
