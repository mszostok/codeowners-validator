package load

import (
	"context"

	"github.com/mszostok/codeowners-validator/internal/check"
	"github.com/mszostok/codeowners-validator/internal/envconfig"
	"github.com/mszostok/codeowners-validator/internal/github"

	"github.com/pkg/errors"
)

// For now, it is a good enough solution to init checks. Important thing is to do not require env variables
// and do not create clients which will not be used because of the given checker.
//
// MAYBE in the future the https://github.com/uber-go/dig will be used.
func Checks(ctx context.Context, enabledChecks []string, experimentalChecks []string) ([]check.Checker, error) {
	var checks []check.Checker

	if isEnabled(enabledChecks, "syntax") {
		checks = append(checks, check.NewValidSyntax())
	}

	if isEnabled(enabledChecks, "duppatterns") {
		checks = append(checks, check.NewDuplicatedPattern())
	}

	if isEnabled(enabledChecks, "files") {
		checks = append(checks, check.NewFileExist())
	}

	if isEnabled(enabledChecks, "owners") {
		var cfg struct {
			OwnerChecker check.ValidOwnerConfig
			Github       github.ClientConfig
		}
		if err := envconfig.Init(&cfg); err != nil {
			return nil, errors.Wrapf(err, "while loading config for %s", "owners")
		}

		ghClient, err := github.NewClient(ctx, cfg.Github)
		if err != nil {
			return nil, errors.Wrap(err, "while creating GitHub client")
		}

		owners, err := check.NewValidOwner(cfg.OwnerChecker, ghClient)
		if err != nil {
			return nil, errors.Wrap(err, "while enabling 'owners' checker")
		}
		checks = append(checks, owners)
	}

	expChecks, err := loadExperimentalChecks(experimentalChecks)
	if err != nil {
		return nil, errors.Wrap(err, "while loading experimental checks")
	}

	return append(checks, expChecks...), nil
}

func loadExperimentalChecks(experimentalChecks []string) ([]check.Checker, error) {
	var checks []check.Checker

	if contains(experimentalChecks, "notowned") {
		var cfg struct {
			NotOwnedChecker check.NotOwnedFileConfig
		}
		if err := envconfig.Init(&cfg); err != nil {
			return nil, errors.Wrapf(err, "while loading config for %s", "notowned")
		}

		checks = append(checks, check.NewNotOwnedFile(cfg.NotOwnedChecker))
	}

	return checks, nil
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
