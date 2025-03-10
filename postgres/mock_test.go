package postgres_test

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mekramy/gosql/postgres"
)

func NewMockReadable() postgres.Readable {
	return &MockReader{}
}

func NewMockExecutable(expect string) postgres.Executable {
	return &MockExecutable{expect: expect}
}

// MockExecutable is a mock implementation of postgres.Executable
type MockReader struct{}

func (r *MockReader) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}
func (r *MockReader) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return &MockRow{}
}

// MockExecutable is a mock implementation of postgres.Executable
type MockExecutable struct {
	expect string
}

func (e *MockExecutable) Exec(_ context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if e.expect != sql {
		return pgconn.CommandTag{}, fmt.Errorf("expected %s, got %s", e.expect, sql)
	}

	return pgconn.CommandTag{}, nil
}

// MockRow is a mock implementation of pgx.Row
type MockRow struct{}

func (r *MockRow) Scan(dest ...any) error {
	if res, ok := dest[0].(*User); ok {
		res.Id = 1
		res.Name = "John Doe"
		return nil
	} else if res, ok := dest[0].(*int64); ok {
		*res = 100
		return nil
	} else {
		return fmt.Errorf("invalid type for scanner")
	}
}
