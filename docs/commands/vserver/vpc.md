# vpc

Manage VPC (Virtual Private Cloud) networks.

```bash
grn vserver vpc <command> [options]
```

## Commands

| Command | Description |
|---------|-------------|
| [list](#list) | List all VPCs |
| [get](#get) | Get details of a VPC |
| [create](#create) | Create a new VPC |
| [delete](#delete) | Delete a VPC |

---

## list

List all VPCs in your project.

### Synopsis

```
grn vserver vpc list
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
: Filter by VPC name (substring match).

### Examples

```bash
grn vserver vpc list
grn vserver vpc list --name prod --output table
```

---

## get

Get details of a VPC.

### Synopsis

```
grn vserver vpc get --vpc-id <value>
```

### Options

`--vpc-id` (required)
: VPC (network) ID.

### Examples

```bash
grn vserver vpc get --vpc-id net-abc12345-0000-0000-0000-000000000001
```

---

## create

Create a new VPC network.

### Synopsis

```
grn vserver vpc create
    --name <value>
    --cidr <value>
    [--description <value>]
    [--is-default]
    [--dry-run]
```

### Options

`--name` (required)
: VPC name.

`--cidr` (required)
: CIDR block for the VPC, e.g. `10.0.0.0/16`. The CIDR must not overlap with other VPCs in the same project.

`--description` (string)
: VPC description.

`--is-default` (boolean)
: Mark this VPC as the default network.

`--dry-run` (boolean)
: Validate parameters without creating the VPC.

### Examples

```bash
grn vserver vpc create --name prod-vpc --cidr 10.0.0.0/16

grn vserver vpc create \
  --name staging-vpc \
  --cidr 10.1.0.0/16 \
  --description "Staging environment VPC"

# Validate first
grn vserver vpc create --name prod-vpc --cidr 10.0.0.0/16 --dry-run
```

---

## delete

Delete a VPC. Shows a confirmation prompt unless `--force` is used.

### Synopsis

```
grn vserver vpc delete
    --vpc-id <value>
    [--force]
```

### Options

`--vpc-id` (required)
: VPC (network) ID.

`--force` (boolean)
: Skip the confirmation prompt.

### Examples

```bash
# Interactive confirmation
grn vserver vpc delete --vpc-id net-abc12345-0000-0000-0000-000000000001

# No prompt
grn vserver vpc delete --vpc-id net-abc12345-0000-0000-0000-000000000001 --force
```
