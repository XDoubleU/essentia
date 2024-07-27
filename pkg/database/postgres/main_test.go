package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/internal/mocks"
	"github.com/xdoubleu/essentia/pkg/database"
	"github.com/xdoubleu/essentia/pkg/database/postgres"
	sentrytools "github.com/xdoubleu/essentia/pkg/sentry"
)

type pair struct {
	Key   int
	Value string
}

func newPair() pair {
	return pair{
		Key:   0,
		Value: "",
	}
}

func TestSetup(t *testing.T) {
	mockedLogger := mocks.NewMockedLogger()

	db, err := postgres.Connect(
		mockedLogger.Logger(),
		"postgres://postgres@localhost/postgres",
		1,
		"1m",
		5,
		15*time.Second,
		30*time.Second,
	)
	require.Nil(t, err)
	defer db.Close()

	mainTestEnv := database.CreateMainTestEnv(db, postgres.CreatePgxSyncTx)

	ctx := context.Background()
	CreateTable(ctx, t, mainTestEnv.TestDB)

	testEnv := mainTestEnv.SetupSingle()

	OperationsInTx(ctx, t, testEnv.Tx)
	testEnv.TeardownSingle()

	OperationsAfterTx(ctx, t, mainTestEnv.TestDB)

	DropTable(ctx, t, mainTestEnv.TestDB)

	assert.Contains(t, mockedLogger.CapturedLogs(), "connected to database")
}

func CreateTable(ctx context.Context, t *testing.T, db postgres.DB) {
	_, err := db.Exec(
		ctx,
		"CREATE TABLE IF NOT EXISTS pairs (key int PRIMARY KEY, value varchar(255) NOT NULL)",
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating table", err)
	}
}

func DropTable(ctx context.Context, t *testing.T, db postgres.DB) {
	_, err := db.Exec(ctx, "DROP TABLE IF EXISTS pairs")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when dropping table", err)
	}
}

func OperationsInTx(ctx context.Context, t *testing.T, tx postgres.DB) {
	rows, err := tx.Query(
		ctx,
		"INSERT INTO pairs (key, value) VALUES ($1, $2)",
		1,
		"test",
	)
	rows.Close()
	require.Nil(t, err)

	p := newPair()

	// wrap with span to get some additional coverage
	sentry.SetHubOnContext(ctx, sentrytools.MockedSentryHub())
	row := database.WrapWithSpanNoError(
		ctx,
		"test-postgresql",
		tx.QueryRow,
		"SELECT value FROM pairs WHERE key = $1", 1,
	)

	err = row.Scan(&p.Value)
	require.Nil(t, err)
	assert.Equal(t, "test", p.Value)
}

func OperationsAfterTx(ctx context.Context, t *testing.T, db postgres.DB) {
	p := newPair()
	err := db.QueryRow(ctx, "SELECT value FROM pairs WHERE key = 1").Scan(&p.Value)
	assert.ErrorIs(t, err, pgx.ErrNoRows)
}
