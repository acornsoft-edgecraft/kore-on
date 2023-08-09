package templates

const DestroyAllText = `
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
{{  end}}
===========================================================================
Is this ok [y/n]: `
