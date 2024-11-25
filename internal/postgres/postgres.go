package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pgErrCode_NotNullViolation = "23502"
	pgErrCode_FKeyViolation    = "23503"
	pgErrCode_UniqueViolation  = "23505"
)

type DB struct {
	pool *pgxpool.Pool

	URL string
}

func NewDB(url string) *DB {
	db := &DB{
		URL: url,
	}
	return db
}

func (db *DB) Open(ctx context.Context) (err error) {
	if db.URL == "" {
		return errors.New("url required")
	}

	if db.pool, err = pgxpool.New(ctx, db.URL); err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

func (db *DB) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	if db.pool == nil {
		return nil, errors.New("pool required")
	}

	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (db *DB) BeginTx(ctx context.Context) (pgx.Tx, error) {
	if db.pool == nil {
		return nil, errors.New("pool required")
	}

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
