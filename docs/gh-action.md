[← back to docs](./README.md)

<p align="center">
  <h3 align="center">GitHub Action for CODEOWNERS Validator</h3>
  <p align="center">Ensures the correctness of your CODEOWNERS file.</p>
  <p align="center">
    <a href="/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-Apache-brightgreen.svg?style=flat-square"></a>
  </p>
</p>

##
The [Codeowners Validator](https://github.com/mszostok/codeowners-validator) is available as a GitHub Action.
                                                                                   
<p align="center">
    <img src="https://raw.githack.com/mszostok/codeowners-validator/master/docs/assets/action-output.png" width="600px" alt="demo">
</p>


## Usage
 
Create a workflow (eg: `.github/workflows/sanity.yml` see [Creating a Workflow file](https://help.github.com/en/articles/configuring-a-workflow#creating-a-workflow-file))

```yaml
name: "Codeowners Validator"

on:
  schedule:
    # Runs at 08:00 UTC every day
    - cron:  '0 8 * * *'

jobs:
  sanity:
    runs-on: ubuntu-latest
    steps:
      # Checks-out your repository, which is validated in the next step
      - uses: actions/checkout@v2
      - name: GitHub CODEOWNERS Validator
        uses: mszostok/codeowners-validator@v0.6.0
        # input parameters
        with:
          # "The list of checks that will be executed. By default, all checks are executed. Possible values: files,owners,duppatterns,syntax"
          checks: "files,owners,duppatterns,syntax"

          # "The comma-separated list of experimental checks that should be executed. By default, all experimental checks are turned off. Possible values: notowned."
          experimental_checks: "notowned"

          # GitHub access token is required only if the `owners` check is enabled
          github_access_token: "${{ secrets.OWNERS_VALIDATOR_GITHUB_SECRET }}"
          
          # The GitHub base URL for API requests. Defaults to the public GitHub API, but can be set to a domain endpoint to use with GitHub Enterprise.
          github_base_url: "https://api.github.com/"

          # The GitHub upload URL for uploading files. It is taken into account only when the GITHUB_BASE_URL is also set. If only the GITHUB_BASE_URL is provided then this parameter defaults to the GITHUB_BASE_URL value.
          github_upload_url: "https://uploads.github.com/"
        
          # The repository path in which CODEOWNERS file should be validated."
          repository_path: "."
        
          # Defines the level on which the application should treat check issues as failures. Defaults to warning, which treats both errors and warnings as failures, and exits with error code 3. Possible values are error and warning. Default: warning"
          check_failure_level: "warning"
        
          # The comma-separated list of patterns that should be ignored by not-owned-checker. For example, you can specify * and as a result, the * pattern from the CODEOWNERS file will be ignored and files owned by this pattern will be reported as unowned unless a later specific pattern will match that path. It's useful because often we have default owners entry at the begging of the CODOEWNERS file, e.g. * @global-owner1 @global-owner2"
          not_owned_checker_skip_patterns: ""
        
          # The owner and repository name. For example, gh-codeowners/codeowners-samples. Used to check if GitHub team is in the given organization and has permission to the given repository."
          owner_checker_repository: "${{ github.repository }}"
        
          # The comma-separated list of owners that should not be validated. Example: @owner1,@owner2,@org/team1,example@email.com."
          owner_checker_ignored_owners: "@ghost"
```

The best is to run this as a cron job and not only if you applying changes to CODEOWNERS file itself, e.g. the CODEOWNERS file can be invalidate when you removing someone from the organization.

> **Note**: To execute `owners` check you need to create a [GitHub token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/#creating-a-token) and store it as a secret in your repository, see ["Creating and storing encrypted secrets."](https://help.github.com/en/actions/configuring-and-managing-workflows/creating-and-storing-encrypted-secrets). Token requires only read-only scope for your repository.

<!--- example repository when failed -->

## Configuration

For the GitHub Action, use the configuration described in the main README under the [Configuration](../README.md#configuration) section but **specify it as the [Action input parameters](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idstepswith) instead of environment variables**. See the [Usage](#usage) section for the full syntax. 

If you want to use environment variables anyway, you must add the `INPUT_` prefix to each environment variable. For example, `OWNER_CHECKER_IGNORED_OWNERS` becomes `INPUT_OWNER_CHECKER_IGNORED_OWNERS`.
