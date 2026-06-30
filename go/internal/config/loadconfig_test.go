package config

import (
	"os"
	"path/filepath"
	"testing"
)

// isolateConfigEnv points HOME at a temp dir and clears GRN_* env vars so
// LoadConfig reads only the files we create under <tmp>/.greenode.
func isolateConfigEnv(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	for _, k := range []string{
		"GRN_PROFILE", "GRN_ACCESS_KEY_ID", "GRN_SECRET_ACCESS_KEY",
		"GRN_DEFAULT_REGION", "GRN_DEFAULT_PROJECT_ID",
	} {
		t.Setenv(k, "")
	}
	dir := filepath.Join(home, ".greenode")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	return dir
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// Regression: a profile present only in the config file (e.g. just after
// `configure set region`) must load its region/output without erroring, even
// when a credentials file exists but has no section for that profile.
func TestLoadConfigProfileInConfigOnly(t *testing.T) {
	dir := isolateConfigEnv(t)
	writeFile(t, filepath.Join(dir, "credentials"),
		"[default]\nclient_id = AKIA-default\nclient_secret = secret-default\n")
	writeFile(t, filepath.Join(dir, "config"),
		"[profile ghost]\nregion = HCM-3\noutput = table\n")

	cfg, err := LoadConfig("ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("cfg is nil")
	}
	if cfg.Region != "HCM-3" {
		t.Errorf("region = %q, want HCM-3", cfg.Region)
	}
	if cfg.Output != "table" {
		t.Errorf("output = %q, want table", cfg.Output)
	}
	if cfg.ClientID != "" || cfg.ClientSecret != "" {
		t.Errorf("expected empty credentials, got id=%q secret=%q", cfg.ClientID, cfg.ClientSecret)
	}
}

// A profile in neither file, when config files exist, errors like
// `aws configure` ("could not be found").
func TestLoadConfigProfileInNeitherFile(t *testing.T) {
	dir := isolateConfigEnv(t)
	writeFile(t, filepath.Join(dir, "credentials"),
		"[default]\nclient_id = AKIA-default\nclient_secret = secret-default\n")

	if _, err := LoadConfig("ghost"); err == nil {
		t.Error("expected error for profile present in neither file")
	}
}

// A profile present only in credentials loads its credentials with no error and
// applies the default output.
func TestLoadConfigProfileInCredsOnly(t *testing.T) {
	dir := isolateConfigEnv(t)
	writeFile(t, filepath.Join(dir, "credentials"),
		"[work]\nclient_id = AKIA-work\nclient_secret = secret-work\n")

	cfg, err := LoadConfig("work")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ClientID != "AKIA-work" {
		t.Errorf("client_id = %q, want AKIA-work", cfg.ClientID)
	}
	if cfg.Output != "json" {
		t.Errorf("output = %q, want json (default)", cfg.Output)
	}
}

// Fresh machine with no config files: not an error (preserves first-run UX);
// returns a populated default cfg.
func TestLoadConfigNoFiles(t *testing.T) {
	isolateConfigEnv(t)

	cfg, err := LoadConfig("default")
	if err != nil {
		t.Fatalf("unexpected error on fresh machine: %v", err)
	}
	if cfg == nil || cfg.Output != "json" {
		t.Errorf("expected non-nil cfg with default output, got %#v", cfg)
	}
}
