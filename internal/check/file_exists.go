package check

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mszostok/codeowners-validator/internal/ctxutil"

	"github.com/mattn/go-zglob"
	"github.com/pkg/errors"
)

type FileExistsConfig struct {
	FailureLevel SeverityType `envconfig:"default=error"`
}

type FileExist struct {
	cfg FileExistsConfig
}

func NewFileExist(cfg FileExistsConfig) *FileExist {
	return &FileExist{
		cfg: cfg,
	}
}

func (f *FileExist) Check(ctx context.Context, in Input) (Output, error) {
	var bldr OutputBuilder

	errTemplate := "%q does not match any files in repository"
	for _, entry := range in.CodeownersEntries {
		if ctxutil.ShouldExit(ctx) {
			return Output{}, ctx.Err()
		}

		fullPath := filepath.Join(in.RepoDir, f.fnmatchPattern(entry.Pattern))
		matches, err := zglob.Glob(fullPath)
		switch {
		case err == nil:
		case errors.Is(err, os.ErrNotExist):
			msg := fmt.Sprintf(errTemplate, entry.Pattern)
			bldr.ReportIssue(msg, WithEntry(entry), WithSeverity(f.cfg.FailureLevel))
			continue
		default:
			errTemplate := "while checking if there is any file in %s matching pattern %s"
			return Output{}, errors.Wrapf(err, errTemplate, in.RepoDir, entry.Pattern)
		}

		if len(matches) == 0 {
			msg := fmt.Sprintf(errTemplate, entry.Pattern)
			bldr.ReportIssue(msg, WithEntry(entry), WithSeverity(f.cfg.FailureLevel))
		}
	}

	return bldr.Output(), nil
}

func (*FileExist) fnmatchPattern(pattern string) string {
	if len(pattern) >= 2 && pattern[:1] == "*" && pattern[1:2] != "*" {
		return "**/" + pattern
	}

	return pattern
}

func (*FileExist) Name() string {
	return "File Exist Checker"
}
