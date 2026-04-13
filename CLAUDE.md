# CLAUDE.md — GreenNode CLI

## Project overview

GreenNode CLI (`grn`) is a unified command-line tool for managing GreenNode (VNG Cloud) services. Written in Go, distributed as a single binary. VKS (VNG Kubernetes Service) is the first service; other product teams add their own services.

- **Repo**: `vngcloud/greennode-cli`
- **Docs**: https://vngcloud.github.io/greennode-cli/
- **Language**: Go (using cobra CLI framework)
- **Binary**: Single file, zero runtime dependencies

## Project structure

```
go/
├── main.go                          # Entry point
├── cmd/
│   ├── root.go                      # Root command + global flags + --version
│   ├── configure/
│   │   ├── configure.go             # Interactive setup
│   │   ├── list.go                  # grn configure list
│   │   ├── get.go                   # grn configure get
│   │   └── set.go                   # grn configure set
│   └── vks/
│       ├── vks.go                   # VKS parent command
│       ├── helpers.go               # Shared utilities (client, output, parsing)
│       ├── list_clusters.go         # Auto-pagination
│       ├── get_cluster.go
│       ├── create_cluster.go        # --dry-run validation
│       ├── update_cluster.go
│       ├── delete_cluster.go        # Confirm + --force + --dry-run
│       ├── list_nodegroups.go
│       ├── get_nodegroup.go
│       ├── create_nodegroup.go
│       ├── update_nodegroup.go
│       ├── delete_nodegroup.go
│       ├── wait_cluster_active.go   # Polling waiter
│       └── auto_upgrade.go          # Set/delete auto-upgrade
├── internal/
│   ├── config/
│   │   ├── config.go               # Config + credentials loading (INI)
│   │   └── writer.go               # ConfigFileWriter (0600 perms)
│   ├── auth/token.go                # OAuth2 Client Credentials (IAM)
│   ├── client/client.go             # HTTP client with retry + auto-refresh
│   ├── formatter/formatter.go       # JSON/Table/Text + JMESPath
│   └── validator/validator.go       # ID format validation
├── go.mod, go.sum
```

## Code conventions

- All source code text in **English**
- Use cobra for all commands
- Internal packages in `internal/` (not importable externally)
- Commands in `cmd/` following cobra patterns
- Use `cobra.Command` with `RunE` for error handling

## VNG Cloud API quirks

- **IAM API uses camelCase**: `grantType`, `accessToken`, `expiresIn`
- **VKS API pagination is 0-based**: page 0 = first page
- **`--version` conflict**: Use `--k8s-version` for Kubernetes version

## Adding a new command

1. Create file in `cmd/vks/<command_name>.go`
2. Define `cobra.Command` with Use, Short, RunE
3. Add flags via `cmd.Flags()`
4. Register in `cmd/vks/vks.go` init(): `VksCmd.AddCommand(newCmd)`
5. Add `validator.ValidateID()` for any ID args in URLs
6. Add `--dry-run` for create/update/delete commands
7. Add `--force` + confirmation for delete commands

## Adding a new service

1. Create `cmd/<service>/` directory
2. Create parent command file with `cobra.Command`
3. Register in `cmd/root.go` init(): `rootCmd.AddCommand(serviceCmd)`

## Security rules

- **Credential masking**: `configure list` and `configure get` mask client_id/client_secret (last 4 chars only)
- **No credential env vars**: `GRN_CLIENT_ID`/`GRN_CLIENT_SECRET` not supported — file only
- **Input validation**: All cluster-id/nodegroup-id validated via `validator.ValidateID()` before URLs
- **SSL default on**: `--no-verify-ssl` prints warning to stderr
- **Tokens in memory only**: Never written to disk or logged
- **File permissions**: Credentials file created with 0600, directory 0700

## Building

```bash
cd go
CGO_ENABLED=0 go build -o grn .

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o grn-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o grn-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o grn-windows-amd64.exe .
```

## Git workflow

- **Do not auto commit/push** — only change source code, user will ask for commit/push
- **Main branch is protected** — must use PR
- **Changelog**: `./scripts/new-change` for every change
- **Release**: `./scripts/bump-version minor` → `git push && git push --tags`

## Documentation update rule

**After ANY change to business logic, security, configuration, or commands:**

1. Review ALL docs below and update what's affected
2. If unsure whether a doc needs updating, read it and check

**Docs to check:**

- `docs/` (GitHub Pages) — command references, usage guides
- `README.md`
- `CLAUDE.md`
- `CONTRIBUTING.md`
- `docs/DEVELOPMENT.md`
- `./scripts/new-change` — changelog fragment

Code without docs is not done.

## Key files

| File | Purpose |
|------|---------|
| `cmd/root.go` | Root command, global flags, --version |
| `cmd/vks/helpers.go` | Client creation, output formatting, label/taint parsing |
| `internal/config/config.go` | Config loading from ~/.greenode/, REGIONS map |
| `internal/config/writer.go` | INI file writer with 0600 perms |
| `internal/auth/token.go` | TokenManager — OAuth2 with IAM (camelCase) |
| `internal/client/client.go` | HTTP client with retry (3x backoff) + 401 refresh |
| `internal/formatter/formatter.go` | JSON/Table/Text + JMESPath |
| `internal/validator/validator.go` | ID format validation |
