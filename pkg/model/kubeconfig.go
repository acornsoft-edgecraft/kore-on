package model

type KubeConfig struct {
	APIVersion     string        `yaml:"apiVersion"`
	Clusters       []KubeCluster `yaml:"clusters"`
	Contexts       []KubeContext `yaml:"contexts"`
	CurrentContext string        `yaml:"current-context"`
	Kind           string        `yaml:"kind"`
	Preferences    struct {
	} `yaml:"preferences"`
	Users []KubeUser `yaml:"users"`
}

type KubeCluster struct {
	Cluster struct {
		CertificateAuthorityData string `yaml:"certificate-authority-data"`
		Server                   string `yaml:"server"`
	} `yaml:"cluster"`
	Name string `yaml:"name"`
}

type KubeContext struct {
	Context struct {
		Cluster string `yaml:"cluster"`
		User    string `yaml:"user"`
	} `yaml:"context"`
	Name string `yaml:"name"`
}

type KubeUser struct {
	Name string `yaml:"name"`
	User struct {
		ClientCertificateData string `yaml:"client-certificate-data"`
		ClientKeyData         string `yaml:"client-key-data"`
	} `yaml:"user"`
}
