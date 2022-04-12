package version

import (
	"fmt"
	"io"
	"runtime"

	flag "github.com/spf13/pflag"
)

var (
	printVer = flag.BoolP("version", "v", false, "Prints current version.")
	short    = flag.Bool("short", false, "Print just the version number.")
)

// Info contains versioning information.
type Info struct {
	Version   string
	GitCommit string
	BuildDate string
	GoVersion string
	Compiler  string
	Platform  string
}

func Init() {
	flag.Parse()
}

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() *Info {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the settings in ./base.go
	return &Info{
		Version:   version,
		GitCommit: commit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func ShouldPrintVersion() bool {
	return *printVer
}

func PrintVersion(out io.Writer) {
	if *short {
		fmt.Fprintf(out, "Version: %s\n", Get())
	} else {
		fmt.Fprintf(out, "%#v\n", Get())
	}
}

// String returns info as a human-friendly version string.
func (info *Info) String() string {
	return info.Version
}
