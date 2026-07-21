# list-cluster-versions

## Description

List the Kubernetes versions available for VKS clusters. Use the version strings returned here as the value of `--k8s-version` when running [create-cluster](create-cluster.md) or upgrade commands.

## Synopsis

```
grn vks list-cluster-versions
```

## Options

This command takes only the global options.

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

List available Kubernetes versions:

```bash
grn vks list-cluster-versions
```

Extract only the version strings:

```bash
grn vks list-cluster-versions --output json | jq '.[].version'
```
