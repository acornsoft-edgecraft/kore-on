package templates

const ClusterGetKubeconfigText = `
{{- $Master := .KoreOnTemp.NodePool.Master}}
## Inventory for {{.Command}} task.
===========================================================================
Node Name                      IP Address              Private IP Adderss
===========================================================================
{{-  range $index, $data := $Master.IP}}
master-{{$index}}                       {{$data}}                    {{if ne (len $Master.PrivateIP) 0}}{{index $Master.PrivateIP $index}}{{end -}}
{{  end}}
===========================================================================
Is this ok [y/n]: `
