# update-nodegroup-metadata

## Description

Update the labels, tags, and taints of a node group without touching its infrastructure. This is a partial-update (PATCH) operation: only the fields explicitly provided are changed; absent fields are left as-is. Sending an empty value for a field (e.g. `--labels ""`) clears that field.

At least one of `--labels`, `--tags`, or `--taints` must be provided.

Use `--dry-run` to preview the metadata patch without executing the request.

## Synopsis

```
grn vks update-nodegroup-metadata
    --cluster-id <value>
    --nodegroup-id <value>
    [--labels <value>]
    [--tags <value>]
    [--taints <value>]
    [--dry-run]
```

## Options

**`--cluster-id`** (string)

ID of the cluster that owns the node group.

- Required: Yes

**`--nodegroup-id`** (string)

ID of the node group to update.

- Required: Yes

**`--labels`** (map)

Kubernetes node labels as comma-separated `key=value` pairs. Replaces the existing label set on all nodes in the group.

- Required: Conditional — at least one of `--labels`, `--tags`, or `--taints` must be provided.
- Constraints: maximum 50 entries. Each key may have an optional DNS-subdomain prefix (`prefix/name`); the name segment must be 63 characters or less, alphanumeric with `-`, `_`, `.` allowed.
- Syntax: `env=prod,tier=app`

**`--tags`** (map)

Cloud tags for the node group resources as comma-separated `key=value` pairs. Replaces the existing tag set.

- Required: Conditional — at least one of `--labels`, `--tags`, or `--taints` must be provided.
- Constraints: keys and values must be 3–63 characters, alphanumeric with `-`, `_`, `.` allowed, must start and end with alphanumeric. Reserved key prefixes (`vks-`, `vng.vks.`, `vng.vpc.`, `vng.billing.`, `vks-mgmt-`) are not allowed.
- Syntax: `team=platform,cost-center=123`

**`--taints`** (list&lt;string&gt;)

Kubernetes node taints as comma-separated `key=value:effect` entries. Replaces the existing taint set.

- Required: Conditional — at least one of `--labels`, `--tags`, or `--taints` must be provided.
- Constraints: maximum 50 entries.
- Syntax: `dedicated=gpu:NoSchedule`

**`--dry-run`** (boolean)

Print the metadata patch payload without sending the request.

- Required: No
- Default: `false`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

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

Preview the metadata patch (dry run):

```bash
grn vks update-nodegroup-metadata \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --labels env=staging \
  --dry-run
```
