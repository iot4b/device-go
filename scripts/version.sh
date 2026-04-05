#!/usr/bin/env sh
set -eu

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
eval "$("${ROOT_DIR}/scripts/release-vars.sh")"

printf 'Version: %s\n' "${IOT4B_VERSION}"
printf 'GitTag: %s\n' "${IOT4B_GIT_TAG}"
printf 'Commit: %s\n' "${IOT4B_COMMIT}"
printf 'BuildDate: %s\n' "${IOT4B_BUILD_DATE}"
