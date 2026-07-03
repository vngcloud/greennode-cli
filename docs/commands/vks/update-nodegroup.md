# update-nodegroup

## Description

Update a node group's desired node count, security groups, auto-scaling configuration, and upgrade configuration. At least one of `--num-nodes`, `--security-groups`, `--auto-scale`, or `--upgrade-config` must be provided.

To update labels, tags, or taints, use [update-nodegroup-metadata](update-nodegroup-metadata.md) — those fields are deprecated on this command.

Use `--dry-run` to preview the update payload without executing the request.

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

**`--cluster-id`** (string)

ID of the cluster that owns the node group.

- Required: Yes

**`--nodegroup-id`** (string)

ID of the node group to update.

- Required: Yes

**`--num-nodes`** (string)

New desired number of nodes. Parsed as an integer by the CLI.

- Required: Conditional — at least one of `--num-nodes`, `--security-groups`, `--auto-scale`, or `--upgrade-config` must be provided.
- Constraints: 0–10. When `--auto-scale` is also set, must be within `[minSize, maxSize]`.

**`--security-groups`** (list&lt;string&gt;)

Security group IDs to replace the current set, comma-separated.

- Required: Conditional — at least one update flag must be provided.
- Constraints: 1–50 entries.
- Syntax: `secg-aaa111,secg-bbb222`

**`--auto-scale`** (structure)

Auto-scaling configuration for the node group. Accepts shorthand or JSON.

- Required: Conditional — at least one update flag must be provided.
- Members:
    - `minSize` (integer) — minimum number of nodes; minimum value `0`
    - `maxSize` (integer) — maximum number of nodes; minimum value `1`

Shorthand syntax:

```
minSize=2,maxSize=10
```

JSON syntax:

```json
{"minSize": 2, "maxSize": 10}
```

**`--upgrade-config`** (structure)

Upgrade strategy configuration for the node group. Accepts shorthand or JSON.

- Required: Conditional — at least one update flag must be provided.
- Members:
    - `strategy` (string) — upgrade strategy; currently only `SURGE` is supported
    - `maxSurge` (integer) — maximum number of extra nodes added during upgrade; range 1–100
    - `maxUnavailable` (integer) — maximum number of nodes that may be unavailable during upgrade; range 0–100

Shorthand syntax:

```
maxSurge=1,maxUnavailable=0,strategy=SURGE
```

JSON syntax:

```json
{"maxSurge": 1, "maxUnavailable": 0, "strategy": "SURGE"}
```

**`--dry-run`** (boolean)

Print the update payload without sending the request.

- Required: No
- Default: `false`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

Scale a node group to 5 nodes:

```bash
grn vks update-nodegroup \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --num-nodes 5
```

Enable auto-scaling with min/max limits:

```bash
grn vks update-nodegroup \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --auto-scale minSize=2,maxSize=10
```

Set the upgrade configuration using JSON:

```bash
grn vks update-nodegroup \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --upgrade-config '{"maxSurge":2,"maxUnavailable":1,"strategy":"SURGE"}'
```

Preview the update payload (dry run):

```bash
grn vks update-nodegroup \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --num-nodes 3 \
  --dry-run
```
