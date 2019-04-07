package cmds

import (
	"github.com/urfave/cli"
)

type Bootstrap struct {
	DataDir string
	Members cli.StringSlice
}

var BootstrapConfig Bootstrap

func NewBootstrapCommand(action func(*cli.Context) error) cli.Command {
	return cli.Command{
		Name:      "bootstrap",
		Usage:     "Bootstrap a node of a new cluster",
		UsageText: appName + " bootstrap [OPTIONS] [ID] [IP]",
		Action:    action,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "data-dir,d",
				Usage:       "Folder to hold state default /var/lib/rancher/k3s or ${HOME}/.rancher/k3s if not root",
				Destination: &BootstrapConfig.DataDir,
			},
			cli.StringSliceFlag{
				Name:  "members",
				Usage: "Initial members of the cluster",
				Value: &BootstrapConfig.Members,
			},
		},
	}
}
