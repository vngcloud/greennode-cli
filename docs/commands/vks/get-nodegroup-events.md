# get-nodegroup-events

## Description

Get the list of events for a node group. Events can be filtered by action and type, and the results are paginated.

## Synopsis

```
grn vks get-nodegroup-events
    --cluster-id <value>
    --nodegroup-id <value>
    [--action <value>]
    [--type <value>]
    [--page <value>]
    [--page-size <value>]
```

## Options

`--cluster-id` (required)
: The ID of the cluster.

`--nodegroup-id` (required)
: The ID of the node group.

`--action` (optional)
: Filter events by action.

`--type` (optional)
: Filter events by event type.

`--page` (optional)
: Page number to retrieve. Pagination is 0-based (page 0 is the first page). Default: `0`.

`--page-size` (optional)
: Number of events per page. Default: `50`.

## Examples

Get the first page of node group events:

```bash
grn vks get-nodegroup-events \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345
```

Filter events by action and type with a custom page size:

```bash
grn vks get-nodegroup-events \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --action SCALE \
  --type INFO \
  --page 0 \
  --page-size 20
```
