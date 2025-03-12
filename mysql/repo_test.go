package mysql_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/mekramy/gosql/mysql"
)

func TestRepository(t *testing.T) {
	type User struct {
		Id   int    `db:"id"`
		Name string `db:"name"`
	}
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

	t.Run("Init", func(t *testing.T) {
		err = conn.Transaction(ctx, func(tx *sql.Tx) error {
			_, err := mysql.NewCmd(tx).
				Command("DROP TABLE IF EXISTS users CASCADE;").
				Exec(ctx)
			if err != nil {
				return err
			}

			_, err = mysql.NewCmd(tx).
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
		err = conn.Transaction(ctx, func(tx *sql.Tx) error {
			for idx, name := range []string{"John Doe", "Jack Ma"} {
				u := User{
					Id:   idx + 1,
					Name: name,
				}

				_, err = mysql.NewInserter[User](tx).
					Table("users").
					Insert(ctx, u, mysql.OnlyFields("name"))
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
		err = conn.Transaction(ctx, func(tx *sql.Tx) error {
			for idx, name := range []string{"John Doe New", "Jack Ma New"} {
				u := User{
					Id:   idx + 1,
					Name: name,
				}

				_, err = mysql.NewUpdater[User](tx).
					Table("users").
					Where("id = ?", u.Id).
					Update(ctx, u, mysql.SkipFields("id"))
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
		count, err := mysql.NewCounter(conn.Database()).
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
		jack, err := mysql.NewFinder[User](conn.Database()).
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
		users, err := mysql.NewFinder[User](conn.Database()).
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
