package query

import (
	"sync"

	"github.com/mekramy/gofs"
)

// QueryManager defines methods to read SQL queries from the filesystem.
// Queries must be defined by `-- query: name-of-query` comment in SQL files.
type QueryManager interface {
	// Load loads the queries from the storage.
	Load() error

	// Get retrieves the SQL query by name.
	// Queries are resolved by `relative/path/to/root/query` key.
	Get(name string) string

	// Find retrieves the SQL query by name and indicates if it was found.
	// Queries are resolved by `relative/path/to/root/query` key.
	Find(name string) (string, bool)
}

// NewQueryManager creates a new query manager with the provided filesystem and options.
func NewQueryManager(fs gofs.FlexibleFS, options ...Options) (QueryManager, error) {
	q := &query{
		root:    ".",
		ext:     ".sql",
		dev:     false,
		fs:      fs,
		queries: make(map[string]string),
	}

	for _, opt := range options {
		opt(q)
	}

	if err := q.Load(); err != nil {
		return nil, err
	}

	return q, nil
}

type query struct {
	root    string
	ext     string
	dev     bool
	fs      gofs.FlexibleFS
	queries map[string]string
	mutex   sync.RWMutex
}

func (q *query) Load() error {
	// Safe race condition
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Read files from fs
	files, err := q.fs.Lookup(
		q.root,
		extPattern("", q.ext),
	)
	if err != nil {
		return err
	}

	// Parse queries
	for _, file := range files {
		// Read file
		content, err := q.fs.ReadFile(file)
		if err != nil {
			return err
		}

		// Parse queries
		fName := toName(file, q.root, q.ext)
		queries, err := parseQueries(string(content))
		if err != nil {
			return err
		}

		for qName, query := range queries {
			q.queries[fName+"/"+qName] = query
		}
	}

	return nil
}

func (q *query) Get(name string) string {
	if q.dev {
		q.Load()
	}

	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return q.queries[name]
}

func (q *query) Find(name string) (string, bool) {
	if q.dev {
		q.Load()
	}

	q.mutex.RLock()
	defer q.mutex.RUnlock()
	v, ok := q.queries[name]
	return v, ok
}
