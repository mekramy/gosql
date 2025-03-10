package postgres_test

import (
	"fmt"
	"testing"

	"github.com/mekramy/gosql/postgres"
)

type User struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

func (u *User) Transform() error {
	u.Name = u.Name + fmt.Sprintf(" %d", u.Id)
	return nil
}

func TestCommander(t *testing.T) {
	e := NewMockExecutable("UPDATE `users` SET name = $1, family = $2 WHERE id = $3;")
	_, err := postgres.NewCmd(e).
		Command("UPDATE `users` SET name = ?, family = ? WHERE id = ?;").
		Exec(t.Context(), "John", "Doe", 3)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestInserter(t *testing.T) {
	u := User{Name: "John doe"}
	e := NewMockExecutable(`INSERT INTO users ("name") VALUES ($1);`)
	_, err := postgres.NewInserter[User](e).
		Table("users").
		Insert(t.Context(), u, postgres.OnlyFields("name"))
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestUpdater(t *testing.T) {
	u := User{Name: "John doe"}
	e := NewMockExecutable(`UPDATE users SET "name" = $1 WHERE id = $2;`)
	_, err := postgres.NewUpdater[User](e).
		Table("users").
		Where("id = ?", 2).
		Update(t.Context(), u, postgres.SkipFields("id"))
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCounter(t *testing.T) {
	r := NewMockReadable()
	count, err := postgres.NewCounter(r).
		Query("SELECT COUNT(*) FROM users;").
		Count(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 100 {
		t.Fatalf("expected %d, got %d", 100, count)
	}
}
