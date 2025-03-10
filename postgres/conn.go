package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mekramy/gosql"
)

// ConfigModifier modifies database config before creation.
type ConfigModifier func(*pgxpool.Config)

// New creates a new PostgreSQL connection using a DSN string.
// Accepts optional ConfigModifier functions to adjust the configuration before establishing the connection.
func New(ctx context.Context, dsn string, modifiers ...ConfigModifier) (Connection, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	for _, modifier := range modifiers {
		modifier(config)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &pgxConnection{db: pool}, nil
}

// Connection represents a PostgreSQL database connection.
type Connection interface {
	// Database returns the underlying connection pool.
	Database() *pgxpool.Pool

	// Ping verifies the database connection by sending a simple query.
	Ping(ctx context.Context) error

	// Transaction executes a function within a transaction.
	// Commits if successful, rolls back on error.
	Transaction(ctx context.Context, cb func(pgx.Tx) error, options ...pgx.TxOptions) error

	// Exec executes a raw SQL command on the connection's database.
	Exec(ctx context.Context, sql string, arguments ...any) error

	// Scan executes a raw SQL query on the connection's database and return a standard scanner.
	Scan(ctx context.Context, sql string, arguments ...any) (gosql.Scanner, error)

	// Close terminates the database connection pool.
	Close() error
}

type pgxConnection struct {
	db *pgxpool.Pool
}

func (d *pgxConnection) Database() *pgxpool.Pool {
	return d.db
}

func (d *pgxConnection) Ping(ctx context.Context) error {
	return d.db.Ping(ctx)
}

func (d *pgxConnection) Transaction(ctx context.Context, f func(pgx.Tx) error, opts ...pgx.TxOptions) error {
	tx, err := d.db.BeginTx(ctx, parseVariadic(pgx.TxOptions{}, opts...))
	if err != nil {
		return err
	}

	if err := f(tx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (d *pgxConnection) Exec(c context.Context, s string, args ...any) error {
	_, err := d.db.Exec(c, s, args...)
	return err
}

func (d *pgxConnection) Scan(c context.Context, s string, args ...any) (gosql.Scanner, error) {
	rows, err := d.db.Query(c, s, args...)
	if err != nil {
		return nil, err
	}

	return &scanner{rows: rows}, nil
}

func (d *pgxConnection) Close() error {
	d.db.Close()
	return nil
}
