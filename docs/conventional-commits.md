# Conventional Commits

This repository uses Conventional Commits for new changes.

The goal is to keep commit history readable and make automatically generated GitHub release notes more useful.

## Format

Use this structure:

```text
type(scope): short summary
```

Examples:

```text
feat(release): deploy validation server from apt package
fix(apt): generate release metadata from tag version
docs(readme): add opkg installation instructions
chore(ci): rename release jobs to publish and deploy
refactor(ansible): simplify deploy inventory
```

## Recommended Types

- `feat`: user-facing feature or new capability
- `fix`: bug fix
- `docs`: documentation only
- `chore`: maintenance, CI, tooling, repo housekeeping
- `refactor`: code or config restructuring without intended behavior change
- `test`: tests only
- `build`: packaging or build-system changes
- `ci`: GitHub Actions or deployment pipeline changes

## Scope

Use a short scope when it helps:

- `release`
- `apt`
- `opkg`
- `brew`
- `ansible`
- `ci`
- `readme`
- `config`

## Summary Style

- keep the subject short
- use the imperative mood
- do not end the subject with a period
- prefer one logical change per commit

## Why This Helps

GitHub auto-generated release notes do not require Conventional Commits, but they become much easier to scan when commit titles and PR titles follow a consistent format.

Clear Conventional Commit titles make it easier to understand:

- what changed
- which subsystem changed
- whether the change is a feature, fix, docs update, or CI change

## Practical Workflow

Before committing, write the message in Conventional Commits style.

Example:

```sh
git commit -m "fix(opkg): publish packages under /opkg path"
```
