#!/usr/bin/env bash
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
chart="$repo_root/helm-chart"
container_engine="${CONTAINER_ENGINE:-docker}"
nginx_image="${NGINX_IMAGE:-docker.io/library/nginx:1.27.5-alpine}"
tmp_dir="$(mktemp -d)"
site_root="$tmp_dir/site"
containers=()

cleanup() {
  local container
  for container in "${containers[@]}"; do
    "$container_engine" rm -f "$container" >/dev/null 2>&1 || true
  done
  rm -rf "$tmp_dir"
}
trap cleanup EXIT

fail() {
  printf 'FAIL: %s\n' "$*" >&2
  exit 1
}

render_config() {
  local output=$1
  shift
  helm template cors-test "$chart" --show-only templates/nginx-config.yaml "$@" |
    awk '
      /^  nginx.conf: \|$/ { in_config = 1; next }
      in_config && /^    / { sub(/^    /, ""); print; next }
      in_config && /^$/ { print; next }
      in_config { exit }
    ' >"$output"
}

header_count() {
  local config=$1
  local count
  count="$(grep -c 'add_header Access-Control-Allow-Origin' "$config" || true)"
  printf '%s' "$count"
}

assert_header() {
  local headers=$1
  local value=$2
  grep -Fqi "$value" "$headers" || fail "missing response header: $value"
}

assert_no_cors() {
  local headers=$1
  assert_no_header "$headers" 'Access-Control-Allow-Origin' 'protected preview response included Access-Control-Allow-Origin'
}

assert_no_header() {
  local headers=$1
  local header=$2
  local message=$3
  if grep -qi "^$header:" "$headers"; then
    fail "$message"
  fi
}

assert_no_release_redirect() {
  local headers=$1
  if grep -Eqi '^Location: .*/releases/' "$headers"; then
    fail "protected preview redirected into /releases/"
  fi
}

request() {
  local base_url=$1
  local path=$2
  local expected_status=$3
  local expected_body=${4-}
  local name=$5
  local status

  status="$(curl --max-redirs 0 --silent --show-error \
    --dump-header "$tmp_dir/$name.headers" \
    --output "$tmp_dir/$name.body" \
    --write-out '%{http_code}' \
    "$base_url$path")"
  [[ "$status" == "$expected_status" ]] ||
    fail "$path returned $status, expected $expected_status"
  if [[ -n "$expected_body" ]]; then
    local actual_body
    actual_body="$(<"$tmp_dir/$name.body")"
    [[ "$actual_body" == "$expected_body" ]] ||
      fail "$path returned the wrong body: $actual_body"
  fi
}

start_nginx() {
  local config=$1
  local name=$2
  local port_output
  local address

  containers+=("$name")
  "$container_engine" run --rm --detach \
    --name "$name" \
    --publish 127.0.0.1::80 \
    --volume "$config:/etc/nginx/nginx.conf:ro" \
    --volume "$site_root:/usr/share/nginx/html:ro" \
    "$nginx_image" >/dev/null

  port_output="$("$container_engine" port "$name" 80/tcp)"
  address="${port_output%%$'\n'*}"
  nginx_url="http://$address"
}

wait_for_nginx() {
  local base_url=$1
  local attempt
  for attempt in 1 2 3 4 5; do
    if curl --fail --silent --show-error "$base_url/index.html" >/dev/null; then
      return
    fi
    sleep 1
  done
  fail "nginx did not become ready"
}

mkdir -p "$site_root/current" "$site_root/releases/r1/foo" "$site_root/releases/r1/no-index"
printf 'public-index' >"$site_root/current/index.html"
printf 'public-404' >"$site_root/current/404.html"
printf 'public-asset' >"$site_root/current/app.js"
printf '{"kind":"content-index"}' >"$site_root/current/content-index.json"
printf '<a href="/docs">preview</a><script src="/app.js"></script>' >"$site_root/releases/r1/index.html"
printf 'preview-foo-index' >"$site_root/releases/r1/foo/index.html"
printf '<a href="/page-target">preview page</a>' >"$site_root/releases/r1/page.html"
printf 'preview-asset' >"$site_root/releases/r1/app.js"

