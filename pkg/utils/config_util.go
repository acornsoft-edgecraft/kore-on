package utils

import (
	"fmt"
	"io/ioutil"
	"kore-on/pkg/logger"
	"kore-on/pkg/model"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
)

func GetKoreonTomlConfig(koreOnConfigFilePath string) (model.KoreOnToml, error) {

	errorCnt := 0
	// configFullPath := workDir + "/" + conf.KoreonConfigFile

	var c []byte
	var err error

	if !FileExists(koreOnConfigFilePath) {
		logger.Fatal(fmt.Sprintf("%s file is not found. Run koreonctl init first", koreOnConfigFilePath))
		os.Exit(1)
	}

	c, err = ioutil.ReadFile(koreOnConfigFilePath)
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}

	str := string(c)
	str = strings.Replace(str, "\\", "/", -1)
	c = []byte(str)

	var koreonToml = model.KoreOnToml{}
	err = toml.Unmarshal(c, &koreonToml)
	if err != nil {
		logger.Fatal(err.Error())
		errorCnt++
	}

	return koreonToml, err
}

func ValidateKoreonTomlConfig(koreOnConfigFilePath string) (model.KoreOnToml, bool) {
	errorCnt := 0
	koreonToml, _ := GetKoreonTomlConfig(koreOnConfigFilePath)

	// koreonClusterName := koreonToml.KoreOn.ClusterName

	kubernetesPodCidr := koreonToml.Kubernetes.PodCidr
	kubernetesServiceCidr := koreonToml.Kubernetes.ServiceCidr
	// k8sVersion := koreonToml.Kubernetes.Version
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

	// if koreonClusterName == "" {
	// 	logger.Fatal("koreon > cluster-name is required.")
	// 	errorCnt++
	// 	//todo 길이 체크
	// }

	if koreonToml.KoreOn.InstallDir != "" && !strings.HasPrefix(koreonToml.KoreOn.InstallDir, "/") {
		logger.Fatal("koreon > install-dir is Only absolute paths are supported.")
		errorCnt++
	}

	// if k8sVersion == "" {
	// 	logger.Fatal("kubernetes > version is required.")
	// 	errorCnt++
	// } else if !IsSupportK8sVersion(k8sVersion) {
	// 	logger.Fatal(fmt.Sprintf("kubernetes > supported version: %v", conf.SupportK8SVersion))
	// 	errorCnt++
	// }

	if nodePoolSecuritySSHUserID == "" {
		logger.Fatal("node-pool.security > ssh-user-id is required.")
		errorCnt++
	}

	if nodePoolSecurityPrivateKeyPath == "" {
		logger.Fatal("node-pool.security > private-key-path is required.")
		errorCnt++
	}

	if nodePoolMasterLbIP == "" {
		logger.Fatal("node-pool.master > lb-ip is required.")
		errorCnt++
	}

	if len(kubernetesPodCidr) > 0 {
		//todo check cider
	}
	if len(kubernetesServiceCidr) > 0 {
		//todo check cider
	}

	if koreonToml.Kubernetes.Etcd.ExternalEtcd {
		switch etcdCnt {
		case 1, 3, 5:
		default:
			logger.Fatal("Only odd number of etcd nodes are supported.(1, 3, 5)")
			errorCnt++

		}
	}

	if len(nodePoolDataDir) > 0 {
		// todo node pool data dir check
	}

	//storage check
	cnt := checkSharedStorage(koreonToml)
	errorCnt += cnt

	if privateRegistryInstall == true {

		if privateRegistryRegistryIP == "" {
			logger.Fatal("private-registry > registry-ip is required.")
			errorCnt++
		}

		if privateRegistryDataDir == "" {
			logger.Fatal("private-registry > data-dir is required.")
			errorCnt++
		}

		if isPrivateRegistryPublicCert {
			if privateRegistryCrt == "" {
				logger.Fatal("private-registry.cert-file > ssl-certificate is required.")
				errorCnt++
			}

			if privateRegistryKey == "" {
				logger.Fatal("private-registry.cert-file > ssl-certificate-key is required.")
				errorCnt++
			}

			if privateRegistryRegistryDomain == "" {
				logger.Fatal("private-registry > registry-domain is required.")
				errorCnt++
			}
		}
	}

	if koreonToml.KoreOn.ClosedNetwork {
		if koreonToml.KoreOn.LocalRepository == "" && koreonToml.KoreOn.LocalRepositoryArchiveFile == "" {
			logger.Fatal("koreon> local-repository or local-repository-archive-file is required.")
			errorCnt++
		}

		if privateRegistryInstall {
			if koreonToml.PrivateRegistry.RegistryArchiveFile == "" {
				logger.Fatal("private-registry >  registry-archive-file is required.")
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

func checkSharedStorage(koreonToml model.KoreOnToml) int {
	errorCnt := 0

	if koreonToml.SharedStorage.Install == true {

		if koreonToml.SharedStorage.VolumeDir == "" {
			logger.Fatal("shared-storage > volume-dir is required.")
			errorCnt++
		}

		if koreonToml.SharedStorage.VolumeSize < 10 {
			logger.Fatal("shared-storage > volume-size is 10 or more.")
			errorCnt++
		}

		if koreonToml.SharedStorage.StorageIP == "" {
			logger.Fatal("shared-storage > storage-ip is required.")
			errorCnt++
		}

	}

	return errorCnt
}
