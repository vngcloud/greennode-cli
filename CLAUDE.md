# CLAUDE.md — GreenNode CLI

## Project overview

GreenNode CLI (`grn`) is a unified command-line tool for managing GreenNode (VNG Cloud) services. Architecture cloned from AWS CLI with hand-written commands. VKS (VNG Kubernetes Service) is the first service; other product teams add their own services by following the same pattern.

- **Repo**: `vngcloud/greennode-cli`
- **Docs**: https://vngcloud.github.io/greennode-cli/
- **PyPI**: https://pypi.org/project/grncli/

## Code conventions

- All source code text must be in **English** — error messages, descriptions, comments, docstrings, ARG_TABLE help_text
- Follow existing AWS CLI patterns: `CLIDriver` → `ServiceCommand` → `BasicCommand`
- Each command file has one class, one responsibility
- Use `display_output(result, parsed_globals)` helper for API response formatting

## VNG Cloud API quirks

- **IAM API uses camelCase**: `grantType`, `accessToken`, `expiresIn` (not snake_case OAuth2 standard)
- **VKS API pagination is 0-based**: page 0 = first page
- **`--version` conflict**: Use `--k8s-version` for Kubernetes version to avoid clash with global `--version` flag

## Adding a new command

1. Create file in `grncli/customizations/vks/<command_name>.py`
2. Extend `BasicCommand` with `NAME`, `DESCRIPTION`, `ARG_TABLE`
3. Implement `_run_main(self, parsed_args, parsed_globals)`
4. Register in `grncli/customizations/vks/__init__.py`
5. Add `validate_id()` calls for any ID args used in URLs
6. Add `--dry-run` for create/update/delete commands
7. Add `--force` + confirmation prompt for delete commands

## Adding a new service

1. Create `grncli/customizations/<service>/`
2. Write commands extending `BasicCommand`
3. Register in `grncli/handlers.py`
4. See `grncli/customizations/vks/` for reference

## Security rules

- **Credential masking**: `grn configure list` and `grn configure get` must mask `client_id`/`client_secret` (show last 4 chars only)
- **Input validation**: All cluster-id and nodegroup-id args must be validated via `validators.validate_id()` before constructing URLs — prevents path traversal
- **SSL default on**: `--no-verify-ssl` must print warning to stderr
- **Tokens in memory only**: Never write tokens to disk or logs
- **Dependency pinning**: Pin to major versions (`httpx<1.0`, `PyYAML<7.0`)

## Testing

```bash
python -m pytest tests/ -v
```

- Tests must pass on Python 3.10-3.13 × Ubuntu/macOS/Windows
- Skip Unix-only tests on Windows with `@pytest.mark.skipif(platform.system() == 'Windows', ...)`

## Git workflow

- **Do not auto commit/push** — only change source code, user will ask for commit/push when ready
- **Branches**: `main` (production), `develop` (testing), `feat/*` or `fix/*` (feature/bug branches)
- **PRs**: feature → develop (test), feature → main (release-ready)
- **Changelog**: Add fragment via `./scripts/new-change` for every change
- **Release**: `./scripts/bump-version minor` → `git push && git push --tags`
- **Main branch is protected** — cannot push directly, must use PR

## Documentation update rule

**After completing any feature or bugfix, update ALL related documentation before considering the work done:**

1. **GitHub Pages docs** (`docs/`):
   - Add/update command reference page in `docs/commands/vks/<command>.md`
   - Update `docs/commands/vks/index.md` command table
   - Update relevant usage guides if behavior changes (pagination, dry-run, etc.)
   - Update `mkdocs.yml` nav if new pages added

2. **CHANGELOG**: Add changelog fragment via `./scripts/new-change`

3. **Spec** (`docs/superpowers/specs/2026-04-10-greenode-cli-design.md`):
   - Update command list in Section 4
   - Update file structure in Section 2 if new files added

This is not optional. Code without docs is not done.

## Key files

| File | Purpose |
|------|---------|
| `grncli/clidriver.py` | CLIDriver + ServiceCommand — main orchestrator |
| `grncli/session.py` | Config, credentials, region, endpoints, SSL, timeouts |
| `grncli/auth.py` | TokenManager — OAuth2 Client Credentials with IAM |
| `grncli/client.py` | HTTP client with retry (3x backoff) + auto token refresh |
| `grncli/customizations/commands.py` | BasicCommand base class + display_output + help system |
| `grncli/customizations/vks/validators.py` | ID format validation |
| `grncli/data/cli.json` | Global CLI options (AWS CLI style) |
| `mkdocs.yml` | Documentation site config |
| `scripts/bump-version` | Bump version + merge changelog + commit + tag |
| `scripts/new-change` | Create changelog fragment |
