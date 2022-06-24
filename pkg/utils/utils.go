package utils

import (
	"bufio"
	cryptornad "crypto/rand"
	"fmt"
	"github.com/hhkbp2/go-logging"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"kore-on/pkg/conf"
	"kore-on/pkg/model"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
)

var logger = logging.GetLogger("utils")

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func FileSizeAndExists(name string) (int64, bool, error) {
	//workDir, _ := os.Getwd()
	//fmt.Printf("workdir %s\n",workDir)
	var size int64 = 0
	stat, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return size, false, err
		}
	}
	size = stat.Size()

	return size, true, nil
}

func CopyFile(source string, dest string) (err error) {
	//fmt.Printf("copy file source %s dest %s \n",source, dest)
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		fmt.Printf("error %s\n", err)
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)

		if err == nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		} else {
			return err
		}
	} else {
		return err
	}

	return nil
}

func CopyFile0600(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		err = os.Chmod(dest, 0600)
	} else {
		return err
	}

	return nil
}

func CopyFilePreWork(workDir string, koreonToml model.KoreonToml, cmd string) error {

	errorCnt := 0

	os.MkdirAll(conf.KoreonDestDir, os.ModePerm)

	idRsa := workDir + "/" + conf.KoreonDestDir + "/" + conf.IdRsa
	sslRegistryCrt := workDir + "/" + conf.KoreonDestDir + "/" + conf.SSLRegistryCrt
	sslRegistryKey := workDir + "/" + conf.KoreonDestDir + "/" + conf.SSLRegistryKey

	repoBackupTgz := workDir + "/" + conf.KoreonDestDir + "/" + conf.RepoFile
	harborBackupTgz := workDir + "/" + conf.KoreonDestDir + "/" + conf.HarborFile

	nodePoolSecurityPrivateKeyPath := koreonToml.NodePool.Security.PrivateKeyPath

	isPrivateRegistryPublicCert := koreonToml.PrivateRegistry.PublicCert
	regiSslCert := koreonToml.PrivateRegistry.CertFile.SslCertificate
	regiSslCertKey := koreonToml.PrivateRegistry.CertFile.SslCertificateKey

	repoBackupFilePath := koreonToml.Koreon.LocalRepositoryArchiveFile
	harborBackupFilePath := koreonToml.PrivateRegistry.RegistryArchiveFile

	switch cmd {
	case conf.CMD_CREATE, conf.CMD_APPLY, conf.CMD_DESTROY:
		if !FileExists(nodePoolSecurityPrivateKeyPath) {
			PrintError(fmt.Sprintf("private-key-path : %s file is not found", nodePoolSecurityPrivateKeyPath))
			errorCnt++
		}
	default:
	}

	//레지스트리 설치 여부
	if koreonToml.PrivateRegistry.Install {
		//레지스트리 공인 인증서 사용하는 경우 인증서 파일이 있어야 함.
		if isPrivateRegistryPublicCert && (cmd == conf.CMD_CREATE || cmd == conf.CMD_PREPARE_AIREGAP) {

			if !FileExists(regiSslCert) {
				PrintError(fmt.Sprintf("registry ssl-certificate : %s file is not found", regiSslCert))
				errorCnt++
			}

			if !FileExists(regiSslCertKey) {
				PrintError(fmt.Sprintf("registry ssl-certificate-key : %s file is not found", regiSslCertKey))
				errorCnt++
			}
		}
	}

	if koreonToml.Koreon.ClosedNetwork && cmd == conf.CMD_CREATE {
		if koreonToml.Koreon.LocalRepository == "" {
			if !FileExists(repoBackupFilePath) {
				PrintError(fmt.Sprintf("local-repository-archive-file: %s file is not found", repoBackupFilePath))
				errorCnt++
			}
		}

		if koreonToml.PrivateRegistry.Install {
			if !FileExists(harborBackupFilePath) {
				PrintError(fmt.Sprintf("local-repository-archive-file: %s file is not found", harborBackupFilePath))
				errorCnt++
			}
		}
	}

	if errorCnt > 0 {
		os.Exit(1)
	} else {
		//상단은 validation check 만 진행하고 기능 수행은 여기부터 진행함.
		os.Remove(idRsa)
		os.Remove(sslRegistryCrt)
		os.Remove(sslRegistryKey)

		CopyFile0600(koreonToml.NodePool.Security.PrivateKeyPath, idRsa) //private-key-path copy

		if isPrivateRegistryPublicCert && (cmd == conf.CMD_CREATE || cmd == conf.CMD_PREPARE_AIREGAP) {
			CopyFile0600(regiSslCert, sslRegistryCrt)
			CopyFile0600(regiSslCertKey, sslRegistryKey)
		}

		if koreonToml.Koreon.ClosedNetwork && cmd == conf.CMD_CREATE {
			if koreonToml.Koreon.LocalRepository == "" {
				size2, _, err2 := FileSizeAndExists(repoBackupFilePath)
				if err2 != nil {
					PrintError(err2.Error())
					os.Exit(1)
				}
				size, isExist, _ := FileSizeAndExists(repoBackupTgz)
				if !isExist || (size != size2) {
					CopyFile0600(repoBackupFilePath, repoBackupTgz)
				}
			}

			if koreonToml.PrivateRegistry.Install {
				size2, _, err2 := FileSizeAndExists(harborBackupFilePath)
				if err2 != nil {
					PrintError(err2.Error())
					os.Exit(1)
				}
				size, isExist, _ := FileSizeAndExists(harborBackupTgz)
				if !isExist || (size != size2) {
					CopyFile0600(harborBackupFilePath, harborBackupTgz)
				}
			}
		}
	}
	return nil
}

