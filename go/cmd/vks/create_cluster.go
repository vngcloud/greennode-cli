package vks

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/cli"
)

var createClusterCmd = &cobra.Command{
	Use:   "create-cluster",
	Short: "Create a new VKS cluster",
	Long: "Create a new VKS cluster (control plane only). " +
		"Add worker nodes afterwards with 'grn vks create-nodegroup'.",
	RunE: runCreateCluster,
}

func init() {
	f := createClusterCmd.Flags()
	// Cluster settings (required)
	f.String("name", "", "Cluster name (required)")
	f.String("k8s-version", "", "Kubernetes version (required)")
	f.String("network-type", "", "Network type: TIGERA, CILIUM_OVERLAY, CILIUM_NATIVE_ROUTING (required)")
	f.String("vpc-id", "", "VPC ID (required)")
	f.String("subnet-id", "", "Subnet ID (required)")

	for _, name := range []string{"name", "k8s-version", "network-type", "vpc-id", "subnet-id"} {
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
		"subnetId":                   subnetID,
		"enablePrivateCluster":       enablePrivateCluster,
		"releaseChannel":             releaseChannel,
		"enabledBlockStoreCsiPlugin": enabledCSIPlugin,
		"enabledLoadBalancerPlugin":  enabledLBPlugin,
		"enabledServiceEndpoint":     enabledServiceEndpoint,
		"azStrategy":                 azStrategy,
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

	if dryRun {
		return validateCreateCluster(name, networkType, cidr)
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

func validateCreateCluster(name, networkType, cidr string) error {
	clusterNameRE := regexp.MustCompile(`^[a-z0-9][a-z0-9\-]{3,18}[a-z0-9]$`)

	var errors []string

	if !clusterNameRE.MatchString(name) {
		errors = append(errors, fmt.Sprintf(
			"Cluster name '%s' is invalid. Must be 5-20 chars, lowercase alphanumeric and hyphens, start/end with alphanumeric.", name))
	}

	if (networkType == "TIGERA" || networkType == "CILIUM_OVERLAY") && cidr == "" {
		errors = append(errors, fmt.Sprintf("--cidr is required when network-type is %s", networkType))
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
