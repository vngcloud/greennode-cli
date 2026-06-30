package vks

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var updateClusterCmd = &cobra.Command{
	Use:   "update-cluster",
	Short: "Update a VKS cluster",
	RunE:  runUpdateCluster,
}

func init() {
	f := updateClusterCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.String("k8s-version", "", "Kubernetes version (required)")
	f.String("whitelist-node-cidrs", "", "Whitelist CIDRs, comma-separated (required)")
	f.String("load-balancer-plugin", "", "Load balancer plugin (enabled, disabled); unset = unchanged")
	f.String("block-store-csi-plugin", "", "Block store CSI plugin (enabled, disabled); unset = unchanged")
	f.Bool("dry-run", false, "Validate parameters without updating")

	updateClusterCmd.MarkFlagRequired("cluster-id")
	updateClusterCmd.MarkFlagRequired("k8s-version")
	updateClusterCmd.MarkFlagRequired("whitelist-node-cidrs")
}

func runUpdateCluster(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	k8sVersion, _ := cmd.Flags().GetString("k8s-version")
	whitelistCIDRs, _ := cmd.Flags().GetString("whitelist-node-cidrs")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}

	body := map[string]interface{}{
		"version":            k8sVersion,
		"whitelistNodeCIDRs": parseCommaSeparated(whitelistCIDRs),
	}

	// Plugin toggles are only sent when explicitly provided (unset = unchanged).
	if cmd.Flags().Changed("load-balancer-plugin") {
		v, _ := cmd.Flags().GetString("load-balancer-plugin")
		enabled, err := parseToggle("load-balancer-plugin", v)
		if err != nil {
			return err
		}
		body["enabledLoadBalancerPlugin"] = enabled
	}
	if cmd.Flags().Changed("block-store-csi-plugin") {
		v, _ := cmd.Flags().GetString("block-store-csi-plugin")
		enabled, err := parseToggle("block-store-csi-plugin", v)
		if err != nil {
			return err
		}
		body["enabledBlockStoreCsiPlugin"] = enabled
	}

	if dryRun {
		fmt.Println("=== DRY RUN: Update cluster ===")
		fmt.Println()
		fmt.Printf("Cluster ID: %s\n", clusterID)
		fmt.Printf("New version: %s\n", k8sVersion)
		fmt.Printf("Whitelist CIDRs: %s\n", whitelistCIDRs)
		if v, ok := body["enabledLoadBalancerPlugin"]; ok {
			fmt.Printf("Load balancer plugin: %v\n", v)
		}
		if v, ok := body["enabledBlockStoreCsiPlugin"]; ok {
			fmt.Printf("Block store CSI plugin: %v\n", v)
		}
		fmt.Println("\nRun without --dry-run to update.")
		return nil
	}

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Put(fmt.Sprintf("/v1/clusters/%s", clusterID), body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}