func CheckDocker() error {
	//fmt.Println("Checking pre-requisition [" + runtime.GOOS + "]")
	_, err := exec.Command("docker", "-v").Output()

	if err != nil {
		//fmt.Println(err.Error())
		PrintError("docker is not found. Install docker before proceeding")
		PrintError("Visit https://www.docker.com/get-started")
		return err
	}
	return nil
}

// AddressRange returns the first and last addresses in the given CIDR range.
func AddressRange(network *net.IPNet) (net.IP, net.IP) {
	// the first IP is easy
	firstIP := network.IP

	// the last IP is the network address OR NOT the mask address
	prefixLen, bits := network.Mask.Size()
	if prefixLen == bits {
		// Easy!
		// But make sure that our two slices are distinct, since they
		// would be in all other cases.
		lastIP := make([]byte, len(firstIP))
		copy(lastIP, firstIP)
		return firstIP, lastIP
	}

	firstIPInt, bits := ipToInt(firstIP)
	hostLen := uint(bits) - uint(prefixLen)
	lastIPInt := big.NewInt(1)
	lastIPInt.Lsh(lastIPInt, hostLen)
	lastIPInt.Sub(lastIPInt, big.NewInt(1))
	lastIPInt.Or(lastIPInt, firstIPInt)

	return firstIP, intToIP(lastIPInt, bits)
}

func ipToInt(ip net.IP) (*big.Int, int) {
	val := &big.Int{}
	val.SetBytes([]byte(ip))
	if len(ip) == net.IPv4len {
		return val, 32
	} else if len(ip) == net.IPv6len {
		return val, 128
	} else {
		panic(fmt.Errorf("Unsupported address length %d", len(ip)))
	}
}

func intToIP(ipInt *big.Int, bits int) net.IP {
	ipBytes := ipInt.Bytes()
	ret := make([]byte, bits/8)
	// Pack our IP bytes into the end of the return array,
	// since big.Int.Bytes() removes front zero padding.
	for i := 1; i <= len(ipBytes); i++ {
		ret[len(ret)-i] = ipBytes[len(ipBytes)-i]
	}
	return net.IP(ret)
}

