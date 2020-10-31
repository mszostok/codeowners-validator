package check

import (
	"strings"

	"github.com/mszostok/codeowners-validator/pkg/codeowners"
)

var FixtureValidCODEOWNERS = `
		# These owners will be the default owners for everything
		*       @global-owner1 @global-owner2

		# js owner
		*.js    @js-owner

		*.go docs@example.com

		/build/logs/ @doctocat

		/script m.t@g.com
`

func LoadInput(in string) Input {
	r := strings.NewReader(in)

	return Input{
		CodeownersEntries: codeowners.ParseCodeowners(r),
	}
}
