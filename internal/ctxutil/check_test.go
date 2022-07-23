package ctxutil_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	contextutil "go.szostok.io/codeowners-validator/internal/ctxutil"
)

func TestShouldExit(t *testing.T) {
	t.Run("Should notify about exit if context is canceled", func(t *testing.T) {
		// given
		ctx, cancel := context.WithCancel(context.Background())

		// when
		cancel()
		shouldExit := contextutil.ShouldExit(ctx)

		// then
		assert.True(t, shouldExit)
	})

	t.Run("Should return false if context is not canceled", func(t *testing.T) {
		// given
		ctx := context.Background()

		// when
		shouldExit := contextutil.ShouldExit(ctx)

		// then
		assert.False(t, shouldExit)
	})
}
