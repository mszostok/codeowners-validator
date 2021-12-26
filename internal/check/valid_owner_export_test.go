package check

func IsValidOwner(owner string) bool {
	return isEmailAddress(owner) || isGitHubUser(owner) || isGitHubTeam(owner)
}
