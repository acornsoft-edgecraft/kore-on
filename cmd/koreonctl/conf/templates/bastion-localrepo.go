package templates

const BastionLocalRepoText = `
[bastion-local-to-file]
name=bastion-local-repo
baseurl=file://{{.}}
gpgcheck=0
enabled=1`
