package utils

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"kore-on/pkg/conf"
	"kore-on/pkg/model"
	"os"
	"strings"
)

func GetKoreonTomlConfig(workDir string) (model.KoreonToml, error) {

	errorCnt := 0
	configFullPath := workDir + "/" + conf.KoreonConfigFile

	var c []byte
	var err error

	if !FileExists(configFullPath) {
		PrintError(fmt.Sprintf("%s file is not found. Run koreonctl init first", conf.KoreonConfigFile))
		os.Exit(1)
	}

	c, err = ioutil.ReadFile(configFullPath)
	if err != nil {
		PrintError(err.Error())
		os.Exit(1)
	}

	str := string(c)
	str = strings.Replace(str, "\\", "/", -1)
	c = []byte(str)

	var koreonToml = model.KoreonToml{}
	err = toml.Unmarshal(c, &koreonToml)
	if err != nil {
		PrintError(err.Error())
		errorCnt++
	}

	return koreonToml, err
}

func ValidateKoreonTomlConfig(workDir string) (model.KoreonToml, bool) {
	errorCnt := 0
	koreonToml, _ := GetKoreonTomlConfig(workDir)

	koreonClusterName := koreonToml.Koreon.ClusterName

	kubernetesPodCidr := koreonToml.Kubernetes.PodCidr
	kubernetesServiceCidr := koreonToml.Kubernetes.ServiceCidr
	k8sVersion := koreonToml.Kubernetes.Version
	//apiSans := koreonToml.Kubernetes.ApiSans
	etcdCnt := len(koreonToml.Kubernetes.Etcd.IP)

	nodePoolDataDir := koreonToml.NodePool.DataDir
	nodePoolSecuritySSHUserID := koreonToml.NodePool.Security.SSHUserID
	nodePoolSecurityPrivateKeyPath := koreonToml.NodePool.Security.PrivateKeyPath
	nodePoolMasterLbIP := koreonToml.NodePool.Master.LbIP

	privateRegistryInstall := koreonToml.PrivateRegistry.Install
	privateRegistryRegistryIP := koreonToml.PrivateRegistry.RegistryIP
	privateRegistryRegistryDomain := koreonToml.PrivateRegistry.RegistryDomain

	privateRegistryDataDir := koreonToml.PrivateRegistry.DataDir
	isPrivateRegistryPublicCert := koreonToml.PrivateRegistry.PublicCert
	privateRegistryCrt := koreonToml.PrivateRegistry.CertFile.SslCertificate
	privateRegistryKey := koreonToml.PrivateRegistry.CertFile.SslCertificateKey

	if koreonClusterName == "" {
		PrintError("koreon > cluster-name is required.")
		errorCnt++
		//todo 길이 체크
	}

	if koreonToml.Koreon.InstallDir != "" && !strings.HasPrefix(koreonToml.Koreon.InstallDir, "/") {
		PrintError("koreon > install-dir is Only absolute paths are supported.")
		errorCnt++
	}

	if k8sVersion == "" {
		PrintError("kubernetes > version is required.")
		errorCnt++
	} else if !IsSupportK8sVersion(k8sVersion) {
		PrintError(fmt.Sprintf("kubernetes > supported version: %v", conf.SupportK8SVersion))
		errorCnt++
	}

	if nodePoolSecuritySSHUserID == "" {
		PrintError("node-pool.security > ssh-user-id is required.")
		errorCnt++
	}

	if nodePoolSecurityPrivateKeyPath == "" {
		PrintError("node-pool.security > private-key-path is required.")
		errorCnt++
	}

	if nodePoolMasterLbIP == "" {
		PrintError("node-pool.master > lb-ip is required.")
		errorCnt++
	}

	if len(kubernetesPodCidr) > 0 {
		//todo check cider
	}
	if len(kubernetesServiceCidr) > 0 {
		//todo check cider
	}

	switch etcdCnt {
	case 1, 3, 5:
	default:
		PrintError("Only odd number of etcd nodes are supported.(1, 3, 5)")
		errorCnt++

	}

	if len(nodePoolDataDir) > 0 {
		// todo node pool data dir check
	}

	//storage check
	cnt := checkSharedStorage(koreonToml)
	errorCnt += cnt

	if privateRegistryInstall == true {

		if privateRegistryRegistryIP == "" {
			PrintError("private-registry > registry-ip is required.")
			errorCnt++
		}

		if privateRegistryDataDir == "" {
			PrintError("private-registry > data-dir is required.")
			errorCnt++
		}

		if isPrivateRegistryPublicCert {
			if privateRegistryCrt == "" {
				PrintError("private-registry.cert-file > ssl-certificate is required.")
				errorCnt++
			}

			if privateRegistryKey == "" {
				PrintError("private-registry.cert-file > ssl-certificate-key is required.")
				errorCnt++
			}

			if privateRegistryRegistryDomain == "" {
				PrintError("private-registry > registry-domain is required.")
				errorCnt++
			}
		}
	}

	if koreonToml.Koreon.ClosedNetwork {
		if koreonToml.Koreon.LocalRepository == "" && koreonToml.Koreon.LocalRepositoryArchiveFile == "" {
			PrintError("koreon> local-repository or local-repository-archive-file is required.")
			errorCnt++
		}

		if privateRegistryInstall {
			if koreonToml.PrivateRegistry.RegistryArchiveFile == "" {
				PrintError("private-registry >  registry-archive-file is required.")
				errorCnt++
			}
		}
	}

	if errorCnt > 0 {
		logger.Error("there are one or more errors")
		os.Exit(1)
		return koreonToml, false
	}
	return koreonToml, true
}

func checkSharedStorage(koreonToml model.KoreonToml) int {
	errorCnt := 0

	if koreonToml.SharedStorage.Install == true {

		if koreonToml.SharedStorage.VolumeDir == "" {
			PrintError("shared-storage > volume-dir is required.")
			errorCnt++
		}

		if koreonToml.SharedStorage.VolumeSize < 10 {
			PrintError("shared-storage > volume-size is 10 or more.")
			errorCnt++
		}

		if koreonToml.SharedStorage.StorageIP == "" {
			PrintError("shared-storage > storage-ip is required.")
			errorCnt++
		}

	}

	return errorCnt
}
