# FAQ

Frequently asked questions about SOM.

## General

### What is SOM?

SOM (SurrealDB Object Mapper) is a code generation tool that creates type-safe database access code for SurrealDB from Go struct definitions.

### Is SOM production-ready?

SOM is currently in **early development** and considered experimental. It works for many use cases but may have unknown bugs. Use with caution in production.

### What Go version is required?

Go 1.23 or later is required due to heavy use of generics.

### What SurrealDB version is supported?

SOM is tested against SurrealDB 2.3.10. Compatibility with other versions is not guaranteed.

## Technical

### Why does SOM use code generation?

Code generation provides:
- Compile-time type safety
- Zero runtime reflection overhead
- Full IDE autocompletion support
- Better performance

### Why not use the official SurrealDB Go client?

SOM uses [sdbc](https://github.com/go-surreal/sdbc), a custom client optimized for code generation patterns. The official client may be supported in the future.

### Can I use raw SurrealQL queries?

Currently, SOM focuses on the type-safe query builder. Raw query support is planned.

### How do I handle database migrations?

Migration support is planned but not yet implemented. Currently, schema changes require manual handling.

## Troubleshooting

### The generator isn't finding my models

Ensure your structs embed `som.Node` or `som.Edge`:

```go
type User struct {
    som.Node  // Required!
    Name string
}
```

### I get import errors in generated code

Regenerate the code after any model changes:

```bash
go run github.com/go-surreal/som/cmd/som@latest gen ./model ./gen/som
```

## More Questions?

Open a [GitHub Discussion](https://github.com/go-surreal/som/discussions) for additional questions.
