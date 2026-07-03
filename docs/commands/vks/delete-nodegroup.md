# delete-nodegroup

## Description

Delete a specific node group from a VKS cluster. Before executing, the command fetches and displays a preview of the node group that will be removed (name, status, node count).

Unless `--force` is provided, you are prompted to confirm. Use `--dry-run` to see the preview without being prompted and without deleting anything. Use `--force-delete` to instruct the API to perform a forced deletion.

**This action is irreversible.**

## Synopsis

```
grn vks delete-nodegroup
    --cluster-id <value>
    --nodegroup-id <value>
    [--force-delete]
    [--dry-run]
    [--force]
```

## Options

**`--cluster-id`** (string)

ID of the cluster that owns the node group.

- Required: Yes

**`--nodegroup-id`** (string)

ID of the node group to delete.

- Required: Yes

**`--force-delete`** (boolean)

Instruct the API to perform a forced deletion of the node group. Passes `forceDelete=true` as a query parameter to the delete endpoint.

- Required: No
- Default: `false`

**`--dry-run`** (boolean)

Display the node group that would be deleted without sending the delete request.

- Required: No
- Default: `false`

**`--force`** (boolean)

Skip the interactive confirmation prompt and delete immediately. Useful in non-interactive scripts.

- Required: No
- Default: `false`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

Delete a node group interactively (prompts for confirmation):

```bash
grn vks delete-nodegroup \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345
```

Preview what will be deleted without deleting:

```bash
grn vks delete-nodegroup \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --dry-run
```

Delete without confirmation (for use in scripts):

```bash
grn vks delete-nodegroup \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --force
```

Force-delete a stuck node group without confirmation:

```bash
grn vks delete-nodegroup \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --force-delete \
  --force
```
