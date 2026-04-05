#!/usr/bin/env sh
set -eu

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
eval "$("${ROOT_DIR}/scripts/release-vars.sh")"

temp_dir="$(mktemp -d)"
trap 'rm -rf "${temp_dir}"' EXIT INT TERM

binary_path="${temp_dir}/iot4b"

"${ROOT_DIR}/scripts/build.sh" "${binary_path}" >/dev/null

exec "${binary_path}" "$@"
