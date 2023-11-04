package config

import (
	"crypto/tls"
	"database/sql"
	"time"
)

type Database struct {
	Host              string `yaml:"host" validate:"required"`
	Port              int    `yaml:"port" default:"3306" validate:"required"`
	Username          string `yaml:"username" validate:"required"`
	Password          string `yaml:"password" validate:"required"`
	DBName            string `yaml:"name" default:"mysql" validate:"required"`
	SSHTunnel         bool   `yaml:"sshTunnel"`
	TLS               *TLS   `yaml:"tls"`
	MaxOpenConn       int    `yaml:"MAX_OPEN_CONN" default:"10"`
	MaxLifeTimeSecond int    `yaml:"MAX_LIFE_TIME_SECOND" default:"300"`
	MaxIdleConn       int    `yaml:"MAX_IDLE_CONN" default:"1"`
	MaxIdleSecond     int    `yaml:"MAX_IDLE_SECOND" default:"0"`
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
