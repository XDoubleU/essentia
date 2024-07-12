package postgres

import (
	"github.com/XDoubleU/essentia/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MainTestEnv struct {
	TestDB *pgxpool.Pool
}

func SetupGlobal(dbDsn string, dbMaxConns int,
	dbMaxIdletime string) (*MainTestEnv, error) {
	logger.SetLogger(logger.NullLogger)

	testDB, err := Connect(
		dbDsn,
		dbMaxConns,
		dbMaxIdletime,
	)
	if err != nil {
		return nil, err
	}

	mainTestEnv := MainTestEnv{
		TestDB: testDB,
	}

	return &mainTestEnv, nil
}

func TeardownGlobal(mainTestEnv *MainTestEnv) error {
	mainTestEnv.TestDB.Close()
	return nil
}
