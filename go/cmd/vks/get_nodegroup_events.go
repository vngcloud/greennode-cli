package vks

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var getNodegroupEventsCmd = &cobra.Command{
	Use:   "get-nodegroup-events",
	Short: "Get the list of events for a node group",
	RunE:  runGetNodegroupEvents,
}

func init() {
	f := getNodegroupEventsCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.String("nodegroup-id", "", "Node group ID (required)")
	f.String("action", "", "Filter by action")
	f.String("type", "", "Filter by event type")
	f.Int("page", 0, "Page number (0-based)")
	f.Int("page-size", 50, "Page size")

	getNodegroupEventsCmd.MarkFlagRequired("cluster-id")
	getNodegroupEventsCmd.MarkFlagRequired("nodegroup-id")
}

func runGetNodegroupEvents(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	nodegroupID, _ := cmd.Flags().GetString("nodegroup-id")
	action, _ := cmd.Flags().GetString("action")
	eventType, _ := cmd.Flags().GetString("type")
	page, _ := cmd.Flags().GetInt("page")
	pageSize, _ := cmd.Flags().GetInt("page-size")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}
	if err := validator.ValidateID(nodegroupID, "nodegroup-id"); err != nil {
		return err
	}

	changed := map[string]bool{
		"action":    cmd.Flags().Changed("action"),
		"type":      cmd.Flags().Changed("type"),
		"page":      cmd.Flags().Changed("page"),
		"page-size": cmd.Flags().Changed("page-size"),
	}
	params := buildEventsQuery(action, eventType, page, pageSize, changed)

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Get(
		fmt.Sprintf("/v1/clusters/%s/node-groups/%s/events", clusterID, nodegroupID), params,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}
