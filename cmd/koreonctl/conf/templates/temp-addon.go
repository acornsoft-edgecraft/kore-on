package templates

const AddonText = `
{{- $Master := .AddonTemp.Addon.K8sMasterIP}}
## Inventory for {{.Command}} task.
===========================================================================
Node Name                      IP Address              Private IP Adderss
===========================================================================
{{ if ne "" $Master }}
k8s-master-1                 {{$Master}}                    
{{ end }}
===========================================================================
Is this ok [y/N]: `
