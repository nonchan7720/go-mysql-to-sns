package client

import (
	"context"
	"testing"

	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestConnection(t *testing.T) {
	conf := config.Config{
		Database: config.Database{
			Host:     "db",
			Port:     3306,
			Username: "admin",
			Password: "pass1234",
			DBName:   "db",
		},
	}
	c, err := NewDB(&conf)
	require.NoError(t, err)
	o, err := c.Outbox.Query().First(context.Background())
	require.NoError(t, err)
	require.Equal(t, string(o.Payload), `{"test": "value"}`)
}
