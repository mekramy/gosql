package query_test

import (
	"errors"
	"io/fs"
	"net/http"
	"testing"

	"github.com/mekramy/gosql/query"
)

type MockFS struct {
	files map[string]string
}

func (*MockFS) Exists(path string) (bool, error)                        { return false, nil }
func (*MockFS) Open(path string) (fs.File, error)                       { return nil, nil }
func (*MockFS) Search(dir, phrase, ignore, ext string) (*string, error) { return nil, nil }
func (*MockFS) Find(dir, pattern string) (*string, error)               { return nil, nil }
func (*MockFS) FS() fs.FS                                               { return nil }
func (*MockFS) Http() http.FileSystem                                   { return nil }
func (f *MockFS) Lookup(dir, pattern string) ([]string, error) {
	names := make([]string, 0)
	for k := range f.files {
		names = append(names, k)
	}
	return names, nil
}
func (f *MockFS) ReadFile(path string) ([]byte, error) {
	v, ok := f.files[path]
	if !ok {
		return nil, errors.New("file not found")
	}
	return []byte(v), nil
}

func TestFS(t *testing.T) {
	fs := &MockFS{
		files: map[string]string{
			"database/queries/user.sql": `
-- { query: list }
SELECT * FROM users WHERE deleted_at IS NULL;


-- { undefined: unsupported }
SELECT name, family
FROM users
WHERE
	deleted_at IS NULL
	AND age > 18
	AND name ILIKE '%@phrase%';

-- { query: single }
SELECT id, name, age FROM users WHERE @conditions;
			`,
		},
	}

	manager, err := query.NewQueryManager(fs, query.WithRoot("database/queries"))
	if err != nil {
		t.Fatal(err.Error())
	}

	if q := manager.Get("user/unsupported"); q != "" {
		t.Fatal("unsupported section must ignored")
	}

	if _, exists := manager.Find("not-exists-query"); exists {
		t.Fatal("not-exists-query should not exists!")
	}

	expected := `SELECT * FROM users WHERE deleted_at IS NULL;`
	if q := manager.Get("user/list"); q != expected {
		t.Fatalf(`expect "%s", got "%s"`, expected, q)
	}

	expected = `SELECT id, name, age FROM users WHERE deleted_at IS NULL AND (name = ? OR family = ?);`
	if q := manager.Query("user/single").And("deleted_at IS NULL").AndClosure("name = ? OR family = ?").Build(); q != expected {
		t.Fatalf(`expect "%s", got "%s"`, expected, q)
	}
}
