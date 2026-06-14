package kubeconfig

import (
	"os"
	"path/filepath"
	"testing"
)

const incomingKubeconfig = `apiVersion: v1
kind: Config
clusters:
- name: vks-cluster
  cluster:
    server: https://1.2.3.4:6443
    certificate-authority-data: AAAA
contexts:
- name: vks-cluster-ctx
  context:
    cluster: vks-cluster
    user: vks-user
current-context: vks-cluster-ctx
users:
- name: vks-user
  user:
    token: secret-token
`

func TestMergeIntoEmptyFileCreatesContext(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "config")

	res, err := Merge(target, incomingKubeconfig, "vks_c-123", true)
	if err != nil {
		t.Fatalf("Merge error: %v", err)
	}
	if res.ContextName != "vks_c-123" {
		t.Errorf("ContextName = %q, want vks_c-123", res.ContextName)
	}

	cfg, err := Load(target)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.CurrentContext != "vks_c-123" {
		t.Errorf("current-context = %q, want vks_c-123", cfg.CurrentContext)
	}
	if findContext(cfg, "vks_c-123") == nil {
		t.Errorf("merged context not found in %#v", cfg.Contexts)
	}
}

func TestMergePreservesExistingContexts(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "config")
	existing := `apiVersion: v1
kind: Config
clusters:
- name: other
  cluster:
    server: https://9.9.9.9
contexts:
- name: other-ctx
  context:
    cluster: other
    user: other-user
current-context: other-ctx
users:
- name: other-user
  user:
    token: other
`
	if err := os.WriteFile(target, []byte(existing), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := Merge(target, incomingKubeconfig, "vks_c-123", false); err != nil {
		t.Fatalf("Merge error: %v", err)
	}

	cfg, _ := Load(target)
	if findContext(cfg, "other-ctx") == nil {
		t.Errorf("existing context was dropped")
	}
	if findContext(cfg, "vks_c-123") == nil {
		t.Errorf("new context missing")
	}
	if cfg.CurrentContext != "other-ctx" {
		t.Errorf("current-context = %q, want other-ctx (setCurrent=false)", cfg.CurrentContext)
	}
}

func TestMergeOverwritesSameContextName(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "config")

	if _, err := Merge(target, incomingKubeconfig, "vks_c-123", true); err != nil {
		t.Fatal(err)
	}
	if _, err := Merge(target, incomingKubeconfig, "vks_c-123", true); err != nil {
		t.Fatal(err)
	}

	cfg, _ := Load(target)
	count := 0
	for _, c := range cfg.Contexts {
		if c.Name == "vks_c-123" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("context vks_c-123 appears %d times, want 1", count)
	}
}
