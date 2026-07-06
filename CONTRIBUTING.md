# Contributing to GreenNode CLI

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git

### Setup development environment

```bash
git clone https://github.com/vngcloud/greennode-cli.git
cd greennode-cli/go
go build -o grn .
./grn --version
```

### Build

```bash
cd go
CGO_ENABLED=0 go build -o grn .
```

## Adding a new product CLI

Each product is a self-registering command group under `go/cmd/<product>/`, mounted
into the single `grn` binary. Scaffold one in seconds:

```bash
./scripts/new-product <product>      # e.g. vdb  (lowercase, valid Go package name)
```

This generates `go/cmd/<product>/` (parent command, an example command, shared
`createClient`/`outputResult` helpers, a starter test, and a product `CLAUDE.md`),
creates `docs/commands/<product>/`, wires self-registration in `go/cmd/register.go`,
and appends `.github/CODEOWNERS` lines. CI and the conventions test pick up the new
package automatically. Then follow the "Next steps" the script prints (add your
`<product>_endpoint` to `internal/config` REGIONS, replace the example command, etc.).

## Development Workflow

### 1. Create a feature branch

```bash
git checkout main && git pull
git checkout -b feat/your-feature-name
```

### 2. Make changes and build

```bash
cd go
# Write code
CGO_ENABLED=0 go build -o grn .
./grn vks <your-command> --help
```

### 3. Add a changelog entry

Every PR should include a changelog fragment:

```bash
./scripts/new-change -t feature -c vks -d "Add your feature description"
```

Change types: `feature`, `bugfix`, `enhancement`, `api-change`

### 4. Commit and push

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(vks): add describe-events command
fix(auth): fix token refresh race condition
docs(readme): update installation instructions
```

### 5. Create a Pull Request

- PR to `main` when release-ready
- CI must pass before merge
- At least 1 approval required

## Adding a New Command

1. Create `go/cmd/vks/<command_name>.go`
2. Define `cobra.Command` with Use, Short, RunE
3. Register in `go/cmd/vks/vks.go`: `VksCmd.AddCommand(newCmd)`
4. Add `validator.ValidateID()` for any ID args
5. Add `--dry-run` for create/update/delete
6. Add `--force` + confirmation for delete
7. Create `docs/commands/vks/<command-name>.md` — command reference page
8. Add entry to `docs/commands/vks/index.md` table
9. Add nav entry to `mkdocs.yml`

## Adding a New Service

1. Create `go/cmd/<service>/` directory
2. Create parent command with `cobra.Command`
3. Register in `go/cmd/root.go`: `rootCmd.AddCommand(serviceCmd)`

## Code Style

- All source code text in English
- Use cobra patterns for all commands
- Validate user inputs (IDs used in URLs)
- Use `--dry-run` for create/update/delete commands
- Add `--force` to skip confirmation on delete commands

## License

By contributing, you agree that your contributions will be licensed under the [Apache License 2.0](LICENSE).
