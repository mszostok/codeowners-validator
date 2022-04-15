[← back to docs](./README.md)

# Github tokens

The [valid_owner.go](./../internal/check/valid_owner.go) check requires the GitHub token for the following reasons:

1. Information about organization teams and their repositories is not publicly available.
2. If you set GitHub Enterprise base URL, an unauthorized error may occur. 
3. For unauthenticated requests, the rate limit allows for up to 60 requests per hour. Unauthenticated requests are associated with the originating IP address. In a big organization where you have a lot of calls between your infrastructure server and the GitHub site, it is easy to exceed that quota. 

You can either use a personal access token or a Github App.

## GitHub personal access token

Instructions for creating a token can be found [here](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/#creating-a-token). The minimal scope required for the token is **read-only**, but the definition of this scope differs between public and private repositories.

#### Public repositories 

For public repositories, select `public_repo` and `read:org`:

![token-public.png](./assets/token-public.png) 

#### Private repositories 

For private repositories, select `repo` and `read:org`:

![token-public.png](./assets/token-private.png) 

The Codeowners Validator source code is available on GitHub. You can always perform a security audit against its code base and build your own version from the source code if your organization is more strict about the software run in its infrastructure.

## Github App

Here are the steps to create a Github App and use it for this tool:

1. [Create a GitHub App](https://docs.github.com/en/developers/apps/building-github-apps/creating-a-github-app). **Note: your app does not need a callback or a webhook URL**.
2. Add a read-only permission to the "Members" item of organization permissions.
3. [Install the app in your organization](https://docs.github.com/en/developers/apps/managing-github-apps/installing-github-apps)
4. Done! To authenticate with your app, you need three environment variables:
   1. `GITHUB_APP_PRIVATE_KEY`: PEM-format key generated when the app is installed. If you lost it, you can regenerate it ([docs](https://docs.github.com/en/developers/apps/building-github-apps/authenticating-with-github-apps#generating-a-private-key)).
   2. `GITHUB_APP_ID`: Found in the app's "About" page (Organization settings -> Developer settings -> Edit button on your app).
   3. `GITHUB_APP_INSTALLATION_ID`: Found in the URL your organization's app install page (Organization settings -> Github Apps -> Configure button on your app). It's the last number in the URL, ex: `https://github.com/organizations/my-org/settings/installations/1234567890`.
