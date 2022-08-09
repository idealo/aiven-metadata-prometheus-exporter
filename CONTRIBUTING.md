# Contribution Guide

We welcome any contributions whether it is:

- Submitting feedback
- Fixing bugs
- Or implementing a new feature.

Please read this guide before making any contributions.

#### Commit Messages
We utilize [conventional commits](https://www.conventionalcommits.org) for writing proper commit messages.
Please stick to the conventional commit guide when you write your commit message.

#### Submit Feedback
The feedback should be submitted by creating an issue on [GitHub issues](https://github.com/idealo/aiven-metadata-prometheus-exporter/issues).
Select the related template (bug report, feature request or custom) and add the corresponding labels.
 
If you have any more generic-like requests, feel free to open a [discussion](https://github.com/idealo/aiven-metadata-prometheus-exporter/discussions).

#### Fix Bugs
If you want to fix a bug, you may look through the [GitHub issues](https://github.com/idealo/aiven-metadata-prometheus-exporter/issues) for bugs.

#### Implement Features
If you have any ideas or suggestions, please open an issue and describe your idea or feature request.
If you want to implement a specific feature, you may look through the [GitHub issues](https://github.com/idealo/aiven-metadata-prometheus-exporter/issues) for feature requests.

## Pull Requests (PR)
1. Fork the repository and create a new branch based on the `main` branch.
2. For bug fixes and features, please add new tests and add according changes to the documentation, if needed.
3. Do a PR from your new branch to our `main` branch of the [original repo](https://github.com/idealo/aiven-metadata-prometheus-exporter).

## Documentation
- Make sure any new feature has a proper documentation or is described in the [README.md](README.md).

## Testing
We use the standard golang library for testing. Make sure to write tests for any new feature and/or bug fixes.

If you want to contribute to the project, please make sure that
  * changes compile
  * and tests are green

          $ go build -o bin/aiven-metadata-prometheus-exporter
          $ go test -v ./...

  * the changes are reflected on the `/metrics` endpoint

          $ export AIVEN_API_TOKEN=MyToken; bin/aiven-metadata-prometheus-exporter
          $ curl -s localhost:2112/metrics
