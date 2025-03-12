package mysql_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/mekramy/gosql/mysql"
)

func TestConfigBuildDSN(t *testing.T) {
	config := mysql.NewConfig().
		Host("localhost").
		User("root").
		Password("root").
		Database("test")

	dsn := config.Build()
	expectedDSN := "root:root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true"
	if dsn != expectedDSN {
		t.Fatalf("expected %v, got %v", expectedDSN, dsn)
	}
}

func TestNewConnection(t *testing.T) {
	ctx := context.Background()
	config := mysql.NewConfig().
		Host("localhost").
		User("root").
		Password("root").
		Database("test")

	conn, err := mysql.New(
		ctx, config.Build(),
		func(d *sql.DB) {
			d.SetMaxOpenConns(10)
		},
		func(d *sql.DB) {
			if d.Stats().MaxOpenConnections != 10 {
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
	config := mysql.NewConfig().
		Host("localhost").
		User("root").
		Password("root")

	manager := mysql.NewConnectionManager(config)
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
	config := mysql.NewConfig().
		Host("localhost").
		User("root").
		Password("root").
		Database("test")

	conn, err := mysql.New(ctx, config.Build())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer conn.Close()

	err = conn.Transaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS test (id SERIAL PRIMARY KEY, name TEXT)")
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, "INSERT INTO test (name) VALUES (?)", "testname")
		return err
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
