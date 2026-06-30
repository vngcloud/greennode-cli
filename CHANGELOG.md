# Changelog

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
