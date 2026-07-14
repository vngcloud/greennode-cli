package sshkey

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
	Short: "Delete an SSH key",
	RunE:  runDelete,
}

func init() {
	f := deleteCmd.Flags()
	f.String("sshkey-id", "", "SSH key ID (required)")
	f.Bool("force", false, "Skip confirmation prompt")
	deleteCmd.MarkFlagRequired("sshkey-id")
}

func runDelete(cmd *cobra.Command, args []string) error {
	sshKeyID, _ := cmd.Flags().GetString("sshkey-id")
	force, _ := cmd.Flags().GetBool("force")

	if err := validator.ValidateID(sshKeyID, "sshkey-id"); err != nil {
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
		fmt.Println("The following SSH key will be deleted:")
		fmt.Println()
		fmt.Printf("  ID: %s\n", sshKeyID)
		fmt.Println()
		fmt.Println("This action is irreversible.")
		fmt.Print("\nAre you sure you want to delete this SSH key? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	result, err := apiClient.Delete(fmt.Sprintf("/v2/%s/sshKeys/%s", projectID, sshKeyID), nil)
	if err != nil {
		return fmt.Errorf("failed to delete SSH key %s: %w", sshKeyID, err)
	}

	return outputResult(cmd, cfg, result)
}
