package check

import (
	"context"
	"fmt"
	"path/filepath"

	ctxutil "github.com/mszostok/codeowners-validator/internal/context"
	"github.com/pkg/errors"
)

type FileExist struct{}

func NewFileExist() *FileExist {
	return &FileExist{}
}

func (FileExist) Check(ctx context.Context, in Input) (Output, error) {
	var output Output

	for _, entry := range in.CodeownerEntries {
		if ctxutil.ShouldExit(ctx) {
			return Output{}, ctx.Err()
		}

		fullPath := filepath.Join(in.RepoDir, entry.Pattern)
		matches, err := filepath.Glob(fullPath)
		if err != nil {
			return Output{}, errors.Wrapf(err, "while checking if there is any file in %s matching pattern %s", in.RepoDir, entry.Pattern)
		}

		if len(matches) == 0 {
			msg := fmt.Sprintf("%q does not match any files in repository", entry.Pattern)
			output.ReportIssue(msg, WithEntry(entry))
		}
	}

	return output, nil
}

func (FileExist) Name() string {
	return "File Exist Checker"
}
