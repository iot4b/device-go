#!/usr/bin/env bash
set -euo pipefail

# repo root (script is in apt/)
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
pkg_name="iot4bd"

# choose package dir (prefer iot4bd, fallback to iot4b)
if [ -d "$root/apt/$pkg_name" ]; then
  pkg_dir="$root/apt/$pkg_name"
else
  echo "No package directory found. Expected \`iot4bd\` in $root/apt" >&2
  exit 1
fi

# ensure required tools
command -v go >/dev/null 2>&1 || { echo "go not found in PATH" >&2; exit 1; }
command -v dpkg-deb >/dev/null 2>&1 || { echo "dpkg-deb not found in PATH" >&2; exit 1; }

# build go binary from root main.go

outbin="$root/apt/$pkg_name/usr/bin/$pkg_name"
echo "Building Go binary from $root/main.go -> $outbin"
cd "$root"
go build -o "$outbin" ./main.go

# determine output deb name (try to parse DEBIAN/control)
control="$pkg_dir/DEBIAN/control"
out_deb="$root/apt/$(basename "$pkg_dir").deb"
if [ -f "$control" ]; then
  pkg=$(awk -F': ' '/^Package:/ {print $2; exit}' "$control" | tr -d ' \t\r\n' || true)
  arch=$(awk -F': ' '/^Architecture:/ {print $2; exit}' "$control" | tr -d ' \t\r\n' || true)
  if [ -n "$pkg" ]; then
    if [ -n "$arch" ]; then
      out_deb="$root/apt/${pkg}_${arch}.deb"
    else
      out_deb="$root/apt/${pkg}.deb"
    fi
  fi
fi

# build deb
echo "Building deb package -> $out_deb"
dpkg-deb --build "$pkg_dir" "$out_deb"

# remove binary
rm -f "$pkg_dir/usr/bin/$pkg_name"

echo "Done: $out_deb"
