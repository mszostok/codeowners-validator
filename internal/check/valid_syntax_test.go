package check_test

import (
	"context"
	"testing"

	"go.szostok.io/codeowners/internal/api"
	"go.szostok.io/codeowners/internal/check"
	"go.szostok.io/codeowners/internal/ptr"
	"go.szostok.io/codeowners/pkg/codeowners"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidSyntaxChecker(t *testing.T) {
	tests := map[string]struct {
		codeowners string
		issue      *api.Issue
	}{
		"No owners": {
			codeowners: `*`,
			issue:      nil,
		},
		"Bad username": {
			codeowners: `pkg/github.com/** @-`,
			issue: &api.Issue{
				Severity: api.Warning,
				LineNo:   ptr.Uint64Ptr(1),
				Message:  "Owner '@-' does not look like a GitHub username or team name",
			},
		},
		"Bad org": {
			codeowners: `* @bad+org`,
			issue: &api.Issue{
				Severity: api.Warning,
				LineNo:   ptr.Uint64Ptr(1),
				Message:  "Owner '@bad+org' does not look like a GitHub username or team name",
			},
		},
		"Bad team name on first place": {
			codeowners: `* @org/+not+a+good+name`,
			issue: &api.Issue{
				Severity: api.Warning,
				LineNo:   ptr.Uint64Ptr(1),
				Message:  "Owner '@org/+not+a+good+name' does not look like a GitHub username or team name",
			},
		},
		"Bad team name on second place": {
			codeowners: `* @org/hakuna-matata @org/-a-team`,
			issue: &api.Issue{
				Severity: api.Warning,
				LineNo:   ptr.Uint64Ptr(1),
				Message:  "Owner '@org/-a-team' does not look like a GitHub username or team name",
			},
		},
		"Doesn't look like username, team name, nor email": {
			codeowners: `* something_weird`,
			issue: &api.Issue{
				Severity: api.Error,
				LineNo:   ptr.Uint64Ptr(1),
				Message:  "Owner 'something_weird' does not look like an email",
			},
		},
		"Comment in pattern line": {
			codeowners: `* @org/hakuna-matata # this is allowed`,
		},
	}
	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// when
			out, err := check.NewValidSyntax().
				Check(context.Background(), LoadInput(tc.codeowners))

			// then
			require.NoError(t, err)

			assertIssue(t, tc.issue, out.Issues)
		})
	}
}

func TestValidSyntaxZeroValueEntry(t *testing.T) {
	// given
	zeroValueInput := api.Input{
		CodeownersEntries: []codeowners.Entry{
			{
				LineNo:  0,
				Pattern: "",
				Owners:  nil,
			},
		},
	}
	expIssues := []api.Issue{
		{
			LineNo:   ptr.Uint64Ptr(0),
			Severity: api.Error,
			Message:  "Missing pattern",
		},
	}

	// when
	out, err := check.NewValidSyntax().
		Check(context.Background(), zeroValueInput)

	// then
	require.NoError(t, err)

	require.Len(t, out.Issues, len(expIssues))
	assert.EqualValues(t, expIssues, out.Issues)
}
