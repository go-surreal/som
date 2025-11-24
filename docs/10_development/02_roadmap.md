# Roadmap

This page outlines the planned development direction for SOM.

## Current Status

SOM is in **early development** (pre-1.0). The API may change between minor versions.

## Planned Features

### Data Types

- [ ] `map[string]T` support
- [ ] `big.Int` / `big.Float` for arbitrary precision
- [ ] `net.IP` for IP addresses
- [ ] Geometry types for spatial data
- [ ] `som.Email` with validation
- [ ] `som.Password` with hashing
- [ ] `som.Slug` for URL-friendly identifiers

### Query Builder

- [ ] Aggregation functions (sum, avg, count, etc.)
- [ ] Group by support
- [ ] Subqueries
- [ ] Raw query escape hatch

### Relationships

- [ ] Deep loading / eager fetching
- [ ] Cascade delete options
- [ ] Record unions (polymorphic relations)

### Database Features

- [ ] Migration support
- [ ] Schema versioning
- [ ] Views
- [ ] Custom functions
- [ ] Computed fields

### Developer Experience

- [ ] Watch mode for code generation
- [ ] Generation caching for faster rebuilds
- [ ] Better error messages
- [ ] IDE plugin support

## Version Goals

### v0.x (Current)

- Stabilize core functionality
- Expand data type support
- Improve documentation
- Gather community feedback

### v1.0 (Future)

- Stable API
- Comprehensive documentation
- Production-ready status
- Migration tooling

## Contributing

See the [Contributing Guide](01_contributing.md) to help with development.
