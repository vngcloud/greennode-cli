# vServer CLI — product notes

Tier-2 notes for the `grn vserver` command group. The root `CLAUDE.md` has the
repo-wide conventions; this file records what is specific to vServer.

## Client & endpoint (differs from vks)

- Build the client with **`vserverclient.BuildClient(cmd)`**, which returns
  `(client, cfg, error)` — NOT `cli.NewClient`. Most helpers wrap it (see
  `helpers.go` in each subpackage).
- Endpoint: `internal/config` REGIONS key `vserver_endpoint`
  (`https://hcm-3.api.vngcloud.vn/vserver/vserver-gateway`).
- **`project_id` is required** and goes in the URL path. Get it via
  `vserverclient.ProjectID(cfg)` (from `grn configure` or `GRN_DEFAULT_PROJECT_ID`);
  it errors clearly if unset.

## API quirks (differ from vks — do not copy vks assumptions)

- **Pagination is 1-based**: params `page` (from **1**) + **`size`** (note: `size`,
  not `pageSize`).
- Paths are under **`/v2/{projectID}/...`** (servers, volumes, networks, secgroups),
  except **images** which are under `/v1/{projectID}/images/...`.
- Subnets carry both `uuid` and `id`; list/complete code prefers `uuid`, falling
  back to `id`.

## Command style (differs from vks)

- vServer uses **noun-group + verb**: `grn vserver server create`, `volume delete`,
  `secgroup rule create` — versus vks's flat `create-nodegroup`. Both satisfy the
  conventions test (the leaf verb is canonical). Keep new vServer commands in this
  nested style for consistency within the product.
- Destructive commands (`delete`, `server stop`/`reboot`) carry `--dry-run` +
  `--force` (see the conventions test).

## Output

- List/detail commands render selected columns via `vserverclient.OutputWithColumns`
  (see the per-subpackage `helpers.go` column definitions).
