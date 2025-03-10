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

// parseQueries extracts all named queries from the given SQL file content.
func parseQueries(content string) (map[string]string, error) {
	var name, query string
	res := make(map[string]string)

	// Parse query
	rx, err := regexp.Compile(`\s*--\s*\{\s*(query)\s*:\s*([\w\s]+)\s*\}`)
	parseQuery := func(content string) (string, bool) {
		matches := rx.FindStringSubmatch(content)
		if len(matches) == 3 && strings.TrimSpace(matches[2]) != "" {
			return strings.TrimSpace(matches[2]), true
		}
		return "", false
	}
	if err != nil {
		return nil, err
	}

	// Scan lines
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if q, ok := parseQuery(line); ok {
			if name != "" {
				res[name] = strings.TrimRight(query, "\n")
			}
			name = q
			query = ""
		} else if strings.TrimSpace(line) != "" && name != "" {
			query = query + line + "\n"
		}
	}

	if name != "" && query != "" {
		res[name] = strings.TrimRight(query, "\n")
	}

	return res, nil
}
