# create-cluster

## Description

Create a new VKS cluster with an initial default node group. The command provisions both the control plane and the first node group in a single call.

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
    --subnet-id <value>
    --node-group-name <value>
    --flavor-id <value>
    [--os <value>]
    --disk-type <value>
    --ssh-key-id <value>
    [--cidr <value>]
    [--description <value>]
    [--private-cluster <enabled|disabled>]
    [--release-channel <value>]
    [--load-balancer-plugin <enabled|disabled>]
    [--block-store-csi-plugin <enabled|disabled>]
    [--disk-size <value>]
    [--num-nodes <value>]
    [--private-nodes <enabled|disabled>]
    [--security-groups <value>]
    [--labels <value>]
    [--taints <value>]
    [--dry-run]
```

## Options

**Cluster settings**

`--name` (required)
: Cluster name. Must be 5–20 characters, lowercase alphanumeric and hyphens, starting and ending with an alphanumeric character.

`--k8s-version` (required)
: Kubernetes version for the cluster (e.g. `v1.29.1`).

`--network-type` (required)
: Network type for the cluster. Accepted values: `TIGERA`, `CILIUM_OVERLAY`, `CILIUM_NATIVE_ROUTING`.

`--vpc-id` (required)
: VPC ID where the cluster will be provisioned.

`--subnet-id` (required)
: Subnet ID for the cluster control plane and the default node group.

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

**Node group settings**

`--node-group-name` (required)
: Name of the initial node group. Must be 5–15 characters, lowercase alphanumeric and hyphens, starting and ending with an alphanumeric character.

`--flavor-id` (required)
: Flavor (instance type) ID for the nodes.

`--os` (optional, default `ubuntu`)
: Node group OS image. Supported values: `ubuntu`, `linux`, `rocky`.

`--disk-type` (required)
: Disk type ID for the node boot volumes.

`--ssh-key-id` (required)
: SSH key pair ID to inject into each node.

`--disk-size` (optional)
: Boot disk size in GiB. Accepted range: 20–5000. Default: `100`.

`--num-nodes` (optional)
: Number of nodes to create in the default node group. Accepted range: 0–10. Default: `1`.

`--private-nodes` (optional, default `disabled`)
: Private nodes state. `enabled` means nodes will not have public IP addresses. Accepted values: `enabled`, `disabled`.

`--security-groups` (optional)
: Comma-separated list of security group IDs to attach to the nodes (e.g. `sg-aaa111,sg-bbb222`).

`--labels` (optional)
: Comma-separated `key=value` pairs to add as Kubernetes node labels (e.g. `env=prod,tier=app`).

`--taints` (optional)
: Comma-separated node taints in `key=value:effect` format (e.g. `dedicated=gpu:NoSchedule`).

`--dry-run` (optional)
: Validate all parameters and print a report without sending the create request.

## Examples

Create a minimal cluster with CILIUM_NATIVE_ROUTING:

```bash
grn vks create-cluster \
  --name my-cluster \
  --k8s-version v1.29.1 \
  --network-type CILIUM_NATIVE_ROUTING \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001 \
  --node-group-name default-ng \
  --flavor-id flv-2c4g \
  --os ubuntu \
  --disk-type SSD \
  --ssh-key-id key-abc12345-0000-0000-0000-000000000001
```

Create a cluster with TIGERA network type (CIDR required):

```bash
grn vks create-cluster \
  --name prod-cluster \
  --k8s-version v1.29.1 \
  --network-type TIGERA \
  --cidr 10.96.0.0/12 \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001 \
  --node-group-name prod-ng \
  --flavor-id flv-4c8g \
  --os ubuntu \
  --disk-type SSD \
  --disk-size 200 \
  --num-nodes 3 \
  --ssh-key-id key-abc12345-0000-0000-0000-000000000001 \
  --labels env=prod,tier=app \
  --taints dedicated=gpu:NoSchedule
```

Validate parameters without creating (dry run):

```bash
grn vks create-cluster \
  --name my-cluster \
  --k8s-version v1.29.1 \
  --network-type CILIUM_NATIVE_ROUTING \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001 \
  --node-group-name default-ng \
  --flavor-id flv-2c4g \
  --os ubuntu \
  --disk-type SSD \
  --ssh-key-id key-abc12345-0000-0000-0000-000000000001 \
  --dry-run
```
