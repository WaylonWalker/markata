#!/usr/bin/env bash
# Offline smoke tests for new and old GitHub release archive layouts.
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# The GoReleaser brand must be canonical without dropping old archive names.
grep -Fxq 'project_name: markata' "${repo_root}/.goreleaser.yml"
grep -Fxq '      markata-go_' "${repo_root}/.goreleaser.yml"
grep -Fxq '      markata_' "${repo_root}/.goreleaser.yml"
# GitHub repo was renamed; installer must use the canonical release URL.
grep -Fq 'REPO="WaylonWalker/markata"' "${repo_root}/install.sh"
if grep -Fq 'REPO="WaylonWalker/markata-go"' "${repo_root}/install.sh"; then
    echo "Installer unexpectedly targets the legacy repository" >&2
    exit 1
fi
test_dir="$(mktemp -d)"
trap 'rm -rf "${test_dir}"' EXIT
mkdir -p "${test_dir}/bin" "${test_dir}/assets" "${test_dir}/install"

# Fail closed for any URL we have not explicitly made available.
cat > "${test_dir}/bin/curl" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
out=""
url=""
while [ "$#" -gt 0 ]; do
    case "$1" in
        -o) out="$2"; shift 2 ;;
        -*) shift ;;
        *) url="$1"; shift ;;
    esac
done
asset="${url##*/}"
if [ -z "${out}" ] || [ ! -f "${FAKE_RELEASE_DIR}/${asset}" ]; then
    exit 22
fi
cp "${FAKE_RELEASE_DIR}/${asset}" "${out}"
MOCK
chmod +x "${test_dir}/bin/curl"
export PATH="${test_dir}/bin:${PATH}"
export FAKE_RELEASE_DIR="${test_dir}/assets"
export MARKATA_INSTALL_DIR="${test_dir}/install"
export MARKATA_VERSION="v0.99.9"

os=linux
case "$(uname -s)" in
    Linux*) ;;
    *) echo "Installer smoke test only supports Linux" >&2; exit 1 ;;
esac
case "$(uname -m)" in
    x86_64|amd64) arch=x86_64 ;;
    aarch64|arm64) arch=arm64 ;;
    *) echo "Unsupported smoke-test platform" >&2; exit 1 ;;
esac

make_archive() {
    local name="$1"
    local work="${test_dir}/package"
    rm -rf "${work}"
    mkdir -p "${work}"
    cat > "${work}/${name}" <<'BINARY'
#!/bin/sh
if [ "${1:-}" = version ]; then
    echo 'markata 0.99.9'
    echo '  go:        go1.26.0'
else
    echo "Markata CLI"
fi
BINARY
    chmod +x "${work}/${name}"
    tar -czf "${FAKE_RELEASE_DIR}/${name}_0.99.9_${os}_${arch}.tar.gz" -C "${work}" "${name}"
}

make_archive markata
bash "${repo_root}/install.sh" >/dev/null
test -x "${MARKATA_INSTALL_DIR}/markata"
test -x "${MARKATA_INSTALL_DIR}/markata-go"
cmp "${MARKATA_INSTALL_DIR}/markata" "${MARKATA_INSTALL_DIR}/markata-go"

rm "${FAKE_RELEASE_DIR}/markata_0.99.9_${os}_${arch}.tar.gz"
make_archive markata-go
bash "${repo_root}/install.sh" >/dev/null
cmp "${MARKATA_INSTALL_DIR}/markata" "${MARKATA_INSTALL_DIR}/markata-go"

# A Python or unrelated executable must not be overwritten.
printf '#!/bin/sh\necho legacy-python\n' > "${MARKATA_INSTALL_DIR}/markata"
chmod +x "${MARKATA_INSTALL_DIR}/markata"
if bash "${repo_root}/install.sh" > "${test_dir}/collision.log" 2>&1; then
    echo "Installer unexpectedly overwrote conflicting command" >&2
    exit 1
fi
grep -q "Refusing to overwrite non-Go executable" "${test_dir}/collision.log"
grep -q legacy-python "${MARKATA_INSTALL_DIR}/markata"
echo "Canonical and legacy install paths: OK"
