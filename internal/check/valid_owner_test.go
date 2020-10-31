package check

import (
	"context"
	"testing"

	"github.com/mszostok/codeowners-validator/internal/ptr"
	"github.com/stretchr/testify/require"

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
			result := isValidOwner(tc.owner)
			assert.Equal(t, tc.isValid, result)
		})
	}
}

func TestValidOwnerCheckerIgnoredOwner(t *testing.T) {
	t.Run("Should ignore owner", func(t *testing.T) {
		// given
		ownerCheck, err := NewValidOwner(ValidOwnerConfig{
			Repository:    "org/repo",
			IgnoredOwners: []string{"@owner1"},
		}, nil)
		require.NoError(t, err)

		givenCodeowners := `*	@owner1`

		// when
		out, err := ownerCheck.Check(context.Background(), LoadInput(givenCodeowners))

		// then
		require.NoError(t, err)
		assert.Empty(t, out.Issues)
	})

	t.Run("Should ignore user only and check the remaining owners", func(t *testing.T) {
		// given
		ownerCheck, err := NewValidOwner(ValidOwnerConfig{
			Repository:    "org/repo",
			IgnoredOwners: []string{"@owner1"},
		}, nil)
		require.NoError(t, err)

		givenCodeowners := `*	@owner1 badOwner`

		expIssue := Issue{
			Severity: Error,
			LineNo:   ptr.Uint64Ptr(1),
			Message:  `Not valid owner definition "badOwner"`,
		}

		// when
		out, err := ownerCheck.Check(context.Background(), LoadInput(givenCodeowners))

		// then
		require.NoError(t, err)
		require.Len(t, out.Issues, 1)
		require.EqualValues(t, expIssue, out.Issues[0])
	})
}

func isValidOwner(owner string) bool {
	return isEmailAddress(owner) || isGithubUser(owner) || isGithubTeam(owner)
}
