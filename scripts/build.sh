#!/usr/bin/env sh
set -eu

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
eval "$("${ROOT_DIR}/scripts/release-vars.sh")"

output_path="${1:-${ROOT_DIR}/iot4b}"

ld_flags="-s -w -X device-go/packages/buildinfo.Version=${IOT4B_VERSION} -X device-go/packages/buildinfo.Commit=${IOT4B_COMMIT} -X device-go/packages/buildinfo.BuildDate=${IOT4B_BUILD_DATE}"

mkdir -p "$(dirname "${output_path}")"

cd "${ROOT_DIR}"
go build -ldflags="${ld_flags}" -o "${output_path}" ./main.go

echo "Built ${output_path}"
echo "Version: ${IOT4B_VERSION}"
echo "Commit: ${IOT4B_COMMIT}"
echo "BuildDate: ${IOT4B_BUILD_DATE}"
