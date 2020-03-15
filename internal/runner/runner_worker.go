package runner

import (
	"context"
	"sync"
	"time"

	"github.com/mszostok/codeowners-validator/internal/check"
	"github.com/mszostok/codeowners-validator/internal/printer"
	"github.com/mszostok/codeowners-validator/pkg/codeowners"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	// MaxUint defines the max unsigned int value.
	MaxUint = ^uint(0)
	// MaxInt defines the max signed int value.
	MaxInt = int(MaxUint >> 1)
)

// Printer prints the checks results
type Printer interface {
	PrintCheckResult(checkName string, duration time.Duration, checkOut check.Output)
	PrintSummary(allCheck int, failedChecks int)
}

// CheckRunner runs all registered checks in parallel.
// Needs to be initialized via NewCheckRunner func.
type CheckRunner struct {
	m                  sync.RWMutex
	log                logrus.FieldLogger
	codeowners         []codeowners.Entry
	repoPath           string
	treatedAsFailure   check.SeverityType
	checks             []check.Checker
	printer            Printer
	allFoundIssues     map[check.SeverityType]uint32
	notPassedChecksCnt int
}

// NewCheckRunner is a constructor for CheckRunner
func NewCheckRunner(log logrus.FieldLogger, co []codeowners.Entry, repoPath string, treatedAsFailure check.SeverityType, checks ...check.Checker) *CheckRunner {
	return &CheckRunner{
		log:              log.WithField("service", "check:runner"),
		repoPath:         repoPath,
		treatedAsFailure: treatedAsFailure,
		codeowners:       co,
		checks:           checks,

		printer:        &printer.TTYPrinter{},
		allFoundIssues: map[check.SeverityType]uint32{},
	}
}

// Run executes given test in a loop with given throttle
func (r *CheckRunner) Run(ctx context.Context) {
	wg := sync.WaitGroup{}

	// TODO(mszostok): timeout per check?
	wg.Add(len(r.checks))
	for _, c := range r.checks {
		go func(c check.Checker) {
			defer wg.Done()
			startTime := time.Now()
			out, err := c.Check(ctx, check.Input{
				CodeownerEntries: r.codeowners,
				RepoDir:          r.repoPath,
			})
			if err != nil {
				// TODO(mszostok): add err handling (logging it internally is not enough)
				r.log.Errorf(errors.Wrapf(err, "while executing checker %s", c.Name()).Error())
				return
			}

			r.collectMetrics(out)

			r.printer.PrintCheckResult(c.Name(), time.Since(startTime), out)
		}(c)
	}
	wg.Wait()

	r.printer.PrintSummary(len(r.checks), r.notPassedChecksCnt)
}

func (r *CheckRunner) ShouldExitWithCheckFailure() bool {
	higherOccurredIssue := check.SeverityType(MaxInt)
	for key := range r.allFoundIssues {
		if higherOccurredIssue > key {
			higherOccurredIssue = key
		}
	}

	return higherOccurredIssue <= r.treatedAsFailure
}

func (r *CheckRunner) collectMetrics(checkOut check.Output) {
	r.m.Lock()
	defer r.m.Unlock()
	for _, i := range checkOut.Issues {
		r.allFoundIssues[i.Severity]++
	}

	if len(checkOut.Issues) > 0 {
		r.notPassedChecksCnt++
	}
}
