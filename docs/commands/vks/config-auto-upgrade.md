# config-auto-upgrade

## Description

Configure the auto-upgrade schedule for a VKS cluster. Sets the days of the week and the time of day when automatic Kubernetes version upgrades will be performed.

Use `--dry-run` to preview the configuration that would be sent without executing the request.

> **Alias:** `set-auto-upgrade-config` is a deprecated alias retained for backward compatibility. Prefer `config-auto-upgrade`.

## Synopsis

```
grn vks config-auto-upgrade
    --cluster-id <value>
    --weekdays <value>
    --time <value>
    [--dry-run]
```

## Options

**`--cluster-id`** (string)

ID of the cluster to configure auto-upgrade for.

- Required: Yes

**`--weekdays`** (string)

Days of the week on which auto-upgrade will run, comma-separated.

- Required: Yes
- Possible values: `Mon`, `Tue`, `Wed`, `Thu`, `Fri`, `Sat`, `Sun`
- Constraints: one or more values from the allowed set, e.g. `Mon,Wed,Fri`

**`--time`** (string)

Time of day at which auto-upgrade will run, in 24-hour `HH:mm` format.

- Required: Yes
- Constraints: `HH:mm`, e.g. `03:00`

**`--dry-run`** (boolean)

Preview the configuration that would be sent without executing the request.

- Required: No
- Default: `false`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

Set auto-upgrade to run on weekdays at 3 AM:

```bash
grn vks config-auto-upgrade \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --weekdays Mon,Tue,Wed,Thu,Fri \
  --time 03:00
```

Set auto-upgrade to run on weekends at midnight:

```bash
grn vks config-auto-upgrade \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --weekdays Sat,Sun \
  --time 00:00
```

Preview the configuration without applying it:

```bash
grn vks config-auto-upgrade \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --weekdays Mon,Wed,Fri \
  --time 02:00 \
  --dry-run
```
