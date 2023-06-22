package model

type KoreOnToml struct {
	KoreOn struct {
		ClusterInstall   bool   `toml:"cluster-install,omitempty"`
		ClusterName      string `toml:"cluster-name,omitempty"`
		ClusterID        string `toml:"cluster-id,omitempty"`
		InstallDir       string `toml:"install-dir,omitempty"`
		CertValidityDays int    `toml:"cert-validity-days,omitempty"`
		FileName         string
		ImageArchivePath string
		HelmCubeRepoUrl  string
		HelmCubeRepoID   string
		HelmCubeRepoPW   string
		HelmChartProject string
		Update           bool
		Create           bool
		Upgrade          bool
		CommandMode      string

		//#Airgap
		ClosedNetwork              bool   `toml:"closed-network,omitempty"`
		LocalRepositoryInstall     bool   `toml:"local-repository-install,omitempty"`
		LocalRepositoryPort        int    `toml:"local-repository-port,omitempty"`
		LocalRepositoryUrl         string `toml:"local-repository-url,omitempty"`
		LocalRepositoryArchiveFile string `toml:"local-repository-archive-file"`
		DebugMode                  bool   `toml:"debug-mode,omitempty"`
		ClusterApi                 bool   `toml:"cluster-api,omitempty"`
	} `toml:"koreon,omitempty"`

	Kubernetes struct {
		Version          string   `toml:"version,omitempty"`
		ContainerRuntime string   `toml:"container-runtime"`
		KubeProxyMode    string   `toml:"kube-proxy-mode"`
		CalicoVersion    string   `toml:"calico-version"`
		ServiceCidr      string   `toml:"service-cidr,omitempty"`
		PodCidr          string   `toml:"pod-cidr,omitempty"`
		NodePortRange    string   `toml:"node-port-range,omitempty"`
		AuditLogEnable   bool     `toml:"audit-log-enable,omitempty"`
		ApiSans          []string `toml:"api-sans,omitempty"`
		GetKubeConfig    bool

		Calico struct {
			Version   string `toml:"version,omitempty"`
			VxlanMode bool   `toml:"vxlan-mode"`
		} `toml:"calico,omitempty"`

		Etcd struct {
			ExternalEtcd  bool     `toml:"external-etcd,omitempty"`
			IP            []string `toml:"ip"`
			PrivateIP     []string `toml:"private-ip"`
			EncryptSecret bool     `toml:"encrypt-secret,omitempty"`
		} `toml:"etcd,omitempty"`
	} `toml:"kubernetes,omitempty"`

	NodePool struct {
		DataDir string `toml:"data-dir,omitempty"`
		SSHPort int    `toml:"ssh-port,omitempty"`

		// Security struct {
		// 	SSHUserID      string `toml:"ssh-user-id,omitempty"`
		// 	SSHPort        int    `toml:"ssh-port,omitempty"`
		// 	PrivateKeyPath string `toml:"private-key-path,omitempty"`
		// } `toml:"security,omitempty"`

		Master struct {
			Name           string   `toml:"name,omitempty"`
			IP             []string `toml:"ip"`
			PrivateIP      []string `toml:"private-ip"`
			LbIP           string   `toml:"lb-ip,omitempty"`
			LbPort         int      `toml:"lb-port,omitempty"`
			Isolated       bool     `toml:"isolated,omitempty"`
			HaproxyInstall bool     `toml:"haproxy-install,omitempty"`
		} `toml:"master,omitempty"`

		Node StrNode `toml:"node,omitempty"`
	} `toml:"node-pool,omitempty"`

	SharedStorage struct {
		Install    bool   `toml:"install"`
		StorageIP  string `toml:"storage-ip,omitempty"`
		PrivateIP  string `toml:"private-ip,omitempty"`
		VolumeDir  string `toml:"volume-dir,omitempty"`
		VolumeSize int    `toml:"volume-size,omitempty"`
		//StorageType       string `toml:"storage-type,omitempty"`

	} `toml:"shared-storage,omitempty"`

	PrivateRegistry struct {
		Install             bool   `toml:"install,omitempty"`
		RegistryVersion     string `toml:"registry-version,omitempty"`
		RegistryIP          string `toml:"registry-ip,omitempty"`
		RegistryDomain      string `toml:"registry-domain,omitempty"`
		PrivateIP           string `toml:"private-ip,omitempty"`
		DataDir             string `toml:"data-dir,omitempty"`
		RegistryArchiveFile string `toml:"registry-archive-file,omitempty"`
		PublicCert          bool   `toml:"public-cert,omitempty"`
		MirrorUse           bool   `toml:"mirror-use,omitempty"`
		CertFile            struct {
			SslCert    string `toml:"ssl-cert,omitempty"`
			SslCertKey string `toml:"ssl-cert-key,omitempty"`
		} `toml:"cert-file,omitempty"`
	} `toml:"private-registry,omitempty"`

	PrepareAirgap struct {
		K8sVersion      string `toml:"k8s-version,omitempty"`
		RegistryVersion string `toml:"registry-version,omitempty"`
		RegistryIP      string `toml:"registry-ip,omitempty"`
	} `toml:"prepare-airgap,omitempty"`

	SupportVersion struct {
		PackageVersion   PackageVersion
		ImageVersion     ImageVersion
		HelmChartVersion HelmChartVersion
	}

	ListVersion struct {
		ListPackageVersion   ListPackageVersion
		ListImageVersion     ListImageVersion
		ListHelmChartVersion ListHelmChartVersion
	}
}

type StrNode struct {
	Name      []string
	IP        []string `toml:"ip"`
	PrivateIP []string `toml:"private-ip"`
}
