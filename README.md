<p align="center">
  <h3 align="center">GitHub CODEOWNERS Validator</h3>
  <p align="center">Ensures the correctness of your CODEOWNERS file.</p>
  <p align="center">
    <a href="/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-Apache-brightgreen.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/goreleaser/godownloader"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/mszostok/codeowners-validator?style=flat-square"></a>
    <a href="https://travis-ci.org/goreleaser/godownloader"><img alt="Travis" src="https://img.shields.io/travis/com/mszostok/codeowners-validator/master.svg?style=flat-square"></a>
    <!-- <a href="http://godoc.org/github.com/mszostok/codeowners-validator"><img alt="Go Doc" src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square"></a> --> 
  </p>
</p>

---

The Codeowners Validator project validates the GitHub [CODEOWNERS](https://help.github.com/articles/about-code-owners/) file using [those checks](#checks). It supports public and private GitHub repositories and also GitHub Enterprise installations.

![usage](./docs/assets/usage.png)

## Installation

It's highly recommended to install a fixed version of ` codeowners-validator`. Releases are available on the [releases page](https://github.com/mszostok/codeowners-validator/releases).

#### From Release

Here is the recommended way to install `codeowners-validator`:

```bash
# binary installed into ./bin/
curl -sfL https://raw.githubusercontent.com/mszostok/codeowners-validator/master/install.sh| sh -s v0.2.0

# binary installed into $(go env GOPATH)/bin/codeowners-validator
curl -sfL https://raw.githubusercontent.com/mszostok/codeowners-validator/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v0.2.0

# In alpine linux (as it does not come with curl by default)
wget -O - -q https://raw.githubusercontent.com/mszostok/codeowners-validator/master/install.sh| sh -s v0.2.0

# Print version. Add `--short` to print just the version number.
codeowners-validator -v
```

You can also download [latest version](https://github.com/mszostok/codeowners-validator/releases/latest) from release page manually.

#### From Sources

You can install `codeowners-validator` with `env GO111MODULE=off go get -u github.com/mszostok/codeowners-validator`.

> NOTE: please use the latest go to do this, ideally go 1.12 or greater.

This will put `codeowners-validator` in `$(go env GOPATH)/bin`

## Checks

The following checks are enabled by default:

| Name       | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
|------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| duppatterns | **[Duplicated Pattern Checker]** <br /><br /> Reports if CODEOWNERS file contain duplicated lines with the same file pattern.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| files      | **[File Exist Checker]** <br /><br /> Reports if CODEOWNERS file contain lines with the file pattern that do not exist in a given repository.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| owners     | **[Valid Owner Checker]** <br /><br /> Reports if CODEOWNERS file contain invalid owners definition. Allowed owner syntax: `@username`, `@org/team-name` or `user@example.com` <br /> _source: https://help.github.com/articles/about-code-owners/#codeowners-syntax_. <br /> <br /> **Checks:** <br /> &#x09; 1. Check if the owner's definition is valid (is either a GitHub user name, an organization team name or an email address). <br /><br /> 2. Check if a GitHub owner has a GitHub account <br /><br /> 3. Check if a GitHub owner is in a given organization <br /> <br />4. Check if an organization team exists |

The experimental checks are disabled by default:

| Name     | Description                                                                                                                                 |
|----------|---------------------------------------------------------------------------------------------------------------------------------------------|
| notowned | **[Not Owned File Checker]** <br /><br /> Reports if a given repository contain files that do not have specified owners in CODEOWNERS file. |

To enable experimental check set `EXPERIMENTAL_CHECKS=notowned` environment variable. 

Check the [Usage](#usage) section for more info on how to enable and configure given checks.

## Usage

Use the following environment variables to configure the application:

| Name | Default | Description |
|-----|:--------|:------------|
| <tt>REPOSITORY_PATH</tt> <b>*</b> | | The repository path to your repository on your local machine. |
| <tt>GITHUB_ACCESS_TOKEN</tt>| | The GitHub access token. Instruction for creating a token can be found [here](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/#creating-a-token). If not provided then validating owners functionality could not work properly, e.g. you can reach the API calls quota or if you are setting GitHub Enterprise base URL then an unauthorized error can occur. |
| <tt>GITHUB_BASE_URL</tt>| https://api.github.com/ | The GitHub base URL for API requests. Defaults to the public GitHub API, but can be set to a domain endpoint to use with GitHub Enterprise. |
| <tt>GITHUB_UPLOAD_URL</tt> | https://uploads.github.com/ | The GitHub upload URL for uploading files. <br> <br>It is taken into account only when the `GITHUB_BASE_URL` is also set. If only the `GITHUB_BASE_URL` is provided then this parameter defaults to the `GITHUB_BASE_URL` value. |
| <tt>CHECKS</tt>| - |  The list of checks that will be executed. By default, all checks are executed. Possible values: `files`,`owners`,`duppatterns` |
| <tt>EXPERIMENTAL_CHECKS</tt> | - | The comma-separated list of experimental checks that should be executed. By default, all experimental checks are turned off. Possible values: `notowned`.|
| <tt>CHECK_FAILURE_LEVEL</tt> | `warning` | Defines the level on which the application should treat check issues as failures. Defaults to `warning`, which treats both errors and warnings as failures, and exits with error code 3. Possible values are `error` and `warning`. |
| <tt>OWNER_CHECKER_ORGANIZATION_NAME</tt>  <b>*</b>| | The organization name where the repository is created. Used to check if GitHub owner is in the given organization. |
| <tt>NOT_OWNED_CHECKER_SKIP_PATTERNS</tt>| - | The comma-separated list of patterns that should be ignored by `not-owned-checker`. For example, you can specify `*` and as a result, the `*` pattern from the **CODEOWNERS** file will be ignored and files owned by this pattern will be reported as unowned unless a later specific pattern will match that path. It's useful because often we have default owners entry at the begging of the CODOEWNERS file, e.g. `*       @global-owner1 @global-owner2` |

 <b>*</b> - Required

#### Exit status codes

Application exits with different status codes which allow you to easily distinguish between error categories.  

| Code | Description |
|:-----:|:------------|
| **1** | The application startup failed due to the wrong configuration or internal error. |
| **2** | The application was closed because the OS sends a termination signal (SIGINT or SIGTERM). |
| **3** | The CODEOWNERS validation failed - executed checks found some issues. |

## Roadmap

The [codeowners-validator roadmap uses Github milestones](https://github.com/mszostok/codeowners-validator/milestone/1) to track the progress of the project.

They are sorted with priority. First are most important.
