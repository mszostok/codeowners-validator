package check

import (
	"context"

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
		LineNo   uint64
		Message  string
	}

	Input struct {
		RepoDir          string
		CodeownerEntries []codeowners.Entry
	}

	Output struct {
		Issues []Issue
	}
)

type Opt func(*Issue)

func WithSeverity(s SeverityType) Opt {
	return func(i *Issue) {
		i.Severity = s
	}
}

// TODO: decide where to put it
func (out *Output) ReportIssue(e codeowners.Entry, msg string, opts ...Opt) Issue {
	if out == nil { // TODO: error?
		return Issue{}
	}

	i := Issue{
		Severity: Error,
		LineNo:   e.LineNo,
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
	Error SeverityType = iota
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
