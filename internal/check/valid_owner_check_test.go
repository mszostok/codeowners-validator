package check_test

import (
	"testing"

	"github.com/mszostok/codeowners-validator/internal/check"

	"github.com/stretchr/testify/assert"
)

func TestValidOwnerChecker(t *testing.T) {
	tests := map[string]struct {
		user    string
		isValid bool
	}{
		"Invalid Email": {
			user:    `asda.comm`,
			isValid: false,
		},
		"Valid Email": {
			user:    `gmail@gmail.com`,
			isValid: true,
		},
		"Invalid Team": {
			user:    `@org/`,
			isValid: false,
		},
		"Valid Team": {
			user:    `@org/user`,
			isValid: true,
		},
		"Invalid User": {
			user:    `user`,
			isValid: false,
		},
		"Valid User": {
			user:    `@user`,
			isValid: true,
		},
	}
	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// when
			result := isValidUser(tc.user)
			assert.Equal(t, tc.isValid, result)
		})
	}
}

func isValidUser(user string) bool {
	return check.IsEmailAddress(user) || check.IsGithubUser(user) || check.IsGithubTeam(user)
}