func getServiceIP(cidr string, nextStep byte) string {

	_, ipv4Net, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Printf(err.Error())
	}
	startIp, _ := AddressRange(ipv4Net)
	//fmt.Println(fmt.Sprintf("start ip %s", startIp))

	startIp[3] = startIp[3] + nextStep
	//fmt.Printf("start ip %s %v\n", startIp, nextStep)

	return fmt.Sprintf("%s", startIp)

}

func PrintInfo(message string) {
	spData := strings.Split(message, "\n")
	if len(spData) > 0 {
		for i := 0; i < len(spData); i++ {
			fmt.Fprintf(os.Stdout, "%s\n", spData[i])
		}
	} else {
		fmt.Fprintf(os.Stdout, "%s\n", message)
	}
}

func PrintError(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", message)
}

func CreateInventoryFile(destDir string, koreonToml model.KoreonToml, addNodes map[string]string) string {

	inventory := "# Inventory create by koreon\n\n"

	masterIps := koreonToml.NodePool.Master.IP
	nodeIps := koreonToml.NodePool.Node.IP
	registryIp := koreonToml.PrivateRegistry.RegistryIP
	storageIp := koreonToml.SharedStorage.StorageIP
	etcdIps := koreonToml.Kubernetes.Etcd.IP

	masterPrivateIps := koreonToml.NodePool.Master.PrivateIP
	nodePrivateIps := koreonToml.NodePool.Node.PrivateIP
	etcdPrivateIps := koreonToml.Kubernetes.Etcd.PrivateIP
	registryPrivateIp := koreonToml.PrivateRegistry.PrivateIP
	storagePrivateIp := koreonToml.SharedStorage.PrivateIP

	sshPort := conf.SshPort

	if koreonToml.NodePool.Security.SSHPort > 0 {
		sshPort = koreonToml.NodePool.Security.SSHPort
	}

	for i := 0; i < len(masterIps); i++ {
		ip := ""
		if len(masterPrivateIps) > 0 {
			ip = masterPrivateIps[i]
		} else {
			ip = masterIps[i]
		}
		inventory += fmt.Sprintf("master-%v ansible_ssh_host=%s ip=%s ansible_port=%v\n", masterIps[i], masterIps[i], ip, sshPort)
	}

	for i := 0; i < len(nodeIps); i++ {
		ip := ""
		if len(nodePrivateIps) > 0 {
			ip = nodePrivateIps[i]
		} else {
			ip = nodeIps[i]
		}

		inventoryItem := []string{
			fmt.Sprintf("worker-%v ansible_ssh_host=%s ip=%s ansible_port=%v", nodeIps[i], nodeIps[i], ip, sshPort),
		}
		inventoryItem = append(inventoryItem, "\n")
		inventory += strings.Join(inventoryItem, " ")
	}

	for i := 0; i < len(etcdIps); i++ {
		ip := ""
		if len(etcdPrivateIps) > 0 {
			ip = etcdPrivateIps[i]
		} else {
			ip = etcdIps[i]
		}
		inventory += fmt.Sprintf("etcd-%v ansible_ssh_host=%s ip=%s ansible_port=%v\n", etcdIps[i], etcdIps[i], ip, sshPort)
	}

	if koreonToml.PrivateRegistry.Install {
		ip := ""
		if registryPrivateIp != "" {
			ip = registryPrivateIp
		} else {
			ip = registryIp
		}
		inventory += fmt.Sprintf("registry-%v ansible_ssh_host=%s ip=%s ansible_port=%v\n", registryIp, registryIp, ip, sshPort)
	}

	if koreonToml.SharedStorage.Install {
		ip := ""
		if storagePrivateIp != "" {
			ip = storagePrivateIp
		} else {
			ip = storageIp
		}
		inventory += fmt.Sprintf("storage-%v ansible_ssh_host=%s ip=%s ansible_port=%v\n", storageIp, storageIp, ip, sshPort)
	}

	//etcd
	inventory += fmt.Sprintf("\n")
	inventory += fmt.Sprintf("[etcd]\n")

	for i := 0; i < len(etcdIps); i++ {
		inventory += fmt.Sprintf("etcd-%v\n", etcdIps[i])
	}

	//etcd-private
	inventory += fmt.Sprintf("\n")
	inventory += fmt.Sprintf("[etcd-private]\n")

	for i := 0; i < len(etcdIps); i++ {
		inventory += fmt.Sprintf("etcd-%v\n", etcdIps[i])
	}

	//masters
	inventory += fmt.Sprintf("\n")
	inventory += fmt.Sprintf("[masters]\n")

	for i := 0; i < len(masterIps); i++ {
		inventory += fmt.Sprintf("master-%v\n", masterIps[i])
	}

	//sslhost
	inventory += fmt.Sprintf("\n")
	inventory += fmt.Sprintf("[sslhost]\n")

	if masterIps != nil {
		inventory += fmt.Sprintf("master-%v\n", masterIps[0])
	}

	//node
	inventory += fmt.Sprintf("\n")
	inventory += fmt.Sprintf("[node]\n")

	if addNodes != nil && len(addNodes) > 0 {
		for ip, _ := range addNodes {
			inventory += fmt.Sprintf("worker-%v\n", ip)
		}
	} else {
		for j := 0; j < len(nodeIps); j++ {
			inventory += fmt.Sprintf("worker-%v\n", nodeIps[j])
		}
	}

	//registry
	inventory += fmt.Sprintf("\n")
	inventory += fmt.Sprintf("[registry]\n")
	if koreonToml.PrivateRegistry.Install {
		if koreonToml.PrivateRegistry.Install {
			//inventory += fmt.Sprintf("registry01\n")
			inventory += fmt.Sprintf("registry-%v\n", registryIp)
		}
	}

	//storage
	inventory += fmt.Sprintf("\n")
	inventory += fmt.Sprintf("[storage]\n")
	if koreonToml.SharedStorage.Install {
		if koreonToml.SharedStorage.Install {
			//inventory += fmt.Sprintf("storage01\n")
			inventory += fmt.Sprintf("storage-%v\n", storageIp)
		}
	}

	inventory += fmt.Sprintf("\n")
	inventory += fmt.Sprintf("[cluster:children]\n")
	inventory += fmt.Sprintf("masters\n")
	inventory += fmt.Sprintf("node\n")

	b := []byte(inventory)
	inventoryPath := destDir + "/" + conf.KoreonDestDir
	os.MkdirAll(inventoryPath, os.ModePerm)
	ioutil.WriteFile(inventoryPath+"/"+conf.KoreonInventoryIni, b, 0600)

	return inventoryPath + "/" + conf.KoreonInventoryIni
}

