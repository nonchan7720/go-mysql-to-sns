package config

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database  Database    `yaml:"database"`
	SSH       SSH         `yaml:"ssh"`
	Saver     BinlogSaver `yaml:"saver"`
	Publisher Publisher   `yaml:"publisher"`
}

func LoadConfig(filePath string) (*Config, error) {
	f, err := NewExpandEnv(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var config Config
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}
	if err := config.Validation(); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) Connect(ctx context.Context) (*sql.DB, error) {
	mysqlNet := "tcp"
	if c.Database.SSHTunnel {
		var dialFunc mysql.DialContextFunc
		sshClient, err := c.SSH.Conn()
		if err != nil {
			return nil, err
		}
		mysqlNet = "mysql+tcp"
		dialFunc = func(ctx context.Context, addr string) (net.Conn, error) {
			return sshClient.Dial("tcp", addr)
		}
		mysql.RegisterDialContext(mysqlNet, dialFunc)
	}
	mysqlConfig := mysql.Config{
		User:                 c.Database.Username,
		Passwd:               c.Database.Password,
		Addr:                 fmt.Sprintf("%s:%d", c.Database.Host, c.Database.Port),
		Net:                  mysqlNet,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
		TLS:                  c.Database.Tls(),
	}

	dsn := mysqlConfig.FormatDSN()
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetConnMaxLifetime(1 * time.Hour)
	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)
	if err := conn.PingContext(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Config) NewBinlogSyncer(serverId int) (*replication.BinlogSyncer, error) {
	var dialFunc client.Dialer
	if c.Database.SSHTunnel {
		sshClient, err := c.SSH.Conn()
		if err != nil {
			return nil, err
		}
		dialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return sshClient.Dial(network, addr)
		}
	}

	cfg := replication.BinlogSyncerConfig{
		ServerID:  uint32(serverId),
		Flavor:    "mysql",
		Host:      c.Database.Host,
		Port:      uint16(c.Database.Port),
		User:      c.Database.Username,
		Password:  c.Database.Password,
		Dialer:    dialFunc,
		TLSConfig: c.Database.Tls(),
	}
	return replication.NewBinlogSyncer(cfg), nil
}

func (c *Config) Validation() error {
	if err := c.Publisher.Validation(); err != nil {
		return err
	}
	return nil
}
