# Migration Guide

This guide helps you upgrade between SOM versions.

## General Upgrade Steps

1. Review the [CHANGELOG](https://github.com/go-surreal/som/blob/main/CHANGELOG.md) for breaking changes
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
go run github.com/go-surreal/som@latest -i ./model

# Update dependencies
go mod tidy
```

## Version-Specific Guides

### Upgrading to v0.17.x

**Geo types**: Added support for geometry types from `github.com/twpayne/go-geom`.

**SemVer type**: Added `som.SemVer` type with query filters (`Major()`, `Minor()`, `Patch()`).

**Field name override**: Use `som:"custom_name"` tag to override database field names.

**Insert method**: Repos now have an `Insert()` method for bulk record creation.

**Index rebuild**: `RebuildIndexes()` has been replaced with `Index().<IndexName>().Rebuild(ctx)`.

**Driver version**: Requires surrealdb.go client version 1.3.0.

### Upgrading to v0.16.x

**Complex IDs**: Added support for array and object-based IDs via `som.ArrayID` and `som.ObjectID` marker structs.

**CLI flags**: Overhauled CLI command and flags.

### Upgrading to v0.15.x

**SurrealDB v3.0**: Added compatibility with SurrealDB v3.0.

**String record IDs**: Added configurable string ID generator support.

### Upgrading to v0.14.x

**Generic Node**: `som.Node` is now `som.Node[T]` where `T` specifies the ID type.

**Before (v0.13.x)**:
```go
type User struct {
    som.Node
    Name string
}
```

**After (v0.14.x+)**:
```go
type User struct {
    som.Node[som.ULID]
    Name string
}
```

Available ID types: `som.ULID`, `som.UUID`, `som.Rand`.

**Complex IDs**: For array or object-based IDs, create a key struct:
```go
type WeatherKey struct {
    som.ArrayID
    City string
    Date time.Time
}

type Weather struct {
    som.Node[WeatherKey]
    Temperature float64
}
```

**Model imports**: Models now import from the generated package rather than the SOM module:
```go
// Before
import "github.com/go-surreal/som"

// After
import "yourproject/gen/som"
```

### Upgrading to v0.13.x

**Sorting**: The `order` package was renamed to `by`, and `OrderBy()` was renamed to `Order()`:

```go
// Before
import "yourproject/gen/som/order"
query.OrderBy(order.User.Name.Asc())

// After
import "yourproject/gen/som/by"
query.Order(by.User.Name.Asc())
```

**Pagination**: `Offset()` was renamed to `Start()`:

```go
// Before
query.Limit(10).Offset(20).All(ctx)

// After
query.Limit(10).Start(20).All(ctx)
```

**First()**: Changed from `(*M, bool, error)` to `(*M, error)`. Returns `ErrNotFound` instead of `exists` bool:

```go
// Before
user, exists, err := query.First(ctx)
if !exists { ... }

// After
user, err := query.First(ctx)
if errors.Is(err, som.ErrNotFound) { ... }
```

### Upgrading to v0.12.x

**CLI syntax change**: The CLI no longer uses positional arguments.

```bash
# Before
go run github.com/go-surreal/som/cmd/som@latest gen ./model ./gen/som

# After
go run github.com/go-surreal/som@latest -i ./model
```

The package path also changed from `github.com/go-surreal/som/cmd/som` to `github.com/go-surreal/som`.

**Output directory**: Default output is now `gen/som`. Use `-o` flag to override.

### Upgrading to v0.11.x

**Soft Delete**: Added `som.SoftDelete` embedding with `Delete()`, `Restore()`, `Erase()`, and `WithDeleted()`.

**Optimistic Locking**: Added `som.OptimisticLock` embedding with version tracking and `ErrOptimisticLock`.

**Password type**: `som.Password` is now generic: `som.Password[som.Bcrypt]`. Supported algorithms: `Bcrypt`, `Argon2`, `Pbkdf2`, `Scrypt`.

### Upgrading to v0.10.x

No specific migration steps required. Regenerate code after upgrade.

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
