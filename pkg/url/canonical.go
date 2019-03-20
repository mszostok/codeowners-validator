package url

import (
	"strings"
)

func CanonicalPath(path string) string {
	normalized := strings.TrimRight(path, "/")
	if !strings.HasSuffix(normalized, "/") {
		normalized += "/"
	}
	return normalized
}
