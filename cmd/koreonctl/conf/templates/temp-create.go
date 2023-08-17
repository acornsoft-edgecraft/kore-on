package templates

const CreateText = `
{{- $Master := .KoreOnTemp.NodePool.Master}}
{{- $Node := .KoreOnTemp.NodePool.Node}}
{{- $PrepareAirgap := .KoreOnTemp.PrepareAirgap}}
{{- $PrivateRegistry := .KoreOnTemp.PrivateRegistry}}
{{- $SharedStorage := .KoreOnTemp.SharedStorage }}
## Inventory for {{ .Command | ToUpper }} task.
## Setting up the cluster in an {{ if .KoreOnTemp.KoreOn.ClosedNetwork }}air-gapped{{else}}online{{ end }} environment
===========================================================================
Node Name                      IP Address              Private IP Adderss
===========================================================================
{{-  range $index, $data := $Master.IP}}
master-{{ $index |printf "%-*v" 22 }}{{ $data | printf "%-*s" 28 }}{{if ne (len $Master.PrivateIP) 0}}{{index $Master.PrivateIP $index}}{{end -}}
{{  end}}
{{  range $index, $data := $Node.IP }}
node-{{ $index |printf "%-*v" 24 }}{{ $data | printf "%-*s" 28 }}{{if ne (len $Node.PrivateIP) 0}}{{index $Node.PrivateIP $index}}{{end -}} 
{{  end}}
{{  if eq true $PrivateRegistry.Install -}}
{{    if eq true $SharedStorage.Install -}}
{{      if eq $PrivateRegistry.RegistryIP $SharedStorage.StorageIP}}
{{ "node-regi-storage" | printf "%-*s" 29 }}{{ $PrivateRegistry.RegistryIP | printf "%-*s" 28 }}{{if ne "" $PrivateRegistry.PrivateIP}}{{$PrivateRegistry.PrivateIP}}{{end -}}
{{      else}}
{{ "node-regi" | printf "%-*s" 29 }}{{ $PrivateRegistry.RegistryIP | printf "%-*s" 28 }}{{if ne "" $PrivateRegistry.PrivateIP}}{{$PrivateRegistry.PrivateIP}}{{end -}}
{{ "node-storage" | printf "%-*s" 29 }}{{ $SharedStorage.StorageIP | printf "%-*s" 28 }}{{if ne "" $SharedStorage.PrivateIP}}{{$SharedStorage.PrivateIP}}{{end -}}
{{      end -}}
{{    else}}
{{ "node-regi" | printf "%-*s" 29 }}{{ $PrivateRegistry.RegistryIP | printf "%-*s" 28 }}{{if ne "" $PrivateRegistry.PrivateIP}}{{$PrivateRegistry.PrivateIP}}{{end -}}
{{    end -}}
{{  else if eq true $SharedStorage.Install}}
{{ "node-storage" | printf "%-*s" 29 }}{{ $SharedStorage.StorageIP | printf "%-*s" 28 }}{{if ne "" $SharedStorage.PrivateIP}}{{$SharedStorage.PrivateIP}}{{end -}}
{{ else if and (ne "" $PrivateRegistry.RegistryIP)  (eq true $PrivateRegistry.PublicCert) }}
## Private repositories are not installed 
{{ "node-regi" | printf "%-*s" 29 }}{{ $PrivateRegistry.RegistryIP | printf "%-*s" 28 }}{{if ne "" $PrivateRegistry.PrivateIP}}{{$PrivateRegistry.PrivateIP}}{{end -}}
{{ else if and (ne "" $PrivateRegistry.RegistryDomain)  (ne true $PrivateRegistry.PublicCert) }}
## Private repositories are not installed. used domain name
{{ "node-regi" | printf "%-*s" 29 }}{{ $PrivateRegistry.RegistryDomain }}
{{  end}}
===========================================================================
Is this ok [y/n]: `
