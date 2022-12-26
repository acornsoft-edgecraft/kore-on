package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"
	"path/filepath"
	"reflect"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/execute/measure"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Commands structure
type strAddonCmd struct {
	dryRun         bool
	verbose        bool
	installHelm    bool
	helmBinaryFile string
	inventory      string
	tags           string
	playbookFiles  []string
	privateKey     string
	user           string
	extravarsFile  map[string]interface{}
	addonExtravars map[string]interface{}
	result         map[string]interface{}
}

func AddonCmd() *cobra.Command {
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
	cmd.AddCommand(AddonDeleteCmd())

	// SubCommand validation
	utils.CheckCommand(cmd)

	// Default value for command struct
	addon.tags = ""
	addon.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	addon.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/add-on.yaml",
	}
	f := cmd.Flags()
	f.BoolVar(&addon.verbose, "verbose", false, "verbose")
	f.BoolVarP(&addon.dryRun, "dry-run", "d", false, "dryRun")
	f.BoolVar(&addon.installHelm, "install-helm", false, "Helm installation options")
	f.StringVar(&addon.helmBinaryFile, "helm-binary-file", "", "helm binary file")
	f.StringVar(&addon.tags, "tags", addon.tags, "Ansible options tags")
	f.StringVarP(&addon.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&addon.user, "user", "u", "", "login user")

	return cmd
}

func AddonDeleteCmd() *cobra.Command {
	addonDelete := &strAddonCmd{}

	cmd := &cobra.Command{
		Use:   "delete [flags]",
		Short: "Deployment Applications in kubernetes cluster",
		Long: "This command deploys the application to Kubernetes.\n" +
			"Use helm as the package manager for Kubernetes.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addonDelete.run()
		},
	}

	// Default value for command struct
	addonDelete.tags = ""
	addonDelete.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	addonDelete.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/delete-add-on.yaml",
	}
	f := cmd.Flags()
	f.BoolVar(&addonDelete.verbose, "verbose", false, "verbose")
	f.BoolVarP(&addonDelete.dryRun, "dry-run", "d", false, "dryRun")
	f.BoolVar(&addonDelete.installHelm, "install-helm", false, "Helm installation options")
	f.StringVar(&addonDelete.helmBinaryFile, "helm-binary-file", "", "helm binary file")
	f.StringVar(&addonDelete.tags, "tags", addonDelete.tags, "Ansible options tags")
	f.StringVarP(&addonDelete.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&addonDelete.user, "user", "u", "", "login user")

	return cmd
}

