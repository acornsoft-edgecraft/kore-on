/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package k8s

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
)

// Node - Kubernetes Node 정보
type Node struct {
	Name             string `json:"name"`
	Status           string `json:"status"`
	Role             string `json:"role"`
	Age              string `json:"age"`
	Version          string `json:"version"`
	InternalIP       string `json:"internal_ip"`
	ExternalIP       string `json:"external_ip"`
	OSImage          string `json:"os_image"`
	KernelVersion    string `json:"kernel_version"`
	ContainerRuntime string `json:"container_image"`
	AnsibleSshHost   string `json:"ansible_ssh_host"`
}

// isReady - 노드 상태 반환
func isReady(node *v1.Node) string {
	var cond v1.NodeCondition
	for _, c := range node.Status.Conditions {
		if c.Type == v1.NodeReady {
			cond = c
			break
		}
	}

	if cond.Status == v1.ConditionTrue {
		return "Ready"
	} else {
		return "Not ready"
	}
}

// getRoles - 노드의 역할들 반환
func getRoles(node *v1.Node) string {
	var roles []string

	if _, ok := node.Labels["node-role.kubernetes.io/control-plane"]; ok {
		roles = append(roles, "control-plane")
	}
	if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
		roles = append(roles, "master")
	}

	if len(roles) == 0 {
		roles = append(roles, "<none>")
	}

	return strings.Join(roles, ",")
}

// getAge - 노드의 생성이후 시간 반환
func getAge(node *v1.Node) string {
	// 초 / 60 > 분
	// 초 / 60 / 60 > 시간
	// 초 / 60 / 60 / 24 > 일
	var buff bytes.Buffer
	diff := time.Since(node.CreationTimestamp.Time)
	days := diff / (24 * time.Hour)

	if days > 1 {
		buff.WriteString(fmt.Sprintf("%dd", days))
	} else {
		hours := diff % (24 * time.Hour)
		minutes := hours % time.Hour
		// seconds := math.Mod(minutes.Seconds(), 60)

		if hours/time.Hour > 0 {
			buff.WriteString(fmt.Sprintf("%dh", hours/time.Hour))
		}
		if minutes/time.Minute > 0 {
			buff.WriteString(fmt.Sprintf("%dm", minutes/time.Minute))
		}
		// if seconds > 0 {
		// 	buff.WriteString(fmt.Sprintf("%.1fs", seconds))
		// }
	}

	return buff.String()

	// diff := uint64(time.Since(node.CreationTimestamp.Time).Seconds())
	// days := diff / (60*60*24)
	// if days > 1 {
	// 	return strconv.FormatUint(days, 10) + "d"
	// }

	// hours := diff / (60*60)
	// if hours > 1 {
	// 	return strconv.FormatUint(diff, 10) + "h"
	// }
	// minutes := diff * 60
	// if minutes > 1 {
	// 	return strconv.FormatUint(minutes, 10) + "m"
	// }

	// return strconv.FormatUint(minutes*60, 10) + "s"
}

// getKubeletVersion - 노드의 쿠버네티스 버전 반환
func getKubeletVersion(node *v1.Node) string {
	return node.Status.NodeInfo.KubeletVersion
}

// getNodeAddresses - 노드의 Internal/External Address 반환
func getNodeAddresses(addresses []v1.NodeAddress) (string, string) {
	var internal, external string
	for _, addr := range addresses {
		if addr.Type == v1.NodeInternalIP {
			internal = addr.Address
		} else if addr.Type == v1.NodeExternalIP {
			external = addr.Address
		}
	}

	return internal, external
}

// ConvertToNodeList - Kubernetes NodeList를 화면에서 사용할 수 있는 정보로 전환
func ConvertToNodeList(nodeList *v1.NodeList) ([]Node, error) {
	var nodes []Node

	for _, item := range nodeList.Items {
		internalIp, externalIp := getNodeAddresses(item.Status.Addresses)
		node := Node{
			Name:             item.Name,
			Status:           isReady(&item),
			Role:             getRoles(&item),
			Age:              getAge(&item),
			Version:          getKubeletVersion(&item),
			InternalIP:       internalIp,
			ExternalIP:       externalIp,
			OSImage:          item.Status.NodeInfo.OSImage,
			KernelVersion:    item.Status.NodeInfo.KernelVersion,
			ContainerRuntime: item.Status.NodeInfo.ContainerRuntimeVersion,
			AnsibleSshHost:   item.Labels["koreon.acornsoft.io/ansible_ssh_host"],
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}
