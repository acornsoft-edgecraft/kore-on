package templates

const ClusterUpdateText = `
{{- $Master := .KoreOnTemp.NodePool.Master}}
{{- $Node := .KoreOnTemp.NodePool.Node}}
## Inventory for {{.Command}} task.
===========================================================================
Node Name                      IP Address              Private IP Adderss
===========================================================================
Current K8s Cluster
-------------------
{{-  range $index, $data := $Master.IP}}
master-{{$index}}                       {{$data}}                    {{if ne (len $Master.PrivateIP) 0}}{{index $Master.PrivateIP $index}}{{end -}}
{{  end}}
{{  range $index, $data := $Node.IP }}
node-{{$index}}                         {{$data}}                    {{if ne (len $Node.PrivateIP) 0}}{{index $Node.PrivateIP $index}}{{end -}} 
{{  end}}


------------

===========================================================================
Is this ok [y/n]: `
