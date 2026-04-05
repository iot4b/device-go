#!/usr/bin/env sh
set -eu

raw_version="${IOT4B_VERSION:-${VERSION:-}}"
if [ -z "${raw_version}" ]; then
  if git describe --tags --exact-match >/dev/null 2>&1; then
    raw_version="$(git describe --tags --exact-match)"
  elif git describe --tags --abbrev=0 >/dev/null 2>&1; then
    raw_version="$(git describe --tags --abbrev=0)-dev"
  else
    raw_version="v2.0.0-dev"
  fi
fi

version="${raw_version#v}"
git_tag="${raw_version}"
case "${git_tag}" in
  v*) ;;
  *) git_tag="v${version}" ;;
esac

commit="${IOT4B_COMMIT:-}"
if [ -z "${commit}" ]; then
  commit="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"
fi

build_date="${IOT4B_BUILD_DATE:-}"
if [ -z "${build_date}" ]; then
  build_date="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
fi

printf 'IOT4B_VERSION=%s\n' "${version}"
printf 'IOT4B_GIT_TAG=%s\n' "${git_tag}"
printf 'IOT4B_COMMIT=%s\n' "${commit}"
printf 'IOT4B_BUILD_DATE=%s\n' "${build_date}"
