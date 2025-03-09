package query_test

import (
	"testing"

	"github.com/mekramy/gosql/query"
)

func TestConditionBuilder_SQL(t *testing.T) {
	cond := query.NewCondition(query.NumbericResolver)
	cond.And("name = ?", "John").
		AndClosure("age > ? AND age < ?", 9, 31).
		OrIf(false, "age IS NULL").
		OrClosureIf(true, "membership @in", "admin", "manager", "accountant")

	expected := "name = $1 AND (age > $2 AND age < $3) OR (membership IN ($4, $5, $6))"
	if sql := cond.SQL(); sql != expected {
		t.Errorf("Expect %s, got %s", expected, sql)
	}
}

func TestConditionBuilder_Build(t *testing.T) {
	cond := query.NewCondition()
	cond.And("name = ?", "John").
		AndClosure("age > ? AND age < ?", 9, 31).
		OrIf(false, "age IS NULL").
		OrClosureIf(true, "membership @in", "admin", "manager", "accountant")

	expected := "SELECT COUNT(*) FROM `users` WHERE name = ? AND (age > ? AND age < ?) OR (membership IN (?, ?, ?));"
	result := cond.Build("SELECT COUNT(*) FROM `users` @where;")
	if result != expected {
		t.Errorf("Expect %s, got %s", expected, result)
	}
}
