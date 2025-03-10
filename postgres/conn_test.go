package postgres_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mekramy/gosql/postgres"
)

func TestConfigBuildDSN(t *testing.T) {
	config := postgres.NewConfig().
		Host("localhost").
		Port(5432).
		User("postgres").
		Password("password").
		Database("test").
		SSLMode("disable").
		MaxConns(7).
		MinConns(2)

	dsn := config.Build()
	expectedDSN := "postgres://postgres:password@localhost:5432/test?pool_max_conns=7&pool_min_conns=2&sslmode=disable"
	if dsn != expectedDSN {
		t.Fatalf("expected %v, got %v", expectedDSN, dsn)
	}
}

func TestNewConnection(t *testing.T) {
	ctx := context.Background()
	config := postgres.NewConfig().
		Host("localhost").
		Port(5432).
		User("postgres").
		Password("root").
		Database("postgres")

	conn, err := postgres.New(
		ctx, config.Build(),
		func(c *pgxpool.Config) { c.MaxConns = 7 },
		func(c *pgxpool.Config) {
			if c.MaxConns != 7 {
				t.Fatalf("connection modifiers fail")
			}
		},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer conn.Close()
}

func TestConnectionManager(t *testing.T) {
	ctx := context.Background()
	config := postgres.NewConfig().
		Host("localhost").
		Port(5432).
		User("postgres").
		Password("root").
		Database("test").
		SSLMode("disable")

	manager := postgres.NewConnectionManager(config)
	defer manager.Close()

	err := manager.Connect(ctx, "test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, exists := manager.Get("test")
	if !exists {
		t.Fatalf("expected connection to exist")
	}

	err = manager.Remove("test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, exists = manager.Get("test")
	if exists {
		t.Fatalf("expected connection to not exist")
	}
}

func TestTransaction(t *testing.T) {
	ctx := context.Background()
	config := postgres.NewConfig().
		Host("localhost").
		Port(5432).
		User("postgres").
		Password("root").
		Database("test").
		SSLMode("disable")

	conn, err := postgres.New(ctx, config.Build())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer conn.Close()

	err = conn.Transaction(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "CREATE TABLE IF NOT EXISTS test (id SERIAL PRIMARY KEY, name TEXT)")
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "INSERT INTO test (name) VALUES ($1)", "testname")
		return err
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
