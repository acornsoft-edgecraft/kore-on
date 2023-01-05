package templates

const AddonText = `
{{- $Master := .AddonTemp.Addon.K8sMasterIP }}
{{- $Apps := .AddonTemp.Apps }}
{{- $check := false }}
{{- $task := "Installation" }}
{{- if eq "delete" .Command }}
{{- $task = "DELETE" }}
{{- end}}
## Inventory for {{ $task }} task.
===========================================================================
Node Name                      IP Address              Private IP Adderss
===========================================================================
{{- if ne "" $Master }}
k8s-master-1                 {{$Master}}                    
{{ end -}}
===========================================================================

 {{ $task }} Application List
-------------------------------
{{- range $k, $v := $Apps -}}
{{-   range $i, $j := $v -}}
{{-     if eq "Install" $i -}}
{{-        if eq true $j}}
{{-       $check = $j -}}
{{-        else }}
{{-        $check = false -}}
{{-        end }}
{{-     end -}}
{{-   end }}
{{- if eq true $check }}
{{ $k }}
{{- end }}
{{- end }}


> Is this ok [y/n]: `
