package config

import (
	"context"
	"testing"
	"time"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestConnectDB(t *testing.T) {
	var (
		host string
		ctx  context.Context = context.Background()
	)
	if utils.IsCI() {
		host = "localhost"
	} else {
		host = "db"
	}
	conf := Config{
		Database: Database{
			Host:     host,
			Port:     3306,
			Username: "admin",
			Password: "pass1234",
		},
	}
	db, err := conf.Connect(ctx)
	require.Nil(t, err)
	err = db.PingContext(ctx)
	require.Nil(t, err)
}

func TestBinlogSyncer(t *testing.T) {
	var (
		host string
		ctx  context.Context = context.Background()
	)
	if utils.IsCI() {
		host = "localhost"
	} else {
		host = "db"
	}
	conf := Config{
		Database: Database{
			Host:     host,
			Port:     3306,
			Username: "admin",
			Password: "pass1234",
		},
	}
	syncer, err := conf.NewBinlogSyncer(1)
	require.Nil(t, err)
	st, err := syncer.StartSync(mysql.Position{})
	require.Nil(t, err)
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()
	ev, err := st.GetEvent(ctx)
	require.Nil(t, err)
	require.NotNil(t, ev)
}
