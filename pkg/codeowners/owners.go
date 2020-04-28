package codeowners

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
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

	return parseCodeowners(r)
}

// openCodeownersFile finds a CODEOWNERS file and returns content.
// see: https://help.github.com/articles/about-code-owners/#codeowners-file-location
func openCodeownersFile(dir string) (io.Reader, error) {
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

		return fs.Open(f)
	}

	return nil, fmt.Errorf("No CODEOWNERS found in the root, docs/, or .github/ directory of the repository %s", dir)
}

func parseCodeowners(r io.Reader) ([]Entry, error) {
	var e []Entry

	usernameOrTeamRegexp := regexp.MustCompile(`^@(?i:[a-z\d](?:[a-z\d-]){0,37}[a-z\d](/[a-z\d](?:[a-z\d_-]*)[a-z\d])?)$`)
	emailRegexp := regexp.MustCompile(`.+@.+\..+`)

	s := bufio.NewScanner(r)
	no := uint64(0)
	for s.Scan() {
		no++

		line := s.Text()
		if strings.HasPrefix(line, "#") { // comment
			continue
		}

		if len(line) == 0 { // empty
			continue
		}

		fields := strings.Fields(s.Text())

		if len(fields) < 2 {
			return e, fmt.Errorf("line %d does not have 2 or more fields", no)
		}

		// This does syntax validation only
		//
		// Syntax check: all fields are valid team/username identifiers or emails
		// Allowed owner syntax:
		// @username
		// @org/team-name
		// user@example.com
		// source: https://help.github.com/articles/about-code-owners/#codeowners-syntax
		for _, entry := range fields[1:] {
			if strings.HasPrefix(entry, "@") {
				// A valid username/organization name has up to 39 characters (per GitHub Join page)
				// and is matched by the following regex: /^[a-z\d](?:[a-z\d]|-(?=[a-z\d])){0,38}$/i
				// A valid team name consists of alphanumerics, underscores and dashes
				if !usernameOrTeamRegexp.MatchString(entry) {
					return e, fmt.Errorf("entry '%s' on line %d does not look like a GitHub username or team name", entry, no)
				}
			} else {
				// Per: https://davidcel.is/posts/stop-validating-email-addresses-with-regex/
				// just check if there is '@' and a '.' afterwards
				if !emailRegexp.MatchString(entry) {
					return e, fmt.Errorf("entry '%s' on line %d does not look like an email", entry, no)
				}
			}
		}

		e = append(e, Entry{
			Pattern: fields[0],
			Owners:  fields[1:],
			LineNo:  no,
		})
	}

	return e, nil
}
