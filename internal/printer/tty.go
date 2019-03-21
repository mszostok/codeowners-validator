package printer

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mszostok/codeowners-validator/internal/check"
)

type TTYPrinter struct{}

func (tty TTYPrinter) PrintCheckResult(checkName string, duration time.Duration, checkOut check.Output) {
	header := color.New(color.Bold).PrintfFunc()
	issueBody := color.New(color.FgWhite).PrintfFunc()
	okCheck := color.New(color.FgGreen).PrintlnFunc()

	header("==> Executing %s (%v)\n", checkName, duration)
	for _, i := range checkOut.Issues {
		issueSeverity := tty.severityPrintfFunc(i.Severity)

		issueSeverity("    [%s]", strings.ToLower(i.Severity.String()[:3]))
		issueBody(" line %d: %s\n", i.LineNo, i.Message)
	}

	if len(checkOut.Issues) == 0 {
		okCheck("    Check OK")
	}
}

func (TTYPrinter) severityPrintfFunc(severity check.SeverityType) func(format string, a ...interface{}) {
	p := color.New()
	switch severity {
	case check.Warning:
		p.Add(color.FgYellow)
	case check.Error:
		p.Add(color.FgRed)
	}

	return p.PrintfFunc()
}

func (TTYPrinter) PrintSummary(allCheck uint64, failedChecks uint64) {
	failures := "no"
	if failedChecks > 0 {
		failures = fmt.Sprintf("%d", failedChecks)
	}
	fmt.Printf("\n%d check(s) executed, %s failure(s)\n", allCheck, failures)
}
