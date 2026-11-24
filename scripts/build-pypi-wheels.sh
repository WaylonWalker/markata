#!/usr/bin/env bash
# Build official Markata (Go) Python wheels with Simon Willison's go-to-wheel.
# go-to-wheel builds "." from the module root, while Markata's CLI main lives
# under ./cmd/markata. Copy that exact main temporarily; never keep a fork.
set -euo pipefail
cd "$(dirname "$0")/.."

version="${1:?usage: scripts/build-pypi-wheels.sh <PEP440-0.x-version> [output-dir]}"
output="${2:-dist/pypi-wheels}"
if [[ ! "${version}" =~ ^0\.[0-9]+\.[0-9]+((a|b|rc|\.dev)[0-9]+)?$ ]]; then
  echo "Refusing invalid/non-0.x PyPI version: ${version}" >&2
  exit 2
fi
if ! command -v go-to-wheel >/dev/null 2>&1; then
  echo "Install Simon Willison's packaging tool: python -m pip install go-to-wheel==0.2" >&2
  exit 2
fi
entrypoint=markata_pypi_entrypoint.go
if [[ -e "${entrypoint}" ]]; then
  echo "Refusing to overwrite ${entrypoint}" >&2
  exit 2
fi
# An extra root-level Go package might conflict with this temporary main.
shopt -s nullglob
root_files=( ./*.go )
if (( ${#root_files[@]} > 0 )); then
  echo "Refusing root Go build with pre-existing Go sources: ${root_files[*]}" >&2
  exit 2
fi
trap 'rm -f -- "$entrypoint"' EXIT
cp cmd/markata/main.go "${entrypoint}"
sha="$(git rev-parse --short=12 HEAD)"
go-to-wheel . \
  --name markata \
  --entry-point markata \
  --version "${version}" \
  --output-dir "${output}" \
  --description "Markata: plugin-driven static site generator (Go implementation)" \
  --license MIT \
  --author "Waylon Walker" \
  --url "https://github.com/WaylonWalker/markata" \
  --readme README.md \
  --set-version-var "github.com/WaylonWalker/markata/cmd/markata-go/cmd.Version" \
  --ldflags "-X github.com/WaylonWalker/markata/cmd/markata-go/cmd.Commit=${sha}"

# Build failures on individual targets are warnings in go-to-wheel 0.2;
# scripts/check-pypi-wheels.py must validate the complete target matrix.

python scripts/add-python-renderer-to-wheels.py "${output}"
