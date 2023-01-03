package conf

const docker_repo_rhel = `
[Docker-CE-Stable]
async = 1
baseurl = https://download.docker.com/linux/centos/$releasever/$basearch/stable
enabled = 1
gpgcheck = 1
gpgkey = https://download.docker.com/linux/centos/gpg
name = Docker-ce repo
repo_gpgcheck = 0
`
