package gosql

import "context"

// Executable represents an entity that can execute SQL commands.
type Executable interface {
	// Exec executes a SQL command with the given arguments.
	Exec(ctx context.Context, sql string, arguments ...any) error
}

// Scanner represents an entity that can scan through SQL query results.
type Scanner interface {
	// Next prepares the next row for reading. It returns true if there is another row and false if no more rows are available
	Next() bool

	// Scan reads the values from the current row into dest values positionally.
	Scan(dest ...any) error

	// Close closes the scanner.
	Close()
}
