package conf

const (
	KoreOnImageName        = "kore-on"
	KoreOnImage            = "ghcr.io/acornsoft-edgecraft/kore-on:latest"
	KoreOnImageArchive     = "koreon-image.tgz"
	KoreOnKubeConfig       = "acloud-client-kubeconfig"
	KoreOnConfigFile       = "koreon.toml"
	AddOnConfigFile        = "addon.toml"
	KoreOnConfigFileSubDir = "internal/playbooks/koreon-playbook/download"
	HelmCubeRepoUrl        = "https://hcapital-harbor.acloud.run/chartrepo/cube"
)

var Addon = map[string]string{
	"KubeConfigDir": "/etc/kubernetes/acloud",
}
