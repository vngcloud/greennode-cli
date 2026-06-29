# volume-type

Browse available volume types for a zone. Volume type IDs are used when creating volumes and servers.

```bash
grn vserver volume-type <command> [options]
```

## Commands

| Command | Description |
|---------|-------------|
| [list](#list) | List available volume types for a zone |

---

## list

List volume types available in a specific availability zone. Volume types are grouped by storage class (e.g. SSD, NVMe).

### Synopsis

```
grn vserver volume-type list
    --zone-id <value>
    --type <value>
```

### Options

`--zone-id` (required)
: Availability zone ID. Omit to see available zones printed to stderr.

`--type` (required)
: Volume type zone name (e.g. `SSD`, `NVMe`). Omit to see available type names for the zone printed to stderr.

### Examples

```bash
# See what zones are available (omit --zone-id)
grn vserver volume-type list

# See what types are available in a zone (omit --type)
grn vserver volume-type list --zone-id zone-abc123

# List SSD volume types in a zone
grn vserver volume-type list --zone-id zone-abc123 --type SSD

# List NVMe volume types
grn vserver volume-type list --zone-id zone-abc123 --type NVMe
```
