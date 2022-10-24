package cmd

import (
	"context"
	"fmt"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"github.com/spf13/cobra"
)

// Commands structure
type strApplyCmd struct {
	dryRun        bool
	verbose       bool
	step          bool
	inventory     string
	playbookFiles []string
}

func ApplyCmd() *cobra.Command {
	apply := &strApplyCmd{}

	cmd := &cobra.Command{
		Use:          "apply [flags]",
		Short:        "Upgrade kubernetes cluster, registry",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return apply.run()
		},
	}

	apply.inventory = "./internal/playbooks/cubescripts/inventories/inventory.ini"
	apply.playbookFiles = []string{
		"./internal/playbooks/cubescripts/cluster.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&apply.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&apply.step, "step", "", false, "step")
	f.BoolVarP(&apply.dryRun, "dry-run", "d", false, "dryRun")

	return cmd
}

func (c *strApplyCmd) run() error {

	if len(c.playbookFiles) < 1 {
		return fmt.Errorf("[ERROR]: %s", "To run ansible-playbook playbook file path must be specified")
	}

	if len(c.inventory) < 1 {
		return fmt.Errorf("[ERROR]: %s", "To run ansible-playbook an inventory must be specified")
	}

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		// Connection: "ssh",
		PrivateKey: "/Users/dongmook/DEV_WORKS/cert_ssh/acloud/id_rsa",
		User:       "centos",
	}
	// if connectionLocal {
	// 	ansiblePlaybookConnectionOptions.Connection = "local"
	// }

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: c.inventory,
		Verbose:   c.verbose,
	}

	// for keyVar, valueVar := range vars {
	// 	ansiblePlaybookOptions.AddExtraVar(keyVar, valueVar)
	// }

	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:         c.playbookFiles,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
		Exec: execute.NewDefaultExecute(
			execute.WithTransformers(
				results.Prepend("cobra-cmd-ansibleplaybook"),
			),
		),
	}

	options.AnsibleForceColor()

	err := playbook.Run(context.TODO())
	if err != nil {
		return err
	}

	return nil
}
