package check

import (
	"net/mail"
	"strings"
)

func IsEmailAddress(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

func IsGithubTeam(s string) bool {
	hasPrefix := strings.HasPrefix(s, "@")
	containsSlash := strings.Contains(s, "/")
	split := strings.SplitN(s, "/", 3) // 3 is enough to confirm that is invalid + will not overflow the buffer
	return hasPrefix && containsSlash && len(split) == 2 && len(split[1]) > 0
}

func IsGithubUser(s string) bool {
	return !strings.Contains(s, "/") && strings.HasPrefix(s, "@")
}

func IsGithubGhostUser(s string) bool {
	return s == "@ghost"
}
