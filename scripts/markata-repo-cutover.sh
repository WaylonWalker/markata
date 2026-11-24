#!/usr/bin/env bash
# Safe owner-driven GitHub repo swap for Python Markata and Go Markata.
# Reads via gh api by default. Only `apply --execute` renames repositories.
set -euo pipefail

OWNER="WaylonWalker"
PYTHON_REPO_ID="333449359"
GO_REPO_ID="1140034778"

fail() { printf 'ERROR: %s\n' "$*" >&2; exit 1; }
info() { printf '%s\n' "$*"; }

usage() {
  cat <<'EOF'
Usage:
  bash scripts/markata-repo-cutover.sh preflight
  bash scripts/markata-repo-cutover.sh apply --execute
  bash scripts/markata-repo-cutover.sh postcheck

Preflight and postcheck are read-only. Apply requires interactive confirmation
and renames the historical Python repository first, then the Go repository.
Requires an authenticated GitHub CLI with admin rights to both repositories.
EOF
}

require_gh() {
  command -v gh >/dev/null 2>&1 || fail 'Install the GitHub CLI (gh) first.'
  gh auth status >/dev/null 2>&1 || fail 'Run gh auth login with an account that administers both repositories.'
}

repo_id() {
  gh api "repos/${OWNER}/$1" --jq '.id'
}

assert_repo() {
  local name="$1" want="$2" actual
  actual="$(repo_id "$name")" || fail "Cannot read ${OWNER}/${name}; check GitHub credentials and repo visibility."
  [[ "$actual" == "$want" ]] || fail "${OWNER}/${name} has ID ${actual}, expected ${want}. Stopping to protect the wrong repository."
  info "Verified ${OWNER}/${name} (ID ${actual})"
}

assert_unoccupied() {
  local name="$1" err_file
  err_file="$(mktemp)"
  if gh api "repos/${OWNER}/${name}" --jq '.id' >/dev/null 2>"$err_file"; then
    rm -f "$err_file"
    fail "${OWNER}/${name} already exists; do not overwrite or rename it."
  fi
  if ! grep -q 'HTTP 404' "$err_file"; then
    cat "$err_file" >&2
    rm -f "$err_file"
    fail "Cannot prove ${OWNER}/${name} is unoccupied (not an HTTP 404)."
  fi
  rm -f "$err_file"
  info "Verified ${OWNER}/${name} is unoccupied"
}

assert_merged() {
  local source="$1" number="$2" state
  state="$(gh api "repos/${OWNER}/${source}/pulls/${number}" --jq '.merged')" ||
    fail "Could not verify PR #${number} in ${source}"
  [[ "$state" == "true" ]] || fail "PR #${number} in ${source} is not merged."
  info "Verified ${source} PR #${number} merged"
}

preflight() {
  require_gh
  assert_repo markata "$PYTHON_REPO_ID"
  assert_repo markata-go "$GO_REPO_ID"
  assert_unoccupied markata-py
  for number in 175 176; do assert_merged markata "$number"; done
  for number in 1629 1671 1672 1673; do assert_merged markata-go "$number"; done
  info
  info "Preflight passed. THIS DOES NOT rename either repo."
  info "Before apply: confirm Python publishing, Cloudflare Pages and external app"
  info "connections; the existing Python GitHub URL will be claimed by Go."
}

postcheck() {
  require_gh
  assert_repo markata-py "$PYTHON_REPO_ID"
  assert_repo markata "$GO_REPO_ID"
  info
  info "Repository names are correct. Next: migrate Go imports/module paths,"
  info "installer links, Argo source URLs, and Python project links. Keep old"
  info "container aliases and config compatibility through the 1.0 RC."
  info "Do not publish 1.0 or change production sources until validated."
}

apply() {
  [[ "${1:-}" == "--execute" && "$#" -eq 1 ]] ||
    fail "Renaming requires: bash scripts/markata-repo-cutover.sh apply --execute"
  [[ -t 0 ]] || fail 'Apply must be run in an interactive terminal.'
  preflight
  info
  info 'This will rename TWO PUBLIC repositories and immediately change the'
  info 'meaning of github.com/WaylonWalker/markata (Python -> Go).'
  info 'If the second step fails, Python will already have moved to markata-py.'
  local confirmation
  read -r -p "Type 'RENAME MARKATA' to continue: " confirmation
  [[ "$confirmation" == "RENAME MARKATA" ]] || fail "Cancelled; no repositories renamed."

  info "Renaming historical Python markata -> markata-py..."
  gh api -X PATCH "repos/${OWNER}/markata" -f 'name=markata-py' --jq '.full_name'
  assert_repo markata-py "$PYTHON_REPO_ID"

  info "Renaming Go markata-go -> markata..."
  gh api -X PATCH "repos/${OWNER}/markata-go" -f 'name=markata' --jq '.full_name'
  postcheck
}

case "${1:-}" in
  preflight) [[ "$#" -eq 1 ]] || { usage; exit 2; }; preflight ;;
  postcheck) [[ "$#" -eq 1 ]] || { usage; exit 2; }; postcheck ;;
  apply) shift; apply "$@" ;;
  *) usage; exit 2 ;;
esac