render_config "$tmp_dir/default.conf"
render_config "$tmp_dir/enabled.conf" --set site.cors.enabled=true
render_config "$tmp_dir/disabled.conf" --set site.cors.enabled=false
render_config "$tmp_dir/no-content-index.conf" --set-string site.contentIndexPath=
render_config "$tmp_dir/runtime.conf" --set site.use404Page=true

if helm template cors-test "$chart" --show-only templates/nginx-config.yaml \
  --set-string site.contentIndexPath=relative.json >/dev/null 2>&1; then
  fail "invalid Content Index path rendered successfully"
fi

[[ "$(header_count "$tmp_dir/default.conf")" == 4 ]] || fail "default render has the wrong CORS header count"
[[ "$(header_count "$tmp_dir/enabled.conf")" == 4 ]] || fail "enabled render has the wrong CORS header count"
[[ "$(header_count "$tmp_dir/disabled.conf")" == 0 ]] || fail "disabled render includes CORS headers"
[[ "$(header_count "$tmp_dir/no-content-index.conf")" == 3 ]] || fail "empty Content Index path changed site-wide CORS"
[[ "$(header_count "$tmp_dir/runtime.conf")" == 5 ]] || fail "custom 404 render has the wrong CORS header count"
if grep -q 'location = /content-index.json' "$tmp_dir/no-content-index.conf"; then
  fail "empty Content Index path retained the exact location"
fi
grep -Fq 'try_files $uri/index.html @preview_file;' "$tmp_dir/runtime.conf" ||
  fail "preview lookup does not resolve directory indexes first"
grep -Fq 'location @preview_file' "$tmp_dir/runtime.conf" ||
  fail "preview file isolation is missing"
grep -Fq 'location @preview_not_found' "$tmp_dir/runtime.conf" ||
  fail "preview 404 isolation is missing"

"$container_engine" run --rm \
  --volume "$tmp_dir/runtime.conf:/etc/nginx/nginx.conf:ro" \
  "$nginx_image" nginx -t

runtime_container="markata-nginx-cors-$$-$RANDOM"
start_nginx "$tmp_dir/runtime.conf" "$runtime_container"
base_url="$nginx_url"
wait_for_nginx "$base_url"

request "$base_url" /index.html 200 public-index public-html
assert_header "$tmp_dir/public-html.headers" 'Access-Control-Allow-Origin: *'
assert_header "$tmp_dir/public-html.headers" 'Cache-Control: no-cache'
assert_header "$tmp_dir/public-html.headers" 'Cache-Control: no-cache, no-store, must-revalidate'

request "$base_url" /app.js 200 public-asset public-asset
assert_header "$tmp_dir/public-asset.headers" 'Access-Control-Allow-Origin: *'
assert_header "$tmp_dir/public-asset.headers" 'Cache-Control: max-age=31536000'
assert_header "$tmp_dir/public-asset.headers" 'Cache-Control: public, immutable'

request "$base_url" /content-index.json 200 '{"kind":"content-index"}' content-index
assert_header "$tmp_dir/content-index.headers" 'Access-Control-Allow-Origin: *'
rm "$site_root/current/content-index.json"
request "$base_url" /content-index.json 404 public-404 missing-content-index
assert_header "$tmp_dir/missing-content-index.headers" 'Access-Control-Allow-Origin: *'

