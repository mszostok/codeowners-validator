package check_test

import (
	"testing"

	"github.com/mszostok/codeowners-validator/internal/check"

	"github.com/stretchr/testify/assert"
)

func TestAPIBuilder(t *testing.T) {
	var bldr *check.OutputBuilder

	t.Run("Does not panic on ReportIssue when builder is nil", func(t *testing.T) {
		assert.NotPanics(t, func() {
			issue := bldr.ReportIssue("test")
			assert.Nil(t, issue)
		})
	})

	t.Run("Does not panic on Output when builder is nil", func(t *testing.T) {
		assert.NotPanics(t, func() {
			out := bldr.Output()
			assert.Empty(t, out)
		})
	})
}
