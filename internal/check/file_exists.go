package check

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mszostok/codeowners-validator/internal/ctxutil"

	"github.com/mattn/go-zglob"
	"github.com/pkg/errors"
)

type FileExist struct{}

func NewFileExist() *FileExist {
	return &FileExist{}
}

func (f *FileExist) Check(ctx context.Context, in Input) (Output, error) {
	var bldr OutputBuilder

	for _, entry := range in.CodeownersEntries {
		if ctxutil.ShouldExit(ctx) {
			return Output{}, ctx.Err()
		}

		fullPath := filepath.Join(in.RepoDir, f.fnmatchPattern(entry.Pattern))
		matches, err := zglob.Glob(fullPath)

		switch {
		case err == nil:
		case errors.Is(err, os.ErrNotExist):
			msg := fmt.Sprintf("%q does not match any files in repository", entry.Pattern)
			bldr.ReportIssue(msg, WithEntry(entry))
			continue
		default:
			return Output{}, errors.Wrapf(err, "while checking if there is any file in %s matching pattern %s", in.RepoDir, entry.Pattern)
		}

		// Glob returns also matched directories, but only files count
		foundFile := false
		for _, name := range matches {
			fi, err := os.Stat(name)
			if err != nil {
				return Output{}, errors.Wrap(err, "while checking if found path is a file")
			}
			if fi.Mode().IsRegular() {
				foundFile = true
				break
			}
		}

		if !foundFile {
			msg := fmt.Sprintf("%q does not match any files in repository", entry.Pattern)
			bldr.ReportIssue(msg, WithEntry(entry))
		}
	}

	return bldr.Output(), nil
}

func (*FileExist) fnmatchPattern(pattern string) string {
	noOfChars := len(pattern)
	if noOfChars >= 2 && pattern[:1] == "*" && pattern[1:2] != "*" {
		return "**/" + pattern
	}

	if noOfChars > 1 && pattern[:1] != "/" {
		return "**/" + pattern + "**/*"
	}

	lastIdx := strings.LastIndex(pattern, "/")
	// handle such pattern: **/foo @pico
	// we need to add the **/* to find all files, without that Glob
	// will return only the first match which can be still a directory
	if lastIdx != -1 && filepath.Ext(pattern[lastIdx+1:]) == "" {
		return pattern + "**/*"
	}

	return pattern
}

func (*FileExist) Name() string {
	return "File Exist Checker"
}
