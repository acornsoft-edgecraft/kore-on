#!/bin/bash

#if [ "${1:0:1}" = '-' ]; then
#	set -- cube "$@"
#fi
#echo "${DEV_IP} ${DEV_DOMAIN}" >> /etc/hosts

#rm -rf /cube/cubescripts/roles/acloud
#rm -rf /cube/cubescripts/acloud.yml
#mkdir /cube/cubescripts/roles/acloud
#/usr/local/bin/dockerd &



if [ ! -f "/etc/ssh/ssh_host_rsa_key" ]; then
	# generate fresh rsa key
	ssh-keygen -q -f /etc/ssh/ssh_host_rsa_key -N '' -t rsa
fi
if [ ! -f "/etc/ssh/ssh_host_dsa_key" ]; then
	# generate fresh dsa key
	ssh-keygen -q -f /etc/ssh/ssh_host_dsa_key -N '' -t dsa
fi

if [ ! -f "/etc/ssh/ssh_host_ecdsa_key" ]; then
	# generate fresh ecdsa key
	ssh-keygen -q -f /etc/ssh/ssh_host_ecdsa_key -N '' -t ecdsa
fi

if [ ! -f "/etc/ssh/ssh_host_ed25519_key" ]; then
	# generate fresh ed25519 key
	ssh-keygen -q -f /etc/ssh/ssh_host_ed25519_key -N '' -t ed25519
fi

#prepare run dir
if [ ! -d "/var/run/sshd" ]; then
  mkdir -p /var/run/sshd
fi

echo "PasswordAuthentication no" >>  /etc/ssh/sshd_config

sed -i 's/^root.*/root:$6$J62il4lEm77WhBaj$Dbe.W1AaXvczdM0bEVYOcNCSt4V1Xspq.g6JYNQ9MAqWQD.vRQ2GXYGweQiHdgPIwR\/ZeoLsbVcPN.Lfa3WF8.:18001:0:::::/g' /etc/shadow

#expect -c "
#spawn passwd root
#expect \'New password:\'
#send \"@SDSDSDS1234\\r\"
#expect \'Retype password:\'
#send \"@SDSDSDS1234\\r\"
#set timeout 600
#expect eof"


cp  /etc/ssh/ssh_host_rsa_key.pub ~/.ssh/authorized_keys
cp  /etc/ssh/ssh_host_rsa_key ~/.ssh/id_rsa

/usr/sbin/sshd -D &

exec "$@"

