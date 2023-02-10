package templates

const ClusterGetKubeconfigText = `
{{- $Master := .KoreOnTemp.NodePool.Master }}
{{- $private_ip := "" }}
{{- $check_ip := false }}
{{- if eq (len $Master.IP) (len $Master.PrivateIP) }}
{{- $check_ip = true }}
{{- end }}

## Inventory for {{.Command}} task.
{{ printf "%.*s" 64 "======================================================================================================================================================" }}
{{"Node Name"|printf "%-*s" 20}}{{"IP"|printf "%-*s" 22}}{{"Private IP"|printf "%-*s" 22}}
{{ printf "%.*s" 64 "======================================================================================================================================================" }}
{{-  range $index, $data := $Master.IP}}
{{- if $check_ip }}
{{- $private_ip = (index $Master.PrivateIP $index)}}
{{- end }}
node-{{$index|printf "%-*v" 15}}{{$data|printf "%-*s" 22}}{{$private_ip|printf "%-*s" 22}}
{{- break }}
{{-  end}}
{{ printf "%.*s" 64 "======================================================================================================================================================" }}
Is this ok [y/n]: `
