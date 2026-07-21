# get-quota

## Description

Get VKS quota limits and current usage for the authenticated user, including the maximum number of clusters, node groups per cluster, nodes per node group, and the current cluster count.

## Synopsis

```
grn vks get-quota
```

## Options

This command takes only the global options.

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

Get quota:

```bash
grn vks get-quota
```

Output:

```json
{
    "maxClusters": 200,
    "maxNodeGroupsPerCluster": 20,
    "maxNodesPerNodeGroup": 10,
    "numClusters": 4
}
```

Extract only the maximum cluster limit:

```bash
grn vks get-quota --query maxClusters
```

Check remaining cluster capacity:

```bash
grn vks get-quota --output json | jq '.maxClusters - .numClusters'
```
