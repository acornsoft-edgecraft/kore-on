## ssk key generate and copy
ssh-keygen -f /path/to/file -t rsa -N ''
ssh-copy-id -i /path/to/file root@ip

## run knit docker container
docker run -it --name=knit --rm -v ${PWD}:/cube/work regi.acloud.run/library/knit:1.0.0 /bin/bash

mkdir inventory
cp -rfp ../inventory/sample inventory/cube

## check ssh connetions
ansible -i inventory/cube/inventory.ini -u cube --private-key id_rsa all -m ping
ansible -i inventory/cube/inventory.ini -u cube --private-key id_rsa all -m setup

## create cluster
ansible-playbook -i inventory/cube/inventory.ini -u cube --private-key id_rsa ../scripts/cluster.yml

## upgrade cluster
ansible-playbook -i inventory/cube/inventory.ini -u cube --private-key id_rsa ../scripts/upgrade.yml

## add worker node
ansible-playbook -i inventory/cube/inventory.ini -u cube --private-key id_rsa ../scripts/add-node.yml

## remove worker node
ansible-playbook -i inventory/cube/inventory.ini -u cube --private-key id_rsa ../scripts/remove-node.yml

## reset cluster only
ansible-playbook -i inventory/cube/inventory.ini -u cube --private-key id_rsa ../scripts/reset.yml --tags reset-cluster

## reset all
ansible-playbook -i inventory/cube/inventory.ini -u cube --private-key id_rsa ../scripts/reset.yml

