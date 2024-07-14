package database

import (
	"context"
)

// TestEnv encapsulates test transaction.
type TestEnv[TTx MinimalDBTx] struct {
	Tx  TTx
	ctx context.Context
}

// SetupSingle starts test transaction.
func (mainTestEnv *MainTestEnv[TDB, TTx]) SetupSingle() TestEnv[TTx] {
	ctx := context.Background()

	testEnv := TestEnv[TTx]{
		Tx:  mainTestEnv.beginTxFunc(ctx, mainTestEnv.TestDB),
		ctx: ctx,
	}

	return testEnv
}

// TeardownSingle executes rollback of test transaction.
func (testEnv *TestEnv[TTx]) TeardownSingle() {
	err := testEnv.Tx.Rollback(testEnv.ctx)
	if err != nil {
		panic(err)
	}
}
