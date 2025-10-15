package app

import (
	"path/filepath"
	"slices"
)

type IgnoreProcessor struct {
	ignorePatterns []string
}

// NewIgnoreProcessor takes one or more ignore patterns and returns a new IgnoreProcessor instance
func NewIgnoreProcessor(ignorePatterns ...string) (proc *IgnoreProcessor) {
	proc = &IgnoreProcessor{}
	proc.Add(ignorePatterns...)
	return
}

// AddCommonIgnores adds common ignores such as git folder
func (proc *IgnoreProcessor) AddCommonIgnores() {
	proc.Add(".git")
}

// Add takes one or more ignore patterns and adds them to the IgnoreProcessor patterns list
func (proc *IgnoreProcessor) Add(patterns ...string) {
	for _, pat := range patterns {
		proc.ignorePatterns = append(proc.ignorePatterns, filepath.Base(pat))
	}
}

// IsPathShouldBeIgnored accepts an absolute or relative file or directory path and returns true if it should be ignored or false otherwise
func (proc *IgnoreProcessor) IsPathShouldBeIgnored(filePath string) bool {
	return slices.Contains(proc.ignorePatterns, filepath.Base(filePath))
}
