# update-nodegroup-metadata

## Description

Update the labels, tags, and taints of a node group. At least one of `--labels`, `--tags`, or `--taints` must be provided.

## Synopsis

```
grn vks update-nodegroup-metadata
    --cluster-id <value>
    --nodegroup-id <value>
    [--labels <value>]
    [--tags <value>]
    [--taints <value>]
```

## Options

`--cluster-id` (required)
: The ID of the cluster.

`--nodegroup-id` (required)
: The ID of the node group to update.

`--labels` (optional)
: Comma-separated `key=value` pairs to set as Kubernetes node labels (e.g. `env=prod,tier=app`).

`--tags` (optional)
: Comma-separated `key=value` pairs to set as tags (e.g. `team=platform,cost-center=123`).

`--taints` (optional)
: Comma-separated node taints in `key=value:effect` format (e.g. `dedicated=gpu:NoSchedule`).

At least one of `--labels`, `--tags`, or `--taints` must be provided.

## Examples

Update node labels:

```bash
grn vks update-nodegroup-metadata \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --labels env=prod,tier=app
```

Update labels, tags, and taints together:

```bash
grn vks update-nodegroup-metadata \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --labels tier=gpu \
  --tags team=platform \
  --taints dedicated=gpu:NoSchedule
```
