[← back to docs](./README.md)

# Development

This document contains development instructions. Read it to learn how to develop this project.

# Table of Contents

<!-- toc -->

- [Prerequisites](#prerequisites)
- [Dependency management](#dependency-management)
- [Testing](#testing)
  * [Unit tests](#unit-tests)
  * [Lint tests](#lint-tests)
  * [Integration tests](#integration-tests)
- [Build a binary](#build-a-binary)

<!-- tocstop -->

## Prerequisites

* [Go](https://golang.org/dl/) 1.15 or higher
* [Docker](https://www.docker.com/)
* Make

Helper scripts may introduce additional dependencies. However, all helper scripts support the `INSTALL_DEPS` environment variable flag.
By default, this flag is set to `false`. This way, the scripts will try to use the tools installed on your local machine. This helps speed up the development process.
If you do not want to install any additional tools, or you want to ensure reproducible script 
results, export `INSTALL_DEPS=true`. This way, the proper tool version will be automatically installed and used. 

## Dependency management

This project uses `go modules` for dependency management. To install all required dependencies, use the following command:

```bash
go mod download
```

## Testing

### Unit tests

To run all unit tests, execute:

```bash
make test-unit
```

To generate the unit test coverage HTML report, execute: 

```bash
make test-unit-cover-html
```

> **NOTE:** The generated report opens automatically in your default browser.

### Lint tests

To check your code for errors, such as typos, wrong formatting, security issues, etc., execute:

```bash
make test-lint
```

To automatically fix detected lint issues, execute:

```bash
make fix-lint-issues
```

### Integration tests

This project supports the integration tests that are defined in the [tests](../tests) package. The tests are executed against [`gh-codeowners/codeowners-samples`](https://github.com/gh-codeowners/codeowners-samples).

> **CAUTION:** Currently, running the integration tests both on external PRs and locally by external contributors is not supported, as the teams used for testing are visible only to the organization members. 
> At the moment, the `codeowners` repository owner is responsible for running these tests. 

## Build a binary

To generate a binary for this project, execute:
```bash
make build
```

This command generates a binary named `codeowners` in the root directory.

[↑ Back to top](#table-of-contents)
