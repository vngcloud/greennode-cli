package cli

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRegisterServiceAccumulates(t *testing.T) {
	saved := services
	services = nil
	defer func() { services = saved }()

	a := &cobra.Command{Use: "a"}
	b := &cobra.Command{Use: "b"}
	RegisterService(a)
	RegisterService(b)

	got := Services()
	if len(got) != 2 || got[0] != a || got[1] != b {
		t.Errorf("Services() = %v, want [a b] in registration order", got)
	}
}
