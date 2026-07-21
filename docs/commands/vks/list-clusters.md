# list-clusters

## Description

List all VKS clusters. By default, auto-pagination is enabled and the command fetches all pages, returning the complete result set. Pass `--page` to retrieve a single specific page, or `--no-paginate` to return only page 0 without fetching further pages.

## Synopsis

```
grn vks list-clusters
    [--page <value>]
    [--page-size <value>]
    [--no-paginate]
```

## Options

**`--page`** (integer)

Page number to retrieve (0-based). When provided, disables auto-pagination and returns only the requested page.

- Required: No

**`--page-size`** (integer)

Number of clusters to return per page.

- Required: No
- Default: `50`

**`--no-paginate`** (boolean)

Disable auto-pagination and return only page 0. Equivalent to `--page 0`.

- Required: No
- Default: `false`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

List all clusters (auto-paginated):

```bash
grn vks list-clusters
```

List clusters with a custom page size:

```bash
grn vks list-clusters --page-size 20
```

Fetch page 2 (0-based) only:

```bash
grn vks list-clusters --page 2 --page-size 10
```

Return only the first page without auto-pagination:

```bash
grn vks list-clusters --no-paginate
```

Output as JSON:

```bash
grn vks list-clusters --output json
```
