# Release process

The release of the codeowners-validator tool is performed by the [GoReleaser](https://github.com/goreleaser/goreleaser) which builds Go binaries for several platforms and then creates a GitHub release.

**Process**

1. Export GITHUB_TOKEN=`YOUR_GH_TOKEN`

2. Tag commit 
   ```bash
   git tag -a v0.1.0 -m "First release"
   ```         

3. Push tag
    ```
    git push origin v0.1.0
    ```

4. Locally from the root of the repository, run `goreleaser`.
   >**NOTE:** Currently, releases are made with goreleaser in version `0.104.0, commit 7c4352147b6d9636f13d2fc633cfab05d82d929c, built at 2019-03-20T02:18:40Z`
   
5. Recheck release generated on GitHub. 
