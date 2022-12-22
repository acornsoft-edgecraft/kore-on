package model

type AddonToml struct {
	Addon struct {
		K8sMasterIP    string `toml:"k8s-master-ip,omitempty"`
		SSHPort        int    `toml:"ssh-port,omitempty"`
		AddonDataDir   string `toml:"addon-data-dir,omitempty"`
		HelmVersion    string
		HelmInstall    bool
		HelmBinaryFile string
	} `toml:"addon,omitempty"`

	Apps struct {
		CsiDriverNfs  AppsCsiDriverNfs  `toml:"csi-driver-nfs,omitempty"`
		BitnamiNginx  AppsBitnamiNginx  `toml:"bitnami-nginx,omitempty"`
		Elasticsearch AppsElasticsearch `toml:"elasticsearch,omitempty"`
		FluentBit     AppsFluentBit     `toml:"fluent-bit,omitempty"`
		Koreboard     AppsKoreboard     `toml:"koreboard,omitempty"`
	} `toml:"apps,omitempty"`
}

type AppsCsiDriverNfs struct {
	Install         bool   `toml:"install,omitempty"`
	ChartRefName    string `toml:"chart_ref_name,omitempty"`
	ChartRef        string `toml:"chart_ref,omitempty"`
	ChartRefID      string
	ChartRefPW      string
	StorageIP       string `toml:"storage_ip,omitempty"`
	SharedVolumeDir string `toml:"shared_volume_dir,omitempty"`
	NfsVersion      string `toml:"nfs_version,omitempty"`
}

type AppsBitnamiNginx struct {
	Install      bool   `toml:"install,omitempty"`
	ChartRefName string `toml:"chart_ref_name,omitempty"`
	ChartRef     string `toml:"chart_ref,omitempty"`
	Port         string `toml:"port,omitempty"`
}

type AppsElasticsearch struct {
	Install      bool   `toml:"install,omitempty"`
	ChartRefName string `toml:"chart_ref_name,omitempty"`
	ChartRef     string `toml:"chart_ref,omitempty"`
	Values       string `toml:"values,omitempty"`
}
type AppsFluentBit struct {
	Install      bool   `toml:"install,omitempty"`
	ChartRefName string `toml:"chart_ref_name,omitempty"`
	ChartRef     string `toml:"chart_ref,omitempty"`
	Values       string `toml:"values,omitempty"`
}

type AppsKoreboard struct {
	Install      bool   `toml:"install,omitempty"`
	ChartRefName string `toml:"chart_ref_name,omitempty"`
	ChartRef     string `toml:"chart_ref,omitempty"`
	Values       string `toml:"values,omitempty"`
}
