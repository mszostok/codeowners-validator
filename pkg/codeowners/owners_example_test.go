package codeowners_test

import (
	"fmt"

	"go.szostok.io/codeowners/pkg/codeowners"
)

func ExampleNewFromPath() {
	pathToCodeownersFile := "./testdata/"

	entries, err := codeowners.NewFromPath(pathToCodeownersFile)
	if err != nil {
		panic(err)
	}

	for _, e := range entries {
		fmt.Printf("[line] %d: [pattern]: %s [owners]: %v\n", e.LineNo, e.Pattern, e.Owners)
	}

	// Output:
	// [line] 8: [pattern]: * [owners]: [@global-owner1 @global-owner2]
	// [line] 14: [pattern]: *.js [owners]: [@js-owner]
	// [line] 19: [pattern]: *.go [owners]: [docs@example.com]
	// [line] 24: [pattern]: /build/logs/ [owners]: [@doctocat]
	// [line] 29: [pattern]: docs/* [owners]: [docs@example.com]
	// [line] 33: [pattern]: apps/ [owners]: [@octocat]
	// [line] 37: [pattern]: /docs/ [owners]: [@doctocat]
}
