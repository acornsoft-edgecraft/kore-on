package conf

var (
	KoreOnVersion          = "latest"
	KoreOnImageName        = "kore-on"
	KoreOnImage            = "ghcr.io/acornsoft-edgecraft/kore-on:" + KoreOnVersion
	KoreOnImageArchive     = "koreon-image.tgz"
	KoreOnKubeConfig       = "acloud-client-kubeconfig"
	KoreOnConfigFile       = "config/koreon.toml"
	AddOnConfigFile        = "config/addon.toml"
	KoreOnConfigFileSubDir = "internal/playbooks/koreon-playbook/download"
	KoreOnConfigDir        = "internal/playbooks/koreon-playbook/download/config"
	KoreOnArchiveFileDir   = "internal/playbooks/koreon-playbook/download/archive"
	KoreOnLogsDir          = "internal/playbooks/koreon-playbook/download/logs"
	HelmCubeRepoUrl        = "https://hcapital-harbor.acloud.run/chartrepo/cube"
	HelmChartProject       = "helm-charts"
)

var Addon = map[string]string{
	"KubeConfigDir":   "/etc/kubernetes/acloud",
	"AddonConfigFile": "config/addon.toml",
}
