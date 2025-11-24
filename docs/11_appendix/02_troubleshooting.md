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

Also check that:
- Model files are in the input directory
- Files have `.go` extension
- Package compiles without errors

### Import errors in generated code

**Problem**: Generated code has broken imports.

**Solution**:
1. Delete the generated output directory completely
2. Regenerate the code
3. Run `go mod tidy` in your project

```bash
rm -rf ./gen/som
go run github.com/go-surreal/som/cmd/som@latest gen ./model ./gen/som
go mod tidy
```

### Generator crashes or panics

**Problem**: The generator panics or crashes.

**Solution**:
1. Check for syntax errors in your model files
2. Ensure all imports are valid
3. Try with a minimal model to isolate the issue
4. Check for unsupported types (see [Supported Types](../04_data_types/01_primitives.md))
5. Report the issue on GitHub with reproduction steps

### Edge missing In/Out fields

**Problem**: Edge generates but doesn't have In/Out relationships.

**Solution**: Ensure your edge has both fields with proper tags:

```go
type Follows struct {
    som.Edge

    In  *User `som:"in"`   // Must have som:"in" tag
    Out *User `som:"out"`  // Must have som:"out" tag

    // Additional fields...
}
```

## Runtime Issues

### Connection refused

**Problem**: Cannot connect to SurrealDB.

**Solution**:

1. Verify SurrealDB is running:
   ```bash
   surreal start --user root --pass root memory
   ```

2. Check the address format:
   ```go
   client, err := som.NewClient(ctx, som.Config{
       Address:   "ws://localhost:8000",  // WebSocket protocol
       // or
       Address:   "http://localhost:8000", // HTTP protocol
   })
   ```

3. Verify credentials match your SurrealDB setup

4. Check namespace and database exist (SurrealDB creates them automatically, but verify names are correct)

### Record not found after Create

**Problem**: `Read` returns `exists=false` for a record you just created.

**Solution**:

1. Ensure Create succeeded:
   ```go
   err := client.UserRepo().Create(ctx, user)
   if err != nil {
       log.Fatal(err)  // Don't ignore errors
   }
   ```

2. Use the correct ID reference:
   ```go
   // Correct - use ID() method
   retrieved, exists, err := client.UserRepo().Read(ctx, user.ID())

   // Wrong - don't construct IDs manually
   // retrieved, exists, err := client.UserRepo().Read(ctx, "user:123")
   ```

3. Check you're using the same namespace/database

### Type mismatch errors

**Problem**: Errors when reading/writing certain fields.

**Solution**:

1. Regenerate code after any model changes
2. Check for unsupported types:
   - `uint64`, `uintptr` - Not supported
   - `complex64`, `complex128` - Not supported
   - Channels, functions - Not supported

3. Ensure pointer/non-pointer consistency between model and database

### Timestamps not updating

**Problem**: `CreatedAt` or `UpdatedAt` fields are not being set.

**Solution**:

1. Ensure you're embedding `som.Timestamps`:
   ```go
   type User struct {
       som.Node
       som.Timestamps  // Must be embedded
       Name string
   }
   ```

2. Timestamps are readonly - don't try to set them manually
3. Regenerate code if you just added `som.Timestamps`

## Query Issues

### Filter not matching expected records

**Problem**: Query returns wrong results.

**Solution**:

1. Check filter logic - multiple filters are ANDed:
   ```go
   // This is AND (both must be true)
   query.Filter(
       where.User.IsActive.IsTrue(),
       where.User.Age.GreaterThan(18),
   )

   // Use where.Any for OR
   query.Filter(
       where.Any(
           where.User.Role.Equal("admin"),
           where.User.Role.Equal("moderator"),
       ),
   )
   ```

2. Check comparison direction:
   ```go
   // "Age greater than 18" not "18 greater than age"
   where.User.Age.GreaterThan(18)
   ```

3. Check nil handling for optional fields:
   ```go
   // Optional fields might be nil
   where.User.DeletedAt.IsNil()  // Not deleted
   ```

### Live query not receiving updates

**Problem**: Live query channel doesn't receive any updates.

**Solution**:

1. Check error handling:
   ```go
   updates, err := client.UserRepo().Query().Live(ctx)
   if err != nil {
       log.Fatal(err)  // Connection might have failed
   }
   ```

2. Ensure context isn't cancelled:
   ```go
   ctx, cancel := context.WithCancel(context.Background())
   // Don't call cancel() until you want to stop
   ```

3. Check that changes actually occur in the database after subscribing

4. Verify filters match the records being changed

## Performance Issues

### Slow queries

**Solution**:

1. Add filters to reduce result sets:
   ```go
   // Bad - fetches all users
   users, _ := client.UserRepo().Query().All(ctx)

   // Good - fetches only what you need
   users, _ := client.UserRepo().Query().
       Filter(where.User.IsActive.IsTrue()).
       Limit(100).
       All(ctx)
   ```

2. Use `Count` or `Exists` when you don't need full records:
   ```go
   count, _ := client.UserRepo().Query().
       Filter(where.User.IsActive.IsTrue()).
       Count(ctx)
   ```

3. Consider indexing frequently queried fields in SurrealDB

### High memory usage

**Solution**:

1. Use pagination:
   ```go
   for page := 0; ; page++ {
       users, _ := client.UserRepo().Query().
           Limit(100).
           Offset(page * 100).
           All(ctx)
       if len(users) == 0 {
           break
       }
       process(users)
   }
   ```

2. Close unused live query subscriptions by cancelling their context

3. Don't hold references to large result sets longer than needed

### Connection pool exhaustion

**Solution**:

1. Limit concurrent live queries
2. Cancel contexts when subscriptions are no longer needed
3. Consider connection pooling at the infrastructure level

## Getting Help

1. Check existing [GitHub Issues](https://github.com/go-surreal/som/issues)
2. Search [GitHub Discussions](https://github.com/go-surreal/som/discussions)
3. Open a new issue with:
   - SOM version
   - Go version
   - SurrealDB version
   - Minimal reproduction code
   - Expected vs actual behavior
