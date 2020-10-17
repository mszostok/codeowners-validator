## File exits checker

This document describes investigation about [`file exists`](../../../internal/check/file_exists.go) checker which needs to deal with the gitignore pattern syntax 

### Problem

A [CODEOWNERS](https://docs.github.com/en/free-pro-team@latest/github/creating-cloning-and-archiving-repositories/about-code-owners#codeowners-syntax) file uses a pattern that follows the same rules used in [gitignore](https://git-scm.com/docs/gitignore#_pattern_format) files.
The gitignore files support two consecutive asterisks ("**") in patterns that match against the full path name. Unfortunately the core Go library `filepath.Glob` does not support [`**`](https://github.com/golang/go/issues/11862) at all.

This caused that for some patterns the [`file exists`](../../../internal/check/file_exists.go) checker didn't work properly, see [issue#22](https://github.com/mszostok/codeowners-validator/issues/22). 

Additionally, we need to support a single asterisk at the beginning of the pattern. For example, `*.js` should check for all JS files in the whole git repository. To achieve that we need to detect that and change from `*.js` to `**/*.js`.

```go
pattern := "*.js"
if len(pattern) >= 2 && pattern[:1] == "*" && pattern[1:2] != "*" {
		pattern = "**/" + pattern
}
```

### Investigation

Instead of creating a dedicated solution, I decided to search for a custom library that's supporting two consecutive asterisks.
There are a few libraries in open-source that can be used for that purpose. I selected three:
- https://github.com/bmatcuk/doublestar/v2
- https://github.com/mattn/go-zglob
- https://github.com/yargevad/filepathx  

I've tested all libraries and all of them were supporting `**` pattern properly. As a final criterion, I created benchmark tests.

#### Benchmarks
 
Run benchmarks with 1 CPU for 5 seconds:

```bash
go test  -bench=. -benchmem -cpu 1 -benchtime 5s ./file_matcher_libs_bench_test.go

goos: darwin
goarch: amd64
BenchmarkPathx                79          72276938 ns/op         7297258 B/op      40808 allocs/op
BenchmarkZGlob               126          47206545 ns/op          840973 B/op      10550 allocs/op
BenchmarkDoublestar          157          38041578 ns/op         3521379 B/op      22150 allocs/op
```

Run benchmarks with 12 CPU for 5 seconds:
```bash
go test  -bench=. -benchmem -cpu 12 -benchtime 5s ./file_matcher_libs_bench_test.go

goos: darwin
goarch: amd64
BenchmarkPathx-12                     78          73096386 ns/op         7297114 B/op      40807 allocs/op
BenchmarkZGlob-12                    637           9234632 ns/op          914239 B/op      10564 allocs/op
BenchmarkDoublestar-12               151          38372922 ns/op         3522899 B/op      22151 allocs/op
```

#### Summary

With the 1 CPU , the `doublestar` library has the shortest time, but the allocated memory is higher than the `z-glob` library.
With the 12 CPU, the `z-glob` is a winner bot in time and memory allocation. The worst one in each test was the `filepathx` library. 

> **NOTE:** The `z-glob` library has an issue with error handling. I've provided PR for fixing that problem: https://github.com/mattn/go-zglob/pull/37.