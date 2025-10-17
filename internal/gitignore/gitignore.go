package gitignore

// ignorePattern holds a single parsed pattern.
type ignorePattern struct {
	Pattern string
	Negates bool
}

// GitIgnore holds an ordered collection of parsed patterns and a base path.
type GitIgnore struct {
	basePath string
	patterns []ignorePattern
}

// New creates a new GitIgnore with no patterns. The base path is not checked to be valid path.
func NewGitIgnore(basePath string) *GitIgnore {
	return &GitIgnore{basePath: basePath}
}

// Match returns true if a given path matches against a collection of added patterns and false otherwise
func (gi *GitIgnore) Match(path string) bool {
	// FIXME: Not fully implemented yet
	shouldIgnore := false
	for _, pat := range gi.patterns {
		if pat.Pattern == path {
			shouldIgnore = true
		}
	}
	return shouldIgnore
}

// AppendPattern adds a new pattern to a collection of patterns
func (gi *GitIgnore) AppendPattern(pattern string) {
	// FIXME: Not fully implemented yet
	gi.patterns = append(gi.patterns, ignorePattern{Pattern: pattern, Negates: false})
}
