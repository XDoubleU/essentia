package database

import (
	"context"
)

type TestEnv[TTx MinimalDBTx] struct {
	Tx  TTx
	ctx context.Context
}

func (mainTestEnv *MainTestEnv[TDb, TTx]) SetupSingle() TestEnv[TTx] {
	ctx := context.Background()

	testEnv := TestEnv[TTx]{
		Tx:  mainTestEnv.beginTxFunc(ctx, mainTestEnv.TestDB),
		ctx: ctx,
	}

	return testEnv
}

func (testEnv *TestEnv[TTx]) TeardownSingle() {
	err := testEnv.Tx.Rollback(testEnv.ctx)
	if err != nil {
		panic(err)
	}
}
