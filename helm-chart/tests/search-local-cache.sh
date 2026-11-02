#!/usr/bin/env bash
set -euo pipefail

chart="$(dirname "$0")/.."
render() {
  helm template search-cache-test "$chart" --show-only templates/search.yaml "$@"
}

default_render="$(render)"
enabled_render="$(render --set search.localCache.enabled=true)"
readonly_render="$(render --set search.mode=read-only-index --set search.localCache.enabled=true)"

if grep -q 'name: search-cache' <<<"$default_render"; then
  echo "default render unexpectedly contains search-cache" >&2
  exit 1
fi
grep -q 'mountPath: /data/cache' <<<"$enabled_render"
grep -q 'name: init-search-cache' <<<"$enabled_render"
grep -q 'mkdir -p /data/cache/build /data/cache/plugin' <<<"$enabled_render"
test "$(grep -c -- '- name: search-cache' <<<"$enabled_render")" -eq 3
grep -A1 -- '- name: search-cache' <<<"$enabled_render" | grep -q 'emptyDir: {}'
if grep -q 'name: search-cache' <<<"$readonly_render"; then
  echo "read-only-index render unexpectedly contains search-cache" >&2
  exit 1
fi
grep -q 'mountPath: /data/source' <<<"$enabled_render"
grep -q 'mountPath: /data/search' <<<"$enabled_render"
