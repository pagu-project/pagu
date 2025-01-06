## Pagu Deployment

This project includes an automated deployment process for both the `stable` and `latest` versions of the Pagu bots.

### First-Time Setup

#### Database User Setup

To grant the correct privileges to a single database user, execute the following SQL commands.
Make sure to replace `<MYSQL_USER>` and `<MYSQL_PASSWORD>` with the appropriate values for your setup.

```sql
CREATE DATABASE IF NOT EXISTS pagu;
CREATE DATABASE IF NOT EXISTS pagu_staging;

-- Ensure the user exists
CREATE USER IF NOT EXISTS '<MYSQL_USER>'@'%' IDENTIFIED BY '<MYSQL_PASSWORD>';

-- Grant privileges to the user on both databases
GRANT ALL PRIVILEGES ON pagu.* TO '<MYSQL_USER>'@'%';
GRANT ALL PRIVILEGES ON pagu_staging.* TO '<MYSQL_USER>'@'%';

-- Apply the changes
FLUSH PRIVILEGES;
```

#### Docker Network Setup

To enable Docker containers to communicate with each other on the same network, you need to create an external network and share it between the containers.
Use the following command to create the network:

```bash
docker network create pagu_network
```

### Deployment Overview

The deployment system operates as follows:

- **Latest Version**: Triggered when changes are pushed to the `main` branch.
- **Stable Version**: Triggered when a new stable version is released.

### Releasing a Stable Version

To release and deploy a *Stable* version, create a Git tag and push it to the repository. Follow these steps:

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

