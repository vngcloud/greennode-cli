// Package kubeconfig loads, merges, and writes Kubernetes kubeconfig files.
package kubeconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config models the subset of a kubeconfig we manipulate. Cluster/User payloads
// are kept as raw yaml.Node values so unknown fields survive a round-trip.
type Config struct {
	APIVersion     string         `yaml:"apiVersion,omitempty"`
	Kind           string         `yaml:"kind,omitempty"`
	Clusters       []NamedEntry   `yaml:"clusters"`
	Contexts       []NamedEntry   `yaml:"contexts"`
	Users          []NamedEntry   `yaml:"users"`
	CurrentContext string         `yaml:"current-context,omitempty"`
	Preferences    map[string]any `yaml:"preferences,omitempty"`
}

// NamedEntry is a generic {name, <payload>} entry. The payload key differs by
// list (cluster/context/user); each is optional so the struct serves all three.
type NamedEntry struct {
	Name string `yaml:"name"`
	// Cluster and User are value-type yaml.Node (not pointers). yaml.v3 only
	// populates a node's content when decoding into an addressable yaml.Node
	// value; decoding into a *yaml.Node field leaves the node empty, silently
	// dropping the payload. omitempty skips a zero-kind node, so the same
	// struct serves cluster, context, and user lists.
	Cluster yaml.Node    `yaml:"cluster,omitempty"`
	Context *contextBody `yaml:"context,omitempty"`
	User    yaml.Node    `yaml:"user,omitempty"`
}

type contextBody struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

// MergeResult reports what was applied.
type MergeResult struct {
	ContextName string
	Path        string
}

// Load reads a kubeconfig from disk. A missing file yields an empty Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Config{APIVersion: "v1", Kind: "Config"}, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg Config
	if len(data) > 0 {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse kubeconfig %s: %w", path, err)
		}
	}
	if cfg.APIVersion == "" {
		cfg.APIVersion = "v1"
	}
	if cfg.Kind == "" {
		cfg.Kind = "Config"
	}
	return &cfg, nil
}

// Write serializes cfg to path with 0600 perms, creating parent dirs (0700).
func Write(path string, cfg *Config) error {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
	}
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o600)
}

// Merge parses incoming (a full kubeconfig YAML string) and grafts its first
// cluster/context/user into the file at target under the name contextName.
// If setCurrent is true, current-context is set to contextName.
func Merge(target, incoming, contextName string, setCurrent bool) (*MergeResult, error) {
	var src Config
	if err := yaml.Unmarshal([]byte(incoming), &src); err != nil {
		return nil, fmt.Errorf("failed to parse incoming kubeconfig: %w", err)
	}
	if len(src.Clusters) == 0 || len(src.Contexts) == 0 || len(src.Users) == 0 {
		return nil, fmt.Errorf("incoming kubeconfig is missing cluster/context/user entries")
	}

	dst, err := Load(target)
	if err != nil {
		return nil, err
	}

	clusterName := "cluster_" + contextName
	userName := "user_" + contextName

	cluster := src.Clusters[0]
	cluster.Name = clusterName
	user := src.Users[0]
	user.Name = userName
	ctx := NamedEntry{
		Name:    contextName,
		Context: &contextBody{Cluster: clusterName, User: userName},
	}

	dst.Clusters = upsert(dst.Clusters, cluster)
	dst.Users = upsert(dst.Users, user)
	dst.Contexts = upsert(dst.Contexts, ctx)
	if setCurrent {
		dst.CurrentContext = contextName
	}

	if err := Write(target, dst); err != nil {
		return nil, err
	}
	return &MergeResult{ContextName: contextName, Path: target}, nil
}

// upsert replaces an entry with the same name or appends it.
func upsert(list []NamedEntry, e NamedEntry) []NamedEntry {
	for i := range list {
		if list[i].Name == e.Name {
			list[i] = e
			return list
		}
	}
	return append(list, e)
}

// findContext is a helper/accessor used in tests.
func findContext(cfg *Config, name string) *NamedEntry {
	for i := range cfg.Contexts {
		if cfg.Contexts[i].Name == name {
			return &cfg.Contexts[i]
		}
	}
	return nil
}
