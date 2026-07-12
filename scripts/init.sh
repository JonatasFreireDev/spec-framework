#!/usr/bin/env sh
set -eu

VERSION="${SPEC_FRAMEWORK_VERSION:-}"
TARGET="${1:-}"
if [ -z "$VERSION" ]; then
  VERSION="$(curl -fsSL https://api.github.com/repos/JonatasFreireDev/spec-framework/releases/latest | sed -n 's/.*"tag_name": *"v\([^"]*\)".*/\1/p')"
fi

case "$(uname -s)" in
  Linux) os=linux ;;
  Darwin) os=darwin ;;
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
(cd "$tmp" && grep "  $archive\$" checksums.txt | sha256sum -c -)
tar -xzf "$tmp/$archive" -C "$tmp"
mkdir -p "$HOME/.local/bin"
install "$tmp/spec-framework" "$HOME/.local/bin/spec-framework"
PATH="$HOME/.local/bin:$PATH"
export PATH
if [ -n "$TARGET" ]; then spec-framework init "$TARGET"; else spec-framework init; fi
printf 'Spec Framework installed at %s\n' "$HOME/.local/bin/spec-framework"
