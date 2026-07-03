# get-cluster-events

## Description

Retrieve the event log for a VKS cluster. Events record lifecycle actions such as cluster creation, upgrades, and scaling. Results can be filtered by action and event type, and are paginated (0-based page index).

## Synopsis

```
grn vks get-cluster-events
    --cluster-id <value>
    [--action <value>]
    [--type <value>]
    [--page <value>]
    [--page-size <value>]
```

## Options

**`--cluster-id`** (string)

ID of the cluster whose events to retrieve.

- Required: Yes

**`--action`** (string)

Filter events by action name (e.g. `UPGRADE`, `CREATE`, `DELETE`). When omitted, all actions are returned.

- Required: No

**`--type`** (string)

Filter events by event type (e.g. `INFO`, `WARNING`, `ERROR`). When omitted, all types are returned.

- Required: No

**`--page`** (integer)

Page number to retrieve. Pagination is 0-based: page `0` is the first page. Only the specified page is returned; auto-pagination is not performed for this command.

- Required: No
- Default: `0`

**`--page-size`** (integer)

Number of events to return per page.

- Required: No
- Default: `50`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

Get the first page of events for a cluster:

```bash
grn vks get-cluster-events \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345
```

Filter by action and event type with a custom page size:

```bash
grn vks get-cluster-events \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --action UPGRADE \
  --type INFO \
  --page-size 20
```

Retrieve the second page of events:

```bash
grn vks get-cluster-events \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --page 1 \
  --page-size 50
```
