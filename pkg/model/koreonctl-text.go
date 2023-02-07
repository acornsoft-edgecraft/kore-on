package model

import "kore-on/pkg/model/k8s"

type KoreonctlText struct {
	Command     string
	KoreOnTemp  KoreOnToml
	Master      []k8s.Node
	Node        []k8s.Node
	UpdateNode  updateNode
	PrintFormat printFormat
}

type updateNode struct {
	IP        []string
	PrivateIP []string
	Name      []string
}

type printFormat struct {
	Name             int
	Status           int
	Roles            int
	Age              int
	Version          int
	InternalIP       int
	ExternalIP       int
	OsImage          int
	KernelVersion    int
	ContainerRuntime int
}
