package printer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/mszostok/codeowners-validator/internal/check"
)

// writer used for test purpose
var writer io.Writer = os.Stdout

type TTYPrinter struct {
	m sync.RWMutex
}

func (tty *TTYPrinter) PrintCheckResult(checkName string, duration time.Duration, checkOut check.Output) {
	tty.m.Lock()
	defer tty.m.Unlock()

	header := color.New(color.Bold).FprintfFunc()
	issueBody := color.New(color.FgWhite).FprintfFunc()
	okCheck := color.New(color.FgGreen).FprintlnFunc()

	header(writer, "==> Executing %s (%v)\n", checkName, duration)
	for _, i := range checkOut.Issues {
		issueSeverity := tty.severityPrintfFunc(i.Severity)

		issueSeverity(writer, "    [%s]", strings.ToLower(i.Severity.String()[:3]))
		if i.LineNo != nil {
			issueBody(writer, " line %d:", *i.LineNo)
		}
		issueBody(writer, " %s\n", i.Message)
	}

	if len(checkOut.Issues) == 0 {
		okCheck(writer, "    Check OK")
	}
}

func (*TTYPrinter) severityPrintfFunc(severity check.SeverityType) func(w io.Writer, format string, a ...interface{}) {
	p := color.New()
	switch severity {
	case check.Warning:
		p.Add(color.FgYellow)
	case check.Error:
		p.Add(color.FgRed)
	}

	return p.FprintfFunc()
}

func (*TTYPrinter) PrintSummary(allCheck, failedChecks int) {
	failures := "no"
	if failedChecks > 0 {
		failures = fmt.Sprintf("%d", failedChecks)
	}
	fmt.Fprintf(writer, "\n%d check(s) executed, %s failure(s)\n", allCheck, failures)
}
