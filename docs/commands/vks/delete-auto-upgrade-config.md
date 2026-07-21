# delete-auto-upgrade-config

## Description

Delete the auto-upgrade configuration for a VKS cluster. Removing this configuration disables scheduled automatic Kubernetes version upgrades for the cluster.

Before executing, the command displays the cluster ID whose auto-upgrade config will be removed. Unless `--force` is provided, you are prompted to confirm. Use `--dry-run` to see the preview without being prompted and without deleting anything.

**This action is irreversible.**

## Synopsis

```
grn vks delete-auto-upgrade-config
    --cluster-id <value>
    [--dry-run]
    [--force]
```

## Options

**`--cluster-id`** (string)

ID of the cluster whose auto-upgrade configuration will be deleted.

- Required: Yes

**`--dry-run`** (boolean)

Display what would be deleted without sending the delete request.

- Required: No
- Default: `false`

**`--force`** (boolean)

Skip the interactive confirmation prompt and delete immediately. Useful in non-interactive scripts.

- Required: No
- Default: `false`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

Delete auto-upgrade config with confirmation:

```bash
grn vks delete-auto-upgrade-config \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345
```

Preview what will be deleted without deleting:

```bash
grn vks delete-auto-upgrade-config \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --dry-run
```

Delete without confirmation (for use in scripts):

```bash
grn vks delete-auto-upgrade-config \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --force
```
