package gitignore

import "testing"

func TestGitIgnore_Match(t *testing.T) {
	testCases := []struct {
		name     string
		patterns []string
		path     string
		want     bool
	}{
		// --- Basic Matching ---
		{
			name:     "Match a simple file name",
			patterns: []string{"hello.txt"},
			path:     "hello.txt",
			want:     true,
		},
		{
			name:     "Do not match a different file name",
			patterns: []string{"hello.txt"},
			path:     "world.txt",
			want:     false,
		},
		{
			name:     "Simple pattern matches file in subdirectory",
			patterns: []string{"hello.txt"},
			path:     "a/b/hello.txt",
			want:     true,
		},

		// --- Blank Lines and Comments (#) ---
		{
			name:     "Blank lines do not match anything",
			patterns: []string{"", "hello.txt"},
			path:     "",
			want:     false,
		},
		{
			name:     "Lines with only whitespace do not match anything",
			patterns: []string{"   ", "hello.txt"},
			path:     "   ",
			want:     false,
		},
		{
			name:     "Comment lines are ignored",
			patterns: []string{"# this is a comment", "a.txt"},
			path:     "a.txt",
			want:     true,
		},
		{
			name:     "Escaped hash symbol matches a file starting with #",
			patterns: []string{`\#hello.txt`},
			path:     "#hello.txt",
			want:     true,
		},
		{
			name:     "Hash symbol in the middle of a pattern is literal",
			patterns: []string{`hello#world.txt`},
			path:     "hello#world.txt",
			want:     true,
		},

		// --- Trailing Spaces and Escaping ---
		{
			name:     "Trailing spaces are ignored",
			patterns: []string{"hello.txt   "},
			path:     "hello.txt",
			want:     true,
		},
		{
			name:     "Escaped trailing space is not ignored",
			patterns: []string{`hello.txt\ `},
			path:     "hello.txt ",
			want:     true,
		},
		{
			name:     "Escaped trailing spaces (multiple) are not ignored",
			patterns: []string{`config\ \ `},
			path:     "config  ",
			want:     true,
		},
		{
			name:     "Escaped space in the middle of a pattern is preserved",
			patterns: []string{`my\ file.txt`},
			path:     "my file.txt",
			want:     true,
		},

		// --- Negation (!) ---
		{
			name:     "Negation pattern re-includes a previously excluded file",
			patterns: []string{"*.log", "!important.log"},
			path:     "important.log",
			want:     false,
		},
		{
			name:     "Negation has no effect if the file was not previously excluded",
			patterns: []string{"*.txt", "!important.log"},
			path:     "important.log",
			want:     false,
		},
		{
			name:     "Last matching pattern wins: a file is ignored if it matches a later pattern",
			patterns: []string{"!important.log", "*.log"},
			path:     "important.log",
			want:     true,
		},
		{
			name:     "Negation of a file inside an excluded directory has no effect",
			patterns: []string{"logs/", "!logs/important.log"},
			path:     "logs/important.log",
			want:     true,
		},
		{
			name:     "Escaped negation symbol matches a file starting with !",
			patterns: []string{`\!important.file`},
			path:     "!important.file",
			want:     true,
		},

		// --- Path Separators (/) ---
		{
			name:     "Trailing slash matches a directory and everything inside it",
			patterns: []string{"logs/"},
			path:     "logs/today.log",
			want:     true,
		},
		{
			name:     "Trailing slash matches a directory itself",
			patterns: []string{"logs/"},
			path:     "logs",
			want:     true,
		},
		{
			name:     "Trailing slash does not match a file with the same name",
			patterns: []string{"build/"},
			path:     "build",
			// A file named 'build' should not be ignored by a pattern 'build/'.
			// We check this with `want: false` if 'build' were a file.
			want: false, // This assumes path 'build' refers to a file.
		},
		{
			name:     "Leading slash anchors the pattern to the root directory",
			patterns: []string{"/root.log"},
			path:     "root.log",
			want:     true,
		},
		{
			name:     "Leading slash prevents matching in subdirectories",
			patterns: []string{"/logs/today.log"},
			path:     "sub/logs/today.log",
			want:     false,
		},
		{
			name:     "Slash in the middle anchors the pattern to the root directory",
			patterns: []string{"config/routes.rb"},
			path:     "config/routes.rb",
			want:     true,
		},
		{
			name:     "Slash in the middle prevents matching in subdirectories",
			patterns: []string{"config/routes.rb"},
			path:     "staging/config/routes.rb",
			want:     false,
		},
		{
			name:     "Example from docs: 'doc/frotz/' matches a nested directory",
			patterns: []string{"doc/frotz/"},
			path:     "doc/frotz/file.txt",
			want:     true,
		},
		{
			name:     "Example from docs: 'doc/frotz/' does not match at other levels",
			patterns: []string{"doc/frotz/"},
			path:     "a/doc/frotz/file.txt",
			want:     false,
		},

		// --- Wildcards (*, ?, []) ---
		{
			name:     "Asterisk matches all files with an extension at any level",
			patterns: []string{"*.log"},
			path:     "some/path/to/a/file.log",
			want:     true,
		},
		{
			name:     "Asterisk matches files in the root directory",
			patterns: []string{"*.log"},
			path:     "file.log",
			want:     true,
		},
		{
			name:     "Asterisk does not match directory separators",
			patterns: []string{"foo*bar"},
			path:     "foo/baz/bar",
			want:     false,
		},
		{
			name:     "Question mark matches any single character except '/'",
			patterns: []string{"file?.log"},
			path:     "file1.log",
			want:     true,
		},
		{
			name:     "Question mark does not match zero or multiple characters",
			patterns: []string{"file?.log"},
			path:     "file12.log",
			want:     false,
		},
		{
			name:     "Range notation [a-z] matches a character within the range",
			patterns: []string{"log[0-9].txt"},
			path:     "log5.txt",
			want:     true,
		},
		{
			name:     "Range notation does not match a character outside the range",
			patterns: []string{"log[0-9].txt"},
			path:     "logA.txt",
			want:     false,
		},
		{
			name:     "Negated range notation [!a-z] matches a character outside the range",
			patterns: []string{"log[!0-9].txt"},
			path:     "logA.txt",
			want:     true,
		},

		// --- Double Asterisk (**) ---
		{
			name:     "Leading '**' matches in all directories (same as no prefix)",
			patterns: []string{"**/foo"},
			path:     "a/b/foo",
			want:     true,
		},
		{
			name:     "Leading '**' followed by a slash matches file in any directory",
			patterns: []string{"**/logs/debug.log"},
			path:     "app/component/logs/debug.log",
			want:     true,
		},
		{
			name:     "Trailing '/**' matches everything inside a directory",
			patterns: []string{"logs/**"},
			path:     "logs/a/b/c.log",
			want:     true,
		},
		{
			name:     "Trailing '/**' does not match files outside the directory",
			patterns: []string{"logs/**"},
			path:     "other/logs/a.log",
			want:     false,
		},
		{
			name:     "Slash-flanked '**' matches zero directories",
			patterns: []string{"a/**/b"},
			path:     "a/b",
			want:     true,
		},
		{
			name:     "Slash-flanked '**' matches one directory",
			patterns: []string{"a/**/b"},
			path:     "a/x/b",
			want:     true,
		},
		{
			name:     "Slash-flanked '**' matches multiple directories",
			patterns: []string{"a/**/b"},
			path:     "a/x/y/z/b",
			want:     true,
		},
		{
			name:     "Other consecutive asterisks are treated as single asterisks",
			patterns: []string{"a/***/b"},
			path:     "a/xyz/b",
			want:     true,
		},
		{
			name:     "Other consecutive asterisks do not cross directory boundaries",
			patterns: []string{"a/***/b"},
			path:     "a/x/y/b",
			want:     false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			gi := NewGitIgnore(".")
			for _, pattern := range testCase.patterns {
				gi.AppendPattern(pattern)
			}

			got := gi.Match(testCase.path)

			if got != testCase.want {
				t.Errorf("GitIgnore.Match(%q) with patterns %v = %v, want %v", testCase.path, testCase.patterns, got, testCase.want)
			}
		})
	}
}
