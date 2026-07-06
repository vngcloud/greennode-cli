# create-cluster

## Description

Create a new VKS cluster. By default only the control plane is provisioned. Provide `--node-group-name` (together with `--flavor-id`, `--disk-type`, `--ssh-key-id`) to also attach a node group at creation, or add one later with [create-nodegroup](create-nodegroup.md).

Cluster names must be 5–20 characters, lowercase alphanumeric and hyphens, starting and ending with an alphanumeric character. Node group names follow the same pattern with a length of 5–15 characters.

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
    [--node-group-name <value>]
    [--flavor-id <value>]
    [--os <value>]
    [--disk-type <value>]
    [--ssh-key-id <value>]
    [--disk-size <value>]
    [--num-nodes <value>]
    [--private-nodes <enabled|disabled>]
    [--security-groups <value>]
    [--labels <value>]
    [--taints <value>]
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

### Node group settings (optional)

Provide `--node-group-name` to attach a node group at creation (sent as the API's `nodeGroup` object; the deprecated `nodeGroups` array is not used). When set, `--flavor-id`, `--disk-type`, and `--ssh-key-id` are also required. The other node-group flags apply only when a node group is attached.

`--node-group-name` (optional)
: Node group name. Setting this attaches a node group. Must be 5–15 characters, lowercase alphanumeric and hyphens, starting and ending with an alphanumeric character.

`--flavor-id` (required with a node group)
: Flavor (instance type) ID for the nodes.

`--os` (optional, default `ubuntu`)
: Node group OS image. Supported values: `ubuntu`, `linux`, `rocky`.

`--disk-type` (required with a node group)
: Disk type ID for the node boot volumes.

`--ssh-key-id` (required with a node group)
: SSH key pair ID to inject into each node.

`--disk-size` (optional, default `100`)
: Boot disk size in GiB. Accepted range: 20–5000.

`--num-nodes` (optional, default `1`)
: Number of nodes to create. Accepted range: 0–10.

`--private-nodes` (optional, default `disabled`)
: Private nodes state. Accepted values: `enabled`, `disabled`.

`--security-groups` (optional)
: Comma-separated list of security group IDs to attach to the nodes.

`--labels` (optional)
: Comma-separated `key=value` pairs to add as Kubernetes node labels.

`--taints` (optional)
: Comma-separated node taints in `key=value:effect` format.

`--dry-run` (optional)
: Validate all parameters and print a report without sending the create request.

## Examples

Create a cluster (control plane only) with CILIUM_NATIVE_ROUTING:

```bash
grn vks create-cluster \
  --name my-cluster \
  --k8s-version v1.29.1 \
  --network-type CILIUM_NATIVE_ROUTING \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001
```

Create a cluster and attach a node group in one call:

```bash
grn vks create-cluster \
  --name my-cluster \
  --k8s-version v1.29.1 \
  --network-type CILIUM_NATIVE_ROUTING \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001 \
  --node-group-name default-ng \
  --flavor-id flv-2c4g \
  --disk-type SSD \
  --ssh-key-id key-abc12345-0000-0000-0000-000000000001 \
  --num-nodes 3
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
