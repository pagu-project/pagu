# Update Dependencies

This document is about how to update the Pagu project repository to latest version.

### Packages

First of all you need to update golang dependencies to latest version using this commands:

```sh
go get -u ./...
go mod tidy
```

Once all packages got updated, make sure you run `make build` and `make test` commands to
make sure none of previous behaviors are broken.
If any packages had breaking changes or some of them are deprecated,
you need to update the code and use new methods or use another package.

### Go version

You have to update the go version to latest release in [go.mod](../go.mod).
Make sure you are updating version of Golang on [Dockerfile](../deployment/Dockerfile).

> Note: you must run `make build` after this change to make sure everything works smoothly.

### Example Pull Request

Here is an example pull request to find out what you need to update and how to set commit message:
https://github.com/pagu-project/pagu/pull/314
