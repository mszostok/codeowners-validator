package runner

import (
	"context"
	"sync"
	"time"

	"go.szostok.io/codeowners/internal/api"
	"go.szostok.io/codeowners/internal/printer"
	"go.szostok.io/codeowners/pkg/codeowners"

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
	PrintCheckResult(checkName string, duration time.Duration, checkOut api.Output, err error)
	PrintSummary(allCheck int, failedChecks int)
}

// CheckRunner runs all registered checks in parallel.
// Needs to be initialized via NewCheckRunner func.
type CheckRunner struct {
	m                  sync.RWMutex
	log                logrus.FieldLogger
	codeowners         []codeowners.Entry
	repoPath           string
	treatedAsFailure   api.SeverityType
	checks             []api.Checker
	printer            Printer
	allFoundIssues     map[api.SeverityType]uint32
	notPassedChecksCnt int
}

// NewCheckRunner is a constructor for CheckRunner
func NewCheckRunner(log logrus.FieldLogger, co []codeowners.Entry, repoPath string, treatedAsFailure api.SeverityType, checks ...api.Checker) *CheckRunner {
	return &CheckRunner{
		log:              log.WithField("service", "check:runner"),
		repoPath:         repoPath,
		treatedAsFailure: treatedAsFailure,
		codeowners:       co,
		checks:           checks,

		printer:        &printer.TTYPrinter{},
		allFoundIssues: map[api.SeverityType]uint32{},
	}
}

// Run executes given test in a loop with given throttle
func (r *CheckRunner) Run(ctx context.Context) {
	wg := sync.WaitGroup{}

	// TODO(mszostok): timeout per check?
	wg.Add(len(r.checks))
	for _, c := range r.checks {
		go func(c api.Checker) {
			defer wg.Done()
			startTime := time.Now()
			out, err := c.Check(ctx, api.Input{
				CodeownersEntries: r.codeowners,
				RepoDir:           r.repoPath,
			})

			r.collectMetrics(out, err)
			r.printer.PrintCheckResult(c.Name(), time.Since(startTime), out, err)
		}(c)
	}
	wg.Wait()

	r.printer.PrintSummary(len(r.checks), r.notPassedChecksCnt)
}

func (r *CheckRunner) ShouldExitWithCheckFailure() bool {
	higherOccurredIssue := api.SeverityType(MaxInt)
	for key := range r.allFoundIssues {
		if higherOccurredIssue > key {
			higherOccurredIssue = key
		}
	}

	return higherOccurredIssue <= r.treatedAsFailure
}

func (r *CheckRunner) collectMetrics(checkOut api.Output, err error) {
	r.m.Lock()
	defer r.m.Unlock()
	for _, i := range checkOut.Issues {
		r.allFoundIssues[i.Severity]++
	}

	if err != nil {
		r.allFoundIssues[api.Error]++
	}

	if len(checkOut.Issues) > 0 || err != nil {
		r.notPassedChecksCnt++
	}
}
