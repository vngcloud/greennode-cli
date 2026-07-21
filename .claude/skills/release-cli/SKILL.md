---
name: release-cli
description: Use when releasing greennode-cli, cutting a version, tagging, or publishing binaries ("release the CLI", "cut a release", "bump version", "publish grn"). Everything is release-please driven — never bump versions or tag by hand.
---

# Releasing greennode-cli

Versioning and `CHANGELOG.md` are automated by **release-please**. **Never** edit
the version (`go/cmd/root.go` `cliVersion`), edit `CHANGELOG.md`, or create tags by
hand — release-please derives all of it from commit messages.

## How a release happens

1. `feat:` / `fix:` PRs merge to `main` (squash-merged; the PR title is the commit
   message, enforced by the `Conventional Commits title` check).
   - `fix:` → patch, `feat:` → minor, `feat!:` / `BREAKING CHANGE:` → major.
2. release-please opens/refreshes a **`chore: release main`** PR that bumps
   `go/cmd/root.go` and updates `CHANGELOG.md`.
3. **Squash-merge that release PR** to publish: it tags `vX.Y.Z`, creates the GitHub
   Release, and `release.yml` builds + attaches the multi-platform binaries.

## Verify

```bash
gh pr list                                  # find the "chore: release main" PR
gh run list --workflow=release-please.yml   # after merge: release + build
gh release view vX.Y.Z                       # binaries attached
```

## Troubleshooting

- **Release PR checks don't start** → default `GITHUB_TOKEN` limitation; close &
  reopen the PR once, or set a `RELEASE_PLEASE_TOKEN` PAT (contents + pull-requests
  write) so checks auto-run.
- **Wrong version bump** → fix the offending PR title; never edit the version file.
- **Emergency/manual**: `release.yml` still accepts a tag push (`git tag vX.Y.Z &&
  git push --tags`) or `workflow_dispatch`. Prefer release-please for normal releases.
