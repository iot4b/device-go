# Scripts

This directory contains helper scripts built around one shared responsibility:

- calculate release/build metadata
- propagate that metadata into local builds and runs
- reuse the same metadata in packaging and release automation

The central piece is `release-vars.sh`. The other scripts either consume that metadata directly or use it to generate release artifacts.

## Shared Build Metadata

### `release-vars.sh`

Calculates the shared metadata used across local scripts, package builds, and CI:

- `IOT4B_VERSION`
- `IOT4B_GIT_TAG`
- `IOT4B_COMMIT`
- `IOT4B_BUILD_DATE`

This is the single source of truth for release/build metadata in the repository.

## Local Build And Run

### `build.sh`

Builds the `iot4b` binary locally with embedded build metadata from `release-vars.sh`.

Example:

```sh
./scripts/build.sh
./scripts/build.sh /tmp/iot4b
```

### `run.sh`

Builds a temporary binary with embedded build metadata and runs it.

Example:

```sh
./scripts/run.sh
./scripts/run.sh --version
./scripts/run.sh setup
```

### `version.sh`

Prints the calculated build metadata without building the binary.

Example:

```sh
./scripts/version.sh
```

Use these scripts instead of plain `go run` or raw `go build` when you want the binary to contain correct build metadata.

You can override the defaults for local testing:

```sh
IOT4B_VERSION=2.1.0-rc.1 ./scripts/build.sh
IOT4B_VERSION=2.1.0-rc.1 ./scripts/run.sh --version
```

Supported override variables:

- `IOT4B_VERSION`
- `IOT4B_COMMIT`
- `IOT4B_BUILD_DATE`

## Release And Packaging

### `generate-homebrew-formula.sh`

Generates the Homebrew formula content for `iot4b`.

Inputs:

- release version
- source archive URL
- source archive SHA256

The release workflow uses this script to produce `Formula/iot4b.rb` for the tap repository.

## Notes

- Prefer these scripts over ad hoc local commands when working on versioned local builds and release-related tasks.
- The same metadata flow is shared by local development, package builds, and the GitHub release workflow.
