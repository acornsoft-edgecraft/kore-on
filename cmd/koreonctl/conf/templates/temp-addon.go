package templates

const AddonText = `
{{- $Master := .AddonTemp.Addon.K8sMasterIP }}
{{- $Apps := .AddonTemp.Apps }}
## Inventory for {{.Command}} task.
===========================================================================
Node Name                      IP Address              Private IP Adderss
===========================================================================
{{- if ne "" $Master }}
k8s-master-1                 {{$Master}}                    
{{ end -}}
===========================================================================

 Installation Application List
-------------------------------
{{ range $k, $v := $Apps }}
{{- $k }}
{{ end }}


> Is this ok [y/N]: `
