# get-cluster

## Description

Get detailed information about a specific VKS cluster, including its status, Kubernetes version, network configuration, node group count, and plugin state.

## Synopsis

```
grn vks get-cluster
    --cluster-id <value>
```

## Options

**`--cluster-id`** (string)

ID of the cluster to retrieve.

- Required: Yes

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

Get cluster details:

```bash
grn vks get-cluster --cluster-id cls-abc12345-6789-def0-1234-abcdef012345
```

Get cluster details as JSON:

```bash
grn vks get-cluster --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 --output json
```

Extract only the cluster status:

```bash
grn vks get-cluster --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 --output json \
  | jq '.status'
```
