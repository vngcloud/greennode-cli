# create-cluster

## Description

Create a new VKS cluster. By default only the control plane is provisioned; add worker nodes afterwards with [create-nodegroup](create-nodegroup.md).

When `--network-type` is `TIGERA` or `CILIUM_OVERLAY`, `--cidr` is required. The load balancer and block store CSI plugins are enabled by default.

Use `--dry-run` to validate parameters without sending a create request.

## Synopsis

```
grn vks create-cluster
    --name <value>
    --k8s-version <value>
    --network-type <value>
    --vpc-id <value>
    [--subnet-id <value>]
    [--list-subnet-ids <value>]
    [--cidr <value>]
    [--description <value>]
    [--private-cluster <enabled|disabled>]
    [--release-channel <value>]
    [--load-balancer-plugin <enabled|disabled>]
    [--block-store-csi-plugin <enabled|disabled>]
    [--service-endpoint <enabled|disabled>]
    [--az-strategy <value>]
    [--secondary-subnets <value>]
    [--node-netmask-size <value>]
    [--auto-upgrade-config <value>]
    [--auto-healing-config <value>]
    [--dry-run]
```

## Options

**`--name`** (string)

Name of the cluster.

- Required: Yes
- Constraints: 5–20 characters; lowercase letters, digits, and hyphens; must start and end with a letter or digit.

**`--k8s-version`** (string)

Kubernetes version for the cluster (e.g. `v1.29.13-vks.1740045600`).

- Required: Yes
- See available versions with [list-cluster-versions](list-cluster-versions.md).

**`--network-type`** (string)

Cluster network plugin.

- Required: Yes
- Possible values: `TIGERA`, `CILIUM_OVERLAY`, `CILIUM_NATIVE_ROUTING`

**`--vpc-id`** (string)

ID of the VPC to provision the cluster in.

- Required: Yes

**`--subnet-id`** (string)

Subnet ID for the cluster control plane.

- Required: No
- Provide `--subnet-id`, `--list-subnet-ids`, or neither — the server validates the combination.

**`--list-subnet-ids`** (list&lt;string&gt;)

Subnet IDs for the cluster, comma-separated.

- Required: No
- Syntax: `sub-aaa,sub-bbb`

**`--cidr`** (string)

Pod CIDR block (e.g. `10.96.0.0/12`).

- Required: Conditional — required when `--network-type` is `TIGERA` or `CILIUM_OVERLAY`.

**`--description`** (string)

Human-readable description for the cluster.

- Required: No

**`--private-cluster`** (string)

Control-plane accessibility. `enabled` makes the control plane unreachable from the public internet.

- Required: No
- Default: `disabled`
- Possible values: `enabled`, `disabled`

**`--release-channel`** (string)

Release channel for automatic upgrades.

- Required: No
- Default: `STABLE`
- Possible values: `RAPID`, `STABLE`

**`--load-balancer-plugin`** (string)

Load balancer plugin state.

- Required: No
- Default: `enabled`
- Possible values: `enabled`, `disabled`

**`--block-store-csi-plugin`** (string)

Block store CSI plugin state.

- Required: No
- Default: `enabled`
- Possible values: `enabled`, `disabled`

**`--service-endpoint`** (string)

Service endpoint state.

- Required: No
- Default: `disabled`
- Possible values: `enabled`, `disabled`

**`--az-strategy`** (string)

Availability-zone strategy for the cluster.

- Required: No
- Default: `SINGLE`

**`--secondary-subnets`** (list&lt;string&gt;)

Secondary subnet IDs, comma-separated.

- Required: No
- Syntax: `sub-aaa,sub-bbb`

**`--node-netmask-size`** (integer)

Node netmask size. Only sent when explicitly provided.

- Required: No

**`--auto-upgrade-config`** (structure)

Auto-upgrade schedule. Accepts shorthand or JSON.

- Required: No
- Members:
    - `weekdays` (string) — days to run auto-upgrade, e.g. `Mon,Wed,Fri`
    - `time` (string) — time of day, 24-hour `HH:mm`, e.g. `03:00`

Shorthand syntax (use JSON when `weekdays` has multiple days, since shorthand splits on commas):

```
time=03:00,weekdays=Mon
```

JSON syntax:

```json
{"weekdays": "Mon,Wed,Fri", "time": "03:00"}
```

**`--auto-healing-config`** (structure)

Auto-healing configuration. Accepts shorthand or JSON.

- Required: No
- Members:
    - `enableAutoHealing` (boolean) — enable or disable auto-healing
    - `maxUnhealthy` (string) — maximum unhealthy nodes, e.g. `20%`
    - `unhealthyRange` (string) — unhealthy node count range, e.g. `[2-5]`
    - `timeoutUnhealthy` (integer) — unhealthy timeout in seconds

Shorthand syntax:

```
enableAutoHealing=true,maxUnhealthy=20%,unhealthyRange=[2-5],timeoutUnhealthy=10
```

JSON syntax:

```json
{"enableAutoHealing": true, "maxUnhealthy": "20%", "unhealthyRange": "[2-5]", "timeoutUnhealthy": 10}
```

**`--dry-run`** (boolean)

Validate parameters and print a report without sending the create request.

- Required: No
- Default: `false`

## Global options

This command also accepts the global options (`--profile`, `--region`, `--output`, `--query`, `--endpoint-url`, `--debug`, …).

## Examples

Create a cluster (control plane only) with CILIUM_NATIVE_ROUTING:

```bash
grn vks create-cluster \
  --name my-cluster \
  --k8s-version v1.29.13-vks.1740045600 \
  --network-type CILIUM_NATIVE_ROUTING \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001
```

Create a cluster with TIGERA (CIDR required) and auto-healing:

```bash
grn vks create-cluster \
  --name prod-cluster \
  --k8s-version v1.29.13-vks.1740045600 \
  --network-type TIGERA \
  --cidr 10.96.0.0/12 \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001 \
  --auto-healing-config 'enableAutoHealing=true,maxUnhealthy=20%,timeoutUnhealthy=10'
```

Validate parameters without creating (dry run):

```bash
grn vks create-cluster \
  --name my-cluster \
  --k8s-version v1.29.13-vks.1740045600 \
  --network-type CILIUM_NATIVE_ROUTING \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --dry-run
```
