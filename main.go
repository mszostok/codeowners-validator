package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/mszostok/codeowners-validator/internal/check"
	"github.com/mszostok/codeowners-validator/internal/runner"
	"github.com/mszostok/codeowners-validator/pkg/codeowners"
	"github.com/mszostok/codeowners-validator/pkg/url"
	"github.com/mszostok/codeowners-validator/pkg/version"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
	"golang.org/x/oauth2"
)

// Config holds the application configuration
type Config struct {
	RepositoryPath string
	Checks         []string `envconfig:"optional"`
	Github         struct {
		AccessToken string `envconfig:"optional"`
		BaseURL     string `envconfig:"optional"`
		UploadURL   string `envconfig:"optional"`
	}
	CheckFailureLevel  check.SeverityType `envconfig:"default=warning"`
	OwnerChecker       check.ValidOwnerCheckerConfig
	NotOwnedChecker    check.NotOwnedFileCheckerConfig
	ExperimentalChecks []string `envconfig:"optional"`
}

func main() {
	if version.ShouldPrintVersion() {
		version.PrintVersion(os.Stdout)
		os.Exit(0)
	}

	var cfg Config
	err := envconfig.Init(&cfg)
	fatalOnError(err)

	log := logrus.New()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	cancelOnInterrupt(ctx, cancelFunc)

	// init codeowners entries
	codeownersEntries, err := codeowners.NewFromPath(cfg.RepositoryPath)
	fatalOnError(err)

	// init checks
	var checks []check.Checker

	if isEnabled(cfg.Checks, "duppattern") {
		checks = append(checks, check.NewDuplicatedPattern())
	}

	if isEnabled(cfg.Checks, "files") {
		checks = append(checks, check.NewFileExist())
	}

	if isEnabled(cfg.Checks, "owners") {
		// init GitHub client
		httpClient := http.DefaultClient
		if cfg.Github.AccessToken != "" {
			httpClient = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: cfg.Github.AccessToken},
			))
		}

		ghClient, err := newGithubClient(cfg, httpClient)
		fatalOnError(err)

		checks = append(checks, check.NewValidOwner(cfg.OwnerChecker, ghClient))
	}

	if contains(cfg.ExperimentalChecks, "notowned") {
		checks = append(checks, check.NewNotOwnedFile(cfg.NotOwnedChecker))
	}

	// run check runner
	absRepoPath, err := filepath.Abs(cfg.RepositoryPath)
	fatalOnError(err)

	checkRunner := runner.NewCheckRunner(log, codeownersEntries, absRepoPath, cfg.CheckFailureLevel, checks...)
	checkRunner.Run(ctx)

	if ctx.Err() != nil {
		log.Error("Application was interrupted by operating system")
		os.Exit(2)
	}
	if checkRunner.ShouldExitWithCheckFailure() {
		os.Exit(3)
	}
}

func isEnabled(checks []string, name string) bool {
	// if a user does not specify concrete checks then all checks are enabled
	if len(checks) == 0 {
		return true
	}

	if contains(checks, name) {
		return true
	}
	return false
}

func contains(checks []string, name string) bool {
	for _, c := range checks {
		if c == name {
			return true
		}
	}
	return false
}

func newGithubClient(cfg Config, httpClient *http.Client) (ghClient *github.Client, err error) {
	baseURL, uploadURL := cfg.Github.BaseURL, cfg.Github.UploadURL

	if baseURL != "" {
		if uploadURL == "" { // often the baseURL are same as the uploadURL, so we do not require to provide both of them
			uploadURL = baseURL
		}

		bURL, uURL := url.CanonicalPath(baseURL), url.CanonicalPath(uploadURL)
		ghClient, err = github.NewEnterpriseClient(bURL, uURL, httpClient)

	} else {
		ghClient = github.NewClient(httpClient)
	}

	return
}

func fatalOnError(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

// cancelOnInterrupt calls cancel func when os.Interrupt or SIGTERM is received
func cancelOnInterrupt(ctx context.Context, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-ctx.Done():
		case <-c:
			cancel()
			<-c
			os.Exit(1) // second signal. Exit directly.
		}
	}()
}
