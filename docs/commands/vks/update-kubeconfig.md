# update-kubeconfig

## Description

Fetch the cluster kubeconfig and merge it into your local kubeconfig file, then (by default) set it as the current context. This is similar to `aws eks update-kubeconfig`.

The kubeconfig must already be `ACTIVE`. If no kubeconfig exists yet (status `NONE`), run `grn vks generate-kubeconfig --cluster-id <value>` first and wait until it becomes active.

The target file is resolved in this order: the `--kubeconfig` flag, then the first entry of `$KUBECONFIG`, then `~/.kube/config`. The merged context is named `vks_<cluster-id>` by default; override it with `--alias`.

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

`--cluster-id` (required)
: The ID of the cluster.

`--kubeconfig` (optional)
: Path to the kubeconfig file to update. Defaults to `$KUBECONFIG` (first entry) or `~/.kube/config`.

`--alias` (optional)
: Context name to use for the merged cluster. Default: `vks_<cluster-id>`.

`--no-set-context` (optional)
: Do not set the merged context as the current context.

`--dry-run` (optional)
: Print what would be written without modifying the kubeconfig file.

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
