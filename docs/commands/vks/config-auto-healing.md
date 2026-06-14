# config-auto-healing

## Description

Configure auto-healing for a VKS cluster. Auto-healing automatically replaces unhealthy nodes to keep the cluster in a working state.

## Synopsis

```
grn vks config-auto-healing
    --cluster-id <value>
    --enable-auto-healing
    [--max-unhealthy <value>]
    [--unhealthy-range <value>]
    [--timeout-unhealthy <value>]
```

## Options

`--cluster-id` (required)
: The ID of the cluster.

`--enable-auto-healing` (required)
: Whether to enable auto-healing. Pass `--enable-auto-healing` to enable, or `--enable-auto-healing=false` to disable.

`--max-unhealthy` (optional)
: Maximum number (or percentage) of unhealthy nodes tolerated, e.g. `30%`.

`--unhealthy-range` (optional)
: The unhealthy range threshold.

`--timeout-unhealthy` (optional)
: Time in seconds a node may stay unhealthy before it is replaced.

## Examples

Enable auto-healing with default thresholds:

```bash
grn vks config-auto-healing \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --enable-auto-healing
```

Enable auto-healing with custom thresholds:

```bash
grn vks config-auto-healing \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --enable-auto-healing \
  --max-unhealthy 30% \
  --timeout-unhealthy 300
```

Disable auto-healing:

```bash
grn vks config-auto-healing \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --enable-auto-healing=false
```