func CreateBasicYaml(destDir string, koreonToml model.KoreonToml, command string) string {
	var allYaml = model.BasicYaml{}

	//default values
	//koreon
	allYaml.Provider = false
	allYaml.CloudProvider = "onpremise"
	allYaml.InstallDir = "/var/lib/koreon"
	allYaml.CertValidityDays = 365

	//#kubernetes
	allYaml.ServiceIPRange = "10.96.0.0/12"
	allYaml.PodIPRange = "10.32.0.0/12"
	allYaml.NodePortRange = "30000-32767"
	allYaml.LbPort = 6443
	allYaml.KubeProxyMode = "iptables"      // iptable, ipvs
	allYaml.ContainerRuntime = "containerd" // docker,containerd
	allYaml.ApiSans = koreonToml.Kubernetes.ApiSans

	//NodePool
	allYaml.DataRootDir = "/data"

	allYaml.ClusterName = koreonToml.Koreon.ClusterName

	clusterID, _ := NewUUID()
	allYaml.AuditLogEnable = koreonToml.Kubernetes.AuditLogEnable
	allYaml.ClusterID = clusterID

	//koreon
	//
	if koreonToml.Koreon.CertValidityDays > 0 {
		allYaml.CertValidityDays = koreonToml.Koreon.CertValidityDays
	}

	allYaml.ClosedNetwork = koreonToml.Koreon.ClosedNetwork
	if koreonToml.Koreon.InstallDir != "" {
		allYaml.InstallDir = koreonToml.Koreon.InstallDir
	}

	allYaml.LocalRepository = koreonToml.Koreon.LocalRepository

	//k8s
	//
	allYaml.K8SVersion = koreonToml.Kubernetes.Version

	if koreonToml.Kubernetes.ServiceCidr != "" {
		allYaml.ServiceIPRange = koreonToml.Kubernetes.ServiceCidr
	}

	if koreonToml.Kubernetes.PodCidr != "" {
		allYaml.PodIPRange = koreonToml.Kubernetes.PodCidr
	}

	if koreonToml.Kubernetes.NodePortRange != "" {
		allYaml.NodePortRange = koreonToml.Kubernetes.NodePortRange
	}

	allYaml.AuditLogEnable = koreonToml.Kubernetes.AuditLogEnable

	if koreonToml.Kubernetes.KubeProxyMode != "" {
		allYaml.KubeProxyMode = koreonToml.Kubernetes.KubeProxyMode
	}

	if koreonToml.Kubernetes.ContainerRuntime != "" {
		allYaml.ContainerRuntime = koreonToml.Kubernetes.ContainerRuntime
	}

	//vxlan-mode
	allYaml.KubeProxyMode = koreonToml.Kubernetes.KubeProxyMode

	//nodepool
	if koreonToml.NodePool.DataDir != "" {
		allYaml.DataRootDir = koreonToml.NodePool.DataDir
	}

	if len(koreonToml.NodePool.Master.PrivateIP) == len(koreonToml.NodePool.Master.IP) {
		allYaml.APILbIP = fmt.Sprintf("https://%s:%d", koreonToml.NodePool.Master.PrivateIP[0], allYaml.LbPort)
	} else {
		allYaml.APILbIP = fmt.Sprintf("https://%s:%d", koreonToml.NodePool.Master.IP[0], allYaml.LbPort)
	}

	if koreonToml.NodePool.Master.LbIP == "" {
		allYaml.LbIP = koreonToml.NodePool.Master.IP[0]
	} else {
		allYaml.LbIP = koreonToml.NodePool.Master.LbIP
	}

	allYaml.MasterIsolated = koreonToml.NodePool.Master.Isolated
	allYaml.Haproxy = koreonToml.NodePool.Master.HaproxyInstall //# Set False When Already Physical Loadbalancer Available"

	//storage
	allYaml.StorageInstall = koreonToml.SharedStorage.Install
	allYaml.NfsIP = koreonToml.SharedStorage.StorageIP
	allYaml.NfsVolumeDir = koreonToml.SharedStorage.VolumeDir

	//registry
	isPrivateRegistryPubicCert := koreonToml.PrivateRegistry.PublicCert

	//os.MkdirAll(sshPath, os.ModePerm)
	//CopyFile(conf.KoreonDestDir+"/"+conf.IdRsa, sshPath+"/"+conf.IdRsa)
	//CopyFile(conf.KoreonDestDir+"/"+"id_rsa.pub", sshPath+"/id_rsa.pub")

	registryIP := koreonToml.PrivateRegistry.RegistryIP
	registryDomain := koreonToml.PrivateRegistry.RegistryIP

	allYaml.RegistryInstall = koreonToml.PrivateRegistry.Install
	allYaml.RegistryDataDir = koreonToml.PrivateRegistry.DataDir
	allYaml.Registry = registryIP
	allYaml.RegistryDomain = strings.Replace(registryDomain, "https://", "", -1)
	allYaml.RegistryPublicCert = isPrivateRegistryPubicCert

	if koreonToml.PrivateRegistry.RegistryDomain != "" {
		registryDomain = koreonToml.PrivateRegistry.RegistryDomain
	}

	switch command {
	case conf.CMD_PREPARE_AIREGAP:
		allYaml.ClosedNetwork = false
		allYaml.ArchiveRepo = true
	default:

	}

	if koreonToml.Koreon.ClosedNetwork {
		allYaml.LocalRepository = fmt.Sprintf("http://%s:8080", koreonToml.PrivateRegistry.RegistryIP)
		allYaml.LocalRepositoryArchiveFile = conf.RepoFile
		allYaml.RegistryArchiveFile = conf.HarborFile
	}

	b, _ := yaml.Marshal(allYaml)

	allYamlPath := destDir + "/" + conf.KoreonDestDir
	os.MkdirAll(allYamlPath, os.ModePerm)
	ioutil.WriteFile(allYamlPath+"/"+conf.KoreonBasicYaml, b, 0600)

	return allYamlPath + "/" + conf.KoreonBasicYaml
}

