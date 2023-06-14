package api

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"go.szostok.io/codeowners/internal/ptr"
	"go.szostok.io/codeowners/pkg/codeowners"
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
		RepoDir           string
		CodeownersEntries []codeowners.Entry
	}

	Output struct {
		Issues []Issue
	}

	OutputBuilder struct {
		mux    sync.Mutex
		issues []Issue
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

func (bldr *OutputBuilder) ReportIssue(msg string, opts ...ReportIssueOpt) *OutputBuilder {
	if bldr == nil { // TODO: error?
		return nil
	}

	i := Issue{
		Severity: Error,
		Message:  msg,
	}

	for _, opt := range opts {
		opt(&i)
	}

	bldr.mux.Lock()
	defer bldr.mux.Unlock()
	bldr.issues = append(bldr.issues, i)

	return bldr
}

func (bldr *OutputBuilder) Output() Output {
	if bldr == nil {
		return Output{}
	}
	return Output{Issues: bldr.issues}
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

func (s *SeverityType) Set(in string) error {
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

func (s *SeverityType) Type() string {
	return "SeverityType"
}
