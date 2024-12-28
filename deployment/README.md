## Auto-Deployment for Pagu Project

This project includes an automated deployment process for
both the `stable` and `latest` versions of the Pagu Discord and Telegram bots.

### Deployment Overview

The deployment system uses the following mechanisms:

- **Stable Version**: Activated when a Git tag is pushed to the repository.
- **Latest Version**: Activated when changes are pushed to the `main` branch.

### How to Create a Tag:

To create a tag and push it to the repository, follow these steps:

1. Ensure that the origin is set to the current repository, not your fork.
2. Ensure that the Pagu's [version](../version.go) is updated.
3. Run the following commands:

```bash
VERSION=0.x.y # Replace x and y with the latest version numbers
git pull origin
git tag -s -a v${VERSION} -m "Version ${VERSION}"
git push origin v${VERSION}
```
