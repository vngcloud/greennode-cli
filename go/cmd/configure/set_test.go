package configure

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// newConfigureTestCmd wires a root command with the persistent `profile` flag
// (registered on rootCmd in production) so the configure subcommands resolve it
// exactly as they do at runtime.
func newConfigureTestCmd() *cobra.Command {
	root := &cobra.Command{Use: "grn"}
	root.PersistentFlags().String("profile", "", "")
	root.AddCommand(ConfigureCmd)
	return root
}

// Regression: `configure set <key> <value> --profile <new>` for a profile that
// does not exist in the credentials file must not panic. LoadConfig returns
// (nil, err) for an unknown profile; set previously dereferenced the nil *Config.
func TestSetRegionOnNonExistentProfile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	root := newConfigureTestCmd()
	root.SetArgs([]string{"configure", "set", "region", "HCM-3", "--profile", "prod-hcm-qc3"})
	if err := root.Execute(); err != nil {
		t.Fatalf("set region on new profile failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(home, ".greenode", "config"))
	if err != nil {
		t.Fatalf("config file not written: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "prod-hcm-qc3") || !strings.Contains(content, "HCM-3") {
		t.Errorf("config missing profile/region; got:\n%s", content)
	}
}

// Regression: `configure list --profile <new>` for an unknown profile must not
// panic and should render unset values.
func TestListOnNonExistentProfile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	root := newConfigureTestCmd()
	root.SetArgs([]string{"configure", "list", "--profile", "ghost"})
	if err := root.Execute(); err != nil {
		t.Fatalf("list on new profile failed: %v", err)
	}
}
