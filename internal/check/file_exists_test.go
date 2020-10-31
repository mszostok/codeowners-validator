package check_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mszostok/codeowners-validator/internal/check"
	"github.com/mszostok/codeowners-validator/internal/ptr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFileExists validates that file exists checker supports
// syntax used by CODEOWNERS. As the CODEOWNERS file uses a pattern that
// follows the same rules used in gitignore files, the test cases cover
// patterns from this document: https://git-scm.com/docs/gitignore#_pattern_format
func TestFileExists(t *testing.T) {
	tests := map[string]struct {
		codeownersInput string
		expectedIssues  []check.Issue
		paths           []string
	}{
		"Should found JS file": {
			codeownersInput: `
					*.js @pico
			`,
			paths: []string{
				"/somewhere/over/the/rainbow/here/it/is.js",
				"/somewhere/not/here/it/is.go",
			},
		},
		"Should match directory 'foo' anywhere": {
			codeownersInput: `
					**/foo @pico
			`,
			paths: []string{
				"/somewhere/over/the/foo/here/it/is.js",
			},
		},
		"Should match file 'foo' anywhere": {
			codeownersInput: `
					**/foo.js @pico
			`,
			paths: []string{
				"/somewhere/over/the/rainbow/here/it/foo.js",
			},
		},
		"Should match directory 'bar' anywhere that is directly under directory 'foo'": {
			codeownersInput: `
					**/foo/bar @bello
			`,
			paths: []string{
				"/somewhere/over/the/foo/bar/it/is.js",
			},
		},
		"Should match file 'bar' anywhere that is directly under directory 'foo'": {
			codeownersInput: `
					**/foo/bar.js @bello
			`,
			paths: []string{
				"/somewhere/over/the/foo/bar.js",
			},
		},
		"Should match all files inside directory 'abc'": {
			codeownersInput: `
					abc/** @bello
			`,
			paths: []string{
				"/abc/over/the/rainbow/bar.js",
			},
		},
		"Should match 'a/b', 'a/x/b', 'a/x/y/b' and so on": {
			codeownersInput: `
					a/**/b @bello
			`,
			paths: []string{
				"a/somewhere/over/the/b/foo.js",
			},
		},
		// https://github.community/t/codeowners-file-with-a-not-file-type-condition/1423
		"Should not match with negation pattern": {
			codeownersInput: `
					!/codeowners-validator @pico
			`,
			paths: []string{
				"/somewhere/over/the/rainbow/here/it/is.js",
			},
			expectedIssues: []check.Issue{
				newErrIssue(`"!/codeowners-validator" does not match any files in repository`),
			},
		},
		"Should not found JS file": {
			codeownersInput: `
					*.js @pico
			`,
			expectedIssues: []check.Issue{
				newErrIssue(`"*.js" does not match any files in repository`),
			},
		},
		"Should not match directory 'foo' anywhere": {
			codeownersInput: `
					**/foo @pico
			`,
			expectedIssues: []check.Issue{
				newErrIssue(`"**/foo" does not match any files in repository`),
			},
		},
		"Should not match file 'foo' anywhere": {
			codeownersInput: `
					**/foo.js @pico
			`,
			expectedIssues: []check.Issue{
				newErrIssue(`"**/foo.js" does not match any files in repository`),
			},
		},
		"Should no match directory 'bar' anywhere that is directly under directory 'foo'": {
			codeownersInput: `
					**/foo/bar @bello
			`,
			expectedIssues: []check.Issue{
				newErrIssue(`"**/foo/bar" does not match any files in repository`),
			},
		},
		"Should not match file 'bar' anywhere that is directly under directory 'foo'": {
			codeownersInput: `
					**/foo/bar.js @bello
			`,
			expectedIssues: []check.Issue{
				newErrIssue(`"**/foo/bar.js" does not match any files in repository`),
			},
		},
		"Should not match all files inside directory 'abc'": {
			codeownersInput: `
					abc/** @bello
			`,
			expectedIssues: []check.Issue{
				newErrIssue(`"abc/**" does not match any files in repository`),
			},
		},
		"Should not match 'a/**/b'": {
			codeownersInput: `
					a/**/b @bello
			`,
			expectedIssues: []check.Issue{
				newErrIssue(`"a/**/b" does not match any files in repository`),
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// given
			tmp, err := ioutil.TempDir("", "file-checker")
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, os.RemoveAll(tmp))
			}()

			initFSStructure(t, tmp, tc.paths)

			fchecker := check.NewFileExist()

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
			defer cancel()

			// when
			in := check.LoadInput(tc.codeownersInput)
			in.RepoDir = tmp
			out, err := fchecker.Check(ctx, in)

			// then
			require.NoError(t, err)
			assert.ElementsMatch(t, tc.expectedIssues, out.Issues)
		})
	}
}

func TestFileExistCheckFileSystemFailure(t *testing.T) {
	// given
	tmpdir, err := ioutil.TempDir("", "file-checker")
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, os.RemoveAll(tmpdir))
	}()

	err = os.MkdirAll(filepath.Join(tmpdir, "foo"), 0222)
	require.NoError(t, err)

	in := check.LoadInput("* @pico")
	in.RepoDir = tmpdir

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	// when
	out, err := check.NewFileExist().Check(ctx, in)

	// then
	require.Error(t, err)
	assert.Empty(t, out)
}

func newErrIssue(msg string) check.Issue {
	return check.Issue{
		Severity: check.Error,
		LineNo:   ptr.Uint64Ptr(2),
		Message:  msg,
	}
}
func initFSStructure(t *testing.T, base string, paths []string) {
	t.Helper()

	for _, p := range paths {
		if filepath.Ext(p) == "" {
			err := os.MkdirAll(filepath.Join(base, p), 0755)
			require.NoError(t, err)
		} else {
			dir := filepath.Dir(p)

			err := os.MkdirAll(filepath.Join(base, dir), 0755)
			require.NoError(t, err)

			err = ioutil.WriteFile(filepath.Join(base, p), []byte("hakuna-matata"), 0600)
			require.NoError(t, err)
		}
	}
}
