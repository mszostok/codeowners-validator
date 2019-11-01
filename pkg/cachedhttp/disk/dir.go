// Heavily inspired by: https://github.com/kubernetes/kubernetes/blob/30c7df5cd822067016640aa267714204ac089172/staging/src/k8s.io/cli-runtime/pkg/genericclioptions/config_flags.go#L350
package disk

import (
	"github.com/pkg/errors"
	"path/filepath"
	"regexp"
	"strings"
)

// ComputeCacheDir takes the parentDir and the host and comes up with a "usually non-colliding" name.
func ComputeCacheDir(parentDir, host string) (string, error) {
	// overlyCautiousIllegalFileCharacters matches characters that *might* not be supported.  Windows is really restrictive, so this is really restrictive
	overlyCautiousIllegalFileCharacters, err := regexp.Compile(`[^(\w/\.)]`)
	if err != nil {
		return "", errors.Wrap(err, "while compiling regex for computing cache dir")
	}

	// strip the optional scheme from host if its there:
	schemelessHost := strings.Replace(strings.Replace(host, "https://", "", 1), "http://", "", 1)
	// now do a simple collapse of non-AZ09 characters.  Collisions are possible but unlikely.  Even if we do collide the problem is short lived
	safeHost := overlyCautiousIllegalFileCharacters.ReplaceAllString(schemelessHost, "_")
	return filepath.Join(parentDir, safeHost), nil
}
