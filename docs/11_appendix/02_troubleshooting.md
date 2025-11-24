# Troubleshooting

Common issues and solutions when working with SOM.

## Code Generation Issues

### "No models found"

**Problem**: The generator doesn't find any models.

**Solution**: Ensure your structs embed `som.Node` or `som.Edge`:

```go
type User struct {
    som.Node  // This is required
    Name string
}
```

### Import errors in generated code

**Problem**: Generated code has broken imports.

**Solution**:
1. Delete the generated output directory
2. Regenerate the code
3. Run `go mod tidy` in your project

### Generator crashes

**Problem**: The generator panics or crashes.

**Solution**:
1. Check for syntax errors in your model files
2. Ensure all imports are valid
3. Try with a minimal model to isolate the issue
4. Report the issue on GitHub with reproduction steps

## Runtime Issues

### Connection refused

**Problem**: Cannot connect to SurrealDB.

**Solution**:
1. Verify SurrealDB is running
2. Check the address (ws://, wss://, http://, https://)
3. Verify credentials (username/password)
4. Check namespace and database names

```go
client, err := som.NewClient(ctx, som.Config{
    Address:   "ws://localhost:8000",  // Correct protocol and port
    Username:  "root",
    Password:  "root",
    Namespace: "test",
    Database:  "test",
})
```

### Record not found

**Problem**: `Read` returns nil or error for existing records.

**Solution**:
1. Verify the ID format matches SurrealDB's format
2. Check you're connected to the correct namespace/database
3. Ensure the record was actually created

### Type mismatch errors

**Problem**: Errors when reading/writing certain fields.

**Solution**:
1. Ensure field types match supported types
2. Check for unsupported types (`uint64`, `complex64`, etc.)
3. Regenerate code after model changes

## Performance Issues

### Slow queries

**Solution**:
1. Add appropriate filters to reduce result sets
2. Use `Limit` for pagination
3. Consider indexing frequently queried fields in SurrealDB

### High memory usage

**Solution**:
1. Use pagination instead of loading all records
2. Process results in batches
3. Close unused live query subscriptions

## Getting Help

1. Check existing [GitHub Issues](https://github.com/go-surreal/som/issues)
2. Search [GitHub Discussions](https://github.com/go-surreal/som/discussions)
3. Open a new issue with reproduction steps
