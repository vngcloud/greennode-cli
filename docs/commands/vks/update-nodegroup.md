# update-nodegroup

## Description

Update a node group's node count, security groups, auto-scaling configuration, and upgrade configuration.

To update labels, tags, or taints, use `grn vks update-nodegroup-metadata` — those fields are deprecated on `update-nodegroup`.

Use `--dry-run` to preview the update payload without executing it.

## Synopsis

```
grn vks update-nodegroup
    --cluster-id <value>
    --nodegroup-id <value>
    [--num-nodes <value>]
    [--security-groups <value>]
    [--auto-scale <value>]
    [--upgrade-config <value>]
    [--dry-run]
```

## Options

`--cluster-id` (required)
: ID of the cluster that owns the node group.

`--nodegroup-id` (required)
: ID of the node group to update.

`--num-nodes` (optional)
: New desired number of nodes for the node group.

`--security-groups` (optional)
: Comma-separated list of security group IDs to replace the current set.

`--auto-scale` (optional)
: Auto-scale configuration. Shorthand `minSize=2,maxSize=10` or JSON `{"minSize":2,"maxSize":10}`.

`--upgrade-config` (optional)
: Upgrade configuration. Shorthand `maxSurge=1,maxUnavailable=0,strategy=SURGE` or JSON `{"maxSurge":1,"maxUnavailable":0,"strategy":"SURGE"}`.

`--dry-run` (optional)
: Print the update payload without sending the request.

## Examples

Scale a node group to 5 nodes:

```bash
grn vks update-nodegroup \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --num-nodes 5
```

Set auto-scaling limits (shorthand or JSON):

```bash
grn vks update-nodegroup \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --auto-scale minSize=2,maxSize=10
```

Set the upgrade configuration:

```bash
grn vks update-nodegroup \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --upgrade-config '{"maxSurge":2,"maxUnavailable":1,"strategy":"SURGE"}'
```

To update labels, tags, or taints, use `update-nodegroup-metadata` (those fields are deprecated on `update-nodegroup`):

```bash
grn vks update-nodegroup-metadata \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --labels env=prod,tier=app \
  --taints dedicated=gpu:NoSchedule
```

Preview the update payload (dry run):

```bash
grn vks update-nodegroup \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --num-nodes 3 \
  --dry-run
```
