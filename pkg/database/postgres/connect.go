package postgres

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect connects to postgres and returns a [*pgxpool.Pool].
//
// The provided arguments are:
//   - dsn: the url to reach the database
//   - maxConns: the maximum amount of open connections
//   - maxIdleTime: the maximum idle time of an open connection
//   - connectTimeout: the timeout on connecting to the database
//   - sleepBeforeRetry: duration to sleep before trying to connect again
//   - maxRetryDuration: total amount of time to try and achieve a database connection
func Connect(
	logger *log.Logger,
	dsn string,
	maxConns int,
	maxIdleTime string,
	connectTimeout string,
	sleepBeforeRetry time.Duration,
	maxRetryDuration time.Duration,
) (*pgxpool.Pool, error) {
	connString, err := setupConnString(dsn, maxConns, maxIdleTime, connectTimeout)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	ctxTimeout, err := strconv.ParseInt(connectTimeout, 10, 64)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	for time.Now().Compare(start.Add(maxRetryDuration)) < 0 {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(ctxTimeout)*time.Second,
		)
		defer cancel()

		err = db.Ping(ctx)
		if err == nil {
			break
		}

		logger.
			Printf(
				"can't connect to database (%v), retrying in %s",
				err,
				sleepBeforeRetry,
			)
		time.Sleep(sleepBeforeRetry)
	}

	if err != nil {
		return nil, fmt.Errorf("can't connect to database (%w)", err)
	}

	return db, nil
}

func setupConnString(
	dsn string,
	maxConns int,
	maxIdleTime string,
	connectTimeout string) (string, error) {
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}

	queryValues := parsedURL.Query()

	queryValues.Add("pool_max_conns", strconv.Itoa(maxConns))
	queryValues.Add("pool_max_conn_idle_time", maxIdleTime)
	queryValues.Add("connect_timeout", connectTimeout)

	parsedURL.RawQuery = queryValues.Encode()

	return parsedURL.String(), nil
}
