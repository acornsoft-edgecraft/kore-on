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
	config *ProviderOpenstack
}

type ProviderOpenstack struct {
	LogLevel string
}

func (o *options) run(ctx context.Context) error {
	return nil
}

func NewProviderOpenstackCommand() *cobra.Command {
	opts := &options{}

	cmd := &cobra.Command{
		Use:   Name,
		Short: "Launch the " + Name,
		Long:  Name + " serves webhook endpoints for resources in the garden cluster.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.run(cmd.Context())
		},
	}

	return cmd
}
