package database

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type ConnectionFactoryConfig struct {
	MasterConnectionString string
	SyncConnectionString   string
	AsyncConnectionString  string
}

type ConnectionFactory struct {
	masterDb *sql.DB
	syncDb   *sql.DB
	asyncDb  *sql.DB
}

const driver = "pgx"

func NewConnectionFactory(cfg ConnectionFactoryConfig) *ConnectionFactory {
	ctx := context.Background()

	masterDb, err := sql.Open(driver, cfg.MasterConnectionString)

	if err != nil {
		panic(err)
	}

	err = masterDb.PingContext(ctx)

	if err != nil {
		panic(err)
	}

	syncDb, err := sql.Open(driver, cfg.SyncConnectionString)

	if err != nil {
		panic(err)
	}

	err = syncDb.PingContext(ctx)

	if err != nil {
		panic(err)
	}

	asyncDb, err := sql.Open(driver, cfg.AsyncConnectionString)

	if err != nil {
		panic(err)
	}

	err = asyncDb.PingContext(ctx)

	if err != nil {
		panic(err)
	}

	return &ConnectionFactory{
		masterDb: masterDb,
		syncDb:   syncDb,
		asyncDb:  asyncDb,
	}
}

func (f *ConnectionFactory) GetMaster() *sql.DB {
	return f.masterDb
}

func (f *ConnectionFactory) GetSync() *sql.DB {
	return f.syncDb
}

func (f *ConnectionFactory) GetAsync() *sql.DB {
	return f.asyncDb
}
