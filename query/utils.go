package query

import (
	"bufio"
	"path/filepath"
	"regexp"
	"strings"
)

// parseVariadic returns the first value from the variadic parameter `vals` if it exists,
// otherwise it returns the default value `def`.
func parseVariadic[T any](def T, vals ...T) T {
	if len(vals) > 0 {
		return vals[0]
	}
	return def
}

// normalizePath normalizes the given path segments
// by joining them and converting to a slashed separator.
func normalizePath(path ...string) string {
	return filepath.ToSlash(filepath.Clean(filepath.Join(path...)))
}

// extPattern creates a regular expression pattern
// to match paths with the specified extension.
func extPattern(path, ext string) string {
	if path == "" {
		return ".*" + regexp.QuoteMeta(ext)
	}

	return "^" + regexp.QuoteMeta(path) + ".*" + regexp.QuoteMeta(ext)
}

// toName converts a file path to
// a name by removing the root and extension.
func toName(path, root, ext string) string {
	if path == "" {
		return ""
	}

	path = strings.TrimPrefix(path, root)
	path = strings.TrimSuffix(path, ext)
	return normalizePath(path)
}

// parseQuery extracts the query name from the given
// content if it follows the format "--query:name".
func parseQuery(content string) (string, bool) {
	content = strings.ReplaceAll(content, " ", "")
	content = strings.ReplaceAll(content, "\t", "")
	parts := strings.Split(content, ":")
	if len(parts) != 2 || parts[0] != "--query" || parts[1] == "" {
		return "", false
	}

	return parts[1], true
}

// parseQueries extracts all named queries from the given SQL file content.
func parseQueries(content string) (map[string]string, error) {
	res := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Split(bufio.ScanLines)

	name := ""
	query := ""
	for scanner.Scan() {
		line := scanner.Text()
		if q, isNew := parseQuery(line); isNew {
			if name != "" {
				res[name] = strings.TrimRight(query, "\n")
			}
			name = q
			query = ""
		} else if strings.TrimSpace(line) != "" {
			query = query + line + "\n"
		}
	}

	if name != "" && query != "" {
		res[name] = strings.TrimRight(query, "\n")
	}

	return res, nil
}
