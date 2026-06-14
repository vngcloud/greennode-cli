package cli

import "github.com/spf13/cobra"

// services holds top-level product service commands registered at init time.
var services []*cobra.Command

// RegisterService records a product's top-level cobra command. Each product
// package calls this from its init(); root wiring iterates Services() to mount
// them, so adding a product never requires editing root.go.
func RegisterService(cmd *cobra.Command) {
	services = append(services, cmd)
}

// Services returns a copy of the registered service commands in registration
// order. A copy is returned so callers cannot mutate the registry's backing slice.
func Services() []*cobra.Command {
	out := make([]*cobra.Command, len(services))
	copy(out, services)
	return out
}
