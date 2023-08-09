package templates

const DestroyRegistryText = `
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
node-regi                      {{$PrivateRegistry.RegistryIP}}                    {{if ne "" $PrivateRegistry.PrivateIP}}{{$PrivateRegistry.PrivateIP}}{{end -}}
{{      end -}}
{{    else}}
node-regi                               {{$PrivateRegistry.RegistryIP}}           {{if ne "" $PrivateRegistry.PrivateIP}}{{$PrivateRegistry.PrivateIP}}{{end -}}
{{    end -}}
{{  end }}
===========================================================================
Is this ok [y/n]: `
