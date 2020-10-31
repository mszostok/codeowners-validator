<h1>
    <img alt="logo" src="./docs/assets/logo-small.png" width="28px" />
    Codeowners Validator - contributing & development
</h1>

This document contains contribution guidelines and tips for this repository. Read it to learn how to develop the Codeowners Validator project.

## Table of Contents

<!-- toc -->

<!-- tocstop -->

## Prerequisites

* [Go](https://golang.org/dl/) 1.15 or higher
* [Docker](https://www.docker.com/)
* Make

Helper scripts may introduce additional dependencies. However, all helper scripts support the `INSTALL_DEPS` environment variable flag.
By default, this flag is set to `false`. This way, the scripts will try to use the tools installed on your local machine. This helps speed up the development process.
If you do not want to install any additional tools, or you want to ensure reproducible script 
results, export `INSTALL_DEPS=true`. This way, the proper tool version will be automatically installed and used. 

## Contributing

### Issues

> **NOTE:** Before adding a new issue, check the [existing issues](https://github.com/mszostok/codeowners-validator/issues) to avoid duplicates. 

If you'd like to see a certain feature added to the project, use this [feature request template](./.github/ISSUE_TEMPLATE/feature_request.md) to create an issue. 

Similarly, if you spot a bug, use this [bug report template](./.github/ISSUE_TEMPLATE/bug_report.md) to let us know!

### PRs

> **NOTE:** Before proposing changes, check the [existing issues](https://github.com/mszostok/codeowners-validator/issues) and make sure the work has not already been started to avoid duplication. 

To start contributing, follow these steps: 

1. Fork the `codeowners-validator` repository.
2. Clone the repository locally. 
    > **TIP:** This project uses Go modules, so you can check it out locally wherever you want. It doesn't need to be checked out in `$GOPATH`.
3. Set the `codeowners-validator` repository as upstream:
    ```bash
    git remote add upstream git@github.com:mszostok/codeowners-validator.git
    ```
4. Fetch all the remote branches for this repository:
    ```bash
    git fetch --all 
    ```
5. Set the master branch to point to upstream:
    ```bash
    git branch -u upstream/master master
    ```

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

This project supports the integration tests that are defined in the [test](./test) package. The tests are executed against [`gh-codeowners/codeowners-samples`](https://github.com/gh-codeowners/codeowners-samples).

> **CAUTION:** Currently, running the integration tests both on external PRs and locally by external contributors is not supported, as the teams used for testing are visible only to the organization members. 
> At the moment, the `codeowners-validator` repository owner is responsible for running these tests. 

## Build a binary

To generate a binary for this project, execute:
```bash
make build
```

This command generates a binary named `codeowners-validator` in the root directory.

[↑ Back to top](#table-of-contents)