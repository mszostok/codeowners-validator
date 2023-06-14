package config

import "go.szostok.io/codeowners/internal/api"

const (
	DefaultConfigFilename = "codeowners-config.yaml"
	EnvPrefix             = "CODEOWNERS"
)

// Config holds the application configuration
type Config struct {
	Checks                           []string         `mapstructure:"checks"`
	CheckFailureLevel                api.SeverityType `mapstructure:"check-failure-level"`
	ExperimentalChecks               []string         `mapstructure:"experimental-checks"`
	GithubAccessToken                string           `mapstructure:"github-access-token"`
	GithubBaseURL                    string           `mapstructure:"github-base-url"`
	GithubUploadURL                  string           `mapstructure:"github-upload-url"`
	GithubAppID                      int64            `mapstructure:"github-app-id"`
	GithubAppInstallationID          int64            `mapstructure:"github-app-installation-id"`
	GithubAppPrivateKey              string           `mapstructure:"github-app-private-key"`
	NotOwnedCheckerSkipPatterns      []string         `mapstructure:"not-owned-checker-skip-patterns"`
	NotOwnedCheckerSubdirectories    []string         `mapstructure:"not-owned-checker-subdirectories"`
	NotOwnedCheckerTrustWorkspace    bool             `mapstructure:"not-owned-checker-trust-workspace"`
	OwnerCheckerRepository           string           `mapstructure:"owner-checker-repository"`
	OwnerCheckerIgnoredOwners        []string         `mapstructure:"owner-checker-ignored-owners"`
	OwnerCheckerAllowUnownedPatterns bool             `mapstructure:"owner-checker-allow-unowned-patterns"`
	OwnerCheckerOwnersMustBeTeams    bool             `mapstructure:"owner-checker-owners-must-be-teams"`
	RepositoryPath                   string           `mapstructure:"repository-path"`
}
