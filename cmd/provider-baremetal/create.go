package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"
	"os"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Commands structure
type strCreateCmd struct {
	dryRun        bool
	verbose       bool
	step          bool
	inventory     string
	tags          string
	playbookFiles []string
	privateKey    string
	user          string
	extravars     map[string]interface{}
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

	// Default value for command struct
	create.tags = ""
	create.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	create.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/cluster.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&create.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&create.step, "step", "", false, "step")
	f.BoolVarP(&create.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&create.inventory, "inventory", "i", create.inventory, "Specify ansible playbook inventory")
	f.StringVar(&create.tags, "tags", create.tags, "Ansible options tags")
	f.StringVarP(&create.privateKey, "private-key", "p", "", "Specify ansible playbook privateKey")
	f.StringVarP(&create.user, "user", "u", "", "SSH login user")

	return cmd
}

func (c *strCreateCmd) run() error {
	koreOnConfigFileName := viper.GetString("KoreOn.KoreOnConfigFile")
	koreOnConfigFilePath := utils.IskoreOnConfigFilePath(koreOnConfigFileName)
	koreonToml, value := utils.ValidateKoreonTomlConfig(koreOnConfigFilePath, "create")

	if value {
		b, err := json.Marshal(koreonToml)
		if err != nil {
			logger.Fatal(err)
			os.Exit(1)
		}
		if err := json.Unmarshal(b, &c.extravars); err != nil {
			logger.Fatal(err.Error())
			os.Exit(1)
		}
	}

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

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		PrivateKey: c.privateKey,
		User:       c.user,
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: c.inventory,
		Verbose:   c.verbose,
		Tags:      c.tags,
		ExtraVars: c.extravars,
	}

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
