package model

type AddonToml struct {
	Addon struct {
		K8sMasterIP    string `toml:"k8s-master-ip,omitempty"`
		SSHPort        int    `toml:"ssh-port,omitempty"`
		AddonDataDir   string `toml:"addon-data-dir,omitempty"`
		ClosedNetwork  bool   `toml:"closed-network,omitempty"`
		KubeConfig     string
		HelmVersion    string
		HelmInstall    bool
		HelmBinaryFile string
		WorkDir        string
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
	Install          bool   `toml:"install,omitempty"`
	ChartRefName     string `toml:"chart_ref_name,omitempty"`
	ChartRef         string `toml:"chart_ref,omitempty"`
	ChartName        string `toml:"chart_name,omitempty"`
	ReleaseNamespace string `toml:"chart_name,omitempty"`
	ChartVersion     string `toml:"chart_version,omitempty"`
	ChartRefID       string
	ChartRefPW       string
	Values           string `toml:"values,omitempty"`
	ValuesFile       string `toml:"values_file,omitempty"`
}

type AppsBitnamiNginx struct {
	Install          bool   `toml:"install,omitempty"`
	ChartRefName     string `toml:"chart_ref_name,omitempty"`
	ChartRef         string `toml:"chart_ref,omitempty"`
	ChartName        string `toml:"chart_name,omitempty"`
	ReleaseNamespace string `toml:"chart_name,omitempty"`
	ChartVersion     string `toml:"chart_version,omitempty"`
	Values           string `toml:"values,omitempty"`
	ValuesFile       string `toml:"values_file,omitempty"`
}

type AppsElasticsearch struct {
	Install          bool   `toml:"install,omitempty"`
	ChartRefName     string `toml:"chart_ref_name,omitempty"`
	ChartRef         string `toml:"chart_ref,omitempty"`
	ChartName        string `toml:"chart_name,omitempty"`
	ReleaseNamespace string `toml:"chart_name,omitempty"`
	ChartVersion     string `toml:"chart_version,omitempty"`
	Values           string `toml:"values,omitempty"`
	ValuesFile       string `toml:"values_file,omitempty"`
}
type AppsFluentBit struct {
	Install          bool   `toml:"install,omitempty"`
	ChartRefName     string `toml:"chart_ref_name,omitempty"`
	ChartRef         string `toml:"chart_ref,omitempty"`
	ChartName        string `toml:"chart_name,omitempty"`
	ReleaseNamespace string `toml:"chart_name,omitempty"`
	ChartVersion     string `toml:"chart_version,omitempty"`
	Values           string `toml:"values,omitempty"`
	ValuesFile       string `toml:"values_file,omitempty"`
}

type AppsKoreboard struct {
	Install          bool   `toml:"install,omitempty"`
	ChartRefName     string `toml:"chart_ref_name,omitempty"`
	ChartRef         string `toml:"chart_ref,omitempty"`
	ChartName        string `toml:"chart_name,omitempty"`
	ReleaseNamespace string `toml:"chart_name,omitempty"`
	Values           string `toml:"values,omitempty"`
	ValuesFile       string `toml:"values_file,omitempty"`
}
