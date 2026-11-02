#!/usr/bin/env bash
set -euo pipefail

chart="$(dirname "$0")/.."
default_render="$(helm template site-volume-test "$chart" --show-only templates/pvcs.yaml)"
bound_render="$(helm template site-volume-test "$chart" --show-only templates/pvcs.yaml --set storage.site.volumeName=pv-fast-site)"
existing_claim_pvcs="$(helm template site-volume-test "$chart" --show-only templates/pvcs.yaml --set-string storage.site.existingClaimName=site-volume-test-site)"
existing_claim_site="$(helm template site-volume-test "$chart" --show-only templates/site.yaml --set-string namespaceOverride=existing-site --set-string storage.site.existingClaimName=site-volume-test-site)"

if grep -q 'volumeName:' <<<"$default_render"; then
  echo "Default PVC unexpectedly pins a volume" >&2
  exit 1
fi
grep -q 'volumeName: "pv-fast-site"' <<<"$bound_render"
if grep -q 'name: "site-volume-test-notes-site-pvc"' <<<"$existing_claim_pvcs"; then
  echo "Existing-claim mode unexpectedly rendered a new site PVC" >&2
  exit 1
fi
grep -q 'claimName: "site-volume-test-site"' <<<"$existing_claim_site"
grep -q 'namespace: "existing-site"' <<<"$existing_claim_site"
