# secgroup

Manage security groups and their inbound/outbound rules.

```bash
grn vserver secgroup <command> [options]
grn vserver secgroup rule <command> [options]
```

## Commands

| Command | Description |
|---------|-------------|
| [list](#list) | List all security groups |
| [get](#get) | Get details of a security group |
| [create](#create) | Create a new security group |
| [delete](#delete) | Delete a security group |
| [rule list](#rule-list) | List rules in a security group |
| [rule get](#rule-get) | Get details of a rule |
| [rule create](#rule-create) | Add a rule to a security group |
| [rule delete](#rule-delete) | Delete a rule from a security group |

---

## list

List all security groups in your project.

### Synopsis

```
grn vserver secgroup list
    [--page <value>]
    [--page-size <value>]
    [--name <value>]
```

### Options

`--page` (integer)
: Page number, 1-based. Default: `1`.

`--page-size` (integer)
: Number of items per page. Default: `50`.

`--name` (string)
: Filter by security group name (substring match).

### Examples

```bash
grn vserver secgroup list
grn vserver secgroup list --name web --output table
```

---

## get

Get details of a security group.

### Synopsis

```
grn vserver secgroup get --secgroup-id <value>
```

### Options

`--secgroup-id` (required)
: Security group ID.

### Examples

```bash
grn vserver secgroup get --secgroup-id sg-abc12345-0000-0000-0000-000000000001
```

---

## create

Create a new security group.

### Synopsis

```
grn vserver secgroup create
    --name <value>
    [--description <value>]
```

### Options

`--name` (required)
: Security group name.

`--description` (string)
: Security group description.

### Examples

```bash
grn vserver secgroup create --name web-sg --description "Allow HTTP and HTTPS traffic"
```

---

## delete

Delete a security group. Shows a confirmation prompt unless `--force` is used.

### Synopsis

```
grn vserver secgroup delete
    --secgroup-id <value>
    [--force]
```

### Options

`--secgroup-id` (required)
: Security group ID.

`--force` (boolean)
: Skip the confirmation prompt.

### Examples

```bash
grn vserver secgroup delete --secgroup-id sg-abc12345-0000-0000-0000-000000000001
grn vserver secgroup delete --secgroup-id sg-abc12345-0000-0000-0000-000000000001 --force
```

---

## rule list

List all rules in a security group.

### Synopsis

```
grn vserver secgroup rule list
    --secgroup-id <value>
    [--page <value>]
    [--page-size <value>]
```

### Options

`--secgroup-id` (required)
: Security group ID.

`--page` (integer)
: Page number, 1-based. Default: `1`.

`--page-size` (integer)
: Number of items per page. Default: `50`.

### Examples

```bash
grn vserver secgroup rule list --secgroup-id sg-abc12345-0000-0000-0000-000000000001
```

---

## rule get

Get details of a security group rule.

### Synopsis

```
grn vserver secgroup rule get
    --secgroup-id <value>
    --rule-id <value>
```

### Options

`--secgroup-id` (required)
: Security group ID.

`--rule-id` (required)
: Security group rule ID.

### Examples

```bash
grn vserver secgroup rule get \
  --secgroup-id sg-abc12345-0000-0000-0000-000000000001 \
  --rule-id rule-abc12345-0000-0000-0000-000000000001
```

---

## rule create

Add a new inbound or outbound rule to a security group.

### Synopsis

```
grn vserver secgroup rule create
    --secgroup-id <value>
    --direction <value>
    --protocol <value>
    --port-range-min <value>
    --port-range-max <value>
    --ether-type <value>
    --remote-ip-prefix <value>
    [--remote-group-id <value>]
    [--description <value>]
```

### Options

`--secgroup-id` (required)
: Security group ID to add the rule to.

`--direction` (required)
: Traffic direction. Accepted values: `ingress`, `egress`.

`--protocol` (required)
: Network protocol. Accepted values: `tcp`, `udp`, `icmp`, `any`.

`--port-range-min` (required for tcp/udp)
: Minimum port number. Not valid for `icmp` or `any`.

`--port-range-max` (required for tcp/udp)
: Maximum port number. Must be ≥ `--port-range-min`. Not valid for `icmp` or `any`.

`--ether-type` (required)
: IP version. Accepted values: `IPv4`, `IPv6`. Default: `IPv4`.

`--remote-ip-prefix` (required)
: Remote CIDR block, e.g. `0.0.0.0/0` (all IPv4) or `192.168.1.0/24`.

`--remote-group-id` (string)
: Remote security group ID. Use instead of `--remote-ip-prefix` to allow traffic from another security group.

`--description` (string)
: Rule description.

### Examples

```bash
# Allow all inbound HTTP traffic
grn vserver secgroup rule create \
  --secgroup-id sg-abc12345-0000-0000-0000-000000000001 \
  --direction ingress \
  --protocol tcp \
  --port-range-min 80 \
  --port-range-max 80 \
  --ether-type IPv4 \
  --remote-ip-prefix 0.0.0.0/0

# Allow HTTPS from a specific CIDR
grn vserver secgroup rule create \
  --secgroup-id sg-abc12345-0000-0000-0000-000000000001 \
  --direction ingress \
  --protocol tcp \
  --port-range-min 443 \
  --port-range-max 443 \
  --ether-type IPv4 \
  --remote-ip-prefix 203.0.113.0/24

# Allow all ICMP (ping) inbound
grn vserver secgroup rule create \
  --secgroup-id sg-abc12345-0000-0000-0000-000000000001 \
  --direction ingress \
  --protocol icmp \
  --ether-type IPv4 \
  --remote-ip-prefix 0.0.0.0/0

# Allow all outbound traffic
grn vserver secgroup rule create \
  --secgroup-id sg-abc12345-0000-0000-0000-000000000001 \
  --direction egress \
  --protocol any \
  --ether-type IPv4 \
  --remote-ip-prefix 0.0.0.0/0
```

---

## rule delete

Delete a rule from a security group.

### Synopsis

```
grn vserver secgroup rule delete
    --secgroup-id <value>
    --rule-id <value>
```

### Options

`--secgroup-id` (required)
: Security group ID.

`--rule-id` (required)
: Security group rule ID.

### Examples

```bash
grn vserver secgroup rule delete \
  --secgroup-id sg-abc12345-0000-0000-0000-000000000001 \
  --rule-id rule-abc12345-0000-0000-0000-000000000001
```
