package conf

var Version = "unknown_version"
var CommitId = "unknown_commitid"
var BuildDate = "unknown_builddate"

const (
	KoreonImageName      = "koreon"
	KoreonImage          = "regi.k3.acornsoft.io/k3lab/koreon:1.1.3"
	KoreonKubeConfigPath = "/etc/kubernetes/acloud"
	KoreonKubeConfig     = "acloud-client-kubeconfig"
	KoreonConfigFile     = "koreon.toml"
	KoreonBasicYaml      = "basic.yaml"
	KoreonInventoryIni   = "inventory.ini"
	KoreonDestDir        = ".koreon"

	CreateYaml        = "/koreon/scripts/cluster.yml"
	AddNodeYaml       = "/koreon/scripts/add-node.yml"
	RemoveNodeYaml    = "/koreon/scripts/remove-node.yml"
	UpgradeYaml       = "/koreon/scripts/upgrade.yml"
	ResetYaml         = "/koreon/scripts/reset.yml"
	PreDestroyYaml    = "/koreon/scripts/pre-destroy.yml"
	Inventory         = "/koreon/inventory/sample"
	InventoryIni      = "/koreon/inventory/sample/inventory.ini"
	PrepareAirgapYaml = "/koreon/scripts/prepare-repository.yml"
	BasicYaml         = "/koreon/inventory/sample/group_vars/all/basic.yml"
	WorkDir           = "/koreon/work"
	SshPort           = 22
)

var SupportK8SVersion = []string{
	"1.19.10", "1.19.11", "1.19.12",
	"1.20.6", "1.20.7", "1.20.8",
	"1.21.0", "1.21.1", "1.21.2",
}

const (
	SUCCESS_FORMAT = "\033[1;32m%s\033[0m\n"
	STATUS_FORMAT  = "\033[1;32m%s\033[0m"
	ERROR_FORMAT   = "\x1B[1;3;31m%s\x1B[0m\n"
	CHECK_FORMAT   = "\033[1;34m%s\033[0m"
)

const (
	CMD_INIT            = "init"
	CMD_CREATE          = "create"
	CMD_APPLY           = "apply"
	CMD_DESTROY         = "destroy"
	CMD_VERSION         = "version"
	CMD_PREPARE_AIREGAP = "prepare-airgap"
)

const (
	RepoFile       = "repo-backup.tgz"
	HarborFile     = "harbor-backup.tgz"
	SSLRegistryCrt = "harbor.crt"
	SSLRegistryKey = "harbor.key"
	IdRsa          = "id_rsa"
	DockerBin      = "/usr/local/bin/docker"
)
