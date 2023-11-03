package client

import (
	"context"
	"testing"

	"github.com/nonchan7720/go-storage-to-messenger/pkg/config"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestConnection(t *testing.T) {
	var host string
	if utils.IsCI() {
		host = "localhost"
	} else {
		host = "db"
	}
	conf := config.Config{
		Database: config.Database{
			Host:     host,
			Port:     3306,
			Username: "admin",
			Password: "pass1234",
			DBName:   "db",
		},
	}
	_, err := NewDB(context.Background(), &conf)
	require.NoError(t, err)
	// o, err := c.Outbox.Query().First(context.Background())
	// require.NoError(t, err)
	// require.Equal(t, string(o.Payload), `{"test": "value"}`)
}
