package version

// Base version information.
//
// This is the fallback data used when version information from git is not
// provided via go ldflags.
var (
	version   = "unknown"
	commit    = "unknown"
	buildDate = "unknown"
)
