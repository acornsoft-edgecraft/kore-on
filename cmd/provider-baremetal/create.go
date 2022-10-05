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
type strCreateCmd struct {
	dryRun        bool
	verbose       bool
	step          bool
	inventory     string
	playbookFiles []string
	privateKey    string
}

func CreateCmd() *cobra.Command {
	create := &strCreateCmd{}

	cmd := &cobra.Command{
		Use:          "create [flags]",
		Short:        "Install kubernetes cluster, registry",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create.run()
		},
	}

	create.inventory = "./internal/playbooks/cubescripts/inventories/inventory.ini"
	create.playbookFiles = []string{
		"./internal/playbooks/cubescripts/cluster.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&create.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&create.step, "step", "", false, "step")
	f.BoolVarP(&create.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&create.privateKey, "private-key", "p", "", "Specify ansible playbook privateKey")

	return cmd
}

func (c *strCreateCmd) run() error {

	if len(c.playbookFiles) < 1 {
		return fmt.Errorf("[ERROR]: %s", "To run ansible-playbook playbook file path must be specified")
	}

	if len(c.inventory) < 1 {
		return fmt.Errorf("[ERROR]: %s", "To run ansible-playbook an inventory must be specified")
	}

	// vars, err := varListToMap(extravars)
	// if err != nil {
	// 	return errors.New("(commandHandler)", "Error parsing extra variables", err)
	// }

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		// Connection: "ssh",
		PrivateKey: "/Users/dongmook/DEV_WORKS/cert_ssh/rhel/cloud-user",
		User:       "cloud-user",
	}
	// if connectionLocal {
	// 	ansiblePlaybookConnectionOptions.Connection = "local"
	// }

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: c.inventory,
		Verbose:   c.verbose,
	}

	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:         c.playbookFiles,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
		Exec: execute.NewDefaultExecute(
			execute.WithTransformers(
				results.Prepend("cobra-cmd-ansibleplaybook example"),
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
