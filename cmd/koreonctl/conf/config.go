package conf

var (
	KoreOnVersion          = "latest"
	KoreOnImageName        = "kore-on"
	KoreOnRegistry         = "ghcr.io/acornsoft-edgecraft"
	KoreOnImage            = KoreOnRegistry + "/" + KoreOnImageName + ":" + KoreOnVersion
	KoreOnImageArchive     = KoreOnImageName + "_" + KoreOnVersion + ".tar.gz"
	KoreOnKubeConfig       = "acloud-client-kubeconfig"
	KoreOnConfigFile       = "koreon.toml"
	AddOnConfigFile        = "addon.toml"
	KoreOnConfigFileSubDir = "internal/playbooks/koreon-playbook/download"
	KoreOnArchiveFileDir   = "internal/playbooks/koreon-playbook/download/archive"
	KoreOnConfigDir        = "internal/playbooks/koreon-playbook/download/config"
	KoreOnExtendsFileDir   = "internal/playbooks/koreon-playbook/download/extends"
	KoreOnLogsDir          = "internal/playbooks/koreon-playbook/download/logs"
	HelmCubeRepoUrl        = "https://hcapital-harbor.acloud.run/chartrepo/cube"
	HelmChartProject       = "helm-charts"
)

var Addon = map[string]string{
	"KubeConfigDir":   "/etc/kubernetes/acloud",
	"AddonConfigFile": "addon.toml",
}
