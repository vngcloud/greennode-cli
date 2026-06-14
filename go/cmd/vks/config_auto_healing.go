package vks

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var configAutoHealingCmd = &cobra.Command{
	Use:   "config-auto-healing",
	Short: "Configure auto-healing for a VKS cluster",
	RunE:  runConfigAutoHealing,
}

func init() {
	f := configAutoHealingCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.Bool("enable-auto-healing", false, "Enable auto-healing (required)")
	f.String("max-unhealthy", "", "Max unhealthy nodes, e.g. \"30%\"")
	f.String("unhealthy-range", "", "Unhealthy range")
	f.Int("timeout-unhealthy", 0, "Unhealthy timeout in seconds")

	configAutoHealingCmd.MarkFlagRequired("cluster-id")
	configAutoHealingCmd.MarkFlagRequired("enable-auto-healing")
}

func buildAutoHealingBody(enable bool, maxUnhealthy, unhealthyRange string, timeoutUnhealthy int, changed map[string]bool) map[string]interface{} {
	body := map[string]interface{}{"enableAutoHealing": enable}
	if changed["max-unhealthy"] && maxUnhealthy != "" {
		body["maxUnhealthy"] = maxUnhealthy
	}
	if changed["unhealthy-range"] && unhealthyRange != "" {
		body["unhealthyRange"] = unhealthyRange
	}
	if changed["timeout-unhealthy"] {
		body["timeoutUnhealthy"] = timeoutUnhealthy
	}
	return body
}

func runConfigAutoHealing(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	enable, _ := cmd.Flags().GetBool("enable-auto-healing")
	maxUnhealthy, _ := cmd.Flags().GetString("max-unhealthy")
	unhealthyRange, _ := cmd.Flags().GetString("unhealthy-range")
	timeoutUnhealthy, _ := cmd.Flags().GetInt("timeout-unhealthy")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}

	changed := map[string]bool{
		"max-unhealthy":     cmd.Flags().Changed("max-unhealthy"),
		"unhealthy-range":   cmd.Flags().Changed("unhealthy-range"),
		"timeout-unhealthy": cmd.Flags().Changed("timeout-unhealthy"),
	}
	body := buildAutoHealingBody(enable, maxUnhealthy, unhealthyRange, timeoutUnhealthy, changed)

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Patch(
		fmt.Sprintf("/v1/clusters/%s/auto-healing-config", clusterID), body,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}
