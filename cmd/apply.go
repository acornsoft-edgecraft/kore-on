package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"kore-on/pkg/conf"
	"kore-on/pkg/model"
	"kore-on/pkg/utils"
	"log"
	"os"
	"strconv"
	"strings"
)

type strApplyCmd struct {
	dryRun  bool
	verbose bool
	step    bool
}

func applyCmd() *cobra.Command {
	apply := &strApplyCmd{}
	cmd := &cobra.Command{
		Use:          "apply [flags]",
		Short:        "apply",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return apply.run()
		},
	}
	f := cmd.Flags()
	f.BoolVarP(&apply.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&apply.step, "step", "", false, "step")
	f.BoolVarP(&apply.dryRun, "dry-run", "d", false, "dryRun")
	return cmd
}

func (c *strApplyCmd) run() error {

	if !utils.CheckUserInput("Do you really want to apply? Only 'yes' will be accepted to confirm: ", "yes") {
		fmt.Println("nothing to changed. exit")
		os.Exit(1)
	}
	applyCnt := 0
	isUpgrade := false
	var addMap = make(map[string]string)
	var delMap = make(map[string]string)
	workDir, _ := os.Getwd()
	koreonToml, _ := utils.ValidateKoreonTomlConfig(workDir)

	utils.CopyFilePreWork(workDir, koreonToml, "apply")

	koreonK8sVersion := koreonToml.Kubernetes.Version
	spK8sVersion := strings.Split(koreonK8sVersion, ".")
	koreonMinor, _ := strconv.Atoi(spK8sVersion[1])
	koreonPatch, _ := strconv.Atoi(spK8sVersion[2])

	port := 22
	certFileName := conf.KoreonDestDir + "/" + conf.IdRsa

	client, err := createK8sClient(koreonToml.NodePool.Master.IP[0], port, koreonToml.NodePool.Security.SSHUserID, certFileName)
	if err != nil {
		return err
	}

	minor, patch, _ := getVersion(client)

	if minor == koreonMinor {
		if patch < koreonPatch {
			//업그레이드 가능
			isUpgrade = true
		} else if patch > koreonPatch {
			return fmt.Errorf("downgrade not supported. current version 1.%d.%d, toml > kubernetes > version: %s\n", minor, patch, koreonK8sVersion)
		}
	} else if minor+1 == koreonMinor {
		//업그레이드 가능
		isUpgrade = true
	} else if minor+1 < koreonMinor {
		//2단계로 여서 업그레이드 불가능
		return fmt.Errorf("upgrade not supported. current version 1.%d.%d, toml > kubernetes > version: %s\n", minor, patch, koreonK8sVersion)

	} else if minor > koreonMinor {
		return fmt.Errorf("downgrade not supported. current version 1.%d.%d, toml > kubernetes > version: %s\n", minor, patch, koreonK8sVersion)
	}

	utils.CheckDocker()

	utils.CopyFilePreWork(workDir, koreonToml, "apply")

	var kubeNodes = getNodes(client)

	//추가 목록
	for i := 0; i < len(koreonToml.NodePool.Node.IP); i++ {

		ip := koreonToml.NodePool.Node.IP[i]
		privateIp := koreonToml.NodePool.Node.IP[i]
		if len(koreonToml.NodePool.Node.IP) == len(koreonToml.NodePool.Node.PrivateIP) {
			privateIp = koreonToml.NodePool.Node.PrivateIP[i]
		}

		exists := false
		for name, _ := range kubeNodes {
			svrIp := kubeNodes[name]
			if strings.Contains(svrIp, ip) || strings.Contains(svrIp, privateIp) {
				exists = true
			}
		}

		if !exists {
			addMap[ip] = ""
		}
	}

	//삭제 목록
	for name, _ := range kubeNodes {
		svrIp := kubeNodes[name]
		exists := false
		for j := 0; j < len(koreonToml.NodePool.Node.IP); j++ {
			ip := koreonToml.NodePool.Node.IP[j]
			privateIp := koreonToml.NodePool.Node.IP[j]
			if len(koreonToml.NodePool.Node.IP) == len(koreonToml.NodePool.Node.PrivateIP) {
				privateIp = koreonToml.NodePool.Node.PrivateIP[j]
			}

			if strings.Contains(svrIp, ip) || strings.Contains(svrIp, privateIp) {
				exists = true
			}
		}
		if !exists {
			delMap[name] = svrIp
		}
	}

	//etcd check
	for delNodeName, _ := range delMap {
		svrIp := kubeNodes[delNodeName]
		for j := 0; j < len(koreonToml.Kubernetes.Etcd.IP); j++ {
			if koreonToml.Kubernetes.Etcd.IP[j] == svrIp {
				utils.PrintInfo(fmt.Sprintf(conf.ERROR_FORMAT, fmt.Sprintf("Delete worker node running etcd is not allowed.")))
				return nil
			}
		}
	}

	if koreonToml.Koreon.DebugMode {
		fmt.Println(fmt.Sprintf("add node %v", addMap))
		fmt.Println(fmt.Sprintf("del node %v", delMap))
		fmt.Println(fmt.Sprintf("upgrade kubernetes 1.%v.%v -> %v", minor, patch, koreonK8sVersion))
	}

	if len(addMap) > 0 {
		applyCnt += 1
	}
	if len(delMap) > 0 {
		applyCnt += 1
	}
	if isUpgrade {
		applyCnt += 1
	}

	if applyCnt == 0 {
		utils.PrintInfo(fmt.Sprintf("There is no worker node to added or deleted."))
	} else if applyCnt > 1 && c.step {
		utils.PrintInfo(fmt.Sprintf("The step option is only available in one function."))
	} else {
		// 공통
		sshId := koreonToml.NodePool.Security.SSHUserID
		basicFilePath := utils.CreateBasicYaml(workDir, koreonToml, conf.CMD_APPLY)

		//노드 추가
		if len(addMap) > 0 {
			inventoryFilePath := utils.CreateInventoryFile(workDir, koreonToml, addMap)
			err := addNode(workDir, inventoryFilePath, basicFilePath, sshId, addMap, c, koreonToml.Koreon.DebugMode)
			if err != nil {
				log.Printf("Command finished with error: %v", err)
			} else {
				fmt.Println(fmt.Sprintf("end add node %v", addMap))
			}
		}

		//노드 삭제
		if len(delMap) > 0 {
			inventoryFilePath := utils.CreateInventoryFile(workDir, koreonToml, nil)
			err := removeNode(workDir, inventoryFilePath, basicFilePath, sshId, delMap, c, koreonToml.Koreon.DebugMode)
			if err != nil {
				log.Printf("Command finished with error: %v", err)
			} else {
				fmt.Println(fmt.Sprintf("end add node %v", addMap))
			}
		}

		//업그레이드
		if isUpgrade {
			inventoryFilePath := utils.CreateInventoryFile(workDir, koreonToml, nil)
			err := upgrade(workDir, inventoryFilePath, basicFilePath, sshId, c, koreonToml.Koreon.DebugMode)
			if err != nil {
				log.Printf("Command finished with error: %v", err)
			}
		}
	}
	return nil
}

