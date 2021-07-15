package check

import (
	"context"
	"net/http"
	"net/mail"
	"strings"

	"github.com/mszostok/codeowners-validator/internal/ctxutil"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

type ValidOwnerConfig struct {
	Repository string
	// IgnoredOwners contains a list of owners that should not be validated.
	// Defaults to @ghost.
	// More info about the @ghost user: https://docs.github.com/en/free-pro-team@latest/github/setting-up-and-managing-your-github-user-account/deleting-your-user-account
	// Tip on how @ghost can be used: https://github.community/t5/How-to-use-Git-and-GitHub/CODEOWNERS-file-with-a-NOT-file-type-condition/m-p/31013/highlight/true#M8523
	IgnoredOwners []string `envconfig:"default=@ghost"`
}

// ValidOwner validates each owner
type ValidOwner struct {
	ghClient    *github.Client
	orgMembers  *map[string]struct{}
	orgName     string
	orgTeams    []*github.Team
	orgRepoName string
	ignOwners   map[string]struct{}
}

// NewValidOwner returns new instance of the ValidOwner
func NewValidOwner(cfg ValidOwnerConfig, ghClient *github.Client) (*ValidOwner, error) {
	split := strings.Split(cfg.Repository, "/")
	if len(split) != 2 {
		return nil, errors.Errorf("Wrong repository name. Expected pattern 'owner/repository', got '%s'", cfg.Repository)
	}

	ignOwners := map[string]struct{}{}
	for _, n := range cfg.IgnoredOwners {
		ignOwners[n] = struct{}{}
	}

	return &ValidOwner{
		ghClient:    ghClient,
		orgName:     split[0],
		orgRepoName: split[1],
		ignOwners:   ignOwners,
	}, nil
}

// Check if defined owners are the valid ones.
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
func (v *ValidOwner) Check(ctx context.Context, in Input) (Output, error) {
	var bldr OutputBuilder

	checkedOwners := map[string]struct{}{}

	for _, entry := range in.CodeownersEntries {
		for _, ownerName := range entry.Owners {
			if ctxutil.ShouldExit(ctx) {
				return Output{}, ctx.Err()
			}

			if v.isIgnoredOwner(ownerName) {
				continue
			}

			if _, alreadyChecked := checkedOwners[ownerName]; alreadyChecked {
				continue
			}

			validFn := v.selectValidateFn(ownerName)
			if err := validFn(ctx, ownerName); err != nil {
				bldr.ReportIssue(err.msg, WithEntry(entry))
				if err.permanent { // Doesn't make sense to process further
					return bldr.Output(), nil
				}
			}
			checkedOwners[ownerName] = struct{}{}
		}
	}

	return bldr.Output(), nil
}

func isEmailAddress(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

func isGithubTeam(s string) bool {
	hasPrefix := strings.HasPrefix(s, "@")
	containsSlash := strings.Contains(s, "/")
	split := strings.SplitN(s, "/", 3) // 3 is enough to confirm that is invalid + will not overflow the buffer
	return hasPrefix && containsSlash && len(split) == 2 && len(split[1]) > 0
}

func isGithubUser(s string) bool {
	return !strings.Contains(s, "/") && strings.HasPrefix(s, "@")
}

func (v *ValidOwner) isIgnoredOwner(name string) bool {
	_, found := v.ignOwners[name]
	return found
}

func (v *ValidOwner) selectValidateFn(name string) func(context.Context, string) *validateError {
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
			return newValidateError("Not valid owner definition %q", name)
		}
	}
}

func (v *ValidOwner) initOrgListTeams(ctx context.Context) *validateError {
	var teams []*github.Team
	req := &github.ListOptions{
		PerPage: 100,
	}
	for {
		resultPage, resp, err := v.ghClient.Teams.ListTeams(ctx, v.orgName, req)
		if err != nil { // TODO(mszostok): implement retry?
			switch err := err.(type) {
			case *github.ErrorResponse:
				if err.Response.StatusCode == http.StatusUnauthorized {
					return newValidateError("Teams for organization %q could not be queried. Requires GitHub authorization.", v.orgName)
				}
				return newValidateError("HTTP error occurred while calling GitHub: %v", err)
			case *github.RateLimitError:
				return newValidateError("GitHub rate limit reached: %v", err.Message)
			default:
				return newValidateError("Unknown error occurred while calling GitHub: %v", err)
			}
		}
		teams = append(teams, resultPage...)
		if resp.NextPage == 0 {
			break
		}
		req.Page = resp.NextPage
	}

	v.orgTeams = teams

	return nil
}

