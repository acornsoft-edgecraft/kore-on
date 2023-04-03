package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"kore-on/cmd/koreonctl/conf/templates"
	"kore-on/pkg/logger"
	"kore-on/pkg/model"
	"kore-on/pkg/utils"
	"os"
	"text/template"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Commands structure
type strDestroyCmd struct {
	dryRun        bool
	verbose       bool
	inventory     string
	tags          string
	playbookFiles []string
	privateKey    string
	user          string
	extravars     map[string]interface{}
}

func DestroyCmd() *cobra.Command {
	destroy := &strDestroyCmd{}

	cmd := &cobra.Command{
		Use:          "destroy [flags]",
		Short:        "Delete kubernetes cluster, registry",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroy.run()
		},
	}

	// SubCommand add
	cmd.AddCommand(
		destroyPrepareAirGapCmd(),
		destroyClusterCmd(),
		destroyRegistryCmd(),
		destroyStorageCmd(),
	)

	// SubCommand validation
	utils.CheckCommand(cmd)

	// Default value for command struct
	destroy.tags = "reset-all"
	destroy.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	destroy.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/reset.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&destroy.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&destroy.dryRun, "dry-run", "d", false, "dryRun")
	// f.StringVarP(&destroy.inventory, "inventory", "i", destroy.inventory, "Specify ansible playbook inventory")
	f.StringVar(&destroy.tags, "tags", destroy.tags, "Ansible options tags")
	f.StringVarP(&destroy.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&destroy.user, "user", "u", "", "login user")

	return cmd
}

func destroyPrepareAirGapCmd() *cobra.Command {
	destroyPrepareAirGapCmd := &strDestroyCmd{}

	cmd := &cobra.Command{
		Use:          "prepare-airgap [flags]",
		Short:        "Destroy prepare-airgap",
		Long:         "This command deletes the registry of the prepare-airgap host and deletes related directories.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroyPrepareAirGapCmd.run()
		},
	}

	destroyPrepareAirGapCmd.tags = "reset-prepare-airgap"
	destroyPrepareAirGapCmd.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	destroyPrepareAirGapCmd.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/reset.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&destroyPrepareAirGapCmd.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&destroyPrepareAirGapCmd.dryRun, "dry-run", "d", false, "dryRun")
	// f.StringVarP(&destroyPrepareAirGapCmd.inventory, "inventory", "i", destroyPrepareAirGapCmd.inventory, "Specify ansible playbook inventory")
	f.StringVar(&destroyPrepareAirGapCmd.tags, "tags", destroyPrepareAirGapCmd.tags, "Ansible options tags")
	f.StringVarP(&destroyPrepareAirGapCmd.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&destroyPrepareAirGapCmd.user, "user", "u", "", "login user")

	return cmd
}

func destroyClusterCmd() *cobra.Command {
	destroyClusterCmd := &strDestroyCmd{}

	cmd := &cobra.Command{
		Use:          "cluster [flags]",
		Short:        "Destroy cluster",
		Long:         "This command only deletes the Kubernetes cluster.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroyClusterCmd.run()
		},
	}

	destroyClusterCmd.tags = "reset-cluster"
	destroyClusterCmd.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	destroyClusterCmd.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/reset.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&destroyClusterCmd.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&destroyClusterCmd.dryRun, "dry-run", "d", false, "dryRun")
	// f.StringVarP(&destroyClusterCmd.inventory, "inventory", "i", destroyClusterCmd.inventory, "Specify ansible playbook inventory")
	f.StringVar(&destroyClusterCmd.tags, "tags", destroyClusterCmd.tags, "Ansible options tags")
	f.StringVarP(&destroyClusterCmd.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&destroyClusterCmd.user, "user", "u", "", "login user")

	return cmd
}

