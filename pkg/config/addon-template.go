package config

const AddonTemplate = `
[addon]
## Required
## - k8s-master-ip: K8s control plane node ip address. (Deployment runs on this node.)
##					If you want to deploy locally, you must use the --kubeconfig option.
## -
k8s-master-ip = ""

[apps.csi-driver-nfs]
## Required
## * If all fields are omitted, default values are distributed.
## - storage-ip: Storage node ip address.
## - volume-dir: Storage node data directory. (default: /data/storage)
## - nfs_version: Nfs-server version.
## -
storage_ip = ""
shared_volume_dir = ""
nfs_version = ""

[apps.bitnami-nginx]
## Required
## * If all fields are omitted, default values are distributed.
`
