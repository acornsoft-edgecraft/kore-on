package templates

const DestroyClusterText = `
{{- $Master := .KoreOnTemp.NodePool.Master}}
{{- $Node := .KoreOnTemp.NodePool.Node}}
{{- $PrepareAirgap := .KoreOnTemp.PrepareAirgap}}
{{- $PrivateRegistry := .KoreOnTemp.PrivateRegistry}}
{{- $SharedStorage := .KoreOnTemp.SharedStorage }}
## Inventory for {{.Command}} task.
===========================================================================
Node Name                      IP Address              Private IP Adderss
===========================================================================
{{-  range $index, $data := $Master.IP}}
master-{{$index}}                       {{$data}}                    {{if ne (len $Master.PrivateIP) 0}}{{index $Master.PrivateIP $index}}{{end -}}
{{  end}}
{{  range $index, $data := $Node.IP }}
node-{{$index}}                         {{$data}}                    {{if ne (len $Node.PrivateIP) 0}}{{index $Node.PrivateIP $index}}{{end -}} 
{{  end}}
===========================================================================
Is this ok [y/n]: `