func destroyRegistryCmd() *cobra.Command {
	destroyRegistryCmd := &strDestroyCmd{}

	cmd := &cobra.Command{
		Use:          "registry [flags]",
		Short:        "Destroy registry",
		Long:         "This command deletes the installed registry(harbor) and deletes related services and directories.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroyRegistryCmd.run()
		},
	}

	destroyRegistryCmd.tags = "reset-registry"
	destroyRegistryCmd.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	destroyRegistryCmd.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/reset.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&destroyRegistryCmd.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&destroyRegistryCmd.dryRun, "dry-run", "d", false, "dryRun")
	// f.StringVarP(&destroyRegistryCmd.inventory, "inventory", "i", destroyRegistryCmd.inventory, "Specify ansible playbook inventory")
	f.StringVar(&destroyRegistryCmd.tags, "tags", destroyRegistryCmd.tags, "Ansible options tags")
	f.StringVarP(&destroyRegistryCmd.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&destroyRegistryCmd.user, "user", "u", "", "login user")

	return cmd
}

func destroyStorageCmd() *cobra.Command {
	destroyStorageCmd := &strDestroyCmd{}

	cmd := &cobra.Command{
		Use:          "storage [flags]",
		Short:        "Destroy storage",
		Long:         "This command deletes the installed storage(NFS) and deletes related services and directories.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroyStorageCmd.run()
		},
	}

	destroyStorageCmd.tags = "reset-storage"
	destroyStorageCmd.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	destroyStorageCmd.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/reset.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&destroyStorageCmd.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&destroyStorageCmd.dryRun, "dry-run", "d", false, "dryRun")
	// f.StringVarP(&destroyStorageCmd.inventory, "inventory", "i", destroyStorageCmd.inventory, "Specify ansible playbook inventory")
	f.StringVar(&destroyStorageCmd.tags, "tags", destroyStorageCmd.tags, "Ansible options tags")
	f.StringVarP(&destroyStorageCmd.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&destroyStorageCmd.user, "user", "u", "", "login user")

	return cmd
}

func (c *strDestroyCmd) run() error {
	koreOnConfigFileName := viper.GetString("KoreOn.KoreOnConfigFile")
	koreOnConfigFilePath := utils.IskoreOnConfigFilePath(koreOnConfigFileName)
	koreonToml, value := utils.ValidateKoreonTomlConfig(koreOnConfigFilePath, c.tags)

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

	// Make provision data
	data := model.KoreonctlText{}
	data.KoreOnTemp = koreonToml
	data.Command = c.tags
	koreonToml.KoreOn.CommandMode = "destroy"

	// Processing template
	var textVar string
	switch data.Command {
	case "reset-all":
		textVar = templates.DestroyAllText
	case "reset-cluster":
		textVar = templates.DestroyClusterText
	case "reset-registry":
		textVar = templates.DestroyRegistryText
	case "reset-storage":
		textVar = templates.DestroyStorageText
	case "reset-prepare-airgap":
		textVar = templates.DestroyPrepareAirgapText
		koreonToml.KoreOn.ClosedNetwork = false
	}

	koreonctlText := template.New("DestroyText")
	temp, err := koreonctlText.Parse(textVar)
	if err != nil {
		logger.Errorf("Template has errors. cause(%s)", err.Error())
		return err
	}

	// TODO: 진행상황을 어떻게 클라이언트에 보여줄 것인가?
	var buff bytes.Buffer
	err = temp.Execute(&buff, data)
	if err != nil {
		logger.Errorf("Template execution failed. cause(%s)", err.Error())
		return err
	}

	if !utils.CheckUserInput(buff.String(), "y") {
		fmt.Println("nothing to changed. exit")
		os.Exit(1)
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

	if len(c.tags) > 1 && c.tags == "all" {
		c.tags = ""
	}

	utils.ValidateKoreonTomlConfig(koreOnConfigFilePath, c.tags)

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
				results.Prepend("Destroy Cluster"),
			),
		),
	}

	options.AnsibleForceColor()

	err = playbook.Run(context.TODO())
	if err != nil {
		return err
	}

	return nil
}
