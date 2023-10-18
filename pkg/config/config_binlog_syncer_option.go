package config

import (
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/siddontang/go-log/loggers"
)

type Option interface {
	apply(cfg *replication.BinlogSyncerConfig)
}

type optionFn func(cfg *replication.BinlogSyncerConfig)

func (fn optionFn) apply(cfg *replication.BinlogSyncerConfig) {
	fn(cfg)
}

func WithDiscardGTID() Option {
	return optionFn(func(cfg *replication.BinlogSyncerConfig) {
		cfg.DiscardGTIDSet = true
	})
}

func WithSemiSyncEnable() Option {
	return optionFn(func(cfg *replication.BinlogSyncerConfig) {
		cfg.SemiSyncEnabled = true
	})
}

func WithLogger(logger loggers.Advanced) Option {
	return optionFn(func(cfg *replication.BinlogSyncerConfig) {
		cfg.Logger = logger
	})
}

func WithFlavorForMariaDB() Option {
	return optionFn(func(cfg *replication.BinlogSyncerConfig) {
		cfg.Flavor = "mariadb"
	})
}

func WithDisableRetrySync() Option {
	return optionFn(func(cfg *replication.BinlogSyncerConfig) {
		cfg.DisableRetrySync = true
	})
}

func WithMaxReconnectAttempts(value int) Option {
	return optionFn(func(cfg *replication.BinlogSyncerConfig) {
		cfg.MaxReconnectAttempts = value
	})
}
