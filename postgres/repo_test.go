package postgres_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/mekramy/gosql/postgres"
)

func TestRepository(t *testing.T) {
	type User struct {
		Id   int    `db:"id"`
		Name string `db:"name"`
	}
	ctx := context.Background()
	config := postgres.NewConfig().
		Host("localhost").
		Port(5432).
		User("postgres").
		Password("root").
		Database("test")

	conn, err := postgres.New(ctx, config.Build())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("Init", func(t *testing.T) {
		err := conn.Transaction(ctx, func(tx pgx.Tx) error {
			_, err := postgres.NewCmd(tx).
				Command("DROP TABLE IF EXISTS users CASCADE;").
				Exec(ctx)
			if err != nil {
				return err
			}

			_, err = postgres.NewCmd(tx).
				Command("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT);").
				Exec(ctx)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("Insert", func(t *testing.T) {
		err := conn.Transaction(ctx, func(tx pgx.Tx) error {
			for idx, name := range []string{"John Doe", "Jack Ma"} {
				u := User{
					Id:   idx + 1,
					Name: name,
				}

				_, err = postgres.NewInserter[User](tx).
					Table("users").
					Insert(ctx, u, postgres.OnlyFields("name"))
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		err := conn.Transaction(ctx, func(tx pgx.Tx) error {
			for idx, name := range []string{"John Doe New", "Jack Ma New"} {
				u := User{
					Id:   idx + 1,
					Name: name,
				}

				_, err = postgres.NewUpdater[User](tx).
					Table("users").
					Where("id = ?", u.Id).
					Update(ctx, u, postgres.SkipFields("id"))
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("Count", func(t *testing.T) {
		count, err := postgres.NewCounter(conn.Database()).
			Query("SELECT COUNT(*) FROM users;").
			Count(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if count != 2 {
			t.Fatalf("expected 2 users, got %d", count)
		}
	})

	t.Run("Single", func(t *testing.T) {
		jack, err := postgres.NewFinder[User](conn.Database()).
			Query("SELECT * FROM users WHERE id = ?;").
			Struct(ctx, 2)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if jack == nil {
			t.Fatal("expected user, got nil")
		}

		if jack.Name != "Jack Ma New" {
			t.Fatalf(`expected "Jack Ma New", got %s`, jack.Name)
		}
	})

	t.Run("Multiple", func(t *testing.T) {
		users, err := postgres.NewFinder[User](conn.Database()).
			Query("SELECT * FROM users;").
			Structs(ctx)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(users) != 2 {
			t.Fatalf("expected 2 user, got %d", len(users))
		}
	})

}
