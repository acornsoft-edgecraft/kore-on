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
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Commands structure
type strClusterUpdateCmd struct {
	dryRun        bool
	verbose       bool
	inventory     string
	tags          string
	playbookFiles []string
	privateKey    string
	user          string
	command       string
	kubeconfig    string
	extravars     map[string]interface{}
}

var err error

func ClusterUpdateCmd() *cobra.Command {
	clusterUpdate := &strClusterUpdateCmd{}

	cmd := &cobra.Command{
		Use:          "update [flags]",
		Short:        "Update kubernetes cluster(node scale in/out)",
		Long:         "This command update the Kubernetes cluster nodes (node scale in/out)",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return clusterUpdate.run()
		},
	}

	// SubCommand add
	cmd.AddCommand(
		GetKubeConfigCmd(),
	)

	// SubCommand validation
	utils.CheckCommand(cmd)

	// Default value for command struct
	clusterUpdate.command = "update"
	clusterUpdate.tags = ""
	clusterUpdate.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	clusterUpdate.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/cluster-update.yaml",
	}

	f := cmd.Flags()
	f.BoolVar(&clusterUpdate.verbose, "verbose", false, "verbose")
	f.BoolVarP(&clusterUpdate.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVar(&clusterUpdate.tags, "tags", clusterUpdate.tags, "Ansible options tags")
	f.StringVarP(&clusterUpdate.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&clusterUpdate.user, "user", "u", "", "login user")
	f.StringVar(&clusterUpdate.kubeconfig, "kubeconfig", "", "get kubeconfig")

	return cmd
}

func GetKubeConfigCmd() *cobra.Command {
	getKubeConfig := &strClusterUpdateCmd{}

	cmd := &cobra.Command{
		Use:          "get-kubeconfig [flags]",
		Short:        "Get Kubeconfig file",
		Long:         "This command get kubeconfig file in k8s controlplane node.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getKubeConfig.run()
		},
	}

	getKubeConfig.command = "get-kubeconfig"
	getKubeConfig.tags = ""
	getKubeConfig.inventory = "./internal/playbooks/koreon-playbook/inventory/inventory.ini"
	getKubeConfig.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/cluster-get.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&getKubeConfig.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&getKubeConfig.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&getKubeConfig.inventory, "inventory", "i", getKubeConfig.inventory, "Specify ansible playbook inventory")
	f.StringVar(&getKubeConfig.tags, "tags", getKubeConfig.tags, "Ansible options tags")
	f.StringVarP(&getKubeConfig.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&getKubeConfig.user, "user", "u", "", "login user")

	return cmd
}

func (c *strClusterUpdateCmd) run() error {
	koreOnConfigFileName := viper.GetString("KoreOn.KoreOnConfigFile")
	koreOnConfigFilePath := utils.IskoreOnConfigFilePath(koreOnConfigFileName)
	koreonToml, errBool := utils.ValidateKoreonTomlConfig(koreOnConfigFilePath, "cluster-update")
	if !errBool {
		message := "Settings are incorrect. Please check the 'korean.toml' file!!"
		logger.Fatal(fmt.Errorf("%s", message))
	}
	if c.command == "get-kubeconfig" {
		koreonToml.Kubernetes.GetKubeConfig = true
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

	if c.command != "get-kubeconfig" && len(c.kubeconfig) < 1 {
		return fmt.Errorf("[ERROR]: %s", "To run this ansible-playbook an kubeconfig option must be specified.\n You can get kubeconfig with 'get-kubeconfig' command")
	}

	if c.command != "get-kubeconfig" {
		// Get k8s clientset
		kubeconfigPath, _ := filepath.Abs(c.kubeconfig)
		client, err := createK8sClient(kubeconfigPath)
		if err != nil {
			logger.Fatal(err)
		}

		// Get K8s Cluster Nodes
		kubeNodes := getNodes(client)
		fmt.Println(kubeNodes)
	}

	// Make provision data
	data := model.KoreonctlText{}
	data.KoreOnTemp = koreonToml

	// Processing template
	koreonctlText := template.New("ClusterUpdateText")
	var tempText = ""
	if c.command == "get-kubeconfig" {
		data.Command = "Get Kubeconfig"
		tempText = templates.ClusterGetKubeconfigText
	}
	if c.command == "update" {
		data.Command = "cluster-update"
		tempText = templates.ClusterUpdateText
	}
	temp, err := koreonctlText.Parse(tempText)
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
				results.Prepend("Update Cluster"),
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

func getNodes(client *kubernetes.Clientset) map[string]string {
	var m = make(map[string]string)
	ctx := context.Background()

	listOpts := metav1.ListOptions{
		LabelSelector: "node-role.kubernetes.io/master!=",
	}

	nodes, err := client.CoreV1().Nodes().List(ctx, listOpts)
	if err != nil {
		fmt.Printf("[ERROR] fail to get nodes: %s", err.Error())
		return nil
	}

	for i := 0; i < len((*nodes).Items); i++ {
		node := (*nodes).Items[i]
		ansibleIP := node.Labels["koreon.acornsoft.io/ansible_ssh_host"]
		m[node.Name] = fmt.Sprintf("%s:%s:%s", findNodeInternalIP(node.Status), findNodeExternalIP(node.Status), ansibleIP)
	}
	return m
}

func getVersion(client *kubernetes.Clientset) (int, int, error) {

	serverVersion, err := client.ServerVersion()
	if err != nil {
		fmt.Printf("[ERROR] fail to get nodes: %s", err.Error())
		return -1, -1, err
	}

	spVersion := strings.Split(serverVersion.GitVersion, ".")

	minor, err := strconv.Atoi(spVersion[1])
	if err != nil {
		fmt.Printf("[ERROR] fail to get nodes: %s", err.Error())
		return -1, -1, err
	}

	patch, err := strconv.Atoi(spVersion[2])
	if err != nil {
		fmt.Printf("[ERROR] fail to get nodes: %s", err.Error())
		return -1, -1, err
	}

	return minor, patch, nil

}

func createK8sClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func findNodeInternalIP(status v12.NodeStatus) string {
	for i := 0; i < len(status.Addresses); i++ {
		if status.Addresses[i].Type == v12.NodeInternalIP {
			return status.Addresses[i].Address
		}
	}

	return "<none>"
}

func findNodeExternalIP(status v12.NodeStatus) string {
	for i := 0; i < len(status.Addresses); i++ {
		if status.Addresses[i].Type == v12.NodeExternalIP {
			return status.Addresses[i].Address
		}
	}

	return "<none>"
}
