package gosql

import "context"

// Executable represents an entity that can execute SQL commands.
type Executable interface {
	// Exec executes a SQL command with the given arguments.
	Exec(ctx context.Context, sql string, arguments ...any) error
}

// Scanner represents an entity that can retrieving result from sql rows.
type Scanner interface {
	// Scan executes a SQL query and returns the result rows.
	Scan(ctx context.Context, sql string, arguments ...any) (Rows, error)
}

// Rows represents the result set of a SQL query.
type Rows interface {
	// Next prepares the next row for reading. Returns true if there is another row.
	Next() bool

	// Scan reads the current row's columns into dest.
	Scan(dest ...any) error

	// Close closes the Rows, preventing further enumeration.
	Close()
}