preview_root_body='<a href="/__preview/r1/docs">preview</a><script src="/__preview/r1/app.js"></script>'
request "$base_url" /__preview/r1/index.html 200 "$preview_root_body" preview-index
request "$base_url" /__preview/r1/ 200 "$preview_root_body" preview-root
request "$base_url" /__preview/r1/foo 200 preview-foo-index preview-foo
request "$base_url" /__preview/r1/foo/ 200 preview-foo-index preview-foo-slash
request "$base_url" /__preview/r1/page 200 '<a href="/__preview/r1/page-target">preview page</a>' preview-page
request "$base_url" /__preview/r1/app.js 200 preview-asset preview-asset
request "$base_url" /__preview/r1/missing 404 '' preview-missing
request "$base_url" /__preview/r1/no-index 404 '' preview-no-index
request "$base_url" /__preview/r1/no-index/ 404 '' preview-no-index-slash
request "$base_url" /__preview/r1 302 '' preview-redirect

for name in preview-index preview-root preview-foo preview-foo-slash preview-page preview-asset preview-missing preview-no-index preview-no-index-slash preview-redirect; do
  assert_no_cors "$tmp_dir/$name.headers"
  assert_no_release_redirect "$tmp_dir/$name.headers"
done
for name in preview-index preview-root preview-foo preview-foo-slash preview-page preview-asset; do
  assert_no_header "$tmp_dir/$name.headers" 'Cache-Control' 'protected preview response used a public cache policy'
done
for name in preview-missing preview-no-index preview-no-index-slash; do
  [[ "$(<"$tmp_dir/$name.body")" != public-404 ]] ||
    fail "protected preview miss used the public 404 page"
done
grep -Eqi '^Location: (https?://[^/]+)?/__preview/r1/' "$tmp_dir/preview-redirect.headers" ||
  fail "bare preview redirect left the protected preview namespace"

"$container_engine" rm -f "$runtime_container" >/dev/null
containers=()
printf '{"kind":"content-index"}' >"$site_root/current/content-index.json"

edge_config="$tmp_dir/edge-cache.conf"
render_config "$edge_config" --set site.cacheHtmlAtEdge=true
edge_container="markata-nginx-edge-cache-$$-$RANDOM"
"$container_engine" run --rm \
  --volume "$edge_config:/etc/nginx/nginx.conf:ro" \
  "$nginx_image" nginx -t
start_nginx "$edge_config" "$edge_container"
edge_url="$nginx_url"
wait_for_nginx "$edge_url"
request "$edge_url" /index.html 200 public-index edge-cache-html
assert_header "$tmp_dir/edge-cache-html.headers" 'Cache-Control: public, max-age=0, s-maxage=60'
assert_no_header "$tmp_dir/edge-cache-html.headers" 'Cache-Control: no-store' 'edge-cache HTML response remained uncacheable'
request "$edge_url" /__preview/r1/index.html 200 "$preview_root_body" edge-cache-preview
assert_no_header "$tmp_dir/edge-cache-preview.headers" 'Cache-Control' 'protected preview received the public edge-cache policy'

"$container_engine" rm -f "$edge_container" >/dev/null
containers=()

disabled_container="markata-nginx-no-cors-$$-$RANDOM"
start_nginx "$tmp_dir/disabled.conf" "$disabled_container"
disabled_url="$nginx_url"
wait_for_nginx "$disabled_url"
request "$disabled_url" /index.html 200 public-index disabled-html
request "$disabled_url" /app.js 200 public-asset disabled-asset
request "$disabled_url" /content-index.json 200 '{"kind":"content-index"}' disabled-content-index
request "$disabled_url" /__preview/r1/missing 404 '' disabled-preview-missing
request "$disabled_url" /__preview/r1/no-index/ 404 '' disabled-preview-no-index
assert_no_cors "$tmp_dir/disabled-html.headers"
assert_no_cors "$tmp_dir/disabled-asset.headers"
assert_no_cors "$tmp_dir/disabled-content-index.headers"
assert_no_cors "$tmp_dir/disabled-preview-missing.headers"
assert_no_cors "$tmp_dir/disabled-preview-no-index.headers"

bash "$repo_root/helm-chart/tests/builder-admin-workspace.sh"
printf 'Helm nginx CORS and Builder Admin workspace tests passed.\n'
