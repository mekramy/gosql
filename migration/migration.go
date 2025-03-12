package migration

import (
	"github.com/mekramy/gofs"
)

// NewMigration initializes a migration with the specified database source, filesystem, and options.
func NewMigration(db MigrationSource, fs gofs.FlexibleFS, options ...Options) (Migration, error) {
	mig := &migration{
		root:      ".",
		extention: ".sql",
		dev:       false,
		files:     make(sortableFiles, 0),
		fs:        fs,
		db:        db,
	}

	for _, opt := range options {
		opt(mig)
	}

	if err := mig.Load(); err != nil {
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

	// Initialize sets up the database migration table.
	Initialize() error

	// Summary returns an overview of the migration.
	Summary() (Summary, error)

	// StageSummary returns an overview of a specific migration stage.
	StageSummary(stage string) (Summary, error)

	// Up applies a migration stages.
	Up(options ...MigrationOption) ([]MigrationResult, error)

	// Down rolls back migration stages.
	Down(options ...MigrationOption) ([]MigrationResult, error)
}

type MigrationResult struct {
	Stage string
	Name  string
}
