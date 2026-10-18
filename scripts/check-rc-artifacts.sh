#!/usr/bin/env bash
# Validate snapshot release assets without publishing any GitHub release or tag.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST="${1:-${ROOT}/dist}"

module="$(cd "$ROOT" && go list -m)"
[[ "$module" == "github.com/WaylonWalker/markata" ]] || {
  printf 'Expected canonical Go module, got %s\n' "$module" >&2
  exit 1
}
[[ -f "$DIST/checksums.txt" ]] || {
  printf 'No checksum manifest in %s\n' "$DIST" >&2
  exit 1
}

for os in linux darwin windows; do
  ext=tar.gz
  if [[ "$os" == windows ]]; then ext=zip; fi

  for executable in markata markata-go; do
    archive="$(find "$DIST" -maxdepth 1 -type f -name "${executable}_*_${os}_x86_64.${ext}" -print -quit)"
    [[ -n "$archive" ]] || {
      printf 'Missing %s release archive for %s/amd64\n' "$executable" "$os" >&2
      exit 1
    }
    printf 'Found %s\n' "$(basename "$archive")"
    expected="$executable"
    if [[ "$os" == windows ]]; then
      expected="${executable}.exe"
      unzip -Z1 "$archive" | grep -E "(^|/)${expected}$" >/dev/null || {
        printf 'Missing executable %s in %s\n' "$expected" "$archive" >&2
        exit 1
      }
    else
      tar -tzf "$archive" | grep -E "(^|/)${expected}$" >/dev/null || {
        printf 'Missing executable %s in %s\n' "$expected" "$archive" >&2
        exit 1
      }
    fi
  done
done

(cd "$DIST" && sha256sum -c checksums.txt)
printf 'Snapshot release packaging validated for both identities on linux/darwin/windows amd64\n'
