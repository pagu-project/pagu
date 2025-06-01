# Pagu Deployment

This project includes an automated deployment process for both the `stable` and `latest` versions of the Pagu bots.

## First-Time Setup

### Database User Setup

To grant the correct privileges to a database user, execute the following SQL commands.
Replace `<MYSQL_USER>`, `<USER_PASSWORD>`, `<MYSQL_READONLY>`, and `<READONLY_PASSWORD>`
with appropriate values for your setup.

The read-only user is granted permission to read the database but has no write or update privileges.

```sql
CREATE DATABASE IF NOT EXISTS pagu;
CREATE DATABASE IF NOT EXISTS pagu_staging;

-- Ensure the user exists
CREATE USER IF NOT EXISTS '<MYSQL_USER>'@'%' IDENTIFIED BY '<USER_PASSWORD>';
CREATE USER IF NOT EXISTS '<MYSQL_READONLY>'@'%' IDENTIFIED BY '<READONLY_PASSWORD>';

-- Grant all privileges to the main user on both databases.
GRANT ALL PRIVILEGES ON pagu.* TO '<MYSQL_USER>'@'%';
GRANT ALL PRIVILEGES ON pagu_staging.* TO '<MYSQL_USER>'@'%';

-- Grant select privilege to the read-only user.
GRANT SELECT ON pagu.* TO '<MYSQL_READONLY>'@'%';
GRANT SELECT ON pagu_staging.* TO '<MYSQL_READONLY>'@'%';

-- Apply the changes
FLUSH PRIVILEGES;
```

You can check users by:

```sql
SELECT * FROM mysql.user;
```

### Docker Network Setup

To enable Docker containers to communicate on the same network,
create an external Docker network with the following command:

```bash
docker network create pagu_network
```

### Docker To Host Port Forwarding

Copy the `docker2host.service` file to the user's systemd unit directory on the target server.
This is typically located at `~/.config/systemd/user/` on most Linux distributions.

Next, enable and start the service:

```bash
chmod +x docker2host.sh

systemctl --user enable docker2host.service
systemctl --user start  docker2host.service
```

To ensure the service continues running even after you log out, enable lingering for the user account:

```bash
sudo loginctl enable-linger
sudo loginctl enable-linger <USERNAME>
```

## Deployment Overview

The deployment system operates as follows:

- **Latest Version**: Triggered when changes are pushed to the `main` branch.
- **Stable Version**: Triggered when a new stable version is released.

## Releasing a Stable Version

To release and deploy a *Stable* version, create a Git tag and push it to the repository. Follow these steps:

1. Ensure that the `origin` remote points to the current repository, not your fork.
2. Verify that Pagu's [version](../internal/version/version.go) is updated and matches the release version.
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

Once a stable version is released, immediately update the [version.go](../internal/version/version.go) file and open a Pull Request.
For reference, you can check this [Pull Request](https://github.com/pagu-project/pagu/pull/215).

