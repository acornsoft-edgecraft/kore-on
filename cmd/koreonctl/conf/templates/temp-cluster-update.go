package templates

const ClusterUpdateText = `
{{- $Master := .Master}}
{{- $Node := .Node}}
{{- $UpdateNode := .UpdateNode}}
{{- $text := "DELETE"}}
{{- $command := .Command}}
{{- $master_len := .Master | maxLength}}
{{- $node_len := .Node | maxLength}}
{{- $update_node_len := .UpdateNode | maxLength}}
{{- $cluster_len := clusterLength $master_len $node_len}}
{{- $total := total $update_node_len.Name  $update_node_len.IP $update_node_len.PrivateIP}}

Cluster Nodes
--------------
{{"NAME"|printf "%-*s" $cluster_len.Name}}{{"STATUS"|printf "%-*s" $cluster_len.Status}}{{"ROLES"|printf "%-*s" $cluster_len.Role}}{{"AGE"|printf "%-*s" $cluster_len.Age}}{{"VERSION"|printf "%-*s" $cluster_len.Version}}{{"INTERNAL-IP"|printf "%-*s" $cluster_len.InternalIP}}{{"EXTERNAL-IP"|printf "%-*s" $cluster_len.ExternalIP}}{{"OS-IMAGE"|printf "%-*s" $cluster_len.OSImage}}{{"KERNEL-VERSION"|printf "%-*s" $cluster_len.KernelVersion}}{{"CONTAINER-RUNTIME"|printf "%-*s" $cluster_len.ContainerRuntime}}
{{- range $index, $data := $Master }}
{{$data.Name|printf "%-*s" $cluster_len.Name}}{{$data.Status|printf "%-*s" $cluster_len.Status}}{{$data.Role|printf "%-*s" $cluster_len.Role}}{{$data.Age|printf "%-*s" $cluster_len.Age}}{{$data.Version|printf "%-*s" $cluster_len.Version}}{{$data.InternalIP|printf "%-*s" $cluster_len.InternalIP}}{{$data.ExternalIP|printf "%-*s" $cluster_len.ExternalIP}}{{$data.OSImage|printf "%-*s" $cluster_len.OSImage}}{{$data.KernelVersion|printf "%-*s" $cluster_len.KernelVersion}}{{$data.ContainerRuntime|printf "%-*s" $cluster_len.ContainerRuntime}}
{{- end}}
{{- range $index, $data := $Node }}
{{$data.Name|printf "%-*s" $cluster_len.Name}}{{$data.Status|printf "%-*s" $cluster_len.Status}}{{$data.Role|printf "%-*s" $cluster_len.Role}}{{$data.Age|printf "%-*s" $cluster_len.Age}}{{$data.Version|printf "%-*s" $cluster_len.Version}}{{$data.InternalIP|printf "%-*s" $cluster_len.InternalIP}}{{$data.ExternalIP|printf "%-*s" $cluster_len.ExternalIP}}{{$data.OSImage|printf "%-*s" $cluster_len.OSImage}}{{$data.KernelVersion|printf "%-*s" $cluster_len.KernelVersion}}{{$data.ContainerRuntime|printf "%-*s" $cluster_len.ContainerRuntime}}
{{- end}}


Update Nodes ({{ $command }})
----------------------
{{ printf "%.*s" $total "======================================================================================================================================================" }}
{{- if eq $command $text }}
{{"Node Name"|printf "%-*s" $update_node_len.Name}}{{"IP"|printf "%-*s" $update_node_len.IP}}{{"Private IP"|printf "%-*s" $update_node_len.PrivateIP}}
{{- else }}
{{"Node Name"}}{{printf "%-*s" 6 ""}}{{"IP"|printf "%-*s" $update_node_len.IP}}{{"Private IP"|printf "%-*s" $update_node_len.PrivateIP}}
{{- end}}
{{ printf "%.*s" $total "======================================================================================================================================================" }}
{{- range $index, $data := $UpdateNode.IP }}
{{-   if eq $command $text }}
{{(index $UpdateNode.Name $index)|printf "%-*s" $update_node_len.Name}}{{$data|printf "%-*s" $update_node_len.IP}}{{(index $UpdateNode.PrivateIP $index)|printf "%-*s" $update_node_len.PrivateIP}}
{{-   else }}
node-{{$index|printf "%-*v" 10 }}{{$data|printf "%-*s" $update_node_len.IP}}{{(index $UpdateNode.PrivateIP $index)|printf "%-*s" $update_node_len.PrivateIP}}
{{-   end }}
{{- end}}
{{ printf "%.*s" $total "======================================================================================================================================================" }}
Is this ok [y/n]: `
