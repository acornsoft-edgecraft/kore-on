package cmd

import (
	"fmt"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type strCreateCmd struct {
	dryRun     bool
	verbose    bool
	privateKey string
	user       string
}

func createCmd() *cobra.Command {
	create := &strCreateCmd{}

	cmd := &cobra.Command{
		Use:          "create [flags]",
		Short:        "Install kubernetes cluster, registry",
		Long:         "This command installs the Kubernetes cluster and registry.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create.run()
		},
	}

	// SubCommand add
	cmd.AddCommand(emptyCmd())

	// SubCommand validation
	utils.CheckCommand(cmd)

	f := cmd.Flags()
	f.BoolVar(&create.verbose, "vvv", false, "verbose")
	f.BoolVarP(&create.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&create.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&create.user, "user", "u", "", "login user")

	return cmd
}

func (c *strCreateCmd) run() error {

	workDir, _ := os.Getwd()
	var err error = nil
	logger.Infof("Start provisioning for cloud infrastructure")

	if err = c.create(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strCreateCmd) create(workDir string) error {
	// Doker check
	utils.CheckDocker()

	koreonImageName := viper.GetString("KoreOn.KoreOnImageName")
	koreOnImage := viper.GetString("KoreOn.KoreOnImage")
	koreOnConfigFileName := viper.GetString("KoreOn.KoreOnConfigFile")
	koreOnConfigFilePath := viper.GetString("KoreOn.KoreOnConfigFileSubDir")

	koreonToml, err := utils.GetKoreonTomlConfig(workDir + "/" + koreOnConfigFileName)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	// // Make provision data
	// data := model.KoreonctlText{}
	// data.KoreOnTemp = koreonToml
	// data.Command = "create"

	// // Processing template
	// koreonctlText := template.New("CreateText")
	// temp, err := koreonctlText.Parse(templates.CreateText)
	// if err != nil {
	// 	logger.Errorf("Template has errors. cause(%s)", err.Error())
	// 	return err
	// }

	// // TODO: 진행상황을 어떻게 클라이언트에 보여줄 것인가?
	// var buff bytes.Buffer
	// err = temp.Execute(&buff, data)
	// if err != nil {
	// 	logger.Errorf("Template execution failed. cause(%s)", err.Error())
	// 	return err
	// }

	// if !utils.CheckUserInput(buff.String(), "y") {
	// 	fmt.Println("nothing to changed. exit")
	// 	os.Exit(1)
	// }

	commandArgs := []string{
		"docker",
		"run",
		"--rm",
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
		"create",
	}

	if c.privateKey != "" {
		key := filepath.Base(c.privateKey)
		keyPath, _ := filepath.Abs(c.privateKey)
		commandArgsVol = append(commandArgsVol, "--mount")
		commandArgsVol = append(commandArgsVol, fmt.Sprintf("type=bind,source=%s,target=/home/%s,readonly", keyPath, key))
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
