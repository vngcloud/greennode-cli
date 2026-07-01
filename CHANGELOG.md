# Changelog

## 1.4.0

### Features
* **output**: Implement the --color flag (on/off/auto, like aws): color status values (ACTIVE=green, ERROR/FAILED=red, CREATING/PENDING=yellow) in text and table output. auto colors only when stdout is a terminal and NO_COLOR is unset; JSON output is never colored

### Enhancements
* **core**: Add --dry-run to destructive commands that lacked it (vServer server stop/reboot/delete, volume/vpc/subnet/secgroup/secgroup-rule delete, and vks delete-auto-upgrade-config) and unify the preview + confirmation flow across delete commands via shared helpers

### Bug Fixes
* **configure**: Mask credential values in 'configure set' output so client_id/client_secret are no longer echoed in plaintext to stdout (consistent with 'configure list'); non-sensitive values are still shown
* **vks**: grn vks wait now aborts immediately on a permanent error (HTTP 403/401/400, or 404 for an active waiter) instead of polling until timeout; transient errors (5xx, network) still retry

### API Changes
* **vks**: Rename 'set-auto-upgrade-config' to 'config-auto-upgrade'; the old name still works as a deprecated alias

## 1.3.1

### Enhancements
* **core**: On error, print a single concise 'Error: ...' line instead of cobra's duplicated error plus a full usage dump (use 'grn <command> --help' for usage)
* **vks**: Document the --unhealthy-range format in config-auto-healing help: expects "[min-max]" (e.g. "[2-5]")

### Bug Fixes
* **configure**: Fix panic (nil pointer) when running 'grn configure set/list --profile <name>' for a profile that does not exist yet. LoadConfig now reads the credentials and config files independently, so a profile created via 'configure set' (config file only) loads its region/output/project_id correctly, and 'configure get'/'configure list' on a truly missing profile report a clear 'profile does not exist' error (exit 1, like 'aws configure') instead of crashing
* **core**: Honor --cli-connect-timeout: the TCP connect and TLS handshake are now bounded by the flag (previously it was accepted but ignored, so a slow/unreachable endpoint hung for ~127s regardless). Wired via the HTTP transport's dialer in both VKS and vServer clients
* **output**: Reject an invalid --output or --color value (e.g. a typo like 'tabel') with a clear error, a 'maybe you meant' suggestion, and a non-zero exit, instead of silently falling back

## 1.3.0

### Features
* **vks**: Add struct-valued flags (shorthand or JSON) to create-nodegroup: --tags, --secondary-subnets, --auto-scale, --placement-group, --upgrade-config

### API Changes
* **vks**: update-nodegroup: drop deprecated --labels/--taints (use update-nodegroup-metadata); replace --auto-scale-min/max and --upgrade-strategy/max-surge/max-unavailable with struct flags --auto-scale and --upgrade-config

## 1.2.0

### Bug Fixes
* **configure**: Fix panic (nil pointer) when running 'grn configure --profile <name>' for a profile that does not exist yet; start from empty defaults so the profile can be created
* **vks**: Return a non-zero exit code and 'unknown command' error for invalid subcommands of grn vks/vserver (and their groups) instead of silently printing help with exit 0

### API Changes
* **vks**: Remove get-nodegroup-events; network-type CALICO->TIGERA; --os adds rocky; replace boolean toggle flags with --private-cluster/--private-nodes/--load-balancer-plugin/--block-store-csi-plugin <enabled|disabled>

## 1.1.0

### Features
* **vserver**: Add vServer CLI: vpc, subnet, secgroup, volume, volume-type, flavor, image, and server commands

### API Changes
* **vks**: Align node group commands with latest API: remove --image-id (no longer in the API), add --os (ubuntu/linux, default ubuntu) for create-cluster/create-nodegroup; update-nodegroup no longer requires/sends image
