package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/mszostok/codeowners-validator/internal/check"
	"github.com/mszostok/codeowners-validator/internal/printer"
	"github.com/mszostok/codeowners-validator/internal/runner"
	"github.com/mszostok/codeowners-validator/pkg/codeowners"
	"github.com/mszostok/codeowners-validator/pkg/url"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
	"golang.org/x/oauth2"
)

// Config holds the application configuration
type Config struct {
	RepositoryPath string
	Github         struct {
		AccessToken string `envconfig:"optional"`
		BaseURL     string `envconfig:"optional"`
		UploadURL   string `envconfig:"optional"`
	}
	ValidOwnerChecker check.ValidOwnerCheckerConfig
}

func main() {
	var cfg Config
	err := envconfig.Init(&cfg)
	fatalOnError(err)

	log := logrus.New()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	cancelOnInterrupt(ctx, cancelFunc)

	// init GitHub client
	httpClient := http.DefaultClient
	if cfg.Github.AccessToken != "" {
		httpClient = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.Github.AccessToken},
		))
	}

	ghClient, err := newGithubClient(cfg, httpClient)
	fatalOnError(err)

	// init codeowners entries
	codeownersEntries, err := codeowners.NewFromPath(cfg.RepositoryPath)
	fatalOnError(err)

	// gather checks
	checks := []check.Checker{
		check.NewFileExist(),
		check.NewValidOwner(cfg.ValidOwnerChecker, ghClient),
	}

	// run check runner
	absRepoPath, err := filepath.Abs(cfg.RepositoryPath)
	fatalOnError(err)

	checkRunner := runner.NewCheckRunner(log, &printer.TTYPrinter{}, codeownersEntries, absRepoPath, checks...)
	checkRunner.Run(ctx)
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
