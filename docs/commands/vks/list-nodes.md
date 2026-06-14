# list-nodes

## Description

List the nodes belonging to a node group, including their status and details. Results are paginated.

## Synopsis

```
grn vks list-nodes
    --cluster-id <value>
    --nodegroup-id <value>
    [--page <value>]
    [--page-size <value>]
```

## Options

`--cluster-id` (required)
: ID of the cluster that owns the node group.

`--nodegroup-id` (required)
: ID of the node group whose nodes to list.

`--page` (optional, default 0)
: Page number. Pagination is 0-based (page 0 is the first page).

`--page-size` (optional, default 50)
: Number of nodes per page.

## Examples

List nodes in a node group:

```bash
grn vks list-nodes \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345
```

List the second page, 20 nodes per page:

```bash
grn vks list-nodes \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --nodegroup-id ng-abc12345-6789-def0-1234-abcdef012345 \
  --page 1 --page-size 20
```
