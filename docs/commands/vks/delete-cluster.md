# delete-cluster

## Description

Delete a VKS cluster and all of its associated node groups. Before executing, the command fetches and displays a preview showing the cluster name, status, version, node count, and the list of node groups that will be removed.

Unless `--force` is provided, you are prompted to confirm. Use `--dry-run` to see the preview without being prompted and without deleting anything.

**This action is irreversible.**

## Synopsis

```
grn vks delete-cluster
    --cluster-id <value>
    [--dry-run]
    [--force]
```

## Options

**`--cluster-id`** (string)

ID of the cluster to delete.

- Required: Yes

**`--dry-run`** (boolean)

Display the resources that would be deleted without sending the delete request.

- Required: No
- Default: `false`

**`--force`** (boolean)

Skip the interactive confirmation prompt and delete immediately. Useful in non-interactive scripts.

- Required: No
- Default: `false`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

Delete a cluster interactively (prompts for confirmation):

```bash
grn vks delete-cluster --cluster-id cls-abc12345-6789-def0-1234-abcdef012345
```

Preview what will be deleted without deleting:

```bash
grn vks delete-cluster \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --dry-run
```

Delete without confirmation (for use in scripts):

```bash
grn vks delete-cluster \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --force
```
