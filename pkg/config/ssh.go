package config

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type SSH struct {
	PrivateKey string `yaml:"privateKey"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Username   string `yaml:"username"`
}

func (conf *SSH) Conn() (*ssh.Client, error) {
	sshKey, err := os.ReadFile(conf.PrivateKey)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(sshKey)
	if err != nil {
		return nil, err
	}
	sshConf := &ssh.ClientConfig{
		User: conf.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port), sshConf)
}
