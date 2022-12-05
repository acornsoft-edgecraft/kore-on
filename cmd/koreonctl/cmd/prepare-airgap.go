package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"kore-on/pkg/logger"
	"kore-on/pkg/model"
	"kore-on/pkg/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"kore-on/cmd/koreonctl/conf"
	"kore-on/cmd/koreonctl/conf/templates"

	"github.com/spf13/cobra"
)

type strAirGapCmd struct {
	dryRun     bool
	verbose    bool
	privateKey string
	user       string
	command    string
}

func airGapCmd() *cobra.Command {
	prepareAirgap := &strAirGapCmd{}

	cmd := &cobra.Command{
		Use:          "prepare-airgap [flags]",
		Short:        "Preparing a kubernetes cluster and registry for Air gap network",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return prepareAirgap.run()
		},
	}

	cmd.AddCommand(downLoadArchiveCmd())

	f := cmd.Flags()
	f.BoolVarP(&prepareAirgap.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&prepareAirgap.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&prepareAirgap.user, "user", "u", "", "login user")

	return cmd
}

func downLoadArchiveCmd() *cobra.Command {
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

	downLoadArchive.command = "download-archive"

	f := cmd.Flags()
	f.BoolVarP(&downLoadArchive.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&downLoadArchive.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&downLoadArchive.user, "user", "u", "", "login user")

	return cmd
}

func (c *strAirGapCmd) run() error {

	//if !utils.CheckUserInput("Do you really want to create? Only 'yes' will be accepted to confirm: ", "yes") {
	//	fmt.Println("nothing to changed. exit")
	//	os.Exit(1)
	//}

	workDir, _ := os.Getwd()
	var err error = nil
	logger.Infof("Start provisioning for preparing a kubernetes cluster and registry")

	if err = c.airgap(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strAirGapCmd) airgap(workDir string) error {
	// Doker check
	utils.CheckDocker()

	koreonImageName := conf.KoreOnImageName
	koreOnImage := conf.KoreOnImage
	koreOnConfigFileName := conf.KoreOnConfigFile
	koreOnConfigFilePath := conf.KoreOnConfigFileSubDir

	koreonToml, err := utils.GetKoreonTomlConfig(workDir + "/" + koreOnConfigFileName)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
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

	commandArgs := []string{
		"docker",
		"run",
		"--pull",
		"--privileged",
		"-it",
	}

	if !koreonToml.KoreOn.ClosedNetwork {
		commandArgs = append(commandArgs, "--pull")
		commandArgs = append(commandArgs, "always")
	}

	commandArgsVol := []string{
		"-v",
		fmt.Sprintf("%s:%s", workDir, "/"+koreOnConfigFilePath),
	}

	commandArgsKoreonctl := []string{
		koreOnImage,
		"./" + koreonImageName,
		"prepare-airgap",
	}

	// docker commands
	if c.privateKey != "" {
		key := filepath.Base(c.privateKey)
		keyPath, _ := filepath.Abs(c.privateKey)
		commandArgsVol = append(commandArgsVol, "--mount")
		commandArgsVol = append(commandArgsVol, fmt.Sprintf("type=bind,source=%s,target=/home/%s,readonly", keyPath, key))
	}

	//- koreonctl commands
	if c.command == "download-archive" {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "download-archive")
	}

	if c.verbose {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--verbose")
	}

	if c.dryRun {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--dry-run")
	}

	if c.privateKey != "" {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--private-key")
		key := filepath.Base(c.privateKey)
		commandArgsKoreonctl = append(commandArgsKoreonctl, "/home/"+key)
	} else {
		logger.Fatal(fmt.Errorf("[ERROR]: %s", "To run ansible-playbook an privateKey must be specified"))
	}

	if c.user != "" {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--user")
		commandArgsKoreonctl = append(commandArgsKoreonctl, c.user)
	} else {
		logger.Fatal(fmt.Errorf("[ERROR]: %s", "To run ansible-playbook an ssh login user must be specified"))
	}
	//-end koreonctl commands

	commandArgs = append(commandArgs, commandArgsVol...)
	commandArgs = append(commandArgs, commandArgsKoreonctl...)

	binary, lookErr := exec.LookPath("docker")
	if lookErr != nil {
		logger.Fatal(lookErr)
	}

	err = syscall.Exec(binary, commandArgs, os.Environ())
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	return nil
}
