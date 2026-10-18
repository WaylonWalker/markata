#!/usr/bin/env python3
"""Markata docs live builder for Argo: Git polling, atomic releases, authenticated UI.

This builds docs-site/ (Python Markdown renderer), NOT the main Markata CLI
site. Source is fetched from the canonical public repository on every change.
"""
from __future__ import annotations

import html
import ipaddress
import json
import os
import re
import shutil
import subprocess
import sys
import threading
import time
from datetime import datetime, timezone
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from urllib.parse import parse_qs, urlsplit

REMOTE = os.getenv("DOCS_REPOSITORY", "https://github.com/WaylonWalker/markata.git")
BRANCH = os.getenv("DOCS_BRANCH", "main")
SOURCE = Path(os.getenv("DOCS_SOURCE_DIR", "/data/source"))
SITE = Path(os.getenv("DOCS_SITE_DIR", "/data/site"))
PUBLIC_ORIGIN = os.getenv("DOCS_BUILDER_ORIGIN", "https://build.markata.waylonwalker.com")
POLL = max(20, int(os.getenv("DOCS_POLL_SECONDS", "45")))
KEEP = max(2, int(os.getenv("DOCS_KEEP_RELEASES", "7")))
TIMEOUT = max(60, int(os.getenv("DOCS_BUILD_TIMEOUT_SECONDS", "900")))
TRUSTED_PROXY_CIDRS = [
    ipaddress.ip_network(s.strip())
    for s in os.getenv("DOCS_TRUSTED_PROXY_CIDRS", "10.42.0.0/24,10.42.2.0/24").split(",")
    if s.strip()
]
STATE_FILE = SITE / ".docs-builder" / "state.json"
LOG_DIR = SITE / ".docs-builder" / "logs"
VERSIONS = SITE / "releases"
CURRENT = SITE / "current"


def stamp() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def command(args: list[str], *, cwd: Path | None = None, timeout: int = 60, env=None) -> str:
    result = subprocess.run(
        args, cwd=cwd, env=env, text=True,
        stdout=subprocess.PIPE, stderr=subprocess.STDOUT, timeout=timeout, check=True,
    )
    return result.stdout


def safe_release(value: str) -> bool:
    return bool(re.fullmatch(r"[A-Za-z0-9][A-Za-z0-9._-]{0,120}", value))


