package vks

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var upgradeNodegroupVersionCmd = &cobra.Command{
	Use:   "upgrade-nodegroup-version",
	Short: "Upgrade the Kubernetes version of a node group",
	RunE:  runUpgradeNodegroupVersion,
}

func init() {
	f := upgradeNodegroupVersionCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.String("nodegroup-id", "", "Node group ID (required)")
	f.String("k8s-version", "", "Target Kubernetes version (required)")

	upgradeNodegroupVersionCmd.MarkFlagRequired("cluster-id")
	upgradeNodegroupVersionCmd.MarkFlagRequired("nodegroup-id")
	upgradeNodegroupVersionCmd.MarkFlagRequired("k8s-version")
}

func buildUpgradeNodegroupBody(k8sVersion string) map[string]interface{} {
	return map[string]interface{}{"kubernetesVersion": k8sVersion}
}

func runUpgradeNodegroupVersion(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	nodegroupID, _ := cmd.Flags().GetString("nodegroup-id")
	k8sVersion, _ := cmd.Flags().GetString("k8s-version")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}
	if err := validator.ValidateID(nodegroupID, "nodegroup-id"); err != nil {
		return err
	}

	body := buildUpgradeNodegroupBody(k8sVersion)

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Post(
		fmt.Sprintf("/v1/clusters/%s/node-groups/%s/upgrade-version", clusterID, nodegroupID), body,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}
