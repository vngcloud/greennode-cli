package vks

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listClusterVersionsCmd = &cobra.Command{
	Use:   "list-cluster-versions",
	Short: "List available Kubernetes versions for VKS clusters",
	RunE:  runListClusterVersions,
}

func runListClusterVersions(cmd *cobra.Command, args []string) error {
	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Get("/v1/cluster-versions", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}