func (c *strAddonCmd) run() error {
	addonConfigFileName := viper.GetString("Addon.AddonConfigFile")
	addonPath := utils.IskoreOnConfigFilePath(addonConfigFileName)
	addonToml, err := utils.GetAddonTomlConfig(addonPath)
	if err != nil {
		logger.Fatal(err)
	} else {
		// Install Helm
		if c.installHelm {
			addonToml.Addon.HelmInstall = c.installHelm
		}
		if c.helmBinaryFile != "" {
			addonToml.Addon.HelmBinaryFile = c.helmBinaryFile
		}

		// Prompt user for more input
		if addonToml.Apps.CsiDriverNfs.Install {
			id := utils.InputPrompt("# Enter the username for the private registry.\nusername:")
			addonToml.Apps.CsiDriverNfs.ChartRefID = base64.StdEncoding.EncodeToString([]byte(id))

			pw := utils.SensitivePrompt("# Enter the password for the private registry.\npassword:")
			addonToml.Apps.CsiDriverNfs.ChartRefPW = base64.StdEncoding.EncodeToString([]byte(pw))
		}

		addonToml.Addon.HelmVersion = utils.IsSupportVersion("", "SupportHelmVersion")
		if addonToml.Addon.AddonDataDir == "" {
			addonToml.Addon.AddonDataDir = "/data/addon"
		}

		b, err := json.Marshal(addonToml)
		if err != nil {
			logger.Fatal(err)
		}
		if err := json.Unmarshal(b, &c.addonExtravars); err != nil {
			logger.Fatal(err.Error())
		}

		result := make(map[string]interface{})
		// for k, v := range c.extravars {
		// 	if _, ok := c.extravars[k]; ok {
		// 		result[k] = v
		// 	}
		// }
		for k, v := range c.addonExtravars {
			if _, ok := c.addonExtravars[k]; ok {
				result[k] = v
			}
		}
		c.result = result
		// resultFiles := make(map[string]interface{})
		// for k, v := range c.addonExtravars {
		// 	if k == "Apps" {
		// 		for i, j := range v.(map[string]interface{}) {
		// 			if addonToml.Apps[i].Install {

		// 			}
		// 			for t, d := range j.(map[string]interface{}) {
		// 				fmt.Println("==== ", addonToml.Apps.FluentBit)
		// 				if t == "Install" && d != true {
		// 					break
		// 				}
		// 				fmt.Println("====================================-")
		// 				fmt.Println("key == ", i)

		// 				if t == "Values" {
		// 					key := fmt.Sprintf("%s_%s", k, i)
		// 					resultFiles[key] = d
		// 					fmt.Println("values == ", d)
		// 				}
		// 				fmt.Println("xxxxxxxxxxxxx")
		// 			}
		// 		}
		// 	}
		// }

		// setValueFile(addonToml.Apps)
		aaa := getValuesFromInterface(&addonToml.Apps, "FluentBit")
		fmt.Println(aaa)
		logger.Fatal()

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
		Inventory:     c.inventory,
		Verbose:       c.verbose,
		Tags:          c.tags,
		ExtraVars:     c.result,
		ExtraVarsFile: []string{"@internal/playbooks/koreon-playbook/download/test-values.yaml"},
	}

	executorTimeMeasurement := measure.NewExecutorTimeMeasurement(
		execute.NewDefaultExecute(
			execute.WithEnvVar("ANSIBLE_FORCE_COLOR", "true"),
			execute.WithTransformers(
				utils.OutputColored(),
				results.Prepend("Addon deployment in cluster"),
				// results.LogFormat(results.DefaultLogFormatLayout, results.Now),
			),
		),
		// measure.WithShowDuration(),
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

func setValueFile(s interface{}) error {
	// v := reflect.ValueOf(s)
	// vt := v.Type()
	// if !v.CanAddr() {
	// 	return fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	// }

	filename, _ := filepath.Abs("internal/playbooks/koreon-playbook/download/koreboard-values.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Fatal(err)
	}

	var values map[string]interface{}

	err = yaml.Unmarshal(yamlFile, &values)
	if err != nil {
		logger.Fatal(err)
	}
	m3 := map[string]interface{}{
		"aaa": values,
	}

	fmt.Println(m3)

	// val := reflect.ValueOf(s)

	// // If it's an interface or a pointer, unwrap it.
	// if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
	// 	val = val.Elem()
	// } else {
	// 	return fmt.Errorf("[ERROR]: %s", "must be a struct")
	// }
	// valNumFields := val.NumField()

	// for i := 0; i < valNumFields; i++ {
	// 	field := val.Field(i)
	// 	fieldKind := field.Kind()

	// 	// Check if it's a pointer to a struct.
	// 	if fieldKind == reflect.Ptr && field.Elem().Kind() == reflect.Struct {
	// 		if field.CanInterface() {
	// 			// Recurse using an interface of the field.
	// 			fmt.Println("field ==", field.Interface())
	// 		}

	// 		// Move onto the next field.
	// 		continue
	// 	}

	// 	// Check if it's a struct value.
	// 	if fieldKind == reflect.Struct {
	// 		if field.CanAddr() && field.Addr().CanInterface() {
	// 			// Recurse using an interface of the pointer value of the field.
	// 			fmt.Println("field pointer ==", field.Addr().Interface())
	// 		}

	// 		// Move onto the next field.
	// 		continue
	// 	}

	// 	// Check if it's a string or a pointer to a string.
	// 	if fieldKind == reflect.String || (fieldKind == reflect.Ptr && field.Elem().Kind() == reflect.String) {
	// 		typeField := val.Type().Field(i)

	// 		fmt.Println("typeField == ", typeField)

	// 		// Set the string value to the sanitized string if it's allowed.
	// 		// It should always be allowed at this point.
	// 		// if field.CanSet() {
	// 		// 	field.SetString(policy.Sanitize(field.String()))
	// 		// }

	// 		continue
	// 	}
	// }

	return nil
}

// GetValuesFromInterface - 지정된 Interface 형식의 Structure에서 지정한 FIeld 이름을 기준으로 값을 Array로 반환 (using Reflect)
func getValuesFromInterface(val interface{}, fields ...string) []interface{} {
	returnArray := make([]interface{}, 0)
	structVal := reflect.ValueOf(val).Elem()

	for i := 0; i < structVal.NumField(); i++ {
		fmt.Println("=== ", structVal.Field(i))
	}

	for _, name := range fields {
		field := structVal.FieldByName(name).Interface()
		returnArray = append(returnArray, field)
	}
	return returnArray
}
