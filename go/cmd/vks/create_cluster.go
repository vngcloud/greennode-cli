package vks

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/cli"
)

var createClusterCmd = &cobra.Command{
	Use:   "create-cluster",
	Short: "Create a new VKS cluster",
	Long: "Create a new VKS cluster. By default only the control plane is provisioned; " +
		"pass --node-group-name (with --flavor-id, --disk-type, --ssh-key-id) to also attach a " +
		"node group at creation, or add one later with 'grn vks create-nodegroup'.",
	RunE: runCreateCluster,
}

func init() {
	f := createClusterCmd.Flags()
	// Cluster settings (required)
	f.String("name", "", "Cluster name (required)")
	f.String("k8s-version", "", "Kubernetes version (required)")
	f.String("network-type", "", "Network type: TIGERA, CILIUM_OVERLAY, CILIUM_NATIVE_ROUTING (required)")
	f.String("vpc-id", "", "VPC ID (required)")
	f.String("subnet-id", "", "Subnet ID")

	for _, name := range []string{"name", "k8s-version", "network-type", "vpc-id"} {
		createClusterCmd.MarkFlagRequired(name)
	}

	// Cluster settings (optional)
	f.String("cidr", "", "CIDR block (required for TIGERA and CILIUM_OVERLAY)")
	f.String("description", "", "Cluster description")
	f.String("private-cluster", "disabled", "Private cluster (enabled, disabled)")
	f.String("release-channel", "STABLE", "Release channel (RAPID, STABLE)")
	f.String("load-balancer-plugin", "enabled", "Load balancer plugin (enabled, disabled)")
	f.String("block-store-csi-plugin", "enabled", "Block store CSI plugin (enabled, disabled)")
	f.String("service-endpoint", "disabled", "Service endpoint (enabled, disabled)")
	f.String("az-strategy", "SINGLE", "Availability zone strategy")
	f.String("secondary-subnets", "", "Secondary subnet IDs (comma-separated)")
	f.String("list-subnet-ids", "", "Subnet IDs for the cluster (comma-separated)")
	f.Int("node-netmask-size", 0, "Node netmask size")
	f.String("auto-upgrade-config", "", "Auto-upgrade config (shorthand time=03:00,weekdays=Mon or JSON; use JSON for multiple weekdays)")
	f.String("auto-healing-config", "", "Auto-healing config (shorthand enableAutoHealing=true,maxUnhealthy=20%,unhealthyRange=[2-5],timeoutUnhealthy=10 or JSON)")

	// Node group settings (optional; a node group is attached only when
	// --node-group-name is set, and then --flavor-id/--disk-type/--ssh-key-id
	// are required too). Sent as the singular `nodeGroup` object.
	f.String("node-group-name", "", "Node group name (attaches a node group at creation when set)")
	f.String("flavor-id", "", "Node group flavor ID (required when attaching a node group)")
	f.String("os", "ubuntu", "Node group OS image (ubuntu, linux, rocky)")
	f.String("disk-type", "", "Node group disk type ID (required when attaching a node group)")
	f.String("ssh-key-id", "", "Node group SSH key ID (required when attaching a node group)")
	f.Int("disk-size", 100, "Node group disk size in GiB (20-5000)")
	f.Int("num-nodes", 1, "Node group number of nodes (0-10)")
	f.String("private-nodes", "disabled", "Node group private nodes (enabled, disabled)")
	f.String("security-groups", "", "Node group security group IDs (comma-separated)")
	f.String("labels", "", "Node group labels as key=value pairs (comma-separated)")
	f.String("taints", "", "Node group taints as key=value:effect (comma-separated)")

	f.Bool("dry-run", false, "Validate parameters without creating the cluster")
}

