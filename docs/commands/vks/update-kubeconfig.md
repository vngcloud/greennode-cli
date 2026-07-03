# update-kubeconfig

## Description

Fetch the cluster kubeconfig and merge it into your local kubeconfig file. By default the merged context is set as the current context. This is similar to `aws eks update-kubeconfig`.

The kubeconfig must already be `ACTIVE`. If no kubeconfig exists yet (status `NONE`), run [generate-kubeconfig](generate-kubeconfig.md) first and wait until it becomes active.

The target file is resolved in this order: the `--kubeconfig` flag, then the first entry of `$KUBECONFIG`, then `~/.kube/config`. The merged context is named `vks_<cluster-id>` by default; override it with `--alias`.

Use `--dry-run` to preview what would be written without modifying any file.

## Synopsis

```
grn vks update-kubeconfig
    --cluster-id <value>
    [--kubeconfig <value>]
    [--alias <value>]
    [--no-set-context]
    [--dry-run]
```

## Options

**`--cluster-id`** (string)

ID of the cluster whose kubeconfig to fetch and merge.

- Required: Yes

**`--kubeconfig`** (string)

Path to the kubeconfig file to update.

- Required: No
- Default: first entry of `$KUBECONFIG`, or `~/.kube/config`

**`--alias`** (string)

Context name to use for the merged cluster entry.

- Required: No
- Default: `vks_<cluster-id>`

**`--no-set-context`** (boolean)

Do not set the merged context as the current context after merging.

- Required: No
- Default: `false`

**`--dry-run`** (boolean)

Print what would be written without modifying the kubeconfig file.

- Required: No
- Default: `false`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

Merge the cluster kubeconfig into the default file and set it as the current context:

```bash
grn vks update-kubeconfig \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345
```

This creates a context named `vks_cls-abc12345-6789-def0-1234-abcdef012345`.

Use a custom context name and a specific kubeconfig file, without switching the current context:

```bash
grn vks update-kubeconfig \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --kubeconfig ./my-kubeconfig.yaml \
  --alias prod-cluster \
  --no-set-context
```

Preview the changes without writing the file:

```bash
grn vks update-kubeconfig \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --dry-run
```

If the kubeconfig does not exist yet, generate it first:

```bash
grn vks generate-kubeconfig --cluster-id cls-abc12345-6789-def0-1234-abcdef012345
# Wait until the kubeconfig becomes ACTIVE, then:
grn vks update-kubeconfig --cluster-id cls-abc12345-6789-def0-1234-abcdef012345
```
