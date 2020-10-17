package check

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/mszostok/codeowners-validator/internal/ptr"
	"github.com/mszostok/codeowners-validator/pkg/codeowners"
)

type (
	// Checker allows to execute check in a generic way
	Checker interface {
		Check(ctx context.Context, in Input) (Output, error)
		Name() string
	}

	Issue struct {
		Severity SeverityType // enum // default error
		LineNo   *uint64
		Message  string
	}

	Input struct {
		RepoDir          string
		CodeownerEntries []codeowners.Entry
	}

	Output struct {
		mux    sync.Mutex
		Issues []Issue
	}
)

type ReportIssueOpt func(*Issue)

func WithSeverity(s SeverityType) ReportIssueOpt {
	return func(i *Issue) {
		i.Severity = s
	}
}

func WithEntry(e codeowners.Entry) ReportIssueOpt {
	return func(i *Issue) {
		i.LineNo = ptr.Uint64Ptr(e.LineNo)
	}
}

func (out *Output) ReportIssue(msg string, opts ...ReportIssueOpt) Issue {
	out.mux.Lock()
	defer out.mux.Unlock()
	if out == nil { // TODO: error?
		return Issue{}
	}

	i := Issue{
		Severity: Error,
		Message:  msg,
	}

	for _, opt := range opts {
		opt(&i)
	}

	out.Issues = append(out.Issues, i)

	return i
}

type SeverityType int

const (
	Error SeverityType = iota + 1
	Warning
)

func (s SeverityType) String() string {
	switch s {
	case Error:
		return "Error"
	case Warning:
		return "Warning"
	default:
		return ""
	}
}

// Unmarshal provides custom parsing of severity type.
// Implements envconfig.Unmarshal interface.
func (s *SeverityType) Unmarshal(in string) error {
	switch strings.ToLower(in) {
	case "error", "err":
		*s = Error
	case "warning", "warn":
		*s = Warning
	default:
		return fmt.Errorf("not a valid severity type: %q", in)
	}

	return nil
}
