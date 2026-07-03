# update-cluster

## Description

Update a VKS cluster's Kubernetes version, node CIDR whitelist, and plugin configuration. The cluster ID, target Kubernetes version, and at least one whitelist CIDR are always required by the API, even when the intent is only to toggle a plugin.

Plugin flags (`--load-balancer-plugin`, `--block-store-csi-plugin`) are only sent when explicitly provided; omitting them leaves the current plugin state unchanged.

Use `--dry-run` to preview the update payload without executing the request.

## Synopsis

```
grn vks update-cluster
    --cluster-id <value>
    --k8s-version <value>
    --whitelist-node-cidrs <value>
    [--load-balancer-plugin <enabled|disabled>]
    [--block-store-csi-plugin <enabled|disabled>]
    [--dry-run]
```

## Options

**`--cluster-id`** (string)

ID of the cluster to update.

- Required: Yes

**`--k8s-version`** (string)

Target Kubernetes version (e.g. `v1.29.1`). Must be the same or a higher version than the cluster's current version.

- Required: Yes
- Constraints: 1–50 characters.
- See available versions with [list-cluster-versions](list-cluster-versions.md).

**`--whitelist-node-cidrs`** (list&lt;string&gt;)

CIDRs allowed to communicate with cluster nodes, comma-separated. At least one value is required.

- Required: Yes
- Constraints: 1–30 entries.
- Syntax: `10.0.0.0/8,192.168.0.0/16`

**`--load-balancer-plugin`** (string)

Load balancer plugin state. When omitted, the current state is left unchanged.

- Required: No
- Possible values: `enabled`, `disabled`

**`--block-store-csi-plugin`** (string)

Block store CSI plugin state. When omitted, the current state is left unchanged.

- Required: No
- Possible values: `enabled`, `disabled`

**`--dry-run`** (boolean)

Print the update payload without sending the request.

- Required: No
- Default: `false`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

Upgrade the Kubernetes version and set whitelist CIDRs:

```bash
grn vks update-cluster \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --k8s-version v1.30.0 \
  --whitelist-node-cidrs 10.0.0.0/8,192.168.0.0/16
```

Update cluster and disable the load balancer plugin:

```bash
grn vks update-cluster \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --k8s-version v1.29.1 \
  --whitelist-node-cidrs 10.0.0.0/8 \
  --load-balancer-plugin disabled
```

Preview what would be sent (dry run):

```bash
grn vks update-cluster \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --k8s-version v1.30.0 \
  --whitelist-node-cidrs 10.0.0.0/8 \
  --dry-run
```
