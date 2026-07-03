# config-auto-healing

## Description

Configure auto-healing for a VKS cluster. Auto-healing automatically replaces unhealthy nodes to keep the cluster in a working state.

Only `--enable-auto-healing` is sent unconditionally. The optional threshold flags (`--max-unhealthy`, `--unhealthy-range`, `--timeout-unhealthy`) are only included in the request when explicitly provided on the command line.

Use `--dry-run` to preview the configuration that would be sent without executing the request.

## Synopsis

```
grn vks config-auto-healing
    --cluster-id <value>
    --enable-auto-healing <value>
    [--max-unhealthy <value>]
    [--unhealthy-range <value>]
    [--timeout-unhealthy <value>]
    [--dry-run]
```

## Options

**`--cluster-id`** (string)

ID of the cluster to configure auto-healing for.

- Required: Yes

**`--enable-auto-healing`** (boolean)

Enable or disable auto-healing. Pass `true` to enable or `false` to disable.

- Required: Yes
- Possible values: `true`, `false`

**`--max-unhealthy`** (string)

Maximum proportion of unhealthy nodes tolerated before auto-healing is triggered. Accepts a percentage string.

- Required: No
- Constraints: percentage string, e.g. `30%`

**`--unhealthy-range`** (string)

Unhealthy node count range. When the number of unhealthy nodes falls within this range, auto-healing is triggered.

- Required: No
- Constraints: bracket-enclosed range string, e.g. `[2-5]`

**`--timeout-unhealthy`** (integer)

Time in seconds that a node may remain unhealthy before it is replaced.

- Required: No
- Default: `0`

**`--dry-run`** (boolean)

Preview the configuration that would be sent without executing the request.

- Required: No
- Default: `false`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

Enable auto-healing with default thresholds:

```bash
grn vks config-auto-healing \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --enable-auto-healing true
```

Enable auto-healing with custom thresholds:

```bash
grn vks config-auto-healing \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --enable-auto-healing true \
  --max-unhealthy 30% \
  --unhealthy-range '[2-5]' \
  --timeout-unhealthy 300
```

Disable auto-healing:

```bash
grn vks config-auto-healing \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --enable-auto-healing false
```

Preview the configuration without applying it:

```bash
grn vks config-auto-healing \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --enable-auto-healing true \
  --max-unhealthy 20% \
  --dry-run
```
