package codeowners_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/mszostok/codeowners-validator/pkg/codeowners"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sampleCodeownerFile = `
# Sample codeowner file
*	@everyone

src/**	@org/hakuna-matata @pico-bello
pkg/github.com/**	@myk

`

func TestParseCodeownersFailure(t *testing.T) {
	// given
	givenCodeownerPath := "workspace/go/repo-name"
	badInputs := []string{
		`# one token
*
`,
		`# bad username
* @myk
pkg/github.com/** @-
`,
		`# bad org
* @bad+org
`,
		`# comment mid line
* @org/hakuna-matata # this should be an error
`,
		`# bad team name second place
* @org/hakuna-matata @org/-a-team
`,
		`# bad team first
* @org/+not+a+good+name
`,
		`# doesn't look like username, team name, nor email
* something_weird
`,
	}

	for _, input := range badInputs {
		tFS := afero.NewMemMapFs()
		revert := codeowners.SetFS(tFS)
		defer revert()

		f, _ := tFS.Create(path.Join(givenCodeownerPath, "CODEOWNERS"))
		_, err := f.WriteString(input)
		require.NoError(t, err)

		// when
		_, err = codeowners.NewFromPath(givenCodeownerPath)

		// then
		require.Error(t, err)
	}
}

func TestParseCodeownersSuccess(t *testing.T) {
	// given
	givenCodeownerPath := "workspace/go/repo-name"
	expEntries := []codeowners.Entry{
		{
			LineNo:  3,
			Pattern: "*",
			Owners:  []string{"@everyone"},
		},
		{
			LineNo:  5,
			Pattern: "src/**",
			Owners:  []string{"@org/hakuna-matata", "@pico-bello"},
		},
		{
			LineNo:  6,
			Pattern: "pkg/github.com/**",
			Owners:  []string{"@myk"},
		},
	}

	tFS := afero.NewMemMapFs()
	revert := codeowners.SetFS(tFS)
	defer revert()

	f, _ := tFS.Create(path.Join(givenCodeownerPath, "CODEOWNERS"))
	_, err := f.WriteString(sampleCodeownerFile)
	require.NoError(t, err)

	// when
	entries, err := codeowners.NewFromPath(givenCodeownerPath)

	// then
	require.NoError(t, err)
	assert.Len(t, entries, len(expEntries))
	for _, expEntry := range expEntries {
		assert.Contains(t, entries, expEntry)
	}
}

func TestFindCodeownersFileSuccess(t *testing.T) {
	tests := map[string]struct {
		basePath string
	}{
		"Should find the CODEOWNERS at root": {
			basePath: "/workspace/go/repo-name1/",
		},
		"Should find the CODEOWNERS in docs/": {
			basePath: "/workspace/go/repo-name2/docs/",
		},
		"Should find the CODEOWNERS IN .github": {
			basePath: "/workspace/go/repo-name3/.github/",
		},
		"Should manage situation without trailing slash": {
			basePath: "/workspace/go/repo-name3/.github",
		},
	}
	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// given
			tFS := afero.NewMemMapFs()
			revert := codeowners.SetFS(tFS)
			defer revert()

			_, err := tFS.Create(path.Join(tc.basePath, "CODEOWNERS"))
			require.NoError(t, err)

			// when
			entry, err := codeowners.NewFromPath(tc.basePath)

			// then
			require.NoError(t, err)
			require.Empty(t, entry)
		})
	}
}

func TestFindCodeownersFileFailure(t *testing.T) {
	// given
	tFS := afero.NewMemMapFs()
	revert := codeowners.SetFS(tFS)
	defer revert()

	givenRepoPath := "/workspace/go/repo-without-codeowners/"
	expErrMsg := fmt.Sprintf("No CODEOWNERS found in the root, docs/, or .github/ directory of the repository %s", givenRepoPath)

	// when
	entries, err := codeowners.NewFromPath(givenRepoPath)

	// then
	assert.EqualError(t, err, expErrMsg)
	assert.Nil(t, entries)
}
