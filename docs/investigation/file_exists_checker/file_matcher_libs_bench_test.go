// Always record the result of func execution to prevent
// the compiler eliminating the function call.
// Always store the result to a package level variable
// so the compiler cannot eliminate the Benchmark itself.
package file_exists_checker

import (
	"fmt"
	"log"
	"os"
	"path"
	"testing"

	"github.com/bmatcuk/doublestar/v2"
	"github.com/mattn/go-zglob"
	"github.com/yargevad/filepathx"
)

var pattern string
func init() {
	curDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	pattern =  path.Join(curDir, "..", "..", "**", "*.md")
	fmt.Println(pattern)
}

var pathx []string

func BenchmarkPathx(b *testing.B) {
	var r []string
	for n := 0; n < b.N; n++ {
		r, _ = filepathx.Glob(pattern)
	}
	pathx = r
}

var zGlob []string

func BenchmarkZGlob(b *testing.B) {
	var r []string
	for n := 0; n < b.N; n++ {
		r, _ = zglob.Glob(pattern)
	}
	zGlob = r
}

var double []string

func BenchmarkDoublestar(b *testing.B) {
	var r []string
	for n := 0; n < b.N; n++ {
		r, _ = doublestar.Glob(pattern)
	}
	double = r
}
