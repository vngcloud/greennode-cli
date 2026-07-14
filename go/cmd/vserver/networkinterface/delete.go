package networkinterface

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
	Short: "Delete a network interface",
	RunE:  runDelete,
}

func init() {
	f := deleteCmd.Flags()
	f.String("network-interface-id", "", "Network interface ID (required)")
	f.Bool("force", false, "Skip confirmation prompt")
	deleteCmd.MarkFlagRequired("network-interface-id")
}

func runDelete(cmd *cobra.Command, args []string) error {
	interfaceID, _ := cmd.Flags().GetString("network-interface-id")
	force, _ := cmd.Flags().GetBool("force")

	if err := validator.ValidateID(interfaceID, "network-interface-id"); err != nil {
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
		fmt.Println("The following network interface will be deleted:")
		fmt.Println()
		fmt.Printf("  ID: %s\n", interfaceID)
		fmt.Println()
		fmt.Println("This action is irreversible.")
		fmt.Print("\nAre you sure you want to delete this network interface? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	result, err := apiClient.Delete(fmt.Sprintf("/v2/%s/network-interfaces-elastic/%s", projectID, interfaceID), nil)
	if err != nil {
		return fmt.Errorf("failed to delete network interface %s: %w", interfaceID, err)
	}

	return outputResult(cmd, cfg, result)
}
