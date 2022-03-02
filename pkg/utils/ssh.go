package utils

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"time"
)

type SSH struct {
	IP      string
	User    string
	Cert    string
	Port    int
	session *ssh.Session
	client  *ssh.Client
}

func (S *SSH) readPublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func (S *SSH) Connect() {
	var sshConfig *ssh.ClientConfig
	var auth []ssh.AuthMethod

	auth = []ssh.AuthMethod{
		S.readPublicKeyFile(S.Cert),
	}

	sshConfig = &ssh.ClientConfig{
		User: S.User,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * 3,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", S.IP, S.Port), sshConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
		client.Close()
		return
	}

	S.session = session
	S.client = client
}

func (S *SSH) RunCmd(cmd string) string {
	out, err := S.session.CombinedOutput(cmd)
	if err != nil {
		fmt.Println(err)
	}
	return string(out)
}

func (S *SSH) Close() {
	S.session.Close()
	S.client.Close()
}
