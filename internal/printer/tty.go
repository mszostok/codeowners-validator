package printer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"go.szostok.io/codeowners/internal/api"
)

// writer used for test purpose
var writer io.Writer = os.Stdout

type TTYPrinter struct {
	m sync.RWMutex
}

func (tty *TTYPrinter) PrintCheckResult(checkName string, duration time.Duration, checkOut api.Output, checkErr error) {
	tty.m.Lock()
	defer tty.m.Unlock()

	header := color.New(color.Bold).FprintfFunc()
	issueBody := color.New(color.FgWhite).FprintfFunc()
	okCheck := color.New(color.FgGreen).FprintlnFunc()
	errCheck := color.New(color.FgRed).FprintfFunc()

	header(writer, "==> Executing %s (%v)\n", checkName, duration)
	for _, i := range checkOut.Issues {
		issueSeverity := tty.severityPrintfFunc(i.Severity)

		issueSeverity(writer, "    [%s]", strings.ToLower(i.Severity.String()[:3]))
		if i.LineNo != nil {
			issueBody(writer, " line %d:", *i.LineNo)
		}
		issueBody(writer, " %s\n", i.Message)
	}

	switch {
	case checkErr == nil && len(checkOut.Issues) == 0:
		okCheck(writer, "    Check OK")
	case checkErr != nil:
		errCheck(writer, "    [Internal Error]")
		issueBody(writer, " %s\n", checkErr)
	}
}

func (*TTYPrinter) severityPrintfFunc(severity api.SeverityType) func(w io.Writer, format string, a ...interface{}) {
	p := color.New()
	switch severity {
	case api.Warning:
		p.Add(color.FgYellow)
	case api.Error:
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
