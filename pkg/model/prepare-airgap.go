package model

// validate: The first argument is a variable of SupportVersion,
// and the second argument is a variable of the list of supported versions.
type PackageVersion struct {
	Containerd    string `validate:"containerd,SupportContainerdVersion"`
	DockerCompose string `validate:"docker-compose,SupportDockerComposeVersion"`
	Crictl        string `validate:"crictl,SupportCrictlVersion"`
	Etcd          string `validate:"etcd,SupportEtcdVersion"`
	Helm          string `validate:"helm,SupportHelmVersion"`
	CalicoCtl     string `validate:"calicoctl,SupportCalicoCtlVersion"`
	ClusterCtl    string `validate:"clusterctl,SupportClusterCtlVersion"`
}

type ImageVersion struct {
	Calico        string `validate:"calico,SupportCalicoVersion"`
	Coredns       string `validate:"coredns,SupportCorednsVersion"`
	MetricsServer string `validate:"metrics-server,SupportMetricsServerVersion"`
	Pause         string `validate:"pause,SupportPauseVersion"`
	DnsUtils      string `validate:"dns-utils,SupportDnsUtilsVersion"`
}

type HelmChartVersion struct {
	CsiDriverNfs string `validate:"csi-driver-nfs,ChartCsiDriverNfsVersion"`
	Koreboard    string `validate:"koreboard,ChartKoreboardVersion"`
}

// List Versions
type ListPackageVersion struct {
	Containerd    map[string][]string `validate:"containerd,SupportContainerdVersion"`
	DockerCompose map[string][]string `validate:"docker-compose,SupportDockerComposeVersion"`
	Crictl        map[string][]string `validate:"crictl,SupportCrictlVersion"`
	Etcd          map[string][]string `validate:"etcd,SupportEtcdVersion"`
	Helm          map[string][]string `validate:"helm,SupportHelmVersion"`
	CalicoCtl     map[string][]string `validate:"calicoctl,SupportCalicoCtlVersion"`
	ClusterCtl    map[string][]string `validate:"clusterctl,SupportClusterCtlVersion"`
}

type ListImageVersion struct {
	Calico        map[string][]string `validate:"calico,SupportCalicoVersion"`
	Coredns       map[string][]string `validate:"coredns,SupportCorednsVersion"`
	MetricsServer map[string][]string `validate:"metrics-server,SupportMetricsServerVersion"`
	Pause         map[string][]string `validate:"pause,SupportPauseVersion"`
	DnsUtils      map[string][]string `validate:"dns-utils,SupportDnsUtilsVersion"`
}

type ListHelmChartVersion struct {
	CsiDriverNfs map[string][]string `validate:"csi-driver-nfs,ChartCsiDriverNfsVersion"`
	Koreboard    map[string][]string `validate:"koreboard,ChartKoreboardVersion"`
}
