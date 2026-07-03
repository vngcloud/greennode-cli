# create-cluster

## Description

Create a new VKS cluster (control plane only). This command provisions the cluster itself; add worker nodes afterwards with [create-nodegroup](create-nodegroup.md).

Cluster names must be 5–20 characters, lowercase alphanumeric and hyphens, starting and ending with an alphanumeric character.

When `--network-type` is `TIGERA` or `CILIUM_OVERLAY`, the `--cidr` option is required. By default, both the load balancer plugin and the block store CSI plugin are enabled; use `--load-balancer-plugin disabled` or `--block-store-csi-plugin disabled` to turn them off.

Use `--dry-run` to validate all parameters without sending a create request.

## Synopsis

```
grn vks create-cluster
    --name <value>
    --k8s-version <value>
    --network-type <value>
    --vpc-id <value>
    [--subnet-id <value>]
    [--cidr <value>]
    [--description <value>]
    [--private-cluster <enabled|disabled>]
    [--release-channel <value>]
    [--load-balancer-plugin <enabled|disabled>]
    [--block-store-csi-plugin <enabled|disabled>]
    [--service-endpoint <enabled|disabled>]
    [--az-strategy <value>]
    [--secondary-subnets <value>]
    [--list-subnet-ids <value>]
    [--node-netmask-size <value>]
    [--auto-upgrade-config <value>]
    [--auto-healing-config <value>]
    [--dry-run]
```

## Options

`--name` (required)
: Cluster name. Must be 5–20 characters, lowercase alphanumeric and hyphens, starting and ending with an alphanumeric character.

`--k8s-version` (required)
: Kubernetes version for the cluster (e.g. `v1.29.1`).

`--network-type` (required)
: Network type for the cluster. Accepted values: `TIGERA`, `CILIUM_OVERLAY`, `CILIUM_NATIVE_ROUTING`.

`--vpc-id` (required)
: VPC ID where the cluster will be provisioned.

`--subnet-id` (optional)
: Subnet ID for the cluster control plane. Optional per the API — provide either `--subnet-id` or `--list-subnet-ids` (or neither); the server validates.

`--cidr` (optional)
: Pod CIDR block. Required when `--network-type` is `TIGERA` or `CILIUM_OVERLAY` (e.g. `10.96.0.0/12`).

`--description` (optional)
: Human-readable description for the cluster.

`--private-cluster` (optional, default `disabled`)
: Control plane accessibility. `enabled` makes the control plane inaccessible from the public internet. Accepted values: `enabled`, `disabled`.

`--release-channel` (optional)
: Release channel for automatic upgrades. Accepted values: `RAPID`, `STABLE`. Default: `STABLE`.

`--load-balancer-plugin` (optional, default `enabled`)
: Load balancer plugin state. Accepted values: `enabled`, `disabled`.

`--block-store-csi-plugin` (optional, default `enabled`)
: Block store CSI plugin state. Accepted values: `enabled`, `disabled`.

`--service-endpoint` (optional, default `disabled`)
: Service endpoint state. Accepted values: `enabled`, `disabled`.

`--az-strategy` (optional, default `SINGLE`)
: Availability-zone strategy for the cluster.

`--secondary-subnets` (optional)
: Comma-separated list of secondary subnet IDs.

`--list-subnet-ids` (optional)
: Comma-separated list of subnet IDs for the cluster.

`--node-netmask-size` (optional)
: Node netmask size (integer). Only sent when explicitly provided.

`--auto-upgrade-config` (optional)
: Auto-upgrade schedule. Shorthand `time=03:00,weekdays=Mon` or JSON `{"weekdays":"Mon,Wed,Fri","time":"03:00"}`. Use JSON when `weekdays` has multiple days (shorthand splits on commas).

`--auto-healing-config` (optional)
: Auto-healing config. Shorthand `enableAutoHealing=true,maxUnhealthy=20%,unhealthyRange=[2-5],timeoutUnhealthy=10` or JSON. `enableAutoHealing` is a boolean, `timeoutUnhealthy` an integer.

`--dry-run` (optional)
: Validate all parameters and print a report without sending the create request.

## Examples

Create a cluster with CILIUM_NATIVE_ROUTING:

```bash
grn vks create-cluster \
  --name my-cluster \
  --k8s-version v1.29.1 \
  --network-type CILIUM_NATIVE_ROUTING \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001
```

Create a cluster with TIGERA network type (CIDR required):

```bash
grn vks create-cluster \
  --name prod-cluster \
  --k8s-version v1.29.1 \
  --network-type TIGERA \
  --cidr 10.96.0.0/12 \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001
```

Then add a node group:

```bash
grn vks create-nodegroup \
  --cluster-id <cluster-id> \
  --name default-ng \
  --flavor-id flv-2c4g \
  --disk-type SSD \
  --ssh-key-id key-abc12345-0000-0000-0000-000000000001
```

Validate parameters without creating (dry run):

```bash
grn vks create-cluster \
  --name my-cluster \
  --k8s-version v1.29.1 \
  --network-type CILIUM_NATIVE_ROUTING \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001 \
  --dry-run
```
