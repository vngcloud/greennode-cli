# volume

Manage block storage volumes.

```bash
grn vserver volume <command> [options]
```

## Commands

| Command | Description |
|---------|-------------|
| [list](#list) | List all volumes |
| [get](#get) | Get details of a volume |
| [create](#create) | Create a new volume |
| [resize](#resize) | Resize or change the type of a volume |
| [delete](#delete) | Delete a volume |

---

## list

List all volumes in your project.

### Synopsis

```
grn vserver volume list
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
: Filter by volume name (substring match).

### Examples

```bash
grn vserver volume list
grn vserver volume list --name data --output table
```

---

## get

Get details of a volume.

### Synopsis

```
grn vserver volume get --volume-id <value>
```

### Options

`--volume-id` (required)
: Volume UUID.

### Examples

```bash
grn vserver volume get --volume-id vol-abc12345-0000-0000-0000-000000000001
```

---

## create

Create a new block storage volume.

### Synopsis

```
grn vserver volume create
    --name <value>
    --volume-type-id <value>
    --size <value>
    [--zone-id <value>]
    [--description <value>]
    [--encryption-type <value>]
    [--multiattach]
    [--is-poc]
    [--dry-run]
```

### Options

`--name` (required)
: Volume name.

`--volume-type-id` (required)
: Volume type ID. Run `grn vserver volume-type list --zone-id <zone-id> --type SSD` to see options.

`--size` (required)
: Volume size in GiB.

`--zone-id` (string)
: Availability zone ID. Omit to see available zones.

`--description` (string)
: Volume description.

`--encryption-type` (string)
: Encryption type for the volume.

`--multiattach` (boolean)
: Allow the volume to be attached to multiple servers simultaneously.

`--is-poc` (boolean)
: Mark as a proof-of-concept volume.

`--dry-run` (boolean)
: Validate parameters without creating the volume.

### Examples

```bash
grn vserver volume create \
  --name data-vol \
  --volume-type-id vtype-ssd-123 \
  --size 100 \
  --zone-id zone-abc123

# Validate first
grn vserver volume create \
  --name data-vol \
  --volume-type-id vtype-ssd-123 \
  --size 100 \
  --dry-run
```

---

## resize

Resize a volume's size in GiB, change its volume type, or both. Prints the current state and planned change before applying.

At least one of `--size` or `--volume-type-id` must be provided.

### Synopsis

```
grn vserver volume resize
    --volume-id <value>
    [--size <value>]
    [--volume-type-id <value>]
    [--dry-run]
```

### Options

`--volume-id` (required)
: Volume UUID.

`--size` (integer)
: New volume size in GiB. Must be equal to or greater than the current size (volumes cannot be shrunk).

`--volume-type-id` (string)
: New volume type ID. If omitted, the current volume type is preserved.

`--dry-run` (boolean)
: Validate parameters without sending the resize request.

### Examples

```bash
# Expand volume to 200 GiB
grn vserver volume resize \
  --volume-id vol-abc12345-0000-0000-0000-000000000001 \
  --size 200

# Change volume type
grn vserver volume resize \
  --volume-id vol-abc12345-0000-0000-0000-000000000001 \
  --volume-type-id vtype-nvme-456

# Resize and change type at the same time
grn vserver volume resize \
  --volume-id vol-abc12345-0000-0000-0000-000000000001 \
  --size 200 \
  --volume-type-id vtype-nvme-456
```

---

## delete

Delete a volume. Shows a confirmation prompt unless `--force` is used.

### Synopsis

```
grn vserver volume delete
    --volume-id <value>
    [--force]
```

### Options

`--volume-id` (required)
: Volume UUID.

`--force` (boolean)
: Skip the confirmation prompt.

### Examples

```bash
grn vserver volume delete --volume-id vol-abc12345-0000-0000-0000-000000000001

# No prompt
grn vserver volume delete \
  --volume-id vol-abc12345-0000-0000-0000-000000000001 \
  --force
```
