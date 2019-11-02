package check

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/mszostok/codeowners-validator/pkg/codeowners"
	"github.com/pkg/errors"
	"gopkg.in/pipe.v2"

	ctxutil "github.com/mszostok/codeowners-validator/internal/context"
)

type NotOwnedFileCheckerConfig struct {
	SkipPatterns []string `envconfig:"optional"`
}

type NotOwnedFileChecker struct {
	skipPatterns map[string]struct{}
}

func NewNotOwnedFile(cfg NotOwnedFileCheckerConfig) *NotOwnedFileChecker {
	skip := map[string]struct{}{}
	for _, p := range cfg.SkipPatterns {
		skip[p] = struct{}{}
	}

	return &NotOwnedFileChecker{
		skipPatterns: skip,
	}
}

func (c *NotOwnedFileChecker) Check(ctx context.Context, in Input) (output Output, err error) {
	patterns, err := c.patternsToBeIgnored(ctx, in.CodeownerEntries)
	if err != nil {
		return Output{}, err
	}

	err = c.GitCheckStatus(in.RepoDir)
	if err != nil {
		return Output{}, err
	}

	defer func() {
		errReset := c.GitResetCurrentBranch(in.RepoDir)
		if err != nil {
			output = Output{}
			err = multierror.Append(err, errReset).ErrorOrNil()
		}
	}()

	err = c.AppendToGitignoreFile(in.RepoDir, patterns)
	if err != nil {
		return Output{}, err
	}

	err = c.GitRemoveIgnoredFiles(in.RepoDir)
	if err != nil {
		return Output{}, err
	}

	out, err := c.GitListFiles(in.RepoDir)
	if err != nil {
		return Output{}, err
	}

	lsOut := strings.TrimSpace(out)
	if lsOut != "" {
		lines := strings.Split(lsOut, "\n")
		msg := fmt.Sprintf("Found %d not owned files (skipped patterns: %q): \n%s", len(lines), c.skipPatternsList(), c.ListFormatFunc(lines))
		output.ReportIssue(msg, WithSeverity(Warning))
	}

	return output, nil
}

func (c *NotOwnedFileChecker) patternsToBeIgnored(ctx context.Context, entries []codeowners.Entry) ([]string, error) {
	var patterns []string
	for _, entry := range entries {
		if ctxutil.ShouldExit(ctx) {
			return []string{}, ctx.Err()
		}

		if _, found := c.skipPatterns[entry.Pattern]; found {
			continue
		}
		patterns = append(patterns, entry.Pattern)
	}

	return patterns, nil
}

func (c *NotOwnedFileChecker) AppendToGitignoreFile(repoDir string, patterns []string) error {
	f, err := os.OpenFile(path.Join(repoDir, ".gitignore"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	content := strings.Builder{}
	for _, p := range patterns {
		content.WriteString(fmt.Sprintf("%s\n", p))
	}

	_, err = f.WriteString(content.String())
	if err != nil {
		return err
	}
	return nil
}

func (c *NotOwnedFileChecker) GitRemoveIgnoredFiles(repoDir string) error {
	gitrm := pipe.Script(
		pipe.ChDir(repoDir),
		pipe.Line(
			pipe.Exec("git", "ls-files", "-ci", "--exclude-standard", "-z"),
			pipe.Exec("xargs", "-0", "git", "rm", "--cached"),
		),
	)

	_, stderr, err := pipe.DividedOutput(gitrm)
	if err != nil {
		return errors.Wrap(err, string(stderr))
	}
	return nil
}

func (c *NotOwnedFileChecker) GitCheckStatus(repoDir string) error {
	gitstate := pipe.Script(
		pipe.ChDir(repoDir),
		pipe.Exec("git", "status", "--porcelain"),
	)

	out, stderr, err := pipe.DividedOutput(gitstate)
	if err != nil {
		return errors.Wrap(err, string(stderr))
	}

	if string(out) != "" {
		return fmt.Errorf("git state is dirty: commit all changes before executing %q", c.Name())
	}

	return nil
}

func (c *NotOwnedFileChecker) GitResetCurrentBranch(repoDir string) error {
	gitreset := pipe.Script(
		pipe.ChDir(repoDir),
		pipe.Exec("git", "reset", "--hard"),
	)
	_, stderr, err := pipe.DividedOutput(gitreset)
	if err != nil {
		return errors.Wrap(err, string(stderr))
	}
	return nil
}

func (c *NotOwnedFileChecker) GitListFiles(repoDir string) (string, error) {
	gitls := pipe.Script(
		pipe.ChDir(repoDir),
		pipe.Exec("git", "ls-files"),
	)

	stdout, stderr, err := pipe.DividedOutput(gitls)
	if err != nil {
		return "", errors.Wrap(err, string(stderr))
	}

	return string(stdout), nil
}

func (c *NotOwnedFileChecker) skipPatternsList() string {
	list := make([]string, 0, len(c.skipPatterns))
	for k := range c.skipPatterns {
		list = append(list, k)
	}
	return strings.Join(list, ",")
}

// ListFormatFunc is a basic formatter that outputs
// a bullet point list of the pattern.
func (c *NotOwnedFileChecker) ListFormatFunc(es []string) string {
	points := make([]string, len(es))
	for i, err := range es {
		points[i] = fmt.Sprintf("            * %s", err)
	}

	return strings.Join(points, "\n")
}

// Name returns human readable name of the validator
func (NotOwnedFileChecker) Name() string {
	return "[Experimental] Not Owned File Checker"
}
