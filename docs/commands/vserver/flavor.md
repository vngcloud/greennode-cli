# flavor

Browse available vServer instance flavors (CPU/memory/GPU combinations).

```bash
grn vserver flavor <command> [options]
```

## Commands

| Command | Description |
|---------|-------------|
| [list-families](#list-families) | List available instance families |
| [list-codes](#list-codes) | List available CPU platform codes |
| [list](#list) | List flavors for a family and CPU platform |

## Discovery workflow

Flavors are organized by **family** (e.g. `general-purpose`, `memory-optimized`) and **CPU platform code** (e.g. `Intel`, `AMD`). To find a flavor:

```bash
# Step 1: find available families
grn vserver flavor list-families

# Step 2: find CPU platform codes
grn vserver flavor list-codes

# Step 3: list flavors for the chosen family and code
grn vserver flavor list --family <family> --code <code>
```

---

## list-families

List available instance families.

### Synopsis

```
grn vserver flavor list-families
```

### Examples

```bash
grn vserver flavor list-families
grn vserver flavor list-families --output table
```

---

## list-codes

List available CPU platform codes.

### Synopsis

```
grn vserver flavor list-codes
```

### Examples

```bash
grn vserver flavor list-codes
```

---

## list

List available flavors for a specific instance family and CPU platform code. Only flavors with remaining capacity are shown.

### Synopsis

```
grn vserver flavor list
    --family <value>
    --code <value>
    [--zone-id <value>]
    [--page <value>]
    [--page-size <value>]
```

### Options

`--family` (required)
: Instance family name. Run `grn vserver flavor list-families` to see options. Omit to see available families printed to stderr.

`--code` (required)
: CPU platform code. Run `grn vserver flavor list-codes` to see options. Omit to see available codes printed to stderr.

`--zone-id` (string)
: Filter flavors by availability zone.

`--page` (integer)
: Page number, 1-based. Default: `1`.

`--page-size` (integer)
: Number of items per page. Default: `50`.

### Examples

```bash
grn vserver flavor list --family general-purpose --code Intel

grn vserver flavor list \
  --family memory-optimized \
  --code AMD \
  --zone-id zone-abc123 \
  --output table

# JMESPath to see only IDs and vCPU/RAM
grn vserver flavor list \
  --family general-purpose \
  --code Intel \
  --query "data[].{id: flavorId, cpu: cpu, ram: ram}"
```
