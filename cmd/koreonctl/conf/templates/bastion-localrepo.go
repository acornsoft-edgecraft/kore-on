package templates

const BastionLocalRepoText = `
[bastionlocal]
name=bastion local repo
baseurl=file://{{.}}
gpgcheck=0
enabled=0`
