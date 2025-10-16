package app

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// IgnorePattern represents a single gitignore pattern.
type IgnorePattern struct {
	Pattern      string // pattern to be matched
	IsNegative   bool   // true if pattern has ! at the beginning
	IsExact      bool   // pattern contains no wildcards
	IsNotOnlyDir bool   // pattern contains separators but not at the end (directory only)
}

// GitIgnore holds the parsed gitignore patterns.
type GitIgnore struct {
	patterns []IgnorePattern
	basePath string
}

// NewGitIgnoreFromFile creates a new GitIgnore instance from a .gitignore file by its path.
func NewGitIgnoreFromFile(filePath string) (*GitIgnore, error) {
	gi := &GitIgnore{}
	gi.AddPattern(".git/")
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return gi, nil
		}
		return nil, err
	}
	defer file.Close()

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}
	gi.basePath = filepath.Dir(absPath)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := purifyPattern(scanner.Text())
		if line != "" {
			gi.AddPattern(line)
		}
	}
	return gi, scanner.Err()
}

// AddPattern accepts a valid gitignore pattern. Pattern should not be an empty line, should not contain comments, and should contain only forward slashes as path separators
func (gi *GitIgnore) AddPattern(patStr string) {
	pat := IgnorePattern{}
	if strings.Contains(patStr, "!") {
		pat.IsNegative = true
		patStr = strings.TrimPrefix(patStr, "!")
	}
	pat.IsExact = !strings.ContainsAny(patStr, "*?[]")
	pat.IsNotOnlyDir = strings.Contains(patStr, "/") && !strings.HasSuffix(patStr, "/")
	pat.Pattern = patStr
	gi.patterns = append(gi.patterns, pat)
}

// Match determines if a given path should be ignored.
func (gi *GitIgnore) Match(filePath string) bool {
	shouldBeIgnored := false
	relPath, err := filepath.Rel(gi.basePath, filePath)
	if err != nil {
		return false
	}
	relPath = filepath.ToSlash(relPath)

	for _, pat := range gi.patterns {
		if !strings.Contains(pat.Pattern, "/") {
			if match, _ := filepath.Match(pat.Pattern, filepath.Base(relPath)); match {
				shouldBeIgnored = !pat.IsNegative
			}
			continue
		}

		if strings.HasSuffix(pat.Pattern, "/") {
			if info, err := os.Stat(filePath); err != nil || !info.IsDir() {
				continue
			}
		}

		if match, _ := filepath.Match(strings.TrimSuffix(pat.Pattern, "/"), relPath); match {
			shouldBeIgnored = !pat.IsNegative
		}
	}
	return shouldBeIgnored
}

// purifyPattern accepts a gitignore file line as an input string, trims any whitespace, removes comments and returns either a pure gitignore pattern, or an empty line otherwise
func purifyPattern(line string) string {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "#") {
		return ""
	}
	parts := strings.Split(line, "#")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}
