# upgrade-nodegroup-version

## Description

Upgrade the Kubernetes version of a node group within a cluster.

Use `grn vks list-cluster-versions` to find the valid Kubernetes versions before running this command.

## Synopsis

```
grn vks upgrade-nodegroup-version
    --cluster-id <value>
    --nodegroup-id <value>
    --k8s-version <value>
```

## Options

`--cluster-id` (required)
: The ID of the cluster.

`--nodegroup-id` (required)
: The ID of the node group to upgrade.

`--k8s-version` (required)
: The target Kubernetes version. Use `grn vks list-cluster-versions` to find valid versions.

## Examples

Upgrade a node group to a specific Kubernetes version:

```bash
grn vks upgrade-nodegroup-version \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --k8s-version v1.29.0
```
