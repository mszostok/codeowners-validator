package codeowners_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.szostok.io/codeowners/pkg/codeowners"
)

const sampleCodeownerFile = `
# Sample codeowner file
*	@everyone

src/**	@org/hakuna-matata @pico-bello
pkg/github.com/**	@myk
tests/**	@ghost # some comment
internal/**	@ghost #some comment v2

`

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
		{
			LineNo:  7,
			Pattern: "tests/**",
			Owners:  []string{"@ghost"},
		},
		{
			LineNo:  8,
			Pattern: "internal/**",
			Owners:  []string{"@ghost"},
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

func TestMultipleCodeownersFileFailure(t *testing.T) {
	const givenRepoPath = "/workspace/go/repo-without-codeowners/"
	tests := map[string]struct {
		expErrMsg                string
		givenCodeownersLocations []string
	}{
		"Should report that no CODEOWNERS file was found": {
			expErrMsg:                fmt.Sprintf("No CODEOWNERS found in the root, docs/, or .github/ directory of the repository %s", givenRepoPath),
			givenCodeownersLocations: nil,
		},
		"Should report that CODEOWNERS file was found on root and docs/": {
			expErrMsg:                fmt.Sprintf("Multiple CODEOWNERS files found in the ./CODEOWNERS and ./docs/CODEOWNERS locations of the repository %s", givenRepoPath),
			givenCodeownersLocations: []string{"CODEOWNERS", path.Join("docs", "CODEOWNERS")},
		},
		"Should report that CODEOWNERS file was found on root and .github/": {
			expErrMsg:                fmt.Sprintf("Multiple CODEOWNERS files found in the ./CODEOWNERS and ./.github/CODEOWNERS locations of the repository %s", givenRepoPath),
			givenCodeownersLocations: []string{"CODEOWNERS", path.Join(".github/", "CODEOWNERS")},
		},
		"Should report that CODEOWNERS file was found in docs/ and .github/": {
			expErrMsg:                fmt.Sprintf("Multiple CODEOWNERS files found in the ./docs/CODEOWNERS and ./.github/CODEOWNERS locations of the repository %s", givenRepoPath),
			givenCodeownersLocations: []string{path.Join(".github", "CODEOWNERS"), path.Join("docs", "CODEOWNERS")},
		},
		"Should report that CODEOWNERS file was found on root, docs/ and .github/": {
			expErrMsg:                fmt.Sprintf("Multiple CODEOWNERS files found in the ./CODEOWNERS, ./docs/CODEOWNERS, and ./.github/CODEOWNERS locations of the repository %s", givenRepoPath),
			givenCodeownersLocations: []string{"CODEOWNERS", path.Join(".github", "CODEOWNERS"), path.Join("docs", "CODEOWNERS")},
		},
	}
	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// given
			tFS := afero.NewMemMapFs()
			revert := codeowners.SetFS(tFS)
			defer revert()

			for _, location := range tc.givenCodeownersLocations {
				_, err := tFS.Create(path.Join(givenRepoPath, location))
				require.NoError(t, err)
			}

			// when
			entries, err := codeowners.NewFromPath(givenRepoPath)

			// then
			assert.EqualError(t, err, tc.expErrMsg)
			assert.Nil(t, entries)
		})
	}
}
