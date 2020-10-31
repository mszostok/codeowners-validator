package check_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mszostok/codeowners-validator/internal/check"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRespectingCanceledContext(t *testing.T) {
	must := func(checker check.Checker, err error) check.Checker {
		require.NoError(t, err)
		return checker
	}

	checkers := []check.Checker{
		check.NewDuplicatedPattern(),
		check.NewFileExist(),
		check.NewValidSyntax(),
		check.NewNotOwnedFile(check.NotOwnedFileConfig{}),
		must(check.NewValidOwner(check.ValidOwnerConfig{Repository: "org/repo"}, nil)),
	}

	for _, checker := range checkers {
		sut := checker
		t.Run(checker.Name(), func(t *testing.T) {
			// given: canceled context
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			// when
			out, err := sut.Check(ctx, check.LoadInput(check.FixtureValidCODEOWNERS))

			// then
			assert.True(t, errors.Is(err, context.Canceled))
			assert.Empty(t, out)
		})
	}
}
