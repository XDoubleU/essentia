package postgres

import (
	"context"
)

type Env struct {
	TestTx  SyncTx
	TestCtx context.Context
}

func SetupSingle(mainTestEnv *MainTestEnv) Env {
	testCtx := context.Background()
	testTx := CreateSyncTx(testCtx, mainTestEnv.TestDB)

	testEnv := Env{
		TestTx:  testTx,
		TestCtx: testCtx,
	}

	return testEnv
}

func TeardownSingle(testEnv Env) {
	err := testEnv.TestTx.Rollback(testEnv.TestCtx)
	if err != nil {
		panic(err)
	}
}
