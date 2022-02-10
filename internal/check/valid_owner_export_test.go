package check

func IsValidOwner(owner string) bool {
	return isEmailAddress(owner) || isGithubUser(owner) || isGithubTeam(owner)
}
