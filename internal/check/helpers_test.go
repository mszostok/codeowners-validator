package check_test

import (
	"strings"
	"testing"

	"github.com/mszostok/codeowners-validator/internal/check"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mszostok/codeowners-validator/pkg/codeowners"
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

func LoadInput(in string) check.Input {
	r := strings.NewReader(in)

	return check.Input{
		CodeownersEntries: codeowners.ParseCodeowners(r),
	}
}

func assertIssue(t *testing.T, expIssue *check.Issue, gotIssues []check.Issue) {
	t.Helper()

	if expIssue != nil {
		require.Len(t, gotIssues, 1)
		assert.EqualValues(t, *expIssue, gotIssues[0])
	} else {
		assert.Empty(t, gotIssues)
	}
}
