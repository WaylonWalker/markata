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
  --set-string builderAdmin.editor.gitAuth.repository=https://github.com/example/site
  --set-string builderAdmin.editor.gitAuth.branch=main
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
publisher_container="$(sed -n '/^[[:space:]]*- name: "editor-test-notes-editor-git-publisher"$/,/readinessProbe:/p' "$tmp_dir/all.yaml")"
editor_container="$(sed -n '/^[[:space:]]*- name: "editor-test-notes-editor"$/,/readinessProbe:/p' "$tmp_dir/all.yaml")"
key_initializer="$(sed -n '/^[[:space:]]*- name: prepare-editor-git-key$/,/^[[:space:]]*- name: editor-git-publish-key$/p' "$tmp_dir/all.yaml")"
grep -q 'editor-git-publish-key' <<<"$publisher_container" || fail "publisher does not mount the private-key copy"
grep -q 'editor-ssh-known-hosts' <<<"$publisher_container" || fail "publisher does not mount verified known_hosts"
grep -q -- 'https://github.com/example/site' <<<"$publisher_container" || fail "publisher origin allowlist is missing"
grep -q 'refs/heads/main' <<<"$publisher_container" || fail "publisher branch allowlist is missing"
grep -q 'editor-git-publisher-socket' <<<"$editor_container" || fail "editor does not mount the publisher socket"
if grep -Eq 'name: editor-git-publish(-key)?$|editor-ssh-known-hosts|GIT_SSH_COMMAND' <<<"$editor_container"; then
  fail "editor container has direct Git credential access"
fi
grep -q 'name: editor-git-publish$' <<<"$key_initializer" || fail "init container does not mount the source Secret"
if grep -q 'ports:' <<<"$publisher_container"; then
  fail "Git publisher must not expose a network port"
fi
grep -q 'name: prepare-editor-git-key' "$tmp_dir/all.yaml" || fail "private key copy init container is missing"
grep -q 'install -m 0400 /run/markata-git-secret/id_ed25519 /run/markata-git/id_ed25519 && chown 65534:1000 /run/markata-git/id_ed25519' "$tmp_dir/all.yaml" || fail "private key copy must set OpenSSH-safe permissions and publisher ownership"
grep -q 'runAsUser: 65534' <<<"$publisher_container" || fail "publisher must use the passwd-resolvable nobody user"
grep -q 'runAsGroup: 1000' <<<"$publisher_container" || fail "publisher must share the editor group for the local socket"
grep -A1 'name: HOME' <<<"$publisher_container" | grep -q 'value: /tmp' || fail "publisher HOME must use its writable temporary volume"
grep -A8 'startupProbe:' <<<"$editor_container" | grep -q 'failureThreshold: 120' || fail "editor startup probe must allow ten minutes for workspace discovery"
grep -q 'secretName: "editor-key"' "$tmp_dir/all.yaml" || fail "publisher Git deploy key secret is not mounted"
grep -A8 'name: editor-git-publish$' "$tmp_dir/all.yaml" | grep -Eq 'defaultMode: (256|0400)' || fail "source deploy key is not owner-readable only"
if grep -q 'GIT_SSH_COMMAND' "$tmp_dir/all.yaml"; then
  fail "SSH private-key command is exposed in the pod environment"
fi
grep -q 'github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl' "$tmp_dir/all.yaml" || fail "editor known_hosts entry is missing"
grep -q 'X-Hlab-User-Id: ""' "$tmp_dir/all.yaml" || fail "editor middleware does not strip spoofed user ID"
grep -q 'X-Hlab-Scopes: ""' "$tmp_dir/all.yaml" || fail "editor middleware does not strip spoofed scopes"
grep -q 'addAuthCookiesToResponse:' "$tmp_dir/all.yaml" || fail "editor middleware does not forward the Hlab app session cookie"
grep -q -- '- "__Host-hlab-app-session"' "$tmp_dir/all.yaml" || fail "editor middleware omits the browser session handoff cookie"
grep -q 'router.middlewares:' "$tmp_dir/all.yaml" || fail "editor ingress does not attach ForwardAuth"
grep -q 'port: 8766' "$tmp_dir/all.yaml" || fail "editor network policy does not allow Traefik to the editor port"

