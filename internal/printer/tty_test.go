package printer

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/mszostok/codeowners-validator/internal/check"
	"github.com/mszostok/codeowners-validator/internal/ptr"

	"github.com/sebdah/goldie/v2"
)

func TestTTYPrinterPrintCheckResult(t *testing.T) {
	t.Run("Should print all reported issues", func(t *testing.T) {
		// given
		tty := TTYPrinter{}

		buff := &bytes.Buffer{}
		restore := overrideWriter(buff)
		defer restore()

		// when
		tty.PrintCheckResult("Foo Checker", time.Second, check.Output{
			Issues: []check.Issue{
				{
					Severity: check.Error,
					LineNo:   ptr.Uint64Ptr(42),
					Message:  "Simulate error in line 42",
				},
				{
					Severity: check.Warning,
					LineNo:   ptr.Uint64Ptr(2020),
					Message:  "Simulate warning in line 2020",
				},
				{
					Severity: check.Error,
					Message:  "Error without line number",
				},
				{
					Severity: check.Warning,
					Message:  "Warning without line number",
				},
			},
		})

		// then
		g := goldie.New(t, goldie.WithNameSuffix(".golden.txt"))
		g.Assert(t, t.Name(), buff.Bytes())
	})

	t.Run("Should print OK status on empty issues list", func(t *testing.T) {
		// given
		tty := TTYPrinter{}

		buff := &bytes.Buffer{}
		restore := overrideWriter(buff)
		defer restore()

		// when
		tty.PrintCheckResult("Foo Checker", time.Second, check.Output{
			Issues: nil,
		})

		// then
		g := goldie.New(t, goldie.WithNameSuffix(".golden.txt"))
		g.Assert(t, t.Name(), buff.Bytes())
	})
}

func TestTTYPrinterPrintSummary(t *testing.T) {
	t.Run("Should print number of failures", func(t *testing.T) {
		// given
		tty := TTYPrinter{}

		buff := &bytes.Buffer{}
		restore := overrideWriter(buff)
		defer restore()

		// when
		tty.PrintSummary(20, 10)

		// then
		g := goldie.New(t, goldie.WithNameSuffix(".golden.txt"))
		g.Assert(t, t.Name(), buff.Bytes())
	})

	t.Run("Should print 'no' when there is no failures", func(t *testing.T) {
		// given
		tty := TTYPrinter{}

		buff := &bytes.Buffer{}
		restore := overrideWriter(buff)
		defer restore()

		// when
		tty.PrintSummary(20, 0)

		// then
		g := goldie.New(t, goldie.WithNameSuffix(".golden.txt"))
		g.Assert(t, t.Name(), buff.Bytes())
	})
}

func overrideWriter(in io.Writer) func() {
	old := writer
	writer = in
	return func() { writer = old }
}
