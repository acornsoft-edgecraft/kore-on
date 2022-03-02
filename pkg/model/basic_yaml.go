package model

type BasicYaml struct {
	Provider      bool   `yaml:"provider"`
	CloudProvider string `yaml:"cloud_provider"`
	ClusterName   string `yaml:"cluster_name"`

	//# install directories
	InstallDir  string `yaml:"install_dir"`
	DataRootDir string `yaml:"data_root_dir"`

	//# kubernetes options
	K8SVersion     string   `yaml:"k8s_version" comment:"kubernetes options"`
	ClusterID      string   `yaml:"cluster_id"`
	APILbIP        string   `yaml:"api_lb_ip"`
	LbIP           string   `yaml:"lb_ip"`
	LbPort         int      `yaml:"lb_port"`
	PodIPRange     string   `yaml:"pod_ip_range"`
	ServiceIPRange string   `yaml:"service_ip_range"`
	NodePortRange  string   `yaml:"node_port_range"`
	ApiSans        []string `yaml:"api_sans"`

	//# for air gap installation
	ClosedNetwork              bool   `yaml:"closed_network"`
	LocalRepository            string `yaml:"local_repository"`
	LocalRepositoryArchiveFile string `yaml:"local_repository_archive_file"`

	//# option for master isolation
	MasterIsolated   bool `yaml:"master_isolated"`
	AuditLogEnable   bool `yaml:"audit_log_enable"`
	CertValidityDays int  `yaml:"cert_validity_days"`

	//# container runtime [containerd | docker]
	ContainerRuntime string `yaml:"container_runtime"`
	KubeProxyMode    string `yaml:"kube_proxy_mode"`

	//# kube-proxy mode [iptables | ipvs]
	RegistryInstall     bool   `yaml:"registry_install"`
	RegistryDataDir     string `yaml:"registry_data_dir"`
	Registry            string `yaml:"registry"`
	RegistryDomain      string `yaml:"registry_domain"`
	RegistryPublicCert  bool   `yaml:"registry_public_cert"`
	RegistryArchiveFile string `yaml:"registry_archive_file"`

	//# option for harbor registry
	StorageInstall bool `yaml:"storage_install"`

	//# option for NFS storage
	NfsIP        string `yaml:"nfs_ip"`
	NfsVolumeDir string `yaml:"nfs_volume_dir"`

	//# for internal load-balancer
	Haproxy bool `yaml:"haproxy"`

	//# Calico network mode
	VxlanMode bool `yaml:"vxlan_mode"`

	//# option for preparing local-repo and registry (do not modify when fully understand this flag)
	ArchiveRepo bool `yaml:"archive_repo"`
}
