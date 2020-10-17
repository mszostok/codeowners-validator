package check_test

import (
	"context"
	"testing"

	"github.com/mszostok/codeowners-validator/internal/check"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuplicatedPattern(t *testing.T) {
	tests := map[string]struct {
		codeownersInput string
		expectedIssues  []check.Issue
	}{
		"Should report info about duplicated entries": {
			codeownersInput: `
					*       @global-owner1 @global-owner2
					
					/build/logs/ @doctocat
					/build/logs/ @doctocat

					/script @mszostok
					/script m.t@g.com
			`,
			expectedIssues: []check.Issue{
				{
					Severity: check.Error,
					LineNo:   nil,
					Message: `Pattern "/build/logs/" is defined 2 times in lines: 
            * 4: with owners: [@doctocat]
            * 5: with owners: [@doctocat]`,
				},
				{
					Severity: check.Error,
					LineNo:   nil,
					Message: `Pattern "/script" is defined 2 times in lines: 
            * 7: with owners: [@mszostok]
            * 8: with owners: [m.t@g.com]`,
				},
			},
		},
		"Should not report any issues with correct CODEOWNERS file": {
			codeownersInput: validCODEOWNERS,
			expectedIssues:  nil,
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// given
			sut := check.NewDuplicatedPattern()

			// when
			out, err := sut.Check(context.TODO(), loadInput(tc.codeownersInput))

			// then
			require.NoError(t, err)
			assert.ElementsMatch(t, tc.expectedIssues, out.Issues)
		})
	}
}
