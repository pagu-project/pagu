## Pagu Deployment

This project includes an automated deployment process for both the `stable` and `latest` versions of the Pagu bots.

### Deployment Overview

The deployment system operates as follows:

- **Stable Version**: Triggered when a new stable version is released.
- **Latest Version**: Triggered when changes are pushed to the `main` branch.

### Releasing a Stable Version

To release and deploy a stable version, create a Git tag and push it to the repository. Follow these steps:

1. Ensure that the `origin` remote points to the current repository, not your fork.
2. Verify that Pagu's [version](../version.go) is updated and matches the release version.
3. Run the following commands:

```bash
VERSION=0.x.y # Should match the release version
git checkout main
git pull
git tag -s -a v${VERSION} -m "Version ${VERSION}"
git push origin v${VERSION}
```

After creating the tag, the stable version will be released, and the deployment process will be triggered automatically.

### Updating the Working Version

Once a stable version is released, immediately update the [version.go](../version.go) file and open a Pull Request.  
For reference, you can check this [Pull Request](https://github.com/pagu-project/pagu/pull/215).
