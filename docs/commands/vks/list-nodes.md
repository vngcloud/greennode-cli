# list-nodes

## Description

List the nodes belonging to a node group, including their status and details. Results are paginated (0-based page index). Unlike `list-clusters` and `list-nodegroups`, this command does not auto-paginate; pass `--page` to navigate through result pages.

## Synopsis

```
grn vks list-nodes
    --cluster-id <value>
    --nodegroup-id <value>
    [--page <value>]
    [--page-size <value>]
```

## Options

**`--cluster-id`** (string)

ID of the cluster that owns the node group.

- Required: Yes

**`--nodegroup-id`** (string)

ID of the node group whose nodes to list.

- Required: Yes

**`--page`** (integer)

Page number to retrieve. Pagination is 0-based: page `0` is the first page.

- Required: No
- Default: `0`

**`--page-size`** (integer)

Number of nodes to return per page.

- Required: No
- Default: `50`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

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
