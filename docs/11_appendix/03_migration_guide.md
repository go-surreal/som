# Migration Guide

This guide helps you upgrade between SOM versions.

## General Upgrade Steps

1. Review the [Changelog](../10_development/03_changelog.md) for breaking changes
2. Update your SOM version
3. Regenerate all code
4. Fix any compilation errors
5. Run your test suite

## Regenerating Code

After any SOM upgrade:

```bash
# Delete old generated code
rm -rf ./gen/som

# Regenerate with new version
go run github.com/go-surreal/som/cmd/som@latest gen ./model ./gen/som

# Update dependencies
go mod tidy
```

## Version-Specific Guides

### Upgrading to v0.10.x

No specific migration steps required. Regenerate code after upgrade.

### Future Versions

Migration guides will be added here as new versions are released with breaking changes.

## Common Migration Issues

### Changed Method Signatures

If method signatures change between versions:

1. Regenerate code
2. Update call sites to match new signatures
3. IDE "Find Usages" can help locate all call sites

### Renamed Packages

If generated package names change:

1. Update all imports in your code
2. Use IDE refactoring tools for bulk updates

### Removed Features

If a feature is removed:

1. Check the changelog for replacement approaches
2. Open an issue if no migration path exists

## Database Schema Changes

SOM doesn't currently handle database schema migrations. For schema changes:

1. Back up your data
2. Update your model structs
3. Regenerate SOM code
4. Apply schema changes to SurrealDB manually

## Getting Help

If you encounter migration issues not covered here:

1. Check [GitHub Issues](https://github.com/go-surreal/som/issues) for known problems
2. Open a new issue describing your migration scenario
