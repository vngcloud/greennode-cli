# VKS CLI — product notes

Tier-2 notes for the `grn vks` command group. The root `CLAUDE.md` has the
repo-wide conventions; this file records what is specific to VKS.

## Client & endpoint

- `createClient(cmd)` → `cli.NewClient(cmd, "vks")`. The endpoint is resolved from
  `internal/config` REGIONS key `vks_endpoint`
  (`https://vks.api.vngcloud.vn` HCM-3, `https://vks-han-1.api.vngcloud.vn` HAN).
- Auth: shared GreenNode IAM bearer token (see `internal/auth`).

## API quirks

- **Pagination is 0-based**: params `page` (from 0) + `pageSize`. The `list-*`
  commands default `--page -1`, a sentinel meaning "auto-paginate all pages".
- Paths are under `/v1/clusters/...`.
- Cluster/node-group **status** values: `CREATING`, `ACTIVE`, `ERROR` (waiters key
  on these; see `wait.go`).
- **Kubeconfig is asynchronous**: `generate-kubeconfig` (POST) requests generation,
  then `update-kubeconfig` (GET) fetches + merges once status is `ACTIVE`
  (`NONE`/`CREATING`/`ERROR` are the other states).

## Command specifics

- `create-cluster` creates the **control plane only**. The API's `nodeGroups` array
  is deprecated — do not send it; add workers with `create-nodegroup`.
- `update-nodegroup`: all body fields optional. Labels/tags/taints are **deprecated**
  on this endpoint — use `update-nodegroup-metadata`.
- `config-auto-upgrade` (the former `set-auto-upgrade-config` is a deprecated alias).
- `--k8s-version` values come from `list-cluster-versions`.
- Network types: `TIGERA`, `CILIUM_OVERLAY`, `CILIUM_NATIVE_ROUTING`; `--cidr` is
  required for `TIGERA`/`CILIUM_OVERLAY`. OS images: `ubuntu`, `linux`, `rocky`.

## Struct-valued flags (shorthand `k=v,k2=v2` or JSON)

- `--auto-scale` → `autoScaleConfig` (int: `minSize`, `maxSize`)
- `--upgrade-config` → `upgradeConfig` (int: `maxSurge`, `maxUnavailable`; string `strategy`)
- `--placement-group` → `placementGroupConfigDto` (strings)
- `--auto-upgrade-config` → `autoUpgradeConfig` (`weekdays` — use JSON for multiple days; `time`)
- `--auto-healing-config` → `autoHealingConfig` (bool `enableAutoHealing`, int `timeoutUnhealthy`)
