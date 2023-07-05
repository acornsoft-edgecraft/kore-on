package templates

const PrepareAirgapText = `
{{- $Master := .KoreOnTemp.NodePool.Master}}
{{- $Node := .KoreOnTemp.NodePool.Node}}
{{- $PrepareAirgap := .KoreOnTemp.PrepareAirgap}}
{{- $PrivateRegistry := .KoreOnTemp.PrivateRegistry}}
{{- $SharedStorage := .KoreOnTemp.SharedStorage }}
## Inventory for {{.Command}} task.
===========================================================================
Node Name                      IP Address              Private IP Adderss
===========================================================================
{{ "prepare-airgap-node" | printf "%-*s" 29 }}{{ $PrepareAirgap.RegistryIP }}
===========================================================================
Is this ok [y/n]: `
