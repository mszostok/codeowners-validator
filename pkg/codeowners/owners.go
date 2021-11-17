package codeowners

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/spf13/afero"
)

// Used for testing purposes
var fs = afero.NewOsFs()

// Entry contains owners for a given pattern
type Entry struct {
	LineNo  uint64
	Pattern string
	Owners  []string
}

func (e Entry) String() string {
	return fmt.Sprintf("line %d: %s\t%v", e.LineNo, e.Pattern, strings.Join(e.Owners, ", "))
}

// NewFromPath returns entries from codeowners
func NewFromPath(path string) ([]Entry, error) {
	r, err := openCodeownersFile(path)
	if err != nil {
		return nil, err
	}

	return ParseCodeowners(r), nil
}

// openCodeownersFile finds a CODEOWNERS file and returns content.
// see: https://help.github.com/articles/about-code-owners/#codeowners-file-location
func openCodeownersFile(dir string) (io.Reader, error) {
	var fileReader io.Reader = nil

	for _, p := range []string{".", "docs", ".github"} {
		pth := path.Join(dir, p)
		exists, err := afero.DirExists(fs, pth)
		if err != nil {
			return nil, err
		}

		if !exists {
			continue
		}

		f := path.Join(pth, "CODEOWNERS")
		_, err = fs.Stat(f)
		switch {
		case err == nil:
		case os.IsNotExist(err):
			continue
		default:
			return nil, err
		}
		if fileReader != nil {
			return nil, fmt.Errorf("Multiple CODEOWNERS files found in root, docs/, or .github/ directory of the repository %s", dir)
		}
		fileReader, err = fs.Open(f)
	}
	if fileReader != nil {
		return fileReader, nil
	}

	return nil, fmt.Errorf("No CODEOWNERS found in the root, docs/, or .github/ directory of the repository %s", dir)
}

func ParseCodeowners(r io.Reader) []Entry {
	var e []Entry
	s := bufio.NewScanner(r)
	no := uint64(0)
	for s.Scan() {
		no++
		fields := strings.Fields(s.Text())

		if len(fields) == 0 { // empty
			continue
		}

		if strings.HasPrefix(fields[0], "#") { // comment
			continue
		}

		e = append(e, Entry{
			Pattern: fields[0],
			Owners:  fields[1:],
			LineNo:  no,
		})
	}

	return e
}
