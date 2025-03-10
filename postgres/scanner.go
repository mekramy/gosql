package postgres

import "github.com/jackc/pgx/v5"

type scanner struct {
	rows pgx.Rows
}

func (s *scanner) Next() bool {
	return s.rows.Next()
}

func (s *scanner) Scan(dest ...any) error {
	return s.rows.Scan(dest...)
}

func (s *scanner) Close() {
	s.rows.Close()
}
