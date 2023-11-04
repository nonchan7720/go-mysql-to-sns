package config

import (
	"crypto/tls"
	"database/sql"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Database struct {
	Host              string `yaml:"host"`
	Port              int    `yaml:"port" default:"3306"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	DBName            string `yaml:"name" default:"mysql"`
	SSHTunnel         bool   `yaml:"sshTunnel"`
	TLS               *TLS   `yaml:"tls"`
	MaxOpenConn       int    `yaml:"maxOpenConn" default:"10"`
	MaxLifeTimeSecond int    `yaml:"maxLifeTimeSecond" default:"300"`
	MaxIdleConn       int    `yaml:"maxIdleConn" default:"1"`
	MaxIdleSecond     int    `yaml:"maxIdleSecond" default:"0"`
}

var (
	_ validation.Validatable = (*Database)(nil)
)

func (d *Database) Tls() *tls.Config {
	if d.TLS == nil {
		return nil
	}
	return d.TLS.Config()
}

func (d Database) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Host, validation.Required),
		validation.Field(&d.Port, validation.Required),
		validation.Field(&d.Username, validation.Required),
		validation.Field(&d.Password, validation.Required),
		validation.Field(&d.DBName, validation.Required),
	)
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
