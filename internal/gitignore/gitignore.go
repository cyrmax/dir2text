package gitignore

// GitIgnore holds an ordered collection of parsed patterns and a base path.
type GitIgnore struct {
	// TODO: to be implemented
}

// New creates a new GitIgnore with no patterns.
func NewGitIgnore(basePath string) *GitIgnore {
	// TODO: implementation
	return &GitIgnore{}
}

// Match returns true if a given path matches against a collection of added patterns and false otherwise
func (gi *GitIgnore) Match(path string) bool {
	// TODO: implementation
	return false
}

// AppendPattern adds a new pattern to a collection of patterns
func (gi *GitIgnore) AppendPattern(pattern string) {
	// TODO: implementation

}
