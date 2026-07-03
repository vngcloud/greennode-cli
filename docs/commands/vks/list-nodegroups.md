# list-nodegroups

## Description

List all node groups belonging to a cluster. By default, auto-pagination is enabled and the command fetches all pages, returning the complete result set. Pass `--page` to retrieve a single specific page, or `--no-paginate` to return only page 0.

## Synopsis

```
grn vks list-nodegroups
    --cluster-id <value>
    [--page <value>]
    [--page-size <value>]
    [--no-paginate]
```

## Options

**`--cluster-id`** (string)

ID of the cluster whose node groups to list.

- Required: Yes

**`--page`** (integer)

Page number to retrieve (0-based). When provided, disables auto-pagination and returns only the requested page.

- Required: No

**`--page-size`** (integer)

Number of node groups to return per page.

- Required: No
- Default: `50`

**`--no-paginate`** (boolean)

Disable auto-pagination and return only page 0. Equivalent to `--page 0`.

- Required: No
- Default: `false`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

List all node groups for a cluster:

```bash
grn vks list-nodegroups --cluster-id cls-abc12345-6789-def0-1234-abcdef012345
```

List node groups with a custom page size:

```bash
grn vks list-nodegroups \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --page-size 20
```

Fetch only page 0:

```bash
grn vks list-nodegroups \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --no-paginate
```

Output as JSON and count node groups:

```bash
grn vks list-nodegroups \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --output json | jq '.items | length'
```
