package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/creasty/defaults"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type SSH struct {
	PrivateKey        string   `yaml:"privateKey"`
	Host              string   `yaml:"host"`
	Port              int      `yaml:"port"`
	Username          string   `yaml:"username"`
	HostKeyAlgorithms []string `yaml:"hostKeyAlgorithms" default:"[\"ssh-ed25519\"]"`
	KnownHosts        string   `yaml:"knownHosts"`
}

func (s *SSH) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(s); err != nil {
		return err
	}
	s.KnownHosts = filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	type plain SSH
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}
	return nil
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
	hostKeyCallback, err := knownhosts.New(conf.KnownHosts)
	if err != nil {
		return nil, err
	}

	sshConf := &ssh.ClientConfig{
		User: conf.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback:   hostKeyCallback,
		HostKeyAlgorithms: conf.HostKeyAlgorithms,
	}
	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port), sshConf)
}
