package templates

const DestroyStorageText = `
{{- $Master := .KoreOnTemp.NodePool.Master}}
{{- $Node := .KoreOnTemp.NodePool.Node}}
{{- $PrepareAirgap := .KoreOnTemp.PrepareAirgap}}
{{- $PrivateRegistry := .KoreOnTemp.PrivateRegistry}}
{{- $SharedStorage := .KoreOnTemp.SharedStorage }}
## Inventory for {{.Command}} task.
===========================================================================
Node Name                      IP Address              Private IP Adderss
===========================================================================
{{-  if eq true $PrivateRegistry.Install -}}
{{    if eq true $SharedStorage.Install -}}
{{      if eq $PrivateRegistry.RegistryIP $SharedStorage.StorageIP}}
node-regi-storage            {{$PrivateRegistry.RegistryIP}}
{{      else}}
node-storage                   {{$SharedStorage.StorageIP}}                    {{if ne "" $SharedStorage.PrivateIP}}{{$SharedStorage.PrivateIP}}{{end -}}
{{      end -}}
{{    end -}}
{{  else if eq true $SharedStorage.Install}}
node-storage                   {{$SharedStorage.StorageIP}}                    {{if ne "" $SharedStorage.PrivateIP}}{{$SharedStorage.PrivateIP}}{{end -}}
{{  end }}
===========================================================================
Is this ok [y/n]: `
