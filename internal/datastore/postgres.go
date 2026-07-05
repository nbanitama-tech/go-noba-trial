package datastore

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func OpenPostgres(ctx context.Context, databaseURL string) (*sql.DB, error) {
	const (
		maxAttempts = 10
		baseDelay   = 500 * time.Millisecond
	)

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err := sql.Open("postgres", databaseURL)
		if err != nil {
			lastErr = err
			sleep(ctx, baseDelay, attempt)
			continue
		}

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(30 * time.Minute)

		if err := db.PingContext(ctx); err != nil {
			lastErr = err
			_ = db.Close()
			sleep(ctx, baseDelay, attempt)
			continue
		}

		return db, nil
	}

	return nil, lastErr
}

func sleep(ctx context.Context, baseDelay time.Duration, attempt int) {
	delay := baseDelay * time.Duration(attempt)
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}
