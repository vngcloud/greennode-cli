package placementgroup

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
	Short: "Delete a placement group",
	RunE:  runDelete,
}

func init() {
	f := deleteCmd.Flags()
	f.String("placement-group-id", "", "Placement group (server group) ID (required)")
	f.Bool("force", false, "Skip confirmation prompt")
	deleteCmd.MarkFlagRequired("placement-group-id")
}

func runDelete(cmd *cobra.Command, args []string) error {
	groupID, _ := cmd.Flags().GetString("placement-group-id")
	force, _ := cmd.Flags().GetBool("force")

	if err := validator.ValidateID(groupID, "placement-group-id"); err != nil {
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
		fmt.Println("The following placement group will be deleted:")
		fmt.Println()
		fmt.Printf("  ID: %s\n", groupID)
		fmt.Println()
		fmt.Println("This action is irreversible.")
		fmt.Print("\nAre you sure you want to delete this placement group? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	result, err := apiClient.Delete(fmt.Sprintf("/v2/%s/serverGroups/%s", projectID, groupID), nil)
	if err != nil {
		return fmt.Errorf("failed to delete placement group %s: %w", groupID, err)
	}

	return outputResult(cmd, cfg, result)
}
