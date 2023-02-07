package templates

const ClusterUpdateText = `
{{- $Master := .Master}}
{{- $Node := .Node}}
{{- $UpdateNode := .UpdateNode}}
{{- $text := "DELETE"}}
{{- $command := .Command}}
{{- $name_len := 0}}
{{- range $index, $data := $Master }}
{{- if lt $name_len (len $data.Name) }}
{{-   $name_len =  len $data.Name }}
{{- end }}
{{- end }}
{{- range $index, $data := $Node }}
{{-   $name_len =  len $data.Name }}
{{- end }}

Cluster Nodes
--------------
{{"NAME"|printf "%-20s"}}{{"STATUS"|printf "%-11s"}}{{"ROLES"|printf "%-23s"}}{{"AGE"|printf "%-8s"}}{{"VERSION"|printf "%-10s"}}{{"INTERNAL-IP"|printf "%-13s"}}{{"EXTERNAL-IP"|printf "%-13s"}}{{"OS-IMAGE"|printf "%-16s"}}{{"KERNEL-VERSION"|printf "%-28s"}}{{"CONTAINER-RUNTIME"|printf "%-20s"}}
{{- range $index, $data := $Master }}
{{$data.Name|printf "%-20s"}}{{$data.Status|printf "%-11s"}}{{$data.Role|printf "%-23s"}}{{$data.Age|printf "%-8s"}}{{$data.Version|printf "%-10s"}}{{$data.InternalIP|printf "%-13s"}}{{$data.ExternalIP|printf "%-13s"}}{{$data.OSImage|printf "%-16s"}}{{$data.KernelVersion|printf "%-28s"}}{{$data.ContainerRuntime|printf "%-20s"}}
{{- end}}
{{- range $index, $data := $Node }}
{{$data.Name|printf "%-20s"}}{{$data.Status|printf "%-11s"}}{{$data.Role|printf "%-23s"}}{{$data.Age|printf "%-8s"}}{{$data.Version|printf "%-10s"}}{{$data.InternalIP|printf "%-13s"}}{{$data.ExternalIP|printf "%-13s"}}{{$data.OSImage|printf "%-16s"}}{{$data.KernelVersion|printf "%-28s"}}{{$data.ContainerRuntime|printf "%-20s"}}
{{- end}}


Update Nodes ({{ $command }})
----------------------
{{ printf "%.*s" 20 "======================================================================================================================================================" }}
===========================================================================
Node Name                      IP Address              Private IP Adderss
===========================================================================
{{- range $index, $data := $UpdateNode.IP }}
{{-   if eq $command $text }}
{{(index $UpdateNode.Name $index)}}      { $name_len }                 {{$data}}             {{(index $UpdateNode.PrivateIP $index)}}
{{-   else }}
node-{{$index}}                       {{$data}}             {{(index $UpdateNode.PrivateIP $index)}}
{{-   end }}
{{- end}}
===========================================================================
Is this ok [y/n]: `
