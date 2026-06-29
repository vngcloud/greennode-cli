# image

Browse available vServer OS and GPU images.

```bash
grn vserver image <command> [options]
```

## Commands

| Command | Description |
|---------|-------------|
| [list](#list) | List available images |

---

## list

List available images by type. Use the image ID with `grn vserver server create --image-id`.

### Synopsis

```
grn vserver image list
    --type <value>
    [--page <value>]
    [--page-size <value>]
    [--image-version <value>]
```

### Options

`--type` (required)
: Image type. Accepted values: `os`, `gpu`.

`--page` (integer)
: Page number, 1-based. Default: `1`.

`--page-size` (integer)
: Number of items per page. Default: `50`.

`--image-version` (string)
: Filter by image version (client-side substring match, case-insensitive).

### Examples

```bash
# List OS images
grn vserver image list --type os

# List GPU images
grn vserver image list --type gpu

# Filter to Ubuntu images
grn vserver image list --type os --image-version ubuntu

# Table output with JMESPath to extract ID and name
grn vserver image list --type os \
  --query "images[].{id: id, name: name, version: imageVersion}" \
  --output table
```
