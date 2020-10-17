package check

import (
	"context"
	"fmt"
	"strings"

	ctxutil "github.com/mszostok/codeowners-validator/internal/context"
	"github.com/mszostok/codeowners-validator/pkg/codeowners"
)

// DuplicatedPattern validates if CODEOWNERS file does not contain
// the duplicated lines with the same file pattern.
type DuplicatedPattern struct{}

// NewDuplicatedPattern returns instance of the DuplicatedPattern
func NewDuplicatedPattern() *DuplicatedPattern {
	return &DuplicatedPattern{}
}

// Check searches for doubles paths(patterns) in CODEOWNERS file.
func (d *DuplicatedPattern) Check(ctx context.Context, in Input) (Output, error) {
	var bldr OutputBuilder

	// TODO(mszostok): decide if the `CodeownersEntries` entry by default should be
	//  indexed by pattern (`map[string][]codeowners.Entry{}`)
	//  Required changes in pkg/codeowners/owners.go.
	patterns := map[string][]codeowners.Entry{}
	for _, entry := range in.CodeownersEntries {
		if ctxutil.ShouldExit(ctx) {
			return Output{}, ctx.Err()
		}

		patterns[entry.Pattern] = append(patterns[entry.Pattern], entry)
	}

	for name, entries := range patterns {
		if len(entries) > 1 {
			msg := fmt.Sprintf("Pattern %q is defined %d times in lines: \n%s", name, len(entries), d.listFormatFunc(entries))
			bldr.ReportIssue(msg)
		}
	}

	return bldr.Output(), nil
}

// listFormatFunc is a basic formatter that outputs a bullet point list of the pattern.
func (d *DuplicatedPattern) listFormatFunc(es []codeowners.Entry) string {
	points := make([]string, len(es))
	for i, err := range es {
		points[i] = fmt.Sprintf("            * %d: with owners: %s", err.LineNo, err.Owners)
	}

	return strings.Join(points, "\n")
}

// Name returns human readable name of the validator.
func (DuplicatedPattern) Name() string {
	return "Duplicated Pattern Checker"
}
