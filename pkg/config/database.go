package config

import (
	"crypto/tls"
	"database/sql"
	"time"

	"github.com/creasty/defaults"
)

type Database struct {
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	DBName            string `yaml:"name"`
	SSHTunnel         bool   `yaml:"sshTunnel"`
	TLS               *TLS   `yaml:"tls"`
	MaxOpenConn       int    `yaml:"MAX_OPEN_CONN" default:"10"`
	MaxLifeTimeSecond int    `yaml:"MAX_LIFE_TIME_SECOND" default:"300"`
	MaxIdleConn       int    `yaml:"MAX_IDLE_CONN" default:"1"`
	MaxIdleSecond     int    `yaml:"MAX_IDLE_SECOND" default:"0"`
}

func (d *Database) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(d); err != nil {
		return err
	}
	type plain Database
	if err := unmarshal((*plain)(d)); err != nil {
		return err
	}
	return nil
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

func setDB(db *sql.DB, cfg Database) {
	if cfg.MaxIdleSecond > 0 {
		db.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleSecond) * time.Second)
	}
	if cfg.MaxLifeTimeSecond > 0 {
		db.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTimeSecond) * time.Second)
	}
	if cfg.MaxIdleConn > 0 {
		db.SetMaxIdleConns(int(cfg.MaxIdleConn))
	}
	if cfg.MaxOpenConn > 0 {
		db.SetMaxOpenConns(int(cfg.MaxOpenConn))
	}
}
