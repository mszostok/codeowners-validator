package check

import (
	"context"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/google/go-github/github"
	ctxutil "github.com/mszostok/codeowners-validator/internal/context"
)

type ValidOwnerCheckerConfig struct {
	OrganizationName string
}

// ValidOwnerChecker validates each owner
type ValidOwnerChecker struct {
	ghClient   *github.Client
	orgMembers *map[string]struct{}
	orgName    string
}

// NewValidOwner returns new instance of the ValidOwnerChecker
func NewValidOwner(cfg ValidOwnerCheckerConfig, ghClient *github.Client) *ValidOwnerChecker {
	return &ValidOwnerChecker{
		ghClient: ghClient,
		orgName:  cfg.OrganizationName,
	}
}

// Check check if defined owners are the valid ones.
// Allowed owner syntax:
// @username
// @org/team-name
// user@example.com
// source: https://help.github.com/articles/about-code-owners/#codeowners-syntax
//
// Checks:
// - if owner is one of: github user, org team, email address
// - if github user then check if have github account
// - if github user then check if he/she is in organization
// - if org team then check if exists in organization
func (v *ValidOwnerChecker) Check(ctx context.Context, in Input) (Output, error) {
	var output Output
	checkedOwners := map[string]struct{}{}

	for _, entry := range in.CodeownerEntries {
		for _, ownerName := range entry.Owners {
			if ctxutil.ShouldExit(ctx) {
				return Output{}, ctx.Err()
			}

			if _, alreadyChecked := checkedOwners[ownerName]; alreadyChecked {
				continue
			}

			validFn := v.selectValidateFn(ownerName)
			if err := validFn(ctx, ownerName); err != nil {
				output.ReportIssue(err.Msg, WithSeverity(err.Severity), WithEntry(entry))
				if err.RateLimitReached { // Doesn't make sense to process further. TODO(mszostok): change for more generic solution like, `IsPermanentError`
					return output, nil
				}
			}
			checkedOwners[ownerName] = struct{}{}
		}
	}

	return output, nil
}

func (v *ValidOwnerChecker) selectValidateFn(name string) func(context.Context, string) *validateError {
	switch {
	case isGithubUser(name):
		return v.validateGithubUser
	case isGithubTeam(name):
		return v.validateTeam
	case isEmailAddress(name):
		// TODO(mszostok): try to check if e-mail really exists
		return func(context.Context, string) *validateError { return nil }
	default:
		return func(_ context.Context, name string) *validateError {
			return &validateError{fmt.Sprintf("Not valid owner definition %q", name), Error, false}
		}
	}
}

type validateError struct {
	Msg              string
	Severity         SeverityType
	RateLimitReached bool
}

func (v *ValidOwnerChecker) validateTeam(ctx context.Context, name string) *validateError {
	parts := strings.SplitN(name, "/", 2)
	org := parts[0]
	org = strings.TrimPrefix(org, "@")
	team := parts[1]

	allTeams, _, err := v.ghClient.Teams.ListTeams(ctx, org, nil)
	if err != nil { // TODO(mszostok): implement retry?
		switch err := err.(type) {
		case *github.ErrorResponse:
			if err.Response.StatusCode == http.StatusUnauthorized {
				return &validateError{fmt.Sprintf("Team %q could not be check. Requires GitHub authorization.", name), Warning, false}
			}
			return &validateError{fmt.Sprintf("HTTP error occurred while calling GitHub: %v", err), Error, false}
		case *github.RateLimitError:
			return &validateError{fmt.Sprintf("GitHub rate limit reached: %v", err.Message), Warning, true}
		default:
			return &validateError{fmt.Sprintf("Unknown error occurred while calling GitHub: %v", err), Error, false}
		}
	}

	teamExists := func() bool {
		for _, v := range allTeams {
			if v.GetSlug() == team {
				return true
			}
		}
		return false
	}

	teamHasPermissions := func() bool {
		for _, v := range allTeams {
			if v.GetPermission() != "pull" {
				return true
			}
		}
		return false
	}

	if !teamExists() {
		return &validateError{fmt.Sprintf("Team %q does not exits in organization %q or has no permissions associated with the repository.", team, org), Warning, false}
	}

	if !teamHasPermissions() {
		return &validateError{fmt.Sprintf("Team %q doesn't have write access to the repository.", team), Warning, false}
	}

	return nil
}

func (v *ValidOwnerChecker) validateGithubUser(ctx context.Context, name string) *validateError {
	if v.orgMembers == nil { //TODO(mszostok): lazy init, make it more robust.
		if err := v.initOrgListMembers(ctx); err != nil {
			return &validateError{fmt.Sprintf("Cannot initialize organization member list: %v", err), Error, false}
		}
	}

	userName := strings.TrimPrefix(name, "@")
	_, _, err := v.ghClient.Users.Get(ctx, userName)
	if err != nil { // TODO(mszostok): implement retry?
		switch err := err.(type) {
		case *github.ErrorResponse:
			if err.Response.StatusCode == http.StatusNotFound {
				return &validateError{fmt.Sprintf("User %q does not have github account", name), Error, false}
			}
			return &validateError{fmt.Sprintf("HTTP error occurred while calling GitHub: %v", err), Error, false}
		case *github.RateLimitError:
			return &validateError{fmt.Sprintf("GitHub rate limit reached: %v", err.Message), Warning, true}
		default:
			return &validateError{fmt.Sprintf("Unknown error occurred while calling GitHub: %v", err), Error, false}
		}
	}

	_, isMember := (*v.orgMembers)[userName]
	if !isMember {
		return &validateError{fmt.Sprintf("User %q is not a member of the organization", name), Error, false}
	}

	return nil
}

// There is a method to check if user is a org member
//  client.Organizations.IsMember(context.Background(), "org-name", "user-name")
// But latency is too huge for checking each single user independent
// better and faster is to ask for all members and cache them.
func (v *ValidOwnerChecker) initOrgListMembers(ctx context.Context) error {
	opt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allMembers []*github.User
	for {
		users, resp, err := v.ghClient.Organizations.ListMembers(ctx, v.orgName, opt)
		if err != nil {
			return err
		}
		allMembers = append(allMembers, users...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	v.orgMembers = &map[string]struct{}{}
	for _, u := range allMembers {
		(*v.orgMembers)[u.GetLogin()] = struct{}{}
	}

	return nil
}

func isEmailAddress(s string) bool {
	_, err := mail.ParseAddress(s)

	return err == nil
}

func isGithubTeam(s string) bool {
	hasPrefix := strings.HasPrefix(s, "@")
	containsSlash := strings.Contains(s, "/")
	splited := strings.SplitN(s, "/", 3) // 3 is enough to confirm that is invalid + will not overflow the buffer

	if hasPrefix && containsSlash && len(splited) == 2 {
		return true
	}

	return false
}

func isGithubUser(s string) bool {
	if strings.HasPrefix(s, "@") && !strings.Contains(s, "/") {
		return true
	}
	return false
}

// Name returns human readable name of the validator
func (ValidOwnerChecker) Name() string {
	return "Valid Owner Checker"
}
