package vserver

import (
	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/cmd/vserver/dhcp"
	"github.com/vngcloud/greennode-cli/cmd/vserver/flavor"
	"github.com/vngcloud/greennode-cli/cmd/vserver/floatingip"
	"github.com/vngcloud/greennode-cli/cmd/vserver/image"
	"github.com/vngcloud/greennode-cli/cmd/vserver/networkinterface"
	"github.com/vngcloud/greennode-cli/cmd/vserver/placementgroup"
	"github.com/vngcloud/greennode-cli/cmd/vserver/secgroup"
	"github.com/vngcloud/greennode-cli/cmd/vserver/server"
	"github.com/vngcloud/greennode-cli/cmd/vserver/sshkey"
	"github.com/vngcloud/greennode-cli/cmd/vserver/subnet"
	"github.com/vngcloud/greennode-cli/cmd/vserver/userimage"
	"github.com/vngcloud/greennode-cli/cmd/vserver/volume"
	"github.com/vngcloud/greennode-cli/cmd/vserver/volumetype"
	"github.com/vngcloud/greennode-cli/cmd/vserver/vpc"
	"github.com/vngcloud/greennode-cli/internal/cli"
)

// VServerCmd is the parent command for all vServer subcommands.
var VServerCmd = &cobra.Command{
	Use:   "vserver",
	Short: "VNG Virtual Server (vServer) commands",
	Long:  "Manage vServer instances and related resources.",
	// Reject unknown subcommands (nested groups don't error by default in cobra).
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	VServerCmd.AddCommand(server.ServerCmd)
	VServerCmd.AddCommand(volume.VolumeCmd)
	VServerCmd.AddCommand(vpc.VpcCmd)
	VServerCmd.AddCommand(subnet.SubnetCmd)
	VServerCmd.AddCommand(secgroup.SecgroupCmd)
	VServerCmd.AddCommand(flavor.FlavorCmd)
	VServerCmd.AddCommand(volumetype.VolumeTypeCmd)
	VServerCmd.AddCommand(image.ImageCmd)
	VServerCmd.AddCommand(sshkey.SSHKeyCmd)
	VServerCmd.AddCommand(placementgroup.PlacementGroupCmd)
	VServerCmd.AddCommand(userimage.UserImageCmd)
	VServerCmd.AddCommand(floatingip.FloatingIPCmd)
	VServerCmd.AddCommand(networkinterface.NetworkInterfaceCmd)
	VServerCmd.AddCommand(dhcp.DhcpCmd)

	registerCompletions()
	cli.RegisterService(VServerCmd)
}
