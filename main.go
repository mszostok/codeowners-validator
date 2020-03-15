package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/mszostok/codeowners-validator/internal/check"
	"github.com/mszostok/codeowners-validator/internal/load"
	"github.com/mszostok/codeowners-validator/internal/runner"
	"github.com/mszostok/codeowners-validator/pkg/codeowners"
	"github.com/mszostok/codeowners-validator/pkg/version"

	"github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
)

// Config holds the application configuration
type Config struct {
	RepositoryPath     string
	CheckFailureLevel  check.SeverityType `envconfig:"default=warning"`
	Checks             []string           `envconfig:"optional"`
	ExperimentalChecks []string           `envconfig:"optional"`
}

func main() {
	version.Init()
	if version.ShouldPrintVersion() {
		version.PrintVersion(os.Stdout)
		os.Exit(0)
	}

	var cfg Config
	err := envconfig.Init(&cfg)
	exitOnError(err)

	log := logrus.New()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	cancelOnInterrupt(ctx, cancelFunc)

	// init codeowners entries
	codeownersEntries, err := codeowners.NewFromPath(cfg.RepositoryPath)
	exitOnError(err)

	// init checks
	checks, err := load.Checks(ctx, cfg.Checks, cfg.ExperimentalChecks)
	exitOnError(err)

	// run check runner
	absRepoPath, err := filepath.Abs(cfg.RepositoryPath)
	exitOnError(err)

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

func exitOnError(err error) {
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
