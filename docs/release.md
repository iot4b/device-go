# Release Process

This document describes how to publish a new `iot4b` release and trigger the automated deployment flow.

## Prerequisites

Before creating a release tag, make sure the repository is configured with the required GitHub Actions settings:

- Repository variable `HOMEBREW_TAP_REPO`
  - expected value: `iot4b/homebrew-tap`
- Repository variable `HOMEBREW_TAP_BRANCH`
  - expected value: `main`
- Repository secret `HOMEBREW_TAP_TOKEN`
  - must have permission to push to the Homebrew tap repository

The release workflow also expects:

- a working `self-hosted` GitHub Actions runner
- the APT repository available at `http://repo.iot4b.co/apt`
- the OPKG repository available at `http://repo.iot4b.co/opkg`
- the public APT key available at `http://repo.iot4b.co/apt/iot4b.asc`
- commit and PR titles written in a consistent style, preferably Conventional Commits

## How To Create A Release

1. Make sure the branch you are releasing from is ready.
2. Create a SemVer tag with the `v` prefix.

Example:

```sh
git tag v2.0.0
git push origin v2.0.0
```

Pushing the tag triggers the GitHub Actions workflow in `.github/workflows/release.yml`.

## What The Release Workflow Does

The release workflow runs on the `self-hosted` runner and has 2 jobs:

### `publish`

This job:

- checks out the repository
- prepares `IOT4B_VERSION`, `IOT4B_COMMIT`, and `IOT4B_BUILD_DATE`
- runs `go test ./...`
- builds `.deb` and `.ipk` packages
- publishes APT and OPKG repositories on the server
- generates checksums
- creates a GitHub Release for the tag with auto-generated release notes
- generates the Homebrew formula
- updates the Homebrew tap repository

### `deploy`

This job runs only after `publish` succeeds.

It:

- runs the Ansible playbook in `ansible/playbooks/deploy.yml`
- configures the APT repository on the Linux validation server if needed
- installs or upgrades the `iot4b` package through APT
- verifies that `iot4b.service` is running

## Post-Release Checks

After the workflow finishes, verify the following:

- a GitHub Release exists for the tag
- the release contains `.deb`, `.ipk`, and `checksums.txt`
- `http://repo.iot4b.co/apt` contains updated APT metadata
- `http://repo.iot4b.co/opkg` contains updated OPKG metadata and installer script
- the Homebrew tap repository contains the updated `Formula/iot4b.rb`
- the Linux validation server has the expected package version installed
- `systemctl status iot4b` is healthy on the Linux validation server

## Notes

- The automated Linux validation deployment happens only on tag push, not on push to `main`.
- The Linux validation server is updated from the published APT repository, not from a copied binary.
- Tag format should always be `vX.Y.Z`.
- Release notes are generated automatically by GitHub from the tag range and repository history.
- For better generated notes, use Conventional Commits in commit and PR titles.
- See [conventional-commits.md](conventional-commits.md).
