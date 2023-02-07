package kubemethod

import (
	"context"
	"fmt"
	"kore-on/pkg/model/k8s"
	"strconv"
	"strings"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// GetNodeList - 해당 클러스터의 노드 리스트 반환
func GetNodeList(client *kubernetes.Clientset) ([]k8s.Node, error) {
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return k8s.ConvertToNodeList(nodes)
}

func CreateK8sClient(config *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func GetVersion(client *kubernetes.Clientset) (int, int, error) {

	serverVersion, err := client.ServerVersion()
	if err != nil {
		fmt.Printf("[ERROR] fail to get nodes: %s", err.Error())
		return -1, -1, err
	}

	spVersion := strings.Split(serverVersion.GitVersion, ".")

	minor, err := strconv.Atoi(spVersion[1])
	if err != nil {
		fmt.Printf("[ERROR] fail to get nodes: %s", err.Error())
		return -1, -1, err
	}

	patch, err := strconv.Atoi(spVersion[2])
	if err != nil {
		fmt.Printf("[ERROR] fail to get nodes: %s", err.Error())
		return -1, -1, err
	}

	return minor, patch, nil

}
