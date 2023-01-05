package config

const AddonTemplate = `
##################################################################################
## You can check supported applications with the 'list' command  	        			##
## and add applications to install. It can be used as below.				          	##
##                                                                              ##
## ※ If both "values" and "value_file" exist, "values" is used.				        	##
## -- Sample --													                                      	##
## [apps.application-name]													                          	##
## install = true																                                ##
## chart_ref_name = "xxx"												                            		##
## chart_ref = "https://helm-chart-address or helm-package-address(or path)"	  ##
## values="""																	                                  ##
## helm-chart-values															                              ##
## """																	                                    		##
## value_file = "helm-chart-values file path"								                  	##
##################################################################################

[addon]
## Required
## - k8s-master-ip: K8s control plane node ip address. (Deployment runs on this node.)
##					If you want to deploy locally, you must use the --kubeconfig option.
## -
## Optional
## - ssh-port: K8s Controlplane Node ssh port (default: 22)
## - addon-data-dir: addon data(helm vales, k8s deployment yaml) dir (default: "/data/addon") 
## -
#k8s-master-ip = "x.x.x.x"
#ssh-port = 22
#addon-data-dir = "/data/addon"
#closed-network = true

[apps.csi-driver-nfs]
## Required
## - install: Choose to proceed with installation.
## - chart_ref_name: helm chart repo name.
## - chart_ref: helm chart repository url.
## - chart_name: deployment chart name.
## - values_file: chart values file path (If both "values" and "value_file" exist, "values" is used.	)
## - values: chart values (If both "values" and "value_file" exist, "values" is used.	)
## -
## Optional
## - chart_version: deploy chart version (default: "latest") 
## -
#install = true
#chart_ref_name = "helm-charts"
#chart_ref = "https://192.168.77.119/chartrepo/helm-charts"
#chart_name = "csi-driver-nfs"
#chart_version = "<chart version>"
#values_file = "./csi-driver-nfs-values.yaml"
#values = """
storageClass:
  create: true
  parameters:
    mountOptions:
    - nfsvers=4.1
    server: 192.168.77.119
    share: /data/storage
"""

[apps.koreboard]
## Required
## - install: Choose to proceed with installation.
## - chart_ref_name: helm chart repo name.
## - chart_ref: helm chart repository url.
## - chart_name: deployment chart name.
## - values_file: chart values file path (If both "values" and "value_file" exist, "values" is used.	)
## - values: chart values (If both "values" and "value_file" exist, "values" is used.	)
## -
## Optional
## - chart_version: deploy chart version (default: "latest") 
## -
#install = true
#chart_ref_name = "<chart repo name>"
#chart_ref = "<chart repo url>"
#chart_name = "<chart name>"
#chart_version = "<chart version>"
#values_file = "<values.yaml to path>"
#values = """
"""
`
