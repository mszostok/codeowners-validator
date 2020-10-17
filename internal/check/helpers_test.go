package check_test

import (
	"strings"

	"github.com/mszostok/codeowners-validator/internal/check"
	"github.com/mszostok/codeowners-validator/pkg/codeowners"
)

var validCODEOWNERS = `
		# These owners will be the default owners for everything
		*       @global-owner1 @global-owner2

		# js owner
		*.js    @js-owner

		*.go docs@example.com

		/build/logs/ @doctocat

		/script m.t@g.com
`

func loadInput(in string) check.Input {
	r := strings.NewReader(in)

	return check.Input{
		CodeownersEntries: codeowners.ParseCodeowners(r),
	}
}
