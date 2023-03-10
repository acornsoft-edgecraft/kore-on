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

type strAddonCmd struct {
	dryRun         bool
	verbose        bool
	privateKey     string
	user           string
	installHelm    bool
	helmBinaryFile string
}

func addonCmd() *cobra.Command {
	addon := &strAddonCmd{}

	cmd := &cobra.Command{
		Use:   "addon [flags]",
		Short: "Deployment Applications in kubernetes cluster",
		Long: "This command deploys the application to Kubernetes.\n" +
			"Use helm as the package manager for Kubernetes.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addon.run()
		},
	}

	// SubCommand add
	cmd.AddCommand(addonInitCmd())

	// SubCommand validation
	utils.CheckCommand(cmd)

	f := cmd.Flags()
	f.BoolVar(&addon.verbose, "vvv", false, "verbose")
	f.BoolVarP(&addon.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&addon.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&addon.user, "user", "u", "", "login user")
	f.BoolVar(&addon.installHelm, "install-helm", false, "Helm installation options")
	f.StringVar(&addon.helmBinaryFile, "helm-binary-file", "", "helm binary file")

	return cmd
}

func (c *strAddonCmd) run() error {

	workDir, _ := os.Getwd()
	var err error = nil
	logger.Infof("Start deployment for k8s cluster")

	if err = c.addon(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strAddonCmd) addon(workDir string) error {
	// Doker check
	utils.CheckDocker()

	koreonImageName := viper.GetString("KoreOn.KoreOnImageName")
	koreOnImage := viper.GetString("KoreOn.KoreOnImage")
	koreOnConfigFilePath := viper.GetString("KoreOn.KoreOnConfigFileSubDir")
	// koreonImageName := conf.KoreOnImageName
	// koreOnImage := conf.KoreOnImage
	// koreOnConfigFilePath := conf.KoreOnConfigFileSubDir

	addonToml, err := utils.GetAddonTomlConfig(workDir + "/" + viper.GetString("Addon.AddOnConfigFile"))
	if err != nil {
		logger.Fatal(err)
	}

	commandArgs := []string{
		"docker",
		"run",
		"--rm",
		"--privileged",
		"-it",
	}

	if !addonToml.Addon.ClosedNetwork {
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
		"addon",
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
