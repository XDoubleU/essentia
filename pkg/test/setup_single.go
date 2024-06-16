package test

import (
	"context"
	"sync"
)

type Env struct {
	TestTx  Tx
	TestCtx context.Context
}

func SetupSingle(mainTestEnv *MainTestEnv) Env {
	testCtx := context.Background()

	var testTx Tx
	var mu sync.Mutex
	for {
		tx, err := mainTestEnv.TestDB.Begin(testCtx)
		if err == nil {
			testTx = Tx{
				tx: tx,
				mu: &mu,
			}
			break
		}
	}

	testEnv := Env{
		TestTx:  testTx,
		TestCtx: testCtx,
	}

	return testEnv
}

func TeardownSingle(testEnv Env) {
	for {
		if testEnv.TestTx.mu.TryLock() {
			break
		}
	}
	defer testEnv.TestTx.mu.Unlock()

	err := testEnv.TestTx.tx.Rollback(testEnv.TestCtx)
	if err != nil {
		panic(err)
	}
}
