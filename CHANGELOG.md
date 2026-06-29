# Changelog

## 1.1.0

### Features
* **vserver**: Add vServer CLI: vpc, subnet, secgroup, volume, volume-type, flavor, image, and server commands

### API Changes
* **vks**: Align node group commands with latest API: remove --image-id (no longer in the API), add --os (ubuntu/linux, default ubuntu) for create-cluster/create-nodegroup; update-nodegroup no longer requires/sends image
