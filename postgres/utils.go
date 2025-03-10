package postgres

// parseVariadic returns the first value from the variadic parameter `vals` if it exists,
// otherwise it returns the default value `def`.
func parseVariadic[T any](def T, vals ...T) T {
	if len(vals) > 0 {
		return vals[0]
	}
	return def
}
