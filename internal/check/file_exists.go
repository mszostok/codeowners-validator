package check

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.szostok.io/codeowners/internal/api"
	"go.szostok.io/codeowners/internal/ctxutil"

	"github.com/mattn/go-zglob"
	"github.com/pkg/errors"
)

type FileExist struct{}

func NewFileExist() *FileExist {
	return &FileExist{}
}

func (f *FileExist) Check(ctx context.Context, in api.Input) (api.Output, error) {
	var bldr api.OutputBuilder

	for _, entry := range in.CodeownersEntries {
		if ctxutil.ShouldExit(ctx) {
			return api.Output{}, ctx.Err()
		}

		fullPath := filepath.Join(in.RepoDir, f.fnmatchPattern(entry.Pattern))
		matches, err := zglob.Glob(fullPath)
		switch {
		case err == nil:
		case errors.Is(err, os.ErrNotExist):
			msg := fmt.Sprintf("%q does not match any files in repository", entry.Pattern)
			bldr.ReportIssue(msg, api.WithEntry(entry))
			continue
		default:
			return api.Output{}, errors.Wrapf(err, "while checking if there is any file in %s matching pattern %s", in.RepoDir, entry.Pattern)
		}

		if len(matches) == 0 {
			msg := fmt.Sprintf("%q does not match any files in repository", entry.Pattern)
			bldr.ReportIssue(msg, api.WithEntry(entry))
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
