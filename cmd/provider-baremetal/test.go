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
	"github.com/apenella/go-ansible/pkg/execute/measure"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Commands structure
type strTestCmd struct {
	dryRun         bool
	verbose        bool
	step           bool
	inventory      string
	tags           string
	playbookFiles  []string
	privateKey     string
	user           string
	extravars      map[string]interface{}
	addonExtravars map[string]interface{}
	result         map[string]interface{}
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

	// SubCommand add
	cmd.AddCommand(emptyCmd())

	// SubCommand validation
	utils.CheckCommand(cmd)

	test.tags = ""
	test.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	test.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/z-test-create-os-image.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&test.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&test.step, "step", "", false, "step")
	f.BoolVarP(&test.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&test.inventory, "inventory", "i", test.inventory, "Specify ansible playbook inventory")
	f.StringVarP(&test.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&test.user, "user", "u", "", "login user")
	f.StringVar(&test.tags, "tags", "", "Ansible options tags")

	return cmd
}

func (c *strTestCmd) run() error {
	koreOnConfigFileName := viper.GetString("KoreOn.KoreOnConfigFile")
	koreOnConfigFilePath := utils.IskoreOnConfigFilePath(koreOnConfigFileName)
	koreonToml, value := utils.ValidateKoreonTomlConfig(koreOnConfigFilePath, "create")
	koreonToml.KoreOn.FileName = koreOnConfigFileName
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
	data.Command = "create"

	// Processing template
	koreonctlText := template.New("CreateText")
	temp, err := koreonctlText.Parse(templates.CreateText)
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

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		PrivateKey: c.privateKey,
		User:       c.user,
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: c.inventory,
		Verbose:   c.verbose,
		Tags:      c.tags,
		ExtraVars: c.extravars,
		// ExtraVars: c.result,
	}

	executorTimeMeasurement := measure.NewExecutorTimeMeasurement(
		execute.NewDefaultExecute(
			execute.WithEnvVar("ANSIBLE_FORCE_COLOR", "true"),
			execute.WithTransformers(
				utils.OutputColored(),
				results.Prepend("cobra-cmd-ansibleplaybook example"),
				// results.LogFormat(results.DefaultLogFormatLayout, results.Now),
			),
		),
		measure.WithShowDuration(),
	)

	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:         c.playbookFiles,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
		Exec:              executorTimeMeasurement,
	}

	options.AnsibleForceColor()

	err = playbook.Run(context.TODO())
	if err != nil {
		return err
	}

	return nil
}
