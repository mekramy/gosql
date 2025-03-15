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

// WithExtension sets the file extension for migration files.
func WithExtension(ext string) Options {
	ext = strings.TrimSpace(ext)
	ext = strings.TrimLeft(ext, ".")
	return func(q *migration) {
		if ext != "" {
			q.ext = ext
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

// OnlyFiles specifies the files to include in the migration.
func OnlyFiles(files ...string) MigrationOption {
	return func(o *migrationOption) {
		o.only.Add(files...)
	}
}

// SkipFiles specifies the files to exclude from the migration.
func SkipFiles(files ...string) MigrationOption {
	return func(o *migrationOption) {
		o.exclude.Add(files...)
	}
}

func newOption() *migrationOption {
	return &migrationOption{
		only:    &optionSet{elements: make(map[string]struct{})},
		exclude: &optionSet{elements: make(map[string]struct{})},
	}
}

type migrationOption struct {
	only    *optionSet
	exclude *optionSet
}

type optionSet struct {
	elements map[string]struct{}
}

func (s *optionSet) Add(elements ...string) {
	for _, element := range elements {
		s.elements[element] = struct{}{}
	}
}

func (s *optionSet) Size() int {
	return len(s.elements)
}

func (s *optionSet) Elements() []string {
	keys := make([]string, 0, len(s.elements))
	for key := range s.elements {
		keys = append(keys, key)
	}
	return keys
}
