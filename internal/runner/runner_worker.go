package runner

import (
	"context"
	"sync"
	"time"

	"github.com/mszostok/codeowners-validator/internal/check"
	"github.com/mszostok/codeowners-validator/pkg/codeowners"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Printer prints the checks results
type Printer interface {
	PrintCheckResult(checkName string, duration time.Duration, checkOut check.Output)
	PrintSummary(allCheck uint64, failedChecks uint64)
}

// CheckRunner runs all registered checks in parallel.
// Needs to be initialized via NewCheckRunner func.
type CheckRunner struct {
	log      logrus.FieldLogger
	checks   []check.Checker
	repoPath string
	c        []codeowners.Entry
	printer  Printer
}

// NewCheckRunner is a constructor for CheckRunner
func NewCheckRunner(log logrus.FieldLogger, printer Printer, c []codeowners.Entry, repoPath string, checks ...check.Checker) *CheckRunner {
	return &CheckRunner{
		log:      log.WithField("service", "check:runner"),
		checks:   checks,
		c:        c,
		printer:  printer,
		repoPath: repoPath,
	}
}

// Run executes given test in a loop with given throttle
func (r *CheckRunner) Run(ctx context.Context) {
	// timeout per check?

	wg := sync.WaitGroup{}
	wg.Add(len(r.checks))
	failedCnt := uint64(0)
	for _, c := range r.checks {
		go func(c check.Checker) {
			defer wg.Done()
			startTime := time.Now()
			out, err := c.Check(ctx, check.Input{
				CodeownerEntries: r.c,
				RepoDir:          r.repoPath,
			})
			if err != nil {
				r.log.Errorf(errors.Wrapf(err, "while executing checker %s", c.Name()).Error())
				return
			}

			if len(out.Issues) != 0 {
				failedCnt++
			}

			r.printer.PrintCheckResult(c.Name(), time.Since(startTime), out)
		}(c)
	}

	wg.Wait()

	r.printer.PrintSummary(uint64(len(r.checks)), failedCnt)
}
