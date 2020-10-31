<h1>
    <img alt="logo" src="./docs/assets/logo-small.png" width="28px" />
    Codeowners Validator - contributing
</h1>

🎉🚀🤘 Thanks for your interest in the Codeowners Validator project!🤘🚀🎉 

This document contains contribution guidelines for this repository. Read it before you start contributing.

## Contributing

Before proposing or adding changes, check the [existing issues](https://github.com/mszostok/codeowners-validator/issues) and make sure the discussion/work has not already been started to avoid duplication. 

If you'd like to see a new feature implemented, use this [feature request template](./.github/ISSUE_TEMPLATE/feature_request.md) to create an issue. 

Similarly, if you spot a bug, use this [bug report template](./.github/ISSUE_TEMPLATE/bug_report.md) to let us know!

### Ready for action? Start developing! 

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

You're all set! 🚀 Read the [development](./docs/development.md) document for further instructions.
