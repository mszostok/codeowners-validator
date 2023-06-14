package check_test

import (
	"context"
	"testing"

	"go.szostok.io/codeowners/internal/api"
	"go.szostok.io/codeowners/internal/check"
	"go.szostok.io/codeowners/internal/config"

	"github.com/stretchr/testify/require"

	"go.szostok.io/codeowners/internal/ptr"

	"github.com/stretchr/testify/assert"
)

func TestValidOwnerChecker(t *testing.T) {
	tests := map[string]struct {
		owner   string
		isValid bool
	}{
		"Invalid Email": {
			owner:   `asda.comm`,
			isValid: false,
		},
		"Valid Email": {
			owner:   `gmail@gmail.com`,
			isValid: true,
		},
		"Invalid Team": {
			owner:   `@org/`,
			isValid: false,
		},
		"Valid Team": {
			owner:   `@org/user`,
			isValid: true,
		},
		"Invalid User": {
			owner:   `user`,
			isValid: false,
		},
		"Valid User": {
			owner:   `@user`,
			isValid: true,
		},
	}
	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// when
			result := check.IsValidOwner(tc.owner)
			assert.Equal(t, tc.isValid, result)
		})
	}
}

func TestValidOwnerCheckerIgnoredOwner(t *testing.T) {
	t.Run("Should ignore owner", func(t *testing.T) {
		// given
		ownerCheck, err := check.NewValidOwner(&config.Config{
			OwnerCheckerRepository:    "org/repo",
			OwnerCheckerIgnoredOwners: []string{"@owner1"},
		}, nil, true)
		require.NoError(t, err)

		givenCodeowners := `*	@owner1`

		// when
		out, err := ownerCheck.Check(context.Background(), LoadInput(givenCodeowners))

		// then
		require.NoError(t, err)
		assert.Empty(t, out.Issues)
	})

	t.Run("Should ignore user only and check the remaining owners", func(t *testing.T) {
		tests := map[string]struct {
			codeowners           string
			issue                *api.Issue
			allowUnownedPatterns bool
		}{
			"No owners": {
				codeowners: `*`,
				issue: &api.Issue{
					Severity: api.Warning,
					LineNo:   ptr.Uint64Ptr(1),
					Message:  "Missing owner, at least one owner is required",
				},
			},
			"Bad owner definition": {
				codeowners: `*	badOwner @owner1`,
				issue: &api.Issue{
					Severity: api.Error,
					LineNo:   ptr.Uint64Ptr(1),
					Message:  `Not valid owner definition "badOwner"`,
				},
			},
			"No owners but allow empty": {
				codeowners:           `*`,
				issue:                nil,
				allowUnownedPatterns: true,
			},
		}
		for tn, tc := range tests {
			t.Run(tn, func(t *testing.T) {
				// given
				ownerCheck, err := check.NewValidOwner(&config.Config{
					OwnerCheckerRepository:           "org/repo",
					OwnerCheckerAllowUnownedPatterns: tc.allowUnownedPatterns,
					OwnerCheckerIgnoredOwners:        []string{"@owner1"},
				}, nil, true)
				require.NoError(t, err)

				// when
				out, err := ownerCheck.Check(context.Background(), LoadInput(tc.codeowners))

				// then
				require.NoError(t, err)
				assertIssue(t, tc.issue, out.Issues)
			})
		}
	})
}

func TestValidOwnerCheckerOwnersMustBeTeams(t *testing.T) {
	tests := map[string]struct {
		codeowners           string
		issue                *api.Issue
		allowUnownedPatterns bool
	}{
		"Bad owner definition": {
			codeowners: `*	@owner1`,
			issue: &api.Issue{
				Severity: api.Error,
				LineNo:   ptr.Uint64Ptr(1),
				Message:  `Only team owners allowed and "@owner1" is not a team`,
			},
		},
		"No owners but allow empty": {
			codeowners:           `*`,
			issue:                nil,
			allowUnownedPatterns: true,
		},
	}
	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// given
			ownerCheck, err := check.NewValidOwner(&config.Config{
				OwnerCheckerRepository:           "org/repo",
				OwnerCheckerAllowUnownedPatterns: tc.allowUnownedPatterns,
				OwnerCheckerOwnersMustBeTeams:    true,
			}, nil, true)
			require.NoError(t, err)

			// when
			out, err := ownerCheck.Check(context.Background(), LoadInput(tc.codeowners))

			// then
			require.NoError(t, err)
			assertIssue(t, tc.issue, out.Issues)
		})
	}
}
