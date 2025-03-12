package migration

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mekramy/gosql/postgres"
)

// Implement MigrationSource
type postgresSource struct {
	conn postgres.Connection
}

func (ps *postgresSource) Transaction(c context.Context, cb func(ExecutableScanner) error) error {
	tx, err := ps.conn.Database().BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return err
	}

	px := &postgresTx{tx: tx}
	if err := cb(px); err != nil {
		tx.Rollback(c)
		return err
	}

	return tx.Commit(c)
}

func (ps *postgresSource) Exec(c context.Context, s string, args ...any) error {
	_, err := ps.conn.Database().Exec(c, s, args...)
	return err
}

func (ps *postgresSource) Scan(c context.Context, s string, args ...any) (Rows, error) {
	rows, err := ps.conn.Database().Query(c, s, args...)
	if err != nil {
		return nil, err
	}

	return &postgresRows{rows: rows}, nil
}

// Implement ExecutableScanner for transaction
type postgresTx struct {
	tx pgx.Tx
}

func (px *postgresTx) Exec(c context.Context, s string, args ...any) error {
	_, err := px.tx.Exec(c, s, args...)
	return err
}

func (px *postgresTx) Scan(c context.Context, s string, args ...any) (Rows, error) {
	rows, err := px.tx.Query(c, s, args...)
	if err != nil {
		return nil, err
	}

	return &postgresRows{rows: rows}, nil
}

// Implement Scanner row
type postgresRows struct {
	rows pgx.Rows
}

func (ps *postgresRows) Next() bool {
	return ps.rows.Next()
}

func (ps *postgresRows) Scan(dest ...any) error {
	return ps.rows.Scan(dest...)
}

func (ps *postgresRows) Close() {
	ps.rows.Close()
}
