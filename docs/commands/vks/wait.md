# wait

## Description

Poll until a VKS cluster or node group reaches a desired state. Exits with code `255` if the resource reaches a terminal failure state (`ERROR`/`FAILED`) or the wait times out. Use after asynchronous operations like create or delete.

## Synopsis

```
grn vks wait cluster-active    --cluster-id <value> [--delay <s>] [--max-attempts <n>]
grn vks wait cluster-deleted   --cluster-id <value> [--delay <s>] [--max-attempts <n>]
grn vks wait nodegroup-active  --cluster-id <value> --nodegroup-id <value> [--delay <s>] [--max-attempts <n>]
grn vks wait nodegroup-deleted --cluster-id <value> --nodegroup-id <value> [--delay <s>] [--max-attempts <n>]
```

## Options

`--cluster-id` (required)
: The ID of the cluster.

`--nodegroup-id` (required for `nodegroup-active` and `nodegroup-deleted`)
: The ID of the node group.

`--delay` (optional)
: Seconds between polls. Default: `30`.

`--max-attempts` (optional)
: Maximum poll attempts before the waiter times out. Default: `40` for `cluster-active`, `cluster-deleted`, and `nodegroup-deleted`; `80` for `nodegroup-active`.

## Exit Codes

| Code | Meaning |
|------|---------|
| `0`   | Resource reached the desired state |
| `255` | Resource reached `ERROR` or `FAILED` status, or waiter timed out |

## Examples

Wait for a cluster to become active:

```bash
grn vks wait cluster-active --cluster-id k8s-xxxxx
```

Wait for a cluster to be fully deleted after `delete-cluster`:

```bash
grn vks delete-cluster --cluster-id k8s-xxxxx --force
grn vks wait cluster-deleted --cluster-id k8s-xxxxx
```

Wait for a node group to be deleted:

```bash
grn vks wait nodegroup-deleted --cluster-id k8s-xxxxx --nodegroup-id ng-xxxxx
```
