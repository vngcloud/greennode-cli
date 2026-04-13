package configure

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/config"
)

var validRegions = []string{"HCM-3", "HAN"}
var validOutputs = []string{"json", "text", "table"}

// ConfigureCmd is the `grn configure` command.
var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure GreenNode CLI credentials and settings",
	Long: `Interactive setup for GreenNode CLI.

Prompts for Client ID, Client Secret, Region, and Output format.
Saves credentials to ~/.greenode/credentials and config to ~/.greenode/config.`,
	Run: runConfigure,
}

func init() {
	ConfigureCmd.AddCommand(listCmd)
	ConfigureCmd.AddCommand(getCmd)
	ConfigureCmd.AddCommand(setCmd)
}

func runConfigure(cmd *cobra.Command, args []string) {
	profile := cmd.Flag("profile").Value.String()
	if profile == "" {
		profile = os.Getenv("GRN_PROFILE")
	}
	if profile == "" {
		profile = "default"
	}

	// Load existing config for defaults
	cfg, _ := config.LoadConfig(profile)

	reader := bufio.NewReader(os.Stdin)

	clientID := promptWithDefault(reader, "Client ID", maskCred(cfg.ClientID))
	clientSecret := promptWithDefault(reader, "Client Secret", maskCred(cfg.ClientSecret))
	region := promptWithDefault(reader, "Default region name", cfg.Region)
	output := promptWithDefault(reader, "Default output format", cfg.Output)

	// If user entered masked value or empty, keep original
	if clientID == maskCred(cfg.ClientID) || clientID == "" {
		clientID = cfg.ClientID
	}
	if clientSecret == maskCred(cfg.ClientSecret) || clientSecret == "" {
		clientSecret = cfg.ClientSecret
	}

	// Validate region
	if !contains(validRegions, region) {
		fmt.Fprintf(os.Stderr, "Warning: invalid region '%s', using default 'HCM-3'\n", region)
		region = "HCM-3"
	}

	// Validate output
	if !contains(validOutputs, output) {
		fmt.Fprintf(os.Stderr, "Warning: invalid output format '%s', using default 'json'\n", output)
		output = "json"
	}

	writer := config.NewConfigFileWriter()

	if err := writer.WriteCredentials(profile, clientID, clientSecret); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving credentials: %v\n", err)
		os.Exit(1)
	}

	if err := writer.WriteConfig(profile, region, output); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Configuration saved successfully.")
}

func promptWithDefault(reader *bufio.Reader, prompt, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultVal)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

func maskCred(value string) string {
	if value == "" {
		return ""
	}
	return config.MaskCredential(value)
}

func contains(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}
