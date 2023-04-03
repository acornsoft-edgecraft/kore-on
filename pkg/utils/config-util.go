package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kore-on/pkg/logger"
	"kore-on/pkg/model"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
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
	// default values
	koreonToml.NodePool.Master.HaproxyInstall = true
	koreonToml.Kubernetes.AuditLogEnable = true

	err = toml.Unmarshal(c, &koreonToml)
	if err != nil {
		logger.Fatal(err.Error())
		errorCnt++
	}

	return koreonToml, err
}

func GetAddonTomlConfig(path string) (model.AddonToml, error) {

	errorCnt = 0
	// configFullPath := workDir + "/" + conf.KoreonConfigFile

	var c []byte
	var err error

	if !FileExists(path) {
		logger.Fatal(fmt.Sprintf("%s file is not found. Run koreonctl addon init first", path))
		os.Exit(1)
	}

	c, err = ioutil.ReadFile(path)
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}

	str := string(c)
	str = strings.Replace(str, "\\", "/", -1)
	c = []byte(str)

	var addonToml = model.AddonToml{}
	err = toml.Unmarshal(c, &addonToml)
	if err != nil {
		logger.Fatal(err.Error())
		errorCnt++
	}

	return addonToml, err
}

func ValidateKoreonTomlConfig(koreOnConfigFilePath string, cmd string) (model.KoreOnToml, bool) {
	var koreon_toml model.KoreOnToml
	errorCnt = 0
	koreonToml, _ := GetKoreonTomlConfig(koreOnConfigFilePath)

	confK8sVersion := "SupportK8sVersion"
	confHarborVersion := "SupportHarborVersion"
	nodePoolSSHPort := koreonToml.NodePool.SSHPort
	// confCalicoVersion := "Support.SupportCalicoVersion"
	// confCorednsVersion := "Support.SupportCorednsVersion"
	// confDockerVersion := "Support.SupportDockerVersion"
	// confDockerComposeVersion := "Support.SupportDockerComposeVersion"
	koreonToml.KoreOn.ImageArchive = viper.GetString("KoreOn.KoreOnImageArchive")
	koreonToml.KoreOn.HelmCubeRepoUrl = viper.GetString("KoreOn.HelmCubeRepoUrl")

	if nodePoolSSHPort == 0 {
		// todo node pool ssh port check
		koreonToml.NodePool.SSHPort = 22
	}

	if cmd == "prepare-airgap" {
		k8sVersion := koreonToml.PrepareAirgap.K8sVersion
		registryIP := koreonToml.PrepareAirgap.RegistryIP
		registryVersion := koreonToml.PrepareAirgap.RegistryVersion
		koreon_toml.KoreOn.HelmCubeRepoUrl = viper.GetString("KoreOn.HelmCubeRepoUrl")

		supportK8sVersion := IsSupportVersion(k8sVersion, confK8sVersion)
		supportHarborVersion := IsSupportVersion(registryVersion, confHarborVersion)
		if registryIP == "" {
			logger.Fatal(fmt.Sprintln("Prepare Air Gap > Registry IP Address is required."))
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
		// Get helm chart package support version
		supportHelmChartList := GetSupportVersion(supportK8sVersion, "helm_chart_package")
		if supportPackageList == nil {
			logger.Fatal("Prepare Air Gap > Support helm chart package version not found.:\n")
			errorCnt++
		}

		// Set image support version
		k8sSupportImagesVersion, err := setField(&koreon_toml.SupportVersion.ImageVersion, supportK8sList)
		if err != nil {
			logger.Fatal(err)
			errorCnt++
		} else if err := json.Unmarshal(k8sSupportImagesVersion, &koreon_toml.ListVersion); err != nil {
			logger.Fatal(err)
		}

		// Set package support version
		packageSupportVersion, err := setField(&koreon_toml.SupportVersion.PackageVersion, supportPackageList)
		if err != nil {
			logger.Fatal(err)
			errorCnt++
		} else if err := json.Unmarshal(packageSupportVersion, &koreon_toml.ListVersion); err != nil {
			logger.Fatal(err)
		}

		// Set package support version
		helmChartSupportVersion, err := setField(&koreon_toml.SupportVersion.HelmChartVersion, supportHelmChartList)
		if err != nil {
			logger.Fatal(err)
			errorCnt++
		} else if err := json.Unmarshal(helmChartSupportVersion, &koreon_toml.ListVersion); err != nil {
			logger.Fatal(err)
		}

		koreonToml = koreon_toml
		koreonToml.KoreOn.ImageArchive = viper.GetString("KoreOn.KoreOnImageArchive")
		// koreonToml.SupportVersion.ImageVersion.Calico = IsSupportVersion(fmt.Sprintf("%v", supportK8sList["calico"]), confCalicoVersion)
		// koreonToml.SupportVersion.ImageVersion.Coredns = IsSupportVersion(fmt.Sprintf("%v", supportK8sList["coredns"]), confCorednsVersion)

		// koreonToml.SupportVersion.PackageVersion.Docker = IsSupportVersion(fmt.Sprintf("%v", supportHarborList["docker"]), confDockerVersion)
		// koreonToml.SupportVersion.PackageVersion.DockerCompose = IsSupportVersion(fmt.Sprintf("%v", supportHarborList["docker-compose"]), confDockerComposeVersion)

	} else if cmd == "create" {
		kubernetesPodCidr := koreonToml.Kubernetes.PodCidr
		kubernetesServiceCidr := koreonToml.Kubernetes.ServiceCidr
		k8sVersion := koreonToml.Kubernetes.Version
		etcdCnt := len(koreonToml.Kubernetes.Etcd.IP)
		masterIP := koreonToml.NodePool.Master.IP
		workerIP := koreonToml.NodePool.Node.IP
		etcdPrivateIpCnt := len(koreonToml.Kubernetes.Etcd.PrivateIP)
		nodePoolDataDir := koreonToml.NodePool.DataDir

		privateRegistryInstall := koreonToml.PrivateRegistry.Install
		privateRegistryRegistryIP := koreonToml.PrivateRegistry.RegistryIP
		privateRegistryRegistryVersion := koreonToml.PrivateRegistry.RegistryVersion
		privateRegistryRegistryDomain := koreonToml.PrivateRegistry.RegistryDomain
		isPrivateRegistryPublicCert := koreonToml.PrivateRegistry.PublicCert
		privateRegistryCrt := koreonToml.PrivateRegistry.CertFile.CaCert
		// privateRegistryKey := koreonToml.PrivateRegistry.CertFile.CaCertKey

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

		if len(masterIP) < 0 {
			logger.Fatal("NodePool > K8s Control Plane node is required.")
		} else {
			//todo check masterIP
			if len(workerIP) < 0 {
				logger.Fatal("NodePool > K8s Worker node is required.")
			}
		}

		if len(kubernetesPodCidr) > 0 {
			//todo check cider
		}

		if len(kubernetesServiceCidr) > 0 {
			//todo check cider
		}

		if koreonToml.Kubernetes.Etcd.ExternalEtcd {
			if etcdPrivateIpCnt == 0 {
				koreonToml.Kubernetes.Etcd.PrivateIP = koreonToml.Kubernetes.Etcd.IP
			}
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

			if privateRegistryRegistryDomain == "" {
				koreonToml.PrivateRegistry.RegistryDomain = privateRegistryRegistryIP
			}
		}

		supportHarborVersion := IsSupportVersion(privateRegistryRegistryVersion, confHarborVersion)
		if privateRegistryRegistryVersion == "" {
			koreonToml.PrivateRegistry.RegistryVersion = supportHarborVersion
			logger.Warn("Private Registry > Harbor version is required. Last version", koreonToml.PrivateRegistry.RegistryVersion, "applied automatically.")
		} else {
			koreonToml.PrivateRegistry.RegistryVersion = supportHarborVersion
		}

		if isPrivateRegistryPublicCert {
			if privateRegistryCrt == "" {
				logger.Fatal("private-registry.cert-file > ca-cert is required.")
			}

			// if privateRegistryKey == "" {
			// 	logger.Fatal("private-registry.cert-file > ssl-certificate-key is required.")
			// }
		}

		if koreonToml.KoreOn.ClosedNetwork {
			if koreonToml.KoreOn.LocalRepositoryInstall {
				if koreonToml.KoreOn.LocalRepositoryArchiveFile == "" {
					logger.Fatal("koreon> When installing a local repository, the local-repository-archive-file entry is required.")
				} else {
					localRepositoryArchiveFile := filepath.Base(koreonToml.KoreOn.LocalRepositoryArchiveFile)
					k8sVersionCheck := strings.Split(localRepositoryArchiveFile, "-")
					if supportK8sVersion != k8sVersionCheck[2] {
						logger.Fatalf("Check the kubernetes installation version.\nIs the version you are trying to install '%s' correct? If different, re-enter the kubernetes.version entry", k8sVersionCheck[2])
					}
					koreonToml.KoreOn.LocalRepositoryArchiveFile = localRepositoryArchiveFile
				}
			} else {
				if koreonToml.KoreOn.LocalRepositoryUrl == "" {
					logger.Fatal("koreon> If you are not installing a local repository, the local-repository-url entry is required.")
				}
				if koreonToml.KoreOn.LocalRepositoryArchiveFile != "" {
					logger.Fatal("koreon> If you are not installing a local repository, the local-repository-archive-file entry should be empty.")
				}
			}

			if privateRegistryInstall {
				if koreonToml.PrivateRegistry.RegistryArchiveFile == "" {
					logger.Fatal("private-registry >  registry-archive-file is required.")
				} else {
					registryArchiveFile := filepath.Base(koreonToml.PrivateRegistry.RegistryArchiveFile)
					harborVersionCheck := strings.Split(registryArchiveFile, "-")
					if supportHarborVersion != harborVersionCheck[1] {
						logger.Fatalf("Check the private registry installation version.\nIs the version you are trying to install '%s' correct? If different, re-enter the registry-archive-file entry", harborVersionCheck[1])
					}
					koreonToml.PrivateRegistry.RegistryArchiveFile = registryArchiveFile
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
		k8sSupportImagesVersion, err := setField(&koreonToml.SupportVersion.ImageVersion, supportK8sList)
		if err != nil {
			logger.Fatal(err)
			errorCnt++
		} else if err := json.Unmarshal(k8sSupportImagesVersion, &koreon_toml.ListVersion); err != nil {
			logger.Fatal(err)
		}

		// Set package support version
		packageSupportVersion, err := setField(&koreonToml.SupportVersion.PackageVersion, supportPackageList)
		if err != nil {
			logger.Fatal(err)
			errorCnt++
		} else if err := json.Unmarshal(packageSupportVersion, &koreon_toml.ListVersion); err != nil {
			logger.Fatal(err)
		}

		koreonToml.PrepareAirgap = koreon_toml.PrepareAirgap
	} else if cmd == "cluster-update" {
		kubernetesPodCidr := koreonToml.Kubernetes.PodCidr
		kubernetesServiceCidr := koreonToml.Kubernetes.ServiceCidr
		k8sVersion := koreonToml.Kubernetes.Version
		workerIP := koreonToml.NodePool.Node.IP
		nodePoolDataDir := koreonToml.NodePool.DataDir

		// privateRegistryInstall := koreonToml.PrivateRegistry.Install
		// privateRegistryRegistryIP := koreonToml.PrivateRegistry.RegistryIP
		// privateRegistryRegistryVersion := koreonToml.PrivateRegistry.RegistryVersion
		// privateRegistryRegistryDomain := koreonToml.PrivateRegistry.RegistryDomain
		// isPrivateRegistryPublicCert := koreonToml.PrivateRegistry.PublicCert
		// privateRegistryCrt := koreonToml.PrivateRegistry.CertFile.CaCert
		// privateRegistryKey := koreonToml.PrivateRegistry.CertFile.CaCertKey

		supportK8sVersion := IsSupportVersion(k8sVersion, confK8sVersion)

		// if koreonToml.KoreOn.InstallDir != "" && !strings.HasPrefix(koreonToml.KoreOn.InstallDir, "/") {
		// 	logger.Fatal("koreon > install-dir is Only absolute paths are supported.")
		// 	errorCnt++
		// }

		if k8sVersion == "" {
			koreonToml.Kubernetes.Version = supportK8sVersion
			logger.Warn("kubernetes > Kubernetes version is required. Last version", koreonToml.Kubernetes.Version, "applied automatically.")
		} else {
			koreonToml.Kubernetes.Version = supportK8sVersion
		}

		if len(workerIP) < 0 {
			logger.Fatal("NodePool > K8s Worker node is required.")
		}

		if len(kubernetesPodCidr) > 0 {
			//todo check cider
		}

		if len(kubernetesServiceCidr) > 0 {
			//todo check cider
		}

		if len(nodePoolDataDir) > 0 {
			// todo node pool data dir check
		}

		if koreonToml.KoreOn.ClosedNetwork {
			if koreonToml.KoreOn.LocalRepositoryInstall {
				if koreonToml.KoreOn.LocalRepositoryArchiveFile == "" {
					logger.Fatal("koreon> When installing a local repository, the local-repository-archive-file entry is required.")
				} else {
					localRepositoryArchiveFile := filepath.Base(koreonToml.KoreOn.LocalRepositoryArchiveFile)
					k8sVersionCheck := strings.Split(localRepositoryArchiveFile, "-")
					if supportK8sVersion != k8sVersionCheck[2] {
						logger.Fatalf("Check the kubernetes installation version.\nIs the version you are trying to install '%s' correct? If different, re-enter the kubernetes.version entry", k8sVersionCheck[2])
					}
					koreonToml.KoreOn.LocalRepositoryArchiveFile = localRepositoryArchiveFile
				}
			} else {
				if koreonToml.KoreOn.LocalRepositoryUrl == "" {
					logger.Fatal("koreon> If you are not installing a local repository, the local-repository-url entry is required.")
				}
				if koreonToml.KoreOn.LocalRepositoryArchiveFile != "" {
					logger.Fatal("koreon> If you are not installing a local repository, the local-repository-archive-file entry should be empty.")
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
		k8sSupportImagesVersion, err := setField(&koreonToml.SupportVersion.ImageVersion, supportK8sList)
		if err != nil {
			logger.Fatal(err)
			errorCnt++
		} else if err := json.Unmarshal(k8sSupportImagesVersion, &koreon_toml.ListVersion); err != nil {
			logger.Fatal(err)
		}

		// Set package support version
		packageSupportVersion, err := setField(&koreonToml.SupportVersion.PackageVersion, supportPackageList)
		if err != nil {
			logger.Fatal(err)
			errorCnt++
		} else if err := json.Unmarshal(packageSupportVersion, &koreon_toml.ListVersion); err != nil {
			logger.Fatal(err)
		}

		koreonToml.PrepareAirgap = koreon_toml.PrepareAirgap
	} else if cmd == "reset-prepare-airgap" {
		registryIP := koreonToml.PrepareAirgap.RegistryIP

		if registryIP == "" {
			logger.Fatal(fmt.Sprintln("Destroy: Prepare Air Gap > Registry IP Address is required."))
		} else {
			koreon_toml.PrepareAirgap = koreonToml.PrepareAirgap
		}

		koreonToml = koreon_toml
	} else if cmd == "reset-all" {
		koreonToml.PrepareAirgap.RegistryIP = ""
	} else if cmd == "reset-cluster" {
		koreonToml.PrepareAirgap.RegistryIP = ""
	} else if cmd == "reset-registry" {
		koreonToml.PrepareAirgap.RegistryIP = ""
	} else if cmd == "reset-storage" {
		koreonToml.PrepareAirgap.RegistryIP = ""
	} else if cmd == "add-on" {
		k8sVersion := koreonToml.PrepareAirgap.K8sVersion

		supportK8sVersion := IsSupportVersion(k8sVersion, confK8sVersion)

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
		k8sSupportImagesVersion, err := setField(&koreonToml.SupportVersion.ImageVersion, supportK8sList)
		if err != nil {
			logger.Fatal(k8sSupportImagesVersion)
			errorCnt++
		} else if err := json.Unmarshal(k8sSupportImagesVersion, &koreon_toml.ListVersion); err != nil {
			logger.Fatal(err)
		}

		// Set package support version
		packageSupportVersion, err := setField(&koreonToml.SupportVersion.PackageVersion, supportPackageList)
		if err != nil {
			logger.Fatal(packageSupportVersion)
			errorCnt++
		} else if err := json.Unmarshal(packageSupportVersion, &koreon_toml.ListVersion); err != nil {
			logger.Fatal(err)
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
	errorCnt = 0

	if koreonToml.SharedStorage.Install == true {
		if koreonToml.SharedStorage.StorageIP == "" {
			logger.Fatal("shared-storage > storage-ip is required.")
			errorCnt++
		}
	}

	return errorCnt
}

func setField(item interface{}, supportList map[string]interface{}) ([]byte, error) {
	v := reflect.ValueOf(item).Elem()
	if !v.CanAddr() {
		return nil, fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}

	result := make(map[string]interface{})
	versions := make(map[string]interface{})

	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag.Get("validate")
		r := strings.Split(tag, ",")
		if len(r) != 2 {
			return nil, fmt.Errorf("tag entry error in %s field", typeField.Name)
		}
		value := IsSupportVersion(fmt.Sprintf("%v", supportList[string(r[0])]), r[1])
		v.Field(i).SetString(value)

		// list Versions
		versions[typeField.Name] = ListSupportVersion(r[1])
		result["List"+v.Type().Name()] = versions
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return data, nil
}