func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(cryptornad.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func CheckUserInput(prompt string, checkWord string) bool {
	var res string
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	buf, _ := reader.ReadString('\n')

	if runtime.GOOS == "windows" {
		res = strings.Split(buf, "\r\n")[0]
	} else {
		res = strings.Split(buf, "\n")[0]
	}

	if res == checkWord {
		return true
	}

	return false
}

func IsSupportK8sVersion(version string) bool {
	isSupport := false
	for _, v := range conf.SupportK8SVersion {
		if v == version {
			isSupport = true
			break
		}
	}
	return isSupport
}

func ExecKoreonCmd(commandArgs []string, isStep bool, debugMode bool) error {
	var err error

	if isStep {

		if debugMode {
			fmt.Printf("syscall.Exec name %s \n", commandArgs[0])
			fmt.Printf("syscall.Exec arg %s \n", commandArgs[1:])
		}

		err = syscall.Exec(commandArgs[0], commandArgs, os.Environ())
	} else {
		if debugMode {
			fmt.Printf("exec.Command name %s \n", commandArgs[0])
			fmt.Printf("exec.Command arg %s \n", commandArgs[1:])
		}
		//
		cmd := exec.Command(commandArgs[0], commandArgs[1:]...)

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}

		chStdOut := make(chan string)
		chStdErr := make(chan string)

		go func() {
			merged := io.Reader(stdout)
			scanner := bufio.NewScanner(merged)
			for scanner.Scan() {
				m := scanner.Text()
				//fmt.Println("gofunc-stdout " + m)
				text := strings.TrimSpace(m)
				chStdOut <- text
			}

		}()

		go func() {
			merged := io.Reader(stderr)
			scanner := bufio.NewScanner(merged)
			for scanner.Scan() {
				m := scanner.Text()
				//fmt.Println("gofunc-stderr "+ m)
				text := strings.TrimSpace(m)
				chStdErr <- text
			}
		}()

		go func() {
			for {
				select {
				case line := <-chStdOut:
					text := strings.TrimSpace(line)

					if strings.Contains(text, "Failed to connect") {
						fmt.Fprintf(os.Stderr, "%s\n", text)
					} else {
						fmt.Println(fmt.Sprintf("%s", text))
					}
				case line := <-chStdErr:
					text := strings.TrimSpace(line)
					fmt.Fprintf(os.Stderr, "%s\n", text)

				}
			}
		}()

		err = cmd.Start()
		if err != nil {
			logger.Errorf("command start: %s", err.Error())
			return err
		}

		err = cmd.Wait()
		if err != nil {
			logger.Errorf("command wait: %s", err.Error())
			return err
		}
	}

	return err
}

