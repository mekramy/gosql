package migration

import (
	"strings"

	"github.com/mekramy/goutils"
)

type Options func(*migration)

// WithRoot sets the root directory for migration files.
func WithRoot(root string) Options {
	root = goutils.NormalizePath(root)
	return func(q *migration) {
		if root != "" {
			q.root = root
		} else {
			q.root = "."
		}
	}
}

// WithExtension sets the file extension for migration files, adding a leading dot if missing.
func WithExtension(ext string) Options {
	ext = strings.TrimSpace(ext)
	if ext != "" && !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	return func(q *migration) {
		if ext != "" {
			q.extention = ext
		}
	}
}

// WithEnv sets the environment mode for migrations.
// Enables development mode if true, causing Load() to be called on each migration run.
func WithEnv(isDev bool) Options {
	return func(q *migration) {
		q.dev = isDev
	}
}

type MigrationOption func(*migrationOption)

// WithStage specifies the stages to include in the migration.
func WithStage(stages ...string) MigrationOption {
	return func(o *migrationOption) {
		o.stages = append(o.stages, stages...)
	}
}

// OnlyFiles specifies the files to include in the migration.
func OnlyFiles(files ...string) MigrationOption {
	return func(o *migrationOption) {
		o.only = append(o.only, files...)
	}
}

// SkipFiles specifies the files to exclude from the migration.
func SkipFiles(files ...string) MigrationOption {
	return func(o *migrationOption) {
		o.exclude = append(o.exclude, files...)
	}
}

type migrationOption struct {
	stages  []string
	only    []string
	exclude []string
}
