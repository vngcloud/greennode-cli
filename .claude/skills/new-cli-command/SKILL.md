---
name: new-cli-command
description: Use when adding a new command or a new product command group to greennode-cli ("add a command", "tạo CLI cho <product>", "implement grn <product>", "new subcommand"). Guides scaffolding and the command implementation flow.
---

# Adding to greennode-cli

`grn` is a single Go binary; each product is a self-registering command group under
`go/cmd/<product>/`. Read the root `CLAUDE.md` (repo conventions) and the target
product's `go/cmd/<product>/CLAUDE.md` (API quirks) first.

## New product CLI → scaffold

```bash
./scripts/new-product <product>      # lowercase, valid Go package name (e.g. vdb)
```

Generates the package (parent command, example command, helpers, starter test,
product CLAUDE.md), docs, self-registration in `go/cmd/register.go`, and CODEOWNERS
lines. Then: add `<product>_endpoint` to `internal/config` REGIONS, replace the
`list-examples` starter, and fill in the product CLAUDE.md.

## New command in an existing product (TDD)

1. Add `go/cmd/<product>/<verb>_<noun>.go`; register it in the product's parent
   `init()`.
2. Follow the conventions (the `go/cmd/conventions_test.go` test enforces them):
   - Name is `verb-noun` (canonical verbs: list/get/create/update/delete/configure/
     upgrade/generate/validate/wait). vServer uses `noun verb` nesting — match the
     product's existing style.
   - `validator.ValidateID(id, "flag")` before interpolating an ID into a URL.
   - Struct-valued flags via `cli.ParseStructFlag` / `ParseStructFlagTyped`
     (shorthand `k=v,k2=v2` or JSON).
   - Destructive commands (`delete`, `stop`, `reboot`): add `--dry-run` and
     `--force`; use `cli.DryRunNotice` / `cli.Confirm`.
   - Output via `cli.Output` (product helper `outputResult`).
3. Write table-driven tests (use `httptest`); build/test with the macOS linker
   workaround: `CGO_ENABLED=1 go test -ldflags='-linkmode=external' ./...`.
4. Add an AWS-CLI-style doc page under `docs/commands/<product>/` (type, Required,
   Default, Possible values, Constraints; Members + Shorthand/JSON for struct flags)
   and link it in that product's `index.md`.

## PR

Open a PR with a **Conventional Commit** title (`feat(...)`, `fix(...)`) — it becomes
the squash-merge commit release-please reads. Do not edit the version or CHANGELOG.
