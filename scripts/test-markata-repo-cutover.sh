#!/usr/bin/env bash
# Offline/read-only regression checks for the guarded owner rename helper.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf -- "$TMP"' EXIT

cat > "$TMP/gh" <<'MOCK_GH'
#!/usr/bin/env bash
set -euo pipefail
[[ "$1" == "auth" && "${2:-}" == "status" ]] && exit 0
[[ "$1" == "api" ]] || { echo "unexpected gh command" >&2; exit 1; }
# The preflight and postcheck modes must never invoke a write operation.
[[ "${2:-}" != "-X" ]] || { echo "unexpected write operation" >&2; exit 1; }

route="${2:-}"
scenario="${MOCK_SCENARIO:-pre}"
case "$route" in
  repos/WaylonWalker/markata)
    if [[ "$scenario" == "post" ]]; then
      echo 1140034778
    elif [[ "${MOCK_WRONG_ID:-}" == 1 ]]; then
      echo 9999
    else
      echo 333449359
    fi
    ;;
  repos/WaylonWalker/markata-go)
    echo 1140034778
    ;;
  repos/WaylonWalker/markata-py)
    if [[ "$scenario" == "post" ]]; then
      echo 333449359
    elif [[ "${MOCK_NON_404:-}" == 1 ]]; then
      echo 'gh: permission denied (HTTP 403)' >&2
      exit 1
    else
      echo 'gh: Not Found (HTTP 404)' >&2
      exit 1
    fi
    ;;
  repos/WaylonWalker/markata/pulls/175|repos/WaylonWalker/markata/pulls/176|\
  repos/WaylonWalker/markata-go/pulls/1629|repos/WaylonWalker/markata-go/pulls/1671|\
  repos/WaylonWalker/markata-go/pulls/1672|repos/WaylonWalker/markata-go/pulls/1673)
    if [[ "${MOCK_OPEN_PR:-}" == 1 && "$route" == */1673 ]]; then
      echo false
    else
      echo true
    fi
    ;;
  *)
    echo "Unexpected route: $route" >&2
    exit 1
    ;;
esac
MOCK_GH
chmod +x "$TMP/gh"

run_success() {
  local label="$1"
  shift
  if env PATH="$TMP:$PATH" "$@" >/dev/null; then
    printf 'PASS: %s\n' "$label"
  else
    printf 'FAIL: %s should succeed\n' "$label" >&2
    exit 1
  fi
}

run_failure() {
  local label="$1"
  shift
  if env PATH="$TMP:$PATH" "$@" >/dev/null 2>&1; then
    printf 'FAIL: %s should have been rejected\n' "$label" >&2
    exit 1
  else
    printf 'PASS: %s rejected\n' "$label"
  fi
}

run_success 'read-only preflight' env MOCK_SCENARIO=pre bash "$ROOT/scripts/markata-repo-cutover.sh" preflight
run_success 'read-only postcheck' env MOCK_SCENARIO=post bash "$ROOT/scripts/markata-repo-cutover.sh" postcheck
run_failure 'wrong Python repository ID' env MOCK_WRONG_ID=1 bash "$ROOT/scripts/markata-repo-cutover.sh" preflight
run_failure 'missing merged prerequisite PR' env MOCK_OPEN_PR=1 bash "$ROOT/scripts/markata-repo-cutover.sh" preflight
run_failure '403 cannot be treated as unoccupied' env MOCK_NON_404=1 bash "$ROOT/scripts/markata-repo-cutover.sh" preflight
run_failure 'apply without --execute' bash "$ROOT/scripts/markata-repo-cutover.sh" apply
run_failure 'apply without terminal' bash "$ROOT/scripts/markata-repo-cutover.sh" apply --execute
