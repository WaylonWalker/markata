#!/usr/bin/env bash
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
chart="$repo_root/helm-chart"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

fail() {
  printf 'FAIL: %s\n' "$*" >&2
  exit 1
}

common=(
  --set-string project_identifier=editor-test
  --set builderAdmin.enabled=true
  --set builderAdmin.ingress.enabled=true
  --set builderAdmin.ingress.auth.enabled=true
  --set-string builderAdmin.ingress.auth.internalUrl=http://hlab-auth.default.svc.cluster.local:8000
  --set-string builderAdmin.ingress.host=build.example.com
  --set-string builderAdmin.ingress.ingressClassName=traefik
  --set builderAdmin.ingress.tls.enabled=true
  --set-string builderAdmin.ingress.tls.secretName=build-tls
  --set builderAdmin.networkPolicy.enabled=true
  --set-string 'builderAdmin.auth.trustedProxyCIDRs[0]=10.42.0.0/24'
  --set builderAdmin.editor.enabled=true
  --set-string builderAdmin.editor.requiredScope=markata-editor:admin
  --set-string builderAdmin.editor.gitAuth.existingSecretName=editor-key
  --set-string 'builderAdmin.editor.gitAuth.knownHosts=github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl'
  --set-string builderAdmin.editor.ingress.host=edit.example.com
  --set-string builderAdmin.editor.ingress.ingressClassName=traefik
  --set-string builderAdmin.editor.ingress.tls.secretName=editor-tls
)

helm template editor-test "$chart" "${common[@]}" >"$tmp_dir/all.yaml"
grep -A20 'name: "editor-test-notes-editor"' "$tmp_dir/all.yaml" | grep -q 'port: 8766' || fail "editor service port missing"
grep -A20 'name: "editor-test-notes-editor"' "$tmp_dir/all.yaml" | grep -q 'targetPort: editor-http' || fail "editor service target port missing"
grep -q -- '--public-origin' "$tmp_dir/all.yaml" || fail "editor public origin missing"
grep -q -- 'https://edit.example.com' "$tmp_dir/all.yaml" || fail "editor origin does not match ingress host"
grep -q -- '--required-scope' "$tmp_dir/all.yaml" || fail "editor scope argument missing"
grep -q -- 'markata-editor:admin' "$tmp_dir/all.yaml" || fail "editor writer scope missing"
grep -q 'secretName: "editor-key"' "$tmp_dir/all.yaml" || fail "editor Git deploy key secret is not mounted"
grep -A8 'name: editor-git-publish' "$tmp_dir/all.yaml" | grep -Eq 'defaultMode: (292|0444)' || fail "editor Git deploy key is not readable by the non-root editor container"
grep -q 'StrictHostKeyChecking=yes' "$tmp_dir/all.yaml" || fail "editor SSH does not verify the GitHub host key"
grep -q 'github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl' "$tmp_dir/all.yaml" || fail "editor known_hosts entry is missing"
grep -q 'X-Hlab-User-Id: ""' "$tmp_dir/all.yaml" || fail "editor middleware does not strip spoofed user ID"
grep -q 'X-Hlab-Scopes: ""' "$tmp_dir/all.yaml" || fail "editor middleware does not strip spoofed scopes"
grep -q 'addAuthCookiesToResponse:' "$tmp_dir/all.yaml" || fail "editor middleware does not forward the Hlab app session cookie"
grep -q -- '- "__Host-hlab-app-session"' "$tmp_dir/all.yaml" || fail "editor middleware omits the browser session handoff cookie"
grep -q 'router.middlewares:' "$tmp_dir/all.yaml" || fail "editor ingress does not attach ForwardAuth"
grep -q 'port: 8766' "$tmp_dir/all.yaml" || fail "editor network policy does not allow Traefik to the editor port"

if helm template editor-test "$chart" "${common[@]}" \
  --set-string builderAdmin.editor.requiredScope=markata-editor:read \
  >/dev/null 2>&1; then
  fail "viewer scope was accepted for hosted editor writes"
fi

if helm template editor-test "$chart" "${common[@]}" \
  --set-string builderAdmin.editor.ingress.host=build.example.com \
  >/dev/null 2>&1; then
  fail "editor ingress accepted the Builder Admin host"
fi

if helm template editor-test "$chart" \
  --set builderAdmin.editor.enabled=true \
  >/dev/null 2>&1; then
  fail "hosted editor rendered without protected Builder Admin infrastructure"
fi

printf 'Builder Admin hosted editor Helm tests passed.\n'