class Builder:
    def __init__(self):
        LOG_DIR.mkdir(parents=True, exist_ok=True)
        VERSIONS.mkdir(parents=True, exist_ok=True)
        self.lock = threading.RLock()
        self.wakeup = threading.Event()
        self.running = False
        self.status = "idle"
        self.error = ""
        self.history = []
        self.last_sha = ""
        self.last_run = ""
        if STATE_FILE.is_file():
            try:
                data = json.loads(STATE_FILE.read_text())
                self.last_sha = str(data.get("last_sha", ""))
                self.last_run = str(data.get("last_run", ""))
                self.history = data.get("history", [])[-30:]
            except (ValueError, OSError):
                pass

    def save(self):
        payload = json.dumps({
            "last_sha": self.last_sha, "last_run": self.last_run, "history": self.history[-30:],
        }, indent=2)
        tmp = STATE_FILE.with_suffix(".tmp")
        tmp.write_text(payload)
        os.replace(tmp, STATE_FILE)

    def log_record(self, status, sha, message):
        self.last_run = stamp()
        self.history.append({"time": self.last_run, "status": status, "sha": sha, "message": message[:400]})
        self.history = self.history[-30:]
        self.save()

    def queue(self):
        self.wakeup.set()

    def publish(self, source: Path, sha: str):
        name = datetime.now(timezone.utc).strftime("%Y%m%dT%H%M%SZ") + "-" + sha[:10]
        candidate = VERSIONS / (".building-" + name)
        final = VERSIONS / name
        if candidate.exists():
            shutil.rmtree(candidate)
        shutil.copytree(source, candidate, symlinks=False)
        for expected in ("index.html", "docs/quickstart/index.html", "docs/guides/authoring/index.html", "assets/site.css"):
            if not (candidate / expected).is_file():
                raise RuntimeError("Missing generated documentation: " + expected)
        os.replace(candidate, final)
        self.set_current(name)
        self.prune()
        return name

    def set_current(self, name: str):
        if not safe_release(name) or not (VERSIONS / name / "index.html").is_file():
            raise ValueError("Release not found")
        tmp = SITE / ".current-next"
        if tmp.is_symlink():
            tmp.unlink()
        tmp.symlink_to("releases/" + name)
        os.replace(tmp, CURRENT)

    def prune(self):
        all_releases = sorted((p for p in VERSIONS.iterdir() if p.is_dir() and not p.name.startswith(".building-")), key=lambda p: p.name, reverse=True)
        active = CURRENT.resolve() if CURRENT.is_symlink() else None
        for old in all_releases[KEEP:]:
            if old.resolve() != active:
                shutil.rmtree(old)

    def build(self, force=False):
        with self.lock:
            if self.running:
                return
            self.running = True
            self.status = "checking"
            self.error = ""
        sha = "unknown"
        try:
            remote_sha = command(["git", "ls-remote", REMOTE, f"refs/heads/{BRANCH}"], timeout=45).split()[0]
            sha = remote_sha
            if not force and remote_sha == self.last_sha and CURRENT.is_symlink() and CURRENT.exists():
                with self.lock:
                    self.status = "up to date"
                return
            if not (SOURCE / ".git").exists():
                SOURCE.parent.mkdir(parents=True, exist_ok=True)
                if SOURCE.exists():
                    shutil.rmtree(SOURCE)
                command(["git", "clone", "--depth=1", "--branch", BRANCH, REMOTE, str(SOURCE)], timeout=180)
            else:
                command(["git", "fetch", "--depth=1", "origin", BRANCH], cwd=SOURCE, timeout=120)
                command(["git", "reset", "--hard", "FETCH_HEAD"], cwd=SOURCE, timeout=90)
            with self.lock:
                self.status = "building"
            env = os.environ.copy()
            env.update(DOCS_BASE_PATH="/", PYTHONDONTWRITEBYTECODE="1")
            output = command([sys.executable, "scripts/build_docs_site.py"], cwd=SOURCE, timeout=TIMEOUT, env=env)
            name = self.publish(SOURCE / "docs-site" / "_site", sha)
            log_file = LOG_DIR / (name + ".log")
            log_file.write_text(output)
            with self.lock:
                self.last_sha = sha
                self.log_record("success", sha, f"Published {name}")
                self.status = "published"
        except Exception as exc:
            with self.lock:
                self.error = str(exc)
                self.status = "failed"
                self.log_record("failed", sha, str(exc))
            print(f"[builder] {stamp()} build failed: {exc}", file=sys.stderr, flush=True)
        finally:
            with self.lock:
                self.running = False

    def run(self):
        # Ensure initial build even when the image changes but the same Git SHA is
        # present; build(force=False) skips only if a valid persisted current exists.
        while True:
            force = self.wakeup.is_set()
            self.wakeup.clear()
            self.build(force=force)
            self.wakeup.wait(POLL)

    def rollback(self, name):
        with self.lock:
            if self.running:
                raise RuntimeError("Build in progress")
            self.set_current(name)
            self.log_record("rollback", self.last_sha, "Restored " + name)
            self.status = "rolled back"


