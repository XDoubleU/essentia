package database

import (
	"context"

	"github.com/XDoubleU/essentia/pkg/logger"
)

type MinimalDB interface {
	Close()
}

type MainTestEnv[TDb MinimalDB, TTx MinimalDBTx] struct {
	TestDB      TDb
	beginTxFunc func(ctx context.Context, db TDb) TTx
}

func SetupGlobal[TDb MinimalDB, TTx MinimalDBTx](
	testDB TDb,
	beginTxFunc func(ctx context.Context, db TDb) TTx,
) (*MainTestEnv[TDb, TTx], error) {
	logger.SetLogger(logger.NullLogger)

	mainTestEnv := MainTestEnv[TDb, TTx]{
		TestDB:      testDB,
		beginTxFunc: beginTxFunc,
	}

	return &mainTestEnv, nil
}

func (mainTestEnv *MainTestEnv[TDb, TTx]) TeardownGlobal() error {
	mainTestEnv.TestDB.Close()
	return nil
}
