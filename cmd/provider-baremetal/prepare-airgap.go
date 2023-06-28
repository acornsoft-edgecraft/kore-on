package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"kore-on/cmd/koreonctl/conf"
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
)

// Commands structure
type strAirGapCmd struct {
	dryRun        bool
	verbose       bool
	inventory     string
	tags          string
	playbookFiles []string
	privateKey    string
	user          string
	command       string
	extravars     map[string]interface{}
}

func AirGapCmd() *cobra.Command {
	prepareAirgap := &strAirGapCmd{}

	cmd := &cobra.Command{
		Use:          "prepare-airgap [flags]",
		Short:        "Preparing a kubernetes cluster and registry for AirGap network",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return prepareAirgap.run()
		},
	}

	// SubCommand add
	cmd.AddCommand(
		DownLoadArchiveCmd(),
		ImageUploadCmd(),
	)

	// SubCommand validation
	utils.CheckCommand(cmd)

	// Default value for command struct
	prepareAirgap.tags = ""
	prepareAirgap.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	prepareAirgap.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/prepare-airgap.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&prepareAirgap.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&prepareAirgap.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVar(&prepareAirgap.tags, "tags", prepareAirgap.tags, "Ansible options tags")
	f.StringVarP(&prepareAirgap.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&prepareAirgap.user, "user", "u", "", "login user")

	return cmd
}

func DownLoadArchiveCmd() *cobra.Command {
	downLoadArchive := &strAirGapCmd{}

	cmd := &cobra.Command{
		Use:          "download-archive [flags]",
		Short:        "Download archive files to localhost",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return downLoadArchive.run()
		},
	}

	// Default value for command struct
	downLoadArchive.tags = ""
	downLoadArchive.command = "downLoad-archive"
	downLoadArchive.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	downLoadArchive.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/download-archive-to-local.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&downLoadArchive.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&downLoadArchive.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVar(&downLoadArchive.tags, "tags", downLoadArchive.tags, "Ansible options tags")
	f.StringVarP(&downLoadArchive.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&downLoadArchive.user, "user", "u", "", "login user")

	return cmd
}

func ImageUploadCmd() *cobra.Command {
	imageUpload := &strAirGapCmd{}

	cmd := &cobra.Command{
		Use:          "image-upload [flags]",
		Short:        "Images Pull and Push to private registry",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return imageUpload.run()
		},
	}

	// Default value for command struct
	imageUpload.tags = ""
	imageUpload.command = "image-upload"
	imageUpload.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	imageUpload.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/prepare-airgap-pull-push-image.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&imageUpload.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&imageUpload.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVar(&imageUpload.tags, "tags", imageUpload.tags, "Ansible options tags")
	f.StringVarP(&imageUpload.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&imageUpload.user, "user", "u", "", "login user")

	return cmd
}

func (c *strAirGapCmd) run() error {
	koreOnConfigFileName := conf.KoreOnConfigFile
	koreOnConfigFilePath := utils.IskoreOnConfigFilePath(koreOnConfigFileName)
	koreonToml, errBool := utils.ValidateKoreonTomlConfig(koreOnConfigFilePath, "prepare-airgap")
	if !errBool {
		message := "Settings are incorrect. Please check the 'korean.toml' file!!"
		logger.Fatal(fmt.Errorf("%s", message))
	}

	// Make provision data
	data := model.KoreonctlText{}
	data.KoreOnTemp = koreonToml
	data.Command = "prepare-airgap"

	// Processing template
	koreonctlText := template.New("PrepareAirgapText")
	temp, err := koreonctlText.Parse(templates.PrepareAirgapText)
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

	// if c.command == "" {
	// 	// Prompt login
	// 	id := utils.InputPrompt("\n## To helm chart pull csi-driver-nfs, you need to login as a private repository (Helm Chart) user.\nusername:")
	// 	pw := utils.SensitivePrompt("password:")
	// 	koreonToml.KoreOn.HelmChartProject = viper.GetString("KoreOn.HelmChartProject")
	// 	koreonToml.KoreOn.HelmCubeRepoID = base64.StdEncoding.EncodeToString([]byte(id))
	// 	koreonToml.KoreOn.HelmCubeRepoPW = base64.StdEncoding.EncodeToString([]byte(pw))

	// 	commandArgs := "helm registry login " + koreonToml.KoreOn.HelmCubeRepoUrl +
	// 		" --username " + id +
	// 		" --password " + pw

	// 	err = checkHelmRepoLogin(id, pw, commandArgs)
	// 	if err != nil {
	// 		str := fmt.Sprintf("%s", err)
	// 		fi := strings.Index(str, "Error")
	// 		li := strings.LastIndex(str, "\"")
	// 		err = fmt.Errorf(str[fi : li+1])
	// 		logger.Fatal(err)
	// 	} else {
	// 		fmt.Println("Login Succeeded!!")
	// 	}
	// }

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

	b, err := json.Marshal(koreonToml)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
	if err := json.Unmarshal(b, &c.extravars); err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
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
				results.Prepend("Prepare AirGap"),
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
