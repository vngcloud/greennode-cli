# generate-kubeconfig

## Description

Request the VKS API to generate (or renew) a kubeconfig for a cluster.

This operation is asynchronous: the server accepts the request and generates the kubeconfig in the background. Once the kubeconfig status becomes `ACTIVE`, run [update-kubeconfig](update-kubeconfig.md) to fetch it and merge it into your local kubeconfig file.

Use `--dry-run` to validate parameters and preview the request without sending it.

## Synopsis

```
grn vks generate-kubeconfig
    --cluster-id <value>
    [--expiration-days <value>]
    [--dry-run]
```

## Options

**`--cluster-id`** (string)

ID of the cluster for which to generate the kubeconfig.

- Required: Yes

**`--expiration-days`** (integer)

Number of days until the generated kubeconfig expires.

- Required: No
- Default: `30`

**`--dry-run`** (boolean)

Validate parameters and print a report without sending the generation request.

- Required: No
- Default: `false`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

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

Preview the request without sending it:

```bash
grn vks generate-kubeconfig \
  --cluster-id cls-abc12345-6789-def0-1234-abcdef012345 \
  --dry-run
```
