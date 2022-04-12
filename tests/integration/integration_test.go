//go:build integration
// +build integration

package integration

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

const (
	binaryPathEnvName                       = "BINARY_PATH"
	codeownersSamplesRepo                   = "https://github.com/gh-codeowners/codeowners-samples.git"
	caseInsensitiveOrgCodeownersSamplesRepo = "https://github.com/GitHubCODEOWNERS/codeowners-samples.git"
)

var repositories = []struct {
	name string
	repo string
}{
	{
		name: "gh-codeowners",
		repo: codeownersSamplesRepo,
	},
	{
		name: "GitHubCODEOWNERS",
		repo: caseInsensitiveOrgCodeownersSamplesRepo,
	},
}

// TestCheckHappyPath tests that codeowners-validator reports no issues for valid CODEOWNERS file.
//
// This test is based on golden file.
// If the `-test.update-golden` flag is set then the actual content is written
// to the golden file.
//
// To update golden file, run:
//   UPDATE_GOLDEN=true make test-integration
func TestCheckSuccess(t *testing.T) {
	type (
		Envs     map[string]string
		testCase []struct {
			name   string
			repo   string
			envs   Envs
			skipOS string
		}
	)

	t.Run("offline checks", func(t *testing.T) {
		for _, repoTC := range repositories {
			// given
			repoDir, cleanup := CloneRepo(t, repoTC.repo, "happy-path")

			tests := testCase{
				{
					name: "files",
					envs: Envs{
						"CHECKS": "files",
					},
				},
				{
					name: "duppatterns",
					envs: Envs{
						"CHECKS": "duppatterns",
					},
				},
				{
					name: "notowned",
					envs: Envs{
						"PATH":                os.Getenv("PATH"), // need to be set to find the `git` binary
						"CHECKS":              "disable-all",
						"EXPERIMENTAL_CHECKS": "notowned",
					},
					skipOS: "windows",
				},
			}
			for _, tc := range tests {
				t.Run(fmt.Sprintf("%s/%s", repoTC.name, tc.name), func(t *testing.T) {
					if runtime.GOOS == tc.skipOS {
						t.Skip("this test is marked as skipped for this OS")
					}

					binaryPath := os.Getenv(binaryPathEnvName)
					codeownersCmd := Exec().
						Binary(binaryPath).
						// codeowners-validator basic config
						WithEnv("REPOSITORY_PATH", repoDir)

					for k, v := range tc.envs {
						codeownersCmd.WithEnv(k, v)
					}

					// when
					result := codeownersCmd.AwaitResultAtMost(3 * time.Minute)

					// then
					assert.Equal(t, 0, result.ExitCode)
					normalizedOutput := normalizeTimeDurations(result.Stdout)

					g := goldie.New(t, goldie.WithNameSuffix(".golden.txt"))
					g.Assert(t, t.Name(), []byte(normalizedOutput))
				})
			}

			cleanup()
		}
	})

	t.Run("online checks", func(t *testing.T) {
		tests := testCase{
			{
				name: "gh-codeowners/owners",
				envs: Envs{
					"CHECKS":                   "owners",
					"OWNER_CHECKER_REPOSITORY": "gh-codeowners/codeowners-samples",
					"GITHUB_ACCESS_TOKEN":      os.Getenv("GITHUB_TOKEN"),
				},
				repo: codeownersSamplesRepo,
			},
			{
				name: "GitHubCODEOWNERS/owners",
				envs: Envs{
					"CHECKS":                   "owners",
					"OWNER_CHECKER_REPOSITORY": "GitHubCODEOWNERS/codeowners-samples",
					"GITHUB_ACCESS_TOKEN":      os.Getenv("GITHUB_TOKEN"),
				},
				repo: caseInsensitiveOrgCodeownersSamplesRepo,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given
				repoDir, cleanup := CloneRepo(t, tc.repo, "happy-path")
				defer cleanup()

				if runtime.GOOS == tc.skipOS {
					t.Skip("this test is marked as skipped for this OS")
				}

				binaryPath := os.Getenv(binaryPathEnvName)
				codeownersCmd := Exec().
					Binary(binaryPath).
					// codeowners-validator basic config
					WithEnv("REPOSITORY_PATH", repoDir)

				for k, v := range tc.envs {
					codeownersCmd.WithEnv(k, v)
				}

				// when
				result := codeownersCmd.AwaitResultAtMost(3 * time.Minute)

				// then
				assert.Equal(t, 0, result.ExitCode)
				normalizedOutput := normalizeTimeDurations(result.Stdout)

				g := goldie.New(t, goldie.WithNameSuffix(".golden.txt"))
				g.Assert(t, t.Name(), []byte(normalizedOutput))
			})
		}
	})
}

