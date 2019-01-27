package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/mszostok/codeowners-validator/internal/check"
	"github.com/mszostok/codeowners-validator/pkg/codeowners"
	"github.com/mszostok/codeowners-validator/internal/printer"
	"github.com/mszostok/codeowners-validator/internal/runner"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
	"golang.org/x/oauth2"
)

type Config struct {
	RepositoryPath    string
	GithubAccessToken string
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
	if cfg.GithubAccessToken != "" {
		httpClient = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.GithubAccessToken},
		))
	}

	ghClient := github.NewClient(httpClient)

	// init codeowners entries
	codeownersEntries, err := codeowners.NewFromPath(cfg.RepositoryPath)
	if err != nil {
		log.Fatal(err)
	}

	// gather checks
	checks := []check.Checker{
		check.NewFileExist(),
		check.NewValidOwner(cfg.ValidOwnerChecker, ghClient),
	}

	// run check runner
	absRepoPath, err := filepath.Abs(cfg.RepositoryPath)
	if err != nil {
		log.Fatal(err)
	}
	checkRunner := runner.NewCheckRunner(log, &printer.TTYPrinter{}, codeownersEntries, absRepoPath, checks...)
	checkRunner.Run(ctx)
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
