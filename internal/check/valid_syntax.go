package check

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/mszostok/codeowners-validator/internal/ctxutil"
)

var (
	// A valid username/organization name has up to 39 characters (per GitHub Join page)
	// and is matched by the following regex: /^[a-z\d](?:[a-z\d]|-(?=[a-z\d])){0,38}$/i
	// A valid team name consists of alphanumerics, underscores and dashes
	usernameOrTeamRegexp = regexp.MustCompile(`^@(?i:[a-z\d](?:[a-z\d-]){0,37}[a-z\d](/[a-z\d](?:[a-z\d_-]*)[a-z\d])?)$`)

	// Per: https://davidcel.is/posts/stop-validating-email-addresses-with-regex/
	// just check if there is '@' and a '.' afterwards
	emailRegexp = regexp.MustCompile(`.+@.+\..+`)
)

// ValidSyntax provides a syntax validation for CODEOWNERS file.
//
// If any line in your CODEOWNERS file contains invalid syntax, the file will not be detected and will
// not be used to request reviews. Invalid syntax includes inline comments and user or team names that do not exist on GitHub.
type ValidSyntax struct{}

// NewValidSyntax returns new ValidSyntax instance.
func NewValidSyntax() *ValidSyntax {
	return &ValidSyntax{}
}

// Check for syntax issues in your CODEOWNERS file.
func (ValidSyntax) Check(ctx context.Context, in Input) (Output, error) {
	var bldr OutputBuilder

	for _, entry := range in.CodeownersEntries {
		if ctxutil.ShouldExit(ctx) {
			return Output{}, ctx.Err()
		}

		if entry.Pattern == "" {
			bldr.ReportIssue("Missing pattern", WithEntry(entry))
		}

	ownersLoop:
		for _, item := range entry.Owners {
			switch {
			case strings.EqualFold(item, "#"):
				msg := "Comment (# sign) is not allowed in line with pattern entry. The correct format is: pattern owner1 ... ownerN"
				bldr.ReportIssue(msg, WithEntry(entry))
				break ownersLoop // no need to check for the rest items in this line, as the whole line is already marked as invalid
			case strings.HasPrefix(item, "@"):
				if !usernameOrTeamRegexp.MatchString(item) {
					msg := fmt.Sprintf("Owner '%s' does not look like a GitHub username or team name", item)
					bldr.ReportIssue(msg, WithEntry(entry), WithSeverity(Warning))
				}
			default:
				if !emailRegexp.MatchString(item) {
					msg := fmt.Sprintf("Owner '%s' does not look like an email", item)
					bldr.ReportIssue(msg, WithEntry(entry))
				}
			}
		}
	}

	return bldr.Output(), nil
}

func (ValidSyntax) Name() string {
	return "Valid Syntax Checker"
}
