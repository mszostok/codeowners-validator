package check

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	ctxutil "github.com/mszostok/codeowners-validator/internal/context"

	"github.com/mattn/go-zglob"
	"github.com/pkg/errors"
)

type FileExist struct{}

func NewFileExist() *FileExist {
	return &FileExist{}
}

func (f *FileExist) Check(ctx context.Context, in Input) (Output, error) {
	var output Output

	for _, entry := range in.CodeownerEntries {
		if ctxutil.ShouldExit(ctx) {
			return Output{}, ctx.Err()
		}

		fullPath := filepath.Join(in.RepoDir, f.fnmatchPattern(entry.Pattern))
		matches, err := zglob.Glob(fullPath)
		switch {
		case err == nil:
		case errors.Is(err, os.ErrNotExist):
			msg := fmt.Sprintf("%q does not match any files in repository", entry.Pattern)
			output.ReportIssue(msg, WithEntry(entry))
			continue
		default:
			return Output{}, errors.Wrapf(err, "while checking if there is any file in %s matching pattern %s", in.RepoDir, entry.Pattern)
		}

		if len(matches) == 0 {
			msg := fmt.Sprintf("%q does not match any files in repository", entry.Pattern)
			output.ReportIssue(msg, WithEntry(entry))
		}
	}

	return output, nil
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
