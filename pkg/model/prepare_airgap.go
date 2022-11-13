package model

// validate: The first argument is a variable of SupportVersion,
// and the second argument is a variable of the list of supported versions.
type PackageVersion struct {
	Containerd    string `validate:"containerd,SupportContainerdVersion"`
	DockerCompose string `validate:"docker-compose,SupportDockerComposeVersion"`
}

type ImageVersion struct {
	Calico        string `validate:"calico,SupportCalicoVersion"`
	Coredns       string `validate:"coredns,SupportCorednsVersion"`
	MetricsServer string `validate:"metrics-server,SupportMetricsServerVersion"`
	Pause         string `validate:"pause,SupportPauseVersion"`
}