helm template editor-test "$chart" "${common[@]}" \
  --set builderAdmin.editor.dropperUploader.enabled=true \
  --set-string builderAdmin.editor.dropperUploader.existingSecretName=editor-dropper \
  --set-string builderAdmin.editor.dropperUploader.tokenKey=token \
  >"$tmp_dir/with-dropper.yaml"
dropper_uploader="$(sed -n '/^[[:space:]]*- name: "editor-test-notes-editor-dropper-uploader"$/,/readinessProbe:/p' "$tmp_dir/with-dropper.yaml")"
dropper_editor="$(sed -n '/^[[:space:]]*- name: "editor-test-notes-editor"$/,/readinessProbe:/p' "$tmp_dir/with-dropper.yaml")"
grep -q 'name: DROPPER_TOKEN' <<<"$dropper_uploader" || fail "Dropper uploader does not receive the token"
grep -q 'secretKeyRef:' <<<"$dropper_uploader" || fail "Dropper uploader token is not read from a Secret"
grep -q 'editor-dropper' <<<"$dropper_uploader" || fail "Dropper uploader Secret reference is missing"
grep -q -- 'https://dropper.waylonwalker.com/api/upload?format=json' <<<"$dropper_uploader" || fail "Dropper endpoint is not fixed in uploader config"
grep -q 'editor-dropper-uploader-socket' <<<"$dropper_uploader" || fail "Dropper uploader socket volume is missing"
if grep -q 'ports:' <<<"$dropper_uploader"; then
  fail "Dropper uploader must not expose a network port"
fi
if grep -Eq 'name: DROPPER_TOKEN|secretKeyRef:' <<<"$dropper_editor"; then
  fail "editor container has direct Dropper credential access"
fi
if grep -Eq 'mountPath: /data/source|editor-git-publish|editor-ssh-known-hosts' <<<"$dropper_uploader"; then
  fail "Dropper uploader has source or Git credential access"
fi
grep -q -- '--dropper-uploader-socket' <<<"$dropper_editor" || fail "editor does not use the credential-isolated uploader"
grep -q 'runAsUser: 65534' <<<"$dropper_uploader" || fail "Dropper uploader must run as nobody"
grep -q 'readOnlyRootFilesystem: true' <<<"$dropper_uploader" || fail "Dropper uploader root filesystem must be read-only"

if helm template editor-test "$chart" "${common[@]}" \
  --set builderAdmin.editor.dropperUploader.enabled=true \
  >/dev/null 2>&1; then
  fail "Dropper uploader rendered without a Secret reference"
fi

if helm template editor-test "$chart" "${common[@]}" \
  --set-string 'builderAdmin.editor.extraEnv[0].name=DROPPER_TOKEN' \
  --set-string 'builderAdmin.editor.extraEnv[0].value=must-not-leak' \
  >/dev/null 2>&1; then
  fail "editor accepted a direct Dropper token environment variable"
fi

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

if helm template editor-test "$chart" "${common[@]}" \
  --set-string builderAdmin.editor.gitAuth.repository= \
  --set-string builderAdmin.sourceGit.repository= \
  >/dev/null 2>&1; then
  fail "hosted publisher rendered without an explicit GitHub repository allowlist"
fi

helm template editor-test "$chart" "${common[@]}" \
  --set builderAdmin.sourceGit.enabled=true \
  --set-string builderAdmin.sourceGit.repository=https://github.com/example/source \
  --set-string builderAdmin.editor.gitAuth.repository= \
  --set-string builderAdmin.editor.gitAuth.branch= \
  >"$tmp_dir/source-git-defaults.yaml"
grep -q -- 'https://github.com/example/source' "$tmp_dir/source-git-defaults.yaml" || fail "publisher did not default to sourceGit.repository"
grep -q 'refs/heads/main' "$tmp_dir/source-git-defaults.yaml" || fail "publisher did not default to sourceGit.branch"

printf 'Builder Admin hosted editor Helm tests passed.\n'
