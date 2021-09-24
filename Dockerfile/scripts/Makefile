encrypt:
	echo "Encrypt files ..."
	ansible-vault encrypt --vault-password-file=/Users/cloud/.ansible/cube.key roles/sslcert/tasks/main.yml
	ansible-vault encrypt --vault-password-file=/Users/cloud/.ansible/cube.key roles/sslcert/templates/openssl.conf.j2
	ansible-vault encrypt --vault-password-file=/Users/cloud/.ansible/cube.key roles/acloud/vars/main.yaml

decrypt:
	echo "Decrypt files ..."
	ansible-vault decrypt --vault-password-file=/Users/cloud/.ansible/cube.key roles/sslcert/tasks/main.yml
	ansible-vault decrypt --vault-password-file=/Users/cloud/.ansible/cube.key roles/sslcert/templates/openssl.conf.j2
	ansible-vault decrypt --vault-password-file=/Users/cloud/.ansible/cube.key roles/acloud/vars/main.yaml
