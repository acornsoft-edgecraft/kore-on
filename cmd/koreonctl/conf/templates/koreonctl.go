package templates

const KoreonctlText = `
{{- $Master := .KoreOnTemp.NodePool.Master}}
{{- $Node := .KoreOnTemp.NodePool.Node}}
{{- $PrepareAirgap := .KoreOnTemp.PrepareAirgap}}
{{- $PrivateRegistry := .KoreOnTemp.PrivateRegistry}}
{{- $SharedStorage := .KoreOnTemp.SharedStorage }}
## Inventory for {{.Command}} task.
===========================================================================
Node Name                      IP Address              Private IP Adderss
===========================================================================
{{- if ne "prepare-airgap" .Command -}}
{{  range $index, $data := $Master.IP}}
master-{{$index}}                       {{$data}}                    {{if ne (len $Master.PrivateIP) 0}}{{index $Master.PrivateIP $index}}{{end -}}
{{  end}}
{{  range $index, $data := $Node.IP }}
node-{{$index}}                         {{$data}}                    {{if ne (len $Node.PrivateIP) 0}}{{index $Node.PrivateIP $index}}{{end -}} 
{{  end}}
{{  if eq true $PrivateRegistry.Install -}}
{{    if eq true $SharedStorage.Install -}}
{{      if eq $PrivateRegistry.RegistryIP $SharedStorage.StorageIP}}
node-regi-storage            {{$PrivateRegistry.RegistryIP}}
{{      else}}
node-regi                      {{$PrivateRegistry.RegistryIP}}                    {{if ne "" $PrivateRegistry.PrivateIP}}{{$PrivateRegistry.PrivateIP}}{{end}}
node-storage                   {{$SharedStorage.StorageIP}}                    {{if ne "" $SharedStorage.PrivateIP}}{{$SharedStorage.PrivateIP}}{{end -}}
{{      end}}
{{    else}}
node-regi                               {{$PrivateRegistry.RegistryIP}}           {{if ne "" $PrivateRegistry.PrivateIP}}{{$PrivateRegistry.PrivateIP}}{{end -}}
{{    end -}}
{{  else if eq true $SharedStorage.Install}}
node-storage                   {{$SharedStorage.StorageIP}}                    {{if ne "" $SharedStorage.PrivateIP}}{{$SharedStorage.PrivateIP}}{{end -}}
{{  end -}}
{{else}}
prepare-airgap-node                      {{$PrepareAirgap.RegistryIP}}                   
{{end -}}
===========================================================================
Is this ok [y/N]: `
