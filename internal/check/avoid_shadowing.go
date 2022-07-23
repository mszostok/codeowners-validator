package check

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"go.szostok.io/codeowners-validator/internal/ctxutil"
	"go.szostok.io/codeowners-validator/pkg/codeowners"
)

type AvoidShadowing struct{}

func NewAvoidShadowing() *AvoidShadowing {
	return &AvoidShadowing{}
}

func (c *AvoidShadowing) Check(ctx context.Context, in Input) (output Output, err error) {
	var bldr OutputBuilder

	previousEntries := []codeowners.Entry{}
	for _, entry := range in.CodeownersEntries {
		if ctxutil.ShouldExit(ctx) {
			return Output{}, ctx.Err()
		}
		re, err := wildCardToRegexp(endWithSlash(entry.Pattern))
		if err != nil {
			return Output{}, errors.Wrapf(err, "while compiling pattern %s into a regexp", entry.Pattern)
		}
		shadowed := []codeowners.Entry{}
		for _, previous := range previousEntries {
			if re.MatchString(endWithSlash(previous.Pattern)) {
				shadowed = append(shadowed, previous)
			}
		}
		if len(shadowed) > 0 {
			msg := fmt.Sprintf("Pattern %q shadows the following patterns:\n%s\nEntries should go from least-specific to most-specific.", entry.Pattern, c.listFormatFunc(shadowed))
			bldr.ReportIssue(msg, WithEntry(entry))
		}
		previousEntries = append(previousEntries, entry)
	}

	return bldr.Output(), nil
}

// listFormatFunc is a basic formatter that outputs a bullet point list of the pattern.
func (c *AvoidShadowing) listFormatFunc(es []codeowners.Entry) string {
	points := make([]string, len(es))
	for i, err := range es {
		points[i] = fmt.Sprintf("            * %d: %q", err.LineNo, err.Pattern)
	}

	return strings.Join(points, "\n")
}

// Name returns human readable name of the validator
func (AvoidShadowing) Name() string {
	return "[Experimental] Avoid Shadowing Checker"
}

// endWithSlash adds a trailing slash to a string if it doesn't already end with one.
// This is useful when matching CODEOWNERS pattern because the trailing slash is optional.
func endWithSlash(s string) string {
	if !strings.HasSuffix(s, "/") {
		return s + "/"
	}
	return s
}

// wildCardToRegexp converts a wildcard pattern to a regular expression pattern.
func wildCardToRegexp(pattern string) (*regexp.Regexp, error) {
	var result strings.Builder
	for i, literal := range strings.Split(pattern, "*") {
		// Replace * with .*
		if i > 0 {
			result.WriteString(".*")
		}

		// Quote any regular expression meta characters in the
		// literal text.
		result.WriteString(regexp.QuoteMeta(literal))
	}
	return regexp.Compile("^" + result.String() + "$")
}
