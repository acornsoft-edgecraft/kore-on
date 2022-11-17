package utils

import (
	"fmt"
	"io/ioutil"
	"kore-on/pkg/logger"
	"kore-on/pkg/model"
	"os"
	"reflect"
	"strings"

	"github.com/pelletier/go-toml"
)

var errorCnt = 0

func GetKoreonTomlConfig(koreOnConfigFilePath string) (model.KoreOnToml, error) {

	errorCnt = 0
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

func ValidateKoreonTomlConfig(koreOnConfigFilePath string, cmd string) (model.KoreOnToml, bool) {
	var koreon_toml model.KoreOnToml
	errorCnt = 0
	koreonToml, _ := GetKoreonTomlConfig(koreOnConfigFilePath)

	confK8sVersion := "SupportK8sVersion"
	confHarborVersion := "SupportHarborVersion"
	// confCalicoVersion := "Support.SupportCalicoVersion"
	// confCorednsVersion := "Support.SupportCorednsVersion"
	// confDockerVersion := "Support.SupportDockerVersion"
	// confDockerComposeVersion := "Support.SupportDockerComposeVersion"

	if cmd == "prepare-airgap" {
		k8sVersion := koreonToml.PrepareAirgap.K8sVersion
		registryIP := koreonToml.PrepareAirgap.RegistryIP
		registryVersion := koreonToml.PrepareAirgap.RegistryVersion

		supportK8sVersion := IsSupportVersion(k8sVersion, confK8sVersion)
		supportHarborVersion := IsSupportVersion(registryVersion, confHarborVersion)
		if registryIP == "" {
			logger.Fatal(fmt.Sprintf("Prepare Air Gap > Kubernetes version is required.\nK8s supported version lists:\n %v", ListSupportVersion(confK8sVersion)))
			errorCnt++
		} else {
			koreon_toml.PrepareAirgap.RegistryIP = registryIP
		}

		if k8sVersion == "" {
			koreon_toml.PrepareAirgap.K8sVersion = supportK8sVersion
			logger.Warn("Prepare Air Gap > Kubernetes version is required. Last version", koreonToml.PrepareAirgap.K8sVersion, "applied automatically.")
		} else {
			koreon_toml.PrepareAirgap.K8sVersion = supportK8sVersion
		}

		if registryVersion == "" {
			koreon_toml.PrepareAirgap.RegistryVersion = supportHarborVersion
			logger.Warn("Prepare Air Gap > Harbor version is required. Last version", koreonToml.PrepareAirgap.RegistryVersion, "applied automatically.")
		} else {
			koreon_toml.PrepareAirgap.RegistryVersion = supportHarborVersion
		}

		// Get image support version
		supportK8sList := GetSupportVersion(supportK8sVersion, "k8s_support_image")
		if supportK8sList == nil {
			logger.Fatal("Prepare Air Gap > Support package and container image not found.:\n")
			errorCnt++
		}
		// Get package support version
		supportPackageList := GetSupportVersion(supportK8sVersion, "k8s_support_package")
		if supportPackageList == nil {
			logger.Fatal("Prepare Air Gap > Support package and container image not found.:\n")
			errorCnt++
		}

		// Set image support version
		k8sSupportImagesVersion := setField(&koreon_toml.SupportVersion.ImageVersion, supportK8sList)
		if k8sSupportImagesVersion != nil {
			logger.Fatal(k8sSupportImagesVersion)
			errorCnt++
		}

		// Set package support version
		packageSupportVersion := setField(&koreon_toml.SupportVersion.PackageVersion, supportPackageList)
		if packageSupportVersion != nil {
			logger.Fatal(packageSupportVersion)
			errorCnt++
		}

		koreonToml = koreon_toml
		// koreonToml.SupportVersion.ImageVersion.Calico = IsSupportVersion(fmt.Sprintf("%v", supportK8sList["calico"]), confCalicoVersion)
		// koreonToml.SupportVersion.ImageVersion.Coredns = IsSupportVersion(fmt.Sprintf("%v", supportK8sList["coredns"]), confCorednsVersion)

		// koreonToml.SupportVersion.PackageVersion.Docker = IsSupportVersion(fmt.Sprintf("%v", supportHarborList["docker"]), confDockerVersion)
		// koreonToml.SupportVersion.PackageVersion.DockerCompose = IsSupportVersion(fmt.Sprintf("%v", supportHarborList["docker-compose"]), confDockerComposeVersion)

	} else if cmd == "create" {
		kubernetesPodCidr := koreonToml.Kubernetes.PodCidr
		kubernetesServiceCidr := koreonToml.Kubernetes.ServiceCidr
		k8sVersion := koreonToml.Kubernetes.Version
		etcdCnt := len(koreonToml.Kubernetes.Etcd.IP)
		etcdPrivateIpCnt := len(koreonToml.Kubernetes.Etcd.PrivateIP)
		nodePoolDataDir := koreonToml.NodePool.DataDir

		privateRegistryInstall := koreonToml.PrivateRegistry.Install
		privateRegistryRegistryIP := koreonToml.PrivateRegistry.RegistryIP
		privateRegistryRegistryVersion := koreonToml.PrivateRegistry.RegistryVersion
		privateRegistryRegistryDomain := koreonToml.PrivateRegistry.RegistryDomain
		privateRegistryDataDir := koreonToml.PrivateRegistry.DataDir
		isPrivateRegistryPublicCert := koreonToml.PrivateRegistry.PublicCert
		privateRegistryCrt := koreonToml.PrivateRegistry.CertFile.SslCertificate
		privateRegistryKey := koreonToml.PrivateRegistry.CertFile.SslCertificateKey

		supportK8sVersion := IsSupportVersion(k8sVersion, confK8sVersion)

		if koreonToml.KoreOn.InstallDir != "" && !strings.HasPrefix(koreonToml.KoreOn.InstallDir, "/") {
			logger.Fatal("koreon > install-dir is Only absolute paths are supported.")
			errorCnt++
		}

		if k8sVersion == "" {
			koreonToml.Kubernetes.Version = supportK8sVersion
			logger.Warn("kubernetes > Kubernetes version is required. Last version", koreonToml.Kubernetes.Version, "applied automatically.")
		} else {
			koreonToml.Kubernetes.Version = supportK8sVersion
		}

		// Get image support version
		// confCalicoVersion := "KoreOn.SupportCalicoVersion"
		// calicoVersion := GetSupportVersion(supportK8sVersion, "k8s")
		// supportCalicoVersion := IsSupportVersion(calicoVersion, confCalicoVersion)

		// if calicoVersion == "" {
		// 	koreonToml.Kubernetes.Calico.Version = supportCalicoVersion
		// 	logger.Warn("kubernetes > Calico version is required. Last version", koreonToml.Kubernetes.Calico.Version, "applied automatically.")
		// } else if supportCalicoVersion == "" {
		// 	logger.Fatal(fmt.Sprintf("kubernetes > Calico supported version lists:\n %v", ListSupportVersion(confCalicoVersion)))
		// 	errorCnt++
		// } else {
		// 	koreonToml.Kubernetes.Calico.Version = supportCalicoVersion
		// }

		// if nodePoolSecuritySSHUserID == "" {
		// 	logger.Fatal("node-pool.security > ssh-user-id is required.")
		// 	errorCnt++
		// }

		// if nodePoolSecurityPrivateKeyPath == "" {
		// 	logger.Fatal("node-pool.security > private-key-path is required.")
		// 	errorCnt++
		// }

		// if nodePoolMasterLbIP == "" {
		// 	logger.Fatal("node-pool.master > lb-ip is required.")
		// 	errorCnt++
		// }

		if len(kubernetesPodCidr) > 0 {
			//todo check cider
		}
		if len(kubernetesServiceCidr) > 0 {
			//todo check cider
		}

		if koreonToml.Kubernetes.Etcd.ExternalEtcd {
			if etcdCnt != etcdPrivateIpCnt && etcdCnt > 0 && etcdPrivateIpCnt > 0 {
				logger.Fatal("etcd nodes IP address and private ip address needs")
				errorCnt++
			}
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

		if privateRegistryInstall {

			if privateRegistryRegistryIP == "" {
				logger.Fatal("private-registry > registry-ip is required.")
				errorCnt++
			}

			supportHarborVersion := IsSupportVersion(privateRegistryRegistryVersion, confHarborVersion)
			if privateRegistryRegistryVersion == "" {
				koreonToml.PrivateRegistry.RegistryVersion = supportHarborVersion
				logger.Warn("Private Registry > Harbor version is required. Last version", koreonToml.PrivateRegistry.RegistryVersion, "applied automatically.")
			} else {
				koreonToml.PrivateRegistry.RegistryVersion = supportHarborVersion
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
			if koreonToml.KoreOn.LocalRepositoryInstall {
				if koreonToml.KoreOn.LocalRepositoryArchiveFile == "" {
					logger.Fatal("koreon> When installing a local repository, the local-repository-archive-file entry is required.")
					errorCnt++
				}
			} else {
				if koreonToml.KoreOn.LocalRepositoryUrl == "" {
					logger.Fatal("koreon> If you are not installing a local repository, the local-repository-url entry is required.")
					errorCnt++
				}
				if koreonToml.KoreOn.LocalRepositoryArchiveFile != "" {
					logger.Fatal("koreon> If you are not installing a local repository, the local-repository-archive-file entry should be empty.")
					errorCnt++
				}
			}

			if privateRegistryInstall {
				if koreonToml.PrivateRegistry.RegistryArchiveFile == "" {
					logger.Fatal("private-registry >  registry-archive-file is required.")
					errorCnt++
				}
			}
		}

		// Get image support version
		supportK8sList := GetSupportVersion(supportK8sVersion, "k8s_support_image")
		if supportK8sList == nil {
			logger.Fatal("Prepare Air Gap > Support package and container image not found.:\n")
			errorCnt++
		}
		// Get package support version
		supportPackageList := GetSupportVersion(supportK8sVersion, "k8s_support_package")
		if supportPackageList == nil {
			logger.Fatal("Prepare Air Gap > Support package and container image not found.:\n")
			errorCnt++
		}

		// Set image support version
		k8sSupportImagesVersion := setField(&koreonToml.SupportVersion.ImageVersion, supportK8sList)
		if k8sSupportImagesVersion != nil {
			logger.Fatal(k8sSupportImagesVersion)
			errorCnt++
		}

		// Set package support version
		packageSupportVersion := setField(&koreonToml.SupportVersion.PackageVersion, supportPackageList)
		if packageSupportVersion != nil {
			logger.Fatal(packageSupportVersion)
			errorCnt++
		}

		koreonToml.PrepareAirgap = koreon_toml.PrepareAirgap
	}

	if errorCnt > 0 {
		logger.Error("there are one or more errors")
		os.Exit(1)
		return koreonToml, false
	}
	return koreonToml, true
}

func checkSharedStorage(koreonToml model.KoreOnToml) int {
	errorCnt = 0

	if koreonToml.SharedStorage.Install == true {

		if koreonToml.SharedStorage.VolumeDir == "" {
			logger.Fatal("shared-storage > volume-dir is required.")
			errorCnt++
		}

		// if koreonToml.SharedStorage.VolumeSize < 10 {
		// 	logger.Fatal("shared-storage > volume-size is 10 or more.")
		// 	errorCnt++
		// }

		if koreonToml.SharedStorage.StorageIP == "" {
			logger.Fatal("shared-storage > storage-ip is required.")
			errorCnt++
		}

	}

	return errorCnt
}

func setField(item interface{}, supportList map[string]interface{}) error {
	v := reflect.ValueOf(item).Elem()
	if !v.CanAddr() {
		return fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}

	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag.Get("validate")
		r := strings.Split(tag, ",")
		if len(r) != 2 {
			return fmt.Errorf("tag entry error in %s field", typeField.Name)
		}
		value := IsSupportVersion(fmt.Sprintf("%v", supportList[string(r[0])]), r[1])
		v.Field(i).SetString(value)
	}
	return nil
}
