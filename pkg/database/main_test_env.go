package database

import (
	"context"
)

// MainTestEnv takes care of the database
// and transactions during tests.
type MainTestEnv[TDB any, TTx MinimalDBTx] struct {
	TestDB      TDB
	beginTxFunc func(ctx context.Context, db TDB) TTx
}

// CreateMainTestEnv creates a new [MainTestEnv].
func CreateMainTestEnv[TDB any, TTx MinimalDBTx](
	testDB TDB,
	beginTxFunc func(ctx context.Context, db TDB) TTx,
) *MainTestEnv[TDB, TTx] {
	return &MainTestEnv[TDB, TTx]{
		TestDB:      testDB,
		beginTxFunc: beginTxFunc,
	}
}
