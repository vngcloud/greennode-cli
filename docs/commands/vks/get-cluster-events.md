# get-cluster-events

## Description

Get the list of events for a VKS cluster. Events can be filtered by action and type, and the results are paginated.

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

`--cluster-id` (required)
: The ID of the cluster.

`--action` (optional)
: Filter events by action.

`--type` (optional)
: Filter events by event type.

`--page` (optional)
: Page number to retrieve. Pagination is 0-based (page 0 is the first page). Default: `0`.

`--page-size` (optional)
: Number of events per page. Default: `50`.

## Examples

Get the first page of cluster events:

```bash
grn vks get-cluster-events \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345
```

Filter events by action and type with a custom page size:

```bash
grn vks get-cluster-events \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --action UPGRADE \
  --type INFO \
  --page 0 \
  --page-size 20
```
