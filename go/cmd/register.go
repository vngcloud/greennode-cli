package cmd

// Product services self-register via their package init(). Blank-importing them
// here is the ONLY change needed to mount a new product CLI — root.go iterates
// the registry and never needs editing.
import (
	_ "github.com/vngcloud/greennode-cli/cmd/vks"
	_ "github.com/vngcloud/greennode-cli/internal/resources/vserver"
	// New products add a line here, e.g.:
	// _ "github.com/vngcloud/greennode-cli/cmd/vserver"
)
