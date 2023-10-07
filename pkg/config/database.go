package config

import "crypto/tls"

type Database struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	SSHTunnel bool   `yaml:"sshTunnel"`
	TLS       *TLS   `yaml:"tls"`
}

func (d *Database) Tls() *tls.Config {
	if d.TLS == nil {
		return nil
	}
	return d.TLS.Config()
}

type TLS struct {
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify"`
	SeverName          string `yaml:"serverName"`
}

func (t *TLS) Config() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: t.InsecureSkipVerify,
		ServerName:         t.SeverName,
	}
}