class Handler(BaseHTTPRequestHandler):
    server: "BuilderServer"

    def allowed(self):
        try:
            ip = ipaddress.ip_address(self.client_address[0])
            return any(ip in net for net in TRUSTED_PROXY_CIDRS) and bool(self.headers.get("X-Hlab-User-Id", "").strip())
        except ValueError:
            return False

    def send_text(self, code, value, content_type="text/plain; charset=utf-8"):
        encoded = value.encode("utf-8")
        self.send_response(code)
        self.send_header("Content-Type", content_type)
        self.send_header("Cache-Control", "no-store")
        self.send_header("X-Content-Type-Options", "nosniff")
        self.send_header("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; form-action 'self'; base-uri 'none'; frame-ancestors 'none'")
        self.send_header("Content-Length", str(len(encoded)))
        self.end_headers()
        self.wfile.write(encoded)

    def do_GET(self):
        route = urlsplit(self.path).path
        if route == "/healthz":
            self.send_text(200, "ok\n")
            return
        if not self.allowed():
            self.send_text(403, "Access denied\n")
            return
        b = self.server.builder
        if route != "/":
            self.send_text(404, "Not found")
            return
        with b.lock:
            state, sha, error, last, running = b.status, b.last_sha, b.error, b.last_run, b.running
            history = list(b.history[-12:])
        active = CURRENT.resolve().name if CURRENT.exists() else "(no published release)"
        versions = sorted([p.name for p in VERSIONS.iterdir() if p.is_dir() and safe_release(p.name)], reverse=True)[:KEEP+2]
        esc = html.escape
        rows = "".join(f"<tr><td>{esc(r['time'])}</td><td>{esc(r['status'])}</td><td><code>{esc(r['sha'][:12])}</code></td><td>{esc(r['message'])}</td></tr>" for r in reversed(history))
        rollback_forms = "".join(
            f'<form method="post" action="/rollback"><input type="hidden" name="release" value="{esc(name)}"><button type="submit">Restore {esc(name)}</button></form>'
            for name in versions if name != active
        )
        body = f"""<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Markata Docs Builder</title><style>body{{max-width:860px;margin:3rem auto;padding:0 1rem;background:#121820;color:#e0e7e7;font:16px/1.6 system-ui}}
a{{color:#83ccbf}}code,pre{{font-family:monospace}}pre{{overflow-x:auto}}button{{padding:.5rem .8rem;margin:.25rem;background:#254c46;border:1px solid #4a8178;color:white;cursor:pointer}}
table{{width:100%;border-collapse:collapse;font-size:.85rem}}td,th{{padding:.6rem;border-bottom:1px solid #3d4b52;text-align:left;vertical-align:top}}
section{{margin:2rem 0}}.error{{color:#f9b4b4}}</style></head><body>
<h1>Markata Docs Builder</h1><p>Builds the documentation site from <code>WaylonWalker/markata</code> and publishes atomic releases to the homelab.</p>
<section><p>Status: <strong>{esc(state)}</strong>{' (busy)' if running else ''}</p>
<p>Commit: <code>{esc(sha[:12] or '(none)')}</code> · Last run: {esc(last or '(none)')} · Live release: <code>{esc(active)}</code></p>
<p class="error">{esc(error)}</p><form method="post" action="/build"><button type="submit">Build now</button></form>
<p><a href="https://markata.waylonwalker.com/">Open published docs ↗</a></p></section>
<section><h2>Previous builds</h2><table><thead><tr><th>When</th><th>Result</th><th>Commit</th><th>Message</th></tr></thead><tbody>{rows}</tbody></table></section>
<section><h2>Rollback</h2>{rollback_forms}</section></body></html>"""
        self.send_text(200, body, "text/html; charset=utf-8")

    def do_POST(self):
        if not self.allowed():
            self.send_text(403, "Access denied")
            return
        if self.headers.get("Origin", "") != PUBLIC_ORIGIN:
            self.send_text(403, "Invalid origin")
            return
        length = int(self.headers.get("Content-Length", "0"))
        if length > 4096 or length < 0:
            self.send_text(413, "Request too large")
            return
        form = parse_qs(self.rfile.read(length).decode("utf-8"))
        route = urlsplit(self.path).path
        if route == "/build":
            self.server.builder.queue()
        elif route == "/rollback":
            try:
                self.server.builder.rollback(form.get("release", [""])[0])
            except (ValueError, RuntimeError) as exc:
                self.send_text(409, str(exc))
                return
        else:
            self.send_text(404, "Not found")
            return
        self.send_response(303)
        self.send_header("Location", "/")
        self.send_header("Content-Length", "0")
        self.end_headers()


class BuilderServer(ThreadingHTTPServer):
    def __init__(self, addr, builder):
        super().__init__(addr, Handler)
        self.builder = builder


if __name__ == "__main__":
    builder = Builder()
    threading.Thread(target=builder.run, daemon=True).start()
    port = int(os.getenv("DOCS_BUILDER_PORT", "8080"))
    print(f"Docs builder listening at 0.0.0.0:{port}, polling every {POLL}s", flush=True)
    BuilderServer(("0.0.0.0", port), builder).serve_forever()
