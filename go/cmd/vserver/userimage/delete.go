package userimage

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
	Short: "Delete a user image",
	RunE:  runDelete,
}

func init() {
	f := deleteCmd.Flags()
	f.String("user-image-id", "", "User image ID (required)")
	f.Bool("force", false, "Skip confirmation prompt")
	deleteCmd.MarkFlagRequired("user-image-id")
}

func runDelete(cmd *cobra.Command, args []string) error {
	imageID, _ := cmd.Flags().GetString("user-image-id")
	force, _ := cmd.Flags().GetBool("force")

	if err := validator.ValidateID(imageID, "user-image-id"); err != nil {
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
		fmt.Println("The following user image will be deleted:")
		fmt.Println()
		fmt.Printf("  ID: %s\n", imageID)
		fmt.Println()
		fmt.Println("This action is irreversible.")
		fmt.Print("\nAre you sure you want to delete this user image? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	result, err := apiClient.Delete(fmt.Sprintf("/v2/%s/user-images/%s", projectID, imageID), nil)
	if err != nil {
		return fmt.Errorf("failed to delete user image %s: %w", imageID, err)
	}

	return outputResult(cmd, cfg, result)
}
