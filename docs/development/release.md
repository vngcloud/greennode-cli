# Release Process

Versioning and the changelog are automated by [release-please](https://github.com/googleapis/release-please).
**Never edit the version (`go/cmd/root.go`) or `CHANGELOG.md` by hand**, and never
create tags manually — release-please derives everything from commit messages.

## During development

- Branch from `main`, open a PR with a **Conventional Commits** title
  (`feat:`, `fix:`, `docs:`, `feat!:`, …). The `Conventional Commits title` check
  enforces this.
- PRs are **squash-merged**, so the PR title becomes the commit message that
  release-please reads.

Bump rules: `fix:` → patch, `feat:` → minor, `feat!:` / `BREAKING CHANGE:` → major.
Wrong bump? Fix the PR title — do not touch the version file.

## Cutting a release

1. As `feat`/`fix` PRs merge to `main`, release-please opens/refreshes a
   **`chore: release main`** PR that bumps the version in `go/cmd/root.go` and
   updates `CHANGELOG.md`.
2. **Squash-merge that release PR.** release-please then tags `vX.Y.Z` and creates
   the GitHub Release; the release workflow builds the multi-platform binaries and
   attaches them.

> If the release PR's checks don't start (a limitation of the default
> `GITHUB_TOKEN`), close & reopen it once — or configure a `RELEASE_PLEASE_TOKEN`
> PAT (contents: write, pull-requests: write) so checks run automatically.

## Manual / emergency release

`release.yml` still accepts a tag push (`git tag vX.Y.Z && git push --tags`) or a
`workflow_dispatch`, which builds the binaries and creates the release directly.
Prefer the release-please flow for normal releases.

## CI/CD workflows

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `run-tests.yml` | PR + push | Build + test Go binary |
| `pr-title.yml` | PR opened/edited | Enforce Conventional Commits PR title |
| `release-please.yml` | Push to `main` | Maintain the release PR; on merge, tag + release + build binaries |
| `release.yml` | Tag push `v*`, dispatch, or called by release-please | Build multi-platform binaries + attach to the release |
| `deploy-docs.yml` | Push to main (docs/) | Deploy documentation to GitHub Pages |
| `stale.yml` | Daily schedule | Auto-close stale issues |
