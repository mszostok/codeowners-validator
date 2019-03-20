package url_test

import (
	"github.com/mszostok/codeowners-validator/pkg/url"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanonicalURLPath(t *testing.T) {
	tests := map[string]struct {
		givenPath string
		expPath   string
	}{
		"no trailing slash": {
			givenPath: "https://api.github.com",
			expPath:   "https://api.github.com/",
		},
		"multiple trailing slashes": {
			givenPath: "https://api.github.com///////////////",
			expPath:   "https://api.github.com/",
		},
		"single trailing slash": {
			givenPath: "https://api.github.com/",
			expPath:   "https://api.github.com/",
		},
	}
	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// when
			normalizedPath := url.CanonicalPath(tc.givenPath)

			// then
			assert.Equal(t, tc.expPath, normalizedPath)
		})
	}
}
