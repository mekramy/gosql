package migration

import (
	"github.com/mekramy/gofs"
)

// NewMigration initializes a migration with the specified database source, filesystem, and options.
func NewMigration(db MigrationSource, fs gofs.FlexibleFS, options ...Options) (Migration, error) {
	mig := &migration{
		root:  ".",
		ext:   "sql",
		dev:   false,
		files: make(sortableFiles, 0),
		fs:    fs,
		db:    db,
	}

	for _, opt := range options {
		opt(mig)
	}

	if err := mig.Load(); err != nil {
		return nil, err
	}

	if err := mig.Initialize(); err != nil {
		return nil, err
	}

	return mig, nil
}

// Migration defines the interface for managing database migrations.
// It includes methods for loading migration stages, initializing the migration table,
// retrieving summaries, and applying or rolling back migration stages.
type Migration interface {
	// Load loads migration stages from the filesystem and caches them.
	Load() error

	// Root returns the root directory of migration files.
	Root() string

	// Extension returns the file extension for migration files.
	Extension() string

	// IsDev indicates if it is in development mode.
	IsDev() bool

	// Initialize sets up the database migration table.
	Initialize() error

	// Summary returns an overview of the migration.
	Summary() (Summary, error)

	// Up applies migration stages.
	Up(stages []string, options ...MigrationOption) (Summary, error)

	// Down rolls back migration stages.
	Down(stages []string, options ...MigrationOption) (Summary, error)

	// Refresh rolls back and reapplies migration stages.
	Refresh(stages []string, options ...MigrationOption) (Summary, error)
}
