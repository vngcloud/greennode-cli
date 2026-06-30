# Changelog

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