func addNode(workDir string, inventoryFilePath string, basicFilePath string, sshId string, addMap map[string]string, c *strApplyCmd, debugMode bool) error {
	itOption := "-t"
	if c.step {
		itOption = "-it"
	}
	commandArgs := []string{
		"docker",
		"run",
		"--name",
		conf.KoreonImageName,
		"--rm",
		"--privileged",
		itOption,
		"-v",
		fmt.Sprintf("%s:%s", workDir, conf.WorkDir),
		"-v",
		fmt.Sprintf("%s:%s", inventoryFilePath, conf.InventoryIni),
		"-v",
		fmt.Sprintf("%s:%s", basicFilePath, conf.BasicYaml),
		conf.KoreonImage,
		"ansible-playbook",
		"-i",
		conf.InventoryIni,
		"-u",
		sshId,
		"--private-key",
		conf.KoreonDestDir + "/" + conf.IdRsa,
		conf.AddNodeYaml,
	}

	if c.verbose {
		commandArgs = append(commandArgs, "-v")
	}

	if c.step {
		commandArgs = append(commandArgs, "--step")
	}

	if c.dryRun {
		commandArgs = append(commandArgs, "-C")
		commandArgs = append(commandArgs, "-D")
	}

	err := utils.ExecKoreonCmd(commandArgs, c.step, debugMode)

	return err
}

