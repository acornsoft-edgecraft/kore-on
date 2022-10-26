package model

type KoreOnToml struct {
	KoreOn struct {
		ClusterName      string `toml:"cluster-name,omitempty"`
		ClusterID        string `toml:"cluster-id,omitempty"`
		InstallDir       string `toml:"install-dir,omitempty"`
		CertValidityDays int    `toml:"cert-validity-days,omitempty"`

		//#Airgap
		ClosedNetwork              bool   `toml:"closed-network,omitempty"`
		LocalRepository            string `toml:"local-repository,omitempty"`
		LocalRepositoryArchiveFile string `toml:"local-repository-archive-file"`
		DebugMode                  bool   `toml:"debug-mode,omitempty"`
	} `toml:"koreon,omitempty"`

	Kubernetes struct {
		Version          string   `toml:"version,omitempty"`
		ContainerRuntime string   `toml:"container-runtime"`
		KubeProxyMode    string   `toml:"kube-proxy-mode"`
		CalicoVersion    string   `toml:"calico-version"`
		ServiceCidr      string   `toml:"service-cidr,omitempty"`
		PodCidr          string   `toml:"pod-cidr,omitempty"`
		NodePortRange    string   `toml:"node-port-range,omitempty"`
		AuditLogEnable   bool     `toml:"audit-log-enable"`
		ApiSans          []string `toml:"api-sans,omitempty"`

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

		Security struct {
			SSHUserID      string `toml:"ssh-user-id,omitempty"`
			SSHPort        int    `toml:"ssh-port,omitempty"`
			PrivateKeyPath string `toml:"private-key-path,omitempty"`
		} `toml:"security,omitempty"`

		Master struct {
			Name           string   `toml:"name,omitempty"`
			IP             []string `toml:"ip,omitempty"`
			PrivateIP      []string `toml:"private-ip,omitempty"`
			LbIP           string   `toml:"lb-ip,omitempty"`
			LbPort         int      `toml:"lb-port,omitempty"`
			Isolated       bool     `toml:"isolated"`
			HaproxyInstall bool     `toml:"haproxy-install"`
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
		Install             bool   `toml:"install"`
		RegistryVersion     string `toml:"registry-version,omitempty"`
		RegistryIP          string `toml:"registry-ip,omitempty"`
		RegistryDomain      string `toml:"registry-domain,omitempty"`
		PrivateIP           string `toml:"private-ip,omitempty"`
		DataDir             string `toml:"data-dir,omitempty"`
		RegistryArchiveFile string `toml:"registry-archive-file"`
		PublicCert          bool   `toml:"public-cert"`
		CertFile            struct {
			SslCertificate    string `toml:"ssl-certificate,omitempty"`
			SslCertificateKey string `toml:"ssl-certificate-key,omitempty"`
		} `toml:"cert-file,omitempty"`
	} `toml:"private-registry,omitempty"`
}

type StrNode struct {
	IP        []string `toml:"ip,omitempty"`
	PrivateIP []string `toml:"private-ip,omitempty"`
}
