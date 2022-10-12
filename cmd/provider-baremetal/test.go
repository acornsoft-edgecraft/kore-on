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
type strTestCmd struct {
	dryRun        bool
	verbose       bool
	step          bool
	inventory     string
	tags          string
	playbookFiles []string
	privateKey    string
	user          string
}

func TestCmd() *cobra.Command {
	test := &strTestCmd{}

	cmd := &cobra.Command{
		Use:          "test [flags]",
		Short:        "Install kubernetes cluster, registry",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return test.run()
		},
	}

	test.inventory = "./internal/playbooks/koreon-playbook/inventories/inventory-redhat/static-inventory.ini"
	test.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/z-test-deploy.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&test.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&test.step, "step", "", false, "step")
	f.BoolVarP(&test.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&test.inventory, "inventory", "i", test.inventory, "Specify ansible playbook inventory")
	f.StringVarP(&test.privateKey, "private-key", "p", "", "Specify ansible playbook privateKey")
	f.StringVarP(&test.user, "user", "u", "", "SSH login user")
	f.StringVar(&test.tags, "tags", "", "Ansible options tags")

	return cmd
}

func (c *strTestCmd) run() error {

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

	// vars, err := varListToMap(extravars)
	// if err != nil {
	// 	return errors.New("(commandHandler)", "Error parsing extra variables", err)
	// }

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		// Connection: "local",
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
