# server

Manage vServer virtual machine instances.

```bash
grn vserver server <command> [options]
```

## Commands

| Command | Description |
|---------|-------------|
| [list](#list) | List all vServer instances |
| [get](#get) | Get details of an instance |
| [create](#create) | Create a new instance |
| [start](#start) | Start a stopped instance |
| [stop](#stop) | Stop a running instance |
| [reboot](#reboot) | Reboot an instance |
| [resize](#resize) | Change an instance to a different flavor |
| [delete](#delete) | Delete an instance |

---

## list

List all vServer instances in your project.

### Synopsis

```
grn vserver server list
    [--page <value>]
    [--page-size <value>]
    [--no-paginate]
    [--name <value>]
```

### Options

`--page` (integer)
: Page number, 1-based. Default: `1`.

`--page-size` (integer)
: Number of items per page. Default: `50`.

`--no-paginate` (boolean)
: Disable pagination and return all results.

`--name` (string)
: Filter by server name (substring match).

### Examples

```bash
# List all servers
grn vserver server list

# Filter by name
grn vserver server list --name web

# Table output
grn vserver server list --output table

# JMESPath filter — names and statuses only
grn vserver server list --query "listData[].{name: name, status: status}"
```

---

## get

Get full details of a vServer instance.

### Synopsis

```
grn vserver server get
    --server-id <value>
```

### Options

`--server-id` (required)
: Server UUID.

### Examples

```bash
grn vserver server get --server-id srv-abc12345-0000-0000-0000-000000000001
```

---

## create

Create a new vServer instance. Several prerequisites are needed — run the discovery commands first if you don't have the IDs yet.

### Synopsis

```
grn vserver server create
    --name <value>
    --flavor-id <value>
    --image-id <value>
    --network-id <value>
    --subnet-id <value>
    --root-disk-type-id <value>
    --zone-id <value>
    [--root-disk-size <value>]
    [--root-disk-encryption-type <value>]
    [--encryption-volume]
    [--data-disk-type-id <value>]
    [--data-disk-size <value>]
    [--data-disk-encryption-type <value>]
    [--data-disk-name <value>]
    [--attach-floating]
    [--security-group <value>]
    [--ssh-key-id <value>]
    [--user-name <value>]
    [--user-password <value>]
    [--expire-password]
    [--server-group-id <value>]
    [--host-group-id <value>]
    [--enable-backup]
    [--backup-instance-point-id <value>]
    [--snapshot-instance-point-id <value>]
    [--period <value>]
    [--is-poc]
    [--is-enable-auto-renew]
    [--os-licence]
    [--user-data <value>]
    [--user-data-base64-encoded]
    [--dry-run]
```

### Options

**Instance settings**

`--name` (required)
: Server name. Must be 5–65 characters, alphanumeric, hyphens, and underscores, starting and ending with an alphanumeric character.

`--flavor-id` (required)
: Flavor (instance type) ID. Run `grn vserver flavor list --family <family> --code <code>` to browse options.

`--image-id` (required)
: OS image ID. Run `grn vserver image list --type os` to browse options.

`--zone-id` (required)
: Availability zone ID. Omit this flag to see available zones printed to stderr.

**Networking**

`--network-id` (required)
: VPC (network) ID. Run `grn vserver vpc list` to browse options.

`--subnet-id` (required)
: Subnet ID within the VPC. Run `grn vserver subnet list --vpc-id <vpc-id>` to browse options.

`--attach-floating` (boolean)
: Attach a floating (public) IP to the server. Default: `false`.

`--security-group` (string)
: Comma-separated list of security group IDs to attach.

**Root disk**

`--root-disk-type-id` (required)
: Volume type ID for the root disk. Run `grn vserver volume-type list --zone-id <zone-id> --type SSD` to browse options.

`--root-disk-size` (integer)
: Root disk size in GiB. Minimum: `20`. Default: `20`.

`--root-disk-encryption-type` (string)
: Encryption type for the root disk.

`--encryption-volume` (boolean)
: Encrypt the root volume.

**Data disk (optional)**

`--data-disk-type-id` (string)
: Volume type ID for an optional data disk.

`--data-disk-size` (integer)
: Data disk size in GiB. Set to `0` to skip the data disk.

`--data-disk-encryption-type` (string)
: Encryption type for the data disk.

`--data-disk-name` (string)
: Name for the data disk.

**Authentication**

`--ssh-key-id` (string)
: SSH key pair ID to inject into the server.

`--user-name` (string)
: OS login username.

`--user-password` (string)
: OS login password.

`--expire-password` (boolean)
: Force a password change on first login. Default: `true`.

**Placement**

`--server-group-id` (string)
: Server group ID for placement affinity/anti-affinity policy.

`--host-group-id` (string)
: Dedicated host group ID.

**Billing**

`--period` (integer)
: Billing period in months. Default: `1`.

`--is-poc` (boolean)
: Mark as a proof-of-concept (PoC) instance.

`--is-enable-auto-renew` (boolean)
: Enable auto-renewal.

`--os-licence` (boolean)
: Include OS licence in billing.

**Backup and restore**

`--enable-backup` (boolean)
: Enable backup for the server.

`--backup-instance-point-id` (string)
: Restore from a backup instance point.

`--snapshot-instance-point-id` (string)
: Restore from a snapshot instance point.

**User data**

`--user-data` (string)
: Cloud-init user data script.

`--user-data-base64-encoded` (boolean)
: Indicate that `--user-data` is already base64-encoded.

**Other**

`--dry-run` (boolean)
: Validate all parameters and print a report without creating the server.

### Examples

```bash
# Minimal server
grn vserver server create \
  --name my-server \
  --flavor-id flv-2c4g \
  --image-id img-ubuntu-22-04 \
  --network-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001 \
  --root-disk-type-id vtype-ssd-123 \
  --zone-id zone-abc123

# Server with floating IP, security group, and SSH key
grn vserver server create \
  --name web-server \
  --flavor-id flv-4c8g \
  --image-id img-ubuntu-22-04 \
  --network-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001 \
  --root-disk-type-id vtype-ssd-123 \
  --root-disk-size 50 \
  --zone-id zone-abc123 \
  --attach-floating \
  --security-group sg-aaa111,sg-bbb222 \
  --ssh-key-id key-abc12345-0000-0000-0000-000000000001

# Dry-run to validate parameters
grn vserver server create \
  --name my-server \
  --flavor-id flv-2c4g \
  --image-id img-ubuntu-22-04 \
  --network-id net-abc12345-0000-0000-0000-000000000001 \
  --subnet-id sub-abc12345-0000-0000-0000-000000000001 \
  --root-disk-type-id vtype-ssd-123 \
  --zone-id zone-abc123 \
  --dry-run
```

---

## start

Start a stopped vServer instance.

### Synopsis

```
grn vserver server start --server-id <value>
```

### Options

`--server-id` (required)
: Server UUID.

### Examples

```bash
grn vserver server start --server-id srv-abc12345-0000-0000-0000-000000000001
```

---

## stop

Stop a running vServer instance.

### Synopsis

```
grn vserver server stop --server-id <value>
```

### Options

`--server-id` (required)
: Server UUID.

### Examples

```bash
grn vserver server stop --server-id srv-abc12345-0000-0000-0000-000000000001
```

---

## reboot

Reboot a vServer instance.

### Synopsis

```
grn vserver server reboot --server-id <value>
```

### Options

`--server-id` (required)
: Server UUID.

### Examples

```bash
grn vserver server reboot --server-id srv-abc12345-0000-0000-0000-000000000001
```

---

## resize

Change a vServer instance to a different flavor (instance type). The server must be stopped before resizing.

### Synopsis

```
grn vserver server resize
    --server-id <value>
    --flavor-id <value>
```

### Options

`--server-id` (required)
: Server UUID.

`--flavor-id` (required)
: New flavor ID. Run `grn vserver flavor list --family <family> --code <code>` to browse options.

### Examples

```bash
grn vserver server resize \
  --server-id srv-abc12345-0000-0000-0000-000000000001 \
  --flavor-id flv-8c16g
```

---

## delete

Delete a vServer instance. Shows a preview and asks for confirmation unless `--force` is used.

### Synopsis

```
grn vserver server delete
    --server-id <value>
    [--delete-all-volumes]
    [--force]
```

### Options

`--server-id` (required)
: Server UUID.

`--delete-all-volumes` (boolean)
: Also delete all volumes attached to the server.

`--force` (boolean)
: Skip the confirmation prompt.

### Examples

```bash
# Interactive confirmation
grn vserver server delete --server-id srv-abc12345-0000-0000-0000-000000000001

# Delete server and all its volumes, no prompt
grn vserver server delete \
  --server-id srv-abc12345-0000-0000-0000-000000000001 \
  --delete-all-volumes \
  --force
```
