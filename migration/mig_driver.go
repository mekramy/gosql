package migration

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/mekramy/gofs"
)

type migration struct {
	root      string
	extention string
	dev       bool
	files     sortableFiles
	fs        gofs.FlexibleFS
	db        MigrationSource
	mutex     sync.RWMutex
}

func (m *migration) Load() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Locate files matching the specified extension
	files, err := m.fs.Lookup(m.root, ".*"+regexp.QuoteMeta(m.extention))
	if err != nil {
		return err
	}

	// Parse and cache migration stages from each file
	m.files = make(sortableFiles, 0)
	for _, file := range files {
		content, err := m.fs.ReadFile(file)
		if err != nil {
			return err
		}

		file := newMigrationFile(file, string(content))
		if file != nil {
			m.files = append(m.files, *file)
		}
	}

	sort.Sort(m.files)
	return nil
}

func (m *migration) Initialize() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return m.db.Exec(
		ctx,
		`CREATE TABLE IF NOT EXISTS migrations (
			name VARCHAR(100) NOT NULL,
			stage VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY(name, stage)
		);`,
	)
}

func (m *migration) Summary() (Summary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := m.db.Scan(ctx, `SELECT name, stage FROM migrations;`)
	if err != nil {
		return nil, err
	}

	result := make(Summary, 0)
	defer rows.Close()
	for rows.Next() {
		var name, stage string
		err := rows.Scan(&name, &stage)
		if err != nil {
			return nil, err
		}
		result = append(result, Migrated{
			Name:  name,
			Stage: stage,
		})
	}
	return result, nil
}

func (m *migration) StageSummary(stage string) (Summary, error) {
	summaries, err := m.Summary()
	if err != nil {
		return nil, err
	}

	return summaries.ForStage(stage), nil
}

func (m *migration) Up(stage string, options ...MigrationOption) ([]string, error) {
	// Hot reload on dev mode
	if m.dev {
		if err := m.Load(); err != nil {
			return nil, err
		}
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Create option
	option := &migrationOption{
		only:    make([]string, 0),
		exclude: make([]string, 0),
	}
	for _, opt := range options {
		opt(option)
	}

	// Skip migrated files
	migrated, err := m.StageSummary(stage)
	if err != nil {
		return nil, err
	}

	// Filter files
	files := m.files.
		Filter([]string{}, migrated.Names()).
		Filter(option.only, option.exclude)
	if files.Len() == 0 {
		return nil, nil
	}

	// Execute scripts
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	result := make([]string, 0)
	err = m.db.Transaction(ctx, func(tx ExecutableScanner) error {
		for _, file := range files {
			script, ok := file.UpScript(stage)
			if !ok || len(script) == 0 {
				continue
			}

			err := tx.Exec(ctx, script)
			if err != nil {
				return fmt.Errorf("%s: %w", file.name, err)
			}

			err = tx.Exec(
				ctx,
				fmt.Sprintf(`INSERT INTO migrations (name, stage) VALUES ('%s', '%s');`, file.name, stage),
			)
			if err != nil {
				return fmt.Errorf("%s: %w", file.name, err)
			}

			result = append(result, file.name)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *migration) Down(stage string, options ...MigrationOption) ([]string, error) {
	// Hot reload on dev mode
	if m.dev {
		if err := m.Load(); err != nil {
			return nil, err
		}
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Create option
	option := &migrationOption{
		only:    make([]string, 0),
		exclude: make([]string, 0),
	}
	for _, opt := range options {
		opt(option)
	}

	// Skip migrated files
	migrated, err := m.StageSummary(stage)
	if err != nil {
		return nil, err
	}

	// Filter files
	files := m.files.Reverse().
		Filter(migrated.Names(), []string{}).
		Filter(option.only, option.exclude)
	if files.Len() == 0 {
		return nil, nil
	}

	// Execute scripts
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	result := make([]string, 0)
	err = m.db.Transaction(ctx, func(tx ExecutableScanner) error {
		for _, file := range files {
			script, ok := file.DownScript(stage)
			if !ok || len(script) == 0 {
				continue
			}

			err := tx.Exec(ctx, script)
			if err != nil {
				return fmt.Errorf("%s: %w", file.name, err)
			}

			err = tx.Exec(
				ctx,
				fmt.Sprintf(`DELETE FROM migrations WHERE name = '%s' AND stage = '%s';`, file.name, stage),
			)
			if err != nil {
				return fmt.Errorf("%s: %w", file.name, err)
			}

			result = append(result, file.name)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
