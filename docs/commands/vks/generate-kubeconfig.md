# generate-kubeconfig

## Description

Request the VKS API to generate (or renew) a kubeconfig for a cluster.

This operation is **asynchronous**: the server accepts the request (HTTP 202) and generates the kubeconfig in the background. Once the kubeconfig becomes `ACTIVE`, run `grn vks update-kubeconfig` to fetch it and merge it into your local kubeconfig file.

## Synopsis

```
grn vks generate-kubeconfig
    --cluster-id <value>
    [--expiration-days <value>]
```

## Options

`--cluster-id` (required)
: The ID of the cluster.

`--expiration-days` (optional)
: Number of days until the generated kubeconfig expires. Default: `30`.

## Examples

Request a kubeconfig with the default 30-day expiration:

```bash
grn vks generate-kubeconfig \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345
```

Request a kubeconfig with a custom expiration, then fetch it once it is active:

```bash
grn vks generate-kubeconfig \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --expiration-days 90

# Once the kubeconfig is ACTIVE:
grn vks update-kubeconfig \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345
```