func CheckError(err error) {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func checkIp(public string, private string) error {
	r, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if len(public) == 0 {
		return fmt.Errorf("public ip of node not exists")
	} else if !r.MatchString(public) {
		return fmt.Errorf("public ip invalid: %s\n", public)
	}
	if len(private) == 0 {
		return fmt.Errorf("private ip of node not exists")
	} else if !r.MatchString(private) {
		return fmt.Errorf("private ip invalid: %s\n", private)
	}

	return nil
}

func ParentPath() string {
	wd, err := os.Getwd()
	if err != nil {
		PrintError(err.Error())
		os.Exit(1)
	}
	parent := filepath.Dir(wd)
	return parent
}

//func GetCubeToml(workDir string) (model.KoreonToml, error) {
//	var koreonToml = model.KoreonToml{}
//
//	if !FileExists(conf.KoreonConfigFile) {
//		return koreonToml, fmt.Errorf("file is not found")
//	}
//
//	c, err := ioutil.ReadFile(workDir + "/" + conf.KoreonConfigFile)
//	if err != nil {
//		//PrintError(err.Error())
//		return koreonToml, err
//	}
//
//	str := string(c)
//	str = strings.Replace(str, "\\", "/", -1)
//	c = []byte(str)
//
//	err = toml.Unmarshal(c, &koreonToml)
//	if err != nil {
//		PrintError(err.Error())
//		return koreonToml, err
//	}
//
//	return koreonToml, nil
//}

//const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//var Spnr = spinner.New(spinner.CharSets[9], 100*time.Millisecond)

//var knownProviders = []string{
//	"gcp",
//	"azure",
//	"aws",
//	"onpremise",
//	"aliyun",
//	"eks",
//	"aks",
//	"gke",
//	"tke",
//	"tencent",
//	"diamanti",
//}
//
//var localProviders = []string{
//	"virtualbox",
//	"minikube",
//}

//func CopyDir(source string, dest string) (err error) {
//
//	// get properties of source dir
//	sourceinfo, err := os.Stat(source)
//	if err != nil {
//		return err
//	}
//
//	// create dest dir
//
//	err = os.MkdirAll(dest, sourceinfo.Mode())
//	if err != nil {
//		return err
//	}
//
//	directory, _ := os.Open(source)
//
//	objects, err := directory.Readdir(-1)
//
//	for _, obj := range objects {
//
//		sourcefilepointer := source + "/" + obj.Name()
//
//		destinationfilepointer := dest + "/" + obj.Name()
//
//		if obj.IsDir() {
//			// create sub-directories - recursively
//			err = CopyDir(sourcefilepointer, destinationfilepointer)
//			if err != nil {
//				logger.Error(err)
//			}
//		} else {
//			// perform copy
//			err = CopyFile(sourcefilepointer, destinationfilepointer)
//			if err != nil {
//				logger.Error(err)
//			}
//		}
//
//	}
//	return
//}
//
//func ReadFile(filePath string, buf *[]byte) {
//
//	file, err := os.Open(filePath)
//	CheckError(err)
//
//	defer file.Close()
//
//	fi, err := file.Stat()
//	CheckError(err)
//
//	*buf = make([]byte, fi.Size())
//
//	_, err = file.Read(*buf)
//	CheckError(err)
//}
//
//func WriteFile(filePath string, buf *[]byte) error {
//	f, err := os.OpenFile(filePath, os.O_RDWR, 0644)
//	if err != nil {
//		logger.Errorf("error while open file: %s\n", err.Error())
//		return err
//	}
//	defer f.Close()
//
//	err = f.Truncate(0)
//	if err != nil {
//		logger.Errorf("fail to write file[1]: %s\n", err.Error())
//		return err
//	}
//
//	_, err = f.WriteString(string(*buf))
//	if err != nil {
//		logger.Errorf("fail to write file[2]: %s\n", err.Error())
//		return err
//	}
//
//	return nil
//}
//
//func WriteFileString(filePath string, content string) error {
//	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
//	if err != nil {
//		logger.Errorf("error while open file: %s\n", err.Error())
//		return err
//	}
//	defer f.Close()
//
//	err = f.Truncate(0)
//	if err != nil {
//		logger.Errorf("fail to write file[1]: %s\n", err.Error())
//		return err
//	}
//
//	_, err = f.WriteString(content)
//	if err != nil {
//		logger.Errorf("fail to write file[2]: %s\n", err.Error())
//		return err
//	}
//
//	return nil
//}
