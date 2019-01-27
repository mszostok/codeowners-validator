# Codeowners Validator

## Overview

The Codeowners Validator project validates the GitHub [CODEOWNERS](https://help.github.com/articles/about-code-owners/) file.

## Local Installation

`go get -u github.com/mszostok/codeowners-validator`

## Usage

Use the following environment variables to configure the application:

| Name | Required | Default | Description |
|-----|---------|--------|------------|
| **REPOSITORY_PATH** | Yes | - | The repository path to your repository on your local machine. |
| **GITHUB_ACCESS_TOKEN** | Yes | - | The GitHub access token. Instruction for creating token can be found [here](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/#creating-a-token)|
| **VALID_OWNER_CHECKER_ORGANIZATION_NAME** | Yes | - | The organization name where the repository is created. Used to check if GitHub owner is in the given organization. |

## Roadmap

_Sorted with priority. First - most important._

* [ ] Possibility to execute validator online. Automatically integrates with your GitHub account and allows you to check any repository online without the need to download and execute binary locally.
* [ ] Possibility to use the GitHub URL instead of the path to the local repository.
* [ ] Offline mode - execute all checks which not require internet connection against your local repository
* [ ] Investigate the [Go Plugins](https://golang.org/pkg/plugin/). Implement if it will simplify extending this tool with other checks.
* [ ] Move to [cobra](https://github.com/spf13/cobra/) library.
* [ ] Add test coverage.
* [ ] Add support for configuration via YAML file.
* [ ] Move dep to go modules 


Checks:

* [ ] Unowned files (Not defined owners for given files) 
* [ ] Find doubles paths
* [x] File/directory does not exits
* [x] Works with teams as well
* [x] Support for private repos (see below)
* [x] User exist or not
* [x] Validate owners
  * [x] check if owner definition is valid (GitHub user, organization team, email address)
  * [x] check if GitHub owner have GitHub account
  * [x] check if GitHub owner is in the given organization
  * [x] check if Organization team exist