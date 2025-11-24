# Contributing

We welcome contributions to SOM! This guide explains how to get involved.

## Getting Started

1. Fork the repository to your own namespace
2. Clone your fork locally
3. Create a feature branch for your changes

## Development Setup

```bash
# Clone the repository
git clone https://github.com/YOUR_USERNAME/som.git
cd som

# Install dependencies
go mod download

# Run tests
go test ./...
```

## Making Changes

1. Write your code
2. Add or update tests as needed
3. Ensure all tests pass
4. Commit your changes

## Commit Messages

Commit messages follow the [Conventional Commits](https://www.conventionalcommits.org) specification:

```
feat: add new filter operation
fix: correct edge traversal bug
docs: update installation guide
refactor: simplify query builder
test: add repository tests
```

During development, commit messages can be informal. The final PR will be squash merged with a properly formatted message.

## Pull Requests

1. Push your branch to your fork
2. Create a pull request to the main repository
3. Ensure the PR title follows Conventional Commits format
4. Wait for review and address any feedback

## Labels

PRs use labels that correspond to commit types:
- `feat` - New features
- `fix` - Bug fixes
- `docs` - Documentation changes
- `refactor` - Code refactoring
- `test` - Test additions/changes

## Code Style

- Follow standard Go conventions
- Run `gofmt` and `golangci-lint`
- Write clear, self-documenting code
- Add comments only where logic isn't obvious

## Questions?

Open a [GitHub Discussion](https://github.com/go-surreal/som/discussions) for questions or ideas.
