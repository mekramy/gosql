package query

import "strings"

type Options func(*query)

// WithRoot set root directory for query files in FS.
// The default root is "."
func WithRoot(root string) Options {
	root = normalizePath(root)
	return func(q *query) {
		if root != "" {
			q.root = root + "/"
		} else {
			q.root = "."
		}
	}
}

// WithExtension sets the file extension for the query files.
// The default extension is ".sql".
func WithExtension(ext string) Options {
	ext = strings.TrimSpace(ext)
	if ext != "" && !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return func(q *query) {
		if ext != "" {
			q.ext = ext
		}
	}
}

// WithEnv sets the environment to dev or production.
// In dev mode Load() called on each compile.
// CAUTION: disable development mode on production.
func WithEnv(isDev bool) Options {
	return func(q *query) {
		q.dev = isDev
	}
}
