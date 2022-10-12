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
type strDestroyCmd struct {
	dryRun        bool
	verbose       bool
	step          bool
	inventory     string
	tags          string
	playbookFiles []string
	privateKey    string
	user          string
}

func DestroyCmd() *cobra.Command {
	destroy := &strDestroyCmd{}

	cmd := &cobra.Command{
		Use:          "destroy [flags]",
		Short:        "Install kubernetes cluster, registry",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroy.run()
		},
	}
	// Default value for command struct
	destroy.tags = ""
	destroy.inventory = "./internal/playbooks/koreon-playbook/inventories/inventory-redhat/static-inventory.ini"
	destroy.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/reset.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&destroy.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&destroy.step, "step", "", false, "step")
	f.BoolVarP(&destroy.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&destroy.inventory, "inventory", "i", destroy.inventory, "Specify ansible playbook inventory")
	f.StringVarP(&destroy.privateKey, "private-key", "p", "", "Specify ansible playbook privateKey")
	f.StringVarP(&destroy.user, "user", "u", "", "SSH login user")
	f.StringVar(&destroy.tags, "tags", destroy.tags, "Ansible options tags")

	return cmd
}

func (c *strDestroyCmd) run() error {

	if len(c.playbookFiles) < 1 {
		return fmt.Errorf("[ERROR]: %s", "To run ansible-playbook playbook file path must be specified")
	}

	if len(c.inventory) < 1 {
		return fmt.Errorf("[ERROR]: %s", "To run ansible-playbook an inventory must be specified")
	}

	if len(c.privateKey) < 1 {
		return fmt.Errorf("[ERROR]: %s", "To run ansible-playbook an privateKey must be specified")
	}

	if len(c.user) < 1 {
		return fmt.Errorf("[ERROR]: %s", "To run ansible-playbook an ssh login user must be specified")
	}

	if len(c.tags) > 1 && c.tags == "all" {
		c.tags = ""
	}

	// vars, err := varListToMap(extravars)
	// if err != nil {
	// 	return errors.New("(commandHandler)", "Error parsing extra variables", err)
	// }

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		// Connection: "ssh",
		PrivateKey: c.privateKey,
		User:       c.user,
	}
	// if connectionLocal {
	// 	ansiblePlaybookConnectionOptions.Connection = "local"
	// }

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: c.inventory,
		Verbose:   c.verbose,
		Tags:      c.tags,
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
