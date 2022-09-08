package app

import (
	"context"

	"github.com/spf13/cobra"
)

const Name = "koreon-provider"

type options struct {
	// configFile is the location of the Gardener controller manager's configuration file.
	configFile string

	// config is the decoded admission controller config.
	config *ProviderBaremetal
}

type ProviderBaremetal struct {
	LogLevel string
}

func (o *options) run(ctx context.Context) error {
	return nil
}

func NewProviderBaremetalCommand() *cobra.Command {
	opts := &options{}

	cmd := &cobra.Command{
		Use:   Name,
		Short: "Launch the " + Name,
		Long:  Name + " automates k8s installation tasks for baremetal.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.run(cmd.Context())
		},
	}

	return cmd
}
