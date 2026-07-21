# wait

## Description

Poll until a VKS cluster or node group reaches a desired state. Exits with code `255` if the resource reaches a terminal failure state (`ERROR` or `FAILED`) or the wait times out. Use after asynchronous operations such as create or delete.

`wait` is a parent command. Invoke one of its four subcommands depending on the resource type and desired state:

| Subcommand | Waits until |
|---|---|
| `cluster-active` | Cluster status is `ACTIVE` |
| `cluster-deleted` | Cluster no longer exists (HTTP 404) |
| `nodegroup-active` | Node group status is `ACTIVE` |
| `nodegroup-deleted` | Node group no longer exists (HTTP 404) |

## cluster-active

### Synopsis

```
grn vks wait cluster-active
    --cluster-id <value>
    [--delay <value>]
    [--max-attempts <value>]
```

### Options

**`--cluster-id`** (string)

ID of the cluster to poll.

- Required: Yes

**`--delay`** (integer)

Seconds to wait between poll attempts.

- Required: No
- Default: `30`

**`--max-attempts`** (integer)

Maximum number of poll attempts before the waiter times out.

- Required: No
- Default: `40`

## cluster-deleted

### Synopsis

```
grn vks wait cluster-deleted
    --cluster-id <value>
    [--delay <value>]
    [--max-attempts <value>]
```

### Options

**`--cluster-id`** (string)

ID of the cluster to poll.

- Required: Yes

**`--delay`** (integer)

Seconds to wait between poll attempts.

- Required: No
- Default: `30`

**`--max-attempts`** (integer)

Maximum number of poll attempts before the waiter times out.

- Required: No
- Default: `40`

## nodegroup-active

### Synopsis

```
grn vks wait nodegroup-active
    --cluster-id <value>
    --nodegroup-id <value>
    [--delay <value>]
    [--max-attempts <value>]
```

### Options

**`--cluster-id`** (string)

ID of the cluster that owns the node group.

- Required: Yes

**`--nodegroup-id`** (string)

ID of the node group to poll.

- Required: Yes

**`--delay`** (integer)

Seconds to wait between poll attempts.

- Required: No
- Default: `30`

**`--max-attempts`** (integer)

Maximum number of poll attempts before the waiter times out.

- Required: No
- Default: `80`

## nodegroup-deleted

### Synopsis

```
grn vks wait nodegroup-deleted
    --cluster-id <value>
    --nodegroup-id <value>
    [--delay <value>]
    [--max-attempts <value>]
```

### Options

**`--cluster-id`** (string)

ID of the cluster that owns the node group.

- Required: Yes

**`--nodegroup-id`** (string)

ID of the node group to poll.

- Required: Yes

**`--delay`** (integer)

Seconds to wait between poll attempts.

- Required: No
- Default: `30`

**`--max-attempts`** (integer)

Maximum number of poll attempts before the waiter times out.

- Required: No
- Default: `40`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

Wait for a cluster to become active after creation:

```bash
grn vks wait cluster-active \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345
```

Wait for a cluster to be fully deleted:

```bash
grn vks delete-cluster --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 --force
grn vks wait cluster-deleted \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345
```

Wait for a node group to become active, polling every 15 seconds:

```bash
grn vks wait nodegroup-active \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --delay 15
```

Wait for a node group to be deleted:

```bash
grn vks wait nodegroup-deleted \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345
```
