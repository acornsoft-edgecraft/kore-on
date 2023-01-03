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
type strNodeCmd struct {
	dryRun        bool
	verbose       bool
	inventory     string
	tags          string
	playbookFiles []string
	privateKey    string
	user          string
	extravars     map[string]interface{}
}

func NodeCmd() *cobra.Command {
	node := &strNodeCmd{}

	cmd := &cobra.Command{
		Use:          "node [flags]",
		Short:        "Update kubernetes cluster nodes",
		Long:         "",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		addNodeCmd(),
	)
	// SubCommand validation
	utils.CheckCommand(cmd)

	f := cmd.Flags()
	f.BoolVar(&node.verbose, "verbose", false, "verbose")
	f.BoolVarP(&node.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVar(&node.tags, "tags", node.tags, "Ansible options tags")
	f.StringVarP(&node.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&node.user, "user", "u", "", "login user")

	return cmd
}

func addNodeCmd() *cobra.Command {
	addNode := &strNodeCmd{}

	cmd := &cobra.Command{
		Use:          "add [flags]",
		Short:        "Add Node in kubernetes cluster",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addNode.run()
		},
	}

	// Default value for command struct
	addNode.tags = ""
	addNode.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	addNode.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/add-node.yaml",
	}

	f := cmd.Flags()
	f.BoolVar(&addNode.verbose, "verbose", false, "verbose")
	f.BoolVarP(&addNode.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVar(&addNode.tags, "tags", addNode.tags, "Ansible options tags")
	f.StringVarP(&addNode.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&addNode.user, "user", "u", "", "login user")

	return cmd
}

func (c *strNodeCmd) run() error {
	koreOnConfigFileName := viper.GetString("KoreOn.KoreOnConfigFile")
	koreOnConfigFilePath := utils.IskoreOnConfigFilePath(koreOnConfigFileName)
	koreonToml, value := utils.ValidateKoreonTomlConfig(koreOnConfigFilePath, "node")

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
				results.Prepend("Create Cluster"),
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
