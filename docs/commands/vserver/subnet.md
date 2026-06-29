# subnet

Manage subnets within a VPC.

```bash
grn vserver subnet <command> [options]
```

## Commands

| Command | Description |
|---------|-------------|
| [list](#list) | List all subnets in a VPC |
| [get](#get) | Get details of a subnet |
| [create](#create) | Create a new subnet |
| [delete](#delete) | Delete a subnet |

---

## list

List all subnets within a VPC.

### Synopsis

```
grn vserver subnet list
    --vpc-id <value>
    [--page <value>]
    [--page-size <value>]
```

### Options

`--vpc-id` (required)
: VPC (network) ID.

`--page` (integer)
: Page number, 1-based. Default: `1`.

`--page-size` (integer)
: Number of items per page. Default: `50`.

### Examples

```bash
grn vserver subnet list --vpc-id net-abc12345-0000-0000-0000-000000000001
grn vserver subnet list --vpc-id net-abc12345-0000-0000-0000-000000000001 --output table
```

---

## get

Get details of a subnet.

### Synopsis

```
grn vserver subnet get
    --vpc-id <value>
    --subnet-id <value>
```

### Options

`--vpc-id` (required)
: VPC (network) ID that the subnet belongs to.

`--subnet-id` (required)
: Subnet UUID.

### Examples

```bash
grn vserver subnet get \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001
```

---

## create

Create a new subnet inside a VPC.

### Synopsis

```
grn vserver subnet create
    --vpc-id <value>
    --cidr <value>
    --zone-id <value>
    [--name <value>]
```

### Options

`--vpc-id` (required)
: VPC (network) ID to create the subnet in.

`--cidr` (required)
: CIDR block for the subnet, e.g. `10.0.1.0/24`. Must be within the VPC CIDR range and must not overlap with other subnets.

`--zone-id` (required)
: Availability zone ID. Omit this flag to see available zones printed to stderr.

`--name` (string)
: Subnet name.

### Examples

```bash
grn vserver subnet create \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --cidr 10.0.1.0/24 \
  --zone-id zone-abc123 \
  --name subnet-app
```

---

## delete

Delete a subnet. Shows a confirmation prompt unless `--force` is used.

### Synopsis

```
grn vserver subnet delete
    --vpc-id <value>
    --subnet-id <value>
    [--force]
```

### Options

`--vpc-id` (required)
: VPC (network) ID that the subnet belongs to.

`--subnet-id` (required)
: Subnet UUID.

`--force` (boolean)
: Skip the confirmation prompt.

### Examples

```bash
grn vserver subnet delete \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001

# No prompt
grn vserver subnet delete \
  --vpc-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001 \
  --force
```
