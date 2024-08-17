package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/pkg/database/postgres"
	"github.com/xdoubleu/essentia/pkg/logging"
)

func setup(t *testing.T) *postgres.PgxSyncTx {
	t.Helper()

	logger := logging.NewNopLogger()

	db, err := postgres.Connect(
		logger,
		"postgres://postgres@localhost/postgres",
		25, //nolint:mnd //no magic number
		"15m",
		30,             //nolint:mnd //no magic number
		30*time.Second, //nolint:mnd //no magic number
		5*time.Minute,  //nolint:mnd //no magic number
	)
	if err != nil {
		panic(err)
	}

	return postgres.CreatePgxSyncTx(context.Background(), db)
}

func TestPing(t *testing.T) {
	tx := setup(t)
	defer tx.Rollback(context.Background())

	db := postgres.NewSpanDB(tx)

	err := db.Ping(context.Background())
	assert.Nil(t, err)
}

func TestExec(t *testing.T) {
	tx := setup(t)
	defer tx.Rollback(context.Background())

	db := postgres.NewSpanDB(tx)

	_, err := db.Exec(context.Background(), "")
	assert.Nil(t, err)
}

func TestQuery(t *testing.T) {
	tx := setup(t)
	defer tx.Rollback(context.Background())

	db := postgres.NewSpanDB(tx)

	_, err := db.Exec(
		context.Background(),
		"CREATE TABLE kv (key VARCHAR(255), value VARCHAR(255));",
	)
	require.Nil(t, err)

	_, err = db.Exec(
		context.Background(),
		"INSERT INTO kv (key, value) VALUES ('key1', 'value1'), ('key2', 'value2');",
	)
	require.Nil(t, err)

	rows, err := db.Query(context.Background(), "SELECT key, value FROM kv;")
	require.Nil(t, err)

	results := make([][]string, 2)
	results[0] = make([]string, 2)
	results[1] = make([]string, 2)

	count := 0
	for rows.Next() {
		err = rows.Scan(&results[count][0], &results[count][1])
		require.Nil(t, err)

		count++
	}

	require.Nil(t, err)
	assert.Equal(t, "key1", results[0][0])
	assert.Equal(t, "value1", results[0][1])
	assert.Equal(t, "key2", results[1][0])
	assert.Equal(t, "value2", results[1][1])
}

func TestQueryRow(t *testing.T) {
	tx := setup(t)
	defer tx.Rollback(context.Background())

	db := postgres.NewSpanDB(tx)

	_, err := db.Exec(
		context.Background(),
		"CREATE TABLE kv (key VARCHAR(255), value VARCHAR(255));",
	)
	require.Nil(t, err)

	_, err = db.Exec(
		context.Background(),
		"INSERT INTO kv (key, value) VALUES ('key1', 'value1');",
	)
	require.Nil(t, err)

	var key string
	var value string
	err = db.QueryRow(context.Background(), "SELECT key, value FROM kv WHERE key = 'key1';").
		Scan(&key, &value)

	require.Nil(t, err)
	assert.Equal(t, "key1", key)
	assert.Equal(t, "value1", value)
}

func TestBegin(t *testing.T) {
	tx := setup(t)
	defer tx.Rollback(context.Background())

	db := postgres.NewSpanDB(tx)

	_, err := db.Begin(context.Background())
	assert.Nil(t, err)
}
