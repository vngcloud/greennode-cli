package floatingip

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
	Short: "Delete a floating IP",
	RunE:  runDelete,
}

func init() {
	f := deleteCmd.Flags()
	f.String("floating-ip-id", "", "Floating IP ID (required)")
	f.Bool("force", false, "Skip confirmation prompt")
	deleteCmd.MarkFlagRequired("floating-ip-id") //nolint:errcheck
}

func runDelete(cmd *cobra.Command, args []string) error {
	ipID, _ := cmd.Flags().GetString("floating-ip-id")
	force, _ := cmd.Flags().GetBool("force")

	if err := validator.ValidateID(ipID, "floating-ip-id"); err != nil {
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
		fmt.Println("The following floating IP will be deleted:")
		fmt.Println()
		fmt.Printf("  ID: %s\n", ipID)
		fmt.Println()
		fmt.Println("This action is irreversible.")
		fmt.Print("\nAre you sure you want to delete this floating IP? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	result, err := apiClient.Delete(fmt.Sprintf("/v2/%s/wanIps/%s", projectID, ipID), nil)
	if err != nil {
		return fmt.Errorf("failed to delete floating IP %s: %w", ipID, err)
	}

	return outputResult(cmd, cfg, result)
}
