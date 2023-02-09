package templates

const ClusterGetKubeconfigText = `
{{- $Master := .KoreOnTemp.NodePool.Master}}
{{- $master_len := $Master | maxLength}}
{{- $total := total $master_len.Name  $master_len.IP $master_len.PrivateIP}}

## Inventory for {{.Command}} task.
{{ printf "%.*s" $total "======================================================================================================================================================" }}
{{"Node Name"|printf "%-*s" $master_len.Name}}{{"IP Address"|printf "%-*s" $master_len.IP}}{{"Private IP Adderss"|printf "%-*s" $master_len.PrivateIP}}
{{ printf "%.*s" $total "======================================================================================================================================================" }}
{{-  range $index, $data := $Master}}
node-{{$index|printf "%-*s" 10}}{{$data.IP|printf "%-*s" $master_len.IP}}{{$data.PrivateIP|printf "%-*s" $master_len.PrivateIP}}
{{ break }}
{{  end}}
{{ printf "%.*s" $total "======================================================================================================================================================" }}
Is this ok [y/n]: `