func (v *ValidOwner) validateTeam(ctx context.Context, name string) *validateError {
	if v.orgTeams == nil {
		if err := v.initOrgListTeams(ctx); err != nil {
			return err.AsPermanent()
		}
	}

	// GitHub normalizes name before comparison
	name = strings.ToLower(name)
	// called after validation it's safe to work on `parts` slice
	parts := strings.SplitN(name, "/", 2)
	org := parts[0]
	org = strings.TrimPrefix(org, "@")
	team := parts[1]

	if org != v.orgName {
		return newValidateError("Team %q does not belongs to %q organization.", team, v.orgName)
	}

	teamExists := func() bool {
		for _, v := range v.orgTeams {
			if v.GetSlug() == team {
				return true
			}
		}
		return false
	}

	if !teamExists() {
		return newValidateError("Team %q does not exist in organization %q.", team, org)
	}

	// repo contains the permissions for the team slug given
	// TODO(mszostok): Switch to GraphQL API, see:
	//   https://github.com/mszostok/codeowners-validator/pull/62#discussion_r561273525
	repo, _, err := v.ghClient.Teams.IsTeamRepoBySlug(ctx, v.orgName, team, org, v.orgRepoName)
	if err != nil { // TODO(mszostok): implement retry?
		switch err := err.(type) {
		case *github.ErrorResponse:
			if err.Response.StatusCode == http.StatusUnauthorized {
				return newValidateError(
					"Team permissions information for %q/%q could not be queried. Requires GitHub authorization.",
					org, v.orgRepoName)
			} else if err.Response.StatusCode == http.StatusNotFound {
				return newValidateError(
					"Team %q does not have permissions associated with the repository %q.",
					team, v.orgRepoName)
			} else {
				return newValidateError("HTTP error occurred while calling GitHub: %v", err)
			}
		case *github.RateLimitError:
			return newValidateError("GitHub rate limit reached: %v", err.Message)
		default:
			return newValidateError("Unknown error occurred while calling GitHub: %v", err)
		}
	}

	teamHasWritePermission := func() bool {
		for k, v := range repo.GetPermissions() {
			if !v {
				continue
			}

			switch k {
			case
				"admin",
				"maintain",
				"push":
				return true
			case
				"pull",
				"triage":
			}
		}

		return false
	}

	if !teamHasWritePermission() {
		return newValidateError(
			"Team %q cannot review PRs on %q as neither it nor any parent team has write permissions.",
			team, v.orgRepoName)
	}

	return nil
}

func (v *ValidOwner) validateGithubUser(ctx context.Context, name string) *validateError {
	if v.orgMembers == nil { //TODO(mszostok): lazy init, make it more robust.
		if err := v.initOrgListMembers(ctx); err != nil {
			return newValidateError("Cannot initialize organization member list: %v", err).AsPermanent()
		}
	}

	userName := strings.TrimPrefix(name, "@")
	_, _, err := v.ghClient.Users.Get(ctx, userName)
	if err != nil { // TODO(mszostok): implement retry?
		switch err := err.(type) {
		case *github.ErrorResponse:
			if err.Response.StatusCode == http.StatusNotFound {
				return newValidateError("User %q does not have github account", name)
			}
			return newValidateError("HTTP error occurred while calling GitHub: %v", err).AsPermanent()
		case *github.RateLimitError:
			return newValidateError("GitHub rate limit reached: %v", err.Message).AsPermanent()
		default:
			return newValidateError("Unknown error occurred while calling GitHub: %v", err).AsPermanent()
		}
	}

	_, isMember := (*v.orgMembers)[userName]
	if !isMember {
		return newValidateError("User %q is not a member of the organization", name)
	}

	return nil
}

// There is a method to check if user is a org member
//  client.Organizations.IsMember(context.Background(), "org-name", "user-name")
// But latency is too huge for checking each single user independent
// better and faster is to ask for all members and cache them.
func (v *ValidOwner) initOrgListMembers(ctx context.Context) error {
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

// Name returns human readable name of the validator
func (ValidOwner) Name() string {
	return "Valid Owner Checker"
}
