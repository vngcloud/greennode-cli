# upgrade-nodegroup-version

## Description

Upgrade the Kubernetes version of a node group within a cluster. The target version must be compatible with the cluster's current control-plane version.

Use [list-cluster-versions](list-cluster-versions.md) to find valid Kubernetes versions before running this command.

Use `--dry-run` to preview the upgrade payload without executing the request.

## Synopsis

```
grn vks upgrade-nodegroup-version
    --cluster-id <value>
    --nodegroup-id <value>
    --k8s-version <value>
    [--dry-run]
```

## Options

**`--cluster-id`** (string)

ID of the cluster that owns the node group.

- Required: Yes

**`--nodegroup-id`** (string)

ID of the node group to upgrade.

- Required: Yes

**`--k8s-version`** (string)

Target Kubernetes version for the node group (e.g. `v1.29.13-vks.1740045600`).

- Required: Yes
- See available versions with [list-cluster-versions](list-cluster-versions.md).

**`--dry-run`** (boolean)

Print the upgrade payload without sending the request.

- Required: No
- Default: `false`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

Upgrade a node group to a specific Kubernetes version:

```bash
grn vks upgrade-nodegroup-version \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --k8s-version v1.29.13-vks.1740045600
```

Preview the upgrade payload (dry run):

```bash
grn vks upgrade-nodegroup-version \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --k8s-version v1.29.13-vks.1740045600 \
  --dry-run
```
