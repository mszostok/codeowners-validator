package load

import (
	"context"

	"go.szostok.io/codeowners/internal/api"
	"go.szostok.io/codeowners/internal/check"
	"go.szostok.io/codeowners/internal/config"
	"go.szostok.io/codeowners/internal/github"

	"github.com/pkg/errors"
	"github.com/vrischmann/envconfig"
)

// For now, it is a good enough solution to init checks. Important thing is to do not require env variables
// and do not create clients which will not be used because of the given checker.
//
// MAYBE in the future the https://github.com/uber-go/dig will be used.
func Checks(ctx context.Context, cfg *config.Config) ([]api.Checker, error) {
	var checks []api.Checker

	if isEnabled(cfg.Checks, "syntax") {
		checks = append(checks, check.NewValidSyntax())
	}

	if isEnabled(cfg.Checks, "duppatterns") {
		checks = append(checks, check.NewDuplicatedPattern())
	}

	if isEnabled(cfg.Checks, "files") {
		checks = append(checks, check.NewFileExist())
	}

	if isEnabled(cfg.Checks, "owners") {
		ghClient, isApp, err := github.NewClient(ctx, cfg)
		if err != nil {
			return nil, errors.Wrap(err, "while creating GitHub client")
		}

		owners, err := check.NewValidOwner(cfg, ghClient, !isApp)
		if err != nil {
			return nil, errors.Wrap(err, "while enabling 'owners' checker")
		}

		if err := owners.CheckSatisfied(ctx); err != nil {
			return nil, errors.Wrap(err, "while checking if 'owners' checker is satisfied")
		}

		checks = append(checks, owners)
	}

	expChecks, err := loadExperimentalChecks(cfg.ExperimentalChecks)
	if err != nil {
		return nil, errors.Wrap(err, "while loading experimental checks")
	}

	return append(checks, expChecks...), nil
}

func loadExperimentalChecks(experimentalChecks []string) ([]api.Checker, error) {
	var checks []api.Checker

	if contains(experimentalChecks, "notowned") {
		var cfg struct {
			NotOwnedChecker check.NotOwnedFileConfig
		}
		if err := envconfig.Init(&cfg); err != nil {
			return nil, errors.Wrapf(err, "while loading config for %s", "notowned")
		}

		checks = append(checks, check.NewNotOwnedFile(cfg.NotOwnedChecker))
	}

	if contains(experimentalChecks, "avoid-shadowing") {
		checks = append(checks, check.NewAvoidShadowing())
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
