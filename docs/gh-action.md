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
        uses: mszostok/codeowners-validator@v0.5.0
        with:
          checks: "files,owners,duppatterns"
          experimental_checks: "notowned"
          # GitHub access token is required only if the `owners` check is enabled
          github_access_token: "${{ secrets.OWNERS_VALIDATOR_GITHUB_SECRET }}"
```

The best is to run this as a cron job and not only if you applying changes to CODEOWNERS file itself, e.g. the CODEOWNERS file can be invalidate when you removing someone from the organization.

> **Note**: To execute `owners` check you need to create a [GitHub token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/#creating-a-token) and store it as a secret in your repository, see ["Creating and storing encrypted secrets."](https://help.github.com/en/actions/configuring-and-managing-workflows/creating-and-storing-encrypted-secrets). Token requires only read-only scope for your repository.

<!--- example repository when failed -->