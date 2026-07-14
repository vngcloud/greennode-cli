package dhcp

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a DHCP option",
	RunE:  runDelete,
}

func init() {
	f := deleteCmd.Flags()
	f.String("dhcp-option-id", "", "DHCP option ID (required)")
	f.Bool("force", false, "Skip confirmation prompt")
	deleteCmd.MarkFlagRequired("dhcp-option-id")
}

func runDelete(cmd *cobra.Command, args []string) error {
	dhcpID, _ := cmd.Flags().GetString("dhcp-option-id")
	force, _ := cmd.Flags().GetBool("force")

	if err := validator.ValidateID(dhcpID, "dhcp-option-id"); err != nil {
		return err
	}

	apiClient, cfg, err := createClient(cmd)
	if err != nil {
		return err
	}

	projectID, err := getProjectID(cfg)
	if err != nil {
		return err
	}

	if !force {
		fmt.Println("The following DHCP option will be deleted:")
		fmt.Println()
		fmt.Printf("  ID: %s\n", dhcpID)
		fmt.Println()
		fmt.Println("This action is irreversible.")
		fmt.Print("\nAre you sure you want to delete this DHCP option? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	result, err := apiClient.Delete(fmt.Sprintf("/v2/%s/dhcp_option/%s", projectID, dhcpID), nil)
	if err != nil {
		return fmt.Errorf("failed to delete DHCP option %s: %w", dhcpID, err)
	}

	return outputResult(cmd, cfg, result)
}