func runCreateCluster(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	k8sVersion, _ := cmd.Flags().GetString("k8s-version")
	networkType, _ := cmd.Flags().GetString("network-type")
	vpcID, _ := cmd.Flags().GetString("vpc-id")
	subnetID, _ := cmd.Flags().GetString("subnet-id")
	cidr, _ := cmd.Flags().GetString("cidr")
	description, _ := cmd.Flags().GetString("description")
	releaseChannel, _ := cmd.Flags().GetString("release-channel")
	azStrategy, _ := cmd.Flags().GetString("az-strategy")
	secondarySubnets, _ := cmd.Flags().GetString("secondary-subnets")
	listSubnetIDs, _ := cmd.Flags().GetString("list-subnet-ids")
	autoUpgradeStr, _ := cmd.Flags().GetString("auto-upgrade-config")
	autoHealingStr, _ := cmd.Flags().GetString("auto-healing-config")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	// Parse enabled/disabled toggle flags.
	privateClusterVal, _ := cmd.Flags().GetString("private-cluster")
	lbPluginVal, _ := cmd.Flags().GetString("load-balancer-plugin")
	csiPluginVal, _ := cmd.Flags().GetString("block-store-csi-plugin")
	serviceEndpointVal, _ := cmd.Flags().GetString("service-endpoint")
	enablePrivateCluster, err := parseToggle("private-cluster", privateClusterVal)
	if err != nil {
		return err
	}
	enabledLBPlugin, err := parseToggle("load-balancer-plugin", lbPluginVal)
	if err != nil {
		return err
	}
	enabledCSIPlugin, err := parseToggle("block-store-csi-plugin", csiPluginVal)
	if err != nil {
		return err
	}
	enabledServiceEndpoint, err := parseToggle("service-endpoint", serviceEndpointVal)
	if err != nil {
		return err
	}

	// Build cluster body. Node groups are created separately via
	// 'grn vks create-nodegroup'.
	body := map[string]interface{}{
		"name":                       name,
		"version":                    k8sVersion,
		"networkType":                networkType,
		"vpcId":                      vpcID,
		"enablePrivateCluster":       enablePrivateCluster,
		"releaseChannel":             releaseChannel,
		"enabledBlockStoreCsiPlugin": enabledCSIPlugin,
		"enabledLoadBalancerPlugin":  enabledLBPlugin,
		"enabledServiceEndpoint":     enabledServiceEndpoint,
		"azStrategy":                 azStrategy,
	}

	if subnetID != "" {
		body["subnetId"] = subnetID
	}
	if cidr != "" {
		body["cidr"] = cidr
	}
	if description != "" {
		body["description"] = description
	}
	if secondarySubnets != "" {
		body["secondarySubnets"] = parseCommaSeparated(secondarySubnets)
	}
	if listSubnetIDs != "" {
		body["listSubnetIds"] = parseCommaSeparated(listSubnetIDs)
	}
	if cmd.Flags().Changed("node-netmask-size") {
		nodeNetmaskSize, _ := cmd.Flags().GetInt("node-netmask-size")
		body["nodeNetmaskSize"] = nodeNetmaskSize
	}
	if autoUpgradeStr != "" {
		uc, err := cli.ParseStructFlag(autoUpgradeStr)
		if err != nil {
			return fmt.Errorf("--auto-upgrade-config: %w", err)
		}
		body["autoUpgradeConfig"] = uc
	}
	if autoHealingStr != "" {
		hc, err := cli.ParseStructFlagTyped(autoHealingStr, []string{"timeoutUnhealthy"}, []string{"enableAutoHealing"})
		if err != nil {
			return fmt.Errorf("--auto-healing-config: %w", err)
		}
		body["autoHealingConfig"] = hc
	}

	// Optionally attach a node group. Triggered when any node-group flag is set;
	// sent as the singular `nodeGroup` object (the API's `nodeGroups` field is
	// deprecated).
	ngFlags := []string{"node-group-name", "flavor-id", "os", "disk-type", "ssh-key-id",
		"disk-size", "num-nodes", "private-nodes", "security-groups", "labels", "taints"}
	wantsNodeGroup := false
	for _, fl := range ngFlags {
		if cmd.Flags().Changed(fl) {
			wantsNodeGroup = true
			break
		}
	}
	var ngName string
	var ngDiskSize, ngNumNodes int
	if wantsNodeGroup {
		ngName, _ = cmd.Flags().GetString("node-group-name")
		flavorID, _ := cmd.Flags().GetString("flavor-id")
		diskType, _ := cmd.Flags().GetString("disk-type")
		sshKeyID, _ := cmd.Flags().GetString("ssh-key-id")

		var missing []string
		if ngName == "" {
			missing = append(missing, "--node-group-name")
		}
		if flavorID == "" {
			missing = append(missing, "--flavor-id")
		}
		if diskType == "" {
			missing = append(missing, "--disk-type")
		}
		if sshKeyID == "" {
			missing = append(missing, "--ssh-key-id")
		}
		if len(missing) > 0 {
			return fmt.Errorf("attaching a node group requires: %s", strings.Join(missing, ", "))
		}

		osImage, _ := cmd.Flags().GetString("os")
		ngDiskSize, _ = cmd.Flags().GetInt("disk-size")
		ngNumNodes, _ = cmd.Flags().GetInt("num-nodes")
		privateNodesVal, _ := cmd.Flags().GetString("private-nodes")
		nodeSecurityGroups, _ := cmd.Flags().GetString("security-groups")
		labels, _ := cmd.Flags().GetString("labels")
		taints, _ := cmd.Flags().GetString("taints")
		enablePrivateNodes, err := parseToggle("private-nodes", privateNodesVal)
		if err != nil {
			return err
		}

		nodeGroup := map[string]interface{}{
			"name":               ngName,
			"flavorId":           flavorID,
			"os":                 osImage,
			"diskSize":           ngDiskSize,
			"diskType":           diskType,
			"numNodes":           ngNumNodes,
			"enablePrivateNodes": enablePrivateNodes,
			"sshKeyId":           sshKeyID,
			"upgradeConfig": map[string]interface{}{
				"maxSurge":       1,
				"maxUnavailable": 0,
				"strategy":       "SURGE",
			},
			"securityGroups": []string{},
		}
		if subnetID != "" {
			nodeGroup["subnetId"] = subnetID
		}
		if nodeSecurityGroups != "" {
			nodeGroup["securityGroups"] = parseCommaSeparated(nodeSecurityGroups)
		}
		if labels != "" {
			nodeGroup["labels"] = parseLabels(labels)
		}
		if taints != "" {
			nodeGroup["taints"] = parseTaints(taints)
		}
		body["nodeGroup"] = nodeGroup
	}

	if dryRun {
		return validateCreateCluster(name, networkType, cidr, wantsNodeGroup, ngName, ngDiskSize, ngNumNodes)
	}

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Post("/v1/clusters", body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}

func validateCreateCluster(name, networkType, cidr string, wantsNodeGroup bool, ngName string, diskSize, numNodes int) error {
	clusterNameRE := regexp.MustCompile(`^[a-z0-9][a-z0-9\-]{3,18}[a-z0-9]$`)

	var errors []string

	if !clusterNameRE.MatchString(name) {
		errors = append(errors, fmt.Sprintf(
			"Cluster name '%s' is invalid. Must be 5-20 chars, lowercase alphanumeric and hyphens, start/end with alphanumeric.", name))
	}

	if (networkType == "TIGERA" || networkType == "CILIUM_OVERLAY") && cidr == "" {
		errors = append(errors, fmt.Sprintf("--cidr is required when network-type is %s", networkType))
	}

	if wantsNodeGroup {
		ngNameRE := regexp.MustCompile(`^[a-z0-9][a-z0-9-]{3,13}[a-z0-9]$`)
		if !ngNameRE.MatchString(ngName) {
			errors = append(errors, fmt.Sprintf(
				"Node group name '%s' is invalid. Must be 5-15 chars, lowercase alphanumeric and hyphens, start/end with alphanumeric.", ngName))
		}
		if diskSize < 20 || diskSize > 5000 {
			errors = append(errors, fmt.Sprintf("Disk size %d out of range (20-5000 GiB)", diskSize))
		}
		if numNodes < 0 || numNodes > 10 {
			errors = append(errors, fmt.Sprintf("Number of nodes %d out of range (0-10)", numNodes))
		}
	}

	fmt.Println("=== DRY RUN: Validation results ===")
	fmt.Println()
	if len(errors) > 0 {
		fmt.Printf("Found %d error(s):\n", len(errors))
		for _, e := range errors {
			fmt.Printf("  - %s\n", e)
		}
		os.Exit(1)
	}

	fmt.Println("All parameters are valid. Run without --dry-run to create the cluster.")
	return nil
}
