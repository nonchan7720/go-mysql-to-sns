package client

import (
	"database/sql"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/ent"
)

var (
	currentDB     *sql.DB
	currentClient *ent.Client
)

func SetDB(db *sql.DB) {
	currentDB = db
	SetClient(newEntClient(db))
}

func GetDB() *sql.DB {
	return currentDB
}

func GetClient() *ent.Client {
	client := currentClient
	if client == nil {
		drv := entsql.OpenDB("mysql", currentDB)
		client = ent.NewClient(ent.Driver(drv))
	}
	return client
}

func SetClient(client *ent.Client) {
	currentClient = client
}

func newEntClient(db *sql.DB) *ent.Client {
	drv := entsql.OpenDB("mysql", currentDB)
	return ent.NewClient(ent.Driver(drv))
}