func removeNode(workDir string, inventoryFilePath string, basicFilePath string, sshId string, delMap map[string]string, c *strApplyCmd, debugMode bool) error {

	itOption := "-t"
	if c.step {
		itOption = "-it"
	}

	for delNodeName, ip := range delMap {
		delIp := strings.Split(ip, ":")[2]
		if delIp == "<none>" || delIp == "" {
			delIp = strings.Split(ip, ":")[1]
		}
		if delIp == "<none>" || delIp == "" {
			delIp = strings.Split(ip, ":")[0]
		}

		commandArgs := []string{
			"docker",
			"run",
			"--name",
			conf.KoreonImageName,
			"--rm",
			"--privileged",
			itOption,
			"-v",
			fmt.Sprintf("%s:%s", workDir, conf.WorkDir),
			"-v",
			fmt.Sprintf("%s:%s", inventoryFilePath, conf.InventoryIni),
			"-v",
			fmt.Sprintf("%s:%s", basicFilePath, conf.BasicYaml),
			conf.KoreonImage,
			"ansible-playbook",
			"-i",
			conf.InventoryIni,
			"-u",
			sshId,
			"--private-key",
			conf.KoreonDestDir + "/" + conf.IdRsa,
			"-e",
			fmt.Sprintf("remove_node_name=%s", delNodeName),
			"-e",
			fmt.Sprintf("target=%s", delIp),
			conf.RemoveNodeYaml,
		}

		if c.verbose {
			commandArgs = append(commandArgs, "-v")
		}

		if c.step {
			commandArgs = append(commandArgs, "--step")
		}

		if c.dryRun {
			commandArgs = append(commandArgs, "-C")
			commandArgs = append(commandArgs, "-D")
		}

		err := utils.ExecKoreonCmd(commandArgs, c.step, debugMode)
		if err != nil {
			return err
		}
	}
	return nil
}

func upgrade(workDir string, inventoryFilePath string, basicFilePath string, sshId string, c *strApplyCmd, debugMode bool) error {

	itOption := "-t"
	if c.step {
		itOption = "-it"
	}

	commandArgs := []string{
		"docker",
		"run",
		"--name",
		conf.KoreonImageName,
		"--rm",
		"--privileged",
		itOption,
		"-v",
		fmt.Sprintf("%s:%s", workDir, conf.WorkDir),
		"-v",
		fmt.Sprintf("%s:%s", inventoryFilePath, conf.InventoryIni),
		"-v",
		fmt.Sprintf("%s:%s", basicFilePath, conf.BasicYaml),
		conf.KoreonImage,
		"ansible-playbook",
		"-i",
		conf.InventoryIni,
		"-u",
		sshId,
		"--private-key",
		conf.KoreonDestDir + "/" + conf.IdRsa,
		conf.UpgradeYaml,
	}

	if c.verbose {
		commandArgs = append(commandArgs, "-v")
	}

	if c.step {
		commandArgs = append(commandArgs, "--step")
	}

	if c.dryRun {
		commandArgs = append(commandArgs, "-C")
		commandArgs = append(commandArgs, "-D")
	}

	err := utils.ExecKoreonCmd(commandArgs, c.step, debugMode)

	return err
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

func createK8sClient(ip string, port int, user string, certFileName string) (*kubernetes.Clientset, error) {

	client := &utils.SSH{
		IP:   ip,
		Port: port,
		User: user,
		Cert: certFileName,
	}
	client.Connect()
	aa := client.RunCmd(fmt.Sprintf("cat %s/%s", conf.KoreonKubeConfigPath, conf.KoreonKubeConfig))
	client.Close()

	y := model.KubeConfig{}

	err := yaml.Unmarshal([]byte(aa), &y)
	if err != nil {
		return nil, err
	}

	var bb = &restclient.Config{}

	bb.Host = y.Clusters[0].Cluster.Server

	certData, _ := base64.StdEncoding.DecodeString(y.Users[0].User.ClientCertificateData)
	keyData, _ := base64.StdEncoding.DecodeString(y.Users[0].User.ClientKeyData)
	caData, _ := base64.StdEncoding.DecodeString(y.Clusters[0].Cluster.CertificateAuthorityData)

	bb.TLSClientConfig.CertData = certData
	bb.TLSClientConfig.KeyData = keyData
	bb.TLSClientConfig.CAData = caData

	clientset, err := kubernetes.NewForConfig(bb)
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
