package model

type Cluster struct {
	Controlplane []Controlplane
	WorkerNode   []WorkerNode
	AddNode      []WorkerNode
	DeleteNode   []WorkerNode
}

type Controlplane struct {
	NodeName       string
	NodeExternalIP string
	NodeInternalIP string
}

type WorkerNode struct {
	NodeName       string
	NodeExternalIP string
	NodeInternalIP string
}

type AddNode struct {
	NodeName       string
	NodeExternalIP string
	NodeInternalIP string
}

type DeleteNode struct {
	NodeName       string
	NodeExternalIP string
	NodeInternalIP string
}
