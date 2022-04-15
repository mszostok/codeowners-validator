package check_test

import (
	"context"
	"testing"

	"github.com/mszostok/codeowners-validator/internal/check"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAvoidShadowing(t *testing.T) {
	getuint64 := func(i int) *uint64 {
		u := uint64(i)
		return &u
	}

	tests := map[string]struct {
		codeownersInput string
		expectedIssues  []check.Issue
	}{
		"Should report info about shadowed entries": {
			codeownersInput: `
					/build/logs/ @doctocat
					/script      @mszostok

					# Shadows
					*            @s1
					/s*/         @s2
					/s*          @s3
					/b*          @s4
					/b*/logs     @s5

					# OK
					/b*/other    @o1
					/script/*	 @o2
			`,
			expectedIssues: []check.Issue{
				{
					Severity: check.Error,
					LineNo:   getuint64(6),
					Message: `Pattern "*" shadows the following patterns:
            * 2: "/build/logs/"
            * 3: "/script"
Entries should go from least-specific to most-specific.`,
				},
				{
					Severity: check.Error,
					LineNo:   getuint64(7),
					Message: `Pattern "/s*/" shadows the following patterns:
            * 3: "/script"
Entries should go from least-specific to most-specific.`,
				},
				{
					Severity: check.Error,
					LineNo:   getuint64(8),
					Message: `Pattern "/s*" shadows the following patterns:
            * 3: "/script"
            * 7: "/s*/"
Entries should go from least-specific to most-specific.`,
				},
				{
					Severity: check.Error,
					LineNo:   getuint64(9),
					Message: `Pattern "/b*" shadows the following patterns:
            * 2: "/build/logs/"
Entries should go from least-specific to most-specific.`,
				},
				{
					Severity: check.Error,
					LineNo:   getuint64(10),
					Message: `Pattern "/b*/logs" shadows the following patterns:
            * 2: "/build/logs/"
Entries should go from least-specific to most-specific.`,
				},
			},
		},
		"Should not report any issues with correct CODEOWNERS file": {
			codeownersInput: FixtureValidCODEOWNERS,
			expectedIssues:  nil,
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// given
			sut := check.NewAvoidShadowing()

			// when
			out, err := sut.Check(context.TODO(), LoadInput(tc.codeownersInput))

			// then
			require.NoError(t, err)
			assert.ElementsMatch(t, tc.expectedIssues, out.Issues)
		})
	}
}
