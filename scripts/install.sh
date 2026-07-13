#!/usr/bin/env sh
set -eu

VERSION="${SPEC_FRAMEWORK_VERSION:-}"
if [ -z "$VERSION" ]; then
  VERSION="$(curl -fsSL https://api.github.com/repos/JonatasFreireDev/spec-framework/releases/latest | sed -n 's/.*"tag_name": *"v\([^"]*\)".*/\1/p')"
fi

case "$(uname -s)" in
  Linux) os=linux ; checksum='sha256sum -c' ;;
  Darwin) os=darwin ; checksum='shasum -a 256 -c' ;;
  *) echo "unsupported operating system" >&2; exit 1 ;;
esac
case "$(uname -m)" in
  arm64|aarch64) arch=arm64 ;;
  x86_64|amd64) arch=amd64 ;;
  *) echo "unsupported architecture" >&2; exit 1 ;;
esac

base="https://github.com/JonatasFreireDev/spec-framework/releases/download/v$VERSION"
archive="spec-framework_${VERSION}_${os}_${arch}.tar.gz"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
curl -fsSLo "$tmp/$archive" "$base/$archive"
curl -fsSLo "$tmp/checksums.txt" "$base/checksums.txt"
(cd "$tmp" && grep "  $archive\$" checksums.txt | $checksum)
tar -xzf "$tmp/$archive" -C "$tmp"
mkdir -p "$HOME/.local/bin"
install "$tmp/spec-framework" "$HOME/.local/bin/spec-framework"
PATH="$HOME/.local/bin:$PATH"
export PATH
spec-framework version
installed_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
printf '{\n  "schema_version": 1,\n  "managed_by": "spec-framework-installer",\n  "version": "%s",\n  "executable": "%s",\n  "path_entry": "%s",\n  "installed_at": "%s"\n}\n' \
  "$VERSION" "$HOME/.local/bin/spec-framework" "$HOME/.local/bin" "$installed_at" > "$HOME/.local/bin/install.json"
printf "Spec Framework installed at %s. Run 'spec-framework init' explicitly when you want to initialize a product.\n" "$HOME/.local/bin/spec-framework"