// TestCheckFailures tests that codeowners-validator reports issues for not valid CODEOWNERS file.
//
// This test is based on golden file.
// If the `-test.update-golden` flag is set then the actual content is written
// to the golden file.
//
// To update golden file, run:
//   UPDATE_GOLDEN=true make test-integration
func TestCheckFailures(t *testing.T) {
	type Envs map[string]string
	tests := []struct {
		name string
		envs Envs
	}{
		{
			name: "files",
			envs: Envs{
				"CHECKS": "files",
			},
		},
		{
			name: "owners",
			envs: Envs{
				"CHECKS":                               "owners",
				"OWNER_CHECKER_REPOSITORY":             "gh-codeowners/codeowners-samples",
				"OWNER_CHECKER_ALLOW_UNOWNED_PATTERNS": "false",
				"GITHUB_ACCESS_TOKEN":                  os.Getenv("GITHUB_TOKEN"),
			},
		},
		{
			name: "duppatterns",
			envs: Envs{
				"CHECKS": "duppatterns",
			},
		},
		{
			name: "notowned",
			envs: Envs{
				"PATH":                            os.Getenv("PATH"), // need to be set to find the `git` binary
				"CHECKS":                          "disable-all",
				"EXPERIMENTAL_CHECKS":             "notowned",
				"NOT_OWNED_CHECKER_SKIP_PATTERNS": "*",
			},
		},
		{
			name: "notowned_sub_dirs",
			envs: Envs{
				"PATH":                             os.Getenv("PATH"), // need to be set to find the `git` binary
				"CHECKS":                           "disable-all",
				"EXPERIMENTAL_CHECKS":              "notowned",
				"NOT_OWNED_CHECKER_SKIP_PATTERNS":  "*",
				"NOT_OWNED_CHECKER_SUBDIRECTORIES": "notowned/dir",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			repoDir, cleanup := CloneRepo(t, codeownersSamplesRepo, "failures")
			defer cleanup()

			binaryPath := os.Getenv(binaryPathEnvName)

			codeownersCmd := Exec().
				Binary(binaryPath).
				// codeowners-validator basic config
				WithEnv("REPOSITORY_PATH", repoDir)

			for k, v := range tc.envs {
				codeownersCmd.WithEnv(k, v)
			}

			// when
			result := codeownersCmd.AwaitResultAtMost(3 * time.Minute)

			// then
			assert.Equal(t, 3, result.ExitCode)

			normalizedOutput := normalizeTimeDurations(result.Stdout)

			g := goldie.New(t, goldie.WithNameSuffix(".golden.txt"))
			g.Assert(t, t.Name(), []byte(normalizedOutput))
		})
	}
}

// To update golden file, run:
//  TEST=TestOwnerCheckAuthZAndAuthN TOKEN_WITH_NO_SCOPES=<token_with_no_scopes> UPDATE_GOLDEN=true make test-integration
func TestOwnerCheckAuthZAndAuthN(t *testing.T) {
	t.Parallel()

	t.Run("token not specified", func(t *testing.T) {
		t.Parallel()

		// given
		codeownersCmd := Exec().
			Binary(os.Getenv(binaryPathEnvName)).
			WithEnv("REPOSITORY_PATH", "not-needed").
			WithEnv("CHECKS", "owners").
			WithEnv("OWNER_CHECKER_REPOSITORY", "gh-codeowners/codeowners-samples")

		// when
		result := codeownersCmd.AwaitResultAtMost(3 * time.Second)

		// then
		assert.Equal(t, 1, result.ExitCode)

		normalizedOutput := normalizeLogTime(result.Stderr)
		g := goldie.New(t, goldie.WithNameSuffix(".golden.txt"))
		g.Assert(t, t.Name(), []byte(normalizedOutput))

	})

	t.Run("invalid token specified", func(t *testing.T) {
		t.Parallel()

		// given
		codeownersCmd := Exec().
			Binary(os.Getenv(binaryPathEnvName)).
			WithEnv("REPOSITORY_PATH", "not-needed").
			WithEnv("CHECKS", "owners").
			WithEnv("OWNER_CHECKER_REPOSITORY", "gh-codeowners/codeowners-samples").
			WithEnv("GITHUB_ACCESS_TOKEN", "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")

		// when
		result := codeownersCmd.AwaitResultAtMost(3 * time.Second)

		// then
		assert.Equal(t, 1, result.ExitCode)

		normalizedOutput := normalizeLogTime(result.Stderr)
		g := goldie.New(t, goldie.WithNameSuffix(".golden.txt"))
		g.Assert(t, t.Name(), []byte(normalizedOutput))
	})

	t.Run("token specified but without necessary scopes", func(t *testing.T) {
		t.Parallel()

		// given
		codeownersCmd := Exec().
			Binary(os.Getenv(binaryPathEnvName)).
			WithEnv("REPOSITORY_PATH", "not-needed").
			WithEnv("CHECKS", "owners").
			WithEnv("OWNER_CHECKER_REPOSITORY", "gh-codeowners/codeowners-samples").
			WithEnv("GITHUB_ACCESS_TOKEN", os.Getenv("TOKEN_WITH_NO_SCOPES"))

		// when
		result := codeownersCmd.AwaitResultAtMost(3 * time.Second)

		// then
		assert.Equal(t, 1, result.ExitCode)

		normalizedOutput := normalizeLogTime(result.Stderr)
		g := goldie.New(t, goldie.WithNameSuffix(".golden.txt"))
		g.Assert(t, t.Name(), []byte(normalizedOutput))
	})

	t.Run("token specified but without necessary scopes and against private repo", func(t *testing.T) {
		t.Parallel()

		// given
		codeownersCmd := Exec().
			Binary(os.Getenv(binaryPathEnvName)).
			WithEnv("REPOSITORY_PATH", "not-needed").
			WithEnv("CHECKS", "owners").
			WithEnv("OWNER_CHECKER_REPOSITORY", "gh-codeowners/private-repo").
			WithEnv("GITHUB_ACCESS_TOKEN", os.Getenv("TOKEN_WITH_NO_SCOPES"))

		// when
		result := codeownersCmd.AwaitResultAtMost(3 * time.Second)

		// then
		assert.Equal(t, 1, result.ExitCode)

		normalizedOutput := normalizeLogTime(result.Stderr)
		g := goldie.New(t, goldie.WithNameSuffix(".golden.txt"))
		g.Assert(t, t.Name(), []byte(normalizedOutput))
	})
}

func TestMultipleChecksSuccess(t *testing.T) {
	t.Skip("not implemented yet")
}
