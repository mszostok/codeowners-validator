package check_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.szostok.io/codeowners/internal/api"

	"go.szostok.io/codeowners/pkg/codeowners"
)

var FixtureValidCODEOWNERS = `
		# These owners will be the default owners for everything
		*       @global-owner1 @global-owner2

		# js owner
		*.js    @js-owner

		*.go docs@example.com

		/build/logs/ @doctocat

		/script m.t@g.com
`

func LoadInput(in string) api.Input {
	r := strings.NewReader(in)

	return api.Input{
		CodeownersEntries: codeowners.ParseCodeowners(r),
	}
}

func assertIssue(t *testing.T, expIssue *api.Issue, gotIssues []api.Issue) {
	t.Helper()

	if expIssue != nil {
		require.Len(t, gotIssues, 1)
		assert.EqualValues(t, *expIssue, gotIssues[0])
	} else {
		assert.Empty(t, gotIssues)
	}
}
