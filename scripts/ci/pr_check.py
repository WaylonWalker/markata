#!/usr/bin/env python3
"""Local merge preflight and shared fast PR Go checks."""
import argparse
import os
from pathlib import Path
import subprocess
import sys
import tempfile

DOCS = ("docs/", "docs-site/")
EDITORS = ("tests/editors/", "pkg/workbench/", "pkg/lsp/", "python/", "examples/wasm-workbench/")
RENDER = ("scripts/rendering-contract/", "spec/rendering-contract/", "pkg/renderingcontract/", "pkg/palettes/")


def run(cmd, **kwargs):
    print("+", " ".join(map(str, cmd)), flush=True)
    return subprocess.run(cmd, check=True, **kwargs)


def affected(paths, prefixes=(), exact=()):
    return any(p in exact or p.startswith(prefixes) for p in paths)


def docs_only(paths):
    return bool(paths) and all(p in ("README.md", "CONTRIBUTING.md") or p.startswith(DOCS) for p in paths)


def changed(base):
    run(["git", "rev-parse", "--verify", base], stdout=subprocess.DEVNULL)
    diff = run(["git", "diff", "--name-only", "-z", base, "--"], capture_output=True).stdout
    diff += run(["git", "ls-files", "--others", "--exclude-standard", "-z"], capture_output=True).stdout
    return sorted({p.decode("utf-8", "surrogateescape") for p in diff.split(b"\0") if p})


def format_go(paths):
    files = [p for p in paths if p.endswith(".go") and Path(p).is_file()]
    if files:
        result = run(["gofmt", "-l", *files], capture_output=True, text=True)
        if result.stdout.strip():
            raise RuntimeError("gofmt needed:\n" + result.stdout)


def go_smoke():
    with tempfile.TemporaryDirectory(prefix="markata-pr-") as tmp:
        binary = str(Path(tmp) / "markata-go")
        run(["go", "build", "-o", binary, "./cmd/markata-go"])
        run([binary, "version"])
        run([binary, "--help"], stdout=subprocess.DEVNULL)
        env = dict(os.environ)
        env["MARKATA_GO_TEST_BINARY"] = binary
        env["MARKATA_GO_TEST_BINARY_REVISION"] = run(
            ["git", "rev-parse", "HEAD"], capture_output=True, text=True).stdout.strip()
        test_args = ["-short", "-timeout=15m", "-skip",
                     "^TestBuildLab_LinkedAndFixtureMutationsCharacterizeProduct$", "./..."]
        if timing_dir := os.environ.get("MARKATA_CI_TIMINGS_DIR"):
            # Capture package timings in CI without an extra run of the suite.
            run([sys.executable, "scripts/ci/go_test_timing.py",
                 "--label", "Fast PR short Go suite",
                 "--output", str(Path(timing_dir) / "pr-go.json"),
                 "--", "go", "test", "-json", *test_args], env=env)
        else:
            run(["go", "test", *test_args], env=env)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--base", default="origin/main")
    mode = parser.add_mutually_exclusive_group()
    mode.add_argument("--only-format", action="store_true")
    mode.add_argument("--only-go", action="store_true")
    parser.add_argument("--browser", action="store_true")
    parser.add_argument("--print-plan", action="store_true")
    opt = parser.parse_args()
    if opt.only_go:
        if not opt.print_plan:
            go_smoke()
        return
    files = changed(opt.base)
    if not files:
        raise RuntimeError("no changed paths against " + opt.base)
    editor = affected(files, EDITORS, ("scripts/wasm-workbench.py", "examples/fastapi_editor.py"))
    config = affected(files, ("pkg/config/", "internal/configdocs/"))
    rendering = affected(files, RENDER, ("go.mod", "go.sum"))
    print(f"Changed: {len(files)}; docs_only={docs_only(files)} editor={editor} config={config} rendering={rendering}")
    if opt.print_plan:
        return
    if not opt.only_format:
        run(["git", "diff", "--check", opt.base, "--"])
    format_go(files)
    if opt.only_format:
        return
    if not docs_only(files):
        go_smoke()
        run(["just", "lint-new"])
    if config:
        run(["go", "test", "./internal/configdocs"])
    if rendering:
        run(["go", "run", "./scripts/rendering-contract", "--check"])
    if editor:
        run([sys.executable, "-m", "unittest", "discover", "-s", "tests/editors", "-p", "test_*.py", "-v"])
        if opt.browser:
            with tempfile.TemporaryDirectory(prefix="markata-editor-") as tmp:
                binary = str(Path(tmp) / "workbench")
                run(["go", "build", "-o", binary, "./cmd/markata-workbench"])
                run([sys.executable, "scripts/wasm-workbench.py"])
                run([sys.executable, "tests/editors/test_parity.py", "--go-binary", binary, "--wasm-output", ".wasm-workbench"])
        else:
            print("Editor change: also run just pr-check --browser before merge")
    print("Local checks passed; inspect full CI on the latest PR SHA.")


if __name__ == "__main__":
    try:
        main()
    except (RuntimeError, OSError, subprocess.CalledProcessError) as error:
        print(f"PR preflight failed: {error}", file=sys.stderr)
        sys.exit(1)
